package flow

import (
	"bfv-bot/bot/private"
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"unicode/utf8"
)

func AddBlackStep(msg *req.MsgData, flowData *flow.PrivateFlow) {
	if flowData.Step == 2 {
		if utf8.RuneCountInString(msg.RawMessage) > 20 {
			private.SendPrivateMsg(msg.UserID, "[添加黑名单] 原因不能超过20个字")
			return
		}

		name := flowData.Content[0]
		delete(PrivateFlowable, msg.UserID)
		id, err := dbService.AddBlack(name, msg.RawMessage)
		if err != nil {
			private.SendPrivateMsg(msg.UserID, "[添加黑名单] 添加失败")
			return
		} else {
			private.SendPrivateMsg(msg.UserID, "[添加黑名单] 黑名单添加成功, id: "+id)
			return
		}
	}
}
