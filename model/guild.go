package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Guild struct {
	ID          discord.GuildID
	Description string
}

func (g Guild) String() string {
	return fmt.Sprintf("%d: %s", g.ID, g.Description)
}

type Guilds []Guild

func (g Guilds) String() string {
	var sb strings.Builder
	for _, guild := range g {
		sb.WriteString(guild.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
