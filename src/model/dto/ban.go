package dto

type BanPlayerReq struct {
	PlayerId int64  `json:"player_id"`
	Name     string `json:"name"`
}

type BanResp struct {
	Success int    `json:"success"`
	Error   int    `json:"error"`
	Code    string `json:"code"`
}
