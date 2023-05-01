package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// TODO: handlers ask for bot too often, instead redefine handlers as func return funcs belonging to the bot
var (
	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// BUTTONS
		JOIN_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue, ok := Bot.QueueChannels.Get(i.Interaction.ChannelID)
			if !ok {
				fmtResponse(s, i.Interaction, JOIN_Q_ERR, discordgo.MessageFlagsEphemeral)
				return
			}
			//fmtResponse(s, i.Interaction, JOIN_Q_SUCCESS, discordgo.MessageFlagsEphemeral)
			queue.AddPlayerToQueue(i.Member.User)
			msg := GetQueueMessage(queue.Queue)
			EditQueueMsg(s, i.Interaction, msg)
			if len(queue.Queue) == queue.MaxPlayers {
				queue.CreateMatch()
			}
		},

		LEAVE_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmtResponse(s, i.Interaction, "left q", 0)
		},

		Q_READY: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// get the match object
			match, _ := Bot.Matches.Get(i.ChannelID)

			// get the user who clicked the buttons voice state
			ustate, err := Bot.State.VoiceState(i.GuildID, i.Member.User.ID)
			if err != nil ||
				ustate.ChannelID != match.VCs.Lobby_vc.ID {
				Bot.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "``You must join the match lobby VC before readying up``",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return

			}
			match.ReadyQueue = append(match.ReadyQueue, i.Member.User)
			if len(match.ReadyQueue) != mp {
				// TODO: if not full ready queue, initial message should still update to reflect ready count
				fmtResponse(s, i.Interaction, "wow u readied up nice", discordgo.MessageFlagsEphemeral)
			} else {
				msg := &discordgo.MessageSend{
					Embeds:     []*discordgo.MessageEmbed{MakePicksEmbedMessage(match)},
					Components: getPicksMessage(match.Captains[0].Username+"'s", match.Players).Components,
				}
				EditMatchMsg(s, i.Interaction, msg)
			}
		},

		// END BUTTONS

		// SELECTS
		PLAYER_PICK: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			match, _ := Bot.Matches.Get(i.Interaction.ChannelID)
			picked, _ := Bot.GetUser(strings.Split(i.Interaction.MessageComponentData().Values[0], "_")[0])
			match.AddToTeam(picked)

			cap := match.Captains[0].Username + "'s" + fmt.Sprintf("%v", match.PickOrder)
			if match.PickOrder {
				cap = match.Captains[1].Username + "'s" + fmt.Sprintf("%v", match.PickOrder)
			}
			EditMatchMsg(s, i.Interaction, &discordgo.MessageSend{
				Embeds:     []*discordgo.MessageEmbed{MakePicksEmbedMessage(match)},
				Components: getPicksMessage(cap, match.Players).Components,
			})
		},
		// END SELECTS
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		READY: func(s *discordgo.Session, i *discordgo.InteractionCreate) {},

		ADD_SERV: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Add server to db
			if err := Bot.Db.AddServer(i.GuildID); err != nil {
				fmtResponse(s, i.Interaction, SERVER_EXISTS_ERR, 0)
				return
			}

			// Get the server, create struct to hold data
			g, _ := Bot.GetGuild(i.GuildID)
			newServ := &pugServer{Guild: g}
			Bot.Servers.Set(g.ID, newServ)

			// Acknowledge the interaction
			fmtResponse(s, i.Interaction, SERVER_ADD_SUCCESS, 0)
		},

		ADD_QUEUE: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Add channel to db
			if err := Bot.Db.AddChannel(i.ChannelID, i.GuildID); err != nil {
				fmtResponse(s, i.Interaction, CHANNEL_EXISTS_ERR, 0)
				return
			}

			// Get channel and create struct to hold data
			c, _ := Bot.GetChannel(i.ChannelID)
			newQ := &queueChan{Chan: c, MaxPlayers: mp}
			Bot.QueueChannels.Set(c.ID, newQ)

			// Acknowledge interaction
			fmtResponse(s, i.Interaction, CHANNEL_ADD_SUCCESS, 0)
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

	JOIN_Q_ERR     = "Error joining the queue"
	JOIN_Q_SUCCESS = "You joined the queue"
)
