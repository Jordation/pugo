package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	TestMode   = flag.Bool("bot-mode", true, "should the bot be in testing or production mode")
	MaxPlayers = flag.Int("ppg", 6, "maximum players per game")
	Token      = flag.String("bot-token", "", "The token for the discord bot")
	Prefix     = flag.String("bot-prefix", "", "The prefix the bot will respond to")
	BotCfg     BotConfig
	Bot        Pugo
)

func init() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &BotCfg)
	log.Info("[CONFIG SCANNED]: Prefix ", BotCfg.Prefix)
}
func InitBot() {
	var (
		queueChannels = &QC_SM{}
		lobbies       = &L_SM{}
	)

	queueChannels.Make()
	lobbies.Make()
	Bot.QueueChannels = queueChannels
	Bot.Lobbies = lobbies

	sesh, err := dgo.New("Bot " + BotCfg.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := sesh.User("@me")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range BotCfg.ListenRooms {
		queueChannels.Set(v, &QueueChannel{MsgTicker: 5})
	}

	Bot.Session = sesh
	Bot.Self = bot
	Bot.TESTING_MODE = *TestMode
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
	InitBot()
	if err := Bot.Session.Open(); err != nil {
		log.Fatal(err)
	}

	for chanID, channel := range Bot.QueueChannels.m {
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
