package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/base_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/cctp_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	solanashared "github.com/smartcontractkit/chainlink/deployment"
	ccipshared "github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
)

type TokenPoolView struct {
	PoolType             string                                     `json:"poolType,omitempty"`
	PoolMetadata         string                                     `json:"poolMetadata,omitempty"`
	UpgradeAuthority     string                                     `json:"upgradeAuthority,omitempty"`
	TokenPoolChainConfig map[uint64]map[string]TokenPoolChainConfig `json:"chainConfig,omitempty"`
	TokenPoolState       map[string]TokenPoolState                  `json:"state,omitempty"`
}

type TokenPoolState struct {
	PDA                   string   `json:"pda,omitempty"`
	TokenProgram          string   `json:"tokenProgram,omitempty"`
	Mint                  string   `json:"mint,omitempty"`
	Decimals              uint8    `json:"decimals,omitempty"`
	PoolSigner            string   `json:"poolSigner,omitempty"`
	PoolTokenAccount      string   `json:"poolTokenAccount,omitempty"`
	Owner                 string   `json:"owner,omitempty"`
	ProposedOwner         string   `json:"proposedOwner,omitempty"`
	RateLimitAdmin        string   `json:"rateLimitAdmin,omitempty"`
	RouterOnrampAuthority string   `json:"routerOnrampAuthority,omitempty"`
	Router                string   `json:"router,omitempty"`
	Rebalancer            string   `json:"rebalancer,omitempty"`
	CanAcceptLiquidity    bool     `json:"canAcceptLiquidity,omitempty"`
	ListEnabled           bool     `json:"listEnabled,omitempty"`
	AllowList             []string `json:"allowList,omitempty"`
	RmnRemote             string   `json:"rmnRemote,omitempty"`
}

type TokenPoolChainConfig struct {
	PDA               string                        `json:"pda,omitempty"`
	PoolAddresses     []string                      `json:"poolAddresses,omitempty"`
	TokenAddress      string                        `json:"tokenAddress,omitempty"`
	Decimals          uint8                         `json:"decimals,omitempty"`
	InboundRateLimit  TokenPoolRateLimitTokenBucket `json:"inboundRateLimit"`
	OutboundRateLimit TokenPoolRateLimitTokenBucket `json:"outboundRateLimit"`
	CCTPChainConfig   *cctp_token_pool.CctpChain    `json:"cctpChainConfig,omitempty"`
}

type TokenPoolRateLimitTokenBucket struct {
	Tokens      uint64 `json:"tokens"`
	LastUpdated uint64 `json:"lastUpdated"`
	Enabled     bool   `json:"enabled"`
	Capacity    uint64 `json:"capacity"`
	Rate        uint64 `json:"rate"`
}

func GenerateTokenPoolView(chain cldf_solana.Chain, program solana.PublicKey, remoteChains []uint64, tokens []solana.PublicKey, poolType string, poolMetadata string) (TokenPoolView, error) {
	view := TokenPoolView{}
	view.PoolType = poolType
	progDataAddr, err := solanashared.GetProgramDataAddress(chain.Client, program)
	if err != nil {
		return view, fmt.Errorf("failed to get program data address for program %s: %w", program.String(), err)
	}
	authority, _, err := solanashared.GetUpgradeAuthority(chain.Client, progDataAddr)
	if err != nil {
		return view, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataAddr.String(), err)
	}
	view.UpgradeAuthority = authority.String()
	view.PoolMetadata = poolMetadata
	view.TokenPoolState = make(map[string]TokenPoolState)
	view.TokenPoolChainConfig = make(map[uint64]map[string]TokenPoolChainConfig)
	for _, remote := range remoteChains {
		view.TokenPoolChainConfig[remote] = make(map[string]TokenPoolChainConfig)
		// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
		for _, token := range tokens {
			remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(remote, token, program)
			if baseConfig, cctpConfig, err := fetchChainConfig(chain, remoteChainConfigPDA, poolType); err == nil && baseConfig != nil {
				view.TokenPoolChainConfig[remote][token.String()] = TokenPoolChainConfig{
					PDA:           remoteChainConfigPDA.String(),
					PoolAddresses: make([]string, len(baseConfig.Remote.PoolAddresses)),
					TokenAddress:  shared.GetAddressFromBytes(remote, baseConfig.Remote.TokenAddress.Address),
					Decimals:      baseConfig.Remote.Decimals,
					InboundRateLimit: TokenPoolRateLimitTokenBucket{
						Tokens:      baseConfig.InboundRateLimit.Tokens,
						LastUpdated: baseConfig.InboundRateLimit.LastUpdated,
						Enabled:     baseConfig.InboundRateLimit.Cfg.Enabled,
						Capacity:    baseConfig.InboundRateLimit.Cfg.Capacity,
						Rate:        baseConfig.InboundRateLimit.Cfg.Rate},
					OutboundRateLimit: TokenPoolRateLimitTokenBucket{
						Tokens:      baseConfig.OutboundRateLimit.Tokens,
						LastUpdated: baseConfig.OutboundRateLimit.LastUpdated,
						Enabled:     baseConfig.OutboundRateLimit.Cfg.Enabled,
						Capacity:    baseConfig.OutboundRateLimit.Cfg.Capacity,
						Rate:        baseConfig.OutboundRateLimit.Cfg.Rate},
					CCTPChainConfig: cctpConfig,
				}
				for i, addr := range baseConfig.Remote.PoolAddresses {
					view.TokenPoolChainConfig[remote][token.String()].PoolAddresses[i] = shared.GetAddressFromBytes(remote, addr.Address)
				}
			}
		}
	}
	// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
	for _, token := range tokens {
		programData := solTestTokenPool.State{}
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(token, program)
		if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &programData); err == nil {
			view.TokenPoolState[token.String()] = TokenPoolState{
				PDA:                   poolConfigPDA.String(),
				TokenProgram:          programData.Config.TokenProgram.String(),
				Mint:                  programData.Config.Mint.String(),
				Decimals:              programData.Config.Decimals,
				PoolSigner:            programData.Config.PoolSigner.String(),
				PoolTokenAccount:      programData.Config.PoolTokenAccount.String(),
				Owner:                 programData.Config.Owner.String(),
				ProposedOwner:         programData.Config.ProposedOwner.String(),
				RateLimitAdmin:        programData.Config.RateLimitAdmin.String(),
				RouterOnrampAuthority: programData.Config.RouterOnrampAuthority.String(),
				Router:                programData.Config.Router.String(),
				Rebalancer:            programData.Config.Rebalancer.String(),
				CanAcceptLiquidity:    programData.Config.CanAcceptLiquidity,
				ListEnabled:           programData.Config.ListEnabled,
				AllowList:             make([]string, len(programData.Config.AllowList)),
				RmnRemote:             programData.Config.RmnRemote.String(),
			}
			for i, addr := range programData.Config.AllowList {
				view.TokenPoolState[token.String()].AllowList[i] = addr.String()
			}
		}
	}
	return view, nil
}

func fetchChainConfig(chain cldf_solana.Chain, chainConfigPDA solana.PublicKey, poolType string) (*base_token_pool.BaseChain, *cctp_token_pool.CctpChain, error) {
	switch poolType {
	case ccipshared.BurnMintTokenPool.String(), ccipshared.LockReleaseTokenPool.String():
		var remoteChainConfigAccount solTestTokenPool.ChainConfig
		if err := chain.GetAccountDataBorshInto(context.Background(), chainConfigPDA, &remoteChainConfigAccount); err != nil {
			return nil, nil, err
		}
		return &remoteChainConfigAccount.Base, nil, nil
	case ccipshared.CCTPTokenPool.String():
		var remoteChainConfigAccount cctp_token_pool.ChainConfig
		if err := chain.GetAccountDataBorshInto(context.Background(), chainConfigPDA, &remoteChainConfigAccount); err != nil {
			return nil, nil, err
		}
		return &remoteChainConfigAccount.Base, &remoteChainConfigAccount.Cctp, nil
	default:
		return nil, nil, fmt.Errorf("unsupported token pool type %s", poolType)
	}
}
