package solana

import (
	"context"
	// "errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

// ADD REMOTE CHAIN
type AddRemoteChainToSolanaConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]RemoteChainConfigSolana
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *cs.MCMSConfig
}

type RemoteChainConfigSolana struct {
	// source
	EnabledAsSource bool
	// destination
	RouterDestinationConfig    solRouter.DestChainConfig
	FeeQuoterDestinationConfig solFeeQuoter.DestChainConfig
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
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return err
	}
	mcmState, err := commonState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return fmt.Errorf("error loading MCMS state for chain %d: %w", cfg.ChainSelector, err)
	}
	if err := commoncs.ValidateOwnershipSolana(e.GetContext(), cfg.MCMS != nil, e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey(), mcmState.TimelockProgram, mcmState.TimelockSeed, chainState.Router); err != nil {
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
		var destChainStateAccount solRouter.DestChain
		err = chain.GetAccountDataBorshInto(context.Background(), routerDestChainPDA, &destChainStateAccount)
		if err == nil {
			return fmt.Errorf("remote %d is already configured on solana chain %d", remote, cfg.ChainSelector)
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
	err = doAddRemoteChainToSolana(e, s, cfg.ChainSelector, cfg.UpdatesByChain, ab)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: ab}, err
	}
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func doAddRemoteChainToSolana(
	e deployment.Environment,
	s cs.CCIPOnChainState,
	chainSel uint64,
	updates map[uint64]RemoteChainConfigSolana,
	ab deployment.AddressBook) error {
	chain := e.SolChains[chainSel]
	ccipRouterID := s.SolChains[chainSel].Router
	feeQuoterID := s.SolChains[chainSel].FeeQuoter
	offRampID := s.SolChains[chainSel].OffRamp
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

		lookUpTableEntries = append(lookUpTableEntries,
			fqDestChainPDA,
			routerDestChainPDA,
			offRampSourceChainPDA,
		)

		solRouter.SetProgramID(ccipRouterID)
		routerIx, err := solRouter.NewAddChainSelectorInstruction(
			remoteChainSel,
			update.RouterDestinationConfig,
			routerDestChainPDA,
			s.SolChains[chainSel].RouterConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return fmt.Errorf("failed to generate instructions: %w", err)
		}

		solFeeQuoter.SetProgramID(feeQuoterID)
		feeQuoterIx, err := solFeeQuoter.NewAddDestChainInstruction(
			remoteChainSel,
			update.FeeQuoterDestinationConfig,
			s.SolChains[chainSel].FeeQuoterConfigPDA,
			fqDestChainPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return fmt.Errorf("failed to generate instructions: %w", err)
		}

		solOffRamp.SetProgramID(offRampID)
		validSourceChainConfig := solOffRamp.SourceChainConfig{
			OnRamp:    [2][64]byte{onRampBytes, [64]byte{}},
			IsEnabled: update.EnabledAsSource,
		}
		offRampIx, err := solOffRamp.NewAddSourceChainInstruction(
			remoteChainSel,
			validSourceChainConfig,
			offRampSourceChainPDA,
			s.SolChains[chainSel].OffRampConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return fmt.Errorf("failed to generate instructions: %w", err)
		}

		err = chain.Confirm([]solana.Instruction{routerIx, feeQuoterIx, offRampIx})
		if err != nil {
			return fmt.Errorf("failed to confirm instructions: %w", err)
		}

		tv := deployment.NewTypeAndVersion(cs.RemoteDest, deployment.Version1_0_0)
		remoteChainSelStr := strconv.FormatUint(remoteChainSel, 10)
		tv.AddLabel(remoteChainSelStr)
		err = ab.Save(chainSel, routerDestChainPDA.String(), tv)
		if err != nil {
			return fmt.Errorf("failed to save dest chain state to address book: %w", err)
		}

		tv = deployment.NewTypeAndVersion(cs.RemoteSource, deployment.Version1_0_0)
		tv.AddLabel(remoteChainSelStr)
		err = ab.Save(chainSel, offRampSourceChainPDA.String(), tv)
		if err != nil {
			return fmt.Errorf("failed to save source chain state to address book: %w", err)
		}
	}

	if err := solCommonUtil.ExtendLookupTable(
		e.GetContext(),
		chain.Client,
		s.SolChains[chainSel].OfframpAddressLookupTable,
		*chain.DeployerKey,
		lookUpTableEntries,
	); err != nil {
		return fmt.Errorf("failed to extend lookup table: %w", err)
	}

	return nil
}
