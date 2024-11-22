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
	RemotePoolAddresses       [][]byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var BurnWithFromMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"localTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004a6338038062004a63833981016040819052620000359162000869565b8484848484336000816200005c57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008f576200008f8162000165565b50506001600160a01b0385161580620000af57506001600160a01b038116155b80620000c257506001600160a01b038216155b15620000e1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0385811660805282811660c05260ff851660a052600480546001600160a01b031916918316919091179055825115801560e0526200013b576040805160008152602081019091526200013b9084620001df565b506200015a935050506001600160a01b0387169050306000196200033c565b505050505062000ac0565b336001600160a01b038216036200018f57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e05162000200576040516335f4a7b360e01b815260040160405180910390fd5b60005b82518110156200028b57600083828151811062000224576200022462000994565b602090810291909101015190506200023e60028262000422565b1562000281576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000203565b5060005b815181101562000337576000828281518110620002b057620002b062000994565b6020026020010151905060006001600160a01b0316816001600160a01b031603620002dc57506200032e565b620002e960028262000442565b156200032c576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6001016200028f565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa1580156200038e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003b49190620009aa565b620003c09190620009da565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b179091529192506200041c918691906200045916565b50505050565b600062000439836001600160a01b0384166200052e565b90505b92915050565b600062000439836001600160a01b03841662000632565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656490820152600090620004a8906001600160a01b03851690849062000684565b805190915015620003375780806020019051810190620004c99190620009f0565b620003375760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084015b60405180910390fd5b60008181526001830160205260408120548015620006275760006200055560018362000a1b565b85549091506000906200056b9060019062000a1b565b9050808214620005d75760008660000182815481106200058f576200058f62000994565b9060005260206000200154905080876000018481548110620005b557620005b562000994565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620005eb57620005eb62000a31565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200043c565b60009150506200043c565b60008181526001830160205260408120546200067b575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200043c565b5060006200043c565b60606200069584846000856200069d565b949350505050565b606082471015620007005760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000525565b600080866001600160a01b031685876040516200071e919062000a6d565b60006040518083038185875af1925050503d80600081146200075d576040519150601f19603f3d011682016040523d82523d6000602084013e62000762565b606091505b509092509050620007768783838762000781565b979650505050505050565b60608315620007f5578251600003620007ed576001600160a01b0385163b620007ed5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000525565b508162000695565b6200069583838151156200080c5781518083602001fd5b8060405162461bcd60e51b815260040162000525919062000a8b565b6001600160a01b03811681146200083e57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b8051620008648162000828565b919050565b600080600080600060a086880312156200088257600080fd5b85516200088f8162000828565b8095505060208087015160ff81168114620008a957600080fd5b60408801519095506001600160401b0380821115620008c757600080fd5b818901915089601f830112620008dc57600080fd5b815181811115620008f157620008f162000841565b8060051b604051601f19603f8301168101818110858211171562000919576200091962000841565b60405291825284820192508381018501918c8311156200093857600080fd5b938501935b828510156200096157620009518562000857565b845293850193928501926200093d565b809850505050505050620009786060870162000857565b9150620009886080870162000857565b90509295509295909350565b634e487b7160e01b600052603260045260246000fd5b600060208284031215620009bd57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b808201808211156200043c576200043c620009c4565b60006020828403121562000a0357600080fd5b8151801515811462000a1457600080fd5b9392505050565b818103818111156200043c576200043c620009c4565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000a6457818101518382015260200162000a4a565b50506000910152565b6000825162000a8181846020870162000a47565b9190910192915050565b602081526000825180602084015262000aac81604085016020870162000a47565b601f01601f19169190910160400192915050565b60805160a05160c05160e051613f0062000b636000396000818161054801528181611c5401526126a5015260008181610522015281816118980152611f400152600081816102d901528181610ba201528181611a4101528181611afb01528181611b2f01528181611b600152611ba70152600081816102400152818161029501528181610701015281816120c30152818161263b01526128900152613f006000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80639a4575b911610104578063c0d78655116100a2578063dc0bd97111610071578063dc0bd97114610520578063e0351e1314610546578063e8a1da171461056c578063f2fde38b1461057f57600080fd5b8063c0d78655146104d2578063c4bffe2b146104e5578063c75eea9c146104fa578063cf7401f31461050d57600080fd5b8063acfecf91116100de578063acfecf911461041f578063af58d59f14610432578063b0f479a1146104a1578063b7946580146104bf57600080fd5b80639a4575b9146103ca578063a42a7b8b146103ea578063a7cd63b71461040a57600080fd5b806354c8a4f31161017157806379ba50971161014b57806379ba50971461037e5780637d54534e146103865780638926f54f146103995780638da5cb5b146103ac57600080fd5b806354c8a4f31461033857806362ddd3c41461034d5780636d3d1a581461036057600080fd5b8063240028e8116101ad578063240028e81461028557806324f65ee7146102d257806339077537146103035780634c5ef0ed1461032557600080fd5b806301ffc9a7146101d4578063181f5a77146101fc57806321df0da71461023e575b600080fd5b6101e76101e2366004613050565b610592565b60405190151581526020015b60405180910390f35b60408051808201909152601f81527f4275726e5769746846726f6d4d696e74546f6b656e506f6f6c20312e352e310060208201525b6040516101f391906130f6565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b6101e761029336600461312b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f3565b610316610311366004613148565b610677565b604051905181526020016101f3565b6101e76103333660046131a1565b610846565b61034b610346366004613270565b610890565b005b61034b61035b3660046131a1565b61090b565b60095473ffffffffffffffffffffffffffffffffffffffff16610260565b61034b6109a8565b61034b61039436600461312b565b610a76565b6101e76103a73660046132dc565b610af7565b60015473ffffffffffffffffffffffffffffffffffffffff16610260565b6103dd6103d83660046132f7565b610b0e565b6040516101f39190613332565b6103fd6103f83660046132dc565b610be7565b6040516101f39190613389565b610412610d52565b6040516101f3919061340b565b61034b61042d3660046131a1565b610d63565b6104456104403660046132dc565b610e7b565b6040516101f3919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610260565b6102316104cd3660046132dc565b610f50565b61034b6104e036600461312b565b611000565b6104ed6110db565b6040516101f39190613465565b6104456105083660046132dc565b611193565b61034b61051b3660046135ed565b611265565b7f0000000000000000000000000000000000000000000000000000000000000000610260565b7f00000000000000000000000000000000000000000000000000000000000000006101e7565b61034b61057a366004613270565b6112e9565b61034b61058d36600461312b565b6117fb565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061062557507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061067157507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261068f8261180f565b60006106e860608401356106e36106a960c0870187613632565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611a3392505050565b611af7565b905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f19610736606086016040870161312b565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401600060405180830381600087803b1580156107a357600080fd5b505af11580156107b7573d6000803e3d6000fd5b506107cc92505050606084016040850161312b565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f08360405161082a91815260200190565b60405180910390a3604080516020810190915290815292915050565b6000610888838360405161085b929190613697565b604080519182900390912067ffffffffffffffff8716600090815260076020529190912060050190611be7565b949350505050565b610898611bff565b61090584848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611c5292505050565b50505050565b610913611bff565b61091c83610af7565b610963576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6109a38383838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611e0892505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146109f9576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a7e611bff565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b6000610671600567ffffffffffffffff8416611be7565b6040805180820190915260608082526020820152610b2b82611f02565b610b38826060013561208e565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610b928460200160208101906104cd91906132dc565b8152602001610bdf6040805160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260609101604051602081830303815290604052905090565b905292915050565b67ffffffffffffffff8116600090815260076020526040812060609190610c1090600501612130565b90506000815167ffffffffffffffff811115610c2e57610c2e6134a7565b604051908082528060200260200182016040528015610c6157816020015b6060815260200190600190039081610c4c5790505b50905060005b8251811015610d4a5760086000848381518110610c8657610c866136a7565b602002602001015181526020019081526020016000208054610ca7906136d6565b80601f0160208091040260200160405190810160405280929190818152602001828054610cd3906136d6565b8015610d205780601f10610cf557610100808354040283529160200191610d20565b820191906000526020600020905b815481529060010190602001808311610d0357829003601f168201915b5050505050828281518110610d3757610d376136a7565b6020908102919091010152600101610c67565b509392505050565b6060610d5e6002612130565b905090565b610d6b611bff565b610d7483610af7565b610db6576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161095a565b610df68282604051610dc9929190613697565b604080519182900390912067ffffffffffffffff861660009081526007602052919091206005019061213d565b610e32578282826040517f74f23c7c00000000000000000000000000000000000000000000000000000000815260040161095a93929190613772565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610e6e929190613796565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261067190612149565b67ffffffffffffffff81166000908152600760205260409020600401805460609190610f7b906136d6565b80601f0160208091040260200160405190810160405280929190818152602001828054610fa7906136d6565b8015610ff45780601f10610fc957610100808354040283529160200191610ff4565b820191906000526020600020905b815481529060010190602001808311610fd757829003601f168201915b50505050509050919050565b611008611bff565b73ffffffffffffffffffffffffffffffffffffffff8116611055576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110e96005612130565b90506000815167ffffffffffffffff811115611107576111076134a7565b604051908082528060200260200182016040528015611130578160200160208202803683370190505b50905060005b825181101561118c57828181518110611151576111516136a7565b602002602001015182828151811061116b5761116b6136a7565b67ffffffffffffffff90921660209283029190910190910152600101611136565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261067190612149565b60095473ffffffffffffffffffffffffffffffffffffffff1633148015906112a5575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156112de576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b6109a38383836121fb565b6112f1611bff565b60005b838110156114de576000858583818110611310576113106136a7565b905060200201602081019061132591906132dc565b905061133c600567ffffffffffffffff831661213d565b61137e576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b67ffffffffffffffff811660009081526007602052604081206113a390600501612130565b905060005b815181101561140f576114068282815181106113c6576113c66136a7565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161213d90919063ffffffff16565b506001016113a8565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906114786004830182612fe3565b600582016000818161148a828261301d565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916945060200192506114cc915050565b60405180910390a150506001016112f4565b5060005b818110156117f45760008383838181106114fe576114fe6136a7565b905060200281019061151091906137aa565b61151990613876565b905061152a816060015160006122e5565b611539816080015160006122e5565b806040015151600003611578576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516115909060059067ffffffffffffffff16612422565b6115d55780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161095a565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a90910299909316171790941695909517909255909202909117600382015590820151600482019061175890826139ed565b5060005b82602001515181101561179c57611794836000015184602001518381518110611787576117876136a7565b6020026020010151611e08565b60010161175c565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c282600001518360400151846060015185608001516040516117e29493929190613b07565b60405180910390a150506001016114e2565b5050505050565b611803611bff565b61180c8161242e565b50565b61182261029360a083016080840161312b565b6118815761183660a082016080830161312b565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161095a565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6118cd60408401602085016132dc565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa15801561193e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119629190613ba0565b15611999576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119b16119ac60408301602084016132dc565b6124f2565b6119d16119c460408301602084016132dc565b61033360a0840184613632565b611a16576119e260a0820182613632565b6040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161095a929190613796565b61180c611a2960408301602084016132dc565b8260600135612618565b60008151600003611a6557507f0000000000000000000000000000000000000000000000000000000000000000919050565b8151602014611aa257816040517f953576f700000000000000000000000000000000000000000000000000000000815260040161095a91906130f6565b600082806020019051810190611ab89190613bbd565b905060ff81111561067157826040517f953576f700000000000000000000000000000000000000000000000000000000815260040161095a91906130f6565b60007f000000000000000000000000000000000000000000000000000000000000000060ff168260ff1603611b2d575081610671565b7f000000000000000000000000000000000000000000000000000000000000000060ff168260ff161115611ba157611b857f000000000000000000000000000000000000000000000000000000000000000083613c05565b611b9090600a613d3e565b611b9a9084613d4d565b9050610671565b611bcb827f0000000000000000000000000000000000000000000000000000000000000000613c05565b611bd690600a613d3e565b611be09084613d88565b9392505050565b60008181526001830160205260408120541515611be0565b60015473ffffffffffffffffffffffffffffffffffffffff163314611c50576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f0000000000000000000000000000000000000000000000000000000000000000611ca9576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611d3f576000838281518110611cc957611cc96136a7565b60200260200101519050611ce781600261265f90919063ffffffff16565b15611d365760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611cac565b5060005b81518110156109a3576000828281518110611d6057611d606136a7565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611da45750611e00565b611daf600282612681565b15611dfe5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611d43565b8051600003611e43576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff8416600090815260079092526040909120611e759060050182612422565b611eaf5782826040517f393b8ad200000000000000000000000000000000000000000000000000000000815260040161095a929190613d9f565b6000818152600860205260409020611ec783826139ed565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea83604051610e6e91906130f6565b611f1561029360a083016080840161312b565b611f295761183660a082016080830161312b565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb611f7560408401602085016132dc565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015611fe6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061200a9190613ba0565b15612041576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612059612054606083016040840161312b565b6126a3565b61207161206c60408301602084016132dc565b612722565b61180c61208460408301602084016132dc565b8260600135612870565b6040517f9dc29fac000000000000000000000000000000000000000000000000000000008152306004820152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639dc29fac90604401600060405180830381600087803b15801561211c57600080fd5b505af11580156117f4573d6000803e3d6000fd5b60606000611be0836128b4565b6000611be0838361290f565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526121d782606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426121bb9190613dc2565b85608001516fffffffffffffffffffffffffffffffff16612a02565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61220483610af7565b612246576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161095a565b6122518260006122e5565b67ffffffffffffffff831660009081526007602052604090206122749083612a2a565b61227f8160006122e5565b67ffffffffffffffff831660009081526007602052604090206122a59060020182612a2a565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516122d893929190613dd5565b60405180910390a1505050565b8151156123b05781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff1610158061233b575060408201516fffffffffffffffffffffffffffffffff16155b1561237457816040517f8020d12400000000000000000000000000000000000000000000000000000000815260040161095a9190613e58565b80156123ac576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff161515806123e9575060208201516fffffffffffffffffffffffffffffffff1615155b156123ac57816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161095a9190613e58565b6000611be08383612bcc565b3373ffffffffffffffffffffffffffffffffffffffff82160361247d576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6124fb81610af7565b61253d576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156125bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125e09190613ba0565b61180c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b67ffffffffffffffff821660009081526007602052604090206123ac90600201827f0000000000000000000000000000000000000000000000000000000000000000612c1b565b6000611be08373ffffffffffffffffffffffffffffffffffffffff841661290f565b6000611be08373ffffffffffffffffffffffffffffffffffffffff8416612bcc565b7f00000000000000000000000000000000000000000000000000000000000000001561180c576126d4600282612f9e565b61180c576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161095a565b61272b81610af7565b61276d576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa1580156127e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061280a9190613e94565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461180c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b67ffffffffffffffff821660009081526007602052604090206123ac90827f0000000000000000000000000000000000000000000000000000000000000000612c1b565b606081600001805480602002602001604051908101604052809291908181526020018280548015610ff457602002820191906000526020600020905b8154815260200190600101908083116128f05750505050509050919050565b600081815260018301602052604081205480156129f8576000612933600183613dc2565b855490915060009061294790600190613dc2565b90508082146129ac576000866000018281548110612967576129676136a7565b906000526020600020015490508087600001848154811061298a5761298a6136a7565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806129bd576129bd613eb1565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610671565b6000915050610671565b6000612a2185612a128486613d88565b612a1c9087613ee0565b612fcd565b95945050505050565b8154600090612a5390700100000000000000000000000000000000900463ffffffff1642613dc2565b90508015612af55760018301548354612a9b916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612a02565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612b1b916fffffffffffffffffffffffffffffffff9081169116612fcd565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906122d8908490613e58565b6000818152600183016020526040812054612c1357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610671565b506000610671565b825474010000000000000000000000000000000000000000900460ff161580612c42575081155b15612c4c57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612c9290700100000000000000000000000000000000900463ffffffff1642613dc2565b90508015612d525781831115612cd4576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612d0e9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612a02565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612e095773ffffffffffffffffffffffffffffffffffffffff8416612db1576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161095a565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161095a565b84831015612f1c5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612e4d9082613dc2565b612e57878a613dc2565b612e619190613ee0565b612e6b9190613d4d565b905073ffffffffffffffffffffffffffffffffffffffff8616612ec4576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161095a565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161095a565b612f268584613dc2565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611be0565b6000818310612fdc5781611be0565b5090919050565b508054612fef906136d6565b6000825580601f10612fff575050565b601f01602090049060005260206000209081019061180c9190613037565b508054600082559060005260206000209081019061180c91905b5b8082111561304c5760008155600101613038565b5090565b60006020828403121561306257600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611be057600080fd5b6000815180845260005b818110156130b85760208185018101518683018201520161309c565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611be06020830184613092565b73ffffffffffffffffffffffffffffffffffffffff8116811461180c57600080fd5b60006020828403121561313d57600080fd5b8135611be081613109565b60006020828403121561315a57600080fd5b813567ffffffffffffffff81111561317157600080fd5b82016101008185031215611be057600080fd5b803567ffffffffffffffff8116811461319c57600080fd5b919050565b6000806000604084860312156131b657600080fd5b6131bf84613184565b9250602084013567ffffffffffffffff808211156131dc57600080fd5b818601915086601f8301126131f057600080fd5b8135818111156131ff57600080fd5b87602082850101111561321157600080fd5b6020830194508093505050509250925092565b60008083601f84011261323657600080fd5b50813567ffffffffffffffff81111561324e57600080fd5b6020830191508360208260051b850101111561326957600080fd5b9250929050565b6000806000806040858703121561328657600080fd5b843567ffffffffffffffff8082111561329e57600080fd5b6132aa88838901613224565b909650945060208701359150808211156132c357600080fd5b506132d087828801613224565b95989497509550505050565b6000602082840312156132ee57600080fd5b611be082613184565b60006020828403121561330957600080fd5b813567ffffffffffffffff81111561332057600080fd5b820160a08185031215611be057600080fd5b60208152600082516040602084015261334e6060840182613092565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612a218282613092565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b828110156133fe577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526133ec858351613092565b945092850192908501906001016133b2565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561345957835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613427565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561345957835167ffffffffffffffff1683529284019291840191600101613481565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156134f9576134f96134a7565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613546576135466134a7565b604052919050565b801515811461180c57600080fd5b80356fffffffffffffffffffffffffffffffff8116811461319c57600080fd5b60006060828403121561358e57600080fd5b6040516060810181811067ffffffffffffffff821117156135b1576135b16134a7565b60405290508082356135c28161354e565b81526135d06020840161355c565b60208201526135e16040840161355c565b60408201525092915050565b600080600060e0848603121561360257600080fd5b61360b84613184565b925061361a856020860161357c565b9150613629856080860161357c565b90509250925092565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261366757600080fd5b83018035915067ffffffffffffffff82111561368257600080fd5b60200191503681900382131561326957600080fd5b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c908216806136ea57607f821691505b602082108103613723577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff84168152604060208201526000612a21604083018486613729565b602081526000610888602083018486613729565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee18336030181126137de57600080fd5b9190910192915050565b600082601f8301126137f957600080fd5b813567ffffffffffffffff811115613813576138136134a7565b61384460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016134ff565b81815284602083860101111561385957600080fd5b816020850160208301376000918101602001919091529392505050565b6000610120823603121561388957600080fd5b6138916134d6565b61389a83613184565b815260208084013567ffffffffffffffff808211156138b857600080fd5b9085019036601f8301126138cb57600080fd5b8135818111156138dd576138dd6134a7565b8060051b6138ec8582016134ff565b918252838101850191858101903684111561390657600080fd5b86860192505b83831015613942578235858111156139245760008081fd5b6139323689838a01016137e8565b835250918601919086019061390c565b808789015250505050604086013592508083111561395f57600080fd5b505061396d368286016137e8565b604083015250613980366060850161357c565b60608201526139923660c0850161357c565b608082015292915050565b601f8211156109a3576000816000526020600020601f850160051c810160208610156139c65750805b601f850160051c820191505b818110156139e5578281556001016139d2565b505050505050565b815167ffffffffffffffff811115613a0757613a076134a7565b613a1b81613a1584546136d6565b8461399d565b602080601f831160018114613a6e5760008415613a385750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556139e5565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613abb57888601518255948401946001909101908401613a9c565b5085821015613af757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613b2b81840187613092565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613b699050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152612a21565b600060208284031215613bb257600080fd5b8151611be08161354e565b600060208284031215613bcf57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff828116828216039081111561067157610671613bd6565b600181815b80851115613c7757817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613c5d57613c5d613bd6565b80851615613c6a57918102915b93841c9390800290613c23565b509250929050565b600082613c8e57506001610671565b81613c9b57506000610671565b8160018114613cb15760028114613cbb57613cd7565b6001915050610671565b60ff841115613ccc57613ccc613bd6565b50506001821b610671565b5060208310610133831016604e8410600b8410161715613cfa575081810a610671565b613d048383613c1e565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613d3657613d36613bd6565b029392505050565b6000611be060ff841683613c7f565b600082613d83577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b808202811582820484141761067157610671613bd6565b67ffffffffffffffff831681526040602082015260006108886040830184613092565b8181038181111561067157610671613bd6565b67ffffffffffffffff8416815260e08101613e2160208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152610888565b6060810161067182848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060208284031215613ea657600080fd5b8151611be081613109565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b8082018082111561067157610671613bd656fea164736f6c6343000818000a",
}

var BurnWithFromMintTokenPoolABI = BurnWithFromMintTokenPoolMetaData.ABI

var BurnWithFromMintTokenPoolBin = BurnWithFromMintTokenPoolMetaData.Bin

func DeployBurnWithFromMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnWithFromMintTokenPool, error) {
	parsed, err := BurnWithFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnWithFromMintTokenPoolBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnWithFromMintTokenPool.Contract.GetTokenDecimals(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnWithFromMintTokenPool.Contract.GetTokenDecimals(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsRemotePool(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsRemotePool(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
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

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)

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
