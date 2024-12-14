package cmd

import (
	"bfv-bot/model/common/req"
	"bfv-bot/service"
	"github.com/gin-gonic/gin"
)

var (
	// dbService private
	dbService   = service.ServiceGroup.DbService
	cronService = service.ServiceGroup.CronService

	// 私聊命令映射
	privateCommandMap = make(map[string]func(*req.MsgData, *gin.Context, string, string))

	// 私聊Op命令映射
	privateOpCommandMap = make(map[string]func(*req.MsgData, *gin.Context, string, string))

	// 私聊快捷命令
	privateQuickCommandMap = make(map[string]func(*req.MsgData, *gin.Context, string))

	// 群聊命令映射
	groupCommandMap = make(map[string]func(*req.MsgData, *gin.Context, string, string))

	// 群聊短命令映射
	groupShortCommandMap = make(map[string]bool)

	// 群聊快捷命令 没有key value形式的
	groupQuickCommandMap = make(map[string]func(*req.MsgData, *gin.Context, string))
)
