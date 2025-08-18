package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var FundTimelockChangeset = cldf.CreateChangeSet(fundTimelockLogic, fundTimelockPrecondition)

func fundTimelockPrecondition(e cldf.Environment, cfg types.FundTimelockConfig) error {
	ctx := e.GetContext()
	return ValidateFundTimelockConfig(ctx, e, cfg)
}

func fundTimelockLogic(e cldf.Environment, cfg types.FundTimelockConfig) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	lggr.Infow("Funding timelock contracts",
		"chains", len(cfg.FundingByChain))

	evmChains := e.BlockChains.EVMChains()

	for chainSelector, amount := range cfg.FundingByChain {
		chain, exists := evmChains[chainSelector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		mutableDS := datastore.NewMemoryDataStore()
		if e.DataStore != nil {
			if err := mutableDS.Merge(e.DataStore); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore: %w", err)
			}
		}

		deps := VaultDeps{
			Chain:     chain,
			Auth:      chain.DeployerKey,
			DataStore: mutableDS,
		}

		fundInput := FundTimelockInput{
			ChainSelector: chainSelector,
			Amount:        amount,
		}

		fundReport, err := operations.ExecuteOperation(
			e.OperationsBundle, FundTimelockOp, deps, fundInput,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fund timelock on chain %d: %w", chainSelector, err)
		}

		lggr.Infow("Timelock funded successfully",
			"chain", chainSelector,
			"amount", amount.String(),
			"txHash", fundReport.Output.TxHash.Hex())
	}

	lggr.Infow("All timelock funding completed successfully",
		"chains", len(cfg.FundingByChain))

	var outputDataStore datastore.MutableDataStore
	if e.DataStore != nil {
		finalDS := datastore.NewMemoryDataStore()
		if err := finalDS.Merge(e.DataStore); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore: %w", err)
		}
		outputDataStore = finalDS
	}

	return cldf.ChangesetOutput{
		DataStore: outputDataStore,
	}, nil
}

func GetTimelockBalances(e cldf.Environment, chainSelectors []uint64) (map[uint64]*types.TimelockNativeBalanceInfo, error) {
	balances := make(map[uint64]*types.TimelockNativeBalanceInfo)
	evmChains := e.BlockChains.EVMChains()

	for _, chainSelector := range chainSelectors {
		chain, exists := evmChains[chainSelector]
		if !exists {
			return nil, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		timelockAddr, err := GetContractAddress(e.DataStore, chainSelector, commontypes.RBACTimelock)
		if err != nil {
			e.Logger.Debugf("Timelock not found for chain %d, skipping balance check: %v", chainSelector, err)
			continue
		}

		timelockAddress := common.HexToAddress(timelockAddr)

		balance, err := chain.Client.BalanceAt(e.GetContext(), timelockAddress, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get timelock balance for chain %d: %w", chainSelector, err)
		}

		balances[chainSelector] = &types.TimelockNativeBalanceInfo{
			TimelockAddr: timelockAddress.Hex(),
			Balance:      balance,
		}
	}

	return balances, nil
}

// CalculateFundingRequirements calculates how much funding each timelock needs for planned transfers
func CalculateFundingRequirements(e cldf.Environment, cfg types.BatchNativeTransferConfig) (map[uint64]*FundingRequirement, error) {
	requirements := make(map[uint64]*FundingRequirement)

	chainSelectors := make([]uint64, 0, len(cfg.TransfersByChain))
	for chainSelector := range cfg.TransfersByChain {
		chainSelectors = append(chainSelectors, chainSelector)
	}

	currentBalances, err := GetTimelockBalances(e, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to get current timelock balances: %w", err)
	}

	for chainSelector, transfers := range cfg.TransfersByChain {
		totalRequired := big.NewInt(0)
		for _, transfer := range transfers {
			totalRequired.Add(totalRequired, transfer.Amount)
		}

		currentBalance := currentBalances[chainSelector].Balance

		requirements[chainSelector] = &FundingRequirement{
			ChainSelector:  chainSelector,
			CurrentBalance: currentBalance,
			RequiredAmount: totalRequired,
			TransferCount:  len(transfers),
		}
	}

	return requirements, nil
}

type FundingRequirement struct {
	ChainSelector  uint64   `json:"chain_selector"`
	CurrentBalance *big.Int `json:"current_balance"`
	RequiredAmount *big.Int `json:"required_amount"`
	TransferCount  int      `json:"transfer_count"`
}
