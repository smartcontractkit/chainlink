package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	le "github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
)

/**
 * LogEmitterContract: A wrapper for interaction with the LogEmitter smart contract.
 * Used primarily in integration tests to trigger EVM events for testing data consumption.
 */
type LogEmitterContract struct {
	address  common.Address
	client   *seth.Client
	instance *le.LogEmitter
	l        zerolog.Logger
}

func (e *LogEmitterContract) Address() common.Address {
	return e.address
}

[Image of Ethereum Smart Contract Event Emission and Logging Architecture]

/**
 * EmitLogIntsFromKey: Emits a log with a slice of integers using a specific transaction key.
 * Converts standard int slice to []*big.Int for EVM compatibility.
 */
func (e *LogEmitterContract) EmitLogIntsFromKey(ints []int, keyNum int) (*types.Transaction, error) {
	bigInts := make([]*big.Int, len(ints))
	for i, v := range ints {
		bigInts[i] = big.NewInt(int64(v))
	}
	
	// Decode is used to unwrap Seth's internal transaction representation to standard Geth types.
	tx, err := e.client.Decode(e.instance.EmitLog1(e.client.NewTXKeyOpts(keyNum), bigInts))
	if err != nil {
		return nil, fmt.Errorf("failed to emit log ints: %w", err)
	}

	return tx.Transaction, nil
}

func (e *LogEmitterContract) EmitLogInts(ints []int) (*types.Transaction, error) {
	return e.EmitLogIntsFromKey(ints, 0)
}

/**
 * EmitLogIntMultiIndexed: Triggers an event with multiple indexed parameters.
 * Indexed parameters are critical for off-chain filtering via 'Topics'.
 */
func (e *LogEmitterContract) EmitLogIntMultiIndexedFromKey(ints int, ints2 int, count int, keyNum int) (*types.Transaction, error) {
	tx, err := e.client.Decode(e.instance.EmitLog4(
		e.client.NewTXKeyOpts(keyNum), 
		big.NewInt(int64(ints)), 
		big.NewInt(int64(ints2)), 
		big.NewInt(int64(count)),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to emit multi-indexed logs: %w", err)
	}

	return tx.Transaction, nil
}

[Image of Ethereum Event Logs and Topics structure]

/**
 * DeployLogEmitterContractFromKey: Handles the deployment of a new LogEmitter contract instance.
 * It leverages Seth Client for gas management and transaction signing.
 */
func DeployLogEmitterContractFromKey(l zerolog.Logger, client *seth.Client, keyNum int) (LogEmitter, error) {
	abi, err := le.LogEmitterMetaData.GetAbi()
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("failed to get LogEmitter ABI: %w", err)
	}

	// DeployContract performs the deployment and waits for confirmation.
	data, err := client.DeployContract(
		client.NewTXKeyOpts(keyNum), 
		"LogEmitter", 
		*abi, 
		common.FromHex(le.LogEmitterMetaData.Bin),
	)
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("deployment failed: %w", err)
	}

	instance, err := le.NewLogEmitter(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &LogEmitterContract{}, fmt.Errorf("failed to instantiate LogEmitter: %w", err)
	}

	l.Info().Str("Address", data.Address.Hex()).Msg("LogEmitter Contract Deployed")

	return &LogEmitterContract{
		client:   client,
		instance: instance,
		address:  data.Address,
		l:        l,
	}, nil
}
