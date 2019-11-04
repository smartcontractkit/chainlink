package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleep_Perform(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	adapter := adapters.Sleep{}
	err := json.Unmarshal([]byte(`{"until": 1332151919}`), &adapter)
	require.NoError(t, err)

	result := adapter.Perform(models.RunInput{}, store)
	require.NoError(t, result.Error())
	assert.Equal(t, string(models.RunStatusPendingSleep), string(result.Status()))
}
