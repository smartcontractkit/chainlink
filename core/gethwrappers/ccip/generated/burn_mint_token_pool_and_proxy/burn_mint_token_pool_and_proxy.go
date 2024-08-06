// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_token_pool_and_proxy

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

var BurnMintTokenPoolAndProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"oldPool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"LegacyPoolChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"onRampAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"prevPool\",\"type\":\"address\"}],\"name\":\"setPreviousPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200492538038062004925833981016040819052620000349162000554565b83838383838383833380600081620000935760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c657620000c6816200017e565b5050506001600160a01b0384161580620000e757506001600160a01b038116155b80620000fa57506001600160a01b038216155b1562000119576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200016c576040805160008152602081019091526200016c908462000229565b505050505050505050505050620006b2565b336001600160a01b03821603620001d85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c0516200024a576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002d55760008382815181106200026e576200026e62000664565b602090810291909101015190506200028860028262000386565b15620002cb576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016200024d565b5060005b815181101562000381576000828281518110620002fa57620002fa62000664565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000326575062000378565b62000333600282620003a6565b1562000376576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002d9565b505050565b60006200039d836001600160a01b038416620003bd565b90505b92915050565b60006200039d836001600160a01b038416620004c1565b60008181526001830160205260408120548015620004b6576000620003e46001836200067a565b8554909150600090620003fa906001906200067a565b9050818114620004665760008660000182815481106200041e576200041e62000664565b906000526020600020015490508087600001848154811062000444576200044462000664565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200047a576200047a6200069c565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a0565b6000915050620003a0565b60008181526001830160205260408120546200050a57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a0565b506000620003a0565b6001600160a01b03811681146200052957600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b80516200054f8162000513565b919050565b600080600080608085870312156200056b57600080fd5b8451620005788162000513565b602086810151919550906001600160401b03808211156200059857600080fd5b818801915088601f830112620005ad57600080fd5b815181811115620005c257620005c26200052c565b8060051b604051601f19603f83011681018181108582111715620005ea57620005ea6200052c565b60405291825284820192508381018501918b8311156200060957600080fd5b938501935b828510156200063257620006228562000542565b845293850193928501926200060e565b809850505050505050620006496040860162000542565b9150620006596060860162000542565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003a057634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c0516141ef620007366000396000818161050301528181611a7701526124c90152600081816104dd0152818161180f0152611d2a0152600081816102260152818161027b0152818161076001528181610d9a0152818161172f01528181611c4a01528181611e300152818161245f01526126b401526141ef6000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80639766b93211610104578063c0d78655116100a2578063db6327dc11610071578063db6327dc146104c8578063dc0bd971146104db578063e0351e1314610501578063f2fde38b1461052757600080fd5b8063c0d786551461047a578063c4bffe2b1461048d578063c75eea9c146104a2578063cf7401f3146104b557600080fd5b8063a8d87a3b116100de578063a8d87a3b146103c7578063af58d59f146103da578063b0f479a114610449578063b79465801461046757600080fd5b80639766b9321461037f5780639a4575b914610392578063a7cd63b7146103b257600080fd5b80636d3d1a58116101715780637d54534e1161014b5780637d54534e1461032857806383826b2b1461033b5780638926f54f1461034e5780638da5cb5b1461036157600080fd5b80636d3d1a58146102ef57806378a010b21461030d57806379ba50971461032057600080fd5b806321df0da7116101ad57806321df0da714610224578063240028e81461026b57806339077537146102b857806354c8a4f3146102da57600080fd5b806301ffc9a7146101d45780630a2fd493146101fc578063181f5a771461021c575b600080fd5b6101e76101e2366004613161565b61053a565b60405190151581526020015b60405180910390f35b61020f61020a3660046131c0565b61061f565b6040516101f39190613249565b61020f6106cf565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b6101e7610279366004613289565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102cb6102c63660046132a6565b6106eb565b604051905181526020016101f3565b6102ed6102e836600461332e565b610878565b005b60085473ffffffffffffffffffffffffffffffffffffffff16610246565b6102ed61031b36600461339a565b6108f3565b6102ed610a67565b6102ed610336366004613289565b610b64565b6101e761034936600461341d565b610bb3565b6101e761035c3660046131c0565b610c80565b60005473ffffffffffffffffffffffffffffffffffffffff16610246565b6102ed61038d366004613289565b610c97565b6103a56103a0366004613454565b610d26565b6040516101f3919061348f565b6103ba610e96565b6040516101f391906134ef565b6102466103d53660046131c0565b503090565b6103ed6103e83660046131c0565b610ea7565b6040516101f3919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610246565b61020f6104753660046131c0565b610f7c565b6102ed610488366004613289565b610fa7565b61049561107b565b6040516101f39190613549565b6103ed6104b03660046131c0565b611133565b6102ed6104c3366004613700565b611205565b6102ed6104d6366004613745565b61128e565b7f0000000000000000000000000000000000000000000000000000000000000000610246565b7f00000000000000000000000000000000000000000000000000000000000000006101e7565b6102ed610535366004613289565b611714565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf0000000000000000000000000000000000000000000000000000000014806105cd57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061061957507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061064a90613787565b80601f016020809104026020016040519081016040528092919081815260200182805461067690613787565b80156106c35780601f10610698576101008083540402835291602001916106c3565b820191906000526020600020905b8154815290600101906020018083116106a657829003601f168201915b50505050509050919050565b6040518060600160405280602381526020016141c06023913981565b60408051602081019091526000815261070b61070683613876565b611728565b60095473ffffffffffffffffffffffffffffffffffffffff166107d6576040517f40c10f19000000000000000000000000000000000000000000000000000000008152336004820152606083013560248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906340c10f1990604401600060405180830381600087803b1580156107b957600080fd5b505af11580156107cd573d6000803e3d6000fd5b505050506107e7565b6107e76107e283613876565b611959565b6107f76060830160408401613289565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0846060013560405161085991815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6108806119f2565b6108ed84848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611a7592505050565b50505050565b6108fb6119f2565b61090483610c80565b61094b576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b67ffffffffffffffff83166000908152600760205260408120600401805461097290613787565b80601f016020809104026020016040519081016040528092919081815260200182805461099e90613787565b80156109eb5780601f106109c0576101008083540402835291602001916109eb565b820191906000526020600020905b8154815290600101906020018083116109ce57829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610a1a8385836139bb565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610a5993929190613ad5565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ae8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610942565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b6c6119f2565b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600073ffffffffffffffffffffffffffffffffffffffff8216301480610c795750600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169281019290925273ffffffffffffffffffffffffffffffffffffffff848116602484015216906383826b2b90604401602060405180830381865afa158015610c55573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c799190613b39565b9392505050565b6000610619600567ffffffffffffffff8416611c2b565b610c9f6119f2565b6009805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f22891015b60405180910390a15050565b6040805180820190915260608082526020820152610d4b610d4683613b56565b611c43565b60095473ffffffffffffffffffffffffffffffffffffffff16610e10576040517f42966c68000000000000000000000000000000000000000000000000000000008152606083013560048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906342966c6890602401600060405180830381600087803b158015610df357600080fd5b505af1158015610e07573d6000803e3d6000fd5b50505050610e21565b610e21610e1c83613b56565b611e0d565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610e7b84602001602081019061047591906131c0565b81526040805160208181019092526000815291015292915050565b6060610ea26002611f27565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261061990611f34565b67ffffffffffffffff8116600090815260076020526040902060050180546060919061064a90613787565b610faf6119f2565b73ffffffffffffffffffffffffffffffffffffffff8116610ffc576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849101610d1a565b606060006110896005611f27565b90506000815167ffffffffffffffff8111156110a7576110a761358b565b6040519080825280602002602001820160405280156110d0578160200160208202803683370190505b50905060005b825181101561112c578281815181106110f1576110f1613bf8565b602002602001015182828151811061110b5761110b613bf8565b67ffffffffffffffff909216602092830291909101909101526001016110d6565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261061990611f34565b60085473ffffffffffffffffffffffffffffffffffffffff163314801590611245575060005473ffffffffffffffffffffffffffffffffffffffff163314155b1561127e576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610942565b611289838383611fe6565b505050565b6112966119f2565b60005b818110156112895760008383838181106112b5576112b5613bf8565b90506020028101906112c79190613c27565b6112d090613c65565b90506112e581608001518260200151156120d0565b6112f88160a001518260200151156120d0565b8060200151156115f457805161131a9060059067ffffffffffffffff16612209565b61135f5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610942565b60408101515115806113745750606081015151155b156113ab576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c1790911696151502959095179098559081015194015193811693169091029190911760038201559151909190600482019061158c9082613d19565b50606082015160058201906115a19082613d19565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506115e79493929190613e33565b60405180910390a161170b565b805161160c9060059067ffffffffffffffff16612215565b6116515780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610942565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906116ba6004830182613113565b6116c8600583016000613113565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b50600101611299565b61171c6119f2565b61172581612221565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff9081169116146117bd5760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610942565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa15801561186b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061188f9190613b39565b156118c6576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6118d38160200151612316565b60006118e2826020015161061f565b9050805160001480611906575080805190602001208260a001518051906020012014155b15611943578160a001516040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016109429190613249565b6119558260200151836060015161243c565b5050565b6009548151606083015160208401516040517f8627fad600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90941693638627fad6936119bd9390923392600401613ecc565b600060405180830381600087803b1580156119d757600080fd5b505af11580156119eb573d6000803e3d6000fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611a73576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610942565b565b7f0000000000000000000000000000000000000000000000000000000000000000611acc576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611b62576000838281518110611aec57611aec613bf8565b60200260200101519050611b0a81600261248390919063ffffffff16565b15611b595760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611acf565b5060005b8151811015611289576000828281518110611b8357611b83613bf8565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611bc75750611c23565b611bd26002826124a5565b15611c215760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611b66565b60008181526001830160205260408120541515610c79565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611cd85760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610942565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611d86573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611daa9190613b39565b15611de1576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611dee81604001516124c7565b611dfb8160200151612546565b61172581602001518260600151612694565b6009546060820151611e5a9173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116929116906126d8565b60095460408083015183516060850151602086015193517f9687544500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90951694639687544594611ec294939291600401613f2d565b6000604051808303816000875af1158015611ee1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526119559190810190613f8d565b60606000610c7983612765565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611fc282606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611fa6919061402a565b85608001516fffffffffffffffffffffffffffffffff166127c0565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b611fef83610c80565b612031576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610942565b61203c8260006120d0565b67ffffffffffffffff8316600090815260076020526040902061205f90836127ea565b61206a8160006120d0565b67ffffffffffffffff8316600090815260076020526040902061209090600201826127ea565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516120c39392919061403d565b60405180910390a1505050565b8151156121975781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612126575060408201516fffffffffffffffffffffffffffffffff16155b1561215f57816040517f8020d12400000000000000000000000000000000000000000000000000000000815260040161094291906140c0565b8015611955576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff161515806121d0575060208201516fffffffffffffffffffffffffffffffff1615155b1561195557816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161094291906140c0565b6000610c79838361298c565b6000610c7983836129db565b3373ffffffffffffffffffffffffffffffffffffffff8216036122a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610942565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61231f81610c80565b612361576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610942565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156123e0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124049190613b39565b611725576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610942565b67ffffffffffffffff8216600090815260076020526040902061195590600201827f0000000000000000000000000000000000000000000000000000000000000000612ace565b6000610c798373ffffffffffffffffffffffffffffffffffffffff84166129db565b6000610c798373ffffffffffffffffffffffffffffffffffffffff841661298c565b7f000000000000000000000000000000000000000000000000000000000000000015611725576124f8600282612e51565b611725576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610942565b61254f81610c80565b612591576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610942565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa15801561260a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061262e91906140fc565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611725576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610942565b67ffffffffffffffff8216600090815260076020526040902061195590827f0000000000000000000000000000000000000000000000000000000000000000612ace565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052611289908490612e80565b6060816000018054806020026020016040519081016040528092919081815260200182805480156106c357602002820191906000526020600020905b8154815260200190600101908083116127a15750505050509050919050565b60006127df856127d08486614119565b6127da9087614130565b612f8c565b90505b949350505050565b815460009061281390700100000000000000000000000000000000900463ffffffff164261402a565b905080156128b5576001830154835461285b916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166127c0565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546128db916fffffffffffffffffffffffffffffffff9081169116612f8c565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906120c39084906140c0565b60008181526001830160205260408120546129d357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610619565b506000610619565b60008181526001830160205260408120548015612ac45760006129ff60018361402a565b8554909150600090612a139060019061402a565b9050818114612a78576000866000018281548110612a3357612a33613bf8565b9060005260206000200154905080876000018481548110612a5657612a56613bf8565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612a8957612a89614143565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610619565b6000915050610619565b825474010000000000000000000000000000000000000000900460ff161580612af5575081155b15612aff57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612b4590700100000000000000000000000000000000900463ffffffff164261402a565b90508015612c055781831115612b87576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612bc19083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166127c0565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612cbc5773ffffffffffffffffffffffffffffffffffffffff8416612c64576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610942565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610942565b84831015612dcf5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612d00908261402a565b612d0a878a61402a565b612d149190614130565b612d1e9190614172565b905073ffffffffffffffffffffffffffffffffffffffff8616612d77576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610942565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610942565b612dd9858461402a565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610c79565b6000612ee2826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16612fa29092919063ffffffff16565b8051909150156112895780806020019051810190612f009190613b39565b611289576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610942565b6000818310612f9b5781610c79565b5090919050565b60606127e28484600085856000808673ffffffffffffffffffffffffffffffffffffffff168587604051612fd691906141ad565b60006040518083038185875af1925050503d8060008114613013576040519150601f19603f3d011682016040523d82523d6000602084013e613018565b606091505b509150915061302987838387613034565b979650505050505050565b606083156130ca5782516000036130c35773ffffffffffffffffffffffffffffffffffffffff85163b6130c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610942565b50816127e2565b6127e283838151156130df5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109429190613249565b50805461311f90613787565b6000825580601f1061312f575050565b601f01602090049060005260206000209081019061172591905b8082111561315d5760008155600101613149565b5090565b60006020828403121561317357600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610c7957600080fd5b803567ffffffffffffffff811681146131bb57600080fd5b919050565b6000602082840312156131d257600080fd5b610c79826131a3565b60005b838110156131f65781810151838201526020016131de565b50506000910152565b600081518084526132178160208601602086016131db565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610c7960208301846131ff565b73ffffffffffffffffffffffffffffffffffffffff8116811461172557600080fd5b80356131bb8161325c565b60006020828403121561329b57600080fd5b8135610c798161325c565b6000602082840312156132b857600080fd5b813567ffffffffffffffff8111156132cf57600080fd5b82016101008185031215610c7957600080fd5b60008083601f8401126132f457600080fd5b50813567ffffffffffffffff81111561330c57600080fd5b6020830191508360208260051b850101111561332757600080fd5b9250929050565b6000806000806040858703121561334457600080fd5b843567ffffffffffffffff8082111561335c57600080fd5b613368888389016132e2565b9096509450602087013591508082111561338157600080fd5b5061338e878288016132e2565b95989497509550505050565b6000806000604084860312156133af57600080fd5b6133b8846131a3565b9250602084013567ffffffffffffffff808211156133d557600080fd5b818601915086601f8301126133e957600080fd5b8135818111156133f857600080fd5b87602082850101111561340a57600080fd5b6020830194508093505050509250925092565b6000806040838503121561343057600080fd5b613439836131a3565b915060208301356134498161325c565b809150509250929050565b60006020828403121561346657600080fd5b813567ffffffffffffffff81111561347d57600080fd5b820160a08185031215610c7957600080fd5b6020815260008251604060208401526134ab60608401826131ff565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526134e682826131ff565b95945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561353d57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161350b565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561353d57835167ffffffffffffffff1683529284019291840191600101613565565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff811182821017156135de576135de61358b565b60405290565b60405160c0810167ffffffffffffffff811182821017156135de576135de61358b565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561364e5761364e61358b565b604052919050565b801515811461172557600080fd5b80356131bb81613656565b80356fffffffffffffffffffffffffffffffff811681146131bb57600080fd5b6000606082840312156136a157600080fd5b6040516060810181811067ffffffffffffffff821117156136c4576136c461358b565b60405290508082356136d581613656565b81526136e36020840161366f565b60208201526136f46040840161366f565b60408201525092915050565b600080600060e0848603121561371557600080fd5b61371e846131a3565b925061372d856020860161368f565b915061373c856080860161368f565b90509250925092565b6000806020838503121561375857600080fd5b823567ffffffffffffffff81111561376f57600080fd5b61377b858286016132e2565b90969095509350505050565b600181811c9082168061379b57607f821691505b6020821081036137d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600067ffffffffffffffff8211156137f4576137f461358b565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261383157600080fd5b813561384461383f826137da565b613607565b81815284602083860101111561385957600080fd5b816020850160208301376000918101602001919091529392505050565b6000610100823603121561388957600080fd5b6138916135ba565b823567ffffffffffffffff808211156138a957600080fd5b6138b536838701613820565b83526138c3602086016131a3565b60208401526138d46040860161327e565b6040840152606085013560608401526138ef6080860161327e565b608084015260a085013591508082111561390857600080fd5b61391436838701613820565b60a084015260c085013591508082111561392d57600080fd5b61393936838701613820565b60c084015260e085013591508082111561395257600080fd5b5061395f36828601613820565b60e08301525092915050565b601f821115611289576000816000526020600020601f850160051c810160208610156139945750805b601f850160051c820191505b818110156139b3578281556001016139a0565b505050505050565b67ffffffffffffffff8311156139d3576139d361358b565b6139e7836139e18354613787565b8361396b565b6000601f841160018114613a395760008515613a035750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556119eb565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015613a885786850135825560209485019460019092019101613a68565b5086821015613ac3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b604081526000613ae860408301866131ff565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060208284031215613b4b57600080fd5b8151610c7981613656565b600060a08236031215613b6857600080fd5b60405160a0810167ffffffffffffffff8282108183111715613b8c57613b8c61358b565b816040528435915080821115613ba157600080fd5b50613bae36828601613820565b825250613bbd602084016131a3565b60208201526040830135613bd08161325c565b6040820152606083810135908201526080830135613bed8161325c565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec1833603018112613c5b57600080fd5b9190910192915050565b60006101408236031215613c7857600080fd5b613c806135e4565b613c89836131a3565b8152613c9760208401613664565b6020820152604083013567ffffffffffffffff80821115613cb757600080fd5b613cc336838701613820565b60408401526060850135915080821115613cdc57600080fd5b50613ce936828601613820565b606083015250613cfc366080850161368f565b6080820152613d0e3660e0850161368f565b60a082015292915050565b815167ffffffffffffffff811115613d3357613d3361358b565b613d4781613d418454613787565b8461396b565b602080601f831160018114613d9a5760008415613d645750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556139b3565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613de757888601518255948401946001909101908401613dc8565b5085821015613e2357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613e57818401876131ff565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613e959050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e08301526134e6565b60a081526000613edf60a08301876131ff565b73ffffffffffffffffffffffffffffffffffffffff8616602084015284604084015267ffffffffffffffff841660608401528281036080840152600081526020810191505095945050505050565b73ffffffffffffffffffffffffffffffffffffffff8516815260a060208201526000613f5c60a08301866131ff565b60408301949094525067ffffffffffffffff9190911660608201528082036080909101526000815260200192915050565b600060208284031215613f9f57600080fd5b815167ffffffffffffffff811115613fb657600080fd5b8201601f81018413613fc757600080fd5b8051613fd561383f826137da565b818152856020838501011115613fea57600080fd5b6134e68260208301602086016131db565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561061957610619613ffb565b67ffffffffffffffff8416815260e0810161408960208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526127e2565b6060810161061982848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561410e57600080fd5b8151610c798161325c565b808202811582820484141761061957610619613ffb565b8082018082111561061957610619613ffb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000826141a8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008251613c5b8184602087016131db56fe4275726e4d696e74546f6b656e506f6f6c416e6450726f787920312e352e302d646576a164736f6c6343000818000a",
}

var BurnMintTokenPoolAndProxyABI = BurnMintTokenPoolAndProxyMetaData.ABI

var BurnMintTokenPoolAndProxyBin = BurnMintTokenPoolAndProxyMetaData.Bin

func DeployBurnMintTokenPoolAndProxy(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnMintTokenPoolAndProxy, error) {
	parsed, err := BurnMintTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintTokenPoolAndProxyBin), backend, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintTokenPoolAndProxy{address: address, abi: *parsed, BurnMintTokenPoolAndProxyCaller: BurnMintTokenPoolAndProxyCaller{contract: contract}, BurnMintTokenPoolAndProxyTransactor: BurnMintTokenPoolAndProxyTransactor{contract: contract}, BurnMintTokenPoolAndProxyFilterer: BurnMintTokenPoolAndProxyFilterer{contract: contract}}, nil
}

type BurnMintTokenPoolAndProxy struct {
	address common.Address
	abi     abi.ABI
	BurnMintTokenPoolAndProxyCaller
	BurnMintTokenPoolAndProxyTransactor
	BurnMintTokenPoolAndProxyFilterer
}

type BurnMintTokenPoolAndProxyCaller struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolAndProxyTransactor struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolAndProxyFilterer struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolAndProxySession struct {
	Contract     *BurnMintTokenPoolAndProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolAndProxyCallerSession struct {
	Contract *BurnMintTokenPoolAndProxyCaller
	CallOpts bind.CallOpts
}

type BurnMintTokenPoolAndProxyTransactorSession struct {
	Contract     *BurnMintTokenPoolAndProxyTransactor
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolAndProxyRaw struct {
	Contract *BurnMintTokenPoolAndProxy
}

type BurnMintTokenPoolAndProxyCallerRaw struct {
	Contract *BurnMintTokenPoolAndProxyCaller
}

type BurnMintTokenPoolAndProxyTransactorRaw struct {
	Contract *BurnMintTokenPoolAndProxyTransactor
}

func NewBurnMintTokenPoolAndProxy(address common.Address, backend bind.ContractBackend) (*BurnMintTokenPoolAndProxy, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintTokenPoolAndProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintTokenPoolAndProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxy{address: address, abi: abi, BurnMintTokenPoolAndProxyCaller: BurnMintTokenPoolAndProxyCaller{contract: contract}, BurnMintTokenPoolAndProxyTransactor: BurnMintTokenPoolAndProxyTransactor{contract: contract}, BurnMintTokenPoolAndProxyFilterer: BurnMintTokenPoolAndProxyFilterer{contract: contract}}, nil
}

func NewBurnMintTokenPoolAndProxyCaller(address common.Address, caller bind.ContractCaller) (*BurnMintTokenPoolAndProxyCaller, error) {
	contract, err := bindBurnMintTokenPoolAndProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyCaller{contract: contract}, nil
}

func NewBurnMintTokenPoolAndProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintTokenPoolAndProxyTransactor, error) {
	contract, err := bindBurnMintTokenPoolAndProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyTransactor{contract: contract}, nil
}

func NewBurnMintTokenPoolAndProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintTokenPoolAndProxyFilterer, error) {
	contract, err := bindBurnMintTokenPoolAndProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyFilterer{contract: contract}, nil
}

func bindBurnMintTokenPoolAndProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintTokenPoolAndProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPoolAndProxy.Contract.BurnMintTokenPoolAndProxyCaller.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.BurnMintTokenPoolAndProxyTransactor.contract.Transfer(opts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.BurnMintTokenPoolAndProxyTransactor.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPoolAndProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.contract.Transfer(opts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetAllowList(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetAllowList(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetAllowListEnabled(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetAllowListEnabled(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getOnRamp", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetOnRamp(&_BurnMintTokenPoolAndProxy.CallOpts, arg0)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetOnRamp(arg0 uint64) (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetOnRamp(&_BurnMintTokenPoolAndProxy.CallOpts, arg0)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRateLimitAdmin(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRemotePool(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRemotePool(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRemoteToken(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRemoteToken(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRmnProxy(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRmnProxy(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRouter(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetRouter(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetSupportedChains(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetSupportedChains(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) GetToken() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetToken(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) GetToken() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.GetToken(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) IsOffRamp(opts *bind.CallOpts, sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "isOffRamp", sourceChainSelector, offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsOffRamp(&_BurnMintTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) IsOffRamp(sourceChainSelector uint64, offRamp common.Address) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsOffRamp(&_BurnMintTokenPoolAndProxy.CallOpts, sourceChainSelector, offRamp)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsSupportedChain(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsSupportedChain(&_BurnMintTokenPoolAndProxy.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsSupportedToken(&_BurnMintTokenPoolAndProxy.CallOpts, token)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.IsSupportedToken(&_BurnMintTokenPoolAndProxy.CallOpts, token)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) Owner() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.Owner(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) Owner() (common.Address, error) {
	return _BurnMintTokenPoolAndProxy.Contract.Owner(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SupportsInterface(&_BurnMintTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SupportsInterface(&_BurnMintTokenPoolAndProxy.CallOpts, interfaceId)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintTokenPoolAndProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPoolAndProxy.Contract.TypeAndVersion(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyCallerSession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPoolAndProxy.Contract.TypeAndVersion(&_BurnMintTokenPoolAndProxy.CallOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.AcceptOwnership(&_BurnMintTokenPoolAndProxy.TransactOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.AcceptOwnership(&_BurnMintTokenPoolAndProxy.TransactOpts)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_BurnMintTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ApplyAllowListUpdates(&_BurnMintTokenPoolAndProxy.TransactOpts, removes, adds)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ApplyChainUpdates(&_BurnMintTokenPoolAndProxy.TransactOpts, chains)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ApplyChainUpdates(&_BurnMintTokenPoolAndProxy.TransactOpts, chains)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.LockOrBurn(&_BurnMintTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.LockOrBurn(&_BurnMintTokenPoolAndProxy.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ReleaseOrMint(&_BurnMintTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.ReleaseOrMint(&_BurnMintTokenPoolAndProxy.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) SetPreviousPool(opts *bind.TransactOpts, prevPool common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "setPreviousPool", prevPool)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetPreviousPool(&_BurnMintTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) SetPreviousPool(prevPool common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetPreviousPool(&_BurnMintTokenPoolAndProxy.TransactOpts, prevPool)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_BurnMintTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRateLimitAdmin(&_BurnMintTokenPoolAndProxy.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRemotePool(&_BurnMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRemotePool(&_BurnMintTokenPoolAndProxy.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRouter(&_BurnMintTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.SetRouter(&_BurnMintTokenPoolAndProxy.TransactOpts, newRouter)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.TransferOwnership(&_BurnMintTokenPoolAndProxy.TransactOpts, to)
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPoolAndProxy.Contract.TransferOwnership(&_BurnMintTokenPoolAndProxy.TransactOpts, to)
}

type BurnMintTokenPoolAndProxyAllowListAddIterator struct {
	Event *BurnMintTokenPoolAndProxyAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyAllowListAdd)
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
		it.Event = new(BurnMintTokenPoolAndProxyAllowListAdd)
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

func (it *BurnMintTokenPoolAndProxyAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyAllowListAddIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyAllowListAddIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyAllowListAdd)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAndProxyAllowListAdd, error) {
	event := new(BurnMintTokenPoolAndProxyAllowListAdd)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyAllowListRemoveIterator struct {
	Event *BurnMintTokenPoolAndProxyAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyAllowListRemove)
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
		it.Event = new(BurnMintTokenPoolAndProxyAllowListRemove)
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

func (it *BurnMintTokenPoolAndProxyAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyAllowListRemoveIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyAllowListRemoveIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyAllowListRemove)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAndProxyAllowListRemove, error) {
	event := new(BurnMintTokenPoolAndProxyAllowListRemove)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyBurnedIterator struct {
	Event *BurnMintTokenPoolAndProxyBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyBurned)
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
		it.Event = new(BurnMintTokenPoolAndProxyBurned)
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

func (it *BurnMintTokenPoolAndProxyBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolAndProxyBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyBurnedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyBurned)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseBurned(log types.Log) (*BurnMintTokenPoolAndProxyBurned, error) {
	event := new(BurnMintTokenPoolAndProxyBurned)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyChainAddedIterator struct {
	Event *BurnMintTokenPoolAndProxyChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyChainAdded)
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
		it.Event = new(BurnMintTokenPoolAndProxyChainAdded)
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

func (it *BurnMintTokenPoolAndProxyChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainAddedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyChainAddedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyChainAdded)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseChainAdded(log types.Log) (*BurnMintTokenPoolAndProxyChainAdded, error) {
	event := new(BurnMintTokenPoolAndProxyChainAdded)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyChainConfiguredIterator struct {
	Event *BurnMintTokenPoolAndProxyChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyChainConfigured)
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
		it.Event = new(BurnMintTokenPoolAndProxyChainConfigured)
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

func (it *BurnMintTokenPoolAndProxyChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainConfiguredIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyChainConfiguredIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyChainConfigured)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseChainConfigured(log types.Log) (*BurnMintTokenPoolAndProxyChainConfigured, error) {
	event := new(BurnMintTokenPoolAndProxyChainConfigured)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyChainRemovedIterator struct {
	Event *BurnMintTokenPoolAndProxyChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyChainRemoved)
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
		it.Event = new(BurnMintTokenPoolAndProxyChainRemoved)
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

func (it *BurnMintTokenPoolAndProxyChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainRemovedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyChainRemovedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyChainRemoved)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseChainRemoved(log types.Log) (*BurnMintTokenPoolAndProxyChainRemoved, error) {
	event := new(BurnMintTokenPoolAndProxyChainRemoved)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyConfigChangedIterator struct {
	Event *BurnMintTokenPoolAndProxyConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyConfigChanged)
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
		it.Event = new(BurnMintTokenPoolAndProxyConfigChanged)
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

func (it *BurnMintTokenPoolAndProxyConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyConfigChangedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyConfigChangedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyConfigChanged)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseConfigChanged(log types.Log) (*BurnMintTokenPoolAndProxyConfigChanged, error) {
	event := new(BurnMintTokenPoolAndProxyConfigChanged)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyLegacyPoolChangedIterator struct {
	Event *BurnMintTokenPoolAndProxyLegacyPoolChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyLegacyPoolChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyLegacyPoolChanged)
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
		it.Event = new(BurnMintTokenPoolAndProxyLegacyPoolChanged)
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

func (it *BurnMintTokenPoolAndProxyLegacyPoolChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyLegacyPoolChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyLegacyPoolChanged struct {
	OldPool common.Address
	NewPool common.Address
	Raw     types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterLegacyPoolChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyLegacyPoolChangedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyLegacyPoolChangedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "LegacyPoolChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "LegacyPoolChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyLegacyPoolChanged)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseLegacyPoolChanged(log types.Log) (*BurnMintTokenPoolAndProxyLegacyPoolChanged, error) {
	event := new(BurnMintTokenPoolAndProxyLegacyPoolChanged)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "LegacyPoolChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyLockedIterator struct {
	Event *BurnMintTokenPoolAndProxyLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyLocked)
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
		it.Event = new(BurnMintTokenPoolAndProxyLocked)
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

func (it *BurnMintTokenPoolAndProxyLockedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolAndProxyLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyLockedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyLocked)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseLocked(log types.Log) (*BurnMintTokenPoolAndProxyLocked, error) {
	event := new(BurnMintTokenPoolAndProxyLocked)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyMintedIterator struct {
	Event *BurnMintTokenPoolAndProxyMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyMinted)
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
		it.Event = new(BurnMintTokenPoolAndProxyMinted)
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

func (it *BurnMintTokenPoolAndProxyMintedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolAndProxyMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyMintedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyMinted)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseMinted(log types.Log) (*BurnMintTokenPoolAndProxyMinted, error) {
	event := new(BurnMintTokenPoolAndProxyMinted)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator struct {
	Event *BurnMintTokenPoolAndProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyOwnershipTransferRequested)
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
		it.Event = new(BurnMintTokenPoolAndProxyOwnershipTransferRequested)
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

func (it *BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyOwnershipTransferRequested)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolAndProxyOwnershipTransferRequested, error) {
	event := new(BurnMintTokenPoolAndProxyOwnershipTransferRequested)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyOwnershipTransferredIterator struct {
	Event *BurnMintTokenPoolAndProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyOwnershipTransferred)
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
		it.Event = new(BurnMintTokenPoolAndProxyOwnershipTransferred)
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

func (it *BurnMintTokenPoolAndProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolAndProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyOwnershipTransferredIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyOwnershipTransferred)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolAndProxyOwnershipTransferred, error) {
	event := new(BurnMintTokenPoolAndProxyOwnershipTransferred)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyReleasedIterator struct {
	Event *BurnMintTokenPoolAndProxyReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyReleased)
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
		it.Event = new(BurnMintTokenPoolAndProxyReleased)
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

func (it *BurnMintTokenPoolAndProxyReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolAndProxyReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyReleasedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyReleased)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseReleased(log types.Log) (*BurnMintTokenPoolAndProxyReleased, error) {
	event := new(BurnMintTokenPoolAndProxyReleased)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyRemotePoolSetIterator struct {
	Event *BurnMintTokenPoolAndProxyRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyRemotePoolSet)
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
		it.Event = new(BurnMintTokenPoolAndProxyRemotePoolSet)
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

func (it *BurnMintTokenPoolAndProxyRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolAndProxyRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyRemotePoolSetIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyRemotePoolSet)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseRemotePoolSet(log types.Log) (*BurnMintTokenPoolAndProxyRemotePoolSet, error) {
	event := new(BurnMintTokenPoolAndProxyRemotePoolSet)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyRouterUpdatedIterator struct {
	Event *BurnMintTokenPoolAndProxyRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyRouterUpdated)
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
		it.Event = new(BurnMintTokenPoolAndProxyRouterUpdated)
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

func (it *BurnMintTokenPoolAndProxyRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyRouterUpdatedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyRouterUpdatedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyRouterUpdated)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolAndProxyRouterUpdated, error) {
	event := new(BurnMintTokenPoolAndProxyRouterUpdated)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAndProxyTokensConsumedIterator struct {
	Event *BurnMintTokenPoolAndProxyTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAndProxyTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAndProxyTokensConsumed)
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
		it.Event = new(BurnMintTokenPoolAndProxyTokensConsumed)
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

func (it *BurnMintTokenPoolAndProxyTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAndProxyTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAndProxyTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyTokensConsumedIterator, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAndProxyTokensConsumedIterator{contract: _BurnMintTokenPoolAndProxy.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPoolAndProxy.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAndProxyTokensConsumed)
				if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxyFilterer) ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolAndProxyTokensConsumed, error) {
	event := new(BurnMintTokenPoolAndProxyTokensConsumed)
	if err := _BurnMintTokenPoolAndProxy.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintTokenPoolAndProxy.abi.Events["AllowListAdd"].ID:
		return _BurnMintTokenPoolAndProxy.ParseAllowListAdd(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["AllowListRemove"].ID:
		return _BurnMintTokenPoolAndProxy.ParseAllowListRemove(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["Burned"].ID:
		return _BurnMintTokenPoolAndProxy.ParseBurned(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["ChainAdded"].ID:
		return _BurnMintTokenPoolAndProxy.ParseChainAdded(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["ChainConfigured"].ID:
		return _BurnMintTokenPoolAndProxy.ParseChainConfigured(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["ChainRemoved"].ID:
		return _BurnMintTokenPoolAndProxy.ParseChainRemoved(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["ConfigChanged"].ID:
		return _BurnMintTokenPoolAndProxy.ParseConfigChanged(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["LegacyPoolChanged"].ID:
		return _BurnMintTokenPoolAndProxy.ParseLegacyPoolChanged(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["Locked"].ID:
		return _BurnMintTokenPoolAndProxy.ParseLocked(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["Minted"].ID:
		return _BurnMintTokenPoolAndProxy.ParseMinted(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintTokenPoolAndProxy.ParseOwnershipTransferRequested(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintTokenPoolAndProxy.ParseOwnershipTransferred(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["Released"].ID:
		return _BurnMintTokenPoolAndProxy.ParseReleased(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["RemotePoolSet"].ID:
		return _BurnMintTokenPoolAndProxy.ParseRemotePoolSet(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["RouterUpdated"].ID:
		return _BurnMintTokenPoolAndProxy.ParseRouterUpdated(log)
	case _BurnMintTokenPoolAndProxy.abi.Events["TokensConsumed"].ID:
		return _BurnMintTokenPoolAndProxy.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintTokenPoolAndProxyAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnMintTokenPoolAndProxyAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnMintTokenPoolAndProxyBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnMintTokenPoolAndProxyChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnMintTokenPoolAndProxyChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnMintTokenPoolAndProxyChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnMintTokenPoolAndProxyConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnMintTokenPoolAndProxyLegacyPoolChanged) Topic() common.Hash {
	return common.HexToHash("0x81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f228")
}

func (BurnMintTokenPoolAndProxyLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnMintTokenPoolAndProxyMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnMintTokenPoolAndProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintTokenPoolAndProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintTokenPoolAndProxyReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnMintTokenPoolAndProxyRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (BurnMintTokenPoolAndProxyRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnMintTokenPoolAndProxyTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnMintTokenPoolAndProxy *BurnMintTokenPoolAndProxy) Address() common.Address {
	return _BurnMintTokenPoolAndProxy.address
}

type BurnMintTokenPoolAndProxyInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetOnRamp(opts *bind.CallOpts, arg0 uint64) (common.Address, error)

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

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAndProxyAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAndProxyAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolAndProxyBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnMintTokenPoolAndProxyBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnMintTokenPoolAndProxyChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnMintTokenPoolAndProxyChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnMintTokenPoolAndProxyChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnMintTokenPoolAndProxyConfigChanged, error)

	FilterLegacyPoolChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyLegacyPoolChangedIterator, error)

	WatchLegacyPoolChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyLegacyPoolChanged) (event.Subscription, error)

	ParseLegacyPoolChanged(log types.Log) (*BurnMintTokenPoolAndProxyLegacyPoolChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolAndProxyLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnMintTokenPoolAndProxyLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolAndProxyMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnMintTokenPoolAndProxyMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolAndProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolAndProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolAndProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolAndProxyOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolAndProxyReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnMintTokenPoolAndProxyReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolAndProxyRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*BurnMintTokenPoolAndProxyRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolAndProxyRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolAndProxyTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAndProxyTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolAndProxyTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
