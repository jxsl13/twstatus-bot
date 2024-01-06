package bot

import (
	"context"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var (
	helpLines = []string{
		"*Usage:*",
		"This bot requires you to initially choose a channel to post the status updates to.",
		"You can do this by using the `/add-channel` command.",
		"Afterwards you have to add tracking of individual servers to the specified channel.",
		"You can do this by using the `/add-tracking address:<ipv4:port or [ipv6]:port>` command.",
		"Lastly, you need to start the bot for the specified channel by using the `/start` command.",
		"In case that you want to stop the bot for a specific channel, use the `/stop` command.",
		"",
		"*Commands:*",
		"`/add-channel` - adds a channel to the list of channels that are being updated",
		"`/add-tracking` - adds a server to the list of tracked servers for the specified channel",
		"If you want to remove a specific tracking, just manually delete the message that was created by the bot.",
		"`/start` - starts the bot for the specified channel",
		"`/stop` - stops the bot for the specified channel",
		"`/list-channels` - lists all channels that are registered for the currend Discord server",
		"`/list-flags` - list all flags that are available for the `/add-flag-mapping`command",
		"`/add-flag-mapping` - allows to ad a custom emoji for any player flag.",
	}
	helpText = strings.Join(helpLines, "\n")
)

func (b *Bot) help(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString(helpText),
		Flags:   discord.EphemeralMessage,
	}
}
