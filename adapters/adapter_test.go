package adapters_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
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
	adapter, err := adapters.For(task, store)
	adapter.Perform(models.RunResult{}, store)
	assert.Nil(t, err)
}

func TestAdapterFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	bt := models.NewBridgeType()
	bt.Name = "rideShare"
	u, err := url.Parse("https://dUber.eth")
	assert.Nil(t, err)
	bt.URL = models.WebURL{u}
	assert.Nil(t, store.Save(bt))

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
			raw := json.RawMessage{}
			assert.Nil(t, json.Unmarshal([]byte(`{}`), &raw))
			task := models.Task{
				Type:   test.bridgeName,
				Params: raw,
			}
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
