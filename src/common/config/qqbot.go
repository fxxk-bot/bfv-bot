package config

type QQBot struct {
	// 地址
	Address            string  `mapstructure:"address" yaml:"address"`
	Qq                 string  `mapstructure:"qq" yaml:"qq"`
	SuperAdminQq       int64   `mapstructure:"super-admin-qq" yaml:"super-admin-qq"`
	AdminQq            []int64 `mapstructure:"admin-qq" yaml:"admin-qq"`
	AdminGroup         []int64 `mapstructure:"admin-group" yaml:"admin-group"`
	WelcomeMsg         string  `mapstructure:"welcome-msg" yaml:"welcome-msg"`
	ShowPlayerBaseInfo bool    `mapstructure:"show-player-base-info" yaml:"show-player-base-info"`
	// 激活的群号
	ActiveGroup                     []int64                `mapstructure:"active-group" yaml:"active-group"`
	MuteGroup                       MuteGroupConfig        `mapstructure:"mute-group" yaml:"mute-group"`
	CustomCommandKey                CustomCommandKeyConfig `mapstructure:"custom-command-key" yaml:"custom-command-key"`
	BotToBot                        BotToBotConfig         `mapstructure:"bot-bot" yaml:"bot-bot"`
	EnableAutoBindGameId            bool                   `mapstructure:"enable-auto-bind-gameid" yaml:"enable-auto-bind-gameid"`
	EnableAutoKickErrorNickname     bool                   `mapstructure:"enable-auto-kick-error-nickname" yaml:"enable-auto-kick-error-nickname"`
	EnablePlayerlistShowGroupMember bool                   `mapstructure:"enable-playerlist-show-group-member" yaml:"enable-playerlist-show-group-member"`
	EnableRejectJoinRequest         bool                   `mapstructure:"enable-reject-join-request" yaml:"enable-reject-join-request"`
	EnableRejectZeroRankJoinRequest bool                   `mapstructure:"enable-reject-zero-rank-join-request" yaml:"enable-reject-zero-rank-join-request"`
	// 私有属性
	activeGroupMap map[int64]bool
	adminQqMap     map[int64]bool
	adminGroupMap  map[int64]bool
}

type CustomCommandKeyConfig struct {
	Cx          []string `mapstructure:"cx" yaml:"cx"`
	C           []string `mapstructure:"c" yaml:"c"`
	Platoon     []string `mapstructure:"platoon" yaml:"platoon"`
	Banlog      []string `mapstructure:"banlog" yaml:"banlog"`
	Bind        []string `mapstructure:"bind" yaml:"bind"`
	Help        []string `mapstructure:"help" yaml:"help"`
	GroupServer []string `mapstructure:"group-server" yaml:"group-server"`
	Server      []string `mapstructure:"server" yaml:"server"`
	Data        []string `mapstructure:"data" yaml:"data"`
	Task        []string `mapstructure:"task" yaml:"task"`
	Playerlist  []string `mapstructure:"playerlist" yaml:"playerlist"`
	GroupMember []string `mapstructure:"group-member" yaml:"group-member"`
}

type MuteGroupConfig struct {
	Enable      bool       `mapstructure:"enable" yaml:"enable"`
	Start       MuteConfig `mapstructure:"start" yaml:"start"`
	End         MuteConfig `mapstructure:"end" yaml:"end"`
	ActiveGroup []int64    `mapstructure:"active-group" yaml:"active-group"`
}

type MuteConfig struct {
	Time string `mapstructure:"time" yaml:"time"`
	Msg  string `mapstructure:"msg" yaml:"msg"`
}

type BotToBotConfig struct {
	Enable   bool     `mapstructure:"enable" yaml:"enable"`
	BotQq    int64    `mapstructure:"bot-qq" yaml:"bot-qq"`
	Interval int      `mapstructure:"interval" yaml:"interval"`
	Msg      []string `mapstructure:"msg" yaml:"msg"`
}

func (q *QQBot) InitMap() {
	q.activeGroupMap = make(map[int64]bool)
	for _, item := range q.ActiveGroup {
		q.activeGroupMap[item] = true
	}

	q.adminQqMap = make(map[int64]bool)
	for _, item := range q.AdminQq {
		q.adminQqMap[item] = true
	}
	q.adminQqMap[q.SuperAdminQq] = true

	q.adminGroupMap = make(map[int64]bool)
	for _, item := range q.AdminGroup {
		q.adminGroupMap[item] = true
	}

}

func (q *QQBot) IsActiveGroup(item int64) bool {
	_, exists := q.activeGroupMap[item]
	return exists
}

func (q *QQBot) IsActiveAdminQq(item int64) bool {
	_, exists := q.adminQqMap[item]
	return exists
}

func (q *QQBot) IsActiveAdminGroup(item int64) bool {
	_, exists := q.adminGroupMap[item]
	return exists
}
