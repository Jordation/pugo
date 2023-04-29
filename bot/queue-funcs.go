package bot

import (
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

func (q *queueChan) SendQueueMessage() {
	Bot.sendSOmeBUttons(q.Chan.ID)
}

func (q *queueChan) AddPlayerToQueue(u *discordgo.User) {
	q.Queue = append(q.Queue, u)
}
func (q *queueChan) RemovePlayerFromQueue(u *discordgo.User) {
	for i, v := range q.Queue {
		if v.ID == u.ID {
			q.Queue = append(q.Queue[:i], q.Queue[i+1:]...)
			return
		}
	}
}
func (q *queueChan) CreateMatch() {

	log.Info(" queue filled, match beginning with: ", q.Queue)

	chanCfg := GetChannelConfig(discordgo.ChannelTypeGuildText, q.Queue, 0)
	newMatchChan, err := Bot.GuildChannelCreateComplex(q.Chan.GuildID, *chanCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Get captains,
	c1, c2 := getCaptains(mp, q.Queue)

	// Cleanup the queue for picking phase later
	q.RemovePlayerFromQueue(c1)
	q.RemovePlayerFromQueue(c2)

	// instantiate the new match with captains, place them on their respective teams
	nm := &liveMatch{Players: q.Queue, Chan: newMatchChan, MatchName: newMatchChan.Name}
	nm.Captains = append(nm.Captains, c1, c2)
	nm.Team1 = append(nm.Team1, c1)
	nm.Team2 = append(nm.Team2, c2)

	nm.Start()
}
