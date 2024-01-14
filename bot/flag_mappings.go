package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

type AddFlagMappingParams struct {
	Abbr  string `discord:"abbr"`
	Emoji string `discord:"emoji"`
}

func (b *Bot) listFlagMappings(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.db.Lock()
	defer b.db.Unlock()

	channelID := optionalChannelID(data)
	mappings, err := dao.ListFlagMappings(ctx,
		b.queries,
		data.Event.GuildID, channelID,
	)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(mappings.String()),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) addFlagMapping(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	opts := data.Options
	var params AddFlagMappingParams
	err := opts.Unmarshal(&params)
	if err != nil {
		return errorResponse(err)
	}

	b.db.Lock()
	defer b.db.Unlock()

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	queries := b.queries.WithTx(tx)

	flag, err := dao.GetFlagByAbbr(ctx, queries, params.Abbr)
	if err != nil {
		return errorResponse(err)
	}

	mapping := model.FlagMapping{
		GuildID:   data.Event.GuildID,
		ChannelID: optionalChannelID(data),
		FlagID:    flag.ID,
		Abbr:      flag.Abbr,
		Emoji:     params.Emoji,
	}

	err = dao.AddFlagMapping(ctx, queries, mapping)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Added flag mapping: %s", mapping.String())
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

type RemoveFlagMappingParams struct {
	Abbr string `discord:"abbr"`
}

func (b *Bot) removeFlagMapping(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	opts := data.Options
	var params RemoveFlagMappingParams
	err := opts.Unmarshal(&params)
	if err != nil {
		return errorResponse(err)
	}

	b.db.Lock()
	defer b.db.Unlock()

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	queries := b.queries.WithTx(tx)

	err = dao.RemoveFlagMapping(
		ctx,
		queries,
		data.Event.GuildID,
		optionalChannelID(data),
		params.Abbr,
	)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Removed flag mapping: %s", params.Abbr)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
