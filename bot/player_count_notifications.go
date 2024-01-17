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

	userTarget := model.MessageUserTarget{
		UserID: e.UserID,
		MessageTarget: model.MessageTarget{
			ChannelTarget: model.ChannelTarget{
				GuildID:   e.GuildID,
				ChannelID: e.ChannelID,
			},
			MessageID: e.MessageID,
		},
	}

	n := model.PlayerCountNotification{
		MessageUserTarget: userTarget,
		Threshold:         val,
	}

	q, closer, err := b.TxQueries(b.ctx)
	if err != nil {
		log.Printf("failed to get transaction queries for player count notification: %v", err)
		return
	}
	defer func() {
		err = closer(err)
		if err != nil {
			log.Printf("failed to close transaction queries for player count notification: %v", err)
		}
	}()

	pcn, err := dao.GetPlayerCountNotification(b.ctx, q, userTarget)
	if err != nil {
		// not found, just insert
		if errors.Is(err, dao.ErrNotFound) {
			err = dao.SetPlayerCountNotification(b.ctx, q, n)
			if err != nil {
				log.Printf("failed to set player count notification(%s -> %s): %v", n.MessageTarget, n.UserID, err)
				return
			}
			log.Printf("added %d player count notification for user %s and message %s", n.Threshold, e.UserID, n.MessageTarget)
			return
		} else {
			log.Printf("failed to get player count notification(%s -> %s): %v", n.MessageTarget, n.UserID, err)
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

	err = dao.SetPlayerCountNotification(b.ctx, q, n)
	if err != nil {
		log.Printf("failed to set player count notification(%s -> %s): %v", n.MessageTarget, n.UserID, err)
		return
	}

	log.Printf("added player count notification for user %s and message %s", e.UserID, n.MessageTarget)
}

func (b *Bot) handleRemovePlayerCountNotifications(e *gateway.MessageReactionRemoveEvent) {
	val, found := reactionPlayerCountNotificationMap[e.Emoji.APIString()]
	if !found || b.userID == e.UserID {
		return
	}

	userTarget := model.MessageUserTarget{
		UserID: e.UserID,
		MessageTarget: model.MessageTarget{
			ChannelTarget: model.ChannelTarget{
				GuildID:   e.GuildID,
				ChannelID: e.ChannelID,
			},
			MessageID: e.MessageID,
		},
	}

	n := model.PlayerCountNotification{
		MessageUserTarget: userTarget,
		Threshold:         val,
	}

	q, closer, err := b.ConnQueries(b.ctx)
	if err != nil {
		log.Printf("failed to get connection queries for player count notification: %v", err)
		return
	}
	defer closer()

	err = dao.RemovePlayerCountNotification(b.ctx, q, n)
	if err != nil {
		log.Printf("failed to remove player count notification(%s -> %s): %v", n.MessageTarget, n.UserID, err)
		return
	}
	log.Printf("removed %d player count notification for user %s and message %s", val, e.UserID, n.MessageTarget)
}
