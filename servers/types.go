package servers

type ServerList struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Addresses []string  `json:"addresses"`
	Location  *Location `json:"location,omitempty"`
	Info      Info      `json:"info"`
}

type Info struct {
	Name            string           `json:"name"`
	GameType        string           `json:"game_type"`
	Map             Map              `json:"map"`
	Version         string           `json:"version"`
	Passworded      bool             `json:"passworded"`
	MaxClients      int64            `json:"max_clients"`
	MaxPlayers      int64            `json:"max_players"`
	Clients         []Client         `json:"clients,omitempty"`
	ClientScoreKind *ClientScoreKind `json:"client_score_kind,omitempty"`
	ServerSignature *string          `json:"server_signature,omitempty"`
	AltamedaNet     *bool            `json:"altameda_net,omitempty"`
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

type Map struct {
	Name   string  `json:"name"`
	Sha256 *string `json:"sha256,omitempty"`
	Size   *int64  `json:"size,omitempty"`
}

type ClientScoreKind string

const (
	Points ClientScoreKind = "points"
	Time   ClientScoreKind = "time"
)

type Location string
