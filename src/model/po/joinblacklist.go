package po

type JoinBlackList struct {
	Qq     int64  `gorm:"primaryKey;column:qq"`
	Reason string `gorm:"column:reason"`
}

func (JoinBlackList) TableName() string {
	return "join_blacklist"
}
