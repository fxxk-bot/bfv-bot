package config

type Mysql struct {
	// 服务器地址
	Url string `mapstructure:"url" yaml:"url"`
	// 端口
	Port string `mapstructure:"port" yaml:"port"`
	// 高级配置
	Config string `mapstructure:"config" yaml:"config"`
	// 数据库名
	DbName string `mapstructure:"db-name" yaml:"db-name"`
	// 数据库用户名
	Username string `mapstructure:"username" yaml:"username"`
	// 数据库密码
	Password string `mapstructure:"password" yaml:"password"`
	// 全局表前缀，单独定义TableName则不生效
	Prefix string `mapstructure:"prefix" yaml:"prefix"`
	// 是否开启全局禁用复数，true表示开启
	Singular bool `mapstructure:"singular" yaml:"singular"`
	// 数据库引擎，默认InnoDB
	Engine string `mapstructure:"engine" yaml:"engine" default:"InnoDB"`
	// 空闲中的最大连接数
	MaxIdleConns int `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`
	// 打开到数据库的最大连接数
	MaxOpenConns int `mapstructure:"max-open-conns" yaml:"max-open-conns"`
	// 是否开启Gorm全局日志
	LogMode string `mapstructure:"log-mode" yaml:"log-mode"`
	// 是否通过zap写入日志文件
	LogZap bool `mapstructure:"log-zap" yaml:"log-zap"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Url + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}
