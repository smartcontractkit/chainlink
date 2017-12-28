package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/h2non/gock"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/create_jobs.json")
	resp, _ := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	respJSON := cltest.JobJSONFromResponse(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	var j models.Job
	app.Store.One("ID", respJSON.ID, &j)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")

	adapter1, _ := adapters.For(j.Tasks[0], app.Store)
	httpGet := adapter1.(*adapters.HttpGet)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")

	adapter2, _ := adapters.For(j.Tasks[1], app.Store)
	jsonParse := adapter2.(*adapters.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter4, _ := adapters.For(j.Tasks[3], app.Store)
	signTx := adapter4.(*adapters.EthSignAndSendTx)
	assert.Equal(t, signTx.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, signTx.FunctionID, "12345679")

	var initr models.Initiator
	app.Store.One("JobID", j.ID, &initr)
	assert.Equal(t, "cron", initr.Type)
	assert.Equal(t, "* * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestCreateJobsIntegration(t *testing.T) {
	RegisterTestingT(t)

	app := cltest.NewApplicationWithKeyStore()
	eth := app.MockEthClient()
	server := app.NewServer()
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
	rawTxResp := `0x6798b8110efe9c191a978d75954d0fbdd53bd866f7534fa0228802fa89d27b83`
	eth.Register("eth_sendRawTransaction", rawTxResp)
	eth.Register("eth_getTransactionReceipt", types.Receipt{TxHash: common.StringToHash(rawTxResp)})
	eth.Register("eth_blockNumber", "0x0111")

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/create_jobs.json")
	resp, err := cltest.BasicAuthPost(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	assert.Nil(t, err)
	defer resp.Body.Close()
	respJSON := cltest.JobJSONFromResponse(resp.Body)

	app.Start()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		app.Store.Where("JobID", respJSON.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	app.Scheduler.Stop()

	var job models.Job
	err = app.Store.One("ID", respJSON.ID, &job)
	assert.Nil(t, err)

	jobRuns, err = app.Store.JobRunsFor(job)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobRuns))
	jobRun := jobRuns[0]
	assert.Equal(t, tickerResponse, jobRun.TaskRuns[0].Result.Value())
	assert.Equal(t, "10583.75", jobRun.TaskRuns[1].Result.Value())
	assert.Equal(t, rawTxResp, jobRun.TaskRuns[3].Result.Value())
	assert.Equal(t, rawTxResp, jobRun.Result.Value())
}

func TestCreateInvalidJobs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplication()
	server := app.NewServer()
	defer app.Stop()

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/create_invalid_jobs.json")
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

	jsonStr := cltest.LoadJSON("../internal/fixtures/web/create_invalid_cron.json")
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
