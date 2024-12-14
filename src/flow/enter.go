package flow

import (
	"bfv-bot/model/common/req"
	"bfv-bot/model/flow"
	"bfv-bot/service"
)

// PrivateFlowable 私聊业务流
var PrivateFlowable = make(map[int64]flow.PrivateFlow)

// GroupFlowable 群聊业务流
var GroupFlowable = make(map[string]flow.GroupFlow)

// GroupStepMap 群聊业务处理映射
var GroupStepMap = make(map[string]func(msg *req.MsgData, flowData *flow.GroupFlow))

// PrivateStepMap 私聊业务处理映射
var PrivateStepMap = make(map[string]func(msg *req.MsgData, flowData *flow.PrivateFlow))

var (
	// dbService private
	dbService = service.ServiceGroup.DbService
)
