package bot

import (
	"errors"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

func optionalChannelID(data cmdroute.CommandData) (channelId discord.ChannelID) {
	channelId = data.Event.ChannelID
	if s, _ := data.Options.Find(channelOptionName).SnowflakeValue(); s != 0 {
		channelId = discord.ChannelID(s)
	}
	return channelId
}

func ErrIsNotFound(err error) bool {
	var herr *httputil.HTTPError
	if !errors.As(err, &herr) {
		return false
	}

	return herr.Code == http.StatusNotFound
}
