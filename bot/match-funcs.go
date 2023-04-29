package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (m *liveMatch) Start() {
	// Start() is called _after_ the match is instantiated with all required data to begin

	// err can occur at each of 4 channel creations, clean up if any fail
	matchVcs, err := MakeMatchVoiceChans(m)
	m.VCs = matchVcs
	if err != nil {
		m.CleanupChannels()
	}

	// TODO: clean this up into a util
	rdyBtn := getButton("ready", Q_READY, discordgo.SuccessButton)
	rdyBtn.Emoji = discordgo.ComponentEmoji{Name: "✅"}
	newmessage := discordgo.MessageSend{
		Content: "ping all players msg",
		Embeds:  []*discordgo.MessageEmbed{MakePicksEmbedMessage(m)},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{rdyBtn},
			},
		},
	}

	_, err = Bot.ChannelMessageSendComplex(m.Chan.ID, &newmessage)
	if err != nil {
		log.Error("msg fail to make ", err)
	}
	// TODO: figure out message edit logic
	// ready checks -> alternating pick offers / updated embed
}
func (m *liveMatch) StartPicks() {

}
func (m *liveMatch) Cancel() {}

const (
	ALL_USERS_OPTION = iota
	PLAYERS_OPTION
	VIEWERS_OPTION
	TEAM1_OPTION
	TEAM2_OPTION
)

func (m *liveMatch) GetUsers(option int) []*discordgo.User {
	switch option {
	case ALL_USERS_OPTION:
		return append(append(m.Captains, m.Players...), m.Viewers...)
	case PLAYERS_OPTION:
		return append(m.Team1, m.Team2...)
	case VIEWERS_OPTION:
		return m.Viewers
	case TEAM1_OPTION:
		return m.Team1
	case TEAM2_OPTION:
		return m.Team2
	default:
		return nil
	}
}

func (m *liveMatch) CleanupChannels() {
	if m.VCs.Lobby_vc != nil {
		Bot.ChannelDelete(m.VCs.Lobby_vc.ID)
	}
	if m.VCs.Team_1_vc != nil {
		Bot.ChannelDelete(m.VCs.Team_1_vc.ID)
	}
	if m.VCs.Team_2_vc != nil {
		Bot.ChannelDelete(m.VCs.Team_2_vc.ID)
	}
	if m.VCs.Viewer_vc != nil {
		Bot.ChannelDelete(m.VCs.Viewer_vc.ID)
	}
}