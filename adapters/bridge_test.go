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

func TestBridge_Perform_transitionsTo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		status     models.RunStatus
		wantStatus models.RunStatus
	}{
		{"from pending bridge", models.RunStatusPendingBridge, models.RunStatusInProgress},
		{"from errored", models.RunStatusErrored, models.RunStatusErrored},
		{"from pending confirmations", models.RunStatusPendingConfirmations, models.RunStatusPendingBridge},
		{"from in progress", models.RunStatusInProgress, models.RunStatusPendingBridge},
		{"from completed", models.RunStatusCompleted, models.RunStatusCompleted},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending": true}`)
			bt := cltest.NewBridgeType("auctionBidding", mock.URL)
			ba := &adapters.Bridge{BridgeType: bt}

			input := models.RunResult{
				Data:   cltest.JSONFromString(`{"value":"100"}`),
				Status: test.status,
			}

			result := ba.Perform(input, store)

			assert.Equal(t, `{"value":"100"}`, result.Data.String())
			assert.Equal(t, test.wantStatus, result.Status)
			if test.wantStatus.Errored() || test.wantStatus.Completed() {
				assert.Equal(t, input, result)
			}
		})
	}
}

func TestBridge_Perform_startANewRun(t *testing.T) {
	t.Parallel()
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

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, ensureCalled := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) {
					assert.JSONEq(t, wantedBody, body)
				})
			defer ensureCalled()

			bt := cltest.NewBridgeType("auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: bt}
			input := cltest.RunResultWithValue("lot 49")
			input.JobRunID = runID

			result := eb.Perform(input, store)
			val, _ := result.Get("value")
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Status.PendingBridge())
		})
	}
}
