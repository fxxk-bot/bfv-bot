package initialize

import (
	"bfv-bot/common/global"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func Viper() *viper.Viper {
	v := viper.New()

	args := os.Args
	if len(args) == 1 {
		panic("缺少配置文件路径")
	}
	v.SetConfigFile(args[1])
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s \n", err))
	}
	// 解析 配置文件的值注入结构体
	if err = v.Unmarshal(&global.GConfig); err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s \n", err))
	}

	return v
}
