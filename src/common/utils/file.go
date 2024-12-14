package utils

import (
	"bfv-bot/common/global"
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateOutputDir() error {
	dir := fmt.Sprintf(global.GConfig.Server.Output+"/%s/", GetDate())
	// 创建文件夹
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(dir, 0755)
		if err != nil {
			global.GLog.Error("os.Mkdir", zap.Error(err))
			return err
		}
	}
	return nil
}

// 根据文件扩展名确定 MIME 类型
func getMimeType(extension string) string {
	switch extension {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

func ImageToBase64(filePath string) (string, error) {
	// 读取文件内容
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 获取文件扩展名并确定 MIME 类型
	ext := filepath.Ext(filePath)
	mimeType := getMimeType(ext)

	// 将文件内容转换为 Base64 编码并添加前缀
	base64Data := base64.StdEncoding.EncodeToString(fileData)
	base64String := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return base64String, nil
}
