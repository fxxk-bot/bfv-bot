package dto

import "time"

type BfBanStatusResp struct {
	Code string          `json:"code"`
	Data BfBanStatusData `json:"data"`
}
type BfBanStatusData struct {
	ID              int         `json:"id"`
	OriginName      string      `json:"originName"`
	OriginUserID    string      `json:"originUserId"`
	OriginPersonaID string      `json:"originPersonaId"`
	Games           []string    `json:"games"`
	CheatMethods    []string    `json:"cheatMethods"`
	AvatarLink      string      `json:"avatarLink"`
	ViewNum         int         `json:"viewNum"`
	CommentsNum     int         `json:"commentsNum"`
	Status          int         `json:"status"`
	CreateTime      time.Time   `json:"createTime"`
	UpdateTime      time.Time   `json:"updateTime"`
	AppealStatus    interface{} `json:"appealStatus"`
}
