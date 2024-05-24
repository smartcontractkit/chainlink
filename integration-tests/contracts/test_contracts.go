package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
)

type LegacyLogEmitterContract struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *le.LogEmitter
	l        zerolog.Logger
}

func (e *LegacyLogEmitterContract) Address() common.Address {
	return e.address
}

func (e *LegacyLogEmitterContract) EmitLogInts(ints []int) (*types.Transaction, error) {
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

func (e *LegacyLogEmitterContract) EmitLogIntsIndexed(ints []int) (*types.Transaction, error) {
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

func (e *LegacyLogEmitterContract) EmitLogIntMultiIndexed(ints int, ints2 int, count int) (*types.Transaction, error) {
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

func (e *LegacyLogEmitterContract) EmitLogStrings(strings []string) (*types.Transaction, error) {
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

func (e *LegacyLogEmitterContract) EmitLogInt(payload int) (*types.Transaction, error) {
	return e.EmitLogInts([]int{payload})
}

func (e *LegacyLogEmitterContract) EmitLogIntIndexed(payload int) (*types.Transaction, error) {
	return e.EmitLogIntsIndexed([]int{payload})
}

func (e *LegacyLogEmitterContract) EmitLogString(strings string) (*types.Transaction, error) {
	return e.EmitLogStrings([]string{strings})
}

func (e *LegacyLogEmitterContract) EmitLogIntsFromKey(_ []int, _ int) (*types.Transaction, error) {
	panic("only Seth-based contracts support this method")
}
func (e *LegacyLogEmitterContract) EmitLogIntsIndexedFromKey(_ []int, _ int) (*types.Transaction, error) {
	panic("only Seth-based contracts support this method")
}
func (e *LegacyLogEmitterContract) EmitLogIntMultiIndexedFromKey(_ int, _ int, _ int, _ int) (*types.Transaction, error) {
	panic("only Seth-based contracts support this method")
}
func (e *LegacyLogEmitterContract) EmitLogStringsFromKey(_ []string, _ int) (*types.Transaction, error) {
	panic("only Seth-based contracts support this method")
}
