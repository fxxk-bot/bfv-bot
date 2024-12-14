package flow

type PrivateFlow struct {
	// 业务key
	Key string
	// 完整流程指令
	Content []string
	// 业务当前步骤数
	Step int
	// 触发步骤的qq
	Qq int64
	// 触发步骤的消息id
	MsgId int64
	// 启用步骤时间 超过1分钟终止 毫秒时间戳
	ActiveTime int64
}
