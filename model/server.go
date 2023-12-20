package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/jxsl13/twstatus-bot/servers"
)

type Server struct {
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
	Clients      []Client // serialized as json into database
}

func (s Server) ClientsJSON() []byte {
	data, _ := json.Marshal(s.Clients)
	return data
}

func (s Server) ProtocolsJSON() []byte {
	data, _ := json.Marshal(s.Protocols)
	return data
}

type Client struct {
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	Country  int    `json:"country"`
	Score    int    `json:"score"`
	IsPlayer bool   `json:"is_player"`
	Skin     *Skin  `json:"skin,omitempty"`
	Afk      *bool  `json:"afk,omitempty"`
	Team     *int   `json:"team,omitempty"`
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

// expands the servers.Server DTO into a slice of Server models
func NewServersFromDTO(servers []servers.Server) ([]Server, error) {
	result := make([]Server, 0, len(servers))

	for _, server := range servers {
		info := server.Info
		passworded := 0
		if info.Passworded {
			passworded = 1
		}

		scoreKind := "points"
		if info.ClientScoreKind != nil {
			scoreKind = string(*info.ClientScoreKind)
		}

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
				Address:      addr,
				Protocols:    protocols,
				Name:         info.Name,
				Gametype:     info.GameType,
				Passworded:   passworded,
				Map:          info.Map.Name,
				MapSha256Sum: info.Map.Sha256,
				MapSize:      info.Map.Size,
				Version:      info.Version,
				MaxClients:   int(info.MaxClients),
				MaxPlayers:   int(info.MaxPlayers),
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
