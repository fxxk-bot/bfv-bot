package dto

type CheckPlayerResp struct {
	Success int             `json:"success"`
	Code    string          `json:"code"`
	Data    CheckPlayerData `json:"data"`
}
type CheckPlayerData struct {
	PersonaID int64 `json:"personaId"`
	PID       string
	Name      string `json:"name"`
	UserID    int64  `json:"userId"`
}
