package bot

import (
	"errors"
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

func (b *Bot) handleAddPlayerCountNotifications(e *gateway.MessageReactionAddEvent) {
	val, found := reactionPlayerCountNotificationMap[e.Emoji.APIString()]
	if !found {
		return
	}

	userTarget := model.UserTarget{
		UserID: e.UserID,
		Target: model.Target{
			GuildID:   e.GuildID,
			ChannelID: e.ChannelID,
			MessageID: e.MessageID,
		},
	}

	n := model.PlayerCountNotification{
		UserTarget: userTarget,
		Threshold:  val,
	}

	b.db.Lock()
	defer b.db.Unlock()

	tx, closer, err := b.Tx(b.ctx)
	if err != nil {
		log.Printf("failed to get transaction: %v", err)
		return
	}
	defer func() {
		err = closer(err)
		if err != nil {
			log.Printf("failed to close transaction: %v", err)
		}
	}()

	pcn, err := dao.GetPlayerCountNotification(b.ctx, tx, userTarget)
	if err != nil {
		// not found, just insert
		if errors.Is(err, dao.ErrNotFound) {
			err = dao.SetPlayerCountNotification(b.ctx, tx, n)
			if err != nil {
				log.Printf("failed to set player count notification(%s -> %s): %v", n.Target, n.UserID, err)
				return
			}
			log.Printf("added %d player count notification for user %s and message %s", n.Threshold, e.UserID, n.Target)
			return
		} else {
			log.Printf("failed to get player count notification(%s -> %s): %v", n.Target, n.UserID, err)
			return
		}
	}

	// already exists, update
	if pcn.Threshold == n.Threshold {
		return
	}

	prevEmoji, ok := reactionPlayerCountNotificationReverseMap[pcn.Threshold]
	if !ok {
		panic("failed to get emoji for player count notification: map must contain value to emoji mapping")
	}

	err = b.state.DeleteUserReaction(e.ChannelID, e.MessageID, e.UserID, prevEmoji)
	if err != nil {
		log.Printf("failed to delete previous reaction: %v", err)
		return
	}

	err = dao.SetPlayerCountNotification(b.ctx, tx, n)
	if err != nil {
		log.Printf("failed to set player count notification(%s -> %s): %v", n.Target, n.UserID, err)
		return
	}

	log.Printf("added player count notification for user %s and message %s", e.UserID, n.Target)
}

func (b *Bot) handleRemovePlayerCountNotifications(e *gateway.MessageReactionRemoveEvent) {
	val, found := reactionPlayerCountNotificationMap[e.Emoji.APIString()]
	if !found || b.userId == e.UserID {
		return
	}

	userTarget := model.UserTarget{
		UserID: e.UserID,
		Target: model.Target{
			GuildID:   e.GuildID,
			ChannelID: e.ChannelID,
			MessageID: e.MessageID,
		},
	}

	n := model.PlayerCountNotification{
		UserTarget: userTarget,
		Threshold:  val,
	}

	b.db.Lock()
	defer b.db.Unlock()

	err := dao.RemovePlayerCountNotification(b.ctx, b.db, n)
	if err != nil {
		log.Printf("failed to remove player count notification(%s -> %s): %v", n.Target, n.UserID, err)
		return
	}
	log.Printf("removed %d player count notification for user %s and message %s", val, e.UserID, n.Target)
}
