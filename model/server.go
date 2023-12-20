package model

import (
	"encoding/json"

	"github.com/jxsl13/twstatus-bot/servers"
)

type Server struct {
	Address      string
	Name         string
	Gametype     string
	Passworded   int
	Map          string
	MapSha256Sum *string
	MapSize      *int64
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

type Skin struct {
	Name       *string `json:"name,omitempty"`
	ColorBody  *int64  `json:"color_body,omitempty"`
	ColorFeet  *int64  `json:"color_feet,omitempty"`
	Body       *Part   `json:"body,omitempty"`
	Marking    *Part   `json:"marking,omitempty"`
	Decoration *Part   `json:"decoration,omitempty"`
	Hands      *Part   `json:"hands,omitempty"`
	Feet       *Part   `json:"feet,omitempty"`
	Eyes       *Part   `json:"eyes,omitempty"`
}

type Part struct {
	Name  string `json:"name"`
	Color *int64 `json:"color,omitempty"`
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

		for _, addr := range server.Addresses {
			server := Server{
				Address:      addr,
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
		Skin: &Skin{
			Name:      client.Skin.Name,
			ColorBody: client.Skin.ColorBody,
			ColorFeet: client.Skin.ColorFeet,
			Body: &Part{
				Name:  client.Skin.Body.Name,
				Color: client.Skin.Body.Color,
			},
			Marking: &Part{
				Name:  client.Skin.Marking.Name,
				Color: client.Skin.Marking.Color,
			},
			Decoration: &Part{
				Name:  client.Skin.Decoration.Name,
				Color: client.Skin.Decoration.Color,
			},
			Hands: &Part{
				Name:  client.Skin.Hands.Name,
				Color: client.Skin.Hands.Color,
			},
			Feet: &Part{
				Name:  client.Skin.Feet.Name,
				Color: client.Skin.Feet.Color,
			},
			Eyes: &Part{
				Name:  client.Skin.Eyes.Name,
				Color: client.Skin.Eyes.Color,
			},
		},
		Afk:  client.Afk,
		Team: client.Team,
	}
}
