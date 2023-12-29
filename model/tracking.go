package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/markdown"
)

// Tracking is a struct that represents a tracking message which contains
// a single server's status.
type Tracking struct {
	GuildID   discord.GuildID
	ChannelID discord.ChannelID
	Address   string // ipv4:port or [ipv6]:port
	MessageID discord.MessageID
}

type Target struct {
	GuildID   discord.GuildID
	ChannelID discord.ChannelID
	MessageID discord.MessageID
}

func (t *Target) Equals(other Target) bool {
	return t.GuildID == other.GuildID && t.ChannelID == other.ChannelID && t.MessageID == other.MessageID
}

type ServerStatus struct {
	Target Target

	Address      string
	Protocols    []string
	Name         string
	Gametype     string
	Passworded   int
	Map          string
	MapSha256Sum *string
	MapSize      *int
	Version      string
	MaxClients   int
	MaxPlayers   int
	ScoreKind    string
	Clients      ClientStatusList
}

func (ss ServerStatus) String() string {
	var sb strings.Builder

	header := fmt.Sprintf("%s (%d/%d)",
		markdown.Escape(ss.Name),
		len(ss.Clients),
		ss.MaxPlayers,
	)
	header = markdown.WrapInFat(header)

	sb.WriteString(header)
	sb.WriteString("\n")
	sb.WriteString(ss.Clients.Format(ss.ScoreKind))

	return sb.String()
}

type ClientStatusList []ClientStatus

func (clients ClientStatusList) Format(scoreKind string) string {
	const maxCharacters = 2000 - 128

	var sb strings.Builder

	longestName := 0
	longestClan := 0
	for _, client := range clients {
		if len([]rune(client.Name)) > longestName {
			longestName = len([]rune(client.Name))
		}

		if len([]rune(client.Clan)) > longestClan {
			longestName = len([]rune(client.Clan))
		}
	}

	nameFmtStr := fmt.Sprintf("%%-%ds", longestName)
	clanFmtStr := fmt.Sprintf("%%-%ds", longestClan)
	if scoreKind == "time" {
		// from smallest to biggest
		for i := len(clients) - 1; i >= 0; i-- {
			client := clients[i]
			line := client.Format(nameFmtStr, clanFmtStr, scoreKind)

			// discord character limit
			if sb.Len()+len(line) > maxCharacters {
				additional := i + 1
				if additional > 0 {
					sb.WriteString(fmt.Sprintf("... and %d more\n", additional))
				}
				break
			} else {
				sb.WriteString(line)
			}
		}
	} else {
		// normal score points ordered from biggest to smallest
		for i := 0; i < len(clients); i++ {
			client := clients[i]
			line := client.Format(nameFmtStr, clanFmtStr, scoreKind)

			// discord character limit
			if sb.Len()+len(line) > maxCharacters {
				additional := len(clients) - i
				if additional > 0 {
					sb.WriteString(fmt.Sprintf("... and %d more\n", len(clients)-i))
				}
				break
			} else {
				sb.WriteString(line)
			}
		}
	}

	return sb.String()
}

type ClientStatus struct {
	Name      string
	Clan      string
	Country   int
	Score     int
	IsPlayer  bool
	FlagAbbr  string
	FlagEmoji string // mapped emoji
}

func (cs ClientStatus) Format(nameFormat, clanFormat, scoreKind string) string {
	name := markdown.WrapInInlineCodeBlock(fmt.Sprintf(nameFormat, cs.Name))
	clan := markdown.WrapInInlineCodeBlock(fmt.Sprintf(clanFormat, cs.Clan))

	score := ""
	if scoreKind == "time" {
		if cs.Score == 0 || cs.Score == -9999 || cs.Score == math.MaxInt { // no times or invalid times
			score = ""
		} else {
			score = "(" + (time.Second * time.Duration(cs.Score)).String() + ")"
		}
	} else {
		score = "(" + strconv.Itoa(cs.Score) + ")"
	}

	return fmt.Sprintf("%s %s %s %s\n", cs.FlagEmoji, name, clan, score)
}
