package orm_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func TestORM_AllNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobs := cltest.AllJobs(t, store)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORM_NodeVersion(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ver, err := store.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Contains(t, ver.Version, "random")

	require.NoError(t, store.UpsertNodeVersion(models.NewNodeVersion("9.9.8")))

	ver, err = store.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Equal(t, "9.9.8", ver.Version)

	require.NoError(t, store.UpsertNodeVersion(models.NewNodeVersion("9.9.8")))
	require.NoError(t, store.UpsertNodeVersion(models.NewNodeVersion("9.9.7")))
	require.NoError(t, store.UpsertNodeVersion(models.NewNodeVersion("9.9.9")))

	ver, err = store.FindLatestNodeVersion()

	require.NoError(t, err)
	require.NotNil(t, ver)
	require.Equal(t, "9.9.9", ver.Version)
}

func TestORM_CreateJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * *")
	store.CreateJob(&j1)

	j2, err := store.FindJobSpec(j1.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	j1.Initiators[0].CreatedAt = j2.Initiators[0].CreatedAt
	j1.Initiators[0].UpdatedAt = j2.Initiators[0].UpdatedAt
	assert.Equal(t, j1.ID, j2.ID)
	assert.Equal(t, j1.Initiators[0], j2.Initiators[0])
	assert.Equal(t, j2.ID, j2.Initiators[0].JobSpecID)
}

func TestORM_Unscoped(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := store.ORM
	job := cltest.NewJob()
	err := orm.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		require.NoError(t, orm.CreateJob(&job))
		require.NoError(t, db.Delete(&job).Error)
		require.Error(t, db.First(&job).Error)
		err := store.ORM.Unscoped().RawDBWithAdvisoryLock(func(db *gorm.DB) error {
			require.NoError(t, db.First(&job).Error)
			return nil
		})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}

func TestORM_ShowJobWithMultipleTasks(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJob()
	job.Tasks = []models.TaskSpec{
		{Type: models.MustNewTaskType("task1")},
		{Type: models.MustNewTaskType("task2")},
		{Type: models.MustNewTaskType("task3")},
		{Type: models.MustNewTaskType("task4")},
	}
	assert.NoError(t, store.CreateJob(&job))

	orm := store.ORM
	retrievedJob, err := orm.FindJobSpec(job.ID)
	require.NoError(t, err)
	require.Len(t, retrievedJob.Tasks, 4)
	assert.Equal(t, string(retrievedJob.Tasks[0].Type), "task1")
	assert.Equal(t, string(retrievedJob.Tasks[1].Type), "task2")
	assert.Equal(t, string(retrievedJob.Tasks[2].Type), "task3")
	assert.Equal(t, string(retrievedJob.Tasks[3].Type), "task4")
}

func TestORM_CreateExternalInitiator(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	exi2, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.Equal(t, `ERROR: duplicate key value violates unique constraint "external_initiators_name_key" (SQLSTATE 23505)`, store.CreateExternalInitiator(exi2).Error())
}

func TestORM_DeleteExternalInitiator(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	_, err = store.FindExternalInitiator(token)
	require.NoError(t, err)

	err = store.DeleteExternalInitiator(exi.Name)
	require.NoError(t, err)

	_, err = store.FindExternalInitiator(token)
	require.Error(t, err)

	require.NoError(t, store.CreateExternalInitiator(exi))
}

func TestORM_ArchiveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithSchedule("* * * * *")
	require.NoError(t, store.CreateJob(&job))

	run := cltest.NewJobRun(job)
	require.NoError(t, store.CreateJobRun(&run))

	require.NoError(t, store.ArchiveJob(job.ID))

	require.Error(t, utils.JustError(store.FindJobSpec(job.ID)))
	require.Error(t, utils.JustError(store.FindJobRun(run.ID)))

	store.ORM.DB = store.DB.Unscoped().Session(&gorm.Session{})
	require.NoError(t, utils.JustError(store.FindJobSpec(job.ID)))
	require.NoError(t, utils.JustError(store.FindJobRun(run.ID)))
}

func TestORM_CreateJobRun_CreatesRunRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	rr := models.NewRunRequest(models.JSON{})
	currentHeight := big.NewInt(0)
	run, _ := services.NewRun(&job, &job.Initiators[0], currentHeight, rr, store.Config, store.ORM, new(mocks.Client), time.Now())
	require.NoError(t, store.CreateJobRun(run))

	requestCount, err := store.ORM.CountOf(&models.RunRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}

func TestORM_SaveJobRun_JobRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	t.Run("does not error on a job run with no task runs", func(t *testing.T) {
		job := cltest.NewJobWithWebInitiator()
		require.NoError(t, store.CreateJob(&job))
		rr := models.NewRunRequest(models.JSON{})
		currentHeight := big.NewInt(0)
		run, _ := services.NewRun(&job, &job.Initiators[0], currentHeight, rr, store.Config, store.ORM, new(mocks.Client), time.Now())
		require.NoError(t, store.CreateJobRun(run))
		run.TaskRuns = []models.TaskRun{}

		require.NoError(t, store.SaveJobRun(run))
	})

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	rr := models.NewRunRequest(models.JSON{})
	currentHeight := big.NewInt(0)
	run, _ := services.NewRun(&job, &job.Initiators[0], currentHeight, rr, store.Config, store.ORM, new(mocks.Client), time.Now())
	require.NoError(t, store.CreateJobRun(run))

	t.Run("if no results exist already, inserts them", func(t *testing.T) {
		require.NoError(t, store.SaveJobRun(run))

		require.True(t, run.ResultID.Valid)
		require.Equal(t, run.ResultID.Int64, run.Result.ID)
		require.Equal(t, models.JSON{}, run.Result.Data)
		require.False(t, run.Result.ErrorMessage.Valid)

		require.Len(t, run.TaskRuns, 1)

		tr := run.TaskRuns[0]
		require.True(t, tr.ResultID.Valid)
		require.Equal(t, tr.ResultID.Int64, tr.Result.ID)
		require.Equal(t, models.JSON{}, tr.Result.Data)
		require.False(t, tr.Result.ErrorMessage.Valid)

		loadedRun, err := store.FindJobRun(run.ID)
		require.NoError(t, err)
		require.Equal(t, *run, loadedRun)
	})

	t.Run("if results exist already, updates all run results for job run and task runs", func(t *testing.T) {
		run.Result.Data = cltest.JSONFromString(t, `{"foo": 42}`)
		run.Result.ErrorMessage = null.StringFrom(`something exploded`)

		run.TaskRuns[0].Result.Data = cltest.JSONFromString(t, `{"bar": 3.14}`)
		run.TaskRuns[0].Result.ErrorMessage = null.StringFrom(`something else exploded`)

		require.NoError(t, store.SaveJobRun(run))

		require.True(t, run.ResultID.Valid)
		require.Equal(t, run.ResultID.Int64, run.Result.ID)
		require.Equal(t, cltest.JSONFromString(t, `{"foo": 42}`), run.Result.Data)
		require.Equal(t, "something exploded", run.Result.ErrorMessage.String)

		require.Len(t, run.TaskRuns, 1)

		tr := run.TaskRuns[0]
		require.True(t, tr.ResultID.Valid)
		require.Equal(t, tr.ResultID.Int64, tr.Result.ID)
		require.Equal(t, cltest.JSONFromString(t, `{"bar": 3.14}`), tr.Result.Data)
		require.Equal(t, "something else exploded", tr.Result.ErrorMessage.String)

		loadedRun, err := store.FindJobRun(run.ID)
		require.NoError(t, err)
		require.Equal(t, *run, loadedRun)
	})

	t.Run("returns optimistic update failure if job run does not exist at all with that ID", func(t *testing.T) {
		run2 := &models.JobRun{ID: uuid.NewV4(), Status: models.RunStatusUnstarted}
		err := store.SaveJobRun(run2)

		require.Error(t, err)
		require.Equal(t, orm.ErrOptimisticUpdateConflict, errors.Cause(err))
	})

	t.Run("returns error if one of the task runs has not been inserted", func(t *testing.T) {
		ts := models.TaskSpec{Type: adapters.TaskTypeNoOp, JobSpecID: job.ID}
		require.NoError(t, store.DB.Create(&ts).Error)

		tr := models.TaskRun{ID: uuid.NewV4(), Status: models.RunStatusErrored, TaskSpecID: ts.ID, JobRunID: run.ID}
		run.TaskRuns = append(run.TaskRuns, tr)

		err := store.SaveJobRun(run)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("SaveJobRun failed: failed to insert run_result; task run with id %s was missing", tr.ID))
	})

	t.Run("if one task run result exists and one does not, does a mixture of inserts and updates", func(t *testing.T) {
		tr := run.TaskRuns[1]
		require.NoError(t, store.DB.Save(&tr).Error)

		run.TaskRuns[0].Result.Data = cltest.JSONFromString(t, `{"baz": 100}`)
		run.TaskRuns[0].Result.ErrorMessage = null.String{}
		run.TaskRuns[1].Result.ErrorMessage = null.StringFrom(`oh dear`)

		require.NoError(t, store.SaveJobRun(run))

		require.Len(t, run.TaskRuns, 2)

		tr = run.TaskRuns[0]
		assert.True(t, tr.ResultID.Valid)
		assert.Equal(t, tr.ResultID.Int64, tr.Result.ID)
		assert.Equal(t, cltest.JSONFromString(t, `{"baz": 100}`), tr.Result.Data)
		assert.False(t, tr.Result.ErrorMessage.Valid)

		tr = run.TaskRuns[1]
		assert.True(t, tr.ResultID.Valid)
		assert.Equal(t, tr.ResultID.Int64, tr.Result.ID)
		assert.Equal(t, cltest.JSONFromString(t, ``), tr.Result.Data)
		assert.Equal(t, "oh dear", tr.Result.ErrorMessage.String)

		loadedRun, err := store.FindJobRun(run.ID)
		require.NoError(t, err)
		require.Equal(t, *run, loadedRun)
	})

	t.Run("updates fields on the job run", func(t *testing.T) {
		finishedAt := null.TimeFrom(time.Unix(42, 0))
		status := models.RunStatusPendingSleep
		creationHeight := utils.NewBigI(43)
		observedHeight := utils.NewBigI(44)
		payment := assets.NewLink(45)

		run.Status = status
		run.FinishedAt = finishedAt
		run.CreationHeight = creationHeight
		run.ObservedHeight = observedHeight
		run.Payment = payment

		require.NoError(t, store.SaveJobRun(run))

		require.Equal(t, finishedAt, run.FinishedAt)
		require.Equal(t, status, run.Status)
		require.Equal(t, creationHeight, run.CreationHeight)
		require.Equal(t, observedHeight, run.ObservedHeight)
		require.Equal(t, payment, run.Payment)

		loadedRun, err := store.FindJobRun(run.ID)
		require.NoError(t, err)
		require.Equal(t, *run, loadedRun)
	})

	t.Run("updates fields on the task run", func(t *testing.T) {
		status := models.RunStatusPendingConnection

		run.TaskRuns[0].Status = status

		require.NoError(t, store.SaveJobRun(run))

		require.Equal(t, status, run.TaskRuns[0].Status)

		loadedRun, err := store.FindJobRun(run.ID)
		require.NoError(t, err)
		require.Equal(t, *run, loadedRun)
	})

	t.Run("inserted sync_event", func(t *testing.T) {
		se := models.SyncEvent{}
		err := store.DB.Order("id desc").First(&se).Error
		require.NoError(t, err)

		assert.Contains(t, se.Body, job.ID.String())
		assert.Contains(t, se.Body, run.ID.String())
	})

	t.Run("returns error if task run result is not preloaded", func(t *testing.T) {
		run.TaskRuns[1].Result = models.RunResult{}
		err := store.SaveJobRun(run)

		require.Error(t, err)
		require.Contains(t, err.Error(), "expected TaskRun.Result to be preloaded")
	})

	t.Run("returns error if job run result is not preloaded", func(t *testing.T) {
		run.Result = models.RunResult{}
		err := store.SaveJobRun(run)

		require.Error(t, err)
		require.Contains(t, err.Error(), "expected JobRun.Result to be preloaded")
	})
}

func TestORM_SaveJobRun_OptimisticLockFailure(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jr := cltest.CreateJobRunWithStatus(t, store, job, models.RunStatusUnstarted)

	// Something else updated it
	require.NoError(t, store.DB.Exec(`UPDATE job_runs SET updated_at = '1942-01-01'`).Error)

	err := store.SaveJobRun(&jr)
	require.Error(t, err)
	assert.True(t, errors.Is(err, orm.ErrOptimisticUpdateConflict))
}

func TestORM_SaveJobRun_ArchivedDoesNotRevertDeletedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	jr := cltest.CreateJobRunWithStatus(t, store, job, models.RunStatusUnstarted)

	require.NoError(t, store.ArchiveJob(job.ID))

	jr.SetStatus(models.RunStatusInProgress)
	require.NoError(t, store.SaveJobRun(&jr))

	require.Error(t, utils.JustError(store.FindJobRun(jr.ID)))
	require.NoError(t, utils.JustError(store.Unscoped().FindJobRun(jr.ID)))
}

func TestORM_SaveJobRun_Cancelled(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.ORM.SetLogging(true)

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	jr := cltest.NewJobRun(job)
	require.NoError(t, store.CreateJobRun(&jr))

	jr.SetStatus(models.RunStatusInProgress)
	require.NoError(t, store.SaveJobRun(&jr))

	jr.SetStatus(models.RunStatusCancelled)
	require.NoError(t, store.SaveJobRun(&jr))

	// Set a previous updated at to simulate a conflict
	jr.UpdatedAt = time.Unix(42, 0)
	jr.SetStatus(models.RunStatusInProgress)
	err := store.SaveJobRun(&jr)
	require.Error(t, err)
	assert.True(t, errors.Is(err, orm.ErrOptimisticUpdateConflict))
}

func TestORM_JobRunsFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jr1 := cltest.NewJobRun(job)
	jr1.CreatedAt = time.Now().AddDate(0, 0, -1)
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := cltest.NewJobRun(job)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, store.CreateJobRun(&jr2))
	jr3 := cltest.NewJobRun(job)
	jr3.CreatedAt = time.Now().AddDate(0, 0, -9)
	require.NoError(t, store.CreateJobRun(&jr3))

	runs, err := store.JobRunsFor(job.ID)
	assert.NoError(t, err)
	actual := []uuid.UUID{runs[0].ID, runs[1].ID, runs[2].ID}
	assert.Equal(t, []uuid.UUID{jr2.ID, jr1.ID, jr3.ID}, actual)

	limRuns, limErr := store.JobRunsFor(job.ID, 2)
	assert.NoError(t, limErr)
	limActual := []uuid.UUID{limRuns[0].ID, limRuns[1].ID}
	assert.Equal(t, []uuid.UUID{jr2.ID, jr1.ID}, limActual)

	_, limZeroErr := store.JobRunsFor(job.ID, 0)
	assert.NoError(t, limZeroErr)
	limZeroActual := []uuid.UUID{}
	assert.Equal(t, []uuid.UUID{}, limZeroActual)
}

func TestORM_LinkEarnedFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	jr1 := cltest.NewJobRun(job)
	jr1.TaskRuns[0].Status = models.RunStatusCompleted
	jr1.SetStatus(models.RunStatusCompleted)
	jr1.Payment = assets.NewLink(2)
	require.NoError(t, store.CreateJobRun(&jr1))

	jr2 := cltest.NewJobRun(job)
	jr2.TaskRuns[0].Status = models.RunStatusCompleted
	jr2.SetStatus(models.RunStatusCompleted)
	jr2.Payment = assets.NewLink(3)
	require.NoError(t, store.CreateJobRun(&jr2))

	jr3 := cltest.NewJobRun(job)
	jr3.TaskRuns[0].Status = models.RunStatusCompleted
	jr3.SetStatus(models.RunStatusCompleted)
	jr3.Payment = assets.NewLink(5)
	jr3.FinishedAt = null.TimeFrom(time.Now())
	require.NoError(t, store.CreateJobRun(&jr3))

	jr4 := cltest.NewJobRun(job)
	jr4.TaskRuns[0].Status = models.RunStatusCompleted
	jr4.SetStatus(models.RunStatusCompleted)
	jr4.Payment = assets.NewLink(5)
	jr4.FinishedAt = null.Time{}
	require.NoError(t, store.CreateJobRun(&jr4))

	jr5 := cltest.NewJobRun(job)
	jr5.SetStatus(models.RunStatusCancelled)
	jr5.Payment = assets.NewLink(5)
	require.NoError(t, store.CreateJobRun(&jr5))

	totalEarned, err := store.LinkEarnedFor(&job)
	require.NoError(t, err)
	assert.Equal(t, assets.NewLink(10), totalEarned)
}

func TestORM_JobRunsSortedFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	includedJob := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&includedJob))

	excludedJob := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&excludedJob))

	jr1 := cltest.NewJobRun(includedJob)
	jr1.CreatedAt = time.Now().AddDate(0, 0, -1)
	jr1.Status = models.RunStatusCompleted
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := cltest.NewJobRun(includedJob)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	jr2.Status = models.RunStatusErrored
	require.NoError(t, store.CreateJobRun(&jr2))

	excludedJobRun := cltest.NewJobRun(excludedJob)
	excludedJobRun.CreatedAt = time.Now().AddDate(0, 0, -9)
	require.NoError(t, store.CreateJobRun(&excludedJobRun))

	runs, count, completedCount, errorCount, err := store.JobRunsSortedFor(includedJob.ID, orm.Descending, 0, 100)
	assert.NoError(t, err)
	require.Equal(t, 2, count)
	require.Equal(t, 1, completedCount)
	require.Equal(t, 1, errorCount)
	actual := []uuid.UUID{runs[0].ID, runs[1].ID} // doesn't include excludedJobRun
	assert.Equal(t, []uuid.UUID{jr2.ID, jr1.ID}, actual)
}

func TestORM_UnscopedJobRunsWithStatus_Happy(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.CreateJob(&j))
	npr := cltest.NewJobRun(j)
	require.NoError(t, store.CreateJobRun(&npr))

	statuses := []models.RunStatus{
		models.RunStatusPendingBridge,
		models.RunStatusPendingIncomingConfirmations,
		models.RunStatusPendingOutgoingConfirmations,
		models.RunStatusCompleted}

	var seedIds []uuid.UUID
	for _, status := range statuses {
		run := cltest.NewJobRun(j)
		run.SetStatus(status)
		require.NoError(t, store.CreateJobRun(&run))
		seedIds = append(seedIds, run.ID)
	}

	tests := []struct {
		name     string
		statuses []models.RunStatus
		expected []uuid.UUID
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]uuid.UUID{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingIncomingConfirmations, models.RunStatusPendingOutgoingConfirmations},
			[]uuid.UUID{seedIds[0], seedIds[1], seedIds[2]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []uuid.UUID{}
			for _, jr := range pending {
				pendingIDs = append(pendingIDs, jr.ID)
			}
			assert.ElementsMatch(t, pendingIDs, test.expected)
		})
	}
}

func TestORM_UnscopedJobRunsWithStatus_Deleted(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.CreateJob(&j))
	npr := cltest.NewJobRun(j)
	require.NoError(t, store.CreateJobRun(&npr))

	statuses := []models.RunStatus{
		models.RunStatusPendingBridge,
		models.RunStatusPendingOutgoingConfirmations,
		models.RunStatusPendingIncomingConfirmations,
		models.RunStatusPendingConnection,
		models.RunStatusCompleted}

	var seedIds []uuid.UUID
	for _, status := range statuses {
		run := cltest.NewJobRun(j)
		run.SetStatus(status)
		require.NoError(t, store.CreateJobRun(&run))
		seedIds = append(seedIds, run.ID)
	}

	require.NoError(t, store.ArchiveJob(j.ID))

	tests := []struct {
		name     string
		statuses []models.RunStatus
		expected []uuid.UUID
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]uuid.UUID{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{
				models.RunStatusPendingBridge,
				models.RunStatusPendingOutgoingConfirmations,
				models.RunStatusPendingIncomingConfirmations,
				models.RunStatusPendingConnection},
			[]uuid.UUID{seedIds[0], seedIds[1], seedIds[2], seedIds[3]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []uuid.UUID{}
			for _, jr := range pending {
				pendingIDs = append(pendingIDs, jr.ID)
			}
			assert.ElementsMatch(t, pendingIDs, test.expected)
		})
	}
}

func TestORM_UnscopedJobRunsWithStatus_OrdersByCreatedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.CreateJob(&j))

	newPending := cltest.NewJobRun(j)
	newPending.SetStatus(models.RunStatusPendingSleep)
	newPending.CreatedAt = time.Now().Add(10 * time.Second)
	require.NoError(t, store.CreateJobRun(&newPending))

	oldPending := cltest.NewJobRun(j)
	oldPending.SetStatus(models.RunStatusPendingSleep)
	oldPending.CreatedAt = time.Now()
	require.NoError(t, store.CreateJobRun(&oldPending))

	runs := cltest.MustAllJobsWithStatus(t, store, models.RunStatusInProgress, models.RunStatusPendingSleep)
	require.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, oldPending.ID)
	assert.Equal(t, runs[1].ID, newPending.ID)
}

func TestORM_AnyJobWithType(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	js := cltest.NewJobWithWebInitiator()
	js.Tasks = []models.TaskSpec{{Type: models.MustNewTaskType("bridgetestname")}}
	assert.NoError(t, store.CreateJob(&js))
	found, err := store.AnyJobWithType("bridgetestname")
	assert.NoError(t, err)
	assert.Equal(t, found, true)
	found, err = store.AnyJobWithType("somethingelse")
	assert.NoError(t, err)
	assert.Equal(t, found, false)

}

func TestORM_JobRunsCountFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	job2 := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job2))

	assert.NotEqual(t, job.ID, job2.ID)

	completedRun := cltest.NewJobRun(job)
	run2 := cltest.NewJobRun(job)
	run3 := cltest.NewJobRun(job2)

	assert.NoError(t, store.CreateJobRun(&completedRun))
	assert.NoError(t, store.CreateJobRun(&run2))
	assert.NoError(t, store.CreateJobRun(&run3))

	count, err := store.JobRunsCountFor(job.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	count, err = store.JobRunsCountFor(job2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestORM_FindBridge(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, store.CreateBridgeType(&bt))

	cases := []struct {
		description string
		name        models.TaskType
		want        models.BridgeType
		errored     bool
	}{
		{"actual external adapter", bt.Name, bt, false},
		{"core adapter", "ethtx", models.BridgeType{}, true},
		{"non-existent adapter", "nonExistent", models.BridgeType{}, true},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			tt, err := store.FindBridge(test.name)
			tt.CreatedAt = test.want.CreatedAt
			tt.UpdatedAt = test.want.UpdatedAt
			assert.Equal(t, test.want, tt)
			assert.Equal(t, test.errored, err != nil)
		})
	}
}

func TestORM_FindBridgesByNames(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bt1 := models.BridgeType{}
	bt1.Name = models.MustNewTaskType("bridge1")
	bt1.URL = cltest.WebURL(t, "http://bridge1.com")
	require.NoError(t, store.CreateBridgeType(&bt1))

	bt2 := models.BridgeType{}
	bt2.Name = models.MustNewTaskType("bridge2")
	bt2.URL = cltest.WebURL(t, "http://bridge2.com")
	require.NoError(t, store.CreateBridgeType(&bt2))

	cases := []struct {
		description string
		arguments   []string
		expectation []models.BridgeType
		errored     bool
	}{
		{"finds one bridge", []string{"bridge1"}, []models.BridgeType{bt1}, false},
		{"finds multiple bridges", []string{"bridge1", "bridge2"}, []models.BridgeType{bt1, bt2}, false},
		{"errors on duplicates", []string{"bridge1", "bridge1"}, nil, true},
		{"errors on non-existent bridge names", []string{"bridge1", "doesnotexist"}, nil, true},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			bridges, err := store.FindBridgesByNames(test.arguments)
			assert.Equal(t, test.errored, err != nil)
			if test.expectation != nil {
				require.Len(t, bridges, len(test.expectation))
				for i, bridge := range test.expectation {
					bridges[i].CreatedAt = bridge.CreatedAt
					bridges[i].UpdatedAt = bridge.UpdatedAt
					assert.Equal(t, bridge, bridges[i])
				}
			}
		})
	}
}

func TestORM_PendingBridgeType_alreadyCompleted(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	_, bt := cltest.NewBridgeType(t)
	require.NoError(t, store.CreateBridgeType(bt))

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	run := cltest.NewJobRun(job)
	require.NoError(t, store.CreateJobRun(&run))

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	executor := services.NewRunExecutor(store, new(mocks.Client), keyStore, pusher)
	require.NoError(t, executor.Execute(run.ID))

	cltest.WaitForJobRunStatus(t, store, run, models.RunStatusCompleted)

	_, err := store.PendingBridgeType(run)
	assert.Error(t, err)
}

func TestORM_PendingBridgeType_success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, bt := cltest.NewBridgeType(t)
	require.NoError(t, store.CreateBridgeType(bt))

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{{Type: bt.Name}}
	assert.NoError(t, store.CreateJob(&job))

	unfinishedRun := cltest.NewJobRun(job)
	retrievedBt, err := store.PendingBridgeType(unfinishedRun)
	assert.NoError(t, err)
	retrievedBt.CreatedAt = bt.CreatedAt
	retrievedBt.UpdatedAt = bt.UpdatedAt
	assert.Equal(t, retrievedBt, *bt)
}

func TestORM_MarkRan(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	js := models.NewJob()
	require.NoError(t, store.CreateJob(&js))
	initr := models.Initiator{
		JobSpecID: js.ID,
		Type:      models.InitiatorRunAt,
		InitiatorParams: models.InitiatorParams{
			Time: models.NewAnyTime(time.Now()),
		},
	}

	require.NoError(t, store.CreateInitiator(&initr))

	assert.NoError(t, store.MarkRan(initr, true))
	ir, err := store.FindInitiator(initr.ID)
	assert.NoError(t, err)
	assert.True(t, ir.Ran)

	assert.Error(t, store.MarkRan(initr, true))
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	user1 := cltest.MustNewUser(t, "test1@email1.net", "password1")
	user2 := cltest.MustNewUser(t, "test2@email2.net", "password2")
	user2.CreatedAt = time.Now().Add(-24 * time.Hour)

	require.NoError(t, store.SaveUser(&user1))
	require.NoError(t, store.SaveUser(&user2))

	actual, err := store.FindUser()
	require.NoError(t, err)
	assert.Equal(t, user1.Email, actual.Email)
	assert.Equal(t, user1.HashedPassword, actual.HashedPassword)
}

func TestORM_AuthorizedUserWithSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		sessionID       string
		sessionDuration time.Duration
		wantError       bool
		wantEmail       string
	}{
		{"authorized", "correctID", cltest.MustParseDuration(t, "3m"), false, "have@email"},
		{"expired", "correctID", cltest.MustParseDuration(t, "0m"), true, ""},
		{"incorrect", "wrong", cltest.MustParseDuration(t, "3m"), true, ""},
		{"empty", "", cltest.MustParseDuration(t, "3m"), true, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			user := cltest.MustNewUser(t, "have@email", "password")
			require.NoError(t, store.SaveUser(&user))

			prevSession := cltest.NewSession("correctID")
			prevSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "2m"))
			require.NoError(t, store.DB.Save(&prevSession).Error)

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := store.ORM.AuthorizedUserWithSession(test.sessionID, test.sessionDuration)
			assert.Equal(t, test.wantEmail, actual.Email)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var bumpedSession models.Session
				err = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
					return db.First(&bumpedSession, "ID = ?", prevSession.ID).Error
				})
				require.NoError(t, err)
				assert.Equal(t, expectedTime[0:13], utils.ISO8601UTC(bumpedSession.LastUsed)[0:13]) // only compare up to the hour
			}
		})
	}
}

func TestORM_DeleteUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, err := store.FindUser()
	require.NoError(t, err)

	err = store.DeleteUser()
	require.NoError(t, err)

	_, err = store.FindUser()
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	session := models.NewSession()
	require.NoError(t, store.DB.Save(&session).Error)

	err := store.DeleteUserSession(session.ID)
	require.NoError(t, err)

	_, err = store.FindUser()
	require.NoError(t, err)

	sessions, err := postgres.Sessions(store.DB, 0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	initial := cltest.MustRandomUser()
	require.NoError(t, store.SaveUser(&initial))

	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"correct", initial.Email, cltest.Password, true},
		{"incorrect email", "bogus@town.org", cltest.Password, false},
		{"incorrect pwd", initial.Email, "jamaicandundada", false},
		{"incorrect both", "dudus@coke.ja", "jamaicandundada", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionRequest := models.SessionRequest{
				Email:    test.email,
				Password: test.password,
			}

			sessionID, err := store.CreateSession(sessionRequest)
			if test.wantSession {
				require.NoError(t, err)
				assert.NotEmpty(t, sessionID)
			} else {
				require.Error(t, err)
				assert.Empty(t, sessionID)
			}
		})
	}
}

func TestORM_AllSyncEvents(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	explorerClient := synchronization.NewExplorerClient(cltest.MustParseURL("http://localhost"), "", "")
	err := explorerClient.Start()
	require.NoError(t, err)
	defer explorerClient.Close()

	statsPusher := synchronization.NewStatsPusher(store.DB, explorerClient)
	require.NoError(t, statsPusher.Start())
	defer statsPusher.Close()

	// Create two events via job run callback
	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	require.NoError(t, store.ORM.CreateJob(&job))

	oldIncompleteRun := cltest.NewJobRun(job)
	oldIncompleteRun.SetStatus(models.RunStatusInProgress)
	err = store.CreateJobRun(&oldIncompleteRun)
	require.NoError(t, err)

	newCompletedRun := cltest.NewJobRun(job)
	newCompletedRun.SetStatus(models.RunStatusCompleted)
	err = store.CreateJobRun(&newCompletedRun)
	require.NoError(t, err)

	events := []models.SyncEvent{}
	err = statsPusher.AllSyncEvents(func(event models.SyncEvent) error {
		events = append(events, event)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, events, 2)
	assert.Greater(t, events[1].ID, events[0].ID)
}

func TestBulkDeleteRuns(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	var resultCount int64
	var taskCount int64
	var runCount int64

	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		job := cltest.NewJobWithWebInitiator()
		require.NoError(t, store.ORM.CreateJob(&job))

		// bulk delete should not delete these because they match the updated before
		// but none of the statuses
		oldIncompleteRun := cltest.NewJobRun(job)
		oldIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 17}`)}
		err := store.ORM.CreateJobRun(&oldIncompleteRun)
		require.NoError(t, err)
		db.Model(&oldIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete *SHOULD* delete these because they match one of the statuses
		// and the updated before
		oldCompletedRun := cltest.NewJobRun(job)
		oldCompletedRun.TaskRuns[0].Status = models.RunStatusCompleted
		oldCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 19}`)}
		oldCompletedRun.SetStatus(models.RunStatusCompleted)
		err = store.ORM.CreateJobRun(&oldCompletedRun)
		require.NoError(t, err)
		db.Model(&oldCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete should not delete these because they match one of the
		// statuses but not the updated before
		newCompletedRun := cltest.NewJobRun(job)
		newCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 23}`)}
		newCompletedRun.SetStatus(models.RunStatusCompleted)
		err = store.ORM.CreateJobRun(&newCompletedRun)
		require.NoError(t, err)
		db.Model(&newCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		// bulk delete should not delete these because none of their attributes match
		newIncompleteRun := cltest.NewJobRun(job)
		newIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 71}`)}
		newIncompleteRun.SetStatus(models.RunStatusCompleted)
		err = store.ORM.CreateJobRun(&newIncompleteRun)
		require.NoError(t, err)
		db.Model(&newIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		err = postgres.BulkDeleteRuns(store.DB, &models.BulkDeleteRunRequest{
			Status:        []models.RunStatus{models.RunStatusCompleted},
			UpdatedBefore: cltest.ParseISO8601(t, "2018-01-15T00:00:00Z"),
		})
		require.NoError(t, err)

		err = db.Model(&models.JobRun{}).Count(&runCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, int(runCount))

		err = db.Model(&models.TaskRun{}).Count(&taskCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, int(taskCount))

		err = db.Model(&models.RunResult{}).Count(&resultCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 6, int(resultCount))

		return nil
	})
	require.NoError(t, err)
}

const linkEthTxWithTaskRunQuery = `
INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES (?, ?)
`

func TestORM_RemoveUnstartedTransaction_RemoveByEthTx(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobSpec := cltest.NewJobWithRunLogInitiator()
	require.NoError(t, store.CreateJob(&jobSpec))

	runRequest := models.NewRunRequest(models.JSON{})
	require.NoError(t, store.DB.Create(runRequest).Error)
	unstartedJobRun := cltest.NewJobRun(jobSpec)
	unstartedJobRun.RunRequest = *runRequest
	unstartedJobRun.Status = models.RunStatusInProgress
	require.NoError(t, store.CreateJobRun(&unstartedJobRun))

	runRequest = models.NewRunRequest(models.JSON{})
	require.NoError(t, store.DB.Create(runRequest).Error)
	startedJobRun := cltest.NewJobRun(jobSpec)
	startedJobRun.RunRequest = *runRequest
	startedJobRun.Status = models.RunStatusInProgress
	require.NoError(t, store.CreateJobRun(&startedJobRun))

	key := cltest.MustInsertRandomKey(t, store.DB)
	ethTx := cltest.NewEthTx(t, key.Address.Address())
	require.NoError(t, store.DB.Create(&ethTx).Error)

	ethTxAttempt := cltest.NewEthTxAttempt(t, ethTx.ID)
	require.NoError(t, store.DB.Create(&ethTxAttempt).Error)
	require.NoError(t, store.DB.Exec(linkEthTxWithTaskRunQuery, unstartedJobRun.TaskRuns[0].ID, ethTx.ID).Error)

	assert.NoError(t, store.RemoveUnstartedTransactions())

	jobRuns, err := store.JobRunsFor(jobSpec.ID, 10)
	require.NoError(t, err)
	require.Len(t, jobRuns, 1, "expected only one JobRun to be left in the db")
	assert.Equal(t, models.RunStatusInProgress, jobRuns[0].Status)

	taskRuns := []models.TaskRun{}
	require.NoError(t, store.DB.Find(&taskRuns).Error)
	assert.Len(t, taskRuns, 1, "expected only one TaskRun to be left in the db")

	runRequests := []models.RunRequest{}
	require.NoError(t, store.DB.Find(&runRequests).Error)
	assert.Len(t, runRequests, 1, "expected only one RunRequests to be left in the db")

	ethTxes := []bulletprooftxmanager.EthTx{}
	require.NoError(t, store.DB.Find(&ethTxes).Error)
	assert.Len(t, ethTxes, 1, "expected only one EthTx to be left in the db")

	ethTxAttempts := []bulletprooftxmanager.EthTxAttempt{}
	require.NoError(t, store.DB.Find(&ethTxAttempts).Error)
	assert.Len(t, ethTxAttempts, 1, "expected only one EthTxAttempt to be left in the db")
}

func TestORM_RemoveUnstartedTransaction_RemoveByJobRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobSpec := cltest.NewJobWithRunLogInitiator()
	require.NoError(t, store.CreateJob(&jobSpec))

	runRequest := models.NewRunRequest(models.JSON{})
	require.NoError(t, store.DB.Create(runRequest).Error)

	unstartedJobRun := cltest.NewJobRun(jobSpec)
	unstartedJobRun.RunRequest = *runRequest
	unstartedJobRun.Status = models.RunStatusUnstarted
	require.NoError(t, store.CreateJobRun(&unstartedJobRun))

	runRequest = models.NewRunRequest(models.JSON{})
	require.NoError(t, store.DB.Create(runRequest).Error)

	startedJobRun := cltest.NewJobRun(jobSpec)
	startedJobRun.RunRequest = *runRequest
	startedJobRun.Status = models.RunStatusInProgress
	require.NoError(t, store.CreateJobRun(&startedJobRun))

	assert.NoError(t, store.RemoveUnstartedTransactions())

	jobRuns, err := store.JobRunsFor(jobSpec.ID, 10)
	require.NoError(t, err)
	require.Len(t, jobRuns, 1, "expected only one JobRun to be left in the db")
	assert.Equal(t, models.RunStatusInProgress, jobRuns[0].Status)

	taskRuns := []models.TaskRun{}
	require.NoError(t, store.DB.Find(&taskRuns).Error)
	assert.Len(t, taskRuns, 1, "expected only one TaskRun to be left in the db")

	runRequests := []models.RunRequest{}
	require.NoError(t, store.DB.Find(&runRequests).Error)
	assert.Len(t, runRequests, 1, "expected only one RunRequest to be left in the db")
}

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewEthTxAttempt(t, tx2.ID)
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, store.DB.Create(&attempt).Error)

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = bulletprooftxmanager.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, store.DB.Save(&tx3).Error)

	count, err := store.CountOf(bulletprooftxmanager.EthTx{})
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := store.EthTransactionsWithAttempts(0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
	assert.Equal(t, int64(0), *txs[1].Nonce, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].EthTxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].EthTxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].EthTxAttempts[0].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].EthTxAttempts[1].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")

	txs, count, err = store.EthTransactionsWithAttempts(0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
}

func TestORM_UpdateBridgeType(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	firstBridge := &models.BridgeType{
		Name: "UniqueName",
		URL:  cltest.WebURL(t, "http:/oneurl.com"),
	}

	require.NoError(t, store.CreateBridgeType(firstBridge))

	updateBridge := &models.BridgeTypeRequest{
		URL: cltest.WebURL(t, "http:/updatedurl.com"),
	}

	require.NoError(t, store.UpdateBridgeType(firstBridge, updateBridge))

	foundbridge, err := store.FindBridge("UniqueName")
	require.NoError(t, err)
	require.Equal(t, updateBridge.URL, foundbridge.URL)
}

func TestJobs_All(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	fmJob := cltest.NewJobWithFluxMonitorInitiator()
	runlogJob := cltest.NewJobWithRunLogInitiator()

	require.NoError(t, store.CreateJob(&fmJob))
	require.NoError(t, store.CreateJob(&runlogJob))

	var returned []*models.JobSpec
	err := store.Jobs(func(j *models.JobSpec) bool {
		// deliberately take pointer to ensure we receive new one per callback
		// checking against go gotcha:
		// https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		returned = append(returned, j)
		return true
	})
	require.NoError(t, err)
	var actual []string
	for _, j := range returned {
		actual = append(actual, j.ID.String())
	}

	var expectation []string
	for _, js := range cltest.AllJobs(t, store) {
		expectation = append(expectation, js.ID.String())
	}
	assert.ElementsMatch(t, expectation, actual)
}

func TestJobs_ScopedInitiator(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	fmJob := cltest.NewJobWithFluxMonitorInitiator()
	runlogJob := cltest.NewJobWithRunLogInitiator()
	twoInitrJob := cltest.NewJobWithFluxMonitorInitiator()
	nextinitr := cltest.NewJobWithFluxMonitorInitiator().Initiators[0]
	twoInitrJob.Initiators = append(twoInitrJob.Initiators, nextinitr)

	require.NoError(t, store.CreateJob(&fmJob))
	require.NoError(t, store.CreateJob(&runlogJob))
	require.NoError(t, store.CreateJob(&twoInitrJob))

	var actual []string
	err := store.Jobs(func(j *models.JobSpec) bool {
		actual = append(actual, j.ID.String())
		return true
	}, models.InitiatorFluxMonitor)
	require.NoError(t, err)

	expectation := []string{fmJob.ID.String(), twoInitrJob.ID.String()}
	assert.ElementsMatch(t, expectation, actual)
}

// TestJobs_SQLiteBatchSizeIntegrity verifies the BatchSize is safe for SQLite
// to handle.  Problems were experienced earlier with a size of 1001.
func TestJobs_SQLiteBatchSizeIntegrity(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	archivedJob := cltest.NewJobWithFluxMonitorInitiator()
	archivedJob.DeletedAt = gorm.DeletedAt{Valid: true, Time: time.Now()}
	require.NoError(t, store.CreateJob(&archivedJob))

	jobs := []models.JobSpec{}
	jobNumber := int(postgres.BatchSize*2 + 1)
	for i := 0; i < jobNumber; i++ {
		job := cltest.NewJobWithFluxMonitorInitiator()
		require.NoError(t, store.CreateJob(&job))
		jobs = append(jobs, job)
	}
	assert.Len(t, jobs, jobNumber)

	counter := 0
	err := store.Jobs(func(j *models.JobSpec) bool {
		counter++
		return true
	}, models.InitiatorFluxMonitor)
	require.NoError(t, err)

	assert.Equal(t, jobNumber, counter)
}

func TestORM_EthTaskRunTx(t *testing.T) {
	t.Parallel()

	// NOTE: Must sidestep transactional tests since we rely on transaction
	// rollback due to constraint violation for this function
	tc, orm, cleanup := heavyweight.FullTestORM(t, "eth_task_run_transactions", true, true)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(t, tc)
	store.ORM = orm
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	sharedTaskRunID, _ := cltest.MustInsertTaskRun(t, store)

	t.Run("creates eth_task_run_transaction and eth_tx", func(t *testing.T) {
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{0, 1, 2}
		gasLimit := uint64(42)

		err := store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: sharedTaskRunID}, fromAddress, toAddress, encodedPayload, gasLimit)
		require.NoError(t, err)

		etrt, err := store.FindEthTaskRunTxByTaskRunID(sharedTaskRunID)
		require.NoError(t, err)

		assert.Equal(t, sharedTaskRunID, etrt.TaskRunID)
		require.NotNil(t, etrt.EthTx)
		assert.Nil(t, etrt.EthTx.Nonce)
		assert.Equal(t, fromAddress, etrt.EthTx.FromAddress)
		assert.Equal(t, toAddress, etrt.EthTx.ToAddress)
		assert.Equal(t, encodedPayload, etrt.EthTx.EncodedPayload)
		assert.Equal(t, gasLimit, etrt.EthTx.GasLimit)
		assert.Equal(t, bulletprooftxmanager.EthTxUnstarted, etrt.EthTx.State)

		// Do it again to test idempotence
		err = store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: sharedTaskRunID}, fromAddress, toAddress, encodedPayload, gasLimit)
		require.NoError(t, err)

		// Ensure it didn't leave a stray EthTx hanging around
		store.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
			var count int64
			require.NoError(t, db.Table("eth_txes").Count(&count).Error)
			assert.Equal(t, 1, int(count))
			return nil
		})
	})

	t.Run("returns error if eth_task_run_transaction already exists with this task run ID but has different values", func(t *testing.T) {
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{3, 2, 1}
		gasLimit := uint64(24)

		err := store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: sharedTaskRunID}, fromAddress, toAddress, encodedPayload, gasLimit)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transaction already exists for task run ID")
	})

	t.Run("does not return error on re-insert if only the gas limit changed", func(t *testing.T) {
		taskRunID, _ := cltest.MustInsertTaskRun(t, store)
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{0, 1, 2}
		firstGasLimit := uint64(42)

		// First insert
		err := store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: taskRunID}, fromAddress, toAddress, encodedPayload, firstGasLimit)
		require.NoError(t, err)

		secondGasLimit := uint64(99)

		// Second insert
		err = store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: taskRunID}, fromAddress, toAddress, encodedPayload, secondGasLimit)
		require.NoError(t, err)

		etrt, err := store.FindEthTaskRunTxByTaskRunID(taskRunID)
		require.NoError(t, err)

		// But the second insert did not change the gas limit
		assert.Equal(t, firstGasLimit, etrt.EthTx.GasLimit)
	})

	t.Run("returns error if fromAddress does not correspond to a key", func(t *testing.T) {
		taskRunID, _ := cltest.MustInsertTaskRun(t, store)
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{0, 1, 2}
		gasLimit := uint64(42)

		err := store.IdempotentInsertEthTaskRunTx(models.EthTxMeta{TaskRunID: taskRunID}, cltest.NewAddress(), toAddress, encodedPayload, gasLimit)
		assert.Error(t, err)
		assert.EqualError(t, err, "ERROR: insert or update on table \"eth_txes\" violates foreign key constraint \"eth_txes_from_address_fkey\" (SQLSTATE 23503)")
	})
}

func TestORM_FindJobWithErrorsPreloadsJobSpecErrors(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job1))
	job2 := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job2))

	description1, description2 := "description 1", "description 2"

	store.UpsertErrorFor(job1.ID, description1)
	store.UpsertErrorFor(job1.ID, description2)

	job1, err := store.FindJobWithErrors(job1.ID)
	require.NoError(t, err)
	job2, err = store.FindJobWithErrors(job2.ID)
	require.NoError(t, err)

	assert.Len(t, job1.Errors, 2)
	assert.Len(t, job2.Errors, 0)

	assert.Equal(t, job1.Errors[0].Description, description1)
	assert.Equal(t, job1.Errors[1].Description, description2)
}

func TestORM_UpsertErrorFor_Happy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job1 := cltest.NewJob()
	job2 := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job1))
	require.NoError(t, store.CreateJob(&job2))

	description1, description2 := "description 1", "description 2"

	store.UpsertErrorFor(job1.ID, description1)

	tests := []struct {
		jobID               models.JobID
		description         string
		expectedOccurrences uint
	}{
		{
			job1.ID,
			description1,
			2, // duplicate
		},
		{
			job1.ID,
			description2,
			1,
		},
		{
			job2.ID,
			description1,
			1,
		},
		{
			job2.ID,
			description2,
			1,
		},
	}

	for _, tt := range tests {
		test := tt
		testName := fmt.Sprintf(`Create JobSpecError with ID %v and description "%s"`, test.jobID, test.description)
		t.Run(testName, func(t *testing.T) {
			store.UpsertErrorFor(test.jobID, test.description)
			jse, err := store.FindJobSpecError(test.jobID, test.description)
			require.NoError(t, err)
			require.Equal(t, test.expectedOccurrences, jse.Occurrences)
			if test.expectedOccurrences > 1 {
				require.True(t, jse.CreatedAt.Before(jse.UpdatedAt))
			} else {
				require.Equal(t, jse.CreatedAt, jse.UpdatedAt)
			}
		})
	}
}

func TestORM_UpsertErrorFor_Error(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJob()
	require.NoError(t, store.CreateJob(&job))
	description := "description"
	store.UpsertErrorFor(job.ID, description)

	tests := []struct {
		name        string
		jobID       models.JobID
		description string
	}{
		{
			"missing job",
			models.NewJobID(),
			description,
		},
		{
			"missing description",
			job.ID,
			"",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store.UpsertErrorFor(test.jobID, test.description)
		})
	}
}

func TestORM_FindOrCreateFluxMonitorRoundStats(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := cltest.NewAddress()
	var roundID uint32 = 1

	fmrs, err := store.FindOrCreateFluxMonitorRoundStats(address, roundID)
	require.NoError(t, err)
	require.Equal(t, roundID, fmrs.RoundID)
	require.Equal(t, address, fmrs.Aggregator)

	count, err := store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 1, count)

	fmrs, err = store.FindOrCreateFluxMonitorRoundStats(address, roundID)
	require.NoError(t, err)
	require.Equal(t, roundID, fmrs.RoundID)
	require.Equal(t, address, fmrs.Aggregator)

	count, err = store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestORM_DeleteFluxMonitorRoundsBackThrough(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := cltest.NewAddress()

	for round := uint32(0); round < 10; round++ {
		_, err := store.FindOrCreateFluxMonitorRoundStats(address, round)
		require.NoError(t, err)
	}

	count, err := store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 10, count)

	err = store.DeleteFluxMonitorRoundsBackThrough(cltest.NewAddress(), 5)
	require.NoError(t, err)

	count, err = store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 10, count)

	err = store.DeleteFluxMonitorRoundsBackThrough(address, 5)
	require.NoError(t, err)

	count, err = store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 5, count)
}

func TestORM_MostRecentFluxMonitorRoundID(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := cltest.NewAddress()

	for round := uint32(0); round < 10; round++ {
		_, err := store.FindOrCreateFluxMonitorRoundStats(address, round)
		require.NoError(t, err)
	}

	count, err := store.ORM.CountOf(&models.FluxMonitorRoundStats{})
	require.NoError(t, err)
	require.Equal(t, 10, count)

	roundID, err := store.MostRecentFluxMonitorRoundID(cltest.NewAddress())
	require.Error(t, err)
	require.Equal(t, uint32(0), roundID)

	roundID, err = store.MostRecentFluxMonitorRoundID(address)
	require.NoError(t, err)
	require.Equal(t, uint32(9), roundID)
}

func TestORM_UpdateFluxMonitorRoundStats(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	address := cltest.NewAddress()
	var roundID uint32 = 1
	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	for expectedCount := uint64(1); expectedCount < 4; expectedCount++ {
		jobRun := cltest.NewJobRun(job)
		require.NoError(t, store.CreateJobRun(&jobRun))
		err := store.UpdateFluxMonitorRoundStats(address, roundID, jobRun.ID)
		require.NoError(t, err)
		fmrs, err := store.FindOrCreateFluxMonitorRoundStats(address, roundID)
		require.NoError(t, err)
		require.Equal(t, expectedCount, fmrs.NumSubmissions)
		require.True(t, fmrs.JobRunID.Valid)
		require.Equal(t, jobRun.ID, fmrs.JobRunID.UUID)
	}
}
