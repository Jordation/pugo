package main

import (
	"time"

	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (qc *QueueChannel) UserInQueue(u *dgo.User) bool {
	if Bot.TESTING_MODE {
		return false
	}
	for _, v := range qc.Queue {
		if v.ID == u.ID {
			return true
		}
	}
	return false
}
func (qc *QueueChannel) AddUserToQueue(u *dgo.User, i *dgo.Interaction) {
	if qc.UserInQueue(u) {
		Bot.DirectNotifyUser(u, qc.Channel, i, ALREADY_IN_QUEUE_J)
		return
	}

	log.Info("[USER JOIN QUEUE]: UID:", i.Member.User.ID)

	qc.Queue = append(qc.Queue, u)
	if len(qc.Queue) == *MaxPlayers {
		qc.StartNewLobby()
	}

	log.Info("[QUEUE STATE]: ", qc.Queue)
}
func (qc *QueueChannel) RemoveUserFromQueue(u *dgo.User, i *dgo.Interaction) {
	if qc.UserInQueue(u) {
		Bot.DirectNotifyUser(u, qc.Channel, i, ALREADY_IN_QUEUE_L)
		return
	}

	log.Info("[USER LEAVE QUEUE]: UID:", i.Member.User.ID)

	for i, v := range qc.Queue {
		if v.ID == u.ID {
			qc.Queue = append(qc.Queue[:i], qc.Queue[i+1:]...)
			return
		}
	}

	log.Info("[QUEUE STATE]: ", qc.Queue)
}
func (qc *QueueChannel) InitQueueChannel() {
	go func() {
		for {
			if qc.MsgTicker >= 5 {
				qc.SendQueueOptions()
				qc.MsgTicker = 0
			}
			time.Sleep(time.Second * 1)
		}
	}()
}
func (qc *QueueChannel) SendQueueOptions() {
	components := []dgo.MessageComponent{
		&dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				GetButton(JOIN_Q, JOIN_Q_ID, dgo.PrimaryButton, false),
				GetButton(LEAVE_Q, LEAVE_Q_ID, dgo.DangerButton, false),
			},
		},
	}

	Bot.ChannelMessageDelete(qc.Channel.ID, qc.LastBotMsgId)
	m, _ := Bot.ChannelMessageSendComplex(qc.Channel.ID, &dgo.MessageSend{
		Content:    "im sending this message as a test",
		Components: components,
	})
	if m != nil {
		qc.LastBotMsgId = m.ID
	}
}
func (qc *QueueChannel) StartNewLobby() {
	var (
		nl             = GetLobby(qc.Channel)
		cap1id, cap2id = GetCaptainIds()
		team_number    = 1
	)
	for i, v := range qc.Queue {
		if i != cap1id && i != cap2id {
			// If its not a captain ID, append to players list
			nl.Players = append(nl.Players, v)
		} else {
			// else append captain to corresponding team
			nl.Captains = append(nl.Captains, v)
			switch team_number {
			case 1:
				nl.Game.Team1 = append(nl.Game.Team1, v)
				team_number++
			case 2:
				nl.Game.Team2 = append(nl.Game.Team2, v)
			default:
				log.Fatal("wtf")
			}
		}
	}
	go nl.StartPickPhase()
	qc.Queue = make([]*dgo.User, 0)
}
