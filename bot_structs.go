package main

import dgo "github.com/bwmarrin/discordgo"

type Pugo struct {
	*dgo.Session
	Self          *dgo.User
	QueueChannels *QC_SM
	Lobbies       *L_SM
	TESTING_MODE  bool
}

type QueueChannel struct {
	Channel      *dgo.Channel
	Queue        []*dgo.User
	MsgTicker    int
	LastBotMsgId string
}

type Lobby struct {
	GldId         string
	ParentChannel *dgo.Channel
	Channel       *dgo.Channel
	Captains      []*dgo.User
	Players       []*dgo.User
	Viewers       []*dgo.User
	Game          *Game
	PickCount     int
	PickOrder     bool
}

type Game struct {
	Team1  []*dgo.User
	Team2  []*dgo.User
	Result int
}

type ChannelConfig struct {
	AllowedUsers []*dgo.User
	GldID        string
	DgoCfg       dgo.GuildChannelCreateData
}

type BotConfig struct {
	Token       string   `json:"token"`
	Prefix      string   `json:"prefix"`
	ListenRooms []string `json:"listen_rooms"`
}
