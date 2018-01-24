package cltest

import (
	"crypto/rand"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

func NewJob() *models.Job {
	j := models.NewJob()
	j.Tasks = []models.Task{{Type: "NoOp"}}
	return j
}

func NewJobWithSchedule(sched string) *models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorCron, Schedule: models.Cron(sched)}}
	return j
}

func NewJobWithWebInitiator() *models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorWeb}}
	return j
}

func NewJobWithLogInitiator() *models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		Type:    models.InitiatorEthLog,
		Address: NewEthAddress(),
	}}
	return j
}

func NewTx(from common.Address, sentAt uint64) *models.Tx {
	return &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     []byte{},
		Value:    big.NewInt(0),
		GasLimit: big.NewInt(250000),
	}
}

func CreateTxAndAttempt(
	store *store.Store,
	from common.Address,
	sentAt uint64,
) *models.Tx {
	tx := NewTx(from, sentAt)
	if err := store.Save(tx); err != nil {
		logger.Fatal(err)
	}
	_, err := store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	if err != nil {
		logger.Fatal(err)
	}
	return tx
}

func NewTxHash() common.Hash {
	b := make([]byte, 32)
	rand.Read(b)
	return common.BytesToHash(b)
}

func NewEthAddress() common.Address {
	b := make([]byte, 20)
	rand.Read(b)
	return common.BytesToAddress(b)
}

func NewCustomTaskType() *models.CustomTaskType {
	tt := &models.CustomTaskType{}
	tt.Name = "fixtureCustomTaskType"
	u, err := url.Parse("https://d.example.eth")
	if err != nil {
		panic(err)
	}
	tt.URL = models.WebURL{u}
	return tt
}
