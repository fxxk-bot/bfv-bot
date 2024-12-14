package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/model/po"
)

func InitDb() {
	if global.GConfig.Server.DbType == "mysql" {
		global.GDb = GormMysql()
	} else if global.GConfig.Server.DbType == "sqlite" {
		global.GDb = GormSqlite()
	} else {
		panic("db-type error")
	}

	// 自动创建表
	err := global.GDb.AutoMigrate(&po.Blacklist{}, &po.Ignorelist{}, &po.Sensitive{}, &po.Bind{}, &po.CardCheck{})
	if err != nil {
		panic(err)
	}
}
