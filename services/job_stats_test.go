package services_test

import (
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobStats_AllJobSpecStats(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	j1, _ := cltest.NewJobWithWebInitiator()
	j1.Initiators[0].Ran = true
	err := app.Store.SaveJob(&j1)
	assert.NoError(t, err)

	j2, initr := cltest.NewJobWithWebInitiator()
	j2.Initiators[0].Ran = true
	err = app.Store.SaveJob(&j2)
	assert.NoError(t, err)

	jr := j2.NewRun(initr)
	jr.ID = "run"
	jr.Status = "completed"
	tp, err := jr.TaskRuns[0].Task.Params.Add("url", "https://chain.link")
	jr.TaskRuns[0].Task.Params = tp
	assert.NoError(t, err)

	err = app.Store.Save(&jr)
	assert.NoError(t, err)

	jss, err := services.AllJobSpecStats(app.Store, []models.JobSpec{j1, j2})
	assert.NoError(t, err)

	assert.Len(t, jss.JobSpecCounts, 2)

	assert.Equal(t, jss.JobSpecCounts[0].ID, j1.ID)
	assert.Equal(t, jss.JobSpecCounts[0].AdaptorCount["noop"], 1)

	assert.Equal(t, jss.JobSpecCounts[1].ID, j2.ID)
	assert.Equal(t, jss.JobSpecCounts[1].AdaptorCount["noop"], 1)
	assert.Equal(t, jss.JobSpecCounts[1].ParamCount["url"][0].Value, "https://chain.link")
	assert.Equal(t, jss.JobSpecCounts[1].ParamCount["url"][0].Count, 1)
}

func TestJobStats_NoAccount(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j, _ := cltest.NewJobWithWebInitiator()
	j.Initiators[0].Ran = true
	err := store.SaveJob(&j)
	assert.NoError(t, err)

	_, err = services.AllJobSpecStats(store, []models.JobSpec{j})
	assert.Error(t, err)
}
