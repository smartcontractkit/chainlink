package web_test

import (
	"bytes"
	"testing"

	"github.com/h2non/gock"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobSchedulerIntegration(t *testing.T) {
	RegisterTestingT(t)

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.Start()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/scheduler_job.json")
	resp := cltest.BasicAuthPost(app.Server.URL+"/v2/jobs", "application/json", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")
	respJSON := cltest.JobJSONFromResponse(resp.Body)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", respJSON.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	var initr models.Initiator
	app.Store.One("JobID", respJSON.ID, &initr)
	assert.Equal(t, "cron", initr.Type)
	assert.Equal(t, "* * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestCreateJobIntegration(t *testing.T) {
	RegisterTestingT(t)

	config := cltest.NewConfig()
	cltest.AddPrivateKey(config, "../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	app, cleanup := cltest.NewApplicationWithConfig(config)
	assert.Nil(t, app.Store.KeyStore.Unlock(cltest.Password))
	eth := app.MockEthClient()
	server := app.Server
	app.Start()
	defer cleanup()

	defer cltest.CloseGock(t)
	gock.EnableNetworking()

	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	gock.New("https://www.bitstamp.net").
		Get("/api/ticker/").
		Reply(200).
		JSON(tickerResponse)

	eth.Register("eth_getTransactionCount", `0x0100`)
	hash, err := utils.StringToHash("0x83c52c31cd40a023728fbc21a570316acd4f90525f81f1d7c477fd958ffa467f")
	assert.Nil(t, err)
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.EthMinConfirmations
	eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	eth.Register("eth_sendRawTransaction", hash)
	eth.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})
	eth.Register("eth_blockNumber", utils.Uint64ToHex(safe))
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{
		Hash:        hash,
		BlockNumber: confirmed,
	})

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")
	resp := cltest.BasicAuthPost(
		server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	defer resp.Body.Close()
	jobID := cltest.JobJSONFromResponse(resp.Body).ID

	url := server.URL + "/v2/jobs/" + jobID + "/runs"
	resp = cltest.BasicAuthPost(url, "application/json", &bytes.Buffer{})
	jrID := cltest.JobJSONFromResponse(resp.Body).ID

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", jobID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	var job models.Job
	err = app.Store.One("ID", jobID, &job)
	assert.Nil(t, err)

	jobRuns, err = app.Store.JobRunsFor(job)
	assert.Nil(t, err)
	jobRun := jobRuns[0]
	assert.Equal(t, jrID, jobRun.ID)
	Eventually(func() string {
		assert.Nil(t, app.Store.One("ID", jobRun.ID, &jobRun))
		return jobRun.Status
	}).Should(Equal("completed"))
	assert.Equal(t, tickerResponse, jobRun.TaskRuns[0].Result.Value())
	assert.Equal(t, "10583.75", jobRun.TaskRuns[1].Result.Value())
	assert.Equal(t, hash.String(), jobRun.TaskRuns[3].Result.Value())
	assert.Equal(t, hash.String(), jobRun.Result.Value())

	assert.True(t, eth.AllCalled())
}

func TestCreateJobWithRunAtIntegration(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.InstantClock()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/run_at_job.json")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")
	var j models.Job
	app.Store.One("ID", respJSON.ID, &j)

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, "runAt", initr.Type)
	assert.Equal(t, "2018-01-08T18:12:01Z", initr.Time.ISO8601())

	app.Start()
	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", respJSON.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestCreateJobWithEthLogIntegration(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	eth := app.MockEthClient()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/eth_log_job.json")
	address, _ := utils.StringToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")
	var j models.Job
	app.Store.One("ID", respJSON.ID, &j)

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, "ethLog", initr.Type)
	assert.Equal(t, address, initr.Address)

	logs := make(chan store.EventLog, 1)
	eth.RegisterSubscription("logs", logs)
	app.Start()

	logs <- store.EventLog{Address: address}

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", respJSON.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
