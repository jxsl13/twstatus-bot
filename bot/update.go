package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/servers"
)

func (b *Bot) updateServers(ctx context.Context) (src, dst int, err error) {
	start := time.Now()
	servers, err := servers.GetAllServers()
	if err != nil {
		return 0, 0, err
	}

	serverList, err := model.NewServersFromDTO(servers)
	if err != nil {
		return 0, 0, err
	}

	b.db.Lock()
	defer b.db.Unlock()

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

	src = len(servers)
	dst = len(serverList)

	log.Printf("updated %d source to %d target servers in %s", src, dst, time.Since(start))
	return src, dst, nil
}

func (b *Bot) activeServers(ctx context.Context) (m map[model.Target]model.ServerStatus, err error) {
	b.db.Lock()
	defer b.db.Unlock()

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = closer(err)
	}()

	servers, err := dao.ActiveServers(ctx, tx)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (b *Bot) updateDiscordMessages(ctx context.Context) (int, error) {
	servers, err := b.activeServers(ctx)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	l := len(servers)

	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(l)
	for target, server := range servers {
		go func(target model.Target, status model.ServerStatus) {
			defer wg.Done()
			err := b.updateDiscordMessage(target, status)
			if err != nil {
				log.Printf("failed to update discord message for %v: %v", target, err)
			}
		}(target, server)
	}
	wg.Wait()

	log.Printf("updated %d discord messages in %s", l, time.Since(start))
	return l, err
}

func (b *Bot) updateDiscordMessage(target model.Target, status model.ServerStatus) error {

	content := status.String()
	_, err := b.state.EditMessage(
		target.ChannelID,
		target.MessageID,
		content,
	)
	if err == nil {
		return nil
	}

	var herr *httputil.HTTPError
	if !errors.As(err, &herr) {
		return err
	}

	if herr.Status != http.StatusNotFound {
		return herr
	}

	// message was somehow deleted without us noticing
	// remove tracking for that message
	b.db.Lock()
	defer b.db.Unlock()

	err = dao.RemoveTrackingByMessageID(b.ctx, b.db, target.GuildID, target.MessageID)
	if err != nil {
		return fmt.Errorf("failed to remove tracking of message id: %s: %w", target.MessageID, err)
	}

	log.Printf("removed tracking for guild %d and message id: %s (reason: 'message not found')", target.GuildID, target.MessageID)
	return nil
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
