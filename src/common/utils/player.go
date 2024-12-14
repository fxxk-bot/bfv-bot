package utils

import (
	"bfv-bot/common/cons"
	"bfv-bot/common/global"
	"bfv-bot/model/dto"
	"bytes"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-rod/rod/lib/proto"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func setFontSize(dc *gg.Context, points float64) {
	_ = dc.LoadFontFace(global.GConfig.Server.Font, points)
}
func countFilesInDir() int {
	var count int
	_ = filepath.Walk(global.GConfig.Server.Resource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func drawOverview(data *dto.BfvAllData, nickname string, bfvrobot int, bfban int) (string, error) {
	if err := recover(); err != nil {
		global.GLog.Error("drawOverview recover", zap.Any("info", err))
	}

	err := CreateOutputDir()
	if err != nil {
		global.GLog.Error("os.Mkdir", zap.Error(err))
		return "", errors.New("服务器文件异常")
	}

	count := countFilesInDir()

	randomInt := RandomInt(count)
	global.GLog.Debug("", zap.Int("图片数量", count), zap.Int("随机图片序号", randomInt))
	// 打开一个图片文件
	file, err := os.Open(global.GConfig.Server.Resource + fmt.Sprintf("/%d.jpg", randomInt))
	if err != nil {
		global.GLog.Error("os.Open", zap.Error(err))
		return "", errors.New("背景图片打开失败")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.GLog.Error("file.Close", zap.Error(err))
		}
	}(file)

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		global.GLog.Error("image.Decode", zap.Error(err))
		return "", errors.New("背景图片解码异常")
	}

	// 创建一个新的上下文
	dc := gg.NewContextForImage(img)

	// 设置绘图颜色
	dc.SetColor(color.RGBA{1, 1, 1, 170})

	// 基础信息圆角矩形
	dc.DrawRoundedRectangle(30, 15, 700, 120, 10)
	dc.Fill()

	// bfban状态圆角矩形
	dc.DrawRoundedRectangle(750, 15, 430, 120, 10)
	dc.Fill()

	// 概览数据圆角矩形
	dc.DrawRoundedRectangle(30, 145, 700, 150, 10)
	dc.Fill()

	// 解包武器圆角矩形
	dc.DrawRoundedRectangle(750, 145, 430, 150, 10)
	dc.Fill()

	// 枪械数据圆角矩形
	dc.DrawRoundedRectangle(30, 305, 700, 400, 10)
	dc.Fill()

	// 载具数据圆角矩形
	dc.DrawRoundedRectangle(750, 305, 430, 400, 10)
	dc.Fill()

	dc.SetColor(color.White)
	dc.SetLineWidth(5)
	dc.DrawLine(960, 22, 960, 33)
	dc.DrawLine(960, 44, 960, 58)
	dc.DrawLine(960, 69, 960, 80)
	dc.DrawLine(960, 91, 960, 105)
	dc.DrawLine(960, 116, 960, 127)
	dc.SetLineWidth(1)

	// 设置字体
	if err := dc.LoadFontFace(global.GConfig.Server.Font, 15); err != nil {
		global.GLog.Error("dc.LoadFontFace", zap.Error(err))
		return "", errors.New("字体异常")
	}
	// 设置文字颜色
	dc.SetColor(color.White)

	// 当前时间
	dc.DrawString(GetDateTime(), 1000, 720)

	// id
	setFontSize(dc, 38)
	dc.DrawString(nickname, 50, 60)

	// 等级
	setFontSize(dc, 30)
	dc.DrawString(fmt.Sprintf("等级: %d", data.Rank), 50, 110)

	// 游戏时间
	setFontSize(dc, 30)
	dc.DrawString(fmt.Sprintf("游戏时间: %s", ConvertSecondsToHoursString(data.TimePlayed)),
		220, 110)

	// bfban
	setFontSize(dc, 30)
	dc.DrawString("BFBAN.COM", 765, 55)
	setFontSize(dc, 37)
	bfbanDesc := cons.BfbanStatusMap[bfban]

	if bfban == -2 {
		// 查询失败 白色
		dc.SetColor(color.White)
		dc.DrawString(bfbanDesc, 775, 115)
	} else if bfban == -1 {
		// 绿玩 绿色
		dc.SetHexColor("#19be6b")
		dc.DrawString(bfbanDesc, 790, 115)
	} else if bfban == 0 {
		// 需要注意 黄色
		dc.SetHexColor("#fff200")
		dc.DrawString(bfbanDesc, 790, 115)
	} else if bfban == 1 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		dc.DrawString(bfbanDesc, 810, 115)
	} else if bfban == 2 {
		// 需要注意 黄色
		dc.SetHexColor("#fff200")
		dc.DrawString(bfbanDesc, 795, 115)
	} else if bfban == 3 {
		// 自证 蓝色
		dc.SetHexColor("#1874bf")
		dc.DrawString(bfbanDesc, 770, 115)
	} else if bfban == 4 {
		// 需要注意 黄色
		dc.SetHexColor("#fff200")
		dc.DrawString(bfbanDesc, 780, 115)
	} else if bfban == 5 {
		// 需要注意 黄色
		dc.SetHexColor("#fff200")
		dc.DrawString(bfbanDesc, 795, 115)
	} else if bfban == 6 {
		// 更多管理投票
		setFontSize(dc, 30)
		dc.SetHexColor("#fff200")
		dc.DrawString("需要更多", 790, 92)
		dc.DrawString("管理投票", 790, 120)
	} else if bfban == 7 {
		// 绿玩 绿色
		dc.SetHexColor("#19be6b")
		dc.DrawString(bfbanDesc, 825, 115)
	} else if bfban == 8 {
		// 需要注意 黄色
		dc.SetHexColor("#fff200")
		dc.DrawString(bfbanDesc, 810, 115)
	}

	// 机器人社区状态
	dc.SetColor(color.White)
	setFontSize(dc, 30)
	dc.DrawString("BFV ROBOT", 985, 55)
	setFontSize(dc, 37)

	bfvRobotDesc := cons.RobotStatusMap[bfvrobot]

	if bfvrobot == -2 {
		// 查询失败 白色
		dc.SetColor(color.White)
		dc.DrawString(bfvRobotDesc, 990, 115)
	} else if bfvrobot == -1 {
		// 绿玩 绿色
		dc.SetHexColor("#19be6b")
		dc.DrawString(bfvRobotDesc, 1010, 115)
	} else if bfvrobot == 0 {
		// 绿玩 绿色
		dc.SetHexColor("#19be6b")
		dc.DrawString(bfvRobotDesc, 995, 115)
	} else if bfvrobot == 1 {
		// 需要注意 黄色
		setFontSize(dc, 26)
		dc.SetHexColor("#fff200")
		dc.DrawString("举报证据不足", 990, 92)
		dc.DrawString("[无效举报]", 1000, 120)
	} else if bfvrobot == 2 {
		// 挂壁 红色
		setFontSize(dc, 28)
		dc.SetHexColor("#ed4014")
		dc.DrawString(bfvRobotDesc, 985, 115)
	} else if bfvrobot == 3 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		dc.DrawString("[嫌疑或实锤]", 1002, 120)
	} else if bfvrobot == 4 {
		// 白名单 蓝色
		dc.SetHexColor("#1874bf")
		setFontSize(dc, 24)
		dc.DrawString("全局白名单", 1010, 92)
		dc.DrawString("[刷枪或其它自证]", 980, 120)
	} else if bfvrobot == 5 {
		// 白名单 蓝色
		dc.SetHexColor("#1874bf")
		setFontSize(dc, 24)
		dc.DrawString("全局白名单", 1010, 92)
		dc.DrawString("[Moss自证]", 1010, 120)
	} else if bfvrobot == 6 {
		// 绿色
		dc.SetHexColor("#19be6b")
		setFontSize(dc, 24)
		dc.DrawString("当前数据正常", 995, 92)
		setFontSize(dc, 15)
		dc.DrawString("(曾经有武器数据异常记录)", 980, 120)
	} else if bfvrobot == 7 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		dc.DrawString("[服主添加]", 1010, 120)
	} else if bfvrobot == 8 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("永久全局黑名单", 986, 92)
		dc.DrawString("[涉及底线的问题]", 980, 120)
	} else if bfvrobot == 9 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		dc.DrawString("[辱华涉政]", 1010, 120)
	} else if bfvrobot == 10 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		dc.DrawString("[检查组添加]", 1002, 120)
	} else if bfvrobot == 11 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		dc.DrawString("[不受欢迎的玩家]", 978, 120)
	} else if bfvrobot == 12 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		setFontSize(dc, 20)
		dc.DrawString("[机器人自动反外挂]", 982, 120)
	} else if bfvrobot == 13 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		setFontSize(dc, 17)
		dc.DrawString("[社区举报(武器数据异常)]", 975, 120)
	} else if bfvrobot == 14 {
		// 挂壁 红色
		dc.SetHexColor("#ed4014")
		setFontSize(dc, 24)
		dc.DrawString("全局黑名单", 1010, 92)
		setFontSize(dc, 20)
		dc.DrawString("[社区举报(BTR异常)", 982, 120)
	} else if bfvrobot == 15 {
		// 白名单 蓝色
		dc.SetHexColor("#1874bf")
		setFontSize(dc, 24)
		dc.DrawString("临时白名单", 1010, 92)
		dc.DrawString("[自证]", 1035, 120)
	}

	// 解包武器
	dc.SetColor(color.White)
	setFontSize(dc, 33)
	dc.DrawString("解包武器", 770, 185)

	// 解包数据
	unpackWeapon := data.UnpackWeapon
	// 按击杀排序取前6
	sort.Slice(unpackWeapon, func(i, j int) bool {
		return unpackWeapon[i].Kills > unpackWeapon[j].Kills
	})

	// 取前 3 个结果
	var topUnpackWeapon []dto.UnpackWeapon
	for _, v := range unpackWeapon {
		if v.Kills > 0 {
			topUnpackWeapon = append(topUnpackWeapon, v)
		}
		if len(topUnpackWeapon) >= 3 {
			break
		}
	}

	indexX := 770.0
	indexY := 220.0
	for _, weapon := range topUnpackWeapon {
		setFontSize(dc, 20)
		dc.DrawString(weapon.Name, indexX, indexY)
		dc.DrawString(fmt.Sprintf("击杀: %d", weapon.Kills), indexX+300, indexY)
		indexY += 30
	}

	dc.SetColor(color.White)
	// 胜率
	setFontSize(dc, 30)
	dc.DrawString(fmt.Sprintf("胜率: %s", data.WinPercent),
		530, 110)

	// kd
	setFontSize(dc, 28)
	dc.DrawString(fmt.Sprintf("KD: %s", data.KillDeath),
		50, 180)

	// 救援数
	dc.DrawString(fmt.Sprintf("救援数: %d", data.Revives),
		270, 185)

	// 击杀
	dc.DrawString(fmt.Sprintf("击杀: %d", data.Kills),
		520, 185)

	// KPM
	dc.DrawString(fmt.Sprintf("KPM: %s", data.KillsPerMinute),
		50, 235)

	// 最高连杀
	dc.DrawString(fmt.Sprintf("最高连杀: %d", data.HighestKillStreak),
		270, 235)

	// 死亡
	dc.DrawString(fmt.Sprintf("死亡: %d", data.Deaths),
		520, 235)

	// SPM
	dc.DrawString(fmt.Sprintf("SPM: %.2f", data.ScorePerMinute),
		50, 285)

	headshots := 0.0
	if data.Kills != 0 {
		headshots = float64(data.Headshots) / float64(data.Kills) * float64(100)
	}

	// 爆头率
	dc.DrawString(fmt.Sprintf("爆头率: %.2f%%", headshots),
		270, 285)

	sumRounds := data.Loses + data.Wins

	avgKills := 0
	if sumRounds != 0 {
		avgKills = data.Kills / sumRounds
	}

	// 场均击杀
	dc.DrawString(fmt.Sprintf("场均击杀: %d", avgKills),
		520, 285)

	// 武器数据
	weapons := data.Weapons
	// 按击杀排序取前6
	sort.Slice(weapons, func(i, j int) bool {
		return weapons[i].Kills > weapons[j].Kills
	})

	// 取前 6 个结果
	var topWeapons []dto.Weapons
	for _, w := range weapons {
		if w.Kills > 0 {
			topWeapons = append(topWeapons, w)
		}
		if len(topWeapons) >= 6 {
			break
		}
	}
	topWeaponsLen := len(topWeapons)
	indexX = 50.0
	indexY = 340.0
	for i, weapon := range topWeapons {
		setFontSize(dc, 20)
		dc.SetColor(color.White)
		// 武器名称
		dc.DrawString(weapon.Name, indexX, indexY)

		setFontSize(dc, 15)
		// KPM
		dc.DrawString(fmt.Sprintf("KPM: %s", weapon.KillsPerMinute),
			indexX+380, indexY-7)
		// 爆头率
		dc.DrawString(fmt.Sprintf("爆头率: %s", weapon.Headshots),
			indexX+510, indexY-7)
		// 击杀
		dc.DrawString(fmt.Sprintf("击杀: %d", weapon.Kills),
			indexX, indexY+23)
		// 效率
		dc.DrawString(fmt.Sprintf("效率: %s", weapon.HitVKills),
			indexX+380, indexY+23)
		// 命中率
		dc.DrawString(fmt.Sprintf("命中率: %s", weapon.Accuracy),
			indexX+510, indexY+23)

		if i+1 != topWeaponsLen {
			dc.SetColor(color.Gray{128})
			dc.DrawLine(indexX, indexY+33, indexX+650, indexY+33)
			dc.Stroke()
		}
		indexY += 65
	}

	// 载具数据
	vehicles := data.Vehicles
	// 按击杀/摧毁排序取前6
	sort.Slice(vehicles, func(i, j int) bool {
		if vehicles[i].Kills == vehicles[j].Kills {
			return vehicles[i].Destroy > vehicles[j].Destroy
		}
		return vehicles[i].Kills > vehicles[j].Kills
	})

	// 取前 6 个结果
	var topVehicles []dto.Vehicles
	for _, v := range vehicles {
		if v.Kills > 0 || v.Destroy > 0 {
			topVehicles = append(topVehicles, v)
		}
		if len(topVehicles) >= 6 {
			break
		}
	}
	topVehiclesLen := len(topVehicles)
	indexX = 770.0
	indexY = 340.0
	for i, vehicle := range topVehicles {
		setFontSize(dc, 20)
		dc.SetColor(color.White)
		// 载具名称
		dc.DrawString(vehicle.Name, indexX, indexY)

		setFontSize(dc, 15)
		// KPM
		dc.DrawString(fmt.Sprintf("KPM: %s", vehicle.KillsPerMinute),
			indexX+300, indexY-5)
		// 击杀
		dc.DrawString(fmt.Sprintf("击杀: %d", vehicle.Kills),
			indexX, indexY+23)
		// 摧毁
		dc.DrawString(fmt.Sprintf("摧毁: %d", vehicle.Destroy),
			indexX+300, indexY+23)

		if i+1 != topVehiclesLen {
			dc.SetColor(color.Gray{128})
			dc.DrawLine(indexX, indexY+33, indexX+400, indexY+33)
			dc.Stroke()
		}
		indexY += 65
	}

	imagePath := fmt.Sprintf(global.GConfig.Server.Output+"/%s/%d_%s_cx.png", GetDate(), data.PersonaID, GetUUID())
	// 保存结果图片
	outputFile, err := os.Create(imagePath)
	if err != nil {
		global.GLog.Error("os.Create", zap.Error(err))
		return "", errors.New("临时图片保存异常")
	}

	err = dc.EncodePNG(outputFile)
	if err != nil {
		global.GLog.Error("dc.EncodePNG", zap.Error(err))
		return "", errors.New("图片编码异常")
	}

	err = outputFile.Close()
	if err != nil {
		global.GLog.Error("outputFile.Close()", zap.Error(err))
	}

	// 压缩
	// 打开保存的 PNG 文件
	pngFile, err := os.Open(imagePath)
	if err != nil {
		global.GLog.Error("os.Open", zap.Error(err))
		return "", errors.New("图片压缩异常")
	}

	// 解码 PNG 文件
	img, _, err = image.Decode(pngFile)
	if err != nil {
		global.GLog.Error("image.Decode", zap.Error(err))
		return "", errors.New("图片二次解码异常")
	}

	// 压缩图片
	resizedImg := resize.Resize(1220, 728, img, resize.Lanczos3)

	// 创建要保存的 JPEG 文件
	imageJpgPath := fmt.Sprintf(global.GConfig.Server.Output+"/%s/%d_%s_cx.jpg", GetDate(), data.PersonaID, GetUUID())
	out, err := os.Create(imageJpgPath)
	if err != nil {
		global.GLog.Error("os.Create", zap.Error(err))
		return "", errors.New("查询结果图片异常")
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			global.GLog.Error("out.Close", zap.Error(err))
		}
	}(out)

	// 以 JPEG 格式保存压缩后的图片
	err = jpeg.Encode(out, resizedImg, nil)
	if err != nil {
		global.GLog.Error("jpeg.Encode", zap.Error(err))
		return "", errors.New("查询结果图片保存异常")
	}
	defer func() {
		err = pngFile.Close()
		if err != nil {
			global.GLog.Error("pngFile.Close()", zap.Error(err))
		}

		// 删除png文件
		err = os.Remove(imagePath)
		if err != nil {
			global.GLog.Error("os.RemoveBlack", zap.Error(err))
		}
	}()
	return imageJpgPath, nil
}

func QueryAndStore(nickname string, queryType int) (string, error) {

	err, checkData := CheckPlayer(nickname)
	if err != nil {
		return "", err
	}

	// 机器人社区状态
	bfvrobotData := dto.BfvRobotStatusData{ReasonStatus: -2, ReasonStatusName: "查询失败"}
	// bfban状态
	bfbanData := dto.BfBanStatusData{Status: -2}
	// 代表战排
	tag := ""
	// 生涯数据

	err, allData := GetPlayerData(checkData.PID)
	if err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	wg.Add(3)
	err = global.GPool.Submit(func() {
		defer func() {
			wg.Done()
		}()
		// bfban
		err, status := GetBfBanStatus(checkData.PID)
		if err != nil {
			return
		}
		bfbanData = status
	})
	if err != nil {
		return "", errors.New("协程池异常")
	}
	err = global.GPool.Submit(func() {
		defer func() {
			wg.Done()
		}()
		// bfvrobot
		err, status := GetBfvRobotStatus(checkData.PID)
		if err != nil {
			return
		}
		bfvrobotData = status
	})
	if err != nil {
		return "", errors.New("协程池异常")
	}

	err = global.GPool.Submit(func() {
		defer func() {
			wg.Done()
		}()
		err, tagResult := GetActiveTag(checkData.PersonaID)
		if err != nil {
			return
		}
		tag = tagResult
	})
	if err != nil {
		return "", errors.New("协程池异常")
	}

	wg.Wait()

	finalName := checkData.Name
	if tag != "" {
		finalName = "[" + tag + "] " + finalName
	}
	if queryType == 1 {
		return drawOverview(&allData, finalName, bfvrobotData.ReasonStatus, bfbanData.Status)
	} else {
		bfbanMap := make(map[string]interface{})
		bfbanMap["Status"] = bfbanData.Status
		bfbanMap["Name"] = cons.BfbanStatusMap[bfbanData.Status]
		return drawFullData(&allData, finalName, bfvrobotData, bfbanMap)
	}
}

func drawFullData(baseData *dto.BfvAllData, name string,
	robot dto.BfvRobotStatusData, bfban map[string]interface{}) (string, error) {

	if err := recover(); err != nil {
		global.GLog.Error("drawFullData recover", zap.Any("info", err))
	}

	err := CreateOutputDir()
	if err != nil {
		global.GLog.Error("os.Mkdir", zap.Error(err))
		return "", errors.New("服务器文件异常")
	}

	data := make(map[string]interface{})

	data["PlayerName"] = name
	data["PersonaID"] = baseData.PersonaID
	data["Rank"] = baseData.Rank
	data["Kills"] = baseData.Kills
	data["KillDeath"] = baseData.KillDeath
	data["KillsPerMinute"] = baseData.KillsPerMinute
	data["ScorePerMinute"] = baseData.ScorePerMinute

	var hs float64
	if baseData.Kills == 0 {
		hs = float64(0)
	} else {
		hs = float64(baseData.Headshots) / float64(baseData.Kills)
	}
	data["HS"] = fmt.Sprintf("%.2f%%", hs*100.0)
	data["Revives"] = baseData.Revives
	data["TimePlayed"] = ConvertSecondsToHoursString(baseData.TimePlayed)

	data["Bfban"] = bfban
	data["Bfvrobot"] = robot

	// 解包数据
	unpackWeapon := baseData.UnpackWeapon

	sort.Slice(unpackWeapon, func(i, j int) bool {
		return unpackWeapon[i].Kills > unpackWeapon[j].Kills
	})

	data["UnpackWeapon"] = unpackWeapon

	// 武器数据
	weapons := baseData.Weapons

	sort.Slice(weapons, func(i, j int) bool {
		return weapons[i].Kills > weapons[j].Kills
	})

	data["Weapons"] = weapons

	gadgets := baseData.Gadgets

	sort.Slice(gadgets, func(i, j int) bool {
		return gadgets[i].Kills > gadgets[j].Kills
	})

	data["Gadgets"] = gadgets

	vehicles := baseData.Vehicles

	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].Kills > vehicles[j].Kills
	})

	data["Vehicles"] = vehicles

	data["Time"] = GetDateTime()

	htmlContent, err := os.ReadFile(global.GConfig.Server.Template.Data)
	if err != nil {
		global.GLog.Error("os.ReadFile Template data", zap.Error(err))
		return "", errors.New("模板文件读取失败")
	}

	t, err := template.New("webpage").Parse(string(htmlContent))
	if err != nil {
		global.GLog.Error("template.Parse", zap.Error(err))
		return "", errors.New("模板文件解析失败")
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		global.GLog.Error("t.Execute", zap.Error(err))
		return "", errors.New("模板文件渲染失败")
	}

	imagePath := fmt.Sprintf(global.GConfig.Server.Output+"/%s/%d_%s_data.jpg", GetDate(),
		baseData.PersonaID, GetUUID())

	htmlPath := imagePath + ".html"
	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return "", errors.New("模板文件创建失败")
	}

	_, err = htmlFile.Write(buf.Bytes())
	if err != nil {
		return "", errors.New("模板文件写入失败")
	}

	page := global.GRodBrowser.MustPage("file://" + htmlPath)
	defer page.MustClose()
	page.MustSetViewport(800, 1, 1, false)
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

	return imagePath, nil
}
