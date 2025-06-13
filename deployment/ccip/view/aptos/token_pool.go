package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	module_rate_limiter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool/rate_limiter"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type TokenPoolView struct {
	aptosCommon.ContractMetaData

	Token              string
	RemoteChainConfigs map[uint64]RemoteChainConfig
	AllowList          []string
	AllowListEnabled   bool
}

type RemoteChainConfig struct {
	RemoteTokenAddress        string
	RemotePoolAddresses       []string
	InboundRateLimiterConfig  RateLimiterConfig
	OutboundRateLimiterConfig RateLimiterConfig
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  uint64
	Rate      uint64
}

type PoolInterface interface {
	Owner(opts *bind.CallOpts) (aptos.AccountAddress, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	GetToken(opts *bind.CallOpts) (aptos.AccountAddress, error)
	GetAllowlistEnabled(opts *bind.CallOpts) (bool, error)
	GetAllowlist(opts *bind.CallOpts) ([]aptos.AccountAddress, error)
	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)
	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)
	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)
	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (module_rate_limiter.TokenBucket, error)
	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (module_rate_limiter.TokenBucket, error)
}

func GenerateTokenPoolView(chain cldf_aptos.Chain, address aptos.AccountAddress, boundTokenPoolModule PoolInterface) (TokenPoolView, error) {
	owner, err := boundTokenPoolModule.Owner(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get owner of token pool: %w", err)
	}
	typeAndVersion, err := boundTokenPoolModule.TypeAndVersion(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get type and version of token pool: %w", err)
	}
	token, err := boundTokenPoolModule.GetToken(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get token of token pool: %w", err)
	}
	allowlistEnabled, err := boundTokenPoolModule.GetAllowlistEnabled(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist of token pool: %w", err)
	}
	allowlist, err := boundTokenPoolModule.GetAllowlist(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist of token pool: %w", err)
	}
	allowListStrings := make([]string, len(allowlist))
	for i, address := range allowlist {
		allowListStrings[i] = address.StringLong()
	}
	remoteChains, err := boundTokenPoolModule.GetSupportedChains(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get remote chains of token pool: %w", err)
	}
	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(remoteChains))
	for _, selector := range remoteChains {
		remotePools, err := boundTokenPoolModule.GetRemotePools(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remote chains of token pool for chain %d: %w", selector, err)
		}
		remotePoolStrings := make([]string, len(remotePools))
		for i, remotePool := range remotePools {
			remotePoolStrings[i] = string(remotePool)
		}
		remoteToken, err := boundTokenPoolModule.GetRemoteToken(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remote token of token pool for chain %d: %w", selector, err)
		}
		inboundState, err := boundTokenPoolModule.GetCurrentInboundRateLimiterState(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get inbound rate limiter state of token pool for chain %d: %w", selector, err)
		}
		outboundState, err := boundTokenPoolModule.GetCurrentOutboundRateLimiterState(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get outbound rate limiter state of token pool for chain %d: %w", selector, err)
		}
		remoteChainConfigs[selector] = RemoteChainConfig{
			RemoteTokenAddress:  string(remoteToken),
			RemotePoolAddresses: remotePoolStrings,
			InboundRateLimiterConfig: RateLimiterConfig{
				IsEnabled: inboundState.IsEnabled,
				Capacity:  inboundState.Capacity,
				Rate:      inboundState.Rate,
			},
			OutboundRateLimiterConfig: RateLimiterConfig{
				IsEnabled: outboundState.IsEnabled,
				Capacity:  outboundState.Capacity,
				Rate:      outboundState.Rate,
			},
		}
	}

	return TokenPoolView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address:        address.StringLong(),
			Owner:          owner.StringLong(),
			TypeAndVersion: typeAndVersion,
		},
		Token:              token.StringLong(),
		RemoteChainConfigs: remoteChainConfigs,
		AllowList:          allowListStrings,
		AllowListEnabled:   allowlistEnabled,
	}, nil
}
