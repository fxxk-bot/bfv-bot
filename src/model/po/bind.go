package po

type Bind struct {
	Qq   int64  `gorm:"primaryKey;column:qq"`
	Name string `gorm:"column:name"`
	Pid  string `gorm:"column:pid"`
}

func (Bind) TableName() string {
	return "bind"
}
