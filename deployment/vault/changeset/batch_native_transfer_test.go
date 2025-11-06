package changeset

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
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
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

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
			t.Parallel()

			err := ValidateBatchNativeTransferConfig(env.GetContext(), *env, tt.config)

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
	t.Parallel()

	rt, err := runtime.New(t.Context())
	require.NoError(t, err)

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

	initialTask := runtime.ChangesetTask(SetWhitelistChangeset, initialConfig)
	err = rt.Exec(initialTask)

	require.NoError(t, err)
	require.NotNil(t, rt.State().Outputs[initialTask.ID()].DataStore)

	whitelist, err := GetWhitelistedAddresses(rt.Environment(), []uint64{testChainID})
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

	removeTask := runtime.ChangesetTask(SetWhitelistChangeset, updatedConfig)
	err = rt.Exec(removeTask)

	require.NoError(t, err)
	require.NotNil(t, rt.State().Outputs[removeTask.ID()].DataStore)

	updatedWhitelist, err := GetWhitelistedAddresses(rt.Environment(), []uint64{testChainID})
	require.NoError(t, err)
	require.Len(t, updatedWhitelist[testChainID], 1)
	require.Equal(t, testAddr2, updatedWhitelist[testChainID][0].Address)
}

func TestBatchNativeTransferIntegration(t *testing.T) {
	t.Parallel()

	t.Run("full workflow with MCMS setup", func(t *testing.T) {
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulatedN(t, 2),
		))
		require.NoError(t, err)

		chainSelectors := make([]uint64, 0)
		for chainSel := range rt.Environment().BlockChains.EVMChains() {
			chainSelectors = append(chainSelectors, chainSel)
		}
		require.Len(t, chainSelectors, 2, "Need 2 chains for testing")

		setupMCMSInfrastructure(t, rt, chainSelectors)
		fundDeployerAccounts(t, rt.Environment(), chainSelectors)
		setupWhitelist(t, rt, chainSelectors...)
		fundTimelockContracts(t, rt, chainSelectors...)
		executeBatchTransfersWithMCMS(t, rt, chainSelectors...)
	})

	t.Run("direct execution without MCMS", func(t *testing.T) {
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulatedN(t, 2),
		))
		require.NoError(t, err)

		chainSelectors := getChainSelectors(rt.Environment())
		testChain := chainSelectors[0]
		setupWhitelist(t, rt, testChain)
		executeDirectTransfers(t, rt, testChain)
	})
}

func setupMCMSInfrastructure(t *testing.T, rt *runtime.Runtime, chainSelectors []uint64) {
	t.Log("Setting up MCMS infrastructure with real deployment")

	timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, sel := range chainSelectors {
		t.Logf("Enabling MCMS on chain %d", sel)
		timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
	}

	err := rt.Exec(
		runtime.ChangesetTask(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			timelockCfgs,
		),
	)
	require.NoError(t, err)

	for _, chainSel := range chainSelectors {
		timelockAddr, err := GetContractAddress(rt.State().DataStore, chainSel, commontypes.RBACTimelock)
		require.NoError(t, err)
		t.Logf("Timelock deployed on chain %d with address %s", chainSel, timelockAddr)

		proposerAddr, err := GetContractAddress(rt.State().DataStore, chainSel, commontypes.ProposerManyChainMultisig)
		require.NoError(t, err)
		t.Logf("Proposer deployed on chain %d with address %s", chainSel, proposerAddr)
	}

	t.Log("MCMS deployment completed successfully")
}

// fundDeployerAccounts ensures deployer accounts have sufficient native tokens
func fundDeployerAccounts(t *testing.T, env cldf.Environment, chainSelectors []uint64) {
	t.Log("Funding deployer accounts")

	for _, chainSel := range chainSelectors {
		chain := env.BlockChains.EVMChains()[chainSel]

		balance, err := chain.Client.BalanceAt(testutils.Context(t), chain.DeployerKey.From, nil)
		require.NoError(t, err)

		minBalance := HundredETH
		require.GreaterOrEqual(t, balance.Cmp(minBalance), 0, "Deployer account has insufficient balance on chain %d: balance=%s, required=%s", chainSel, balance.String(), minBalance.String())

		t.Logf("Deployer account funded on chain %d with balance %s", chainSel, balance.String())
	}
}

func setupWhitelist(t *testing.T, rt *runtime.Runtime, chainSelectors ...uint64) {
	t.Log("Setting up whitelist")

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

	setTask := runtime.ChangesetTask(SetWhitelistChangeset, whitelistConfig)
	err := rt.Exec(setTask)

	require.NoError(t, err)
	require.NotNil(t, rt.State().Outputs[setTask.ID()].DataStore)

	totalAddresses := 0
	for _, addresses := range whitelistByChain {
		totalAddresses += len(addresses)
	}

	t.Logf("Whitelist configured | chains: %d, total_addresses: %d", len(whitelistByChain), totalAddresses)
}

// fundTimelockContracts funds timelock contracts with native tokens
func fundTimelockContracts(t *testing.T, rt *runtime.Runtime, chainSelectors ...uint64) {
	t.Log("Funding timelock contracts")

	fundingConfig := types.FundTimelockConfig{
		FundingByChain: make(map[uint64]*big.Int),
	}

	fundingAmount := TenETH

	timelockBalances, err := GetTimelockBalances(rt.Environment(), chainSelectors)
	require.NoError(t, err, "Failed to get timelock balances - contracts may not be deployed")

	for _, chainSel := range chainSelectors {
		balance, exists := timelockBalances[chainSel]
		require.True(t, exists, "Timelock balance info not found for chain %d", chainSel)
		t.Logf("Found timelock to fund on chain %d with address %s and current balance %s", chainSel, balance.TimelockAddr, balance.Balance.String())

		fundingConfig.FundingByChain[chainSel] = fundingAmount
	}

	err = rt.Exec(
		runtime.ChangesetTask(FundTimelockChangeset, fundingConfig),
	)

	require.NoError(t, err)

	timelockBalances, err = GetTimelockBalances(rt.Environment(), chainSelectors)
	require.NoError(t, err, "Failed to get timelock balances - contracts may not be deployed")

	for _, chainSel := range chainSelectors {
		balance, exists := timelockBalances[chainSel]
		require.True(t, exists, "Timelock balance info not found for chain %d", chainSel)
		t.Logf("Found timelock to fund on chain %d with address %s and current balance %s", chainSel, balance.TimelockAddr, balance.Balance.String())

		fundingConfig.FundingByChain[chainSel] = fundingAmount
	}

	t.Log("Timelock contracts funded successfully")
}

func executeBatchTransfersWithMCMS(t *testing.T, rt *runtime.Runtime, chainSelectors ...uint64) {
	t.Log("Executing batch transfers with MCMS")

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

	transferTask := runtime.ChangesetTask(BatchNativeTransferChangeset, transferConfig)
	err := rt.Exec(transferTask)
	require.NoError(t, err)

	output := rt.State().Outputs[transferTask.ID()]

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

	err = rt.Exec(
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)

	t.Log("MCMS proposal executed successfully")

	verifyTransferExecution(t, rt.Environment(), transferConfig, chainSelectors)
}

func executeDirectTransfers(t *testing.T, rt *runtime.Runtime, chainSelector uint64) {
	t.Log("Executing direct transfers without MCMS")

	recipient := common.HexToAddress(testRecipientAddr1)
	transferAmount := OneETH

	chain := rt.Environment().BlockChains.EVMChains()[chainSelector]
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

	output, err := BatchNativeTransferChangeset.Apply(rt.Environment(), transferConfig)
	require.NoError(t, err)

	require.NotNil(t, output.Reports, "Should have execution reports")

	finalBalance, err := chain.Client.BalanceAt(testutils.Context(t), recipient, nil)
	require.NoError(t, err)

	expectedBalance := big.NewInt(0).Add(initialBalance, transferAmount)
	require.Equal(t, expectedBalance, finalBalance, "Recipient balance should increase by transfer amount")

	t.Log("Direct transfers executed and verified successfully")
}

// verifyTransferExecution verifies that transfers were executed by checking recipient balances
func verifyTransferExecution(t *testing.T, env cldf.Environment, config types.BatchNativeTransferConfig, chainSelectors []uint64) {
	lggr := env.Logger
	t.Log("Verifying transfer execution")

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
