package web_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotsController_CreateSnapshot_V1_Format(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobWithAssignmentViaWeb(t, app, "../internal/fixtures/web/v1_format_job.json")

	path := "/v1/assignments/" + j.ID + "/snapshots"
	resp, cleanup := client.Post(path, bytes.NewBuffer([]byte{}))
	defer cleanup()

	runID := cltest.ParseCommonJSON(resp.Body).ID

	assert.NotNil(t, runID)
}

func TestSnapshotsController_CreateSnapshot_V1_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	url := "/v1/assignments/" + "badid" + "/snapshots"
	resp, cleanup := client.Post(url, bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestSnapshotsController_CreateSnapshot_V1_LateJob(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobWithAssignmentViaWeb(t, app, "../internal/fixtures/web/v1_format_job_past_endat_time.json")

	url := "/v1/assignments/" + j.ID + "/snapshots"
	resp, cleanup := client.Post(url, bytes.NewBuffer([]byte{}))
	defer cleanup()

	assert.Equal(t, 500, resp.StatusCode, "Response should be server error")
}

func TestSnapshotsController_ShowSnapshot_V1_Format(t *testing.T) {
	t.Parallel()

	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, 200, "GET", tickerResponse)
	defer assertCalled()

	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
	eth := app.MockEthClient()

	newHeads := make(chan models.BlockHeader, 10)
	eth.RegisterSubscription("newHeads", newHeads)
	eth.Register("eth_getTransactionCount", `0x0100`)
	hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	sentAt := uint64(23456)

	eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	eth.Register("eth_sendRawTransaction", hash)
	eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})

	assert.Nil(t, app.Start())
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.CreateMockAssignmentViaWeb(t, app, mockServer.URL)

	url := "/v1/assignments/" + j.ID + "/snapshots"
	resp, cleanup := client.Post(url, bytes.NewBuffer([]byte{}))
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	runID := cltest.ParseCommonJSON(resp.Body).ID

	cltest.WaitForRuns(t, j, app.Store, 1)
	jr, err := app.Store.FindJobRun(runID)
	assert.NoError(t, err)

	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	resp2, cleanup := client.Get("/v1/snapshots/" + runID)
	defer cleanup()
	assert.Equal(t, 200, resp2.StatusCode, "Response should be successful")

	var ss models.Snapshot
	assert.Nil(t, json.Unmarshal(cltest.ParseResponseBody(resp2), &ss))

	assert.Equal(t, ss.ID, jr.Result.JobRunID)
	assert.Equal(t, ss.Details.Result.String(), jr.Result.Data.Result.String())
	assert.Empty(t, ss.Error.String)
	assert.False(t, ss.Pending)
}

func TestSnapshotsController_ShowSnapshot_V1_NotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v1/snapshots/" + "badid")
	defer cleanup()
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}
