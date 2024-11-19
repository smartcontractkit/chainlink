// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_token_pool

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

var BurnMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"remotePairHash\",\"type\":\"bytes32\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"remotePairHash\",\"type\":\"bytes32\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200474e3803806200474e8339810160408190526200003491620004e5565b83838383336000816200005a57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008d576200008d8162000140565b50506001600160a01b0384161580620000ad57506001600160a01b038116155b80620000c057506001600160a01b038216155b15620000df576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200013257604080516000815260208101909152620001329084620001ba565b505050505050505062000643565b336001600160a01b038216036200016a57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c051620001db576040516335f4a7b360e01b815260040160405180910390fd5b60005b825181101562000266576000838281518110620001ff57620001ff620005f5565b602090810291909101015190506200021960028262000317565b156200025c576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101620001de565b5060005b8151811015620003125760008282815181106200028b576200028b620005f5565b6020026020010151905060006001600160a01b0316816001600160a01b031603620002b7575062000309565b620002c460028262000337565b1562000307576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6001016200026a565b505050565b60006200032e836001600160a01b0384166200034e565b90505b92915050565b60006200032e836001600160a01b03841662000452565b6000818152600183016020526040812054801562000447576000620003756001836200060b565b85549091506000906200038b906001906200060b565b9050808214620003f7576000866000018281548110620003af57620003af620005f5565b9060005260206000200154905080876000018481548110620003d557620003d5620005f5565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200040b576200040b6200062d565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000331565b600091505062000331565b60008181526001830160205260408120546200049b5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000331565b50600062000331565b6001600160a01b0381168114620004ba57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b8051620004e081620004a4565b919050565b60008060008060808587031215620004fc57600080fd5b84516200050981620004a4565b602086810151919550906001600160401b03808211156200052957600080fd5b818801915088601f8301126200053e57600080fd5b815181811115620005535762000553620004bd565b8060051b604051601f19603f830116810181811085821117156200057b576200057b620004bd565b60405291825284820192508381018501918b8311156200059a57600080fd5b938501935b82851015620005c357620005b385620004d3565b845293850193928501926200059f565b809850505050505050620005da60408601620004d3565b9150620005ea60608601620004d3565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b818103818111156200033157634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05161408e620006c06000396000818161050801528181611f1501526128a90152600081816104e201528181611c0201526121c8015260008181610231015281816102860152818161065b01528181611b22015281816120e8015281816122da015281816126fd0152612a94015261408e6000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80639a4575b9116100f9578063c0d7865511610097578063cf7401f311610071578063cf7401f3146104cd578063dc0bd971146104e0578063e0351e1314610506578063f2fde38b1461052c57600080fd5b8063c0d7865514610492578063c4bffe2b146104a5578063c75eea9c146104ba57600080fd5b8063acfecf91116100d3578063acfecf91146103df578063af58d59f146103f2578063b0f479a114610461578063b79465801461047f57600080fd5b80639a4575b91461038a578063a42a7b8b146103aa578063a7cd63b7146103ca57600080fd5b806354c8a4f31161016657806379ba50971161014057806379ba50971461033e5780637d54534e146103465780638926f54f146103595780638da5cb5b1461036c57600080fd5b806354c8a4f3146102fa57806362ddd3c41461030d5780636d3d1a581461032057600080fd5b8063240028e811610197578063240028e81461027657806339077537146102c35780634d5e7b47146102e557600080fd5b806301ffc9a7146101be578063181f5a77146101e657806321df0da71461022f575b600080fd5b6101d16101cc366004613187565b61053f565b60405190151581526020015b60405180910390f35b6102226040518060400160405280601781526020017f4275726e4d696e74546f6b656e506f6f6c20312e352e3000000000000000000081525081565b6040516101dd919061322d565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101dd565b6101d1610284366004613272565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102d66102d136600461328f565b610624565b604051905181526020016101dd565b6102f86102f3366004613317565b6107aa565b005b6102f8610308366004613317565b610e34565b6102f861031b36600461339b565b610eaf565b600a5473ffffffffffffffffffffffffffffffffffffffff16610251565b6102f861104c565b6102f8610354366004613272565b61111a565b6101d161036736600461341e565b61119b565b60015473ffffffffffffffffffffffffffffffffffffffff16610251565b61039d610398366004613439565b6111b2565b6040516101dd9190613474565b6103bd6103b836600461341e565b611259565b6040516101dd91906134d4565b6103d261134c565b6040516101dd9190613556565b6102f86103ed36600461339b565b61135d565b61040561040036600461341e565b611694565b6040516101dd919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610251565b61022261048d36600461341e565b611769565b6102f86104a0366004613272565b611819565b6104ad6118f4565b6040516101dd91906135b0565b6104056104c836600461341e565b6119ac565b6102f86104db36600461370d565b611a7e565b7f0000000000000000000000000000000000000000000000000000000000000000610251565b7f00000000000000000000000000000000000000000000000000000000000000006101d1565b6102f861053a366004613272565b611b07565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf0000000000000000000000000000000000000000000000000000000014806105d257507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061061e57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261064461063f836137fd565b611b1b565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f196106906060850160408601613272565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260608501356024820152604401600060405180830381600087803b15801561070057600080fd5b505af1158015610714573d6000803e3d6000fd5b50610729925050506060830160408401613272565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0846060013560405161078b91815260200190565b60405180910390a3506040805160208101909152606090910135815290565b6107b2611d31565b60005b83811015610a335760008585838181106107d1576107d16138f2565b90506020020160208101906107e6919061341e565b90506107fd600567ffffffffffffffff8316611d84565b610844576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201526024015b60405180910390fd5b67ffffffffffffffff8116600090815260076020908152604080832060050180548251818502810185019093528083529192909190849084015b8282101561092a57838290600052602060002001805461089d90613921565b80601f01602080910402602001604051908101604052809291908181526020018280546108c990613921565b80156109165780601f106108eb57610100808354040283529160200191610916565b820191906000526020600020905b8154815290600101906020018083116108f957829003601f168201915b50505050508152602001906001019061087e565b50505050905060005b81518110156109735761096a61096284848481518110610955576109556138f2565b6020026020010151611d97565b600890611d84565b50600101610933565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906109dc60048301826130a7565b6109ea6005830160006130e1565b505060405167ffffffffffffffff831681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a150506001016107b5565b5060005b81811015610e2d576000838383818110610a5357610a536138f2565b9050602002810190610a659190613974565b610a6e906139b2565b9050610a7f81606001516000611dca565b610a8e81608001516000611dca565b8051610aa69060059067ffffffffffffffff16611f07565b610aeb5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161083b565b6020810151511580610b005750604081015151155b15610b37576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252606083810180516020908101516fffffffffffffffffffffffffffffffff90811660808087019182524263ffffffff1660a0808901829052865151151560c08a01528651860151851660e08a0152955189015184166101008901529187528751948501885280890180518501518416865285850192909252815151151585890152815184015183168587015290518701519091169083015280840191909152848401518385015283516000808252918101909452919290830191610c16565b6060815260200190600190039081610c015790505b509052815167ffffffffffffffff1660009081526007602090815260409182902083518051825482850151838701511515740100000000000000000000000000000000000000009081027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff63ffffffff9384167001000000000000000000000000000000009081027fffffffffffffffffffffffff00000000000000000000000000000000000000009687166fffffffffffffffffffffffffffffffff98891617178216929092178855606080880151608098890151881684029088161760018a01558b8a0151805160028b0180549c830151838f015115159097029690971685029b909716968816969096179990991716919091179091559481015192015181169093029216919091176003820155908201516004820190610d599082613aa5565b5060608201518051610d759160058401916020909101906130ff565b5050815167ffffffffffffffff1660009081526007602090815260408220818501516005909101805460018101825590845291909220019150610db89082613aa5565b50610dd6610dce82600001518360200151611d97565b600890611f07565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28160000151826040015183606001518460800151604051610e1c9493929190613bbf565b60405180910390a150600101610a37565b5050505050565b610e3c611d31565b610ea984848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611f1392505050565b50505050565b610eb7611d31565b610ec08361119b565b610f02576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161083b565b6000819003610f3d576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610f7f8484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611d9792505050565b9050610f8c600882611f07565b610fc8578383836040517f393b8ad200000000000000000000000000000000000000000000000000000000815260040161083b93929190613ca1565b67ffffffffffffffff8416600090815260076020908152604082206005018054600181018255908352912001610fff838583613cc5565b508367ffffffffffffffff167f3c7ec26195bfbae0c857332cd648e0121ebc1a16c59eec923f3158b48f1261a384848460405161103e93929190613ddf565b60405180910390a250505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461109d576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611122611d31565b600a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b600061061e600567ffffffffffffffff84166120c9565b60408051808201909152606080825260208201526111d76111d283613e03565b6120e1565b6111e482606001356122ab565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a2604051806040016040528061123e84602001602081019061048d919061341e565b81526040805160208181019092526000815291015292915050565b67ffffffffffffffff81166000908152600760209081526040808320600501805482518185028101850190935280835260609492939192909184015b828210156113415783829060005260206000200180546112b490613921565b80601f01602080910402602001604051908101604052809291908181526020018280546112e090613921565b801561132d5780601f106113025761010080835404028352916020019161132d565b820191906000526020600020905b81548152906001019060200180831161131057829003601f168201915b505050505081526020019060010190611295565b505050509050919050565b60606113586002612347565b905090565b611365611d31565b61136e8361119b565b6113b0576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161083b565b60006113f28484848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611d9792505050565b90506113ff600882611d84565b61143b578383836040517f74f23c7c00000000000000000000000000000000000000000000000000000000815260040161083b93929190613ca1565b67ffffffffffffffff8416600090815260076020908152604080832060050180548251818502810185019093528083529192909190849084015b8282101561152157838290600052602060002001805461149490613921565b80601f01602080910402602001604051908101604052809291908181526020018280546114c090613921565b801561150d5780601f106114e25761010080835404028352916020019161150d565b820191906000526020600020905b8154815290600101906020018083116114f057829003601f168201915b505050505081526020019060010190611475565b50505050905060008484604051611539929190613e7f565b6040518091039020905060005b82518110156116455781838281518110611562576115626138f2565b6020026020010151805190602001200361163d5782600184516115859190613ebe565b81518110611595576115956138f2565b6020026020010151600760008967ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050182815481106115d8576115d86138f2565b9060005260206000200190816115ee9190613aa5565b5067ffffffffffffffff8716600090815260076020526040902060050180548061161a5761161a613ed1565b60019003818190600052602060002001600061163691906130a7565b9055611645565b600101611546565b508567ffffffffffffffff167fb50f71ac1e206d236b676fba3b59ebc718039069aa6771810520848abedd68aa86868660405161168493929190613ddf565b60405180910390a2505050505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261061e90612354565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061179490613921565b80601f01602080910402602001604051908101604052809291908181526020018280546117c090613921565b801561180d5780601f106117e25761010080835404028352916020019161180d565b820191906000526020600020905b8154815290600101906020018083116117f057829003601f168201915b50505050509050919050565b611821611d31565b73ffffffffffffffffffffffffffffffffffffffff811661186e576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006119026005612347565b90506000815167ffffffffffffffff811115611920576119206135f2565b604051908082528060200260200182016040528015611949578160200160208202803683370190505b50905060005b82518110156119a55782818151811061196a5761196a6138f2565b6020026020010151828281518110611984576119846138f2565b67ffffffffffffffff9092166020928302919091019091015260010161194f565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261061e90612354565b600a5473ffffffffffffffffffffffffffffffffffffffff163314801590611abe575060015473ffffffffffffffffffffffffffffffffffffffff163314155b15611af7576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161083b565b611b02838383612406565b505050565b611b0f611d31565b611b18816124f0565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611bb05760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161083b565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611c5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c829190613f00565b15611cb9576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611cc681602001516125b4565b611ce3611cdb82602001518360a00151611d97565b6008906120c9565b611d1f578060a001516040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161083b919061322d565b611b18816020015182606001516126da565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d82576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000611d908383612721565b9392505050565b60008282604051602001611dac929190613f1d565b60405160208183030381529060405280519060200120905092915050565b815115611e955781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580611e20575060408201516fffffffffffffffffffffffffffffffff16155b15611e5957816040517f8020d12400000000000000000000000000000000000000000000000000000000815260040161083b9190613f40565b8015611e91576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580611ece575060208201516fffffffffffffffffffffffffffffffff1615155b15611e9157816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161083b9190613f40565b6000611d908383612814565b7f0000000000000000000000000000000000000000000000000000000000000000611f6a576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015612000576000838281518110611f8a57611f8a6138f2565b60200260200101519050611fa881600261286390919063ffffffff16565b15611ff75760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611f6d565b5060005b8151811015611b02576000828281518110612021576120216138f2565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361206557506120c1565b612070600282612885565b156120bf5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101612004565b60008181526001830160205260408120541515611d90565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff9081169116146121765760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161083b565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015612224573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122489190613f00565b1561227f576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61228c81604001516128a7565b6122998160200151612926565b611b1881602001518260600151612a74565b6040517f42966c68000000000000000000000000000000000000000000000000000000008152600481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906342966c6890602401600060405180830381600087803b15801561233357600080fd5b505af1158015610e2d573d6000803e3d6000fd5b60606000611d9083612ab8565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526123e282606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426123c69190613ebe565b85608001516fffffffffffffffffffffffffffffffff16612b13565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61240f8361119b565b612451576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161083b565b61245c826000611dca565b67ffffffffffffffff8316600090815260076020526040902061247f9083612b3d565b61248a816000611dca565b67ffffffffffffffff831660009081526007602052604090206124b09060020182612b3d565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516124e393929190613f7c565b60405180910390a1505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361253f576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6125bd8161119b565b6125ff576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161083b565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa15801561267e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126a29190613f00565b611b18576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161083b565b67ffffffffffffffff82166000908152600760205260409020611e9190600201827f0000000000000000000000000000000000000000000000000000000000000000612cdf565b6000818152600183016020526040812054801561280a576000612745600183613ebe565b855490915060009061275990600190613ebe565b90508082146127be576000866000018281548110612779576127796138f2565b906000526020600020015490508087600001848154811061279c5761279c6138f2565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806127cf576127cf613ed1565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061061e565b600091505061061e565b600081815260018301602052604081205461285b5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561061e565b50600061061e565b6000611d908373ffffffffffffffffffffffffffffffffffffffff8416612721565b6000611d908373ffffffffffffffffffffffffffffffffffffffff8416612814565b7f000000000000000000000000000000000000000000000000000000000000000015611b18576128d8600282613062565b611b18576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161083b565b61292f8161119b565b612971576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161083b565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa1580156129ea573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a0e9190613fff565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611b18576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161083b565b67ffffffffffffffff82166000908152600760205260409020611e9190827f0000000000000000000000000000000000000000000000000000000000000000612cdf565b60608160000180548060200260200160405190810160405280929190818152602001828054801561180d57602002820191906000526020600020905b815481526020019060010190808311612af45750505050509050919050565b6000612b3285612b23848661401c565b612b2d9087614033565b613091565b90505b949350505050565b8154600090612b6690700100000000000000000000000000000000900463ffffffff1642613ebe565b90508015612c085760018301548354612bae916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612b13565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612c2e916fffffffffffffffffffffffffffffffff9081169116613091565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906124e3908490613f40565b825474010000000000000000000000000000000000000000900460ff161580612d06575081155b15612d1057505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612d5690700100000000000000000000000000000000900463ffffffff1642613ebe565b90508015612e165781831115612d98576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612dd29083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612b13565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612ecd5773ffffffffffffffffffffffffffffffffffffffff8416612e75576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161083b565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161083b565b84831015612fe05760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612f119082613ebe565b612f1b878a613ebe565b612f259190614033565b612f2f9190614046565b905073ffffffffffffffffffffffffffffffffffffffff8616612f88576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161083b565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161083b565b612fea8584613ebe565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611d90565b60008183106130a05781611d90565b5090919050565b5080546130b390613921565b6000825580601f106130c3575050565b601f016020900490600052602060002090810190611b189190613155565b5080546000825590600052602060002090810190611b18919061316a565b828054828255906000526020600020908101928215613145579160200282015b8281111561314557825182906131359082613aa5565b509160200191906001019061311f565b5061315192915061316a565b5090565b5b808211156131515760008155600101613156565b8082111561315157600061317e82826130a7565b5060010161316a565b60006020828403121561319957600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611d9057600080fd5b6000815180845260005b818110156131ef576020818501810151868301820152016131d3565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d9060208301846131c9565b73ffffffffffffffffffffffffffffffffffffffff81168114611b1857600080fd5b803561326d81613240565b919050565b60006020828403121561328457600080fd5b8135611d9081613240565b6000602082840312156132a157600080fd5b813567ffffffffffffffff8111156132b857600080fd5b82016101008185031215611d9057600080fd5b60008083601f8401126132dd57600080fd5b50813567ffffffffffffffff8111156132f557600080fd5b6020830191508360208260051b850101111561331057600080fd5b9250929050565b6000806000806040858703121561332d57600080fd5b843567ffffffffffffffff8082111561334557600080fd5b613351888389016132cb565b9096509450602087013591508082111561336a57600080fd5b50613377878288016132cb565b95989497509550505050565b803567ffffffffffffffff8116811461326d57600080fd5b6000806000604084860312156133b057600080fd5b6133b984613383565b9250602084013567ffffffffffffffff808211156133d657600080fd5b818601915086601f8301126133ea57600080fd5b8135818111156133f957600080fd5b87602082850101111561340b57600080fd5b6020830194508093505050509250925092565b60006020828403121561343057600080fd5b611d9082613383565b60006020828403121561344b57600080fd5b813567ffffffffffffffff81111561346257600080fd5b820160a08185031215611d9057600080fd5b60208152600082516040602084015261349060608401826131c9565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526134cb82826131c9565b95945050505050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613549577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526135378583516131c9565b945092850192908501906001016134fd565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156135a457835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613572565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156135a457835167ffffffffffffffff16835292840192918401916001016135cc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff81118282101715613645576136456135f2565b60405290565b60405160a0810167ffffffffffffffff81118282101715613645576136456135f2565b8015158114611b1857600080fd5b80356fffffffffffffffffffffffffffffffff8116811461326d57600080fd5b6000606082840312156136ae57600080fd5b6040516060810181811067ffffffffffffffff821117156136d1576136d16135f2565b60405290508082356136e28161366e565b81526136f06020840161367c565b60208201526137016040840161367c565b60408201525092915050565b600080600060e0848603121561372257600080fd5b61372b84613383565b925061373a856020860161369c565b9150613749856080860161369c565b90509250925092565b600082601f83011261376357600080fd5b813567ffffffffffffffff8082111561377e5761377e6135f2565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156137c4576137c46135f2565b816040528381528660208588010111156137dd57600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000610100823603121561381057600080fd5b613818613621565b823567ffffffffffffffff8082111561383057600080fd5b61383c36838701613752565b835261384a60208601613383565b602084015261385b60408601613262565b60408401526060850135606084015261387660808601613262565b608084015260a085013591508082111561388f57600080fd5b61389b36838701613752565b60a084015260c08501359150808211156138b457600080fd5b6138c036838701613752565b60c084015260e08501359150808211156138d957600080fd5b506138e636828601613752565b60e08301525092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061393557607f821691505b60208210810361396e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee18336030181126139a857600080fd5b9190910192915050565b600061012082360312156139c557600080fd5b6139cd61364b565b6139d683613383565b8152602083013567ffffffffffffffff808211156139f357600080fd5b6139ff36838701613752565b60208401526040850135915080821115613a1857600080fd5b50613a2536828601613752565b604083015250613a38366060850161369c565b6060820152613a4a3660c0850161369c565b608082015292915050565b601f821115611b02576000816000526020600020601f850160051c81016020861015613a7e5750805b601f850160051c820191505b81811015613a9d57828155600101613a8a565b505050505050565b815167ffffffffffffffff811115613abf57613abf6135f2565b613ad381613acd8454613921565b84613a55565b602080601f831160018114613b265760008415613af05750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613a9d565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613b7357888601518255948401946001909101908401613b54565b5085821015613baf57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613be3818401876131c9565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613c219050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e08301526134cb565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff84168152604060208201526000612b32604083018486613c58565b67ffffffffffffffff831115613cdd57613cdd6135f2565b613cf183613ceb8354613921565b83613a55565b6000601f841160018114613d435760008515613d0d5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355610e2d565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015613d925786850135825560209485019460019092019101613d72565b5086821015613dcd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b604081526000613df3604083018587613c58565b9050826020830152949350505050565b600060a08236031215613e1557600080fd5b613e1d61364b565b823567ffffffffffffffff811115613e3457600080fd5b613e4036828601613752565b825250613e4f60208401613383565b60208201526040830135613e6281613240565b6040820152606083810135908201526080830135613a4a81613240565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561061e5761061e613e8f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600060208284031215613f1257600080fd5b8151611d908161366e565b67ffffffffffffffff83168152604060208201526000612b3560408301846131c9565b6060810161061e82848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b67ffffffffffffffff8416815260e08101613fc860208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152612b35565b60006020828403121561401157600080fd5b8151611d9081613240565b808202811582820484141761061e5761061e613e8f565b8082018082111561061e5761061e613e8f565b60008261407c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea164736f6c6343000818000a",
}

var BurnMintTokenPoolABI = BurnMintTokenPoolMetaData.ABI

var BurnMintTokenPoolBin = BurnMintTokenPoolMetaData.Bin

func DeployBurnMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnMintTokenPool, error) {
	parsed, err := BurnMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintTokenPoolBin), backend, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintTokenPool{address: address, abi: *parsed, BurnMintTokenPoolCaller: BurnMintTokenPoolCaller{contract: contract}, BurnMintTokenPoolTransactor: BurnMintTokenPoolTransactor{contract: contract}, BurnMintTokenPoolFilterer: BurnMintTokenPoolFilterer{contract: contract}}, nil
}

type BurnMintTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnMintTokenPoolCaller
	BurnMintTokenPoolTransactor
	BurnMintTokenPoolFilterer
}

type BurnMintTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolSession struct {
	Contract     *BurnMintTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolCallerSession struct {
	Contract *BurnMintTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnMintTokenPoolTransactorSession struct {
	Contract     *BurnMintTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolRaw struct {
	Contract *BurnMintTokenPool
}

type BurnMintTokenPoolCallerRaw struct {
	Contract *BurnMintTokenPoolCaller
}

type BurnMintTokenPoolTransactorRaw struct {
	Contract *BurnMintTokenPoolTransactor
}

func NewBurnMintTokenPool(address common.Address, backend bind.ContractBackend) (*BurnMintTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPool{address: address, abi: abi, BurnMintTokenPoolCaller: BurnMintTokenPoolCaller{contract: contract}, BurnMintTokenPoolTransactor: BurnMintTokenPoolTransactor{contract: contract}, BurnMintTokenPoolFilterer: BurnMintTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnMintTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnMintTokenPoolCaller, error) {
	contract, err := bindBurnMintTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolCaller{contract: contract}, nil
}

func NewBurnMintTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintTokenPoolTransactor, error) {
	contract, err := bindBurnMintTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolTransactor{contract: contract}, nil
}

func NewBurnMintTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintTokenPoolFilterer, error) {
	contract, err := bindBurnMintTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolFilterer{contract: contract}, nil
}

func bindBurnMintTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPool.Contract.GetAllowList(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPool.Contract.GetAllowList(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPool.Contract.GetAllowListEnabled(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPool.Contract.GetAllowListEnabled(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRateLimitAdmin(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRateLimitAdmin(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemotePools(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemotePools(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemoteToken(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemoteToken(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRmnProxy(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRmnProxy(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRouter(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRouter(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPool.Contract.GetSupportedChains(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPool.Contract.GetSupportedChains(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetToken(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetToken(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedChain(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedChain(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedToken(&_BurnMintTokenPool.CallOpts, token)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedToken(&_BurnMintTokenPool.CallOpts, token)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) Owner() (common.Address, error) {
	return _BurnMintTokenPool.Contract.Owner(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnMintTokenPool.Contract.Owner(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPool.Contract.SupportsInterface(&_BurnMintTokenPool.CallOpts, interfaceId)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPool.Contract.SupportsInterface(&_BurnMintTokenPool.CallOpts, interfaceId)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPool.Contract.TypeAndVersion(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPool.Contract.TypeAndVersion(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AcceptOwnership(&_BurnMintTokenPool.TransactOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AcceptOwnership(&_BurnMintTokenPool.TransactOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AddRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AddRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyChainUpdates(&_BurnMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyChainUpdates(&_BurnMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.LockOrBurn(&_BurnMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.LockOrBurn(&_BurnMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ReleaseOrMint(&_BurnMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ReleaseOrMint(&_BurnMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.RemoveRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.RemoveRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRateLimitAdmin(&_BurnMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRateLimitAdmin(&_BurnMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRouter(&_BurnMintTokenPool.TransactOpts, newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRouter(&_BurnMintTokenPool.TransactOpts, newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.TransferOwnership(&_BurnMintTokenPool.TransactOpts, to)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.TransferOwnership(&_BurnMintTokenPool.TransactOpts, to)
}

type BurnMintTokenPoolAllowListAddIterator struct {
	Event *BurnMintTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAllowListAdd)
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
		it.Event = new(BurnMintTokenPoolAllowListAdd)
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

func (it *BurnMintTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAllowListAddIterator{contract: _BurnMintTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAllowListAdd)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAllowListAdd, error) {
	event := new(BurnMintTokenPoolAllowListAdd)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAllowListRemoveIterator struct {
	Event *BurnMintTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAllowListRemove)
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
		it.Event = new(BurnMintTokenPoolAllowListRemove)
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

func (it *BurnMintTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAllowListRemoveIterator{contract: _BurnMintTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAllowListRemove)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAllowListRemove, error) {
	event := new(BurnMintTokenPoolAllowListRemove)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolBurnedIterator struct {
	Event *BurnMintTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolBurned)
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
		it.Event = new(BurnMintTokenPoolBurned)
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

func (it *BurnMintTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolBurnedIterator{contract: _BurnMintTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolBurned)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseBurned(log types.Log) (*BurnMintTokenPoolBurned, error) {
	event := new(BurnMintTokenPoolBurned)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainAddedIterator struct {
	Event *BurnMintTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainAdded)
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
		it.Event = new(BurnMintTokenPoolChainAdded)
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

func (it *BurnMintTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainAddedIterator{contract: _BurnMintTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainAdded)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnMintTokenPoolChainAdded, error) {
	event := new(BurnMintTokenPoolChainAdded)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainConfiguredIterator struct {
	Event *BurnMintTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainConfigured)
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
		it.Event = new(BurnMintTokenPoolChainConfigured)
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

func (it *BurnMintTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainConfiguredIterator{contract: _BurnMintTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainConfigured)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnMintTokenPoolChainConfigured, error) {
	event := new(BurnMintTokenPoolChainConfigured)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainRemovedIterator struct {
	Event *BurnMintTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainRemoved)
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
		it.Event = new(BurnMintTokenPoolChainRemoved)
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

func (it *BurnMintTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainRemovedIterator{contract: _BurnMintTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainRemoved)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnMintTokenPoolChainRemoved, error) {
	event := new(BurnMintTokenPoolChainRemoved)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolConfigChangedIterator struct {
	Event *BurnMintTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolConfigChanged)
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
		it.Event = new(BurnMintTokenPoolConfigChanged)
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

func (it *BurnMintTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolConfigChangedIterator{contract: _BurnMintTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolConfigChanged)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnMintTokenPoolConfigChanged, error) {
	event := new(BurnMintTokenPoolConfigChanged)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolLockedIterator struct {
	Event *BurnMintTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolLocked)
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
		it.Event = new(BurnMintTokenPoolLocked)
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

func (it *BurnMintTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolLockedIterator{contract: _BurnMintTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolLocked)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseLocked(log types.Log) (*BurnMintTokenPoolLocked, error) {
	event := new(BurnMintTokenPoolLocked)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolMintedIterator struct {
	Event *BurnMintTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolMinted)
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
		it.Event = new(BurnMintTokenPoolMinted)
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

func (it *BurnMintTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolMintedIterator{contract: _BurnMintTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolMinted)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseMinted(log types.Log) (*BurnMintTokenPoolMinted, error) {
	event := new(BurnMintTokenPoolMinted)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnMintTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnMintTokenPoolOwnershipTransferRequested)
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

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolOwnershipTransferRequestedIterator{contract: _BurnMintTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolOwnershipTransferRequested)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnMintTokenPoolOwnershipTransferRequested)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolOwnershipTransferredIterator struct {
	Event *BurnMintTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnMintTokenPoolOwnershipTransferred)
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

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolOwnershipTransferredIterator{contract: _BurnMintTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolOwnershipTransferred)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolOwnershipTransferred, error) {
	event := new(BurnMintTokenPoolOwnershipTransferred)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnMintTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnMintTokenPoolRateLimitAdminSet)
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

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRateLimitAdminSetIterator{contract: _BurnMintTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRateLimitAdminSet)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnMintTokenPoolRateLimitAdminSet, error) {
	event := new(BurnMintTokenPoolRateLimitAdminSet)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolReleasedIterator struct {
	Event *BurnMintTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolReleased)
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
		it.Event = new(BurnMintTokenPoolReleased)
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

func (it *BurnMintTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolReleasedIterator{contract: _BurnMintTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolReleased)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseReleased(log types.Log) (*BurnMintTokenPoolReleased, error) {
	event := new(BurnMintTokenPoolReleased)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRemotePoolAddedIterator struct {
	Event *BurnMintTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnMintTokenPoolRemotePoolAdded)
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

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	RemotePairHash      [32]byte
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRemotePoolAddedIterator{contract: _BurnMintTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRemotePoolAdded)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnMintTokenPoolRemotePoolAdded, error) {
	event := new(BurnMintTokenPoolRemotePoolAdded)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnMintTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnMintTokenPoolRemotePoolRemoved)
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

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	RemotePairHash      [32]byte
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRemotePoolRemovedIterator{contract: _BurnMintTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRemotePoolRemoved)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnMintTokenPoolRemotePoolRemoved, error) {
	event := new(BurnMintTokenPoolRemotePoolRemoved)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRouterUpdatedIterator struct {
	Event *BurnMintTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRouterUpdated)
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
		it.Event = new(BurnMintTokenPoolRouterUpdated)
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

func (it *BurnMintTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRouterUpdatedIterator{contract: _BurnMintTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRouterUpdated)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolRouterUpdated, error) {
	event := new(BurnMintTokenPoolRouterUpdated)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolTokensConsumedIterator struct {
	Event *BurnMintTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolTokensConsumed)
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
		it.Event = new(BurnMintTokenPoolTokensConsumed)
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

func (it *BurnMintTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolTokensConsumedIterator{contract: _BurnMintTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolTokensConsumed)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolTokensConsumed, error) {
	event := new(BurnMintTokenPoolTokensConsumed)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintTokenPool *BurnMintTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnMintTokenPool.ParseAllowListAdd(log)
	case _BurnMintTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnMintTokenPool.ParseAllowListRemove(log)
	case _BurnMintTokenPool.abi.Events["Burned"].ID:
		return _BurnMintTokenPool.ParseBurned(log)
	case _BurnMintTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnMintTokenPool.ParseChainAdded(log)
	case _BurnMintTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnMintTokenPool.ParseChainConfigured(log)
	case _BurnMintTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnMintTokenPool.ParseChainRemoved(log)
	case _BurnMintTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnMintTokenPool.ParseConfigChanged(log)
	case _BurnMintTokenPool.abi.Events["Locked"].ID:
		return _BurnMintTokenPool.ParseLocked(log)
	case _BurnMintTokenPool.abi.Events["Minted"].ID:
		return _BurnMintTokenPool.ParseMinted(log)
	case _BurnMintTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnMintTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintTokenPool.ParseOwnershipTransferred(log)
	case _BurnMintTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnMintTokenPool.ParseRateLimitAdminSet(log)
	case _BurnMintTokenPool.abi.Events["Released"].ID:
		return _BurnMintTokenPool.ParseReleased(log)
	case _BurnMintTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnMintTokenPool.ParseRemotePoolAdded(log)
	case _BurnMintTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnMintTokenPool.ParseRemotePoolRemoved(log)
	case _BurnMintTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnMintTokenPool.ParseRouterUpdated(log)
	case _BurnMintTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnMintTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnMintTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnMintTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnMintTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnMintTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnMintTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnMintTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnMintTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnMintTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnMintTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnMintTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnMintTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x3c7ec26195bfbae0c857332cd648e0121ebc1a16c59eec923f3158b48f1261a3")
}

func (BurnMintTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0xb50f71ac1e206d236b676fba3b59ebc718039069aa6771810520848abedd68aa")
}

func (BurnMintTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnMintTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnMintTokenPool *BurnMintTokenPool) Address() common.Address {
	return _BurnMintTokenPool.address
}

type BurnMintTokenPoolInterface interface {
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

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnMintTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnMintTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnMintTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnMintTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnMintTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnMintTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnMintTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnMintTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnMintTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnMintTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnMintTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
