package model_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/servers"
	"github.com/jxsl13/twstatus-bot/testutils"
	"github.com/stretchr/testify/require"
)

func TestDtoToEntityMapping(t *testing.T) {
	name := testutils.FilePath("../testdata/issue-011-dto.json")
	data, err := os.ReadFile(name)
	require.NoError(t, err)

	ss := []servers.Server{}
	err = json.Unmarshal(data, &ss)
	require.NoError(t, err)
	require.Greater(t, len(ss), 0)

	sl, err := model.NewServersFromDTO(ss)
	require.NoError(t, err)
	require.Greater(t, len(ss), 0)

	servers := map[string][]model.Server{}
	for _, s := range sl {
		servers[s.Address] = append(servers[s.Address], s)
	}

	for address, servers := range servers {
		if len(servers) == 1 {
			continue
		}
		t.Errorf("address %s has %d entries", address, len(servers))
		d, _ := json.MarshalIndent(servers, "", "  ")
		t.Errorf("servers: %s", string(d))
	}

}
