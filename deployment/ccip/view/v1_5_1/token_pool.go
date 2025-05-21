package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type TokenPoolContract interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
	GetToken(opts *bind.CallOpts) (common.Address, error)
	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)
	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)
	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)
	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)
}

func GetCurrentInboundRateLimiterState(t TokenPoolContract, remoteChainSelector uint64) (token_pool.RateLimiterTokenBucket, error) {
	switch v := t.(type) {
	case *burn_mint_token_pool.BurnMintTokenPool:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *burn_from_mint_token_pool.BurnFromMintTokenPool:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *burn_with_from_mint_token_pool.BurnWithFromMintTokenPool:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *lock_release_token_pool.LockReleaseTokenPool:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *usdc_token_pool.USDCTokenPool:
		state, err := v.GetCurrentInboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	default:
		return token_pool.RateLimiterTokenBucket{}, fmt.Errorf("unknown type %T", t)
	}
}

func GetCurrentOutboundRateLimiterState(t TokenPoolContract, remoteChainSelector uint64) (token_pool.RateLimiterTokenBucket, error) {
	switch v := t.(type) {
	case *burn_mint_token_pool.BurnMintTokenPool:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *burn_from_mint_token_pool.BurnFromMintTokenPool:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *burn_with_from_mint_token_pool.BurnWithFromMintTokenPool:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *lock_release_token_pool.LockReleaseTokenPool:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	case *usdc_token_pool.USDCTokenPool:
		state, err := v.GetCurrentOutboundRateLimiterState(nil, remoteChainSelector)
		return token_pool.RateLimiterTokenBucket(state), err
	default:
		return token_pool.RateLimiterTokenBucket{}, fmt.Errorf("unknown type %T", t)
	}
}

type RemoteChainConfig struct {
	// RemoteTokenAddress is the raw representation of the token address on the remote chain.
	RemoteTokenAddress string
	// RemotePoolAddresses are raw addresses of valid token pools on the remote chain.
	RemotePoolAddresses []string
	// InboundRateLimiterConfig is the rate limiter config for inbound transfers from the remote chain.
	InboundRateLimterConfig token_pool.RateLimiterConfig
	// OutboundRateLimiterConfig is the rate limiter config for outbound transfers to the remote chain.
	OutboundRateLimiterConfig token_pool.RateLimiterConfig
}

type PoolView struct {
	TokenPoolView
	LockReleaseTokenPoolView
	USDCTokenPoolView
}

type TokenPoolView struct {
	types.ContractMetaData
	Token              common.Address               `json:"token"`
	TokenPriceFeed     common.Address               `json:"tokenPriceFeed"`
	RemoteChainConfigs map[uint64]RemoteChainConfig `json:"remoteChainConfigs"`
	AllowList          []common.Address             `json:"allowList"`
	AllowListEnabled   bool                         `json:"allowListEnable"`
}

type USDCTokenPoolView struct {
	TokenMessenger     common.Address                                 `json:"tokenMessenger,omitempty"`
	MessageTransmitter common.Address                                 `json:"messageTransmitter,omitempty"`
	LocalDomain        uint32                                         `json:"localDomain,omitempty"`
	ChainToDomain      map[uint64]usdc_token_pool.USDCTokenPoolDomain `json:"chainToDomain,omitempty"`
}

type LockReleaseTokenPoolView struct {
	TokenPoolView
	AcceptLiquidity bool           `json:"acceptLiquidity,omitempty"`
	Rebalancer      common.Address `json:"rebalancer,omitempty"`
}

func GenerateTokenPoolView(pool TokenPoolContract, priceFeed common.Address) (TokenPoolView, error) {
	owner, err := pool.Owner(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	typeAndVersion, err := pool.TypeAndVersion(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	token, err := pool.GetToken(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	allowList, err := pool.GetAllowList(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	allowListEnabled, err := pool.GetAllowListEnabled(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	remoteChains, err := pool.GetSupportedChains(nil)
	if err != nil {
		return TokenPoolView{}, err
	}
	remoteChainConfigs := make(map[uint64]RemoteChainConfig)
	for _, remoteChain := range remoteChains {
		remotePools, err := pool.GetRemotePools(nil, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		remoteToken, err := pool.GetRemoteToken(nil, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		inboundState, err := GetCurrentInboundRateLimiterState(pool, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		outboundState, err := GetCurrentOutboundRateLimiterState(pool, remoteChain)
		if err != nil {
			return TokenPoolView{}, err
		}
		remoteChainConfigs[remoteChain] = RemoteChainConfig{
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
		for i, remotePool := range remotePools {
			remoteChainConfigs[remoteChain].RemotePoolAddresses[i] = shared.GetAddressFromBytes(remoteChain, remotePool)
		}
	}

	return TokenPoolView{
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

func GenerateLockReleaseTokenPoolView(pool *lock_release_token_pool.LockReleaseTokenPool, priceFeed common.Address) (PoolView, error) {
	basePoolView, err := GenerateTokenPoolView(pool, priceFeed)
	if err != nil {
		return PoolView{}, err
	}
	poolView := PoolView{
		TokenPoolView: basePoolView,
	}
	acceptLiquidity, err := pool.CanAcceptLiquidity(nil)
	if err != nil {
		return poolView, err
	}
	rebalancer, err := pool.GetRebalancer(nil)
	if err != nil {
		return poolView, err
	}
	poolView.LockReleaseTokenPoolView = LockReleaseTokenPoolView{
		TokenPoolView:   basePoolView,
		AcceptLiquidity: acceptLiquidity,
		Rebalancer:      rebalancer,
	}
	return poolView, nil
}

func GenerateUSDCTokenPoolView(pool *usdc_token_pool.USDCTokenPool) (PoolView, error) {
	basePoolView, err := GenerateTokenPoolView(pool, common.Address{})
	if err != nil {
		return PoolView{}, err
	}
	poolView := PoolView{
		TokenPoolView: basePoolView,
	}
	tokenMessenger, err := pool.ITokenMessenger(nil)
	if err != nil {
		return poolView, err
	}
	messageTransmitter, err := pool.IMessageTransmitter(nil)
	if err != nil {
		return poolView, err
	}
	localDomain, err := pool.ILocalDomainIdentifier(nil)
	if err != nil {
		return poolView, err
	}
	chainToDomain := make(map[uint64]usdc_token_pool.USDCTokenPoolDomain)
	remoteChains := maps.Keys(basePoolView.RemoteChainConfigs)
	for _, chainSel := range remoteChains {
		domain, err := pool.GetDomain(nil, chainSel)
		if err != nil {
			return poolView, err
		}
		chainToDomain[chainSel] = domain
	}
	poolView.USDCTokenPoolView = USDCTokenPoolView{
		TokenMessenger:     tokenMessenger,
		MessageTransmitter: messageTransmitter,
		LocalDomain:        localDomain,
		ChainToDomain:      chainToDomain,
	}
	return poolView, nil
}
