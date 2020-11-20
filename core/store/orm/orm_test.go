package orm_test

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestORM_AllNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobs := cltest.AllJobs(t, store)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORM_CreateJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * *")
	store.CreateJob(&j1)

	j2, err := store.FindJob(j1.ID)
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
	err := orm.RawDB(func(db *gorm.DB) error {
		require.NoError(t, orm.CreateJob(&job))
		require.NoError(t, db.Delete(&job).Error)
		require.Error(t, db.First(&job).Error)
		err := store.ORM.Unscoped().RawDB(func(db *gorm.DB) error {
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
		models.TaskSpec{Type: models.MustNewTaskType("task1")},
		models.TaskSpec{Type: models.MustNewTaskType("task2")},
		models.TaskSpec{Type: models.MustNewTaskType("task3")},
		models.TaskSpec{Type: models.MustNewTaskType("task4")},
	}
	assert.NoError(t, store.CreateJob(&job))

	orm := store.ORM
	retrievedJob, err := orm.FindJob(job.ID)
	assert.NoError(t, err)
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
	require.Equal(t, store.CreateExternalInitiator(exi).Error(), `pq: duplicate key value violates unique constraint "external_initiators_name_key"`)
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

	require.Error(t, utils.JustError(store.FindJob(job.ID)))
	require.Error(t, utils.JustError(store.FindJobRun(run.ID)))

	orm := store.ORM.Unscoped()
	require.NoError(t, utils.JustError(orm.FindJob(job.ID)))
	require.NoError(t, utils.JustError(orm.FindJobRun(run.ID)))
}

func TestORM_CreateJobRun_CreatesRunRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	rr := models.NewRunRequest(models.JSON{})
	currentHeight := big.NewInt(0)
	run, _ := services.NewRun(&job, &job.Initiators[0], currentHeight, rr, store.Config, store.ORM, time.Now())
	require.NoError(t, store.CreateJobRun(run))

	requestCount, err := store.ORM.CountOf(&models.RunRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}

func TestORM_SaveJobRun_OnConstraintViolationOtherThanOptimisticLockFailureReturnsError(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jr := cltest.CreateJobRunWithStatus(t, store, job, models.RunStatusUnstarted)

	jr.InitiatorID = 0
	jr.Initiator = models.Initiator{}
	err := store.SaveJobRun(&jr)
	assert.EqualError(t, err, "pq: insert or update on table \"job_runs\" violates foreign key constraint \"fk_job_runs_initiator_id\"")
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

	// Save the updated at before saving with cancelled
	updatedAt := jr.UpdatedAt

	jr.SetStatus(models.RunStatusCancelled)
	require.NoError(t, store.SaveJobRun(&jr))

	// Restore the previous updated at to simulate a conflict
	jr.UpdatedAt = updatedAt
	jr.SetStatus(models.RunStatusInProgress)
	assert.Equal(t, orm.ErrOptimisticUpdateConflict, store.SaveJobRun(&jr))
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
	actual := []*models.ID{runs[0].ID, runs[1].ID, runs[2].ID}
	assert.Equal(t, []*models.ID{jr2.ID, jr1.ID, jr3.ID}, actual)

	limRuns, limErr := store.JobRunsFor(job.ID, 2)
	assert.NoError(t, limErr)
	limActual := []*models.ID{limRuns[0].ID, limRuns[1].ID}
	assert.Equal(t, []*models.ID{jr2.ID, jr1.ID}, limActual)

	_, limZeroErr := store.JobRunsFor(job.ID, 0)
	assert.NoError(t, limZeroErr)
	limZeroActual := []*models.ID{}
	assert.Equal(t, []*models.ID{}, limZeroActual)
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
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := cltest.NewJobRun(includedJob)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, store.CreateJobRun(&jr2))

	excludedJobRun := cltest.NewJobRun(excludedJob)
	excludedJobRun.CreatedAt = time.Now().AddDate(0, 0, -9)
	require.NoError(t, store.CreateJobRun(&excludedJobRun))

	runs, count, err := store.JobRunsSortedFor(includedJob.ID, orm.Descending, 0, 100)
	assert.NoError(t, err)
	require.Equal(t, 2, count)
	actual := []*models.ID{runs[0].ID, runs[1].ID} // doesn't include excludedJobRun
	assert.Equal(t, []*models.ID{jr2.ID, jr1.ID}, actual)
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

	var seedIds []*models.ID
	for _, status := range statuses {
		run := cltest.NewJobRun(j)
		run.SetStatus(status)
		require.NoError(t, store.CreateJobRun(&run))
		seedIds = append(seedIds, run.ID)
	}

	tests := []struct {
		name     string
		statuses []models.RunStatus
		expected []*models.ID
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]*models.ID{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingIncomingConfirmations, models.RunStatusPendingOutgoingConfirmations},
			[]*models.ID{seedIds[0], seedIds[1], seedIds[2]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []*models.ID{}
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

	var seedIds []*models.ID
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
		expected []*models.ID
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]*models.ID{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{
				models.RunStatusPendingBridge,
				models.RunStatusPendingOutgoingConfirmations,
				models.RunStatusPendingIncomingConfirmations,
				models.RunStatusPendingConnection},
			[]*models.ID{seedIds[0], seedIds[1], seedIds[2], seedIds[3]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []*models.ID{}
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
	js.Tasks = []models.TaskSpec{models.TaskSpec{Type: models.MustNewTaskType("bridgetestname")}}
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

	_, bt := cltest.NewBridgeType(t)
	require.NoError(t, store.CreateBridgeType(bt))

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	run := cltest.NewJobRun(job)
	require.NoError(t, store.CreateJobRun(&run))

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	executor := services.NewRunExecutor(store, pusher)
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
	job.Tasks = []models.TaskSpec{models.TaskSpec{Type: bt.Name}}
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
			require.NoError(t, store.SaveSession(&prevSession))

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := store.ORM.AuthorizedUserWithSession(test.sessionID, test.sessionDuration)
			assert.Equal(t, test.wantEmail, actual.Email)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var bumpedSession models.Session
				err = store.ORM.RawDB(func(db *gorm.DB) error {
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

	_, err = store.DeleteUser()
	require.NoError(t, err)

	_, err = store.FindUser()
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	session := models.NewSession()
	require.NoError(t, store.SaveSession(&session))

	err := store.DeleteUserSession(session.ID)
	require.NoError(t, err)

	_, err = store.FindUser()
	require.NoError(t, err)

	sessions, err := store.Sessions(0, 10)
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

	orm := store.ORM
	statsPusher := synchronization.NewStatsPusher(orm, explorerClient)
	require.NoError(t, statsPusher.Start())
	defer statsPusher.Close()

	// Create two events via job run callback
	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	require.NoError(t, store.ORM.CreateJob(&job))

	oldIncompleteRun := cltest.NewJobRun(job)
	oldIncompleteRun.SetStatus(models.RunStatusInProgress)
	err = orm.CreateJobRun(&oldIncompleteRun)
	require.NoError(t, err)

	newCompletedRun := cltest.NewJobRun(job)
	newCompletedRun.SetStatus(models.RunStatusCompleted)
	err = orm.CreateJobRun(&newCompletedRun)
	require.NoError(t, err)

	events := []models.SyncEvent{}
	err = orm.AllSyncEvents(func(event models.SyncEvent) error {
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

	var resultCount int
	var taskCount int
	var runCount int
	orm := store.ORM

	err := orm.RawDB(func(db *gorm.DB) error {
		job := cltest.NewJobWithWebInitiator()
		require.NoError(t, store.ORM.CreateJob(&job))

		// bulk delete should not delete these because they match the updated before
		// but none of the statuses
		oldIncompleteRun := cltest.NewJobRun(job)
		oldIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 17}`)}
		err := orm.CreateJobRun(&oldIncompleteRun)
		require.NoError(t, err)
		db.Model(&oldIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete *SHOULD* delete these because they match one of the statuses
		// and the updated before
		oldCompletedRun := cltest.NewJobRun(job)
		oldCompletedRun.TaskRuns[0].Status = models.RunStatusCompleted
		oldCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 19}`)}
		oldCompletedRun.SetStatus(models.RunStatusCompleted)
		err = orm.CreateJobRun(&oldCompletedRun)
		require.NoError(t, err)
		db.Model(&oldCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete should not delete these because they match one of the
		// statuses but not the updated before
		newCompletedRun := cltest.NewJobRun(job)
		newCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 23}`)}
		newCompletedRun.SetStatus(models.RunStatusCompleted)
		err = orm.CreateJobRun(&newCompletedRun)
		require.NoError(t, err)
		db.Model(&newCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		// bulk delete should not delete these because none of their attributes match
		newIncompleteRun := cltest.NewJobRun(job)
		newIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 71}`)}
		newIncompleteRun.SetStatus(models.RunStatusCompleted)
		err = orm.CreateJobRun(&newIncompleteRun)
		require.NoError(t, err)
		db.Model(&newIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		err = store.ORM.BulkDeleteRuns(&models.BulkDeleteRunRequest{
			Status:        []models.RunStatus{models.RunStatusCompleted},
			UpdatedBefore: cltest.ParseISO8601(t, "2018-01-15T00:00:00Z"),
		})
		require.NoError(t, err)

		err = db.Model(&models.JobRun{}).Count(&runCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, runCount)

		err = db.Model(&models.TaskRun{}).Count(&taskCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, taskCount)

		err = db.Model(&models.RunResult{}).Count(&resultCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, resultCount)

		return nil
	})
	require.NoError(t, err)
}

func TestORM_KeysOrdersByCreatedAtAsc(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := store.ORM

	testJSON := cltest.JSONFromString(t, "{}")

	earlierAddress := cltest.DefaultKeyAddressEIP55
	earlier := models.Key{Address: earlierAddress, JSON: testJSON}

	require.NoError(t, orm.CreateKeyIfNotExists(earlier))
	time.Sleep(10 * time.Millisecond)

	laterAddress, err := models.NewEIP55Address("0xBB68588621f7E847070F4cC9B9e70069BA55FC5A")
	require.NoError(t, err)
	later := models.Key{Address: laterAddress, JSON: testJSON}

	require.NoError(t, orm.CreateKeyIfNotExists(later))

	keys, err := store.SendKeys()
	require.NoError(t, err)

	require.Len(t, keys, 2)

	assert.Equal(t, keys[0].Address, earlierAddress)
	assert.Equal(t, keys[1].Address, laterAddress)
}

func TestORM_SendKeys(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := store.ORM

	testJSON := cltest.JSONFromString(t, "{}")

	sendingAddress := cltest.DefaultKeyAddressEIP55
	sending := models.Key{Address: sendingAddress, JSON: testJSON}

	require.NoError(t, orm.CreateKeyIfNotExists(sending))
	time.Sleep(10 * time.Millisecond)

	fundingAddress, err := models.NewEIP55Address("0xBB68588621f7E847070F4cC9B9e70069BA55FC5A")
	require.NoError(t, err)
	funding := models.Key{Address: fundingAddress, JSON: testJSON, IsFunding: true}

	require.NoError(t, orm.CreateKeyIfNotExists(funding))

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	keys, err = store.SendKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
}

func TestORM_SyncDbKeyStoreToDisk(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := store.ORM
	require.NoError(t, store.KeyStore.Unlock(cltest.Password))

	keysDir := store.Config.KeysDir()
	// Clear out the fixture
	require.NoError(t, os.RemoveAll(keysDir))
	require.NoError(t, store.DeleteKey(cltest.DefaultKeyAddress[:]))
	// Fixture key is deleted
	dbkeys, err := store.SendKeys()
	require.NoError(t, err)
	require.Len(t, dbkeys, 0)

	seed, err := models.NewKeyFromFile(fmt.Sprintf("../../internal/fixtures/keys/%s", cltest.DefaultKeyFixtureFileName))
	require.NoError(t, err)
	require.NoError(t, orm.CreateKeyIfNotExists(seed))

	require.True(t, isDirEmpty(t, keysDir))
	err = orm.ClobberDiskKeyStoreWithDBKeys(keysDir)
	require.NoError(t, err)

	dbkeys, err = store.SendKeys()
	require.NoError(t, err)
	require.Len(t, dbkeys, 1)

	diskkeys, err := utils.FilesInDir(keysDir)
	require.NoError(t, err)
	require.Len(t, diskkeys, 1)

	key := dbkeys[0]
	content, err := utils.FileContents(filepath.Join(keysDir, diskkeys[0]))
	require.NoError(t, err)
	assert.Equal(t, key.JSON.String(), content)
}

const linkEthTxWithTaskRunQuery = `
INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)
`

func TestORM_RemoveUnstartedTransaction(t *testing.T) {
	storeInstance, cleanup := cltest.NewStore(t)
	defer cleanup()
	ormInstance := storeInstance.ORM

	jobSpec := cltest.NewJobWithRunLogInitiator()
	require.NoError(t, storeInstance.CreateJob(&jobSpec))

	for _, status := range []models.RunStatus{
		"in_progress",
		"unstarted",
	} {
		jobRun := cltest.NewJobRun(jobSpec)
		jobRun.Status = status
		jobRun.TaskRuns = []models.TaskRun{
			{
				ID:         models.NewID(),
				Status:     models.RunStatusUnstarted,
				TaskSpecID: jobSpec.Tasks[0].ID,
			},
		}
		runRequest := models.NewRunRequest(models.JSON{})
		require.NoError(t, storeInstance.DB.Create(&runRequest).Error)
		jobRun.RunRequest = *runRequest
		require.NoError(t, storeInstance.CreateJobRun(&jobRun))

		key := cltest.MustInsertRandomKey(t, storeInstance)
		ethTx := cltest.NewEthTx(t, storeInstance, key.Address.Address())
		ethTx.State = models.EthTxState(status)
		if status == "in_progress" {
			var nonce int64 = 1
			ethTx.Nonce = &nonce
		}
		require.NoError(t, storeInstance.DB.Save(&ethTx).Error)

		ethTxAttempt := cltest.NewEthTxAttempt(t, ethTx.ID)
		ethTxAttempt.State = models.EthTxAttemptInProgress
		require.NoError(t, storeInstance.DB.Save(&ethTxAttempt).Error)

		require.NoError(t, storeInstance.DB.Exec(linkEthTxWithTaskRunQuery, jobRun.TaskRuns[0].ID.UUID(), ethTx.ID).Error)
	}

	assert.NoError(t, ormInstance.RemoveUnstartedTransactions())

	jobRuns, err := ormInstance.JobRunsFor(jobSpec.ID, 10)
	assert.NoError(t, err)
	assert.Len(t, jobRuns, 1, "expected only one JobRun to be left in the db")
	assert.Equal(t, jobRuns[0].Status, models.RunStatusInProgress)

	taskRuns := []models.TaskRun{}
	assert.NoError(t, storeInstance.DB.Find(&taskRuns).Error)
	assert.Len(t, taskRuns, 1, "expected only one TaskRun to be left in the db")

	runRequests := []models.RunRequest{}
	assert.NoError(t, storeInstance.DB.Find(&runRequests).Error)
	assert.Len(t, runRequests, 1, "expected only one RunRequest to be left in the db")

	ethTxes := []models.EthTx{}
	assert.NoError(t, storeInstance.DB.Find(&ethTxes).Error)
	assert.Len(t, ethTxes, 1, "expected only one EthTx to be left in the db")

	ethTxAttempts := []models.EthTxAttempt{}
	assert.NoError(t, storeInstance.DB.Find(&ethTxAttempts).Error)
	assert.Len(t, ethTxAttempts, 1, "expected only one EthTxAttempt to be left in the db")
}

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	from := cltest.DefaultKeyAddress
	cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewEthTxAttempt(t, tx2.ID)
	attempt.State = models.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, store.DB.Create(&attempt).Error)

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, store, from)
	tx3.State = models.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, store.DB.Save(&tx3).Error)

	count, err := store.CountOf(models.EthTx{})
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

func isDirEmpty(t *testing.T, dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		t.Fatal(err)
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		return true
	}

	return false
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
	archivedJob.DeletedAt = cltest.NullableTime(time.Now())
	require.NoError(t, store.CreateJob(&archivedJob))

	jobs := []models.JobSpec{}
	jobNumber := orm.BatchSize*2 + 1
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

func TestORM_Heads_Chain(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// A competing chain existed from block num 3 to 4
	var baseOfForkHash common.Hash
	var longestChainHeadHash common.Hash
	var parentHash *common.Hash
	for idx := 0; idx < 8; idx++ {
		h := *cltest.Head(idx)
		if parentHash != nil {
			h.ParentHash = *parentHash
		}
		parentHash = &h.Hash
		if idx == 2 {
			baseOfForkHash = h.Hash
		} else if idx == 7 {
			longestChainHeadHash = h.Hash
		}
		assert.Nil(t, store.IdempotentInsertHead(h))
	}

	competingHead1 := *cltest.Head(3)
	competingHead1.ParentHash = baseOfForkHash
	assert.Nil(t, store.IdempotentInsertHead(competingHead1))
	competingHead2 := *cltest.Head(4)
	competingHead2.ParentHash = competingHead1.Hash
	assert.Nil(t, store.IdempotentInsertHead(competingHead2))

	// Query for the top of the longer chain does not include the competing chain
	h, err := store.Chain(longestChainHeadHash, 12)
	require.NoError(t, err)
	assert.Equal(t, longestChainHeadHash, h.Hash)
	count := 1
	for {
		if h.Parent == nil {
			break
		}
		require.NotEqual(t, competingHead1.Hash, h.Hash)
		require.NotEqual(t, competingHead2.Hash, h.Hash)
		h = *h.Parent
		count++
	}
	assert.Equal(t, 8, count)

	// If we set the limit lower we get fewer heads in chain
	h, err = store.Chain(longestChainHeadHash, 2)
	require.NoError(t, err)
	assert.Equal(t, longestChainHeadHash, h.Hash)
	count = 1
	for {
		if h.Parent == nil {
			break
		}
		h = *h.Parent
		count++
	}
	assert.Equal(t, 2, count)

	// If we query for the top of the competing chain we get its parents
	head, err := store.Chain(competingHead2.Hash, 12)
	require.NoError(t, err)
	assert.Equal(t, competingHead2.Hash, head.Hash)
	require.NotNil(t, head.Parent)
	assert.Equal(t, competingHead1.Hash, head.Parent.Hash)
	require.NotNil(t, head.Parent.Parent)
	assert.Equal(t, baseOfForkHash, head.Parent.Parent.Hash)
	assert.NotNil(t, head.Parent.Parent.Parent) // etc...

	// Returns error if hash has no matches
	_, err = store.Chain(cltest.NewHash(), 12)
	require.Error(t, err)
}

func TestORM_Heads_IdempotentInsertHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Returns nil when inserting first head
	head := *cltest.Head(0)
	require.NoError(t, store.IdempotentInsertHead(head))

	// Head is inserted
	foundHead, err := store.LastHead()
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)

	// Returns nil when inserting same head again
	require.NoError(t, store.IdempotentInsertHead(head))

	// Head is still inserted
	foundHead, err = store.LastHead()
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)
}

func TestORM_EthTaskRunTx(t *testing.T) {
	t.Parallel()

	// NOTE: Must sidestep transactional tests since we rely on transaction
	// rollback due to constraint violation for this function
	tc, orm, cleanup := cltest.BootstrapThrowawayORM(t, "eth_task_run_transactions", true, true)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(tc)
	store.ORM = orm
	defer cleanup()

	sharedTaskRunID := cltest.MustInsertTaskRun(t, store)
	keys, err := orm.SendKeys()
	require.NoError(t, err)
	fromAddress := keys[0].Address.Address()

	t.Run("creates eth_task_run_transaction and eth_tx", func(t *testing.T) {
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{0, 1, 2}
		gasLimit := uint64(42)

		err := store.IdempotentInsertEthTaskRunTx(sharedTaskRunID, fromAddress, toAddress, encodedPayload, gasLimit)
		require.NoError(t, err)

		etrt, err := store.FindEthTaskRunTxByTaskRunID(sharedTaskRunID.UUID())
		require.NoError(t, err)

		assert.Equal(t, sharedTaskRunID.UUID(), etrt.TaskRunID)
		require.NotNil(t, etrt.EthTx)
		assert.Nil(t, etrt.EthTx.Nonce)
		assert.Equal(t, fromAddress, etrt.EthTx.FromAddress)
		assert.Equal(t, toAddress, etrt.EthTx.ToAddress)
		assert.Equal(t, encodedPayload, etrt.EthTx.EncodedPayload)
		assert.Equal(t, gasLimit, etrt.EthTx.GasLimit)
		assert.Equal(t, models.EthTxUnstarted, etrt.EthTx.State)

		// Do it again to test idempotence
		err = store.IdempotentInsertEthTaskRunTx(sharedTaskRunID, fromAddress, toAddress, encodedPayload, gasLimit)
		require.NoError(t, err)

		// Ensure it didn't leave a stray EthTx hanging around
		store.RawDB(func(db *gorm.DB) error {
			var count int
			require.NoError(t, db.Table("eth_txes").Count(&count).Error)
			assert.Equal(t, 1, count)
			return nil
		})
	})

	t.Run("returns error if eth_task_run_transaction already exists with this task run ID but has different values", func(t *testing.T) {
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{3, 2, 1}
		gasLimit := uint64(24)

		err := store.IdempotentInsertEthTaskRunTx(sharedTaskRunID, fromAddress, toAddress, encodedPayload, gasLimit)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transaction already exists for task run ID")
	})

	t.Run("does not return error on re-insert if only the gas limit changed", func(t *testing.T) {
		taskRunID := cltest.MustInsertTaskRun(t, store)
		toAddress := cltest.NewAddress()
		encodedPayload := []byte{0, 1, 2}
		firstGasLimit := uint64(42)

		// First insert
		err := store.IdempotentInsertEthTaskRunTx(taskRunID, fromAddress, toAddress, encodedPayload, firstGasLimit)
		require.NoError(t, err)

		secondGasLimit := uint64(99)

		// Second insert
		err = store.IdempotentInsertEthTaskRunTx(taskRunID, fromAddress, toAddress, encodedPayload, secondGasLimit)
		require.NoError(t, err)

		etrt, err := store.FindEthTaskRunTxByTaskRunID(taskRunID.UUID())
		require.NoError(t, err)

		// But the second insert did not change the gas limit
		assert.Equal(t, firstGasLimit, etrt.EthTx.GasLimit)
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
		jobID               *models.ID
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
		jobID       *models.ID
		description string
	}{
		{
			"missing job",
			models.NewID(),
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
		require.Equal(t, jobRun.ID, fmrs.JobRunID)
	}
}

func TestORM_GetRoundRobinAddress(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	fundingKey := models.Key{Address: models.EIP55Address(cltest.NewAddress().Hex()), JSON: cltest.JSONFromString(t, `{"key": 2}`), IsFunding: true}
	k0Address := cltest.DefaultKey
	k1 := models.Key{Address: models.EIP55Address(cltest.NewAddress().Hex()), JSON: cltest.JSONFromString(t, `{"key": 1}`)}
	k2 := models.Key{Address: models.EIP55Address(cltest.NewAddress().Hex()), JSON: cltest.JSONFromString(t, `{"key": 2}`)}

	require.NoError(t, store.CreateKeyIfNotExists(fundingKey))
	require.NoError(t, store.CreateKeyIfNotExists(k1))
	require.NoError(t, store.CreateKeyIfNotExists(k2))

	t.Run("with no address filter, rotates between all addresses", func(t *testing.T) {
		address, err := store.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k0Address, address.Hex())

		address, err = store.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k1.Address.Hex(), address.Hex())

		address, err = store.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k2.Address.Hex(), address.Hex())

		address, err = store.GetRoundRobinAddress()
		require.NoError(t, err)
		assert.Equal(t, k0Address, address.Hex())
	})

	t.Run("with address filter, rotates between given addresses", func(t *testing.T) {
		addresses := []common.Address{k1.Address.Address(), k2.Address.Address()}

		address, err := store.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k1.Address.Hex(), address.Hex())

		address, err = store.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k2.Address.Hex(), address.Hex())

		address, err = store.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k1.Address.Hex(), address.Hex())

		address, err = store.GetRoundRobinAddress(addresses...)
		require.NoError(t, err)
		assert.Equal(t, k2.Address.Hex(), address.Hex())
	})

	t.Run("with address filter when no address matches", func(t *testing.T) {
		_, err := store.GetRoundRobinAddress([]common.Address{cltest.NewAddress()}...)
		require.Error(t, err)
		require.Equal(t, "no keys available", err.Error())
	})
}

func TestORM_MarkLogConsumed(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := store.ORM

	blockHash := cltest.NewHash()
	logIndex := uint(42)
	job := cltest.MustInsertJobSpec(t, store)
	blockNumber := uint64(142)

	require.NoError(t, orm.MarkLogConsumed(blockHash, logIndex, job.ID, blockNumber))

	res, err := orm.DB.DB().Exec(`SELECT * FROM log_consumptions;`)
	require.NoError(t, err)
	rowsaffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsaffected)
}
