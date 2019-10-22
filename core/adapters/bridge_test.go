package adapters_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_PerformEmbedsParamsInData(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

	data := ""
	token := ""
	mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(t, b)
			data = body.Get("data").String()
			token = h.Get("Authorization")
		},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: *bt, Params: params}

	input := cltest.NewRunInputWithResult("100")
	result := ba.Perform(input, store)
	require.NoError(t, result.Error())
	assert.Equal(t, `{"bodyParam":true,"result":"100"}`, data)
	assert.Equal(t, "Bearer "+bt.OutgoingToken, token)
}

func TestBridge_PerformAcceptsNonJsonObjectResponses(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))
	jobRunID := models.NewID()

	mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", fmt.Sprintf(`{"jobRunID": "%s", "data": 251990120, "statusCode": 200}`, jobRunID.String()),
		func(h http.Header, b string) {},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: *bt, Params: params}

	input := *models.NewRunInput(jobRunID, cltest.JSONFromString(t, `{"jobRunID": "jobID", "data": 251990120, "statusCode": 200}`), models.RunStatusUnstarted)
	result := ba.Perform(input, store)
	require.NoError(t, result.Error())
	resultString, err := result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "251990120", resultString)
}

func TestBridge_Perform_transitionsTo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		status     models.RunStatus
		wantStatus models.RunStatus
		result     string
	}{
		{"from pending bridge", models.RunStatusPendingBridge, models.RunStatusInProgress, `{"result":"100"}`},
		{"from errored", models.RunStatusErrored, models.RunStatusErrored, ""},
		{"from in progress", models.RunStatusInProgress, models.RunStatusPendingBridge, ""},
		{"from completed", models.RunStatusCompleted, models.RunStatusCompleted, `{"result":"100"}`},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", "")

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`)
			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			ba := &adapters.Bridge{BridgeType: *bt}

			input := *models.NewRunInputWithResult(models.NewID(), "100", test.status)
			result := ba.Perform(input, store)

			assert.Equal(t, test.result, result.Data().String())
			assert.Equal(t, test.wantStatus, result.Status())
			if test.wantStatus.Completed() {
				assert.Equal(t, input.Data(), result.Data())
				assert.Equal(t, input.Status(), result.Status())
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
		{"success", http.StatusOK, "purchased", false, false, `{"data":{"result": "purchased"}}`},
		{"run error", http.StatusOK, "", true, false, `{"error": "overload", "data": {}}`},
		{"server error", http.StatusBadRequest, "", true, false, `bad request`},
		{"server error", http.StatusInternalServerError, "", true, false, `big error`},
		{"JSON parse error", http.StatusOK, "", true, false, `}`},
		{"pending response", http.StatusOK, "", false, true, `{"pending":true}`},
		{"unsetting result", http.StatusOK, "", false, false, `{"data":{"result":null}}`},
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
			eb := &adapters.Bridge{BridgeType: *bt}

			input := *models.NewRunInput(runID, cltest.JSONFromString(t, `{"result": "lot 49"}`), models.RunStatusUnstarted)
			result := eb.Perform(input, store)
			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Status().PendingBridge())
		})
	}
}

func TestBridge_Perform_responseURL(t *testing.T) {
	input := cltest.NewRunInputWithResult("lot 49")

	t.Parallel()
	cases := []struct {
		name          string
		configuredURL models.WebURL
		want          string
	}{
		{
			name:          "basic URL",
			configuredURL: cltest.WebURL(t, "https://chain.link"),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"},"responseURL":"https://chain.link/v2/runs/%s"}`, input.JobRunID().String(), input.JobRunID().String()),
		},
		{
			name:          "blank URL",
			configuredURL: cltest.WebURL(t, ""),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"}}`, input.JobRunID().String()),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()
			store.Config.Set("BRIDGE_RESPONSE_URL", test.configuredURL)

			mock, ensureCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", ``,
				func(_ http.Header, body string) {
					assert.JSONEq(t, test.want, body)
				})
			defer ensureCalled()

			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: *bt}
			eb.Perform(input, store)
		})
	}
}
