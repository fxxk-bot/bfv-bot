package dto

type BfBanBatchResp struct {
	Success int              `json:"success"`
	Code    string           `json:"code"`
	Data    []BfBanBatchData `json:"data"`
}
type BfBanBatchData struct {
	PersonaID int64 `json:"personaId"`
	Status    int   `json:"status"`
}
