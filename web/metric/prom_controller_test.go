package metric_test

import (
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"testing"
)

func BenchmarkPromController_Show(b *testing.B) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := app.NewHTTPClient()
	setupPromControllerShow(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/metrics", map[string]string{
			"User-Agent": "Prometheus",
			"Authorization": "Bearer " + app.Config.MetricsBearerToken,
		})
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestPromController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := app.NewHTTPClient()

	_, _, err := setupPromControllerShow(app)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/metrics", map[string]string{
		"User-Agent": "Prometheus",
		"Authorization": "Bearer " + app.Config.MetricsBearerToken,
	})
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
}

func setupPromControllerShow(app *cltest.TestApplication) (*models.JobSpec, *models.JobSpec, error) {
	j1, _ := cltest.NewJobWithWebInitiator()
	j1.Initiators[0].Ran = true
	merr := app.Store.SaveJob(&j1)
	j2, initr := cltest.NewJobWithWebInitiator()
	j2.Initiators[0].Ran = true
	merr = multierr.Append(merr, app.Store.SaveJob(&j2))

	jr2 := j2.NewRun(initr)
	jr2.ID = "runB"
	jr2.Status = "completed"
	tp, err := jr2.TaskRuns[0].Task.Params.Add("url", "https://chain.link")
	merr = multierr.Append(merr, err)
	jr2.TaskRuns[0].Task.Params = tp
	merr = multierr.Append(merr, app.Store.Save(&jr2))
	return &j1, &j2, err
}
