package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	config := cltest.NewConfig()
	task := models.Task{Type: "NoOp"}
	adapter, err := adapters.For(task, config)
	adapter.Perform(models.RunResult{})
	assert.Nil(t, err)
	rval := adapter.(*adapters.NoOp).Config
	assert.NotEqual(t, "", rval.EthereumURL)
}
