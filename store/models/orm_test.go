package models_test

import (
	"encoding/hex"
	"math/big"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestWhereNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := models.NewJob()
	jobs := []models.JobSpec{j1}

	err := store.Where("ID", "bogus", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestAllNotFound(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	var jobs []models.JobSpec
	err := store.All(&jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORMSaveJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * *")
	store.SaveJob(&j1)

	j2, _ := store.FindJob(j1.ID)
	assert.Equal(t, j1.ID, j2.ID)

	assert.Equal(t, j2.ID, j2.Initiators[0].JobID)

	var initr models.Initiator
	store.One("JobID", j1.ID, &initr)
	assert.Equal(t, models.Cron("* * * * *"), initr.Schedule)
}

func TestPendingJobRuns(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j := models.NewJob()
	assert.Nil(t, store.SaveJob(&j))
	npr := j.NewRun()
	assert.Nil(t, store.Save(&npr))

	pr := j.NewRun()
	pr.Status = models.StatusPending
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

func TestCreatingTx(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	from := common.HexToAddress("0x2C83ACd90367e7E0D3762eA31aC77F18faecE874")
	to := common.HexToAddress("0x4A7d17De4B3eC94c59BF07764d9A6e97d92A547A")
	value := new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)
	nonce := uint64(1232421)
	gasLimit := uint64(50000)
	data, err := hex.DecodeString("0987612345abcdef")
	assert.Nil(t, err)

	_, err = store.CreateTx(from, nonce, to, data, value, gasLimit)
	assert.Nil(t, err)

	txs := []models.Tx{}
	assert.Nil(t, store.Where("Nonce", nonce, &txs))
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
	assert.Nil(t, err)
	tt.URL = models.WebURL{u}
	assert.Nil(t, store.Save(&tt))

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
