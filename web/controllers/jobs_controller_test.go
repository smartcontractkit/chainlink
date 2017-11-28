package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/models/tasks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type JobJSON struct {
	ID string `json:"id"`
}

func TestCreateJobs(t *testing.T) {
	db := cltest.SetUpDB()
	defer cltest.TearDownDB()
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	jsonStr := cltest.LoadJSON("./fixtures/create_jobs.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	var respJSON JobJSON
	json.Unmarshal(b, &respJSON)

	var j models.Job
	db.One("ID", respJSON.ID, &j)
	sched := j.Schedule
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
	assert.Equal(t, "* 7 * * *", string(sched.Cron), "Wrong cron schedule saved")
	assert.Equal(t, (*models.Time)(nil), sched.StartAt, "Wrong start at saved")
	endAt := models.Time{cltest.TimeParse("2019-11-27T23:05:49Z")}
	assert.Equal(t, endAt, *sched.EndAt, "Wrong end at saved")
	runAt0 := models.Time{cltest.TimeParse("2018-11-27T23:05:49Z")}
	assert.Equal(t, runAt0, sched.RunAt[0], "Wrong run at saved")

	httpGet := j.Tasks[0].Adapter.(*tasks.HttpGet)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")

	jsonParse := j.Tasks[1].Adapter.(*tasks.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	bytes32 := j.Tasks[2].Adapter.(*tasks.EthBytes32)
	assert.Equal(t, bytes32.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, bytes32.FunctionID, "12345679")
}

func TestCreateInvalidJobs(t *testing.T) {
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	jsonStr := cltest.LoadJSON("./fixtures/create_invalid_jobs.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["IdoNotExist is not a supported adapter type"]}`, string(body), "Response should return JSON")
}

func TestCreateInvalidCron(t *testing.T) {
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	jsonStr := cltest.LoadJSON("./fixtures/create_invalid_cron.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`, string(body), "Response should return JSON")
}

func TestShowJobs(t *testing.T) {
	db := cltest.SetUpDB()
	defer cltest.TearDownDB()
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "9 9 9 9 6"}

	db.Save(&j)

	resp, err := http.Get(server.URL + "/jobs/" + j.ID)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var respJob models.Job
	json.Unmarshal(b, &respJob)
	assert.Equal(t, respJob.Schedule, j.Schedule, "should have the same schedule")
}

func TestShowNotFoundJobs(t *testing.T) {
	cltest.SetUpDB()
	defer cltest.TearDownDB()
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	resp, err := http.Get(server.URL + "/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}
