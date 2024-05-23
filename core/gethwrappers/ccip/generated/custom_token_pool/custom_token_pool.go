// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package custom_token_pool

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

var CustomTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRatelimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousPoolAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SynthBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SynthMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chains\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePool\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destPoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"setRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162003ad638038062003ad6833981016040819052620000349162000534565b6040805160008082526020820190925284915083833380600081620000a05760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000d357620000d38162000186565b5050506001600160a01b0384161580620000f457506001600160a01b038116155b806200010757506001600160a01b038216155b1562000126576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0384811660805282811660a052600480546001600160a01b031916918316919091179055825115801560c052620001795760408051600081526020810190915262000179908462000231565b50505050505050620005d6565b336001600160a01b03821603620001e05760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000097565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c05162000252576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002dd57600083828151811062000276576200027662000588565b60209081029190910101519050620002906002826200038e565b15620002d3576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000255565b5060005b81518110156200038957600082828151811062000302576200030262000588565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200032e575062000380565b6200033b600282620003ae565b156200037e576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620002e1565b505050565b6000620003a5836001600160a01b038416620003c5565b90505b92915050565b6000620003a5836001600160a01b038416620004c9565b60008181526001830160205260408120548015620004be576000620003ec6001836200059e565b855490915060009062000402906001906200059e565b90508181146200046e57600086600001828154811062000426576200042662000588565b90600052602060002001549050808760000184815481106200044c576200044c62000588565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004825762000482620005c0565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003a8565b6000915050620003a8565b60008181526001830160205260408120546200051257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003a8565b506000620003a8565b6001600160a01b03811681146200053157600080fd5b50565b6000806000606084860312156200054a57600080fd5b835162000557816200051b565b60208501519093506200056a816200051b565b60408501519092506200057d816200051b565b809150509250925092565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003a857634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05161349162000645600039600081816104400152818161141a0152611d5701526000818161041a0152818161126201526116ce0152600081816101d1015281816102260152818161069d015281816115eb01528181611ced0152611f4701526134916000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c8063a7cd63b7116100d8578063c75eea9c1161008c578063e0351e1311610066578063e0351e131461043e578063f2fde38b14610464578063f6e2145e1461047757600080fd5b8063c75eea9c146103f2578063cf7401f314610405578063dc0bd9711461041857600080fd5b8063b0f479a1116100bd578063b0f479a1146103ac578063c0d78655146103ca578063c4bffe2b146103dd57600080fd5b8063a7cd63b714610328578063af58d59f1461033d57600080fd5b806354c8a4f31161013a5780638926f54f116101145780638926f54f146102d75780638da5cb5b146102ea5780639a4575b91461030857600080fd5b806354c8a4f3146102a757806378a010b2146102bc57806379ba5097146102cf57600080fd5b806321df0da71161016b57806321df0da7146101cf578063240028e8146102165780634059f55b1461026357600080fd5b806301ffc9a7146101875780630a2fd493146101af575b600080fd5b61019a6101953660046126ea565b61048a565b60405190151581526020015b60405180910390f35b6101c26101bd366004612749565b61056f565b6040516101a691906127c8565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101a6565b61019a610224366004612808565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b610276610271366004612825565b61061f565b60408051825173ffffffffffffffffffffffffffffffffffffffff16815260209283015192810192909252016101a6565b6102ba6102b53660046128ac565b6106cd565b005b6102ba6102ca366004612918565b610748565b6102ba6108bc565b61019a6102e5366004612749565b6109b9565b60005473ffffffffffffffffffffffffffffffffffffffff166101f1565b61031b61031636600461299b565b6109d0565b6040516101a691906129d6565b610330610a68565b6040516101a69190612a36565b61035061034b366004612749565b610a79565b6040516101a6919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff166101f1565b6102ba6103d8366004612808565b610b4e565b6103e5610c29565b6040516101a69190612a90565b610350610400366004612749565b610ce1565b6102ba610413366004612bec565b610db3565b7f00000000000000000000000000000000000000000000000000000000000000006101f1565b7f000000000000000000000000000000000000000000000000000000000000000061019a565b6102ba610472366004612808565b610dcb565b6102ba610485366004612c31565b610ddf565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061051d57507fffffffff0000000000000000000000000000000000000000000000000000000082167f773a5d4500000000000000000000000000000000000000000000000000000000145b8061056957507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061059a90612c73565b80601f01602080910402602001604051908101604052809291908181526020018280546105c690612c73565b80156106135780601f106105e857610100808354040283529160200191610613565b820191906000526020600020905b8154815290600101906020018083116105f657829003601f168201915b50505050509050919050565b604080518082019091526000808252602082015261064461063f83612d71565b611224565b604051606083013581527fbb0b72e5f44e331506684da008a30e10d50658c29d8159f6c6ab40bf1e52e6009060200160405180910390a1506040805180820190915273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152606090910135602082015290565b6106d5611395565b6107428484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061141892505050565b50505050565b610750611395565b610759836109b9565b6107a0576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b67ffffffffffffffff8316600090815260076020526040812060040180546107c790612c73565b80601f01602080910402602001604051908101604052809291908181526020018280546107f390612c73565b80156108405780601f1061081557610100808354040283529160200191610840565b820191906000526020600020905b81548152906001019060200180831161082357829003601f168201915b5050505067ffffffffffffffff861660009081526007602052604090209192505060040161086f838583612ea4565b508367ffffffffffffffff167fdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf8285856040516108ae93929190612fbf565b60405180910390a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461093d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610797565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610569600567ffffffffffffffff84166115ce565b60408051808201909152606080825260208201526109f56109f083613023565b6115e9565b604051606083013581527f02992093bca69a36949677658a77d359b510dc6232c68f9f118f7c0127a1b1479060200160405180910390a16040518060400160405280610a4d8460200160208101906101bd9190612749565b81526040805160208181019092526000815291015292915050565b6060610a7460026117b1565b905090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526003909101548084166060830152919091049091166080820152610569906117be565b610b56611395565b73ffffffffffffffffffffffffffffffffffffffff8116610ba3576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b60606000610c3760056117b1565b90506000815167ffffffffffffffff811115610c5557610c55612ad2565b604051908082528060200260200182016040528015610c7e578160200160208202803683370190505b50905060005b8251811015610cda57828181518110610c9f57610c9f6130aa565b6020026020010151828281518110610cb957610cb96130aa565b67ffffffffffffffff90921660209283029190910190910152600101610c84565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff161515948201949094526001909101548084166060830152919091049091166080820152610569906117be565b610dbb611395565b610dc6838383611870565b505050565b610dd3611395565b610ddc8161195a565b50565b610de7611395565b60005b81811015610dc6576000838383818110610e0657610e066130aa565b9050602002810190610e1891906130d9565b610e2190613117565b9050610e368160600151826020015115611a4f565b610e498160800151826020015115611a4f565b806020015115611112578051610e6b9060059067ffffffffffffffff16611b8c565b610eb05780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610797565b806040015151600003610eef576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252606083810180516020908101516fffffffffffffffffffffffffffffffff9081168486019081524263ffffffff9081166080808901829052865151151560a0808b01919091528751870151861660c08b015296518a0151851660e08a015292885288519586018952828a01805186015185168752868601919091528051511515868a015280518501518416868801525188015183168583015283870194855288880151878901908152895167ffffffffffffffff1660009081526007865289902088518051825482890151838e01519289167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316177001000000000000000000000000000000009188168202177fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff908116740100000000000000000000000000000000000000009415158502178655848d015194890151948a16948a168202949094176001860155995180516002860180549b8301519f830151918b169b9093169a909a179d9096168a029c909c179091169615150295909517909855948501519401519381169316909102919091176003820155915190919060048201906110c5908261319b565b50508151606083015160808401516040517f0f135cbb9afa12a8bf3bbd071c117bcca4ddeca6160ef7f33d012a81b9c0c4719450611105939291906132b5565b60405180910390a161121b565b805161112a9060059067ffffffffffffffff16611b98565b61116f5780516040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610797565b805167ffffffffffffffff16600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906111d8600483018261269c565b5050805160405167ffffffffffffffff90911681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599169060200160405180910390a15b50600101610dea565b60208101516040517f58babe3300000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906358babe3390602401602060405180830381865afa1580156112be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112e29190613338565b15611319576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113268160200151611ba4565b611333816020015161056f565b80519060200120816080015180519060200120146113835780608001516040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161079791906127c8565b610ddc81602001518260600151611cca565b60005473ffffffffffffffffffffffffffffffffffffffff163314611416576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610797565b565b7f000000000000000000000000000000000000000000000000000000000000000061146f576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b825181101561150557600083828151811061148f5761148f6130aa565b602002602001015190506114ad816002611d1190919063ffffffff16565b156114fc5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611472565b5060005b8151811015610dc6576000828281518110611526576115266130aa565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361156a57506115c6565b611575600282611d33565b156115c45760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611509565b600081815260018301602052604081205415155b9392505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16816080015173ffffffffffffffffffffffffffffffffffffffff16146116905760808101516040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610797565b60208101516040517f58babe3300000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906358babe3390602401602060405180830381865afa15801561172a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061174e9190613338565b15611785576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117928160400151611d55565b61179f8160200151611dd9565b610ddc81602001518260600151611f27565b606060006115e283611f6b565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261184c82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426118309190613384565b85608001516fffffffffffffffffffffffffffffffff16611fc6565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b611879836109b9565b6118bb576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610797565b6118c6826000611a4f565b67ffffffffffffffff831660009081526007602052604090206118e99083611ff0565b6118f4816000611a4f565b67ffffffffffffffff8316600090815260076020526040902061191a9060020182611ff0565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b83838360405161194d939291906132b5565b60405180910390a1505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036119d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610797565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b815115611b1a5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580611aa5575060408201516fffffffffffffffffffffffffffffffff16155b15611ade57816040517f70505e560000000000000000000000000000000000000000000000000000000081526004016107979190613397565b8015611b16576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580611b53575060208201516fffffffffffffffffffffffffffffffff1615155b15611b1657816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016107979190613397565b60006115e28383612192565b60006115e283836121e1565b611bad816109b9565b611bef576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610797565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015611c6e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c929190613338565b610ddc576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b67ffffffffffffffff82166000908152600760205260409020611b1690600201827f00000000000000000000000000000000000000000000000000000000000000006122d4565b60006115e28373ffffffffffffffffffffffffffffffffffffffff84166121e1565b60006115e28373ffffffffffffffffffffffffffffffffffffffff8416612192565b7f00000000000000000000000000000000000000000000000000000000000000008015611d8a5750611d88600282612657565b155b15610ddc576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610797565b611de2816109b9565b611e24576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610797565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015611e9d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ec191906133d3565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610ddc576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610797565b67ffffffffffffffff82166000908152600760205260409020611b1690827f00000000000000000000000000000000000000000000000000000000000000006122d4565b60608160000180548060200260200160405190810160405280929190818152602001828054801561061357602002820191906000526020600020905b815481526020019060010190808311611fa75750505050509050919050565b6000611fe585611fd684866133f0565b611fe09087613407565b612686565b90505b949350505050565b815460009061201990700100000000000000000000000000000000900463ffffffff1642613384565b905080156120bb5760018301548354612061916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416611fc6565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546120e1916fffffffffffffffffffffffffffffffff9081169116612686565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199061194d908490613397565b60008181526001830160205260408120546121d957508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610569565b506000610569565b600081815260018301602052604081205480156122ca576000612205600183613384565b855490915060009061221990600190613384565b905081811461227e576000866000018281548110612239576122396130aa565b906000526020600020015490508087600001848154811061225c5761225c6130aa565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061228f5761228f61341a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610569565b6000915050610569565b825474010000000000000000000000000000000000000000900460ff1615806122fb575081155b1561230557505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061234b90700100000000000000000000000000000000900463ffffffff1642613384565b9050801561240b578183111561238d576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546123c79083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16611fc6565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156124c25773ffffffffffffffffffffffffffffffffffffffff841661246a576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610797565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610797565b848310156125d55760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff169060009082906125069082613384565b612510878a613384565b61251a9190613407565b6125249190613449565b905073ffffffffffffffffffffffffffffffffffffffff861661257d576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610797565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610797565b6125df8584613384565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156115e2565b600081831061269557816115e2565b5090919050565b5080546126a890612c73565b6000825580601f106126b8575050565b601f016020900490600052602060002090810190610ddc91905b808211156126e657600081556001016126d2565b5090565b6000602082840312156126fc57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146115e257600080fd5b803567ffffffffffffffff8116811461274457600080fd5b919050565b60006020828403121561275b57600080fd5b6115e28261272c565b6000815180845260005b8181101561278a5760208185018101518683018201520161276e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006115e26020830184612764565b73ffffffffffffffffffffffffffffffffffffffff81168114610ddc57600080fd5b8035612744816127db565b60006020828403121561281a57600080fd5b81356115e2816127db565b60006020828403121561283757600080fd5b813567ffffffffffffffff81111561284e57600080fd5b820160e081850312156115e257600080fd5b60008083601f84011261287257600080fd5b50813567ffffffffffffffff81111561288a57600080fd5b6020830191508360208260051b85010111156128a557600080fd5b9250929050565b600080600080604085870312156128c257600080fd5b843567ffffffffffffffff808211156128da57600080fd5b6128e688838901612860565b909650945060208701359150808211156128ff57600080fd5b5061290c87828801612860565b95989497509550505050565b60008060006040848603121561292d57600080fd5b6129368461272c565b9250602084013567ffffffffffffffff8082111561295357600080fd5b818601915086601f83011261296757600080fd5b81358181111561297657600080fd5b87602082850101111561298857600080fd5b6020830194508093505050509250925092565b6000602082840312156129ad57600080fd5b813567ffffffffffffffff8111156129c457600080fd5b820160a081850312156115e257600080fd5b6020815260008251604060208401526129f26060840182612764565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612a2d8282612764565b95945050505050565b6020808252825182820181905260009190848201906040850190845b81811015612a8457835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612a52565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015612a8457835167ffffffffffffffff1683529284019291840191600101612aac565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff81118282101715612b2457612b24612ad2565b60405290565b60405160a0810167ffffffffffffffff81118282101715612b2457612b24612ad2565b8015158114610ddc57600080fd5b80356fffffffffffffffffffffffffffffffff8116811461274457600080fd5b600060608284031215612b8d57600080fd5b6040516060810181811067ffffffffffffffff82111715612bb057612bb0612ad2565b6040529050808235612bc181612b4d565b8152612bcf60208401612b5b565b6020820152612be060408401612b5b565b60408201525092915050565b600080600060e08486031215612c0157600080fd5b612c0a8461272c565b9250612c198560208601612b7b565b9150612c288560808601612b7b565b90509250925092565b60008060208385031215612c4457600080fd5b823567ffffffffffffffff811115612c5b57600080fd5b612c6785828601612860565b90969095509350505050565b600181811c90821680612c8757607f821691505b602082108103612cc0577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600082601f830112612cd757600080fd5b813567ffffffffffffffff80821115612cf257612cf2612ad2565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715612d3857612d38612ad2565b81604052838152866020858801011115612d5157600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060e08236031215612d8357600080fd5b612d8b612b01565b823567ffffffffffffffff80821115612da357600080fd5b612daf36838701612cc6565b8352612dbd6020860161272c565b6020840152612dce604086016127fd565b6040840152606085013560608401526080850135915080821115612df157600080fd5b612dfd36838701612cc6565b608084015260a0850135915080821115612e1657600080fd5b612e2236838701612cc6565b60a084015260c0850135915080821115612e3b57600080fd5b50612e4836828601612cc6565b60c08301525092915050565b601f821115610dc6576000816000526020600020601f850160051c81016020861015612e7d5750805b601f850160051c820191505b81811015612e9c57828155600101612e89565b505050505050565b67ffffffffffffffff831115612ebc57612ebc612ad2565b612ed083612eca8354612c73565b83612e54565b6000601f841160018114612f225760008515612eec5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355612fb8565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015612f715786850135825560209485019460019092019101612f51565b5086821015612fac577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b604081526000612fd26040830186612764565b82810360208401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050949350505050565b600060a0823603121561303557600080fd5b61303d612b2a565b823567ffffffffffffffff81111561305457600080fd5b61306036828601612cc6565b82525061306f6020840161272c565b60208201526040830135613082816127db565b604082015260608381013590820152608083013561309f816127db565b608082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261310d57600080fd5b9190910192915050565b6000610120823603121561312a57600080fd5b613132612b2a565b61313b8361272c565b8152602083013561314b81612b4d565b6020820152604083013567ffffffffffffffff81111561316a57600080fd5b61317636828601612cc6565b6040830152506131893660608501612b7b565b606082015261309f3660c08501612b7b565b815167ffffffffffffffff8111156131b5576131b5612ad2565b6131c9816131c38454612c73565b84612e54565b602080601f83116001811461321c57600084156131e65750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555612e9c565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156132695788860151825594840194600190910190840161324a565b50858210156132a557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b67ffffffffffffffff8416815260e0810161330160208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152611fe8565b60006020828403121561334a57600080fd5b81516115e281612b4d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561056957610569613355565b6060810161056982848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b6000602082840312156133e557600080fd5b81516115e2816127db565b808202811582820484141761056957610569613355565b8082018082111561056957610569613355565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b60008261347f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea164736f6c6343000818000a",
}

var CustomTokenPoolABI = CustomTokenPoolMetaData.ABI

var CustomTokenPoolBin = CustomTokenPoolMetaData.Bin

func DeployCustomTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *CustomTokenPool, error) {
	parsed, err := CustomTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CustomTokenPoolBin), backend, token, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CustomTokenPool{address: address, abi: *parsed, CustomTokenPoolCaller: CustomTokenPoolCaller{contract: contract}, CustomTokenPoolTransactor: CustomTokenPoolTransactor{contract: contract}, CustomTokenPoolFilterer: CustomTokenPoolFilterer{contract: contract}}, nil
}

type CustomTokenPool struct {
	address common.Address
	abi     abi.ABI
	CustomTokenPoolCaller
	CustomTokenPoolTransactor
	CustomTokenPoolFilterer
}

type CustomTokenPoolCaller struct {
	contract *bind.BoundContract
}

type CustomTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type CustomTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type CustomTokenPoolSession struct {
	Contract     *CustomTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CustomTokenPoolCallerSession struct {
	Contract *CustomTokenPoolCaller
	CallOpts bind.CallOpts
}

type CustomTokenPoolTransactorSession struct {
	Contract     *CustomTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type CustomTokenPoolRaw struct {
	Contract *CustomTokenPool
}

type CustomTokenPoolCallerRaw struct {
	Contract *CustomTokenPoolCaller
}

type CustomTokenPoolTransactorRaw struct {
	Contract *CustomTokenPoolTransactor
}

func NewCustomTokenPool(address common.Address, backend bind.ContractBackend) (*CustomTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(CustomTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCustomTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPool{address: address, abi: abi, CustomTokenPoolCaller: CustomTokenPoolCaller{contract: contract}, CustomTokenPoolTransactor: CustomTokenPoolTransactor{contract: contract}, CustomTokenPoolFilterer: CustomTokenPoolFilterer{contract: contract}}, nil
}

func NewCustomTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*CustomTokenPoolCaller, error) {
	contract, err := bindCustomTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolCaller{contract: contract}, nil
}

func NewCustomTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*CustomTokenPoolTransactor, error) {
	contract, err := bindCustomTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolTransactor{contract: contract}, nil
}

func NewCustomTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*CustomTokenPoolFilterer, error) {
	contract, err := bindCustomTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolFilterer{contract: contract}, nil
}

func bindCustomTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CustomTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CustomTokenPool *CustomTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CustomTokenPool.Contract.CustomTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_CustomTokenPool *CustomTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.CustomTokenPoolTransactor.contract.Transfer(opts)
}

func (_CustomTokenPool *CustomTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.CustomTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_CustomTokenPool *CustomTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CustomTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_CustomTokenPool *CustomTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.contract.Transfer(opts)
}

func (_CustomTokenPool *CustomTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _CustomTokenPool.Contract.GetAllowList(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _CustomTokenPool.Contract.GetAllowList(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _CustomTokenPool.Contract.GetAllowListEnabled(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _CustomTokenPool.Contract.GetAllowListEnabled(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _CustomTokenPool.Contract.GetCurrentInboundRateLimiterState(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _CustomTokenPool.Contract.GetCurrentInboundRateLimiterState(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _CustomTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _CustomTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getRemotePool", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _CustomTokenPool.Contract.GetRemotePool(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetRemotePool(remoteChainSelector uint64) ([]byte, error) {
	return _CustomTokenPool.Contract.GetRemotePool(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _CustomTokenPool.Contract.GetRmnProxy(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _CustomTokenPool.Contract.GetRmnProxy(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetRouter() (common.Address, error) {
	return _CustomTokenPool.Contract.GetRouter(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _CustomTokenPool.Contract.GetRouter(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _CustomTokenPool.Contract.GetSupportedChains(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _CustomTokenPool.Contract.GetSupportedChains(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) GetToken() (common.Address, error) {
	return _CustomTokenPool.Contract.GetToken(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _CustomTokenPool.Contract.GetToken(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _CustomTokenPool.Contract.IsSupportedChain(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _CustomTokenPool.Contract.IsSupportedChain(&_CustomTokenPool.CallOpts, remoteChainSelector)
}

func (_CustomTokenPool *CustomTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _CustomTokenPool.Contract.IsSupportedToken(&_CustomTokenPool.CallOpts, token)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _CustomTokenPool.Contract.IsSupportedToken(&_CustomTokenPool.CallOpts, token)
}

func (_CustomTokenPool *CustomTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) Owner() (common.Address, error) {
	return _CustomTokenPool.Contract.Owner(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) Owner() (common.Address, error) {
	return _CustomTokenPool.Contract.Owner(&_CustomTokenPool.CallOpts)
}

func (_CustomTokenPool *CustomTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CustomTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CustomTokenPool *CustomTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CustomTokenPool.Contract.SupportsInterface(&_CustomTokenPool.CallOpts, interfaceId)
}

func (_CustomTokenPool *CustomTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CustomTokenPool.Contract.SupportsInterface(&_CustomTokenPool.CallOpts, interfaceId)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_CustomTokenPool *CustomTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _CustomTokenPool.Contract.AcceptOwnership(&_CustomTokenPool.TransactOpts)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CustomTokenPool.Contract.AcceptOwnership(&_CustomTokenPool.TransactOpts)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_CustomTokenPool *CustomTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ApplyAllowListUpdates(&_CustomTokenPool.TransactOpts, removes, adds)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ApplyAllowListUpdates(&_CustomTokenPool.TransactOpts, removes, adds)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_CustomTokenPool *CustomTokenPoolSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ApplyChainUpdates(&_CustomTokenPool.TransactOpts, chains)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) ApplyChainUpdates(chains []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ApplyChainUpdates(&_CustomTokenPool.TransactOpts, chains)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_CustomTokenPool *CustomTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.LockOrBurn(&_CustomTokenPool.TransactOpts, lockOrBurnIn)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.LockOrBurn(&_CustomTokenPool.TransactOpts, lockOrBurnIn)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_CustomTokenPool *CustomTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ReleaseOrMint(&_CustomTokenPool.TransactOpts, releaseOrMintIn)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.ReleaseOrMint(&_CustomTokenPool.TransactOpts, releaseOrMintIn)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_CustomTokenPool *CustomTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetChainRateLimiterConfig(&_CustomTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetChainRateLimiterConfig(&_CustomTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "setRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_CustomTokenPool *CustomTokenPoolSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetRemotePool(&_CustomTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) SetRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetRemotePool(&_CustomTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_CustomTokenPool *CustomTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetRouter(&_CustomTokenPool.TransactOpts, newRouter)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.SetRouter(&_CustomTokenPool.TransactOpts, newRouter)
}

func (_CustomTokenPool *CustomTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_CustomTokenPool *CustomTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.TransferOwnership(&_CustomTokenPool.TransactOpts, to)
}

func (_CustomTokenPool *CustomTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _CustomTokenPool.Contract.TransferOwnership(&_CustomTokenPool.TransactOpts, to)
}

type CustomTokenPoolAllowListAddIterator struct {
	Event *CustomTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolAllowListAdd)
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
		it.Event = new(CustomTokenPoolAllowListAdd)
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

func (it *CustomTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*CustomTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolAllowListAddIterator{contract: _CustomTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolAllowListAdd)
				if err := _CustomTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*CustomTokenPoolAllowListAdd, error) {
	event := new(CustomTokenPoolAllowListAdd)
	if err := _CustomTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolAllowListRemoveIterator struct {
	Event *CustomTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolAllowListRemove)
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
		it.Event = new(CustomTokenPoolAllowListRemove)
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

func (it *CustomTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*CustomTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolAllowListRemoveIterator{contract: _CustomTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolAllowListRemove)
				if err := _CustomTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*CustomTokenPoolAllowListRemove, error) {
	event := new(CustomTokenPoolAllowListRemove)
	if err := _CustomTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolBurnedIterator struct {
	Event *CustomTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolBurned)
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
		it.Event = new(CustomTokenPoolBurned)
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

func (it *CustomTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*CustomTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolBurnedIterator{contract: _CustomTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolBurned)
				if err := _CustomTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseBurned(log types.Log) (*CustomTokenPoolBurned, error) {
	event := new(CustomTokenPoolBurned)
	if err := _CustomTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolChainAddedIterator struct {
	Event *CustomTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolChainAdded)
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
		it.Event = new(CustomTokenPoolChainAdded)
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

func (it *CustomTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*CustomTokenPoolChainAddedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolChainAddedIterator{contract: _CustomTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolChainAdded)
				if err := _CustomTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseChainAdded(log types.Log) (*CustomTokenPoolChainAdded, error) {
	event := new(CustomTokenPoolChainAdded)
	if err := _CustomTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolChainConfiguredIterator struct {
	Event *CustomTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolChainConfigured)
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
		it.Event = new(CustomTokenPoolChainConfigured)
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

func (it *CustomTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*CustomTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolChainConfiguredIterator{contract: _CustomTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolChainConfigured)
				if err := _CustomTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseChainConfigured(log types.Log) (*CustomTokenPoolChainConfigured, error) {
	event := new(CustomTokenPoolChainConfigured)
	if err := _CustomTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolChainRemovedIterator struct {
	Event *CustomTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolChainRemoved)
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
		it.Event = new(CustomTokenPoolChainRemoved)
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

func (it *CustomTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*CustomTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolChainRemovedIterator{contract: _CustomTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolChainRemoved)
				if err := _CustomTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseChainRemoved(log types.Log) (*CustomTokenPoolChainRemoved, error) {
	event := new(CustomTokenPoolChainRemoved)
	if err := _CustomTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolConfigChangedIterator struct {
	Event *CustomTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolConfigChanged)
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
		it.Event = new(CustomTokenPoolConfigChanged)
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

func (it *CustomTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*CustomTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolConfigChangedIterator{contract: _CustomTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolConfigChanged)
				if err := _CustomTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseConfigChanged(log types.Log) (*CustomTokenPoolConfigChanged, error) {
	event := new(CustomTokenPoolConfigChanged)
	if err := _CustomTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolLockedIterator struct {
	Event *CustomTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolLocked)
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
		it.Event = new(CustomTokenPoolLocked)
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

func (it *CustomTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*CustomTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolLockedIterator{contract: _CustomTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolLocked)
				if err := _CustomTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseLocked(log types.Log) (*CustomTokenPoolLocked, error) {
	event := new(CustomTokenPoolLocked)
	if err := _CustomTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolMintedIterator struct {
	Event *CustomTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolMinted)
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
		it.Event = new(CustomTokenPoolMinted)
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

func (it *CustomTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*CustomTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolMintedIterator{contract: _CustomTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolMinted)
				if err := _CustomTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseMinted(log types.Log) (*CustomTokenPoolMinted, error) {
	event := new(CustomTokenPoolMinted)
	if err := _CustomTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolOwnershipTransferRequestedIterator struct {
	Event *CustomTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolOwnershipTransferRequested)
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
		it.Event = new(CustomTokenPoolOwnershipTransferRequested)
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

func (it *CustomTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CustomTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolOwnershipTransferRequestedIterator{contract: _CustomTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolOwnershipTransferRequested)
				if err := _CustomTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*CustomTokenPoolOwnershipTransferRequested, error) {
	event := new(CustomTokenPoolOwnershipTransferRequested)
	if err := _CustomTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolOwnershipTransferredIterator struct {
	Event *CustomTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolOwnershipTransferred)
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
		it.Event = new(CustomTokenPoolOwnershipTransferred)
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

func (it *CustomTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CustomTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolOwnershipTransferredIterator{contract: _CustomTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolOwnershipTransferred)
				if err := _CustomTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*CustomTokenPoolOwnershipTransferred, error) {
	event := new(CustomTokenPoolOwnershipTransferred)
	if err := _CustomTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolReleasedIterator struct {
	Event *CustomTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolReleased)
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
		it.Event = new(CustomTokenPoolReleased)
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

func (it *CustomTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*CustomTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolReleasedIterator{contract: _CustomTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolReleased)
				if err := _CustomTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseReleased(log types.Log) (*CustomTokenPoolReleased, error) {
	event := new(CustomTokenPoolReleased)
	if err := _CustomTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolRemotePoolSetIterator struct {
	Event *CustomTokenPoolRemotePoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolRemotePoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolRemotePoolSet)
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
		it.Event = new(CustomTokenPoolRemotePoolSet)
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

func (it *CustomTokenPoolRemotePoolSetIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolRemotePoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolRemotePoolSet struct {
	RemoteChainSelector uint64
	PreviousPoolAddress []byte
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*CustomTokenPoolRemotePoolSetIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolRemotePoolSetIterator{contract: _CustomTokenPool.contract, event: "RemotePoolSet", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "RemotePoolSet", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolRemotePoolSet)
				if err := _CustomTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseRemotePoolSet(log types.Log) (*CustomTokenPoolRemotePoolSet, error) {
	event := new(CustomTokenPoolRemotePoolSet)
	if err := _CustomTokenPool.contract.UnpackLog(event, "RemotePoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolRouterUpdatedIterator struct {
	Event *CustomTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolRouterUpdated)
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
		it.Event = new(CustomTokenPoolRouterUpdated)
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

func (it *CustomTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*CustomTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolRouterUpdatedIterator{contract: _CustomTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolRouterUpdated)
				if err := _CustomTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*CustomTokenPoolRouterUpdated, error) {
	event := new(CustomTokenPoolRouterUpdated)
	if err := _CustomTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolSynthBurnedIterator struct {
	Event *CustomTokenPoolSynthBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolSynthBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolSynthBurned)
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
		it.Event = new(CustomTokenPoolSynthBurned)
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

func (it *CustomTokenPoolSynthBurnedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolSynthBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolSynthBurned struct {
	Amount *big.Int
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterSynthBurned(opts *bind.FilterOpts) (*CustomTokenPoolSynthBurnedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "SynthBurned")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolSynthBurnedIterator{contract: _CustomTokenPool.contract, event: "SynthBurned", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchSynthBurned(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolSynthBurned) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "SynthBurned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolSynthBurned)
				if err := _CustomTokenPool.contract.UnpackLog(event, "SynthBurned", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseSynthBurned(log types.Log) (*CustomTokenPoolSynthBurned, error) {
	event := new(CustomTokenPoolSynthBurned)
	if err := _CustomTokenPool.contract.UnpackLog(event, "SynthBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolSynthMintedIterator struct {
	Event *CustomTokenPoolSynthMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolSynthMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolSynthMinted)
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
		it.Event = new(CustomTokenPoolSynthMinted)
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

func (it *CustomTokenPoolSynthMintedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolSynthMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolSynthMinted struct {
	Amount *big.Int
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterSynthMinted(opts *bind.FilterOpts) (*CustomTokenPoolSynthMintedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "SynthMinted")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolSynthMintedIterator{contract: _CustomTokenPool.contract, event: "SynthMinted", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchSynthMinted(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolSynthMinted) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "SynthMinted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolSynthMinted)
				if err := _CustomTokenPool.contract.UnpackLog(event, "SynthMinted", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseSynthMinted(log types.Log) (*CustomTokenPoolSynthMinted, error) {
	event := new(CustomTokenPoolSynthMinted)
	if err := _CustomTokenPool.contract.UnpackLog(event, "SynthMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CustomTokenPoolTokensConsumedIterator struct {
	Event *CustomTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CustomTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CustomTokenPoolTokensConsumed)
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
		it.Event = new(CustomTokenPoolTokensConsumed)
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

func (it *CustomTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *CustomTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CustomTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_CustomTokenPool *CustomTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*CustomTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _CustomTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &CustomTokenPoolTokensConsumedIterator{contract: _CustomTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_CustomTokenPool *CustomTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _CustomTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CustomTokenPoolTokensConsumed)
				if err := _CustomTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_CustomTokenPool *CustomTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*CustomTokenPoolTokensConsumed, error) {
	event := new(CustomTokenPoolTokensConsumed)
	if err := _CustomTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CustomTokenPool *CustomTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CustomTokenPool.abi.Events["AllowListAdd"].ID:
		return _CustomTokenPool.ParseAllowListAdd(log)
	case _CustomTokenPool.abi.Events["AllowListRemove"].ID:
		return _CustomTokenPool.ParseAllowListRemove(log)
	case _CustomTokenPool.abi.Events["Burned"].ID:
		return _CustomTokenPool.ParseBurned(log)
	case _CustomTokenPool.abi.Events["ChainAdded"].ID:
		return _CustomTokenPool.ParseChainAdded(log)
	case _CustomTokenPool.abi.Events["ChainConfigured"].ID:
		return _CustomTokenPool.ParseChainConfigured(log)
	case _CustomTokenPool.abi.Events["ChainRemoved"].ID:
		return _CustomTokenPool.ParseChainRemoved(log)
	case _CustomTokenPool.abi.Events["ConfigChanged"].ID:
		return _CustomTokenPool.ParseConfigChanged(log)
	case _CustomTokenPool.abi.Events["Locked"].ID:
		return _CustomTokenPool.ParseLocked(log)
	case _CustomTokenPool.abi.Events["Minted"].ID:
		return _CustomTokenPool.ParseMinted(log)
	case _CustomTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _CustomTokenPool.ParseOwnershipTransferRequested(log)
	case _CustomTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _CustomTokenPool.ParseOwnershipTransferred(log)
	case _CustomTokenPool.abi.Events["Released"].ID:
		return _CustomTokenPool.ParseReleased(log)
	case _CustomTokenPool.abi.Events["RemotePoolSet"].ID:
		return _CustomTokenPool.ParseRemotePoolSet(log)
	case _CustomTokenPool.abi.Events["RouterUpdated"].ID:
		return _CustomTokenPool.ParseRouterUpdated(log)
	case _CustomTokenPool.abi.Events["SynthBurned"].ID:
		return _CustomTokenPool.ParseSynthBurned(log)
	case _CustomTokenPool.abi.Events["SynthMinted"].ID:
		return _CustomTokenPool.ParseSynthMinted(log)
	case _CustomTokenPool.abi.Events["TokensConsumed"].ID:
		return _CustomTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CustomTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (CustomTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (CustomTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (CustomTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x0f135cbb9afa12a8bf3bbd071c117bcca4ddeca6160ef7f33d012a81b9c0c471")
}

func (CustomTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (CustomTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (CustomTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (CustomTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (CustomTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (CustomTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (CustomTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (CustomTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (CustomTokenPoolRemotePoolSet) Topic() common.Hash {
	return common.HexToHash("0xdb4d6220746a38cbc5335f7e108f7de80f482f4d23350253dfd0917df75a14bf")
}

func (CustomTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (CustomTokenPoolSynthBurned) Topic() common.Hash {
	return common.HexToHash("0x02992093bca69a36949677658a77d359b510dc6232c68f9f118f7c0127a1b147")
}

func (CustomTokenPoolSynthMinted) Topic() common.Hash {
	return common.HexToHash("0xbb0b72e5f44e331506684da008a30e10d50658c29d8159f6c6ab40bf1e52e600")
}

func (CustomTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_CustomTokenPool *CustomTokenPool) Address() common.Address {
	return _CustomTokenPool.address
}

type CustomTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRemotePool(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

	IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, chains []TokenPoolChainUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*CustomTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*CustomTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*CustomTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*CustomTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*CustomTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*CustomTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*CustomTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*CustomTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*CustomTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*CustomTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*CustomTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*CustomTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*CustomTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*CustomTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*CustomTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*CustomTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*CustomTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*CustomTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CustomTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*CustomTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CustomTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*CustomTokenPoolOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*CustomTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*CustomTokenPoolReleased, error)

	FilterRemotePoolSet(opts *bind.FilterOpts, remoteChainSelector []uint64) (*CustomTokenPoolRemotePoolSetIterator, error)

	WatchRemotePoolSet(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolRemotePoolSet, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolSet(log types.Log) (*CustomTokenPoolRemotePoolSet, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*CustomTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*CustomTokenPoolRouterUpdated, error)

	FilterSynthBurned(opts *bind.FilterOpts) (*CustomTokenPoolSynthBurnedIterator, error)

	WatchSynthBurned(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolSynthBurned) (event.Subscription, error)

	ParseSynthBurned(log types.Log) (*CustomTokenPoolSynthBurned, error)

	FilterSynthMinted(opts *bind.FilterOpts) (*CustomTokenPoolSynthMintedIterator, error)

	WatchSynthMinted(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolSynthMinted) (event.Subscription, error)

	ParseSynthMinted(log types.Log) (*CustomTokenPoolSynthMinted, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*CustomTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *CustomTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*CustomTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
