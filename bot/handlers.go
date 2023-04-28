package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ready": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},

		"add-server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "``Server registered successfully``",
				},
			})

			if err := Bot.Db.AddServer(i.GuildID); err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "``Server already registered.``",
				})
			}

		},

		"add-queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			if err := Bot.Db.AddChannel(i.ChannelID, i.GuildID); err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "There was an error initializing the queue channel",
				})
			}
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
