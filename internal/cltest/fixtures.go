package cltest

import (
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

func NewJob() models.Job {
	j := models.NewJob()
	j.Tasks = []models.Task{{Type: "NoOp"}}
	return j
}

func NewJobWithSchedule(sched string) models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: "cron", Schedule: models.Cron(sched)}}
	return j
}

func NewJobWithWebInitiator() models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: "web"}}
	return j
}

func NewTx(from string, sentAt uint64) *models.Tx {
	return &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     "deadbeef",
		Value:    big.NewInt(0),
		GasLimit: big.NewInt(250000),
	}
}

func CreateTxAndAttempt(
	store *store.Store,
	from string,
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

func NewTxHash() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hexutil.Encode(b)
}

func NewEthAddress() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hexutil.Encode(b)
}
