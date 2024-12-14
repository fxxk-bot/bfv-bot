package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/service"
	"fmt"
)

// LoadBlackList 加载黑名单
func LoadBlackList() {
	global.GBlackListMap = service.ServiceGroup.QueryAllBlackList()
	global.GLog.Info(fmt.Sprintf("加载到了 %d 条黑名单", len(global.GBlackListMap)))
}
