package utils

import (
	"bfv-bot/bot/group"
	"bfv-bot/common/cons"
	"bfv-bot/common/global"
	"bfv-bot/common/utils/cache"
	"bfv-bot/model/dto"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"
)

func GetPlayerList(name string) (error, string) {

	if err := recover(); err != nil {
		global.GLog.Error("GetPlayerList", zap.Any("info", err))
	}

	err := CreateOutputDir()
	if err != nil {
		global.GLog.Error("os.Mkdir", zap.Error(err))
		return errors.New("服务器文件异常"), ""
	}

	err, servers := GetBfvRobotServerByName(name)
	if err != nil {
		return errors.New("搜索服务器失败"), ""
	}
	if len(servers) == 0 {
		return errors.New("未搜索到服务器"), ""
	}

	if len(servers) > 1 {
		return errors.New("搜索到的服务器过多"), ""
	}
	server := servers[0]
	gameId := server.GameID

	data := make(map[string]interface{})

	var serverInfo dto.TemplatePlayerServerInfoModel
	serverInfo.ServerName = server.ServerName
	serverInfo.MapName = server.MapName
	serverInfo.MapMode = server.MapMode

	path, err := cache.Get(server.URL)
	if err != nil {
		return errors.New("资源文件读取失败"), ""
	}

	base64, err := ImageToBase64(path)
	if err != nil {
		return errors.New("资源文件编码失败"), ""
	}

	serverInfo.ImageBase64 = base64

	data["ServerInfo"] = serverInfo

	err, playerlist := GetServerPlayerByBfvRobot(strconv.FormatInt(gameId, 10))
	if err != nil {
		return errors.New("获取服务器玩家列表失败"), ""
	}

	personaIdArr := make([]int64, 0)
	for _, item := range playerlist.TeamOne {
		personaIdArr = append(personaIdArr, item.PersonaID)
	}
	for _, item := range playerlist.TeamTwo {
		personaIdArr = append(personaIdArr, item.PersonaID)
	}

	bfvrobotMap := make(map[int64]dto.BfvRobotBatchStatsData)
	bfbanMap := make(map[int64]dto.BfBanBatchData)

	wg := sync.WaitGroup{}
	wg.Add(2)
	// base stat
	_ = global.GPool.Submit(func() {
		defer wg.Done()
		err, bfvRobotBatchStats := GetBfvRobotBatchStats(personaIdArr)
		if err != nil {
			global.GLog.Error("GetBfvRobotBatchStats", zap.Any("info", err))
			return
		}
		for _, item := range bfvRobotBatchStats {
			bfvrobotMap[item.PersonaID] = item
		}
	})

	// bfban
	_ = global.GPool.Submit(func() {
		defer wg.Done()

		err, bfBanBatchStatus := GetBfBanBatchStatus(personaIdArr)
		if err != nil {
			global.GLog.Error("GetBfBanBatchStatus", zap.Any("info", err))
			return
		}
		for _, item := range bfBanBatchStatus {
			bfbanMap[item.PersonaID] = item
		}
	})

	wg.Wait()

	if len(bfvrobotMap) == 0 {
		return errors.New("玩家信息获取失败"), ""
	}

	cardMap := make(map[string]bool)
	if global.GConfig.QQBot.EnablePlayerlistShowGroupMember {
		cardMap = group.GetActiveGroupMemberCardMap()
	}

	unixMicro := time.Now().UnixMicro()

	var teamOne dto.TemplatePlayerTeamModel
	teamOne.TeamName = "Team A"
	teamOneList := make([]dto.PlayerlistData, 0)

	for _, item := range playerlist.TeamOne {

		if len(teamOneList) > 32 {
			break
		}

		var target dto.PlayerlistData
		if item.Platoon != "" {
			target.Name = "[" + item.Platoon + "] " + item.Name
		} else {
			target.Name = item.Name
		}
		if global.GConfig.QQBot.EnablePlayerlistShowGroupMember {
			_, ok := cardMap[item.Name]
			target.IsGroupMember = ok
		}

		duration := AbsoluteDurationMinute(item.Join, unixMicro)
		target.Join = strconv.Itoa(duration)

		statsData, ok := bfvrobotMap[item.PersonaID]
		if ok {
			target.KillDeath = statsData.KillDeath
			target.KillsPerMinute = statsData.KillsPerMinute
			target.Rank = statsData.Rank
		}

		bfbanData, ok := bfbanMap[item.PersonaID]
		if ok {
			target.BfBanStatus = bfbanData.Status
			target.BfBanStatusName = cons.BfbanStatusMap[bfbanData.Status]
		}

		teamOneList = append(teamOneList, target)
	}
	teamOne.List = teamOneList

	var teamTwo dto.TemplatePlayerTeamModel

	teamTwo.TeamName = "Team B"
	teamTwoList := make([]dto.PlayerlistData, 0)

	for _, item := range playerlist.TeamTwo {

		if len(teamTwoList) > 32 {
			break
		}

		var target dto.PlayerlistData

		if item.Platoon != "" {
			target.Name = "[" + item.Platoon + "] " + item.Name
		} else {
			target.Name = item.Name
		}

		if global.GConfig.QQBot.EnablePlayerlistShowGroupMember {
			_, ok := cardMap[item.Name]
			target.IsGroupMember = ok
		}

		duration := AbsoluteDurationMinute(item.Join, unixMicro)
		target.Join = strconv.Itoa(duration)

		statsData, ok := bfvrobotMap[item.PersonaID]
		if ok {
			target.KillDeath = statsData.KillDeath
			target.KillsPerMinute = statsData.KillsPerMinute
			target.Rank = statsData.Rank
		}

		bfbanData, ok := bfbanMap[item.PersonaID]
		if ok {
			target.BfBanStatus = bfbanData.Status
			target.BfBanStatusName = cons.BfbanStatusMap[bfbanData.Status]
		}

		teamTwoList = append(teamTwoList, target)
	}
	teamTwo.List = teamTwoList

	data["TeamOne"] = teamOne
	data["TeamTwo"] = teamTwo

	htmlContent, err := os.ReadFile(global.GConfig.Server.Template.Playerlist)
	if err != nil {
		global.GLog.Error("os.ReadFile", zap.Error(err))
		return errors.New("模板文件读取失败"), ""
	}
	t, err := template.New("playerlist").Parse(string(htmlContent))
	if err != nil {
		return errors.New("模板文件解析失败"), ""
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		return errors.New("模板文件渲染失败"), ""
	}

	imagePath := fmt.Sprintf(global.GConfig.Server.Output+"/%s/%d_%s_playerlist.jpg", GetDate(), gameId, GetUUID())
	htmlPath := imagePath + ".html"

	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return errors.New("模板文件创建失败"), ""
	}

	_, err = htmlFile.Write(buf.Bytes())
	if err != nil {
		return errors.New("模板文件写入失败"), ""
	}

	// rod渲染
	page := global.GRodBrowser.MustPage("file://" + htmlPath)
	defer page.MustClose()
	page.MustSetViewport(1920, 1080, 1, false)
	page.MustWaitStable()
	err = page.WaitLoad()
	if err != nil {
		global.GLog.Error("page.WaitLoad", zap.Error(err))
		return errors.New("模板渲染异常"), ""
	}

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:      proto.PageCaptureScreenshotFormatJpeg,
		FromSurface: false,
		Clip: &proto.PageViewport{
			X:      15,
			Y:      15,
			Width:  1905,
			Height: 1065,
			Scale:  1.0,
		},
	})
	if err != nil {
		global.GLog.Error("page.Screenshot", zap.Error(err))
		return errors.New("模板渲染结果获取异常"), ""
	}
	jpgFile, err := os.Create(imagePath)
	if err != nil {
		global.GLog.Error("os.Create", zap.Error(err), zap.String("path", imagePath))
		return errors.New("模板渲染结果存储路径读取失败"), ""
	}
	_, err = jpgFile.Write(img)
	if err != nil {
		global.GLog.Error("file.Write", zap.Error(err))
		return errors.New("模板渲染结果写入失败"), ""
	}

	defer func() {
		err := jpgFile.Close()
		if err != nil {
			global.GLog.Error("file.Close()", zap.Error(err))
			return
		}
		err = htmlFile.Close()
		if err != nil {
			global.GLog.Error("file.Close()", zap.Error(err))
			return
		}
		err = os.Remove(htmlPath)
		if err != nil {
			global.GLog.Error("os.RemoveBlack", zap.Error(err))
		}
	}()

	return nil, imagePath
}
