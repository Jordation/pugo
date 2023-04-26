package main

import (
	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func (b *Pugo) DirectNotifyUser(
	u *dgo.User,
	c *dgo.Channel,
	i *dgo.Interaction,
	msg string,
) {
	b.FollowupMessageCreate(i, true, &dgo.WebhookParams{
		Content: msg,
		Flags:   dgo.MessageFlagsEphemeral,
	})
}

func (b *Pugo) CreateTextChannel(conf *ChannelConfig) (*dgo.Channel, error) {
	return b.GuildChannelCreateComplex(conf.GldID, conf.DgoCfg)
}

func (b *Pugo) HandleQueueMessages(s *dgo.Session, m *dgo.MessageCreate) {
	if m.Author.ID == b.Self.ID {
		return
	}

	var queueIds []string
	for k := range b.QueueChannels.m {
		queueIds = append(queueIds, k)
	}
	if !slices.Contains(queueIds, m.ChannelID) {
		return
	}

	v, ok := b.QueueChannels.Get(m.ChannelID)
	if ok {
		log.Info("[MESSAGE IN Q CHANNEL]")
		v.MsgTicker++
	}
}

func (b *Pugo) HandleButtonPress(s *dgo.Session, i *dgo.InteractionCreate) {
	if i.Type == dgo.InteractionMessageComponent &&
		i.MessageComponentData().ComponentType == dgo.ButtonComponent {
		bid := i.MessageComponentData().CustomID

		if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseDeferredMessageUpdate,
		}); err != nil {
			log.Error("[INTERACTION]: ", err)
		}

		switch bid {
		case JOIN_Q_ID:
			v, ok := b.QueueChannels.Get(i.ChannelID)
			if !ok {
				log.Error("[ERR]: Queue channel not found")
			}
			v.AddUserToQueue(i.Member.User, i.Interaction)
		case LEAVE_Q_ID:
			v, ok := b.QueueChannels.Get(i.ChannelID)
			if !ok {
				log.Error("[ERR]: Queue channel not found")
			}
			v.RemoveUserFromQueue(i.Member.User, i.Interaction)

		default:
			log.Info("[BTN CLICK]: no case matched")
		}

	}
}

func (b *Pugo) HandleSelectPlayer(s *dgo.Session, i *dgo.InteractionCreate) {
	if i.Type == dgo.InteractionMessageComponent &&
		i.MessageComponentData().ComponentType == dgo.SelectMenuComponent {
		selId := i.MessageComponentData().CustomID
		s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{Type: dgo.InteractionResponseChannelMessageWithSource})
		switch selId {
		case CPT_PICK:
			lobby, ok := b.Lobbies.Get(i.ChannelID)
			if !ok {
				log.Error("[ERR]: Lobby channel not found")
				return
			}

			lobby.AddToTeam(i.Member.User, lobby.PickOrder)
			lobby.PickCount++

			if lobby.PickCount%2 != 0 {
				lobby.PickOrder = !lobby.PickOrder
			}
			lobby.SendPickOptions()

			/*
				 TODO: This is for real usage, currently just adds self to team
				lobby.AddToTeam(GetUser(i.Interaction.MessageComponentData().Values[0]), lobby.PickOrder)
				if there is only 1 player left in the pool, auto assign to last time, end picks
			*/
		}
	}
}
