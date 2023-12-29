package bot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/servers"
)

func (b *Bot) updateServers(ctx context.Context) (src, dst int, err error) {
	servers, err := servers.GetAllServers()
	if err != nil {
		return 0, 0, err
	}

	serverList, err := model.NewServersFromDTO(servers)
	if err != nil {
		return 0, 0, err
	}

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		err = closer(err)
	}()

	err = dao.SetServers(ctx, tx, serverList)
	if err != nil {
		return 0, 0, err
	}
	return len(servers), len(serverList), nil
}

func (b *Bot) updateDiscordMessages(ctx context.Context) (int, error) {
	servers, err := dao.ActiveServers(ctx, b.db)
	if err != nil {
		return 0, err
	}
	l := len(servers)

	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(l)
	for k, v := range servers {
		go func(target model.Target, status model.ServerStatus) {
			defer wg.Done()
			err := b.updateDiscordMessage(target, status)
			if err != nil {
				log.Printf("failed to update discord message for %v: %v", target, err)
			}
		}(k, v)
	}
	wg.Wait()
	dur := time.Since(start)
	log.Printf("updated %d discord messages in %s", l, dur)
	return l, err
}

func (b *Bot) updateDiscordMessage(target model.Target, status model.ServerStatus) error {

	_, err := b.state.EditMessage(
		target.ChannelID,
		target.MessageID,
		status.String(),
	)
	return err
}

func (b *Bot) updateServerListCommand(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	start := time.Now()
	src, dst, err := b.updateServers(ctx)
	if err != nil {
		return errorResponse(err)
	}
	dur := time.Since(start)

	msg := fmt.Sprintf("Updated %d source to %d target servers in %s", src, dst, dur)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) updateDiscordMessagesCommand(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	if !b.IsSuperAdmin(data.Event.SenderID()) {
		return ErrAccessForbidden()
	}

	start := time.Now()
	numMsgs, err := b.updateDiscordMessages(ctx)
	if err != nil {
		return errorResponse(err)
	}
	dur := time.Since(start)

	msg := fmt.Sprintf("Updated %d discord server status messages in %s", numMsgs, dur)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
