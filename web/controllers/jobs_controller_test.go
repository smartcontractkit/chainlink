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

	jsonStr := []byte(`{"tasks":[{"type": "HttpGet", "params": {"endpoint": "https://bitstamp.net/api/ticker/"}}, {"type": "JsonParse", "params": {"path": ["last"]}}], "schedule": "* 7 * * *","version":"1.0.0"}`)
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
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
	assert.Equal(t, j.Schedule, "* 7 * * *", "Wrong schedule saved")

	httpGet := j.Tasks[0].Adapter.(*tasks.HttpGet)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")

	jsonParse := j.Tasks[1].Adapter.(*tasks.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})
}

func TestCreateInvalidJobs(t *testing.T) {
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	jsonStr := []byte(`{"tasks":[{"type": "ethereumBytes32", "params": {}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["ethereumBytes32 is not a supported adapter type"]}`, string(body), "Response should return JSON")
}

func TestShowJobs(t *testing.T) {
	db := cltest.SetUpDB()
	defer cltest.TearDownDB()
	server := cltest.SetUpWeb()
	defer cltest.TearDownWeb()

	j := models.NewJob()
	j.Schedule = "9 9 9 9 9"

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
