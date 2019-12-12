package orm_test

import (
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/assets"
	"chainlink/core/internal/cltest"
	"chainlink/core/services"
	"chainlink/core/services/synchronization"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestORM_WhereNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j1 := models.NewJob()
	jobs := []models.JobSpec{j1}

	err := store.Where("ID", models.NewID().String(), &jobs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

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

func TestORM_ArchiveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithSchedule("* * * * *")
	require.NoError(t, store.CreateJob(&job))

	init := job.Initiators[0]
	run := job.NewRun(init)
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

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	requestCount, err := store.ORM.CountOf(&models.RunRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}

func TestORM_SaveJobRun_ArchivedDoesNotRevertDeletedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	require.NoError(t, store.ArchiveJob(job.ID))

	jr.Status = models.RunStatusInProgress
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

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	jr.Status = models.RunStatusInProgress
	require.NoError(t, store.SaveJobRun(&jr))

	// Save the updated at before saving with cancelled
	updatedAt := jr.UpdatedAt

	jr.Status = models.RunStatusCancelled
	require.NoError(t, store.SaveJobRun(&jr))

	// Restore the previous updated at to simulate a conflict
	jr.UpdatedAt = updatedAt
	jr.Status = models.RunStatusInProgress
	assert.Equal(t, orm.OptimisticUpdateConflictError, store.SaveJobRun(&jr))
}

func TestORM_JobRunsFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	i := job.Initiators[0]
	jr1 := job.NewRun(i)
	jr1.CreatedAt = time.Now().AddDate(0, 0, -1)
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := job.NewRun(i)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, store.CreateJobRun(&jr2))
	jr3 := job.NewRun(i)
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

	initr := job.Initiators[0]
	jr1 := job.NewRun(initr)
	jr1.Status = models.RunStatusCompleted
	jr1.Payment = assets.NewLink(2)
	jr1.FinishedAt = null.TimeFrom(time.Now())
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := job.NewRun(initr)
	jr2.Status = models.RunStatusCompleted
	jr2.Payment = assets.NewLink(3)
	jr2.FinishedAt = null.TimeFrom(time.Now())
	require.NoError(t, store.CreateJobRun(&jr2))
	jr3 := job.NewRun(initr)
	jr3.Status = models.RunStatusCompleted
	jr3.Payment = assets.NewLink(5)
	jr3.FinishedAt = null.TimeFrom(time.Now())
	require.NoError(t, store.CreateJobRun(&jr3))
	jr4 := job.NewRun(initr)
	jr4.Status = models.RunStatusCompleted
	jr4.Payment = assets.NewLink(5)
	require.NoError(t, store.CreateJobRun(&jr4))
	jr5 := job.NewRun(initr)
	jr5.Status = models.RunStatusCancelled
	jr5.Payment = assets.NewLink(5)
	jr5.FinishedAt = null.TimeFrom(time.Now())
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

	i := includedJob.Initiators[0]
	jr1 := includedJob.NewRun(i)
	jr1.CreatedAt = time.Now().AddDate(0, 0, -1)
	require.NoError(t, store.CreateJobRun(&jr1))
	jr2 := includedJob.NewRun(i)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, store.CreateJobRun(&jr2))

	excludedJobRun := excludedJob.NewRun(excludedJob.Initiators[0])
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
	i := j.Initiators[0]
	npr := j.NewRun(i)
	require.NoError(t, store.CreateJobRun(&npr))

	statuses := []models.RunStatus{
		models.RunStatusPendingBridge,
		models.RunStatusPendingConfirmations,
		models.RunStatusCompleted}

	var seedIds []*models.ID
	for _, status := range statuses {
		run := j.NewRun(i)
		run.Status = status
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
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingConfirmations},
			[]*models.ID{seedIds[0], seedIds[1]},
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
	i := j.Initiators[0]
	npr := j.NewRun(i)
	require.NoError(t, store.CreateJobRun(&npr))

	statuses := []models.RunStatus{
		models.RunStatusPendingBridge,
		models.RunStatusPendingConfirmations,
		models.RunStatusPendingConnection,
		models.RunStatusCompleted}

	var seedIds []*models.ID
	for _, status := range statuses {
		run := j.NewRun(i)
		run.Status = status
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
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingConfirmations, models.RunStatusPendingConnection},
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

func TestORM_UnscopedJobRunsWithStatus_OrdersByCreatedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.CreateJob(&j))
	i := j.Initiators[0]

	newPending := j.NewRun(i)
	newPending.Status = models.RunStatusPendingSleep
	newPending.CreatedAt = time.Now().Add(10 * time.Second)
	require.NoError(t, store.CreateJobRun(&newPending))

	oldPending := j.NewRun(i)
	oldPending.Status = models.RunStatusPendingSleep
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

	completedRun := job.NewRun(job.Initiators[0])
	run2 := job.NewRun(job.Initiators[0])
	run3 := job2.NewRun(job2.Initiators[0])

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

func TestORM_CreateTx(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	transaction := cltest.NewTransaction(9182731)

	tx, err := store.CreateTx(transaction)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 0)

	txs := []models.Tx{}
	assert.NoError(t, store.Where("Nonce", transaction.Nonce, &txs))
	require.Len(t, txs, 1)
	ntx := txs[0]

	assert.NotNil(t, ntx.ID)
	assert.NotEmpty(t, ntx.From)
	assert.NotEmpty(t, ntx.To)
	assert.NotEmpty(t, ntx.Data)
	assert.NotEmpty(t, ntx.Nonce)
	assert.NotEmpty(t, ntx.Value.ToInt())
	assert.NotEmpty(t, ntx.GasLimit)
}

func TestORM_CreateTx_WithSurrogateIDIsIdempotent(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	newNonce := uint64(13)

	transaction := cltest.NewTransaction(11)
	transaction.SurrogateID = null.StringFrom("9182323")
	tx1, err := store.CreateTx(transaction)
	assert.NoError(t, err)

	transaction2 := cltest.NewTransaction(newNonce)
	transaction2.SurrogateID = null.StringFrom("9182323")
	tx2, err := store.CreateTx(transaction2)
	assert.NoError(t, err)

	// IDs should be the same because only record should ever be created
	assert.Equal(t, tx1.ID, tx2.ID)

	// New nonce should be saved
	assert.Equal(t, newNonce, tx2.Nonce)

	// New nonce should change the hash
	assert.Equal(t, transaction2.Hash, tx2.Hash)
}

func TestORM_AddTxAttempt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	transaction := cltest.NewTransaction(0)

	tx, err := store.CreateTx(transaction)
	assert.NoError(t, err)

	txAttempt, err := store.AddTxAttempt(tx, transaction)
	assert.NoError(t, err)
	require.Len(t, tx.Attempts, 1)
	assert.Equal(t, tx.ID, txAttempt.TxID)
	assert.Equal(t, tx.Attempts[0], txAttempt)

	transaction = cltest.NewTransaction(1)
	txAttempt, err = store.AddTxAttempt(tx, transaction)
	assert.NoError(t, err)
	require.Len(t, tx.Attempts, 2)
	assert.Equal(t, tx.ID, txAttempt.TxID)
	assert.Equal(t, tx.Attempts[1], txAttempt)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, tx.Hash, txAttempt.Hash)

	// Another attempt with exact same EthTx still generates a new attempt record
	txAttempt, err = store.AddTxAttempt(tx, transaction)
	assert.NoError(t, err)

	require.Len(t, tx.Attempts, 3)
	assert.Equal(t, tx.ID, txAttempt.TxID)
	assert.Equal(t, tx.Attempts[2], txAttempt)

	transaction = cltest.NewTransaction(3)

	// Another attempt with new EthTx updates Tx hash/rawTx etc.
	txAttempt, err = store.AddTxAttempt(tx, transaction)
	assert.NoError(t, err)

	require.Len(t, tx.Attempts, 4)
	assert.Equal(t, tx.ID, txAttempt.TxID)
	assert.Equal(t, tx.Attempts[3], txAttempt)
	assert.Equal(t, tx.Hash, txAttempt.Hash)
	assert.Equal(t, tx.SignedRawTx, txAttempt.SignedRawTx)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, tx.Hash, txAttempt.Hash)
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
			assert.Equal(t, test.want, tt)
			assert.Equal(t, test.errored, err != nil)
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
	initr := job.Initiators[0]

	run := job.NewRun(initr)
	require.NoError(t, store.CreateJobRun(&run))

	executor := services.NewRunExecutor(store)
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
	initr := job.Initiators[0]

	unfinishedRun := job.NewRun(initr)
	retrievedBt, err := store.PendingBridgeType(unfinishedRun)
	assert.NoError(t, err)
	assert.Equal(t, retrievedBt, *bt)
}

func TestORM_GetLastNonce_StormNotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())
	store := app.Store

	account := cltest.GetAccountAddress(t, store)
	nonce, err := store.GetLastNonce(account)

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)
}

func TestORM_GetLastNonce_Valid(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	manager := store.TxManager
	ethMock := app.MockCallerSubscriberClient()
	one := uint64(1)

	ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(one))
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	ethMock.Register("eth_chainId", store.Config.ChainID())

	assert.NoError(t, app.StartAndConnect())

	to := cltest.NewAddress()
	_, err := manager.CreateTx(to, []byte{})
	assert.NoError(t, err)

	account := cltest.GetAccountAddress(t, store)
	nonce, err := store.GetLastNonce(account)

	assert.NoError(t, err)
	assert.Equal(t, one, nonce)
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

	assert.NoError(t, store.MarkRan(&initr, true))
	ir, err := store.FindInitiator(initr.ID)
	assert.NoError(t, err)
	assert.True(t, ir.Ran)

	assert.Error(t, store.MarkRan(&initr, true))
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	user1 := cltest.MustUser("test1@email1.net", "password1")
	user2 := cltest.MustUser("test2@email2.net", "password2")
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

			user := cltest.MustUser("have@email", "password")
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
	user := cltest.MustUser("test1@email1.net", "password1")
	require.NoError(t, store.SaveUser(&user))

	_, err := store.DeleteUser()
	require.NoError(t, err)

	_, err = store.FindUser()
	require.Error(t, err)
}

func TestORM_DeleteUserSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	user := cltest.MustUser("test1@email1.net", "password1")
	require.NoError(t, store.SaveUser(&user))

	session := models.NewSession()
	require.NoError(t, store.SaveSession(&session))

	err := store.DeleteUserSession(session.ID)
	require.NoError(t, err)

	user, err = store.FindUser()
	require.NoError(t, err)

	sessions, err := store.Sessions(0, 10)
	assert.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"correct", cltest.APIEmail, cltest.Password, true},
		{"incorrect email", "bogus@town.org", cltest.Password, false},
		{"incorrect pwd", cltest.APIEmail, "jamaicandundada", false},
		{"incorrect both", "dudus@coke.ja", "jamaicandundada", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			initial := cltest.MustUser(cltest.APIEmail, cltest.Password)
			require.NoError(t, store.SaveUser(&initial))

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

func TestORM_DeleteTransaction(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)
	transaction := cltest.NewTransaction(0)
	require.NoError(t, utils.JustError(store.AddTxAttempt(tx, transaction)))

	require.NoError(t, store.DeleteTransaction(tx))

	_, err = store.FindTx(tx.ID)
	require.Error(t, err)
}

func TestORM_AllSyncEvents(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := store.ORM
	synchronization.NewStatsPusher(orm, cltest.MustParseURL("http://localhost"), "", "")

	// Create two events via job run callback
	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	require.NoError(t, store.ORM.CreateJob(&job))
	initiator := job.Initiators[0]

	oldIncompleteRun := job.NewRun(initiator)
	oldIncompleteRun.Status = models.RunStatusInProgress
	err := orm.CreateJobRun(&oldIncompleteRun)
	require.NoError(t, err)

	newCompletedRun := job.NewRun(initiator)
	newCompletedRun.Status = models.RunStatusCompleted
	err = orm.CreateJobRun(&newCompletedRun)
	require.NoError(t, err)

	events := []models.SyncEvent{}
	err = orm.AllSyncEvents(func(event *models.SyncEvent) error {
		events = append(events, *event)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, events, 2)
	assert.Greater(t, events[1].ID, events[0].ID)
}

func TestBulkDeleteRuns(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := store.ORM

	err := orm.RawDB(func(db *gorm.DB) error {
		job := cltest.NewJobWithWebInitiator()
		job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
		require.NoError(t, store.ORM.CreateJob(&job))
		initiator := job.Initiators[0]

		// bulk delete should not delete these because they match the updated before
		// but none of the statuses
		oldIncompleteRun := job.NewRun(initiator)
		oldIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 17}`)}
		oldIncompleteRun.Status = models.RunStatusInProgress
		err := orm.CreateJobRun(&oldIncompleteRun)
		require.NoError(t, err)
		db.Model(&oldIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete *SHOULD* delete these because they match one of the statuses
		// and the updated before
		oldCompletedRun := job.NewRun(initiator)
		oldCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 19}`)}
		oldCompletedRun.Status = models.RunStatusCompleted
		err = orm.CreateJobRun(&oldCompletedRun)
		require.NoError(t, err)
		db.Model(&oldCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-01T00:00:00Z"))

		// bulk delete should not delete these because they match one of the
		// statuses but not the updated before
		newCompletedRun := job.NewRun(initiator)
		newCompletedRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 23}`)}
		newCompletedRun.Status = models.RunStatusCompleted
		err = orm.CreateJobRun(&newCompletedRun)
		require.NoError(t, err)
		db.Model(&newCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		// bulk delete should not delete these because none of their attributes match
		newIncompleteRun := job.NewRun(initiator)
		newIncompleteRun.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result": 71}`)}
		newIncompleteRun.Status = models.RunStatusCompleted
		err = orm.CreateJobRun(&newIncompleteRun)
		require.NoError(t, err)
		db.Model(&newIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601(t, "2018-01-30T00:00:00Z"))

		err = store.ORM.BulkDeleteRuns(&models.BulkDeleteRunRequest{
			Status:        []models.RunStatus{models.RunStatusCompleted},
			UpdatedBefore: cltest.ParseISO8601(t, "2018-01-15T00:00:00Z"),
		})

		require.NoError(t, err)

		var runCount int
		err = db.Model(&models.JobRun{}).Count(&runCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, runCount)

		var taskCount int
		err = db.Model(&models.TaskRun{}).Count(&taskCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, taskCount)

		var resultCount int
		err = db.Model(&models.RunResult{}).Count(&resultCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, resultCount)

		var requestCount int
		err = db.Model(&models.RunRequest{}).Count(&requestCount).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, requestCount)
		return nil
	})
	require.NoError(t, err)
}

func TestORM_FindTxAttempt_CurrentAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	txAttempt, err := store.FindTxAttempt(tx.Attempts[0].Hash)
	require.NoError(t, err)

	assert.Equal(t, tx.ID, txAttempt.ID)
	assert.Equal(t, tx.Confirmed, txAttempt.Confirmed)
	assert.Equal(t, tx.Hash, txAttempt.Hash)
	assert.Equal(t, tx.GasPrice, txAttempt.GasPrice)
	assert.Equal(t, tx.SentAt, txAttempt.SentAt)
	assert.Equal(t, tx.SignedRawTx, txAttempt.SignedRawTx)
}

func TestORM_FindTxAttempt_PastAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)
	transaction := cltest.NewTransaction(0)
	require.NoError(t, utils.JustError(store.AddTxAttempt(tx, transaction)))

	txAttempt, err := store.FindTxAttempt(tx.Attempts[0].Hash)
	require.NoError(t, err)

	assert.Equal(t, tx.ID, txAttempt.TxID)
	assert.Equal(t, tx.Confirmed, txAttempt.Confirmed)
	assert.NotEqual(t, tx.Hash, txAttempt.Hash)
	assert.NotEqual(t, tx.GasPrice, txAttempt.GasPrice)
	assert.NotEqual(t, tx.SentAt, txAttempt.SentAt)
	assert.NotEqual(t, tx.SignedRawTx, txAttempt.SignedRawTx)
}

func TestORM_FindTxByAttempt_CurrentAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(t, store)

	createdTx := cltest.CreateTx(t, store, from, 1)
	fetchedTx, fetchedTxAttempt, err := store.FindTxByAttempt(createdTx.Hash)

	assert.Equal(t, createdTx.ID, fetchedTx.ID)
	assert.Equal(t, createdTx.From, fetchedTx.From)
	assert.Equal(t, createdTx.To, fetchedTx.To)
	assert.Equal(t, createdTx.Nonce, fetchedTx.Nonce)
	assert.Equal(t, createdTx.Value, fetchedTx.Value)
	assert.Equal(t, createdTx.GasLimit, fetchedTx.GasLimit)
	assert.Equal(t, createdTx.Confirmed, fetchedTx.Confirmed)
	assert.Equal(t, createdTx.Hash, fetchedTx.Hash)
	assert.Equal(t, createdTx.GasPrice, fetchedTx.GasPrice)
	assert.Equal(t, createdTx.SentAt, fetchedTx.SentAt)

	assert.Equal(t, createdTx.ID, fetchedTxAttempt.ID)
	assert.Equal(t, createdTx.Confirmed, fetchedTxAttempt.Confirmed)
	assert.Equal(t, createdTx.Hash, fetchedTxAttempt.Hash)
	assert.Equal(t, createdTx.GasPrice, fetchedTxAttempt.GasPrice)
	assert.Equal(t, createdTx.SentAt, fetchedTxAttempt.SentAt)
}

func TestORM_FindTxByAttempt_PastAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(t, store)
	createdTx := cltest.CreateTx(t, store, from, 1)
	pastTxAttempt := createdTx.Attempts[0]

	transaction := cltest.NewTransaction(0)
	require.NoError(t, utils.JustError(store.AddTxAttempt(createdTx, transaction)))

	fetchedTx, pastTxAttempt, err := store.FindTxByAttempt(pastTxAttempt.Hash)
	require.NoError(t, err)

	assert.Equal(t, createdTx.ID, fetchedTx.ID)
	assert.Equal(t, createdTx.From, fetchedTx.From)
	assert.Equal(t, createdTx.To, fetchedTx.To)
	assert.Equal(t, createdTx.Nonce, fetchedTx.Nonce)
	assert.Equal(t, createdTx.Value, fetchedTx.Value)
	assert.Equal(t, createdTx.GasLimit, fetchedTx.GasLimit)
	assert.Equal(t, createdTx.Confirmed, fetchedTx.Confirmed)
	assert.Equal(t, createdTx.Hash, fetchedTx.Hash)
	assert.Equal(t, createdTx.GasPrice, fetchedTx.GasPrice)
	assert.Equal(t, createdTx.SentAt, fetchedTx.SentAt)

	assert.Equal(t, createdTx.ID, pastTxAttempt.TxID)
	assert.NotEqual(t, createdTx.Hash, pastTxAttempt.Hash)
	assert.NotEqual(t, createdTx.GasPrice, pastTxAttempt.GasPrice)
	assert.NotEqual(t, createdTx.SentAt, pastTxAttempt.SentAt)
	assert.NotEqual(t, createdTx.SignedRawTx, pastTxAttempt.SignedRawTx)
}

func TestORM_DeduceDialect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, connection string
		expect           orm.DialectName
		wantError        bool
	}{
		{"windows full path", `D:/node-0/node/db.sqlite3`, `sqlite3`, false},
		{"relative file", "db.sqlite", "sqlite3", false},
		{"relative dir path", "store/db/here", "sqlite3", false},
		{"file url", "file://host/path", "sqlite3", false},
		{"sqlite url", "sqlite:///path/to/sqlite.db", "", true},
		{"sqlite3 url", "sqlite3:///path/to/sqlite.db", "", true},
		{"postgres url", "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full", "postgres", false},
		{"postgresql url", "postgresql://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full", "postgres", false},
		{"postgres string", "user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := orm.DeduceDialect(test.connection)
			assert.Equal(t, test.expect, actual)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}

func TestORM_SyncDbKeyStoreToDisk(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := store.ORM

	seed, err := models.NewKeyFromFile("../../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	require.NoError(t, err)
	require.NoError(t, orm.FirstOrCreateKey(seed))

	keysDir := store.Config.KeysDir()
	require.True(t, isDirEmpty(t, keysDir))
	require.NoError(t, orm.ClobberDiskKeyStoreWithDBKeys(keysDir))

	dbkeys, err := store.Keys()
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

func TestORM_UnconfirmedTxAttempts(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	t.Run("tx #1, 4 attempts", func(t *testing.T) {
		transaction := cltest.NewTransaction(0, 0)
		transaction.SurrogateID = null.StringFrom("0")
		tx, err := store.CreateTx(transaction)
		require.NoError(t, err)

		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 1)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 2)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 3)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)
		require.Len(t, tx.Attempts, 4)

		tx.Attempts[0].GasPrice = utils.NewBig(big.NewInt(1111))
		tx.Attempts[1].GasPrice = utils.NewBig(big.NewInt(2222))
		tx.Attempts[2].GasPrice = utils.NewBig(big.NewInt(3333))
		tx.Attempts[3].GasPrice = utils.NewBig(big.NewInt(4444))

		err = store.ORM.RawDB(func(db *gorm.DB) error {
			return db.Save(&tx).Error
		})
		require.NoError(t, err)
	})

	t.Run("tx #2, 3 attempts", func(t *testing.T) {
		transaction := cltest.NewTransaction(0)
		transaction.SurrogateID = null.StringFrom("1")
		tx, err := store.CreateTx(transaction)
		require.NoError(t, err)

		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 1)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 2)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)
		require.Len(t, tx.Attempts, 3)

		tx.Attempts[0].GasPrice = utils.NewBig(big.NewInt(5555))
		tx.Attempts[1].GasPrice = utils.NewBig(big.NewInt(6666))
		tx.Attempts[2].GasPrice = utils.NewBig(big.NewInt(7777))

		err = store.ORM.RawDB(func(db *gorm.DB) error {
			return db.Save(&tx).Error
		})
		require.NoError(t, err)
	})

	t.Run("tx #2, 2 attempts", func(t *testing.T) {
		transaction := cltest.NewTransaction(0)
		transaction.SurrogateID = null.StringFrom("2")
		tx, err := store.CreateTx(transaction)
		require.NoError(t, err)

		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		transaction = cltest.NewTransaction(0, 1)
		_, err = store.AddTxAttempt(tx, transaction)
		require.NoError(t, err)

		// This tx's attempts should not appear in the results
		tx.Confirmed = true

		err = store.ORM.RawDB(func(db *gorm.DB) error {
			return db.Save(&tx).Error
		})
		require.NoError(t, err)
	})

	attempts, err := store.ORM.UnconfirmedTxAttempts()
	require.NoError(t, err)

	assert.Len(t, attempts, 7)
}
