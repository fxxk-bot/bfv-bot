package initialize

import (
	"bfv-bot/common/global"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"strings"
)

func Ai() {

	if global.GConfig.Ai.Enable {
		qianfan.GetConfig().AccessKey = global.GConfig.Ai.AccessKey
		qianfan.GetConfig().SecretKey = global.GConfig.Ai.SecretKey

		// 支持的模型
		list := global.GAi.ModelList()
		result := strings.Join(list, ", ")
		global.GLog.Debug("支持的模型")
		global.GLog.Debug(result)

		// 目前speed模型和lite模型是免费使用 Yi-34B-Chat(第三方模型)限时免费
		global.GAi = qianfan.NewChatCompletion(qianfan.WithModel(global.GConfig.Ai.ModelName))
	}
}
