package adapters_test

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	task := models.TaskSpec{Type: adapters.TaskTypeNoOp}
	adapter, err := adapters.For(task, store)
	adapter.Perform(models.RunResult{}, nil)
	assert.NoError(t, err)
}

func TestAdapterFor(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	bt := cltest.NewBridgeType("rideShare", "https://dUber.eth")
	assert.Nil(t, store.Save(&bt))

	cases := []struct {
		bridgeName string
		want       string
		errored    bool
	}{
		{"NoOp", "*adapters.NoOp", false},
		{"EthTx", "*adapters.EthTx", false},
		{"nonExistent", "<nil>", true},
		{bt.Name.String(), "*adapters.Bridge", false},
		{bt.Name.String(), "*adapters.Bridge", false},
	}

	for _, test := range cases {
		t.Run(test.want, func(t *testing.T) {
			task := models.TaskSpec{Type: models.NewTaskType(test.bridgeName)}
			adapter, err := adapters.For(task, store)
			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				wa, ok := adapter.(adapters.MinConfsWrappedAdapter)
				if ok {
					assert.Equal(t, test.want, reflect.TypeOf(wa.Adapter).String())
				} else {
					assert.Equal(t, test.want, reflect.TypeOf(adapter).String())
				}
			}
		})
	}
}
