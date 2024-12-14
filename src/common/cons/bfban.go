package cons

// BfbanStatusMap 状态枚举
var BfbanStatusMap = map[int]string{
	-2: "查询失败",
	-1: "无记录",
	0:  "待处理",
	1:  "石锤",
	2:  "待自证",
	3:  "Moss自证",
	4:  "无效举报",
	5:  "讨论中",
	6:  "需要更多管理投票",
	7:  "无",
	8:  "刷枪",
}
