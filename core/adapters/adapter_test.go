package adapters_test

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	task := models.TaskSpec{Type: adapters.TaskTypeNoOp}
	adapter, err := adapters.For(task, store.Config, store.ORM)
	adapter.Perform(models.RunInput{}, nil, nil)
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
		name        string
		bridgeName  string
		wantType    string
		wantErrored bool
	}{
		{"adapter not found", "nonExistent", "<nil>", true},
		{"noop", "NoOp", "*adapters.NoOp", false},
		{"ethtx", "EthTx", "*adapters.EthTx", false},
		{"bridge mixed case", "rideShare", "*adapters.Bridge", false},
		{"bridge lower case", "rideshare", "*adapters.Bridge", false},
	}

	for _, test := range cases {
		t.Run(test.wantType, func(t *testing.T) {
			task := models.TaskSpec{Type: models.MustNewTaskType(test.bridgeName)}
			adapter, err := adapters.For(task, store.Config, store.ORM)
			if test.wantErrored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantType, reflect.TypeOf(adapter.BaseAdapter).String())
			}
		})
	}
}
