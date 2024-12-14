package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/service"
)

// LoadIgnoreList 一分钟加载一次忽略名单
func LoadIgnoreList() {
	global.GIgnoreListMap = service.ServiceGroup.DbService.QueryAllIgnoreList()
}
