package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Channel struct {
	GuildID discord.GuildID   `discord:"guild_id"`
	ID      discord.ChannelID `discord:"id"`
	Running bool              `discord:"running"`
}

func (c Channel) String() string {
	return fmt.Sprintf("<#%d>", c.ID)
}

func (c Channel) StatusString() string {
	active := "inactive"
	if c.Running {
		active = "active"
	}
	return fmt.Sprintf("%s (%s)", c.String(), active)
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

func (c Channels) StatusString() string {
	if len(c) == 0 {
		return "no channels"
	}
	var sb strings.Builder
	for _, channel := range c {
		sb.WriteString(channel.StatusString())
		sb.WriteString("\n")
	}
	return sb.String()
}
