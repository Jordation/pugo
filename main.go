package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type PugBot struct {
	Self          *discordgo.User
	Session       *discordgo.Session
	QueueChannels map[string]*QueueChannel
}

type QueueChannel struct {
	Channel                                   *discordgo.Channel
	Queue                                     []*discordgo.User
	Games                                     []*Game
	MessagesSinceISentOneAboutJoiningTheQueue int
	LastMessageId                             string
}

type NewGame struct {
	ActiveLobby *Lobby
	ActiveGame  *Game
}

type Lobby struct {
	Captains []*discordgo.User
	Players  []*discordgo.User
	GameId   int
}

type Game struct {
	Team1  []*discordgo.User
	Team2  []*discordgo.User
	Result int
	GameId int
}

type Config struct {
	Token       string   `json:"token"`
	Prefix      string   `json:"prefix"`
	ListenRooms []string `json:"listen_rooms"`
}

var (
	MaxPlayers = flag.Int("ppg", 2, "maximum players per game")
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
	bot, err := discordgo.New("Bot " + cfg.Token)
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
}

func GetChannel(cID string) *discordgo.Channel {
	c, err := Bot.Session.Channel(cID)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func StartBot() {
	GetPugBot()
	// look into lazygit
	if err := Bot.Session.Open(); err != nil {
		log.Fatal(err)
	}

	for chanID, channel := range Bot.QueueChannels {
		channel.Channel = GetChannel(chanID)
		channel.StartQueueChannel()
	}
	Bot.Session.AddHandler(Bot.HandleQueueMessages)
	select {}
}

func main() {
	StartBot()
	fmt.Println("Shutting down ")
}
