package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/orm"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JobJSON struct {
	ID string `json:"id"`
}

func TestCreateJobs(t *testing.T) {
	server := httptest.NewServer(web.Router())
	defer server.Close()
	orm.InitTest()
	defer orm.Close()
	db := orm.GetDB()

	jsonStr := []byte(`{"subtasks":[{"adapterType": "httpJSON", "adapterParams": {"endpoint": "https://bitstamp.net/api/ticker/", "fields": ["last"]}}], "schedule": "* 7 * * *","version":"1.0.0"}`)
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
}

func TestCreateInvalidJobs(t *testing.T) {
	server := httptest.NewServer(web.Router())
	defer server.Close()

	jsonStr := []byte(`{"subtasks":[{"adapterType": "ethereumBytes32", "adapterParams": {}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["\"ethereumBytes32\" is not a supported adapter type."]}`, string(body), "Repsonse should return JSON")
}
