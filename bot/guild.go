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
	d "github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

func (b *Bot) listGuilds(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.IsSuperAdmin(data) {
		return ErrAccessForbidden()
	}

	dao, closer, err := b.ConnDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	guilds, err := dao.ListGuilds(ctx)
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

	dao, closer, err := b.ConnDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	err = dao.AddGuild(ctx, model.Guild{
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

	guild, err := dao.RemoveGuild(ctx, id)
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
	dao, closer, err := b.ConnDAO(b.ctx)
	if err != nil {
		b.l.Errorf("failed to create transaction for addition of guild %s: %v", e.ID, err)
		return
	}
	defer closer()

	err = dao.AddGuild(b.ctx, model.Guild{
		ID:          e.ID,
		Description: e.Name,
	})
	if err != nil && !errors.Is(err, d.ErrAlreadyExists) {
		b.l.Errorf("failed to add guild %d (%s): %v", e.ID, e.Name, err)
	} else if errors.Is(err, d.ErrAlreadyExists) {
		log.Printf("guild %d (%s) already exists", e.ID, e.Name)
	} else {
		log.Printf("added guild %d (%s)", e.ID, e.Name)
	}
}

func (b *Bot) handleRemoveGuild(e *gateway.GuildDeleteEvent) {
	dao, closer, err := b.TxDAO(b.ctx)
	if err != nil {
		b.l.Errorf("failed to create transaction for deletion of guild %d: %v", e.ID, err)
		return
	}
	defer func() {
		err = closer(err)
		if err != nil {
			b.l.Errorf("failed to close transaction for deletion of guild %d: %v", e.ID, err)
		}
	}()

	guild, err := dao.RemoveGuild(b.ctx, e.ID)
	if err != nil {
		b.l.Errorf("failed to remove guild %d (%s): %v", e.ID, guild.Description, err)
	} else {
		log.Printf("removed guild %d (%s)", e.ID, guild.Description)
	}
}
