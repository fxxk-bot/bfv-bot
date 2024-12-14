package dto

type TofResp struct {
	Success int     `json:"success"`
	Code    string  `json:"code"`
	Data    TofData `json:"data"`
}
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Requirements struct {
	Desc          string `json:"desc"`
	RequiredValue string `json:"requiredValue"`
	Code          string `json:"code"`
}
type Rewards struct {
	ItemType  string `json:"itemType"`
	AssetGUID string `json:"assetGuid"`
	Quantity  string `json:"quantity"`
}
type Achievement struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Desc          string         `json:"desc"`
	Dependencies  []string       `json:"dependencies"`
	Requirements  []Requirements `json:"requirements"`
	Rewards       []Rewards      `json:"rewards"`
	Score         int            `json:"score"`
	Image         string         `json:"image"`
	Rarity        string         `json:"rarity"`
	PublishedAtMs interface{}    `json:"publishedAtMs"`
	StartsAtMs    interface{}    `json:"startsAtMs"`
	ExpiresAtMs   interface{}    `json:"expiresAtMs"`
	AwardGroup    string         `json:"awardGroup"`
}
type StoryEvents struct {
	ID          string      `json:"id"`
	ImageURL    string      `json:"imageUrl"`
	Position    Position    `json:"position"`
	Achievement Achievement `json:"achievement"`
	IconSize    string      `json:"iconSize"`
}
type Weeks struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	StartTimestamp      string        `json:"startTimestamp"`
	EndTimestamp        string        `json:"endTimestamp"`
	ImageURL            string        `json:"imageUrl"`
	RewardImageURL      string        `json:"rewardImageUrl"`
	RewardBackgroundURL string        `json:"rewardBackgroundUrl"`
	StoryEvents         []StoryEvents `json:"storyEvents"`
}
type Events struct {
	ID          string  `json:"id"`
	EventType   string  `json:"eventType"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	Weeks       []Weeks `json:"weeks"`
}
type TofData struct {
	ID                  string        `json:"id"`
	StartTimestamp      string        `json:"startTimestamp"`
	EndTimestamp        string        `json:"endTimestamp"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	IconURL             string        `json:"iconUrl"`
	Events              []Events      `json:"events"`
	VideoURL            string        `json:"videoUrl"`
	TrackedAchievements []interface{} `json:"trackedAchievements"`
	NoProgression       bool          `json:"noProgression"`
}

type Rows struct {
	Ys []int
	Xs []int
}

type Node struct {
	Name         string
	Position     Position
	Dependencies string
	Requirements []string
	Rewards      []string
}
