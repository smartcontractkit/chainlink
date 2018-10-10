package orm_test

import (
	"encoding/hex"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhereNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := models.NewJob()
	jobs := []models.JobSpec{j1}

	err := store.Where("ID", "bogus", &jobs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestAllNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	var jobs []models.JobSpec
	err := store.All(&jobs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORM_SaveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1, initr := cltest.NewJobWithSchedule("* * * * *")
	store.SaveJob(&j1)

	j2, _ := store.FindJob(j1.ID)
	assert.Equal(t, j1.ID, j2.ID)
	assert.Equal(t, j1.Initiators[0], j2.Initiators[0])
	assert.Equal(t, j2.ID, j2.Initiators[0].JobID)

	assert.NoError(t, store.One("JobID", j1.ID, &initr))
	assert.Equal(t, models.Cron("* * * * *"), initr.Schedule)
}

func TestJobRunsFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, i := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.SaveJob(&job))
	jr1 := job.NewRun(i)
	jr1.CreatedAt = time.Now().AddDate(0, 0, -1)
	require.NoError(t, store.Save(&jr1))
	jr2 := job.NewRun(i)
	jr2.CreatedAt = time.Now().AddDate(0, 0, 1)
	require.NoError(t, store.Save(&jr2))
	jr3 := job.NewRun(i)
	jr3.CreatedAt = time.Now().AddDate(0, 0, -9)
	require.NoError(t, store.Save(&jr3))

	runs, err := store.JobRunsFor(job.ID)
	assert.NoError(t, err)
	actual := []string{runs[0].ID, runs[1].ID, runs[2].ID}
	assert.Equal(t, []string{jr2.ID, jr1.ID, jr3.ID}, actual)
}

func TestORM_SaveServiceAgreement(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic",
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"HttpGet","url":"https://bitstamp.net/api/ticker/"},{"type":"JsonParse","path":["last"]},{"type":"EthBytes32"},{"type":"EthTx"}]}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa, err := cltest.ServiceAgreementFromString(test.input)
			assert.NoError(t, err)

			assert.NoError(t, store.SaveServiceAgreement(&sa))

			sa, err = store.FindServiceAgreement(sa.ID)
			assert.NoError(t, err)
			_, err = store.FindJob(sa.JobSpecID)
			assert.NoError(t, err)
		})
	}
}

func TestJobRunsWithStatus(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j, i := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.SaveJob(&j))
	npr := j.NewRun(i)
	assert.NoError(t, store.Save(&npr))

	statuses := []models.RunStatus{
		models.RunStatusPendingBridge,
		models.RunStatusPendingConfirmations,
		models.RunStatusCompleted}
	var seedIds []string
	for _, status := range statuses {
		run := j.NewRun(i)
		run.Status = status
		assert.NoError(t, store.Save(&run))
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

			pending, err := store.JobRunsWithStatus(test.statuses...)
			assert.NoError(t, err)

			pendingIDs := []string{}
			for _, jr := range pending {
				pendingIDs = append(pendingIDs, jr.ID)
			}
			assert.ElementsMatch(t, pendingIDs, test.expected)
		})
	}
}

func TestAnyJobWithType(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	js, _ := cltest.NewJobWithWebInitiator()
	js.Tasks = []models.TaskSpec{models.TaskSpec{Type: models.MustNewTaskType("bridgetestname")}}
	assert.NoError(t, store.Save(&js))
	found, err := store.AnyJobWithType("bridgetestname")
	assert.NoError(t, err)
	assert.Equal(t, found, true)
	found, err = store.AnyJobWithType("somethingelse")
	assert.NoError(t, err)
	assert.Equal(t, found, false)

}

func TestJobRunsCountFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	job, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.SaveJob(&job))
	job2, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.SaveJob(&job2))

	assert.NotEqual(t, job.ID, job2.ID)

	completedRun := job.NewRun(initr)
	run2 := job.NewRun(initr)
	run3 := job2.NewRun(initr)

	assert.NoError(t, store.Save(&completedRun))
	assert.NoError(t, store.Save(&run2))
	assert.NoError(t, store.Save(&run3))

	count, err := store.JobRunsCountFor(job.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	count, err = store.JobRunsCountFor(job2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCreatingTx(t *testing.T) {
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

	_, err = store.CreateTx(from, nonce, to, data, value, gasLimit)
	assert.NoError(t, err)

	txs := []models.Tx{}
	assert.NoError(t, store.Where("Nonce", nonce, &txs))
	assert.Equal(t, 1, len(txs))
	tx := txs[0]

	assert.NotNil(t, tx.ID)
	assert.Equal(t, from, tx.From)
	assert.Equal(t, to, tx.To)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, nonce, tx.Nonce)
	assert.Equal(t, value, tx.Value)
	assert.Equal(t, gasLimit, tx.GasLimit)
}

func TestFindBridge(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	tt := models.BridgeType{}
	tt.Name = models.MustNewTaskType("solargridreporting")
	tt.URL = cltest.WebURL("https://denergy.eth")
	assert.NoError(t, store.Save(&tt))

	cases := []struct {
		description string
		name        string
		want        models.BridgeType
		errored     bool
	}{
		{"actual external adapter", tt.Name.String(), tt, false},
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

	bt := cltest.NewBridgeType()
	assert.NoError(t, store.Save(&bt))

	job, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, store.SaveJob(&job))

	run := job.NewRun(initr)
	assert.NoError(t, store.Save(&run))

	store.RunChannel.Send(run.ID)
	cltest.WaitForJobRunStatus(t, store, run, models.RunStatusCompleted)

	_, err := store.PendingBridgeType(run)
	assert.Error(t, err)
}

func TestORM_PendingBridgeType_success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	bt := cltest.NewBridgeType()
	assert.NoError(t, store.Save(&bt))

	job, initr := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{models.TaskSpec{Type: bt.Name}}
	assert.NoError(t, store.SaveJob(&job))

	unfinishedRun := job.NewRun(initr)
	retrievedBt, err := store.PendingBridgeType(unfinishedRun)
	assert.NoError(t, err)
	assert.Equal(t, bt, retrievedBt)
}

func TestORM_GetLastNonce_StormNotFound(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store

	account := cltest.GetAccountAddress(store)
	nonce, err := store.GetLastNonce(account)

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)
}

func TestORM_GetLastNonce_Valid(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager
	ethMock := app.MockEthClient()
	one := uint64(1)

	ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(one))
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(one))
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	assert.NoError(t, app.Start())

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

	_, initr := cltest.NewJobWithRunAtInitiator(time.Now())
	assert.NoError(t, store.Save(&initr))

	assert.NoError(t, store.MarkRan(&initr))
	var ir models.Initiator
	assert.NoError(t, store.One("ID", initr.ID, &ir))
	assert.True(t, ir.Ran)
}

func TestORM_FindUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	user1 := cltest.MustUser("test1@email1.net", "password1")
	user2 := cltest.MustUser("test2@email2.net", "password2")
	user2.CreatedAt = models.Time{time.Now().Add(-24 * time.Hour)}

	require.NoError(t, store.Save(&user1))
	require.NoError(t, store.Save(&user2))

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
	require.NoError(t, store.Save(&user))

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
			prevSession.LastUsed = models.Time{time.Now().Add(-cltest.MustParseDuration("2m"))}
			require.NoError(t, store.Save(&prevSession))

			expectedTime := models.Time{time.Now()}.HumanString()
			actual, err := store.ORM.AuthorizedUserWithSession(test.sessionID, test.sessionDuration)
			assert.Equal(t, test.wantEmail, actual.Email)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var bumpedSession models.Session
				err = store.One("ID", prevSession.ID, &bumpedSession)
				require.NoError(t, err)
				assert.Equal(t, expectedTime[0:13], bumpedSession.LastUsed.HumanString()[0:13]) // only compare up to the hour
			}
		})
	}
}

func TestORM_DeleteUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	user := cltest.MustUser("test1@email1.net", "password1")
	require.NoError(t, store.Save(&user))

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
	require.NoError(t, store.Save(&user))

	session := models.NewSession()
	require.NoError(t, store.Save(&session))

	err := store.DeleteUserSession(session.ID)
	require.NoError(t, err)

	user, err = store.FindUser()
	require.NoError(t, err)

	var sessions []models.Session
	err = store.All(&sessions)
	require.NoError(t, err)
	require.Empty(t, sessions)
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	initial := cltest.MustUser(cltest.APIEmail, cltest.Password)
	require.NoError(t, store.Save(&initial))

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

func TestORM_AllInBatches_DifferentBatchSizes(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	var expected []string
	for i := 0; i < 5; i++ {
		job, _ := cltest.NewJobWithWebInitiator()
		require.NoError(t, store.SaveJob(&job))
		expected = append(expected, job.ID)
	}

	tests := []struct {
		name      string
		batchSize int
	}{
		{"smaller", 3},
		{"equal", 5},
		{"larger", 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bucket []models.JobSpec // Need this to tell storm which bucket to use (via reflection)
			var actual []string

			err := store.AllInBatches(&bucket, func(j models.JobSpec) bool {
				actual = append(actual, j.ID)
				return true
			}, test.batchSize)

			require.NoError(t, err)
			sort.Strings(expected)
			sort.Strings(actual)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestORM_AllInBatches_EarlyExit(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	var jobids []string
	for i := 0; i < 2; i++ {
		job, _ := cltest.NewJobWithWebInitiator()
		require.NoError(t, store.SaveJob(&job))
		jobids = append(jobids, job.ID)
	}

	var bucket []models.JobSpec
	var actual []string
	err := store.AllInBatches(&bucket, func(j models.JobSpec) bool {
		actual = append(actual, j.ID)
		return false
	})

	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))
	assert.Contains(t, jobids, actual[0])
}

func TestORM_AllInBatches_Empty(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	var bucket []models.JobSpec
	var actual []string
	err := store.AllInBatches(&bucket, func(j models.JobSpec) bool {
		actual = append(actual, j.ID)
		return true
	})

	require.NoError(t, err)
	assert.Empty(t, actual)
}

func TestORM_AllInBatches_IncorrectCallback(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, _ := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.SaveJob(&job))

	var bucket []models.JobSpec

	tests := []struct {
		name     string
		callback interface{}
		err      error
	}{
		{"missing bool", func(models.JobSpec) {}, orm.ErrorInvalidCallbackSignature},
		{"wrong rval", func(models.JobSpec) int { return 0 }, orm.ErrorInvalidCallbackSignature},
		{"mismatched bucket and model", func(models.BridgeType) bool { return true }, orm.ErrorInvalidCallbackModel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.AllInBatches(&bucket, test.callback)
			require.Equal(t, test.err, err)
		})
	}
}
