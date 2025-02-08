package solana

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	chainsel "github.com/smartcontractkit/chain-selectors"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"

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
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, deployment.SPLTokens, deployment.SPL2022Tokens)
	}
	return programID, nil
}

// GetPoolType returns the token pool type constant for the given string
func GetPoolType(poolType string) (solTokenPool.PoolType, error) {
	poolTypes := map[string]solTokenPool.PoolType{
		"LockAndRelease": solTokenPool.LockAndRelease_PoolType,
		"BurnAndMint":    solTokenPool.BurnAndMint_PoolType,
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
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
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
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	return nil
}

func validateRouterConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
	if chainState.Router.IsZero() {
		return fmt.Errorf("router not found in existing state, deploy the router first for chain %d", chain.Selector)
	}
	// addressing errcheck in the next PR
	var routerConfigAccount solRouter.Config
	err := chain.GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func validateFeeQuoterConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
	if chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("fee quoter not found in existing state, deploy the fee quoter first for chain %d", chain.Selector)
	}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	return nil
}

// ADD REMOTE CHAIN
type AddRemoteChainToSolanaConfig struct {
	ChainSelector uint64
	// UpdatesByChain is a mapping of SVM chain selector -> remote chain selector -> remote chain config update
	UpdatesByChain map[uint64]RemoteChainConfigSolana
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
	if err := commoncs.ValidateOwnershipSolana(e.GetContext(), cfg.MCMS != nil, e.SolChains[cfg.ChainSelector].DeployerKey.PublicKey(), chainState.Timelock, chainState.Router); err != nil {
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

	for remoteChainSel, update := range updates {
		var onRampBytes [64]byte
		// already verified, skipping errcheck
		remoteChainFamily, _ := chainsel.GetSelectorFamily(remoteChainSel)
		switch remoteChainFamily {
		case chainsel.FamilySolana:
			return fmt.Errorf("support for solana chain as remote chain is not implemented yet %d", remoteChainSel)
		case chainsel.FamilyEVM:
			onRampAddress := s.Chains[remoteChainSel].OnRamp.Address().String()
			if onRampAddress == "" {
				return fmt.Errorf("onramp address not found for chain %d", remoteChainSel)
			}
			addressBytes := []byte(onRampAddress)
			copy(onRampBytes[:], addressBytes)
		}

		// verified while loading state
		fqEvmDestChainPDA, _, _ := solState.FindFqDestChainPDA(remoteChainSel, feeQuoterID)
		destChainStatePDA, _ := solState.FindDestChainStatePDA(remoteChainSel, ccipRouterID)
		offRampSourceChainPDA, _, _ := solState.FindOfframpSourceChainPDA(remoteChainSel, s.SolChains[chainSel].OffRamp)

		solRouter.SetProgramID(ccipRouterID)
		routerIx, err := solRouter.NewAddChainSelectorInstruction(
			remoteChainSel,
			update.RouterDestinationConfig,
			destChainStatePDA,
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
			fqEvmDestChainPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return fmt.Errorf("failed to generate instructions: %w", err)
		}

		solOffRamp.SetProgramID(s.SolChains[chainSel].OffRamp)
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
		err = ab.Save(chainSel, destChainStatePDA.String(), tv)
		if err != nil {
			return fmt.Errorf("failed to save dest chain state to address book: %w", err)
		}

		// TODO
		// tv = deployment.NewTypeAndVersion(cs.RemoteSource, deployment.Version1_0_0)
		// tv.AddLabel(remoteChainSelStr)
		// err = ab.Save(chainSel, sourceChainStatePDA.String(), tv)
		// if err != nil {
		// 	return fmt.Errorf("failed to save source chain state to address book: %w", err)
		// }
	}

	return nil
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
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := cfg.Validate(e, state); err != nil {
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
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get don id for chain %d: %w", remote, err)
		}
		args, err := internal.BuildSetOCR3ConfigArgsSolana(donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build set ocr3 config args: %w", err)
		}
		// TODO: check if ocr3 has already been set
		// set, err := isOCR3ConfigSetSolana(e.Logger, e.Chains[remote], state.Chains[remote].OffRamp, args)
		var instructions []solana.Instruction
		offRampConfigPDA := solChains[remote].OffRampConfigPDA
		offRampStatePDA := solChains[remote].OffRampStatePDA
		solOffRamp.SetProgramID(solChains[remote].OffRamp)
		for _, arg := range args {
			instruction, err := solOffRamp.NewSetOcrConfigInstruction(
				arg.OCRPluginType,
				solOffRamp.Ocr3ConfigInfo{
					ConfigDigest:                   arg.ConfigDigest,
					F:                              arg.F,
					IsSignatureVerificationEnabled: btoi(arg.IsSignatureVerificationEnabled),
				},
				arg.Signers,
				arg.Transmitters,
				offRampConfigPDA,
				offRampStatePDA,
				e.SolChains[remote].DeployerKey.PublicKey(),
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			instructions = append(instructions, instruction)
		}
		if cfg.MCMS == nil {
			if err := e.SolChains[remote].Confirm(instructions); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}
	return deployment.ChangesetOutput{}, nil
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

	tokenPool := chainState.TokenPool
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
	}
	chain := e.SolChains[cfg.ChainSelector]
	var poolConfigAccount solTokenPool.Config
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err == nil {
		return fmt.Errorf("token pool config already exists for (mint: %s, pool: %s)", tokenPubKey.String(), tokenPool.String())
	}
	return nil
}

func AddTokenPool(e deployment.Environment, cfg TokenPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
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
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create associated token account for tokenpool (mint: %s, pool: %s): %w", tokenPubKey.String(), chainState.TokenPool.String(), err)
	}

	solTokenPool.SetProgramID(chainState.TokenPool)
	// initialize token pool for token
	poolInitI, err := solTokenPool.NewInitializeInstruction(
		poolType,
		rampAuthorityPubKey,
		poolConfigPDA,
		tokenPubKey,
		poolSigner,
		authorityPubKey, // this is assumed to be chain.DeployerKey for now (owner of token pool)
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	// make pool mint_authority for token (required for burn/mint)
	authI, err := solTokenUtil.SetTokenMintAuthority(
		tokenprogramID,
		poolSigner,
		tokenPubKey,
		authorityPubKey,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	instructions := []solana.Instruction{createI, poolInitI, authI}

	// add signer here if authority is different from deployer key
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
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
	RemoteConfig        solTokenPool.RemoteConfig
	InboundRateLimit    solTokenPool.RateLimitConfig
	OutboundRateLimit   solTokenPool.RateLimitConfig
}

func (cfg RemoteChainTokenPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, solana.MustPublicKeyFromBase58(cfg.TokenPubKey)); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy token pool for chain %d", cfg.ChainSelector)
	}

	chain := e.SolChains[cfg.ChainSelector]
	tokenPool := chainState.TokenPool

	// check if pool config exists (cannot do remote setup without it)
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
	}
	var poolConfigAccount solTokenPool.Config
	if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err != nil {
		return fmt.Errorf("token pool config not found (mint: %s, pool: %s): %w", tokenPubKey.String(), chainState.TokenPool.String(), err)
	}

	// check if existing pool setup already has this remote chain configured
	remoteChainConfigPDA, _, err := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, tokenPool)
	if err != nil {
		return fmt.Errorf("failed to get token pool remote chain config pda (remoteSelector: %d, mint: %s, pool: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String(), err)
	}
	var remoteChainConfigAccount solTokenPool.ChainConfig
	if err := chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount); err == nil {
		return fmt.Errorf("remote chain config already exists for (remoteSelector: %d, mint: %s, pool: %s)", cfg.RemoteChainSelector, tokenPubKey.String(), tokenPool.String())
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
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.TokenPool)

	solTokenPool.SetProgramID(chainState.TokenPool)
	ixConfigure, err := solTokenPool.NewInitChainRemoteConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig,
		poolConfigPDA,
		remoteChainConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	ixRates, err := solTokenPool.NewSetChainRateLimitInstruction(
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
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	instructions := []solana.Instruction{ixConfigure, ixRates}
	err = chain.Confirm(instructions)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

// ADD BILLING TOKEN
type BillingTokenConfig struct {
	ChainSelector    uint64
	TokenPubKey      string
	TokenProgramName string
	Config           solFeeQuoter.BillingTokenConfig
}

func (cfg BillingTokenConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	if _, err := GetTokenProgramID(cfg.TokenProgramName); err != nil {
		return err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return err
	}
	// check if already setup
	billingConfigPDA, _, err := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find billing token config pda (mint: %s, feeQuoter: %s): %w", tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	if err := chain.GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount); err == nil {
		return fmt.Errorf("billing token config already exists for (mint: %s, feeQuoter: %s)", tokenPubKey.String(), chainState.FeeQuoter.String())
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
	// verified
	tokenprogramID, _ := GetTokenProgramID(cfg.TokenProgramName)
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)

	// addressing errcheck in the next PR
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(chainState.Router)
	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenprogramID, tokenPubKey, billingSignerPDA)

	e.Logger.Infow("chainState.FeeQuoterConfigPDA", "feeQuoterConfigPDA", chainState.FeeQuoterConfigPDA.String())
	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	ixConfig, cerr := solFeeQuoter.NewAddBillingTokenConfigInstruction(
		cfg.Config,
		chainState.FeeQuoterConfigPDA,
		tokenBillingPDA,
		tokenprogramID,
		tokenPubKey,
		token2022Receiver,
		chain.DeployerKey.PublicKey(), // ccip admin
		billingSignerPDA,
		ata.ProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if cerr != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", cerr)
	}

	instructions := []solana.Instruction{ixConfig}
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Billing token added", "chainSelector", cfg.ChainSelector, "tokenPubKey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

// // ADD BILLING TOKEN FOR REMOTE CHAIN
type BillingTokenForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	Config              solFeeQuoter.TokenTransferFeeConfig
	TokenPubKey         string
}

func (cfg BillingTokenForRemoteChainConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateFeeQuoterConfig(chain, chainState); err != nil {
		return fmt.Errorf("fee quoter validation failed: %w", err)
	}
	// check if desired state already exists
	remoteBillingPDA, _, err := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)
	if err != nil {
		return fmt.Errorf("failed to find remote billing token config pda for (remoteSelector: %d, mint: %s, feeQuoter: %s): %w", cfg.RemoteChainSelector, tokenPubKey.String(), chainState.FeeQuoter.String(), err)
	}
	var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
	if err := chain.GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount); err == nil {
		return fmt.Errorf("billing token config already exists for (remoteSelector: %d, mint: %s, feeQuoter: %s)", cfg.RemoteChainSelector, tokenPubKey.String(), chainState.FeeQuoter.String())
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
	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(cfg.RemoteChainSelector, tokenPubKey, chainState.FeeQuoter)

	// fee_quoter.NewSetTokenTransferFeeConfigInstruction
	ix, err := solFeeQuoter.NewSetTokenTransferFeeConfigInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.Config,
		chainState.FeeQuoterConfigPDA,
		remoteBillingPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	instructions := []solana.Instruction{ix}
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Token billing set for remote chain", "chainSelector ", cfg.ChainSelector, "remoteChainSelector ", cfg.RemoteChainSelector, "tokenPubKey", tokenPubKey.String())
	return deployment.ChangesetOutput{}, nil
}

// // TOKEN ADMIN REGISTRY
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

	if cfg.TokenAdminRegistryAdmin == "" {
		return errors.New("token admin registry admin is required")
	}

	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
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
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
		return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), chainState.Router.String())
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
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	tokenAdminRegistryAdmin := solana.MustPublicKeyFromBase58(cfg.TokenAdminRegistryAdmin)

	var instruction *solRouter.Instruction
	var err error
	switch cfg.RegisterType {
	// the ccip admin signs and makes tokenAdminRegistryAdmin the authority of the tokenAdminRegistry PDA
	case ViaGetCcipAdminInstruction:
		instruction, err = solRouter.NewCcipAdminProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			chainState.RouterConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // (ccip admin)
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	case ViaOwnerInstruction:
		// the token mint authority signs and makes itself the authority of the tokenAdminRegistry PDA
		instruction, err = solRouter.NewOwnerProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // admin of the tokenAdminRegistry PDA
			chainState.RouterConfigPDA,
			tokenAdminRegistryPDA, // this gets created
			tokenPubKey,
			chain.DeployerKey.PublicKey(), // (token mint authority) becomes the authority of the tokenAdminRegistry PDA
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
		}
	}
	// if we want to have a different authority, we will need to add the corresponding singer here
	// for now we are assuming both token owner and ccip admin will always be deployer key
	instructions := []solana.Instruction{instruction}
	if err := chain.Confirm(instructions); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
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
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}

	currentRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.CurrentRegistryAdminPrivateKey)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	if currentRegistryAdminPrivateKey.PublicKey().Equals(newRegistryAdminPubKey) {
		return fmt.Errorf("new registry admin public key (%s) cannot be the same as current registry admin public key (%s) for token %s",
			newRegistryAdminPubKey.String(),
			currentRegistryAdminPrivateKey.PublicKey().String(),
			tokenPubKey.String(),
		)
	}

	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot transfer admin role", tokenPubKey.String(), chainState.Router.String())
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

	currentRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.CurrentRegistryAdminPrivateKey)
	newRegistryAdminPubKey := solana.MustPublicKeyFromBase58(cfg.NewRegistryAdminPublicKey)

	ix1, err := solRouter.NewTransferAdminRoleTokenAdminRegistryInstruction(
		newRegistryAdminPubKey,
		chainState.RouterConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		currentRegistryAdminPrivateKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}
	instructions := []solana.Instruction{ix1}
	// the existing authority will have to sign the transfer
	if err := chain.Confirm(instructions, solCommonUtil.AddSigners(currentRegistryAdminPrivateKey)); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
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
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
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
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot accept admin role", tokenPubKey.String(), chainState.Router.String())
	}
	// check if accepting admin is the pending admin
	newRegistryAdminPrivateKey := solana.MustPrivateKeyFromBase58(cfg.NewRegistryAdminPrivateKey)
	newRegistryAdminPublicKey := newRegistryAdminPrivateKey.PublicKey()
	if !tokenAdminRegistryAccount.PendingAdministrator.Equals(newRegistryAdminPublicKey) {
		return fmt.Errorf("new admin public key (%s) does not match pending registry admin role (%s) for token %s",
			newRegistryAdminPublicKey.String(),
			tokenAdminRegistryAccount.PendingAdministrator.String(),
			tokenPubKey.String(),
		)
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

	ix1, err := solRouter.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		chainState.RouterConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		newRegistryAdminPrivateKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
	}

	instructions := []solana.Instruction{ix1}
	// the new authority will have to sign the acceptance
	if err := chain.Confirm(instructions, solCommonUtil.AddSigners(newRegistryAdminPrivateKey)); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

type TokenPoolLookupTableConfig struct {
	ChainSelector uint64
	TokenPubKey   string
	TokenProgram  string // this can go as a address book tag
}

func (cfg TokenPoolLookupTableConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the token pool first for chain %d", cfg.ChainSelector)
	}

	// TODO: do we need to validate if everything that goes into the lookup table is already created ?
	return nil
}

func AddTokenPoolLookupTable(e deployment.Environment, cfg TokenPoolLookupTableConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chain := e.SolChains[cfg.ChainSelector]
	ctx := e.GetContext()
	client := chain.Client
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	authorityPrivKey := chain.DeployerKey // assuming the authority is the deployer key
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)

	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	tokenPoolChainConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, chainState.TokenPool)
	tokenPoolSigner, _ := solTokenUtil.TokenPoolSignerAddress(tokenPubKey, chainState.TokenPool)
	tokenProgram, _ := GetTokenProgramID(cfg.TokenProgram)
	poolTokenAccount, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgram, tokenPubKey, tokenPoolSigner)
	feeTokenConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, chainState.FeeQuoter)

	// the 'table' address is not derivable
	// but this will be stored in tokenAdminRegistryPDA as a part of the SetPool changeset
	// and tokenAdminRegistryPDA is derivable using token and router address
	table, err := solCommonUtil.CreateLookupTable(ctx, client, *authorityPrivKey)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	list := solana.PublicKeySlice{
		table,                   // 0
		tokenAdminRegistryPDA,   // 1
		chainState.TokenPool,    // 2
		tokenPoolChainConfigPDA, // 3 - writable
		poolTokenAccount,        // 4 - writable
		tokenPoolSigner,         // 5
		tokenProgram,            // 6
		tokenPubKey,             // 7 - writable
		feeTokenConfigPDA,       // 8
	}
	if err = solCommonUtil.ExtendLookupTable(ctx, client, table, *authorityPrivKey, list); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to extend lookup table for token pool (mint: %s): %w", tokenPubKey.String(), err)
	}
	if err := solCommonUtil.AwaitSlotChange(ctx, client); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to await slot change while extending lookup table: %w", err)
	}
	newAddressBook := deployment.NewMemoryAddressBook()
	tv := deployment.NewTypeAndVersion(cs.TokenPoolLookupTable, deployment.Version1_0_0)
	tv.Labels.Add(tokenPubKey.String())
	if err := newAddressBook.Save(cfg.ChainSelector, table.String(), tv); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to save tokenpool address lookup table: %w", err)
	}
	return deployment.ChangesetOutput{
		AddressBook: newAddressBook,
	}, nil
}

type SetPoolConfig struct {
	ChainSelector                     uint64
	TokenPubKey                       string
	TokenAdminRegistryAdminPrivateKey string
	PoolLookupTable                   string
	WritableIndexes                   []uint8
}

func (cfg SetPoolConfig) Validate(e deployment.Environment) error {
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	if err := commonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
		return err
	}
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	if chainState.TokenPool.IsZero() {
		return fmt.Errorf("token pool not found in existing state, deploy the token pool first for chain %d", cfg.ChainSelector)
	}
	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), chainState.Router.String(), err)
	}
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err != nil {
		return fmt.Errorf("token admin registry not found for (mint: %s, router: %s), cannot set pool", tokenPubKey.String(), chainState.Router.String())
	}
	return nil
}

// this sets the writable indexes of the token pool lookup table
func SetPool(e deployment.Environment, cfg SetPoolConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chain := e.SolChains[cfg.ChainSelector]
	state, _ := cs.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	tokenPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, chainState.Router)
	tokenAdminRegistryAdminPrivKey := solana.MustPrivateKeyFromBase58(cfg.TokenAdminRegistryAdminPrivateKey)
	lookupTablePubKey := solana.MustPublicKeyFromBase58(cfg.PoolLookupTable)

	base := solRouter.NewSetPoolInstruction(
		cfg.WritableIndexes,
		routerConfigPDA,
		tokenAdminRegistryPDA,
		tokenPubKey,
		lookupTablePubKey,
		tokenAdminRegistryAdminPrivKey.PublicKey(),
	)

	base.AccountMetaSlice = append(base.AccountMetaSlice, solana.Meta(lookupTablePubKey))
	instruction, err := base.ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(tokenAdminRegistryAdminPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// TODO:
// NewAppendRemotePoolAddressesInstruction (https://smartcontract-it.atlassian.net/browse/INTAUTO-436)
// OffRamp updates
