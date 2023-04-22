package main

import (
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// while the queue isnt full
// post the interactions message
// wait for 5 messages to be sent
// re-post

// IDs for buttons/forms
const (
	JOIN_Q    = "Join"
	JOIN_Q_ID = "JQID"

	LEAVE_Q    = "Leave"
	LEAVE_Q_ID = "LQID"

	CPT_PICK = "CaptainPickChoice"
)

func GetNewGame() *Game {
	return &Game{
		Lobby: &ActiveLobby{},
		Match: &ActiveMatch{},
	}
}

func GetCaptainIds() (int, int) {
	pcount := *MaxPlayers - 1
	rand.Seed(time.Now().UnixNano())

	n1 := rand.Intn(pcount)
	n2 := rand.Intn(pcount)
	if n1 == n2 {
		return GetCaptainIds()
	}

	return n1, n2
}

func GetButton(l, id string, s dgo.ButtonStyle, d bool) *dgo.Button {
	return &dgo.Button{
		Label:    l,
		Style:    s,
		CustomID: id,
		Disabled: d,
	}
}

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

		default:
			log.Info("[BTN CLICK]: no case matched")
		}

	}
}

func (qc *QueueChannel) UserInQueue(u *dgo.User) bool {
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
		if !Bot.TESTING_MODE {
			return
		}
	}

	qc.Queue = append(qc.Queue, u)
	if len(qc.Queue) == *MaxPlayers {
		qc.InitNewGame()
	}
}

func (qc *QueueChannel) RemoveUser(u *dgo.User, i *dgo.Interaction) {
	if qc.UserInQueue(u) {
		Bot.DirectNotifyUser(u, qc.Channel, i, ALREADY_IN_QUEUE_L)
		if !Bot.TESTING_MODE {
			return
		}
	}

	for i, v := range qc.Queue {
		if v.ID == u.ID {
			qc.Queue = append(qc.Queue[:i], qc.Queue[i+1:]...)
		}
	}
}

func (qc *QueueChannel) InitQueueChannel() {
	qc.Games = make(map[uuid.UUID]*Game)
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
	components := []dgo.MessageComponent{
		&dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				GetButton(JOIN_Q, JOIN_Q_ID, dgo.PrimaryButton, false),
				GetButton(LEAVE_Q, LEAVE_Q_ID, dgo.DangerButton, false),
			},
		},
	}
	Bot.Session.ChannelMessageDelete(qc.Channel.ID, qc.LastMessageId)
	m, _ := Bot.Session.ChannelMessageSendComplex(qc.Channel.ID, &dgo.MessageSend{
		Content:    "im sending this message as a test",
		Components: components,
	})
	if m != nil {
		qc.LastMessageId = m.ID
	}
}

func (al *ActiveLobby) SendPickOptions(u *dgo.User) {
	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.SelectMenu{
					MenuType:    dgo.UserSelectMenu,
					CustomID:    CPT_PICK,
					Placeholder: "Pick player",
				},
			},
		},
	}
	_, _ = Bot.Session.ChannelMessageSendComplex(al.Channel.ID, &dgo.MessageSend{
		Content:    "Message about picks",
		Components: components,
	})
}

// TODO: create text channel which runs the lobby and disallows user messages for that channel
func (qc *QueueChannel) InitNewGame() {
	var (
		gameId         = uuid.New()
		newGame        = GetNewGame()
		capId1, capId2 = GetCaptainIds()
	)
	for i, v := range qc.Queue {
		if i == capId1 || i == capId2 {
			newGame.Lobby.Captains = append(newGame.Lobby.Captains, v)
		} else {
			newGame.Lobby.Players = append(newGame.Lobby.Players, v)
		}
	}

	//TODO: make new chan for game
	newGame.Lobby.Channel = qc.Channel
	qc.Games[gameId] = newGame
	newGame.StartPicks()
}

func (g *Game) StartPicks() {
	for i := 0; i < len(g.Lobby.Players); i++ {
		for _, v := range g.Lobby.Captains {
			g.Lobby.SendPickOptions(v)
			// offer choice of pick from pool, remove picked from pool
			// wait for choice to select before going to next
			// filter the selections available by the active pool
		}
	}
}
func (g *Game) StartMatch() {}

func (b *PugBot) DirectNotifyUser(
	u *dgo.User,
	c *dgo.Channel,
	i *dgo.Interaction,
	msg string,
) {
	b.Session.FollowupMessageCreate(i, true, &dgo.WebhookParams{
		Content: msg,
		Flags:   dgo.MessageFlagsEphemeral,
	})
}
