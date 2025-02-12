// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_with_from_mint_token_pool_and_proxy

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

var BurnWithFromMintTokenPoolAndProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"oldPool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"LegacyPoolChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"onRampAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPreviousPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"prevPool\",\"type\":\"address\"}],\"name\":\"setPreviousPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004dde38038062004dde8339810160408190526200003491620008cc565b83838383838383833380600081620000935760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c657620000c6816200019b565b5050506001600160a01b0384161580620000e757506001600160a01b038116155b80620000fa57506001600160a01b038216155b1562000119576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200016c576040805160008152602081019091526200016c908462000246565b50620001919650506001600160a01b038a169450309350600019925050620003a39050565b5050505062000b08565b336001600160a01b03821603620001f55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c05162000267576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002f25760008382815181106200028b576200028b620009dc565b60209081029190910101519050620002a560028262000489565b15620002e8576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016200026a565b5060005b81518110156200039e576000828281518110620003175762000317620009dc565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000343575062000395565b62000350600282620004a9565b1562000393576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002f6565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015620003f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200041b9190620009f2565b62000427919062000a22565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b179091529192506200048391869190620004c016565b50505050565b6000620004a0836001600160a01b03841662000591565b90505b92915050565b6000620004a0836001600160a01b03841662000695565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564908201526000906200050f906001600160a01b038516908490620006e7565b8051909150156200039e578080602001905181019062000530919062000a38565b6200039e5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016200008a565b600081815260018301602052604081205480156200068a576000620005b860018362000a63565b8554909150600090620005ce9060019062000a63565b90508082146200063a576000866000018281548110620005f257620005f2620009dc565b9060005260206000200154905080876000018481548110620006185762000618620009dc565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200064e576200064e62000a79565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620004a3565b6000915050620004a3565b6000818152600183016020526040812054620006de57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620004a3565b506000620004a3565b6060620006f8848460008562000700565b949350505050565b606082471015620007635760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016200008a565b600080866001600160a01b0316858760405162000781919062000ab5565b60006040518083038185875af1925050503d8060008114620007c0576040519150601f19603f3d011682016040523d82523d6000602084013e620007c5565b606091505b509092509050620007d987838387620007e4565b979650505050505050565b606083156200085857825160000362000850576001600160a01b0385163b620008505760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016200008a565b5081620006f8565b620006f883838151156200086f5781518083602001fd5b8060405162461bcd60e51b81526004016200008a919062000ad3565b6001600160a01b0381168114620008a157600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b8051620008c7816200088b565b919050565b60008060008060808587031215620008e357600080fd5b8451620008f0816200088b565b602086810151919550906001600160401b03808211156200091057600080fd5b818801915088601f8301126200092557600080fd5b8151818111156200093a576200093a620008a4565b8060051b604051601f19603f83011681018181108582111715620009625762000962620008a4565b60405291825284820192508381018501918b8311156200098157600080fd5b938501935b82851015620009aa576200099a85620008ba565b8452938501939285019262000986565b809850505050505050620009c160408601620008ba565b9150620009d160608601620008ba565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b60006020828403121562000a0557600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b80820180821115620004a357620004a362000a0c565b60006020828403121562000a4b57600080fd5b8151801515811462000a5c57600080fd5b9392505050565b81810381811115620004a357620004a362000a0c565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000aac57818101518382015260200162000a92565b50506000910152565b6000825162000ac981846020870162000a8f565b9190910192915050565b602081526000825180602084015262000af481604085016020870162000a8f565b601f01601f19169190910160400192915050565b60805160a05160c05161425262000b8c6000396000818161052c01528181611ad60152612528015260008181610506015281816118690152611d89015260008181610231015281816102860152818161076801528181610df40152818161178901528181611ca901528181611e8f015281816124be015261271301526142526000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80639a4575b911610104578063c0d78655116100a2578063db6327dc11610071578063db6327dc146104f1578063dc0bd97114610504578063e0351e131461052a578063f2fde38b1461055057600080fd5b8063c0d78655146104a3578063c4bffe2b146104b6578063c75eea9c146104cb578063cf7401f3146104de57600080fd5b8063a8d87a3b116100de578063a8d87a3b146103f0578063af58d59f14610403578063b0f479a114610472578063b79465801461049057600080fd5b80639a4575b91461039d578063a2b261d8146103bd578063a7cd63b7146103db57600080fd5b80636d3d1a581161017c57806383826b2b1161014b57806383826b2b146103465780638926f54f146103595780638da5cb5b1461036c5780639766b9321461038a57600080fd5b80636d3d1a58146102fa57806378a010b21461031857806379ba50971461032b5780637d54534e1461033357600080fd5b806321df0da7116101b857806321df0da71461022f578063240028e81461027657806339077537146102c357806354c8a4f3146102e557600080fd5b806301ffc9a7146101df5780630a2fd49314610207578063181f5a7714610227575b600080fd5b6101f26101ed3660046131c0565b610563565b60405190151581526020015b60405180910390f35b61021a61021536600461321f565b610648565b6040516101fe91906132a8565b61021a6106f8565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fe565b6101f26102843660046132e8565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102d66102d1366004613305565b610714565b604051905181526020016101fe565b6102f86102f336600461338d565b6108cc565b005b60085473ffffffffffffffffffffffffffffffffffffffff16610251565b6102f86103263660046133f9565b610947565b6102f8610abb565b6102f86103413660046132e8565b610bb8565b6101f261035436600461347c565b610c07565b6101f261036736600461321f565b610cd4565b60005473ffffffffffffffffffffffffffffffffffffffff16610251565b6102f86103983660046132e8565b610ceb565b6103b06103ab3660046134b3565b610d7a565b6040516101fe91906134ee565b60095473ffffffffffffffffffffffffffffffffffffffff16610251565b6103e3610ef0565b6040516101fe919061354e565b6102516103fe36600461321f565b503090565b61041661041136600461321f565b610f01565b6040516101fe919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610251565b61021a61049e36600461321f565b610fd6565b6102f86104b13660046132e8565b611001565b6104be6110d5565b6040516101fe91906135a8565b6104166104d936600461321f565b61118d565b6102f86104ec36600461375f565b61125f565b6102f86104ff3660046137a4565b6112e8565b7f0000000000000000000000000000000000000000000000000000000000000000610251565b7f00000000000000000000000000000000000000000000000000000000000000006101f2565b6102f861055e3660046132e8565b61176e565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf0000000000000000000000000000000000000000000000000000000014806105f657507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061064257507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b67ffffffffffffffff81166000908152600760205260409020600401805460609190610673906137e6565b80601f016020809104026020016040519081016040528092919081815260200182805461069f906137e6565b80156106ec5780601f106106c1576101008083540402835291602001916106ec565b820191906000526020600020905b8154815290600101906020018083116106cf57829003601f168201915b50505050509050919050565b60405180606001604052806027815260200161421f6027913981565b60408051602081019091526000815261073461072f836138d5565b611782565b60095473ffffffffffffffffffffffffffffffffffffffff1661082a5773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f1961079d60608501604086016132e8565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260608501356024820152604401600060405180830381600087803b15801561080d57600080fd5b505af1158015610821573d6000803e3d6000fd5b5050505061083b565b61083b610836836138d5565b6119b3565b61084b60608301604084016132e8565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f084606001356040516108ad91815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6108d4611a51565b61094184848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611ad492505050565b50505050565b61094f611a51565b61095883610cd4565b61099f576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b67ffffffffffffffff8316600090815260076020526040812060040180546109c6906137e6565b80601f01602080910402602001604051908101604052809291908181526020018280546109f2906137e6565b8015610a3f5780601f10610a1457610100808354040283529160200191610a3f565b820191906000526020600020905b815481529060010190602001808311610a2257829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610a6e838583613a1a565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610aad93929190613b34565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b3c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610996565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610bc0611a51565b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600073ffffffffffffffffffffffffffffffffffffffff8216301480610ccd5750600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169281019290925273ffffffffffffffffffffffffffffffffffffffff848116602484015216906383826b2b90604401602060405180830381865afa158015610ca9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ccd9190613b98565b9392505050565b6000610642600567ffffffffffffffff8416611c8a565b610cf3611a51565b6009805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f22891015b60405180910390a15050565b6040805180820190915260608082526020820152610d9f610d9a83613bb5565b611ca2565b60095473ffffffffffffffffffffffffffffffffffffffff16610e6a576040517f79cc6790000000000000000000000000000000000000000000000000000000008152306004820152606083013560248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906379cc679090604401600060405180830381600087803b158015610e4d57600080fd5b505af1158015610e61573d6000803e3d6000fd5b50505050610e7b565b610e7b610e7683613bb5565b611e6c565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610ed584602001602081019061049e919061321f565b81526040805160208181019092526000815291015292915050565b6060610efc6002611f86565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261064290611f93565b67ffffffffffffffff81166000908152600760205260409020600501805460609190610673906137e6565b611009611a51565b73ffffffffffffffffffffffffffffffffffffffff8116611056576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849101610d6e565b606060006110e36005611f86565b90506000815167ffffffffffffffff811115611101576111016135ea565b60405190808252806020026020018201604052801561112a578160200160208202803683370190505b50905060005b82518110156111865782818151811061114b5761114b613c57565b602002602001015182828151811061116557611165613c57565b67ffffffffffffffff90921660209283029190910190910152600101611130565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261064290611f93565b60085473ffffffffffffffffffffffffffffffffffffffff16331480159061129f575060005473ffffffffffffffffffffffffffffffffffffffff163314155b156112d8576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610996565b6112e3838383612045565b505050565b6112f0611a51565b60005b818110156112e357600083838381811061130f5761130f613c57565b90506020028101906113219190613c86565b61132a90613cc4565b905061133f816080015182602001511561212f565b6113528160a0015182602001511561212f565b80602001511561164e5780516113749060059067ffffffffffffffff16612268565b6113b95780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610996565b60408101515115806113ce5750606081015151155b15611405576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c179091169615150295909517909855908101519401519381169316909102919091176003820155915190919060048201906115e69082613d78565b50606082015160058201906115fb9082613d78565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506116419493929190613e92565b60405180910390a1611765565b80516116669060059067ffffffffffffffff16612274565b6116ab5780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610996565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906117146004830182613172565b611722600583016000613172565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b506001016112f3565b611776611a51565b61177f81612280565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff9081169116146118175760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610996565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa1580156118c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118e99190613b98565b15611920576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61192d8160200151612375565b600061193c8260200151610648565b9050805160001480611960575080805190602001208260a001518051906020012014155b1561199d578160a001516040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161099691906132a8565b6119af8260200151836060015161249b565b5050565b60095481516040808401516060850151602086015192517f8627fad600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90951694638627fad694611a1c9490939291600401613f2b565b600060405180830381600087803b158015611a3657600080fd5b505af1158015611a4a573d6000803e3d6000fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611ad2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610996565b565b7f0000000000000000000000000000000000000000000000000000000000000000611b2b576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611bc1576000838281518110611b4b57611b4b613c57565b60200260200101519050611b698160026124e290919063ffffffff16565b15611bb85760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611b2e565b5060005b81518110156112e3576000828281518110611be257611be2613c57565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611c265750611c82565b611c31600282612504565b15611c805760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611bc5565b60008181526001830160205260408120541515610ccd565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611d375760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610996565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611de5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e099190613b98565b15611e40576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611e4d8160400151612526565b611e5a81602001516125a5565b61177f816020015182606001516126f3565b6009546060820151611eb99173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811692911690612737565b60095460408083015183516060850151602086015193517f9687544500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90951694639687544594611f2194939291600401613f8c565b6000604051808303816000875af1158015611f40573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526119af9190810190613fec565b60606000610ccd836127c4565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261202182606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426120059190614089565b85608001516fffffffffffffffffffffffffffffffff1661281f565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61204e83610cd4565b612090576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610996565b61209b82600061212f565b67ffffffffffffffff831660009081526007602052604090206120be9083612849565b6120c981600061212f565b67ffffffffffffffff831660009081526007602052604090206120ef9060020182612849565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516121229392919061409c565b60405180910390a1505050565b8151156121f65781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612185575060408201516fffffffffffffffffffffffffffffffff16155b156121be57816040517f8020d124000000000000000000000000000000000000000000000000000000008152600401610996919061411f565b80156119af576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff1615158061222f575060208201516fffffffffffffffffffffffffffffffff1615155b156119af57816040517fd68af9cc000000000000000000000000000000000000000000000000000000008152600401610996919061411f565b6000610ccd83836129eb565b6000610ccd8383612a3a565b3373ffffffffffffffffffffffffffffffffffffffff8216036122ff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610996565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61237e81610cd4565b6123c0576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610996565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa15801561243f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124639190613b98565b61177f576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610996565b67ffffffffffffffff821660009081526007602052604090206119af90600201827f0000000000000000000000000000000000000000000000000000000000000000612b2d565b6000610ccd8373ffffffffffffffffffffffffffffffffffffffff8416612a3a565b6000610ccd8373ffffffffffffffffffffffffffffffffffffffff84166129eb565b7f00000000000000000000000000000000000000000000000000000000000000001561177f57612557600282612eb0565b61177f576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610996565b6125ae81610cd4565b6125f0576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610996565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612669573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061268d919061415b565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461177f576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610996565b67ffffffffffffffff821660009081526007602052604090206119af90827f0000000000000000000000000000000000000000000000000000000000000000612b2d565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526112e3908490612edf565b6060816000018054806020026020016040519081016040528092919081815260200182805480156106ec57602002820191906000526020600020905b8154815260200190600101908083116128005750505050509050919050565b600061283e8561282f8486614178565b612839908761418f565b612feb565b90505b949350505050565b815460009061287290700100000000000000000000000000000000900463ffffffff1642614089565b9050801561291457600183015483546128ba916fffffffffffffffffffffffffffffffff8082169281169185917001000000000000000000000000000000009091041661281f565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b6020820151835461293a916fffffffffffffffffffffffffffffffff9081169116612feb565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199061212290849061411f565b6000818152600183016020526040812054612a3257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610642565b506000610642565b60008181526001830160205260408120548015612b23576000612a5e600183614089565b8554909150600090612a7290600190614089565b9050808214612ad7576000866000018281548110612a9257612a92613c57565b9060005260206000200154905080876000018481548110612ab557612ab5613c57565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612ae857612ae86141a2565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610642565b6000915050610642565b825474010000000000000000000000000000000000000000900460ff161580612b54575081155b15612b5e57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612ba490700100000000000000000000000000000000900463ffffffff1642614089565b90508015612c645781831115612be6576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612c209083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1661281f565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612d1b5773ffffffffffffffffffffffffffffffffffffffff8416612cc3576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610996565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610996565b84831015612e2e5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612d5f9082614089565b612d69878a614089565b612d73919061418f565b612d7d91906141d1565b905073ffffffffffffffffffffffffffffffffffffffff8616612dd6576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610996565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610996565b612e388584614089565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610ccd565b6000612f41826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166130019092919063ffffffff16565b8051909150156112e35780806020019051810190612f5f9190613b98565b6112e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610996565b6000818310612ffa5781610ccd565b5090919050565b60606128418484600085856000808673ffffffffffffffffffffffffffffffffffffffff168587604051613035919061420c565b60006040518083038185875af1925050503d8060008114613072576040519150601f19603f3d011682016040523d82523d6000602084013e613077565b606091505b509150915061308887838387613093565b979650505050505050565b606083156131295782516000036131225773ffffffffffffffffffffffffffffffffffffffff85163b613122576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610996565b5081612841565b612841838381511561313e5781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161099691906132a8565b50805461317e906137e6565b6000825580601f1061318e575050565b601f01602090049060005260206000209081019061177f91905b808211156131bc57600081556001016131a8565b5090565b6000602082840312156131d257600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610ccd57600080fd5b803567ffffffffffffffff8116811461321a57600080fd5b919050565b60006020828403121561323157600080fd5b610ccd82613202565b60005b8381101561325557818101518382015260200161323d565b50506000910152565b6000815180845261327681602086016020860161323a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610ccd602083018461325e565b73ffffffffffffffffffffffffffffffffffffffff8116811461177f57600080fd5b803561321a816132bb565b6000602082840312156132fa57600080fd5b8135610ccd816132bb565b60006020828403121561331757600080fd5b813567ffffffffffffffff81111561332e57600080fd5b82016101008185031215610ccd57600080fd5b60008083601f84011261335357600080fd5b50813567ffffffffffffffff81111561336b57600080fd5b6020830191508360208260051b850101111561338657600080fd5b9250929050565b600080600080604085870312156133a357600080fd5b843567ffffffffffffffff808211156133bb57600080fd5b6133c788838901613341565b909650945060208701359150808211156133e057600080fd5b506133ed87828801613341565b95989497509550505050565b60008060006040848603121561340e57600080fd5b61341784613202565b9250602084013567ffffffffffffffff8082111561343457600080fd5b818601915086601f83011261344857600080fd5b81358181111561345757600080fd5b87602082850101111561346957600080fd5b6020830194508093505050509250925092565b6000806040838503121561348f57600080fd5b61349883613202565b915060208301356134a8816132bb565b809150509250929050565b6000602082840312156134c557600080fd5b813567ffffffffffffffff8111156134dc57600080fd5b820160a08185031215610ccd57600080fd5b60208152600082516040602084015261350a606084018261325e565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152613545828261325e565b95945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561359c57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161356a565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561359c57835167ffffffffffffffff16835292840192918401916001016135c4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff8111828210171561363d5761363d6135ea565b60405290565b60405160c0810167ffffffffffffffff8111828210171561363d5761363d6135ea565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156136ad576136ad6135ea565b604052919050565b801515811461177f57600080fd5b803561321a816136b5565b80356fffffffffffffffffffffffffffffffff8116811461321a57600080fd5b60006060828403121561370057600080fd5b6040516060810181811067ffffffffffffffff82111715613723576137236135ea565b6040529050808235613734816136b5565b8152613742602084016136ce565b6020820152613753604084016136ce565b60408201525092915050565b600080600060e0848603121561377457600080fd5b61377d84613202565b925061378c85602086016136ee565b915061379b85608086016136ee565b90509250925092565b600080602083850312156137b757600080fd5b823567ffffffffffffffff8111156137ce57600080fd5b6137da85828601613341565b90969095509350505050565b600181811c908216806137fa57607f821691505b602082108103613833577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600067ffffffffffffffff821115613853576138536135ea565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261389057600080fd5b81356138a361389e82613839565b613666565b8181528460208386010111156138b857600080fd5b816020850160208301376000918101602001919091529392505050565b600061010082360312156138e857600080fd5b6138f0613619565b823567ffffffffffffffff8082111561390857600080fd5b6139143683870161387f565b835261392260208601613202565b6020840152613933604086016132dd565b60408401526060850135606084015261394e608086016132dd565b608084015260a085013591508082111561396757600080fd5b6139733683870161387f565b60a084015260c085013591508082111561398c57600080fd5b6139983683870161387f565b60c084015260e08501359150808211156139b157600080fd5b506139be3682860161387f565b60e08301525092915050565b601f8211156112e3576000816000526020600020601f850160051c810160208610156139f35750805b601f850160051c820191505b81811015613a12578281556001016139ff565b505050505050565b67ffffffffffffffff831115613a3257613a326135ea565b613a4683613a4083546137e6565b836139ca565b6000601f841160018114613a985760008515613a625750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355611a4a565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015613ae75786850135825560209485019460019092019101613ac7565b5086821015613b22577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b604081526000613b47604083018661325e565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060208284031215613baa57600080fd5b8151610ccd816136b5565b600060a08236031215613bc757600080fd5b60405160a0810167ffffffffffffffff8282108183111715613beb57613beb6135ea565b816040528435915080821115613c0057600080fd5b50613c0d3682860161387f565b825250613c1c60208401613202565b60208201526040830135613c2f816132bb565b6040820152606083810135908201526080830135613c4c816132bb565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec1833603018112613cba57600080fd5b9190910192915050565b60006101408236031215613cd757600080fd5b613cdf613643565b613ce883613202565b8152613cf6602084016136c3565b6020820152604083013567ffffffffffffffff80821115613d1657600080fd5b613d223683870161387f565b60408401526060850135915080821115613d3b57600080fd5b50613d483682860161387f565b606083015250613d5b36608085016136ee565b6080820152613d6d3660e085016136ee565b60a082015292915050565b815167ffffffffffffffff811115613d9257613d926135ea565b613da681613da084546137e6565b846139ca565b602080601f831160018114613df95760008415613dc35750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613a12565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613e4657888601518255948401946001909101908401613e27565b5085821015613e8257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613eb68184018761325e565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613ef49050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152613545565b60a081526000613f3e60a083018761325e565b73ffffffffffffffffffffffffffffffffffffffff8616602084015284604084015267ffffffffffffffff841660608401528281036080840152600081526020810191505095945050505050565b73ffffffffffffffffffffffffffffffffffffffff8516815260a060208201526000613fbb60a083018661325e565b60408301949094525067ffffffffffffffff9190911660608201528082036080909101526000815260200192915050565b600060208284031215613ffe57600080fd5b815167ffffffffffffffff81111561401557600080fd5b8201601f8101841361402657600080fd5b805161403461389e82613839565b81815285602083850101111561404957600080fd5b61354582602083016020860161323a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156106425761064261405a565b67ffffffffffffffff8416815260e081016140e860208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152612841565b6060810161064282848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561416d57600080fd5b8151610ccd816132bb565b80820281158282048414176106425761064261405a565b808201808211156106425761064261405a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600082614207577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008251613cba81846020870161323a56fe4275726e5769746846726f6d4d696e74546f6b656e506f6f6c416e6450726f787920312e352e30a164736f6c6343000818000a",
}

var BurnWithFromMintTokenPoolAndProxyABI = BurnWithFromMintTokenPoolAndProxyMetaData.ABI

var BurnWithFromMintTokenPoolAndProxyBin = BurnWithFromMintTokenPoolAndProxyMetaData.Bin

func DeployBurnWithFromMintTokenPoolAndProxy(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnWithFromMintTokenPoolAndProxy, error) {
	parsed, err := BurnWithFromMintTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnWithFromMintTokenPoolAndProxyBin), backend, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnWithFromMintTokenPoolAndProxy{address: address, abi: *parsed, BurnWithFromMintTokenPoolAndProxyCaller: BurnWithFromMintTokenPoolAndProxyCaller{contract: contract}, BurnWithFromMintTokenPoolAndProxyTransactor: BurnWithFromMintTokenPoolAndProxyTransactor{contract: contract}, BurnWithFromMintTokenPoolAndProxyFilterer: BurnWithFromMintTokenPoolAndProxyFilterer{contract: contract}}, nil
}

type BurnWithFromMintTokenPoolAndProxy struct {
	address common.Address
	abi     abi.ABI
	BurnWithFromMintTokenPoolAndProxyCaller
	BurnWithFromMintTokenPoolAndProxyTransactor
	BurnWithFromMintTokenPoolAndProxyFilterer
}

type BurnWithFromMintTokenPoolAndProxyCaller struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolAndProxyTransactor struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolAndProxyFilterer struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolAndProxySession struct {
	Contract     *BurnWithFromMintTokenPoolAndProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolAndProxyCallerSession struct {
	Contract *BurnWithFromMintTokenPoolAndProxyCaller
	CallOpts bind.CallOpts
}

type BurnWithFromMintTokenPoolAndProxyTransactorSession struct {
	Contract     *BurnWithFromMintTokenPoolAndProxyTransactor
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolAndProxyRaw struct {
	Contract *BurnWithFromMintTokenPoolAndProxy
}

type BurnWithFromMintTokenPoolAndProxyCallerRaw struct {
	Contract *BurnWithFromMintTokenPoolAndProxyCaller
}

type BurnWithFromMintTokenPoolAndProxyTransactorRaw struct {
	Contract *BurnWithFromMintTokenPoolAndProxyTransactor
}

func NewBurnWithFromMintTokenPoolAndProxy(address common.Address, backend bind.ContractBackend) (*BurnWithFromMintTokenPoolAndProxy, error) {
	abi, err := abi.JSON(strings.NewReader(BurnWithFromMintTokenPoolAndProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnWithFromMintTokenPoolAndProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxy{address: address, abi: abi, BurnWithFromMintTokenPoolAndProxyCaller: BurnWithFromMintTokenPoolAndProxyCaller{contract: contract}, BurnWithFromMintTokenPoolAndProxyTransactor: BurnWithFromMintTokenPoolAndProxyTransactor{contract: contract}, BurnWithFromMintTokenPoolAndProxyFilterer: BurnWithFromMintTokenPoolAndProxyFilterer{contract: contract}}, nil
}

func NewBurnWithFromMintTokenPoolAndProxyCaller(address common.Address, caller bind.ContractCaller) (*BurnWithFromMintTokenPoolAndProxyCaller, error) {
	contract, err := bindBurnWithFromMintTokenPoolAndProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyCaller{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolAndProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnWithFromMintTokenPoolAndProxyTransactor, error) {
	contract, err := bindBurnWithFromMintTokenPoolAndProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyTransactor{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolAndProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnWithFromMintTokenPoolAndProxyFilterer, error) {
	contract, err := bindBurnWithFromMintTokenPoolAndProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyFilterer{contract: contract}, nil
}

func bindBurnWithFromMintTokenPoolAndProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnWithFromMintTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.BurnWithFromMintTokenPoolAndProxyCaller.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.BurnWithFromMintTokenPoolAndProxyTransactor.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.BurnWithFromMintTokenPoolAndProxyTransactor.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetAllowList(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetAllowList(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getOnRamp", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetOnRamp(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, arg0)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetOnRamp(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, arg0)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetPreviousPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getPreviousPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetPreviousPool() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetPreviousPool(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetPreviousPool() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetPreviousPool(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRemotePool(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRemotePool(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRemoteToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRemoteToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRmnProxy(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRmnProxy(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRouter(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetRouter(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetSupportedChains(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetSupportedChains(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.GetToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "isOffRamp", sourceChainSelector, offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsOffRamp(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsOffRamp(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsSupportedChain(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsSupportedChain(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsSupportedToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, token)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.IsSupportedToken(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, token)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.Owner(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.Owner(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SupportsInterface(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SupportsInterface(&_BurnWithFromMintTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPoolAndProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.TypeAndVersion(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyCallerSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.TypeAndVersion(&_BurnWithFromMintTokenPoolAndProxy.CallOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "acceptOwnership")
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.AcceptOwnership(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.AcceptOwnership(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, chains)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, chains)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.LockOrBurn(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.LockOrBurn(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) SetPreviousPool(opts *bind.TransactOpts, prevPool common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "setPreviousPool", prevPool)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetPreviousPool(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetPreviousPool(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRemotePool(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRemotePool(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRouter(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.SetRouter(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.TransferOwnership(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, to)
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPoolAndProxy.Contract.TransferOwnership(&_BurnWithFromMintTokenPoolAndProxy.TransactOpts, to)
}

type BurnWithFromMintTokenPoolAndProxyAllowListAddIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyAllowListAdd)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyAllowListAdd)
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

func (it *BurnWithFromMintTokenPoolAndProxyAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyAllowListAddIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyAllowListAddIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyAllowListAdd)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAndProxyAllowListAdd, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyAllowListAdd)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyAllowListRemove)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyAllowListRemove)
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

func (it *BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyAllowListRemove)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAndProxyAllowListRemove, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyAllowListRemove)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyBurnedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyBurned)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyBurned)
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

func (it *BurnWithFromMintTokenPoolAndProxyBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolAndProxyBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyBurnedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyBurned)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolAndProxyBurned, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyBurned)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyChainAddedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyChainAdded)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyChainAdded)
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

func (it *BurnWithFromMintTokenPoolAndProxyChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainAddedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyChainAddedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyChainAdded)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainAdded, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyChainAdded)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyChainConfigured)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyChainConfigured)
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

func (it *BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyChainConfigured)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainConfigured, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyChainConfigured)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyChainRemovedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyChainRemoved)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyChainRemoved)
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

func (it *BurnWithFromMintTokenPoolAndProxyChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainRemovedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyChainRemovedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyChainRemoved)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainRemoved, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyChainRemoved)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyConfigChangedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyConfigChanged)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyConfigChanged)
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

func (it *BurnWithFromMintTokenPoolAndProxyConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyConfigChangedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyConfigChangedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyConfigChanged)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolAndProxyConfigChanged, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyConfigChanged)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged)
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

func (it *BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged struct {
	OldPool common.Address
	NewPool common.Address
	Raw     types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterLegacyPoolChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "LegacyPoolChanged", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseLegacyPoolChanged(log types.Log) (*BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyLockedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyLocked)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyLocked)
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

func (it *BurnWithFromMintTokenPoolAndProxyLockedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolAndProxyLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyLockedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyLocked)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolAndProxyLocked, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyLocked)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyMintedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyMinted)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyMinted)
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

func (it *BurnWithFromMintTokenPoolAndProxyMintedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolAndProxyMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyMintedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyMinted)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolAndProxyMinted, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyMinted)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested)
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

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferred)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferred)
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

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferred)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferred, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyOwnershipTransferred)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyReleasedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyReleased)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyReleased)
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

func (it *BurnWithFromMintTokenPoolAndProxyReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolAndProxyReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyReleasedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyReleased)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolAndProxyReleased, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyReleased)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyRemotePoolSet)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyRemotePoolSet)
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

func (it *BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyRemotePoolSet)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseRemotePoolSet(log types.Log) (*BurnWithFromMintTokenPoolAndProxyRemotePoolSet, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyRemotePoolSet)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyRouterUpdated)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyRouterUpdated)
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

func (it *BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyRouterUpdated)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolAndProxyRouterUpdated, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyRouterUpdated)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator struct {
	Event *BurnWithFromMintTokenPoolAndProxyTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAndProxyTokensConsumed)
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
		it.Event = new(BurnWithFromMintTokenPoolAndProxyTokensConsumed)
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

func (it *BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAndProxyTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator{contract: _BurnWithFromMintTokenPoolAndProxy.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPoolAndProxy.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAndProxyTokensConsumed)
				if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxyFilterer) ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolAndProxyTokensConsumed, error) {
	event := new(BurnWithFromMintTokenPoolAndProxyTokensConsumed)
	if err := _BurnWithFromMintTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["AllowListAdd"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseAllowListAdd(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["AllowListRemove"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseAllowListRemove(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["Burned"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseBurned(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["ChainAdded"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseChainAdded(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["ChainConfigured"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseChainConfigured(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["ChainRemoved"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseChainRemoved(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["ConfigChanged"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseConfigChanged(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["LegacyPoolChanged"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseLegacyPoolChanged(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["Locked"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseLocked(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["Minted"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseMinted(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseOwnershipTransferRequested(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["OwnershipTransferred"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseOwnershipTransferred(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["Released"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseReleased(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["RemotePoolSet"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseRemotePoolSet(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["RouterUpdated"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseRouterUpdated(log)
	case _BurnWithFromMintTokenPoolAndProxy.abi.Events["TokensConsumed"].ID:
		return _BurnWithFromMintTokenPoolAndProxy.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnWithFromMintTokenPoolAndProxyAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnWithFromMintTokenPoolAndProxyAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnWithFromMintTokenPoolAndProxyBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnWithFromMintTokenPoolAndProxyChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnWithFromMintTokenPoolAndProxyChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnWithFromMintTokenPoolAndProxyChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnWithFromMintTokenPoolAndProxyConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged) Topic() common.Hash {
	return common.HexToHash("0x81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f228")
}

func (BurnWithFromMintTokenPoolAndProxyLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnWithFromMintTokenPoolAndProxyMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnWithFromMintTokenPoolAndProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnWithFromMintTokenPoolAndProxyReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnWithFromMintTokenPoolAndProxyRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (BurnWithFromMintTokenPoolAndProxyRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnWithFromMintTokenPoolAndProxyTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnWithFromMintTokenPoolAndProxy *BurnWithFromMintTokenPoolAndProxy) Address() common.Address {
	return _BurnWithFromMintTokenPoolAndProxy.address
}

type BurnWithFromMintTokenPoolAndProxyInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error)

	GetPreviousPool(opts *bind.CallOpts) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

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

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetPreviousPool(opts *bind.TransactOpts, prevPool common.Address) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAndProxyAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAndProxyAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolAndProxyBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolAndProxyBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolAndProxyChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolAndProxyConfigChanged, error)

	FilterLegacyPoolChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyLegacyPoolChangedIterator, error)

	WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error)

	ParseLegacyPoolChanged(log types.Log) (*BurnWithFromMintTokenPoolAndProxyLegacyPoolChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolAndProxyLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolAndProxyLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolAndProxyMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolAndProxyMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolAndProxyOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolAndProxyReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolAndProxyReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolAndProxyRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*BurnWithFromMintTokenPoolAndProxyRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolAndProxyRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAndProxyTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAndProxyTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolAndProxyTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
