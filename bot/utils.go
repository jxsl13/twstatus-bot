package bot

import (
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

func optionalChannelID(data cmdroute.CommandData) (channelId discord.ChannelID) {
	channelId = data.Event.ChannelID
	if s, _ := data.Options.Find(channelOptionName).SnowflakeValue(); s != 0 {
		channelId = discord.ChannelID(s)
	}
	return channelId
}
