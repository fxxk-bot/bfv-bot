package dto

type BfvRobotStatusResp struct {
	Success int                `json:"success"`
	Data    BfvRobotStatusData `json:"data"`
}

type BfvRobotStatusData struct {
	PersonaID           int64  `json:"personaId"`
	OperationStatus     int    `json:"operationStatus"`
	OperationStatusName string `json:"operationStatusName"`
	ReasonStatus        int    `json:"reasonStatus"`
	ReasonStatusName    string `json:"reasonStatusName"`
}
