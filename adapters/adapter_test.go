package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	store := cltest.NewStore()
	defer store.Close()

	task := models.Task{Type: "NoOp"}
	adapter, err := adapters.For(task)
	adapter.Perform(models.RunResult{}, store)
	assert.Nil(t, err)
}
