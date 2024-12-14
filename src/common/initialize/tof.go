package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
)

func InitTofData() {
	err, data := utils.GetTof()

	if err != nil {
		global.GLog.Error("tof数据读取失败")
	}

	global.GTofData = data
}
