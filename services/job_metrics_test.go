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

	jsm, err := services.AllJobSpecMetrics(app.Store, []models.JobSpec{j1, j2})
	assert.NoError(t, err)

	assert.Len(t, jsm, 2)

	assert.Equal(t, jsm[0].ID, j1.ID)
	assert.Equal(t, jsm[0].AdaptorCount["noop"], 1)

	assert.Equal(t, jsm[1].ID, j2.ID)
	assert.Equal(t, jsm[1].AdaptorCount["noop"], 1)
	assert.Equal(t, jsm[1].ParamCount["url"][0].Value, "https://chain.link")
	assert.Equal(t, jsm[1].ParamCount["url"][0].Count, 1)
}