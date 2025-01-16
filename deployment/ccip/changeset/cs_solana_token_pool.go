package changeset

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	"github.com/smartcontractkit/chainlink/deployment"

	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

type AddTokenPoolConfig struct {
	ChainSelector    uint64
	PoolType         string
	RampAuthority    string
	Authority        string
	TokenName        string
	TokenProgramName string
}

var _ deployment.ChangeSet[AddTokenPoolConfig] = AddTokenPool

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

type SetupTokenPoolForChainConfig struct {
	ChainSelector       uint64
	RemoteChainSelector uint64
	TokenName           string
	TokenProgramName    string
	// TODO: maybe change this to native types instead of using token_pool types
	RemoteConfig      token_pool.RemoteConfig
	InboundRateLimit  token_pool.RateLimitConfig
	OutboundRateLimit token_pool.RateLimitConfig
}

var _ deployment.ChangeSet[SetupTokenPoolForChainConfig] = SetupTokenPoolForChain

func SetupTokenPoolForChain(e deployment.Environment, cfg SetupTokenPoolForChainConfig) (deployment.ChangesetOutput, error) {
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

// Add billing changesets
// Everything required for router
// Add logs
