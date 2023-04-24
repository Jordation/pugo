package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/jordation/go-pug/syncmap"
	log "github.com/sirupsen/logrus"
)

type PugBot struct {
	Self          *dgo.User
	Session       *dgo.Session
	QueueChannels map[string]*QueueChannel
	TESTING_MODE  bool
	PlayerMap     syncmap.UserMap
}

type QueueChannel struct {
	Channel                                   *dgo.Channel
	Queue                                     []*dgo.User
	Lobbies                                   map[uuid.UUID]*ActiveLobby
	MessagesSinceISentOneAboutJoiningTheQueue int
	LastMessageId                             string
}

type ActiveLobby struct {
	GldId     string
	Channel   *dgo.Channel
	Captains  []*dgo.User
	Players   []*dgo.User
	Match     *ActiveMatch
	PickOrder bool
	PickCount int
}

type ActiveMatch struct {
	Team1  []*dgo.User
	Team2  []*dgo.User
	Result int
}

type ChannelConfig struct {
	AllowedUsers []*dgo.User
	GldId        string
	dgoConfig    dgo.GuildChannelCreateData
}

type Config struct {
	Token       string   `json:"token"`
	Prefix      string   `json:"prefix"`
	ListenRooms []string `json:"listen_rooms"`
}

var (
	TestMode   = flag.Bool("bot-mode", true, "should the bot be in testing or production mode")
	MaxPlayers = flag.Int("ppg", 6, "maximum players per game")
	Token      = flag.String("bot-token", "", "The token for the discord bot")
	Prefix     = flag.String("bot-prefix", "", "The prefix the bot will respond to")
	cfg        Config
	Bot        PugBot
)

func init() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &cfg)
	log.Info("[CONFIG SCANNED]: Prefix ", cfg.Prefix)
}
func GetPugBot() {
	bot, err := dgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal(err)
	}
	botUser, err := bot.User("@me")
	if err != nil {
		log.Fatal(err)
	}
	queueChannels := make(map[string]*QueueChannel)
	for _, v := range cfg.ListenRooms {
		queueChannels[v] = &QueueChannel{MessagesSinceISentOneAboutJoiningTheQueue: 5}
	}
	Bot.Session = bot
	Bot.Self = botUser
	Bot.QueueChannels = queueChannels
	Bot.TESTING_MODE = *TestMode
	Bot.PlayerMap.ActiveLobbies = make(map[string]uuid.UUID)
}

func GetUser(uID string) *dgo.User {
	log.Info("[GETTING USER]: ", uID)
	u, err := Bot.Session.User(uID)
	if err != nil {
		log.Fatal(err)
	}
	return u
}
func GetChannel(cID string) *dgo.Channel {
	c, err := Bot.Session.Channel(cID)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func StartBot() {
	GetPugBot()
	if err := Bot.Session.Open(); err != nil {
		log.Fatal(err)
	}

	for chanID, channel := range Bot.QueueChannels {
		channel.Channel = GetChannel(chanID)
		channel.InitQueueChannel()
	}
	Bot.Session.AddHandler(Bot.HandleQueueMessages)
	Bot.Session.AddHandler(Bot.HandleButtonPress)
	Bot.Session.AddHandler(Bot.HandleSelectPlayer)
	select {}
}

func main() {
	StartBot()
	fmt.Println("Shutting down ")
}
