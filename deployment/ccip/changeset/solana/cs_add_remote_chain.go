package solana

import (
	"context"

	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// ADD REMOTE CHAIN
type AddRemoteChainToSolanaConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]RemoteChainConfigSolana
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *cs.MCMSConfig
	// Public key of program authorities. Depending on when this changeset is called, some may be under
	// the control of the deployer, and some may be under the control of the timelock. (e.g. during new offramp deploy)
	RouterAuthority    solana.PublicKey
	FeeQuoterAuthority solana.PublicKey
	OffRampAuthority   solana.PublicKey
}

type RemoteChainConfigSolana struct {
	// source
	EnabledAsSource bool
	// destination
	RouterDestinationConfig    solRouter.DestChainConfig
	FeeQuoterDestinationConfig solFeeQuoter.DestChainConfig
	// We have different instructions for add vs update, so we need to know which one to use
	IsUpdate bool
}

func (cfg AddRemoteChainToSolanaConfig) Validate(e deployment.Environment) error {
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	if err := validateOffRampConfig(chain, chainState); err != nil {
		return err
	}
	if err := ValidateMCMSConfig(e, cfg.ChainSelector, cfg.MCMS); err != nil {
		return err
	}
	routerUsingMCMS := cfg.MCMS != nil && !cfg.RouterAuthority.IsZero()
	feeQuoterUsingMCMS := cfg.MCMS != nil && !cfg.FeeQuoterAuthority.IsZero()
	offRampUsingMCMS := cfg.MCMS != nil && !cfg.OffRampAuthority.IsZero()
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cs.ValidateOwnershipSolana(&e, chain, routerUsingMCMS, e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey(), chainState.Router, cs.Router); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if err := cs.ValidateOwnershipSolana(&e, chain, feeQuoterUsingMCMS, e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey(), chainState.FeeQuoter, cs.FeeQuoter); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	if err := cs.ValidateOwnershipSolana(&e, chain, offRampUsingMCMS, e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey(), chainState.OffRamp, cs.OffRamp); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	var routerConfigAccount solRouter.Config
	// already validated that router config exists
	_ = chain.GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)

	supportedChains := state.SupportedChains()
	for remote := range cfg.UpdatesByChain {
		if _, ok := supportedChains[remote]; !ok {
			return fmt.Errorf("remote chain %d is not supported", remote)
		}
		if remote == routerConfigAccount.SvmChainSelector {
			return fmt.Errorf("cannot add remote chain %d with same chain selector as current chain %d", remote, cfg.ChainSelector)
		}
		if err := state.ValidateRamp(remote, cs.OnRamp); err != nil {
			return err
		}
		routerDestChainPDA, err := solState.FindDestChainStatePDA(remote, chainState.Router)
		if err != nil {
			return fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		if !cfg.UpdatesByChain[remote].IsUpdate {
			var destChainStateAccount solRouter.DestChain
			err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
			if err == nil {
				return fmt.Errorf("remote %d is already configured on solana chain %d", remote, cfg.ChainSelector)
			}
		}
	}
	return nil
}

// Adds new remote chain configurations
func AddRemoteChainToSolana(e deployment.Environment, cfg AddRemoteChainToSolanaConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ab := deployment.NewMemoryAddressBook()
	txns, err := doAddRemoteChainToSolana(e, s, cfg, ab)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}

	// create proposals for ixns
	if len(txns) > 0 {
		timelocks := map[uint64]string{}
		proposers := map[uint64]string{}
		inspectors := map[uint64]sdk.Inspector{}
		batches := make([]mcmsTypes.BatchOperation, 0)
		chain := e.SolChains[cfg.ChainSelector]
		addresses, _ := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
		mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)

		timelocks[cfg.ChainSelector] = mcmsSolana.ContractAddress(
			mcmState.TimelockProgram,
			mcmsSolana.PDASeed(mcmState.TimelockSeed),
		)
		proposers[cfg.ChainSelector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
		inspectors[cfg.ChainSelector] = mcmsSolana.NewInspector(chain.Client)
		batches = append(batches, mcmsTypes.BatchOperation{
			ChainSelector: mcmsTypes.ChainSelector(cfg.ChainSelector),
			Transactions:  txns,
		})
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e.GetContext(),
			timelocks,
			proposers,
			inspectors,
			batches,
			"proposal to add remote chains to Solana",
			cfg.MCMS.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           ab,
		}, nil
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToSolana(
	e deployment.Environment,
	s cs.CCIPOnChainState,
	cfg AddRemoteChainToSolanaConfig,
	ab deployment.AddressBook) ([]mcmsTypes.Transaction, error) {
	txns := make([]mcmsTypes.Transaction, 0)
	ixns := make([]solana.Instruction, 0)
	chainSel := cfg.ChainSelector
	updates := cfg.UpdatesByChain
	chain := e.SolChains[chainSel]
	ccipRouterID := s.SolChains[chainSel].Router
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	offRampID := s.SolChains[chainSel].OffRamp
	routerUsingMCMS := cfg.MCMS != nil && !cfg.RouterAuthority.IsZero()
	feeQuoterUsingMCMS := cfg.MCMS != nil && !cfg.FeeQuoterAuthority.IsZero()
	offRampUsingMCMS := cfg.MCMS != nil && !cfg.OffRampAuthority.IsZero()
	lookUpTableEntries := make([]solana.PublicKey, 0)

	for remoteChainSel, update := range updates {
		var onRampBytes [64]byte
		// already verified, skipping errcheck
		remoteChainFamily, _ := chainsel.GetSelectorFamily(remoteChainSel)
		var addressBytes []byte
		switch remoteChainFamily {
		case chainsel.FamilySolana:
			addressBytes, _ = s.SolChains[remoteChainSel].OnRampBytes()
		case chainsel.FamilyEVM:
			addressBytes, _ = s.Chains[remoteChainSel].OnRampBytes()
		}
		addressBytes = common.LeftPadBytes(addressBytes, 64)
		copy(onRampBytes[:], addressBytes)

		// verified while loading state
		fqDestChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		routerDestChainPDA, _ := solState.FindDestChainStatePDA(remoteChainSel, ccipRouterID)
		offRampSourceChainPDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, s.SolChains[chainSel].OffRamp)

		if !update.IsUpdate {
			lookUpTableEntries = append(lookUpTableEntries,
				fqDestChainPDA,
				routerDestChainPDA,
				offRampSourceChainPDA,
			)
		}

		solRouter.SetProgramID(ccipRouterID)
		var authority solana.PublicKey
		if routerUsingMCMS {
			authority = cfg.RouterAuthority
		} else {
			authority = chain.DeployerKey.PublicKey()
		}
		var routerIx solana.Instruction
		var err error
		if update.IsUpdate {
			routerIx, err = solRouter.NewUpdateDestChainConfigInstruction(
				remoteChainSel,
				update.RouterDestinationConfig,
				routerDestChainPDA,
				s.SolChains[chainSel].RouterConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
		} else {
			routerIx, err = solRouter.NewAddChainSelectorInstruction(
				remoteChainSel,
				update.RouterDestinationConfig,
				routerDestChainPDA,
				s.SolChains[chainSel].RouterConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
		}
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if routerUsingMCMS {
			tx, err := BuildMCMSTxn(routerIx, ccipRouterID.String(), cs.Router)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, routerIx)
		}

		solFeeQuoter.SetProgramID(feeQuoterID)
		if feeQuoterUsingMCMS {
			authority = cfg.RouterAuthority
		} else {
			authority = chain.DeployerKey.PublicKey()
		}
		var feeQuoterIx solana.Instruction
		if update.IsUpdate {
			feeQuoterIx, err = solFeeQuoter.NewUpdateDestChainConfigInstruction(
				remoteChainSel,
				update.FeeQuoterDestinationConfig,
				s.SolChains[chainSel].FeeQuoterConfigPDA,
				fqDestChainPDA,
				authority,
			).ValidateAndBuild()
		} else {
			feeQuoterIx, err = solFeeQuoter.NewAddDestChainInstruction(
				remoteChainSel,
				update.FeeQuoterDestinationConfig,
				s.SolChains[chainSel].FeeQuoterConfigPDA,
				fqDestChainPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
		}
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if feeQuoterUsingMCMS {
			tx, err := BuildMCMSTxn(feeQuoterIx, feeQuoterID.String(), cs.FeeQuoter)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, feeQuoterIx)
		}

		solOffRamp.SetProgramID(offRampID)
		validSourceChainConfig := solOffRamp.SourceChainConfig{
			OnRamp:    [2][64]byte{onRampBytes, [64]byte{}},
			IsEnabled: update.EnabledAsSource,
		}
		if offRampUsingMCMS {
			authority = cfg.RouterAuthority
		} else {
			authority = chain.DeployerKey.PublicKey()
		}
		var offRampIx solana.Instruction
		if update.IsUpdate {
			offRampIx, err = solOffRamp.NewUpdateSourceChainConfigInstruction(
				remoteChainSel,
				validSourceChainConfig,
				offRampSourceChainPDA,
				s.SolChains[chainSel].OffRampConfigPDA,
				authority,
			).ValidateAndBuild()
		} else {
			offRampIx, err = solOffRamp.NewAddSourceChainInstruction(
				remoteChainSel,
				validSourceChainConfig,
				offRampSourceChainPDA,
				s.SolChains[chainSel].OffRampConfigPDA,
				authority,
				solana.SystemProgramID,
			).ValidateAndBuild()
		}
		if err != nil {
			return txns, fmt.Errorf("failed to generate instructions: %w", err)
		}
		if offRampUsingMCMS {
			tx, err := BuildMCMSTxn(offRampIx, offRampID.String(), cs.OffRamp)
			if err != nil {
				return txns, fmt.Errorf("failed to create transaction: %w", err)
			}
			txns = append(txns, *tx)
		} else {
			ixns = append(ixns, offRampIx)
		}
		if len(ixns) > 0 {
			err = chain.Confirm(ixns)
			if err != nil {
				return txns, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
		if !update.IsUpdate {
			tv := deployment.NewTypeAndVersion(cs.RemoteDest, deployment.Version1_0_0)
			remoteChainSelStr := strconv.FormatUint(remoteChainSel, 10)
			tv.AddLabel(remoteChainSelStr)
			err = ab.Save(chainSel, routerDestChainPDA.String(), tv)
			if err != nil {
				return txns, fmt.Errorf("failed to save dest chain state to address book: %w", err)
			}

			tv = deployment.NewTypeAndVersion(cs.RemoteSource, deployment.Version1_0_0)
			tv.AddLabel(remoteChainSelStr)
			err = ab.Save(chainSel, offRampSourceChainPDA.String(), tv)
			if err != nil {
				return txns, fmt.Errorf("failed to save source chain state to address book: %w", err)
			}
		}
	}

	if len(lookUpTableEntries) > 0 {
		addressLookupTable, err := cs.FetchOfframpLookupTable(e.GetContext(), chain, offRampID)
		if err != nil {
			return txns, fmt.Errorf("failed to get offramp reference addresses: %w", err)
		}

		if err := solCommonUtil.ExtendLookupTable(
			e.GetContext(),
			chain.Client,
			addressLookupTable,
			*chain.DeployerKey,
			lookUpTableEntries,
		); err != nil {
			return txns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	return txns, nil
}
