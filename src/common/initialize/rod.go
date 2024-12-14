package initialize

import (
	"bfv-bot/common/global"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func InitRod() {
	sandbox := launcher.New().NoSandbox(true)
	launch := sandbox.MustLaunch()
	global.GRodBrowser = rod.New().ControlURL(launch).MustConnect()
}
