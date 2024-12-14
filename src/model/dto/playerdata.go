package dto

type BfvAllResp struct {
	Success int        `json:"success"`
	Code    string     `json:"code"`
	Data    BfvAllData `json:"data"`
}

type Weapons struct {
	Name            string `json:"name"`
	Damage          int    `json:"damage,omitempty"`
	ShotsHit        int    `json:"shotsHit"`
	HeadshotKills   int    `json:"headshotKills,omitempty"`
	TimeEquipped    int    `json:"timeEquipped,omitempty"`
	ShotsFired      int    `json:"shotsFired,omitempty"`
	Kills           int    `json:"kills"`
	Score           int    `json:"score,omitempty"`
	KillsPerMinute  string `json:"killsPerMinute,omitempty"`
	Accuracy        string `json:"accuracy"`
	Headshots       string `json:"headshots"`
	HitVKills       string `json:"hitVKills"`
	DamagePerHit    string `json:"damagePerHit,omitempty"`
	DamagePerMinute string `json:"damagePerMinute,omitempty"`
}
type Vehicles struct {
	Name            string `json:"name"`
	Damage          int    `json:"damage,omitempty"`
	TimeEquipped    int    `json:"timeEquipped"`
	Kills           int    `json:"kills"`
	KillsPerMinute  string `json:"killsPerMinute"`
	Accuracy        string `json:"accuracy"`
	Headshots       string `json:"headshots"`
	HitVKills       string `json:"hitVKills"`
	DamagePerHit    string `json:"damagePerHit,omitempty"`
	DamagePerMinute string `json:"damagePerMinute,omitempty"`
	Destroy         int    `json:"destroy,omitempty"`
	Destroyed       int    `json:"destroyed,omitempty"`
}
type Gadgets struct {
	Name           string `json:"name"`
	TimeEquipped   int    `json:"timeEquipped"`
	Kills          int    `json:"kills"`
	KillsPerMinute string `json:"killsPerMinute"`
	Accuracy       string `json:"accuracy"`
	Headshots      string `json:"headshots"`
	HitVKills      string `json:"hitVKills"`
}

type UnpackWeapon struct {
	Name           string `json:"name"`
	ShotsHit       int    `json:"shotsHit,omitempty"`
	TimeEquipped   int    `json:"timeEquipped"`
	ShotsFired     int    `json:"shotsFired,omitempty"`
	Kills          int    `json:"kills"`
	Score          int    `json:"score,omitempty"`
	KillsPerMinute string `json:"killsPerMinute"`
	Accuracy       string `json:"accuracy"`
	Headshots      string `json:"headshots"`
	HitVKills      string `json:"hitVKills"`
}
type BfvAllData struct {
	Weapons           []Weapons      `json:"weapons"`
	Vehicles          []Vehicles     `json:"vehicles"`
	Gadgets           []Gadgets      `json:"gadgets"`
	UnpackWeapon      []UnpackWeapon `json:"unpackWeapon"`
	PersonaID         int64          `json:"personaId"`
	Kills             int            `json:"kills"`
	Deaths            int            `json:"deaths"`
	KillAssists       int            `json:"killAssists"`
	Heals             int            `json:"heals"`
	Revives           int            `json:"revives"`
	Headshots         int            `json:"headshots"`
	Rank              int            `json:"rank"`
	Wins              int            `json:"wins"`
	RoundsPlayed      int            `json:"roundsPlayed"`
	Loses             int            `json:"loses"`
	HighestKillStreak int            `json:"highestKillStreak"`
	AwardScore        int            `json:"awardScore"`
	BonusScore        int            `json:"bonusScore"`
	SquadScore        int            `json:"squadScore"`
	TotalScore        int            `json:"totalScore"`
	KillDeath         string         `json:"killDeath"`
	TimePlayed        int            `json:"timePlayed"`
	KillsPerMinute    string         `json:"killsPerMinute"`
	WinPercent        string         `json:"winPercent"`
	ScorePerMinute    float64        `json:"scorePerMinute"`
}
