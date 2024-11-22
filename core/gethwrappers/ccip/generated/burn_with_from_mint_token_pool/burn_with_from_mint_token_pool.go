// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_with_from_mint_token_pool

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
	RemotePoolAddress         []byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var BurnWithFromMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620047b9380380620047b983398101604081905262000034916200085d565b83838383336000816200005a57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008d576200008d8162000159565b50506001600160a01b0384161580620000ad57506001600160a01b038116155b80620000c057506001600160a01b038216155b15620000df576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200013257604080516000815260208101909152620001329084620001d3565b506200014f925050506001600160a01b0385163060001962000330565b5050505062000a99565b336001600160a01b038216036200018357604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c051620001f4576040516335f4a7b360e01b815260040160405180910390fd5b60005b82518110156200027f5760008382815181106200021857620002186200096d565b602090810291909101015190506200023260028262000416565b1562000275576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101620001f7565b5060005b81518110156200032b576000828281518110620002a457620002a46200096d565b6020026020010151905060006001600160a01b0316816001600160a01b031603620002d0575062000322565b620002dd60028262000436565b1562000320576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010162000283565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa15801562000382573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003a8919062000983565b620003b49190620009b3565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b1790915291925062000410918691906200044d16565b50505050565b60006200042d836001600160a01b03841662000522565b90505b92915050565b60006200042d836001600160a01b03841662000626565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564908201526000906200049c906001600160a01b03851690849062000678565b8051909150156200032b5780806020019051810190620004bd9190620009c9565b6200032b5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084015b60405180910390fd5b600081815260018301602052604081205480156200061b57600062000549600183620009f4565b85549091506000906200055f90600190620009f4565b9050808214620005cb5760008660000182815481106200058357620005836200096d565b9060005260206000200154905080876000018481548110620005a957620005a96200096d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620005df57620005df62000a0a565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000430565b600091505062000430565b60008181526001830160205260408120546200066f5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000430565b50600062000430565b606062000689848460008562000691565b949350505050565b606082471015620006f45760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000519565b600080866001600160a01b0316858760405162000712919062000a46565b60006040518083038185875af1925050503d806000811462000751576040519150601f19603f3d011682016040523d82523d6000602084013e62000756565b606091505b5090925090506200076a8783838762000775565b979650505050505050565b60608315620007e9578251600003620007e1576001600160a01b0385163b620007e15760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000519565b508162000689565b620006898383815115620008005781518083602001fd5b8060405162461bcd60e51b815260040162000519919062000a64565b6001600160a01b03811681146200083257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b805162000858816200081c565b919050565b600080600080608085870312156200087457600080fd5b845162000881816200081c565b602086810151919550906001600160401b0380821115620008a157600080fd5b818801915088601f830112620008b657600080fd5b815181811115620008cb57620008cb62000835565b8060051b604051601f19603f83011681018181108582111715620008f357620008f362000835565b60405291825284820192508381018501918b8311156200091257600080fd5b938501935b828510156200093b576200092b856200084b565b8452938501939285019262000917565b80985050505050505062000952604086016200084b565b915062000962606086016200084b565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b6000602082840312156200099657600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b808201808211156200043057620004306200099d565b600060208284031215620009dc57600080fd5b81518015158114620009ed57600080fd5b9392505050565b818103818111156200043057620004306200099d565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000a3d57818101518382015260200162000a23565b50506000910152565b6000825162000a5a81846020870162000a20565b9190910192915050565b602081526000825180602084015262000a8581604085016020870162000a20565b601f01601f19169190910160400192915050565b60805160a05160c051613ca362000b166000396000818161050101528181611bd701526125bf0152600081816104db015281816118ce0152611e8a01526000818161022a0152818161027f01528181610654015281816117ee01528181611daa01528181611fa2015281816123b801526127aa0152613ca36000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80639a4575b9116100f9578063c0d7865511610097578063cf7401f311610071578063cf7401f3146104c6578063dc0bd971146104d9578063e0351e13146104ff578063f2fde38b1461052557600080fd5b8063c0d786551461048b578063c4bffe2b1461049e578063c75eea9c146104b357600080fd5b8063acfecf91116100d3578063acfecf91146103d8578063af58d59f146103eb578063b0f479a11461045a578063b79465801461047857600080fd5b80639a4575b914610383578063a42a7b8b146103a3578063a7cd63b7146103c357600080fd5b806354c8a4f31161016657806379ba50971161014057806379ba5097146103375780637d54534e1461033f5780638926f54f146103525780638da5cb5b1461036557600080fd5b806354c8a4f3146102f357806362ddd3c4146103065780636d3d1a581461031957600080fd5b8063240028e811610197578063240028e81461026f57806339077537146102bc5780634d5e7b47146102de57600080fd5b806301ffc9a7146101be578063181f5a77146101e657806321df0da714610228575b600080fd5b6101d16101cc366004612dcf565b610538565b60405190151581526020015b60405180910390f35b60408051808201909152601f81527f4275726e5769746846726f6d4d696e74546f6b656e506f6f6c20312e352e300060208201525b6040516101dd9190612e75565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101dd565b6101d161027d366004612eba565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102cf6102ca366004612ed7565b61061d565b604051905181526020016101dd565b6102f16102ec366004612f5f565b6107a3565b005b6102f1610301366004612f5f565b610cc4565b6102f1610314366004612fe3565b610d3f565b60095473ffffffffffffffffffffffffffffffffffffffff1661024a565b6102f1610ecc565b6102f161034d366004612eba565b610f9a565b6101d1610360366004613066565b61101b565b60015473ffffffffffffffffffffffffffffffffffffffff1661024a565b610396610391366004613081565b611032565b6040516101dd91906130bc565b6103b66103b1366004613066565b6110d9565b6040516101dd919061311c565b6103cb611244565b6040516101dd919061319e565b6102f16103e6366004612fe3565b611255565b6103fe6103f9366004613066565b611360565b6040516101dd919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff1661024a565b61021b610486366004613066565b611435565b6102f1610499366004612eba565b6114e5565b6104a66115c0565b6040516101dd91906131f8565b6103fe6104c1366004613066565b611678565b6102f16104d4366004613355565b61174a565b7f000000000000000000000000000000000000000000000000000000000000000061024a565b7f00000000000000000000000000000000000000000000000000000000000000006101d1565b6102f1610533366004612eba565b6117d3565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf0000000000000000000000000000000000000000000000000000000014806105cb57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061061757507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261063d61063883613445565b6117e7565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f196106896060850160408601612eba565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260608501356024820152604401600060405180830381600087803b1580156106f957600080fd5b505af115801561070d573d6000803e3d6000fd5b50610722925050506060830160408401612eba565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0846060013560405161078491815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6107ab611a19565b60005b8381101561099d5760008585838181106107ca576107ca61353a565b90506020020160208101906107df9190613066565b90506107f6600567ffffffffffffffff8316611a6c565b61083d576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024015b60405180910390fd5b67ffffffffffffffff8116600090815260076020526040812061086290600501611a7f565b905060005b81518110156108ce576108c58282815181106108855761088561353a565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600501611a6c90919063ffffffff16565b50600101610867565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906109376004830182612d62565b60058201600081816109498282612d9c565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169450602001925061098b915050565b60405180910390a150506001016107ae565b5060005b81811015610cbd5760008383838181106109bd576109bd61353a565b90506020028101906109cf9190613569565b6109d8906135a7565b90506109e981606001516000611a8c565b6109f881608001516000611a8c565b8051610a109060059067ffffffffffffffff16611bc9565b610a555780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610834565b6020810151511580610a6a5750604081015151155b15610aa1576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a909102999093161717909416959095179092559092029091176003820155908201516004820190610c2490826136ed565b506020808301518051910120610c3e906005830190611bc9565b5060208083015180518183012060009081526008909252604090912090610c6590826136ed565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28260000151836040015184606001518560800151604051610cab9493929190613807565b60405180910390a150506001016109a1565b5050505050565b610ccc611a19565b610d3984848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611bd592505050565b50505050565b610d47611a19565b610d508361101b565b610d92576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610834565b6000819003610dcd576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e0d8282604051610de09291906138a0565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190611bc9565b610e49578282826040517f393b8ad2000000000000000000000000000000000000000000000000000000008152600401610834939291906138f9565b8181600860008585604051610e5f9291906138a0565b604051809103902081526020019081526020016000209182610e8292919061391d565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8383604051610ebf929190613a37565b60405180910390a2505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610f1d576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610fa2611a19565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b6000610617600567ffffffffffffffff8416611d8b565b604080518082019091526060808252602082015261105761105283613a4b565b611da3565b6110648260600135611f6d565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a260405180604001604052806110be8460200160208101906104869190613066565b81526040805160208181019092526000815291015292915050565b67ffffffffffffffff811660009081526007602052604081206060919061110290600501611a7f565b90506000815167ffffffffffffffff8111156111205761112061323a565b60405190808252806020026020018201604052801561115357816020015b606081526020019060019003908161113e5790505b50905060005b825181101561123c57600860008483815181106111785761117861353a565b6020026020010151815260200190815260200160002080546111999061364a565b80601f01602080910402602001604051908101604052809291908181526020018280546111c59061364a565b80156112125780601f106111e757610100808354040283529160200191611212565b820191906000526020600020905b8154815290600101906020018083116111f557829003601f168201915b50505050508282815181106112295761122961353a565b6020908102919091010152600101611159565b509392505050565b60606112506002611a7f565b905090565b61125d611a19565b6112668361101b565b6112a8576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610834565b6112e882826040516112bb9291906138a0565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190611a6c565b611324578282826040517f74f23c7c000000000000000000000000000000000000000000000000000000008152600401610834939291906138f9565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610ebf929190613a37565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260039091015480841660608301529190910490911660808201526106179061200f565b67ffffffffffffffff811660009081526007602052604090206004018054606091906114609061364a565b80601f016020809104026020016040519081016040528092919081815260200182805461148c9061364a565b80156114d95780601f106114ae576101008083540402835291602001916114d9565b820191906000526020600020905b8154815290600101906020018083116114bc57829003601f168201915b50505050509050919050565b6114ed611a19565b73ffffffffffffffffffffffffffffffffffffffff811661153a576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006115ce6005611a7f565b90506000815167ffffffffffffffff8111156115ec576115ec61323a565b604051908082528060200260200182016040528015611615578160200160208202803683370190505b50905060005b8251811015611671578281815181106116365761163661353a565b60200260200101518282815181106116505761165061353a565b67ffffffffffffffff9092166020928302919091019091015260010161161b565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526106179061200f565b60095473ffffffffffffffffffffffffffffffffffffffff16331480159061178a575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156117c3576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610834565b6117ce8383836120c1565b505050565b6117db611a19565b6117e4816121ab565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161461187c5760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610834565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa15801561192a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061194e9190613ac7565b15611985576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611992816020015161226f565b60a081015180516020918201208183015167ffffffffffffffff166000908152600790925260409091206119cb91600590910190611d8b565b611a07578060a001516040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016108349190612e75565b6117e481602001518260600151612395565b60015473ffffffffffffffffffffffffffffffffffffffff163314611a6a576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000611a7883836123dc565b9392505050565b60606000611a78836124cf565b815115611b575781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580611ae2575060408201516fffffffffffffffffffffffffffffffff16155b15611b1b57816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016108349190613ae4565b8015611b53576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580611b90575060208201516fffffffffffffffffffffffffffffffff1615155b15611b5357816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016108349190613ae4565b6000611a78838361252a565b7f0000000000000000000000000000000000000000000000000000000000000000611c2c576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611cc2576000838281518110611c4c57611c4c61353a565b60200260200101519050611c6a81600261257990919063ffffffff16565b15611cb95760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611c2f565b5060005b81518110156117ce576000828281518110611ce357611ce361353a565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611d275750611d83565b611d3260028261259b565b15611d815760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611cc6565b60008181526001830160205260408120541515611a78565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611e385760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610834565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611ee6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f0a9190613ac7565b15611f41576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611f4e81604001516125bd565b611f5b816020015161263c565b6117e48160200151826060015161278a565b6040517f9dc29fac000000000000000000000000000000000000000000000000000000008152306004820152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639dc29fac90604401600060405180830381600087803b158015611ffb57600080fd5b505af1158015610cbd573d6000803e3d6000fd5b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261209d82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426120819190613b4f565b85608001516fffffffffffffffffffffffffffffffff166127ce565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b6120ca8361101b565b61210c576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610834565b612117826000611a8c565b67ffffffffffffffff8316600090815260076020526040902061213a90836127f8565b612145816000611a8c565b67ffffffffffffffff8316600090815260076020526040902061216b90600201826127f8565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b83838360405161219e93929190613b62565b60405180910390a1505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036121fa576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6122788161101b565b6122ba576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610834565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015612339573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235d9190613ac7565b6117e4576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610834565b67ffffffffffffffff82166000908152600760205260409020611b5390600201827f000000000000000000000000000000000000000000000000000000000000000061299a565b600081815260018301602052604081205480156124c5576000612400600183613b4f565b855490915060009061241490600190613b4f565b90508082146124795760008660000182815481106124345761243461353a565b90600052602060002001549050808760000184815481106124575761245761353a565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061248a5761248a613be5565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610617565b6000915050610617565b6060816000018054806020026020016040519081016040528092919081815260200182805480156114d957602002820191906000526020600020905b81548152602001906001019080831161250b5750505050509050919050565b600081815260018301602052604081205461257157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610617565b506000610617565b6000611a788373ffffffffffffffffffffffffffffffffffffffff84166123dc565b6000611a788373ffffffffffffffffffffffffffffffffffffffff841661252a565b7f0000000000000000000000000000000000000000000000000000000000000000156117e4576125ee600282612d1d565b6117e4576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610834565b6126458161101b565b612687576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610834565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612700573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127249190613c14565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146117e4576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610834565b67ffffffffffffffff82166000908152600760205260409020611b5390827f000000000000000000000000000000000000000000000000000000000000000061299a565b60006127ed856127de8486613c31565b6127e89087613c48565b612d4c565b90505b949350505050565b815460009061282190700100000000000000000000000000000000900463ffffffff1642613b4f565b905080156128c35760018301548354612869916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166127ce565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546128e9916fffffffffffffffffffffffffffffffff9081169116612d4c565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199061219e908490613ae4565b825474010000000000000000000000000000000000000000900460ff1615806129c1575081155b156129cb57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612a1190700100000000000000000000000000000000900463ffffffff1642613b4f565b90508015612ad15781831115612a53576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612a8d9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166127ce565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612b885773ffffffffffffffffffffffffffffffffffffffff8416612b30576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610834565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610834565b84831015612c9b5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612bcc9082613b4f565b612bd6878a613b4f565b612be09190613c48565b612bea9190613c5b565b905073ffffffffffffffffffffffffffffffffffffffff8616612c43576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610834565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610834565b612ca58584613b4f565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611a78565b6000818310612d5b5781611a78565b5090919050565b508054612d6e9061364a565b6000825580601f10612d7e575050565b601f0160209004906000526020600020908101906117e49190612db6565b50805460008255906000526020600020908101906117e491905b5b80821115612dcb5760008155600101612db7565b5090565b600060208284031215612de157600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611a7857600080fd5b6000815180845260005b81811015612e3757602081850181015186830182015201612e1b565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611a786020830184612e11565b73ffffffffffffffffffffffffffffffffffffffff811681146117e457600080fd5b8035612eb581612e88565b919050565b600060208284031215612ecc57600080fd5b8135611a7881612e88565b600060208284031215612ee957600080fd5b813567ffffffffffffffff811115612f0057600080fd5b82016101008185031215611a7857600080fd5b60008083601f840112612f2557600080fd5b50813567ffffffffffffffff811115612f3d57600080fd5b6020830191508360208260051b8501011115612f5857600080fd5b9250929050565b60008060008060408587031215612f7557600080fd5b843567ffffffffffffffff80821115612f8d57600080fd5b612f9988838901612f13565b90965094506020870135915080821115612fb257600080fd5b50612fbf87828801612f13565b95989497509550505050565b803567ffffffffffffffff81168114612eb557600080fd5b600080600060408486031215612ff857600080fd5b61300184612fcb565b9250602084013567ffffffffffffffff8082111561301e57600080fd5b818601915086601f83011261303257600080fd5b81358181111561304157600080fd5b87602082850101111561305357600080fd5b6020830194508093505050509250925092565b60006020828403121561307857600080fd5b611a7882612fcb565b60006020828403121561309357600080fd5b813567ffffffffffffffff8111156130aa57600080fd5b820160a08185031215611a7857600080fd5b6020815260008251604060208401526130d86060840182612e11565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526131138282612e11565b95945050505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613191577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261317f858351612e11565b94509285019290850190600101613145565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156131ec57835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016131ba565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156131ec57835167ffffffffffffffff1683529284019291840191600101613214565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff8111828210171561328d5761328d61323a565b60405290565b60405160a0810167ffffffffffffffff8111828210171561328d5761328d61323a565b80151581146117e457600080fd5b80356fffffffffffffffffffffffffffffffff81168114612eb557600080fd5b6000606082840312156132f657600080fd5b6040516060810181811067ffffffffffffffff821117156133195761331961323a565b604052905080823561332a816132b6565b8152613338602084016132c4565b6020820152613349604084016132c4565b60408201525092915050565b600080600060e0848603121561336a57600080fd5b61337384612fcb565b925061338285602086016132e4565b915061339185608086016132e4565b90509250925092565b600082601f8301126133ab57600080fd5b813567ffffffffffffffff808211156133c6576133c661323a565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561340c5761340c61323a565b8160405283815286602085880101111561342557600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000610100823603121561345857600080fd5b613460613269565b823567ffffffffffffffff8082111561347857600080fd5b6134843683870161339a565b835261349260208601612fcb565b60208401526134a360408601612eaa565b6040840152606085013560608401526134be60808601612eaa565b608084015260a08501359150808211156134d757600080fd5b6134e33683870161339a565b60a084015260c08501359150808211156134fc57600080fd5b6135083683870161339a565b60c084015260e085013591508082111561352157600080fd5b5061352e3682860161339a565b60e08301525092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261359d57600080fd5b9190910192915050565b600061012082360312156135ba57600080fd5b6135c2613293565b6135cb83612fcb565b8152602083013567ffffffffffffffff808211156135e857600080fd5b6135f43683870161339a565b6020840152604085013591508082111561360d57600080fd5b5061361a3682860161339a565b60408301525061362d36606085016132e4565b606082015261363f3660c085016132e4565b608082015292915050565b600181811c9082168061365e57607f821691505b602082108103613697577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156117ce576000816000526020600020601f850160051c810160208610156136c65750805b601f850160051c820191505b818110156136e5578281556001016136d2565b505050505050565b815167ffffffffffffffff8111156137075761370761323a565b61371b81613715845461364a565b8461369d565b602080601f83116001811461376e57600084156137385750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556136e5565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156137bb5788860151825594840194600190910190840161379c565b50858210156137f757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff8716835280602084015261382b81840187612e11565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506138699050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152613113565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff841681526040602082015260006127ed6040830184866138b0565b67ffffffffffffffff8311156139355761393561323a565b61394983613943835461364a565b8361369d565b6000601f84116001811461399b57600085156139655750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355610cbd565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156139ea57868501358255602094850194600190920191016139ca565b5086821015613a25577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b6020815260006127f06020830184866138b0565b600060a08236031215613a5d57600080fd5b613a65613293565b823567ffffffffffffffff811115613a7c57600080fd5b613a883682860161339a565b825250613a9760208401612fcb565b60208201526040830135613aaa81612e88565b604082015260608381013590820152608083013561363f81612e88565b600060208284031215613ad957600080fd5b8151611a78816132b6565b6060810161061782848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561061757610617613b20565b67ffffffffffffffff8416815260e08101613bae60208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526127f0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600060208284031215613c2657600080fd5b8151611a7881612e88565b808202811582820484141761061757610617613b20565b8082018082111561061757610617613b20565b600082613c91577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea164736f6c6343000818000a",
}

var BurnWithFromMintTokenPoolABI = BurnWithFromMintTokenPoolMetaData.ABI

var BurnWithFromMintTokenPoolBin = BurnWithFromMintTokenPoolMetaData.Bin

func DeployBurnWithFromMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnWithFromMintTokenPool, error) {
	parsed, err := BurnWithFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnWithFromMintTokenPoolBin), backend, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnWithFromMintTokenPool{address: address, abi: *parsed, BurnWithFromMintTokenPoolCaller: BurnWithFromMintTokenPoolCaller{contract: contract}, BurnWithFromMintTokenPoolTransactor: BurnWithFromMintTokenPoolTransactor{contract: contract}, BurnWithFromMintTokenPoolFilterer: BurnWithFromMintTokenPoolFilterer{contract: contract}}, nil
}

type BurnWithFromMintTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnWithFromMintTokenPoolCaller
	BurnWithFromMintTokenPoolTransactor
	BurnWithFromMintTokenPoolFilterer
}

type BurnWithFromMintTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolSession struct {
	Contract     *BurnWithFromMintTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolCallerSession struct {
	Contract *BurnWithFromMintTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnWithFromMintTokenPoolTransactorSession struct {
	Contract     *BurnWithFromMintTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolRaw struct {
	Contract *BurnWithFromMintTokenPool
}

type BurnWithFromMintTokenPoolCallerRaw struct {
	Contract *BurnWithFromMintTokenPoolCaller
}

type BurnWithFromMintTokenPoolTransactorRaw struct {
	Contract *BurnWithFromMintTokenPoolTransactor
}

func NewBurnWithFromMintTokenPool(address common.Address, backend bind.ContractBackend) (*BurnWithFromMintTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnWithFromMintTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnWithFromMintTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPool{address: address, abi: abi, BurnWithFromMintTokenPoolCaller: BurnWithFromMintTokenPoolCaller{contract: contract}, BurnWithFromMintTokenPoolTransactor: BurnWithFromMintTokenPoolTransactor{contract: contract}, BurnWithFromMintTokenPoolFilterer: BurnWithFromMintTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnWithFromMintTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnWithFromMintTokenPoolCaller, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolCaller{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnWithFromMintTokenPoolTransactor, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolTransactor{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnWithFromMintTokenPoolFilterer, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolFilterer{contract: contract}, nil
}

func bindBurnWithFromMintTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnWithFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowList(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowList(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemotePools(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemotePools(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRouter(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRouter(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetToken(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetToken(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintTokenPool.CallOpts, token)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintTokenPool.CallOpts, token)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.Owner(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.Owner(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.SupportsInterface(&_BurnWithFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.SupportsInterface(&_BurnWithFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintTokenPool.TransactOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintTokenPool.TransactOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AddRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AddRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.LockOrBurn(&_BurnWithFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.LockOrBurn(&_BurnWithFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.RemoveRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.RemoveRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRouter(&_BurnWithFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRouter(&_BurnWithFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.TransferOwnership(&_BurnWithFromMintTokenPool.TransactOpts, to)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.TransferOwnership(&_BurnWithFromMintTokenPool.TransactOpts, to)
}

type BurnWithFromMintTokenPoolAllowListAddIterator struct {
	Event *BurnWithFromMintTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAllowListAdd)
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
		it.Event = new(BurnWithFromMintTokenPoolAllowListAdd)
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

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAllowListAddIterator{contract: _BurnWithFromMintTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAllowListAdd)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAllowListAdd, error) {
	event := new(BurnWithFromMintTokenPoolAllowListAdd)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAllowListRemoveIterator struct {
	Event *BurnWithFromMintTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAllowListRemove)
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
		it.Event = new(BurnWithFromMintTokenPoolAllowListRemove)
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

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAllowListRemoveIterator{contract: _BurnWithFromMintTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAllowListRemove)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAllowListRemove, error) {
	event := new(BurnWithFromMintTokenPoolAllowListRemove)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolBurnedIterator struct {
	Event *BurnWithFromMintTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolBurned)
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
		it.Event = new(BurnWithFromMintTokenPoolBurned)
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

func (it *BurnWithFromMintTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolBurnedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolBurned)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolBurned, error) {
	event := new(BurnWithFromMintTokenPoolBurned)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainAddedIterator struct {
	Event *BurnWithFromMintTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainAdded)
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
		it.Event = new(BurnWithFromMintTokenPoolChainAdded)
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

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainAddedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainAdded)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolChainAdded, error) {
	event := new(BurnWithFromMintTokenPoolChainAdded)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainConfiguredIterator struct {
	Event *BurnWithFromMintTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainConfigured)
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
		it.Event = new(BurnWithFromMintTokenPoolChainConfigured)
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

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainConfiguredIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainConfigured)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolChainConfigured, error) {
	event := new(BurnWithFromMintTokenPoolChainConfigured)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainRemovedIterator struct {
	Event *BurnWithFromMintTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainRemoved)
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
		it.Event = new(BurnWithFromMintTokenPoolChainRemoved)
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

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainRemovedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainRemoved)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolChainRemoved, error) {
	event := new(BurnWithFromMintTokenPoolChainRemoved)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolConfigChangedIterator struct {
	Event *BurnWithFromMintTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolConfigChanged)
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
		it.Event = new(BurnWithFromMintTokenPoolConfigChanged)
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

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolConfigChangedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolConfigChanged)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolConfigChanged, error) {
	event := new(BurnWithFromMintTokenPoolConfigChanged)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolLockedIterator struct {
	Event *BurnWithFromMintTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolLocked)
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
		it.Event = new(BurnWithFromMintTokenPoolLocked)
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

func (it *BurnWithFromMintTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolLockedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolLocked)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolLocked, error) {
	event := new(BurnWithFromMintTokenPoolLocked)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolMintedIterator struct {
	Event *BurnWithFromMintTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolMinted)
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
		it.Event = new(BurnWithFromMintTokenPoolMinted)
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

func (it *BurnWithFromMintTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolMintedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolMinted)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolMinted, error) {
	event := new(BurnWithFromMintTokenPoolMinted)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnWithFromMintTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
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

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolOwnershipTransferredIterator struct {
	Event *BurnWithFromMintTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferred)
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

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolOwnershipTransferredIterator{contract: _BurnWithFromMintTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolOwnershipTransferred)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferred, error) {
	event := new(BurnWithFromMintTokenPoolOwnershipTransferred)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnWithFromMintTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnWithFromMintTokenPoolRateLimitAdminSet)
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

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRateLimitAdminSetIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRateLimitAdminSet)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnWithFromMintTokenPoolRateLimitAdminSet, error) {
	event := new(BurnWithFromMintTokenPoolRateLimitAdminSet)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolReleasedIterator struct {
	Event *BurnWithFromMintTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolReleased)
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
		it.Event = new(BurnWithFromMintTokenPoolReleased)
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

func (it *BurnWithFromMintTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolReleasedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolReleased)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolReleased, error) {
	event := new(BurnWithFromMintTokenPoolReleased)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRemotePoolAddedIterator struct {
	Event *BurnWithFromMintTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnWithFromMintTokenPoolRemotePoolAdded)
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

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRemotePoolAddedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRemotePoolAdded)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolAdded, error) {
	event := new(BurnWithFromMintTokenPoolRemotePoolAdded)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnWithFromMintTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnWithFromMintTokenPoolRemotePoolRemoved)
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

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRemotePoolRemovedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRemotePoolRemoved)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolRemoved, error) {
	event := new(BurnWithFromMintTokenPoolRemotePoolRemoved)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRouterUpdatedIterator struct {
	Event *BurnWithFromMintTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRouterUpdated)
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
		it.Event = new(BurnWithFromMintTokenPoolRouterUpdated)
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

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRouterUpdatedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRouterUpdated)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolRouterUpdated, error) {
	event := new(BurnWithFromMintTokenPoolRouterUpdated)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolTokensConsumedIterator struct {
	Event *BurnWithFromMintTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolTokensConsumed)
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
		it.Event = new(BurnWithFromMintTokenPoolTokensConsumed)
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

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolTokensConsumedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolTokensConsumed)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolTokensConsumed, error) {
	event := new(BurnWithFromMintTokenPoolTokensConsumed)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnWithFromMintTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnWithFromMintTokenPool.ParseAllowListAdd(log)
	case _BurnWithFromMintTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnWithFromMintTokenPool.ParseAllowListRemove(log)
	case _BurnWithFromMintTokenPool.abi.Events["Burned"].ID:
		return _BurnWithFromMintTokenPool.ParseBurned(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnWithFromMintTokenPool.ParseChainAdded(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnWithFromMintTokenPool.ParseChainConfigured(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnWithFromMintTokenPool.ParseChainRemoved(log)
	case _BurnWithFromMintTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnWithFromMintTokenPool.ParseConfigChanged(log)
	case _BurnWithFromMintTokenPool.abi.Events["Locked"].ID:
		return _BurnWithFromMintTokenPool.ParseLocked(log)
	case _BurnWithFromMintTokenPool.abi.Events["Minted"].ID:
		return _BurnWithFromMintTokenPool.ParseMinted(log)
	case _BurnWithFromMintTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnWithFromMintTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnWithFromMintTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnWithFromMintTokenPool.ParseOwnershipTransferred(log)
	case _BurnWithFromMintTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnWithFromMintTokenPool.ParseRateLimitAdminSet(log)
	case _BurnWithFromMintTokenPool.abi.Events["Released"].ID:
		return _BurnWithFromMintTokenPool.ParseReleased(log)
	case _BurnWithFromMintTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnWithFromMintTokenPool.ParseRemotePoolAdded(log)
	case _BurnWithFromMintTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnWithFromMintTokenPool.ParseRemotePoolRemoved(log)
	case _BurnWithFromMintTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnWithFromMintTokenPool.ParseRouterUpdated(log)
	case _BurnWithFromMintTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnWithFromMintTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnWithFromMintTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnWithFromMintTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnWithFromMintTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnWithFromMintTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnWithFromMintTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnWithFromMintTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnWithFromMintTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnWithFromMintTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnWithFromMintTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnWithFromMintTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnWithFromMintTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnWithFromMintTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnWithFromMintTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnWithFromMintTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnWithFromMintTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnWithFromMintTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnWithFromMintTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPool) Address() common.Address {
	return _BurnWithFromMintTokenPool.address
}

type BurnWithFromMintTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

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

	AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnWithFromMintTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
