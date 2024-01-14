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

	b.db.Lock()
	defer b.db.Unlock()

	guilds, err := dao.ListGuilds(ctx, b.queries)
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

	b.db.Lock()
	defer b.db.Unlock()

	err = dao.AddGuild(ctx, b.queries, model.Guild{
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

	guild, err := dao.RemoveGuild(ctx, queries, id)
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
	b.db.Lock()
	defer b.db.Unlock()

	err := dao.AddGuild(b.ctx, b.queries, model.Guild{
		ID:          e.ID,
		Description: e.Name,
	})
	if err != nil && !errors.Is(err, dao.ErrAlreadyExists) {
		log.Printf("failed to add guild %d (%s): %v", e.ID, e.Name, err)
	} else if errors.Is(err, dao.ErrAlreadyExists) {
		log.Printf("guild %d (%s) already exists", e.ID, e.Name)
	} else {
		log.Printf("added guild %d (%s)", e.ID, e.Name)
	}
}

func (b *Bot) handleRemoveGuild(e *gateway.GuildDeleteEvent) {
	b.db.Lock()
	defer b.db.Unlock()

	tx, closer, err := b.Tx(b.ctx)
	if err != nil {
		log.Printf("failed to create transaction for deletion of guild %d: %v", e.ID, err)
		return
	}
	defer func() {
		err = closer(err)
		if err != nil {
			log.Printf("failed to close transactionfor deletion of guild %d: %v", e.ID, err)
		}
	}()

	queries := b.queries.WithTx(tx)

	guild, err := dao.RemoveGuild(b.ctx, queries, e.ID)
	if err != nil {
		log.Printf("failed to remove guild %d (%s): %v", e.ID, guild.Description, err)
	} else {
		log.Printf("removed guild %d (%s)", e.ID, guild.Description)
	}
}
