package flow

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"time"
)

func CleanExpiredPrivateFlow() {

	deleteList := make([]int64, 0)

	now := time.Now().UnixMilli()
	for key, value := range PrivateFlowable {
		if now-(60*1000) > value.ActiveTime {
			private.SendPrivateReplyMsg(key, value.MsgId, "当前对话已超时, 请重新发起")
			deleteList = append(deleteList, key)
		}
	}

	for _, item := range deleteList {
		delete(PrivateFlowable, item)
	}

}

func CleanExpiredGroupFlow() {

	deleteList := make([]string, 0)

	now := time.Now().UnixMilli()
	for key, value := range GroupFlowable {
		if now-(60*1000) > value.ActiveTime {
			group.SendGroupReplyMsg(value.GroupId, value.MsgId, "当前对话已超时, 请重新发起")
			deleteList = append(deleteList, key)
		}
	}

	for _, item := range deleteList {
		DeleteGroupStepByKey(item)
	}

}
