package utils

import (
	"bfv-bot/common/global"
	"bfv-bot/model/dto"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetTaskAndCache 获取任务信息并缓存
// offset 0 当前周任务 -1上周任务 1下周任务
func GetTaskAndCache(offset int) (string, error) {

	if err := recover(); err != nil {
		global.GLog.Error("GetTaskAndCache", zap.Any("info", err))
	}

	err := CreateOutputDir()
	if err != nil {
		global.GLog.Error("os.Mkdir", zap.Error(err))
		return "", errors.New("服务器文件异常")
	}

	var currentWeek int
	now := time.Now()
	for index, item := range global.GTofData.Events[0].Weeks {
		start, _ := strconv.ParseInt(item.StartTimestamp, 10, 64)
		end, _ := strconv.ParseInt(item.EndTimestamp, 10, 64)
		if start < now.UnixMilli() && now.UnixMilli() < end {
			currentWeek = index
		}
	}

	cacheKey := currentWeek + offset
	if cacheKey < 0 {
		return "", errors.New("前面没有数据了")
	}

	if len(global.GTofData.Events[0].Weeks) <= cacheKey {
		return "", errors.New("后面没有数据了")
	}

	// 读取已经渲染过的图片路径返回
	path, ok := global.GTofDataCache.Load(cacheKey)
	if ok {
		return path.(string), nil
	}

	week := global.GTofData.Events[0].Weeks[cacheKey]

	data := make(map[string]interface{})

	start, _ := strconv.ParseInt(week.StartTimestamp, 10, 64)
	end, _ := strconv.ParseInt(week.EndTimestamp, 10, 64)

	data["StartTime"] = FormatTimestamp(start)
	data["EndTime"] = FormatTimestamp(end)

	nodeData := make(map[string]dto.Node)
	idToPosition := make(map[string]dto.Position)

	for _, item := range week.StoryEvents {
		idToPosition[item.Achievement.ID] = item.Position
	}

	for _, item := range week.StoryEvents {
		var node dto.Node
		node.Name = item.Achievement.Name
		node.Position = item.Position

		dependencies := make([]dto.Position, 0)
		for _, item := range item.Achievement.Dependencies {
			position, ok := idToPosition[item]
			if ok {
				dependencies = append(dependencies, position)
			}
		}
		jsonData, _ := json.Marshal(dependencies)
		node.Dependencies = base64.StdEncoding.EncodeToString(jsonData)

		requirements := make([]string, 0)
		for _, item := range item.Achievement.Requirements {
			desc := strings.ReplaceAll(item.Desc, "{0:d}", item.RequiredValue)
			requirements = append(requirements, desc)
		}
		node.Requirements = requirements

		rewards := make([]string, 0)
		for _, item := range item.Achievement.Rewards {
			var unit string
			if item.ItemType == "grindCurrency" {
				unit = "连队币"
			} else if item.ItemType == "premiumCurrency" {
				unit = "金币"
			} else {
				unit = ""
			}
			rewards = append(rewards, fmt.Sprintf("%s %s", item.Quantity, unit))
		}
		node.Rewards = rewards

		nodeData[fmt.Sprintf("n%d_%d", item.Position.X, item.Position.Y)] = node
	}

	ys := []int{1, 2, 3}
	xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["Rows"] = dto.Rows{Xs: xs, Ys: ys}

	for _, x := range xs {
		for _, y := range ys {
			mapKey := fmt.Sprintf("n%d_%d", x, y)
			_, ok := nodeData[mapKey]
			if !ok {
				nodeData[mapKey] = dto.Node{Name: "",
					Position:     dto.Position{X: x, Y: y},
					Dependencies: "",
					Requirements: make([]string, 0),
					Rewards:      make([]string, 0),
				}
			}
		}
	}

	data["NodeData"] = nodeData

	htmlContent, err := os.ReadFile(global.GConfig.Server.Template.Task)
	if err != nil {
		global.GLog.Error("os.ReadFile", zap.Error(err))
		return "", errors.New("模板文件读取失败")
	}
	t, err := template.New("task").Parse(string(htmlContent))
	if err != nil {
		return "", errors.New("模板文件解析失败")
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		return "", errors.New("模板文件渲染失败")
	}

	imagePath := fmt.Sprintf(global.GConfig.Server.Output+"/%s/%d_%s_task.jpg", GetDate(), cacheKey, GetUUID())
	htmlPath := imagePath + ".html"
	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return "", errors.New("模板文件创建失败")
	}

	_, err = htmlFile.Write(buf.Bytes())
	if err != nil {
		return "", errors.New("模板文件写入失败")
	}

	// rod渲染

	page := global.GRodBrowser.MustPage("data:text/html;charset=utf-8," + buf.String())
	defer page.MustClose()
	page.MustSetViewport(1920, 1076, 1, false)
	page.MustWaitStable()
	err = page.WaitLoad()
	if err != nil {
		global.GLog.Error("page.WaitLoad", zap.Error(err))
		return "", errors.New("模板渲染异常")
	}

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:      proto.PageCaptureScreenshotFormatJpeg,
		FromSurface: false,
	})
	if err != nil {
		global.GLog.Error("page.Screenshot", zap.Error(err))
		return "", errors.New("模板渲染结果获取异常")
	}
	file, err := os.Create(imagePath)
	if err != nil {
		global.GLog.Error("os.Create", zap.Error(err), zap.String("path", imagePath))
		return "", errors.New("模板渲染结果存储路径读取失败")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.GLog.Error("file.Close()", zap.Error(err))
		}
		err = htmlFile.Close()
		if err != nil {
			global.GLog.Error("htmlFile.Close()", zap.Error(err))
			return
		}
		err = os.Remove(htmlPath)
		if err != nil {
			global.GLog.Error("os.Remove(htmlPath)", zap.Error(err))
			return
		}
	}(file)

	_, err = file.Write(img)
	if err != nil {
		global.GLog.Error("file.Write", zap.Error(err))
		return "", errors.New("模板渲染结果写入失败")
	}

	global.GTofDataCache.Store(cacheKey, imagePath)

	return imagePath, nil
}
