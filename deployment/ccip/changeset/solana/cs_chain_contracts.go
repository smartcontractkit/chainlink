package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

var _ deployment.ChangeSet[AddRemoteChainToSolanaConfig] = AddRemoteChainToSolana
var _ deployment.ChangeSet[TokenPoolConfig] = AddTokenPool
var _ deployment.ChangeSet[RemoteChainTokenPoolConfig] = SetupTokenPoolForRemoteChain
var _ deployment.ChangeSet[cs.SetOCR3OffRampConfig] = SetOCR3ConfigSolana
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingToken
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain
var _ deployment.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ deployment.ChangeSet[TransferAdminRoleTokenAdminRegistryConfig] = TransferAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[AcceptAdminRoleTokenAdminRegistryConfig] = AcceptAdminRoleTokenAdminRegistry

// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName string) (solana.PublicKey, error) {
	tokenPrograms := map[string]solana.PublicKey{
		deployment.SPLTokens:     solana.TokenProgramID, // not used yet
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
	if tokenPubKey.Equals(chainState.LinkToken) || tokenPubKey.Equals(chainState.WSOL) {
		return nil
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
	// addressing errcheck in the next PR
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	var routerConfigAccount solRouter.Config
	err := chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, deploy the prerequisites first %d", chain.Selector)
	}
	return nil
}

// ADD REMOTE CHAIN
type AddRemoteChainToSolanaConfig struct {
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]map[uint64]RemoteChainConfigSolana
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *cs.MCMSConfig
}

// https://github.com/smartcontractkit/chainlink-ccip/blob/771fb9957d818253d833431e7e980669984e1d6a/chains/solana/gobindings/ccip_router/types.go#L1141
// https://github.com/smartcontractkit/chainlink-ccip/blob/771fb9957d818253d833431e7e980669984e1d6a/chains/solana/contracts/tests/ccip/ccip_router_test.go#L130
// We are not using solRouter.SourceChainConfig because that would involve the user
// converting the onRamp address into [2][64]byte{} which is not intuitive.
// The solRouter.DestChainConfig on the other hand has a lot of fields and most of them are uint
// So we are using that directly instead of copying over the fields here to reduce
// overhead cost if that type is bumped in chainlink-ccip
type RemoteChainConfigSolana struct {
	// source
	EnabledAsSource bool
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
		var routerConfigAccount solRouter.Config
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

	for remoteChainSel, update := range updates {
		var onRampBytes [64]byte
		// already verified, skipping errcheck
		remoteChainFamily, _ := chainsel.GetSelectorFamily(remoteChainSel)
		switch remoteChainFamily {
		case chainsel.FamilySolana:
			return deployment.ChangesetOutput{}, fmt.Errorf("cannot add solana chain as remote chain")
		case chainsel.FamilyEVM:
			onRampAddress := s.Chains[remoteChainSel].OnRamp.Address().String()
			if onRampAddress == "" {
				return deployment.ChangesetOutput{}, fmt.Errorf("onramp address not found for chain %d", remoteChainSel)
			}
			addressBytes := []byte(onRampAddress)
			copy(onRampBytes[:], addressBytes)
		}

		validSourceChainConfig := solRouter.SourceChainConfig{
			OnRamp:    [2][64]byte{onRampBytes, [64]byte{}},
			IsEnabled: update.EnabledAsSource,
		}
		// addressing errcheck in the next PR
		routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterID)
		destChainStatePDA, _ := solState.FindDestChainStatePDA(remoteChainSel, ccipRouterID)
		sourceChainStatePDA, _ := solState.FindSourceChainStatePDA(remoteChainSel, ccipRouterID)

		instruction, err := solRouter.NewAddChainSelectorInstruction(
			remoteChainSel,
			validSourceChainConfig,
			update.DestinationConfig,
			sourceChainStatePDA,
			destChainStatePDA,
			routerConfigPDA,
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
		// addressing errcheck in the next PR
		routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterID)
		routerStatePDA, _, _ := solState.FindStatePDA(ccipRouterID)
		for _, arg := range args {
			instruction, err := solRouter.NewSetOcrConfigInstruction(
				arg.OCRPluginType,
				solRouter.Ocr3ConfigInfo{
					ConfigDigest:                   arg.ConfigDigest,
					F:                              arg.F,
					IsSignatureVerificationEnabled: btoi(arg.IsSignatureVerificationEnabled),
				},
				arg.Signers,
				arg.Transmitters,
				routerConfigPDA,
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
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the token pool first for chain %d", cfg.ChainSelector)
	}
	if _, err := GetPoolType(cfg.PoolType); err != nil {
		return err
	}
	if _, err := GetTokenProgramID(cfg.TokenProgramName); err != nil {
		return err
	}

	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, state.SolChains[cfg.ChainSelector].TokenPool)
	if err != nil {
		return err
	}
	chain := e.SolChains[cfg.ChainSelector]
	var poolConfigAccount token_pool.Config
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err == nil {
		return fmt.Errorf("token pool config already exists for token %s", tokenPubKey.String())
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
	authorityPubKey := solana.MustPublicKeyFromBase58(cfg.Authority)
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	// verified
	tokenprogramID, _ := GetTokenProgramID(cfg.TokenProgramName)
	poolType, _ := GetPoolType(cfg.PoolType)
	poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.TokenPool)
	poolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, chainState.TokenPool)

	// addressing errcheck in the next PR
	rampAuthorityPubKey, _, _ := solState.FindExternalExecutionConfigPDA(chainState.Router)

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
		authorityPubKey, // this is assumed to be chain.DeployerKey for now (owner of token pool)
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
	if err := chain.Confirm(instructions); err != nil {
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
	if err := commonValidation(e, cfg.ChainSelector, solana.MustPublicKeyFromBase58(cfg.TokenPubKey)); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the prerequisites first for chain %d", cfg.ChainSelector)
	}

	chain := e.SolChains[cfg.ChainSelector]

	// check if pool config exists (cannot do remote setup without it)
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, state.SolChains[cfg.ChainSelector].TokenPool)
	if err != nil {
		return err
	}
	var poolConfigAccount token_pool.Config
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err != nil {
		return fmt.Errorf("token pool config not found, call AddTokenPool first for chain %d", cfg.ChainSelector)
	}

	// check if existing pool setup already has this remote chain configured
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.TokenPool)
	if err != nil {
		return err
	}
	var remoteChainConfigAccount token_pool.ChainConfig
	if err := chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount); err == nil {
		return fmt.Errorf("remote chain config already exists for token %s", tokenPubKey.String())
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

	// verified
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
type BillingTokenConfig struct {
	ChainSelector    uint64
	TokenPubKey      string
	TokenProgramName string
	Config           solRouter.BillingTokenConfig
}

func (cfg BillingTokenConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	_, err = GetTokenProgramID(cfg.TokenProgramName)
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
	billingConfigPDA, _, err := solState.FindFeeBillingTokenConfigPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return err
	}
	var token0ConfigAccount solRouter.BillingTokenConfigWrapper
	err = chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount)
	if err == nil {
		return fmt.Errorf("billing token config already exists for token %s", tokenPubKey.String())
	}
	return nil
}

func AddBillingToken(e deployment.Environment, cfg BillingTokenConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	solRouter.SetProgramID(chainState.Router)

	// verified
	tokenprogramID, _ := GetTokenProgramID(cfg.TokenProgramName)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	billingConfigPDA, _, _ := solState.FindFeeBillingTokenConfigPDA(tokenPubKey, chainState.Router)

	// addressing errcheck in the next PR
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenPubKey, billingSignerPDA)

	ixConfig, cerr := solRouter.NewAddBillingTokenConfigInstruction(
		cfg.Config,
		routerConfigPDA,
		billingConfigPDA,
		tokenprogramID,
		tokenPubKey,
		token2022Receiver,
		chain.DeployerKey.PublicKey(), // ccip admin
		billingSignerPDA,
		ata.ProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if cerr != nil {
		return deployment.ChangesetOutput{}, cerr
	}

	instructions := []solana.Instruction{ixConfig}
	err := chain.Confirm(instructions)
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
	Config              solRouter.TokenBilling
	TokenPubKey         string
}

func (cfg BillingTokenForRemoteChainConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return fmt.Errorf("router validation failed: %w", err)
	}
	// check if desired state already exists
	remoteBillingPDA, _, err := solState.FindCcipTokenpoolBillingPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.Router)
	if err != nil {
		return err
	}
	var remoteBillingAccount solRouter.TokenBilling
	err = chain.GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount)
	if err == nil {
		return fmt.Errorf("billing token config already exists for token %s", tokenPubKey.String())
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
	// verified
	remoteBillingPDA, _, _ := solState.FindCcipTokenpoolBillingPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.Router)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	ix, err := solRouter.NewSetTokenBillingInstruction(
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
	ChainSelector           uint64
	TokenPubKey             string
	TokenAdminRegistryAdmin string
	RegisterType            RegisterTokenAdminRegistryType
}

func (cfg RegisterTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	if cfg.RegisterType != ViaGetCcipAdminInstruction && cfg.RegisterType != ViaOwnerInstruction {
		return fmt.Errorf("invalid register type, valid types are %d and %d", ViaGetCcipAdminInstruction, ViaOwnerInstruction)
	}

	if cfg.RegisterType == ViaOwnerInstruction && cfg.TokenAdminRegistryAdmin != "" {
		return fmt.Errorf("token admin registry should be empty for via owner instruction")
	}

	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return err
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	err = chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	if err == nil {
		return fmt.Errorf("token admin registry already exists for token %s", tokenPubKey.String())
	}
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

	// verified
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)

	var instruction *solRouter.Instruction
	var err error
	switch cfg.RegisterType {
	// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
	case ViaGetCcipAdminInstruction:
		tokenAdminRegistryAdmin := solana.MustPublicKeyFromBase58(cfg.TokenAdminRegistryAdmin)
		instruction, err = solRouter.NewRegisterTokenAdminRegistryViaGetCcipAdminInstruction(
			tokenPubKey,
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			routerConfigPDA,
			tokenAdminRegistryPDA,         // this gets created
			chain.DeployerKey.PublicKey(), // (ccip admin)
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	case ViaOwnerInstruction:
		// the token mint authority signs and makes itself the authority of the tokenAdminRegistry PDA
		instruction, err = solRouter.NewRegisterTokenAdminRegistryViaOwnerInstruction(
			routerConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // (token mint authority) becomes the authority of the tokenAdminRegistry PDA
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}
	// if we want to have a different authority, we will need to add the corresponding singer here
	// for now we are assuming both token owner and ccip admin will always be deployer key
	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// TRANSFER AND ACCEPT TOKEN ADMIN REGISTRY
type TransferAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector                  uint64
	TokenPubKey                    string
	NewRegistryAdminPublicKey      string
	CurrentRegistryAdminPrivateKey string
}

func (cfg TransferAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}

	currentRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.CurrentRegistryAdminPrivateKey)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	if currentRegistryAdminPrivateKey.PublicKey().Equals(newRegistryAdminPubKey) {
		return fmt.Errorf("new registry admin public key cannot be the same as current registry admin public key for token %s", tokenPubKey.String())
	}

	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return err
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	err = chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	if err != nil {
		return fmt.Errorf("token admin registry not found for token %s, cannot transfer admin role", tokenPubKey.String())
	}
	// check if passed admin is the current admin
	if !tokenAdminRegistryAccount.Administrator.Equals(currentRegistryAdminPrivateKey.PublicKey()) {
		return fmt.Errorf("current registry admin private key does not match administrator for token %s", tokenPubKey.String())
	}
	return nil
}

func TransferAdminRoleTokenAdminRegistry(e deployment.Environment, cfg TransferAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	currentRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.CurrentRegistryAdminPrivateKey)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	ix1, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		newRegistryAdminPubKey,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		currentRegistryAdminPrivateKey.PublicKey(), // as we are assuming this is the default authority for everything in the beginning
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	instructions := []solana.Instruction{ix1}
	// the existing authority will have to sign the transfer
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(currentRegistryAdminPrivateKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

type AcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector              uint64
	TokenPubKey                string
	NewRegistryAdminPrivateKey string
}

func (cfg AcceptAdminRoleTokenAdminRegistryConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	err := commonValidation(e, cfg.ChainSelector, tokenPubKey)
	if err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return err
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	err = chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	if err != nil {
		return fmt.Errorf("token admin registry not found for token %s, cannot accept admin role", tokenPubKey.String())
	}
	// check if accepting admin is the pending admin
	newRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.NewRegistryAdminPrivateKey)
	newRegistryAdminPublicKey := newRegistryAdminPrivateKey.PublicKey()
	if !tokenAdminRegistryAccount.PendingAdministrator.Equals(newRegistryAdminPublicKey) {
		return fmt.Errorf("new admin public key does not match pending registry admin role for token %s", tokenPubKey.String())
	}
	return nil
}

func AcceptAdminRoleTokenAdminRegistry(e deployment.Environment, cfg AcceptAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	newRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.NewRegistryAdminPrivateKey)

	// verified
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)

	ix1, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		newRegistryAdminPrivateKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{ix1}
	// the new authority will have to sign the acceptance
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(newRegistryAdminPrivateKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// TODO (all look up table related changesets):
// Update look up tables with tokens and pools
// Set Pool (https://smartcontract-it.atlassian.net/browse/INTAUTO-437)
// NewAppendRemotePoolAddressesInstruction (https://smartcontract-it.atlassian.net/browse/INTAUTO-436)
