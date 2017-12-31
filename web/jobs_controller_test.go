package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")
	resp, _ := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	var j models.Job
	app.Store.One("ID", respJSON.ID, &j)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")

	adapter1, _ := adapters.For(j.Tasks[0])
	httpGet := adapter1.(*adapters.HttpGet)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1])
	jsonParse := adapter2.(*adapters.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3])
	signTx := adapter4.(*adapters.EthTx)
	assert.Equal(t, signTx.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, signTx.FunctionID, "12345679")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, "web", initr.Type)
}

func TestCreateJobSchedulerIntegration(t *testing.T) {
	RegisterTestingT(t)

	app := cltest.NewApplication()
	server := app.NewServer()
	app.Start()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/no_op_job.json")
	resp, err := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	assert.Nil(t, err)
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
	app := cltest.NewApplicationWithConfig(config)
	app.Store.KeyStore.Unlock("password")
	eth := app.MockEthClient()
	server := app.NewServer()
	app.Start()
	defer app.Stop()

	err := app.Store.KeyStore.Unlock("password")
	assert.Nil(t, err)

	defer cltest.CloseGock(t)
	gock.EnableNetworking()

	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	gock.New("https://www.bitstamp.net").
		Get("/api/ticker/").
		Reply(200).
		JSON(tickerResponse)

	eth.Register("eth_getTransactionCount", `0x0100`)
	txid := `0x83c52c31cd40a023728fbc21a570316acd4f90525f81f1d7c477fd958ffa467f`
	confed := uint64(23456)
	eth.Register("eth_sendRawTransaction", txid)
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{})
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{TXID: txid, BlockNumber: confed})
	eth.Register("eth_getTransactionReceipt", store.TxReceipt{TXID: txid, BlockNumber: confed})
	eth.Register("eth_blockNumber", utils.Uint64ToHex(confed+config.EthConfMin-1))
	eth.Register("eth_blockNumber", utils.Uint64ToHex(confed+config.EthConfMin))

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json")
	resp, err := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	assert.Nil(t, err)
	defer resp.Body.Close()
	jobID := cltest.JobJSONFromResponse(resp.Body).ID

	url := server.URL + "/jobs/" + jobID + "/runs"
	resp, err = cltest.BasicAuthPost(url, "application/json", &bytes.Buffer{})
	assert.Nil(t, err)
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
	assert.Equal(t, txid, jobRun.TaskRuns[3].Result.Value())
	assert.Equal(t, txid, jobRun.Result.Value())
	assert.True(t, eth.AllCalled())
}

func TestCreateInvalidJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_job.json")
	resp, err := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["IdoNotExist is not a supported adapter type"]}`, string(body), "Response should return JSON")
}

func TestCreateInvalidCron(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/invalid_cron.json")
	resp, err := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`, string(body), "Response should return JSON")
}

func TestShowJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := cltest.NewJobWithSchedule("9 9 9 9 6")
	app.Store.Save(&j)

	resp, err := cltest.BasicAuthGet(server.URL + "/jobs/" + j.ID)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var respJob models.Job
	json.Unmarshal(b, &respJob)
	assert.Equal(t, respJob.Initiators[0].Schedule, j.Initiators[0].Schedule, "should have the same schedule")
}

func TestShowNotFoundJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	resp, err := cltest.BasicAuthGet(server.URL + "/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}

func TestShowJobUnauthenticated(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	resp, err := http.Get(server.URL + "/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Response should be forbidden")
}
