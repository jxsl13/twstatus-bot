package main_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jxsl13/twstatus-bot/sqlc"
	"github.com/jxsl13/twstatus-bot/testutils"
	"github.com/stretchr/testify/require"
)

func TestDuplicateKeyViolationIssue011(t *testing.T) {
	name := testutils.FilePath("testdata/issue-011.json")
	data, err := os.ReadFile(name)
	require.NoError(t, err)

	servers := []sqlc.InsertActiveServersParams{}

	err = json.Unmarshal(data, &servers)
	require.NoError(t, err)

	m := map[string][]sqlc.InsertActiveServersParams{}
	for _, server := range servers {
		m[server.Address] = append(m[server.Address], server)
	}
	found := false
	for address, servers := range m {
		if len(servers) == 1 {
			continue
		}
		found = true

		t.Logf("address %s has %d entries", address, len(servers))
		d, _ := json.MarshalIndent(servers, "", "  ")
		t.Logf("servers: %s", string(d))
	}

	require.True(t, found)
}
