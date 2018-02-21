package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestBridge_Perform(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantExists  bool
		wantErrored bool
		wantPending bool
		response    string
	}{
		{"success", 200, "purchased", true, false, false, `{"data":{"value": "purchased"}}`},
		{"run error", 200, "", false, true, false, `{"error": "overload", "data": {}}`},
		{"server error", 400, "", false, true, false, `bad request`},
		{"server error", 500, "", false, true, false, `big error`},
		{"JSON parse error", 200, "", false, true, false, `}`},
		{"pending response", 200, "", false, false, true, `{"pending":true}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	wantedBody := `{"value":"lot 49"}`

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) { assert.Equal(t, wantedBody, body) })
			defer cleanup()

			bt := cltest.NewBridgeType("auctionBidding", mock.URL)
			eb := &adapters.Bridge{bt}
			input := models.RunResultWithValue("lot 49")

			result := eb.Perform(input, store)
			val, _ := result.Get("value")
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantExists, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Pending)
		})
	}
}
