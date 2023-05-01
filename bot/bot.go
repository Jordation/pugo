package bot

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/jordation/go-pug/db"
)

type botService struct {
	*discordgo.Session
	Db            *db.DatabaseConnection
	Servers       *db.ConcurrentMap[string, *pugServer]
	QueueChannels *db.ConcurrentMap[string, *queueChan]
	Matches       *db.ConcurrentMap[string, *liveMatch]
	AppID         string
}

type pugServer struct {
	Guild *discordgo.Guild
}

type queueChan struct {
	Chan  *discordgo.Channel
	Queue []*discordgo.User

	MaxPlayers    int
	MessageTicker int
	InitialMsg    *discordgo.Message
}

type liveMatch struct {
	Chan *discordgo.Channel

	VCs        *vcs
	Captains   []*discordgo.User
	Players    []*discordgo.User
	ReadyQueue []*discordgo.User
	Viewers    []*discordgo.User

	Team1 []*discordgo.User
	Team2 []*discordgo.User

	PickOrder bool
	MatchName string
}
type vcs struct {
	Lobby_vc  *discordgo.Channel
	Team_1_vc *discordgo.Channel
	Team_2_vc *discordgo.Channel
	Viewer_vc *discordgo.Channel
}
type config struct {
	Token string
	AppID string
}

var (
	Bot = &botService{}
	mp  = 4
	cfg config
)

func GetBotService() *botService {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, &cfg)

	Bot.Session, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal(err)
	}
	Bot.AppID = cfg.AppID

	if err := Bot.Open(); err != nil {
		log.Fatal(err)
	}

	Bot.Servers = db.NewConcurrentMap[string, *pugServer]()
	Bot.QueueChannels = db.NewConcurrentMap[string, *queueChan]()
	Bot.Matches = db.NewConcurrentMap[string, *liveMatch]()

	db := db.GetDb(true)
	Bot.Db = db

	Bot.BeginListening()
	Bot.StateEnabled = true
	return Bot
}
