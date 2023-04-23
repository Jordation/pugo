package main

import (
	"math/rand"
	"strconv"
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

func GetNewLobby() *ActiveLobby {
	return &ActiveLobby{
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
			return
		}
	}
}

func (qc *QueueChannel) InitQueueChannel() {
	qc.Lobbies = make(map[uuid.UUID]*ActiveLobby)
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

// TODO: create text channel which runs the lobby and disallows user messages for that channel
func (qc *QueueChannel) InitNewGame() {
	var (
		gameId         = uuid.New()
		newLobby       = GetNewLobby()
		capId1, capId2 = GetCaptainIds()
	)
	playerIds := make([]string, 0)
	for i, v := range qc.Queue {
		team_number := 1
		if i == capId1 || i == capId2 {
			newLobby.Captains = append(newLobby.Captains, v)

			switch team_number {
			case 1:
				newLobby.Match.Team1 = append(newLobby.Match.Team1, v)
				playerIds = append(playerIds, v.ID)
			case 2:
				newLobby.Match.Team2 = append(newLobby.Match.Team2, v)
				playerIds = append(playerIds, v.ID)
			default:
				log.Fatal("wtf")
			}

			team_number++
		} else {
			newLobby.Players = append(newLobby.Players, v)
			playerIds = append(playerIds, v.ID)
		}
	}
	Bot.PlayerMap.Set(gameId, playerIds...)
	//TODO: make new chan for game
	// make cancel game remove user set from map
	newLobby.Channel = qc.Channel
	qc.Lobbies[gameId] = newLobby
	newLobby.StartPicks()
}

func (al *ActiveLobby) StartPicks() {
	al.SendPickOptions(al.Captains[0])
}

func (al *ActiveLobby) MapRemainingPicks() (res []dgo.SelectMenuOption) {

	for _, player := range al.Players {
		res = append(res, dgo.SelectMenuOption{
			Label: player.Username, // + role + etc...
			Value: player.ID + strconv.Itoa(rand.Intn(100)),
		})
	}
	log.Info("[REMAINING PICKS]: ", res)
	return res
}

// bool so i can flip it in the parent function each time captain swaps
func (al *ActiveLobby) AddToTeam(u *dgo.User, t bool) {
	log.Info("[ADDING TO TEAM]: ", u.Username, t)
	// flipped so defaults to starting captain first
	if !t {
		al.Match.Team1 = append(al.Match.Team1, u)
	} else {
		al.Match.Team2 = append(al.Match.Team2, u)
	}
	al.RemovePickedUser(u)
}
func (al *ActiveLobby) RemovePickedUser(u *dgo.User) {
	log.Info("[REMOVING FROM PICK POOL]: ", u.Username)
	for i, user := range al.Players {
		if user.ID == u.ID {
			al.Players = append(al.Players[:i], al.Players[i+1:]...)
			return
		}
	}
}

func (al *ActiveLobby) SendPickOptions(u *dgo.User) {
	log.Info("[SENDING PICK OPTIONS]")
	components := []dgo.MessageComponent{
		dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.SelectMenu{
					MenuType:    dgo.StringSelectMenu,
					CustomID:    CPT_PICK,
					Placeholder: "Pick player",
					Options:     al.MapRemainingPicks(),
				},
			},
		},
	}
	_, err := Bot.Session.ChannelMessageSendComplex(al.Channel.ID, &dgo.MessageSend{
		Content:    "Message about picks",
		Components: components,
	})
	if err != nil {
		log.Error(err)
	}
}

func (am *ActiveMatch) Start() {
	log.Info("OMG im starting the game", "\n[TEAM 1]: ", am.Team1, "\n[TEAM 2]: ", am.Team2)
}

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
