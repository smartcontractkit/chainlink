package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	task := models.Task{Type: "NoOp"}
	adapter, err := adapters.For(task)
	adapter.Perform(models.RunResult{}, store)
	assert.Nil(t, err)
}
