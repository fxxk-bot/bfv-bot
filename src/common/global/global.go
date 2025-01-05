package global

import (
	"bfv-bot/common/config"
	"bfv-bot/model/dto"
	"bfv-bot/model/po"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/feiin/sensitivewords"
	"github.com/go-rod/rod"
	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var (
	GConfig           config.Config
	GLog              *zap.Logger
	GDb               *gorm.DB
	GCron             *cron.Cron
	GPool             *ants.Pool
	GSensitive        *sensitivewords.SensitiveWords
	GAi               *qianfan.ChatCompletion
	GBlackListMap     map[string]po.Blacklist
	GJoinBlackListMap map[int64]string
	GBindMap          map[int64]string
	GRodBrowser       *rod.Browser
	GTofData          dto.TofData
	GTofDataCache     sync.Map
	GResourceCache    sync.Map
)
