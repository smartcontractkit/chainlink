package solana

import (
	"context"

	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
)

type TokenPoolView struct {
	PoolType             string                                     `json:"poolType,omitempty"`
	PoolMetadata         string                                     `json:"poolMetadata,omitempty"`
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
	InboundRateLimit  TokenPoolRateLimitTokenBucket `json:"inboundRateLimit,omitempty"`
	OutboundRateLimit TokenPoolRateLimitTokenBucket `json:"outboundRateLimit,omitempty"`
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
	// skip this to avoid incurring fees on every state generation
	// switch poolType {
	// case solTestTokenPool.BurnAndMint_PoolType.String():
	// 	solBurnMintTokenPool.SetProgramID(program)
	// 	ixn, err := solBurnMintTokenPool.NewTypeVersionInstruction(solana.SysVarClockPubkey).ValidateAndBuild()
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to build instruction: %w", err)
	// 	}
	// 	result, err := solCommonUtil.SendAndConfirmWithLookupTables(context.Background(), chain.Client, []solana.Instruction{ixn}, *chain.DeployerKey, rpc.CommitmentConfirmed, nil)
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to confirm instruction: %w", err)
	// 	}
	// 	output, err := solCommonUtil.ExtractTypedReturnValue(context.Background(), result.Meta.LogMessages, program.String(), func(b []byte) string {
	// 		return string(b[4:])
	// 	})
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to extract typed return value: %w", err)
	// 	}
	// 	view.TypeAndVersion = output
	// case solTestTokenPool.LockAndRelease_PoolType.String():
	// 	solLockReleaseTokenPool.SetProgramID(program)
	// 	ixn, err := solLockReleaseTokenPool.NewTypeVersionInstruction(solana.SysVarClockPubkey).ValidateAndBuild()
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to build instruction: %w", err)
	// 	}
	// 	result, err := solCommonUtil.SendAndConfirmWithLookupTables(context.Background(), chain.Client, []solana.Instruction{ixn}, *chain.DeployerKey, rpc.CommitmentConfirmed, nil)
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to confirm instruction: %w", err)
	// 	}
	// 	output, err := solCommonUtil.ExtractTypedReturnValue(context.Background(), result.Meta.LogMessages, program.String(), func(b []byte) string {
	// 		return string(b[4:])
	// 	})
	// 	if err != nil {
	// 		return view, fmt.Errorf("failed to extract typed return value: %w", err)
	// 	}
	// 	view.TypeAndVersion = output
	// default:
	// 	return view, fmt.Errorf("unknown pool type %s", poolType)
	// }
	view.PoolMetadata = poolMetadata
	view.TokenPoolState = make(map[string]TokenPoolState)
	view.TokenPoolChainConfig = make(map[uint64]map[string]TokenPoolChainConfig)
	for _, remote := range remoteChains {
		view.TokenPoolChainConfig[remote] = make(map[string]TokenPoolChainConfig)
		// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
		for _, token := range tokens {
			remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(remote, token, program)
			var remoteChainConfigAccount solTestTokenPool.ChainConfig
			if err := chain.GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount); err == nil {
				view.TokenPoolChainConfig[remote][token.String()] = TokenPoolChainConfig{
					PDA:           remoteChainConfigPDA.String(),
					PoolAddresses: make([]string, len(remoteChainConfigAccount.Base.Remote.PoolAddresses)),
					TokenAddress:  shared.GetAddressFromBytes(remote, remoteChainConfigAccount.Base.Remote.TokenAddress.Address),
					Decimals:      remoteChainConfigAccount.Base.Remote.Decimals,
					InboundRateLimit: TokenPoolRateLimitTokenBucket{
						Tokens:      remoteChainConfigAccount.Base.InboundRateLimit.Tokens,
						LastUpdated: remoteChainConfigAccount.Base.InboundRateLimit.LastUpdated,
						Enabled:     remoteChainConfigAccount.Base.InboundRateLimit.Cfg.Enabled,
						Capacity:    remoteChainConfigAccount.Base.InboundRateLimit.Cfg.Capacity,
						Rate:        remoteChainConfigAccount.Base.InboundRateLimit.Cfg.Rate},
					OutboundRateLimit: TokenPoolRateLimitTokenBucket{
						Tokens:      remoteChainConfigAccount.Base.OutboundRateLimit.Tokens,
						LastUpdated: remoteChainConfigAccount.Base.OutboundRateLimit.LastUpdated,
						Enabled:     remoteChainConfigAccount.Base.OutboundRateLimit.Cfg.Enabled,
						Capacity:    remoteChainConfigAccount.Base.OutboundRateLimit.Cfg.Capacity,
						Rate:        remoteChainConfigAccount.Base.OutboundRateLimit.Cfg.Rate},
				}
				for i, addr := range remoteChainConfigAccount.Base.Remote.PoolAddresses {
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
