package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var BeginListenOnce sync.Once

func (b *botService) GetGuild(gID string) (g *discordgo.Guild, err error) {
	log.Info("[GETTING GUILD] ", gID)
	return b.Guild(gID)

}

func (b *botService) GetUser(uID string) (g *discordgo.User, err error) {
	log.Info("[GETTING USER]: ", uID)
	return b.User(uID)
}

func (b *botService) GetChannel(cID string) (g *discordgo.Channel, err error) {
	log.Info("[GETTING CHANNEL]: ", cID)
	return b.Channel(cID)
}

func (b *botService) DirectNotifyUser(
	u *discordgo.User,
	c *discordgo.Channel,
	i *discordgo.Interaction,
	msg string,
) {
	b.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: msg,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

/*
Get all previously registered servers and queue channels and initilize the bot
with maps containing those results
*/
func (b *botService) BeginListening() {
	BeginListenOnce.Do(func() {
		log.Info("[GETTING INITIAL CHANS/SERVERS]:")

		servers, queues := b.Db.GetRegisteredIds()

		for _, serv := range servers {
			guild, err := b.Guild(serv.ServerID)
			if err != nil {
				log.Error("server not found, removing from db and skipping")
				b.Db.RemoveServer(serv.ServerID)
				continue
			}
			b.Servers.Set(serv.ServerID, &pugServer{Guild: guild})
		}

		for _, queue := range queues {
			ch, err := b.Channel(queue.ChanID)
			if err != nil {
				log.Error("channel not found, removing from db and skipping")
				b.Db.RemoveChannel(queue.ChanID)
				continue
			}
			qc := &queueChan{Chan: ch, MaxPlayers: mp}
			b.QueueChannels.Set(queue.ChanID, qc)
			qc.SendQueueMessage()
		}

		if len(queues) == 0 && len(servers) == 0 {
			log.Info("no configured channels yet :)")
		}
	})
}
