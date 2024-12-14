package cache

import (
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path"
)

func Get(url string) (string, error) {
	value, ok := global.GResourceCache.Load(url)

	if ok {
		return value.(string), nil
	}
	resourcePath, err := put(url)
	if err != nil {
		return "", err
	}
	global.GResourceCache.Store(url, resourcePath)
	return resourcePath, nil
}

func put(resourceUrl string) (string, error) {

	u, err := url.Parse(resourceUrl)
	if err != nil {
		global.GLog.Error("url.Parse(resourceUrl)", zap.String("resourceUrl", resourceUrl), zap.Error(err))
		return "", errors.New("图片获取失败")
	}

	hash := md5.New()
	hash.Write([]byte(resourceUrl))
	md5Value := hex.EncodeToString(hash.Sum(nil))

	fileName := path.Base(u.Path)
	fileExt := path.Ext(fileName)

	cachePath := global.GConfig.Server.ResourcesCache + "/" + md5Value + fileExt

	// 获取文件信息
	_, err = os.Stat(cachePath)

	if err == nil {
		return cachePath, nil
	}

	resp, err := http.Get(resourceUrl, nil)
	if err != nil {
		global.GLog.Error("http.Get", zap.String("resourceUrl", resourceUrl), zap.Error(err))
		return "", errors.New("资源获取失败")
	}

	file, err := os.Create(cachePath)
	if err != nil {
		global.GLog.Error("os.Create", zap.Error(err), zap.String("path", cachePath))
		return "", errors.New("资源缓存失败")
	}
	defer file.Close()

	_, err = file.Write([]byte(resp))
	if err != nil {
		global.GLog.Error("file.Write", zap.Error(err))
		return "", errors.New("资源缓存写入失败")
	}

	return cachePath, nil
}
