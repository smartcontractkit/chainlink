package adapters_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestBridge_PerformEmbedsParamsInData(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

	data := ""
	token := ""
	mock, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending": true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(t, b)
			data = body.Get("data").String()
			token = h.Get("Authorization")
		},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: bt, Params: &params}

	input := models.RunResult{
		Data:   cltest.JSONFromString(t, `{"result":"100"}`),
		Status: models.RunStatusUnstarted,
	}
	ba.Perform(input, store)

	assert.Equal(t, `{"bodyParam":true,"result":"100"}`, data)
	assert.Equal(t, "Bearer "+bt.OutgoingToken, token)
}

func TestBridge_PerformAcceptsNonJsonObjectResponses(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

	mock, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", fmt.Sprintf(`{"jobRunID": "%s", "data": 251990120, "statusCode": 200}`, models.NewID()),
		func(h http.Header, b string) {},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: bt, Params: &params}

	input := models.RunResult{
		Data:   cltest.JSONFromString(t, `{"jobRunID": "jobID", "data": 251990120, "statusCode": 200}`),
		Status: models.RunStatusUnstarted,
	}
	result := ba.Perform(input, store)
	assert.NoError(t, result.GetError())
	assert.Equal(t, "251990120", result.Data.Get("result").String())
}

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

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", "")

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending": true}`)
			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			ba := &adapters.Bridge{BridgeType: bt}

			input := models.RunResult{
				Data:   cltest.JSONFromString(t, `{"result":"100"}`),
				Status: test.status,
			}

			result := ba.Perform(input, store)

			assert.Equal(t, `{"result":"100"}`, result.Data.String())
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
		{"success", 200, "purchased", false, false, `{"data":{"result": "purchased"}}`},
		{"run error", 200, "lot 49", true, false, `{"error": "overload", "data": {}}`},
		{"server error", 400, "lot 49", true, false, `bad request`},
		{"server error", 500, "lot 49", true, false, `big error`},
		{"JSON parse error", 200, "lot 49", true, false, `}`},
		{"pending response", 200, "lot 49", false, true, `{"pending":true}`},
		{"unsetting result", 200, "", false, false, `{"data":{"result":null}}`},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", "")
	runID := models.NewID()
	wantedBody := fmt.Sprintf(`{"id":"%v","data":{"result":"lot 49"}}`, runID)

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, ensureCalled := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(_ http.Header, body string) {
					assert.JSONEq(t, wantedBody, body)
				})
			defer ensureCalled()

			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: bt}
			input := cltest.RunResultWithResult("lot 49")
			input.CachedJobRunID = runID

			result := eb.Perform(input, store)
			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Status.PendingBridge())
		})
	}
}

func TestBridge_Perform_responseURL(t *testing.T) {
	input := cltest.RunResultWithResult("lot 49")
	input.CachedJobRunID = models.NewID()

	t.Parallel()
	cases := []struct {
		name          string
		configuredURL models.WebURL
		want          string
	}{
		{
			name:          "basic URL",
			configuredURL: cltest.WebURL(t, "https://chain.link"),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"},"responseURL":"https://chain.link/v2/runs/%s"}`, input.CachedJobRunID, input.CachedJobRunID),
		},
		{
			name:          "blank URL",
			configuredURL: cltest.WebURL(t, ""),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"}}`, input.CachedJobRunID),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()
			store.Config.Set("BRIDGE_RESPONSE_URL", test.configuredURL)

			mock, ensureCalled := cltest.NewHTTPMockServer(t, 200, "POST", ``,
				func(_ http.Header, body string) {
					assert.JSONEq(t, test.want, body)
				})
			defer ensureCalled()

			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: bt}
			eb.Perform(input, store)
		})
	}
}
