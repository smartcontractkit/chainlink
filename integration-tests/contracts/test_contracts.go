package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
)

type LogEmitterContract struct {
	address  common.Address
	client   *seth.Client
	instance *le.LogEmitter
	l        zerolog.Logger
}

func (e *LogEmitterContract) Address() common.Address {
	return e.address
}

func (e *LogEmitterContract) EmitLogIntsFromKey(ints []int, keyNum int) (*types.Transaction, error) {
	bigInts := make([]*big.Int, len(ints))
	for i, v := range ints {
		bigInts[i] = big.NewInt(int64(v))
	}
	tx, err := e.client.Decode(e.instance.EmitLog1(e.client.NewTXKeyOpts(keyNum), bigInts))
	if err != nil {
		return nil, err
	}

	return tx.Transaction, nil
}

func (e *LogEmitterContract) EmitLogInts(ints []int) (*types.Transaction, error) {
	return e.EmitLogIntsFromKey(ints, 0)
}

func (e *LogEmitterContract) EmitLogIntsIndexedFromKey(ints []int, keyNum int) (*types.Transaction, error) {
	bigInts := make([]*big.Int, len(ints))
	for i, v := range ints {
		bigInts[i] = big.NewInt(int64(v))
	}
	tx, err := e.client.Decode(e.instance.EmitLog2(e.client.NewTXKeyOpts(keyNum), bigInts))
	if err != nil {
		return nil, err
	}

	return tx.Transaction, nil
}

func (e *LogEmitterContract) EmitLogIntsIndexed(ints []int) (*types.Transaction, error) {
	return e.EmitLogIntsIndexedFromKey(ints, 0)
}

func (e *LogEmitterContract) EmitLogIntMultiIndexedFromKey(ints int, ints2 int, count int, keyNum int) (*types.Transaction, error) {
	tx, err := e.client.Decode(e.instance.EmitLog4(e.client.NewTXKeyOpts(keyNum), big.NewInt(int64(ints)), big.NewInt(int64(ints2)), big.NewInt(int64(count))))
	if err != nil {
		return nil, err
	}

	return tx.Transaction, nil
}

func (e *LogEmitterContract) EmitLogIntMultiIndexed(ints int, ints2 int, count int) (*types.Transaction, error) {
	return e.EmitLogIntMultiIndexedFromKey(ints, ints2, count, 0)
}

func (e *LogEmitterContract) EmitLogStringsFromKey(strings []string, keyNum int) (*types.Transaction, error) {
	tx, err := e.client.Decode(e.instance.EmitLog3(e.client.NewTXKeyOpts(keyNum), strings))
	if err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

func (e *LogEmitterContract) EmitLogStrings(strings []string) (*types.Transaction, error) {
	return e.EmitLogStringsFromKey(strings, 0)
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

func DeployLogEmitterContract(l zerolog.Logger, client *seth.Client) (LogEmitter, error) {
	return DeployLogEmitterContractFromKey(l, client, 0)
}

func DeployLogEmitterContractFromKey(l zerolog.Logger, client *seth.Client, keyNum int) (LogEmitter, error) {
	abi, err := le.LogEmitterMetaData.GetAbi()
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("failed to get LogEmitter ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "LogEmitter", *abi, common.FromHex(le.LogEmitterMetaData.Bin))
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("LogEmitter instance deployment have failed: %w", err)
	}

	instance, err := le.NewLogEmitter(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("failed to instantiate LogEmitter instance: %w", err)
	}

	return &LogEmitterContract{
		client:   client,
		instance: instance,
		address:  data.Address,
		l:        l,
	}, err
}
