package dto

import "time"

type BanLogResp struct {
	Success int          `json:"success"`
	Code    string       `json:"code"`
	Data    []BanLogData `json:"data"`
}
type BanLogData struct {
	PersonaID  string    `json:"personaId"`
	ServerName string    `json:"serverName"`
	Reason     string    `json:"reason"`
	BanType    int       `json:"banType"`
	CreateTime time.Time `json:"createTime"`
}
