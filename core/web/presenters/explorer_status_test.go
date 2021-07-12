package presenters

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExplorerStatus(t *testing.T) {
	t.Parallel()

	es := NewExplorerStatus(synchronization.NoopStatsPusher{})

	b, err := json.Marshal(es)
	require.NoError(t, err)

	expected := `{
		"status": "disconnected",
		"url": ""
	}`

	assert.JSONEq(t, expected, string(b))
}
