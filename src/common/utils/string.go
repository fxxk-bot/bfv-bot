package utils

import (
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
)

var timeRegex = regexp.MustCompile(`^(?:[01]?[0-9]|2[0-3]):[0-5][0-9]$`)

func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func GetUUID() string {
	id := uuid.New()
	idString := strings.ReplaceAll(id.String(), "-", "")
	return idString
}

func GetCommandKeyValue(command string) (string, string) {
	parts := strings.SplitN(command, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	parts = strings.SplitN(command, "＝", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	if strings.HasPrefix(command, ".") && len(command) > 2 {
		runes := []rune(command)
		result := string(runes[1:])
		parts = strings.SplitN(result, " ", 2)
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}
	if strings.HasPrefix(command, "/") && len(command) > 2 {
		runes := []rune(command)
		result := string(runes[1:])
		parts = strings.SplitN(result, " ", 2)
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}
	return "", ""
}

func SplitByColon(str string) (string, string) {
	parts := strings.Split(str, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// IsValidTimeFormat 校验字符串是否符合 "HH:MM" 格式
func IsValidTimeFormat(timeStr string) bool {
	return timeRegex.MatchString(timeStr)
}
