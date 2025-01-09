package initialize

import (
	botPrivate "bfv-bot/bot/private"
	"bfv-bot/cmd"
	"bfv-bot/common/config"
	"bfv-bot/common/global"
)

func InitBot() {

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Banlog {
		cmd.InitBanlogKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Cx {
		cmd.InitCxKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.C {
		cmd.InitCKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Platoon {
		cmd.InitPlatoonKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Bind {
		cmd.InitBindKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Server {
		cmd.InitServerKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Data {
		cmd.InitDataKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Task {
		cmd.InitTaskKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Playerlist {
		cmd.InitPlayerListKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.GroupMember {
		cmd.InitGroupMemberKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Task {
		cmd.InitQuickTaskKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.Help {
		cmd.InitHelpKey(item)
	}

	for _, item := range global.GConfig.QQBot.CustomCommandKey.GroupServer {
		cmd.InitGroupServerKey(item)
	}
	cmd.InitGroupServerKey(global.GConfig.Bfv.GroupName)

	global.GConfig.QQBot.InitMap()
	botPrivate.SendPrivateMsg(global.GConfig.QQBot.SuperAdminQq, "服务启动成功\n"+config.GetVersion())

	if global.GConfig.Bfv.Active {
		botPrivate.SendPrivateMsg(global.GConfig.QQBot.SuperAdminQq, "已自动开启检测功能")
	}
}
