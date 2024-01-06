package model

import (
	"encoding/json"
	"fmt"
	"slices"
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
	Timestamp    time.Time // not used for equality checks
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

func (ss *ServerStatus) Equals(other ServerStatus) bool {
	return ss.Clients.Equals(other.Clients) &&
		ss.Map == other.Map &&
		ss.MaxClients == other.MaxClients &&
		ss.MaxPlayers == other.MaxPlayers &&
		ss.Gametype == other.Gametype &&
		ss.Name == other.Name &&
		ss.Passworded == other.Passworded &&
		slices.Equal(ss.Protocols, other.Protocols) &&
		ss.Version == other.Version &&
		ss.ScoreKind == other.ScoreKind &&
		ss.Address == other.Address &&
		equalPtrType(ss.MapSize, other.MapSize) &&
		equalPtrType(ss.MapSha256Sum, other.MapSha256Sum)
}

func (s *ServerStatus) ProtocolsJSON() []byte {
	data, _ := json.Marshal(s.Protocols)
	return data
}

func (s *ServerStatus) ProtocolsFromJSON(data []byte) error {
	err := json.Unmarshal(data, &s.Protocols)
	if err != nil {
		return fmt.Errorf("failed to unmarshal protocols: %w", err)
	}
	return nil
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

func (c ClientStatusList) Equals(other ClientStatusList) bool {
	if len(c) != len(other) {
		return false
	}

	for i, client := range c {
		if !client.Equals(&other[i]) {
			return false
		}
	}
	return true
}

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

	if len(clients) == 0 {
		return []discord.Embed{
			{
				Fields: []discord.EmbedField{
					{
						Value:  "No players",
						Inline: false,
					},
				},
			},
		}
	}
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
	Team      *int
	FlagAbbr  string
	FlagEmoji string // mapped emoji
}

func (cs *ClientStatus) Equals(other *ClientStatus) bool {
	return cs.Score == other.Score &&
		cs.IsPlayer == other.IsPlayer &&
		equalPtrType(cs.Team, other.Team) &&
		cs.Name == other.Name &&
		cs.Clan == other.Clan &&
		cs.Country == other.Country &&
		cs.FlagAbbr == other.FlagAbbr &&
		cs.FlagEmoji == other.FlagEmoji
}

func (c *ClientStatus) IsSpectator() bool {
	if c.Team != nil {
		return *c.Team == -1
	}

	return !c.IsPlayer && c.Score <= 0
}

func (c *ClientStatus) IsBot() bool {
	return !c.IsPlayer && !c.IsSpectator()
}

func (cs *ClientStatus) FormatScore(scoreKind string) string {
	if scoreKind == "time" {
		if cs.Score == 0 || cs.IsSpectator() { // no times or invalid times
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
	if score != "" {
		robot += " "
	}

	if cs.IsSpectator() {
		robot += ":eye:"
	} else if cs.IsBot() {
		robot += ":robot:"
	}

	return fmt.Sprintf("%s %s %s %s%s", cs.FlagEmoji, name, clan, score, robot)
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

func equalPtrType[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil || a == nil && b != nil {
		return false
	}
	return *a == *b
}
