package dto

type PlayerBaseInfoResp struct {
	Success int                `json:"success"`
	Code    string             `json:"code"`
	Data    PlayerBaseInfoData `json:"data"`
}
type PlayerBaseInfoRank struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}
type PlayerBaseInfoRankProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
type BasicStats struct {
	TimePlayed      int                        `json:"timePlayed"`
	Wins            int                        `json:"wins"`
	Losses          int                        `json:"losses"`
	Kills           int                        `json:"kills"`
	Deaths          int                        `json:"deaths"`
	Kpm             float64                    `json:"kpm"`
	Spm             float64                    `json:"spm"`
	Skill           int                        `json:"skill"`
	SoldierImageURL string                     `json:"soldierImageUrl"`
	Rank            PlayerBaseInfoRank         `json:"rank"`
	RankProgress    PlayerBaseInfoRankProgress `json:"rankProgress"`
}
type Progression struct {
	ValueNeeded   int  `json:"valueNeeded"`
	ValueAcquired int  `json:"valueAcquired"`
	Unlocked      bool `json:"unlocked"`
}
type TidesOfWarInfo struct {
	CurrentRank int         `json:"currentRank"`
	Progression Progression `json:"progression"`
}
type PlayerBaseInfoData struct {
	BasicStats        BasicStats     `json:"basicStats"`
	AwardScore        int            `json:"awardScore"`
	BonusScore        int            `json:"bonusScore"`
	SquadScore        int            `json:"squadScore"`
	AvengerKills      int            `json:"avengerKills"`
	SaviorKills       int            `json:"saviorKills"`
	HighestKillStreak int            `json:"highestKillStreak"`
	DogtagsTaken      int            `json:"dogtagsTaken"`
	RoundsPlayed      int            `json:"roundsPlayed"`
	FlagsCaptured     int            `json:"flagsCaptured"`
	FlagsDefended     int            `json:"flagsDefended"`
	AccuracyRatio     float64        `json:"accuracyRatio"`
	HeadShots         int            `json:"headShots"`
	LongestHeadShot   int            `json:"longestHeadShot"`
	Revives           int            `json:"revives"`
	Heals             int            `json:"heals"`
	Repairs           int            `json:"repairs"`
	SuppressionAssist int            `json:"suppressionAssist"`
	Kdr               float64        `json:"kdr"`
	KillAssists       int            `json:"killAssists"`
	TidesOfWarInfo    TidesOfWarInfo `json:"tidesOfWarInfo"`
	Draws             int            `json:"draws"`
	DetailedStatType  string         `json:"detailedStatType"`
}
