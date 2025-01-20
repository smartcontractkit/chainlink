package changeset

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	"github.com/smartcontractkit/chainlink/deployment"

	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solTestConfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

const (
	ViaGetCcipAdminInstruction RegisterTokenAdminRegistryType = iota
	ViaOwnerInstruction
)

var _ deployment.ChangeSet[AddTokenPoolConfig] = AddTokenPool
var _ deployment.ChangeSet[SetupTokenPoolForRemoteChainConfig] = SetupTokenPoolForRemoteChain
var _ deployment.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ deployment.ChangeSet[TransferAndAcceptAdminRoleTokenAdminRegistryConfig] = TransferAndAcceptAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[UpdateTokenPoolConfig] = UpdateTokenPool

// ADD TOKEN POOL
type AddTokenPoolConfig struct {
	ChainSelector    uint64
	PoolType         string
	RampAuthority    string
	Authority        string
	TokenName        string
	TokenProgramName string
}

func AddTokenPool(e deployment.Environment, cfg AddTokenPoolConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolTokenPool.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}
	token_pool.SetProgramID(chainState.SolTokenPool)

	tokenProgramId, err := deployment.GetTokenProgramID(cfg.TokenProgramName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	poolType, err := deployment.GetPoolType(cfg.PoolType)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	tokenPubKey, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Convert string addresses to public keys
	rampAuthorityPubKey := solana.MustPublicKeyFromBase58(cfg.RampAuthority)
	authorityPubKey := solana.MustPublicKeyFromBase58(cfg.Authority)

	// TODO: this will break if we use a tokenPoolProgram different from the one mentioned in chailink-ccip
	// the programId is hardcoded inside the function
	poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// TODO: this will break if we use a tokenPoolProgram different from the one mentioned in chailink-ccip
	// the programId is hardcoded inside the function
	poolSigner, err := solTokenUtil.TokenPoolSignerAddress(tokenPubKey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// ata for token pool
	createI, tokenPoolATA, err := solTokenUtil.CreateAssociatedTokenAccount(tokenProgramId, tokenPubKey, poolSigner, chain.DeployerKey.PublicKey())
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// initialize token pool for token
	poolInitI, err := token_pool.NewInitializeInstruction(poolType, rampAuthorityPubKey, poolConfig, tokenPubKey, poolSigner, authorityPubKey, solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// make pool mint_authority for token (required for burn/mint)
	authI, err := solTokenUtil.SetTokenMintAuthority(tokenProgramId, poolSigner, tokenPubKey, chain.DeployerKey.PublicKey())
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

// SETUP TOKEN POOL FOR CHAIN
type SetupTokenPoolForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	TokenName           string
	TokenProgramName    string
	// TODO: maybe change this to native types instead of using token_pool types
	RemoteConfig      token_pool.RemoteConfig
	InboundRateLimit  token_pool.RateLimitConfig
	OutboundRateLimit token_pool.RateLimitConfig
}

func SetupTokenPoolForRemoteChain(e deployment.Environment, cfg SetupTokenPoolForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolTokenPool.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}
	token_pool.SetProgramID(chainState.SolTokenPool)

	tokenPubKey, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("ccip_tokenpool_chainconfig"),
			binary.LittleEndian.AppendUint64([]byte{}, cfg.RemoteChainSelector),
			tokenPubKey.Bytes(),
		},
		chainState.SolTokenPool,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	ixConfigure, err := token_pool.NewSetChainRemoteConfigInstruction(cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.RemoteConfig,
		poolConfig,
		chainPDA,
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
		poolConfig,
		chainPDA,
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

// TOKEN ADMIN REGISTRY
type RegisterTokenAdminRegistryType int
type RegisterTokenAdminRegistryConfig struct {
	ChainSelector       uint64
	TokenName           string
	TokenPoolAdmin      string
	AuthorityPrivateKey string
	RegisterType        RegisterTokenAdminRegistryType
}

func RegisterTokenAdminRegistry(e deployment.Environment, cfg RegisterTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolTokenPool.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}

	tokenPubKey, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Convert string addresses to public keys
	authorityPrivKey := solana.MustPrivateKeyFromBase58(cfg.AuthorityPrivateKey)
	var instruction *ccip_router.Instruction

	if cfg.RegisterType == ViaGetCcipAdminInstruction {
		tokenPoolAdminPubKey := solana.MustPublicKeyFromBase58(cfg.TokenPoolAdmin)
		instruction, err = ccip_router.NewRegisterTokenAdminRegistryViaGetCcipAdminInstruction(
			tokenPubKey,
			tokenPoolAdminPubKey,
			GetRouterConfigPDA(chainState.SolCcipRouter),
			GetTokenAdminRegistryPDA(chainState.SolCcipRouter, tokenPubKey),
			authorityPrivKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	} else if cfg.RegisterType == ViaOwnerInstruction {
		instruction, err = ccip_router.NewRegisterTokenAdminRegistryViaOwnerInstruction(
			GetTokenAdminRegistryPDA(chainState.SolCcipRouter, tokenPubKey),
			tokenPubKey,
			authorityPrivKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	} else {
		return deployment.ChangesetOutput{}, fmt.Errorf("Unsupported RegisterType")
	}

	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(authorityPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

type TransferAndAcceptAdminRoleTokenAdminRegistryConfig struct {
	ChainSelector               uint64
	TokenName                   string
	TokenPoolAdminPrivateKey    string
	NewTokenPoolAdminPrivateKey string
}

func TransferAndAcceptAdminRoleTokenAdminRegistry(e deployment.Environment, cfg TransferAndAcceptAdminRoleTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolTokenPool.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}

	tokenPubKey, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Convert string addresses to public keys
	tokenPoolAdminPrivKey := solana.MustPrivateKeyFromBase58(cfg.TokenPoolAdminPrivateKey)
	newTokenPoolAdminPrivKey := solana.MustPrivateKeyFromBase58(cfg.NewTokenPoolAdminPrivateKey)
	ix1, err := ccip_router.NewTransferAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		newTokenPoolAdminPrivKey.PublicKey(),
		GetTokenAdminRegistryPDA(chainState.SolCcipRouter, tokenPubKey),
		tokenPoolAdminPrivKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ix2, err := ccip_router.NewAcceptAdminRoleTokenAdminRegistryInstruction(
		tokenPubKey,
		GetTokenAdminRegistryPDA(chainState.SolCcipRouter, tokenPubKey),
		newTokenPoolAdminPrivKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{ix1, ix2}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(tokenPoolAdminPrivKey, newTokenPoolAdminPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// UPDATE TOKEN POOL
type UpdateTokenPoolConfig struct {
	ChainSelector       uint64
	TokenName           string
	AuthorityPrivateKey string
	PoolLookupTable     string
}

func UpdateTokenPool(e deployment.Environment, cfg UpdateTokenPoolConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolTokenPool.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("token pool not found in existing state, deploy the prerequisites first")
	}

	tokenPubKey, err := deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Convert string addresses to public keys
	authorityPrivKey := solana.MustPrivateKeyFromBase58(cfg.AuthorityPrivateKey)
	lookupTablePubKey := solana.MustPublicKeyFromBase58(cfg.PoolLookupTable)
	base := ccip_router.NewSetPoolInstruction(
		tokenPubKey,
		lookupTablePubKey,
		GetTokenAdminRegistryPDA(chainState.SolCcipRouter, tokenPubKey),
		authorityPrivKey.PublicKey(),
	)
	base.AccountMetaSlice = append(base.AccountMetaSlice, solana.Meta(lookupTablePubKey))
	instruction, err := base.ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	instructions := []solana.Instruction{instruction}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(authorityPrivKey))
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

// BILLING
type AddBillingTokenPoolConfig struct {
	ChainSelector    uint64
	TokenName        string
	TokenProgramName string
	TokenPubKey      string
	Config           ccip_router.BillingTokenConfig
}

func AddBillingToken(e deployment.Environment, cfg AddBillingTokenPoolConfig) (deployment.ChangesetOutput, error) {

	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	if chainState.SolCcipRouter.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("ccip router not found in existing state, deploy the prerequisites first")
	}
	ccip_router.SetProgramID(chainState.SolCcipRouter)

	var tokenPubKey solana.PublicKey
	if cfg.TokenPubKey == "" {
		tokenPubKey, err = deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	} else {
		tokenPubKey = solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	}

	fmt.Println("tokenPubKey", tokenPubKey.String())

	billingConfigPDA, _, _ := solana.FindProgramAddress([][]byte{solTestConfig.BillingTokenConfigPrefix, tokenPubKey.Bytes()}, chainState.SolCcipRouter)
	fmt.Println("billingConfigPDA", billingConfigPDA.String())

	var token0ConfigAccount ccip_router.BillingTokenConfigWrapper
	err = solCommonUtil.GetAccountDataBorshInto(context.Background(), chain.Client, billingConfigPDA, solRpc.CommitmentFinalized, &token0ConfigAccount)
	if err == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("billing token config already exists")
	}
	if err.Error() != "not found" {
		return deployment.ChangesetOutput{}, err
	}

	billingSignerPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("fee_billing_signer")}, chainState.SolCcipRouter)
	fmt.Println("billingSignerPDA", billingSignerPDA.String())

	tokenProgramId, _ := deployment.GetTokenProgramID(cfg.TokenProgramName)
	fmt.Println("tokenProgramId", tokenProgramId.String())

	token2022Receiver, _, _ := solTokenUtil.FindAssociatedTokenAddress(tokenProgramId, tokenPubKey, billingSignerPDA)
	fmt.Println("token2022Receiver", token2022Receiver.String())

	routerConfigPDA := GetRouterConfigPDA(chainState.SolCcipRouter)
	fmt.Println("routerConfigPDA", routerConfigPDA.String())

	fmt.Println("deployerKey", chain.DeployerKey.PublicKey().String())
	fmt.Println("ata.ProgramID", ata.ProgramID)

	cfg.Config.Mint = tokenPubKey

	ixConfig, cerr := ccip_router.NewAddBillingTokenConfigInstruction(
		cfg.Config,
		routerConfigPDA,
		billingConfigPDA,
		tokenProgramId,
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

type BillingTokenForRemoteChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	TokenName           string
	TokenProgramName    string
	Config              ccip_router.TokenBilling
	TokenPubKey         string
}

func AddBillingTokenForRemoteChain(e deployment.Environment, cfg BillingTokenForRemoteChainConfig) (deployment.ChangesetOutput, error) {
	chain, ok := e.SolChains[cfg.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain selector %d not found in environment", cfg.ChainSelector)
	}
	state, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}
	var tokenPubKey solana.PublicKey
	if cfg.TokenPubKey == "" {
		tokenPubKey, err = deployment.FindTokenAddress(e, cfg.ChainSelector, cfg.TokenName)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	} else {
		tokenPubKey = solana.MustPublicKeyFromBase58(cfg.TokenPubKey)
	}
	remoteBillingPDA, _, err := solana.FindProgramAddress([][]byte{[]byte("ccip_tokenpool_billing"), binary.LittleEndian.AppendUint64([]byte{}, cfg.RemoteChainSelector), tokenPubKey.Bytes()}, chainState.SolCcipRouter)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	ix, err := ccip_router.NewSetTokenBillingInstruction(
		cfg.RemoteChainSelector,
		tokenPubKey,
		cfg.Config,
		GetRouterConfigPDA(chainState.SolCcipRouter),
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

// wsol billing (will it work with above billing?)
// add test helpers for funding and approvals ?
/*
t.Run("Billing", func(t *testing.T) {
			ix, err := ccip_router.NewSetTokenBillingInstruction(
			config.EvmChainSelector,
			token0.Mint.PublicKey(),
			ccip_router.TokenBilling{},
			config.RouterConfigPDA,
			token0.Billing[config.EvmChainSelector],
			anotherAdmin.PublicKey(),
			solana.SystemProgramID
			).ValidateAndBuild()
			require.NoError(t, err)
			testutils.SendAndConfirm(ctx, t, solanaGoClient, []solana.Instruction{ix}, anotherAdmin, config.DefaultCommitment)
		})
*/
