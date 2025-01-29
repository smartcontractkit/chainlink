package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

var _ deployment.ChangeSet[AddRemoteChainToSolanaConfig] = AddRemoteChainToSolana
var _ deployment.ChangeSet[TokenPoolConfig] = AddTokenPool

// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName string) (solana.PublicKey, error) {
	tokenPrograms := map[string]solana.PublicKey{
		// "spl-token":      solana.TokenProgramID,
		deployment.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: spl-token, spl-token-2022", programName)
	}
	return programID, nil
}

// GetPoolType returns the token pool type constant for the given string
func GetPoolType(poolType string) (token_pool.PoolType, error) {
	poolTypes := map[string]token_pool.PoolType{
		"LockAndRelease": token_pool.LockAndRelease_PoolType,
		"BurnAndMint":    token_pool.BurnAndMint_PoolType,
	}

	poolTypeConstant, ok := poolTypes[poolType]
	if !ok {
		return 0, fmt.Errorf("invalid pool type: %s. Must be one of: LockAndRelease, BurnAndMint", poolType)
	}
	return poolTypeConstant, nil
}

func commonValidation(e deployment.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	chain, ok := e.SolChains[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return err
	}
	chainState, chainExists := state.SolChains[selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	exists := false
	for _, token := range chainState.SPL2022Tokens {
		if token.Equals(tokenPubKey) {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("token %s not found in existing state, deploy the prerequisites first", tokenPubKey.String())
	}
	return nil
}

func validateRouterConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
	if chainState.Router.IsZero() {
		return fmt.Errorf("ccip router not found in existing state, deploy the prerequisites first chain %d", chain.Selector)
	}
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	var routerConfigAccount ccip_router.Config
	err := chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, deploy the prerequisites first %d", chain.Selector)
	}
	return nil
}

// ADD REMOTE CHAIN
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
		chain := e.SolChains[chainSel]
		if err := validateRouterConfig(chain, chainState); err != nil {
			return err
		}
		if err := commoncs.ValidateOwnershipSolana(e.GetContext(), cfg.MCMS != nil, e.SolChains[chainSel].DeployerKey.PublicKey(), chainState.Timelock, chainState.Router); err != nil {
			return err
		}
		routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
		var routerConfigAccount ccip_router.Config
		// already validated that router config exists
		_ = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)

		for remote := range updates {
			if _, ok := supportedChains[remote]; !ok {
				return fmt.Errorf("remote chain %d is not supported", remote)
			}
			if remote == routerConfigAccount.SvmChainSelector {
				return fmt.Errorf("cannot add remote chain with same chain selector as current chain %d", remote)
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

// SET OCR3 CONFIG
func btoi(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// SetOCR3OffRamp will set the OCR3 offramp for the given chain.
// to the active configuration on CCIPHome. This
// is used to complete the candidate->active promotion cycle, it's
// run after the candidate is confirmed to be working correctly.
// Multichain is especially helpful for NOP rotations where we have
// to touch all the chain to change signers.
func SetOCR3ConfigSolana(e deployment.Environment, cfg cs.SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	solChains := state.SolChains

	// cfg.RemoteChainSels will be a bunch of solana chains
	// can add this in validate
	for _, remote := range cfg.RemoteChainSels {
		donID, err := internal.DonIDForChain(
			state.Chains[cfg.HomeChainSel].CapabilityRegistry,
			state.Chains[cfg.HomeChainSel].CCIPHome,
			remote)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		args, err := internal.BuildSetOCR3ConfigArgsSolana(donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		// TODO: check if ocr3 has already been set
		// set, err := isOCR3ConfigSetSolana(e.Logger, e.Chains[remote], state.Chains[remote].OffRamp, args)
		var instructions []solana.Instruction
		ccipRouterID := solChains[remote].Router
		configPDA, _, _ := solState.FindConfigPDA(ccipRouterID)
		routerStatePDA, _, _ := solState.FindStatePDA(ccipRouterID)
		for _, arg := range args {
			instruction, err := ccip_router.NewSetOcrConfigInstruction(
				arg.OCRPluginType,
				ccip_router.Ocr3ConfigInfo{
					ConfigDigest:                   arg.ConfigDigest,
					F:                              arg.F,
					IsSignatureVerificationEnabled: btoi(arg.IsSignatureVerificationEnabled),
				},
				arg.Signers,
				arg.Transmitters,
				configPDA,
				routerStatePDA,
				e.SolChains[remote].DeployerKey.PublicKey(),
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			instructions = append(instructions, instruction)
		}
		if cfg.MCMS == nil {
			err := e.SolChains[remote].Confirm(instructions)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	return deployment.ChangesetOutput{}, nil

	// TODO: timelock mcms support
}

// ADD TOKEN POOL
type TokenPoolConfig struct {
	ChainSelector    uint64
	PoolType         string
	Authority        string
	TokenPubKey      string
	TokenProgramName string
}

func (cfg TokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the token pool first for chain %d", cfg.ChainSelector)
	}
	_, err = GetPoolType(cfg.PoolType)
	if err != nil {
		return err
	}
	_, err = GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return err
	}
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, state.SolChains[cfg.ChainSelector].TokenPool)
	chain := e.SolChains[cfg.ChainSelector]
	var poolConfigAccount token_pool.Config
	err = chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount)
	if err == nil {
		return fmt.Errorf("token pool config already exists")
	}
	return nil
}

func AddTokenPool(e deployment.Environment, cfg TokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Adding token pool", "chain", cfg.ChainSelector, "pool_type", cfg.PoolType, "authority", cfg.Authority, "token_pubkey", cfg.TokenPubKey, "token_program_name", cfg.TokenProgramName)
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	tokenprogramID, _ := GetTokenProgramID(cfg.TokenProgramName)
	poolType, _ := GetPoolType(cfg.PoolType)
	rampAuthorityPubKey, _, _ := solState.FindExternalExecutionConfigPDA(chainState.Router)
	authorityPubKey := solana.MustPublicKeyFromBase58(cfg.Authority)
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.TokenPool)
	poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, chainState.TokenPool)

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
		poolConfigPDA,
		tokenPubKey,
		poolSigner,
		authorityPubKey, // this is assumed to be chain.DeployerKey for now
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

	// add signer here if authority is different from deployer key
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Created new token pool config", "token_pool_ata", tokenPoolATA.String(), "pool_config", poolConfigPDA.String(), "pool_signer", poolSigner.String())
	e.Logger.Infow("Set mint authority", "poolSigner", poolSigner.String())

	return deployment.ChangesetOutput{}, nil
}

// ADD TOKEN POOL FOR REMOTE CHAIN
type RemoteChainTokenPoolConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	TokenPubKey         string
	RemoteConfig        token_pool.RemoteConfig
	InboundRateLimit    token_pool.RateLimitConfig
	OutboundRateLimit   token_pool.RateLimitConfig
}

func (cfg RemoteChainTokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, solana.MustPublicKeyFromBase58(cfg.TokenPubKey))
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}

	// check if pool config exists (cannot do remote setup without it)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, state.SolChains[cfg.ChainSelector].TokenPool)
	chain := e.SolChains[cfg.ChainSelector]
	var poolConfigAccount token_pool.Config
	err = chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount)
	if err != nil {
		return fmt.Errorf("token pool config not found, call AddTokenPool first")
	}

	// check if existing pool setup already has this remote chain configured
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.TokenPool)
	var remoteChainConfigAccount token_pool.RemoteConfig
	err = chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)
	if err == nil {
		return fmt.Errorf("remote chain config already exists")
	}
	return nil
}

func SetupTokenPoolForRemoteChain(e deployment.Environment, cfg RemoteChainTokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.TokenPool)
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.TokenPool)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	token_pool.SetProgramID(chainState.TokenPool)
	ixConfigure, err := token_pool.NewInitChainRemoteConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	ixRates, err := token_pool.NewSetChainRateLimitInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.InboundRateLimit,
		cfg.OutboundRateLimit,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions := []solana.Instruction{ixConfigure, ixRates}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN
type BillingTokenPoolConfig struct {
	ChainSelector    uint64
	TokenPubKey      string
	TokenProgramName string
	Config           solRouter.BillingTokenConfig
}

func (cfg BillingTokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	// check if already setup
	billingConfigPDA, _, _ := solState.FindFeeBillingTokenConfigPDA(tokenPubKey, chainState.Router)
	var token0ConfigAccount ccip_router.BillingTokenConfigWrapper
	err = chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount)
	if err == nil {
		return fmt.Errorf("billing token config already exists")
	}

	_, err = GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return err
	}
	return nil
}

func AddBillingToken(e deployment.Environment, cfg BillingTokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	ccip_router.SetProgramID(chainState.Router)

	billingConfigPDA, _, _ := solState.FindFeeBillingTokenConfigPDA(tokenPubKey, chainState.Router)
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	tokenprogramID, _ := GetTokenProgramID(cfg.TokenProgramName)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenPubKey, billingSignerPDA)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	cfg.Config.Mint = tokenPubKey

	ixConfig, cerr := ccip_router.NewAddBillingTokenConfigInstruction(
		cfg.Config,
		routerConfigPDA,
		billingConfigPDA,
		tokenprogramID,
		tokenPubKey,
		token2022Receiver,
		chain.DeployerKey.PublicKey(),
		billingSignerPDA,
		ata.ProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if cerr != nil {
		return deployment.ChangesetOutput{}, cerr
	}

	instructions := []solana.Instruction{ixConfig}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Billing token added", "chainSelector", cfg.ChainSelector, "tokenPubKey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN FOR REMOTE CHAIN
type BillingTokenForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              ccip_router.TokenBilling
	TokenPubKey         string
}

func (cfg BillingTokenForRemoteChainConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.Router.IsZero() {
		return fmt.Errorf("router not found in existing state for chain %d", cfg.ChainSelector)
	}
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return fmt.Errorf("router validation failed: %w", err)
	}
	// check if desired state already exists ?
	remoteBillingPDA, _, _ := solState.FindCcipTokenpoolBillingPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.Router)
	var remoteBillingAccount ccip_router.TokenBilling
	err = chain.GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount)
	if err == nil {
		return fmt.Errorf("billing token config already exists")
	}
	return nil
}

func AddBillingTokenForRemoteChain(e deployment.Environment, cfg BillingTokenForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	remoteBillingPDA, _, _ := solState.FindCcipTokenpoolBillingPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.Router)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	ix, err := ccip_router.NewSetTokenBillingInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.Config,
		routerConfigPDA,
		remoteBillingPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions := []solana.Instruction{ix}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Token billing set for remote chain", "chainSelector ", cfg.ChainSelector, "remoteChainSelector ", cfg.RemoteChainSelector, "tokenPubKey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

// TOKEN ADMIN REGISTRY
type RegisterTokenAdminRegistryType int

const (
	ViaGetCcipAdminInstruction RegisterTokenAdminRegistryType = iota
	ViaOwnerInstruction
)

type RegisterTokenAdminRegistryConfig struct {
	ChainSelector  uint64
	TokenPubKey    string
	TokenPoolAdmin string
	RegisterType   RegisterTokenAdminRegistryType
}

func (cfg RegisterTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	err := commonValidation(e, cfg.ChainSelector, solana.MustPublicKeyFromBase58(cfg.TokenPubKey))
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.Router.IsZero() {
		return fmt.Errorf("ccip router not found in existing state, deploy the prerequisites first")
	}
	// check for router setup
	// check if tokenAdminRegistryPDA already exists
	return nil
}

func RegisterTokenAdminRegistry(e deployment.Environment, cfg RegisterTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	// Convert string addresses to public keys
	configPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(chainState.Router, tokenPubKey)
	tokenPoolAdminPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPoolAdmin)
	var instruction *ccip_router.Instruction
	var err error
	switch cfg.RegisterType {
	case ViaGetCcipAdminInstruction:
		instruction, err = solRouter.NewRegisterTokenAdminRegistryViaGetCcipAdminInstruction(
			tokenPubKey,
			tokenPoolAdminPubKey,
			configPDA,
			tokenAdminRegistryPDA,         // this gets created
			chain.DeployerKey.PublicKey(), // authority defaulting to deployer key
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	case ViaOwnerInstruction:
		instruction, err = solRouter.NewRegisterTokenAdminRegistryViaOwnerInstruction(
			configPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // authority defaulting to deployer key
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}
	// if we want to have a different authority, we will need to add the corresponding singer here
	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// TRANSFER AND ACCEPT TOKEN ADMIN REGISTRY
type TransferAndAcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector               uint64
	TokenName                   string
	NewTokenPoolAdminPrivateKey string
}

func (cfg TransferAndAcceptAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	err := commonValidation(e, cfg.ChainSelector, solana.MustPublicKeyFromBase58(cfg.TokenName))
	if err != nil {
		return err
	}
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.Router.IsZero() {
		return fmt.Errorf("ccip router not found in existing state, deploy the prerequisites first")
	}
	// check if router is setup
	return nil
}

func TransferAndAcceptAdminRoleTokenAdminRegistry(e deployment.Environment, cfg TransferAndAcceptAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenName)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(chainState.Router, tokenPubKey)
	configPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	newTokenPoolAdminPrivKey := solana.MustPrivateKeyFromBase58(cfg.NewTokenPoolAdminPrivateKey)
	newTokenPoolAdminPubKey := newTokenPoolAdminPrivKey.PublicKey()

	ix1, err := ccip_router.NewTransferAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		newTokenPoolAdminPubKey,
		configPDA,
		tokenAdminRegistryPDA,
		chain.DeployerKey.PublicKey(), // as we are assuming this is the default authority for everything in the beginning
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ix2, err := ccip_router.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		configPDA,
		tokenAdminRegistryPDA,
		newTokenPoolAdminPubKey,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{ix1, ix2}
	// the new authority will have to sign the acceptance
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(newTokenPoolAdminPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// TODO:
// token admin registry tests
// set pool
// NewAppendRemotePoolAddressesInstruction
