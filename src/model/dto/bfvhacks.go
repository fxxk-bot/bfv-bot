package dto

type BfvHacksResp struct {
	AccountSpecificMessage  interface{}   `json:"account_specific_message"`
	Bfban                   []interface{} `json:"bfban"`
	BfbanStatus             interface{}   `json:"bfban_status"`
	HackLevel               string        `json:"hack_level"`
	HackScore               int           `json:"hack_score"`
	HackScoreCurrent        int           `json:"hack_score_current"`
	LackingStats            bool          `json:"lacking_stats"`
	OriginPersonalID        interface{}   `json:"origin_personal_id"`
	PlayerHandle            string        `json:"player_handle"`
	PlayerID                int64         `json:"player_id"`
	PlayerIsNew             bool          `json:"player_is_new"`
	PotentialHardcorePlayer bool          `json:"potential_hardcore_player"`
	Stats                   Stats         `json:"stats"`
	StatsAll                interface{}   `json:"stats_all"`
	Vehicles                Vehicles      `json:"vehicles"`
	Videos                  []interface{} `json:"videos"`
	Weapons                 []Weapons     `json:"weapons"`
}

type Stats struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Resulte []Resulte `json:"resulte"`
}
type Rank struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}
type RankProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
type Resulte struct {
	TimePlayed       int           `json:"timePlayed"`
	Wins             int           `json:"wins"`
	Losses           int           `json:"losses"`
	Kills            int           `json:"kills"`
	Deaths           int           `json:"deaths"`
	Kpm              float64       `json:"kpm"`
	Spm              float64       `json:"spm"`
	Skill            int           `json:"skill"`
	SoldierImageURL  string        `json:"soldierImageUrl"`
	Rank             Rank          `json:"rank"`
	RankProgress     RankProgress  `json:"rankProgress"`
	FreemiumRank     interface{}   `json:"freemiumRank"`
	Completion       []interface{} `json:"completion"`
	Highlights       interface{}   `json:"highlights"`
	HighlightsByType interface{}   `json:"highlightsByType"`
	EquippedDogtags  interface{}   `json:"equippedDogtags"`
	PersonaID        string        `json:"personaId"`
}
