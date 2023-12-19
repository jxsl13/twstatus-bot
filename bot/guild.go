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

type AddGuildOpts struct { // optional (taken from current guild)
	Description string `discord:"description"` // required
}

func (b *Bot) addGuild(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	var id discord.GuildID
	s, err := data.Options.Find("id").SnowflakeValue()
	if err == nil && s != 0 {
		id = discord.GuildID(s)
	} else {
		id = data.Event.GuildID
	}

	var opts AddGuildOpts
	err = data.Options.Unmarshal(&opts)
	if err != nil {
		return errorResponse(err)
	}

	err = dao.AddGuild(ctx, b.db, model.Guild{
		ID:          id,
		Description: opts.Description,
	})

	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Added guild %d (%s)", data.Event.GuildID, opts.Description)),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) listGuilds(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	guilds, err := dao.ListGuilds(ctx, b.db)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Guilds: \n%s\n", guilds)),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) removeGuild(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	var id discord.GuildID
	s, err := data.Options.Find("id").SnowflakeValue()
	if err == nil && s != 0 {
		id = discord.GuildID(s)
	} else {
		id = data.Event.GuildID
	}

	guild, err := dao.RemoveGuild(ctx, b.db, id)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Removed guild %d (%s)", guild.ID, guild.Description)),
		Flags:   discord.EphemeralMessage,
	}
}
