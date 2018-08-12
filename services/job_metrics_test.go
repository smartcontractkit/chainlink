package services_test

import (
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobMetrics_Get(t *testing.T) {
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

	app.JobMetrics.Start()
	jsm1 := app.JobMetrics.Get(j1.ID)

	assert.Equal(t, jsm1.ID, j1.ID)
	assert.Equal(t, jsm1.AdaptorCount["noop"], 1)

	jsm2 := app.JobMetrics.Get(j2.ID)
	assert.Equal(t, jsm2.ID, j2.ID)
	assert.Equal(t, jsm2.AdaptorCount["noop"], 1)
	assert.Equal(t, jsm2.ParamCount["url"][0].Value, "https://chain.link")
	assert.Equal(t, jsm2.ParamCount["url"][0].Count, 1)
}

func TestJobMetrics_Add(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	j1, initr := cltest.NewJobWithWebInitiator()
	j1.Initiators[0].Ran = true
	err := app.Store.SaveJob(&j1)
	assert.NoError(t, err)

	jr := j1.NewRun(initr)
	jr.ID = "run"
	jr.Status = "completed"
	tp, err := jr.TaskRuns[0].Task.Params.Add("url", "https://chain.link")
	jr.TaskRuns[0].Task.Params = tp
	assert.NoError(t, err)

	err = app.Store.Save(&jr)
	assert.NoError(t, err)

	err = app.JobMetrics.Add(j1)
	assert.NoError(t, err)

	jsm := app.JobMetrics.Get(j1.ID)
	assert.Equal(t, jsm.ID, j1.ID)
	assert.Equal(t, jsm.AdaptorCount["noop"], 1)
	assert.Equal(t, jsm.ParamCount["url"][0].Value, "https://chain.link")
	assert.Equal(t, jsm.ParamCount["url"][0].Count, 1)
}

func TestJobMetrics_AddRun(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	j1, initr := cltest.NewJobWithWebInitiator()
	j1.Initiators[0].Ran = true
	err := app.Store.SaveJob(&j1)
	assert.NoError(t, err)

	err = app.JobMetrics.Add(j1)
	assert.NoError(t, err)

	jr := j1.NewRun(initr)
	jr.ID = "run"
	jr.Status = "completed"
	tp, err := jr.TaskRuns[0].Task.Params.Add("url", "https://chain.link")
	jr.TaskRuns[0].Task.Params = tp
	assert.NoError(t, err)

	err = app.Store.Save(&jr)
	assert.NoError(t, err)

	app.JobMetrics.AddRun(jr)

	jsm := app.JobMetrics.Get(j1.ID)
	assert.Equal(t, jsm.ID, j1.ID)
	assert.Equal(t, jsm.AdaptorCount["noop"], 1)
	assert.Equal(t, jsm.ParamCount["url"][0].Value, "https://chain.link")
	assert.Equal(t, jsm.ParamCount["url"][0].Count, 1)
}
