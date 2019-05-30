package adapters_test

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	task := models.TaskSpec{Type: adapters.TaskTypeNoOp}
	adapter, err := adapters.For(task, store)
	adapter.Perform(models.RunResult{}, nil)
	assert.NoError(t, err)
}

func TestAdapterFor(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "rideShare", "https://dUber.eth")
	bt.MinimumContractPayment = assets.NewLink(10)
	assert.Nil(t, store.CreateBridgeType(bt))

	cases := []struct {
		name                   string
		bridgeName             string
		wantType               string
		wantMinContractPayment *assets.Link
		wantErrored            bool
	}{
		{"adapter not found", "nonExistent", "<nil>", nil, true},
		{"noop", "NoOp", "*adapters.NoOp", assets.NewLink(0), false},
		{"ethtx", "EthTx", "*adapters.EthTx", store.Config.MinimumContractPayment(), false},
		{"bridge mixed case", "rideShare", "*adapters.Bridge", assets.NewLink(10), false},
		{"bridge lower case", "rideshare", "*adapters.Bridge", assets.NewLink(10), false},
	}

	for _, test := range cases {
		t.Run(test.wantType, func(t *testing.T) {
			task := models.TaskSpec{Type: models.MustNewTaskType(test.bridgeName)}
			adapter, err := adapters.For(task, store)
			if test.wantErrored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantType, reflect.TypeOf(adapter.BaseAdapter).String())
				assert.Equal(t, test.wantMinContractPayment, adapter.MinContractPayment())
			}
		})
	}
}
