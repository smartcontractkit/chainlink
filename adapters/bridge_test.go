package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestBridgeAdapterPerform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantExists  bool
		wantErrored bool
		response    string
	}{
		{"success", 200, "purchased", true, false, `{"output":{"value": "purchased"}}`},
		{"run error", 200, "", false, true, `{"error": "overload", "output": {}}`},
		{"server error", 500, "", false, true, `big error`},
		{"JSON parse error", 200, "", false, true, `}`},
		{"pending response", 200, "", false, false, `{"pending":true}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	wantedBody := `{"value":"lot 49"}`

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, wantedBody, test.response)
			defer cleanup()

			bt := cltest.NewBridgeType("auctionBidding", mock.URL)
			eb := &adapters.Bridge{bt}
			input := models.RunResultWithValue("lot 49")

			result := eb.Perform(input, store)
			val, _ := result.Get("value")
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantExists, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}
