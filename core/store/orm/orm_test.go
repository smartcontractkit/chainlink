package orm_test

import (
	"encoding/hex"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_WhereNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := models.NewJob()
	jobs := []models.JobSpec{j1}

	err := store.Where("ID", "bogus", &jobs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORM_AllNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	var jobs []models.JobSpec
	err := store.ORM.DB.Find(&jobs).Error
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORM_CreateJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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
	store, cleanup := cltest.NewStore()
	defer cleanup()

	orm := store.ORM
	job := cltest.NewJob()
	require.NoError(t, orm.CreateJob(&job))
	require.NoError(t, orm.DB.Delete(&job).Error)
	require.Error(t, orm.DB.First(&job).Error)
	orm = store.ORM.Unscoped()
	require.NoError(t, orm.DB.First(&job).Error)
}

func TestORM_ArchiveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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

func TestORM_CreateJobRun_ArchivesRunIfJobArchived(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	require.NoError(t, store.ArchiveJob(job.ID))

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	require.Error(t, utils.JustError(store.FindJobRun(jr.ID)))
	require.NoError(t, utils.JustError(store.Unscoped().FindJobRun(jr.ID)))
}

func TestORM_CreateJobRun_CreatesRunRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	var requestCount int
	err := store.ORM.DB.Model(&models.RunRequest{}).Count(&requestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}

func TestORM_SaveJobRun_DoesNotSaveTaskSpec(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := cltest.NewJobWithSchedule("* * * * *")
	require.NoError(t, store.CreateJob(&job))

	jr := job.NewRun(job.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	var err error
	jr.TaskRuns[0].TaskSpec.Params, err = jr.TaskRuns[0].TaskSpec.Params.Merge(cltest.JSONFromString(t, `{"random": "input"}`))
	require.NoError(t, err)
	require.NoError(t, store.SaveJobRun(&jr))

	retrievedJob, err := store.FindJob(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Tasks, 1)
	require.Len(t, retrievedJob.Tasks, 1)
	assert.JSONEq(
		t,
		coercedJSON(job.Tasks[0].Params.String()),
		retrievedJob.Tasks[0].Params.String())
}

func TestORM_SaveJobRun_ArchivedDoesNotRevertDeletedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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

func coercedJSON(v string) string {
	if v == "" {
		return "{}"
	}
	return v
}

func TestORM_JobRunsFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
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
	actual := []string{runs[0].ID, runs[1].ID, runs[2].ID}
	assert.Equal(t, []string{jr2.ID, jr1.ID, jr3.ID}, actual)
}

func TestORM_JobRunsSortedFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
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
	actual := []string{runs[0].ID, runs[1].ID} // doesn't include excludedJobRun
	assert.Equal(t, []string{jr2.ID, jr1.ID}, actual)
}

func TestORM_UnscopedJobRunsWithStatus_Happy(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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

	var seedIds []string
	for _, status := range statuses {
		run := j.NewRun(i)
		run.Status = status
		require.NoError(t, store.CreateJobRun(&run))
		seedIds = append(seedIds, run.ID)
	}

	tests := []struct {
		name     string
		statuses []models.RunStatus
		expected []string
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]string{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingConfirmations},
			[]string{seedIds[0], seedIds[1]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []string{}
			for _, jr := range pending {
				pendingIDs = append(pendingIDs, jr.ID)
			}
			assert.ElementsMatch(t, pendingIDs, test.expected)
		})
	}
}

func TestORM_UnscopedJobRunsWithStatus_Deleted(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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

	var seedIds []string
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
		expected []string
	}{
		{
			"single status",
			[]models.RunStatus{models.RunStatusPendingBridge},
			[]string{seedIds[0]},
		},
		{
			"multiple status'",
			[]models.RunStatus{models.RunStatusPendingBridge, models.RunStatusPendingConfirmations, models.RunStatusPendingConnection},
			[]string{seedIds[0], seedIds[1], seedIds[2]},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			pending := cltest.MustAllJobsWithStatus(t, store, test.statuses...)

			pendingIDs := []string{}
			for _, jr := range pending {
				pendingIDs = append(pendingIDs, jr.ID)
			}
			assert.ElementsMatch(t, pendingIDs, test.expected)
		})
	}
}

func TestORM_UnscopedJobRunsWithStatus_OrdersByCreatedAt(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
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

func TestORM_CreatingTx(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	from := common.HexToAddress("0x2C83ACd90367e7E0D3762eA31aC77F18faecE874")
	to := common.HexToAddress("0x4A7d17De4B3eC94c59BF07764d9A6e97d92A547A")
	value := new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)
	nonce := uint64(1232421)
	gasLimit := uint64(50000)
	data, err := hex.DecodeString("0987612345abcdef")
	assert.NoError(t, err)

	tx := &models.Tx{
		From:     from,
		To:       to,
		Nonce:    nonce,
		Data:     data,
		Value:    models.NewBig(value),
		GasLimit: gasLimit,
	}

	err = store.CreateTx(tx)
	assert.NoError(t, err)

	txs := []models.Tx{}
	assert.NoError(t, store.Where("Nonce", nonce, &txs))
	require.Len(t, txs, 1)
	ntx := txs[0]

	assert.NotNil(t, ntx.ID)
	assert.Equal(t, from, ntx.From)
	assert.Equal(t, to, ntx.To)
	assert.Equal(t, data, ntx.Data)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Equal(t, value, ntx.Value.ToInt())
	assert.Equal(t, gasLimit, ntx.GasLimit)
}

func TestORM_FindBridge(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL("https://denergy.eth")
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

	store, cleanup := cltest.NewStore()
	defer cleanup()
	jobRunner, cleanup := cltest.NewJobRunner(store)
	defer cleanup()
	jobRunner.Start()

	_, bt := cltest.NewBridgeType()
	require.NoError(t, store.CreateBridgeType(bt))

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	initr := job.Initiators[0]

	run := job.NewRun(initr)
	require.NoError(t, store.CreateJobRun(&run))

	store.RunChannel.Send(run.ID)
	cltest.WaitForJobRunStatus(t, store, run, models.RunStatusCompleted)

	_, err := store.PendingBridgeType(run)
	assert.Error(t, err)
}

func TestORM_PendingBridgeType_success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	_, bt := cltest.NewBridgeType()
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
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store

	account := cltest.GetAccountAddress(store)
	nonce, err := store.GetLastNonce(account)

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)
}

func TestORM_GetLastNonce_Valid(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	manager := store.TxManager
	ethMock := app.MockEthClient()
	one := uint64(1)

	ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(one))
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(one))
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	assert.NoError(t, app.StartAndConnect())

	to := cltest.NewAddress()
	_, err := manager.CreateTx(to, []byte{})
	assert.NoError(t, err)

	account := cltest.GetAccountAddress(store)
	nonce, err := store.GetLastNonce(account)

	assert.NoError(t, err)
	assert.Equal(t, one, nonce)
}

func TestORM_MarkRan(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
	defer cleanup()

	user := cltest.MustUser("have@email", "password")
	require.NoError(t, store.SaveUser(&user))

	tests := []struct {
		name            string
		sessionID       string
		sessionDuration time.Duration
		wantError       bool
		wantEmail       string
	}{
		{"authorized", "correctID", cltest.MustParseDuration("3m"), false, "have@email"},
		{"expired", "correctID", cltest.MustParseDuration("0m"), true, ""},
		{"incorrect", "wrong", cltest.MustParseDuration("3m"), true, ""},
		{"empty", "", cltest.MustParseDuration("3m"), true, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevSession := cltest.NewSession("correctID")
			prevSession.LastUsed = time.Now().Add(-cltest.MustParseDuration("2m"))
			require.NoError(t, store.SaveSession(&prevSession))

			expectedTime := utils.ISO8601UTC(time.Now())
			actual, err := store.ORM.AuthorizedUserWithSession(test.sessionID, test.sessionDuration)
			assert.Equal(t, test.wantEmail, actual.Email)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var bumpedSession models.Session
				err = store.ORM.DB.First(&bumpedSession, "ID = ?", prevSession.ID).Error
				require.NoError(t, err)
				assert.Equal(t, expectedTime[0:13], utils.ISO8601UTC(bumpedSession.LastUsed)[0:13]) // only compare up to the hour
			}
		})
	}
}

func TestORM_DeleteUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
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

	store, cleanup := cltest.NewStore()
	defer cleanup()

	initial := cltest.MustUser(cltest.APIEmail, cltest.Password)
	require.NoError(t, store.SaveUser(&initial))

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
	store, cleanup := cltest.NewStore()
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(store)
	tx := cltest.CreateTxAndAttempt(store, from, 1)
	_, err = store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), 3)
	require.NoError(t, err)

	require.NoError(t, store.DeleteTransaction(tx))

	_, err = store.FindTx(tx.ID)
	assert.Error(t, err)
	attempts, err := store.TxAttemptsFor(tx.ID)
	assert.NoError(t, err)
	assert.Empty(t, attempts)
}

func TestORM_AllSyncEvents(t *testing.T) {
	store, cleanup := cltest.NewStore()
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

	assert.Len(t, events, 2)
	assert.Greater(t, events[1].ID, events[0].ID)
}

func TestBulkDeleteRuns(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	orm := store.ORM
	db := orm.DB
	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	require.NoError(t, store.ORM.CreateJob(&job))
	initiator := job.Initiators[0]

	// matches updated before but none of the statuses
	oldIncompleteRun := job.NewRun(initiator)
	oldIncompleteRun.Status = models.RunStatusInProgress
	err := orm.CreateJobRun(&oldIncompleteRun)
	require.NoError(t, err)
	db.Model(&oldIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-01T00:00:00Z"))

	// matches one of the statuses and the updated before
	oldCompletedRun := job.NewRun(initiator)
	oldCompletedRun.Status = models.RunStatusCompleted
	err = orm.CreateJobRun(&oldCompletedRun)
	require.NoError(t, err)
	db.Model(&oldCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-01T00:00:00Z"))

	// matches one of the statuses but not the updated before
	newCompletedRun := job.NewRun(initiator)
	newCompletedRun.Status = models.RunStatusCompleted
	err = orm.CreateJobRun(&newCompletedRun)
	require.NoError(t, err)
	db.Model(&newCompletedRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-30T00:00:00Z"))

	// matches nothing
	newIncompleteRun := job.NewRun(initiator)
	newIncompleteRun.Status = models.RunStatusCompleted
	err = orm.CreateJobRun(&newIncompleteRun)
	require.NoError(t, err)
	db.Model(&newIncompleteRun).UpdateColumn("updated_at", cltest.ParseISO8601("2018-01-30T00:00:00Z"))

	err = store.ORM.BulkDeleteRuns(&models.BulkDeleteRunRequest{
		Status:        []models.RunStatus{models.RunStatusCompleted},
		UpdatedBefore: cltest.ParseISO8601("2018-01-15T00:00:00Z"),
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
	assert.Equal(t, 6, resultCount)

	var requestCount int
	err = db.Model(&models.RunRequest{}).Count(&requestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, 3, requestCount)
}

func TestORM_FindTxByAttempt_CurrentAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(store)

	tx1 := cltest.CreateTxAndAttempt(store, from, 1)
	tx2, err := store.FindTxByAttempt(tx1.Hash)

	require.Equal(t, tx1.ID, tx2.ID)
	require.Equal(t, tx1.From, tx2.From)
	require.Equal(t, tx1.To, tx2.To)
	require.Equal(t, tx1.Nonce, tx1.Nonce)
	require.Equal(t, tx1.Value, tx2.Value)
	require.Equal(t, tx1.GasLimit, tx2.GasLimit)
	require.Equal(t, tx1.Confirmed, tx2.Confirmed)

	require.Equal(t, tx1.Hash, tx2.Hash)
	require.Equal(t, tx1.GasPrice, tx2.GasPrice)
	require.Equal(t, tx1.Hex, tx2.Hex)
	require.Equal(t, tx1.SentAt, tx2.SentAt)
}

func TestORM_FindTxByAttempt_PastAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	defer cleanup()

	from := cltest.GetAccountAddress(store)
	tx := cltest.CreateTxAndAttempt(store, from, 1)
	tx1 := *tx
	_, err = store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), 3)
	require.NoError(t, err)

	tx2, err := store.FindTxByAttempt(tx1.Hash)
	require.NoError(t, err)

	require.Equal(t, tx.ID, tx2.ID)
	require.Equal(t, tx.From, tx2.From)
	require.Equal(t, tx.To, tx2.To)
	require.Equal(t, tx.Nonce, tx1.Nonce)
	require.Equal(t, tx.Value, tx2.Value)
	require.Equal(t, tx.GasLimit, tx2.GasLimit)
	require.Equal(t, tx.Confirmed, tx2.Confirmed)
	require.NotEqual(t, tx.Hash, tx2.Hash)
	require.NotEqual(t, tx.GasPrice, tx2.GasPrice)
	require.NotEqual(t, tx.Hex, tx2.Hex)
	require.NotEqual(t, tx.SentAt, tx2.SentAt)

	require.Equal(t, tx1.ID, tx2.ID)
	require.Equal(t, tx1.From, tx2.From)
	require.Equal(t, tx1.To, tx2.To)
	require.Equal(t, tx1.Nonce, tx1.Nonce)
	require.Equal(t, tx1.Value, tx2.Value)
	require.Equal(t, tx1.GasLimit, tx2.GasLimit)
	require.Equal(t, tx1.Confirmed, tx2.Confirmed)
	require.Equal(t, tx1.Hash, tx2.Hash)
	require.Equal(t, tx1.GasPrice, tx2.GasPrice)
	require.Equal(t, tx1.Hex, tx2.Hex)
	require.Equal(t, tx1.SentAt, tx2.SentAt)
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
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	store := app.GetStore()
	orm := store.ORM

	seed, err := models.NewKeyFromFile("../../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	require.NoError(t, err)
	require.NoError(t, orm.FirstOrCreateKey(seed))

	keysDir := store.Config.KeysDir()
	require.True(t, isDirEmpty(t, keysDir))
	require.NoError(t, orm.ClobberDiskKeyStoreWithDBKeys(keysDir))

	dbkeys, err := store.Keys()
	require.NoError(t, err)
	assert.Len(t, dbkeys, 1)

	diskkeys, err := utils.FilesInDir(keysDir)
	require.NoError(t, err)
	assert.Len(t, diskkeys, 1)

	key := dbkeys[0]
	content, err := utils.FileContents(filepath.Join(keysDir, diskkeys[0]))
	require.NoError(t, err)
	assert.Equal(t, key.JSON.String(), content)
}

func TestORM_UpdateBridgeType(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	firstBridge := &models.BridgeType{
		Name: "UniqueName",
		URL:  cltest.WebURL("http:/oneurl.com"),
	}

	require.NoError(t, store.CreateBridgeType(firstBridge))

	updateBridge := &models.BridgeTypeRequest{
		URL: cltest.WebURL("http:/updatedurl.com"),
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
