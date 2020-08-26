package adapters_test

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBridge_PerformEmbedsParamsInData(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))

	data := ""
	meta := false
	token := ""
	mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(t, b)
			data = body.Get("data").String()
			meta = body.Get("meta").Exists()
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
	assert.False(t, meta)
	assert.Equal(t, "Bearer "+bt.OutgoingToken, token)
}

func setupJobRunAndStore(t *testing.T, txHash []byte, blockHash []byte) (*store.Store, *models.ID, func()) {
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	app.Store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))
	require.NoError(t, app.Start())
	jr := app.MustCreateJobRun(txHash, blockHash)

	return app.Store, jr.ID, cleanup
}

func TestBridge_IncludesMetaIfJobRunIsInDB(t *testing.T) {
	txHashHex := "d6432b8321d9988e664f23cfce392dff8221da36a44ebb622160156dcef4abb9"
	blockHashHex := "d5150a4f602af1de7ff51f02c5b55b130693596c68f00b7796ac2b0f51175675"
	txHash, _ := hex.DecodeString(txHashHex)
	blockHash, _ := hex.DecodeString(blockHashHex)
	store, jobRunID, cleanup := setupJobRunAndStore(t, txHash, blockHash)
	defer cleanup()

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
	assert.Equal(t, fmt.Sprintf(`{"initiator":{"transactionHash":"0x%s","blockHash":"0x%s"}}`, txHashHex, blockHashHex), meta)
	assert.Equal(t, "Bearer "+bt.OutgoingToken, token)
}

func TestBridge_PerformAcceptsNonJsonObjectResponses(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", cltest.WebURL(t, ""))
	jobRunID := models.NewID()
	taskRunID := *models.NewID()

	mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", fmt.Sprintf(`{"jobRunID": "%s", "data": 251990120, "statusCode": 200}`, jobRunID.String()),
		func(h http.Header, b string) {},
	)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	params := cltest.JSONFromString(t, `{"bodyParam": true}`)
	ba := &adapters.Bridge{BridgeType: *bt, Params: params}

	input := *models.NewRunInput(jobRunID, taskRunID, cltest.JSONFromString(t, `{"jobRunID": "jobID", "data": 251990120, "statusCode": 200}`), models.RunStatusUnstarted)
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

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Config.Set("BRIDGE_RESPONSE_URL", "")

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, _ := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending": true}`)
			_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
			ba := &adapters.Bridge{BridgeType: *bt}

			input := *models.NewRunInputWithResult(models.NewID(), *models.NewID(), "100", test.status)
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
	taskRunID := *models.NewID()
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

			input := *models.NewRunInput(runID, taskRunID, cltest.JSONFromString(t, `{"result": "lot 49"}`), models.RunStatusUnstarted)
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

func TestBridgeResponse_TooLarge(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	config.Set("DEFAULT_HTTP_LIMIT", "1")

	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	largePayload := `{"pending": true}`
	mock, serverCleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", largePayload,
		func(h http.Header, b string) {},
	)
	defer serverCleanup()

	_, bt := cltest.NewBridgeType(t, "auctionBidding", mock.URL)
	ba := &adapters.Bridge{BridgeType: *bt}

	input := cltest.NewRunInputWithResult("100")
	result := ba.Perform(input, store)

	require.Error(t, result.Error())
	assert.Contains(t, result.Error().Error(), "HTTP response too large")
	assert.Equal(t, "", result.Result().String())
}
