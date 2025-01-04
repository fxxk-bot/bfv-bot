package http

import (
	"bfv-bot/common/global"
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func PostByHeader(urlPar string, header map[string]string, data interface{}) (string, error) {
	// 将数据编码为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		global.GLog.Error("json.Marshal", zap.Error(err))
		return "", err
	}

	req, err := http.NewRequest("POST", urlPar, bytes.NewBuffer(jsonData))
	if err != nil {
		global.GLog.Error("http.NewRequest", zap.Error(err))
		return "", err
	}

	// 请求头
	for key, value := range header {
		req.Header.Set(key, value)
	}
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 创建 http.Client 并设置超时时间为 5 秒
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		global.GLog.Error("client.Do", zap.Error(err))
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GLog.Error("Body.Close", zap.Error(err))
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	stringResult := string(body)
	if stringResult == "" {
		return "", errors.New("响应为空")
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		global.GLog.Debug("Post", zap.String("url", urlPar), zap.String("stringResult", stringResult))
	}

	return stringResult, nil
}

func Post(urlPar string, data interface{}) (string, error) {
	m := make(map[string]string)
	m["Content-Type"] = "application/json;charset=UTF-8"
	m["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	return PostByHeader(urlPar, m, data)
}

func getByHeader(baseURL string, header map[string]string, queryParams map[string]string) (string, error) {

	// 创建 URL 对象
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// 设置查询参数
	if queryParams != nil {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// 创建 http.Client 并设置超时时间为 5 秒
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 创建 GET 请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	// 请求头
	for key, value := range header {
		req.Header.Set(key, value)
	}
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	stringResult := string(body)
	if stringResult == "" {
		return "", errors.New("响应为空")
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		global.GLog.Debug("Get", zap.String("url", baseURL), zap.String("stringResult", stringResult))
	}
	return stringResult, nil
}

func Get(baseURL string, queryParams map[string]string) (string, error) {
	m := make(map[string]string)
	m["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	return getByHeader(baseURL, m, queryParams)
}
