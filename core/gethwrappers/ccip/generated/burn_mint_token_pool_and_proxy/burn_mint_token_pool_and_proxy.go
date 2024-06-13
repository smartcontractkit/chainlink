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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRatelimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"oldPool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"LegacyPoolChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getOnRamp\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"onRampAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolPriorTo1_5\",\"name\":\"prevPool\",\"type\":\"address\"}],\"name\":\"setPreviousPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200482338038062004823833981016040819052620000349162000554565b83838383838383833380600081620000935760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c657620000c6816200017e565b5050506001600160a01b0384161580620000e757506001600160a01b038116155b80620000fa57506001600160a01b038216155b1562000119576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200016c576040805160008152602081019091526200016c908462000229565b505050505050505050505050620006b2565b336001600160a01b03821603620001d85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c0516200024a576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002d55760008382815181106200026e576200026e62000664565b602090810291909101015190506200028860028262000386565b15620002cb576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016200024d565b5060005b815181101562000381576000828281518110620002fa57620002fa62000664565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000326575062000378565b62000333600282620003a6565b1562000376576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002d9565b505050565b60006200039d836001600160a01b038416620003bd565b90505b92915050565b60006200039d836001600160a01b038416620004c1565b60008181526001830160205260408120548015620004b6576000620003e46001836200067a565b8554909150600090620003fa906001906200067a565b9050818114620004665760008660000182815481106200041e576200041e62000664565b906000526020600020015490508087600001848154811062000444576200044462000664565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200047a576200047a6200069c565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a0565b6000915050620003a0565b60008181526001830160205260408120546200050a57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a0565b506000620003a0565b6001600160a01b03811681146200052957600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b80516200054f8162000513565b919050565b600080600080608085870312156200056b57600080fd5b8451620005788162000513565b602086810151919550906001600160401b03808211156200059857600080fd5b818801915088601f830112620005ad57600080fd5b815181811115620005c257620005c26200052c565b8060051b604051601f19603f83011681018181108582111715620005ea57620005ea6200052c565b60405291825284820192508381018501918b8311156200060957600080fd5b938501935b828510156200063257620006228562000542565b845293850193928501926200060e565b809850505050505050620006496040860162000542565b9150620006596060860162000542565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003a057634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c0516140ed62000736600039600081816104bc0152818161197001526123c2015260008181610496015281816117080152611c23015260008181610210015281816102650152818161071901528181610d040152818161162801528181611b4301528181611d290152818161235801526125b201526140ed6000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80639a4575b9116100f9578063c4bffe2b11610097578063db6327dc11610071578063db6327dc14610481578063dc0bd97114610494578063e0351e13146104ba578063f2fde38b146104e057600080fd5b8063c4bffe2b14610446578063c75eea9c1461045b578063cf7401f31461046e57600080fd5b8063af58d59f116100d3578063af58d59f14610393578063b0f479a114610402578063b794658014610420578063c0d786551461043357600080fd5b80639a4575b91461034b578063a7cd63b71461036b578063a8d87a3b1461038057600080fd5b806354c8a4f31161016657806383826b2b1161014057806383826b2b146102f45780638926f54f146103075780638da5cb5b1461031a5780639766b9321461033857600080fd5b806354c8a4f3146102c457806378a010b2146102d957806379ba5097146102ec57600080fd5b806321df0da71161019757806321df0da71461020e578063240028e81461025557806339077537146102a257600080fd5b806301ffc9a7146101be5780630a2fd493146101e6578063181f5a7714610206575b600080fd5b6101d16101cc36600461305f565b6104f3565b60405190151581526020015b60405180910390f35b6101f96101f43660046130be565b6105d8565b6040516101dd9190613147565b6101f9610688565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101dd565b6101d1610263366004613187565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102b56102b03660046131a4565b6106a4565b604051905181526020016101dd565b6102d76102d236600461322c565b610831565b005b6102d76102e7366004613298565b6108ac565b6102d7610a20565b6101d161030236600461331b565b610b1d565b6101d16103153660046130be565b610bea565b60005473ffffffffffffffffffffffffffffffffffffffff16610230565b6102d7610346366004613187565b610c01565b61035e610359366004613352565b610c90565b6040516101dd919061338d565b610373610e00565b6040516101dd91906133ed565b61023061038e3660046130be565b503090565b6103a66103a13660046130be565b610e11565b6040516101dd919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610230565b6101f961042e3660046130be565b610ee6565b6102d7610441366004613187565b610f11565b61044e610fe5565b6040516101dd9190613447565b6103a66104693660046130be565b61109d565b6102d761047c3660046135fe565b61116f565b6102d761048f366004613643565b611187565b7f0000000000000000000000000000000000000000000000000000000000000000610230565b7f00000000000000000000000000000000000000000000000000000000000000006101d1565b6102d76104ee366004613187565b61160d565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061058657507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b806105d257507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061060390613685565b80601f016020809104026020016040519081016040528092919081815260200182805461062f90613685565b801561067c5780601f106106515761010080835404028352916020019161067c565b820191906000526020600020905b81548152906001019060200180831161065f57829003601f168201915b50505050509050919050565b6040518060600160405280602381526020016140be6023913981565b6040805160208101909152600081526106c46106bf83613774565b611621565b60085473ffffffffffffffffffffffffffffffffffffffff1661078f576040517f40c10f19000000000000000000000000000000000000000000000000000000008152336004820152606083013560248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906340c10f1990604401600060405180830381600087803b15801561077257600080fd5b505af1158015610786573d6000803e3d6000fd5b505050506107a0565b6107a061079b83613774565b611852565b6107b06060830160408401613187565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0846060013560405161081291815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6108396118eb565b6108a68484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061196e92505050565b50505050565b6108b46118eb565b6108bd83610bea565b610904576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b67ffffffffffffffff83166000908152600760205260408120600401805461092b90613685565b80601f016020809104026020016040519081016040528092919081815260200182805461095790613685565b80156109a45780601f10610979576101008083540402835291602001916109a4565b820191906000526020600020905b81548152906001019060200180831161098757829003601f168201915b5050505067ffffffffffffffff86166000908152600760205260409020919250506004016109d38385836138b9565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610a12939291906139d3565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610aa1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016108fb565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600073ffffffffffffffffffffffffffffffffffffffff8216301480610be35750600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86169281019290925273ffffffffffffffffffffffffffffffffffffffff848116602484015216906383826b2b90604401602060405180830381865afa158015610bbf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be39190613a37565b9392505050565b60006105d2600567ffffffffffffffff8416611b24565b610c096118eb565b6008805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f81accd0a7023865eaa51b3399dd0eafc488bf3ba238402911e1659cfe860f22891015b60405180910390a15050565b6040805180820190915260608082526020820152610cb5610cb083613a54565b611b3c565b60085473ffffffffffffffffffffffffffffffffffffffff16610d7a576040517f42966c68000000000000000000000000000000000000000000000000000000008152606083013560048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906342966c6890602401600060405180830381600087803b158015610d5d57600080fd5b505af1158015610d71573d6000803e3d6000fd5b50505050610d8b565b610d8b610d8683613a54565b611d06565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610de584602001602081019061042e91906130be565b81526040805160208181019092526000815291015292915050565b6060610e0c6002611e20565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260039091015480841660608301529190910490911660808201526105d290611e2d565b67ffffffffffffffff8116600090815260076020526040902060050180546060919061060390613685565b610f196118eb565b73ffffffffffffffffffffffffffffffffffffffff8116610f66576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849101610c84565b60606000610ff36005611e20565b90506000815167ffffffffffffffff81111561101157611011613489565b60405190808252806020026020018201604052801561103a578160200160208202803683370190505b50905060005b82518110156110965782818151811061105b5761105b613af6565b602002602001015182828151811061107557611075613af6565b67ffffffffffffffff90921660209283029190910190910152600101611040565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526105d290611e2d565b6111776118eb565b611182838383611edf565b505050565b61118f6118eb565b60005b818110156111825760008383838181106111ae576111ae613af6565b90506020028101906111c09190613b25565b6111c990613b63565b90506111de8160800151826020015115611fc9565b6111f18160a00151826020015115611fc9565b8060200151156114ed5780516112139060059067ffffffffffffffff16612102565b6112585780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016108fb565b604081015151158061126d5750606081015151155b156112a4576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c179091169615150295909517909855908101519401519381169316909102919091176003820155915190919060048201906114859082613c17565b506060820151600582019061149a9082613c17565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506114e09493929190613d31565b60405180910390a1611604565b80516115059060059067ffffffffffffffff1661210e565b61154a5780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016108fb565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906115b36004830182613011565b6115c1600583016000613011565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b50600101611192565b6116156118eb565b61161e8161211a565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff9081169116146116b65760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016108fb565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611764573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117889190613a37565b156117bf576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117cc816020015161220f565b60006117db82602001516105d8565b90508051600014806117ff575080805190602001208260a001518051906020012014155b1561183c578160a001516040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016108fb9190613147565b61184e82602001518360600151612335565b5050565b6008548151606083015160208401516040517f8627fad600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90941693638627fad6936118b69390923392600401613dca565b600060405180830381600087803b1580156118d057600080fd5b505af11580156118e4573d6000803e3d6000fd5b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461196c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016108fb565b565b7f00000000000000000000000000000000000000000000000000000000000000006119c5576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611a5b5760008382815181106119e5576119e5613af6565b60200260200101519050611a0381600261237c90919063ffffffff16565b15611a525760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016119c8565b5060005b8151811015611182576000828281518110611a7c57611a7c613af6565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611ac05750611b1c565b611acb60028261239e565b15611b1a5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611a5f565b60008181526001830160205260408120541515610be3565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611bd15760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016108fb565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611c7f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ca39190613a37565b15611cda576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611ce781604001516123c0565b611cf48160200151612444565b61161e81602001518260600151612592565b6008546060820151611d539173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116929116906125d6565b60085460408083015183516060850151602086015193517f9687544500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90951694639687544594611dbb94939291600401613e2b565b6000604051808303816000875af1158015611dda573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261184e9190810190613e8b565b60606000610be383612663565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611ebb82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611e9f9190613f28565b85608001516fffffffffffffffffffffffffffffffff166126be565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b611ee883610bea565b611f2a576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024016108fb565b611f35826000611fc9565b67ffffffffffffffff83166000908152600760205260409020611f5890836126e8565b611f63816000611fc9565b67ffffffffffffffff83166000908152600760205260409020611f8990600201826126e8565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b838383604051611fbc93929190613f3b565b60405180910390a1505050565b8151156120905781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff1610158061201f575060408201516fffffffffffffffffffffffffffffffff16155b1561205857816040517f70505e560000000000000000000000000000000000000000000000000000000081526004016108fb9190613fbe565b801561184e576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff161515806120c9575060208201516fffffffffffffffffffffffffffffffff1615155b1561184e57816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016108fb9190613fbe565b6000610be3838361288a565b6000610be383836128d9565b3373ffffffffffffffffffffffffffffffffffffffff821603612199576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016108fb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61221881610bea565b61225a576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024016108fb565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156122d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122fd9190613a37565b61161e576040517f728fe07b0000000000000000000000000000000000000000000000000000000081523360048201526024016108fb565b67ffffffffffffffff8216600090815260076020526040902061184e90600201827f00000000000000000000000000000000000000000000000000000000000000006129cc565b6000610be38373ffffffffffffffffffffffffffffffffffffffff84166128d9565b6000610be38373ffffffffffffffffffffffffffffffffffffffff841661288a565b7f000000000000000000000000000000000000000000000000000000000000000080156123f557506123f3600282612d4f565b155b1561161e576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016108fb565b61244d81610bea565b61248f576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024016108fb565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612508573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061252c9190613ffa565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461161e576040517f728fe07b0000000000000000000000000000000000000000000000000000000081523360048201526024016108fb565b67ffffffffffffffff8216600090815260076020526040902061184e90827f00000000000000000000000000000000000000000000000000000000000000006129cc565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052611182908490612d7e565b60608160000180548060200260200160405190810160405280929190818152602001828054801561067c57602002820191906000526020600020905b81548152602001906001019080831161269f5750505050509050919050565b60006126dd856126ce8486614017565b6126d8908761402e565b612e8a565b90505b949350505050565b815460009061271190700100000000000000000000000000000000900463ffffffff1642613f28565b905080156127b35760018301548354612759916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166126be565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546127d9916fffffffffffffffffffffffffffffffff9081169116612e8a565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990611fbc908490613fbe565b60008181526001830160205260408120546128d1575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105d2565b5060006105d2565b600081815260018301602052604081205480156129c25760006128fd600183613f28565b855490915060009061291190600190613f28565b905081811461297657600086600001828154811061293157612931613af6565b906000526020600020015490508087600001848154811061295457612954613af6565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061298757612987614041565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105d2565b60009150506105d2565b825474010000000000000000000000000000000000000000900460ff1615806129f3575081155b156129fd57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612a4390700100000000000000000000000000000000900463ffffffff1642613f28565b90508015612b035781831115612a85576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612abf9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166126be565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612bba5773ffffffffffffffffffffffffffffffffffffffff8416612b62576040517ff94ebcd100000000000000000000000000000000000000000000000000000000815260048101839052602481018690526044016108fb565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff851660448201526064016108fb565b84831015612ccd5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612bfe9082613f28565b612c08878a613f28565b612c12919061402e565b612c1c9190614070565b905073ffffffffffffffffffffffffffffffffffffffff8616612c75576040517f15279c0800000000000000000000000000000000000000000000000000000000815260048101829052602481018690526044016108fb565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff871660448201526064016108fb565b612cd78584613f28565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610be3565b6000612de0826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16612ea09092919063ffffffff16565b8051909150156111825780806020019051810190612dfe9190613a37565b611182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016108fb565b6000818310612e995781610be3565b5090919050565b60606126e08484600085856000808673ffffffffffffffffffffffffffffffffffffffff168587604051612ed491906140ab565b60006040518083038185875af1925050503d8060008114612f11576040519150601f19603f3d011682016040523d82523d6000602084013e612f16565b606091505b5091509150612f2787838387612f32565b979650505050505050565b60608315612fc8578251600003612fc15773ffffffffffffffffffffffffffffffffffffffff85163b612fc1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016108fb565b50816126e0565b6126e08383815115612fdd5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108fb9190613147565b50805461301d90613685565b6000825580601f1061302d575050565b601f01602090049060005260206000209081019061161e91905b8082111561305b5760008155600101613047565b5090565b60006020828403121561307157600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610be357600080fd5b803567ffffffffffffffff811681146130b957600080fd5b919050565b6000602082840312156130d057600080fd5b610be3826130a1565b60005b838110156130f45781810151838201526020016130dc565b50506000910152565b600081518084526131158160208601602086016130d9565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610be360208301846130fd565b73ffffffffffffffffffffffffffffffffffffffff8116811461161e57600080fd5b80356130b98161315a565b60006020828403121561319957600080fd5b8135610be38161315a565b6000602082840312156131b657600080fd5b813567ffffffffffffffff8111156131cd57600080fd5b82016101008185031215610be357600080fd5b60008083601f8401126131f257600080fd5b50813567ffffffffffffffff81111561320a57600080fd5b6020830191508360208260051b850101111561322557600080fd5b9250929050565b6000806000806040858703121561324257600080fd5b843567ffffffffffffffff8082111561325a57600080fd5b613266888389016131e0565b9096509450602087013591508082111561327f57600080fd5b5061328c878288016131e0565b95989497509550505050565b6000806000604084860312156132ad57600080fd5b6132b6846130a1565b9250602084013567ffffffffffffffff808211156132d357600080fd5b818601915086601f8301126132e757600080fd5b8135818111156132f657600080fd5b87602082850101111561330857600080fd5b6020830194508093505050509250925092565b6000806040838503121561332e57600080fd5b613337836130a1565b915060208301356133478161315a565b809150509250929050565b60006020828403121561336457600080fd5b813567ffffffffffffffff81111561337b57600080fd5b820160a08185031215610be357600080fd5b6020815260008251604060208401526133a960608401826130fd565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526133e482826130fd565b95945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561343b57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613409565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561343b57835167ffffffffffffffff1683529284019291840191600101613463565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff811182821017156134dc576134dc613489565b60405290565b60405160c0810167ffffffffffffffff811182821017156134dc576134dc613489565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561354c5761354c613489565b604052919050565b801515811461161e57600080fd5b80356130b981613554565b80356fffffffffffffffffffffffffffffffff811681146130b957600080fd5b60006060828403121561359f57600080fd5b6040516060810181811067ffffffffffffffff821117156135c2576135c2613489565b60405290508082356135d381613554565b81526135e16020840161356d565b60208201526135f26040840161356d565b60408201525092915050565b600080600060e0848603121561361357600080fd5b61361c846130a1565b925061362b856020860161358d565b915061363a856080860161358d565b90509250925092565b6000806020838503121561365657600080fd5b823567ffffffffffffffff81111561366d57600080fd5b613679858286016131e0565b90969095509350505050565b600181811c9082168061369957607f821691505b6020821081036136d2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600067ffffffffffffffff8211156136f2576136f2613489565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261372f57600080fd5b813561374261373d826136d8565b613505565b81815284602083860101111561375757600080fd5b816020850160208301376000918101602001919091529392505050565b6000610100823603121561378757600080fd5b61378f6134b8565b823567ffffffffffffffff808211156137a757600080fd5b6137b33683870161371e565b83526137c1602086016130a1565b60208401526137d26040860161317c565b6040840152606085013560608401526137ed6080860161317c565b608084015260a085013591508082111561380657600080fd5b6138123683870161371e565b60a084015260c085013591508082111561382b57600080fd5b6138373683870161371e565b60c084015260e085013591508082111561385057600080fd5b5061385d3682860161371e565b60e08301525092915050565b601f821115611182576000816000526020600020601f850160051c810160208610156138925750805b601f850160051c820191505b818110156138b15782815560010161389e565b505050505050565b67ffffffffffffffff8311156138d1576138d1613489565b6138e5836138df8354613685565b83613869565b6000601f84116001811461393757600085156139015750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556118e4565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156139865786850135825560209485019460019092019101613966565b50868210156139c1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b6040815260006139e660408301866130fd565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060208284031215613a4957600080fd5b8151610be381613554565b600060a08236031215613a6657600080fd5b60405160a0810167ffffffffffffffff8282108183111715613a8a57613a8a613489565b816040528435915080821115613a9f57600080fd5b50613aac3682860161371e565b825250613abb602084016130a1565b60208201526040830135613ace8161315a565b6040820152606083810135908201526080830135613aeb8161315a565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec1833603018112613b5957600080fd5b9190910192915050565b60006101408236031215613b7657600080fd5b613b7e6134e2565b613b87836130a1565b8152613b9560208401613562565b6020820152604083013567ffffffffffffffff80821115613bb557600080fd5b613bc13683870161371e565b60408401526060850135915080821115613bda57600080fd5b50613be73682860161371e565b606083015250613bfa366080850161358d565b6080820152613c0c3660e0850161358d565b60a082015292915050565b815167ffffffffffffffff811115613c3157613c31613489565b613c4581613c3f8454613685565b84613869565b602080601f831160018114613c985760008415613c625750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556138b1565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613ce557888601518255948401946001909101908401613cc6565b5085821015613d2157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613d55818401876130fd565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613d939050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e08301526133e4565b60a081526000613ddd60a08301876130fd565b73ffffffffffffffffffffffffffffffffffffffff8616602084015284604084015267ffffffffffffffff841660608401528281036080840152600081526020810191505095945050505050565b73ffffffffffffffffffffffffffffffffffffffff8516815260a060208201526000613e5a60a08301866130fd565b60408301949094525067ffffffffffffffff9190911660608201528082036080909101526000815260200192915050565b600060208284031215613e9d57600080fd5b815167ffffffffffffffff811115613eb457600080fd5b8201601f81018413613ec557600080fd5b8051613ed361373d826136d8565b818152856020838501011115613ee857600080fd5b6133e48260208301602086016130d9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156105d2576105d2613ef9565b67ffffffffffffffff8416815260e08101613f8760208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526126e0565b606081016105d282848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561400c57600080fd5b8151610be38161315a565b80820281158282048414176105d2576105d2613ef9565b808201808211156105d2576105d2613ef9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6000826140a6577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008251613b598184602087016130d956fe4275726e4d696e74546f6b656e506f6f6c416e6450726f787920312e352e302d646576a164736f6c6343000818000a",
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
