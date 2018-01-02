package models_test

import (
	"math/big"
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

func TestCreatingEthTx(t *testing.T) {
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	data := "0987612345abcdef"
	from := "0x2C83ACd90367e7E0D3762eA31aC77F18faecE874"
	to := "0x4A7d17De4B3eC94c59BF07764d9A6e97d92A547A"
	value := new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)
	nonce := uint64(1232421)
	gasLimit := big.NewInt(500000)
	gasPrice := big.NewInt(20 * 1000000000)

	_, err := store.CreateEthTx(from, nonce, to, data, value, gasLimit, gasPrice)
	assert.Nil(t, err)

	txs := []models.EthTx{}
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
	assert.Equal(t, gasPrice, tx.GasPrice)
}
