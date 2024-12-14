package config

type Server struct {
	Port           int      `mapstructure:"port" yaml:"port"`
	GinMode        string   `mapstructure:"gin-mode" yaml:"gin-mode"`
	Resource       string   `mapstructure:"resource" yaml:"resource"`
	Output         string   `mapstructure:"output" yaml:"output"`
	ResourcesCache string   `mapstructure:"resources-cache" yaml:"resources-cache"`
	Font           string   `mapstructure:"font" yaml:"font"`
	Template       Template `mapstructure:"template" yaml:"template"`
	DbType         string   `mapstructure:"db-type" yaml:"db-type"`
}

type Template struct {
	Data       string `mapstructure:"data" yaml:"data"`
	Task       string `mapstructure:"task" yaml:"task"`
	Playerlist string `mapstructure:"playerlist" yaml:"playerlist"`
}
