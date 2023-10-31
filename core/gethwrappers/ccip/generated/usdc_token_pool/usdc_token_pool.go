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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"localDomainIdentifier\",\"type\":\"uint32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expected\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"got\",\"type\":\"uint32\"}],\"name\":\"InvalidDestinationDomain\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.DomainUpdate\",\"name\":\"domain\",\"type\":\"tuple\"}],\"name\":\"InvalidDomain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"}],\"name\":\"InvalidMessageVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expected\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"got\",\"type\":\"uint64\"}],\"name\":\"InvalidNonce\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expected\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"got\",\"type\":\"uint32\"}],\"name\":\"InvalidSourceDomain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"}],\"name\":\"InvalidTokenMessengerVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"NonExistentRamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermissionsError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"RampAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"domain\",\"type\":\"uint64\"}],\"name\":\"UnknownDomain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnlockingUSDCFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"DomainsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"OnRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"SUPPORTED_USDC_VERSION\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"onRamps\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"offRamps\",\"type\":\"tuple[]\"}],\"name\":\"applyRampUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"currentOffRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"currentOnRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getArmProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getDomain\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.Domain\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOffRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOnRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUSDCInterfaceId\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localDomainIdentifier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"isOnRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destinationReceiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"releaseOrMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageTransmitter\",\"type\":\"address\"}],\"internalType\":\"structUSDCTokenPool.USDCConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"domains\",\"type\":\"tuple[]\"}],\"name\":\"setDomains\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOffRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOnRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200495038038062004950833981016040819052620000359162000bce565b83838333806000816200008f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c257620000c28162000150565b5050506001600160a01b038316620000ed576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808416608052811660a052815115801560c0526200012857604080516000815260208101909152620001289083620001fb565b5050506200013c856200036c60201b60201c565b63ffffffff1660e0525062000e1492505050565b336001600160a01b03821603620001aa5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000086565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60c0516200021c576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620002b157600083828151811062000240576200024062000cc4565b602090810291909101015190506200025a60028262000588565b156200029d576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50620002a98162000cf0565b90506200021f565b5060005b815181101562000367576000828281518110620002d657620002d662000cc4565b6020026020010151905060006001600160a01b0316816001600160a01b03160362000302575062000354565b6200030f600282620005a8565b1562000352576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6200035f8162000cf0565b9050620002b5565b505050565b805163ffffffff16156200039f5780516040516334697c6b60e11b815263ffffffff909116600482015260240162000086565b60408101516001600160a01b03161580620003c5575060208101516001600160a01b0316155b15620003e4576040516306b7c75960e31b815260040160405180910390fd5b600081602001516001600160a01b0316639cdbb1816040518163ffffffff1660e01b81526004016020604051808303816000875af11580156200042b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000451919062000d0c565b905063ffffffff81161562000482576040516316ba39c560e31b815263ffffffff8216600482015260240162000086565b600a5464010000000090046001600160a01b031615620004c557600a54608051620004c5916001600160a01b0391821691640100000000909104166000620005bf565b6020820151608051620004e7916001600160a01b0390911690600019620005bf565b8151600a80546020808601805163ffffffff9095166001600160c01b031990931683176401000000006001600160a01b03968716021790935560408087018051600b80546001600160a01b031916918816919091179055815193845293518516918301919091529151909216908201527f33a7d35707e0c8e46d6fa8dd98b73765c14247a559106927070b1cfd2933f4039060600160405180910390a15050565b60006200059f836001600160a01b03841662000709565b90505b92915050565b60006200059f836001600160a01b0384166200080d565b8015806200063d5750604051636eb1769f60e11b81523060048201526001600160a01b03838116602483015284169063dd62ed3e90604401602060405180830381865afa15801562000615573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200063b919062000d2a565b155b620006b15760405162461bcd60e51b815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e636500000000000000000000606482015260840162000086565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b17909152620003679185916200085f16565b60008181526001830160205260408120548015620008025760006200073060018362000d44565b8554909150600090620007469060019062000d44565b9050818114620007b25760008660000182815481106200076a576200076a62000cc4565b906000526020600020015490508087600001848154811062000790576200079062000cc4565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620007c657620007c662000d5a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620005a2565b6000915050620005a2565b60008181526001830160205260408120546200085657508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620005a2565b506000620005a2565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656490820152600090620008ae906001600160a01b03851690849062000930565b805190915015620003675780806020019051810190620008cf919062000d70565b620003675760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000086565b606062000941848460008562000949565b949350505050565b606082471015620009ac5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000086565b600080866001600160a01b03168587604051620009ca919062000dc1565b60006040518083038185875af1925050503d806000811462000a09576040519150601f19603f3d011682016040523d82523d6000602084013e62000a0e565b606091505b50909250905062000a228783838762000a2d565b979650505050505050565b6060831562000aa157825160000362000a99576001600160a01b0385163b62000a995760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000086565b508162000941565b62000941838381511562000ab85781518083602001fd5b8060405162461bcd60e51b815260040162000086919062000ddf565b634e487b7160e01b600052604160045260246000fd5b805163ffffffff8116811462000aff57600080fd5b919050565b6001600160a01b038116811462000b1a57600080fd5b50565b805162000aff8162000b04565b600082601f83011262000b3c57600080fd5b815160206001600160401b038083111562000b5b5762000b5b62000ad4565b8260051b604051601f19603f8301168101818110848211171562000b835762000b8362000ad4565b60405293845285810183019383810192508785111562000ba257600080fd5b83870191505b8482101562000a2257815162000bbe8162000b04565b8352918301919083019062000ba8565b600080600080600085870360e081121562000be857600080fd5b606081121562000bf757600080fd5b50604051606081016001600160401b03808211838310171562000c1e5762000c1e62000ad4565b8160405262000c2d8962000aea565b83526020890151915062000c418262000b04565b8160208401526040890151915062000c598262000b04565b81604084015282975062000c7060608a0162000b1d565b9650608089015192508083111562000c8757600080fd5b505062000c978882890162000b2a565b93505062000ca860a0870162000b1d565b915062000cb860c0870162000aea565b90509295509295909350565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006001820162000d055762000d0562000cda565b5060010190565b60006020828403121562000d1f57600080fd5b6200059f8262000aea565b60006020828403121562000d3d57600080fd5b5051919050565b81810381811115620005a257620005a262000cda565b634e487b7160e01b600052603160045260246000fd5b60006020828403121562000d8357600080fd5b8151801515811462000d9457600080fd5b9392505050565b60005b8381101562000db857818101518382015260200162000d9e565b50506000910152565b6000825162000dd581846020870162000d9b565b9190910192915050565b602081526000825180602084015262000e0081604085016020870162000d9b565b601f01601f19169190910160400192915050565b60805160a05160c05160e051613abc62000e94600039600081816102f90152818161101401528181611be40152611c420152600081816105f501528181610d9601526116c8015260006102bd01526000818161026301528181610f2f01528181611556015281816115a801528181611b080152611d000152613abc6000f3fe608060405234801561001057600080fd5b50600436106101c35760003560e01c80638627fad6116100f9578063b3a3fb4111610097578063d612b94511610071578063d612b94514610542578063dfadfa3514610555578063e0351e13146105f3578063f2fde38b1461061957600080fd5b8063b3a3fb411461046c578063c3f909d41461047f578063c49907b51461052f57600080fd5b806396875445116100d357806396875445146104415780639fdf13ff14610454578063a40e69c71461045c578063a7cd63b71461046457600080fd5b80638627fad6146103fb578063873813141461040e5780638da5cb5b1461042357600080fd5b806354c8a4f3116101665780636f32b872116101405780636f32b8721461035e5780637448b3c7146103715780637787e7ab1461038457806379ba5097146103f357600080fd5b806354c8a4f3146102e15780636b716b0d146102f45780636d1081391461033057600080fd5b80631d7a74a0116101a25780631d7a74a01461024e57806321df0da714610261578063263a890a146102a85780635246492f146102bb57600080fd5b806241d3c1146101c857806301ffc9a7146101dd578063181f5a7714610205575b600080fd5b6101db6101d6366004612da4565b61062c565b005b6101f06101eb366004612e19565b6107d3565b60405190151581526020015b60405180910390f35b6102416040518060400160405280601381526020017f55534443546f6b656e506f6f6c20312e322e300000000000000000000000000081525081565b6040516101fc9190612ec9565b6101f061025c366004612f05565b61082f565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fc565b6101db6102b6366004612ffc565b61083c565b7f0000000000000000000000000000000000000000000000000000000000000000610283565b6101db6102ef366004613098565b610850565b61031b7f000000000000000000000000000000000000000000000000000000000000000081565b60405163ffffffff90911681526020016101fc565b6040517fd6aca1be0000000000000000000000000000000000000000000000000000000081526020016101fc565b6101f061036c366004612f05565b6108cb565b6101db61037f366004613183565b6108d8565b610397610392366004612f05565b610997565b6040516101fc919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b6101db610a75565b6101db610409366004613269565b610b72565b610416610d41565b6040516101fc91906132fc565b60005473ffffffffffffffffffffffffffffffffffffffff16610283565b61024161044f366004613398565b610d52565b61031b600081565b61041661106c565b610416611078565b61039761047a366004612f05565b611084565b6104e960408051606081018252600080825260208201819052918101919091525060408051606081018252600a5463ffffffff8116825273ffffffffffffffffffffffffffffffffffffffff64010000000090910481166020830152600b54169181019190915290565b60408051825163ffffffff16815260208084015173ffffffffffffffffffffffffffffffffffffffff9081169183019190915292820151909216908201526060016101fc565b6101db61053d36600461347c565b611162565b6101db610550366004613183565b611176565b6105c96105633660046134dc565b60408051606080820183526000808352602080840182905292840181905267ffffffffffffffff949094168452600c82529282902082519384018352805484526001015463ffffffff811691840191909152640100000000900460ff1615159082015290565b604080518251815260208084015163ffffffff1690820152918101511515908201526060016101fc565b7f00000000000000000000000000000000000000000000000000000000000000006101f0565b6101db610627366004612f05565b611235565b610634611246565b60005b81811015610795576000838383818110610653576106536134f9565b9050608002018036038101906106699190613528565b805190915015806106865750604081015167ffffffffffffffff16155b156106f557604080517fa087bd2900000000000000000000000000000000000000000000000000000000815282516004820152602083015163ffffffff1660248201529082015167ffffffffffffffff1660448201526060820151151560648201526084015b60405180910390fd5b60408051606080820183528351825260208085015163ffffffff9081168285019081529286015115158486019081529585015167ffffffffffffffff166000908152600c90925293902091518255516001909101805493511515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000909416919092161791909117905561078e816135d3565b9050610637565b507f1889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee5682826040516107c792919061360b565b60405180910390a15050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fd6aca1be0000000000000000000000000000000000000000000000000000000014806108295750610829826112c9565b92915050565b6000610829600783611361565b610844611246565b61084d81611393565b50565b610858611246565b6108c5848480806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250506040805160208088028281018201909352878252909350879250869182918501908490808284376000920191909152506116c692505050565b50505050565b6000610829600483611361565b6108e0611246565b6108e9826108cb565b610937576040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016106ec565b73ffffffffffffffffffffffffffffffffffffffff821660009081526006602052604090206109669082611891565b7f578db78e348076074dbff64a94073a83e9a65aa6766b8c75fdc89282b0e30ed682826040516107c7929190613694565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915273ffffffffffffffffffffffffffffffffffffffff8216600090815260066020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261082990611a40565b60015473ffffffffffffffffffffffffffffffffffffffff163314610af6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016106ec565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b7b3361082f565b610bb1576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610bba83611af2565b60008082806020019051810190610bd19190613731565b91509150600082806020019051810190610beb9190613795565b9050600082806020019051810190610c0391906137d6565b9050610c13816000015183611b2c565b600b54815160208301516040517f57ecfd2800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909316926357ecfd2892610c70929091600401613867565b6020604051808303816000875af1158015610c8f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cb39190613895565b610ce9576040517fbf969f2200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405187815273ffffffffffffffffffffffffffffffffffffffff89169033907f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f09060200160405180910390a3505050505050505050565b6060610d4d6004611cdd565b905090565b6060610d5d336108cb565b610d93576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b877f00000000000000000000000000000000000000000000000000000000000000008015610dc95750610dc7600282611361565b155b15610e18576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016106ec565b67ffffffffffffffff85166000908152600c602090815260409182902082516060810184528154815260019091015463ffffffff81169282019290925264010000000090910460ff16151591810182905290610eac576040517fd201c48a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff871660048201526024016106ec565b610eb587611cea565b6000610ec46020828b8d6138b2565b610ecd916138dc565b600a54602084015184516040517ff856ddb6000000000000000000000000000000000000000000000000000000008152600481018d905263ffffffff90921660248301526044820184905273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116606484015260848301919091529293506000926401000000009092049091169063f856ddb69060a4016020604051808303816000875af1158015610f98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fbc9190613918565b6040518a815290915033907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a260408051808201825267ffffffffffffffff9290921680835263ffffffff7f00000000000000000000000000000000000000000000000000000000000000008116602094850190815283519485019290925290511682820152805180830382018152606090920190529b9a5050505050505050505050565b6060610d4d6007611cdd565b6060610d4d6002611cdd565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915273ffffffffffffffffffffffffffffffffffffffff8216600090815260096020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261082990611a40565b61116a611246565b6108c584848484611d24565b61117e611246565b6111878261082f565b6111d5576040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016106ec565b73ffffffffffffffffffffffffffffffffffffffff821660009081526009602052604090206112049082611891565b7fb3ba339cfbb8ef80d7a29ce5493051cb90e64fcfa85d7124efc1adfa4c68399f82826040516107c7929190613694565b61123d611246565b61084d816122d4565b60005473ffffffffffffffffffffffffffffffffffffffff1633146112c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016106ec565b565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f317fa33400000000000000000000000000000000000000000000000000000000148061082957507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b9392505050565b805163ffffffff16156113dd5780516040517f68d2f8d600000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526024016106ec565b604081015173ffffffffffffffffffffffffffffffffffffffff16158061141c5750602081015173ffffffffffffffffffffffffffffffffffffffff16155b15611453576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000816020015173ffffffffffffffffffffffffffffffffffffffff16639cdbb1816040518163ffffffff1660e01b81526004016020604051808303816000875af11580156114a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ca9190613935565b905063ffffffff811615611512576040517fb5d1ce2800000000000000000000000000000000000000000000000000000000815263ffffffff821660048201526024016106ec565b600a54640100000000900473ffffffffffffffffffffffffffffffffffffffff161561158857600a546115889073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081169164010000000090041660006123c9565b60208201516115ef9073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6123c9565b8151600a80546020808601805163ffffffff9095167fffffffffffffffff000000000000000000000000000000000000000000000000909316831764010000000073ffffffffffffffffffffffffffffffffffffffff968716021790935560408087018051600b80547fffffffffffffffffffffffff000000000000000000000000000000000000000016918816919091179055815193845293518516918301919091529151909216908201527f33a7d35707e0c8e46d6fa8dd98b73765c14247a559106927070b1cfd2933f403906060016107c7565b7f000000000000000000000000000000000000000000000000000000000000000061171d576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82518110156117bb57600083828151811061173d5761173d6134f9565b6020026020010151905061175b81600261258290919063ffffffff16565b156117aa5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b506117b4816135d3565b9050611720565b5060005b815181101561188c5760008282815181106117dc576117dc6134f9565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611820575061187c565b61182b6002826125a4565b1561187a5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b611885816135d3565b90506117bf565b505050565b81546000906118ba90700100000000000000000000000000000000900463ffffffff1642613952565b9050801561195c5760018301548354611902916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166125c6565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354611982916fffffffffffffffffffffffffffffffff90811691166125f0565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990611a33908490613965565b60405180910390a1505050565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611ace82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611ab29190613952565b85608001516fffffffffffffffffffffffffffffffff166125c6565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b33600090815260096020526040902061084d90827f0000000000000000000000000000000000000000000000000000000000000000612606565b600482015163ffffffff811615611b77576040517f68d2f8d600000000000000000000000000000000000000000000000000000000815263ffffffff821660048201526024016106ec565b6008830151600c8401516014850151602085015163ffffffff808516911614611be25760208501516040517fe366a11700000000000000000000000000000000000000000000000000000000815263ffffffff918216600482015290841660248201526044016106ec565b7f000000000000000000000000000000000000000000000000000000000000000063ffffffff168263ffffffff1614611c77576040517f77e4802600000000000000000000000000000000000000000000000000000000815263ffffffff7f000000000000000000000000000000000000000000000000000000000000000081166004830152831660248201526044016106ec565b845167ffffffffffffffff828116911614611cd55784516040517ff917ffea00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff918216600482015290821660248201526044016106ec565b505050505050565b6060600061138c83612989565b33600090815260066020526040902061084d90827f0000000000000000000000000000000000000000000000000000000000000000612606565b611d2c611246565b60005b83811015612049576000858583818110611d4b57611d4b6134f9565b905060a00201803603810190611d6191906139a1565b9050806020015115611f39578051611d7b906004906125a4565b15611eec576040805160a08101825282820180516020908101516fffffffffffffffffffffffffffffffff908116845263ffffffff4281168386019081528451511515868801908152855185015184166060880190815286518901518516608089019081528a5173ffffffffffffffffffffffffffffffffffffffff1660009081526006909752958990209751885493519251151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff939095167001000000000000000000000000000000009081027fffffffffffffffffffffffff0000000000000000000000000000000000000000909516918716919091179390931791909116929092178655905192518216029116176001909201919091558251905191517f0b594bb0555ff7b252e0c789ccc9d8903fec294172064308727d570505cee1ac92611edf9291613694565b60405180910390a1612038565b80516040517fd3eb6bc500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016106ec565b8051611f4790600490612582565b15611feb57805173ffffffffffffffffffffffffffffffffffffffff1660009081526006602052604080822080547fffffffffffffffffffffff00000000000000000000000000000000000000000016815560010191909155815190517f7fd064821314ad863a0714a3f1229375ace6b6427ed5544b7b2ba1c47b1b529491611edf9173ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b80516040517f498f12f600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016106ec565b50612042816135d3565b9050611d2f565b5060005b818110156122cd576000838383818110612069576120696134f9565b905060a0020180360381019061207f91906139a1565b905080602001511561220a578051612099906007906125a4565b15611eec576040805160a08101825282820180516020908101516fffffffffffffffffffffffffffffffff908116845263ffffffff4281168386019081528451511515868801908152855185015184166060880190815286518901518516608089019081528a5173ffffffffffffffffffffffffffffffffffffffff1660009081526009909752958990209751885493519251151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff939095167001000000000000000000000000000000009081027fffffffffffffffffffffffff0000000000000000000000000000000000000000909516918716919091179390931791909116929092178655905192518216029116176001909201919091558251905191517f395b7374909d2b54e5796f53c898ebf41d767c86c78ea86519acf2b805852d88926121fd9291613694565b60405180910390a16122bc565b805161221890600790612582565b15611feb57805173ffffffffffffffffffffffffffffffffffffffff1660009081526009602052604080822080547fffffffffffffffffffffff00000000000000000000000000000000000000000016815560010191909155815190517fcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c916121fd9173ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b506122c6816135d3565b905061204d565b5050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612353576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106ec565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80158061246957506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa158015612443573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061246791906139e6565b155b6124f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e63650000000000000000000060648201526084016106ec565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f095ea7b30000000000000000000000000000000000000000000000000000000017905261188c9084906129e5565b600061138c8373ffffffffffffffffffffffffffffffffffffffff8416612af1565b600061138c8373ffffffffffffffffffffffffffffffffffffffff8416612be4565b60006125e5856125d684866139ff565b6125e09087613a16565b6125f0565b90505b949350505050565b60008183106125ff578161138c565b5090919050565b825474010000000000000000000000000000000000000000900460ff16158061262d575081155b1561263757505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061267d90700100000000000000000000000000000000900463ffffffff1642613952565b9050801561273d57818311156126bf576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546126f99083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166125c6565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156127f45773ffffffffffffffffffffffffffffffffffffffff841661279c576040517ff94ebcd100000000000000000000000000000000000000000000000000000000815260048101839052602481018690526044016106ec565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff851660448201526064016106ec565b848310156129075760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff169060009082906128389082613952565b612842878a613952565b61284c9190613a16565b6128569190613a29565b905073ffffffffffffffffffffffffffffffffffffffff86166128af576040517f15279c0800000000000000000000000000000000000000000000000000000000815260048101829052602481018690526044016106ec565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff871660448201526064016106ec565b6129118584613952565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b6060816000018054806020026020016040519081016040528092919081815260200182805480156129d957602002820191906000526020600020905b8154815260200190600101908083116129c5575b50505050509050919050565b6000612a47826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16612c339092919063ffffffff16565b80519091501561188c5780806020019051810190612a659190613895565b61188c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016106ec565b60008181526001830160205260408120548015612bda576000612b15600183613952565b8554909150600090612b2990600190613952565b9050818114612b8e576000866000018281548110612b4957612b496134f9565b9060005260206000200154905080876000018481548110612b6c57612b6c6134f9565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612b9f57612b9f613a64565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610829565b6000915050610829565b6000818152600183016020526040812054612c2b57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610829565b506000610829565b60606125e88484600085856000808673ffffffffffffffffffffffffffffffffffffffff168587604051612c679190613a93565b60006040518083038185875af1925050503d8060008114612ca4576040519150601f19603f3d011682016040523d82523d6000602084013e612ca9565b606091505b5091509150612cba87838387612cc5565b979650505050505050565b60608315612d5b578251600003612d545773ffffffffffffffffffffffffffffffffffffffff85163b612d54576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016106ec565b50816125e8565b6125e88383815115612d705781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ec9190612ec9565b60008060208385031215612db757600080fd5b823567ffffffffffffffff80821115612dcf57600080fd5b818501915085601f830112612de357600080fd5b813581811115612df257600080fd5b8660208260071b8501011115612e0757600080fd5b60209290920196919550909350505050565b600060208284031215612e2b57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461138c57600080fd5b60005b83811015612e76578181015183820152602001612e5e565b50506000910152565b60008151808452612e97816020860160208601612e5b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061138c6020830184612e7f565b803573ffffffffffffffffffffffffffffffffffffffff81168114612f0057600080fd5b919050565b600060208284031215612f1757600080fd5b61138c82612edc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715612f7257612f72612f20565b60405290565b6040805190810167ffffffffffffffff81118282101715612f7257612f72612f20565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715612fe257612fe2612f20565b604052919050565b63ffffffff8116811461084d57600080fd5b60006060828403121561300e57600080fd5b613016612f4f565b823561302181612fea565b815261302f60208401612edc565b602082015261304060408401612edc565b60408201529392505050565b60008083601f84011261305e57600080fd5b50813567ffffffffffffffff81111561307657600080fd5b6020830191508360208260051b850101111561309157600080fd5b9250929050565b600080600080604085870312156130ae57600080fd5b843567ffffffffffffffff808211156130c657600080fd5b6130d28883890161304c565b909650945060208701359150808211156130eb57600080fd5b506130f88782880161304c565b95989497509550505050565b801515811461084d57600080fd5b80356fffffffffffffffffffffffffffffffff81168114612f0057600080fd5b60006060828403121561314457600080fd5b61314c612f4f565b9050813561315981613104565b815261316760208301613112565b602082015261317860408301613112565b604082015292915050565b6000806080838503121561319657600080fd5b61319f83612edc565b91506131ae8460208501613132565b90509250929050565b600067ffffffffffffffff8211156131d1576131d1612f20565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261320e57600080fd5b813561322161321c826131b7565b612f9b565b81815284602083860101111561323657600080fd5b816020850160208301376000918101602001919091529392505050565b67ffffffffffffffff8116811461084d57600080fd5b600080600080600060a0868803121561328157600080fd5b853567ffffffffffffffff8082111561329957600080fd5b6132a589838a016131fd565b96506132b360208901612edc565b955060408801359450606088013591506132cc82613253565b909250608087013590808211156132e257600080fd5b506132ef888289016131fd565b9150509295509295909350565b6020808252825182820181905260009190848201906040850190845b8181101561334a57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613318565b50909695505050505050565b60008083601f84011261336857600080fd5b50813567ffffffffffffffff81111561338057600080fd5b60208301915083602082850101111561309157600080fd5b600080600080600080600060a0888a0312156133b357600080fd5b6133bc88612edc565b9650602088013567ffffffffffffffff808211156133d957600080fd5b6133e58b838c01613356565b909850965060408a0135955060608a0135915061340182613253565b9093506080890135908082111561341757600080fd5b506134248a828b01613356565b989b979a50959850939692959293505050565b60008083601f84011261344957600080fd5b50813567ffffffffffffffff81111561346157600080fd5b60208301915083602060a08302850101111561309157600080fd5b6000806000806040858703121561349257600080fd5b843567ffffffffffffffff808211156134aa57600080fd5b6134b688838901613437565b909650945060208701359150808211156134cf57600080fd5b506130f887828801613437565b6000602082840312156134ee57600080fd5b813561138c81613253565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006080828403121561353a57600080fd5b6040516080810181811067ffffffffffffffff8211171561355d5761355d612f20565b60405282358152602083013561357281612fea565b6020820152604083013561358581613253565b6040820152606083013561359881613104565b60608201529392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203613604576136046135a4565b5060010190565b6020808252818101839052600090604080840186845b8781101561368757813583528482013561363a81612fea565b63ffffffff16838601528184013561365181613253565b67ffffffffffffffff168385015260608281013561366e81613104565b1515908401526080928301929190910190600101613621565b5090979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff831681526080810161138c60208301848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600082601f8301126136fd57600080fd5b815161370b61321c826131b7565b81815284602083860101111561372057600080fd5b6125e8826020830160208701612e5b565b6000806040838503121561374457600080fd5b825167ffffffffffffffff8082111561375c57600080fd5b613768868387016136ec565b9350602085015191508082111561377e57600080fd5b5061378b858286016136ec565b9150509250929050565b6000604082840312156137a757600080fd5b6137af612f78565b82516137ba81613253565b815260208301516137ca81612fea565b60208201529392505050565b6000602082840312156137e857600080fd5b815167ffffffffffffffff8082111561380057600080fd5b908301906040828603121561381457600080fd5b61381c612f78565b82518281111561382b57600080fd5b613837878286016136ec565b82525060208301518281111561384c57600080fd5b613858878286016136ec565b60208301525095945050505050565b60408152600061387a6040830185612e7f565b828103602084015261388c8185612e7f565b95945050505050565b6000602082840312156138a757600080fd5b815161138c81613104565b600080858511156138c257600080fd5b838611156138cf57600080fd5b5050820193919092039150565b80356020831015610829577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b60006020828403121561392a57600080fd5b815161138c81613253565b60006020828403121561394757600080fd5b815161138c81612fea565b81810381811115610829576108296135a4565b6060810161082982848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060a082840312156139b357600080fd5b6139bb612f4f565b6139c483612edc565b815260208301356139d481613104565b60208201526130408460408501613132565b6000602082840312156139f857600080fd5b5051919050565b8082028115828204841417610829576108296135a4565b80820180821115610829576108296135a4565b600082613a5f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b60008251613aa5818460208701612e5b565b919091019291505056fea164736f6c6343000813000a",
}

var USDCTokenPoolABI = USDCTokenPoolMetaData.ABI

var USDCTokenPoolBin = USDCTokenPoolMetaData.Bin

func DeployUSDCTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, config USDCTokenPoolUSDCConfig, token common.Address, allowlist []common.Address, armProxy common.Address, localDomainIdentifier uint32) (common.Address, *types.Transaction, *USDCTokenPool, error) {
	parsed, err := USDCTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCTokenPoolBin), backend, config, token, allowlist, armProxy, localDomainIdentifier)
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

func (_USDCTokenPool *USDCTokenPoolCaller) SUPPORTEDUSDCVERSION(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "SUPPORTED_USDC_VERSION")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) SUPPORTEDUSDCVERSION() (uint32, error) {
	return _USDCTokenPool.Contract.SUPPORTEDUSDCVERSION(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) SUPPORTEDUSDCVERSION() (uint32, error) {
	return _USDCTokenPool.Contract.SUPPORTEDUSDCVERSION(&_USDCTokenPool.CallOpts)
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

func (_USDCTokenPool *USDCTokenPoolCaller) ILocalDomainIdentifier(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "i_localDomainIdentifier")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) ILocalDomainIdentifier() (uint32, error) {
	return _USDCTokenPool.Contract.ILocalDomainIdentifier(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) ILocalDomainIdentifier() (uint32, error) {
	return _USDCTokenPool.Contract.ILocalDomainIdentifier(&_USDCTokenPool.CallOpts)
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

func (_USDCTokenPool *USDCTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) TypeAndVersion() (string, error) {
	return _USDCTokenPool.Contract.TypeAndVersion(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _USDCTokenPool.Contract.TypeAndVersion(&_USDCTokenPool.CallOpts)
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
	SUPPORTEDUSDCVERSION(opts *bind.CallOpts) (uint32, error)

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

	ILocalDomainIdentifier(opts *bind.CallOpts) (uint32, error)

	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)

	IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

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
