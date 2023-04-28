package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"

	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	TestMode      = flag.Bool("bot-mode", true, "should the bot be in testing or production mode")
	MaxPlayers    = flag.Int("ppg", 4, "maximum players per game")
	Token         = flag.String("bot-token", "", "The token for the discord bot")
	AppId         = flag.String("app-id", "", "The application id")
	CaptainRoleID = flag.String("crId", "1101122692679225465", "The id of the captain role for the server")
	Prefix        = flag.String("bot-prefix", "", "The prefix the bot will respond to")
	BotCfg        BotConfig
	Bot           Pugo
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
	u, err := Bot.User(uID)
	if err != nil {
		log.Fatal(err)
	}
	return u
}
func GetChannel(cID string) *dgo.Channel {
	c, err := Bot.Channel(cID)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func StartBot() {
	InitBot()
	if err := Bot.Open(); err != nil {
		log.Fatal(err)
	}

	for chanID, channel := range Bot.QueueChannels.m {
		channel.Channel = GetChannel(chanID)
		channel.InitQueueChannel()
	}
	Bot.AddHandler(Bot.HandleQueueMessages)
	Bot.AddHandler(Bot.HandleButtonPress)
	Bot.AddHandler(Bot.HandleSelectPlayer)

	Bot.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		switch i.Type {
		case dgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case dgo.InteractionModalSubmit:
		default:
			log.Info("unhandled interaction type")
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := Bot.ApplicationCommandCreate(BotCfg.AppId, "", &cmd)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", rcmd.Name, err)
		}
		cmdIDs[rcmd.ID] = rcmd.Name
	}

	defer Bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	//cleanup
	for id, name := range cmdIDs {
		err := Bot.ApplicationCommandDelete(BotCfg.AppId, "", id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
	log.Info("Graceful shutdown")
}

func main() {
	StartBot()
}

var (
	commands = []dgo.ApplicationCommand{
		{
			Name:        "pick",
			Description: "pick a player for the game",
		},
	}
	commandsHandlers = map[string]func(s *dgo.Session, i *dgo.InteractionCreate){
		"pick": func(s *dgo.Session, i *dgo.InteractionCreate) {
			lobby, ok := Bot.Lobbies.Get(i.ChannelID)
			if ok {
				err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseModal,
					Data: &dgo.InteractionResponseData{
						CustomID: "PICK",
						Title:    "Pick a player",
						Components: []dgo.MessageComponent{
							dgo.SelectMenu{
								MenuType:    dgo.StringSelectMenu,
								CustomID:    "PICKPICK",
								Placeholder: "do something",
								Options:     MapUsersToPickOptions(lobby.Players),
							},
						},
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
	}
)
