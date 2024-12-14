package dto

type TemplatePlayerServerInfoModel struct {
	ServerName  string
	MapName     string
	MapMode     string
	ImageBase64 string
}

type TemplatePlayerTeamModel struct {
	TeamName string
	List     []PlayerlistData
}

type PlayerlistData struct {
	Name            string
	KillDeath       string
	KillsPerMinute  string
	Rank            int
	Join            string
	BfBanStatus     int
	BfBanStatusName string
	IsGroupMember   bool
}
