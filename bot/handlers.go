package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ButtonHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		JOIN_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue, ok := Bot.QueueChannels.Get(i.Interaction.ChannelID)
			if !ok {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "``Error joining queue``",
				})
				return
			}
			fmt.Print(queue.MaxPlayers)
		},
		LEAVE_Q: func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ready": func(s *discordgo.Session, i *discordgo.InteractionCreate) {},

		"add-server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := Bot.Db.AddServer(i.GuildID); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "``Server already registered``",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "``Server registered successfully``",
				},
			})
		},

		"add-queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := Bot.Db.AddChannel(i.ChannelID, i.GuildID); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "``Channel already registered``",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "``Channel registered successfully, hf ;)``",
				},
			})
		},
	}

	Commands = []discordgo.ApplicationCommand{
		{
			Name:        "ready",
			Description: "Ready up for the lobby that just begun",
		},
		{
			Name:        "add-server",
			Description: "Add this server to the pug server list (required to add a queue)",
		},
		{
			Name:        "add-queue",
			Description: "Initialize this channel as a queue channel",
		},
	}
)

const (
	JOIN_Q  = "JOIN_QUEUE_BUTTON"
	LEAVE_Q = "LEAVE_QUEUE_BUTTON"
)

const (
	SERVER_EXISTS_ERR  = "Server already registered"
	SERVER_ADD_SUCCESS = "Server registered successfully"

	CHANNEL_EXISTS_ERR  = "Channel already registered"
	CHANNEL_ADD_SUCCESS = "Channel registered successfully, hf ;)"

	JOIN_Q_ERR = "Error joining the queue"
)
