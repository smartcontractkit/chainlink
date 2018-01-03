package cltest

import (
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

func NewEthTx(from string, sentAt uint64) *models.EthTx {
	return &models.EthTx{
		From:     from,
		Nonce:    0,
		Data:     "deadbeef",
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(20000000000),
		GasLimit: big.NewInt(250000),
		Attempts: []*models.EthTxAttempt{&models.EthTxAttempt{
			TxID:     NewTxID(),
			GasPrice: big.NewInt(20000000000),
			Hex:      "0x0000",
			SentAt:   sentAt,
		}},
	}
}

func NewTxID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hexutil.Encode(b)
}

func NewEthAddress() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hexutil.Encode(b)
}
