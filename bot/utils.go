package bot

import (
	"errors"
	"net/http"
	"time"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

type Backoff struct {
	Backoff time.Duration
	Until   time.Time
}

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

	return herr.Status == http.StatusNotFound
}

func ErrIsAccessDenied(err error) bool {
	var herr *httputil.HTTPError
	if !errors.As(err, &herr) {
		return false
	}

	return herr.Status == http.StatusForbidden
}

// closeTimer should be used as a deferred function
// in order to cleanly shut down a timer
func closeTimer(timer *time.Timer, drained *bool) {
	if drained == nil {
		panic("drained bool pointer is nil")
	}
	if !timer.Stop() {
		if *drained {
			return
		}
		<-timer.C
		*drained = true
	}
}

// resetTimer sets drained to false after resetting the timer.
func resetTimer(timer *time.Timer, duration time.Duration, drained *bool) {
	if drained == nil {
		panic("drained bool pointer is nil")
	}
	if !timer.Stop() {
		if !*drained {
			<-timer.C
		}
	}
	timer.Reset(duration)
	*drained = false
}
