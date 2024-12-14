package dto

type CaptchaResp struct {
	Success int         `json:"success"`
	Code    string      `json:"code"`
	Data    CaptchaData `json:"data"`
}

type CaptchaData struct {
	Hash    string `json:"hash"`
	Content string `json:"content"`
	Type    string `json:"type"`
}
