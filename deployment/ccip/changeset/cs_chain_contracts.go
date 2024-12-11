package changeset

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

var _ deployment.ChangeSet[map[uint64]UpdateOnRampsDestsConfig] = UpdateOnRampsDests

type UpdateOnRampsDestsConfig struct {
	UpdatesByChain map[uint64][]OnRampDestinationUpdate
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be acheived by calling this function multiple times)
	MCMS *MCMSConfig
}

type OnRampDestinationUpdate struct {
	DestinationSelector uint64
	Disable             bool // If true, disables the destination by setting router to 0x0.
	TestRouter          bool // Flag for safety only allow specifying either router or testRouter.
	AllowListEnabled    bool
}

func (cfg UpdateOnRampsDestsConfig) Validate(ctx context.Context, state CCIPOnChainState) error {
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if chainState.TestRouter == nil {
			return fmt.Errorf("missing test router for chain %d", chainSel)
		}
		if chainState.Router == nil {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}
		if chainState.OnRamp == nil {
			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
		}
		for _, update := range updates {
			if update.DestinationSelector == 0 {
				return fmt.Errorf("destination selector cannot be 0")
			}
			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: ctx})
			if err != nil {
				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
			}
			if update.DestinationSelector == sc.ChainSelector {
				return fmt.Errorf("cannot update onramp destination to the same chain")
			}
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[update.DestinationSelector]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", update.DestinationSelector, chainState.OnRamp.Address())
			}
		}
	}
	return nil
}

// UpdateOnRampsDests updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDests(e deployment.Environment, cfg UpdateOnRampsDestsConfig) (deployment.ChangesetOutput, error) {
	s, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e.GetContext(), s); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var txes []*types.Transaction
	var batches []timelock.BatchChainOperation
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.Chains[chainSel].DeployerKey
		txOpts.Context = e.GetContext()
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		onRamp := s.Chains[chainSel].OnRamp
		var args []onramp.OnRampDestChainConfigArgs
		for _, update := range updates {
			router := common.HexToAddress("0x0")
			if !update.Disable {
				if update.TestRouter {
					router = s.Chains[chainSel].TestRouter.Address()
				} else {
					router = s.Chains[chainSel].Router.Address()
				}
			}
			args = append(args, onramp.OnRampDestChainConfigArgs{
				DestChainSelector: update.DestinationSelector,
				Router:            router,
				AllowlistEnabled:  update.AllowListEnabled,
			})
		}
		tx, err := onRamp.ApplyDestChainConfigUpdates(txOpts, args)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoError(e.Chains[chainSel], tx, err); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		} else {
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chainSel),
				Batch: []mcms.Operation{
					{
						To:    onRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	p, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.ChainSelector: s.Chains[cfg.ChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.ChainSelector: s.Chains[cfg.ChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{
			{
				ChainIdentifier: mcms.ChainIdentifier(cfg.ChainSelector),
				Batch: []mcms.Operation{
					{
						To:    onRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			},
		},
		"Update onramp destinations",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}
