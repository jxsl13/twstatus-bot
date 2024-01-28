package model

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jxsl13/twstatus-bot/servers"
	"github.com/jxsl13/twstatus-bot/sqlc"
	"github.com/jxsl13/twstatus-bot/utils"
)

type ServerList []Server

func (sl ServerList) ToSQLC(knownFlags map[int16]bool) ([]sqlc.InsertActiveServersParams, []sqlc.InsertActiveServerClientsParams) {
	servers := make([]sqlc.InsertActiveServersParams, 0, len(sl))
	clients := make([]sqlc.InsertActiveServerClientsParams, 0, len(sl))
	for _, server := range sl {
		s, cs := server.ToSQLC(knownFlags)
		servers = append(servers, s)
		clients = append(clients, cs...)
	}

	return servers, clients
}

type Server struct {
	Timestamp    time.Time
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
	Clients      ClientList // serialized as json into database
}

func (s *Server) Merge(s2 Server) {
	if s.Timestamp.IsZero() {
		*s = s2
		return
	} else if s2.Timestamp.IsZero() {
		return
	} else if s.Timestamp.Before(s2.Timestamp) {
		*s = s2
		return
	} else if s.Timestamp.After(s2.Timestamp) {
		return
	}
	// both time stamps are equal

	s.Address = utils.MergeValue(s.Address, s2.Address)
	s.Protocols = utils.MergeSliceUnique(s.Protocols, s2.Protocols)
	if len(s.Protocols) > 1 {
		sort.Strings(s.Protocols)
	}

	s.Name = utils.MergeValue(s.Name, s2.Name)
	s.Gametype = utils.MergeValue(s.Gametype, s2.Gametype)
	s.Passworded = s.Passworded || s2.Passworded
	s.Map = utils.MergeValue(s.Map, s2.Map)
	s.MapSha256Sum = utils.MergePointer(s.MapSha256Sum, s2.MapSha256Sum)
	s.MapSize = utils.MergePointer(s.MapSize, s2.MapSize)
	s.Version = utils.MergeValue(s.Version, s2.Version)
	s.MaxClients = utils.MergeValue(s.MaxClients, s2.MaxClients)
	s.MaxPlayers = utils.MergeValue(s.MaxPlayers, s2.MaxPlayers)
	s.ScoreKind = utils.MergeValue(s.ScoreKind, s2.ScoreKind)

	if len(s.Clients) < len(s2.Clients) {
		// merging is only necessary when a player connects to the server
		// which means that we have an additional player on one of the multiple protocols
		s.Clients = s.Clients[:0]
		s.Clients = append(s.Clients, s2.Clients...)
	} else if len(s.Clients) == len(s2.Clients) {
		// both servers have the same amount of clients
		// so we need to merge the clients
		for i := range s.Clients {
			s.Clients[i].Merge(s2.Clients[i])
		}
	}
	// s.Clients > s2.Clients
	// no merging necessary
}

func (s *Server) ToSQLC(knownFlags map[int16]bool) (sqlc.InsertActiveServersParams, []sqlc.InsertActiveServerClientsParams) {
	srv := sqlc.InsertActiveServersParams{
		Timestamp: pgtype.Timestamptz{
			Time:  s.Timestamp,
			Valid: true,
		},
		Address:      s.Address,
		Protocols:    s.ProtocolsJSON(),
		Name:         s.Name,
		Gametype:     s.Gametype,
		Passworded:   s.Passworded,
		Map:          s.Map,
		MapSha256sum: s.MapSha256Sum,
		MapSize:      s.MapSize,
		Version:      s.Version,
		MaxClients:   s.MaxClients,
		MaxPlayers:   s.MaxPlayers,
		ScoreKind:    s.ScoreKind,
	}

	clients := s.Clients.ToSQLC(srv.Address, knownFlags)
	return srv, clients
}

func (s *Server) ProtocolsJSON() []byte {
	data, _ := json.Marshal(s.Protocols)
	return data
}

func (s *Server) ProtocolsFromJSON(data []byte) error {
	return json.Unmarshal(data, &s.Protocols)
}

type ClientList []Client

func (cl ClientList) ToSQLC(address string, knownFlags map[int16]bool) []sqlc.InsertActiveServerClientsParams {
	result := make([]sqlc.InsertActiveServerClientsParams, 0, len(cl))
	for _, client := range cl {
		if client.IsConnecting() {
			continue
		}

		if !knownFlags[client.Country] {
			// unknown flags fall back to default flag
			client.Country = -1
		}

		result = append(result, client.ToSQLC(address))
	}
	return result
}

type Client struct {
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	Country  int16  `json:"country"`
	Score    int32  `json:"score"`
	IsPlayer bool   `json:"is_player"`
	Skin     *Skin  `json:"skin,omitempty"`
	Afk      *bool  `json:"afk,omitempty"`
	Team     *int16 `json:"team,omitempty"`
}

func (c *Client) Merge(c2 Client) {
	c.Name = utils.MergeValue(c.Name, c2.Name)
	c.Clan = utils.MergeValue(c.Clan, c2.Clan)
	c.Country = utils.MergeValue(c.Country, c2.Country)
	c.Score = utils.MergeValue(c.Score, c2.Score)
	c.IsPlayer = c.IsPlayer || c2.IsPlayer
	c.Skin.Merge(c2.Skin)
	c.Afk = utils.MergePointer(c.Afk, c2.Afk)
	c.Team = utils.MergePointer(c.Team, c2.Team)
}

func (c *Client) ToSQLC(address string) sqlc.InsertActiveServerClientsParams {
	return sqlc.InsertActiveServerClientsParams{
		Address:   address,
		Name:      c.Name,
		Clan:      c.Clan,
		CountryID: c.Country,
		Score:     c.Score,
		IsPlayer:  c.IsPlayer,
		Team:      c.Team,
	}
}

func (c *Client) IsPlayerInt64() int64 {
	if c.IsPlayer {
		return 1
	}
	return 0
}

func (c *Client) IsConnecting() bool {
	return c.Name == "(connecting)" &&
		c.Score == -1 &&
		c.Clan == "" &&
		!c.IsPlayer
}

type Skin struct {
	Name       *string `json:"name,omitempty"`
	ColorBody  *int32  `json:"color_body,omitempty"`
	ColorFeet  *int32  `json:"color_feet,omitempty"`
	Body       *Part   `json:"body,omitempty"`
	Marking    *Part   `json:"marking,omitempty"`
	Decoration *Part   `json:"decoration,omitempty"`
	Hands      *Part   `json:"hands,omitempty"`
	Feet       *Part   `json:"feet,omitempty"`
	Eyes       *Part   `json:"eyes,omitempty"`
}

func (s *Skin) Merge(s2 *Skin) {
	s.Name = utils.MergePointer(s.Name, s2.Name)
	s.ColorBody = utils.MergePointer(s.ColorBody, s2.ColorBody)
	s.ColorFeet = utils.MergePointer(s.ColorFeet, s2.ColorFeet)
	s.Body.Merge(s2.Body)
	s.Marking.Merge(s2.Marking)
	s.Decoration.Merge(s2.Decoration)
	s.Hands.Merge(s2.Hands)
	s.Feet.Merge(s2.Feet)
	s.Eyes.Merge(s2.Eyes)
}

type Part struct {
	Name  string `json:"name"`
	Color *int32 `json:"color,omitempty"`
}

func (p *Part) Merge(p2 *Part) {
	p.Name = utils.MergeValue(p.Name, p2.Name)
	p.Color = utils.MergePointer(p.Color, p2.Color)
}

var pointGametypes = []string{
	"alien",
	"ball",
	"bomb",
	"bunter",
	"catch",
	"city",
	"ctf",
	"dm",
	"fng",
	"foot",
	"freeze",
	"inf",
	"lms",
	"lts",
	"monster",
	"nodes",
	"rpg",
	"smash",
	"tdm",
	"teemo",
	"town",
	"war3",
	"xpanic",
	"zomb",
}

func isPointGameType(gameType string) bool {
	gameType = strings.ToLower(gameType)
	for _, gt := range pointGametypes {
		if strings.Contains(gameType, gt) {
			return true
		}
	}
	return false
}

func ScoreKindFromDTO(clientScoreKind *string, gameType string) string {
	if clientScoreKind == nil {
		return "points"
	}

	scoreKind := strings.ToLower(*clientScoreKind)
	if strings.Contains(scoreKind, "time") {
		scoreKind = "time"
		if isPointGameType(gameType) {
			scoreKind = "points"
		}
	} else {
		if !strings.Contains(scoreKind, "points") {
			log.Printf("unknown score kind %q", scoreKind)
		}
		scoreKind = "points"
	}

	return scoreKind
}

// expands the servers.Server DTO into a slice of Server models
func NewServersFromDTO(servers []servers.Server) ([]Server, error) {
	timestamp := time.Now()

	duplicates := make(map[string][]Server, len(servers))
	result := make([]Server, 0, len(servers))

	for _, server := range servers {
		info := server.Info

		scoreKind := ScoreKindFromDTO((*string)(info.ClientScoreKind), info.GameType)
		clients := make([]Client, 0, len(info.Clients))
		for _, client := range info.Clients {
			clients = append(clients, ClientFromDTO(client))
		}

		// split server that has multiple addresses and show it as individual server
		// one ipv4 addr, one ipv6 addr
		m := make(map[string][]string, len(server.Addresses))
		for _, addr := range server.Addresses {
			u, err := url.ParseRequestURI(addr)
			if err != nil {
				log.Println(fmt.Errorf("failed to parse address %s: %w", addr, err))
				continue
			}

			if u.Scheme != "" {
				host := u.Host
				scheme := u.Scheme
				m[host] = append(m[host], scheme)
			}
		}

		// same server, different addresses and potentially protocols
		for addr, protocols := range m {
			server := Server{
				Timestamp:    timestamp,
				Address:      addr,
				Protocols:    protocols,
				Name:         info.Name,
				Gametype:     info.GameType,
				Passworded:   info.Passworded,
				Map:          info.Map.Name,
				MapSha256Sum: info.Map.Sha256,
				MapSize:      info.Map.Size,
				Version:      info.Version,
				MaxClients:   info.MaxClients,
				MaxPlayers:   info.MaxPlayers,
				ScoreKind:    scoreKind,
				Clients:      clients,
			}
			duplicates[addr] = append(duplicates[addr], server)
			result = append(result, server)
		}
	}

	if len(duplicates) == len(result) {
		return result, nil
	}

	// deduplicate
	result = result[:0]

	for _, ss := range duplicates {
		if len(ss) == 1 {
			result = append(result, ss[0])
			continue
		}

		// merge duplicates
		var merged Server
		for s := range ss {
			merged.Merge(ss[s])
		}
		result = append(result, merged)
	}

	return result, nil
}

func ClientFromDTO(client servers.Client) Client {
	country := int16(-1)
	if math.MinInt16 <= client.Country && client.Country <= math.MaxInt16 {
		country = int16(client.Country)
	}

	return Client{
		Name:     client.Name,
		Clan:     client.Clan,
		Country:  country,
		Score:    client.Score,
		IsPlayer: client.IsPlayer,
		Skin:     SkinFromDTO(client.Skin),
		Afk:      client.Afk,
		Team:     client.Team,
	}
}

func SkinFromDTO(skin *servers.Skin) *Skin {
	if skin == nil {
		return nil
	}
	return &Skin{
		Name:       skin.Name,
		ColorBody:  skin.ColorBody,
		ColorFeet:  skin.ColorFeet,
		Body:       PartFromDTO(skin.Body),
		Marking:    PartFromDTO(skin.Marking),
		Decoration: PartFromDTO(skin.Decoration),
		Hands:      PartFromDTO(skin.Hands),
		Feet:       PartFromDTO(skin.Feet),
		Eyes:       PartFromDTO(skin.Eyes),
	}
}

func PartFromDTO(part *servers.Part) *Part {
	if part == nil {
		return nil
	}
	return &Part{
		Name:  part.Name,
		Color: part.Color,
	}
}
