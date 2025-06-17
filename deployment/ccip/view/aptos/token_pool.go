package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	module_rate_limiter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool/rate_limiter"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type TokenPoolView struct {
	aptosCommon.ContractMetaData

	Token              string                       `json:"token"`
	RemoteChainConfigs map[uint64]RemoteChainConfig `json:"remoteChainConfigs"`
	AllowList          []string                     `json:"allowList"`
	AllowListEnabled   bool                         `json:"allowListEnabled"`
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
		return TokenPoolView{}, fmt.Errorf("failed to get owner of token pool %s: %w", address.StringLong(), err)
	}
	typeAndVersion, err := boundTokenPoolModule.TypeAndVersion(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get typeAndVersion of token pool %s: %w", address.StringLong(), err)
	}
	token, err := boundTokenPoolModule.GetToken(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get token of token pool %s: %w", address.StringLong(), err)
	}
	allowlistEnabled, err := boundTokenPoolModule.GetAllowlistEnabled(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist of token pool %s: %w", address.StringLong(), err)
	}
	allowlist, err := boundTokenPoolModule.GetAllowlist(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist of token pool %s: %w", address.StringLong(), err)
	}
	allowListStrings := make([]string, len(allowlist))
	for i, address := range allowlist {
		allowListStrings[i] = address.StringLong()
	}
	remoteChains, err := boundTokenPoolModule.GetSupportedChains(nil)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get supportedChains of token pool %s: %w", address.StringLong(), err)
	}
	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(remoteChains))
	for _, selector := range remoteChains {
		remotePools, err := boundTokenPoolModule.GetRemotePools(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remotePools of token pool %s for chain %d: %w", address.StringLong(), selector, err)
		}
		remotePoolStrings := make([]string, len(remotePools))
		for i, remotePool := range remotePools {
			remotePoolStrings[i] = shared.GetAddressFromBytes(selector, remotePool)
		}
		remoteToken, err := boundTokenPoolModule.GetRemoteToken(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remoteToken of token pool %s for chain %d: %w", address.StringLong(), selector, err)
		}
		inboundState, err := boundTokenPoolModule.GetCurrentInboundRateLimiterState(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get inboundRateLimiterState of token pool %s for chain %d: %w", address.StringLong(), selector, err)
		}
		outboundState, err := boundTokenPoolModule.GetCurrentOutboundRateLimiterState(nil, selector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get outboundRateLimiterState of token pool %s for chain %d: %w", address.StringLong(), selector, err)
		}
		remoteChainConfigs[selector] = RemoteChainConfig{
			RemoteTokenAddress:  shared.GetAddressFromBytes(selector, remoteToken),
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
