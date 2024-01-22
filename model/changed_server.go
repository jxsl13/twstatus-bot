package model

import (
	"fmt"
)

type ChangedServerStatus struct {
	Target MessageTarget

	Prev    ServerStatus
	Curr    ServerStatus
	Offline bool
}

func (c *ChangedServerStatus) Content() string {
	if c.Offline {
		return fmt.Sprintf("%s [OFFLINE]", c.Prev.Name)
	}

	header := c.Curr.Header()
	return header
}
