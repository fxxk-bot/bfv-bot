package dto

// ApiPlayersResp 以下是gametools的玩家列表
type ApiPlayersResp struct {
	Serverinfo      ServerinfoData `json:"serverinfo"`
	Teams           []TeamsData    `json:"teams"`
	Que             []QueData      `json:"que"`
	Loading         []interface{}  `json:"loading"`
	UpdateTimestamp int            `json:"update_timestamp"`
}
type ServerinfoData struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Region      string        `json:"region"`
	Country     string        `json:"country"`
	Level       string        `json:"level"`
	Mode        string        `json:"mode"`
	Maps        []string      `json:"maps"`
	Owner       string        `json:"owner"`
	Settings    []interface{} `json:"settings"`
	Servertype  string        `json:"servertype"`
}
type PlayersData struct {
	Rank     int    `json:"rank"`
	Latency  int    `json:"latency"`
	Slot     int    `json:"slot"`
	JoinTime int64  `json:"join_time"`
	UserID   int64  `json:"user_id"`
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Platoon  string `json:"platoon"`
}
type TeamsData struct {
	Teamid    string        `json:"teamid"`
	Players   []PlayersData `json:"players"`
	Key       string        `json:"key"`
	Name      string        `json:"name"`
	ShortName string        `json:"shortName"`
	Image     string        `json:"image"`
	Faction   string        `json:"faction"`
}
type QueData struct {
	Rank     int    `json:"rank"`
	Latency  int    `json:"latency"`
	Slot     int    `json:"slot"`
	JoinTime int64  `json:"join_time"`
	UserID   int64  `json:"user_id"`
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Platoon  string `json:"platoon"`
}

// RobotPlayersResp 以下是bfvrobot的struct
type RobotPlayersResp struct {
	Success int                          `json:"success"`
	Message string                       `json:"message"`
	Data    RobotPlayersRobotPlayersData `json:"data"`
}
type RobotPlayersSoldier struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type RobotPlayersSpectator struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type RobotPlayersQueue struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}
type RobotPlayersSlots struct {
	Soldier   RobotPlayersSoldier   `json:"Soldier"`
	Spectator RobotPlayersSpectator `json:"Spectator"`
	Queue     RobotPlayersQueue     `json:"Queue"`
}
type RobotPlayersTeam1 struct {
	PersonaID int64  `json:"personaId"`
	UserID    int64  `json:"userId"`
	Join      int64  `json:"join"`
	Locale    int    `json:"locale"`
	Name      string `json:"name"`
	Platoon   string `json:"platoon"`
}
type RobotPlayersTeam2 struct {
	PersonaID int64  `json:"personaId"`
	UserID    int64  `json:"userId"`
	Join      int64  `json:"join"`
	Locale    int    `json:"locale"`
	Name      string `json:"name"`
	Platoon   string `json:"platoon"`
}
type RobotPlayersLoading struct {
	PersonaID int64  `json:"personaId"`
	UserID    int64  `json:"userId"`
	Join      int64  `json:"join"`
	Locale    int    `json:"locale"`
	Name      string `json:"name"`
	Platoon   string `json:"platoon"`
}
type RobotPlayersPlayers struct {
	Team1      []RobotPlayersTeam1   `json:"team_1"`
	Team2      []RobotPlayersTeam2   `json:"team_2"`
	TeamX      []interface{}         `json:"team_x"`
	Loading    []RobotPlayersLoading `json:"loading"`
	Spectators []interface{}         `json:"spectators"`
}
type RobotPlayersRobotPlayersData struct {
	GameID       int64               `json:"gameId"`
	ServerName   string              `json:"serverName"`
	Description  string              `json:"description"`
	Region       string              `json:"region"`
	Country      string              `json:"country"`
	Slots        RobotPlayersSlots   `json:"slots"`
	Players      RobotPlayersPlayers `json:"players"`
	Admins       []interface{}       `json:"admins"`
	PlaygroundID string              `json:"playgroundId"`
}

// CommonPlayersResp 以下是通用服务器玩家列表struct
type CommonPlayersResp struct {
	TeamOne []CommonPlayersData
	TeamTwo []CommonPlayersData
	Que     []CommonPlayersData
}

type CommonPlayersData struct {
	UserID    int64
	PersonaID int64
	Name      string
	Platoon   string
	Join      int64
}
