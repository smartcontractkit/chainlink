package web_test

import (
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"testing"
)

func BenchmarkStatsController_Index(b *testing.B) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := app.NewHTTPClient()
	setupStatsControllerIndex(app)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, cleanup := client.Get("/v2/stats")
		defer cleanup()
		assert.Equal(b, 200, resp.StatusCode, "Response should be successful")
	}
}

func TestStatsController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := app.NewHTTPClient()

	j1, j2, err := setupStatsControllerIndex(app)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v2/stats?size=x")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 422)

	resp, cleanup = client.Get("/v2/stats?size=1")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
	body := cltest.ParseResponseBody(resp)

	metaCount, err := cltest.ParseJSONAPIResponseMetaCount(body)
	assert.NoError(t, err)
	assert.Equal(t, 2, metaCount)

	var links jsonapi.Links
	stats := presenters.Stats{}
	err = web.ParsePaginatedResponse(body, &stats, &links)
	assert.NoError(t, err)
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)

	assert.Len(t, stats.JobSpecStats, 1)
	assert.Equal(t, stats.JobSpecStats[0].AdaptorCount["noop"], 1, "Should have noop as an adaptor")
	assert.Equal(t, j1.ID, stats.JobSpecStats[0].ID, "should have the same ID")

	resp, cleanup = client.Get(links["next"].Href)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	stats = presenters.Stats{}
	err = web.ParsePaginatedResponse(cltest.ParseResponseBody(resp), &stats, &links)
	assert.NoError(t, err)
	assert.Empty(t, links["next"])
	assert.NotEmpty(t, links["prev"])

	assert.Len(t, stats.JobSpecStats, 1)
	assert.Equal(t, stats.JobSpecStats[0].RunCount, 1, "Should have a single run")
	assert.Equal(t, stats.JobSpecStats[0].StatusCount["completed"], 1, "Should have a single completed run")
	assert.Equal(t, stats.JobSpecStats[0].AdaptorCount["noop"], 1, "Should have noop as an adaptor")
	assert.Equal(t, stats.JobSpecStats[0].ParamCount["url"][0].Value, "https://chain.link", "Should include the same URL")
	assert.Equal(t, stats.JobSpecStats[0].ParamCount["url"][0].Count, 1, "Should include a url")
	assert.Equal(t, j2.ID, stats.JobSpecStats[0].ID, "should have the same ID")
}

func setupStatsControllerIndex(app *cltest.TestApplication) (*models.JobSpec, *models.JobSpec, error) {
	j1, initr := cltest.NewJobWithWebInitiator()
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
