package model

import "github.com/diamondburned/arikawa/v3/discord"

type MessageMentions map[MessageTarget]Mentions

type Mentions []discord.UserID

func (m Mentions) Equals(other Mentions) bool {
	if len(m) != len(other) {
		return false
	}

	for i := range m {
		if m[i] != other[i] {
			return false
		}
	}

	return true
}
