package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	solanasdk "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type ConfigPerRoleV2 struct {
	Proposer  mcmstypes.Config
	Canceller mcmstypes.Config
	Bypasser  mcmstypes.Config
}

type MCMSConfigV2 struct {
	ConfigsPerChain map[uint64]ConfigPerRoleV2
	ProposalConfig  *proposalutils.TimelockConfig
}

var _ cldf.ChangeSet[MCMSConfigV2] = SetConfigMCMSV2

// Validate checks that the MCMSConfigV2 is valid
func (cfg MCMSConfigV2) Validate(e cldf.Environment, selectors []uint64) error {
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
				err := cfg.ProposalConfig.Validate(e.BlockChains.EVMChains()[chainSelector], *chainState)
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

// setConfigOrTxDataV2 executes set config tx or gets the tx data for the MCMS proposal
func setConfigOrTxDataV2(ctx context.Context, lggr logger.Logger, chain cldf_evm.Chain, cfg mcmstypes.Config, contract *gethwrappers.ManyChainMultiSig, useMCMS bool) (*types.Transaction, error) {
	opts := cldf.SimTransactOpts()
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
		_, err = cldf.ConfirmIfNoErrorWithABI(chain, transaction, gethwrappers.ManyChainMultiSigABI, err)
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

// setConfigPerRoleV2 sets the configuration for each of the MCMS contract roles on the mcmsState.
func setConfigPerRoleV2(ctx context.Context, lggr logger.Logger, chain cldf_evm.Chain, cfg ConfigPerRoleV2, mcmsState *commonState.MCMSWithTimelockState, useMCMS bool) (setConfigTxs, error) {
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

// SetConfigMCMSV2 is a reimplementation of SetConfigMCMS that uses the new MCMS library.
func SetConfigMCMSV2(e cldf.Environment, cfg MCMSConfigV2) (cldf.ChangesetOutput, error) {
	selectors := []uint64{}
	lggr := e.Logger
	ctx := e.GetContext()
	for chainSelector := range cfg.ConfigsPerChain {
		selectors = append(selectors, chainSelector)
	}
	useMCMS := cfg.ProposalConfig != nil
	err := cfg.Validate(e, selectors)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	var batches []mcmstypes.BatchOperation
	timelockAddressesPerChain := map[uint64]string{}
	proposerMcmsPerChain := map[uint64]string{}
	inspectorPerChain, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for chainSelector, c := range cfg.ConfigsPerChain {
		family, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		switch family {
		case chain_selectors.FamilyEVM:
			chain := e.BlockChains.EVMChains()[chainSelector]
			mcmsStatePerChain, err := commonState.MaybeLoadMCMSWithTimelockState(e, []uint64{chainSelector})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			state := mcmsStatePerChain[chainSelector]
			timelockAddressesPerChain[chainSelector] = state.Timelock.Address().Hex()
			if cfg.ProposalConfig != nil {
				mcmsContract, err := cfg.ProposalConfig.MCMBasedOnAction(*state)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contract: %w", err)
				}
				proposerMcmsPerChain[chainSelector] = mcmsContract.Address().Hex()
			}
			setConfigTxsChain, err := setConfigPerRoleV2(ctx, lggr, chain, c, state, useMCMS)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			if useMCMS {
				batch := addTxsToProposalBatchV2(setConfigTxsChain, chainSelector, *state)
				batches = append(batches, batch)
			}
		case chain_selectors.FamilySolana:
			batch, err := setConfigSolana(e, chainSelector, c, timelockAddressesPerChain, proposerMcmsPerChain, useMCMS)
			if err != nil {
				return cldf.ChangesetOutput{}, err
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
		}
		lggr.Infow("SetConfigMCMS proposal created", "proposal", proposal)
		return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}

	return cldf.ChangesetOutput{}, nil
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
	e cldf.Environment, chainSelector uint64, cfg ConfigPerRoleV2,
	timelockAddressesPerChain, proposerMcmsPerChain map[uint64]string, useMCMS bool,
) ([]mcmstypes.BatchOperation, error) {
	chain := e.BlockChains.SolanaChains()[chainSelector]
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

func setConfigForRole(e cldf.Environment, chain cldf_solana.Chain, cfg mcmstypes.Config, mcmAddress string, contractType string, useMCMS bool, timelockSignerPDA solanasdk.PublicKey) (mcmstypes.BatchOperation, error) {
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
