package service

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"bfv-bot/common/cons"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/po"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CronService struct{}

// CheckBlackListAndNotify 判断黑名单
func (c *CronService) CheckBlackListAndNotify() {

	if !global.GConfig.Bfv.Active {
		return
	}

	for _, serverInfo := range global.GConfig.Bfv.Server {

		time.Sleep(2 * time.Second)

		if serverInfo.GetGameId() == "" {
			continue
		}
		err, apiPlayers := utils.GetServerPlayerByGameToolsConvert(serverInfo.GetGameId())
		if err != nil {
			return
		}
		// 黑名单通知
		aTeamNotifyMap := make(map[string]string)
		bTeamNotifyMap := make(map[string]string)

		for _, player := range apiPlayers.TeamOne {

			if _, ok := global.GIgnoreListMap[strings.ToLower(player.Name)]; ok {
				continue
			}

			value, ok := global.GBlackListMap[strconv.FormatInt(player.PersonaID, 10)]
			if ok {
				aTeamNotifyMap[player.Name] = value.Reason
			}
		}

		for _, player := range apiPlayers.TeamTwo {

			if _, ok := global.GIgnoreListMap[strings.ToLower(player.Name)]; ok {
				continue
			}

			value, ok := global.GBlackListMap[strconv.FormatInt(player.PersonaID, 10)]
			if ok {
				bTeamNotifyMap[player.Name] = value.Reason
			}
		}

		aTeamCnt := len(apiPlayers.TeamOne)
		bTeamCnt := len(apiPlayers.TeamTwo)

		queCnt := len(apiPlayers.Que)

		if queCnt > 0 && ((aTeamCnt <= global.GConfig.Bfv.BlockingPlayers && bTeamCnt >= 32) ||
			(aTeamCnt >= 32 && bTeamCnt <= global.GConfig.Bfv.BlockingPlayers)) {
			var builder strings.Builder
			builder.WriteString(fmt.Sprintf("%s\n卡排队了\nA队人数:%d\nB队人数:%d\n及时处理!!要凉服了!!\n",
				serverInfo.ServerName, aTeamCnt, bTeamCnt))
			builder.WriteString("排队ID: \n")
			for i, item := range apiPlayers.Que {
				builder.WriteString(fmt.Sprintf("%d:   ", i+1))
				builder.WriteString(item.Name)
				builder.WriteString("\n")
			}

			group.SendGroupMsgMultiple(global.GConfig.QQBot.AdminGroup, builder.String())
		}

		if len(aTeamNotifyMap) != 0 || len(bTeamNotifyMap) != 0 {
			var builder strings.Builder
			builder.WriteString(serverInfo.ServerName)
			builder.WriteString("\n")
			builder.WriteString("黑名单玩家进入服务器\n")
			builder.WriteString("A队\n")
			for key, value := range aTeamNotifyMap {
				builder.WriteString(fmt.Sprintf("[%s]       %s\n", key, value))
			}

			builder.WriteString("B队\n")
			for key, value := range bTeamNotifyMap {
				builder.WriteString(fmt.Sprintf("[%s]       %s\n", key, value))
			}

			// 获取最终的字符串
			result := builder.String()

			group.SendGroupMsgMultiple(global.GConfig.QQBot.AdminGroup, result)
		}
	}

}

// GetTof 任务数据
func (c *CronService) GetTof() {
	oldData := global.GTofData

	err, data := utils.GetTof()

	if err != nil {
		global.GLog.Error("tof数据读取失败")
	}

	global.GTofData = data

	if oldData.StartTimestamp != global.GTofData.StartTimestamp {
		global.GTofDataCache = sync.Map{}
	}

}

// CheckCard 检测昵称
func (c *CronService) CheckCard() {
	// 查询 next_check_time 小于当前时间戳的
	err, data := ServiceGroup.DbService.QueryCardCheckByTime(time.Now().UnixMilli())
	if err != nil {
		return
	}

	for _, item := range data {
		global.GLog.Info(fmt.Sprintf("开始确认 %d 在 群: %d 的名片", item.Qq, item.GroupId))
		err, memberInfo := group.GetGroupMemberInfo(item.GroupId, item.Qq)
		if err != nil {
			continue
		}
		// 名片为空
		if memberInfo.Card == "" {
			doByStrategy(&item, memberInfo.Card)
			continue
		}
		err, apiResult := utils.CheckPlayer(memberInfo.Card)
		if err != nil {
			// 接口返回{} 判定为id错误 累计错误次数
			// 其他接口异常 不累计错误次数
			if err.Error() == cons.PlayerNotFound {
				doByStrategy(&item, memberInfo.Card)
			}
		} else if apiResult.PID != "" {
			_ = ServiceGroup.DbService.DeleteCardCheck(item.Qq)
			global.GLog.Info(fmt.Sprintf("已确认成功 %d 在 群: %d 的名片", item.Qq, item.GroupId))
		}
	}
}

func doByStrategy(item *po.CardCheck, card string) {
	if item.FailCnt == 1 {
		// 更新下次检测时间
		_ = ServiceGroup.DbService.UpdateCardCheck(item.Qq, 2,
			time.Unix(0, item.NextCheckTime*int64(time.Millisecond)).Add(42*time.Hour).UnixMilli())
		global.GLog.Info(fmt.Sprintf("%d 在 群: %d 的名片 需要再次检测", item.Qq, item.GroupId))
		if card == "" {
			group.SendAtGroupMsg(item.GroupId, item.Qq, "机器人无法确认你提供的ID, 请再次检查并修改你的群名片")
		} else {
			group.SendAtGroupMsg(item.GroupId, item.Qq, "机器人无法确认你提供的ID: ["+card+"]，"+
				"请再次检查并修改你的群名片")
		}
	} else if item.FailCnt >= 2 {
		err := ServiceGroup.DbService.DeleteCardCheck(item.Qq)
		if err != nil {
			global.GLog.Error("DbService.DeleteCardCheck", zap.Error(err))
			return
		}
		// 踢出
		if global.GConfig.QQBot.EnableAutoKickErrorNickname {
			group.SetGroupKick(item.GroupId, item.Qq)
			global.GLog.Info(fmt.Sprintf("%d 在 群: %d 被踢出", item.Qq, item.GroupId))
		}
	}

}

func (c *CronService) StartMute() {
	for _, item := range global.GConfig.QQBot.MuteGroup.ActiveGroup {
		group.SetGroupWholeBan(item, true)
		group.SendGroupMsg(item, global.GConfig.QQBot.MuteGroup.Start.Msg)
	}
	global.GLog.Info("已开启禁言")
}

func (c *CronService) EndMute() {
	for _, item := range global.GConfig.QQBot.MuteGroup.ActiveGroup {
		group.SendGroupMsg(item, global.GConfig.QQBot.MuteGroup.End.Msg)
		group.SetGroupWholeBan(item, false)
	}
	global.GLog.Info("已结束禁言")
}

func (c *CronService) BotToBot() {
	if !global.GConfig.QQBot.BotToBot.Enable {
		return
	}
	private.SendPrivateMsg(global.GConfig.QQBot.BotToBot.BotQq, global.GConfig.QQBot.BotToBot.Msg)
}

func (c *CronService) AutoBindGameId() {

	if !global.GConfig.QQBot.EnableAutoBindGameId && !global.GConfig.Bfv.Active {
		return
	}

	err, servers := utils.GetGameToolsServerByName(global.GConfig.Bfv.GroupUniName)
	if err != nil {
		return
	}

	// 先将配置中的 服务器名称和开服id对应起来
	localServer := make(map[string]string)

	for _, info := range global.GConfig.Bfv.Server {
		localServer[info.ServerName] = info.OwnerId
	}

	// 将api中的服务器名称和gameid对应起来
	apiMap := make(map[string]string)

	for _, item := range servers {
		value, ok := localServer[item.Prefix]
		if ok {
			// 接口返回的开服id必须和配置的pid相同
			if item.OwnerID == value {
				apiMap[item.Prefix] = item.GameID
			}
		}
	}

	// 绑定到运行时配置中
	for _, serverInfo := range global.GConfig.Bfv.Server {
		value, ok := apiMap[serverInfo.ServerName]
		if ok {
			global.GConfig.Bfv.SetGameId(serverInfo.Id, value)
			if serverInfo.GetGameId() != value {
				global.GLog.Info(fmt.Sprintf("[%s] 绑定到了GameId: [%s]", serverInfo.ServerName, value))
			}
		} else if serverInfo.GetGameId() != "" {
			global.GConfig.Bfv.SetGameId(serverInfo.Id, "")
			global.GLog.Info(fmt.Sprintf("[%s] 的GameId: [%s] 被重置",
				serverInfo.ServerName, serverInfo.GetGameId()))
		}
	}

}
