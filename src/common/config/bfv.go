package config

import "bfv-bot/model/dto"

type Bfv struct {
	// 搜索服务器时 群组的唯一名称 例如: miku/7k7k
	GroupUniName string `mapstructure:"group-uni-name" yaml:"group-uni-name"`
	// 群组名称
	GroupName string `mapstructure:"group-name" yaml:"group-name"`
	// 卡排队判断 小于等于这个值 就触发告警
	BlockingPlayers int          `mapstructure:"blocking-players" yaml:"blocking-players"`
	Server          []ServerInfo `mapstructure:"server" yaml:"server"`
	// 启用黑名单检查
	Active bool
}

type ServerInfo struct {
	Id         string  `mapstructure:"id" yaml:"id"`
	OwnerId    string  `mapstructure:"owner-id" yaml:"owner-id"`
	ServerName string  `mapstructure:"server-name" yaml:"server-name"`
	Kpm        float64 `mapstructure:"kpm" yaml:"kpm"`
	MaxRank    float64 `mapstructure:"max-rank" yaml:"max-rank"`
	MinRank    float64 `mapstructure:"min-rank" yaml:"min-rank"`
	// 运行中获取的id
	gameId string
	// bfvrobot token 与服务器绑定
	token string
	// playerMap 服内玩家数据
	playerMap map[int64]dto.GtBatchStatusData
}

func (b *Bfv) ClearGameId() {
	for i := range b.Server {
		b.Server[i].SetGameId("")
	}
}

func (b *Bfv) GetGameInfo(id string) ServerInfo {
	for _, info := range b.Server {
		if info.Id == id {
			return info
		}
	}
	return ServerInfo{}
}

func (b *Bfv) SetToken(id string, token string) {
	for i := range b.Server {
		info := b.Server[i]
		if info.Id == id {
			b.Server[i].SetToken(token)
		}
	}
}

func (b *Bfv) SetGameId(id string, gameid string) {
	for i := range b.Server {
		info := b.Server[i]
		if info.Id == id {
			b.Server[i].SetGameId(gameid)
		}
	}
}

// SetGameId 设置 GameId
func (s *ServerInfo) SetGameId(gameid string) {
	s.gameId = gameid
}

// GetGameId 获取 GameId
func (s *ServerInfo) GetGameId() string {
	return s.gameId
}

// SetToken 设置 GameId
func (s *ServerInfo) SetToken(token string) {
	s.token = token
}

// GetToken 获取 GameId
func (s *ServerInfo) GetToken() string {
	return s.token
}

func (s *ServerInfo) SetPlayerMap(m map[int64]dto.GtBatchStatusData) {
	s.playerMap = m
}

func (s *ServerInfo) GetPlayerMap() map[int64]dto.GtBatchStatusData {
	return s.playerMap
}
