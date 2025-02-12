// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_with_from_mint_rebasing_token_pool

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

var BurnWithFromMintRebasingTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountBurned\",\"type\":\"uint256\"}],\"name\":\"NegativeMintAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004742380380620047428339810160408190526200003491620008c8565b83838383838383833380600081620000935760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c657620000c68162000197565b5050506001600160a01b0384161580620000e757506001600160a01b038116155b80620000fa57506001600160a01b038216155b1562000119576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c0526200016c576040805160008152602081019091526200016c908462000242565b5062000189925050506001600160a01b038516306000196200039f565b505050505050505062000b04565b336001600160a01b03821603620001f15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c05162000263576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002ee576000838281518110620002875762000287620009d8565b60209081029190910101519050620002a160028262000485565b15620002e4576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000266565b5060005b81518110156200039a576000828281518110620003135762000313620009d8565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200033f575062000391565b6200034c600282620004a5565b156200038f576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002f2565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015620003f1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004179190620009ee565b62000423919062000a1e565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b179091529192506200047f91869190620004bc16565b50505050565b60006200049c836001600160a01b0384166200058d565b90505b92915050565b60006200049c836001600160a01b03841662000691565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564908201526000906200050b906001600160a01b038516908490620006e3565b8051909150156200039a57808060200190518101906200052c919062000a34565b6200039a5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016200008a565b6000818152600183016020526040812054801562000686576000620005b460018362000a5f565b8554909150600090620005ca9060019062000a5f565b905080821462000636576000866000018281548110620005ee57620005ee620009d8565b9060005260206000200154905080876000018481548110620006145762000614620009d8565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200064a576200064a62000a75565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200049f565b60009150506200049f565b6000818152600183016020526040812054620006da575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200049f565b5060006200049f565b6060620006f48484600085620006fc565b949350505050565b6060824710156200075f5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016200008a565b600080866001600160a01b031685876040516200077d919062000ab1565b60006040518083038185875af1925050503d8060008114620007bc576040519150601f19603f3d011682016040523d82523d6000602084013e620007c1565b606091505b509092509050620007d587838387620007e0565b979650505050505050565b60608315620008545782516000036200084c576001600160a01b0385163b6200084c5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016200008a565b5081620006f4565b620006f483838151156200086b5781518083602001fd5b8060405162461bcd60e51b81526004016200008a919062000acf565b6001600160a01b03811681146200089d57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b8051620008c38162000887565b919050565b60008060008060808587031215620008df57600080fd5b8451620008ec8162000887565b602086810151919550906001600160401b03808211156200090c57600080fd5b818801915088601f8301126200092157600080fd5b815181811115620009365762000936620008a0565b8060051b604051601f19603f830116810181811085821117156200095e576200095e620008a0565b60405291825284820192508381018501918b8311156200097d57600080fd5b938501935b82851015620009a6576200099685620008b6565b8452938501939285019262000982565b809850505050505050620009bd60408601620008b6565b9150620009cd60608601620008b6565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b60006020828403121562000a0157600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b808201808211156200049f576200049f62000a08565b60006020828403121562000a4757600080fd5b8151801515811462000a5857600080fd5b9392505050565b818103818111156200049f576200049f62000a08565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000aa857818101518382015260200162000a8e565b50506000910152565b6000825162000ac581846020870162000a8b565b9190910192915050565b602081526000825180602084015262000af081604085016020870162000a8b565b601f01601f19169190910160400192915050565b60805160a05160c051613bb362000b8f600039600081816104a901528181611958015261233c015260008181610483015281816117890152611c0e0152600081816102050152818161025a015281816106ca015281816107a50152818161087b015281816116a901528181611b2e01528181611d26015281816122d201526125270152613bb36000f3fe608060405234801561001057600080fd5b50600436106101ae5760003560e01c80639a4575b9116100ee578063c4bffe2b11610097578063db6327dc11610071578063db6327dc1461046e578063dc0bd97114610481578063e0351e13146104a7578063f2fde38b146104cd57600080fd5b8063c4bffe2b14610433578063c75eea9c14610448578063cf7401f31461045b57600080fd5b8063b0f479a1116100c8578063b0f479a1146103ef578063b79465801461040d578063c0d786551461042057600080fd5b80639a4575b91461034b578063a7cd63b71461036b578063af58d59f1461038057600080fd5b806354c8a4f31161015b57806379ba50971161013557806379ba5097146102ff5780637d54534e146103075780638926f54f1461031a5780638da5cb5b1461032d57600080fd5b806354c8a4f3146102b95780636d3d1a58146102ce57806378a010b2146102ec57600080fd5b806321df0da71161018c57806321df0da714610203578063240028e81461024a578063390775371461029757600080fd5b806301ffc9a7146101b35780630a2fd493146101db578063181f5a77146101fb575b600080fd5b6101c66101c1366004612cca565b6104e0565b60405190151581526020015b60405180910390f35b6101ee6101e9366004612d29565b6105c5565b6040516101d29190612da8565b6101ee610675565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101d2565b6101c6610258366004612de8565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b6102aa6102a5366004612e05565b610691565b604051905181526020016101d2565b6102cc6102c7366004612e8d565b610a15565b005b60085473ffffffffffffffffffffffffffffffffffffffff16610225565b6102cc6102fa366004612ef9565b610a90565b6102cc610bff565b6102cc610315366004612de8565b610cfc565b6101c6610328366004612d29565b610d4b565b60005473ffffffffffffffffffffffffffffffffffffffff16610225565b61035e610359366004612f7c565b610d62565b6040516101d29190612fb7565b610373610e09565b6040516101d29190613017565b61039361038e366004612d29565b610e1a565b6040516101d2919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610225565b6101ee61041b366004612d29565b610eef565b6102cc61042e366004612de8565b610f1a565b61043b610ff5565b6040516101d29190613071565b610393610456366004612d29565b6110ad565b6102cc6104693660046131d9565b61117f565b6102cc61047c36600461321e565b611208565b7f0000000000000000000000000000000000000000000000000000000000000000610225565b7f00000000000000000000000000000000000000000000000000000000000000006101c6565b6102cc6104db366004612de8565b61168e565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061057357507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b806105bf57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b67ffffffffffffffff811660009081526007602052604090206004018054606091906105f090613260565b80601f016020809104026020016040519081016040528092919081815260200182805461061c90613260565b80156106695780601f1061063e57610100808354040283529160200191610669565b820191906000526020600020905b81548152906001019060200180831161064c57829003601f168201915b50505050509050919050565b604051806060016040528060278152602001613b806027913981565b6040805160208101909152600081526106b16106ac8361335e565b6116a2565b600073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166370a082316106ff6060860160408701612de8565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401602060405180830381865afa158015610768573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078c9190613453565b905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f196107da6060860160408701612de8565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260608601356024820152604401600060405180830381600087803b15801561084a57600080fd5b505af115801561085e573d6000803e3d6000fd5b50600092505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690506370a082316108b26060870160408801612de8565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401602060405180830381865afa15801561091b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061093f9190613453565b90508181101561099157610953818361349b565b6040517f02164a2d00000000000000000000000000000000000000000000000000000000815260040161098891815260200190565b60405180910390fd5b6109a16060850160408601612de8565b73ffffffffffffffffffffffffffffffffffffffff16337f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f06109e3858561349b565b60405190815260200160405180910390a360405180602001604052808383610a0b919061349b565b9052949350505050565b610a1d6118d3565b610a8a8484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061195692505050565b50505050565b610a986118d3565b610aa183610d4b565b610ae3576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610988565b67ffffffffffffffff831660009081526007602052604081206004018054610b0a90613260565b80601f0160208091040260200160405190810160405280929190818152602001828054610b3690613260565b8015610b835780601f10610b5857610100808354040283529160200191610b83565b820191906000526020600020905b815481529060010190602001808311610b6657829003601f168201915b5050505067ffffffffffffffff8616600090815260076020526040902091925050600401610bb28385836134fe565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf828585604051610bf193929190613618565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610c80576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610988565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610d046118d3565b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60006105bf600567ffffffffffffffff8416611b0c565b6040805180820190915260608082526020820152610d87610d828361367c565b611b27565b610d948260600135611cf1565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610dee84602001602081019061041b9190612d29565b81526040805160208181019092526000815291015292915050565b6060610e156002611d9a565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260039091015480841660608301529190910490911660808201526105bf90611da7565b67ffffffffffffffff811660009081526007602052604090206005018054606091906105f090613260565b610f226118d3565b73ffffffffffffffffffffffffffffffffffffffff8116610f6f576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110036005611d9a565b90506000815167ffffffffffffffff811115611021576110216130b3565b60405190808252806020026020018201604052801561104a578160200160208202803683370190505b50905060005b82518110156110a65782818151811061106b5761106b61371e565b60200260200101518282815181106110855761108561371e565b67ffffffffffffffff90921660209283029190910190910152600101611050565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526105bf90611da7565b60085473ffffffffffffffffffffffffffffffffffffffff1633148015906111bf575060005473ffffffffffffffffffffffffffffffffffffffff163314155b156111f8576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610988565b611203838383611e59565b505050565b6112106118d3565b60005b8181101561120357600083838381811061122f5761122f61371e565b9050602002810190611241919061374d565b61124a9061378b565b905061125f8160800151826020015115611f43565b6112728160a00151826020015115611f43565b80602001511561156e5780516112949060059067ffffffffffffffff1661207c565b6112d95780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610988565b60408101515115806112ee5750606081015151155b15611325576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161012081018252608083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff90811660a0808901829052865151151560c08a01528651860151851660e08a015295518901518416610100890152918752875180860189529489018051850151841686528585019290925281515115158589015281518401518316606080870191909152915188015183168587015283870194855288880151878901908152828a015183890152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff90811674010000000000000000000000000000000000000000941515850217865584890151948d0151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c17909116961515029590951790985590810151940151938116931690910291909117600382015591519091906004820190611506908261383f565b506060820151600582019061151b908261383f565b505081516060830151608084015160a08501516040517f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c295506115619493929190613959565b60405180910390a1611685565b80516115869060059067ffffffffffffffff16612088565b6115cb5780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610988565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906116346004830182612c7c565b611642600583016000612c7c565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b50600101611213565b6116966118d3565b61169f81612094565b50565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff9081169116146117375760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610988565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa1580156117e5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061180991906139f2565b15611840576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61184d8160200151612189565b600061185c82602001516105c5565b9050805160001480611880575080805190602001208260a001518051906020012014155b156118bd578160a001516040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016109889190612da8565b6118cf826020015183606001516122af565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611954576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610988565b565b7f00000000000000000000000000000000000000000000000000000000000000006119ad576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611a435760008382815181106119cd576119cd61371e565b602002602001015190506119eb8160026122f690919063ffffffff16565b15611a3a5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506001016119b0565b5060005b8151811015611203576000828281518110611a6457611a6461371e565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611aa85750611b04565b611ab3600282612318565b15611b025760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611a47565b600081815260018301602052604081205415155b9392505050565b60808101517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614611bbc5760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610988565b60208101516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815260809190911b77ffffffffffffffff000000000000000000000000000000001660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632cbc26bb90602401602060405180830381865afa158015611c6a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c8e91906139f2565b15611cc5576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611cd2816040015161233a565b611cdf81602001516123b9565b61169f81602001518260600151612507565b6040517f9dc29fac000000000000000000000000000000000000000000000000000000008152306004820152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639dc29fac90604401600060405180830381600087803b158015611d7f57600080fd5b505af1158015611d93573d6000803e3d6000fd5b5050505050565b60606000611b208361254b565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611e3582606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611e19919061349b565b85608001516fffffffffffffffffffffffffffffffff166125a6565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b611e6283610d4b565b611ea4576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610988565b611eaf826000611f43565b67ffffffffffffffff83166000908152600760205260409020611ed290836125d0565b611edd816000611f43565b67ffffffffffffffff83166000908152600760205260409020611f0390600201826125d0565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b838383604051611f3693929190613a0f565b60405180910390a1505050565b81511561200a5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580611f99575060408201516fffffffffffffffffffffffffffffffff16155b15611fd257816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016109889190613a92565b80156118cf576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408201516fffffffffffffffffffffffffffffffff16151580612043575060208201516fffffffffffffffffffffffffffffffff1615155b156118cf57816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016109889190613a92565b6000611b208383612772565b6000611b2083836127c1565b3373ffffffffffffffffffffffffffffffffffffffff821603612113576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610988565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61219281610d4b565b6121d4576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610988565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015612253573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061227791906139f2565b61169f576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610988565b67ffffffffffffffff821660009081526007602052604090206118cf90600201827f00000000000000000000000000000000000000000000000000000000000000006128b4565b6000611b208373ffffffffffffffffffffffffffffffffffffffff84166127c1565b6000611b208373ffffffffffffffffffffffffffffffffffffffff8416612772565b7f00000000000000000000000000000000000000000000000000000000000000001561169f5761236b600282612c37565b61169f576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610988565b6123c281610d4b565b612404576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610988565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa15801561247d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124a19190613ace565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461169f576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610988565b67ffffffffffffffff821660009081526007602052604090206118cf90827f00000000000000000000000000000000000000000000000000000000000000006128b4565b60608160000180548060200260200160405190810160405280929190818152602001828054801561066957602002820191906000526020600020905b8154815260200190600101908083116125875750505050509050919050565b60006125c5856125b68486613aeb565b6125c09087613b02565b612c66565b90505b949350505050565b81546000906125f990700100000000000000000000000000000000900463ffffffff164261349b565b9050801561269b5760018301548354612641916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166125a6565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546126c1916fffffffffffffffffffffffffffffffff9081169116612c66565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990611f36908490613a92565b60008181526001830160205260408120546127b9575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105bf565b5060006105bf565b600081815260018301602052604081205480156128aa5760006127e560018361349b565b85549091506000906127f99060019061349b565b905080821461285e5760008660000182815481106128195761281961371e565b906000526020600020015490508087600001848154811061283c5761283c61371e565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061286f5761286f613b15565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105bf565b60009150506105bf565b825474010000000000000000000000000000000000000000900460ff1615806128db575081155b156128e557505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061292b90700100000000000000000000000000000000900463ffffffff164261349b565b905080156129eb578183111561296d576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546129a79083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166125a6565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612aa25773ffffffffffffffffffffffffffffffffffffffff8416612a4a576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610988565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610988565b84831015612bb55760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612ae6908261349b565b612af0878a61349b565b612afa9190613b02565b612b049190613b44565b905073ffffffffffffffffffffffffffffffffffffffff8616612b5d576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610988565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610988565b612bbf858461349b565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611b20565b6000818310612c755781611b20565b5090919050565b508054612c8890613260565b6000825580601f10612c98575050565b601f01602090049060005260206000209081019061169f91905b80821115612cc65760008155600101612cb2565b5090565b600060208284031215612cdc57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611b2057600080fd5b803567ffffffffffffffff81168114612d2457600080fd5b919050565b600060208284031215612d3b57600080fd5b611b2082612d0c565b6000815180845260005b81811015612d6a57602081850181015186830182015201612d4e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611b206020830184612d44565b73ffffffffffffffffffffffffffffffffffffffff8116811461169f57600080fd5b8035612d2481612dbb565b600060208284031215612dfa57600080fd5b8135611b2081612dbb565b600060208284031215612e1757600080fd5b813567ffffffffffffffff811115612e2e57600080fd5b82016101008185031215611b2057600080fd5b60008083601f840112612e5357600080fd5b50813567ffffffffffffffff811115612e6b57600080fd5b6020830191508360208260051b8501011115612e8657600080fd5b9250929050565b60008060008060408587031215612ea357600080fd5b843567ffffffffffffffff80821115612ebb57600080fd5b612ec788838901612e41565b90965094506020870135915080821115612ee057600080fd5b50612eed87828801612e41565b95989497509550505050565b600080600060408486031215612f0e57600080fd5b612f1784612d0c565b9250602084013567ffffffffffffffff80821115612f3457600080fd5b818601915086601f830112612f4857600080fd5b813581811115612f5757600080fd5b876020828501011115612f6957600080fd5b6020830194508093505050509250925092565b600060208284031215612f8e57600080fd5b813567ffffffffffffffff811115612fa557600080fd5b820160a08185031215611b2057600080fd5b602081526000825160406020840152612fd36060840182612d44565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084830301604085015261300e8282612d44565b95945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561306557835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613033565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561306557835167ffffffffffffffff168352928401929184019160010161308d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff81118282101715613106576131066130b3565b60405290565b60405160c0810167ffffffffffffffff81118282101715613106576131066130b3565b801515811461169f57600080fd5b8035612d248161312f565b80356fffffffffffffffffffffffffffffffff81168114612d2457600080fd5b60006060828403121561317a57600080fd5b6040516060810181811067ffffffffffffffff8211171561319d5761319d6130b3565b60405290508082356131ae8161312f565b81526131bc60208401613148565b60208201526131cd60408401613148565b60408201525092915050565b600080600060e084860312156131ee57600080fd5b6131f784612d0c565b92506132068560208601613168565b91506132158560808601613168565b90509250925092565b6000806020838503121561323157600080fd5b823567ffffffffffffffff81111561324857600080fd5b61325485828601612e41565b90969095509350505050565b600181811c9082168061327457607f821691505b6020821081036132ad577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082601f8301126132c457600080fd5b813567ffffffffffffffff808211156132df576132df6130b3565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715613325576133256130b3565b8160405283815286602085880101111561333e57600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000610100823603121561337157600080fd5b6133796130e2565b823567ffffffffffffffff8082111561339157600080fd5b61339d368387016132b3565b83526133ab60208601612d0c565b60208401526133bc60408601612ddd565b6040840152606085013560608401526133d760808601612ddd565b608084015260a08501359150808211156133f057600080fd5b6133fc368387016132b3565b60a084015260c085013591508082111561341557600080fd5b613421368387016132b3565b60c084015260e085013591508082111561343a57600080fd5b50613447368286016132b3565b60e08301525092915050565b60006020828403121561346557600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156105bf576105bf61346c565b601f821115611203576000816000526020600020601f850160051c810160208610156134d75750805b601f850160051c820191505b818110156134f6578281556001016134e3565b505050505050565b67ffffffffffffffff831115613516576135166130b3565b61352a836135248354613260565b836134ae565b6000601f84116001811461357c57600085156135465750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355611d93565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156135cb57868501358255602094850194600190920191016135ab565b5086821015613606577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b60408152600061362b6040830186612d44565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060a0823603121561368e57600080fd5b60405160a0810167ffffffffffffffff82821081831117156136b2576136b26130b3565b8160405284359150808211156136c757600080fd5b506136d4368286016132b3565b8252506136e360208401612d0c565b602082015260408301356136f681612dbb565b604082015260608381013590820152608083013561371381612dbb565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec183360301811261378157600080fd5b9190910192915050565b6000610140823603121561379e57600080fd5b6137a661310c565b6137af83612d0c565b81526137bd6020840161313d565b6020820152604083013567ffffffffffffffff808211156137dd57600080fd5b6137e9368387016132b3565b6040840152606085013591508082111561380257600080fd5b5061380f368286016132b3565b6060830152506138223660808501613168565b60808201526138343660e08501613168565b60a082015292915050565b815167ffffffffffffffff811115613859576138596130b3565b61386d816138678454613260565b846134ae565b602080601f8311600181146138c0576000841561388a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556134f6565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561390d578886015182559484019460019091019084016138ee565b508582101561394957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff8716835280602084015261397d81840187612d44565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506139bb9050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e083015261300e565b600060208284031215613a0457600080fd5b8151611b208161312f565b67ffffffffffffffff8416815260e08101613a5b60208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c08301526125c8565b606081016105bf82848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060208284031215613ae057600080fd5b8151611b2081612dbb565b80820281158282048414176105bf576105bf61346c565b808201808211156105bf576105bf61346c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600082613b7a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fe4275726e5769746846726f6d4d696e745265626173696e67546f6b656e506f6f6c20312e352e30a164736f6c6343000818000a",
}

var BurnWithFromMintRebasingTokenPoolABI = BurnWithFromMintRebasingTokenPoolMetaData.ABI

var BurnWithFromMintRebasingTokenPoolBin = BurnWithFromMintRebasingTokenPoolMetaData.Bin

func DeployBurnWithFromMintRebasingTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnWithFromMintRebasingTokenPool, error) {
	parsed, err := BurnWithFromMintRebasingTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnWithFromMintRebasingTokenPoolBin), backend, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnWithFromMintRebasingTokenPool{address: address, abi: *parsed, BurnWithFromMintRebasingTokenPoolCaller: BurnWithFromMintRebasingTokenPoolCaller{contract: contract}, BurnWithFromMintRebasingTokenPoolTransactor: BurnWithFromMintRebasingTokenPoolTransactor{contract: contract}, BurnWithFromMintRebasingTokenPoolFilterer: BurnWithFromMintRebasingTokenPoolFilterer{contract: contract}}, nil
}

type BurnWithFromMintRebasingTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnWithFromMintRebasingTokenPoolCaller
	BurnWithFromMintRebasingTokenPoolTransactor
	BurnWithFromMintRebasingTokenPoolFilterer
}

type BurnWithFromMintRebasingTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnWithFromMintRebasingTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnWithFromMintRebasingTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnWithFromMintRebasingTokenPoolSession struct {
	Contract     *BurnWithFromMintRebasingTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintRebasingTokenPoolCallerSession struct {
	Contract *BurnWithFromMintRebasingTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnWithFromMintRebasingTokenPoolTransactorSession struct {
	Contract     *BurnWithFromMintRebasingTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintRebasingTokenPoolRaw struct {
	Contract *BurnWithFromMintRebasingTokenPool
}

type BurnWithFromMintRebasingTokenPoolCallerRaw struct {
	Contract *BurnWithFromMintRebasingTokenPoolCaller
}

type BurnWithFromMintRebasingTokenPoolTransactorRaw struct {
	Contract *BurnWithFromMintRebasingTokenPoolTransactor
}

func NewBurnWithFromMintRebasingTokenPool(address common.Address, backend bind.ContractBackend) (*BurnWithFromMintRebasingTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnWithFromMintRebasingTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnWithFromMintRebasingTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPool{address: address, abi: abi, BurnWithFromMintRebasingTokenPoolCaller: BurnWithFromMintRebasingTokenPoolCaller{contract: contract}, BurnWithFromMintRebasingTokenPoolTransactor: BurnWithFromMintRebasingTokenPoolTransactor{contract: contract}, BurnWithFromMintRebasingTokenPoolFilterer: BurnWithFromMintRebasingTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnWithFromMintRebasingTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnWithFromMintRebasingTokenPoolCaller, error) {
	contract, err := bindBurnWithFromMintRebasingTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolCaller{contract: contract}, nil
}

func NewBurnWithFromMintRebasingTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnWithFromMintRebasingTokenPoolTransactor, error) {
	contract, err := bindBurnWithFromMintRebasingTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolTransactor{contract: contract}, nil
}

func NewBurnWithFromMintRebasingTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnWithFromMintRebasingTokenPoolFilterer, error) {
	contract, err := bindBurnWithFromMintRebasingTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolFilterer{contract: contract}, nil
}

func bindBurnWithFromMintRebasingTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnWithFromMintRebasingTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintRebasingTokenPool.Contract.BurnWithFromMintRebasingTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.BurnWithFromMintRebasingTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.BurnWithFromMintRebasingTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintRebasingTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetAllowList(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetAllowList(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRemotePool(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRemotePool(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRouter(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetRouter(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetToken(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.GetToken(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintRebasingTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintRebasingTokenPool.CallOpts, token)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintRebasingTokenPool.CallOpts, token)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) Owner() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.Owner(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.Owner(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SupportsInterface(&_BurnWithFromMintRebasingTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SupportsInterface(&_BurnWithFromMintRebasingTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnWithFromMintRebasingTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintRebasingTokenPool.CallOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintRebasingTokenPool.TransactOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintRebasingTokenPool.TransactOpts)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintRebasingTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintRebasingTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintRebasingTokenPool.TransactOpts, chains)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintRebasingTokenPool.TransactOpts, chains)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.LockOrBurn(&_BurnWithFromMintRebasingTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.LockOrBurn(&_BurnWithFromMintRebasingTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintRebasingTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintRebasingTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintRebasingTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintRebasingTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintRebasingTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintRebasingTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRemotePool(&_BurnWithFromMintRebasingTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRemotePool(&_BurnWithFromMintRebasingTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRouter(&_BurnWithFromMintRebasingTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.SetRouter(&_BurnWithFromMintRebasingTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.TransferOwnership(&_BurnWithFromMintRebasingTokenPool.TransactOpts, to)
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintRebasingTokenPool.Contract.TransferOwnership(&_BurnWithFromMintRebasingTokenPool.TransactOpts, to)
}

type BurnWithFromMintRebasingTokenPoolAllowListAddIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolAllowListAdd)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolAllowListAdd)
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

func (it *BurnWithFromMintRebasingTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolAllowListAddIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolAllowListAdd)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnWithFromMintRebasingTokenPoolAllowListAdd, error) {
	event := new(BurnWithFromMintRebasingTokenPoolAllowListAdd)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolAllowListRemove)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolAllowListRemove)
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

func (it *BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolAllowListRemove)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnWithFromMintRebasingTokenPoolAllowListRemove, error) {
	event := new(BurnWithFromMintRebasingTokenPoolAllowListRemove)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolBurnedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolBurned)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolBurned)
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

func (it *BurnWithFromMintRebasingTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintRebasingTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolBurnedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolBurned)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseBurned(log types.Log) (*BurnWithFromMintRebasingTokenPoolBurned, error) {
	event := new(BurnWithFromMintRebasingTokenPoolBurned)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolChainAddedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolChainAdded)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolChainAdded)
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

func (it *BurnWithFromMintRebasingTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolChainAddedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolChainAdded)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainAdded, error) {
	event := new(BurnWithFromMintRebasingTokenPoolChainAdded)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolChainConfiguredIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolChainConfigured)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolChainConfigured)
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

func (it *BurnWithFromMintRebasingTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolChainConfiguredIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolChainConfigured)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainConfigured, error) {
	event := new(BurnWithFromMintRebasingTokenPoolChainConfigured)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolChainRemovedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolChainRemoved)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolChainRemoved)
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

func (it *BurnWithFromMintRebasingTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolChainRemovedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolChainRemoved)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainRemoved, error) {
	event := new(BurnWithFromMintRebasingTokenPoolChainRemoved)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolConfigChangedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolConfigChanged)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolConfigChanged)
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

func (it *BurnWithFromMintRebasingTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolConfigChangedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolConfigChanged)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnWithFromMintRebasingTokenPoolConfigChanged, error) {
	event := new(BurnWithFromMintRebasingTokenPoolConfigChanged)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolLockedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolLocked)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolLocked)
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

func (it *BurnWithFromMintRebasingTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintRebasingTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolLockedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolLocked)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseLocked(log types.Log) (*BurnWithFromMintRebasingTokenPoolLocked, error) {
	event := new(BurnWithFromMintRebasingTokenPoolLocked)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolMintedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolMinted)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolMinted)
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

func (it *BurnWithFromMintRebasingTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintRebasingTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolMintedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolMinted)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseMinted(log types.Log) (*BurnWithFromMintRebasingTokenPoolMinted, error) {
	event := new(BurnWithFromMintRebasingTokenPoolMinted)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested)
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

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolOwnershipTransferred)
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

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolOwnershipTransferred)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferred, error) {
	event := new(BurnWithFromMintRebasingTokenPoolOwnershipTransferred)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolReleasedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolReleased)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolReleased)
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

func (it *BurnWithFromMintRebasingTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintRebasingTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolReleasedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolReleased)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseReleased(log types.Log) (*BurnWithFromMintRebasingTokenPoolReleased, error) {
	event := new(BurnWithFromMintRebasingTokenPoolReleased)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolRemotePoolSet)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolRemotePoolSet)
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

func (it *BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolRemotePoolSet)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseRemotePoolSet(log types.Log) (*BurnWithFromMintRebasingTokenPoolRemotePoolSet, error) {
	event := new(BurnWithFromMintRebasingTokenPoolRemotePoolSet)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolRouterUpdated)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolRouterUpdated)
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

func (it *BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolRouterUpdated)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnWithFromMintRebasingTokenPoolRouterUpdated, error) {
	event := new(BurnWithFromMintRebasingTokenPoolRouterUpdated)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintRebasingTokenPoolTokensConsumedIterator struct {
	Event *BurnWithFromMintRebasingTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintRebasingTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintRebasingTokenPoolTokensConsumed)
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
		it.Event = new(BurnWithFromMintRebasingTokenPoolTokensConsumed)
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

func (it *BurnWithFromMintRebasingTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintRebasingTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintRebasingTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintRebasingTokenPoolTokensConsumedIterator{contract: _BurnWithFromMintRebasingTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintRebasingTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintRebasingTokenPoolTokensConsumed)
				if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnWithFromMintRebasingTokenPoolTokensConsumed, error) {
	event := new(BurnWithFromMintRebasingTokenPoolTokensConsumed)
	if err := _BurnWithFromMintRebasingTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnWithFromMintRebasingTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseAllowListAdd(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseAllowListRemove(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["Burned"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseBurned(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseChainAdded(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseChainConfigured(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseChainRemoved(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseConfigChanged(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["Locked"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseLocked(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["Minted"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseMinted(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseOwnershipTransferred(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["Released"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseReleased(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["RemotePoolSet"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseRemotePoolSet(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseRouterUpdated(log)
	case _BurnWithFromMintRebasingTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnWithFromMintRebasingTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnWithFromMintRebasingTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnWithFromMintRebasingTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnWithFromMintRebasingTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnWithFromMintRebasingTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnWithFromMintRebasingTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnWithFromMintRebasingTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnWithFromMintRebasingTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnWithFromMintRebasingTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnWithFromMintRebasingTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnWithFromMintRebasingTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnWithFromMintRebasingTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnWithFromMintRebasingTokenPoolRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (BurnWithFromMintRebasingTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnWithFromMintRebasingTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnWithFromMintRebasingTokenPool *BurnWithFromMintRebasingTokenPool) Address() common.Address {
	return _BurnWithFromMintRebasingTokenPool.address
}

type BurnWithFromMintRebasingTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

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

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnWithFromMintRebasingTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnWithFromMintRebasingTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintRebasingTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnWithFromMintRebasingTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnWithFromMintRebasingTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnWithFromMintRebasingTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintRebasingTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnWithFromMintRebasingTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintRebasingTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnWithFromMintRebasingTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintRebasingTokenPoolOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintRebasingTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnWithFromMintRebasingTokenPoolReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintRebasingTokenPoolRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*BurnWithFromMintRebasingTokenPoolRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnWithFromMintRebasingTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintRebasingTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintRebasingTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnWithFromMintRebasingTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
