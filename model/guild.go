package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

type Guild struct {
	ID          discord.GuildID
	Description string
}

// to sqlc
func (g *Guild) ToSQLC() sqlc.AddGuildParams {
	return sqlc.AddGuildParams{
		GuildID:     int64(g.ID),
		Description: g.Description,
	}
}

func (g *Guild) String() string {
	return fmt.Sprintf("`%d`: %s", g.ID, g.Description)
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
