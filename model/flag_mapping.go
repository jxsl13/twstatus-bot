package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type FlagMapping struct {
	GuildID   discord.GuildID
	ChannelID discord.ChannelID
	FlagID    int64
	Abbr      string
	Emoji     string
}

func (f *FlagMapping) String() string {
	return fmt.Sprintf("`%s`: %s", f.Abbr, f.Emoji)
}

type FlagMappings []FlagMapping

func (f FlagMappings) String() string {
	if len(f) == 0 {
		return "no flag mappings"
	}
	var sb strings.Builder
	sb.Grow(len(f) * 16)
	for _, m := range f {
		sb.WriteString(m.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
