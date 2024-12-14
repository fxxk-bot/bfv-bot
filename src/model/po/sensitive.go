package po

type Sensitive struct {
	Id string `gorm:"primaryKey;column:id"`
}

func (Sensitive) TableName() string {
	return "sensitive"
}
