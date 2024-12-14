package dto

type RemoveBanPlayerReq struct {
	PlayerId   string `json:"player_id"`
	Name       string `json:"name"`
	ServerName string `json:"serverName"`
}

type RemoveBanResp struct {
	Success int    `json:"success"`
	Error   int    `json:"error"`
	Code    string `json:"code"`
}
