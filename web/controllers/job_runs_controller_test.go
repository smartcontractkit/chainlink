package controllers_test

import (
	"encoding/json"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type JobRunsJSON struct {
	Runs []JobRun `json:"runs"`
}

type JobRun struct {
	ID string `json:"id"`
}

func TestJobRunsIndex(t *testing.T) {
	store := cltest.Store()
	defer store.Close()
	server := cltest.SetUpWeb(store)
	defer cltest.TearDownWeb()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "9 9 9 9 6"}
	err := store.Save(&j)
	assert.Nil(t, err)
	jr := j.Run()
	err = store.Save(jr)
	assert.Nil(t, err)

	resp, err := http.Get(server.URL + "/jobs/" + j.ID + "/runs")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var respJSON JobRunsJSON
	json.Unmarshal(b, &respJSON)
	assert.Equal(t, 1, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr.ID, respJSON.Runs[0].ID, "expected the run IDs to match")
}
