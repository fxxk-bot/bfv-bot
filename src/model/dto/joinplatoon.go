package dto

type JoinPlatoonResp struct {
	Success int               `json:"success"`
	Code    string            `json:"code"`
	Data    []JoinPlatoonData `json:"data"`
}
type JoinConfig struct {
	CanApplyMembership bool `json:"canApplyMembership"`
	IsFreeJoin         bool `json:"isFreeJoin"`
}
type JoinPlatoonData struct {
	GUID        string     `json:"guid"`
	Name        string     `json:"name"`
	Size        int        `json:"size"`
	JoinConfig  JoinConfig `json:"joinConfig"`
	Description string     `json:"description"`
	Tag         string     `json:"tag"`
	Emblem      string     `json:"emblem"`
	Verified    bool       `json:"verified"`
	CreatorID   string     `json:"creatorId"`
	DateCreated int        `json:"dateCreated"`
}
