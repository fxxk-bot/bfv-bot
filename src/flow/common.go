package flow

import (
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"fmt"
	"time"
)

const Ban = "ban"
const RemoveBan = "remove-ban"
const AddBlack = "addblack"
const AddJoinBlack = "addjoinblack"
const BindToken = "bind-token"
const BindGameID = "bind-gameid"

func init() {
	GroupStepMap[Ban] = BanStep
	GroupStepMap[RemoveBan] = RemoveBanStep

	PrivateStepMap[AddBlack] = AddBlackStep
	PrivateStepMap[AddJoinBlack] = AddJoinBlackStep
	PrivateStepMap[BindToken] = BindTokenStep
	PrivateStepMap[BindGameID] = BindGameIDStep
}

func BuildGroupKey(qq int64, group int64) string {
	return fmt.Sprintf("%d_%d", qq, group)
}

func InitPrivateFlow(qq int64, msgId int64, key string, cmd string) {
	content := make([]string, 0)
	content = append(content, cmd)
	flowable := flow.PrivateFlow{
		Key:        key,
		Content:    content,
		Step:       2,
		Qq:         qq,
		MsgId:      msgId,
		ActiveTime: time.Now().UnixMilli(),
	}
	PrivateFlowable[qq] = flowable
}

func InitGroupFlow(qq int64, group int64, msgId int64, key string, cmd string) {
	content := make([]string, 0)
	content = append(content, cmd)
	flowable := flow.GroupFlow{
		Key:        key,
		Content:    content,
		Step:       1,
		GroupId:    group,
		MsgId:      msgId,
		ActiveTime: time.Now().UnixMilli(),
	}
	GroupFlowable[BuildGroupKey(qq, group)] = flowable
}

func DoPrivateNextStep(msg *req.MsgData) bool {
	value, ok := PrivateFlowable[msg.UserID]
	if ok {
		PrivateStepMap[value.Key](msg, &value)

		value.ActiveTime = time.Now().UnixMilli()
		_, ok := PrivateFlowable[msg.UserID]
		if ok {
			PrivateFlowable[msg.UserID] = value
		}
	}
	return ok
}

func DoGroupNextStep(msg *req.MsgData) bool {
	key := BuildGroupKey(msg.UserID, msg.GroupID)
	value, ok := GroupFlowable[key]
	if ok {
		GroupStepMap[value.Key](msg, &value)

		value.ActiveTime = time.Now().UnixMilli()

		_, ok := GroupFlowable[key]
		if ok {
			GroupFlowable[key] = value
		}
	}
	return ok
}

func IncGroupStep(msg *req.MsgData, content string) {
	key := BuildGroupKey(msg.UserID, msg.GroupID)
	value, ok := GroupFlowable[key]
	if ok {
		value.Step = value.Step + 1
		value.Content = append(value.Content, content)
		GroupFlowable[key] = value
	}
}

func DeleteGroupStep(msg *req.MsgData) {
	delete(GroupFlowable, BuildGroupKey(msg.UserID, msg.GroupID))
}

func DeleteGroupStepByKey(key string) {
	delete(GroupFlowable, key)
}
