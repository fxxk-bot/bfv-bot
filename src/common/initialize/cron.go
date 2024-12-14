package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/flow"
	"bfv-bot/service"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type cronLogger struct {
}

func (l cronLogger) Info(msg string, keysAndValues ...interface{}) {
	global.GLog.Info(msg, zap.Any("keys", keysAndValues))
}

func (l cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	global.GLog.Error(msg, zap.Error(err), zap.Any("keys", keysAndValues))
}

func Cron() {
	// 创建一个新的cron实例
	global.GCron = cron.New(cron.WithChain(cron.Recover(cronLogger{})))

	// 90s 判断一次黑名单
	_, err := global.GCron.AddFunc("@every 90s", service.ServiceGroup.CronService.CheckBlackListAndNotify)
	if err != nil {
		global.GLog.Error("CheckBlackListAndNotify", zap.Error(err))
		return
	}

	_, err = global.GCron.AddFunc("@every 60s", service.ServiceGroup.CronService.CheckCard)
	if err != nil {
		global.GLog.Error("CheckCard", zap.Error(err))
		return
	}

	_, err = global.GCron.AddFunc("01 18 * * ?", service.ServiceGroup.CronService.GetTof)
	if err != nil {
		global.GLog.Error("GetTof", zap.Error(err))
		return
	}

	// 1s
	_, err = global.GCron.AddFunc("@every 1s", flow.CleanExpiredPrivateFlow)
	if err != nil {
		global.GLog.Error("CleanExpiredPrivateFlow", zap.Error(err))
		return
	}

	_, err = global.GCron.AddFunc("@every 1s", flow.CleanExpiredGroupFlow)
	if err != nil {
		global.GLog.Error("CleanExpiredGroupFlow", zap.Error(err))
		return
	}

	_, err = global.GCron.AddFunc(fmt.Sprintf("@every %ds", global.GConfig.QQBot.BotToBot.Interval),
		service.ServiceGroup.CronService.BotToBot)
	if err != nil {
		global.GLog.Error("BotToBot", zap.Error(err))
		return
	}

	if global.GConfig.QQBot.EnableAutoBindGameId {
		_, err = global.GCron.AddFunc("@every 80s", service.ServiceGroup.CronService.AutoBindGameId)
		if err != nil {
			global.GLog.Error("AutoBindGameId", zap.Error(err))
			return
		}
		global.GLog.Info("已启用自动绑定GameId功能")
	}

	if global.GConfig.QQBot.MuteGroup.Enable {
		// 禁言定时器
		if !utils.IsValidTimeFormat(global.GConfig.QQBot.MuteGroup.Start.Time) {
			panic("禁言开始时间配置错误")
		}

		startHour, startMinute := utils.SplitByColon(global.GConfig.QQBot.MuteGroup.Start.Time)

		_, err = global.GCron.AddFunc(fmt.Sprintf("%s %s * * ?", startMinute, startHour),
			service.ServiceGroup.CronService.StartMute)
		if err != nil {
			global.GLog.Error("StartMute", zap.Error(err))
			return
		}

		if !utils.IsValidTimeFormat(global.GConfig.QQBot.MuteGroup.End.Time) {
			panic("禁言结束时间配置错误")
		}

		endHour, endMinute := utils.SplitByColon(global.GConfig.QQBot.MuteGroup.End.Time)

		_, err = global.GCron.AddFunc(fmt.Sprintf("%s %s * * ?", endMinute, endHour),
			service.ServiceGroup.CronService.EndMute)
		if err != nil {
			global.GLog.Error("EndMute", zap.Error(err))
			return
		}
		global.GLog.Info("已启用自动宵禁功能")
	}

	global.GCron.Start()
	global.GLog.Info("定时任务初始化完成")
}
