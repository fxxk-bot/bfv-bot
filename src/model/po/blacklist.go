package po

type Blacklist struct {
	Id     string `gorm:"primaryKey;column:id"`
	Name   string `gorm:"column:name"`
	Reason string `gorm:"column:reason"`
}

func (Blacklist) TableName() string {
	return "blacklist"
}
