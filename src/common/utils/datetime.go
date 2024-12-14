package utils

import (
	"fmt"
	"math"
	"time"
)

func GetDate() string {
	now := time.Now()
	date := now.Format("20060102")
	return date
}

func GetDateTime() string {
	now := time.Now()
	dateTime := now.Format("2006-01-02 15:04:05")
	return dateTime
}

func ConvertSecondsToHoursString(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours())
	return fmt.Sprintf("%d小时", hours)
}

// Format 格式化
func Format(userTime time.Time) string {
	// 固定时区
	shanghaiLocation := time.FixedZone("CST", 8*3600)
	userTime = userTime.In(shanghaiLocation)
	return userTime.Format("2006-01-02 15:04:05")
}

// FormatTimestamp 格式化
func FormatTimestamp(ms int64) string {
	dataTime := time.Unix(ms/1000, (ms%1000)*1000000)
	shanghaiLocation := time.FixedZone("CST", 8*3600)
	dataTime = dataTime.In(shanghaiLocation)
	return dataTime.Format("01月02日 15:04")
}

func AbsoluteDurationMinute(a int64, b int64) int {
	return int(math.Abs(float64((a - b) / 1000000 / 60)))
}
