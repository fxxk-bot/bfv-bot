package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/service"
	"fmt"
	"github.com/feiin/sensitivewords"
)

func InitSensitive() {
	list := service.ServiceGroup.SelectAllSensitive()
	global.GSensitive = sensitivewords.New()
	if len(list) == 0 {
		global.GLog.Info("未加载到敏感词")
		return
	}
	global.GSensitive.AddWords(list...)
	global.GLog.Info(fmt.Sprintf("加载到 %d 条敏感词", len(list)))
}
