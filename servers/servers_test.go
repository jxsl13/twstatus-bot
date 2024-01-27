package servers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jxsl13/twstatus-bot/servers"
	"github.com/jxsl13/twstatus-bot/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllServers(t *testing.T) {
	data, servers, err := servers.GetAllServers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	assert.GreaterOrEqual(t, len(servers), 1)
	assert.GreaterOrEqual(t, len(data), 1)
}

func TestGetAllServerMods(t *testing.T) {
	_, servers, err := servers.GetAllServers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	m := make(map[string]struct{}, 1024)
	pointsM := make(map[string]struct{}, 1024)
	timeM := make(map[string]struct{}, 1024)
	unknownM := make(map[string]struct{}, 1024)
	for _, s := range servers {
		m[strings.ToLower(s.Info.GameType)] = struct{}{}

		if s.Info.ClientScoreKind != nil {
			if strings.HasPrefix(strings.ToLower(string(*s.Info.ClientScoreKind)), "points") {
				pointsM[strings.ToLower(s.Info.GameType)] = struct{}{}
			} else if strings.HasPrefix(strings.ToLower(string(*s.Info.ClientScoreKind)), "time") {
				timeM[strings.ToLower(s.Info.GameType)] = struct{}{}
			} else {
				unknownM[strings.ToLower(s.Info.GameType)] = struct{}{}
			}
		} else {
			unknownM[strings.ToLower(s.Info.GameType)] = struct{}{}
		}
	}

	fmt.Println("All gametypes:")
	fmt.Println(strings.Join(utils.SortedMapKeys(m), "\n"))
	fmt.Println("All gametypes with points:")
	fmt.Println(strings.Join(utils.SortedMapKeys(pointsM), "\n"))
	fmt.Println("All gametypes with time:")
	fmt.Println(strings.Join(utils.SortedMapKeys(timeM), "\n"))
	fmt.Println("All gametypes with unknown:")
	fmt.Println(strings.Join(utils.SortedMapKeys(unknownM), "\n"))
}
