package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

type Channel struct {
	GuildID discord.GuildID
	ID      discord.ChannelID
	Running bool
}

func (c *Channel) ToSQLC() sqlc.AddGuildChannelParams {
	return sqlc.AddGuildChannelParams{
		GuildID:   int64(c.GuildID),
		ChannelID: int64(c.ID),
		Running:   c.Running,
	}
}

func (c *Channel) RunningInt64() int64 {
	if c.Running {
		return 1
	}
	return 0
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
