package main

import (
	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (b *PugBot) HandleQueueMessages(s *dgo.Session, m *dgo.MessageCreate) {
	if m.Author.ID == b.Self.ID {
		return
	}
	v, ok := b.QueueChannels[m.ChannelID]
	if ok {
		log.Info("[MESSAGE IN Q CHANNEL]")
		v.MessagesSinceISentOneAboutJoiningTheQueue++
	}
}

func (b *PugBot) HandleButtonPress(s *dgo.Session, i *dgo.InteractionCreate) {
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
			log.Info("[USER JOIN QUEUE]: UID:", i.Member.User.ID)
			b.QueueChannels[i.ChannelID].AddUserToQueue(i.Member.User, i.Interaction)
		case LEAVE_Q_ID:
			log.Info("[USER LEAVE QUEUE]: UID:", i.Member.User.ID)
			b.QueueChannels[i.ChannelID].RemoveUser(i.Member.User, i.Interaction)
		default:
			log.Info("[BTN CLICK]: no case matched")
		}

	}
}

func (b *PugBot) HandleSelectPlayer(s *dgo.Session, i *dgo.InteractionCreate) {
	if i.Type == dgo.InteractionMessageComponent &&
		i.MessageComponentData().ComponentType == dgo.SelectMenuComponent {
		selId := i.MessageComponentData().CustomID
		if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseDeferredMessageUpdate,
		}); err != nil {
			log.Error("[INTERACTION]: ", err)
		}
		switch selId {
		case CPT_PICK:
			//TODO: fix getuser
			lobbyUuid := b.PlayerMap.Get(i.Member.User.ID)
			lobby := b.QueueChannels[i.ChannelID].Lobbies[lobbyUuid]
			/* 			selectedUserId := i.Interaction.MessageComponentData().Values[0] */
			lobby.AddToTeam(i.Member.User, lobby.PickOrder)
			// send next pick message
			// if there is only 1 player left in the pool, auto assign to last time, end picks
			if len(lobby.Players) != 1 {
				lobby.SendPickOptions(i.Member.User)
			} else {
				lobby.AddToTeam(lobby.Players[0], lobby.PickOrder)
				lobby.Match.Start()
			}
		}
	}
}
