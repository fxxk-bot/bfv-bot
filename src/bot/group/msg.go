package group

import (
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"bfv-bot/model/common/req"
	"strconv"
)

// SendGroupMsg 发送群聊消息
func SendGroupMsg(id int64, content string) {

	msg := req.Message{
		Data: req.Data{Text: content},
		Type: "text",
	}

	data := map[string]interface{}{
		"group_id": id,
		"message":  msg,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_group_msg_async", data)
}

// SendGroupMsgMultiple 发送群聊消息给多个群
func SendGroupMsgMultiple(id []int64, content string) {
	for _, item := range id {
		SendGroupMsg(item, content)
	}
}

// SendAtGroupMsg 发送群聊艾特消息
func SendAtGroupMsg(id int64, userId int64, content string) {
	var msgArr [2]req.Message
	msgArr[0] = req.Message{
		Data: req.Data{Qq: strconv.FormatInt(userId, 10)},
		Type: "at",
	}

	msgArr[1] = req.Message{
		Data: req.Data{Text: " " + content},
		Type: "text",
	}

	data := map[string]interface{}{
		"group_id": id,
		"message":  msgArr,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_group_msg_async", data)
}

func SendGroupReplyMsg(id int64, msgId int64, content string) {
	messages := make([]req.Message, 0)

	msg1 := req.Message{
		Data: req.Data{Id: strconv.FormatInt(msgId, 10)},
		Type: "reply",
	}
	messages = append(messages, msg1)

	msg2 := req.Message{
		Data: req.Data{Text: content},
		Type: "text",
	}
	messages = append(messages, msg2)

	data := map[string]interface{}{
		"group_id": id,
		"message":  messages,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_group_msg_async", data)
}

func SendGroupImageMsg(id int64, content string) {
	messages := make([]req.Message, 0)

	msg := req.Message{
		Data: req.Data{File: content},
		Type: "image",
	}
	messages = append(messages, msg)

	data := map[string]interface{}{
		"group_id": id,
		"message":  messages,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_group_msg_async", data)
}

func SendGroupImageReplyMsg(id int64, msgId int64, content string) {
	messages := make([]req.Message, 0)

	msg1 := req.Message{
		Data: req.Data{Id: strconv.FormatInt(msgId, 10)},
		Type: "reply",
	}
	messages = append(messages, msg1)

	imageMsg := req.Message{
		Data: req.Data{File: content},
		Type: "image",
	}
	messages = append(messages, imageMsg)

	data := map[string]interface{}{
		"group_id": id,
		"message":  messages,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_group_msg_async", data)
}

// DeleteMsg 撤回消息
func DeleteMsg(id int64) {

	data := map[string]interface{}{
		"message_id": id,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/delete_msg_async", data)
}
