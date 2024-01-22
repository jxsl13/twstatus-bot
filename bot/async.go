package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

// start this asynchronously once
func (b *Bot) serverUpdater(duration time.Duration) {
	var (
		timer   = time.NewTimer(0)
		drained = false
	)
	defer closeTimer(timer, &drained)
	for {
		select {
		case <-timer.C:
			drained = true
			// do something
			resetTimer(timer, duration, &drained)
			func() {
				_, _, err := b.updateServers()
				if err != nil {
					b.l.Errorf("failed to update servers: %v", err)
					return
				}

				// publish changed servers
				err = b.changedServers()
				if err != nil {
					b.l.Errorf("failed to get changed server messages from db: %v", err)
					return
				}
			}()
		case <-b.ctx.Done():
			log.Println("closed async goroutine for server and message updates")
			return
		}
	}
}

func (b *Bot) messageUpdater(id int) {
	log.Printf("goroutine %d starting async goroutine for message updates", id)

loop:
	for {
		select {
		case <-b.ctx.Done():
			break loop
		case server, ok := <-b.c:
			if !ok {
				break loop
			}
			err := b.updateDiscordMessage(server)
			if err != nil {
				b.l.Errorf("goroutine %0d: failed to update discord message %v: %v", id, server.Target, err)
			}

		}
	}

	log.Printf("goroutine %d: closed async goroutine for message updates", id)
}

func (b *Bot) notificationUpdater(id int) {
	log.Printf("goroutine %d starting async goroutine for channel notifications", id)

loop:
	for {
		select {
		case <-b.ctx.Done():
			break loop
		case notification, ok := <-b.n:
			if !ok {
				break loop
			}
			err := b.updateChannelNotification(notification)
			if err != nil {
				b.l.Errorf("goroutine %0d: failed to update channel notification %v: %v", id, notification, err)
			}

		}
	}

	log.Printf("goroutine %d: closed async goroutine for channel notifications", id)
}

func (b *Bot) updateChannelNotification(n model.PlayerCountNotificationMessage) (err error) {
	dao, closer, err := b.TxDAO(b.ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction queries for channel notification: %w", err)
	}
	defer func() {
		err = closer(err)
	}()

	// remove previous notification message if exists
	if n.PrevMessageID != 0 {
		// check if message still exists
		err := b.state.DeleteMessage(
			n.ChannelTarget.ChannelID,
			n.PrevMessageID,
			api.AuditLogReason("removing previous channel notification message"),
		)
		if err != nil && !ErrIsNotFound(err) {
			b.l.Errorf("failed to delete previous notification message %s: %v", n.MessageTarget(n.PrevMessageID), err)
			err = nil
		}

		// cleanup database if notification was deleted by some user/admin
		err = dao.RemovePlayerCountNotificationMessage(b.ctx, n.ChannelTarget.ChannelID, n.PrevMessageID)
		if err != nil {
			return fmt.Errorf("failed to remove previous channel notification message from database: %w", err)
		}
	}

	// remove all requests from database
	for _, umt := range n.RemoveUserMessageReactions {
		err = dao.RemovePlayerCountNotificationRequest(b.ctx, model.PlayerCountNotificationRequest{
			MessageUserTarget: model.MessageUserTarget{
				UserID: umt.UserID,
				MessageTarget: model.MessageTarget{
					ChannelTarget: model.ChannelTarget{
						GuildID:   n.ChannelTarget.GuildID,
						ChannelID: n.ChannelTarget.ChannelID,
					},
					MessageID: umt.MessageID,
				},
			},
			Threshold: umt.Threshold,
		})
		if err != nil {
			return fmt.Errorf("failed to remove player count notification request from database: %w", err)
		}
	}

	// delete all reactions from specified messages
	for _, mt := range n.RemoveMessageReactions {
		err = b.state.DeleteReactions(n.ChannelID, mt.MessageID, mt.Reaction())
		if err != nil && !ErrIsNotFound(err) {
			b.l.Errorf("failed to delete reaction %s from message %s: %v", mt.Reaction(), n.MessageTarget(mt.MessageID), err)
			err = nil
		}
	}

	mentionUsers := n.UserIDs
	if len(n.UserIDs) > 100 {
		// we do not expect more than 100 users to be mentioned anyway
		mentionUsers = n.UserIDs[:100]
	}

	// send new message
	msg, err := b.state.SendMessageComplex(n.ChannelTarget.ChannelID, api.SendMessageData{
		Content: n.Format(),
		Flags:   discord.SuppressEmbeds,
		AllowedMentions: &api.AllowedMentions{
			Users: mentionUsers,
		},
	})
	if err != nil {
		return err
	}

	// update database to contain latest notification message for the current channel
	err = dao.AddPlayerCountNotificationMessage(b.ctx, msg.ChannelID, msg.ID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) cacheCleanup(id int) {
	log.Printf("goroutine %d starting async goroutine for cache cleanup", id)
	var (
		cleanupInterval = 20 * b.pollingInterval
		timer           = time.NewTimer(cleanupInterval)
		drained         = false
	)
	defer closeTimer(timer, &drained)
	for {
		select {
		case <-timer.C:
			drained = true
			// do something
			resetTimer(timer, cleanupInterval, &drained)

			size := b.conflictMap.Size()
			if size == 0 {
				// nothing to do
				continue
			}

			now := time.Now()
			log.Printf("cache contains %d entries before cleanup at %s", size, now)
			b.conflictMap.Range(func(key model.MessageTarget, value Backoff) bool {
				// remove expired keys
				if now.After(value.Until) {
					b.conflictMap.Delete(key)
				}
				return true
			})
			log.Printf("cache contains %d entries after cleanup at %s", b.conflictMap.Size(), now)

		case <-b.ctx.Done():
			log.Printf("goroutine %d: closed async goroutine for cache cleanup", id)
			return
		}
	}
}
