package performance

import (
	"io/ioutil"
	"encoding/base64"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/tsenart/vegeta/lib"
	"fmt"
	"time"
	"log"
	"github.com/smartcontractkit/chainlink/store"
	"net/http"
)

// Pre-load the schemas prior to run-time to reduce overhead
func GetSchemas() (schemaBytes [][]byte) {
	schemaPrefix := "../../internal/fixtures/web/"
	schemas      := [...]string{
		"uint256_job.json",
		"hello_world_job.json",
	}

	for _, schema := range schemas {
		byteArray, err := ioutil.ReadFile(schemaPrefix + schema)
		if err != nil {
			log.Fatal(err)
		}
		schemaBytes = append(schemaBytes, byteArray)
	}
	return schemaBytes
}

// Get the basic auth http header
func GetBasicAuthHeader(store *store.Store) (header http.Header) {
	header = http.Header{}
	authString := store.Config.BasicAuthUsername + ":" + store.Config.BasicAuthPassword
	header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))
	return
}

// Create the vegata targets based on the schemas used for creating jobs
func GetCreateJobTargets(app *cltest.TestApplication) (targets []vegeta.Target) {
	schemas := GetSchemas()
	for _, schema := range schemas {
		targets = append(targets, vegeta.Target{
			Method: "POST",
			URL:    fmt.Sprintf("%s/v2/specs", app.Server.URL),
			Body:   schema,
			Header: GetBasicAuthHeader(app.Store),
		})
	}
	return targets
}

// Create the vegata targets for viewing the jobs created
func GetViewJobTargets(app *cltest.TestApplication) (targets []vegeta.Target) {
	jobs, err := app.Store.Jobs()
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range jobs {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/v2/specs/%s", app.Server.URL, job.ID),
			Header: GetBasicAuthHeader(app.Store),
		})
	}
	return targets
}

// Create the vegata targets for creating job runs
func GetJobRunTargets(app *cltest.TestApplication) (targets []vegeta.Target) {
	jobs, err := app.Store.Jobs()
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range jobs {
		targets = append(targets, vegeta.Target{
			Method: "POST",
			URL:    fmt.Sprintf("%s/v2/specs/%s/runs", app.Server.URL, job.ID),
			Header: GetBasicAuthHeader(app.Store),
		})
	}
	return targets
}

// Create the vegata targets for viewing job runs
func GetViewJobRunTargets(app *cltest.TestApplication) (targets []vegeta.Target) {
	jobs, err := app.Store.Jobs()
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range jobs {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/v2/specs/%s/runs", app.Server.URL, job.ID),
			Header: GetBasicAuthHeader(app.Store),
		})
	}
	return targets
}

// Calculate the average latency to complete each job run
func CalculateAverageJobRunLatency(app *cltest.TestApplication) time.Duration {
	WaitForJobRunsToComplete(app)
	jobs, err := app.Store.Jobs()
	var durationSum, durationCount int64
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range jobs {
		jobRuns, err := app.Store.JobRunsFor(job.ID)
		if err != nil {
			log.Fatal(err)
		}
		for _, jobRun := range jobRuns {
			if jobRun.Status == "completed" {
				durationSum += jobRun.CompletedAt.Time.Sub(jobRun.CreatedAt).Nanoseconds()
				durationCount++
			}
		}
	}
	return time.Duration(durationSum/durationCount)
}

// Wait for all the jobs and their tasks to complete after the job runs
func WaitForJobRunsToComplete(app *cltest.TestApplication) {
	jobs, err := app.Store.Jobs()
	if err != nil {
		log.Fatal(err)
	}
	startTime := time.Now()
	for {
		completed := true
		for _, job := range jobs {
			jobRuns, err := app.Store.JobRunsFor(job.ID)
			if err != nil {
				log.Fatal(err)
			}
			if len(jobRuns) == 0 {
				completed = false
				break
			}
			for _, jobRun := range jobRuns {
				if jobRun.Status == "in progress" {
					completed = false
					break
				}
			}
		}
		if completed {
			break
		} else if time.Now().Sub(startTime) > time.Minute * 3 {
			log.Fatal("3 minute timeout hit on waiting for job runs to complete")
		}
	}
}

