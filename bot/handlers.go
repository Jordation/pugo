package bot

import (
	"github.com/bwmarrin/discordgo"
)

// TODO: handlers ask for bot too often, instead redefine handlers as func return funcs belonging to the bot
var (
	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// BUTTONS
		JOIN_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue, ok := Bot.QueueChannels.Get(i.Interaction.ChannelID)
			if !ok {
				interactionRespond(s, i.Interaction, JOIN_Q_ERR, nil)
				return
			}
			interactionRespond(s, i.Interaction, "joined q", nil)
			queue.AddPlayerToQueue(i.Member.User)
			if len(queue.Queue) == queue.MaxPlayers {
				queue.CreateMatch()
			}
		},

		LEAVE_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			interactionRespond(s, i.Interaction, "left q", nil)
		},

		Q_READY: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// get the match object
			match, _ := Bot.Matches.Get(i.ChannelID)

			// get the user who clicked the buttons voice state
			ustate, err := Bot.State.VoiceState(i.GuildID, i.Member.User.ID)
			if ustate.ChannelID != match.VCs.Lobby_vc.ID || err != nil {
				Bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: 1,
				})
			}

		},

		// END BUTTONS

		// SELECTS

		PLAYER_PICK: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},
		// END SELECTS
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		READY: func(s *discordgo.Session, i *discordgo.InteractionCreate) {},

		ADD_SERV: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := Bot.Db.AddServer(i.GuildID); err != nil {
				interactionRespond(s, i.Interaction, SERVER_EXISTS_ERR, nil)
				return
			}
			interactionRespond(s, i.Interaction, SERVER_ADD_SUCCESS, nil)
		},

		ADD_QUEUE: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := Bot.Db.AddChannel(i.ChannelID, i.GuildID); err != nil {
				interactionRespond(s, i.Interaction, CHANNEL_EXISTS_ERR, nil)
				return
			}
			interactionRespond(s, i.Interaction, CHANNEL_ADD_SUCCESS, nil)
		},
	}

	Commands = []discordgo.ApplicationCommand{
		{
			Name:        READY,
			Description: "Ready up for the lobby that just begun",
		},
		{
			Name:        ADD_SERV,
			Description: "Add this server to the pug server list (required to add a queue)",
		},
		{
			Name:        ADD_QUEUE,
			Description: "Initialize this channel as a queue channel",
		},
	}
)

const (
	READY     = "ready"
	ADD_SERV  = "register"
	ADD_QUEUE = "init"
)

const (
	JOIN_Q      = "JOIN_QUEUE_BUTTON"
	LEAVE_Q     = "LEAVE_QUEUE_BUTTON"
	Q_READY     = "QUEUE_READY_UP"
	PLAYER_PICK = "PLAYER_PICK_FOR_MATCH"
)

const (
	SERVER_EXISTS_ERR  = "Server already registered"
	SERVER_ADD_SUCCESS = "Server registered successfully"

	CHANNEL_EXISTS_ERR  = "Error registering queue, channel already registered"
	CHANNEL_ADD_SUCCESS = "Channel registered successfully, hf ;)"

	JOIN_Q_ERR = "Error joining the queue"
)
