package dto

type ServerGameToolsResp struct {
	Servers []ServerGameToolsData `json:"servers"`
}

type ServerGameToolsData struct {
	Prefix       string      `json:"prefix"`
	Description  string      `json:"description"`
	PlayerAmount int         `json:"playerAmount"`
	MaxPlayers   int         `json:"maxPlayers"`
	InSpectator  int         `json:"inSpectator"`
	InQue        int         `json:"inQue"`
	ServerInfo   string      `json:"serverInfo"`
	URL          string      `json:"url"`
	Mode         string      `json:"mode"`
	CurrentMap   string      `json:"currentMap"`
	OwnerID      string      `json:"ownerId"`
	Country      string      `json:"country"`
	Region       string      `json:"region"`
	Platform     string      `json:"platform"`
	ServerID     string      `json:"serverId"`
	IsCustom     bool        `json:"isCustom"`
	SmallMode    string      `json:"smallMode"`
	Teams        interface{} `json:"teams"`
	Official     bool        `json:"official"`
	GameID       string      `json:"gameId"`
}
