package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/model"
)

type AddFlagMappingParams struct {
	Abbr  string `discord:"abbr"`
	Emoji string `discord:"emoji"`
}

func (b *Bot) listFlagMappings(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	dao, closer, err := b.ConnDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	channelID := optionalChannelID(data)
	mappings, err := dao.ListFlagMappings(
		ctx,
		data.Event.GuildID,
		channelID,
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

	dao, closer, err := b.TxDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	flag, err := dao.GetFlagByAbbr(ctx, params.Abbr)
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

	err = dao.AddFlagMapping(ctx, mapping)
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

	dao, closer, err := b.TxDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	err = dao.RemoveFlagMapping(
		ctx,
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
