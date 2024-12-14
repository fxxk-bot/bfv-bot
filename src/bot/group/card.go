package group

import (
	"bfv-bot/common/global"
	"bfv-bot/common/http"
)

func SetCard(groupId int64, userId int64, card string) {
	data := map[string]interface{}{
		"group_id": groupId,
		"user_id":  userId,
		"card":     card,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/set_group_card_async", data)
}
