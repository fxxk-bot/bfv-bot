package po

type CardCheck struct {
	Qq      int64 `gorm:"primaryKey;column:qq"`
	GroupId int64 `gorm:"column:group_id"`
	// id检测失败次数 接口异常不累加 接口返回{}累加
	FailCnt       int   `gorm:"column:fail_cnt"`
	NextCheckTime int64 `gorm:"column:next_check_time"`
}

func (CardCheck) TableName() string {
	return "card_check"
}
