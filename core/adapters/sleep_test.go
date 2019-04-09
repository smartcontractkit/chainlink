package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
)

func TestSleep_Perform(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	adapter := adapters.Sleep{}
	err := json.Unmarshal([]byte(`{"until": 1332151919}`), &adapter)
	assert.NoError(t, err)

	result := adapter.Perform(models.RunResult{}, store)
	assert.Equal(t, string(models.RunStatusPendingSleep), string(result.Status))
}
