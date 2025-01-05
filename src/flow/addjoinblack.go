package flow

import (
	"bfv-bot/bot/private"
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"unicode/utf8"
)

func AddJoinBlackStep(msg *req.MsgData, flowData *flow.PrivateFlow) {
	if flowData.Step == 2 {
		if utf8.RuneCountInString(msg.RawMessage) > 20 {
			private.SendPrivateMsg(msg.UserID, "[添加加群黑名单] 原因不能超过20个字")
			return
		}

		qq := flowData.Content[0]
		delete(PrivateFlowable, msg.UserID)
		err := dbService.AddJoinBlackList(qq, msg.RawMessage)
		if err != nil {
			private.SendPrivateMsg(msg.UserID, "[添加加群黑名单] 添加失败")
			return
		} else {
			private.SendPrivateMsg(msg.UserID, "[添加加群黑名单] 添加成功")
			return
		}
	}
}
