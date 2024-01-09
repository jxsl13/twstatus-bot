package model

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/markdown"
	"github.com/mattn/go-runewidth"
)

// Tracking is a struct that represents a tracking message which contains
// a single server's status.
type Tracking struct {
	MessageTarget
	Address string // ipv4:port or [ipv6]:port
}

type Trackings []Tracking

type ByMessageTargetIDs []MessageTarget

func (a ByMessageTargetIDs) Len() int      { return len(a) }
func (a ByMessageTargetIDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMessageTargetIDs) Less(i, j int) bool {
	return a[i].Less(a[j])
}

type ChannelTarget struct {
	GuildID   discord.GuildID
	ChannelID discord.ChannelID
}

func (a ChannelTarget) Less(other ChannelTarget) bool {
	aGuildId := a.GuildID
	bGuildId := other.GuildID

	if aGuildId < bGuildId {
		return true
	}

	if aGuildId > bGuildId {
		return false
	}

	// guildIds are equal
	return a.ChannelID < other.ChannelID
}

type MessageTarget struct {
	ChannelTarget
	MessageID discord.MessageID
}

func (a MessageTarget) Less(other MessageTarget) bool {
	aGuildId := a.GuildID
	bGuildId := other.GuildID

	if aGuildId < bGuildId {
		return true
	}

	if aGuildId > bGuildId {
		return false
	}

	// guildIds are equal
	aChannelId := a.ChannelID
	bChannelId := other.ChannelID

	if aChannelId < bChannelId {
		return true
	}

	if aChannelId > bChannelId {
		return false
	}

	// channelIds are equal
	return a.MessageID < other.MessageID
}

func (t *MessageTarget) Equals(other MessageTarget) bool {
	return t.GuildID == other.GuildID && t.ChannelID == other.ChannelID && t.MessageID == other.MessageID
}

func (t MessageTarget) String() string {
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

func (ss ServerStatus) ToEmbeds() (embeds []discord.Embed) {
	embeds = ss.Clients.ToEmbeds(ss.ScoreKind)

	return embeds
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

func (clients ClientStatusList) LongestValues() (maxNameLen, maxClanLen int) {
	longestName := 0
	longestClan := 0
	width := 0
	for _, client := range clients {
		width = client.NameLen()
		if width > longestName {
			longestName = width
		}

		width = client.ClanLen()
		if width > longestClan {
			longestClan = width
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
		return []discord.Embed{}
	}
	var (
		embeds               = make([]discord.Embed, 0, len(clients))
		embed  discord.Embed = discord.Embed{
			Type: discord.NormalEmbed,
		}
		namePadding, clanPadding = clients.LongestValues()
		characterCnt             = 0
	)

	clients.Iterate(scoreKind, func(i int, client ClientStatus) bool {
		fields, charLen := client.ToEmbedFields(namePadding, clanPadding, scoreKind)

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

	list := make([]ClientStatus, len(clients))
	copy(list, clients)

	sort.Slice(list, func(i, j int) bool {
		aSpec := list[i].IsSpectator()
		bSpec := list[j].IsSpectator()

		if aSpec && !bSpec {
			return false
		} else if !aSpec && bSpec {
			return true
		}
		if scoreKind == "time" {
			// asc
			return list[i].Score < list[j].Score
		} else {
			// scoreKind == "points"
			// desc
			return list[i].Score > list[j].Score
		}
	})

	for i, client := range list {
		if !f(i, client) {
			return
		}
	}
}

func (clients ClientStatusList) Format(scoreKind string) string {
	const maxCharacters = 2000 - 128

	if len(clients) == 0 {
		return ""
	}

	var (
		sb                       strings.Builder
		namePadding, clanPadding = clients.LongestValues()
	)
	sb.Grow(min((64)*len(clients), maxCharacters))

	clients.Iterate(scoreKind, func(i int, client ClientStatus) bool {
		line := client.Format(namePadding, clanPadding, scoreKind)

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
	const spec = "👁️"
	if scoreKind == "time" {
		if cs.IsSpectator() {
			return spec
		} else if cs.Score == 0 { // no times or invalid times
			return ""
		}
		return (time.Second * time.Duration(cs.Score)).String()
	}

	if cs.IsSpectator() {
		return spec
	}

	return strconv.Itoa(cs.Score)
}

func (cs *ClientStatus) NameLen() int {
	return runewidth.StringWidth(cs.Name)
}

func (cs *ClientStatus) ClanLen() int {
	return runewidth.StringWidth(cs.Clan)
}

func (cs *ClientStatus) FormatName(padding int) string {
	return markdown.WrapInInlineCodeBlock(runewidth.FillRight(cs.Name, padding))
}

func (cs *ClientStatus) FormatClan(padding int) string {
	return markdown.WrapInInlineCodeBlock(runewidth.FillRight(cs.Clan, padding))
}

func (cs *ClientStatus) Format(namePadding, clanPadding int, scoreKind string) string {
	var (
		name  = cs.FormatName(namePadding)
		clan  = cs.FormatClan(clanPadding)
		score = cs.FormatScore(scoreKind)
	)
	// len(flag) == 4
	return fmt.Sprintf("%s %s %s %s", cs.FlagEmoji, name, clan, score)
}

func (cs *ClientStatus) ToEmbedFields(namePadding, clanPadding int, scoreKind string) (fields []discord.EmbedField, charLen int) {

	line := cs.Format(namePadding, clanPadding, scoreKind)
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
