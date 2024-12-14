package dto

type GetGroupMemberInfoResp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Wording string                 `json:"wording"`
	Retcode int                    `json:"retcode"`
	Data    GetGroupMemberInfoData `json:"data,omitempty"`
}

type GetGroupMemberInfoData struct {
	GroupId int64  `json:"group_id"`
	UserId  int64  `json:"user_id"`
	Card    string `json:"card"`
}

type GetGroupMemberListResp struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Wording string                   `json:"wording"`
	Retcode int                      `json:"retcode"`
	Data    []GetGroupMemberListData `json:"data,omitempty"`
}

type GetGroupMemberListData struct {
	GroupId int64  `json:"group_id"`
	UserId  int64  `json:"user_id"`
	Card    string `json:"card"`
}
