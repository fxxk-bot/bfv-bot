package config

type Ai struct {
	Enable    bool   `mapstructure:"enable" yaml:"enable"`
	ModelName string `mapstructure:"model-name" yaml:"model-name"`
	AccessKey string `mapstructure:"access-key" yaml:"access-key"`
	SecretKey string `mapstructure:"secret-key" yaml:"secret-key"`
}
