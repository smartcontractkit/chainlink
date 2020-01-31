package adapters_test

import (
	"fmt"
	"net/http"
	"testing"

	"chainlink/core/adapters"
	"chainlink/core/internal/cltest"
	"chainlink/core/store"
	"chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJobRunAndStore(t *testing.T) (*store.Store, *models.ID, func()) {
	app, cleanup := cltest.NewApplication(t)
	require.NoError(t, app.Start())
	store := app.Store
	jr := app.MustCreateJobRun(cltest.JSONFromString(t, `{"random": "meta"}`))

	return store, jr.ID, cleanup
}

func TestBridge_PerformEmbedsParamsInData(t *testing.T) {
	store, jobRunID, cleanup := setupJobRunAndStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

	data := ""
	meta := ""
	token := ""
	mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(t, b)
			data = body.Get("data").String()
			meta = body.Get("meta").String()
			token = h.Get("Authorization")
		},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: *bt, Params: params}

	input := cltest.NewRunInputWithResultAndJobRunID("100", jobRunID)
	result := ba.Perform(input, store)
	require.NoError(t, result.Error())
	assert.Equal(t, `{"bodyParam":true,"result":"100"}`, data)
	assert.Equal(t, `{"random":"meta"}`, meta)
	assert.Equal(t, "Bearer "+bt.OutgoingToken, token)
}

func TestBridge_PerformAcceptsNonJsonObjectResponses(t *testing.T) {
	store, jobRunID, cleanup := setupJobRunAndStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

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
	assert.Equal(t, "251990120", result.Result().String())
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
		{"from in progress", models.RunStatusInProgress, models.RunStatusPendingBridge, ""},
		{"from completed", models.RunStatusCompleted, models.RunStatusCompleted, `{"result":"100"}`},
	}

	store, jobRunID, cleanup := setupJobRunAndStore(t)
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`)
			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			ba := &adapters.Bridge{BridgeType: *bt}

			input := *models.NewRunInputWithResult(jobRunID, "100", test.status)
			result := ba.Perform(input, store)

			assert.Equal(t, test.result, result.Data().String())
			assert.Equal(t, test.wantStatus, result.Status())
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
		{"success", http.StatusOK, "purchased", false, false, `{"meta":{"random":"meta"},"data":{"result": "purchased"}}`},
		{"run error", http.StatusOK, "", true, false, `{"error": "overload", "meta":{"random":"meta"},"data": {}}`},
		{"server error", http.StatusBadRequest, "", true, false, `bad request`},
		{"server error", http.StatusInternalServerError, "", true, false, `big error`},
		{"JSON parse error", http.StatusOK, "", true, false, `}`},
		{"pending response", http.StatusOK, "", false, true, `{"pending":true}`},
		{"unsetting result", http.StatusOK, "", false, false, `{"meta":{"random":"meta"},"data":{"result":null}}`},
	}

	store, jobRunID, cleanup := setupJobRunAndStore(t)
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))
	defer cleanup()

	wantedBody := fmt.Sprintf(`{"id":"%v","data":{"result":"lot 49"},"meta":{"random":"meta"}}`, jobRunID)

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, ensureCalled := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(_ http.Header, body string) {
					assert.JSONEq(t, wantedBody, body)
				})
			defer ensureCalled()

			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: *bt}

			input := *models.NewRunInput(jobRunID, cltest.JSONFromString(t, `{"result": "lot 49"}`), models.RunStatusUnstarted)
			result := eb.Perform(input, store)
			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, test.wantPending, result.Status().PendingBridge())
		})
	}
}

func TestBridge_Perform_responseURL(t *testing.T) {
	t.Parallel()
	store, jobRunID, cleanup := setupJobRunAndStore(t)
	defer cleanup()
	input := cltest.NewRunInputWithResultAndJobRunID("lot 49", jobRunID)
	cases := []struct {
		name          string
		configuredURL models.WebURL
		want          string
	}{
		{
			name:          "basic URL",
			configuredURL: cltest.WebURL(t, "https://chain.link"),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"}, "meta":{"random":"meta"},"responseURL":"https://chain.link/v2/runs/%s"}`, input.JobRunID().String(), input.JobRunID().String()),
		},
		{
			name:          "blank URL",
			configuredURL: cltest.WebURL(t, ""),
			want:          fmt.Sprintf(`{"id":"%s","data":{"result":"lot 49"}, "meta":{"random":"meta"}}`, input.JobRunID().String()),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			store.Config.Set("BRIDGE_RESPONSE_URL", test.configuredURL)
			mock, ensureCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", ``,
				func(_ http.Header, body string) {
					fmt.Println("body", body)
					assert.JSONEq(t, test.want, body)
				})
			defer ensureCalled()

			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			eb := &adapters.Bridge{BridgeType: *bt}
			res := eb.Perform(input, store)
			fmt.Println(res)
		})
	}
}
