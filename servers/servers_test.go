package servers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllServers(t *testing.T) {
	servers, err := GetAllServers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	assert.GreaterOrEqual(t, len(servers), 1)
}
