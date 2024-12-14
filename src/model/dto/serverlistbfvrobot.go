package dto

type ServerListBfvRobotResp struct {
	Success int                      `json:"success"`
	Code    string                   `json:"code"`
	Data    []ServerListBfvRobotData `json:"data"`
}
type Soldier struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type Spectator struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type Queue struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type Slots struct {
	Soldier   Soldier   `json:"Soldier"`
	Spectator Spectator `json:"Spectator"`
	Queue     Queue     `json:"Queue"`
}
type ServerListBfvRobotData struct {
	MapName     string `json:"mapName"`
	MapMode     string `json:"mapMode"`
	SmallMode   string `json:"smallMode"`
	Description string `json:"description"`
	GameID      int64  `json:"gameId"`
	Official    bool   `json:"official"`
	OwnerID     int64  `json:"ownerId"`
	ServerName  string `json:"serverName"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	Slots       Slots  `json:"slots"`
	URL         string `json:"url"`
}
