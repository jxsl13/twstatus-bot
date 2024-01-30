package servers

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
)

// Allows to change this url
var DDNetHTTPMasterUrl = "https://master1.ddnet.org/ddnet/15/servers.json"

func GetAllServers() ([]byte, []Server, error) {
	return GetServers(DDNetHTTPMasterUrl)
}

func newClient() *resty.Client {
	return resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "twstatus-bot")
}

func GetServers(url string) ([]byte, []Server, error) {
	var result ServerList
	resp, err := newClient().SetDoNotParseResponse(true).R().Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.RawBody().Close()

	data, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return data, nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return data, nil, err
	}

	if resp.IsError() {
		return nil, nil, fmt.Errorf("error while fetching servers: %s", resp.Status())
	}

	return data, result.Servers, nil
}
