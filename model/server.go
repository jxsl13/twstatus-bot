package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/jxsl13/twstatus-bot/servers"
)

type Server struct {
	Timestamp    time.Time
	Address      string
	Protocols    []string
	Name         string
	Gametype     string
	Passworded   int64
	Map          string
	MapSha256Sum *string
	MapSize      *int64
	Version      string
	MaxClients   int64
	MaxPlayers   int64
	ScoreKind    string
	Clients      []Client // serialized as json into database
}

func (s *Server) ProtocolsJSON() []byte {
	data, _ := json.Marshal(s.Protocols)
	return data
}

func (s *Server) ProtocolsFromJSON(data []byte) error {
	return json.Unmarshal(data, &s.Protocols)
}

type Client struct {
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	Country  int64  `json:"country"`
	Score    int64  `json:"score"`
	IsPlayer bool   `json:"is_player"`
	Skin     *Skin  `json:"skin,omitempty"`
	Afk      *bool  `json:"afk,omitempty"`
	Team     *int64 `json:"team,omitempty"`
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

type Part struct {
	Name  string `json:"name"`
	Color *int32 `json:"color,omitempty"`
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
	result := make([]Server, 0, len(servers))

	for _, server := range servers {
		info := server.Info

		scoreKind := ScoreKindFromDTO((*string)(info.ClientScoreKind), info.GameType)
		clients := make([]Client, 0, len(info.Clients))
		for _, client := range info.Clients {
			clients = append(clients, ClientFromDTO(client))
		}

		m := make(map[string][]string, len(server.Addresses))
		for _, addr := range server.Addresses {
			u, err := url.ParseRequestURI(addr)
			if err != nil {
				log.Println(fmt.Errorf("failed to parse address %s: %w", addr, err))
				continue
			}

			if u.Scheme != "" {
				host := u.Host
				m[host] = append(m[host], u.Scheme)
			}
		}

		for addr, protocols := range m {
			server := Server{
				Timestamp:    timestamp,
				Address:      addr,
				Protocols:    protocols,
				Name:         info.Name,
				Gametype:     info.GameType,
				Passworded:   info.PasswordedInt64(),
				Map:          info.Map.Name,
				MapSha256Sum: info.Map.Sha256,
				MapSize:      info.Map.Size,
				Version:      info.Version,
				MaxClients:   info.MaxClients,
				MaxPlayers:   info.MaxPlayers,
				ScoreKind:    scoreKind,
				Clients:      clients,
			}
			result = append(result, server)
		}
	}
	return result, nil
}

func ClientFromDTO(client servers.Client) Client {
	return Client{
		Name:     client.Name,
		Clan:     client.Clan,
		Country:  client.Country,
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
