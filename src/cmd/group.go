package cmd

import (
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/common/req"
	"bfv-bot/model/common/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func init() {
	groupCommandMap["ban"] = ban
	groupCommandMap["removeban"] = removeban
}

func cx(_ *req.MsgData, c *gin.Context, _ string, value string) {

	path, err := utils.QueryAndStore(value, 1)
	if err != nil {
		global.GLog.Error("utils.QueryAndStore, 1", zap.String("name", value), zap.Error(err))
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}
	global.GLog.Info("file:///" + path)
	resp.ImageOk(c, "file:///"+path, value+"的查询结果")
}

func ban(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	if !global.GConfig.QQBot.IsActiveAdminGroup(msg.GroupID) {
		resp.EmptyOk(c)
		return
	}
	resp.ReplyOk(c, "[屏蔽] 已下线")
}

func platoon(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, result := utils.GetJoinPlatoonsByName(value)
	if err != nil {
		global.GLog.Error("utils.GetJoinPlatoonsByName",
			zap.String("name", value), zap.Error(err))
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}
	resp.ReplyOk(c, result)
}

func banlog(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err, result := utils.GetBanLog(value)
	if err != nil {
		global.GLog.Error("utils.GetBanLog",
			zap.String("name", value), zap.Error(err))
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}
	resp.ReplyOk(c, result)
}

func removeban(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	if !global.GConfig.QQBot.IsActiveAdminGroup(msg.GroupID) {
		resp.EmptyOk(c)
		return
	}
	resp.ReplyOk(c, "[解除屏蔽] 已下线")
}

func bind(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, data := utils.CheckPlayer(value)
	if err != nil {
		resp.ReplyOk(c, "绑定失败 "+err.Error())
		return
	}
	err = dbService.AddBind(msg.UserID, value, data.PID)
	if err != nil {
		resp.ReplyOk(c, "绑定失败 "+err.Error())
		return
	} else {
		resp.ReplyOk(c, "绑定成功: "+data.PID)
		return
	}
}

func server(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err, str := utils.GetBfvRobotServer(value, true)
	if err != nil {
		return
	}
	resp.ReplyOk(c, str)
}

func data(_ *req.MsgData, c *gin.Context, _ string, value string) {

	path, err := utils.QueryAndStore(value, 2)
	if err != nil {
		global.GLog.Error("utils.QueryAndStore, 2", zap.String("name", value), zap.Error(err))
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}
	global.GLog.Info("file:///" + path)
	resp.ImageOk(c, "file:///"+path, value+"的完整数据")
}

func task(_ *req.MsgData, c *gin.Context, _ string, value string) {

	offset, err := strconv.Atoi(value)
	if err != nil {
		resp.ReplyOk(c, "必须是数字")
		return
	}
	path, err := utils.GetTaskAndCache(offset)
	if err != nil {
		global.GLog.Error("utils.GetTaskAndCache", zap.String("value", value), zap.Error(err))
		resp.ReplyOk(c, "获取失败: "+err.Error())
		return
	}
	global.GLog.Info("file:///" + path)
	resp.ImageOk(c, "file:///"+path, "周任务数据")
}

func playerlist(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err, path := utils.GetPlayerList(value)
	if err != nil {
		resp.ReplyOk(c, "获取失败: "+err.Error())
		return
	}
	global.GLog.Info("file:///" + path)
	resp.ImageOk(c, "file:///"+path, "服务器玩家列表")
}

func groupMember(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err, s := utils.GerServerGroupMember(value)
	if err != nil {
		resp.ReplyOk(c, err.Error())
		return
	}
	resp.ReplyOk(c, s)
}

func quickTask(msg *req.MsgData, c *gin.Context, key string) {
	task(msg, c, key, "0")
}

func ShortCommandFunction(msg *req.MsgData, c *gin.Context, command string) {
	err, name := dbService.GetBindName(msg.UserID)
	if err != nil {
		resp.ReplyOk(c, "快捷查询失败: "+err.Error())
		return
	}
	groupCommandFunction, groupCommandOk := groupCommandMap[command]
	if groupCommandOk {
		groupCommandFunction(msg, c, command, name)
	}
}

func getGroupServerInfo(_ *req.MsgData, c *gin.Context, _ string) {
	err, result := utils.GetBfvRobotServer(global.GConfig.Bfv.GroupUniName, false)
	if err == nil {
		resp.ReplyOk(c, result)
		return
	} else {
		resp.ReplyOk(c, err.Error())
		return
	}
}

func quickCx(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err, data := utils.CheckPlayer(value)
	if err != nil {
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}

	err, finalMsg := utils.GetBaseInfoAndStatusByName(&data)
	if err != nil {
		resp.ReplyOk(c, "查询失败: "+err.Error())
		return
	}
	resp.ReplyOk(c, "玩家 ["+data.Name+"] 基础数据如下\n\n"+finalMsg)
}

func getGroupHelpInfo(_ *req.MsgData, c *gin.Context, _ string) {
	var builder strings.Builder

	if len(global.GConfig.QQBot.CustomCommandKey.Banlog) != 0 {
		builder.WriteString("玩家屏蔽记录: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Banlog, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Server) != 0 {
		builder.WriteString("服务器查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Server, "/") + "=<name>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Data) != 0 {
		builder.WriteString("完整数据查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Data, "/") + "=<name>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Task) != 0 {
		builder.WriteString("周任务查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Task, "/") + "=<上周: -1, 本周: 0, 下周: 1>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Playerlist) != 0 {
		builder.WriteString("服务器玩家列表查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Playerlist, "/") + "=<服务器>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.GroupMember) != 0 {
		builder.WriteString("查询服务器内群友: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.GroupMember, "/") + "=<服务器>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Cx) != 0 {
		builder.WriteString("战绩查询: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.Cx, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.C) != 0 {
		builder.WriteString("快捷查询: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.C, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Bind) != 0 {
		builder.WriteString("绑定玩家: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Bind, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Platoon) != 0 {
		builder.WriteString("加入的战排: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Platoon, "/") + "=<id>\n")
	}

	builder.WriteString("其他快捷指令: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.GroupServer, "/") +
		"/" + global.GConfig.Bfv.GroupName)

	resp.ReplyOk(c, builder.String())
}

func InitBanlogKey(key string) {
	groupCommandMap[key] = banlog
	groupShortCommandMap[key] = true
}

func InitCxKey(key string) {
	groupCommandMap[key] = cx
	groupShortCommandMap[key] = true
}

func InitCKey(key string) {
	groupShortCommandMap[key] = true
	groupCommandMap[key] = quickCx
}

func InitPlatoonKey(key string) {
	groupCommandMap[key] = platoon
	groupShortCommandMap[key] = true
}

func InitBindKey(key string) {
	groupCommandMap[key] = bind
}

func InitServerKey(key string) {
	groupCommandMap[key] = server
}

func InitDataKey(key string) {
	groupCommandMap[key] = data
	groupShortCommandMap[key] = true
}

func InitTaskKey(key string) {
	groupCommandMap[key] = task
}

func InitPlayerListKey(key string) {
	groupCommandMap[key] = playerlist
}

func InitGroupMemberKey(key string) {
	groupCommandMap[key] = groupMember
}

func InitHelpKey(key string) {
	groupQuickCommandMap[key] = getGroupHelpInfo
}

func InitGroupServerKey(key string) {
	groupQuickCommandMap[key] = getGroupServerInfo
}

func InitQuickTaskKey(key string) {
	groupQuickCommandMap[key] = quickTask
}

func GetGroupCommandFunc(key string) (func(*req.MsgData, *gin.Context, string, string), bool) {
	f, ok := groupCommandMap[key]
	return f, ok
}

func GetGroupShortCommandFunc(key string) (bool, bool) {
	f, ok := groupShortCommandMap[key]
	return f, ok
}

func GetGroupQuickCommandFunc(key string) (func(*req.MsgData, *gin.Context, string), bool) {
	f, ok := groupQuickCommandMap[key]
	return f, ok
}
