package main

import (
	"math/rand"
	"strconv"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"syreclabs.com/go/faker"
)

// IDs for buttons/forms
const (
	JOIN_Q    = "Join"
	JOIN_Q_ID = "JQID"

	LEAVE_Q    = "Leave"
	LEAVE_Q_ID = "LQID"

	CPT_PICK = "CaptainPickChoice"
)

/* General use functions that create values or datastructures to be used throughout the bot */

func GetCaptainIds() (int, int) {
	rand.Seed(time.Now().UnixNano())

	n1 := rand.Intn(*MaxPlayers)
	n2 := rand.Intn(*MaxPlayers)
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

func MapUserPerms(users []*dgo.User, permType int64) (res []*dgo.PermissionOverwrite) {
	for _, v := range users {
		res = append(res, &dgo.PermissionOverwrite{
			ID:    v.ID,
			Type:  dgo.PermissionOverwriteTypeMember,
			Allow: int64(permType),
		})
	}
	return
}

func GetChannelConfig(ctype dgo.ChannelType, GldId string, ValidUsers []*dgo.User) *ChannelConfig {
	switch ctype {
	case dgo.ChannelTypeGuildText:
		return &ChannelConfig{
			GldID: GldId,
			DgoCfg: dgo.GuildChannelCreateData{
				Type:                 ctype,
				PermissionOverwrites: MapUserPerms(ValidUsers, dgo.PermissionViewChannel),
				Name:                 faker.Lorem().Word(),
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
}

func MapUsersToPickoptions(u []*dgo.User) (res []dgo.SelectMenuOption) {
	for _, player := range u {
		res = append(res, dgo.SelectMenuOption{
			Label: player.Username, // + role + etc...
			// TODO: the random int is for self joining test
			// breaks GetUser function in HandleSelectPlayer
			Value: player.ID + strconv.Itoa(rand.Intn(10000)),
		})
	}
	log.Info("[REMAINING PICKS]: ", res)
	return
}
