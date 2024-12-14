package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/service"
	"fmt"
)

func LoadBindName() {
	global.GBindMap = service.ServiceGroup.DbService.QueryAllBind()
	global.GLog.Info(fmt.Sprintf("加载到了 %d 条绑定关系", len(global.GBindMap)))
}
