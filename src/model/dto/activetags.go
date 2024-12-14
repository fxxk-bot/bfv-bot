package dto

type ActiveTagsResp struct {
	Success int               `json:"success"`
	Code    string            `json:"code"`
	Data    map[string]string `json:"data"`
}
