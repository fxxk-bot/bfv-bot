package config

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type Zap struct {
	// 级别
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 日志前缀
	Prefix string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	// 输出
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// 日志文件夹
	Director string `mapstructure:"director" json:"director"  yaml:"director"`
	// 编码级
	EncodeLevel string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`
	// 栈名
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	// 显示行
	ShowLine bool `mapstructure:"show-line" json:"show-line" yaml:"show-line"`
	// 输出控制台
	LogInConsole bool `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
	// 日志留存时间
	MaxAge int `mapstructure:"max-age" yaml:"max-age"`
}

// Levels 根据字符串转化为 zapcore.Levels
func (c *Zap) Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	return levels
}

func (c *Zap) Encoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		NameKey:       "name",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: c.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(c.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
		},
		EncodeLevel:    c.LevelEncoder(),
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	if c.Format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)

}

// LevelEncoder 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author [SliverHorn](https://github.com/SliverHorn)
func (c *Zap) LevelEncoder() zapcore.LevelEncoder {
	switch {
	case c.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case c.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case c.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case c.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}
