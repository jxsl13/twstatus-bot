package bot

import (
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jxsl13/twstatus-bot/dao"
)

func (b *Bot) handleMessageDeletion(e *gateway.MessageDeleteEvent) {

	// delete tracking messages from db in case someone deletes any message
	err := dao.RemoveTrackingByMessageID(b.ctx, b.db, e.GuildID, e.ID)
	if err != nil {
		log.Printf("failed to remove tracking of message id: %s: %v", e.ID, err)
	}
}
