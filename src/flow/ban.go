package flow

import (
	"bfv-bot/bot/group"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"go.uber.org/zap"
	"strconv"
	"unicode/utf8"
)

func BanStep(msg *req.MsgData, flowData *flow.GroupFlow) {

	if flowData.Step == 2 {
		if utf8.RuneCountInString(msg.RawMessage) > 20 {
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 屏蔽原因不能超过20个字")
			return
		}
		IncGroupStep(msg, msg.RawMessage)
		group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 选择服务器")
	} else if flowData.Step == 3 {
		serverInfo := global.GConfig.Bfv.GetGameInfo(msg.RawMessage)
		if serverInfo.GetGameId() == "" {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 未能找到GameID, 流程终止")
			return
		}
		if serverInfo.GetToken() == "" {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 该服务器未配置Token, 流程终止")
			return
		}
		captcha, hash, err := utils.GetCaptchaBase64()
		if err != nil {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 验证码获取失败, 流程终止")
			return
		}
		group.SendGroupImageReplyMsg(msg.GroupID, flowData.MsgId, captcha)

		IncGroupStep(msg, serverInfo.GetToken())
		IncGroupStep(msg, serverInfo.GetGameId())
		IncGroupStep(msg, hash)
	} else if flowData.Step == 6 {

		// 执行屏蔽
		name := flowData.Content[0]
		pidStr := flowData.Content[1]
		reason := flowData.Content[2]
		token := flowData.Content[3]
		gameId := flowData.Content[4]
		captchaHash := flowData.Content[5]
		captchaValue := msg.RawMessage
		if utf8.RuneCountInString(captchaValue) > 20 {
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 输入的验证码过长")
			return
		}

		pid, err := strconv.ParseInt(pidStr, 10, 64)
		if err != nil {
			global.GLog.Error("strconv.ParseInt", zap.Error(err))
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] PID数据异常, 流程终止")
			return
		}
		err, banResult := utils.Ban(captchaValue, captchaHash, gameId, reason, pid, name, token)
		if err != nil {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 由于网络或其他原因导致失败")
			return
		}

		if banResult.Error != 0 {
			if banResult.Code == "captcha.wrong" || banResult.Code == "captcha.bad" {
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 验证码错误")
				return
			} else if banResult.Code == "user.tokenExpired" {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] Token过期")
				return
			} else if banResult.Code == "verifyServer.server not found" {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 无法操作该服务器")
				return
			} else {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 其它失败: "+banResult.Code)
				return
			}
		}

		if banResult.Success == 1 {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 成功")
			return
		}

	}
}
