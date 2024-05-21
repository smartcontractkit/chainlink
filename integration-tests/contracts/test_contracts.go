package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
)

type LogEmitterContract struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *le.LogEmitter
	l        zerolog.Logger
}

func (e *LogEmitterContract) Address() common.Address {
	return e.address
}

func (e *LogEmitterContract) EmitLogInts(ints []int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	bigInts := make([]*big.Int, len(ints))
	for i, v := range ints {
		bigInts[i] = big.NewInt(int64(v))
	}
	tx, err := e.instance.EmitLog1(opts, bigInts)
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *LogEmitterContract) EmitLogIntsIndexed(ints []int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	bigInts := make([]*big.Int, len(ints))
	for i, v := range ints {
		bigInts[i] = big.NewInt(int64(v))
	}
	tx, err := e.instance.EmitLog2(opts, bigInts)
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *LogEmitterContract) EmitLogIntMultiIndexed(ints int, ints2 int, count int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.EmitLog4(opts, big.NewInt(int64(ints)), big.NewInt(int64(ints2)), big.NewInt(int64(count)))
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *LogEmitterContract) EmitLogStrings(strings []string) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.EmitLog3(opts, strings)
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *LogEmitterContract) EmitLogInt(payload int) (*types.Transaction, error) {
	return e.EmitLogInts([]int{payload})
}

func (e *LogEmitterContract) EmitLogIntIndexed(payload int) (*types.Transaction, error) {
	return e.EmitLogIntsIndexed([]int{payload})
}

func (e *LogEmitterContract) EmitLogString(strings string) (*types.Transaction, error) {
	return e.EmitLogStrings([]string{strings})
}
