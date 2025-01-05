package main

import (
	"bfv-bot/common/global"
	"bfv-bot/common/initialize"
	"fmt"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy

func main() {

	// 读取配置
	initialize.Viper()
	// 初始化日志
	global.GLog = initialize.Zap()
	zap.ReplaceGlobals(global.GLog)

	// 初始化数据库链接以及DAO
	initialize.InitDb()

	// 协程池
	initialize.Ants()

	// 初始化路由
	routers := initialize.Routers()

	// 初始化ai
	initialize.Ai()

	// 定时任务
	initialize.Cron()

	// 敏感词系统
	initialize.InitSensitive()

	// 加载快捷查询绑定关系
	initialize.LoadBindName()

	// 黑名单列表
	initialize.LoadBlackList()

	// 加群黑名单列表
	initialize.LoadJoinBlackList()

	// rod
	initialize.InitRod()

	// tof
	initialize.InitTofData()

	// bot
	initialize.InitBot()

	port := fmt.Sprintf(":%d", global.GConfig.Server.Port)

	apiServer := initialize.InitServer(port, routers)

	global.GLog.Info(fmt.Sprintf("服务成功启动在 Port: %d", global.GConfig.Server.Port))

	global.GLog.Error(apiServer.ListenAndServe().Error())
}
