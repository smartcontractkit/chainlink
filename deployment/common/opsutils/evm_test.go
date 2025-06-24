package opsutils

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestCloneTransactOptsWithGas(t *testing.T) {
	orig := &bind.TransactOpts{
		GasLimit: 100,
		GasPrice: big.NewInt(123),
	}
	// Should clone and override both
	cloned := cloneTransactOptsWithGas(orig, 200, 456)
	assert.NotSame(t, orig, cloned)
	assert.Equal(t, uint64(200), cloned.GasLimit)
	assert.Equal(t, big.NewInt(456), cloned.GasPrice)
	// Should not override if zero
	cloned2 := cloneTransactOptsWithGas(orig, 0, 0)
	assert.Equal(t, orig.GasLimit, cloned2.GasLimit)
	assert.Equal(t, orig.GasPrice, cloned2.GasPrice)
	// Nil input
	assert.Nil(t, cloneTransactOptsWithGas(nil, 1, 1))
}

func TestGasBoostConfigsForChainMap(t *testing.T) {
	chainMap := map[uint64]string{1: "a", 2: "b"}
	gasBoostConfigs := map[uint64]commontypes.GasBoostConfig{
		1: {InitialGasLimit: 10},
	}
	cfgs := GasBoostConfigsForChainMap(chainMap, gasBoostConfigs)
	assert.Len(t, cfgs, 2)
	assert.NotNil(t, cfgs[1])
	assert.Nil(t, cfgs[2])
	// Nil configs
	assert.Empty(t, GasBoostConfigsForChainMap[string](chainMap, nil))
	assert.Empty(t, GasBoostConfigsForChainMap[string](nil, gasBoostConfigs))
}

func TestGetBoostedGasForAttempt_DefaultsAndOverrides(t *testing.T) {
	cfg := commontypes.GasBoostConfig{}
	limit, price := getBoostedGasForAttempt(cfg, 0)
	assert.Equal(t, uint64(200_000), limit)
	assert.Equal(t, uint64(20_000_000_000), price)
	limit, price = getBoostedGasForAttempt(cfg, 2)
	assert.Equal(t, uint64(200_000+2*50_000), limit)
	assert.Equal(t, uint64(20_000_000_000+2*10_000_000_000), price)

	cfg = commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	limit, price = getBoostedGasForAttempt(cfg, 3)
	assert.Equal(t, uint64(1000+3*100), limit)
	assert.Equal(t, uint64(2000+3*100), price)
}

func TestRetryDeploymentWithGasBoost(t *testing.T) {
	cfg := &commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	opt := RetryDeploymentWithGasBoost[any](cfg)
	// Should not panic and should be non-nil
	assert.NotNil(t, opt)
	// Should fallback to default if nil
	assert.NotNil(t, RetryDeploymentWithGasBoost[string](nil))
}

func TestAddEVMCallSequenceToCSOutput_SequenceError(t *testing.T) {
	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]EVMCallOutput]{}
	seqErr := errors.New("sequence failed")

	result, err := AddEVMCallSequenceToCSOutput(
		cldf.Environment{},
		csOutput,
		seqReport,
		seqErr,
		nil,
		nil,
		"test",
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute")
	assert.Contains(t, err.Error(), "sequence failed")
	assert.Equal(t, seqReport.ExecutionReports, result.Reports)
}

func TestAddEVMCallSequenceToCSOutput_NoMCMS(t *testing.T) {
	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]EVMCallOutput]{}

	result, err := AddEVMCallSequenceToCSOutput(
		cldf.Environment{},
		csOutput,
		seqReport,
		nil,
		nil,
		nil, // No MCMS config
		"test",
	)

	require.NoError(t, err)
	assert.Equal(t, seqReport.ExecutionReports, result.Reports)
}

func TestAddEVMCallSequenceToCSOutput_AllConfirmed(t *testing.T) {
	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]EVMCallOutput]{}
	mcmsCfg := &proposalutils.TimelockConfig{}

	result, err := AddEVMCallSequenceToCSOutput(
		cldf.Environment{},
		csOutput,
		seqReport,
		nil,
		map[uint64]state.MCMSWithTimelockState{},
		mcmsCfg,
		"test",
	)

	require.NoError(t, err)
	assert.Equal(t, seqReport.ExecutionReports, result.Reports)
	assert.Nil(t, result.MCMSTimelockProposals)
}
func TestNewEVMCallOperation(t *testing.T) {
	version, _ := semver.NewVersion("1.0.0")

	t.Run("ChainSelectorMismatch", func(t *testing.T) {
		op := NewEVMCallOperation[string, any](
			"test",
			version,
			"description",
			"abi",
			cldf.ContractType("TestContract"),
			func(address common.Address, backend bind.ContractBackend) (any, error) {
				return nil, nil
			},
			func(contract any, opts *bind.TransactOpts, input string) (*types.Transaction, error) {
				return nil, nil
			},
		)

		input := EVMCallInput[string]{
			ChainSelector: 123,
			Address:       common.HexToAddress("0x1234"),
		}
		chain := cldf_evm.Chain{Selector: 456}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch between inputted chain selector")
	})

	t.Run("ConstructorError", func(t *testing.T) {
		op := NewEVMCallOperation[string, any](
			"test",
			version,
			"description",
			"abi",
			cldf.ContractType("TestContract"),
			func(address common.Address, backend bind.ContractBackend) (any, error) {
				return nil, errors.New("constructor failed")
			},
			func(contract any, opts *bind.TransactOpts, input string) (*types.Transaction, error) {
				return nil, nil
			},
		)

		input := EVMCallInput[string]{
			ChainSelector: 123,
			Address:       common.HexToAddress("0x1234"),
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create contract instance")
		assert.Contains(t, err.Error(), "constructor failed")
	})

	t.Run("NoSendMode", func(t *testing.T) {
		mockTx := types.NewTransaction(
			0,                             // nonce
			common.HexToAddress("0x1234"), // to address
			big.NewInt(0),                 // value
			21000,                         // gas limit
			big.NewInt(0),                 // gas price
			nil,                           // data
		)
		op := NewEVMCallOperation[string, any](
			"test",
			version,
			"description",
			"abi",
			cldf.ContractType("TestContract"),
			func(address common.Address, backend bind.ContractBackend) (any, error) {
				return struct{}{}, nil
			},
			func(contract any, opts *bind.TransactOpts, input string) (*types.Transaction, error) {
				return mockTx, nil
			},
		)

		input := EVMCallInput[string]{
			ChainSelector: 123,
			Address:       common.HexToAddress("0x1234"),
			NoSend:        true,
			CallInput:     "test input",
		}
		chain := cldf_evm.Chain{Selector: 123}

		output, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.NoError(t, err)
		assert.Equal(t, input.Address, output.Output.To)
		assert.Equal(t, cldf.ContractType("TestContract"), output.Output.ContractType)
		assert.False(t, output.Output.Confirmed)
	})

	t.Run("CustomGasSettings", func(t *testing.T) {
		var capturedOpts *bind.TransactOpts
		mockTx := types.NewTransaction(
			0,                             // nonce
			common.HexToAddress("0x1234"), // to address
			big.NewInt(0),                 // value
			21000,                         // gas limit
			big.NewInt(0),                 // gas price
			nil,                           // data
		)

		op := NewEVMCallOperation[string, any](
			"test",
			version,
			"description",
			"abi",
			cldf.ContractType("TestContract"),
			func(address common.Address, backend bind.ContractBackend) (any, error) {
				return struct{}{}, nil
			},
			func(contract any, opts *bind.TransactOpts, input string) (*types.Transaction, error) {
				capturedOpts = opts
				return mockTx, nil
			},
		)

		input := EVMCallInput[string]{
			ChainSelector: 123,
			Address:       common.HexToAddress("0x1234"),
			GasLimit:      100000,
			GasPrice:      50000000000,
			NoSend:        true, // Use NoSend to avoid confirmation
		}

		deployerKey := &bind.TransactOpts{
			GasLimit: 50000,
			GasPrice: big.NewInt(25000000000),
		}
		chain := cldf_evm.Chain{
			Selector:    123,
			DeployerKey: deployerKey,
		}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.NoError(t, err)

		// In NoSend mode, SimTransactOpts are used instead of custom gas
		assert.NotNil(t, capturedOpts.Signer)
	})
}
func TestNewEVMDeployOperation(t *testing.T) {
	version, _ := semver.NewVersion("1.0.0")
	typeAndVersion := cldf.TypeAndVersion{Type: "TestContract", Version: *version}

	t.Run("ChainSelectorMismatch", func(t *testing.T) {
		deployers := VMDeployers[string]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput string) (common.Address, *types.Transaction, error) {
				return common.Address{}, nil, nil
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 456}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch between inputted chain selector")
	})

	t.Run("EVMDeploymentError", func(t *testing.T) {
		deployers := VMDeployers[string]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput string) (common.Address, *types.Transaction, error) {
				return common.Address{}, nil, errors.New("deployment failed")
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123, IsZkSyncVM: false}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deploy")
		assert.Contains(t, err.Error(), "deployment failed")
	})

	t.Run("ZkSyncVMDeploymentError", func(t *testing.T) {
		deployers := VMDeployers[string]{
			DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, deployInput string) (common.Address, error) {
				return common.Address{}, errors.New("zksync deployment failed")
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123, IsZkSyncVM: true}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deploy")
		assert.Contains(t, err.Error(), "zksync deployment failed")
	})

	t.Run("EVMSuccessfulDeployment", func(t *testing.T) {
		expectedAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
		mockTx := types.NewTransaction(
			0,
			common.HexToAddress("0x1234"),
			big.NewInt(0),
			21000,
			big.NewInt(0),
			nil,
		)

		deployers := VMDeployers[string]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput string) (common.Address, *types.Transaction, error) {
				return expectedAddr, mockTx, nil
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}

		// Mock chain with confirmation method
		chain := cldf_evm.Chain{
			Selector:   123,
			IsZkSyncVM: false,
		}
		// Override Confirm method to avoid nil pointer
		chain.Confirm = func(tx *types.Transaction) (uint64, error) {
			return 0, nil
		}

		output, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.NoError(t, err)
		assert.Equal(t, expectedAddr, output.Output.Address)
		assert.Equal(t, typeAndVersion.String(), output.Output.TypeAndVersion)
	})

	t.Run("ZkSyncVMSuccessfulDeployment", func(t *testing.T) {
		expectedAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

		deployers := VMDeployers[string]{
			DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, deployInput string) (common.Address, error) {
				return expectedAddr, nil
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123, IsZkSyncVM: true}

		output, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.NoError(t, err)
		assert.Equal(t, expectedAddr, output.Output.Address)
		assert.Equal(t, typeAndVersion.String(), output.Output.TypeAndVersion)
	})

	t.Run("EVMConfirmationError", func(t *testing.T) {
		expectedAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
		mockTx := types.NewTransaction(
			0,
			common.HexToAddress("0x1234"),
			big.NewInt(0),
			21000,
			big.NewInt(0),
			nil,
		)

		deployers := VMDeployers[string]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput string) (common.Address, *types.Transaction, error) {
				return expectedAddr, mockTx, nil
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}

		chain := cldf_evm.Chain{
			Selector:   123,
			IsZkSyncVM: false,
		}
		// Mock confirmation failure
		chain.Confirm = func(tx *types.Transaction) (uint64, error) {
			return 1, errors.New("confirmation failed")
		}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to confirm deployment")
		assert.Contains(t, err.Error(), "confirmation failed")
	})

	t.Run("CustomGasSettings", func(t *testing.T) {
		var capturedOpts *bind.TransactOpts
		expectedAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
		mockTx := types.NewTransaction(
			0,
			common.HexToAddress("0x1234"),
			big.NewInt(0),
			21000,
			big.NewInt(0),
			nil,
		)

		deployers := VMDeployers[string]{
			DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput string) (common.Address, *types.Transaction, error) {
				capturedOpts = opts
				return expectedAddr, mockTx, nil
			},
		}

		op := NewEVMDeployOperation[string](
			"test",
			version,
			"description",
			typeAndVersion,
			deployers,
		)

		input := EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
			GasLimit:      100000,
			GasPrice:      50000000000,
		}

		deployerKey := &bind.TransactOpts{
			GasLimit: 50000,
			GasPrice: big.NewInt(25000000000),
		}
		chain := cldf_evm.Chain{
			Selector:    123,
			IsZkSyncVM:  false,
			DeployerKey: deployerKey,
		}
		chain.Confirm = func(tx *types.Transaction) (uint64, error) {
			return 0, nil
		}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.NoError(t, err)

		// Verify custom gas settings were applied
		assert.Equal(t, uint64(100000), capturedOpts.GasLimit)
		assert.Equal(t, big.NewInt(50000000000), capturedOpts.GasPrice)
	})
}
