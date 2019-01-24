package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/scheduler_job.json")

	cltest.WaitForRuns(t, j, app.Store, 1)

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "* * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HelloWorld(t *testing.T) {
	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, 200, "GET", tickerResponse)
	defer assertCalled()

	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
	defer cleanup()
	eth := app.MockEthClient()

	newHeads := make(chan models.BlockHeader)
	attempt1Hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	attempt2Hash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")
	sentAt := uint64(23456)
	confirmed := sentAt + config.EthGasBumpThreshold() + 1
	safe := confirmed + config.MinOutgoingConfirmations() - 1
	unconfirmedReceipt := store.TxReceipt{}
	confirmedReceipt := store.TxReceipt{
		Hash:        attempt1Hash,
		BlockNumber: cltest.Int(confirmed),
	}

	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
	})
	assert.NoError(t, app.Start())
	eth.EventuallyAllCalled(t)

	eth.Context("ethTx.Perform()#1 at block 23456", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
		eth.Register("eth_sendRawTransaction", attempt1Hash) // Initial tx attempt sent
		eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))
	eth.EventuallyAllCalled(t)

	eth.Context("ethTx.Perform()#2 at block 23459", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmed-1))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_sendRawTransaction", attempt2Hash) // Gas bumped tx attempt sent
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(confirmed - 1)} // 23459: For Gas Bump
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	eth.Context("ethTx.Perform()#3 at block 23460", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(confirmed)} // 23460
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	eth.Context("ethTx.Perform()#4 at block 23465", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(safe))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", confirmedReceipt) // confirmed for gas bumped txat
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(safe)} // 23465
	eth.EventuallyAllCalled(t)

	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.TaskRuns[0].Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, tickerResponse, val)
	val, err = jr.TaskRuns[1].Result.Value()
	assert.Equal(t, "10583.75", val)
	assert.NoError(t, err)
	val, err = jr.TaskRuns[3].Result.Value()
	assert.Equal(t, attempt1Hash.String(), val)
	assert.NoError(t, err)
	val, err = jr.Result.Value()
	assert.Equal(t, attempt1Hash.String(), val)
	assert.NoError(t, err)
	assert.Equal(t, jr.Result.CachedJobRunID, jr.ID)

	eth.EventuallyAllCalled(t)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.InstantClock()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/run_at_job.json")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorRunAt, initr.Type)
	assert.Equal(t, "2018-01-08T18:12:01Z", initr.Time.ISO8601())

	app.Start()

	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_EthLog(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
	})
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs <- cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_RunLog(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig()
	defer cfgCleanup()
	config.Set("MIN_INCOMING_CONFIRMATIONS", 6)
	app, cleanup := cltest.NewApplicationWithConfig(config)
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	newHeads := eth.RegisterNewHeads()
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
	})
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/runlog_noop_job.json")
	requiredConfs := 100

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorRunLog, initr.Type)

	logBlockNumber := 1
	logs <- cltest.NewRunLog(j.ID, cltest.NewAddress(), cltest.NewAddress(), logBlockNumber, `{}`)
	cltest.WaitForRuns(t, j, app.Store, 1)

	runs, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)
	jr := runs[0]
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	minConfigHeight := logBlockNumber + int(app.Store.Config.MinIncomingConfirmations())
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(minConfigHeight)}
	<-time.After(time.Second)
	cltest.JobRunStaysPendingConfirmations(t, app.Store, jr)

	safeNumber := logBlockNumber + requiredConfs
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(safeNumber)}
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_EndAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	clock := cltest.UseSettableClock(app.Store)
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/end_at_job.json")
	endAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	assert.Equal(t, endAt, j.EndAt.Time)

	cltest.CreateJobRunViaWeb(t, app, j)

	clock.SetTime(endAt.Add(time.Nanosecond))

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", &bytes.Buffer{})
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode)
	gomega.NewGomegaWithT(t).Consistently(func() []models.JobRun {
		jobRuns, err := app.Store.JobRunsFor(j.ID)
		assert.NoError(t, err)
		return jobRuns
	}).Should(gomega.HaveLen(1))
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	clock := cltest.UseSettableClock(app.Store)
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/start_at_job.json")
	startAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	resp, cleanup := client.Post("/v2/specs/"+j.ID+"/runs", &bytes.Buffer{})
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode)
	cltest.WaitForRuns(t, j, app.Store, 0)

	clock.SetTime(startAt)

	cltest.CreateJobRunViaWeb(t, app, j)
}

func TestIntegration_ExternalAdapter_RunLogInitiated(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	newHeads := make(chan models.BlockHeader, 10)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.RegisterSubscription("newHeads", newHeads)
	})
	app.Start()

	eaValue := "87698118359"
	eaExtra := "other values to be used by external adapters"
	eaResponse := fmt.Sprintf(`{"data":{"value": "%v", "extra": "%v"}}`, eaValue, eaExtra)
	mockServer, ensureRequest := cltest.NewHTTPMockServer(t, 200, "POST", eaResponse)
	defer ensureRequest()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v","confirmations":10}`, mockServer.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/log_initiated_bridge_type_job.json")

	logBlockNumber := 1
	logs <- cltest.NewRunLog(j.ID, cltest.NewAddress(), cltest.NewAddress(), logBlockNumber, `{}`)
	jr := cltest.WaitForRuns(t, j, app.Store, 1)[0]
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 8)}
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 9)}
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, "randomnumber", tr.TaskSpec.Type.String())
	val, err := tr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, eaValue, val)
	res := tr.Result.Get("extra")
	assert.Equal(t, eaExtra, res.String())
}

// This test ensures that the response body of an external adapter are supplied
// as params to the successive task
func TestIntegration_ExternalAdapter_Copy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	bridgeURL := cltest.WebURL("https://test.chain.link/always")
	app.Store.Config.Set("BRIDGE_RESPONSE_URL", bridgeURL)
	app.Start()

	eaPrice := "1234"
	eaQuote := "USD"
	eaResponse := fmt.Sprintf(`{"data":{"price": "%v", "quote": "%v"}}`, eaPrice, eaQuote)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/", r.URL.Path)

		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		body := cltest.JSONFromString(string(b))
		data := body.Get("data")
		assert.True(t, data.Exists())
		bodyParam := data.Get("bodyParam")
		assert.True(t, bodyParam.Exists())
		assert.Equal(t, true, bodyParam.Bool())

		url := body.Get("responseURL")
		assert.Contains(t, url.String(), "https://test.chain.link/always/v2/runs")

		w.WriteHeader(200)
		io.WriteString(w, eaResponse)
	}))
	defer ts.Close()

	bridgeJSON := fmt.Sprintf(`{"name":"assetPrice","url":"%v"}`, ts.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/bridge_type_copy_job.json")
	jr := cltest.WaitForJobRunToComplete(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j, `{"copyPath": ["price"]}`))

	tr := jr.TaskRuns[0]
	assert.Equal(t, "assetprice", tr.TaskSpec.Type.String())
	tr = jr.TaskRuns[1]
	assert.Equal(t, "copy", tr.TaskSpec.Type.String())
	val, err := tr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, eaPrice, val)
	res := tr.Result.Get("quote")
	assert.Equal(t, eaQuote, res.String())
}

// This test ensures that an bridge adapter task is resumed from pending after
// sending out a request to an external adapter and waiting to receive a
// request back
func TestIntegration_ExternalAdapter_Pending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	var bt models.BridgeType
	var j models.JobSpec
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(b)

			jrs := cltest.WaitForRuns(t, j, app.Store, 1)
			jr := jrs[0]
			id := body.Get("id")
			assert.True(t, id.Exists())
			assert.Equal(t, jr.ID, id.String())

			data := body.Get("data")
			assert.True(t, data.Exists())
			assert.Equal(t, data.Type, gjson.JSON)

			token := utils.StripBearer(h.Get("Authorization"))
			assert.Equal(t, bt.OutgoingToken, token)
		})
	defer cleanup()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v"}`, mockServer.URL)
	bt = cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j = cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/random_number_bridge_type_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToPendBridge(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusPendingBridge, tr.Status)
	val, err := tr.Result.Value()
	assert.Error(t, err)
	assert.Equal(t, "", val)

	jr = cltest.UpdateJobRunViaWeb(t, app, jr, `{"data":{"value":"100"}}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	tr = jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusCompleted, tr.Status)
	val, err = tr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestIntegration_WeiWatchers(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	eth := app.MockEthClient()
	eth.RegisterNewHead(1)
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
	})

	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(_ http.Header, body string) {
			marshaledLog, err := json.Marshal(&log)
			assert.NoError(t, err)
			assert.JSONEq(t, string(marshaledLog), body)
		})
	defer cleanup()

	j := cltest.NewJobWithLogInitiator()
	post := cltest.NewTask("httppost", fmt.Sprintf(`{"url":"%v"}`, mockServer.URL))
	tasks := []models.TaskSpec{post}
	j.Tasks = tasks
	j = cltest.CreateJobSpecViaWeb(t, app, j)

	app.Start()
	logs <- log

	jobRuns := cltest.WaitForRuns(t, j, app.Store, 1)
	jr := cltest.WaitForJobRunToComplete(t, app.Store, jobRuns[0])
	assert.Equal(t, jr.Result.CachedJobRunID, jr.ID)
}

func TestIntegration_MultiplierInt256(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/int256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"value":"-10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0674e", val)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/uint256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"value":"10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.Value()
	assert.NoError(t, err)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", val)
}

func TestIntegration_NonceManagement_firstRunWithExistingTXs(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/web_initiated_eth_tx_job.json")

	eth := app.MockEthClient()
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionCount", `0x0100`) // activate account nonce
	})
	require.NoError(t, app.StartAndConnect())

	createCompletedJobRun := func(blockNumber uint64, expectedNonce uint64) {
		hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
		eth.Context("ethTx.Perform()", func(eth *cltest.EthMock) {
			eth.Register("eth_blockNumber", utils.Uint64ToHex(blockNumber))
			eth.Register("eth_sendRawTransaction", hash)

			eth.Register("eth_getTransactionReceipt", store.TxReceipt{
				Hash:        hash,
				BlockNumber: cltest.Int(blockNumber),
			})
			confirmedBlockNumber := blockNumber + app.Store.Config.MinOutgoingConfirmations()
			eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmedBlockNumber))
		})

		jr := cltest.CreateJobRunViaWeb(t, app, j, `{"value":"0x11"}`)
		jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

		attempt := cltest.GetLastTxAttempt(t, app.Store)
		tx, err := app.Store.FindTx(attempt.TxID)
		assert.NoError(t, err)
		assert.Equal(t, expectedNonce, tx.Nonce)
	}

	createCompletedJobRun(100, uint64(0x0100))
	createCompletedJobRun(200, uint64(0x0101))
}

func TestIntegration_CreateServiceAgreement(t *testing.T) {
	t.Parallel()
	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionCount", `0x100`)
	})
	assert.NoError(t, app.Start())
	sa := cltest.FixtureCreateServiceAgreementViaWeb(t, app, "../internal/fixtures/web/noop_agreement.json")

	assert.NotEqual(t, "", sa.ID)
	j := cltest.FindJob(app.Store, sa.JobSpecID)

	assert.Equal(t, assets.NewLink(1000000000000000000), sa.Encumbrance.Payment)
	assert.Equal(t, uint64(300), sa.Encumbrance.Expiration)

	assert.Equal(t, time.Unix(1571523439, 0).UTC(), sa.Encumbrance.EndAt.Time)
	assert.NotEqual(t, "", sa.ID)

	// Request execution of the job associated with this ServiceAgreement
	logs <- cltest.NewServiceAgreementExecutionLog(
		j.ID,
		cltest.NewAddress(),
		cltest.NewAddress(),
		1,
		`{}`)

	runs := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, runs[0])

	eth.EventuallyAllCalled(t)

}

func TestIntegration_BulkDeleteRuns(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	job := cltest.NewJobWithWebInitiator()
	err := app.GetStore().SaveJob(&job)
	require.NoError(t, err)
	completedRun := job.NewRun(job.Initiators[0])
	completedRun.Status = models.RunStatusCompleted
	err = app.GetStore().SaveJobRun(&completedRun)
	require.NoError(t, err)

	client := app.NewHTTPClient()
	request := `{
		"status": ["completed", "errored"],
		"updatedBefore": "2100-11-28T21:24:03Z"
	}`
	resp, cleanup := client.Post("/v2/bulk_delete_runs", bytes.NewBufferString(request))
	defer cleanup()
	assert.Equal(t, 201, resp.StatusCode)
	task := models.BulkDeleteRunTask{}
	err = cltest.ParseJSONAPIResponse(resp, &task)
	assert.NoError(t, err)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		resp, cleanup = client.Get("/v2/bulk_delete_runs/" + task.ID)
		defer cleanup()
		assert.Equal(t, 200, resp.StatusCode)
		err = cltest.ParseJSONAPIResponse(resp, &task)
		assert.NoError(t, err)
		return task.Status == "completed"
	}).Should(gomega.BeTrue())

	runCount, err := app.Store.JobRunsCount()
	assert.NoError(t, err)
	assert.Equal(t, 0, runCount)
}
