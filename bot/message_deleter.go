package bot

import (
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jxsl13/twstatus-bot/dao"
)

func (b *Bot) handleMessageDeletion(e *gateway.MessageDeleteEvent) {
	q, closer, err := b.ConnQueries(b.ctx)
	if err != nil {
		log.Printf("failed to get connection queries for message deletion: %v", err)
		return
	}
	defer closer()

	// delete tracking messages from db in case someone deletes any message
	err = dao.RemoveTrackingByMessageID(b.ctx, q, e.GuildID, e.ID)
	if err != nil {
		log.Printf("failed to remove tracking of guild %s and message id: %s: %v", e.GuildID, e.ID, err)
	}
}
