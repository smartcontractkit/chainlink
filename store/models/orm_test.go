package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestWhereNotFound(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := store.Where("ID", "bogus", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestAllNotFound(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	var jobs []models.Job
	err := store.All(&jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORMSaveJob(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	j1 := cltest.NewJobWithSchedule("* * * * *")
	store.SaveJob(j1)

	var j2 models.Job
	store.One("ID", j1.ID, &j2)
	assert.Equal(t, j1.ID, j2.ID)

	var initr models.Initiator
	store.One("JobID", j1.ID, &initr)
	assert.Equal(t, models.Cron("* * * * *"), initr.Schedule)
}

func TestPendingJobRuns(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	j := models.NewJob()
	assert.Nil(t, store.SaveJob(j))
	npr := j.NewRun()
	assert.Nil(t, store.Save(&npr))

	pr := j.NewRun()
	pr.Status = "pending"
	assert.Nil(t, store.Save(&pr))

	pending, err := store.PendingJobRuns()
	assert.Nil(t, err)
	pendingIDs := []string{}
	for _, jr := range pending {
		pendingIDs = append(pendingIDs, jr.ID)
	}

	assert.Contains(t, pendingIDs, pr.ID)
	assert.NotContains(t, pendingIDs, npr.ID)
}
