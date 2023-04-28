package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

/* General use functions that create values or datastructures to be used throughout the bot */

func MapUserMentions(u []*dgo.User) (res string) {
	for i, v := range u {
		res += v.Username
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
			MapUserMentions(lm.Players) +
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
	res = append(res, &dgo.MessageEmbedField{
		Name:  "Last selection:",
		Value: "```md\n- if i cbf\n```",
	})
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

func getCaptainIds(maxPlayers int) (int, int) {
	rand.Seed(time.Now().UnixNano())

	n1 := rand.Intn(maxPlayers)
	n2 := rand.Intn(maxPlayers)
	if n1 == n2 {
		return getCaptainIds(maxPlayers)
	}
	return n1, n2
}

func getButton(l, id string, s dgo.ButtonStyle, d bool) *dgo.Button {
	return &dgo.Button{
		Label:    l,
		Style:    s,
		CustomID: id,
		Disabled: d,
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

func followUpAndError(s *discordgo.Session, i *discordgo.Interaction, errStr string) {
	s.FollowupMessageCreate(i, true, &dgo.WebhookParams{
		Content: fmt.Sprintf("``Sorry, it looks like there's been an error : %v``", errStr),
	})
}
func followUpMessage(s *discordgo.Session, i *discordgo.Interaction, msg string) {
	s.FollowupMessageCreate(i, true, &dgo.WebhookParams{
		Content: fmt.Sprintf("``%v``", msg),
	})
}

/*
TODO:
func GetChannelConfig(ctype dgo.ChannelType, GldId string, ValidUsers []*dgo.User) *ChannelConfig {
	switch ctype {
	case dgo.ChannelTypeGuildText:
		return &ChannelConfig{
			GldID: GldId,
			DgoCfg: dgo.GuildChannelCreateData{
				Type:                 ctype,
				PermissionOverwrites: MapUserPerms(ValidUsers, dgo.PermissionViewChannel),
				// TODO: make match names
				Name: faker.Lorem().Word(),
			}}

	case dgo.ChannelTypeGuildVoice:
		return &ChannelConfig{
			GldID: GldId,
			DgoCfg: dgo.GuildChannelCreateData{
				Type:                 ctype,
				PermissionOverwrites: MapUserPerms(ValidUsers, dgo.PermissionViewChannel),
				Name:                 faker.Lorem().Word(),
				UserLimit:            *MaxPlayers / 2,
			},
		}
	default:
		log.Error("[UNHANDLED CHANNEL TYPE]: GetChannelConfig")
		return nil
	}
} */

func MapUsersToPickOptions(u []*dgo.User) (res []dgo.SelectMenuOption) {
	for _, player := range u {
		res = append(res, dgo.SelectMenuOption{
			Label: player.Username, // + role + etc...
			// TODO: the random int is for self joining test
			// breaks GetUser function in HandleSelectPlayer
			Value: player.ID + strconv.Itoa(rand.Intn(10000)),
		})
	}
	for _, p := range res {
		log.Info("[REMAINING PICKS]: ", p.Label)
	}
	return
}
