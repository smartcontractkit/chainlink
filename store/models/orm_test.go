package models_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
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

func TestORMSaveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1, initr := cltest.NewJobWithSchedule("* * * * *")
	store.SaveJob(&j1)

	j2, _ := store.FindJob(j1.ID)
	assert.Equal(t, j1.ID, j2.ID)
	assert.NotEqual(t, 0, j2.Initiators[0])
	assert.Equal(t, j2.Initiators[0].ID, j1.Initiators[0].ID)
	assert.Equal(t, j2.ID, j2.Initiators[0].JobID)
	assert.NoError(t, store.One("JobID", j1.ID, &initr))
	assert.Equal(t, models.Cron("* * * * *"), initr.Schedule)
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

func TestBridgeTypeFor(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	tt := models.BridgeType{}
	tt.Name = "solargridreporting"
	u, err := url.Parse("https://denergy.eth")
	assert.NoError(t, err)
	tt.URL = models.WebURL{URL: u}
	assert.NoError(t, store.Save(&tt))

	cases := []struct {
		description string
		name        string
		want        models.BridgeType
		errored     bool
	}{
		{"actual external adapter", tt.Name, tt, false},
		{"core adapter", "ethtx", models.BridgeType{}, true},
		{"non-existent adapter", "nonExistent", models.BridgeType{}, true},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			tt, err := store.BridgeTypeFor(test.name)
			assert.Equal(t, test.want, tt)
			assert.Equal(t, test.errored, err != nil)
		})
	}
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

func TestORM_SaveCreationHeight(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	cases := []struct {
		name            string
		creationHeight  *big.Int
		parameterHeight *big.Int
		wantHeight      *big.Int
	}{
		{"unset", nil, big.NewInt(2), big.NewInt(2)},
		{"set", big.NewInt(1), big.NewInt(2), big.NewInt(1)},
		{"unset and nil", nil, nil, nil},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			jr := job.NewRun(initr)
			if test.creationHeight != nil {
				ch := hexutil.Big(*test.creationHeight)
				jr.CreationHeight = &ch
			}
			assert.NoError(t, store.Save(&jr))

			bn := cltest.IndexableBlockNumber(test.parameterHeight)
			result, err := store.SaveCreationHeight(jr, bn)

			assert.NoError(t, err)
			assert.Equal(t, test.wantHeight, result.CreationHeight.ToInt())
			assert.NoError(t, store.One("ID", jr.ID, &jr))
			assert.Equal(t, test.wantHeight, jr.CreationHeight.ToInt())
		})
	}
}

func TestMarkRan(t *testing.T) {
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

func TestORM_ParseQuery(t *testing.T) {
	tests := []struct {
		jsonFile string
		model    interface{}
		expLen   int
	}{
		{"../../internal/fixtures/web/create_random_number_bridge_type.json", &models.BridgeType{},
			3},
		{"../../internal/fixtures/web/end_at_job.json",
			&models.JobSpec{},
			3},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		params := json.RawMessage(cltest.LoadJSON(test.jsonFile))
		query, err := store.ParseQuery(params, test.model, "Re")
		assert.Error(t, err)
		assert.Equal(t, len(query), 0)
		query, err = store.ParseQuery(params, test.model, "Lte")
		assert.NoError(t, err)
		assert.Equal(t, len(query), test.expLen)
		query, err = store.ParseQuery(params, test.model, "Eq")
		assert.NoError(t, err)
		assert.Equal(t, len(query), test.expLen)
		query, err = store.ParseQuery(params, test.model, "Xx")
		assert.Error(t, err)

	}

	t.Parallel()
}

func TestORM_BuildQuery(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	var params models.QueryObject
	params.Eq = json.RawMessage(`{"url" : "http://www.example.com" }`)
	params.Gte = json.RawMessage(`{ "defaultConfirmations" : 0 }`)
	params.Re = json.RawMessage(`{ "name" : "^test" }`)
	expected_len := 3

	result, err := store.BuildQuery(params, &models.BridgeType{})
	assert.Nil(t, err)
	assert.Equal(t, expected_len, len(result))
	u, err := url.Parse("http://www.example.com")
	assert.Nil(t, err)
	assert.Equal(t, result[0], q.Eq("URL", models.WebURL{URL: u}))
	assert.Equal(t, result[1], q.Gte("DefaultConfirmations", uint64(0)))
	assert.Equal(t, result[2], q.Re("Name", "^test"))

	expected_len = 0
	params = models.QueryObject{Eq: json.RawMessage(`{ "url" : 000 }`)}
	result, err = store.BuildQuery(params, &models.BridgeType{})
	assert.Error(t, err)
	assert.Equal(t, expected_len, len(result))

	params = models.QueryObject{Re: json.RawMessage(`{ "wrongField" : "test" }`)}
	result, err = store.BuildQuery(params, &models.BridgeType{})
	assert.Error(t, err)
	assert.Equal(t, expected_len, len(result))

	params = models.QueryObject{Eq: json.RawMessage(`{ "invalidJSON" : "test" `)}
	result, err = store.BuildQuery(params, &models.BridgeType{})
	assert.Error(t, err)
	assert.Equal(t, expected_len, len(result))

}
func TestORM_AdvancedBridgeSearch(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	var bridges []models.BridgeType

	for i := 0; i < 8; i++ {
		bt := models.BridgeType{Name: fmt.Sprintf("testbridge%v", i),
			URL:                  cltest.WebURL(fmt.Sprintf("https://testing.com/bridges%v", i%2)),
			DefaultConfirmations: uint64(i)}
		bridges = append(bridges, bt)
		err := store.Save(&bt)
		assert.NoError(t, err)
	}

	cases := []struct {
		name         string
		searchParams models.QueryObject
		errored      bool
		expectedLen  int
	}{
		{"value not found",
			models.QueryObject{Eq: json.RawMessage(`{"defaultConfirmations":18}`)},
			true,
			0,
		},
		{"unknown field",
			models.QueryObject{Eq: json.RawMessage(`{"testingBadField":999}`)},
			true,
			0,
		},
		{"empty input, returns all values",
			models.QueryObject{},
			false,
			8,
		},
		{"less than query",
			models.QueryObject{Lt: json.RawMessage(`{"defaultConfirmations":4}`)},
			false,
			4,
		},
		{"regex query",
			models.QueryObject{Re: json.RawMessage(`{"name":"^testbridge[0-3]+"}`)},
			false,
			4,
		},
		{"multiple fields",
			models.QueryObject{Re: json.RawMessage(`{"name":"^test"}`),
				Gt: json.RawMessage(`{"defaultConfirmations":4}`),
				Eq: json.RawMessage(`{"url":"https://testing.com/bridges1"}`)},
			false,
			2,
		},
		{"nested operators",
			models.QueryObject{Re: json.RawMessage(`{"name":"^testbridge[0-3]"}`),
				Or: json.RawMessage(`{"not" : {"re" : {"name":"^testbridge[0-1|6-7]"}}}`),
			},
			false,
			6,
		},
	}

	for _, test := range cases {
		query, err := store.AdvancedBridgeSearch(test.searchParams)
		assert.Equal(t, test.errored, err != nil)
		assert.Equal(t, test.expectedLen, len(query))
	}
}

func TestORM_AdvancedJobRunSearch(t *testing.T) {

	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()
	status := []models.RunStatus{models.RunStatusCompleted,
		models.RunStatusInProgress,
		models.RunStatusErrored}
	for i := 0; i < 6; i++ {
		job, initr := cltest.NewJobWithWebInitiator()
		run := job.NewRun(initr)
		run.Status = status[i%3]
		assert.NoError(t, store.Save(&run))
	}

	cases := []struct {
		name         string
		searchParams models.QueryObject
		errored      bool
		expectedLen  int
	}{
		{"value not found",
			models.QueryObject{Eq: json.RawMessage(`{"id": "zzzzzz"}`)},
			true,
			0,
		},
		{"less than query on date field",
			models.QueryObject{Lt: json.RawMessage(`{"createdAt":"3000-01-01T00:00:00.000Z"}`)},
			false,
			6,
		},
		{"find completed runs",
			models.QueryObject{Eq: json.RawMessage(`{"status":"completed"}`)},
			false,
			2,
		},
		{"find uncompleted runs",
			models.QueryObject{Not: json.RawMessage(`{"Eq" : {"status":"completed"}}`)},
			false,
			4,
		},
		{"find in progress or errored",
			models.QueryObject{Eq: json.RawMessage(`{"status":"in_progress"}`),
				Or: json.RawMessage(`{"eq" : {"status":"errored"}}`)},
			false,
			4,
		},
	}

	for _, test := range cases {
		query, err := store.AdvancedJobRunSearch(test.searchParams)
		assert.Equal(t, test.errored, err != nil)
		assert.Equal(t, test.expectedLen, len(query))
	}
}
