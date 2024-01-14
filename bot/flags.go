package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
)

func (b *Bot) listFlags(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {

	b.db.Lock()
	flags, err := dao.ListFlags(ctx, b.queries)
	b.db.Unlock()

	if err != nil {
		return errorResponse(err)
	}

	var sb strings.Builder
	sb.Grow(len(flags) * 16)

	const maxCharactersPerLine = 64
	lineStart := 0

	for idx, f := range flags {
		if idx > 0 {
			sb.WriteString(" ")
		}

		tag := fmt.Sprintf("`%s`", f.Abbr)
		if sb.Len()+len(tag)-lineStart > maxCharactersPerLine {
			sb.WriteString("\n")
			lineStart = sb.Len()
		}

		sb.WriteString(tag)
	}

	content := sb.String()

	return &api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Flags:   discord.EphemeralMessage,
	}
}
