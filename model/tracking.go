package model

import "github.com/diamondburned/arikawa/v3/discord"

// Tracking is a struct that represents a tracking message which contains
// a single server's status.
type Tracking struct {
	GuildID   discord.GuildID
	ChannelID discord.ChannelID
	Address   string // ipv4:port or [ipv6]:port
	MessageID discord.MessageID
}
