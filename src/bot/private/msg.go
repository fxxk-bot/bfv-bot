package private

import (
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"bfv-bot/model/common/req"
	"strconv"
)

// SendPrivateMsg 发送私聊消息
func SendPrivateMsg(id int64, content string) {
	msg := req.Message{
		Data: req.Data{Text: content},
		Type: "text",
	}

	data := map[string]interface{}{
		"user_id": id,
		"message": msg,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_private_msg_async", data)
}

// SendPrivateMsgMultiple 发送私聊消息给多人
func SendPrivateMsgMultiple(id []int64, content string) {
	for _, item := range id {
		SendPrivateMsg(item, content)
	}
}

func SendPrivateReplyMsg(id int64, msgId int64, content string) {

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
		"user_id": id,
		"message": messages,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/send_private_msg_async", data)
}
