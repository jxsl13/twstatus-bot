package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (b *Bot) updateServers() (src, dst int, err error) {
	start := time.Now()
	servers, err := servers.GetAllServers()
	if err != nil {
		return 0, 0, err
	}
	httpGet := time.Since(start)

	start = time.Now()
	serverList, err := model.NewServersFromDTO(servers)
	if err != nil {
		return 0, 0, err
	}
	convert := time.Since(start)

	start = time.Now()
	err = func(srvs []model.Server) error {
		q, closer, err := b.TxQueries(b.ctx)
		if err != nil {
			return err
		}
		defer func() {
			err = closer(err)
		}()

		err = dao.SetServers(b.ctx, q, srvs)
		if err != nil {
			return err
		}
		return nil
	}(serverList)
	if err != nil {
		return 0, 0, err
	}
	dbSet := time.Since(start)

	src = len(servers)
	dst = len(serverList)

	dur := httpGet + convert + dbSet
	log.Printf("updated %d source to %d target servers in %s", src, dst, dur)
	if dur > b.pollingInterval {
		b.Warnf(`updating servers took longer than the polling interval (%s > %s)
http request took   %s
dto conversion took %s
db transaction took %s
`,
			dur,
			b.pollingInterval,
			httpGet,
			convert,
			dbSet,
		)
	}
	return src, dst, nil
}

func (b *Bot) changedServers() error {
	var producer chan<- model.ChangedServerStatus = b.c

	servers, err := func() (map[model.MessageTarget]model.ChangedServerStatus, error) {
		q, closer, err := b.TxQueries(b.ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = closer(err)
		}()

		return dao.ChangedServers(b.ctx, q)
	}()
	if err != nil {
		return err
	}

	log.Printf("%d server messages require an update", len(servers))
	for _, server := range servers {
		select {
		case producer <- server:
			continue
		case <-b.ctx.Done():
			return b.ctx.Err()
		}
	}
	return nil
}

func (b *Bot) updateDiscordMessage(change model.ChangedServerStatus) (err error) {
	var (
		content string
		embeds  []discord.Embed = []discord.Embed{}
		status                  = change.Curr
		target                  = change.Target
	)
	waitUntil, found := b.conflictMap.Load(target)
	expired := !found || waitUntil.Until.After(time.Now())

	if !expired {
		log.Printf("skipping update of %s, because it was updated recently", target)
		return nil
	}

	if b.useEmbeds {
		// new message format
		content = change.Content()
		embeds = status.ToEmbeds()
	} else {
		// legacy message format
		content = status.String()
	}

	data := api.EditMessageData{
		Content: option.NewNullableString(content),
		Embeds:  &embeds,
	}

	_, err = b.state.EditMessageComplex(
		target.ChannelID,
		target.MessageID,
		data,
	)
	if err == nil {
		// delete in case everything is fine and a backoff was found
		if found {
			b.conflictMap.Delete(target)
		}
		return nil
	}

	var herr *httputil.HTTPError
	if !errors.As(err, &herr) {
		return err
	}

	b.Warnf("failed to update message %s: %v", target, herr)
	editingTooFrequently := herr.Status == http.StatusTooManyRequests && herr.Code == 30046
	if editingTooFrequently {
		b.conflictMap.Compute(target, func(backoff Backoff, loaded bool) (newValue Backoff, delete bool) {
			// not found
			now := time.Now()
			if !loaded {
				// initialize
				return Backoff{
					Backoff: b.pollingInterval,
					Until:   now.Add(b.pollingInterval),
				}, false
			}

			// already exists, increase backoff
			newBackoff := backoff.Backoff * 2
			return Backoff{
				Backoff: newBackoff,
				Until:   now.Add(newBackoff),
			}, false
		})

		return herr
	}

	isNotFound := herr.Status == http.StatusNotFound
	if !isNotFound {
		// is NOT a 404, unknown error
		return herr
	}

	// message will be deleted from db, also remove from cache
	if found {
		b.conflictMap.Delete(target)
	}

	// message was somehow deleted without us noticing
	// remove tracking for that message
	q, closer, err := b.ConnQueries(b.ctx)
	if err != nil {
		return err
	}
	defer closer()

	err = dao.RemoveTrackingByMessageID(b.ctx, q, target.GuildID, target.MessageID)
	if err != nil {
		return fmt.Errorf("failed to remove tracking of message id: %s: %w", target.MessageID, err)
	}

	log.Printf("removed tracking for %s (reason: 'message not found')", target)
	return nil
}

func (b *Bot) updateServerListCommand(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	if !b.IsSuperAdmin(data) {
		return ErrAccessForbidden()
	}

	start := time.Now()
	src, dst, err := b.updateServers()
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
