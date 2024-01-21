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
	"github.com/jxsl13/twstatus-bot/sqlc"
	"github.com/jxsl13/twstatus-bot/utils"
	"github.com/mattn/go-runewidth"
)

// Tracking is a struct that represents a tracking message which contains
// a single server's status.
type Tracking struct {
	MessageTarget
	Address string // ipv4:port or [ipv6]:port
}

// to AddSQLC
func (t Tracking) ToAddSQLC() sqlc.AddTrackingParams {
	return sqlc.AddTrackingParams{
		GuildID:   int64(t.GuildID),
		ChannelID: int64(t.ChannelID),
		MessageID: int64(t.MessageID),
		Address:   t.Address,
	}
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
	Passworded   bool
	Map          string
	MapSha256Sum *string
	MapSize      *int32
	Version      string
	MaxClients   int16
	MaxPlayers   int16
	ScoreKind    string
	Clients      ClientStatusList

	// not relevant for equality checks
	// derived meta data
	Spectators    ClientStatusList
	Teams         map[int16]ClientStatusList
	LongestName   int
	LongestClan   int
	NumPlayers    int // not spectators
	NumSpectators int
}

func (ss *ServerStatus) TotalTeams() int {
	if len(ss.Spectators) > 0 {
		return len(ss.Teams) + 1
	}
	return len(ss.Teams)
}

func (ss *ServerStatus) AddClientStatus(client ClientStatus) {
	ss.Clients = append(ss.Clients, client)

	if ss.Teams == nil {
		ss.Teams = make(map[int16]ClientStatusList, 2)
	}

	if client.IsSpectator() {
		ss.Spectators = append(ss.Spectators, client)
		ss.NumSpectators++
	} else {
		teamID := client.TeamID()
		ss.Teams[teamID] = append(ss.Teams[teamID], client)
		ss.NumPlayers++
	}

	nameLen := client.NameLen()
	if nameLen > ss.LongestName {
		ss.LongestName = nameLen
	}

	clanLen := client.ClanLen()
	if clanLen > ss.LongestClan {
		ss.LongestClan = clanLen
	}
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

func (ss ServerStatus) NameToQuickJoinUrl() string {
	return fmt.Sprintf("[%s](https://ddnet.org/connect-to/?addr=%s)", ss.Name, ss.Address)
}

func (s *ServerStatus) HasV6Protocol() bool {
	for _, protocol := range s.Protocols {
		if strings.Contains(protocol, "0.6") {
			return true
		}
	}
	return false
}

func (ss ServerStatus) Header() string {

	var header string
	// 0.6.x is required for the ddnet client to join
	// via the quick join url.
	if !ss.HasV6Protocol() {
		header = ss.Name
	} else {
		header = ss.NameToQuickJoinUrl()
	}

	add := ""
	if ss.NumSpectators > 0 {
		add = "+" + strconv.Itoa(ss.NumSpectators)
	}

	header = fmt.Sprintf("%s (%d%s/%d)",
		header,
		ss.NumPlayers,
		add,
		ss.MaxPlayers,
	)
	return markdown.WrapInFat(header)
}

func (ss ServerStatus) ToEmbeds() []discord.Embed {
	if len(ss.Clients) == 0 {
		return []discord.Embed{}
	}

	const discordEmbedsLimit = 10
	totalTeams := ss.TotalTeams()
	if ss.ScoreKind == "time" || totalTeams > discordEmbedsLimit || (len(ss.Spectators) == 0 && len(ss.Teams) == 1) {
		return ss.Clients.ToEmbedList(0, ss.LongestName, ss.LongestClan, ss.ScoreKind)
	}

	// scoreKind == "points"
	embeds := make([]discord.Embed, 0, totalTeams)
	var (
		teamIDs = utils.SortedMapKeys(ss.Teams)
		color   discord.Color
		team    ClientStatusList
	)
	for _, teamID := range teamIDs {
		team = ss.Teams[teamID]
		color = teamColors[int(teamID)%maxTeamColors]

		embeds = append(embeds, team.ToEmbedList(color, ss.LongestName, ss.LongestClan, ss.ScoreKind)...)
	}

	embeds = append(embeds, ss.Spectators.ToEmbedList(0, ss.LongestName, ss.LongestClan, ss.ScoreKind)...)
	return embeds
}

func (ss ServerStatus) String() string {
	var sb strings.Builder

	header := ss.Header()
	clients := ss.Clients.Format(ss.LongestName, ss.LongestClan, ss.ScoreKind)
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

/*
Red 0xED4245 | rgb(237,66,69)
Blue 0x3498DB | rgb(52,152,219)
Green 0x57F287 | rgb(87,242,135)
Yellow 0xFEE75C | rgb(254,231,92)
Purple 0x9B59B6 | rgb(155,89,182)
Fuchsia 0xEB459E | rgb(235,69,158)
White 0xFFFFFF | rgb(255,255,255)
Aqua 0x1ABC9C | rgb(26,188,156)
LuminousVividPink 0xE91E63 | rgb(233,30,99)
Gold 0xF1C40F | rgb(241,196,15)
Orange 0xE67E22 | rgb(230,126,34)
Grey 0x95A5A6 | rgb(149,165,166)
Navy 0x34495E | rgb(52,73,94)
DarkAqua 0x11806A | rgb(17,128,106)
DarkGreen 0x1F8B4C | rgb(31,139,76)
DarkBlue 0x206694 | rgb(32,102,148)
DarkPurple 0x71368A | rgb(113,54,138)
DarkVividPink 0xAD1457 | rgb(173,20,87)
DarkGold 0xC27C0E | rgb(194,124,14)
DarkOrange 0xA84300 | rgb(168,67,0)
DarkRed 0x992D22 | rgb(153,45,34)
DarkGrey 0x979C9F | rgb(151,156,159)
DarkerGrey 0x7F8C8D | rgb(127,140,141)
LightGrey 0xBCC0C0 | rgb(188,192,192)
DarkNavy 0x2C3E50 | rgb(44,62,80)
Blurple 0x5865F2 | rgb(88,101,242)
Greyple 0x99AAb5 | rgb(153,170,181)
DarkButNotBlack 0x2C2F33 | rgb(44,47,51)
NotQuiteBlack 0x23272A | rgb(35,39,42)
*/

var teamColors = []discord.Color{
	0xED4245, // red
	0x3498DB, // blue
	0x57F287, // green
	0xFEE75C, // yellow
	0x9B59B6, // purple
	0xEB459E, // fuchsia
	0xFFFFFF, // white
	0x1ABC9C, // aqua
	0xE91E63, // luminous vivid pink
	0xF1C40F, // gold
	0xE67E22, // orange
	0x95A5A6, // grey
	0x34495E, // navy
	0x11806A, // dark aqua
	0x1F8B4C, // dark green
	0x206694, // dark blue
	0x71368A, // dark purple
	0xAD1457, // dark vivid pink
	0xC27C0E, // dark gold
	0xA84300, // dark orange
	0x992D22, // dark red
	0x979C9F, // dark grey
	0x7F8C8D, // darker grey
	0xBCC0C0, // light grey
	0x2C3E50, // dark navy
	0x5865F2, // blurple
	0x99AAb5, // greyple
	0x2C2F33, // dark but not black
	0x23272A, // not quite black
}
var maxTeamColors = len(teamColors)

func (clients ClientStatusList) ToEmbedList(color discord.Color, namePadding, clanPadding int, scoreKind string) []discord.Embed {
	const (
		maxCharacters     = 6000 - 128
		maxFieldsPerEmbed = 25
	)

	if len(clients) == 0 {
		return []discord.Embed{}
	}
	var (
		embeds               = make([]discord.Embed, 0, len(clients))
		embed  discord.Embed = discord.Embed{
			Color: color,
			Type:  discord.NormalEmbed,
		}
		characterCnt = 0
	)

	clients.Iterate(scoreKind, func(i int, client ClientStatus) bool {
		fields, charLen := client.ToEmbedFields(namePadding, clanPadding, scoreKind)

		if len(embed.Fields)+len(fields) > maxFieldsPerEmbed {
			embeds = append(embeds, embed)
			embed = discord.Embed{
				Color: color,
				Type:  discord.NormalEmbed,
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

func (clients ClientStatusList) Format(namePadding, clanPadding int, scoreKind string) string {
	const maxCharacters = 2000 - 128

	if len(clients) == 0 {
		return ""
	}

	var sb strings.Builder
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
	Country   int16
	Score     int32
	IsPlayer  bool
	Team      *int16
	FlagAbbr  string
	FlagEmoji string // mapped emoji
}

func (cs *ClientStatus) Equals(other *ClientStatus) bool {
	equalScore := cs.Score == other.Score
	equalPlayer := cs.IsPlayer == other.IsPlayer
	equalTeam := equalPtrType(cs.Team, other.Team)
	equalName := cs.Name == other.Name
	equalClan := cs.Clan == other.Clan
	equalCountry := cs.Country == other.Country
	equalFlagAbbr := cs.FlagAbbr == other.FlagAbbr
	equalFlagEmoji := cs.FlagEmoji == other.FlagEmoji
	return equalScore &&
		equalPlayer &&
		equalTeam &&
		equalName &&
		equalClan &&
		equalCountry &&
		equalFlagAbbr &&
		equalFlagEmoji

}

func (c *ClientStatus) TeamID() int16 {
	if c.IsSpectator() {
		return -1
	}
	if c.Team != nil {
		return *c.Team
	}
	return 0
}

func (c *ClientStatus) IsSpectator() bool {
	if c.Team != nil {
		return *c.Team < 0 // servers may send weird team ids
	}

	return !c.IsPlayer && c.Score < 0
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

	return strconv.FormatInt(int64(cs.Score), 10)
}

func (cs *ClientStatus) NameLen() int {
	return runewidth.StringWidth(cs.Name)
}

func (cs *ClientStatus) ClanLen() int {
	return runewidth.StringWidth(cs.Clan)
}

func padRight(s string, padding int) string {
	// handle special case where the content of one of these fields is empty
	if s == "" {
		return strings.Repeat(" ", padding)
	}
	return runewidth.FillRight(s, padding)
}

func (cs *ClientStatus) FormatName(padding int) string {
	return markdown.WrapInInlineCodeBlock(padRight(cs.Name, padding))
}

func (cs *ClientStatus) FormatClan(padding int) string {
	return markdown.WrapInInlineCodeBlock(padRight(cs.Clan, padding))
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
