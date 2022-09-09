package contracts

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/log_emitter"
)

type EventType int

const (
	EventTypeInt EventType = iota
	EventTypeIntIndexed
	EventTypeString
)

// EthereumLogEmitter represents LogEmitter contract used for tests
type EthereumLogEmitter struct {
	address      *common.Address
	client       blockchain.EVMClient
	emitter      *log_emitter.LogEmitter
	mu           *sync.Mutex
	transactions []*types.Transaction
}

// EthereumLogEmitterSharedData shared settings for requests
type EthereumLogEmitterSharedData struct {
	EventType           EventType
	EventsPerRequest    int
	ConfirmTransactions bool
}

func (e *EthereumLogEmitter) EmitLogs(d *EthereumLogEmitterSharedData) error {
	var err error
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	logs, stringLogs := make([]*big.Int, 0), make([]string, 0)
	for i := 0; i < d.EventsPerRequest; i++ {
		logs = append(logs, big.NewInt(20))
	}
	for i := 0; i < d.EventsPerRequest; i++ {
		stringLogs = append(stringLogs, "sjdfhskdjvskjdhfsf")
	}
	var tx *types.Transaction
	switch d.EventType {
	case EventTypeInt:
		tx, err = e.emitter.EmitLog1(opts, logs)
	case EventTypeIntIndexed:
		tx, err = e.emitter.EmitLog2(opts, logs)
	case EventTypeString:
		tx, err = e.emitter.EmitLog3(opts, stringLogs)
	}
	if err != nil {
		return err
	}
	e.mu.Lock()
	e.transactions = append(e.transactions, tx)
	e.mu.Unlock()
	if d.ConfirmTransactions {
		return e.client.ProcessTransaction(tx)
	}
	return nil
}

func (e *EthereumLogEmitter) PerformLoad(data interface{}) error {
	d := data.(*EthereumLogEmitterSharedData)
	return e.EmitLogs(d)
}

func (e *EthereumLogEmitter) GetTransactions() []*types.Transaction {
	return e.transactions
}
