package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/burn_mint_token_pool_and_proxy"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	v1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type TokenPoolContract interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
	GetToken(opts *bind.CallOpts) (common.Address, error)
	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)
	GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)
	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)
	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)
}

func GetCurrentInboundRateLimiterState(t TokenPoolContract, remoteChainSelector uint64) (token_pool.RateLimiterTokenBucket, error) {
	switch v := t.(type) {
	case *burn_mint_token_pool_and_proxy.BurnMintTokenPoolAndProxy:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	default:
		return token_pool.RateLimiterTokenBucket{}, fmt.Errorf("unknown type %T", t)
	}
}

func GetCurrentOutboundRateLimiterState(t TokenPoolContract, remoteChainSelector uint64) (token_pool.RateLimiterTokenBucket, error) {
	switch v := t.(type) {
	case *burn_mint_token_pool_and_proxy.BurnMintTokenPoolAndProxy:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	default:
		return token_pool.RateLimiterTokenBucket{}, fmt.Errorf("unknown type %T", t)
	}
}

func GenerateTokenPoolView(pool TokenPoolContract, priceFeed common.Address) (v1_5_1.TokenPoolView, error) {
	owner, err := pool.Owner(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	typeAndVersion, err := pool.TypeAndVersion(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	token, err := pool.GetToken(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	allowList, err := pool.GetAllowList(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	allowListEnabled, err := pool.GetAllowListEnabled(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	remoteChains, err := pool.GetSupportedChains(nil)
	if err != nil {
		return v1_5_1.TokenPoolView{}, err
	}
	remoteChainConfigs := make(map[uint64]v1_5_1.RemoteChainConfig)
	for _, remoteChain := range remoteChains {
		remotePools, err := pool.GetRemotePool(nil, remoteChain)
		if err != nil {
			return v1_5_1.TokenPoolView{}, err
		}
		remoteToken, err := pool.GetRemoteToken(nil, remoteChain)
		if err != nil {
			return v1_5_1.TokenPoolView{}, err
		}
		inboundState, err := GetCurrentInboundRateLimiterState(pool, remoteChain)
		if err != nil {
			return v1_5_1.TokenPoolView{}, err
		}
		outboundState, err := GetCurrentOutboundRateLimiterState(pool, remoteChain)
		if err != nil {
			return v1_5_1.TokenPoolView{}, err
		}
		remoteChainConfigs[remoteChain] = v1_5_1.RemoteChainConfig{
			RemoteTokenAddress:  shared.GetAddressFromBytes(remoteChain, remoteToken),
			RemotePoolAddresses: make([]string, len(remotePools)),
			InboundRateLimterConfig: token_pool.RateLimiterConfig{
				IsEnabled: inboundState.IsEnabled,
				Capacity:  inboundState.Capacity,
				Rate:      inboundState.Rate,
			},
			OutboundRateLimiterConfig: token_pool.RateLimiterConfig{
				IsEnabled: outboundState.IsEnabled,
				Capacity:  outboundState.Capacity,
				Rate:      outboundState.Rate,
			},
		}

		remoteChainConfigs[remoteChain].RemotePoolAddresses[0] = shared.GetAddressFromBytes(remoteChain, remotePools)
	}

	return v1_5_1.TokenPoolView{
		ContractMetaData: types.ContractMetaData{
			TypeAndVersion: typeAndVersion,
			Address:        pool.Address(),
			Owner:          owner,
		},
		Token:              token,
		TokenPriceFeed:     priceFeed,
		RemoteChainConfigs: remoteChainConfigs,
		AllowList:          allowList,
		AllowListEnabled:   allowListEnabled,
	}, nil
}
