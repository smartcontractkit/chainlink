package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/scheduler_job.json")

	cltest.WaitForRuns(t, j, app.Store, 1)

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "* * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HelloWorld(t *testing.T) {
	gock.EnableNetworking()
	defer cltest.CloseGock(t)

	config, _ := cltest.NewConfig()
	cltest.AddPrivateKey(config, "../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	app, cleanup := cltest.NewApplicationWithConfig(config)
	assert.Nil(t, app.Store.KeyStore.Unlock(cltest.Password))
	eth := app.MockEthClient()

	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	gock.New("https://www.bitstamp.net").
		Get("/api/ticker/").
		Reply(200).
		JSON(tickerResponse)

	newHeads := make(chan models.BlockHeader, 10)
	eth.RegisterSubscription("newHeads", newHeads)
	eth.Register("eth_getTransactionCount", `0x0100`)
	hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	sentAt := uint64(23456)
	confirmed := sentAt + config.EthGasBumpThreshold + 1
	safe := confirmed + config.EthMinConfirmations

	eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	eth.Register("eth_sendRawTransaction", hash)
	eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})

	app.Start()
	defer cleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/hello_world_job.json")
	jr := cltest.WaitForJobRunToPend(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))

	eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmed-1))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(confirmed - 1)}

	eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{
		Hash:        hash,
		BlockNumber: cltest.BigHexInt(confirmed),
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(confirmed)}

	eth.Register("eth_blockNumber", utils.Uint64ToHex(safe))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{
		Hash:        hash,
		BlockNumber: cltest.BigHexInt(confirmed),
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(safe)}

	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.TaskRuns[0].Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, tickerResponse, val)
	val, err = jr.TaskRuns[1].Result.Value()
	assert.Equal(t, "10583.75", val)
	assert.Nil(t, err)
	val, err = jr.TaskRuns[3].Result.Value()
	assert.Equal(t, hash.String(), val)
	assert.Nil(t, err)
	val, err = jr.Result.Value()
	assert.Equal(t, hash.String(), val)
	assert.Nil(t, err)
	assert.Equal(t, jr.Result.JobRunID, jr.ID)

	eth.EventuallyAllCalled(t)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.InstantClock()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/run_at_job.json")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
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
	logs := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logs)
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs <- cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_RunLog(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logs)
	app.Start()

	gock.EnableNetworking()
	defer cltest.CloseGock(t)
	gock.New("https://etherprice.com").
		Get("/api").
		Reply(200).
		JSON(`{}`)

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/runlog_random_number_job.json")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, models.InitiatorRunLog, initr.Type)

	logs <- cltest.NewRunLog(j.ID, cltest.NewAddress(), `{"url":"https://etherprice.com/api"}`)

	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_EndAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	clock := cltest.UseSettableClock(app.Store)
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/end_at_job.json")
	endAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	assert.Equal(t, endAt, j.EndAt.Time)

	cltest.CreateJobRunViaWeb(t, app, j)

	clock.SetTime(endAt.Add(time.Nanosecond))

	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", &bytes.Buffer{})
	assert.Equal(t, 500, resp.StatusCode)
	gomega.NewGomegaWithT(t).Consistently(func() []models.JobRun {
		jobRuns, err := app.Store.JobRunsFor(j.ID)
		assert.Nil(t, err)
		return jobRuns
	}).Should(gomega.HaveLen(1))
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	clock := cltest.UseSettableClock(app.Store)
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/start_at_job.json")
	startAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	resp := cltest.BasicAuthPost(url, "application/json", &bytes.Buffer{})
	assert.Equal(t, 500, resp.StatusCode)
	cltest.WaitForRuns(t, j, app.Store, 0)

	clock.SetTime(startAt)

	cltest.CreateJobRunViaWeb(t, app, j)
}

func TestIntegration_ExternalAdapter(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	eaValue := "87698118359"
	eaExtra := "other values to be used by external adapters"
	eaResponse := fmt.Sprintf(`{"data":{"value": "%v", "extra": "%v"}}`, eaValue, eaExtra)
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", eaResponse)
	defer cleanup()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v"}`, mockServer.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/random_number_bridge_type_job.json")
	jr := cltest.WaitForJobRunToComplete(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))

	tr := jr.TaskRuns[0]
	assert.Equal(t, "randomnumber", tr.Task.Type)
	val, err := tr.Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, eaValue, val)
	res, err := tr.Result.Get("extra")
	assert.Nil(t, err)
	assert.Equal(t, eaExtra, res.String())
}

func TestIntegration_ExternalAdapter_Pending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	var j models.JobSpec
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(body string) {
			jrs := cltest.WaitForRuns(t, j, app.Store, 1)
			jr := jrs[0]
			assert.JSONEq(t, fmt.Sprintf(`{"id":"%v","data":{}}`, jr.ID), body)
		})
	defer cleanup()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v"}`, mockServer.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j = cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/random_number_bridge_type_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToPend(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, models.StatusPending, tr.Status)
	val, err := tr.Result.Value()
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	jr = cltest.UpdateJobRunViaWeb(t, app, jr, `{"data":{"value":"100"}}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	tr = jr.TaskRuns[0]
	assert.Equal(t, models.StatusCompleted, tr.Status)
	val, err = tr.Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, "100", val)
}

func TestIntegration_WeiWatchers(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	eth := app.MockEthClient()
	eth.RegisterNewHead(1)
	logs := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logs)

	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(body string) {
			marshaledLog, err := json.Marshal(&log)
			assert.Nil(t, err)
			assert.JSONEq(t, string(marshaledLog), body)
		})
	defer cleanup()

	j := cltest.NewJobWithLogInitiator()
	post := cltest.NewTask("httppost", fmt.Sprintf(`{"url":"%v"}`, mockServer.URL))
	tasks := []models.TaskSpec{post}
	j.Tasks = tasks
	cltest.CreateJobSpecViaWeb(t, app, j)

	app.Start()
	logs <- log

	jobRuns := cltest.WaitForRuns(t, j, app.Store, 1)
	jr := cltest.WaitForJobRunToComplete(t, app.Store, jobRuns[0])
	assert.Equal(t, jr.Result.JobRunID, jr.ID)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	gock.EnableNetworking()
	defer cltest.CloseGock(t)

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	gock.New("https://bitstamp.net").
		Get("/api/ticker").
		Reply(200).
		JSON(`{"last": "10221.30"}`)

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/uint256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", val)
}

func TestIntegration_MultiplierUint256String(t *testing.T) {
	gock.EnableNetworking()
	defer cltest.CloseGock(t)

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	gock.New("https://bitstamp.net").
		Get("/api/ticker").
		Reply(200).
		JSON(`{"last": "10221.30"}`)

	j := cltest.FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/uint256_string_times_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.Value()
	assert.Nil(t, err)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", val)
}
