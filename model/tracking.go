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

func (t *Target) String() string {
	// https://discord.com/channels/628902095747285012/718814596323868766/1190423006590279791
	return fmt.Sprintf("https://discord.com/channels/%d/%d/%d", t.GuildID, t.ChannelID, t.MessageID)
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

func (ss ServerStatus) Header() string {
	header := fmt.Sprintf("[%s](https://ddnet.org/connect-to/?addr=%s) (%d/%d)",
		ss.Name,
		ss.Address,
		len(ss.Clients),
		ss.MaxPlayers,
	)
	header = markdown.WrapInFat(header)
	return header
}

func (ss ServerStatus) ToDiscordMessage() (content string, embeds []discord.Embed) {
	header := ss.Header()
	embeds = ss.Clients.ToEmbeds(ss.ScoreKind)

	return header, embeds
}

func (ss ServerStatus) String() string {
	var sb strings.Builder

	header := ss.Header()
	clients := ss.Clients.Format(ss.ScoreKind)
	sb.WriteString(header)
	sb.WriteString("\n")
	sb.WriteString(clients)

	return sb.String()
}

type ClientStatusList []ClientStatus

func (clients ClientStatusList) LongestNames() (maxNameLen, maxClanLen int) {
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
	return longestName, longestClan
}

func (clients ClientStatusList) ToEmbeds(scoreKind string) []discord.Embed {
	const (
		maxCharacters     = 6000 - 128
		maxFieldsPerEmbed = 25
		maxEmbeds         = 10
	)
	var (
		embeds               = make([]discord.Embed, 0, len(clients))
		embed  discord.Embed = discord.Embed{
			Type: discord.NormalEmbed,
		}
		longestName, longestClan = clients.LongestNames()
		nameFmtStr               = fmt.Sprintf("%%-%ds", longestName)
		clanFmtStr               = fmt.Sprintf("%%-%ds", longestClan)
		characterCnt             = 0
	)
	clients.Iterate(scoreKind, func(i int, client ClientStatus) bool {
		fields, charLen := client.ToEmbedFields(nameFmtStr, clanFmtStr, scoreKind)

		if len(embed.Fields)+len(fields) > maxFieldsPerEmbed {
			embeds = append(embeds, embed)
			embed = discord.Embed{
				Type: discord.NormalEmbed,
			}
		}

		// discord character limit
		if characterCnt+charLen > maxCharacters {
			embed.Fields = append(embed.Fields, discord.EmbedField{
				Value:  fmt.Sprintf("... and %d more", len(clients)-i),
				Inline: false,
			})

			return false
		} else {
			embed.Fields = append(embed.Fields, fields...)
			characterCnt += charLen
		}
		return true
	})

	if len(embed.Fields) > 0 {
		embeds = append(embeds, embed)
	}

	return embeds
}

// The index that is passed to the must not be assumed to be the current position in the list.
// Depending on the scoreKind, the iteration might happend in reverse while the index is still increasing.
func (clients ClientStatusList) Iterate(scoreKind string, f func(idx int, client ClientStatus) bool) {
	var continueIterating bool
	if scoreKind == "time" {
		// from smallest to biggest
		start := len(clients) - 1
		for i := start; i >= 0; i-- {
			continueIterating = f(i-start, clients[i])
			if !continueIterating {
				return
			}
		}
	} else {
		// normal score points ordered from biggest to smallest
		for i := 0; i < len(clients); i++ {
			continueIterating = f(i, clients[i])
			if !continueIterating {
				return
			}
		}
	}
}

func (clients ClientStatusList) Format(scoreKind string) string {
	const maxCharacters = 2000 - 128

	var (
		sb                       strings.Builder
		longestName, longestClan = clients.LongestNames()
		nameFmtStr               = fmt.Sprintf("%%-%ds", longestName)
		clanFmtStr               = fmt.Sprintf("%%-%ds", longestClan)
	)
	sb.Grow((longestName + longestClan + 16) * len(clients))

	clients.Iterate(scoreKind, func(i int, client ClientStatus) bool {
		line := client.Format(nameFmtStr, clanFmtStr, scoreKind)

		// discord character limit
		if sb.Len()+len(line) > maxCharacters {
			additional := len(clients) - i
			if additional > 0 {
				sb.WriteString(fmt.Sprintf("... and %d more\n", len(clients)-i))
			}
			return false
		} else {
			sb.WriteString(line)
		}
		return true
	})

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

func (cs *ClientStatus) FormatScore(scoreKind string) string {
	if scoreKind == "time" {
		if cs.Score == 0 || cs.Score == -9999 || cs.Score == math.MaxInt { // no times or invalid times
			return ""
		} else {
			return (time.Second * time.Duration(cs.Score)).String()
		}
	} else {
		return strconv.Itoa(cs.Score)
	}
}

func (cs *ClientStatus) FormatName(nameFormat string) string {
	return markdown.WrapInInlineCodeBlock(fmt.Sprintf(nameFormat, cs.Name))
}

func (cs *ClientStatus) FormatClan(clanFormat string) string {
	return markdown.WrapInInlineCodeBlock(fmt.Sprintf(clanFormat, cs.Clan))
}

func (cs *ClientStatus) Format(nameFormat, clanFormat, scoreKind string) string {
	var (
		name  = cs.FormatName(nameFormat)
		clan  = cs.FormatClan(clanFormat)
		score = cs.FormatScore(scoreKind)
	)

	if score != "" {
		score = "(" + score + ")"
	}

	robot := ""
	if !cs.IsPlayer {
		if score == "" {
			robot = ":robot:"
		} else {
			robot = " :robot:" // append a robot emoji behind robots
		}
	}

	return fmt.Sprintf("%s %s %s %s%s\n", cs.FlagEmoji, name, clan, score, robot)
}

func (cs *ClientStatus) ToEmbedFields(nameFormat, clanFormat, scoreKind string) (fields []discord.EmbedField, charLen int) {

	line := cs.Format(nameFormat, clanFormat, scoreKind)
	return []discord.EmbedField{
		{
			Value:  line,
			Inline: false,
		},
	}, len(line)
}
