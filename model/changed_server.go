package model

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type ChangedServerStatus struct {
	Prev              ServerStatus
	Curr              ServerStatus
	UserNotifications []discord.UserID
}

func (c *ChangedServerStatus) Content() string {
	header := c.Curr.Header()

	if len(c.UserNotifications) == 0 {
		return header
	}

	const limit = 2000
	sb := strings.Builder{}
	sb.Grow(2000)
	sb.WriteString(header)
	sb.WriteString("\n")

	for _, user := range c.UserNotifications {
		mention := user.Mention()
		if sb.Len()+len(mention) > limit {
			break
		}
		sb.WriteString(mention)
	}
	return sb.String()
}
