package performance

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"fmt"
	"github.com/smartcontractkit/chainlink/store/models"
)

func TestGetSchemas(t *testing.T) {
	schemas := GetSchemas()
	assert.True(t, len(schemas) == 2, "No schemas were returned from file")
}

func TestGetBasicAuthHeader(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	header := GetBasicAuthHeader(app.Store)
	assert.True(t, len(header.Get("Authorization")) > 0, "Authorization header wasn't set")
}

func TestTargets(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	targets := GetCreateJobTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/jobs", app.Server.URL))

	job := models.NewJob()
	app.AddJob(job)
	targets = GetViewJobTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/jobs/%s", app.Server.URL, job.ID))

	job.NewRun()
	targets = GetJobRunTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/jobs/%s/runs", app.Server.URL, job.ID))

	targets = GetViewJobRunTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/jobs/%s/runs", app.Server.URL, job.ID))
}