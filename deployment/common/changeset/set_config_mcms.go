package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type ConfigPerRole struct {
	Proposer  config.Config
	Canceller config.Config
	Bypasser  config.Config
}
type TimelockConfig struct {
	MinDelay time.Duration // delay for timelock worker to execute the transfers.
}

type MCMSConfig struct {
	ConfigsPerChain map[uint64]ConfigPerRole
	ProposalConfig  *TimelockConfig
}

var _ deployment.ChangeSet[MCMSConfig] = SetConfigMCMS

// Validate checks that the MCMSConfig is valid
func (cfg MCMSConfig) Validate(e deployment.Environment, selectors []uint64) error {
	if len(cfg.ConfigsPerChain) == 0 {
		return errors.New("no chain configs provided")
	}
	// configs should have at least one chain
	state, err := MaybeLoadMCMSWithTimelockState(e, selectors)
	if err != nil {
		return err
	}
	for chainSelector, c := range cfg.ConfigsPerChain {
		family, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return err
		}
		if family != chain_selectors.FamilyEVM {
			return fmt.Errorf("chain selector: %d is not an ethereum chain", chainSelector)
		}
		_, ok := e.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain selector: %d not found in environment", chainSelector)
		}
		_, ok = state[chainSelector]
		if !ok {
			return fmt.Errorf("chain selector: %d not found for MCMS state", chainSelector)
		}
		if err := c.Proposer.Validate(); err != nil {
			return err
		}
		if err := c.Canceller.Validate(); err != nil {
			return err
		}
		if err := c.Bypasser.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// setConfigOrTxData executes set config tx or gets the tx data for the MCMS proposal
func setConfigOrTxData(ctx context.Context, lggr logger.Logger, chain deployment.Chain, cfg config.Config, contract *gethwrappers.ManyChainMultiSig, useMCMS bool) (*types.Transaction, error) {
	groupQuorums, groupParents, signerAddresses, signerGroups := cfg.ExtractSetConfigInputs()
	opts := deployment.SimTransactOpts()
	if !useMCMS {
		opts = chain.DeployerKey
	}
	opts.Context = ctx
	tx, err := contract.SetConfig(opts, signerAddresses, signerGroups, groupQuorums, groupParents, false)
	if err != nil {
		return nil, err
	}
	if !useMCMS {
		_, err = deployment.ConfirmIfNoErrorWithABI(chain, tx, gethwrappers.ManyChainMultiSigABI, err)
		if err != nil {
			return nil, err
		}
		lggr.Infow("SetConfigMCMS tx confirmed", "txHash", tx.Hash().Hex())
	}
	return tx, nil
}

type setConfigTxs struct {
	proposerTx  *types.Transaction
	cancellerTx *types.Transaction
	bypasserTx  *types.Transaction
}

// setConfigPerRole sets the configuration for each of the MCMS contract roles on the mcmsState.
func setConfigPerRole(ctx context.Context, lggr logger.Logger, chain deployment.Chain, cfg ConfigPerRole, mcmsState *MCMSWithTimelockState, useMCMS bool) (setConfigTxs, error) {
	// Proposer set config
	proposerTx, err := setConfigOrTxData(ctx, lggr, chain, cfg.Proposer, mcmsState.ProposerMcm, useMCMS)
	if err != nil {
		return setConfigTxs{}, err
	}
	// Canceller set config
	cancellerTx, err := setConfigOrTxData(ctx, lggr, chain, cfg.Canceller, mcmsState.CancellerMcm, useMCMS)
	if err != nil {
		return setConfigTxs{}, err
	}
	// Bypasser set config
	bypasserTx, err := setConfigOrTxData(ctx, lggr, chain, cfg.Bypasser, mcmsState.BypasserMcm, useMCMS)
	if err != nil {
		return setConfigTxs{}, err
	}

	return setConfigTxs{
		proposerTx:  proposerTx,
		cancellerTx: cancellerTx,
		bypasserTx:  bypasserTx,
	}, nil
}

func addTxsToProposalBatch(setConfigTxsChain setConfigTxs, chainSelector uint64, state MCMSWithTimelockState) timelock.BatchChainOperation {
	result := timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(chainSelector),
		Batch:           []mcms.Operation{},
	}
	result.Batch = append(result.Batch, mcms.Operation{
		To:           state.ProposerMcm.Address(),
		Data:         setConfigTxsChain.proposerTx.Data(),
		Value:        big.NewInt(0),
		ContractType: string(commontypes.ProposerManyChainMultisig),
	})
	result.Batch = append(result.Batch, mcms.Operation{
		To:           state.CancellerMcm.Address(),
		Data:         setConfigTxsChain.cancellerTx.Data(),
		Value:        big.NewInt(0),
		ContractType: string(commontypes.CancellerManyChainMultisig),
	})
	result.Batch = append(result.Batch, mcms.Operation{
		To:           state.BypasserMcm.Address(),
		Data:         setConfigTxsChain.bypasserTx.Data(),
		Value:        big.NewInt(0),
		ContractType: string(commontypes.BypasserManyChainMultisig),
	})
	return result
}

// SetConfigMCMS sets the configuration of the MCMS contract on the chain identified by the chainSelector.
// Deprecated: Use SetConfigMCMSV2 instead.
func SetConfigMCMS(e deployment.Environment, cfg MCMSConfig) (deployment.ChangesetOutput, error) {
	selectors := []uint64{}
	lggr := e.Logger
	ctx := e.GetContext()
	for chainSelector := range cfg.ConfigsPerChain {
		selectors = append(selectors, chainSelector)
	}
	useMCMS := cfg.ProposalConfig != nil
	err := cfg.Validate(e, selectors)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	batches := []timelock.BatchChainOperation{}
	timelocksPerChain := map[uint64]common.Address{}
	proposerMcmsPerChain := map[uint64]*gethwrappers.ManyChainMultiSig{}

	mcmsStatePerChain, err := MaybeLoadMCMSWithTimelockState(e, selectors)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSelector, c := range cfg.ConfigsPerChain {
		chain := e.Chains[chainSelector]
		state := mcmsStatePerChain[chainSelector]
		timelocksPerChain[chainSelector] = state.Timelock.Address()
		proposerMcmsPerChain[chainSelector] = state.ProposerMcm
		setConfigTxsChain, err := setConfigPerRole(ctx, lggr, chain, c, state, useMCMS)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if useMCMS {
			batch := addTxsToProposalBatch(setConfigTxsChain, chainSelector, *state)
			batches = append(batches, batch)
		}
	}

	if useMCMS {
		// Create MCMS with timelock proposal
		proposal, err := proposalutils.BuildProposalFromBatches(timelocksPerChain, proposerMcmsPerChain, batches, "Set config proposal", cfg.ProposalConfig.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
		}
		lggr.Infow("SetConfigMCMS proposal created", "proposal", proposal)
		return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{*proposal}}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

// SetConfigMCMSV2 is a reimplementation of SetConfigMCMS that uses the new MCMS library.
func SetConfigMCMSV2(e deployment.Environment, cfg MCMSConfig) (deployment.ChangesetOutput, error) {
	selectors := []uint64{}
	lggr := e.Logger
	ctx := e.GetContext()
	for chainSelector := range cfg.ConfigsPerChain {
		selectors = append(selectors, chainSelector)
	}
	useMCMS := cfg.ProposalConfig != nil
	err := cfg.Validate(e, selectors)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var batches []mcmstypes.BatchOperation
	timelockAddressesPerChain := map[uint64]string{}
	inspectorPerChain := map[uint64]sdk.Inspector{}
	proposerMcmsPerChain := map[uint64]string{}

	mcmsStatePerChain, err := MaybeLoadMCMSWithTimelockState(e, selectors)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSelector, c := range cfg.ConfigsPerChain {
		chain := e.Chains[chainSelector]
		state := mcmsStatePerChain[chainSelector]
		timelockAddressesPerChain[chainSelector] = state.Timelock.Address().Hex()
		proposerMcmsPerChain[chainSelector] = state.ProposerMcm.Address().Hex()
		inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
		setConfigTxsChain, err := setConfigPerRole(ctx, lggr, chain, c, state, useMCMS)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if useMCMS {
			batch := addTxsToProposalBatchV2(setConfigTxsChain, chainSelector, *state)
			batches = append(batches, batch)
		}
	}

	if useMCMS {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(e.GetContext(), timelockAddressesPerChain,
			proposerMcmsPerChain, inspectorPerChain, batches, "Set config proposal", cfg.ProposalConfig.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
		}
		lggr.Infow("SetConfigMCMS proposal created", "proposal", proposal)
		return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

func addTxsToProposalBatchV2(setConfigTxsChain setConfigTxs, chainSelector uint64, state MCMSWithTimelockState) mcmstypes.BatchOperation {
	result := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chainSelector),
		Transactions:  []mcmstypes.Transaction{},
	}

	result.Transactions = append(result.Transactions,
		evm.NewTransaction(state.ProposerMcm.Address(),
			setConfigTxsChain.proposerTx.Data(), big.NewInt(0), string(commontypes.ProposerManyChainMultisig), nil))

	result.Transactions = append(result.Transactions, evm.NewTransaction(state.CancellerMcm.Address(),
		setConfigTxsChain.cancellerTx.Data(), big.NewInt(0), string(commontypes.CancellerManyChainMultisig), nil))

	result.Transactions = append(result.Transactions,
		evm.NewTransaction(state.BypasserMcm.Address(),
			setConfigTxsChain.bypasserTx.Data(), big.NewInt(0), string(commontypes.BypasserManyChainMultisig), nil))
	return result
}
