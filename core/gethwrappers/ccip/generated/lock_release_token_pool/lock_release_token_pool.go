// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lock_release_token_pool

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

var LockReleaseTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"acceptLiquidity\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LiquidityNotAccepted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canAcceptLiquidity\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRebalancer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"provideLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebalancer\",\"type\":\"address\"}],\"name\":\"setRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200499e3803806200499e833981016040819052620000359162000565565b848484833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c3816200017e565b5050506001600160a01b0384161580620000e457506001600160a01b038116155b80620000f757506001600160a01b038216155b1562000116576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c052620001695760408051600081526020810190915262000169908462000229565b5050505090151560e05250620006d692505050565b336001600160a01b03821603620001d85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c0516200024a576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002d55760008382815181106200026e576200026e62000688565b602090810291909101015190506200028860028262000386565b15620002cb576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016200024d565b5060005b815181101562000381576000828281518110620002fa57620002fa62000688565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000326575062000378565b62000333600282620003a6565b1562000376576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002d9565b505050565b60006200039d836001600160a01b038416620003bd565b90505b92915050565b60006200039d836001600160a01b038416620004c1565b60008181526001830160205260408120548015620004b6576000620003e46001836200069e565b8554909150600090620003fa906001906200069e565b9050818114620004665760008660000182815481106200041e576200041e62000688565b906000526020600020015490508087600001848154811062000444576200044462000688565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200047a576200047a620006c0565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a0565b6000915050620003a0565b60008181526001830160205260408120546200050a57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a0565b506000620003a0565b6001600160a01b03811681146200052957600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b80516200054f8162000513565b919050565b805180151581146200054f57600080fd5b600080600080600060a086880312156200057e57600080fd5b85516200058b8162000513565b602087810151919650906001600160401b0380821115620005ab57600080fd5b818901915089601f830112620005c057600080fd5b815181811115620005d557620005d56200052c565b8060051b604051601f19603f83011681018181108582111715620005fd57620005fd6200052c565b60405291825284820192508381018501918c8311156200061c57600080fd5b938501935b828510156200064557620006358562000542565b8452938501939285019262000621565b8099505050505050506200065c6040870162000542565b92506200066c6060870162000554565b91506200067c6080870162000542565b90509295509295909350565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003a057634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e05161422c62000772600039600081816104ef015261173201526000818161059c01528181611cce015261277301526000818161057601528181611aff0152611f84015260008181610290015281816102e50152818161077a0152818161084c015281816108dd015281816117f401528181611a1f01528181611ea401528181612709015261295e015261422c6000f3fe608060405234801561001057600080fd5b50600436106101f05760003560e01c80638da5cb5b1161010f578063c4bffe2b116100a2578063dc0bd97111610071578063dc0bd97114610574578063e0351e131461059a578063eb521a4c146105c0578063f2fde38b146105d357600080fd5b8063c4bffe2b14610526578063c75eea9c1461053b578063cf7401f31461054e578063db6327dc1461056157600080fd5b8063b0f479a1116100de578063b0f479a1146104bc578063b7946580146104da578063bb98546b146104ed578063c0d786551461051357600080fd5b80638da5cb5b146103fa5780639a4575b914610418578063a7cd63b714610438578063af58d59f1461044d57600080fd5b806354c8a4f31161018757806378a010b21161015657806378a010b2146103b957806379ba5097146103cc5780637d54534e146103d45780638926f54f146103e757600080fd5b806354c8a4f31461036257806366320087146103755780636cfd1553146103885780636d3d1a581461039b57600080fd5b806321df0da7116101c357806321df0da71461028e578063240028e8146102d55780633907753714610322578063432a6ba31461034457600080fd5b806301ffc9a7146101f55780630a2fd4931461021d5780630a861f2a1461023d578063181f5a7714610252575b600080fd5b610208610203366004613308565b6105e6565b60405190151581526020015b60405180910390f35b61023061022b366004613367565b610642565b60405161021491906133f0565b61025061024b366004613403565b6106f2565b005b6102306040518060400160405280601e81526020017f4c6f636b52656c65617365546f6b656e506f6f6c20312e352e302d646576000081525081565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610214565b6102086102e3366004613449565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b610335610330366004613466565b6108a3565b60405190518152602001610214565b60095473ffffffffffffffffffffffffffffffffffffffff166102b0565b6102506103703660046134ee565b610999565b61025061038336600461355a565b610a14565b610250610396366004613449565b610af0565b60085473ffffffffffffffffffffffffffffffffffffffff166102b0565b6102506103c7366004613586565b610b3f565b610250610cae565b6102506103e2366004613449565b610dab565b6102086103f5366004613367565b610dfa565b60005473ffffffffffffffffffffffffffffffffffffffff166102b0565b61042b610426366004613609565b610e11565b6040516102149190613644565b610440610eab565b60405161021491906136a4565b61046061045b366004613367565b610ebc565b604051610214919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff166102b0565b6102306104e8366004613367565b610f91565b7f0000000000000000000000000000000000000000000000000000000000000000610208565b610250610521366004613449565b610fbc565b61052e611097565b60405161021491906136fe565b610460610549366004613367565b61114f565b61025061055c366004613866565b611221565b61025061056f3660046138ab565b6112aa565b7f00000000000000000000000000000000000000000000000000000000000000006102b0565b7f0000000000000000000000000000000000000000000000000000000000000000610208565b6102506105ce366004613403565b611730565b6102506105e1366004613449565b61184c565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fe1d4056600000000000000000000000000000000000000000000000000000000148061063c575061063c82611860565b92915050565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061066d906138ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610699906138ed565b80156106e65780601f106106bb576101008083540402835291602001916106e6565b820191906000526020600020905b8154815290600101906020018083116106c957829003601f168201915b50505050509050919050565b60095473ffffffffffffffffffffffffffffffffffffffff16331461074a576040517f8e4a23d60000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015281907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156107d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107fa9190613940565b1015610832576040517fbb55fd2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61087373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163383611944565b604051819033907fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171990600090a350565b6040805160208101909152600081526108c36108be83613a04565b611a18565b61090873ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016336060850135611944565b6109186060830160408401613449565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52846060013560405161097a91815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6109a1611c49565b610a0e84848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611ccc92505050565b50505050565b610a1c611c49565b6040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810182905273ffffffffffffffffffffffffffffffffffffffff831690630a861f2a90602401600060405180830381600087803b158015610a8457600080fd5b505af1158015610a98573d6000803e3d6000fd5b505050508173ffffffffffffffffffffffffffffffffffffffff167f6fa7abcf1345d1d478e5ea0da6b5f26a90eadb0546ef15ed3833944fbfd1db6282604051610ae491815260200190565b60405180910390a25050565b610af8611c49565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610b47611c49565b610b5083610dfa565b610b92576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610741565b67ffffffffffffffff831660009081526007602052604081206004018054610bb9906138ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610be5906138ed565b8015610c325780601f10610c0757610100808354040283529160200191610c32565b820191906000526020600020905b815481529060010190602001808311610c1557829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610c61838583613b49565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610ca093929190613c64565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610d2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610741565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610db3611c49565b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600061063c600567ffffffffffffffff8416611e82565b6040805180820190915260608082526020820152610e36610e3183613cc8565b611e9d565b6040516060830135815233907f9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd600089060200160405180910390a26040518060400160405280610e908460200160208101906104e89190613367565b81526040805160208181019092526000815291015292915050565b6060610eb76002612067565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261063c90612074565b67ffffffffffffffff8116600090815260076020526040902060050180546060919061066d906138ed565b610fc4611c49565b73ffffffffffffffffffffffffffffffffffffffff8116611011576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110a56005612067565b90506000815167ffffffffffffffff8111156110c3576110c3613740565b6040519080825280602002602001820160405280156110ec578160200160208202803683370190505b50905060005b82518110156111485782818151811061110d5761110d613d6a565b602002602001015182828151811061112757611127613d6a565b67ffffffffffffffff909216602092830291909101909101526001016110f2565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261063c90612074565b60085473ffffffffffffffffffffffffffffffffffffffff163314801590611261575060005473ffffffffffffffffffffffffffffffffffffffff163314155b1561129a576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610741565b6112a5838383612126565b505050565b6112b2611c49565b60005b818110156112a55760008383838181106112d1576112d1613d6a565b90506020028101906112e39190613d99565b6112ec90613dd7565b90506113018160800151826020015115612210565b6113148160a00151826020015115612210565b8060200151156116105780516113369060059067ffffffffffffffff16612349565b61137b5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610741565b60408101515115806113905750606081015151155b156113c7576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c179091169615150295909517909855908101519401519381169316909102919091176003820155915190919060048201906115a89082613e8b565b50606082015160058201906115bd9082613e8b565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506116039493929190613fa5565b60405180910390a1611727565b80516116289060059067ffffffffffffffff16612355565b61166d5780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610741565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906116d660048301826132ba565b6116e46005830160006132ba565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b506001016112b5565b7f0000000000000000000000000000000000000000000000000000000000000000611787576040517fe93f8fa400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60095473ffffffffffffffffffffffffffffffffffffffff1633146117da576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610741565b61181c73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612361565b604051819033907fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208890600090a350565b611854611c49565b61185d816123bf565b50565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf0000000000000000000000000000000000000000000000000000000014806118f357507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061063c57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526112a59084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526124b4565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611aad5760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610741565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611b5b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b7f919061403e565b15611bb6576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611bc381602001516125c0565b6000611bd28260200151610642565b9050805160001480611bf6575080805190602001208260a001518051906020012014155b15611c33578160a001516040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161074191906133f0565b611c45826020015183606001516126e6565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611cca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610741565b565b7f0000000000000000000000000000000000000000000000000000000000000000611d23576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611db9576000838281518110611d4357611d43613d6a565b60200260200101519050611d6181600261272d90919063ffffffff16565b15611db05760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611d26565b5060005b81518110156112a5576000828281518110611dda57611dda613d6a565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611e1e5750611e7a565b611e2960028261274f565b15611e785760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611dbd565b600081815260018301602052604081205415155b9392505050565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611f325760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610741565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611fe0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612004919061403e565b1561203b576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6120488160400151612771565b61205581602001516127f0565b61185d8160200151826060015161293e565b60606000611e9683612982565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261210282606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426120e6919061408a565b85608001516fffffffffffffffffffffffffffffffff166129dd565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61212f83610dfa565b612171576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610741565b61217c826000612210565b67ffffffffffffffff8316600090815260076020526040902061219f9083612a07565b6121aa816000612210565b67ffffffffffffffff831660009081526007602052604090206121d09060020182612a07565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516122039392919061409d565b60405180910390a1505050565b8151156122d75781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612266575060408201516fffffffffffffffffffffffffffffffff16155b1561229f57816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016107419190614120565b8015611c45576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff16151580612310575060208201516fffffffffffffffffffffffffffffffff1615155b15611c4557816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016107419190614120565b6000611e968383612ba9565b6000611e968383612bf8565b60405173ffffffffffffffffffffffffffffffffffffffff80851660248301528316604482015260648101829052610a0e9085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401611996565b3373ffffffffffffffffffffffffffffffffffffffff82160361243e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610741565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612516826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16612ceb9092919063ffffffff16565b8051909150156112a55780806020019051810190612534919061403e565b6112a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610741565b6125c981610dfa565b61260b576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610741565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa15801561268a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126ae919061403e565b61185d576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610741565b67ffffffffffffffff82166000908152600760205260409020611c4590600201827f0000000000000000000000000000000000000000000000000000000000000000612cfa565b6000611e968373ffffffffffffffffffffffffffffffffffffffff8416612bf8565b6000611e968373ffffffffffffffffffffffffffffffffffffffff8416612ba9565b7f00000000000000000000000000000000000000000000000000000000000000001561185d576127a260028261307d565b61185d576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610741565b6127f981610dfa565b61283b576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610741565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa1580156128b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128d8919061415c565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461185d576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610741565b67ffffffffffffffff82166000908152600760205260409020611c4590827f0000000000000000000000000000000000000000000000000000000000000000612cfa565b6060816000018054806020026020016040519081016040528092919081815260200182805480156106e657602002820191906000526020600020905b8154815260200190600101908083116129be5750505050509050919050565b60006129fc856129ed8486614179565b6129f79087614190565b6130ac565b90505b949350505050565b8154600090612a3090700100000000000000000000000000000000900463ffffffff164261408a565b90508015612ad25760018301548354612a78916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166129dd565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612af8916fffffffffffffffffffffffffffffffff90811691166130ac565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612203908490614120565b6000818152600183016020526040812054612bf05750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561063c565b50600061063c565b60008181526001830160205260408120548015612ce1576000612c1c60018361408a565b8554909150600090612c309060019061408a565b9050818114612c95576000866000018281548110612c5057612c50613d6a565b9060005260206000200154905080876000018481548110612c7357612c73613d6a565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612ca657612ca66141a3565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061063c565b600091505061063c565b60606129ff84846000856130c2565b825474010000000000000000000000000000000000000000900460ff161580612d21575081155b15612d2b57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612d7190700100000000000000000000000000000000900463ffffffff164261408a565b90508015612e315781831115612db3576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612ded9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166129dd565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612ee85773ffffffffffffffffffffffffffffffffffffffff8416612e90576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610741565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610741565b84831015612ffb5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612f2c908261408a565b612f36878a61408a565b612f409190614190565b612f4a91906141d2565b905073ffffffffffffffffffffffffffffffffffffffff8616612fa3576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610741565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610741565b613005858461408a565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611e96565b60008183106130bb5781611e96565b5090919050565b606082471015613154576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610741565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161317d919061420d565b60006040518083038185875af1925050503d80600081146131ba576040519150601f19603f3d011682016040523d82523d6000602084013e6131bf565b606091505b50915091506131d0878383876131db565b979650505050505050565b6060831561327157825160000361326a5773ffffffffffffffffffffffffffffffffffffffff85163b61326a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610741565b50816129ff565b6129ff83838151156132865781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161074191906133f0565b5080546132c6906138ed565b6000825580601f106132d6575050565b601f01602090049060005260206000209081019061185d91905b8082111561330457600081556001016132f0565b5090565b60006020828403121561331a57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611e9657600080fd5b803567ffffffffffffffff8116811461336257600080fd5b919050565b60006020828403121561337957600080fd5b611e968261334a565b60005b8381101561339d578181015183820152602001613385565b50506000910152565b600081518084526133be816020860160208601613382565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611e9660208301846133a6565b60006020828403121561341557600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461185d57600080fd5b80356133628161341c565b60006020828403121561345b57600080fd5b8135611e968161341c565b60006020828403121561347857600080fd5b813567ffffffffffffffff81111561348f57600080fd5b82016101008185031215611e9657600080fd5b60008083601f8401126134b457600080fd5b50813567ffffffffffffffff8111156134cc57600080fd5b6020830191508360208260051b85010111156134e757600080fd5b9250929050565b6000806000806040858703121561350457600080fd5b843567ffffffffffffffff8082111561351c57600080fd5b613528888389016134a2565b9096509450602087013591508082111561354157600080fd5b5061354e878288016134a2565b95989497509550505050565b6000806040838503121561356d57600080fd5b82356135788161341c565b946020939093013593505050565b60008060006040848603121561359b57600080fd5b6135a48461334a565b9250602084013567ffffffffffffffff808211156135c157600080fd5b818601915086601f8301126135d557600080fd5b8135818111156135e457600080fd5b8760208285010111156135f657600080fd5b6020830194508093505050509250925092565b60006020828403121561361b57600080fd5b813567ffffffffffffffff81111561363257600080fd5b820160a08185031215611e9657600080fd5b60208152600082516040602084015261366060608401826133a6565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084830301604085015261369b82826133a6565b95945050505050565b6020808252825182820181905260009190848201906040850190845b818110156136f257835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016136c0565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156136f257835167ffffffffffffffff168352928401929184019160010161371a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff8111828210171561379357613793613740565b60405290565b60405160c0810167ffffffffffffffff8111828210171561379357613793613740565b801515811461185d57600080fd5b8035613362816137bc565b80356fffffffffffffffffffffffffffffffff8116811461336257600080fd5b60006060828403121561380757600080fd5b6040516060810181811067ffffffffffffffff8211171561382a5761382a613740565b604052905080823561383b816137bc565b8152613849602084016137d5565b602082015261385a604084016137d5565b60408201525092915050565b600080600060e0848603121561387b57600080fd5b6138848461334a565b925061389385602086016137f5565b91506138a285608086016137f5565b90509250925092565b600080602083850312156138be57600080fd5b823567ffffffffffffffff8111156138d557600080fd5b6138e1858286016134a2565b90969095509350505050565b600181811c9082168061390157607f821691505b60208210810361393a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60006020828403121561395257600080fd5b5051919050565b600082601f83011261396a57600080fd5b813567ffffffffffffffff8082111561398557613985613740565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156139cb576139cb613740565b816040528381528660208588010111156139e457600080fd5b836020870160208301376000602085830101528094505050505092915050565b60006101008236031215613a1757600080fd5b613a1f61376f565b823567ffffffffffffffff80821115613a3757600080fd5b613a4336838701613959565b8352613a516020860161334a565b6020840152613a626040860161343e565b604084015260608501356060840152613a7d6080860161343e565b608084015260a0850135915080821115613a9657600080fd5b613aa236838701613959565b60a084015260c0850135915080821115613abb57600080fd5b613ac736838701613959565b60c084015260e0850135915080821115613ae057600080fd5b50613aed36828601613959565b60e08301525092915050565b601f8211156112a5576000816000526020600020601f850160051c81016020861015613b225750805b601f850160051c820191505b81811015613b4157828155600101613b2e565b505050505050565b67ffffffffffffffff831115613b6157613b61613740565b613b7583613b6f83546138ed565b83613af9565b6000601f841160018114613bc75760008515613b915750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355613c5d565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015613c165786850135825560209485019460019092019101613bf6565b5086821015613c51577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b604081526000613c7760408301866133a6565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060a08236031215613cda57600080fd5b60405160a0810167ffffffffffffffff8282108183111715613cfe57613cfe613740565b816040528435915080821115613d1357600080fd5b50613d2036828601613959565b825250613d2f6020840161334a565b60208201526040830135613d428161341c565b6040820152606083810135908201526080830135613d5f8161341c565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec1833603018112613dcd57600080fd5b9190910192915050565b60006101408236031215613dea57600080fd5b613df2613799565b613dfb8361334a565b8152613e09602084016137ca565b6020820152604083013567ffffffffffffffff80821115613e2957600080fd5b613e3536838701613959565b60408401526060850135915080821115613e4e57600080fd5b50613e5b36828601613959565b606083015250613e6e36608085016137f5565b6080820152613e803660e085016137f5565b60a082015292915050565b815167ffffffffffffffff811115613ea557613ea5613740565b613eb981613eb384546138ed565b84613af9565b602080601f831160018114613f0c5760008415613ed65750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613b41565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613f5957888601518255948401946001909101908401613f3a565b5085821015613f9557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613fc9818401876133a6565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506140079050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e083015261369b565b60006020828403121561405057600080fd5b8151611e96816137bc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561063c5761063c61405b565b67ffffffffffffffff8416815260e081016140e960208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526129ff565b6060810161063c82848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561416e57600080fd5b8151611e968161341c565b808202811582820484141761063c5761063c61405b565b8082018082111561063c5761063c61405b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600082614208577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008251613dcd81846020870161338256fea164736f6c6343000818000a",
}

var LockReleaseTokenPoolABI = LockReleaseTokenPoolMetaData.ABI

var LockReleaseTokenPoolBin = LockReleaseTokenPoolMetaData.Bin

func DeployLockReleaseTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, acceptLiquidity bool, router common.Address) (common.Address, *types.Transaction, *LockReleaseTokenPool, error) {
	parsed, err := LockReleaseTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LockReleaseTokenPoolBin), backend, token, allowlist, rmnProxy, acceptLiquidity, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockReleaseTokenPool{address: address, abi: *parsed, LockReleaseTokenPoolCaller: LockReleaseTokenPoolCaller{contract: contract}, LockReleaseTokenPoolTransactor: LockReleaseTokenPoolTransactor{contract: contract}, LockReleaseTokenPoolFilterer: LockReleaseTokenPoolFilterer{contract: contract}}, nil
}

type LockReleaseTokenPool struct {
	address common.Address
	abi     abi.ABI
	LockReleaseTokenPoolCaller
	LockReleaseTokenPoolTransactor
	LockReleaseTokenPoolFilterer
}

type LockReleaseTokenPoolCaller struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolSession struct {
	Contract     *LockReleaseTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolCallerSession struct {
	Contract *LockReleaseTokenPoolCaller
	CallOpts bind.CallOpts
}

type LockReleaseTokenPoolTransactorSession struct {
	Contract     *LockReleaseTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolRaw struct {
	Contract *LockReleaseTokenPool
}

type LockReleaseTokenPoolCallerRaw struct {
	Contract *LockReleaseTokenPoolCaller
}

type LockReleaseTokenPoolTransactorRaw struct {
	Contract *LockReleaseTokenPoolTransactor
}

func NewLockReleaseTokenPool(address common.Address, backend bind.ContractBackend) (*LockReleaseTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(LockReleaseTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLockReleaseTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPool{address: address, abi: abi, LockReleaseTokenPoolCaller: LockReleaseTokenPoolCaller{contract: contract}, LockReleaseTokenPoolTransactor: LockReleaseTokenPoolTransactor{contract: contract}, LockReleaseTokenPoolFilterer: LockReleaseTokenPoolFilterer{contract: contract}}, nil
}

func NewLockReleaseTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*LockReleaseTokenPoolCaller, error) {
	contract, err := bindLockReleaseTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolCaller{contract: contract}, nil
}

func NewLockReleaseTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*LockReleaseTokenPoolTransactor, error) {
	contract, err := bindLockReleaseTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolTransactor{contract: contract}, nil
}

func NewLockReleaseTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*LockReleaseTokenPoolFilterer, error) {
	contract, err := bindLockReleaseTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolFilterer{contract: contract}, nil
}

func bindLockReleaseTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LockReleaseTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolTransactor.contract.Transfer(opts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.contract.Transfer(opts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) CanAcceptLiquidity(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "canAcceptLiquidity")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPool.Contract.CanAcceptLiquidity(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPool.Contract.CanAcceptLiquidity(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetAllowList(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetAllowList(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPool.Contract.GetAllowListEnabled(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPool.Contract.GetAllowListEnabled(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRateLimitAdmin(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRateLimitAdmin(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRebalancer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRebalancer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRebalancer(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRebalancer(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemotePool(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemotePool(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemoteToken(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemoteToken(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRmnProxy(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRmnProxy(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRouter(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRouter(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPool.Contract.GetSupportedChains(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPool.Contract.GetSupportedChains(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetToken(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetToken(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedChain(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedChain(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedToken(&_LockReleaseTokenPool.CallOpts, token)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedToken(&_LockReleaseTokenPool.CallOpts, token)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) Owner() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.Owner(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) Owner() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.Owner(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.SupportsInterface(&_LockReleaseTokenPool.CallOpts, interfaceId)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.SupportsInterface(&_LockReleaseTokenPool.CallOpts, interfaceId)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPool.Contract.TypeAndVersion(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPool.Contract.TypeAndVersion(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AcceptOwnership(&_LockReleaseTokenPool.TransactOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AcceptOwnership(&_LockReleaseTokenPool.TransactOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPool.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPool.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyChainUpdates(&_LockReleaseTokenPool.TransactOpts, chains)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyChainUpdates(&_LockReleaseTokenPool.TransactOpts, chains)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockOrBurn(&_LockReleaseTokenPool.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockOrBurn(&_LockReleaseTokenPool.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "provideLiquidity", amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ProvideLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ProvideLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ReleaseOrMint(&_LockReleaseTokenPool.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ReleaseOrMint(&_LockReleaseTokenPool.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRateLimitAdmin(&_LockReleaseTokenPool.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRateLimitAdmin(&_LockReleaseTokenPool.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRebalancer", rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRebalancer(&_LockReleaseTokenPool.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRebalancer(&_LockReleaseTokenPool.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRouter(&_LockReleaseTokenPool.TransactOpts, newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRouter(&_LockReleaseTokenPool.TransactOpts, newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "transferLiquidity", from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferLiquidity(&_LockReleaseTokenPool.TransactOpts, from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferLiquidity(&_LockReleaseTokenPool.TransactOpts, from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferOwnership(&_LockReleaseTokenPool.TransactOpts, to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferOwnership(&_LockReleaseTokenPool.TransactOpts, to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "withdrawLiquidity", amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.WithdrawLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.WithdrawLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

type LockReleaseTokenPoolAllowListAddIterator struct {
	Event *LockReleaseTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAllowListAdd)
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
		it.Event = new(LockReleaseTokenPoolAllowListAdd)
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

func (it *LockReleaseTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAllowListAddIterator{contract: _LockReleaseTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAllowListAdd)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAllowListAdd, error) {
	event := new(LockReleaseTokenPoolAllowListAdd)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAllowListRemoveIterator struct {
	Event *LockReleaseTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAllowListRemove)
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
		it.Event = new(LockReleaseTokenPoolAllowListRemove)
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

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAllowListRemoveIterator{contract: _LockReleaseTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAllowListRemove)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAllowListRemove, error) {
	event := new(LockReleaseTokenPoolAllowListRemove)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolBurnedIterator struct {
	Event *LockReleaseTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolBurned)
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
		it.Event = new(LockReleaseTokenPoolBurned)
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

func (it *LockReleaseTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolBurnedIterator{contract: _LockReleaseTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolBurned)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseBurned(log types.Log) (*LockReleaseTokenPoolBurned, error) {
	event := new(LockReleaseTokenPoolBurned)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainAddedIterator struct {
	Event *LockReleaseTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainAdded)
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
		it.Event = new(LockReleaseTokenPoolChainAdded)
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

func (it *LockReleaseTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainAddedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainAddedIterator{contract: _LockReleaseTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainAdded)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainAdded(log types.Log) (*LockReleaseTokenPoolChainAdded, error) {
	event := new(LockReleaseTokenPoolChainAdded)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainConfiguredIterator struct {
	Event *LockReleaseTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainConfigured)
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
		it.Event = new(LockReleaseTokenPoolChainConfigured)
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

func (it *LockReleaseTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainConfiguredIterator{contract: _LockReleaseTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainConfigured)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolChainConfigured, error) {
	event := new(LockReleaseTokenPoolChainConfigured)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainRemovedIterator struct {
	Event *LockReleaseTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainRemoved)
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
		it.Event = new(LockReleaseTokenPoolChainRemoved)
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

func (it *LockReleaseTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainRemovedIterator{contract: _LockReleaseTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainRemoved)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolChainRemoved, error) {
	event := new(LockReleaseTokenPoolChainRemoved)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolConfigChangedIterator struct {
	Event *LockReleaseTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolConfigChanged)
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
		it.Event = new(LockReleaseTokenPoolConfigChanged)
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

func (it *LockReleaseTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolConfigChangedIterator{contract: _LockReleaseTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolConfigChanged)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolConfigChanged, error) {
	event := new(LockReleaseTokenPoolConfigChanged)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityAddedIterator struct {
	Event *LockReleaseTokenPoolLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityAdded)
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
		it.Event = new(LockReleaseTokenPoolLiquidityAdded)
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

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityAdded struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityAddedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityAddedIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityAdded)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolLiquidityAdded, error) {
	event := new(LockReleaseTokenPoolLiquidityAdded)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityRemovedIterator struct {
	Event *LockReleaseTokenPoolLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityRemoved)
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
		it.Event = new(LockReleaseTokenPoolLiquidityRemoved)
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

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityRemoved struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityRemovedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityRemovedIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityRemoved)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolLiquidityRemoved, error) {
	event := new(LockReleaseTokenPoolLiquidityRemoved)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityTransferredIterator struct {
	Event *LockReleaseTokenPoolLiquidityTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityTransferred)
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
		it.Event = new(LockReleaseTokenPoolLiquidityTransferred)
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

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityTransferred struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityTransferred(opts *bind.FilterOpts, from []common.Address) (*LockReleaseTokenPoolLiquidityTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityTransferred", fromRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityTransferredIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityTransferred", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityTransferred, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityTransferred", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityTransferred)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityTransferred(log types.Log) (*LockReleaseTokenPoolLiquidityTransferred, error) {
	event := new(LockReleaseTokenPoolLiquidityTransferred)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLockedIterator struct {
	Event *LockReleaseTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLocked)
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
		it.Event = new(LockReleaseTokenPoolLocked)
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

func (it *LockReleaseTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLockedIterator{contract: _LockReleaseTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLocked)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLocked(log types.Log) (*LockReleaseTokenPoolLocked, error) {
	event := new(LockReleaseTokenPoolLocked)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolMintedIterator struct {
	Event *LockReleaseTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolMinted)
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
		it.Event = new(LockReleaseTokenPoolMinted)
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

func (it *LockReleaseTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolMintedIterator{contract: _LockReleaseTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolMinted)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseMinted(log types.Log) (*LockReleaseTokenPoolMinted, error) {
	event := new(LockReleaseTokenPoolMinted)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolOwnershipTransferRequestedIterator struct {
	Event *LockReleaseTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolOwnershipTransferRequested)
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
		it.Event = new(LockReleaseTokenPoolOwnershipTransferRequested)
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

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolOwnershipTransferRequestedIterator{contract: _LockReleaseTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolOwnershipTransferRequested)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolOwnershipTransferRequested, error) {
	event := new(LockReleaseTokenPoolOwnershipTransferRequested)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolOwnershipTransferredIterator struct {
	Event *LockReleaseTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolOwnershipTransferred)
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
		it.Event = new(LockReleaseTokenPoolOwnershipTransferred)
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

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolOwnershipTransferredIterator{contract: _LockReleaseTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolOwnershipTransferred)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolOwnershipTransferred, error) {
	event := new(LockReleaseTokenPoolOwnershipTransferred)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolReleasedIterator struct {
	Event *LockReleaseTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolReleased)
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
		it.Event = new(LockReleaseTokenPoolReleased)
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

func (it *LockReleaseTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolReleasedIterator{contract: _LockReleaseTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolReleased)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseReleased(log types.Log) (*LockReleaseTokenPoolReleased, error) {
	event := new(LockReleaseTokenPoolReleased)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRemotePoolSetIterator struct {
	Event *LockReleaseTokenPoolRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRemotePoolSet)
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
		it.Event = new(LockReleaseTokenPoolRemotePoolSet)
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

func (it *LockReleaseTokenPoolRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRemotePoolSetIterator{contract: _LockReleaseTokenPool.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRemotePoolSet)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRemotePoolSet(log types.Log) (*LockReleaseTokenPoolRemotePoolSet, error) {
	event := new(LockReleaseTokenPoolRemotePoolSet)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRouterUpdatedIterator struct {
	Event *LockReleaseTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRouterUpdated)
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
		it.Event = new(LockReleaseTokenPoolRouterUpdated)
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

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRouterUpdatedIterator{contract: _LockReleaseTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRouterUpdated)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolRouterUpdated, error) {
	event := new(LockReleaseTokenPoolRouterUpdated)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolTokensConsumedIterator struct {
	Event *LockReleaseTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolTokensConsumed)
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
		it.Event = new(LockReleaseTokenPoolTokensConsumed)
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

func (it *LockReleaseTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolTokensConsumedIterator{contract: _LockReleaseTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolTokensConsumed)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolTokensConsumed, error) {
	event := new(LockReleaseTokenPoolTokensConsumed)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LockReleaseTokenPool.abi.Events["AllowListAdd"].ID:
		return _LockReleaseTokenPool.ParseAllowListAdd(log)
	case _LockReleaseTokenPool.abi.Events["AllowListRemove"].ID:
		return _LockReleaseTokenPool.ParseAllowListRemove(log)
	case _LockReleaseTokenPool.abi.Events["Burned"].ID:
		return _LockReleaseTokenPool.ParseBurned(log)
	case _LockReleaseTokenPool.abi.Events["ChainAdded"].ID:
		return _LockReleaseTokenPool.ParseChainAdded(log)
	case _LockReleaseTokenPool.abi.Events["ChainConfigured"].ID:
		return _LockReleaseTokenPool.ParseChainConfigured(log)
	case _LockReleaseTokenPool.abi.Events["ChainRemoved"].ID:
		return _LockReleaseTokenPool.ParseChainRemoved(log)
	case _LockReleaseTokenPool.abi.Events["ConfigChanged"].ID:
		return _LockReleaseTokenPool.ParseConfigChanged(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityAdded"].ID:
		return _LockReleaseTokenPool.ParseLiquidityAdded(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityRemoved"].ID:
		return _LockReleaseTokenPool.ParseLiquidityRemoved(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityTransferred"].ID:
		return _LockReleaseTokenPool.ParseLiquidityTransferred(log)
	case _LockReleaseTokenPool.abi.Events["Locked"].ID:
		return _LockReleaseTokenPool.ParseLocked(log)
	case _LockReleaseTokenPool.abi.Events["Minted"].ID:
		return _LockReleaseTokenPool.ParseMinted(log)
	case _LockReleaseTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _LockReleaseTokenPool.ParseOwnershipTransferRequested(log)
	case _LockReleaseTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _LockReleaseTokenPool.ParseOwnershipTransferred(log)
	case _LockReleaseTokenPool.abi.Events["Released"].ID:
		return _LockReleaseTokenPool.ParseReleased(log)
	case _LockReleaseTokenPool.abi.Events["RemotePoolSet"].ID:
		return _LockReleaseTokenPool.ParseRemotePoolSet(log)
	case _LockReleaseTokenPool.abi.Events["RouterUpdated"].ID:
		return _LockReleaseTokenPool.ParseRouterUpdated(log)
	case _LockReleaseTokenPool.abi.Events["TokensConsumed"].ID:
		return _LockReleaseTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LockReleaseTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (LockReleaseTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (LockReleaseTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (LockReleaseTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (LockReleaseTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (LockReleaseTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (LockReleaseTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (LockReleaseTokenPoolLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (LockReleaseTokenPoolLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (LockReleaseTokenPoolLiquidityTransferred) Topic() common.Hash {
	return common.HexToHash("0x6fa7abcf1345d1d478e5ea0da6b5f26a90eadb0546ef15ed3833944fbfd1db62")
}

func (LockReleaseTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (LockReleaseTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (LockReleaseTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LockReleaseTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LockReleaseTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (LockReleaseTokenPoolRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (LockReleaseTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (LockReleaseTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_LockReleaseTokenPool *LockReleaseTokenPool) Address() common.Address {
	return _LockReleaseTokenPool.address
}

type LockReleaseTokenPoolInterface interface {
	CanAcceptLiquidity(opts *bind.CallOpts) (bool, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

	GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

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

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*LockReleaseTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*LockReleaseTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolConfigChanged, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolLiquidityAdded, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolLiquidityRemoved, error)

	FilterLiquidityTransferred(opts *bind.FilterOpts, from []common.Address) (*LockReleaseTokenPoolLiquidityTransferredIterator, error)

	WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityTransferred, from []common.Address) (event.Subscription, error)

	ParseLiquidityTransferred(log types.Log) (*LockReleaseTokenPoolLiquidityTransferred, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*LockReleaseTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*LockReleaseTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*LockReleaseTokenPoolReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*LockReleaseTokenPoolRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
