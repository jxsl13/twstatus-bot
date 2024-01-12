package model

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

type ChangedServerStatus struct {
	Target MessageTarget

	Prev              ServerStatus
	Curr              ServerStatus
	Offline           bool
	UserNotifications []discord.UserID
}

func (c *ChangedServerStatus) Content() string {
	if c.Offline {
		return fmt.Sprintf("%s [OFFLINE]", c.Prev.Name)
	}

	header := c.Curr.Header()

	if len(c.UserNotifications) == 0 {
		return header
	}

	const limit = 2000
	sb := strings.Builder{}
	sb.Grow(limit)
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
