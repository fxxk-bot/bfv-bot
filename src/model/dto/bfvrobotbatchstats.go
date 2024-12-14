package dto

type BfvRobotBatchStatsResp struct {
	Success int                      `json:"success"`
	Code    string                   `json:"code"`
	Data    []BfvRobotBatchStatsData `json:"data"`
}
type BfvRobotBatchStatsData struct {
	PersonaID         int64   `json:"personaId"`
	Kills             int     `json:"kills"`
	Deaths            int     `json:"deaths"`
	KillAssists       int     `json:"killAssists"`
	Heals             int     `json:"heals"`
	Revives           int     `json:"revives"`
	Headshots         int     `json:"headshots"`
	Rank              int     `json:"rank"`
	Wins              int     `json:"wins"`
	RoundsPlayed      int     `json:"roundsPlayed"`
	Loses             int     `json:"loses"`
	HighestKillStreak int     `json:"highestKillStreak"`
	AwardScore        int     `json:"awardScore"`
	BonusScore        int     `json:"bonusScore"`
	SquadScore        int     `json:"squadScore"`
	TotalScore        int     `json:"totalScore"`
	KillDeath         string  `json:"killDeath"`
	TimePlayed        int     `json:"timePlayed"`
	KillsPerMinute    string  `json:"killsPerMinute"`
	WinPercent        string  `json:"winPercent"`
	ScorePerMinute    float64 `json:"scorePerMinute"`
}
