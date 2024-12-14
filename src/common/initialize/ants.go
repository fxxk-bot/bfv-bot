package initialize

import (
	"bfv-bot/common/global"
	"github.com/panjf2000/ants/v2"
)

func Ants() {
	pool, err := ants.NewPool(3000)
	if err != nil {
		panic(err)
	}
	global.GPool = pool
}
