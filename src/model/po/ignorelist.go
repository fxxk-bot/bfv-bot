package po

type Ignorelist struct {
	Id string `gorm:"primaryKey;column:id"`
}

func (Ignorelist) TableName() string {
	return "ignorelist"
}
