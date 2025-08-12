package changeset

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	testAddr1               = "0x1234567890123456789012345678901234567890"
	testAddr2               = "0x0987654321098765432109876543210987654321"
	zeroAddr                = "0x0000000000000000000000000000000000000000"
	testRecipientAddr1      = "0x742d35cc64ca395db82e2e3e8fa8bc6d1b7c0832"
	testRecipientMultiChain = "0x123456789012345678901234567890123456789a"
	testChainID             = 11155111
)

var (
	OneETH     = big.NewInt(1000000000000000000)
	TenETH     = big.NewInt(0).Mul(OneETH, big.NewInt(10))
	HundredETH = big.NewInt(0).Mul(OneETH, big.NewInt(100))
)

func TestBatchNativeTransferValidation(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 0,
	})
	// Ensure a non-nil datastore for validation that inspects whitelist state
	env.DataStore = datastore.NewMemoryDataStore().Seal()

	tests := []struct {
		name      string
		config    types.BatchNativeTransferConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "address not whitelisted",
			config: types.BatchNativeTransferConfig{
				TransfersByChain: map[uint64][]types.NativeTransfer{
					testChainID: {
						{
							To:     testAddr1,
							Amount: OneETH,
						},
					},
				},
				Description: "Test transfer",
			},
			wantError: true,
			errorMsg:  "is not whitelisted",
		},
		{
			name: "empty transfers",
			config: types.BatchNativeTransferConfig{
				TransfersByChain: map[uint64][]types.NativeTransfer{},
			},
			wantError: true,
			errorMsg:  "transfers_by_chain must not be empty",
		},
		{
			name: "zero amount transfer",
			config: types.BatchNativeTransferConfig{
				TransfersByChain: map[uint64][]types.NativeTransfer{
					testChainID: {
						{
							To:     testAddr1,
							Amount: big.NewInt(0),
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "amount must be positive",
		},
		{
			name: "zero address transfer",
			config: types.BatchNativeTransferConfig{
				TransfersByChain: map[uint64][]types.NativeTransfer{
					testChainID: {
						{
							To:     zeroAddr,
							Amount: OneETH,
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "'to' address cannot be zero address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBatchNativeTransferConfig(env.GetContext(), env, tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSetWhitelist(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 0,
	})

	initialConfig := types.SetWhitelistConfig{
		WhitelistByChain: map[uint64][]types.WhitelistAddress{
			testChainID: {
				{
					Address:     common.HexToAddress(testAddr1).Hex(),
					Description: "Test address 1",
					Labels:      []string{"team", "approved"},
				},
				{
					Address:     common.HexToAddress(testAddr2).Hex(),
					Description: "Test address 2",
					Labels:      []string{"partner", "approved"},
				},
			},
		},
	}

	output, err := SetWhitelistChangeset.Apply(env, initialConfig)
	require.NoError(t, err)
	require.NotNil(t, output.DataStore)

	env.DataStore = output.DataStore.Seal()

	whitelist, err := GetWhitelistedAddresses(env, []uint64{testChainID})
	require.NoError(t, err)
	require.Len(t, whitelist[testChainID], 2)

	// Test removing one address
	updatedConfig := types.SetWhitelistConfig{
		WhitelistByChain: map[uint64][]types.WhitelistAddress{
			testChainID: {
				{
					Address:     common.HexToAddress(testAddr2).Hex(),
					Description: "Test address 2 - kept",
					Labels:      []string{"partner", "approved"},
				},
			},
		},
	}

	output2, err := SetWhitelistChangeset.Apply(env, updatedConfig)
	require.NoError(t, err)
	require.NotNil(t, output2.DataStore)

	env.DataStore = output2.DataStore.Seal()

	updatedWhitelist, err := GetWhitelistedAddresses(env, []uint64{testChainID})
	require.NoError(t, err)
	require.Len(t, updatedWhitelist[testChainID], 1)
	require.Equal(t, testAddr2, updatedWhitelist[testChainID][0].Address)
}

func TestBatchNativeTransferIntegration(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 2, "Need 2 chains for testing")

	t.Run("full workflow with MCMS setup", func(t *testing.T) {
		env = setupMCMSInfrastructure(t, env, chainSelectors)
		fundDeployerAccounts(t, env, chainSelectors)
		env = setupWhitelist(t, env, chainSelectors...)
		env = fundTimelockContracts(t, env, chainSelectors...)
		executeBatchTransfersWithMCMS(t, env, chainSelectors...)
	})

	t.Run("direct execution without MCMS", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
			Chains: 2,
		})
		chainSelectors := getChainSelectors(env)
		testChain := chainSelectors[0]
		env = setupWhitelist(t, env, testChain)
		executeDirectTransfers(t, env, testChain)
	})
}

func setupMCMSInfrastructure(t *testing.T, env cldf.Environment, chainSelectors []uint64) cldf.Environment {
	lggr := env.Logger
	lggr.Info("Setting up MCMS infrastructure with real deployment")

	timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, sel := range chainSelectors {
		t.Logf("Enabling MCMS on chain %d", sel)
		timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
	}

	updatedEnv, err := commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			timelockCfgs,
		),
	)
	require.NoError(t, err)

	for _, chainSel := range chainSelectors {
		timelockAddr, err := GetContractAddress(updatedEnv.DataStore, chainSel, commontypes.RBACTimelock)
		require.NoError(t, err)
		lggr.Infow("Timelock deployed", "chain", chainSel, "address", timelockAddr)

		proposerAddr, err := GetContractAddress(updatedEnv.DataStore, chainSel, commontypes.ProposerManyChainMultisig)
		require.NoError(t, err)
		lggr.Infow("Proposer deployed", "chain", chainSel, "address", proposerAddr)
	}

	lggr.Info("MCMS deployment completed successfully")
	return updatedEnv
}

// fundDeployerAccounts ensures deployer accounts have sufficient native tokens
func fundDeployerAccounts(t *testing.T, env cldf.Environment, chainSelectors []uint64) {
	lggr := env.Logger
	lggr.Info("Funding deployer accounts")

	for _, chainSel := range chainSelectors {
		chain := env.BlockChains.EVMChains()[chainSel]

		balance, err := chain.Client.BalanceAt(testutils.Context(t), chain.DeployerKey.From, nil)
		require.NoError(t, err)

		minBalance := HundredETH
		require.GreaterOrEqual(t, balance.Cmp(minBalance), 0, "Deployer account has insufficient balance on chain %d: balance=%s, required=%s", chainSel, balance.String(), minBalance.String())

		lggr.Infow("Deployer account funded", "chain", chainSel, "balance", balance.String())
	}
}

func setupWhitelist(t *testing.T, env cldf.Environment, chainSelectors ...uint64) cldf.Environment {
	lggr := env.Logger
	lggr.Info("Setting up whitelist")

	whitelistByChain := make(map[uint64][]types.WhitelistAddress)

	for i, chainSelector := range chainSelectors {
		var addr common.Address
		var description string
		var labels []string

		if i == 0 {
			addr = common.HexToAddress(testRecipientAddr1)
			description = "Test recipient 1"
			labels = []string{"test", "recipient"}
		} else {
			addr = common.HexToAddress(testRecipientMultiChain)
			description = fmt.Sprintf("Test recipient %d", i+1)
			labels = []string{"test", "multi-chain"}
		}

		whitelistByChain[chainSelector] = []types.WhitelistAddress{
			{
				Address:     addr.Hex(),
				Description: description,
				Labels:      labels,
			},
		}
	}

	whitelistConfig := types.SetWhitelistConfig{
		WhitelistByChain: whitelistByChain,
	}

	output, err := SetWhitelistChangeset.Apply(env, whitelistConfig)
	require.NoError(t, err)
	require.NotNil(t, output.DataStore)

	// Merge the new whitelist DataStore with existing DataStore (which contains MCMS addresses)
	mergedDS := datastore.NewMemoryDataStore()
	err = mergedDS.Merge(env.DataStore)
	require.NoError(t, err)
	err = mergedDS.Merge(output.DataStore.Seal())
	require.NoError(t, err)

	env.DataStore = mergedDS.Seal()

	totalAddresses := 0
	for _, addresses := range whitelistByChain {
		totalAddresses += len(addresses)
	}
	lggr.Infow("Whitelist configured", "chains", len(whitelistByChain), "total_addresses", totalAddresses)
	return env
}

// fundTimelockContracts funds timelock contracts with native tokens
func fundTimelockContracts(t *testing.T, env cldf.Environment, chainSelectors ...uint64) cldf.Environment {
	lggr := env.Logger
	lggr.Info("Funding timelock contracts")

	fundingConfig := types.FundTimelockConfig{
		FundingByChain: make(map[uint64]*big.Int),
	}

	fundingAmount := TenETH

	timelockBalances, err := GetTimelockBalances(env, chainSelectors)
	require.NoError(t, err, "Failed to get timelock balances - contracts may not be deployed")

	for _, chainSel := range chainSelectors {
		balance, exists := timelockBalances[chainSel]
		require.True(t, exists, "Timelock balance info not found for chain %d", chainSel)
		lggr.Infow("Found timelock to fund", "chain", chainSel, "address", balance.TimelockAddr, "current_balance", balance.Balance.String())

		fundingConfig.FundingByChain[chainSel] = fundingAmount
	}

	output, err := FundTimelockChangeset.Apply(env, fundingConfig)
	require.NoError(t, err)

	if output.DataStore != nil {
		env.DataStore = output.DataStore.Seal()
	}

	lggr.Info("Timelock contracts funded successfully")
	return env
}

func executeBatchTransfersWithMCMS(t *testing.T, env cldf.Environment, chainSelectors ...uint64) {
	lggr := env.Logger
	lggr.Info("Executing batch transfers with MCMS")

	transferConfig := types.BatchNativeTransferConfig{
		TransfersByChain: make(map[uint64][]types.NativeTransfer),
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay: 0,
		},
		Description: "Integration test batch transfer",
	}

	// Add transfers for each chain - use the same addresses as in the whitelist
	transferAmount := OneETH
	for i, chainSel := range chainSelectors {
		var recipientAddr common.Address
		if i == 0 {
			recipientAddr = common.HexToAddress(testRecipientAddr1)
		} else {
			recipientAddr = common.HexToAddress(testRecipientMultiChain)
		}

		transferConfig.TransfersByChain[chainSel] = []types.NativeTransfer{
			{
				To:     recipientAddr.Hex(),
				Amount: transferAmount,
			},
		}
	}

	output, err := BatchNativeTransferChangeset.Apply(env, transferConfig)
	require.NoError(t, err)
	require.NotEmpty(t, output.MCMSTimelockProposals, "Should create MCMS proposals")

	require.Len(t, output.MCMSTimelockProposals, 1, "Should create exactly 1 MCMS proposal for all chains")

	proposal := output.MCMSTimelockProposals[0]
	require.Len(t, proposal.Operations, len(chainSelectors), "Single proposal should contain operations for all %d chains", len(chainSelectors))

	operationChains := make(map[uint64]bool)
	for _, operation := range proposal.Operations {
		operationChains[uint64(operation.ChainSelector)] = true
	}

	for _, expectedChain := range chainSelectors {
		require.True(t, operationChains[expectedChain], "Proposal should contain operation for chain %d", expectedChain)
	}

	lggr.Infow("MCMS proposals created", "count", len(output.MCMSTimelockProposals), "operations_in_proposal", len(proposal.Operations), "chains_in_proposal", len(operationChains))

	for i, proposal := range output.MCMSTimelockProposals {
		require.NotEmpty(t, proposal.Operations, "Proposal %d should have operations", i)

		lggr.Infow("Executing MCMS proposal", "index", i, "operations", len(proposal.Operations))

		mcmProp := proposalutils.SignMCMSTimelockProposal(t, env, &proposal, false)

		err = proposalutils.ExecuteMCMSProposalV2(t, env, mcmProp)
		require.NoError(t, err, "Failed to execute MCMS proposal %d", i)

		err = proposalutils.ExecuteMCMSTimelockProposalV2(t, env, &proposal)
		require.NoError(t, err, "Failed to execute timelock proposal %d", i)

		lggr.Infow("MCMS proposal executed successfully", "index", i)
	}

	verifyTransferExecution(t, env, transferConfig, chainSelectors)
}

func executeDirectTransfers(t *testing.T, env cldf.Environment, chainSelector uint64) {
	lggr := env.Logger
	lggr.Info("Executing direct transfers without MCMS")

	recipient := common.HexToAddress(testRecipientAddr1)
	transferAmount := OneETH

	chain := env.BlockChains.EVMChains()[chainSelector]
	initialBalance, err := chain.Client.BalanceAt(testutils.Context(t), recipient, nil)
	require.NoError(t, err)

	transferConfig := types.BatchNativeTransferConfig{
		TransfersByChain: map[uint64][]types.NativeTransfer{
			chainSelector: {
				{
					To:     recipient.Hex(),
					Amount: transferAmount,
				},
			},
		},
		MCMSConfig:  nil,
		Description: "Direct transfer test",
	}

	output, err := BatchNativeTransferChangeset.Apply(env, transferConfig)
	require.NoError(t, err)
	require.NotNil(t, output.Reports, "Should have execution reports")

	finalBalance, err := chain.Client.BalanceAt(testutils.Context(t), recipient, nil)
	require.NoError(t, err)

	expectedBalance := big.NewInt(0).Add(initialBalance, transferAmount)
	require.Equal(t, expectedBalance, finalBalance, "Recipient balance should increase by transfer amount")

	lggr.Info("Direct transfers executed and verified successfully")
}

// verifyTransferExecution verifies that transfers were executed by checking recipient balances
func verifyTransferExecution(t *testing.T, env cldf.Environment, config types.BatchNativeTransferConfig, chainSelectors []uint64) {
	lggr := env.Logger
	lggr.Info("Verifying transfer execution")

	evmChains := env.BlockChains.EVMChains()

	for _, chainSel := range chainSelectors {
		chain, exists := evmChains[chainSel]
		require.True(t, exists, "Chain %d should exist", chainSel)

		transfers, exists := config.TransfersByChain[chainSel]
		require.True(t, exists, "Transfers should exist for chain %d", chainSel)

		for i, transfer := range transfers {
			balance, err := chain.Client.BalanceAt(testutils.Context(t), common.HexToAddress(transfer.To), nil)
			require.NoError(t, err)

			require.Equal(t, transfer.Amount, balance,
				"Recipient %s on chain %d (transfer %d) should have exactly %s wei, but has %s wei",
				transfer.To, chainSel, i, transfer.Amount.String(), balance.String())

			lggr.Infow("Transfer verified",
				"chain", chainSel,
				"transfer", i,
				"recipient", transfer.To,
				"amount", transfer.Amount.String(),
				"balance", balance.String())
		}
	}

	lggr.Info("All transfers verified successfully")
}

func getChainSelectors(env cldf.Environment) []uint64 {
	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	return chainSelectors
}
