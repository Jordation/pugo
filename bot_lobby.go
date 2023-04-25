package main

import (
	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func GetLobby() *Lobby {
	return &Lobby{
		Game: &Game{},
	}
}

func (l *Lobby) SendPickOptions(c *dgo.User) {
	log.Info("[SENDING PICK OPTIONS]")
	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.SelectMenu{
					MenuType:    dgo.StringSelectMenu,
					CustomID:    CPT_PICK,
					Placeholder: "Pick player",
					Options:     MapUsersToPickoptions(l.GetParticipants()),
				},
			},
		},
	}
	_, err := Bot.Session.ChannelMessageSendComplex(l.Channel.ID, &dgo.MessageSend{
		Content:    "Message about picks",
		Components: components,
	})
	if err != nil {
		log.Error(err)
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
	chanCfg := GetChannelConfig(dgo.ChannelTypeGuildText, l.GldId, l.GetParticipants())
	nc, err := Bot.CreateTextChannel(chanCfg)
	if err != nil {
		log.Error("[Start Picks Channel Create Error]: ", err)
	}
	if nc == nil {
		panic("channel boom")
	}
	l.Channel = GetChannel(nc.ID)
	l.SendPickOptions(l.Captains[0])
}
func (l *Lobby) GetParticipants() (res []*dgo.User) {
	res = append(res, l.Players...)
	res = append(res, l.Captains...)
	res = append(res, l.Viewers...)
	return
}
