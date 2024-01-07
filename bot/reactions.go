package bot

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (b *Bot) handleAddReactions(e *gateway.MessageReactionAddEvent) {
	// check if the reaction was added to a message that is tracked
	// and if the user that added the reaction is the bot itself
	// if so, add the reaction to the message
	// if not, ignore the reaction

	b.handleAddPlayerCountNotifications(e)
}

func (b *Bot) handleRemoveReactions(e *gateway.MessageReactionRemoveEvent) {

	// check if the reaction was removed from a message that is tracked
	// and if the user that removed the reaction is the bot itself
	// if so, remove the reaction from the message
	// if not, ignore the reaction
	b.handleRemovePlayerCountNotifications(e)
}
