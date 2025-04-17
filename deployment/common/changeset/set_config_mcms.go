package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	solanasdk "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// Deprecated: use ConfigPerRoleV2 instead
type ConfigPerRole struct {
	Proposer  config.Config
	Canceller config.Config
	Bypasser  config.Config
}

type ConfigPerRoleV2 struct {
	Proposer  mcmstypes.Config
	Canceller mcmstypes.Config
	Bypasser  mcmstypes.Config
}

type MCMSConfig struct {
	ConfigsPerChain map[uint64]ConfigPerRole
	ProposalConfig  *proposalutils.TimelockConfig
}

type MCMSConfigV2 struct {
	ConfigsPerChain map[uint64]ConfigPerRoleV2
	ProposalConfig  *proposalutils.TimelockConfig
}

var _ deployment.ChangeSet[MCMSConfig] = SetConfigMCMS
var _ deployment.ChangeSet[MCMSConfigV2] = SetConfigMCMSV2

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

// Validate checks that the MCMSConfigV2 is valid
func (cfg MCMSConfigV2) Validate(e deployment.Environment, selectors []uint64) error {
	if len(cfg.ConfigsPerChain) == 0 {
		return errors.New("no chain configs provided")
	}

	err := deployment.ValidateSelectorsInEnvironment(e, selectors)
	if err != nil {
		return err
	}

	for chainSelector, c := range cfg.ConfigsPerChain {
		family, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return err
		}

		switch family {
		case chain_selectors.FamilyEVM:
			state, err := commonState.MaybeLoadMCMSWithTimelockState(e, []uint64{chainSelector})
			if err != nil {
				return err
			}
			chainState, ok := state[chainSelector]
			if !ok {
				return fmt.Errorf("chain selector: %d not found for MCMS state", chainSelector)
			}
			if cfg.ProposalConfig != nil {
				err := cfg.ProposalConfig.Validate(e.Chains[chainSelector], *chainState)
				if err != nil {
					return err
				}
			}
		case chain_selectors.FamilySolana:
			state, err := commonState.MaybeLoadMCMSWithTimelockStateSolana(e, []uint64{chainSelector})
			if err != nil {
				return err
			}
			_, ok := state[chainSelector]
			if !ok {
				return fmt.Errorf("chain selector: %d not found for MCMS state", chainSelector)
			}
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
// Deprecated: Use setConfigOrTxDataV2 instead.
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

// setConfigOrTxDataV2 executes set config tx or gets the tx data for the MCMS proposal
func setConfigOrTxDataV2(ctx context.Context, lggr logger.Logger, chain deployment.Chain, cfg mcmstypes.Config, contract *gethwrappers.ManyChainMultiSig, useMCMS bool) (*types.Transaction, error) {
	opts := deployment.SimTransactOpts()
	if !useMCMS {
		opts = chain.DeployerKey
	}
	opts.Context = ctx

	configurer := evm.NewConfigurer(chain.Client, opts)
	res, err := configurer.SetConfig(ctx, contract.Address().Hex(), &cfg, false)
	if err != nil {
		return nil, err
	}

	transaction := res.RawData.(*types.Transaction)
	if !useMCMS {
		_, err = deployment.ConfirmIfNoErrorWithABI(chain, transaction, gethwrappers.ManyChainMultiSigABI, err)
		if err != nil {
			return nil, err
		}
		lggr.Infow("SetConfigMCMS tx confirmed", "txHash", res.Hash)
	}
	return transaction, nil
}

type setConfigTxs struct {
	proposerTx  *types.Transaction
	cancellerTx *types.Transaction
	bypasserTx  *types.Transaction
}

// setConfigPerRole sets the configuration for each of the MCMS contract roles on the mcmsState.
// Deprecated: Use setConfigPerRoleV2 instead.
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

// setConfigPerRoleV2 sets the configuration for each of the MCMS contract roles on the mcmsState.
func setConfigPerRoleV2(ctx context.Context, lggr logger.Logger, chain deployment.Chain, cfg ConfigPerRoleV2, mcmsState *commonState.MCMSWithTimelockState, useMCMS bool) (setConfigTxs, error) {
	// Proposer set config
	proposerTx, err := setConfigOrTxDataV2(ctx, lggr, chain, cfg.Proposer, mcmsState.ProposerMcm, useMCMS)
	if err != nil {
		return setConfigTxs{}, err
	}
	// Canceller set config
	cancellerTx, err := setConfigOrTxDataV2(ctx, lggr, chain, cfg.Canceller, mcmsState.CancellerMcm, useMCMS)
	if err != nil {
		return setConfigTxs{}, err
	}
	// Bypasser set config
	bypasserTx, err := setConfigOrTxDataV2(ctx, lggr, chain, cfg.Bypasser, mcmsState.BypasserMcm, useMCMS)
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
func SetConfigMCMSV2(e deployment.Environment, cfg MCMSConfigV2) (deployment.ChangesetOutput, error) {
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
	proposerMcmsPerChain := map[uint64]string{}
	inspectorPerChain, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSelector, c := range cfg.ConfigsPerChain {
		family, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		switch family {
		case chain_selectors.FamilyEVM:
			chain := e.Chains[chainSelector]
			mcmsStatePerChain, err := commonState.MaybeLoadMCMSWithTimelockState(e, []uint64{chainSelector})
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			state := mcmsStatePerChain[chainSelector]
			timelockAddressesPerChain[chainSelector] = state.Timelock.Address().Hex()
			if cfg.ProposalConfig != nil {
				mcmsContract, err := cfg.ProposalConfig.MCMBasedOnAction(*state)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contract: %w", err)
				}
				proposerMcmsPerChain[chainSelector] = mcmsContract.Address().Hex()
			}
			setConfigTxsChain, err := setConfigPerRoleV2(ctx, lggr, chain, c, state, useMCMS)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			if useMCMS {
				batch := addTxsToProposalBatchV2(setConfigTxsChain, chainSelector, *state)
				batches = append(batches, batch)
			}
		case chain_selectors.FamilySolana:
			batch, err := setConfigSolana(e, chainSelector, c, timelockAddressesPerChain, proposerMcmsPerChain, useMCMS)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			if useMCMS {
				batches = append(batches, batch...)
			}
		}
	}

	if useMCMS {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(e, timelockAddressesPerChain,
			proposerMcmsPerChain, inspectorPerChain, batches, "Set config proposal", *cfg.ProposalConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
		}
		lggr.Infow("SetConfigMCMS proposal created", "proposal", proposal)
		return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

func addTxsToProposalBatchV2(setConfigTxsChain setConfigTxs, chainSelector uint64, state commonState.MCMSWithTimelockState) mcmstypes.BatchOperation {
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

func setConfigSolana(
	e deployment.Environment, chainSelector uint64, cfg ConfigPerRoleV2,
	timelockAddressesPerChain, proposerMcmsPerChain map[uint64]string, useMCMS bool,
) ([]mcmstypes.BatchOperation, error) {
	chain := e.SolChains[chainSelector]
	mcmsStatePerChain, err := commonState.MaybeLoadMCMSWithTimelockStateSolana(e, []uint64{chainSelector})
	if err != nil {
		return nil, err
	}
	solState := mcmsStatePerChain[chainSelector]
	timelockAddressesPerChain[chainSelector] = solana.ContractAddress(solState.TimelockProgram, solana.PDASeed(solState.TimelockSeed))
	proposerMcmsPerChain[chainSelector] = solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.ProposerMcmSeed))
	cancellerAddress := solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.CancellerMcmSeed))
	bypasserAddress := solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.BypasserMcmSeed))
	proposerAddress := solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.ProposerMcmSeed))

	timelockSignerPDA, err := solana.FindTimelockSignerPDA(solState.TimelockProgram, solana.PDASeed(solState.TimelockSeed))
	if err != nil {
		return nil, err
	}

	batches := []mcmstypes.BatchOperation{}
	// broken into single batch per role (total 3 batches) due to size constraints on solana when all instructions were in the same single batch
	proposerOps, err := setConfigForRole(e, chain, cfg.Proposer, proposerAddress, string(commontypes.ProposerManyChainMultisig), useMCMS, timelockSignerPDA)
	if err != nil {
		return nil, err
	}
	batches = append(batches, proposerOps)

	cancellerOps, err := setConfigForRole(e, chain, cfg.Canceller, cancellerAddress, string(commontypes.CancellerManyChainMultisig), useMCMS, timelockSignerPDA)
	if err != nil {
		return nil, err
	}
	batches = append(batches, cancellerOps)
	bypasserOps, err := setConfigForRole(e, chain, cfg.Bypasser, bypasserAddress, string(commontypes.BypasserManyChainMultisig), useMCMS, timelockSignerPDA)
	if err != nil {
		return nil, err
	}
	batches = append(batches, bypasserOps)

	return batches, nil
}

func setConfigForRole(e deployment.Environment, chain deployment.SolChain, cfg mcmstypes.Config, mcmAddress string, contractType string, useMCMS bool, timelockSignerPDA solanasdk.PublicKey) (mcmstypes.BatchOperation, error) {
	var configurer *solana.Configurer

	if useMCMS {
		configurer = solana.NewConfigurer(chain.Client, *chain.DeployerKey, mcmstypes.ChainSelector(chain.Selector),
			solana.WithDoNotSendInstructionsOnChain(), solana.WithAuthorityAccount(timelockSignerPDA))
	} else {
		configurer = solana.NewConfigurer(chain.Client, *chain.DeployerKey, mcmstypes.ChainSelector(chain.Selector))
	}

	res, err := configurer.SetConfig(e.GetContext(), mcmAddress, &cfg, false)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}

	if useMCMS {
		instructions := res.RawData.([]solanasdk.Instruction)

		txs := make([]mcmstypes.Transaction, 0, len(instructions))
		for _, ix := range instructions {
			tx, err := solana.NewTransactionFromInstruction(ix, contractType, []string{})
			if err != nil {
				return mcmstypes.BatchOperation{}, err
			}
			txs = append(txs, tx)
		}

		return mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chain.Selector),
			Transactions:  txs,
		}, nil
	}

	e.Logger.Infow("SetConfig tx confirmed", "txHash", res.Hash)
	return mcmstypes.BatchOperation{}, nil
}
