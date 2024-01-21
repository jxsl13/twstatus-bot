package bot

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (b *Bot) handleMessageDeletion(e *gateway.MessageDeleteEvent) {
	dao, closer, err := b.ConnDAO(b.ctx)
	if err != nil {
		b.l.Errorf("failed to get connection queries for message deletion: %v", err)
		return
	}
	defer closer()

	// delete tracking messages from db in case someone deletes any message
	err = dao.RemoveTrackingByMessageID(b.ctx, e.GuildID, e.ID)
	if err != nil {
		b.l.Errorf("failed to remove tracking of guild %s and message id: %s: %v", e.GuildID, e.ID, err)
	}
}
