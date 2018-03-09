package adapters_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	t.Parallel()
	task := models.TaskSpec{Type: "NoOp"}
	adapter, err := adapters.For(task, nil)
	adapter.Perform(models.RunResult{}, nil)
	assert.Nil(t, err)
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
		{bt.Name, "*adapters.Bridge", false},
		{strings.ToLower(bt.Name), "*adapters.Bridge", false},
	}

	for _, test := range cases {
		t.Run(test.want, func(t *testing.T) {
			task := models.TaskSpec{Type: test.bridgeName}
			adapter, err := adapters.For(task, store)
			if test.errored {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.want, reflect.TypeOf(adapter).String())
			}
		})
	}
}
