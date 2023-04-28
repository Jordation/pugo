package main

import (
	"flag"
	"os"
	"os/signal"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/jordation/go-pug/bot"
	log "github.com/sirupsen/logrus"
)

var (
	TestMode      = flag.Bool("bot-mode", true, "should the bot be in testing or production mode")
	MaxPlayers    = flag.Int("ppg", 4, "maximum players per game")
	CaptainRoleID = flag.String("crId", "1101122692679225465", "The id of the captain role for the server")
)

func main() {
	// Get the bot, open the session
	Bot := bot.GetBotService()

	defer Bot.Close()

	// Register interaction handlers
	Bot.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		switch i.Type {
		case dgo.InteractionApplicationCommand:
			if h, ok := bot.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case dgo.InteractionModalSubmit:
		default:
			log.Info("unhandled interaction type")
		}
	})

	// Register custom /commands
	cmdIDs := make(map[string]string, len(bot.Commands))
	for _, cmd := range bot.Commands {
		rcmd, err := Bot.ApplicationCommandCreate(Bot.AppID, "", &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", rcmd.Name, err)
		}
		cmdIDs[rcmd.ID] = rcmd.Name
	}

	// Wait for interrupt sig
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	//cleanup
	log.Info("Graceful shutdown Started...")

	for id, name := range cmdIDs {
		err := Bot.ApplicationCommandDelete(Bot.AppID, "", id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

	log.Info("Graceful shutdown complete")
}
