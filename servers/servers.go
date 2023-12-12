package servers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Allows to change this url
var DDNetHTTPMasterUrl = "https://master1.ddnet.org/ddnet/15/servers.json"

func GetAllServers() ([]Server, error) {
	return GetServers(DDNetHTTPMasterUrl)
}

func newClient() *resty.Client {
	return resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "twstatus-bot")
}

func GetServers(url string) ([]Server, error) {
	var result ServerList
	resp, err := newClient().R().SetResult(&result).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error while fetching servers: %s", resp.Status())
	}

	return result.Servers, nil
}
