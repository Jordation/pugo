package main

import (
	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func GetLobby(pc *dgo.Channel) *Lobby {
	return &Lobby{
		ParentChannel: pc,
		Game:          &Game{},
	}
}

func (l *Lobby) GetPingAllString() (res string) {
	users := l.GetParticipants()
	for _, user := range users {
		res += user.Mention() + " "
	}
	return
}

func (l *Lobby) SendPickOptions() {
	var err error
	var cpt string
	log.Info("[SENDING PICK OPTIONS]")

	cpt = l.Captains[1].Username
	if l.PickOrder {
		cpt = l.Captains[0].Username
	}

	PicksMsg := &dgo.MessageSend{
		Content: cpt + "'s Choice of Pick",
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					dgo.SelectMenu{
						MenuType:    dgo.StringSelectMenu,
						CustomID:    CPT_PICK,
						Placeholder: "Pick player",
						Options:     MapUsersToPickOptions(l.Players),
					},
				},
			},
		},
	}

	ComplexMsg := &dgo.MessageSend{
		AllowedMentions: &dgo.MessageAllowedMentions{Parse: []dgo.AllowedMentionType{dgo.AllowedMentionTypeUsers}},
		Content:         "Match Beginning! Attention: " + l.GetPingAllString(),
		Embeds:          []*dgo.MessageEmbed{MakePicksEmbedMessage(l)},
	}

	// handle final man going to final side
	if len(l.Players) == 1 {
		l.AddToTeam(l.Players[0], l.PickOrder)

		// remake the embed message to reflect change
		ComplexMsg.Embeds = []*dgo.MessageEmbed{MakePicksEmbedMessage(l)}
		_, err = Bot.ChannelMessageSendComplex(l.Channel.ID, ComplexMsg)
		if err != nil {
			log.Error(err)
		}

		l.Game.Start()
		return

	} else {
		_, err = Bot.ChannelMessageSendComplex(l.Channel.ID, ComplexMsg)
		if err != nil {
			log.Error(err)
		}

		_, err = Bot.ChannelMessageSendComplex(l.Channel.ID, PicksMsg)
		if err != nil {
			log.Error(err)
		}
	}
}
func (l *Lobby) AddToTeam(u *dgo.User, t bool) {
	log.Info("[ADDING TO TEAM]: ", u.Username)
	// flipped so defaults to starting captain first
	if !t {
		l.Game.Team1 = append(l.Game.Team1, u)
	} else {
		l.Game.Team2 = append(l.Game.Team2, u)
	}
	l.RemovePickedUser(u)
}
func (l *Lobby) RemovePickedUser(u *dgo.User) {
	log.Info("[REMOVING FROM PICK POOL]: ", u.Username)
	for i, user := range l.Players {
		if user.ID == u.ID {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			return
		}
	}
}
func (l *Lobby) StartPickPhase() {
	chanCfg := GetChannelConfig(dgo.ChannelTypeGuildText, l.ParentChannel.GuildID, l.GetParticipants())
	nc, err := Bot.CreateTextChannel(chanCfg)
	if err != nil {
		log.Error("[Start Picks Channel Create Error]: ", err)
	}

	if nc == nil {
		panic("channel boom")
	}

	Bot.Lobbies.Set(nc.ID, l)

	// Sets the match name
	l.MatchName = chanCfg.DgoCfg.Name
	// Sets the channel to the newly created channel
	l.Channel = GetChannel(nc.ID)

	l.SendPickOptions()
}

func (l *Lobby) GetParticipants() (res []*dgo.User) {
	res = append(res, l.Players...)
	res = append(res, l.Captains...)
	res = append(res, l.Viewers...)
	return
}
