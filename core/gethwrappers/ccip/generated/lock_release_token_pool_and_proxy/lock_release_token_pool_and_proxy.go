// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lock_release_token_pool_and_proxy

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type PoolLockOrBurnInV1 struct {
	Receiver            []byte
	RemoteChainSelector uint64
	OriginalSender      common.Address
	Amount              *big.Int
	LocalToken          common.Address
}

type PoolLockOrBurnOutV1 struct {
	DestTokenAddress []byte
	DestPoolData     []byte
}

type PoolReleaseOrMintInV1 struct {
	OriginalSender      []byte
	RemoteChainSelector uint64
	Receiver            common.Address
	Amount              *big.Int
	LocalToken          common.Address
	SourcePoolAddress   []byte
	SourcePoolData      []byte
	OffchainTokenData   []byte
}

type PoolReleaseOrMintOutV1 struct {
	DestinationAmount *big.Int
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

type TokenPoolChainUpdate struct {
	RemoteChainSelector       uint64
	Allowed                   bool
	RemotePoolAddress         []byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var LockReleaseTokenPoolAndProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"acceptLiquidity\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LiquidityNotAccepted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"oldPool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"LegacyPoolChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canAcceptLiquidity\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"onRampAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRebalancer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"provideLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"prevPool\",\"type\":\"address\"}],\"name\":\"setPreviousPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebalancer\",\"type\":\"address\"}],\"name\":\"setRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004fbb38038062004fbb83398101604081905262000035916200056d565b84848483838383833380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c78162000186565b5050506001600160a01b0384161580620000e857506001600160a01b038116155b80620000fb57506001600160a01b038216155b156200011a576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200016d576040805160008152602081019091526200016d908462000231565b5050505094151560e05250620006de9650505050505050565b336001600160a01b03821603620001e05760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c05162000252576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002dd57600083828151811062000276576200027662000690565b60209081029190910101519050620002906002826200038e565b15620002d3576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000255565b5060005b81518110156200038957600082828151811062000302576200030262000690565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200032e575062000380565b6200033b600282620003ae565b156200037e576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002e1565b505050565b6000620003a5836001600160a01b038416620003c5565b90505b92915050565b6000620003a5836001600160a01b038416620004c9565b60008181526001830160205260408120548015620004be576000620003ec600183620006a6565b85549091506000906200040290600190620006a6565b90508181146200046e57600086600001828154811062000426576200042662000690565b90600052602060002001549050808760000184815481106200044c576200044c62000690565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004825762000482620006c8565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a8565b6000915050620003a8565b60008181526001830160205260408120546200051257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a8565b506000620003a8565b6001600160a01b03811681146200053157600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b805162000557816200051b565b919050565b805180151581146200055757600080fd5b600080600080600060a086880312156200058657600080fd5b855162000593816200051b565b602087810151919650906001600160401b0380821115620005b357600080fd5b818901915089601f830112620005c857600080fd5b815181811115620005dd57620005dd62000534565b8060051b604051601f19603f8301168101818110858211171562000605576200060562000534565b60405291825284820192508381018501918c8311156200062457600080fd5b938501935b828510156200064d576200063d856200054a565b8452938501939285019262000629565b80995050505050505062000664604087016200054a565b925062000674606087016200055c565b915062000684608087016200054a565b90509295509295909350565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003a857634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e0516148336200078860003960008181610545015261190b0152600081816105f201528181611ffa0152612bb60152600081816105cc01528181611cd801526122ad0152600081816102ad01528181610302015281816107d0015281816108a20152818161096c015281816119cd01528181611bf801528181611ef3015281816121cd015281816123b301528181612b4c0152612da101526148336000f3fe608060405234801561001057600080fd5b50600436106102415760003560e01c80638da5cb5b11610145578063c0d78655116100bd578063db6327dc1161008c578063e0351e1311610071578063e0351e13146105f0578063eb521a4c14610616578063f2fde38b1461062957600080fd5b8063db6327dc146105b7578063dc0bd971146105ca57600080fd5b8063c0d7865514610569578063c4bffe2b1461057c578063c75eea9c14610591578063cf7401f3146105a457600080fd5b8063a8d87a3b11610114578063b0f479a1116100f9578063b0f479a114610512578063b794658014610530578063bb98546b1461054357600080fd5b8063a8d87a3b14610490578063af58d59f146104a357600080fd5b80638da5cb5b1461042a5780639766b932146104485780639a4575b91461045b578063a7cd63b71461047b57600080fd5b806354c8a4f3116101d857806378a010b2116101a75780637d54534e1161018c5780637d54534e146103f157806383826b2b146104045780638926f54f1461041757600080fd5b806378a010b2146103d657806379ba5097146103e957600080fd5b806354c8a4f31461037f57806366320087146103925780636cfd1553146103a55780636d3d1a58146103b857600080fd5b806321df0da71161021457806321df0da7146102ab578063240028e8146102f2578063390775371461033f578063432a6ba31461036157600080fd5b806301ffc9a7146102465780630a2fd4931461026e5780630a861f2a1461028e578063181f5a77146102a3575b600080fd5b61025961025436600461374b565b61063c565b60405190151581526020015b60405180910390f35b61028161027c3660046137aa565b610698565b6040516102659190613833565b6102a161029c366004613846565b610748565b005b6102816108f9565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610265565b61025961030036600461388c565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b61035261034d3660046138a9565b610915565b60405190518152602001610265565b600a5473ffffffffffffffffffffffffffffffffffffffff166102cd565b6102a161038d366004613931565b610a3e565b6102a16103a036600461399d565b610ab9565b6102a16103b336600461388c565b610b45565b60085473ffffffffffffffffffffffffffffffffffffffff166102cd565b6102a16103e43660046139c9565b610b94565b6102a1610d03565b6102a16103ff36600461388c565b610e00565b610259610412366004613a4c565b610e4f565b6102596104253660046137aa565b610f1c565b60005473ffffffffffffffffffffffffffffffffffffffff166102cd565b6102a161045636600461388c565b610f33565b61046e610469366004613a83565b610fc2565b6040516102659190613abe565b61048361108b565b6040516102659190613b1e565b6102cd61049e3660046137aa565b503090565b6104b66104b13660046137aa565b61109c565b604051610265919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff166102cd565b61028161053e3660046137aa565b611171565b7f0000000000000000000000000000000000000000000000000000000000000000610259565b6102a161057736600461388c565b61119c565b610584611270565b6040516102659190613b78565b6104b661059f3660046137aa565b611328565b6102a16105b2366004613d2f565b6113fa565b6102a16105c5366004613d74565b611483565b7f00000000000000000000000000000000000000000000000000000000000000006102cd565b7f0000000000000000000000000000000000000000000000000000000000000000610259565b6102a1610624366004613846565b611909565b6102a161063736600461388c565b611a25565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fe1d40566000000000000000000000000000000000000000000000000000000001480610692575061069282611a39565b92915050565b67ffffffffffffffff811660009081526007602052604090206004018054606091906106c390613db6565b80601f01602080910402602001604051908101604052809291908181526020018280546106ef90613db6565b801561073c5780601f106107115761010080835404028352916020019161073c565b820191906000526020600020905b81548152906001019060200180831161071f57829003601f168201915b50505050509050919050565b600a5473ffffffffffffffffffffffffffffffffffffffff1633146107a0576040517f8e4a23d60000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015281907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561082c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108509190613e09565b1015610888576040517fbb55fd2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108c973ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163383611b1d565b604051819033907fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171990600090a350565b6040518060600160405280602681526020016148016026913981565b60408051602081019091526000815261093561093083613ebe565b611bf1565b60095473ffffffffffffffffffffffffffffffffffffffff1661099c5761099773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016336060850135611b1d565b6109ad565b6109ad6109a883613ebe565b611e22565b6109bd606083016040840161388c565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f528460600135604051610a1f91815260200190565b60405180910390a3506040805160208101909152606090910135815290565b610a46611f75565b610ab384848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611ff892505050565b50505050565b610ac1611f75565b6040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810182905273ffffffffffffffffffffffffffffffffffffffff831690630a861f2a90602401600060405180830381600087803b158015610b2957600080fd5b505af1158015610b3d573d6000803e3d6000fd5b505050505050565b610b4d611f75565b600a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610b9c611f75565b610ba583610f1c565b610be7576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610797565b67ffffffffffffffff831660009081526007602052604081206004018054610c0e90613db6565b80601f0160208091040260200160405190810160405280929190818152602001828054610c3a90613db6565b8015610c875780601f10610c5c57610100808354040283529160200191610c87565b820191906000526020600020905b815481529060010190602001808311610c6a57829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610cb6838583613ffb565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610cf593929190614116565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610d84576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610797565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610e08611f75565b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600073ffffffffffffffffffffffffffffffffffffffff8216301480610f155750600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169281019290925273ffffffffffffffffffffffffffffffffffffffff848116602484015216906383826b2b90604401602060405180830381865afa158015610ef1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f15919061417a565b9392505050565b6000610692600567ffffffffffffffff84166121ae565b610f3b611f75565b6009805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f22891015b60405180910390a15050565b6040805180820190915260608082526020820152610fe7610fe283614197565b6121c6565b60095473ffffffffffffffffffffffffffffffffffffffff16156110165761101661101183614197565b612390565b6040516060830135815233907f9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd600089060200160405180910390a2604051806040016040528061107084602001602081019061053e91906137aa565b81526040805160208181019092526000815291015292915050565b606061109760026124aa565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526003909101548084166060830152919091049091166080820152610692906124b7565b67ffffffffffffffff811660009081526007602052604090206005018054606091906106c390613db6565b6111a4611f75565b73ffffffffffffffffffffffffffffffffffffffff81166111f1576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849101610fb6565b6060600061127e60056124aa565b90506000815167ffffffffffffffff81111561129c5761129c613bba565b6040519080825280602002602001820160405280156112c5578160200160208202803683370190505b50905060005b8251811015611321578281815181106112e6576112e6614239565b602002602001015182828151811061130057611300614239565b67ffffffffffffffff909216602092830291909101909101526001016112cb565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526001909101548084166060830152919091049091166080820152610692906124b7565b60085473ffffffffffffffffffffffffffffffffffffffff16331480159061143a575060005473ffffffffffffffffffffffffffffffffffffffff163314155b15611473576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b61147e838383612569565b505050565b61148b611f75565b60005b8181101561147e5760008383838181106114aa576114aa614239565b90506020028101906114bc9190614268565b6114c5906142a6565b90506114da8160800151826020015115612653565b6114ed8160a00151826020015115612653565b8060200151156117e957805161150f9060059067ffffffffffffffff1661278c565b6115545780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610797565b60408101515115806115695750606081015151155b156115a0576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c17909116961515029590951790985590810151940151938116931690910291909117600382015591519091906004820190611781908261435a565b5060608201516005820190611796908261435a565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506117dc9493929190614474565b60405180910390a1611900565b80516118019060059067ffffffffffffffff16612798565b6118465780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610797565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906118af60048301826136fd565b6118bd6005830160006136fd565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b5060010161148e565b7f0000000000000000000000000000000000000000000000000000000000000000611960576040517fe93f8fa400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a5473ffffffffffffffffffffffffffffffffffffffff1633146119b3576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b6119f573ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163330846127a4565b604051819033907fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208890600090a350565b611a2d611f75565b611a3681612802565b50565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf000000000000000000000000000000000000000000000000000000001480611acc57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061069257507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b60405173ffffffffffffffffffffffffffffffffffffffff831660248201526044810182905261147e9084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526128f7565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611c865760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610797565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611d34573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d58919061417a565b15611d8f576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611d9c8160200151612a03565b6000611dab8260200151610698565b9050805160001480611dcf575080805190602001208260a001518051906020012014155b15611e0c578160a001516040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016107979190613833565b611e1e82602001518360600151612b29565b5050565b6009548151606083015160208401516040517f8627fad600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90941693638627fad693611e86939092309260040161450d565b600060405180830381600087803b158015611ea057600080fd5b505af1158015611eb4573d6000803e3d6000fd5b5050505060608101516040517f095ea7b300000000000000000000000000000000000000000000000000000000815233600482015260248101919091527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063095ea7b3906044016020604051808303816000875af1158015611f51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e1e919061417a565b60005473ffffffffffffffffffffffffffffffffffffffff163314611ff6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610797565b565b7f000000000000000000000000000000000000000000000000000000000000000061204f576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82518110156120e557600083828151811061206f5761206f614239565b6020026020010151905061208d816002612b7090919063ffffffff16565b156120dc5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101612052565b5060005b815181101561147e57600082828151811061210657612106614239565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361214a57506121a6565b612155600282612b92565b156121a45760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6001016120e9565b60008181526001830160205260408120541515610f15565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161461225b5760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610797565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015612309573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061232d919061417a565b15612364576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6123718160400151612bb4565b61237e8160200151612c33565b611a3681602001518260600151612d81565b60095460608201516123dd9173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811692911690611b1d565b60095460408083015183516060850151602086015193517f9687544500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909516946396875445946124459493929160040161456e565b6000604051808303816000875af1158015612464573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611e1e91908101906145ce565b60606000610f1583612dc5565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261254582606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642612529919061466b565b85608001516fffffffffffffffffffffffffffffffff16612e20565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61257283610f1c565b6125b4576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610797565b6125bf826000612653565b67ffffffffffffffff831660009081526007602052604090206125e29083612e4a565b6125ed816000612653565b67ffffffffffffffff831660009081526007602052604090206126139060020182612e4a565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516126469392919061467e565b60405180910390a1505050565b81511561271a5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff161015806126a9575060408201516fffffffffffffffffffffffffffffffff16155b156126e257816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016107979190614701565b8015611e1e576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff16151580612753575060208201516fffffffffffffffffffffffffffffffff1615155b15611e1e57816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016107979190614701565b6000610f158383612fec565b6000610f15838361303b565b60405173ffffffffffffffffffffffffffffffffffffffff80851660248301528316604482015260648101829052610ab39085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401611b6f565b3373ffffffffffffffffffffffffffffffffffffffff821603612881576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610797565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612959826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661312e9092919063ffffffff16565b80519091501561147e5780806020019051810190612977919061417a565b61147e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610797565b612a0c81610f1c565b612a4e576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610797565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015612acd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612af1919061417a565b611a36576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b67ffffffffffffffff82166000908152600760205260409020611e1e90600201827f000000000000000000000000000000000000000000000000000000000000000061313d565b6000610f158373ffffffffffffffffffffffffffffffffffffffff841661303b565b6000610f158373ffffffffffffffffffffffffffffffffffffffff8416612fec565b7f000000000000000000000000000000000000000000000000000000000000000015611a3657612be56002826134c0565b611a36576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610797565b612c3c81610f1c565b612c7e576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610797565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612cf7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d1b919061473d565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611a36576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b67ffffffffffffffff82166000908152600760205260409020611e1e90827f000000000000000000000000000000000000000000000000000000000000000061313d565b60608160000180548060200260200160405190810160405280929190818152602001828054801561073c57602002820191906000526020600020905b815481526020019060010190808311612e015750505050509050919050565b6000612e3f85612e30848661475a565b612e3a9087614771565b6134ef565b90505b949350505050565b8154600090612e7390700100000000000000000000000000000000900463ffffffff164261466b565b90508015612f155760018301548354612ebb916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612e20565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612f3b916fffffffffffffffffffffffffffffffff90811691166134ef565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612646908490614701565b600081815260018301602052604081205461303357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610692565b506000610692565b6000818152600183016020526040812054801561312457600061305f60018361466b565b85549091506000906130739060019061466b565b90508181146130d857600086600001828154811061309357613093614239565b90600052602060002001549050808760000184815481106130b6576130b6614239565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806130e9576130e9614784565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610692565b6000915050610692565b6060612e428484600085613505565b825474010000000000000000000000000000000000000000900460ff161580613164575081155b1561316e57505050565b825460018401546fffffffffffffffffffffffffffffffff808316929116906000906131b490700100000000000000000000000000000000900463ffffffff164261466b565b9050801561327457818311156131f6576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546132309083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612e20565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b8482101561332b5773ffffffffffffffffffffffffffffffffffffffff84166132d3576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610797565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610797565b8483101561343e5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1690600090829061336f908261466b565b613379878a61466b565b6133839190614771565b61338d91906147b3565b905073ffffffffffffffffffffffffffffffffffffffff86166133e6576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610797565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610797565b613448858461466b565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610f15565b60008183106134fe5781610f15565b5090919050565b606082471015613597576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610797565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516135c091906147ee565b60006040518083038185875af1925050503d80600081146135fd576040519150601f19603f3d011682016040523d82523d6000602084013e613602565b606091505b50915091506136138783838761361e565b979650505050505050565b606083156136b45782516000036136ad5773ffffffffffffffffffffffffffffffffffffffff85163b6136ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610797565b5081612e42565b612e4283838151156136c95781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107979190613833565b50805461370990613db6565b6000825580601f10613719575050565b601f016020900490600052602060002090810190611a3691905b808211156137475760008155600101613733565b5090565b60006020828403121561375d57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610f1557600080fd5b803567ffffffffffffffff811681146137a557600080fd5b919050565b6000602082840312156137bc57600080fd5b610f158261378d565b60005b838110156137e05781810151838201526020016137c8565b50506000910152565b600081518084526138018160208601602086016137c5565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610f1560208301846137e9565b60006020828403121561385857600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff81168114611a3657600080fd5b80356137a58161385f565b60006020828403121561389e57600080fd5b8135610f158161385f565b6000602082840312156138bb57600080fd5b813567ffffffffffffffff8111156138d257600080fd5b82016101008185031215610f1557600080fd5b60008083601f8401126138f757600080fd5b50813567ffffffffffffffff81111561390f57600080fd5b6020830191508360208260051b850101111561392a57600080fd5b9250929050565b6000806000806040858703121561394757600080fd5b843567ffffffffffffffff8082111561395f57600080fd5b61396b888389016138e5565b9096509450602087013591508082111561398457600080fd5b50613991878288016138e5565b95989497509550505050565b600080604083850312156139b057600080fd5b82356139bb8161385f565b946020939093013593505050565b6000806000604084860312156139de57600080fd5b6139e78461378d565b9250602084013567ffffffffffffffff80821115613a0457600080fd5b818601915086601f830112613a1857600080fd5b813581811115613a2757600080fd5b876020828501011115613a3957600080fd5b6020830194508093505050509250925092565b60008060408385031215613a5f57600080fd5b613a688361378d565b91506020830135613a788161385f565b809150509250929050565b600060208284031215613a9557600080fd5b813567ffffffffffffffff811115613aac57600080fd5b820160a08185031215610f1557600080fd5b602081526000825160406020840152613ada60608401826137e9565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152613b1582826137e9565b95945050505050565b6020808252825182820181905260009190848201906040850190845b81811015613b6c57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613b3a565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015613b6c57835167ffffffffffffffff1683529284019291840191600101613b94565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff81118282101715613c0d57613c0d613bba565b60405290565b60405160c0810167ffffffffffffffff81118282101715613c0d57613c0d613bba565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613c7d57613c7d613bba565b604052919050565b8015158114611a3657600080fd5b80356137a581613c85565b80356fffffffffffffffffffffffffffffffff811681146137a557600080fd5b600060608284031215613cd057600080fd5b6040516060810181811067ffffffffffffffff82111715613cf357613cf3613bba565b6040529050808235613d0481613c85565b8152613d1260208401613c9e565b6020820152613d2360408401613c9e565b60408201525092915050565b600080600060e08486031215613d4457600080fd5b613d4d8461378d565b9250613d5c8560208601613cbe565b9150613d6b8560808601613cbe565b90509250925092565b60008060208385031215613d8757600080fd5b823567ffffffffffffffff811115613d9e57600080fd5b613daa858286016138e5565b90969095509350505050565b600181811c90821680613dca57607f821691505b602082108103613e03577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208284031215613e1b57600080fd5b5051919050565b600067ffffffffffffffff821115613e3c57613e3c613bba565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613e7957600080fd5b8135613e8c613e8782613e22565b613c36565b818152846020838601011115613ea157600080fd5b816020850160208301376000918101602001919091529392505050565b60006101008236031215613ed157600080fd5b613ed9613be9565b823567ffffffffffffffff80821115613ef157600080fd5b613efd36838701613e68565b8352613f0b6020860161378d565b6020840152613f1c60408601613881565b604084015260608501356060840152613f3760808601613881565b608084015260a0850135915080821115613f5057600080fd5b613f5c36838701613e68565b60a084015260c0850135915080821115613f7557600080fd5b613f8136838701613e68565b60c084015260e0850135915080821115613f9a57600080fd5b50613fa736828601613e68565b60e08301525092915050565b601f82111561147e576000816000526020600020601f850160051c81016020861015613fdc5750805b601f850160051c820191505b81811015610b3d57828155600101613fe8565b67ffffffffffffffff83111561401357614013613bba565b614027836140218354613db6565b83613fb3565b6000601f84116001811461407957600085156140435750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561410f565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156140c857868501358255602094850194600190920191016140a8565b5086821015614103577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60408152600061412960408301866137e9565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b60006020828403121561418c57600080fd5b8151610f1581613c85565b600060a082360312156141a957600080fd5b60405160a0810167ffffffffffffffff82821081831117156141cd576141cd613bba565b8160405284359150808211156141e257600080fd5b506141ef36828601613e68565b8252506141fe6020840161378d565b602082015260408301356142118161385f565b604082015260608381013590820152608083013561422e8161385f565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec183360301811261429c57600080fd5b9190910192915050565b600061014082360312156142b957600080fd5b6142c1613c13565b6142ca8361378d565b81526142d860208401613c93565b6020820152604083013567ffffffffffffffff808211156142f857600080fd5b61430436838701613e68565b6040840152606085013591508082111561431d57600080fd5b5061432a36828601613e68565b60608301525061433d3660808501613cbe565b608082015261434f3660e08501613cbe565b60a082015292915050565b815167ffffffffffffffff81111561437457614374613bba565b614388816143828454613db6565b84613fb3565b602080601f8311600181146143db57600084156143a55750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610b3d565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561442857888601518255948401946001909101908401614409565b508582101561446457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152614498818401876137e9565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506144d69050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152613b15565b60a08152600061452060a08301876137e9565b73ffffffffffffffffffffffffffffffffffffffff8616602084015284604084015267ffffffffffffffff841660608401528281036080840152600081526020810191505095945050505050565b73ffffffffffffffffffffffffffffffffffffffff8516815260a06020820152600061459d60a08301866137e9565b60408301949094525067ffffffffffffffff9190911660608201528082036080909101526000815260200192915050565b6000602082840312156145e057600080fd5b815167ffffffffffffffff8111156145f757600080fd5b8201601f8101841361460857600080fd5b8051614616613e8782613e22565b81815285602083850101111561462b57600080fd5b613b158260208301602086016137c5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156106925761069261463c565b67ffffffffffffffff8416815260e081016146ca60208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152612e42565b6060810161069282848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561474f57600080fd5b8151610f158161385f565b80820281158282048414176106925761069261463c565b808201808211156106925761069261463c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000826147e9577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000825161429c8184602087016137c556fe4c6f636b52656c65617365546f6b656e506f6f6c416e6450726f787920312e352e302d646576a164736f6c6343000818000a",
}

var LockReleaseTokenPoolAndProxyABI = LockReleaseTokenPoolAndProxyMetaData.ABI

var LockReleaseTokenPoolAndProxyBin = LockReleaseTokenPoolAndProxyMetaData.Bin

func DeployLockReleaseTokenPoolAndProxy(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, acceptLiquidity bool, router common.Address) (common.Address, *types.Transaction, *LockReleaseTokenPoolAndProxy, error) {
	parsed, err := LockReleaseTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LockReleaseTokenPoolAndProxyBin), backend, token, allowlist, rmnProxy, acceptLiquidity, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockReleaseTokenPoolAndProxy{address: address, abi: *parsed, LockReleaseTokenPoolAndProxyCaller: LockReleaseTokenPoolAndProxyCaller{contract: contract}, LockReleaseTokenPoolAndProxyTransactor: LockReleaseTokenPoolAndProxyTransactor{contract: contract}, LockReleaseTokenPoolAndProxyFilterer: LockReleaseTokenPoolAndProxyFilterer{contract: contract}}, nil
}

type LockReleaseTokenPoolAndProxy struct {
	address common.Address
	abi     abi.ABI
	LockReleaseTokenPoolAndProxyCaller
	LockReleaseTokenPoolAndProxyTransactor
	LockReleaseTokenPoolAndProxyFilterer
}

type LockReleaseTokenPoolAndProxyCaller struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolAndProxyTransactor struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolAndProxyFilterer struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolAndProxySession struct {
	Contract     *LockReleaseTokenPoolAndProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolAndProxyCallerSession struct {
	Contract *LockReleaseTokenPoolAndProxyCaller
	CallOpts bind.CallOpts
}

type LockReleaseTokenPoolAndProxyTransactorSession struct {
	Contract     *LockReleaseTokenPoolAndProxyTransactor
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolAndProxyRaw struct {
	Contract *LockReleaseTokenPoolAndProxy
}

type LockReleaseTokenPoolAndProxyCallerRaw struct {
	Contract *LockReleaseTokenPoolAndProxyCaller
}

type LockReleaseTokenPoolAndProxyTransactorRaw struct {
	Contract *LockReleaseTokenPoolAndProxyTransactor
}

func NewLockReleaseTokenPoolAndProxy(address common.Address, backend bind.ContractBackend) (*LockReleaseTokenPoolAndProxy, error) {
	abi, err := abi.JSON(strings.NewReader(LockReleaseTokenPoolAndProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLockReleaseTokenPoolAndProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxy{address: address, abi: abi, LockReleaseTokenPoolAndProxyCaller: LockReleaseTokenPoolAndProxyCaller{contract: contract}, LockReleaseTokenPoolAndProxyTransactor: LockReleaseTokenPoolAndProxyTransactor{contract: contract}, LockReleaseTokenPoolAndProxyFilterer: LockReleaseTokenPoolAndProxyFilterer{contract: contract}}, nil
}

func NewLockReleaseTokenPoolAndProxyCaller(address common.Address, caller bind.ContractCaller) (*LockReleaseTokenPoolAndProxyCaller, error) {
	contract, err := bindLockReleaseTokenPoolAndProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyCaller{contract: contract}, nil
}

func NewLockReleaseTokenPoolAndProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*LockReleaseTokenPoolAndProxyTransactor, error) {
	contract, err := bindLockReleaseTokenPoolAndProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyTransactor{contract: contract}, nil
}

func NewLockReleaseTokenPoolAndProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*LockReleaseTokenPoolAndProxyFilterer, error) {
	contract, err := bindLockReleaseTokenPoolAndProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyFilterer{contract: contract}, nil
}

func bindLockReleaseTokenPoolAndProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LockReleaseTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPoolAndProxy.Contract.LockReleaseTokenPoolAndProxyCaller.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.LockReleaseTokenPoolAndProxyTransactor.contract.Transfer(opts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.LockReleaseTokenPoolAndProxyTransactor.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPoolAndProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.contract.Transfer(opts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) CanAcceptLiquidity(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "canAcceptLiquidity")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.CanAcceptLiquidity(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.CanAcceptLiquidity(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetAllowList(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetAllowList(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetAllowListEnabled(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetAllowListEnabled(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getOnRamp", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetOnRamp(&_LockReleaseTokenPoolAndProxy.CallOpts, arg0)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetOnRamp(&_LockReleaseTokenPoolAndProxy.CallOpts, arg0)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRebalancer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRebalancer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRebalancer(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRebalancer(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRemotePool(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRemotePool(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRemoteToken(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRemoteToken(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRmnProxy(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRmnProxy(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRouter(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetRouter(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetSupportedChains(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetSupportedChains(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetToken(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetToken(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "isOffRamp", sourceChainSelector, offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsOffRamp(&_LockReleaseTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsOffRamp(&_LockReleaseTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsSupportedChain(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsSupportedChain(&_LockReleaseTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsSupportedToken(&_LockReleaseTokenPoolAndProxy.CallOpts, token)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.IsSupportedToken(&_LockReleaseTokenPoolAndProxy.CallOpts, token)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) Owner() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.Owner(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) Owner() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.Owner(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SupportsInterface(&_LockReleaseTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SupportsInterface(&_LockReleaseTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TypeAndVersion(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TypeAndVersion(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "acceptOwnership")
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.AcceptOwnership(&_LockReleaseTokenPoolAndProxy.TransactOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.AcceptOwnership(&_LockReleaseTokenPoolAndProxy.TransactOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ApplyChainUpdates(&_LockReleaseTokenPoolAndProxy.TransactOpts, chains)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ApplyChainUpdates(&_LockReleaseTokenPoolAndProxy.TransactOpts, chains)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.LockOrBurn(&_LockReleaseTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.LockOrBurn(&_LockReleaseTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "provideLiquidity", amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ProvideLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ProvideLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ReleaseOrMint(&_LockReleaseTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.ReleaseOrMint(&_LockReleaseTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetPreviousPool(opts *bind.TransactOpts, prevPool common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setPreviousPool", prevPool)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetPreviousPool(&_LockReleaseTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetPreviousPool(&_LockReleaseTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_LockReleaseTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_LockReleaseTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setRebalancer", rebalancer)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRebalancer(&_LockReleaseTokenPoolAndProxy.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRebalancer(&_LockReleaseTokenPoolAndProxy.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRemotePool(&_LockReleaseTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRemotePool(&_LockReleaseTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "setRouter", newRouter)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRouter(&_LockReleaseTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.SetRouter(&_LockReleaseTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "transferLiquidity", from, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TransferLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, from, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TransferLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, from, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TransferOwnership(&_LockReleaseTokenPoolAndProxy.TransactOpts, to)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.TransferOwnership(&_LockReleaseTokenPoolAndProxy.TransactOpts, to)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactor) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.contract.Transact(opts, "withdrawLiquidity", amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.WithdrawLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, amount)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyTransactorSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.WithdrawLiquidity(&_LockReleaseTokenPoolAndProxy.TransactOpts, amount)
}

type LockReleaseTokenPoolAndProxyAllowListAddIterator struct {
	Event *LockReleaseTokenPoolAndProxyAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyAllowListAdd)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyAllowListAdd)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyAllowListAddIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyAllowListAddIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyAllowListAddIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyAllowListAdd)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAndProxyAllowListAdd, error) {
	event := new(LockReleaseTokenPoolAndProxyAllowListAdd)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyAllowListRemoveIterator struct {
	Event *LockReleaseTokenPoolAndProxyAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyAllowListRemove)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyAllowListRemove)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyAllowListRemoveIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyAllowListRemoveIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyAllowListRemove)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAndProxyAllowListRemove, error) {
	event := new(LockReleaseTokenPoolAndProxyAllowListRemove)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyBurnedIterator struct {
	Event *LockReleaseTokenPoolAndProxyBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyBurnedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolAndProxyBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyBurnedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyBurned)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseBurned(log types.Log) (*LockReleaseTokenPoolAndProxyBurned, error) {
	event := new(LockReleaseTokenPoolAndProxyBurned)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyChainAddedIterator struct {
	Event *LockReleaseTokenPoolAndProxyChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyChainAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyChainAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyChainAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainAddedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyChainAddedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainAdded) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyChainAdded)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseChainAdded(log types.Log) (*LockReleaseTokenPoolAndProxyChainAdded, error) {
	event := new(LockReleaseTokenPoolAndProxyChainAdded)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyChainConfiguredIterator struct {
	Event *LockReleaseTokenPoolAndProxyChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyChainConfigured)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyChainConfigured)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainConfiguredIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyChainConfiguredIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainConfigured) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyChainConfigured)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolAndProxyChainConfigured, error) {
	event := new(LockReleaseTokenPoolAndProxyChainConfigured)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyChainRemovedIterator struct {
	Event *LockReleaseTokenPoolAndProxyChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyChainRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyChainRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyChainRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainRemovedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyChainRemovedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainRemoved) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyChainRemoved)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolAndProxyChainRemoved, error) {
	event := new(LockReleaseTokenPoolAndProxyChainRemoved)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyConfigChangedIterator struct {
	Event *LockReleaseTokenPoolAndProxyConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyConfigChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyConfigChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyConfigChangedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyConfigChangedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyConfigChangedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyConfigChanged) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyConfigChanged)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolAndProxyConfigChanged, error) {
	event := new(LockReleaseTokenPoolAndProxyConfigChanged)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator struct {
	Event *LockReleaseTokenPoolAndProxyLegacyPoolChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyLegacyPoolChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyLegacyPoolChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyLegacyPoolChanged struct {
	OldPool common.Address
	NewPool common.Address
	Raw     types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterLegacyPoolChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "LegacyPoolChanged", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyLegacyPoolChanged)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseLegacyPoolChanged(log types.Log) (*LockReleaseTokenPoolAndProxyLegacyPoolChanged, error) {
	event := new(LockReleaseTokenPoolAndProxyLegacyPoolChanged)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyLiquidityAddedIterator struct {
	Event *LockReleaseTokenPoolAndProxyLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyLiquidityAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyLiquidityAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyLiquidityAdded struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolAndProxyLiquidityAddedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyLiquidityAddedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyLiquidityAdded)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolAndProxyLiquidityAdded, error) {
	event := new(LockReleaseTokenPoolAndProxyLiquidityAdded)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyLiquidityRemovedIterator struct {
	Event *LockReleaseTokenPoolAndProxyLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyLiquidityRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyLiquidityRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyLiquidityRemoved struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolAndProxyLiquidityRemovedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyLiquidityRemovedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyLiquidityRemoved)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolAndProxyLiquidityRemoved, error) {
	event := new(LockReleaseTokenPoolAndProxyLiquidityRemoved)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyLockedIterator struct {
	Event *LockReleaseTokenPoolAndProxyLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyLockedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolAndProxyLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyLockedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyLocked)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseLocked(log types.Log) (*LockReleaseTokenPoolAndProxyLocked, error) {
	event := new(LockReleaseTokenPoolAndProxyLocked)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyMintedIterator struct {
	Event *LockReleaseTokenPoolAndProxyMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyMintedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolAndProxyMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyMintedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyMinted)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseMinted(log types.Log) (*LockReleaseTokenPoolAndProxyMinted, error) {
	event := new(LockReleaseTokenPoolAndProxyMinted)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator struct {
	Event *LockReleaseTokenPoolAndProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyOwnershipTransferRequested)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolAndProxyOwnershipTransferRequested, error) {
	event := new(LockReleaseTokenPoolAndProxyOwnershipTransferRequested)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyOwnershipTransferredIterator struct {
	Event *LockReleaseTokenPoolAndProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolAndProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyOwnershipTransferredIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyOwnershipTransferred)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolAndProxyOwnershipTransferred, error) {
	event := new(LockReleaseTokenPoolAndProxyOwnershipTransferred)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyReleasedIterator struct {
	Event *LockReleaseTokenPoolAndProxyReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyReleasedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolAndProxyReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyReleasedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyReleased)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseReleased(log types.Log) (*LockReleaseTokenPoolAndProxyReleased, error) {
	event := new(LockReleaseTokenPoolAndProxyReleased)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyRemotePoolSetIterator struct {
	Event *LockReleaseTokenPoolAndProxyRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyRemotePoolSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyRemotePoolSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolAndProxyRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyRemotePoolSetIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyRemotePoolSet)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseRemotePoolSet(log types.Log) (*LockReleaseTokenPoolAndProxyRemotePoolSet, error) {
	event := new(LockReleaseTokenPoolAndProxyRemotePoolSet)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyRouterUpdatedIterator struct {
	Event *LockReleaseTokenPoolAndProxyRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyRouterUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyRouterUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyRouterUpdatedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyRouterUpdatedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyRouterUpdated)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolAndProxyRouterUpdated, error) {
	event := new(LockReleaseTokenPoolAndProxyRouterUpdated)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAndProxyTokensConsumedIterator struct {
	Event *LockReleaseTokenPoolAndProxyTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAndProxyTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAndProxyTokensConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(LockReleaseTokenPoolAndProxyTokensConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *LockReleaseTokenPoolAndProxyTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAndProxyTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAndProxyTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyTokensConsumedIterator, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAndProxyTokensConsumedIterator{contract: _LockReleaseTokenPoolAndProxy.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPoolAndProxy.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAndProxyTokensConsumed)
				if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyFilterer) ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolAndProxyTokensConsumed, error) {
	event := new(LockReleaseTokenPoolAndProxyTokensConsumed)
	if err := _LockReleaseTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LockReleaseTokenPoolAndProxy.abi.Events["AllowListAdd"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseAllowListAdd(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["AllowListRemove"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseAllowListRemove(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["Burned"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseBurned(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["ChainAdded"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseChainAdded(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["ChainConfigured"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseChainConfigured(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["ChainRemoved"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseChainRemoved(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["ConfigChanged"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseConfigChanged(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["LegacyPoolChanged"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseLegacyPoolChanged(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["LiquidityAdded"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseLiquidityAdded(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["LiquidityRemoved"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseLiquidityRemoved(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["Locked"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseLocked(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["Minted"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseMinted(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseOwnershipTransferRequested(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["OwnershipTransferred"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseOwnershipTransferred(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["Released"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseReleased(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["RemotePoolSet"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseRemotePoolSet(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["RouterUpdated"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseRouterUpdated(log)
	case _LockReleaseTokenPoolAndProxy.abi.Events["TokensConsumed"].ID:
		return _LockReleaseTokenPoolAndProxy.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LockReleaseTokenPoolAndProxyAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (LockReleaseTokenPoolAndProxyAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (LockReleaseTokenPoolAndProxyBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (LockReleaseTokenPoolAndProxyChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (LockReleaseTokenPoolAndProxyChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (LockReleaseTokenPoolAndProxyChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (LockReleaseTokenPoolAndProxyConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (LockReleaseTokenPoolAndProxyLegacyPoolChanged) Topic() common.Hash {
	return common.HexToHash("0x81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f228")
}

func (LockReleaseTokenPoolAndProxyLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (LockReleaseTokenPoolAndProxyLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (LockReleaseTokenPoolAndProxyLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (LockReleaseTokenPoolAndProxyMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (LockReleaseTokenPoolAndProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LockReleaseTokenPoolAndProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LockReleaseTokenPoolAndProxyReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (LockReleaseTokenPoolAndProxyRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (LockReleaseTokenPoolAndProxyRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (LockReleaseTokenPoolAndProxyTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxy) Address() common.Address {
	return _LockReleaseTokenPoolAndProxy.address
}

type LockReleaseTokenPoolAndProxyInterface interface {
	CanAcceptLiquidity(opts *bind.CallOpts) (bool, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

	GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

	IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetPreviousPool(opts *bind.TransactOpts, prevPool common.Address) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAndProxyAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAndProxyAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolAndProxyBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*LockReleaseTokenPoolAndProxyBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*LockReleaseTokenPoolAndProxyChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolAndProxyChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolAndProxyChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolAndProxyConfigChanged, error)

	FilterLegacyPoolChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyLegacyPoolChangedIterator, error)

	WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error)

	ParseLegacyPoolChanged(log types.Log) (*LockReleaseTokenPoolAndProxyLegacyPoolChanged, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolAndProxyLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolAndProxyLiquidityAdded, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolAndProxyLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolAndProxyLiquidityRemoved, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolAndProxyLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*LockReleaseTokenPoolAndProxyLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolAndProxyMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*LockReleaseTokenPoolAndProxyMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolAndProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolAndProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolAndProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolAndProxyOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolAndProxyReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*LockReleaseTokenPoolAndProxyReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolAndProxyRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*LockReleaseTokenPoolAndProxyRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolAndProxyRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolAndProxyTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAndProxyTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolAndProxyTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
