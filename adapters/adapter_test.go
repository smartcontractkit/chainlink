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

	tt := models.NewCustomTaskType()
	tt.Name = "rideShare"
	u, err := url.Parse("https://dUber.eth")
	assert.Nil(t, err)
	tt.URL = models.WebURL{u}
	assert.Nil(t, store.Save(tt))

	cases := []struct {
		taskType string
		want     string
		errored  bool
	}{
		{"NoOp", "*adapters.NoOp", false},
		{"EthTx", "*adapters.EthTx", false},
		{"nonExistent", "<nil>", true},
		{tt.Name, "*adapters.ExternalBridge", false},
		{strings.ToLower(tt.Name), "*adapters.ExternalBridge", false},
	}

	for _, test := range cases {
		t.Run(test.want, func(t *testing.T) {
			raw := json.RawMessage{}
			assert.Nil(t, json.Unmarshal([]byte(`{}`), &raw))
			task := models.Task{
				Type:   test.taskType,
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
