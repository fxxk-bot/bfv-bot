package dto

type GtBatchStatusResp struct {
	Data []GtBatchStatusData `json:"data"`
}

type GtBatchStatusData struct {
	ID             string `json:"id"`
	Name           string
	Rank           float64 `json:"rank"`
	KillsPerMinute float64 `json:"killsPerMinute"`
}
