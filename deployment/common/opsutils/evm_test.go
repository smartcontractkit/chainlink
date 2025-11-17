package opsutils_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_evm_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestCloneTransactOptsWithGas(t *testing.T) {
	t.Parallel()
	orig := &bind.TransactOpts{
		GasLimit: 100,
		GasPrice: big.NewInt(123),
	}
	// Should clone and override both
	cloned := opsutils.CloneTransactOptsWithGas(orig, 200, 456)
	assert.NotSame(t, orig, cloned)
	assert.Equal(t, uint64(200), cloned.GasLimit)
	assert.Equal(t, big.NewInt(456), cloned.GasPrice)
	// Should not override if zero
	cloned2 := opsutils.CloneTransactOptsWithGas(orig, 0, 0)
	assert.Equal(t, orig.GasLimit, cloned2.GasLimit)
	assert.Equal(t, orig.GasPrice, cloned2.GasPrice)
	// Nil input
	assert.Nil(t, opsutils.CloneTransactOptsWithGas(nil, 1, 1))
}

func TestGasBoostConfigsForChainMap(t *testing.T) {
	t.Parallel()
	chainMap := map[uint64]string{1: "a", 2: "b"}
	gasBoostConfigs := map[uint64]commontypes.GasBoostConfig{
		1: {InitialGasLimit: 10},
	}
	cfgs := opsutils.GasBoostConfigsForChainMap(chainMap, gasBoostConfigs)
	assert.Len(t, cfgs, 2)
	assert.NotNil(t, cfgs[1])
	assert.Nil(t, cfgs[2])
	// Nil configs
	assert.Empty(t, opsutils.GasBoostConfigsForChainMap[string](chainMap, nil))
	assert.Empty(t, opsutils.GasBoostConfigsForChainMap[string](nil, gasBoostConfigs))
}

func TestGetBoostedGasForAttempt_DefaultsAndOverrides(t *testing.T) {
	t.Parallel()
	cfg := commontypes.GasBoostConfig{}
	limit, price := opsutils.GetBoostedGasForAttempt(cfg, 0)
	assert.Equal(t, uint64(200_000), limit)
	assert.Equal(t, uint64(20_000_000_000), price)
	limit, price = opsutils.GetBoostedGasForAttempt(cfg, 2)
	assert.Equal(t, uint64(200_000+2*50_000), limit)
	assert.Equal(t, uint64(20_000_000_000+2*10_000_000_000), price)

	cfg = commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	limit, price = opsutils.GetBoostedGasForAttempt(cfg, 3)
	assert.Equal(t, uint64(1000+3*100), limit)
	assert.Equal(t, uint64(2000+3*100), price)
}

func TestRetryDeploymentWithGasBoost(t *testing.T) {
	t.Parallel()
	cfg := &commontypes.GasBoostConfig{
		InitialGasLimit:   1000,
		GasLimitIncrement: 100,
		InitialGasPrice:   2000,
		GasPriceIncrement: 100,
	}
	opt := opsutils.RetryDeploymentWithGasBoost[any](cfg)
	// Should not panic and should be non-nil
	assert.NotNil(t, opt)
	// Should fallback to default if nil
	assert.NotNil(t, opsutils.RetryDeploymentWithGasBoost[string](nil))
}

func TestAddEVMCallSequenceToCSOutput_SequenceError(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedN(t, 1),
	)
	require.NoError(t, err)

	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]opsutils.EVMCallOutput]{}
	seqErr := errors.New("sequence failed")

	result, err := opsutils.AddEVMCallSequenceToCSOutput(
		*env,
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
	t.Parallel()

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedN(t, 1),
	)
	require.NoError(t, err)

	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]opsutils.EVMCallOutput]{}

	result, err := opsutils.AddEVMCallSequenceToCSOutput(
		*env,
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
	t.Parallel()

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedN(t, 1),
	)
	require.NoError(t, err)

	csOutput := cldf.ChangesetOutput{}
	seqReport := operations.SequenceReport[string, map[uint64][]opsutils.EVMCallOutput]{}
	mcmsCfg := &proposalutils.TimelockConfig{}

	result, err := opsutils.AddEVMCallSequenceToCSOutput(
		*env,
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

func TestAddEVMCallSequenceToCSOutput_ProposalCombination(t *testing.T) {
	t.Parallel()
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
	)
	env := deployedEnvironment.Env

	// Create initial changeset output with existing proposals to test combination logic
	existingProposal1 := mcmslib.TimelockProposal{
		BaseProposal: mcmslib.BaseProposal{
			Description: "First proposal",
		},
		Operations: []mcmstypes.BatchOperation{
			{
				ChainSelector: mcmstypes.ChainSelector(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]),
				Transactions: []mcmstypes.Transaction{
					{
						To:               common.HexToAddress("0x1111111111111111111111111111111111111111").String(),
						Data:             []byte("data1"),
						AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
					},
				},
			},
		},
	}

	existingProposal2 := mcmslib.TimelockProposal{
		BaseProposal: mcmslib.BaseProposal{
			Description: "Second proposal",
		},
		Operations: []mcmstypes.BatchOperation{
			{
				ChainSelector: mcmstypes.ChainSelector(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]),
				Transactions: []mcmstypes.Transaction{
					{
						To:               common.HexToAddress("0x1111112222222222222222222222222222222222").String(),
						Data:             []byte("data2"),
						AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
					},
				},
			},
		},
	}

	csOutput := cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{
			existingProposal1,
			existingProposal2,
		},
	}

	// Create sequence report with unconfirmed calls to generate a new proposal
	chainSel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]
	seqReport := operations.SequenceReport[string, map[uint64][]opsutils.EVMCallOutput]{
		Report: operations.Report[string, map[uint64][]opsutils.EVMCallOutput]{
			Output: map[uint64][]opsutils.EVMCallOutput{
				chainSel: {
					{
						To:           common.HexToAddress("0x3333333333333333333333333333333333333333"),
						Data:         []byte("new_call_data"),
						ContractType: "TestContract",
						Confirmed:    false, // This will create a new proposal
					},
				},
			},
		},
	}

	mcmsCfg := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second, // No delay for testing
		MCMSAction: mcmstypes.TimelockActionSchedule,
	}

	mcmsDescription := "Third proposal"
	// Load onchain state
	chainState, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	t.Logf("mcms state: %+v", chainState.EVMMCMSStateByChain())

	result, err := opsutils.AddEVMCallSequenceToCSOutput(
		env,
		csOutput,
		seqReport,
		nil,
		chainState.EVMMCMSStateByChain(),
		mcmsCfg,
		mcmsDescription,
	)

	require.NoError(t, err)
	assert.Equal(t, seqReport.ExecutionReports, result.Reports)

	// Test the key combination logic:
	// 1. Should have exactly 1 proposal after aggregation
	assert.Len(t, result.MCMSTimelockProposals, 1, "Expected exactly 1 aggregated proposal")

	// 2. Description should be comma-separated combination of all proposals
	aggregatedProposal := result.MCMSTimelockProposals[0]
	expectedDescription := "First proposal, Second proposal, Third proposal"
	assert.Equal(t, expectedDescription, aggregatedProposal.Description,
		"Aggregated proposal should have comma-separated descriptions")

	// 3. Operations should be combined from all proposals
	assert.NotEmpty(t, aggregatedProposal.Operations,
		"Aggregated proposal should contain operations")
}

func TestNewEVMCallOperation(t *testing.T) {
	t.Parallel()
	version, _ := semver.NewVersion("1.0.0")

	t.Run("ChainSelectorMismatch", func(t *testing.T) {
		op := opsutils.NewEVMCallOperation(
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

		input := opsutils.EVMCallInput[string]{
			ChainSelector: 123,
			Address:       common.HexToAddress("0x1234"),
		}
		chain := cldf_evm.Chain{Selector: 456}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch between inputted chain selector")
	})

	t.Run("ConstructorError", func(t *testing.T) {
		op := opsutils.NewEVMCallOperation[string, any](
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

		input := opsutils.EVMCallInput[string]{
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
		op := opsutils.NewEVMCallOperation[string, any](
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

		input := opsutils.EVMCallInput[string]{
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

		op := opsutils.NewEVMCallOperation[string, any](
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

		input := opsutils.EVMCallInput[string]{
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

func TestContractOpts_Validate(t *testing.T) {
	tests := []struct {
		desc       string
		opts       *opsutils.ContractOpts
		isZkSyncVM bool
		err        string
	}{
		{
			desc: "valid evm opts",
			opts: &opsutils.ContractOpts{
				Version:     semver.MustParse("1.0.0"),
				EVMBytecode: []byte{0x01, 0x02, 0x03},
			},
			isZkSyncVM: false,
		},
		{
			desc: "valid zksyncvm opts",
			opts: &opsutils.ContractOpts{
				Version:          semver.MustParse("1.0.0"),
				ZkSyncVMBytecode: []byte{0x05, 0x06, 0x07, 0x08},
			},
			isZkSyncVM: true,
		},
		{
			desc: "nil version",
			opts: &opsutils.ContractOpts{},
			err:  "version must be defined",
		},
		{
			desc: "missing evm bytecode",
			opts: &opsutils.ContractOpts{
				Version: semver.MustParse("1.0.0"),
			},
			isZkSyncVM: false,
			err:        "evm bytecode must be defined",
		},
		{
			desc: "missing zkSyncVM bytecode",
			opts: &opsutils.ContractOpts{
				Version: semver.MustParse("1.0.0"),
			},
			isZkSyncVM: true,
			err:        "zkSyncVM bytecode must be defined",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := test.opts.Validate(test.isZkSyncVM)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.err)
			}
		})
	}
}

func TestNewEVMDeployOperation(t *testing.T) {
	t.Parallel()
	contractType := cldf.ContractType("TestContract")
	version, _ := semver.NewVersion("1.0.0")

	t.Run("ChainSelectorMismatch", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			nil,
			nil,
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 456}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch between inputted chain selector")
	})

	t.Run("ContractMetadata undefined", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			nil,
			nil,
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "contract metadata must be provided for deployment")
	})

	t.Run("ContractOpts not defined", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{},
			nil,
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must define ContractOpts for deployment, no defaults provided")
	})

	t.Run("Default ContractOpts not valid", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{},
			&opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      nil,
				ZkSyncVMBytecode: []byte{0x05, 0x06, 0x07, 0x08},
			},
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			DeployInput:   "test",
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ContractOpts: evm bytecode must be defined")
	})

	t.Run("Inputted ContractOpts not valid", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{},
			nil,
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			ContractOpts: &opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      nil,
				ZkSyncVMBytecode: []byte{0x05, 0x06, 0x07, 0x08},
			},
			DeployInput: "test",
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ContractOpts: evm bytecode must be defined")
	})

	t.Run("ABI parsing failure", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{},
			nil,
			func(string) []any { return nil },
		)

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 123,
			ContractOpts: &opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      []byte{0x01, 0x02, 0x03, 0x04},
				ZkSyncVMBytecode: []byte{0x05, 0x06, 0x07, 0x08},
			},
			DeployInput: "test",
		}
		chain := cldf_evm.Chain{Selector: 123}

		_, err := operations.ExecuteOperation(optest.NewBundle(t), op, chain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse ABI")
	})

	t.Run("EVM deployment failure", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{ABI: "[]"},
			nil,
			func(string) []any { return nil },
		)

		chain, err := cldf_evm_provider.NewSimChainProvider(t, 5009297550715157269,
			cldf_evm_provider.SimChainProviderConfig{
				NumAdditionalAccounts: 1,
			},
		).Initialize(t.Context())
		require.NoError(t, err, "Failed to create SimChainProvider")

		chains := cldf_chain.NewBlockChainsFromSlice(
			[]cldf_chain.BlockChain{chain},
		)
		evmChain := chains.EVMChains()[5009297550715157269]

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 5009297550715157269,
			ContractOpts: &opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      []byte{0x01, 0x02, 0x03, 0x04},
				ZkSyncVMBytecode: []byte{0x05, 0x06, 0x07, 0x08},
			},
			DeployInput: "test",
		}

		_, err = operations.ExecuteOperation(optest.NewBundle(t), op, evmChain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deploy")
	})

	t.Run("EVM confirmation failure", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{ABI: "[]"},
			nil,
			func(string) []any { return nil },
		)

		chain, err := cldf_evm_provider.NewSimChainProvider(t, 5009297550715157269,
			cldf_evm_provider.SimChainProviderConfig{
				NumAdditionalAccounts: 1,
			},
		).Initialize(t.Context())
		require.NoError(t, err, "Failed to create SimChainProvider")

		chains := cldf_chain.NewBlockChainsFromSlice(
			[]cldf_chain.BlockChain{chain},
		)
		evmChain := chains.EVMChains()[5009297550715157269]
		evmChain.Confirm = func(tx *types.Transaction) (uint64, error) {
			return 0, errors.New("confirmation failed")
		}

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 5009297550715157269,
			ContractOpts: &opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      []byte{0x00},
				ZkSyncVMBytecode: []byte{0x00},
			},
			DeployInput: "test",
		}

		_, err = operations.ExecuteOperation(optest.NewBundle(t), op, evmChain, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "confirmation failed")
	})

	t.Run("EVM deployment success", func(t *testing.T) {
		op := opsutils.NewEVMDeployOperation(
			"test",
			version,
			"description",
			contractType,
			&bind.MetaData{ABI: "[]"},
			nil,
			func(string) []any { return nil },
		)

		chain, err := cldf_evm_provider.NewSimChainProvider(t, 5009297550715157269,
			cldf_evm_provider.SimChainProviderConfig{
				NumAdditionalAccounts: 1,
			},
		).Initialize(t.Context())
		require.NoError(t, err, "Failed to create SimChainProvider")

		chains := cldf_chain.NewBlockChainsFromSlice(
			[]cldf_chain.BlockChain{chain},
		)
		evmChain := chains.EVMChains()[5009297550715157269]

		input := opsutils.EVMDeployInput[string]{
			ChainSelector: 5009297550715157269,
			ContractOpts: &opsutils.ContractOpts{
				Version:          semver.MustParse("0.1.0"),
				EVMBytecode:      []byte{0x00},
				ZkSyncVMBytecode: []byte{0x00},
			},
			DeployInput: "test",
		}

		_, err = operations.ExecuteOperation(optest.NewBundle(t), op, evmChain, input)
		require.NoError(t, err)
	})
}
