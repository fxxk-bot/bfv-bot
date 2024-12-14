package config

type Config struct {
	// mysql
	Mysql Mysql `mapstructure:"mysql" yaml:"mysql"`
	// sqlite
	Sqlite Sqlite `mapstructure:"sqlite" yaml:"sqlite"`
	// Zap
	Zap Zap `mapstructure:"zap" yaml:"zap"`
	// Server
	Server Server `mapstructure:"server" yaml:"server"`
	// QQBot
	QQBot QQBot `mapstructure:"qq-bot" yaml:"qq-bot"`
	// Bfv
	Bfv Bfv `mapstructure:"bfv" yaml:"bfv"`
	// ai
	Ai Ai `mapstructure:"ai" yaml:"ai"`
}
