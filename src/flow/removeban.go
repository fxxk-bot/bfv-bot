package flow

import (
	"bfv-bot/bot/group"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"unicode/utf8"
)

func RemoveBanStep(msg *req.MsgData, flowData *flow.GroupFlow) {
	if flowData.Step == 2 {
		serverInfo := global.GConfig.Bfv.GetGameInfo(msg.RawMessage)
		if serverInfo.GetGameId() == "" {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 未能找到游戏ID, 流程终止")
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
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 由于网络或其他原因导致失败, 流程终止")
			return
		}
		group.SendGroupImageMsg(msg.GroupID, captcha)

		IncGroupStep(msg, serverInfo.GetToken())
		IncGroupStep(msg, serverInfo.GetGameId())
		IncGroupStep(msg, serverInfo.ServerName)
		IncGroupStep(msg, hash)

	} else if flowData.Step == 6 {

		// 解除屏蔽
		name := flowData.Content[0]

		pid := flowData.Content[1]
		token := flowData.Content[2]
		gameId := flowData.Content[3]
		serverName := flowData.Content[4]
		captchaHash := flowData.Content[5]
		captchaValue := msg.RawMessage

		if utf8.RuneCountInString(captchaValue) > 20 {
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[屏蔽] 输入的验证码过长")
			return
		}

		err, removeBanResult := utils.RemoveBan(captchaValue, captchaHash, gameId, pid, serverName, name, token)
		if err != nil {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 由于网络或其他原因导致失败")
			return
		}

		if removeBanResult.Error != 0 {
			if removeBanResult.Code == "captcha.wrong" || removeBanResult.Code == "captcha.bad" {
				// 保留group step 以重新输入验证码
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 验证码错误")
				return
			} else if removeBanResult.Code == "user.tokenExpired" {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] Token过期")
				return
			} else if removeBanResult.Code == "verifyServer.server not found" {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 无法操作该服务器")
				return
			} else {
				DeleteGroupStep(msg)
				group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 其它失败: "+removeBanResult.Code)
				return
			}
		}

		if removeBanResult.Success == 1 {
			DeleteGroupStep(msg)
			group.SendGroupReplyMsg(msg.GroupID, flowData.MsgId, "[解除屏蔽] 成功")
			return
		}

	}
}
