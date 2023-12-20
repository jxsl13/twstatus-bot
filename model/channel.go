package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Channel struct {
	GuildID discord.GuildID   `discord:"guild_id"`
	ID      discord.ChannelID `discord:"id"`
	Running int               `discord:"running"`
}

func (c Channel) String() string {
	active := "inactive"
	if c.Running != 0 {
		active = "active"
	}
	return fmt.Sprintf("<#%d> (%s)", c.ID, active)
}

type Channels []Channel

func (c Channels) String() string {
	if len(c) == 0 {
		return "no channels"
	}
	var sb strings.Builder
	for _, channel := range c {
		sb.WriteString(channel.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
