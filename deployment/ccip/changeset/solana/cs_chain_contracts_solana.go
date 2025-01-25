package changeset_solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type AddRemoteChainToSolanaConfig struct {
	// UpdatesByChain is a mapping of source -> dest -> update
	UpdatesByChain map[uint64]map[uint64]RemoteChainConfigSolana
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *cs.MCMSConfig
}

type RemoteChainConfigSolana struct {
	EnabledAsSource      bool
	EnabledAsDestination bool
	// TODO: what if remote chain family is solana ? will this be the router address ?
	RemoteChainOnRampAddress string
	DefaultTxGasLimit        uint32
	MaxPerMsgGasLimit        uint32
	MaxDataBytes             uint32
	MaxNumberOfTokensPerMsg  uint16
	ChainFamilySelector      [4]uint8
}

func (cfg AddRemoteChainToSolanaConfig) Validate(e deployment.Environment) error {
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return err
	}

	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.SolChains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}

		if chainState.Router.IsZero() {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}

		if err := commoncs.ValidateOwnershipSolana(e.GetContext(), cfg.MCMS != nil, e.SolChains[chainSel].DeployerKey.PublicKey(), chainState.Timelock, chainState.Router); err != nil {
			return err
		}

		var routerConfigAccount solRouter.Config
		err = solCommonUtil.GetAccountDataBorshInto(e.GetContext(), e.SolChains[chainSel].Client, cs.GetRouterConfigPDA(chainState.Router), deployment.SolDefaultCommitment, &routerConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to get router config %s: %w", chainState.Router, err)
		}

		for destination := range updates {
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not supported", destination)
			}
			if destination == routerConfigAccount.SolanaChainSelector {
				return fmt.Errorf("cannot add remote chain with same chain selector as current chain %d", destination)
			}
		}
	}

	return nil
}

// AddRemoteChainToSolana adds new remote chain configurations to Solana CCIP routers
func AddRemoteChainToSolana(e deployment.Environment, cfg AddRemoteChainToSolanaConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	s, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSel, updates := range cfg.UpdatesByChain {
		_, err := doAddRemoteChainToSolana(e, s, chainSel, updates)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	return deployment.ChangesetOutput{}, nil
}

func doAddRemoteChainToSolana(e deployment.Environment, s cs.CCIPOnChainState, chainSel uint64, updates map[uint64]RemoteChainConfigSolana) (deployment.ChangesetOutput, error) {
	e.Logger.Infow("Adding remote chain to solana", "chain", chainSel, "updates", updates)
	chain := e.SolChains[chainSel]

	ccipRouterID := s.SolChains[chainSel].Router

	// TODO: will this fail if chain has already been added?
	for destination, update := range updates {
		// TODO: this should be GetSourceChainStatePDA
		sourceChainStatePDA := cs.GetEvmSourceChainStatePDA(ccipRouterID, destination)
		validSourceChainConfig := solRouter.SourceChainConfig{
			OnRamp:    []byte(update.RemoteChainOnRampAddress),
			IsEnabled: update.EnabledAsSource,
		}
		// TODO: this should be GetDestChainStatePDA
		destChainStatePDA := cs.GetEvmDestChainStatePDA(ccipRouterID, destination)
		validDestChainConfig := solRouter.DestChainConfig{
			IsEnabled:               update.EnabledAsDestination,
			DefaultTxGasLimit:       update.DefaultTxGasLimit,
			MaxPerMsgGasLimit:       update.MaxPerMsgGasLimit,
			MaxDataBytes:            update.MaxDataBytes,
			MaxNumberOfTokensPerMsg: update.MaxNumberOfTokensPerMsg,
			// TODO: what if chain family is solana ?
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			ChainFamilySelector: [4]uint8{40, 18, 213, 44},
		}
		instruction, err := solRouter.NewAddChainSelectorInstruction(
			destination,
			validSourceChainConfig,
			validDestChainConfig,
			sourceChainStatePDA,
			destChainStatePDA,
			cs.GetRouterConfigPDA(ccipRouterID),
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}

		err = chain.Confirm([]solana.Instruction{instruction})

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
		e.Logger.Infow("Confirmed instruction", "instruction", instruction)
	}

	return deployment.ChangesetOutput{}, nil
}
