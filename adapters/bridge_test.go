package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestBridgeAdapterPerform(t *testing.T) {
	t.Parallel()

	nilString := cltest.NullString(nil)
	cases := []struct {
		name     string
		status   int
		want     null.String
		errored  bool
		response string
	}{
		{"success", 200, cltest.NullString("purchased"), false, `{"output":{"value": "purchased"}}`},
		{"run error", 200, nilString, true, `{"error": "overload", "output": {}}`},
		{"server error", 500, nilString, true, `big error`},
		{"JSON parse error", 200, nilString, true, `}`},
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

			output := eb.Perform(input, store)
			assert.Equal(t, test.want, output.Output["value"])
			assert.Equal(t, test.errored, output.HasError())
		})
	}
}
