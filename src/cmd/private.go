package cmd

import (
	"bfv-bot/bot/private"
	"bfv-bot/common/global"
	"bfv-bot/flow"
	"bfv-bot/model/common/req"
	"bfv-bot/model/common/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func init() {
	privateCommandMap["addblack"] = addblack
	privateCommandMap["removeblack"] = removeblack
	privateCommandMap["removecardcheck"] = removecardcheck
	privateCommandMap["addsensitive"] = addsensitive
	privateCommandMap["removesensitive"] = removesensitive
	privateCommandMap["addignore"] = addignore
	privateCommandMap["removeignore"] = removeignore
	privateCommandMap["bindtoken"] = bindtoken
	privateCommandMap["bindgameid"] = bindgameid
	privateCommandMap["op"] = op

	privateOpCommandMap["start"] = opStart
	privateOpCommandMap["stop"] = opStop

	privateOpCommandMap["start-broadcast"] = opStartBroadcast
	privateOpCommandMap["stop-broadcast"] = opStopBroadcast

	privateOpCommandMap["checknow"] = opChecknow
	privateOpCommandMap["gameid"] = opGameid
	privateOpCommandMap["token"] = opToken
	privateOpCommandMap["ignorelist"] = opIgnorelist
	privateOpCommandMap["deleteignorelist"] = opDeleteignorelist
	privateOpCommandMap["blacklist"] = opBlacklist

	privateQuickCommandMap["help"] = getPrivateHelpInfo
	privateQuickCommandMap[".help"] = getPrivateHelpInfo

}

func opStart(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.Bfv.Active = true
	resp.ReplyOk(c, "开始检测")
}

func opStop(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.Bfv.Active = false
	global.GConfig.Bfv.ClearGameId()
	resp.ReplyOk(c, "结束检测")
}

func opStartBroadcast(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.QQBot.BotToBot.Enable = true
	resp.ReplyOk(c, "开始喊话")
}

func opStopBroadcast(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.QQBot.BotToBot.Enable = false
	resp.ReplyOk(c, "结束喊话")
}

func opChecknow(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	cronService.CheckBlackListAndNotify()
	resp.ReplyOk(c, "立即检测")
}

func opGameid(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	var builder strings.Builder
	for _, info := range global.GConfig.Bfv.Server {
		builder.WriteString(info.ServerName)
		builder.WriteString("\n")
		if info.GetGameId() == "" {
			builder.WriteString("无")
		} else {
			builder.WriteString(info.GetGameId())
		}
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opToken(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	var builder strings.Builder
	for _, info := range global.GConfig.Bfv.Server {
		builder.WriteString(info.ServerName)
		builder.WriteString("\n")
		if info.GetToken() == "" {
			builder.WriteString("无")
		} else {
			builder.WriteString(info.GetToken())
		}
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opIgnorelist(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	list := dbService.QueryAllIgnoreList()
	var builder strings.Builder
	builder.WriteString("忽略名单\n")
	for key := range list {
		builder.WriteString(key)
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opDeleteignorelist(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	err := dbService.DeleteAllIgnore()
	if err != nil {
		resp.ReplyOk(c, "清空忽略名单失败")
	} else {
		resp.ReplyOk(c, "清空忽略名单成功")
	}
}

func opBlacklist(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	list := dbService.QueryAllBlackList()
	var builder strings.Builder
	builder.WriteString("黑名单\n")
	for key, value := range list {
		builder.WriteString("pid: ")
		builder.WriteString(key)
		builder.WriteString("\t")
		builder.WriteString("id: ")
		builder.WriteString(value.Name)
		builder.WriteString("\t")
		builder.WriteString("原因: ")
		builder.WriteString(value.Reason)
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func addblack(msg *req.MsgData, c *gin.Context, _ string, value string) {

	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.AddBlack, value)
	private.SendPrivateMsg(msg.UserID, "[添加黑名单] 请输入原因")

	resp.EmptyOk(c)
}

func removeblack(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err := dbService.RemoveBlack(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("黑名单用户 [%s] 移除成功", value))
	}
}

func removecardcheck(_ *req.MsgData, c *gin.Context, _ string, value string) {

	qq, _ := strconv.ParseInt(value, 10, 64)
	err := dbService.DeleteCardCheck(qq)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("ID检测 [%s] 移除成功", value))
	}
}

func addsensitive(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err := dbService.AddSensitive(value)
	if err != nil {
		resp.ReplyOk(c, "添加失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("添加成功"))
		global.GSensitive.AddWord(value)
	}
}

func removesensitive(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err := dbService.RemoveSensitive(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, "移除成功, 重启生效")
	}

}

func addignore(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err := dbService.AddIgnore(value)
	if err != nil {
		resp.ReplyOk(c, "添加失败")
	} else {
		resp.ReplyOk(c, "忽略名单添加成功")
	}
}

func removeignore(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err := dbService.RemoveIgnore(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("忽略名单用户 [%s] 移除成功", value))
	}
}

func bindtoken(msg *req.MsgData, c *gin.Context, _ string, value string) {
	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.BindToken, value)
	private.SendPrivateMsg(msg.UserID, "[绑定Token] 请输入服务器ID")

	resp.EmptyOk(c)
}

func bindgameid(msg *req.MsgData, c *gin.Context, _ string, value string) {
	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.BindGameID, value)
	private.SendPrivateMsg(msg.UserID, "[绑定GameID] 请输入服务器ID")

	resp.EmptyOk(c)
}

func op(msg *req.MsgData, c *gin.Context, key string, value string) {
	function, ok := privateOpCommandMap[value]
	if ok {
		function(msg, c, key, value)
		return
	}
	resp.EmptyOk(c)
}

func getPrivateHelpInfo(_ *req.MsgData, c *gin.Context, _ string) {
	var builder strings.Builder
	builder.WriteString("绑定token: bindtoken=<token>\n")
	builder.WriteString("绑定gameid: bindgameid=<gameid>\n")
	builder.WriteString("添加黑名单: addblack=<id>\n")
	builder.WriteString("移除黑名单: removeblack=<id>\n")
	builder.WriteString("移除id检测: removecardcheck=<qq>\n")
	builder.WriteString("添加敏感词: addsensitive=<id>\n")
	builder.WriteString("移除敏感词: removesensitive=<id>\n")
	builder.WriteString("添加忽略名单: addignore=<id>\n")
	builder.WriteString("移除忽略名单: removeignore=<id>\n")
	builder.WriteString("获取游戏id: op=gameid\n")
	builder.WriteString("获取服务器token: op=token\n")
	builder.WriteString("开始检测黑名单: op=start\n")
	builder.WriteString("停止检测黑名单: op=stop\n")
	builder.WriteString("开始喊话: op=start-broadcast\n")
	builder.WriteString("停止喊话: op=stop-broadcast\n")
	builder.WriteString("立即检测黑名单: op=checknow\n")
	builder.WriteString("清空忽略名单: op=deleteignorelist\n")
	builder.WriteString("忽略列表: op=ignorelist\n")
	builder.WriteString("黑名单列表: op=blacklist")
	resp.ReplyOk(c, builder.String())
}

func GetPrivateCommandFunc(key string) (func(*req.MsgData, *gin.Context, string, string), bool) {
	f, ok := privateCommandMap[key]
	return f, ok
}

func GetPrivateQuickCommandFunc(key string) (func(*req.MsgData, *gin.Context, string), bool) {
	f, ok := privateQuickCommandMap[key]
	return f, ok
}
