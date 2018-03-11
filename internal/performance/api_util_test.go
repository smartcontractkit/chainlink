package performance

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"fmt"
	"github.com/smartcontractkit/chainlink/store/models"
	"time"
	"github.com/tsenart/vegeta/lib"
	"log"
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
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/specs", app.Server.URL))

	job := models.NewJob()
	app.AddJob(job)
	targets = GetViewJobTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/specs/%s", app.Server.URL, job.ID))

	job.NewRun()
	targets = GetJobRunTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/specs/%s/runs", app.Server.URL, job.ID))

	targets = GetViewJobRunTargets(app)
	assert.Equal(t, targets[0].URL, fmt.Sprintf("%s/v2/specs/%s/runs", app.Server.URL, job.ID))
}

func TestCalculateAverageJobRunLatency(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	rate := uint64(1)
	duration := time.Second

	createJobTargets := GetCreateJobTargets(app)
	targeter := vegeta.NewStaticTargeter(createJobTargets...)
	attacker := vegeta.NewAttacker()

	for range attacker.Attack(targeter, rate, duration){}

	jobRunTargets := GetJobRunTargets(app)
	targeter = vegeta.NewStaticTargeter(jobRunTargets...)
	attacker = vegeta.NewAttacker()

	for range attacker.Attack(targeter, rate, duration){}

	calculatedAverageLatency := CalculateAverageJobRunLatency(app)

	jobs, err := app.Store.Jobs()
	job := jobs[0]
	if err != nil {
		log.Fatal(err)
	}
	jobRuns, err := app.Store.JobRunsFor(job.ID)
	if err != nil {
		log.Fatal(err)
	}
	jobRun := jobRuns[0]
	averageLatency := jobRun.CompletedAt.Time.Sub(jobRun.CreatedAt)

	assert.Equal(t, averageLatency, calculatedAverageLatency)
}