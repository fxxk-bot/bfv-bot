package utils

import (
	"bfv-bot/bot/group"
	"bfv-bot/common/cons"
	"bfv-bot/common/des"
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"bfv-bot/model/dto"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"html"
	"io"
	"strconv"
	"strings"
	"sync"
)

func CheckPlayer(name string) (error, dto.CheckPlayerData) {
	var data dto.CheckPlayerData

	params := make(map[string]string)
	params["name"] = name
	baseInfo, err := http.Get(cons.CheckPlayer, params)
	if err != nil {
		return errors.New("网络异常"), data
	}

	if baseInfo == "" {
		return errors.New("网络异常"), data
	}

	if baseInfo == "{}" {
		return errors.New(cons.PlayerNotFound), dto.CheckPlayerData{}
	}

	var apiResult dto.CheckPlayerResp

	err = des.StringToStruct(baseInfo, &apiResult)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("未找到玩家, 序列化失败"), data
	}

	if apiResult.Success != 1 || apiResult.Code != "player.success" {
		if apiResult.Code == "player.not_found" {
			return errors.New(cons.PlayerNotFound), data
		}
		return errors.New("未找到玩家, 状态: " + apiResult.Code), data
	}

	if apiResult.Data.PersonaID == 0 {
		return errors.New("未找到玩家"), data
	}
	apiResult.Data.PID = strconv.FormatInt(apiResult.Data.PersonaID, 10)
	return nil, apiResult.Data

}

func GetCaptchaBase64() (string, string, error) {
	get, err := http.Get(cons.Captcha, nil)
	if err != nil {
		global.GLog.Error("Get(cons.Captcha, nil)", zap.Error(err))
		return "", "", errors.New("请求失败")
	}
	var captcha dto.CaptchaResp
	err = des.StringToStruct(get, &captcha)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return "", "", errors.New("序列化失败")
	}
	if captcha.Success != 1 || captcha.Code != "captcha.gen" {
		global.GLog.Error("captcha.Success != 1 || captcha.Code != \"captcha.gen\"", zap.String("get", get))
		return "", "", errors.New("验证码请求失败")
	}
	content := html.UnescapeString(captcha.Data.Content)

	png, err := SvgToPng(strings.NewReader(content))
	if err != nil {
		global.GLog.Error("SvgToPng", zap.Error(err))
		return "", "", errors.New("图片转码失败")
	}
	// 读取 io.Reader 的内容到内存中
	data, err := io.ReadAll(png)
	if err != nil {
		global.GLog.Error("io.ReadAll", zap.Error(err))
		return "", "", errors.New("IO异常")
	}

	// 将数据编码为 Base64
	encoded := "base64://" + base64.StdEncoding.EncodeToString(data)
	return encoded, captcha.Data.Hash, nil
}

func GetBanLogByName(name string) (error, []dto.BanLogData) {

	var emptyData []dto.BanLogData

	err, data := CheckPlayer(name)
	if err != nil {
		return err, emptyData
	}
	params := make(map[string]string)
	params["personaId"] = data.PID
	params["limit"] = "5"
	result, err := http.Get(cons.BanLog, params)
	if err != nil {
		global.GLog.Error("result, err := Get(cons.BanLog, params)", zap.Error(err))
		return errors.New("请求失败"), emptyData
	}

	if result == "" || result == "{}" {
		return errors.New("接口异常"), emptyData
	}
	var banlog dto.BanLogResp
	err = des.StringToStruct(result, &banlog)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("序列化异常"), emptyData
	}

	if banlog.Success != 1 || banlog.Code != "getBannedLogsByPersonaId.success" {
		return errors.New("接口响应码异常"), emptyData
	}

	return nil, banlog.Data
}

func GetBanLog(name string) (error, string) {

	err, data := GetBanLogByName(name)
	if err != nil {
		return err, ""
	}

	builder := strings.Builder{}
	length := len(data)

	if length == 0 {
		return nil, "无记录"
	}

	for i, item := range data {
		builder.WriteString(item.ServerName)
		builder.WriteString("\n")

		if item.BanType == 1 || item.BanType == 4 || item.BanType == 5 || item.BanType == 9 {
			builder.WriteString("[" + cons.BanTypeMap[item.BanType] + "]")
		} else if item.BanType == 2 {
			robotReason, err := strconv.Atoi(item.Reason)
			if err == nil {
				value, ok := cons.RobotStatusMap[robotReason]
				if ok {
					builder.WriteString(cons.BanTypeMap[item.BanType] + ": " + value)
				} else {
					builder.WriteString("[" + cons.BanTypeMap[item.BanType] + "]")
				}
			} else {
				builder.WriteString("[" + cons.BanTypeMap[item.BanType] + "]")
			}
		} else if item.BanType == 6 || item.BanType == 7 || item.BanType == 8 {
			builder.WriteString(cons.BanTypeMap[item.BanType] + ": " + item.Reason)
		} else {
			builder.WriteString(item.Reason)
		}

		builder.WriteString("\n")
		builder.WriteString(Format(item.CreateTime))
		if i != length-1 {
			builder.WriteString("\n")
			builder.WriteString("\n")
		}
	}
	return nil, builder.String()
}

func GetBfvRobotServerByName(name string) (error, []dto.ServerListBfvRobotData) {

	var servers []dto.ServerListBfvRobotData

	params := make(map[string]string)
	params["serverName"] = name
	params["region"] = "all"
	params["limit"] = "20"
	params["lang"] = "zh-CN"

	result, err := http.Get(cons.ServerListBfvRobot, params)
	if err != nil {
		global.GLog.Error("result, err := Get(cons.ServerListBfvRobot, params)", zap.String("api result", result), zap.Error(err))
		return errors.New("接口异常"), servers
	}

	if result == "" || result == "{}" {
		global.GLog.Error("result == \"\" || result == \"{}\" server list empty", zap.Error(err))
		return errors.New("接口内容异常"), servers
	}

	var serverlist dto.ServerListBfvRobotResp
	err = des.StringToStruct(result, &serverlist)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("序列化异常"), servers
	}

	if serverlist.Success != 1 || serverlist.Code != "servers.success" {
		return errors.New("接口响应码异常"), servers
	}

	return nil, serverlist.Data
}

func GetGameToolsServerByName(name string) (error, []dto.ServerGameToolsData) {
	params := make(map[string]string)
	params["name"] = name
	params["region"] = "all"
	params["limit"] = "20"
	params["platform"] = "pc"
	result, err := http.Get(cons.ServerListGameTools, params)
	if err != nil {
		global.GLog.Error("Get(cons.ServerListGameTools, params)", zap.String("api result", result), zap.Error(err))
		return err, nil
	}

	if result == "" || result == "{}" {
		global.GLog.Error("gametools server list empty", zap.Error(err))
		return err, nil
	}

	var serverlist dto.ServerGameToolsResp
	err = des.StringToStruct(result, &serverlist)
	if err != nil {
		global.GLog.Error("StringToStruct(result, &serverlist)", zap.Error(err))
		return err, nil
	}
	return nil, serverlist.Servers
}

func GetBfvRobotServer(name string, groupFlag bool) (error, string) {

	ownerMap := make(map[string]bool)

	for _, info := range global.GConfig.Bfv.Server {
		ownerMap[info.OwnerId] = true
	}

	err, servers := GetBfvRobotServerByName(name)
	if err != nil {
		return err, ""
	}

	length := len(servers)
	if length == 0 {
		return errors.New("没开"), ""
	}

	var builder strings.Builder

	for index, item := range servers {
		pid := strconv.FormatInt(item.OwnerID, 10)
		_, ok := ownerMap[pid]
		if ok || groupFlag {
			builder.WriteString(item.ServerName)
			builder.WriteString("   当前地图: [")
			builder.WriteString(item.MapName)
			builder.WriteString("/")
			builder.WriteString(item.MapMode)
			builder.WriteString("]")
			builder.WriteString("   人数: ")
			builder.WriteString(strconv.Itoa(item.Slots.Soldier.Current))
			builder.WriteString("/")
			builder.WriteString(strconv.Itoa(item.Slots.Soldier.Max))
			if item.Slots.Queue.Current == 0 {
				builder.WriteString("\n")
			} else {
				builder.WriteString("[")
				builder.WriteString(strconv.Itoa(item.Slots.Queue.Current))
				builder.WriteString("]\n")
			}
			if index != length-1 {
				builder.WriteString("\n")
			}
		}
	}
	finalResult := builder.String()

	if finalResult == "" {
		return errors.New("没开"), ""
	}

	return nil, finalResult

}

func GetPlayerData(pid string) (error, dto.BfvAllData) {
	var allData dto.BfvAllData
	params := make(map[string]string)
	params["personaId"] = pid

	baseInfo, err := http.Get(cons.PlayerData, params)
	if err != nil || baseInfo == "" {
		return errors.New("接口返回为空"), allData
	}
	var data dto.BfvAllResp
	err = des.ByteToStruct([]byte(baseInfo), &data)
	if err != nil {
		global.GLog.Error("ByteToStruct", zap.Error(err))
		return errors.New("序列化异常"), allData
	}

	if data.Success != 1 || data.Code != "playerAll.success" {
		return errors.New("接口状态码异常"), allData
	}

	return nil, data.Data
}

func GetBfBanStatus(pid string) (error, dto.BfBanStatusData) {
	bfbanParams := map[string]string{
		"personaId": pid,
	}
	result, err := http.Get(cons.BfBanStatus, bfbanParams)
	if err != nil {
		global.GLog.Error("result, err := Get(cons.BfBanStatus, bfbanParams)", zap.Error(err))
		return errors.New("社区接口异常"), dto.BfBanStatusData{Status: -2}
	}

	var bfbanData dto.BfBanStatusResp
	err = des.StringToStruct(result, &bfbanData)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("序列化异常"), dto.BfBanStatusData{Status: -2}
	}

	if bfbanData.Code == "player.notFound" {
		return nil, dto.BfBanStatusData{Status: -1}
	} else if bfbanData.Code == "player.ok" {
		return nil, bfbanData.Data
	} else {
		return errors.New("接口状态码异常"), dto.BfBanStatusData{Status: -2}
	}
}

func GetBfvRobotStatus(pid string) (error, dto.BfvRobotStatusData) {
	params := map[string]string{
		"personaId": pid,
	}

	result, err := http.Get(cons.BfvRobotStatus, params)
	if err != nil {
		global.GLog.Error("result, err := Get(cons.BfvRobotStatus, params)", zap.Error(err))
		return errors.New("社区接口异常"), dto.BfvRobotStatusData{ReasonStatus: -2}
	}
	var bfvRobotApiData dto.BfvRobotStatusResp
	err = des.StringToStruct(result, &bfvRobotApiData)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("序列化异常"), dto.BfvRobotStatusData{ReasonStatus: -2}
	}

	if bfvRobotApiData.Success == 1 {
		formatInt := strconv.FormatInt(bfvRobotApiData.Data.PersonaID, 10)
		if formatInt == pid {
			return nil, bfvRobotApiData.Data
		} else {
			return errors.New("接口数据异常"), dto.BfvRobotStatusData{ReasonStatus: -2}
		}
	} else {
		return errors.New("接口状态码异常"), dto.BfvRobotStatusData{ReasonStatus: -2}
	}

}

func GetActiveTag(pid int64) (error, string) {
	m := make(map[string]interface{})
	pidList := make([]int64, 1)
	pidList[0] = pid
	m["personaIds"] = pidList
	result, err := http.Post(cons.ActiveTag, m)
	if err != nil {
		global.GLog.Error("result, err := Post(cons.ActiveTag, m)", zap.Error(err))
		return errors.New("请求异常"), ""
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), ""
	}
	var activeTags dto.ActiveTagsResp
	err = des.StringToStruct(result, &activeTags)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &activeTags)", zap.Error(err))
		return errors.New("序列化异常"), ""
	}

	if activeTags.Success != 1 || activeTags.Code != "platoonActiveTags.success" {
		return errors.New("接口状态码异常"), ""
	}
	personaId := strconv.FormatInt(pid, 10)
	return nil, activeTags.Data[personaId]
}

func Ban(captchaValue string, captchaHash string, gameId string, reason string, pid int64,
	name string, token string) (error, dto.BanResp) {

	header := make(map[string]string)
	header["x-access-token"] = token
	header["Content-Type"] = "application/json;charset=UTF-8"
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	data := make(map[string]interface{})

	players := make([]dto.BanPlayerReq, 1)
	players[0] = dto.BanPlayerReq{PlayerId: pid, Name: name}

	data["captcha"] = captchaValue
	data["encryptCaptcha"] = captchaHash
	data["gameId"] = gameId
	data["reason"] = reason
	data["players"] = players

	result, err := http.PostByHeader(cons.Ban, header, data)
	if err != nil {
		global.GLog.Error("result, err := PostByHeader(cons.Ban, header, data)", zap.Error(err))
		return err, dto.BanResp{}
	}

	global.GLog.Info("屏蔽结果: " + result)

	var banResult dto.BanResp
	err = des.StringToStruct(result, &banResult)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.BanResp{}
	}
	return nil, banResult
}

func RemoveBan(captchaValue string, captchaHash string, gameId string, pid string, serverName string,
	name string, token string) (error, dto.RemoveBanResp) {
	header := make(map[string]string)
	header["x-access-token"] = token
	header["Content-Type"] = "application/json;charset=UTF-8"
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/126.0.0.0 Safari/537.36"

	data := make(map[string]interface{})

	player := dto.RemoveBanPlayerReq{PlayerId: pid, Name: name, ServerName: serverName}

	data["captcha"] = captchaValue
	data["encryptCaptcha"] = captchaHash
	data["gameId"] = gameId
	data["player"] = player

	result, err := http.PostByHeader(cons.RemoveBan, header, data)
	if err != nil {
		global.GLog.Error("result, err := PostByHeader(cons.RemoveBan, header, data)", zap.Error(err))
		return err, dto.RemoveBanResp{}
	}
	global.GLog.Info("解除屏蔽结果: " + result)

	var removeBanResult dto.RemoveBanResp
	err = des.StringToStruct(result, &removeBanResult)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.RemoveBanResp{}
	}
	return nil, removeBanResult
}

func GetServerPlayerByGameToolsConvert(gameId string) (error, dto.CommonPlayersResp) {
	err, apiPlayers := GetServerPlayerByGameTools(gameId)
	if err != nil {
		return err, dto.CommonPlayersResp{}
	}
	return nil, gameToolsPlayersConvert(&apiPlayers)
}

func GetServerPlayerByGameTools(gameId string) (error, dto.ApiPlayersResp) {
	params := make(map[string]string)
	params["gameid"] = gameId
	playersStr, err := http.Get(cons.ServerPlayerGameTools, params)
	if err != nil {
		global.GLog.Error("playersStr, err := Get(cons.ServerPlayerGameTools, params) 请求错误", zap.Error(err))
		return err, dto.ApiPlayersResp{}
	}

	if strings.Contains(playersStr, "server not found") {
		global.GLog.Error("playersStr, err := Get(cons.ServerPlayerGameTools, params) 玩家列表获取失败 请求错误",
			zap.Error(err))
		return errors.New("server not found"), dto.ApiPlayersResp{}
	}

	var apiPlayers dto.ApiPlayersResp
	err = des.StringToStruct(playersStr, &apiPlayers)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.ApiPlayersResp{}
	}
	return nil, apiPlayers
}

func GetServerPlayerByBfvRobot(gameId string) (error, dto.CommonPlayersResp) {
	params := make(map[string]string)
	params["gameId"] = gameId
	result, err := http.Get(cons.ServerPlayerBfvRobot, params)
	if err != nil {
		global.GLog.Error("result, err := Get(cons.ServerPlayerBfvRobot, params)", zap.Error(err))
		return errors.New("请求失败"), dto.CommonPlayersResp{}
	}

	if result == "" || result == "{}" {
		return errors.New("接口异常"), dto.CommonPlayersResp{}
	}
	var players dto.RobotPlayersResp
	err = des.StringToStruct(result, &players)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return errors.New("序列化异常"), dto.CommonPlayersResp{}
	}

	if players.Success != 1 || players.Message != "players.successful" {
		return errors.New("接口响应码异常"), dto.CommonPlayersResp{}
	}

	return nil, bfvRobotPlayersConvert(&players.Data)
}

func gameToolsPlayersConvert(source *dto.ApiPlayersResp) dto.CommonPlayersResp {
	commonData := dto.CommonPlayersResp{}

	que := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Que {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PlayerID
		target.Join = item.JoinTime
		target.Platoon = item.Platoon
		que = append(que, target)
	}
	commonData.Que = que

	teamOne := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Teams[0].Players {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PlayerID
		target.Join = item.JoinTime
		target.Platoon = item.Platoon
		teamOne = append(teamOne, target)
	}
	commonData.TeamOne = teamOne

	teamTwo := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Teams[1].Players {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PlayerID
		target.Join = item.JoinTime
		target.Platoon = item.Platoon
		teamTwo = append(teamTwo, target)
	}
	commonData.TeamTwo = teamTwo

	return commonData
}

func bfvRobotPlayersConvert(source *dto.RobotPlayersRobotPlayersData) dto.CommonPlayersResp {
	commonData := dto.CommonPlayersResp{}

	que := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Players.Loading {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PersonaID
		target.Join = item.Join
		target.Platoon = item.Platoon
		que = append(que, target)
	}
	commonData.Que = que

	teamOne := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Players.Team1 {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PersonaID
		target.Join = item.Join
		target.Platoon = item.Platoon
		teamOne = append(teamOne, target)
	}
	commonData.TeamOne = teamOne

	teamTwo := make([]dto.CommonPlayersData, 0)
	for _, item := range source.Players.Team2 {
		target := dto.CommonPlayersData{}
		target.Name = item.Name
		target.PersonaID = item.PersonaID
		target.Join = item.Join
		target.Platoon = item.Platoon
		teamTwo = append(teamTwo, target)
	}
	commonData.TeamTwo = teamTwo

	return commonData
}

func GetTof() (error, dto.TofData) {
	result, err := http.Get(cons.GameTof, nil)
	if err != nil {
		global.GLog.Error("Get(cons.GameTof, nil)", zap.Error(err))
		return errors.New("请求异常"), dto.TofData{}
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), dto.TofData{}
	}
	var tofResp dto.TofResp
	err = des.StringToStruct(result, &tofResp)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &activeTags)", zap.Error(err))
		return errors.New("序列化异常"), dto.TofData{}
	}

	if tofResp.Success != 1 || tofResp.Code != "getTOF.success" {
		return errors.New("接口状态码异常"), dto.TofData{}
	}
	return nil, tofResp.Data
}

func GetJoinPlatoons(pid string) (error, []dto.JoinPlatoonData) {

	params := make(map[string]string)
	params["personaId"] = pid
	result, err := http.Get(cons.JoinPlatoons, params)
	if err != nil {
		global.GLog.Error("Get(cons.JoinPlatoons, nil)", zap.Error(err))
		return errors.New("请求异常"), nil
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), nil
	}
	var joinResp dto.JoinPlatoonResp
	err = des.StringToStruct(result, &joinResp)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &activeTags)", zap.Error(err))
		return errors.New("序列化异常"), nil
	}

	if joinResp.Success != 1 || joinResp.Code != "platoonInfo.success" {
		return errors.New("接口状态码异常"), nil
	}
	return nil, joinResp.Data
}

func GetJoinPlatoonsByName(name string) (error, string) {
	err, data := CheckPlayer(name)
	if err != nil {
		return err, ""
	}
	err, arr := GetJoinPlatoons(data.PID)
	if err != nil {
		return err, ""
	}
	length := len(arr)
	if length == 0 {
		return nil, "玩家未加入战排"
	}

	builder := strings.Builder{}
	builder.WriteString("玩家 [")
	builder.WriteString(data.Name)
	builder.WriteString("] 加入的战排有\n\n")

	for i, item := range arr {
		builder.WriteString("名称: [")
		builder.WriteString(item.Name)
		builder.WriteString("] ")
		builder.WriteString("\tTag: [")
		builder.WriteString(item.Tag)
		builder.WriteString("] ")
		builder.WriteString("\t人数: [")
		builder.WriteString(strconv.Itoa(item.Size))
		builder.WriteString("] \n")
		builder.WriteString("描述: ")
		builder.WriteString(item.Description)
		if i != length-1 {
			builder.WriteString("\n\n")
		}
	}

	return nil, builder.String()
}

func GetPlayerBaseInfo(pid string) (error, dto.PlayerBaseInfoData) {

	params := make(map[string]string)
	params["personaId"] = pid
	result, err := http.Get(cons.PlayerBaseInfo, params)
	if err != nil {
		global.GLog.Error("Get(cons.PlayerBaseInfo, nil)", zap.Error(err))
		return errors.New("请求异常"), dto.PlayerBaseInfoData{}
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), dto.PlayerBaseInfoData{}
	}
	var infoResp dto.PlayerBaseInfoResp
	err = des.StringToStruct(result, &infoResp)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &activeTags)", zap.Error(err))
		return errors.New("序列化异常"), dto.PlayerBaseInfoData{}
	}

	if infoResp.Success != 1 || infoResp.Code != "playerStats.success" {
		return errors.New("接口状态码异常"), dto.PlayerBaseInfoData{}
	}
	return nil, infoResp.Data

}

func GetPlayerBaseInfoByName(pid string) (error, string) {
	err, data := GetPlayerBaseInfo(pid)
	if err != nil {
		return err, ""
	}
	builder := strings.Builder{}

	builder.WriteString("等级: \t")
	builder.WriteString(strconv.Itoa(data.BasicStats.Rank.Number))
	builder.WriteString("\n")

	builder.WriteString("KPM: \t")
	builder.WriteString(fmt.Sprintf("%.2f", data.BasicStats.Kpm))
	builder.WriteString("\n")

	var kd float64
	if data.BasicStats.Deaths == 0 {
		kd = float64(data.BasicStats.Kills)
	} else {
		kd = float64(data.BasicStats.Kills) / float64(data.BasicStats.Deaths)
	}
	builder.WriteString("KD: \t\t")
	builder.WriteString(fmt.Sprintf("%.2f", kd))
	builder.WriteString("\n")

	builder.WriteString("SPM: \t")
	builder.WriteString(fmt.Sprintf("%.2f", data.BasicStats.Spm))
	builder.WriteString("\n")

	var hs float64
	if data.BasicStats.Kills == 0 {
		hs = float64(0)
	} else {
		hs = float64(data.HeadShots) / float64(data.BasicStats.Kills)
	}
	builder.WriteString("爆头率: \t")
	builder.WriteString(fmt.Sprintf("%.2f%%", hs*100.0))

	return nil, builder.String()
}

func GetBaseInfoAndStatusByName(data *dto.CheckPlayerData) (error, string) {
	var baseInfo string
	var bfbanInfo string
	var bfvrobotInfo string
	wg := sync.WaitGroup{}
	wg.Add(3)
	_ = global.GPool.Submit(func() {
		defer wg.Done()
		err, baseInfoMsg := GetPlayerBaseInfoByName(data.PID)
		if err == nil {
			baseInfo = baseInfoMsg
		}
	})

	_ = global.GPool.Submit(func() {
		defer wg.Done()
		err, result := GetBfvRobotStatus(data.PID)
		if err == nil {
			bfvrobotInfo = cons.RobotStatusMap[result.ReasonStatus]
		}
	})

	_ = global.GPool.Submit(func() {
		defer wg.Done()
		err, result := GetBfBanStatus(data.PID)
		if err == nil {
			bfbanInfo = cons.BfbanStatusMap[result.Status]
		}
	})

	wg.Wait()

	if baseInfo == "" || bfbanInfo == "" || bfvrobotInfo == "" {
		return errors.New("查询失败"), ""
	}

	return nil, baseInfo + "\n\n" + "BFBAN: \t" + bfbanInfo + "\n" +
		"BFV ROBOT: \t" + bfvrobotInfo
}

func GerServerGroupMember(name string) (error, string) {
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

	gameId := servers[0].GameID

	err, data := GetServerPlayerByBfvRobot(strconv.FormatInt(gameId, 10))
	if err != nil {
		return errors.New("获取服务器玩家列表失败"), ""
	}
	cardMap := group.GetActiveGroupMemberCardMap()

	builder := strings.Builder{}
	builder.WriteString("当前正在游玩 ")
	builder.WriteString(servers[0].ServerName)
	builder.WriteString(" 的群友:\n")

	teamOneLength := len(data.TeamOne)
	builder.WriteString("Team One\n")
	for index, item := range data.TeamOne {
		_, ok := cardMap[item.Name]
		if ok {
			builder.WriteString("\t")
			builder.WriteString(item.Name)
			if index != teamOneLength-1 {
				builder.WriteString("\n")
			}
		}
	}

	teamTwoLength := len(data.TeamTwo)
	builder.WriteString("\nTeam Two\n")
	for index, item := range data.TeamTwo {
		_, ok := cardMap[item.Name]
		if ok {
			builder.WriteString("\t")
			builder.WriteString(item.Name)
			if index != teamTwoLength-1 {
				builder.WriteString("\n")
			}
		}
	}

	return nil, builder.String()
}

func GetBfBanBatchStatus(pidArr []int64) (error, []dto.BfBanBatchData) {
	params := make(map[string]interface{})
	params["personaIds"] = pidArr
	result, err := http.Post(cons.BfBanBatchStatus, params)
	if err != nil {
		global.GLog.Error("Get(cons.BfBanBatchStatus, params)", zap.Error(err))
		return errors.New("请求异常"), []dto.BfBanBatchData{}
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), []dto.BfBanBatchData{}
	}
	var apiResp dto.BfBanBatchResp
	err = des.StringToStruct(result, &apiResp)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &apiResp)", zap.Error(err))
		return errors.New("序列化异常"), []dto.BfBanBatchData{}
	}

	if apiResp.Success != 1 || apiResp.Code != "playerBatch.ok" {
		return errors.New("接口状态码异常"), []dto.BfBanBatchData{}
	}
	return nil, apiResp.Data
}

func GetBfvRobotBatchStats(pidArr []int64) (error, []dto.BfvRobotBatchStatsData) {
	params := make(map[string]interface{})
	params["personaIds"] = pidArr
	result, err := http.Post(cons.BfvRobotBatchStatus, params)
	if err != nil {
		global.GLog.Error("Get(cons.BfvRobotBatchStatus, params)", zap.Error(err))
		return errors.New("请求异常"), []dto.BfvRobotBatchStatsData{}
	}
	if result == "" || result == "{}" {
		return errors.New("响应异常"), []dto.BfvRobotBatchStatsData{}
	}
	var apiResp dto.BfvRobotBatchStatsResp
	err = des.StringToStruct(result, &apiResp)
	if err != nil {
		global.GLog.Error("err = StringToStruct(result, &apiResp)", zap.Error(err))
		return errors.New("序列化异常"), []dto.BfvRobotBatchStatsData{}
	}

	if apiResp.Success != 1 || apiResp.Code != "playerGrpcStats.success" {
		return errors.New("接口状态码异常"), []dto.BfvRobotBatchStatsData{}
	}
	return nil, apiResp.Data
}

func GetGameToolsBatchStatus(arr []int64) (error, []dto.GtBatchStatusData) {
	result, err := http.Post(cons.GameToolsBatchStatus+"?raw=false&format_values=true", arr)
	if err != nil {
		global.GLog.Error("Get(cons.GameToolsBatchStatus, arr)",
			zap.String("api result", result), zap.Error(err))
		return err, nil
	}

	if result == "" || result == "{}" {
		global.GLog.Error("gametools batch list empty", zap.Error(err))
		return err, nil
	}

	var list dto.GtBatchStatusResp
	err = des.StringToStruct(result, &list)
	if err != nil {
		global.GLog.Error("StringToStruct(result, &list)", zap.Error(err))
		return err, nil
	}
	return nil, list.Data
}
