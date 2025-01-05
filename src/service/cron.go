package service

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"bfv-bot/common/config"
	"bfv-bot/common/cons"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/dto"
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

	for index, serverInfo := range global.GConfig.Bfv.Server {

		time.Sleep(2 * time.Second)

		if serverInfo.GetGameId() == "" {
			global.GLog.Info(serverInfo.ServerName + " gameid为空, 跳过检测")
			continue
		}
		err, apiPlayers := utils.GetServerPlayerByGameToolsConvert(serverInfo.GetGameId())
		if err != nil {
			return
		}
		// 黑名单通知
		aTeamNotifyMap := make(map[string]string)
		bTeamNotifyMap := make(map[string]string)

		playerMap := make(map[int64]string)

		for _, player := range apiPlayers.TeamOne {

			value, ok := global.GBlackListMap[strconv.FormatInt(player.PersonaID, 10)]
			if ok {
				aTeamNotifyMap[player.Name] = value.Reason
			}
			playerMap[player.PersonaID] = player.Name
		}

		for _, player := range apiPlayers.TeamTwo {

			value, ok := global.GBlackListMap[strconv.FormatInt(player.PersonaID, 10)]
			if ok {
				bTeamNotifyMap[player.Name] = value.Reason
			}
			playerMap[player.PersonaID] = player.Name
		}

		checkRankAndKpm(&playerMap, serverInfo, index)

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

// checkRankAndKpm 请求接口并检查等级以及kpm
func checkRankAndKpm(all *map[int64]string, serverInfo config.ServerInfo, index int) {

	if serverInfo.Kpm == 0 && serverInfo.MaxRank == 0 && serverInfo.MinRank == 0 {
		return
	}

	insertions := make([]dto.CheckPlayerData, 0)
	exists := make([]dto.GtBatchStatusData, 0)

	playerMap := serverInfo.GetPlayerMap()
	// 如果上次数据长度为空 视为第一次检测
	if len(playerMap) == 0 {
		for key, value := range *all {
			insertions = append(insertions, dto.CheckPlayerData{PersonaID: key, Name: value})
		}
	} else {
		for key, value := range *all {
			playerMapValue, ok := playerMap[key]
			if ok {
				exists = append(exists, playerMapValue)
			} else {
				insertions = append(insertions, dto.CheckPlayerData{PersonaID: key, Name: value})
			}
		}
	}

	finalResult := make([]dto.GtBatchStatusData, 0)
	if len(insertions) > 16 {
		// 计算中间位置
		mid := len(insertions) / 2

		part := make([][]dto.CheckPlayerData, 2)
		// 分割切片
		part[0] = insertions[:mid]
		part[1] = insertions[mid:]
		wg := sync.WaitGroup{}
		wg.Add(2)

		resultPart := make([][]dto.GtBatchStatusData, 2)

		for i := 0; i < 2; i++ {
			err := global.GPool.Submit(func() {
				defer wg.Done()
				err, data := utils.GetGameToolsBatchStatus(getId(part[i]))
				if err != nil {
					return
				}
				resultPart[i] = data
			})
			if err != nil {
				return
			}
		}
		wg.Wait()
		finalResult = append(finalResult, resultPart[0]...)
		finalResult = append(finalResult, resultPart[1]...)

	} else {
		err, data := utils.GetGameToolsBatchStatus(getId(insertions))
		if err != nil {
			return
		}
		finalResult = data
	}

	finalResult = append(finalResult, exists...)

	finalMap := make(map[int64]dto.GtBatchStatusData)
	for _, data := range finalResult {
		pid, err := strconv.ParseInt(data.ID, 10, 64)
		if err != nil {
			continue
		}
		value, ok := (*all)[pid]
		if ok {
			data.Name = value
		}
		finalMap[pid] = data
	}

	global.GConfig.Bfv.Server[index].SetPlayerMap(finalMap)

	builder := strings.Builder{}

	for _, data := range finalMap {
		loopBuf := strings.Builder{}
		flag := false
		if serverInfo.Kpm != 0 && data.KillsPerMinute > serverInfo.Kpm {
			flag = true
			loopBuf.WriteString("KPM: ")
			loopBuf.WriteString(strconv.FormatFloat(data.KillsPerMinute, 'f', 2, 64))
			loopBuf.WriteString(" 高于设定值,")
		}

		if serverInfo.MaxRank != 0 && data.Rank > serverInfo.MaxRank {
			flag = true
			loopBuf.WriteString("等级: ")
			loopBuf.WriteString(strconv.FormatFloat(data.Rank, 'f', 0, 64))
			loopBuf.WriteString(" 高于设定值,")

		}

		if serverInfo.MinRank != 0 && data.Rank < serverInfo.MinRank {
			flag = true
			loopBuf.WriteString("等级: ")
			loopBuf.WriteString(strconv.FormatFloat(data.Rank, 'f', 0, 64))
			loopBuf.WriteString(" 低于设定值,")
		}

		if flag {
			str := loopBuf.String()
			str = str[:len(str)-1]
			builder.WriteString(data.Name)
			builder.WriteString("\t[")
			builder.WriteString(str)
			builder.WriteString("]")
			builder.WriteString("\n")
		}
	}
	finalStr := builder.String()
	if len(finalStr) > 0 {
		finalStr = finalStr[:len(finalStr)-1]
		group.SendGroupMsgMultiple(global.GConfig.QQBot.AdminGroup, serverInfo.ServerName+"\n\n"+finalStr)
	}
}

func getId(param []dto.CheckPlayerData) []int64 {
	int64s := make([]int64, 0)
	for _, data := range param {
		int64s = append(int64s, data.PersonaID)
	}
	return int64s
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
