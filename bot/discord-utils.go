package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

/* General use functions that create values or datastructures to be used throughout the bot */

func MapUsers(u []*dgo.User, mention bool) (res string) {
	for i, v := range u {
		if mention {
			res += "<@" + v.ID + ">"
		} else {
			res += v.Username
		}
		if i+1 == len(u) {
			continue
		}
		res += ", "
	}
	return
}
func MapTeamDisplay(u []*dgo.User) (res string) {
	res += "```md"
	for _, user := range u {
		res += "\n- " + user.Username
	}

	res += "\n```"
	return
}

func MakePicksEmbedFields(lm *liveMatch) (res []*dgo.MessageEmbedField) {
	res = append(res, &dgo.MessageEmbedField{
		Name: "Available picks",
		Value: "```md\n" +
			MapUsers(lm.Players, false) +
			"\n```",
	})
	res = append(res, &dgo.MessageEmbedField{
		Name:   "Team 1",
		Value:  MapTeamDisplay(lm.Team1),
		Inline: true,
	})
	res = append(res, &dgo.MessageEmbedField{
		Name:   "Team 2",
		Value:  MapTeamDisplay(lm.Team2),
		Inline: true,
	})
	// if the ready queue isnt full, embed shows its status
	if len(lm.ReadyQueue) != mp {
		res = append(res, &dgo.MessageEmbedField{
			Name:  "Waiting for the players to ready up",
			Value: "```md\n- " + strconv.Itoa(len(lm.ReadyQueue)) + " / " + strconv.Itoa(mp) + "\n```",
		})
	} else {
		cap := lm.Captains[0].Username + "'s"
		if lm.PickOrder {
			cap = lm.Captains[1].Username + "'s"
		}
		res = append(res, &dgo.MessageEmbedField{
			Name:  "Picking Phase",
			Value: "``" + cap + " turn to pick.``",
		})
	}

	return
}

func MakePicksEmbedMessage(lm *liveMatch) *dgo.MessageEmbed {
	res := &dgo.MessageEmbed{}

	// TODO: unhardcode this
	res.Author = &dgo.MessageEmbedAuthor{Name: "Division 1"}
	res.Title = "Welcome to match " + lm.MatchName + "!"

	res.Description =
		"```md\n# Welcome! #" +
			"\n\n" +
			"Captains for this match are: " +
			"\n- " + lm.Captains[0].Username +
			"\n- " + lm.Captains[1].Username +
			"\n```"
	res.Fields = MakePicksEmbedFields(lm)
	return res
}

func getPicksMessage(captain string, players []*dgo.User) *dgo.MessageSend {
	return &dgo.MessageSend{
		Content: "``It is " + captain + " turn to pick",
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					dgo.SelectMenu{
						MenuType: dgo.StringSelectMenu,
						CustomID: PLAYER_PICK,
						Options:  MapUsersToPickOptions(players),
					},
				},
			},
		},
	}
}

func getCaptains(maxPlayers int, players []*dgo.User) (c1, c2 *dgo.User) {
	rand.Seed(time.Now().UnixNano())

	n1 := rand.Intn(maxPlayers)
	n2 := rand.Intn(maxPlayers)
	if n1 == n2 {
		return getCaptains(maxPlayers, players)
	}
	return players[n1], players[n2]
}

func getButton(l, id string, s dgo.ButtonStyle) *dgo.Button {
	return &dgo.Button{
		Label:    l,
		Style:    s,
		CustomID: id,
	}
}

func EditMatchMsg(
	s *dgo.Session,
	i *dgo.Interaction,
	msg *dgo.MessageSend,
) {
	s.InteractionRespond(i, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
		Data: &dgo.InteractionResponseData{
			Embeds:     msg.Embeds,
			Content:    msg.Content,
			Components: msg.Components,
		},
	})
}
func EditQueueMsg(
	s *dgo.Session,
	i *dgo.Interaction,
	msg *dgo.MessageSend,
) {
	s.InteractionRespond(i, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseUpdateMessage,
		Data: &dgo.InteractionResponseData{
			Content:    msg.Content,
			Components: msg.Components,
		},
	})
}

func fmtResponse(
	s *dgo.Session,
	i *dgo.Interaction,
	msg string,
	flag dgo.MessageFlags,
) {
	s.InteractionRespond(i, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: fmt.Sprintf("``%v``", msg),
			Flags:   flag,
		},
	})
}

// Use 0 limit for text channel
func GetChannelConfig(
	ctype dgo.ChannelType,
	ValidUsers []*dgo.User,
	userLimit int,
	chanName string,
) *dgo.GuildChannelCreateData {

	switch ctype {
	case dgo.ChannelTypeGuildText:
		return &dgo.GuildChannelCreateData{
			Type:                 ctype,
			PermissionOverwrites: mapUserPerms(ValidUsers, dgo.PermissionViewChannel),
			Name:                 chanName,
		}
	case dgo.ChannelTypeGuildVoice:
		return &dgo.GuildChannelCreateData{
			Type:                 ctype,
			PermissionOverwrites: mapUserPerms(ValidUsers, dgo.PermissionViewChannel),
			Name:                 chanName,
			UserLimit:            userLimit,
		}
	default:
		log.Error("[UNHANDLED CHANNEL TYPE]: GetChannelConfig")
		return nil
	}
}
func mapUserPerms(users []*dgo.User, permType int64) (res []*dgo.PermissionOverwrite) {
	for _, v := range users {
		res = append(res, &dgo.PermissionOverwrite{
			ID:    v.ID,
			Type:  dgo.PermissionOverwriteTypeMember,
			Allow: int64(permType),
		})
	}
	return
}

func MapUsersToPickOptions(u []*dgo.User) (res []dgo.SelectMenuOption) {
	for _, player := range u {
		res = append(res, dgo.SelectMenuOption{
			Label: player.Username, // + role + etc...
			// TODO: the random int is for self joining test
			// breaks GetUser function in HandleSelectPlayer
			// splitting on _ for now
			Value: player.ID + "_" + strconv.Itoa(rand.Intn(10000)),
		})
	}
	for _, p := range res {
		log.Info("[REMAINING PICKS]: ", p.Label)
	}
	return
}

func MakeMatchVoiceChans(m *liveMatch, rdy bool) (*vcs, error) {
	res := &vcs{}

	// lobby vc
	if !rdy {
		matchUsers := m.GetUsers(ALL_USERS_OPTION)
		lbvc := GetChannelConfig(dgo.ChannelTypeGuildVoice, matchUsers, len(matchUsers), "Lobby - "+m.MatchName)
		lbvc_CHAN, err := Bot.GuildChannelCreateComplex(m.Chan.GuildID, *lbvc)
		if err != nil {
			return res, err
		}
		res.Lobby_vc = lbvc_CHAN
		return res, nil
	}

	// Team vcs
	res.Lobby_vc = m.VCs.Lobby_vc
	t1vc := GetChannelConfig(dgo.ChannelTypeGuildVoice, m.GetUsers(TEAM1_OPTION), mp/2, "Team 1 - "+m.MatchName)
	t1vc_CHAN, err := Bot.GuildChannelCreateComplex(m.Chan.GuildID, *t1vc)
	if err != nil {
		return res, err
	}
	res.Team_1_vc = t1vc_CHAN

	t2vc := GetChannelConfig(dgo.ChannelTypeGuildVoice, m.GetUsers(TEAM2_OPTION), mp/2, "Team 2 - "+m.MatchName)
	t2vc_CHAN, err := Bot.GuildChannelCreateComplex(m.Chan.GuildID, *t2vc)
	if err != nil {
		return res, err
	}
	res.Team_2_vc = t2vc_CHAN

	// Viewer vc
	viewerUsers := m.GetUsers(VIEWERS_OPTION)
	if len(viewerUsers) != 0 {
		vwvc := GetChannelConfig(dgo.ChannelTypeGuildVoice, viewerUsers, len(viewerUsers), "Viewers - "+m.MatchName)
		vwvc_CHAN, err := Bot.GuildChannelCreateComplex(m.Chan.GuildID, *vwvc)
		if err != nil {
			return res, err
		}
		res.Viewer_vc = vwvc_CHAN
	}

	return res, nil
}

func GetQueueMessage(queue []*dgo.User) *dgo.MessageSend {
	return &dgo.MessageSend{
		Content: getQueueMessageBody(queue),
		Components: []dgo.MessageComponent{
			dgo.ActionsRow{
				Components: []dgo.MessageComponent{
					getButton("join q", JOIN_Q, dgo.PrimaryButton),
					getButton("leave q", LEAVE_Q, dgo.DangerButton),
				},
			},
		},
	}
}

func getQueueMessageBody(queue []*dgo.User) (content string) {
	content = "```md\n"
	content += "# Welcome! #\n"
	content += "Click the button below to join the queue\n"
	content += "Current queue:\n"
	content += " - " + strconv.Itoa(len(queue)) + "/" + strconv.Itoa(mp)
	content += "\n```"
	return
}
