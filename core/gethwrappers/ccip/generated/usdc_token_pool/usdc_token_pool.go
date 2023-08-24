// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package usdc_token_pool

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

type TokenPoolRampUpdate struct {
	Ramp              common.Address
	Allowed           bool
	RateLimiterConfig RateLimiterConfig
}

type USDCTokenPoolDomain struct {
	AllowedCaller    [32]byte
	DomainIdentifier uint32
	Enabled          bool
}

type USDCTokenPoolDomainUpdate struct {
	AllowedCaller     [32]byte
	DomainIdentifier  uint32
	DestChainSelector uint64
	Enabled           bool
}

type USDCTokenPoolUSDCConfig struct {
	Version            uint32
	TokenMessenger     common.Address
	MessageTransmitter common.Address
}

var USDCTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"NonExistentRamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermissionsError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"RampAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"domain\",\"type\":\"uint64\"}],\"name\":\"UnknownDomain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnlockingUSDCFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"DomainsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"OnRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"onRamps\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"offRamps\",\"type\":\"tuple[]\"}],\"name\":\"applyRampUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"currentOffRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"currentOnRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getArmProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getDomain\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.Domain\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOffRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOnRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUSDCInterfaceId\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"isOnRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destinationReceiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"releaseOrMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"domains\",\"type\":\"tuple[]\"}],\"name\":\"setDomains\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOffRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOnRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162003be138038062003be183398101604081905262000034916200080e565b82828233806000816200008e5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c157620000c18162000145565b5050506001600160a01b038316620000ec576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808416608052811660a052815115801560c0526200012757604080516000815260208101909152620001279083620001f0565b5050506200013b846200036160201b60201c565b5050505062000977565b336001600160a01b038216036200019f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000085565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c05162000211576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002a6576000838281518110620002355762000235620008d8565b602090810291909101015190506200024f60028262000566565b1562000292576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506200029e8162000904565b905062000214565b5060005b81518110156200035c576000828281518110620002cb57620002cb620008d8565b6020026020010151905060006001600160a01b0316816001600160a01b031603620002f7575062000349565b6200030460028262000586565b1562000347576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b620003548162000904565b9050620002aa565b505050565b60408101516001600160a01b0316158062000387575060208101516001600160a01b0316155b15620003a6576040516306b7c75960e31b815260040160405180910390fd5b600a5464010000000090046001600160a01b0316156200044657608051600a5460405163095ea7b360e01b81526001600160a01b03640100000000909204821660048201526000602482015291169063095ea7b3906044016020604051808303816000875af11580156200041e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000444919062000920565b505b608051602082015160405163095ea7b360e01b81526001600160a01b039182166004820152600019602482015291169063095ea7b3906044016020604051808303816000875af11580156200049f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004c5919062000920565b508051600a80546020808501805163ffffffff9095166001600160c01b031990931683176401000000006001600160a01b03968716021790935560408086018051600b80546001600160a01b031916918816919091179055815193845293518516918301919091529151909216908201527f33a7d35707e0c8e46d6fa8dd98b73765c14247a559106927070b1cfd2933f4039060600160405180910390a150565b60006200057d836001600160a01b0384166200059d565b90505b92915050565b60006200057d836001600160a01b038416620006a1565b6000818152600183016020526040812054801562000696576000620005c46001836200094b565b8554909150600090620005da906001906200094b565b905081811462000646576000866000018281548110620005fe57620005fe620008d8565b9060005260206000200154905080876000018481548110620006245762000624620008d8565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200065a576200065a62000961565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000580565b600091505062000580565b6000818152600183016020526040812054620006ea5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000580565b50600062000580565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b03811182821017156200072e576200072e620006f3565b60405290565b6001600160a01b03811681146200074a57600080fd5b50565b80516200075a8162000734565b919050565b600082601f8301126200077157600080fd5b815160206001600160401b0380831115620007905762000790620006f3565b8260051b604051601f19603f83011681018181108482111715620007b857620007b8620006f3565b604052938452858101830193838101925087851115620007d757600080fd5b83870191505b8482101562000803578151620007f38162000734565b83529183019190830190620007dd565b979650505050505050565b60008060008084860360c08112156200082657600080fd5b60608112156200083557600080fd5b506200084062000709565b855163ffffffff811681146200085557600080fd5b81526020860151620008678162000734565b602082015260408601516200087c8162000734565b6040820152935062000891606086016200074d565b60808601519093506001600160401b03811115620008ae57600080fd5b620008bc878288016200075f565b925050620008cd60a086016200074d565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201620009195762000919620008ee565b5060010190565b6000602082840312156200093357600080fd5b815180151581146200094457600080fd5b9392505050565b81810381811115620005805762000580620008ee565b634e487b7160e01b600052603160045260246000fd5b60805160a05160c051613209620009d86000396000818161055401528181610c2b0152611507015260006102530152600081816101f901528181610dc4015281816112d9015281816113b801528181611947015261198e01526132096000f3fe608060405234801561001057600080fd5b50600436106101a25760003560e01c80638627fad6116100ee578063b3a3fb4111610097578063d612b94511610071578063d612b945146104a1578063dfadfa35146104b4578063e0351e1314610552578063f2fde38b1461057857600080fd5b8063b3a3fb41146103cb578063c3f909d4146103de578063c49907b51461048e57600080fd5b806396875445116100c8578063968754451461039b578063a40e69c7146103bb578063a7cd63b7146103c357600080fd5b80638627fad61461035557806387381314146103685780638da5cb5b1461037d57600080fd5b806354c8a4f3116101505780637448b3c71161012a5780637448b3c7146102cb5780637787e7ab146102de57806379ba50971461034d57600080fd5b806354c8a4f3146102775780636d1081391461028a5780636f32b872146102b857600080fd5b806321df0da71161018157806321df0da7146101f7578063263a890a1461023e5780635246492f1461025157600080fd5b806241d3c1146101a757806301ffc9a7146101bc5780631d7a74a0146101e4575b600080fd5b6101ba6101b53660046125fa565b61058b565b005b6101cf6101ca36600461266f565b6106a7565b60405190151581526020015b60405180910390f35b6101cf6101f23660046126da565b610703565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101db565b6101ba61024c3660046127b0565b610710565b7f0000000000000000000000000000000000000000000000000000000000000000610219565b6101ba61028536600461284a565b610724565b6040517fd6aca1be0000000000000000000000000000000000000000000000000000000081526020016101db565b6101cf6102c63660046126da565b61079f565b6101ba6102d9366004612935565b6107ac565b6102f16102ec3660046126da565b610870565b6040516101db919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b6101ba61094e565b6101ba610363366004612a1b565b610a4b565b610370610bd6565b6040516101db9190612aae565b60005473ffffffffffffffffffffffffffffffffffffffff16610219565b6103ae6103a9366004612b4a565b610be7565b6040516101db9190612c57565b610370610ebe565b610370610eca565b6102f16103d93660046126da565b610ed6565b61044860408051606081018252600080825260208201819052918101919091525060408051606081018252600a5463ffffffff8116825273ffffffffffffffffffffffffffffffffffffffff64010000000090910481166020830152600b54169181019190915290565b60408051825163ffffffff16815260208084015173ffffffffffffffffffffffffffffffffffffffff9081169183019190915292820151909216908201526060016101db565b6101ba61049c366004612caf565b610fb4565b6101ba6104af366004612935565b610fc8565b6105286104c2366004612d0f565b60408051606080820183526000808352602080840182905292840181905267ffffffffffffffff949094168452600c82529282902082519384018352805484526001015463ffffffff811691840191909152640100000000900460ff1615159082015290565b604080518251815260208084015163ffffffff1690820152918101511515908201526060016101db565b7f00000000000000000000000000000000000000000000000000000000000000006101cf565b6101ba6105863660046126da565b611087565b610593611098565b60005b818110156106695760008383838181106105b2576105b2612d2c565b9050608002018036038101906105c89190612d5b565b60408051606080820183528351825260208085015163ffffffff9081168285019081529286015115158486019081529585015167ffffffffffffffff166000908152600c90925293902091518255516001909101805493511515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000090941691909216179190911790555061066281612e04565b9050610596565b507f1889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee56828260405161069b929190612e3c565b60405180910390a15050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fd6aca1be0000000000000000000000000000000000000000000000000000000014806106fd57506106fd8261111b565b92915050565b60006106fd6007836111b3565b610718611098565b610721816111e5565b50565b61072c611098565b6107998484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061150592505050565b50505050565b60006106fd6004836111b3565b6107b4611098565b6107bd8261079f565b610810576040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216600090815260066020526040902061083f90826116d0565b7f578db78e348076074dbff64a94073a83e9a65aa6766b8c75fdc89282b0e30ed6828260405161069b929190612ec3565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915273ffffffffffffffffffffffffffffffffffffffff8216600090815260066020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526106fd9061187f565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610807565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a5433610703565b610a8a576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a9383611931565b600081806020019051810190610aa99190612f68565b600b54815160208301516040517f57ecfd2800000000000000000000000000000000000000000000000000000000815293945073ffffffffffffffffffffffffffffffffffffffff909216926357ecfd2892610b08929160040161300f565b6020604051808303816000875af1158015610b27573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b4b9190613034565b610b81576040517fbf969f2200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405184815273ffffffffffffffffffffffffffffffffffffffff86169033907f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f09060200160405180910390a3505050505050565b6060610be2600461196b565b905090565b6060610bf23361079f565b610c28576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b877f00000000000000000000000000000000000000000000000000000000000000008015610c5e5750610c5c6002826111b3565b155b15610cad576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610807565b67ffffffffffffffff85166000908152600c602090815260409182902082516060810184528154815260019091015463ffffffff81169282019290925264010000000090910460ff16151591810182905290610d41576040517fd201c48a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152602401610807565b610d4a87611978565b6000610d596020828b8d613051565b610d629161307b565b600a54602084015184516040517ff856ddb6000000000000000000000000000000000000000000000000000000008152600481018d905263ffffffff90921660248301526044820184905273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116606484015260848301919091529293506000926401000000009092049091169063f856ddb69060a4016020604051808303816000875af1158015610e2d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5191906130b7565b6040518a815290915033907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040805167ffffffffffffffff8316602082015201604051602081830303815290604052945050505050979650505050505050565b6060610be2600761196b565b6060610be2600261196b565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915273ffffffffffffffffffffffffffffffffffffffff8216600090815260096020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526106fd9061187f565b610fbc611098565b610799848484846119b2565b610fd0611098565b610fd982610703565b611027576040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610807565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260096020526040902061105690826116d0565b7fb3ba339cfbb8ef80d7a29ce5493051cb90e64fcfa85d7124efc1adfa4c68399f828260405161069b929190612ec3565b61108f611098565b61072181611f62565b60005473ffffffffffffffffffffffffffffffffffffffff163314611119576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610807565b565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f317fa3340000000000000000000000000000000000000000000000000000000014806106fd57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b9392505050565b604081015173ffffffffffffffffffffffffffffffffffffffff1615806112245750602081015173ffffffffffffffffffffffffffffffffffffffff16155b1561125b576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a54640100000000900473ffffffffffffffffffffffffffffffffffffffff161561134857600a546040517f095ea7b300000000000000000000000000000000000000000000000000000000815264010000000090910473ffffffffffffffffffffffffffffffffffffffff9081166004830152600060248301527f0000000000000000000000000000000000000000000000000000000000000000169063095ea7b3906044016020604051808303816000875af1158015611322573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113469190613034565b505b60208101516040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60248201527f00000000000000000000000000000000000000000000000000000000000000009091169063095ea7b3906044016020604051808303816000875af1158015611403573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114279190613034565b508051600a80546020808501805163ffffffff9095167fffffffffffffffff000000000000000000000000000000000000000000000000909316831764010000000073ffffffffffffffffffffffffffffffffffffffff968716021790935560408086018051600b80547fffffffffffffffffffffffff000000000000000000000000000000000000000016918816919091179055815193845293518516918301919091529151909216908201527f33a7d35707e0c8e46d6fa8dd98b73765c14247a559106927070b1cfd2933f4039060600160405180910390a150565b7f000000000000000000000000000000000000000000000000000000000000000061155c576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82518110156115fa57600083828151811061157c5761157c612d2c565b6020026020010151905061159a81600261205790919063ffffffff16565b156115e95760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506115f381612e04565b905061155f565b5060005b81518110156116cb57600082828151811061161b5761161b612d2c565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361165f57506116bb565b61166a600282612079565b156116b95760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6116c481612e04565b90506115fe565b505050565b81546000906116f990700100000000000000000000000000000000900463ffffffff16426130d4565b9050801561179b5760018301548354611741916fffffffffffffffffffffffffffffffff8082169281169185917001000000000000000000000000000000009091041661209b565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546117c1916fffffffffffffffffffffffffffffffff90811691166120c3565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906118729084906130e7565b60405180910390a1505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261190d82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426118f191906130d4565b85608001516fffffffffffffffffffffffffffffffff1661209b565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b33600090815260096020526040902061072190827f00000000000000000000000000000000000000000000000000000000000000006120d9565b606060006111de8361245c565b33600090815260066020526040902061072190827f00000000000000000000000000000000000000000000000000000000000000006120d9565b6119ba611098565b60005b83811015611cd75760008585838181106119d9576119d9612d2c565b905060a002018036038101906119ef9190613123565b9050806020015115611bc7578051611a0990600490612079565b15611b7a576040805160a08101825282820180516020908101516fffffffffffffffffffffffffffffffff908116845263ffffffff4281168386019081528451511515868801908152855185015184166060880190815286518901518516608089019081528a5173ffffffffffffffffffffffffffffffffffffffff1660009081526006909752958990209751885493519251151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff939095167001000000000000000000000000000000009081027fffffffffffffffffffffffff0000000000000000000000000000000000000000909516918716919091179390931791909116929092178655905192518216029116176001909201919091558251905191517f0b594bb0555ff7b252e0c789ccc9d8903fec294172064308727d570505cee1ac92611b6d9291612ec3565b60405180910390a1611cc6565b80516040517fd3eb6bc500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610807565b8051611bd590600490612057565b15611c7957805173ffffffffffffffffffffffffffffffffffffffff1660009081526006602052604080822080547fffffffffffffffffffffff00000000000000000000000000000000000000000016815560010191909155815190517f7fd064821314ad863a0714a3f1229375ace6b6427ed5544b7b2ba1c47b1b529491611b6d9173ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b80516040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610807565b50611cd081612e04565b90506119bd565b5060005b81811015611f5b576000838383818110611cf757611cf7612d2c565b905060a00201803603810190611d0d9190613123565b9050806020015115611e98578051611d2790600790612079565b15611b7a576040805160a08101825282820180516020908101516fffffffffffffffffffffffffffffffff908116845263ffffffff4281168386019081528451511515868801908152855185015184166060880190815286518901518516608089019081528a5173ffffffffffffffffffffffffffffffffffffffff1660009081526009909752958990209751885493519251151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff939095167001000000000000000000000000000000009081027fffffffffffffffffffffffff0000000000000000000000000000000000000000909516918716919091179390931791909116929092178655905192518216029116176001909201919091558251905191517f395b7374909d2b54e5796f53c898ebf41d767c86c78ea86519acf2b805852d8892611e8b9291612ec3565b60405180910390a1611f4a565b8051611ea690600790612057565b15611c7957805173ffffffffffffffffffffffffffffffffffffffff1660009081526009602052604080822080547fffffffffffffffffffffff00000000000000000000000000000000000000000016815560010191909155815190517fcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c91611e8b9173ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b50611f5481612e04565b9050611cdb565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611fe1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610807565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006111de8373ffffffffffffffffffffffffffffffffffffffff84166124b8565b60006111de8373ffffffffffffffffffffffffffffffffffffffff84166125ab565b60006120ba856120ab8486613168565b6120b5908761317f565b6120c3565b95945050505050565b60008183106120d257816111de565b5090919050565b825474010000000000000000000000000000000000000000900460ff161580612100575081155b1561210a57505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061215090700100000000000000000000000000000000900463ffffffff16426130d4565b905080156122105781831115612192576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546121cc9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1661209b565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156122c75773ffffffffffffffffffffffffffffffffffffffff841661226f576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610807565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610807565b848310156123da5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1690600090829061230b90826130d4565b612315878a6130d4565b61231f919061317f565b6123299190613192565b905073ffffffffffffffffffffffffffffffffffffffff8616612382576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610807565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610807565b6123e485846130d4565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b6060816000018054806020026020016040519081016040528092919081815260200182805480156124ac57602002820191906000526020600020905b815481526020019060010190808311612498575b50505050509050919050565b600081815260018301602052604081205480156125a15760006124dc6001836130d4565b85549091506000906124f0906001906130d4565b905081811461255557600086600001828154811061251057612510612d2c565b906000526020600020015490508087600001848154811061253357612533612d2c565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612566576125666131cd565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106fd565b60009150506106fd565b60008181526001830160205260408120546125f2575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106fd565b5060006106fd565b6000806020838503121561260d57600080fd5b823567ffffffffffffffff8082111561262557600080fd5b818501915085601f83011261263957600080fd5b81358181111561264857600080fd5b8660208260071b850101111561265d57600080fd5b60209290920196919550909350505050565b60006020828403121561268157600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146111de57600080fd5b803573ffffffffffffffffffffffffffffffffffffffff811681146126d557600080fd5b919050565b6000602082840312156126ec57600080fd5b6111de826126b1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715612747576127476126f5565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715612794576127946126f5565b604052919050565b803563ffffffff811681146126d557600080fd5b6000606082840312156127c257600080fd5b6127ca612724565b6127d38361279c565b81526127e1602084016126b1565b60208201526127f2604084016126b1565b60408201529392505050565b60008083601f84011261281057600080fd5b50813567ffffffffffffffff81111561282857600080fd5b6020830191508360208260051b850101111561284357600080fd5b9250929050565b6000806000806040858703121561286057600080fd5b843567ffffffffffffffff8082111561287857600080fd5b612884888389016127fe565b9096509450602087013591508082111561289d57600080fd5b506128aa878288016127fe565b95989497509550505050565b801515811461072157600080fd5b80356fffffffffffffffffffffffffffffffff811681146126d557600080fd5b6000606082840312156128f657600080fd5b6128fe612724565b9050813561290b816128b6565b8152612919602083016128c4565b602082015261292a604083016128c4565b604082015292915050565b6000806080838503121561294857600080fd5b612951836126b1565b915061296084602085016128e4565b90509250929050565b600067ffffffffffffffff821115612983576129836126f5565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126129c057600080fd5b81356129d36129ce82612969565b61274d565b8181528460208386010111156129e857600080fd5b816020850160208301376000918101602001919091529392505050565b67ffffffffffffffff8116811461072157600080fd5b600080600080600060a08688031215612a3357600080fd5b853567ffffffffffffffff80821115612a4b57600080fd5b612a5789838a016129af565b9650612a65602089016126b1565b95506040880135945060608801359150612a7e82612a05565b90925060808701359080821115612a9457600080fd5b50612aa1888289016129af565b9150509295509295909350565b6020808252825182820181905260009190848201906040850190845b81811015612afc57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612aca565b50909695505050505050565b60008083601f840112612b1a57600080fd5b50813567ffffffffffffffff811115612b3257600080fd5b60208301915083602082850101111561284357600080fd5b600080600080600080600060a0888a031215612b6557600080fd5b612b6e886126b1565b9650602088013567ffffffffffffffff80821115612b8b57600080fd5b612b978b838c01612b08565b909850965060408a0135955060608a01359150612bb382612a05565b90935060808901359080821115612bc957600080fd5b50612bd68a828b01612b08565b989b979a50959850939692959293505050565b60005b83811015612c04578181015183820152602001612bec565b50506000910152565b60008151808452612c25816020860160208601612be9565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006111de6020830184612c0d565b60008083601f840112612c7c57600080fd5b50813567ffffffffffffffff811115612c9457600080fd5b60208301915083602060a08302850101111561284357600080fd5b60008060008060408587031215612cc557600080fd5b843567ffffffffffffffff80821115612cdd57600080fd5b612ce988838901612c6a565b90965094506020870135915080821115612d0257600080fd5b506128aa87828801612c6a565b600060208284031215612d2157600080fd5b81356111de81612a05565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060808284031215612d6d57600080fd5b6040516080810181811067ffffffffffffffff82111715612d9057612d906126f5565b60405282358152612da36020840161279c565b60208201526040830135612db681612a05565b60408201526060830135612dc9816128b6565b60608201529392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203612e3557612e35612dd5565b5060010190565b6020808252818101839052600090604080840186845b87811015612eb6578135835263ffffffff612e6e86840161279c565b168584015283820135612e8081612a05565b67ffffffffffffffff1683850152606082810135612e9d816128b6565b1515908401526080928301929190910190600101612e52565b5090979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff83168152608081016111de60208301848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600082601f830112612f2c57600080fd5b8151612f3a6129ce82612969565b818152846020838601011115612f4f57600080fd5b612f60826020830160208701612be9565b949350505050565b600060208284031215612f7a57600080fd5b815167ffffffffffffffff80821115612f9257600080fd5b9083019060408286031215612fa657600080fd5b604051604081018181108382111715612fc157612fc16126f5565b604052825182811115612fd357600080fd5b612fdf87828601612f1b565b825250602083015182811115612ff457600080fd5b61300087828601612f1b565b60208301525095945050505050565b6040815260006130226040830185612c0d565b82810360208401526120ba8185612c0d565b60006020828403121561304657600080fd5b81516111de816128b6565b6000808585111561306157600080fd5b8386111561306e57600080fd5b5050820193919092039150565b803560208310156106fd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b6000602082840312156130c957600080fd5b81516111de81612a05565b818103818111156106fd576106fd612dd5565b606081016106fd82848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060a0828403121561313557600080fd5b61313d612724565b613146836126b1565b81526020830135613156816128b6565b60208201526127f284604085016128e4565b80820281158282048414176106fd576106fd612dd5565b808201808211156106fd576106fd612dd5565b6000826131c8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var USDCTokenPoolABI = USDCTokenPoolMetaData.ABI

var USDCTokenPoolBin = USDCTokenPoolMetaData.Bin

func DeployUSDCTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, config USDCTokenPoolUSDCConfig, token common.Address, allowlist []common.Address, armProxy common.Address) (common.Address, *types.Transaction, *USDCTokenPool, error) {
	parsed, err := USDCTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCTokenPoolBin), backend, config, token, allowlist, armProxy)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &USDCTokenPool{USDCTokenPoolCaller: USDCTokenPoolCaller{contract: contract}, USDCTokenPoolTransactor: USDCTokenPoolTransactor{contract: contract}, USDCTokenPoolFilterer: USDCTokenPoolFilterer{contract: contract}}, nil
}

type USDCTokenPool struct {
	address common.Address
	abi     abi.ABI
	USDCTokenPoolCaller
	USDCTokenPoolTransactor
	USDCTokenPoolFilterer
}

type USDCTokenPoolCaller struct {
	contract *bind.BoundContract
}

type USDCTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type USDCTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type USDCTokenPoolSession struct {
	Contract     *USDCTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type USDCTokenPoolCallerSession struct {
	Contract *USDCTokenPoolCaller
	CallOpts bind.CallOpts
}

type USDCTokenPoolTransactorSession struct {
	Contract     *USDCTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type USDCTokenPoolRaw struct {
	Contract *USDCTokenPool
}

type USDCTokenPoolCallerRaw struct {
	Contract *USDCTokenPoolCaller
}

type USDCTokenPoolTransactorRaw struct {
	Contract *USDCTokenPoolTransactor
}

func NewUSDCTokenPool(address common.Address, backend bind.ContractBackend) (*USDCTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(USDCTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUSDCTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPool{address: address, abi: abi, USDCTokenPoolCaller: USDCTokenPoolCaller{contract: contract}, USDCTokenPoolTransactor: USDCTokenPoolTransactor{contract: contract}, USDCTokenPoolFilterer: USDCTokenPoolFilterer{contract: contract}}, nil
}

func NewUSDCTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*USDCTokenPoolCaller, error) {
	contract, err := bindUSDCTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolCaller{contract: contract}, nil
}

func NewUSDCTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*USDCTokenPoolTransactor, error) {
	contract, err := bindUSDCTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolTransactor{contract: contract}, nil
}

func NewUSDCTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*USDCTokenPoolFilterer, error) {
	contract, err := bindUSDCTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolFilterer{contract: contract}, nil
}

func bindUSDCTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := USDCTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_USDCTokenPool *USDCTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCTokenPool.Contract.USDCTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.USDCTokenPoolTransactor.contract.Transfer(opts)
}

func (_USDCTokenPool *USDCTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.USDCTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.contract.Transfer(opts)
}

func (_USDCTokenPool *USDCTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolCaller) CurrentOffRampRateLimiterState(opts *bind.CallOpts, offRamp common.Address) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "currentOffRampRateLimiterState", offRamp)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) CurrentOffRampRateLimiterState(offRamp common.Address) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.CurrentOffRampRateLimiterState(&_USDCTokenPool.CallOpts, offRamp)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) CurrentOffRampRateLimiterState(offRamp common.Address) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.CurrentOffRampRateLimiterState(&_USDCTokenPool.CallOpts, offRamp)
}

func (_USDCTokenPool *USDCTokenPoolCaller) CurrentOnRampRateLimiterState(opts *bind.CallOpts, onRamp common.Address) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "currentOnRampRateLimiterState", onRamp)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) CurrentOnRampRateLimiterState(onRamp common.Address) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.CurrentOnRampRateLimiterState(&_USDCTokenPool.CallOpts, onRamp)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) CurrentOnRampRateLimiterState(onRamp common.Address) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.CurrentOnRampRateLimiterState(&_USDCTokenPool.CallOpts, onRamp)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetAllowList(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetAllowList(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _USDCTokenPool.Contract.GetAllowListEnabled(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _USDCTokenPool.Contract.GetAllowListEnabled(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetArmProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getArmProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetArmProxy() (common.Address, error) {
	return _USDCTokenPool.Contract.GetArmProxy(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetArmProxy() (common.Address, error) {
	return _USDCTokenPool.Contract.GetArmProxy(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetConfig(opts *bind.CallOpts) (USDCTokenPoolUSDCConfig, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(USDCTokenPoolUSDCConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(USDCTokenPoolUSDCConfig)).(*USDCTokenPoolUSDCConfig)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetConfig() (USDCTokenPoolUSDCConfig, error) {
	return _USDCTokenPool.Contract.GetConfig(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetConfig() (USDCTokenPoolUSDCConfig, error) {
	return _USDCTokenPool.Contract.GetConfig(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetDomain(opts *bind.CallOpts, chainSelector uint64) (USDCTokenPoolDomain, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getDomain", chainSelector)

	if err != nil {
		return *new(USDCTokenPoolDomain), err
	}

	out0 := *abi.ConvertType(out[0], new(USDCTokenPoolDomain)).(*USDCTokenPoolDomain)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetDomain(chainSelector uint64) (USDCTokenPoolDomain, error) {
	return _USDCTokenPool.Contract.GetDomain(&_USDCTokenPool.CallOpts, chainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetDomain(chainSelector uint64) (USDCTokenPoolDomain, error) {
	return _USDCTokenPool.Contract.GetDomain(&_USDCTokenPool.CallOpts, chainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetOffRamps(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getOffRamps")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetOffRamps() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetOffRamps(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetOffRamps() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetOffRamps(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetOnRamps(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getOnRamps")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetOnRamps() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetOnRamps(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetOnRamps() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetOnRamps(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetToken() (common.Address, error) {
	return _USDCTokenPool.Contract.GetToken(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _USDCTokenPool.Contract.GetToken(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetUSDCInterfaceId(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getUSDCInterfaceId")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetUSDCInterfaceId() ([4]byte, error) {
	return _USDCTokenPool.Contract.GetUSDCInterfaceId(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetUSDCInterfaceId() ([4]byte, error) {
	return _USDCTokenPool.Contract.GetUSDCInterfaceId(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "isOffRamp", offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsOffRamp(&_USDCTokenPool.CallOpts, offRamp)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsOffRamp(&_USDCTokenPool.CallOpts, offRamp)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "isOnRamp", onRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsOnRamp(&_USDCTokenPool.CallOpts, onRamp)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsOnRamp(&_USDCTokenPool.CallOpts, onRamp)
}

func (_USDCTokenPool *USDCTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) Owner() (common.Address, error) {
	return _USDCTokenPool.Contract.Owner(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) Owner() (common.Address, error) {
	return _USDCTokenPool.Contract.Owner(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _USDCTokenPool.Contract.SupportsInterface(&_USDCTokenPool.CallOpts, interfaceId)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _USDCTokenPool.Contract.SupportsInterface(&_USDCTokenPool.CallOpts, interfaceId)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_USDCTokenPool *USDCTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AcceptOwnership(&_USDCTokenPool.TransactOpts)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AcceptOwnership(&_USDCTokenPool.TransactOpts)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyAllowListUpdates(&_USDCTokenPool.TransactOpts, removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyAllowListUpdates(&_USDCTokenPool.TransactOpts, removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ApplyRampUpdates(opts *bind.TransactOpts, onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "applyRampUpdates", onRamps, offRamps)
}

func (_USDCTokenPool *USDCTokenPoolSession) ApplyRampUpdates(onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyRampUpdates(&_USDCTokenPool.TransactOpts, onRamps, offRamps)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ApplyRampUpdates(onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyRampUpdates(&_USDCTokenPool.TransactOpts, onRamps, offRamps)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, originalSender common.Address, destinationReceiver []byte, amount *big.Int, destChainSelector uint64, arg4 []byte) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "lockOrBurn", originalSender, destinationReceiver, amount, destChainSelector, arg4)
}

func (_USDCTokenPool *USDCTokenPoolSession) LockOrBurn(originalSender common.Address, destinationReceiver []byte, amount *big.Int, destChainSelector uint64, arg4 []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.LockOrBurn(&_USDCTokenPool.TransactOpts, originalSender, destinationReceiver, amount, destChainSelector, arg4)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) LockOrBurn(originalSender common.Address, destinationReceiver []byte, amount *big.Int, destChainSelector uint64, arg4 []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.LockOrBurn(&_USDCTokenPool.TransactOpts, originalSender, destinationReceiver, amount, destChainSelector, arg4)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, arg0 []byte, receiver common.Address, amount *big.Int, arg3 uint64, extraData []byte) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "releaseOrMint", arg0, receiver, amount, arg3, extraData)
}

func (_USDCTokenPool *USDCTokenPoolSession) ReleaseOrMint(arg0 []byte, receiver common.Address, amount *big.Int, arg3 uint64, extraData []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ReleaseOrMint(&_USDCTokenPool.TransactOpts, arg0, receiver, amount, arg3, extraData)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ReleaseOrMint(arg0 []byte, receiver common.Address, amount *big.Int, arg3 uint64, extraData []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ReleaseOrMint(&_USDCTokenPool.TransactOpts, arg0, receiver, amount, arg3, extraData)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetConfig(opts *bind.TransactOpts, config USDCTokenPoolUSDCConfig) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setConfig", config)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetConfig(config USDCTokenPoolUSDCConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetConfig(&_USDCTokenPool.TransactOpts, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetConfig(config USDCTokenPoolUSDCConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetConfig(&_USDCTokenPool.TransactOpts, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetDomains(opts *bind.TransactOpts, domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setDomains", domains)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetDomains(domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetDomains(&_USDCTokenPool.TransactOpts, domains)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetDomains(domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetDomains(&_USDCTokenPool.TransactOpts, domains)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetOffRampRateLimiterConfig(opts *bind.TransactOpts, offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setOffRampRateLimiterConfig", offRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetOffRampRateLimiterConfig(offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetOffRampRateLimiterConfig(&_USDCTokenPool.TransactOpts, offRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetOffRampRateLimiterConfig(offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetOffRampRateLimiterConfig(&_USDCTokenPool.TransactOpts, offRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetOnRampRateLimiterConfig(opts *bind.TransactOpts, onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setOnRampRateLimiterConfig", onRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetOnRampRateLimiterConfig(onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetOnRampRateLimiterConfig(&_USDCTokenPool.TransactOpts, onRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetOnRampRateLimiterConfig(onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetOnRampRateLimiterConfig(&_USDCTokenPool.TransactOpts, onRamp, config)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_USDCTokenPool *USDCTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.TransferOwnership(&_USDCTokenPool.TransactOpts, to)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.TransferOwnership(&_USDCTokenPool.TransactOpts, to)
}

type USDCTokenPoolAllowListAddIterator struct {
	Event *USDCTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolAllowListAdd)
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
		it.Event = new(USDCTokenPoolAllowListAdd)
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

func (it *USDCTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*USDCTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolAllowListAddIterator{contract: _USDCTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolAllowListAdd)
				if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*USDCTokenPoolAllowListAdd, error) {
	event := new(USDCTokenPoolAllowListAdd)
	if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolAllowListRemoveIterator struct {
	Event *USDCTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolAllowListRemove)
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
		it.Event = new(USDCTokenPoolAllowListRemove)
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

func (it *USDCTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*USDCTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolAllowListRemoveIterator{contract: _USDCTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolAllowListRemove)
				if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*USDCTokenPoolAllowListRemove, error) {
	event := new(USDCTokenPoolAllowListRemove)
	if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolBurnedIterator struct {
	Event *USDCTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolBurned)
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
		it.Event = new(USDCTokenPoolBurned)
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

func (it *USDCTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolBurnedIterator{contract: _USDCTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolBurned)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseBurned(log types.Log) (*USDCTokenPoolBurned, error) {
	event := new(USDCTokenPoolBurned)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolConfigSetIterator struct {
	Event *USDCTokenPoolConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolConfigSet)
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
		it.Event = new(USDCTokenPoolConfigSet)
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

func (it *USDCTokenPoolConfigSetIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolConfigSet struct {
	Arg0 USDCTokenPoolUSDCConfig
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterConfigSet(opts *bind.FilterOpts) (*USDCTokenPoolConfigSetIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolConfigSetIterator{contract: _USDCTokenPool.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigSet) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolConfigSet)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseConfigSet(log types.Log) (*USDCTokenPoolConfigSet, error) {
	event := new(USDCTokenPoolConfigSet)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolDomainsSetIterator struct {
	Event *USDCTokenPoolDomainsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolDomainsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolDomainsSet)
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
		it.Event = new(USDCTokenPoolDomainsSet)
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

func (it *USDCTokenPoolDomainsSetIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolDomainsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolDomainsSet struct {
	Arg0 []USDCTokenPoolDomainUpdate
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterDomainsSet(opts *bind.FilterOpts) (*USDCTokenPoolDomainsSetIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "DomainsSet")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolDomainsSetIterator{contract: _USDCTokenPool.contract, event: "DomainsSet", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchDomainsSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolDomainsSet) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "DomainsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolDomainsSet)
				if err := _USDCTokenPool.contract.UnpackLog(event, "DomainsSet", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseDomainsSet(log types.Log) (*USDCTokenPoolDomainsSet, error) {
	event := new(USDCTokenPoolDomainsSet)
	if err := _USDCTokenPool.contract.UnpackLog(event, "DomainsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolLockedIterator struct {
	Event *USDCTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolLocked)
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
		it.Event = new(USDCTokenPoolLocked)
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

func (it *USDCTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolLockedIterator{contract: _USDCTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolLocked)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseLocked(log types.Log) (*USDCTokenPoolLocked, error) {
	event := new(USDCTokenPoolLocked)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolMintedIterator struct {
	Event *USDCTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolMinted)
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
		it.Event = new(USDCTokenPoolMinted)
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

func (it *USDCTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolMintedIterator{contract: _USDCTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolMinted)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseMinted(log types.Log) (*USDCTokenPoolMinted, error) {
	event := new(USDCTokenPoolMinted)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOffRampAddedIterator struct {
	Event *USDCTokenPoolOffRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOffRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOffRampAdded)
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
		it.Event = new(USDCTokenPoolOffRampAdded)
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

func (it *USDCTokenPoolOffRampAddedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOffRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOffRampAdded struct {
	OffRamp           common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOffRampAdded(opts *bind.FilterOpts) (*USDCTokenPoolOffRampAddedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OffRampAdded")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOffRampAddedIterator{contract: _USDCTokenPool.contract, event: "OffRampAdded", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampAdded) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OffRampAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOffRampAdded)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOffRampAdded(log types.Log) (*USDCTokenPoolOffRampAdded, error) {
	event := new(USDCTokenPoolOffRampAdded)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOffRampConfiguredIterator struct {
	Event *USDCTokenPoolOffRampConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOffRampConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOffRampConfigured)
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
		it.Event = new(USDCTokenPoolOffRampConfigured)
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

func (it *USDCTokenPoolOffRampConfiguredIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOffRampConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOffRampConfigured struct {
	OffRamp           common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOffRampConfigured(opts *bind.FilterOpts) (*USDCTokenPoolOffRampConfiguredIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OffRampConfigured")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOffRampConfiguredIterator{contract: _USDCTokenPool.contract, event: "OffRampConfigured", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOffRampConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampConfigured) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OffRampConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOffRampConfigured)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampConfigured", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOffRampConfigured(log types.Log) (*USDCTokenPoolOffRampConfigured, error) {
	event := new(USDCTokenPoolOffRampConfigured)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOffRampRemovedIterator struct {
	Event *USDCTokenPoolOffRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOffRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOffRampRemoved)
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
		it.Event = new(USDCTokenPoolOffRampRemoved)
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

func (it *USDCTokenPoolOffRampRemovedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOffRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOffRampRemoved struct {
	OffRamp common.Address
	Raw     types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOffRampRemoved(opts *bind.FilterOpts) (*USDCTokenPoolOffRampRemovedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OffRampRemoved")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOffRampRemovedIterator{contract: _USDCTokenPool.contract, event: "OffRampRemoved", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampRemoved) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OffRampRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOffRampRemoved)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOffRampRemoved(log types.Log) (*USDCTokenPoolOffRampRemoved, error) {
	event := new(USDCTokenPoolOffRampRemoved)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOnRampAddedIterator struct {
	Event *USDCTokenPoolOnRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOnRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOnRampAdded)
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
		it.Event = new(USDCTokenPoolOnRampAdded)
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

func (it *USDCTokenPoolOnRampAddedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOnRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOnRampAdded struct {
	OnRamp            common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOnRampAdded(opts *bind.FilterOpts) (*USDCTokenPoolOnRampAddedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OnRampAdded")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOnRampAddedIterator{contract: _USDCTokenPool.contract, event: "OnRampAdded", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOnRampAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampAdded) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OnRampAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOnRampAdded)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampAdded", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOnRampAdded(log types.Log) (*USDCTokenPoolOnRampAdded, error) {
	event := new(USDCTokenPoolOnRampAdded)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOnRampConfiguredIterator struct {
	Event *USDCTokenPoolOnRampConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOnRampConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOnRampConfigured)
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
		it.Event = new(USDCTokenPoolOnRampConfigured)
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

func (it *USDCTokenPoolOnRampConfiguredIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOnRampConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOnRampConfigured struct {
	OnRamp            common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOnRampConfigured(opts *bind.FilterOpts) (*USDCTokenPoolOnRampConfiguredIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OnRampConfigured")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOnRampConfiguredIterator{contract: _USDCTokenPool.contract, event: "OnRampConfigured", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOnRampConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampConfigured) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OnRampConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOnRampConfigured)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampConfigured", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOnRampConfigured(log types.Log) (*USDCTokenPoolOnRampConfigured, error) {
	event := new(USDCTokenPoolOnRampConfigured)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOnRampRemovedIterator struct {
	Event *USDCTokenPoolOnRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOnRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOnRampRemoved)
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
		it.Event = new(USDCTokenPoolOnRampRemoved)
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

func (it *USDCTokenPoolOnRampRemovedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOnRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOnRampRemoved struct {
	OnRamp common.Address
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOnRampRemoved(opts *bind.FilterOpts) (*USDCTokenPoolOnRampRemovedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OnRampRemoved")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOnRampRemovedIterator{contract: _USDCTokenPool.contract, event: "OnRampRemoved", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOnRampRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampRemoved) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OnRampRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOnRampRemoved)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampRemoved", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOnRampRemoved(log types.Log) (*USDCTokenPoolOnRampRemoved, error) {
	event := new(USDCTokenPoolOnRampRemoved)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OnRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOwnershipTransferRequestedIterator struct {
	Event *USDCTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOwnershipTransferRequested)
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
		it.Event = new(USDCTokenPoolOwnershipTransferRequested)
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

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOwnershipTransferRequestedIterator{contract: _USDCTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOwnershipTransferRequested)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*USDCTokenPoolOwnershipTransferRequested, error) {
	event := new(USDCTokenPoolOwnershipTransferRequested)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOwnershipTransferredIterator struct {
	Event *USDCTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOwnershipTransferred)
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
		it.Event = new(USDCTokenPoolOwnershipTransferred)
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

func (it *USDCTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOwnershipTransferredIterator{contract: _USDCTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOwnershipTransferred)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*USDCTokenPoolOwnershipTransferred, error) {
	event := new(USDCTokenPoolOwnershipTransferred)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolReleasedIterator struct {
	Event *USDCTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolReleased)
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
		it.Event = new(USDCTokenPoolReleased)
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

func (it *USDCTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolReleasedIterator{contract: _USDCTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolReleased)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseReleased(log types.Log) (*USDCTokenPoolReleased, error) {
	event := new(USDCTokenPoolReleased)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_USDCTokenPool *USDCTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _USDCTokenPool.abi.Events["AllowListAdd"].ID:
		return _USDCTokenPool.ParseAllowListAdd(log)
	case _USDCTokenPool.abi.Events["AllowListRemove"].ID:
		return _USDCTokenPool.ParseAllowListRemove(log)
	case _USDCTokenPool.abi.Events["Burned"].ID:
		return _USDCTokenPool.ParseBurned(log)
	case _USDCTokenPool.abi.Events["ConfigSet"].ID:
		return _USDCTokenPool.ParseConfigSet(log)
	case _USDCTokenPool.abi.Events["DomainsSet"].ID:
		return _USDCTokenPool.ParseDomainsSet(log)
	case _USDCTokenPool.abi.Events["Locked"].ID:
		return _USDCTokenPool.ParseLocked(log)
	case _USDCTokenPool.abi.Events["Minted"].ID:
		return _USDCTokenPool.ParseMinted(log)
	case _USDCTokenPool.abi.Events["OffRampAdded"].ID:
		return _USDCTokenPool.ParseOffRampAdded(log)
	case _USDCTokenPool.abi.Events["OffRampConfigured"].ID:
		return _USDCTokenPool.ParseOffRampConfigured(log)
	case _USDCTokenPool.abi.Events["OffRampRemoved"].ID:
		return _USDCTokenPool.ParseOffRampRemoved(log)
	case _USDCTokenPool.abi.Events["OnRampAdded"].ID:
		return _USDCTokenPool.ParseOnRampAdded(log)
	case _USDCTokenPool.abi.Events["OnRampConfigured"].ID:
		return _USDCTokenPool.ParseOnRampConfigured(log)
	case _USDCTokenPool.abi.Events["OnRampRemoved"].ID:
		return _USDCTokenPool.ParseOnRampRemoved(log)
	case _USDCTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _USDCTokenPool.ParseOwnershipTransferRequested(log)
	case _USDCTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _USDCTokenPool.ParseOwnershipTransferred(log)
	case _USDCTokenPool.abi.Events["Released"].ID:
		return _USDCTokenPool.ParseReleased(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (USDCTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (USDCTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (USDCTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (USDCTokenPoolConfigSet) Topic() common.Hash {
	return common.HexToHash("0x33a7d35707e0c8e46d6fa8dd98b73765c14247a559106927070b1cfd2933f403")
}

func (USDCTokenPoolDomainsSet) Topic() common.Hash {
	return common.HexToHash("0x1889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee56")
}

func (USDCTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (USDCTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (USDCTokenPoolOffRampAdded) Topic() common.Hash {
	return common.HexToHash("0x395b7374909d2b54e5796f53c898ebf41d767c86c78ea86519acf2b805852d88")
}

func (USDCTokenPoolOffRampConfigured) Topic() common.Hash {
	return common.HexToHash("0xb3ba339cfbb8ef80d7a29ce5493051cb90e64fcfa85d7124efc1adfa4c68399f")
}

func (USDCTokenPoolOffRampRemoved) Topic() common.Hash {
	return common.HexToHash("0xcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c")
}

func (USDCTokenPoolOnRampAdded) Topic() common.Hash {
	return common.HexToHash("0x0b594bb0555ff7b252e0c789ccc9d8903fec294172064308727d570505cee1ac")
}

func (USDCTokenPoolOnRampConfigured) Topic() common.Hash {
	return common.HexToHash("0x578db78e348076074dbff64a94073a83e9a65aa6766b8c75fdc89282b0e30ed6")
}

func (USDCTokenPoolOnRampRemoved) Topic() common.Hash {
	return common.HexToHash("0x7fd064821314ad863a0714a3f1229375ace6b6427ed5544b7b2ba1c47b1b5294")
}

func (USDCTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (USDCTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (USDCTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (_USDCTokenPool *USDCTokenPool) Address() common.Address {
	return _USDCTokenPool.address
}

type USDCTokenPoolInterface interface {
	CurrentOffRampRateLimiterState(opts *bind.CallOpts, offRamp common.Address) (RateLimiterTokenBucket, error)

	CurrentOnRampRateLimiterState(opts *bind.CallOpts, onRamp common.Address) (RateLimiterTokenBucket, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetArmProxy(opts *bind.CallOpts) (common.Address, error)

	GetConfig(opts *bind.CallOpts) (USDCTokenPoolUSDCConfig, error)

	GetDomain(opts *bind.CallOpts, chainSelector uint64) (USDCTokenPoolDomain, error)

	GetOffRamps(opts *bind.CallOpts) ([]common.Address, error)

	GetOnRamps(opts *bind.CallOpts) ([]common.Address, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetUSDCInterfaceId(opts *bind.CallOpts) ([4]byte, error)

	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)

	IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyRampUpdates(opts *bind.TransactOpts, onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, originalSender common.Address, destinationReceiver []byte, amount *big.Int, destChainSelector uint64, arg4 []byte) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, arg0 []byte, receiver common.Address, amount *big.Int, arg3 uint64, extraData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, config USDCTokenPoolUSDCConfig) (*types.Transaction, error)

	SetDomains(opts *bind.TransactOpts, domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error)

	SetOffRampRateLimiterConfig(opts *bind.TransactOpts, offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error)

	SetOnRampRateLimiterConfig(opts *bind.TransactOpts, onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*USDCTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*USDCTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*USDCTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*USDCTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*USDCTokenPoolBurned, error)

	FilterConfigSet(opts *bind.FilterOpts) (*USDCTokenPoolConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*USDCTokenPoolConfigSet, error)

	FilterDomainsSet(opts *bind.FilterOpts) (*USDCTokenPoolDomainsSetIterator, error)

	WatchDomainsSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolDomainsSet) (event.Subscription, error)

	ParseDomainsSet(log types.Log) (*USDCTokenPoolDomainsSet, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*USDCTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*USDCTokenPoolMinted, error)

	FilterOffRampAdded(opts *bind.FilterOpts) (*USDCTokenPoolOffRampAddedIterator, error)

	WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampAdded) (event.Subscription, error)

	ParseOffRampAdded(log types.Log) (*USDCTokenPoolOffRampAdded, error)

	FilterOffRampConfigured(opts *bind.FilterOpts) (*USDCTokenPoolOffRampConfiguredIterator, error)

	WatchOffRampConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampConfigured) (event.Subscription, error)

	ParseOffRampConfigured(log types.Log) (*USDCTokenPoolOffRampConfigured, error)

	FilterOffRampRemoved(opts *bind.FilterOpts) (*USDCTokenPoolOffRampRemovedIterator, error)

	WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOffRampRemoved) (event.Subscription, error)

	ParseOffRampRemoved(log types.Log) (*USDCTokenPoolOffRampRemoved, error)

	FilterOnRampAdded(opts *bind.FilterOpts) (*USDCTokenPoolOnRampAddedIterator, error)

	WatchOnRampAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampAdded) (event.Subscription, error)

	ParseOnRampAdded(log types.Log) (*USDCTokenPoolOnRampAdded, error)

	FilterOnRampConfigured(opts *bind.FilterOpts) (*USDCTokenPoolOnRampConfiguredIterator, error)

	WatchOnRampConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampConfigured) (event.Subscription, error)

	ParseOnRampConfigured(log types.Log) (*USDCTokenPoolOnRampConfigured, error)

	FilterOnRampRemoved(opts *bind.FilterOpts) (*USDCTokenPoolOnRampRemovedIterator, error)

	WatchOnRampRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOnRampRemoved) (event.Subscription, error)

	ParseOnRampRemoved(log types.Log) (*USDCTokenPoolOnRampRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*USDCTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*USDCTokenPoolOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*USDCTokenPoolReleased, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
