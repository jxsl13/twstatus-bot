package bot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

func (b *Bot) listGuilds(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data) {
		return ErrAccessForbidden()
	}

	q, closer, err := b.ConnQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	guilds, err := dao.ListGuilds(ctx, q)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Guilds: \n%s", guilds)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

type AddGuildOpts struct { // optional (taken from current guild)
	Description string `discord:"description"` // required
}

func (b *Bot) addGuildCommand(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data) {
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

	q, closer, err := b.ConnQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	err = dao.AddGuild(ctx, q, model.Guild{
		ID:          id,
		Description: opts.Description,
	})

	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Added guild %d (%s)", data.Event.GuildID, opts.Description)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) removeGuildCommand(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	if !b.IsSuperAdmin(data) {
		return ErrAccessForbidden()
	}

	var id discord.GuildID
	s, err := data.Options.Find("id").SnowflakeValue()
	if err == nil && s != 0 {
		id = discord.GuildID(s)
	} else {
		id = data.Event.GuildID
	}

	q, closer, err := b.TxQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	guild, err := dao.RemoveGuild(ctx, q, id)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Removed guild `%d` (%s)", guild.ID, guild.Description)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) handleAddGuild(e *gateway.GuildCreateEvent) {
	q, closer, err := b.ConnQueries(b.ctx)
	if err != nil {
		b.Errorf("failed to create transaction for addition of guild %s: %v", e.ID, err)
		return
	}
	defer closer()

	err = dao.AddGuild(b.ctx, q, model.Guild{
		ID:          e.ID,
		Description: e.Name,
	})
	if err != nil && !errors.Is(err, dao.ErrAlreadyExists) {
		b.Errorf("failed to add guild %d (%s): %v", e.ID, e.Name, err)
	} else if errors.Is(err, dao.ErrAlreadyExists) {
		log.Printf("guild %d (%s) already exists", e.ID, e.Name)
	} else {
		log.Printf("added guild %d (%s)", e.ID, e.Name)
	}
}

func (b *Bot) handleRemoveGuild(e *gateway.GuildDeleteEvent) {
	q, closer, err := b.TxQueries(b.ctx)
	if err != nil {
		b.Errorf("failed to create transaction for deletion of guild %d: %v", e.ID, err)
		return
	}
	defer func() {
		err = closer(err)
		if err != nil {
			b.Errorf("failed to close transaction for deletion of guild %d: %v", e.ID, err)
		}
	}()

	guild, err := dao.RemoveGuild(b.ctx, q, e.ID)
	if err != nil {
		b.Errorf("failed to remove guild %d (%s): %v", e.ID, guild.Description, err)
	} else {
		log.Printf("removed guild %d (%s)", e.ID, guild.Description)
	}
}
