package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// while the queue isnt full
// post the interactions message
// wait for 5 messages to be sent
// re-post

const (
	JOIN_Q    = "Join"
	JOIN_Q_ID = "JQID"

	LEAVE_Q    = "Leave"
	LEAVE_Q_ID = "LQID"
)

func GetButton(l, id string, s discordgo.ButtonStyle, d bool) *discordgo.Button {
	return &discordgo.Button{
		Label:    l,
		Style:    s,
		CustomID: id,
		Disabled: d,
	}
}

func (b *PugBot) HandleQueueMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.Self.ID {
		return
	}
	v, ok := b.QueueChannels[m.ChannelID]
	if ok {
		log.Info("[MESSAGE IN Q CHANNEL]")
		v.MessagesSinceISentOneAboutJoiningTheQueue++
	}
}

func (b *PugBot) HandleUserJoinQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func (qc *QueueChannel) StartQueueChannel() {
	go func() {
		for {
			if qc.MessagesSinceISentOneAboutJoiningTheQueue >= 5 {
				qc.SendQueueOptions()
				qc.MessagesSinceISentOneAboutJoiningTheQueue = 0
			}
			time.Sleep(time.Second * 5)
		}
	}()

}

func (qc *QueueChannel) SendQueueOptions() {
	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				GetButton(JOIN_Q, JOIN_Q_ID, discordgo.PrimaryButton, false),
				GetButton(LEAVE_Q, LEAVE_Q_ID, discordgo.DangerButton, false),
			},
		},
	}
	Bot.Session.ChannelMessageDelete(qc.Channel.ID, qc.LastMessageId)
	m, _ := Bot.Session.ChannelMessageSendComplex(qc.Channel.ID, &discordgo.MessageSend{
		Content:    "im sending this message as a test",
		Components: components,
	})
	if m != nil {
		qc.LastMessageId = m.ID
	}
}
