package adapters_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestBridge_Perform_fromUnstarted(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		wantPending bool
		response    string
	}{
		{"success", 200, "purchased", false, false, `{"data":{"value": "purchased"}}`},
		{"run error", 200, "lot 49", true, false, `{"error": "overload", "data": {}}`},
		{"server error", 400, "lot 49", true, false, `bad request`},
		{"server error", 500, "lot 49", true, false, `big error`},
		{"JSON parse error", 200, "lot 49", true, false, `}`},
		{"pending response", 200, "lot 49", false, true, `{"pending":true}`},
		{"unsetting value", 200, "", false, false, `{"data":{"value":null}}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	runID := utils.NewBytes32ID()
	wantedBody := fmt.Sprintf(`{"id":"%v","data":{"value":"lot 49"}}`, runID)

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) {
					assert.JSONEq(t, wantedBody, body)
				})
			defer cleanup()

			bt := cltest.NewBridgeType("auctionBidding", mock.URL)
			eb := &adapters.Bridge{bt}
			result := cltest.RunResultWithValue("lot 49")
			result.JobRunID = runID

			result = eb.Perform(result, store)
			val, _ := result.Get("value")
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Pending())
		})
	}
}

func TestBridge_Perform_resuming(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		input      string
		status     models.Status
		want       string
		wantStatus models.Status
	}{
		{"from pending", `{"value":"100","old":"remains"}`, models.StatusPending, `{"value":"100","old":"remains"}`, models.StatusInProgress},
		{"from errored", `{"value":"100","old":"remains"}`, models.StatusErrored, `{"value":"100","old":"remains"}`, models.StatusErrored},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	bt := cltest.NewBridgeType("auctionBidding", "https://notused.example.com")
	assert.Nil(t, store.Save(&bt))
	ba := &adapters.Bridge{bt}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := models.RunResult{
				Data:   cltest.JSONFromString(test.input),
				Status: test.status,
			}

			result := ba.Perform(input, store)

			assert.Equal(t, test.want, result.Data.String())
			assert.Equal(t, test.wantStatus, result.Status)
			if test.wantStatus.Errored() {
				assert.Equal(t, input, result)
			}
		})
	}
}
