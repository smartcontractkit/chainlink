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
}

type PoolLockOrBurnOutV1 struct {
	DestPoolAddress []byte
	DestPoolData    []byte
}

type PoolReleaseOrMintInV1 struct {
	OriginalSender      []byte
	RemoteChainSelector uint64
	Receiver            common.Address
	Amount              *big.Int
	SourcePoolAddress   []byte
	SourcePoolData      []byte
	OffchainTokenData   []byte
}

type PoolReleaseOrMintOutV1 struct {
	LocalToken        common.Address
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
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var LockReleaseTokenPoolAndProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"acceptLiquidity\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRatelimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LiquidityNotAccepted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"oldPool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"LegacyPoolChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canAcceptLiquidity\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getArmProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockReleaseInterfaceId\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"onRampAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRebalancer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destPoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"provideLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"prevPool\",\"type\":\"address\"}],\"name\":\"setPreviousPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebalancer\",\"type\":\"address\"}],\"name\":\"setRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004ccc38038062004ccc83398101604081905262000035916200056e565b84848483838383833380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c78162000173565b5050506001600160a01b0384161580620000e857506001600160a01b038116155b1562000107576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200015a576040805160008152602081019091526200015a90846200021e565b5050505094151560e05250620007059650505050505050565b336001600160a01b03821603620001cd5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c0516200023f576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002d457600083828151811062000263576200026362000691565b602090810291909101015190506200027d6002826200038f565b15620002c0576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50620002cc81620006bd565b905062000242565b5060005b81518110156200038a576000828281518110620002f957620002f962000691565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000325575062000377565b62000332600282620003af565b1562000375576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6200038281620006bd565b9050620002d8565b505050565b6000620003a6836001600160a01b038416620003c6565b90505b92915050565b6000620003a6836001600160a01b038416620004ca565b60008181526001830160205260408120548015620004bf576000620003ed600183620006d9565b85549091506000906200040390600190620006d9565b90508181146200046f57600086600001828154811062000427576200042762000691565b90600052602060002001549050808760000184815481106200044d576200044d62000691565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004835762000483620006ef565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a9565b6000915050620003a9565b60008181526001830160205260408120546200051357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a9565b506000620003a9565b6001600160a01b03811681146200053257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b805162000558816200051c565b919050565b805180151581146200055857600080fd5b600080600080600060a086880312156200058757600080fd5b855162000594816200051c565b602087810151919650906001600160401b0380821115620005b457600080fd5b818901915089601f830112620005c957600080fd5b815181811115620005de57620005de62000535565b8060051b604051601f19603f8301168101818110858211171562000606576200060662000535565b60405291825284820192508381018501918c8311156200062557600080fd5b938501935b828510156200064e576200063e856200054b565b845293850193928501926200062a565b80995050505050505062000665604087016200054b565b925062000675606087016200055d565b915062000685608087016200054b565b90509295509295909350565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201620006d257620006d2620006a7565b5060010190565b81810381811115620003a957620003a9620006a7565b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e0516145326200079a600039600081816105120152611661015260008181610586015281816120a301526122690152600081816103200152818161091e0152610f13015260008181610277015281816107c30152818161089501528181610ab801528181610b950152818161172301528181611f5c0152818161248301526124ca01526145326000f3fe608060405234801561001057600080fd5b506004361061020b5760003560e01c80638926f54f1161012a578063bb98546b116100bd578063cf7401f31161008c578063eb521a4c11610071578063eb521a4c146105aa578063f2fde38b146105bd578063f6e2145e146105d057600080fd5b8063cf7401f314610571578063e0351e131461058457600080fd5b8063bb98546b14610510578063c0d7865514610536578063c4bffe2b14610549578063c75eea9c1461055e57600080fd5b8063a7cd63b7116100f9578063a7cd63b71461045b578063a8d87a3b14610470578063af58d59f14610483578063b0f479a1146104f257600080fd5b80638926f54f146103e95780638bfca18c146103fc5780638da5cb5b1461042a5780639766b9321461044857600080fd5b806354c8a4f3116101a257806379ba50971161017157806379ba50971461039b5780637a5c972d146103a35780637d54534e146103c357806383826b2b146103d657600080fd5b806354c8a4f3146103445780636cfd1553146103575780636d3d1a581461036a57806378a010b21461038857600080fd5b806321df0da7116101de57806321df0da7146102755780634059f55b146102bc578063432a6ba3146103005780635246492f1461031e57600080fd5b806301ffc9a7146102105780630a2fd493146102385780630a861f2a14610258578063181f5a771461026d575b600080fd5b61022361021e3660046134cd565b6105e3565b60405190151581526020015b60405180910390f35b61024b610246366004613527565b61068b565b60405161022f91906135b0565b61026b6102663660046135c3565b61073b565b005b61024b6108ec565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161022f565b6102cf6102ca3660046135dc565b610908565b60408051825173ffffffffffffffffffffffffffffffffffffffff168152602092830151928101929092520161022f565b60095473ffffffffffffffffffffffffffffffffffffffff16610297565b7f0000000000000000000000000000000000000000000000000000000000000000610297565b61026b610352366004613663565b610bc7565b61026b6103653660046136fc565b610c42565b600a5473ffffffffffffffffffffffffffffffffffffffff16610297565b61026b610396366004613719565b610c91565b61026b610e00565b6103b66103b136600461379c565b610efd565b60405161022f91906137d7565b61026b6103d13660046136fc565b6110c8565b6102236103e4366004613837565b611117565b6102236103f7366004613527565b6111e4565b6040517f98a4717700000000000000000000000000000000000000000000000000000000815260200161022f565b60005473ffffffffffffffffffffffffffffffffffffffff16610297565b61026b6104563660046136fc565b6111fb565b61046361128a565b60405161022f919061386e565b61029761047e366004613527565b503090565b610496610491366004613527565b61129b565b60405161022f919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610297565b7f0000000000000000000000000000000000000000000000000000000000000000610223565b61026b6105443660046136fc565b611370565b610551611444565b60405161022f91906138c8565b61049661056c366004613527565b611504565b61026b61057f366004613a50565b6115d6565b7f0000000000000000000000000000000000000000000000000000000000000000610223565b61026b6105b83660046135c3565b61165f565b61026b6105cb3660046136fc565b61177b565b61026b6105de366004613a95565b61178f565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f98a4717700000000000000000000000000000000000000000000000000000000148061067657507fffffffff0000000000000000000000000000000000000000000000000000000082167fe1d4056600000000000000000000000000000000000000000000000000000000145b80610685575061068582611c06565b92915050565b67ffffffffffffffff811660009081526007602052604090206004018054606091906106b690613ad7565b80601f01602080910402602001604051908101604052809291908181526020018280546106e290613ad7565b801561072f5780601f106107045761010080835404028352916020019161072f565b820191906000526020600020905b81548152906001019060200180831161071257829003601f168201915b50505050509050919050565b60095473ffffffffffffffffffffffffffffffffffffffff163314610793576040517f8e4a23d60000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015281907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561081f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108439190613b2a565b101561087b576040517fbb55fd2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108bc73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163383611cea565b604051819033907fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171990600090a350565b6040518060600160405280602681526020016145006026913981565b60408051808201909152600080825260208201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610987573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109ab9190613b43565b156109e2576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109fa6109f56040840160208501613527565b611dbe565b610a54610a0d6040840160208501613527565b610a1a6080850185613b60565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611ee492505050565b610a71610a676040840160208501613527565b8360600135611f39565b60085473ffffffffffffffffffffffffffffffffffffffff16610ae857610ae3610aa160608401604085016136fc565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906060850135611cea565b610af9565b610af9610af483613c61565b611f80565b610b0960608301604084016136fc565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f528460600135604051610b6b91815260200190565b60405180910390a3506040805180820190915273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152606082013560208201525b919050565b610bcf61201e565b610c3c848480806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250506040805160208088028281018201909352878252909350879250869182918501908490808284376000920191909152506120a192505050565b50505050565b610c4a61201e565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610c9961201e565b610ca2836111e4565b610ce4576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161078a565b67ffffffffffffffff831660009081526007602052604081206004018054610d0b90613ad7565b80601f0160208091040260200160405190810160405280929190818152602001828054610d3790613ad7565b8015610d845780601f10610d5957610100808354040283529160200191610d84565b820191906000526020600020905b815481529060010190602001808311610d6757829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610db3838583613d92565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610df293929190613eac565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610e81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161078a565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60408051808201909152606080825260208201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f7c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fa09190613b43565b15610fd7576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610fef610fea60608401604085016136fc565b612267565b6110076110026040840160208501613527565b612315565b61102461101a6040840160208501613527565b8360600135612463565b60085473ffffffffffffffffffffffffffffffffffffffff16156110535761105361104e83613f10565b6124a7565b6040516060830135815233907f9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd600089060200160405180910390a260405180604001604052806110ad8460200160208101906102469190613527565b81526040805160208082019092526000815291015292915050565b6110d061201e565b600a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600073ffffffffffffffffffffffffffffffffffffffff82163014806111dd5750600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169281019290925273ffffffffffffffffffffffffffffffffffffffff848116602484015216906383826b2b90604401602060405180830381865afa1580156111b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111dd9190613b43565b9392505050565b6000610685600567ffffffffffffffff84166125c1565b61120361201e565b6008805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f22891015b60405180910390a15050565b606061129660026125d9565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526003909101548084166060830152919091049091166080820152610685906125e6565b61137861201e565b73ffffffffffffffffffffffffffffffffffffffff81166113c5576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910161127e565b6060600061145260056125d9565b90506000815167ffffffffffffffff8111156114705761147061390a565b604051908082528060200260200182016040528015611499578160200160208202803683370190505b50905060005b82518110156114fd578281815181106114ba576114ba613fa0565b60200260200101518282815181106114d4576114d4613fa0565b67ffffffffffffffff909216602092830291909101909101526114f681613ffe565b905061149f565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526001909101548084166060830152919091049091166080820152610685906125e6565b600a5473ffffffffffffffffffffffffffffffffffffffff163314801590611616575060005473ffffffffffffffffffffffffffffffffffffffff163314155b1561164f576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161078a565b61165a838383612698565b505050565b7f00000000000000000000000000000000000000000000000000000000000000006116b6576040517fe93f8fa400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60095473ffffffffffffffffffffffffffffffffffffffff163314611709576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161078a565b61174b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612782565b604051819033907fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208890600090a350565b61178361201e565b61178c816127e0565b50565b61179761201e565b60005b8181101561165a5760008383838181106117b6576117b6613fa0565b90506020028101906117c89190614036565b6117d190614074565b90506117e681606001518260200151156128d5565b6117f981608001518260200151156128d5565b806020015115611ac457805161181b9060059067ffffffffffffffff16612a0e565b6118605780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161078a565b6040805161012081018252606083810180516020908101516fffffffffffffffffffffffffffffffff90811660808087019182524263ffffffff90811660a0808a01829052875151151560c08b01528751870151861660e08b015296518a015185166101008a015292885288519586018952818a01805186015185168752868601939093528251511515868a01528251850151841686880152915188015183168582015283870194855288880151878901908152848a0151151587890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009189168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff908116740100000000000000000000000000000000000000009415158502178655848d015194880151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9097168a029c909c17909116961515029590951790985594810151940151938116931690910291909117600382015591519091906004820190611a409082614121565b5060609182015160059190910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905581519082015160808301516040517f0f135cbb9afa12a8bf3bbd071c117bcca4ddeca6160ef7f33d012a81b9c0c47193611ab7939092909161423b565b60405180910390a1611bf5565b8051611adc9060059067ffffffffffffffff16612a1a565b611b215780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161078a565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff00000000000000000000000000000000000000000090811682556001820183905560028201805490911690556003810182905590611b8a600483018261347f565b5060050180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b50611bff81613ffe565b905061179a565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf000000000000000000000000000000000000000000000000000000001480611c9957507fffffffff0000000000000000000000000000000000000000000000000000000082167fb323973900000000000000000000000000000000000000000000000000000000145b8061068557507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b60405173ffffffffffffffffffffffffffffffffffffffff831660248201526044810182905261165a9084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152612a26565b611dc7816111e4565b611e09576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161078a565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015611e88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611eac9190613b43565b61178c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161078a565b611eed8261068b565b80519060200120818051906020012014611f3557806040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161078a91906135b0565b5050565b67ffffffffffffffff82166000908152600760205260409020611f3590600201827f0000000000000000000000000000000000000000000000000000000000000000612b32565b60085481516040808401516060850151602086015192517f8627fad600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90951694638627fad694611fe994909392916004016142be565b600060405180830381600087803b15801561200357600080fd5b505af1158015612017573d6000803e3d6000fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461209f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161078a565b565b7f00000000000000000000000000000000000000000000000000000000000000006120f8576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b825181101561219657600083828151811061211857612118613fa0565b60200260200101519050612136816002612eb590919063ffffffff16565b156121855760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5061218f81613ffe565b90506120fb565b5060005b815181101561165a5760008282815181106121b7576121b7613fa0565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036121fb5750612257565b612206600282612ed7565b156122555760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b61226081613ffe565b905061219a565b7f000000000000000000000000000000000000000000000000000000000000000080156122c657506122c460028273ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156111dd565b155b1561178c576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161078a565b61231e816111e4565b612360576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161078a565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa1580156123d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123fd919061431f565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461178c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161078a565b67ffffffffffffffff82166000908152600760205260409020611f3590827f0000000000000000000000000000000000000000000000000000000000000000612b32565b60085460608201516124f49173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811692911690611cea565b60085460408083015183516060850151602086015193517f9687544500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9095169463968754459461255c9493929160040161433c565b6000604051808303816000875af115801561257b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611f35919081019061439c565b600081815260018301602052604081205415156111dd565b606060006111dd83612ef9565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261267482606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642612658919061440a565b85608001516fffffffffffffffffffffffffffffffff16612f54565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b6126a1836111e4565b6126e3576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161078a565b6126ee8260006128d5565b67ffffffffffffffff831660009081526007602052604090206127119083612f7e565b61271c8160006128d5565b67ffffffffffffffff831660009081526007602052604090206127429060020182612f7e565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516127759392919061423b565b60405180910390a1505050565b60405173ffffffffffffffffffffffffffffffffffffffff80851660248301528316604482015260648101829052610c3c9085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401611d3c565b3373ffffffffffffffffffffffffffffffffffffffff82160361285f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161078a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b81511561299c5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff1610158061292b575060408201516fffffffffffffffffffffffffffffffff16155b1561296457816040517f70505e5600000000000000000000000000000000000000000000000000000000815260040161078a919061441d565b8015611f35576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff161515806129d5575060208201516fffffffffffffffffffffffffffffffff1615155b15611f3557816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161078a919061441d565b60006111dd8383613120565b60006111dd838361316f565b6000612a88826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166132629092919063ffffffff16565b80519091501561165a5780806020019051810190612aa69190613b43565b61165a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161078a565b825474010000000000000000000000000000000000000000900460ff161580612b59575081155b15612b6357505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612ba990700100000000000000000000000000000000900463ffffffff164261440a565b90508015612c695781831115612beb576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612c259083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612f54565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612d205773ffffffffffffffffffffffffffffffffffffffff8416612cc8576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161078a565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161078a565b84831015612e335760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612d64908261440a565b612d6e878a61440a565b612d789190614459565b612d82919061446c565b905073ffffffffffffffffffffffffffffffffffffffff8616612ddb576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161078a565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161078a565b612e3d858461440a565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b60006111dd8373ffffffffffffffffffffffffffffffffffffffff841661316f565b60006111dd8373ffffffffffffffffffffffffffffffffffffffff8416613120565b60608160000180548060200260200160405190810160405280929190818152602001828054801561072f57602002820191906000526020600020905b815481526020019060010190808311612f355750505050509050919050565b6000612f7385612f6484866144a7565b612f6e9087614459565b613271565b90505b949350505050565b8154600090612fa790700100000000000000000000000000000000900463ffffffff164261440a565b905080156130495760018301548354612fef916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612f54565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b6020820151835461306f916fffffffffffffffffffffffffffffffff9081169116613271565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199061277590849061441d565b600081815260018301602052604081205461316757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610685565b506000610685565b6000818152600183016020526040812054801561325857600061319360018361440a565b85549091506000906131a79060019061440a565b905081811461320c5760008660000182815481106131c7576131c7613fa0565b90600052602060002001549050808760000184815481106131ea576131ea613fa0565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061321d5761321d6144be565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610685565b6000915050610685565b6060612f768484600085613287565b600081831061328057816111dd565b5090919050565b606082471015613319576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c0000000000000000000000000000000000000000000000000000606482015260840161078a565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161334291906144ed565b60006040518083038185875af1925050503d806000811461337f576040519150601f19603f3d011682016040523d82523d6000602084013e613384565b606091505b5091509150613395878383876133a0565b979650505050505050565b6060831561343657825160000361342f5773ffffffffffffffffffffffffffffffffffffffff85163b61342f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161078a565b5081612f76565b612f76838381511561344b5781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078a91906135b0565b50805461348b90613ad7565b6000825580601f1061349b575050565b601f01602090049060005260206000209081019061178c91905b808211156134c957600081556001016134b5565b5090565b6000602082840312156134df57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146111dd57600080fd5b803567ffffffffffffffff81168114610bc257600080fd5b60006020828403121561353957600080fd5b6111dd8261350f565b60005b8381101561355d578181015183820152602001613545565b50506000910152565b6000815180845261357e816020860160208601613542565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006111dd6020830184613566565b6000602082840312156135d557600080fd5b5035919050565b6000602082840312156135ee57600080fd5b813567ffffffffffffffff81111561360557600080fd5b820160e081850312156111dd57600080fd5b60008083601f84011261362957600080fd5b50813567ffffffffffffffff81111561364157600080fd5b6020830191508360208260051b850101111561365c57600080fd5b9250929050565b6000806000806040858703121561367957600080fd5b843567ffffffffffffffff8082111561369157600080fd5b61369d88838901613617565b909650945060208701359150808211156136b657600080fd5b506136c387828801613617565b95989497509550505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461178c57600080fd5b8035610bc2816136cf565b60006020828403121561370e57600080fd5b81356111dd816136cf565b60008060006040848603121561372e57600080fd5b6137378461350f565b9250602084013567ffffffffffffffff8082111561375457600080fd5b818601915086601f83011261376857600080fd5b81358181111561377757600080fd5b87602082850101111561378957600080fd5b6020830194508093505050509250925092565b6000602082840312156137ae57600080fd5b813567ffffffffffffffff8111156137c557600080fd5b8201608081850312156111dd57600080fd5b6020815260008251604060208401526137f36060840182613566565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084830301604085015261382e8282613566565b95945050505050565b6000806040838503121561384a57600080fd5b6138538361350f565b91506020830135613863816136cf565b809150509250929050565b6020808252825182820181905260009190848201906040850190845b818110156138bc57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161388a565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156138bc57835167ffffffffffffffff16835292840192918401916001016138e4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171561395c5761395c61390a565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156139a9576139a961390a565b604052919050565b801515811461178c57600080fd5b80356fffffffffffffffffffffffffffffffff81168114610bc257600080fd5b6000606082840312156139f157600080fd5b6040516060810181811067ffffffffffffffff82111715613a1457613a1461390a565b6040529050808235613a25816139b1565b8152613a33602084016139bf565b6020820152613a44604084016139bf565b60408201525092915050565b600080600060e08486031215613a6557600080fd5b613a6e8461350f565b9250613a7d85602086016139df565b9150613a8c85608086016139df565b90509250925092565b60008060208385031215613aa857600080fd5b823567ffffffffffffffff811115613abf57600080fd5b613acb85828601613617565b90969095509350505050565b600181811c90821680613aeb57607f821691505b602082108103613b24577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208284031215613b3c57600080fd5b5051919050565b600060208284031215613b5557600080fd5b81516111dd816139b1565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112613b9557600080fd5b83018035915067ffffffffffffffff821115613bb057600080fd5b60200191503681900382131561365c57600080fd5b600067ffffffffffffffff821115613bdf57613bdf61390a565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613c1c57600080fd5b8135613c2f613c2a82613bc5565b613962565b818152846020838601011115613c4457600080fd5b816020850160208301376000918101602001919091529392505050565b600060e08236031215613c7357600080fd5b613c7b613939565b823567ffffffffffffffff80821115613c9357600080fd5b613c9f36838701613c0b565b8352613cad6020860161350f565b6020840152613cbe604086016136f1565b6040840152606085013560608401526080850135915080821115613ce157600080fd5b613ced36838701613c0b565b608084015260a0850135915080821115613d0657600080fd5b613d1236838701613c0b565b60a084015260c0850135915080821115613d2b57600080fd5b50613d3836828601613c0b565b60c08301525092915050565b601f82111561165a57600081815260208120601f850160051c81016020861015613d6b5750805b601f850160051c820191505b81811015613d8a57828155600101613d77565b505050505050565b67ffffffffffffffff831115613daa57613daa61390a565b613dbe83613db88354613ad7565b83613d44565b6000601f841160018114613e105760008515613dda5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355612017565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015613e5f5786850135825560209485019460019092019101613e3f565b5086821015613e9a577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b604081526000613ebf6040830186613566565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060808236031215613f2257600080fd5b6040516080810167ffffffffffffffff8282108183111715613f4657613f4661390a565b816040528435915080821115613f5b57600080fd5b50613f6836828601613c0b565b825250613f776020840161350f565b60208201526040830135613f8a816136cf565b6040820152606092830135928101929092525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361402f5761402f613fcf565b5060010190565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261406a57600080fd5b9190910192915050565b6000610120823603121561408757600080fd5b60405160a0810167ffffffffffffffff82821081831117156140ab576140ab61390a565b816040526140b88561350f565b8352602085013591506140ca826139b1565b81602084015260408501359150808211156140e457600080fd5b506140f136828601613c0b565b60408301525061410436606085016139df565b60608201526141163660c085016139df565b608082015292915050565b815167ffffffffffffffff81111561413b5761413b61390a565b61414f816141498454613ad7565b84613d44565b602080601f8311600181146141a2576000841561416c5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613d8a565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156141ef578886015182559484019460019091019084016141d0565b508582101561422b57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b67ffffffffffffffff8416815260e0810161428760208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152612f76565b60a0815260006142d160a0830187613566565b73ffffffffffffffffffffffffffffffffffffffff8616602084015284604084015267ffffffffffffffff841660608401528281036080840152600081526020810191505095945050505050565b60006020828403121561433157600080fd5b81516111dd816136cf565b73ffffffffffffffffffffffffffffffffffffffff8516815260a06020820152600061436b60a0830186613566565b60408301949094525067ffffffffffffffff9190911660608201528082036080909101526000815260200192915050565b6000602082840312156143ae57600080fd5b815167ffffffffffffffff8111156143c557600080fd5b8201601f810184136143d657600080fd5b80516143e4613c2a82613bc5565b8181528560208385010111156143f957600080fd5b61382e826020830160208601613542565b8181038181111561068557610685613fcf565b6060810161068582848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b8082018082111561068557610685613fcf565b6000826144a2577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b808202811582820484141761068557610685613fcf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000825161406a81846020870161354256fe4c6f636b52656c65617365546f6b656e506f6f6c416e6450726f787920312e352e302d646576a164736f6c6343000813000a",
}

var LockReleaseTokenPoolAndProxyABI = LockReleaseTokenPoolAndProxyMetaData.ABI

var LockReleaseTokenPoolAndProxyBin = LockReleaseTokenPoolAndProxyMetaData.Bin

func DeployLockReleaseTokenPoolAndProxy(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, armProxy common.Address, acceptLiquidity bool, router common.Address) (common.Address, *types.Transaction, *LockReleaseTokenPoolAndProxy, error) {
	parsed, err := LockReleaseTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LockReleaseTokenPoolAndProxyBin), backend, token, allowlist, armProxy, acceptLiquidity, router)
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

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetArmProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getArmProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetArmProxy() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetArmProxy(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetArmProxy() (common.Address, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetArmProxy(&_LockReleaseTokenPoolAndProxy.CallOpts)
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

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCaller) GetLockReleaseInterfaceId(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPoolAndProxy.contract.Call(opts, &out, "getLockReleaseInterfaceId")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxySession) GetLockReleaseInterfaceId() ([4]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetLockReleaseInterfaceId(&_LockReleaseTokenPoolAndProxy.CallOpts)
}

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxyCallerSession) GetLockReleaseInterfaceId() ([4]byte, error) {
	return _LockReleaseTokenPoolAndProxy.Contract.GetLockReleaseInterfaceId(&_LockReleaseTokenPoolAndProxy.CallOpts)
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
	return common.HexToHash("0x0f135cbb9afa12a8bf3bbd071c117bcca4ddeca6160ef7f33d012a81b9c0c471")
}

func (LockReleaseTokenPoolAndProxyChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (LockReleaseTokenPoolAndProxyChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
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

func (_LockReleaseTokenPoolAndProxy *LockReleaseTokenPoolAndProxy) Address() common.Address {
	return _LockReleaseTokenPoolAndProxy.address
}

type LockReleaseTokenPoolAndProxyInterface interface {
	CanAcceptLiquidity(opts *bind.CallOpts) (bool, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetArmProxy(opts *bind.CallOpts) (common.Address, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetLockReleaseInterfaceId(opts *bind.CallOpts) ([4]byte, error)

	GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

	GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

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

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
