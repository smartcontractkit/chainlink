package changeset_solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

var _ deployment.ChangeSet[AddRemoteChainToSolanaConfig] = AddRemoteChainToSolana
var _ deployment.ChangeSet[AddTokenPoolConfig] = AddTokenPool

type AddRemoteChainToSolanaConfig struct {
	// UpdatesByChain is a mapping of source -> dest -> update
	UpdatesByChain map[uint64]map[uint64]RemoteChainConfigSolana
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *cs.MCMSConfig
}

// We are not using solRouter.SourceChainConfig because that would involve the user
// converting the onRamp address into [2][64]byte{} which is not intuitive.
// The solRouter.DestChainConfig on the other hand has a lot of fields and most of them are uint
// So we are using that directly instead of copying over the fields here to reduce
// overhead cost if that type is bumped in chainlink-ccip
type RemoteChainConfigSolana struct {
	// source
	EnabledAsSource bool
	// TODO: what if remote chain family is solana ? will this be the router address ?
	RemoteChainOnRampAddress string

	// destination
	DestinationConfig solRouter.DestChainConfig
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
		configPDA, _, _ := solState.FindConfigPDA(chainState.Router)
		err = solCommonUtil.GetAccountDataBorshInto(e.GetContext(), e.SolChains[chainSel].Client, configPDA, deployment.SolDefaultCommitment, &routerConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to get router config %s: %w", chainState.Router, err)
		}

		for destination := range updates {
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not supported", destination)
			}
			if destination == routerConfigAccount.SvmChainSelector {
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
		sourceChainStatePDA, _ := solState.FindSourceChainStatePDA(destination, ccipRouterID)

		// Convert string address to bytes and pad to 64 bytes
		var onRampBytes [64]byte
		addressBytes := []byte(update.RemoteChainOnRampAddress)
		copy(onRampBytes[:], addressBytes)

		validSourceChainConfig := solRouter.SourceChainConfig{
			OnRamp:    [2][64]byte{onRampBytes, [64]byte{}},
			IsEnabled: update.EnabledAsSource,
		}
		configPDA, _, _ := solState.FindConfigPDA(ccipRouterID)
		destChainStatePDA, _ := solState.FindDestChainStatePDA(destination, ccipRouterID)
		instruction, err := solRouter.NewAddChainSelectorInstruction(
			destination,
			validSourceChainConfig,
			update.DestinationConfig,
			sourceChainStatePDA,
			destChainStatePDA,
			configPDA,
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

// ADD TOKEN POOL
type AddTokenPoolConfig struct {
	ChainSelector    uint64
	PoolType         string
	RampAuthority    string
	Authority        string
	TokenPubKey      string
	TokenProgramName string
}

func (cfg AddTokenPoolConfig) Validate(e deployment.Environment) error {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}
	_, err = deployment.GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return err
	}
	_, err = deployment.GetPoolType(cfg.PoolType)
	if err != nil {
		return err
	}
	return nil
}

func AddTokenPool(e deployment.Environment, cfg AddTokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	tokenprogramID, _ := deployment.GetTokenProgramID(cfg.TokenProgramName)
	poolType, _ := deployment.GetPoolType(cfg.PoolType)

	// Convert string addresses to public keys
	rampAuthorityPubKey := solana.MustPublicKeyFromBase58(cfg.RampAuthority)
	authorityPubKey := solana.MustPublicKeyFromBase58(cfg.Authority)
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	poolConfig, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenprogramID)
	poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, tokenprogramID)

	// ata for token pool
	createI, tokenPoolATA, err := solTokenUtil.CreateAssociatedTokenAccount(
		tokenprogramID,
		tokenPubKey,
		poolSigner,
		chain.DeployerKey.PublicKey(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	token_pool.SetProgramID(chainState.TokenPool)
	// initialize token pool for token
	poolInitI, err := token_pool.NewInitializeInstruction(
		poolType,
		rampAuthorityPubKey,
		poolConfig,
		tokenPubKey,
		poolSigner,
		authorityPubKey,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// make pool mint_authority for token (required for burn/mint)
	authI, err := solTokenUtil.SetTokenMintAuthority(
		tokenprogramID,
		poolSigner,
		tokenPubKey,
		chain.DeployerKey.PublicKey(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions := []solana.Instruction{createI, poolInitI, authI}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Created new token pool config", "token_pool_ata", tokenPoolATA.String(), "pool_config", poolConfig.String(), "pool_signer", poolSigner.String())
	e.Logger.Infow("Set mint authority", "poolSigner", poolSigner.String())

	return deployment.ChangesetOutput{}, nil
}
