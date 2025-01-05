package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/service"
	"fmt"
)

func LoadJoinBlackList() {
	global.GJoinBlackListMap = service.ServiceGroup.DbService.QueryAllJoinBlackList()
	global.GLog.Info(fmt.Sprintf("加载到了 %d 条加群黑名单", len(global.GJoinBlackListMap)))
}
