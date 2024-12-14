package flow

import (
	"bfv-bot/bot/private"
	"bfv-bot/common/global"
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
)

func BindTokenStep(msg *req.MsgData, flowData *flow.PrivateFlow) {
	if flowData.Step == 2 {
		serverInfo := global.GConfig.Bfv.GetGameInfo(msg.RawMessage)
		if serverInfo.ServerName == "" {
			private.SendPrivateMsg(msg.UserID, "[绑定Token] 服务器ID错误")
			return
		}
		token := flowData.Content[0]
		delete(PrivateFlowable, msg.UserID)
		global.GConfig.Bfv.SetToken(msg.RawMessage, token)
		private.SendPrivateMsg(msg.UserID, "[绑定Token] 服务器: "+serverInfo.ServerName+" 绑定成功")
	}
}
