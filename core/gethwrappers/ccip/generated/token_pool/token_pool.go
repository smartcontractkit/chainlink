// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token_pool

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

var TokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"NonExistentRamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermissionsError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"}],\"name\":\"RampAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OffRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"OnRampConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"OnRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"onRamps\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.RampUpdate[]\",\"name\":\"offRamps\",\"type\":\"tuple[]\"}],\"name\":\"applyRampUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"currentOffRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"currentOnRampRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getArmProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOffRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOnRamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"isOnRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"releaseOrMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOffRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setOnRampRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var TokenPoolABI = TokenPoolMetaData.ABI

type TokenPool struct {
	address common.Address
	abi     abi.ABI
	TokenPoolCaller
	TokenPoolTransactor
	TokenPoolFilterer
}

type TokenPoolCaller struct {
	contract *bind.BoundContract
}

type TokenPoolTransactor struct {
	contract *bind.BoundContract
}

type TokenPoolFilterer struct {
	contract *bind.BoundContract
}

type TokenPoolSession struct {
	Contract     *TokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TokenPoolCallerSession struct {
	Contract *TokenPoolCaller
	CallOpts bind.CallOpts
}

type TokenPoolTransactorSession struct {
	Contract     *TokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type TokenPoolRaw struct {
	Contract *TokenPool
}

type TokenPoolCallerRaw struct {
	Contract *TokenPoolCaller
}

type TokenPoolTransactorRaw struct {
	Contract *TokenPoolTransactor
}

func NewTokenPool(address common.Address, backend bind.ContractBackend) (*TokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(TokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenPool{address: address, abi: abi, TokenPoolCaller: TokenPoolCaller{contract: contract}, TokenPoolTransactor: TokenPoolTransactor{contract: contract}, TokenPoolFilterer: TokenPoolFilterer{contract: contract}}, nil
}

func NewTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*TokenPoolCaller, error) {
	contract, err := bindTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenPoolCaller{contract: contract}, nil
}

func NewTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenPoolTransactor, error) {
	contract, err := bindTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenPoolTransactor{contract: contract}, nil
}

func NewTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenPoolFilterer, error) {
	contract, err := bindTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenPoolFilterer{contract: contract}, nil
}

func bindTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TokenPool *TokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenPool.Contract.TokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_TokenPool *TokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPool.Contract.TokenPoolTransactor.contract.Transfer(opts)
}

func (_TokenPool *TokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenPool.Contract.TokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_TokenPool *TokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_TokenPool *TokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPool.Contract.contract.Transfer(opts)
}

func (_TokenPool *TokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_TokenPool *TokenPoolCaller) CurrentOffRampRateLimiterState(opts *bind.CallOpts, offRamp common.Address) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "currentOffRampRateLimiterState", offRamp)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_TokenPool *TokenPoolSession) CurrentOffRampRateLimiterState(offRamp common.Address) (RateLimiterTokenBucket, error) {
	return _TokenPool.Contract.CurrentOffRampRateLimiterState(&_TokenPool.CallOpts, offRamp)
}

func (_TokenPool *TokenPoolCallerSession) CurrentOffRampRateLimiterState(offRamp common.Address) (RateLimiterTokenBucket, error) {
	return _TokenPool.Contract.CurrentOffRampRateLimiterState(&_TokenPool.CallOpts, offRamp)
}

func (_TokenPool *TokenPoolCaller) CurrentOnRampRateLimiterState(opts *bind.CallOpts, onRamp common.Address) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "currentOnRampRateLimiterState", onRamp)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_TokenPool *TokenPoolSession) CurrentOnRampRateLimiterState(onRamp common.Address) (RateLimiterTokenBucket, error) {
	return _TokenPool.Contract.CurrentOnRampRateLimiterState(&_TokenPool.CallOpts, onRamp)
}

func (_TokenPool *TokenPoolCallerSession) CurrentOnRampRateLimiterState(onRamp common.Address) (RateLimiterTokenBucket, error) {
	return _TokenPool.Contract.CurrentOnRampRateLimiterState(&_TokenPool.CallOpts, onRamp)
}

func (_TokenPool *TokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _TokenPool.Contract.GetAllowList(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _TokenPool.Contract.GetAllowList(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _TokenPool.Contract.GetAllowListEnabled(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _TokenPool.Contract.GetAllowListEnabled(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) GetArmProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getArmProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetArmProxy() (common.Address, error) {
	return _TokenPool.Contract.GetArmProxy(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetArmProxy() (common.Address, error) {
	return _TokenPool.Contract.GetArmProxy(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) GetOffRamps(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getOffRamps")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetOffRamps() ([]common.Address, error) {
	return _TokenPool.Contract.GetOffRamps(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetOffRamps() ([]common.Address, error) {
	return _TokenPool.Contract.GetOffRamps(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) GetOnRamps(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getOnRamps")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetOnRamps() ([]common.Address, error) {
	return _TokenPool.Contract.GetOnRamps(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetOnRamps() ([]common.Address, error) {
	return _TokenPool.Contract.GetOnRamps(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) GetToken() (common.Address, error) {
	return _TokenPool.Contract.GetToken(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) GetToken() (common.Address, error) {
	return _TokenPool.Contract.GetToken(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "isOffRamp", offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenPool *TokenPoolSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _TokenPool.Contract.IsOffRamp(&_TokenPool.CallOpts, offRamp)
}

func (_TokenPool *TokenPoolCallerSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _TokenPool.Contract.IsOffRamp(&_TokenPool.CallOpts, offRamp)
}

func (_TokenPool *TokenPoolCaller) IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "isOnRamp", onRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenPool *TokenPoolSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _TokenPool.Contract.IsOnRamp(&_TokenPool.CallOpts, onRamp)
}

func (_TokenPool *TokenPoolCallerSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _TokenPool.Contract.IsOnRamp(&_TokenPool.CallOpts, onRamp)
}

func (_TokenPool *TokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenPool *TokenPoolSession) Owner() (common.Address, error) {
	return _TokenPool.Contract.Owner(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCallerSession) Owner() (common.Address, error) {
	return _TokenPool.Contract.Owner(&_TokenPool.CallOpts)
}

func (_TokenPool *TokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenPool *TokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TokenPool.Contract.SupportsInterface(&_TokenPool.CallOpts, interfaceId)
}

func (_TokenPool *TokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TokenPool.Contract.SupportsInterface(&_TokenPool.CallOpts, interfaceId)
}

func (_TokenPool *TokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_TokenPool *TokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenPool.Contract.AcceptOwnership(&_TokenPool.TransactOpts)
}

func (_TokenPool *TokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenPool.Contract.AcceptOwnership(&_TokenPool.TransactOpts)
}

func (_TokenPool *TokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_TokenPool *TokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _TokenPool.Contract.ApplyAllowListUpdates(&_TokenPool.TransactOpts, removes, adds)
}

func (_TokenPool *TokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _TokenPool.Contract.ApplyAllowListUpdates(&_TokenPool.TransactOpts, removes, adds)
}

func (_TokenPool *TokenPoolTransactor) ApplyRampUpdates(opts *bind.TransactOpts, onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "applyRampUpdates", onRamps, offRamps)
}

func (_TokenPool *TokenPoolSession) ApplyRampUpdates(onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _TokenPool.Contract.ApplyRampUpdates(&_TokenPool.TransactOpts, onRamps, offRamps)
}

func (_TokenPool *TokenPoolTransactorSession) ApplyRampUpdates(onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error) {
	return _TokenPool.Contract.ApplyRampUpdates(&_TokenPool.TransactOpts, onRamps, offRamps)
}

func (_TokenPool *TokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, originalSender common.Address, receiver []byte, amount *big.Int, destChainSelector uint64, extraArgs []byte) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "lockOrBurn", originalSender, receiver, amount, destChainSelector, extraArgs)
}

func (_TokenPool *TokenPoolSession) LockOrBurn(originalSender common.Address, receiver []byte, amount *big.Int, destChainSelector uint64, extraArgs []byte) (*types.Transaction, error) {
	return _TokenPool.Contract.LockOrBurn(&_TokenPool.TransactOpts, originalSender, receiver, amount, destChainSelector, extraArgs)
}

func (_TokenPool *TokenPoolTransactorSession) LockOrBurn(originalSender common.Address, receiver []byte, amount *big.Int, destChainSelector uint64, extraArgs []byte) (*types.Transaction, error) {
	return _TokenPool.Contract.LockOrBurn(&_TokenPool.TransactOpts, originalSender, receiver, amount, destChainSelector, extraArgs)
}

func (_TokenPool *TokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, originalSender []byte, receiver common.Address, amount *big.Int, sourceChainSelector uint64, extraData []byte) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "releaseOrMint", originalSender, receiver, amount, sourceChainSelector, extraData)
}

func (_TokenPool *TokenPoolSession) ReleaseOrMint(originalSender []byte, receiver common.Address, amount *big.Int, sourceChainSelector uint64, extraData []byte) (*types.Transaction, error) {
	return _TokenPool.Contract.ReleaseOrMint(&_TokenPool.TransactOpts, originalSender, receiver, amount, sourceChainSelector, extraData)
}

func (_TokenPool *TokenPoolTransactorSession) ReleaseOrMint(originalSender []byte, receiver common.Address, amount *big.Int, sourceChainSelector uint64, extraData []byte) (*types.Transaction, error) {
	return _TokenPool.Contract.ReleaseOrMint(&_TokenPool.TransactOpts, originalSender, receiver, amount, sourceChainSelector, extraData)
}

func (_TokenPool *TokenPoolTransactor) SetOffRampRateLimiterConfig(opts *bind.TransactOpts, offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "setOffRampRateLimiterConfig", offRamp, config)
}

func (_TokenPool *TokenPoolSession) SetOffRampRateLimiterConfig(offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.Contract.SetOffRampRateLimiterConfig(&_TokenPool.TransactOpts, offRamp, config)
}

func (_TokenPool *TokenPoolTransactorSession) SetOffRampRateLimiterConfig(offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.Contract.SetOffRampRateLimiterConfig(&_TokenPool.TransactOpts, offRamp, config)
}

func (_TokenPool *TokenPoolTransactor) SetOnRampRateLimiterConfig(opts *bind.TransactOpts, onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "setOnRampRateLimiterConfig", onRamp, config)
}

func (_TokenPool *TokenPoolSession) SetOnRampRateLimiterConfig(onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.Contract.SetOnRampRateLimiterConfig(&_TokenPool.TransactOpts, onRamp, config)
}

func (_TokenPool *TokenPoolTransactorSession) SetOnRampRateLimiterConfig(onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error) {
	return _TokenPool.Contract.SetOnRampRateLimiterConfig(&_TokenPool.TransactOpts, onRamp, config)
}

func (_TokenPool *TokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _TokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_TokenPool *TokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenPool.Contract.TransferOwnership(&_TokenPool.TransactOpts, to)
}

func (_TokenPool *TokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenPool.Contract.TransferOwnership(&_TokenPool.TransactOpts, to)
}

type TokenPoolAllowListAddIterator struct {
	Event *TokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolAllowListAdd)
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
		it.Event = new(TokenPoolAllowListAdd)
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

func (it *TokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *TokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*TokenPoolAllowListAddIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &TokenPoolAllowListAddIterator{contract: _TokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *TokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolAllowListAdd)
				if err := _TokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseAllowListAdd(log types.Log) (*TokenPoolAllowListAdd, error) {
	event := new(TokenPoolAllowListAdd)
	if err := _TokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolAllowListRemoveIterator struct {
	Event *TokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolAllowListRemove)
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
		it.Event = new(TokenPoolAllowListRemove)
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

func (it *TokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *TokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*TokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &TokenPoolAllowListRemoveIterator{contract: _TokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *TokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolAllowListRemove)
				if err := _TokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseAllowListRemove(log types.Log) (*TokenPoolAllowListRemove, error) {
	event := new(TokenPoolAllowListRemove)
	if err := _TokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolBurnedIterator struct {
	Event *TokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolBurned)
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
		it.Event = new(TokenPoolBurned)
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

func (it *TokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*TokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolBurnedIterator{contract: _TokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *TokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolBurned)
				if err := _TokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseBurned(log types.Log) (*TokenPoolBurned, error) {
	event := new(TokenPoolBurned)
	if err := _TokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolLockedIterator struct {
	Event *TokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolLocked)
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
		it.Event = new(TokenPoolLocked)
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

func (it *TokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*TokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolLockedIterator{contract: _TokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *TokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolLocked)
				if err := _TokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseLocked(log types.Log) (*TokenPoolLocked, error) {
	event := new(TokenPoolLocked)
	if err := _TokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolMintedIterator struct {
	Event *TokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolMinted)
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
		it.Event = new(TokenPoolMinted)
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

func (it *TokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*TokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolMintedIterator{contract: _TokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *TokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolMinted)
				if err := _TokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseMinted(log types.Log) (*TokenPoolMinted, error) {
	event := new(TokenPoolMinted)
	if err := _TokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOffRampAddedIterator struct {
	Event *TokenPoolOffRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOffRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOffRampAdded)
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
		it.Event = new(TokenPoolOffRampAdded)
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

func (it *TokenPoolOffRampAddedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOffRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOffRampAdded struct {
	OffRamp           common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOffRampAdded(opts *bind.FilterOpts) (*TokenPoolOffRampAddedIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OffRampAdded")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOffRampAddedIterator{contract: _TokenPool.contract, event: "OffRampAdded", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampAdded) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OffRampAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOffRampAdded)
				if err := _TokenPool.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOffRampAdded(log types.Log) (*TokenPoolOffRampAdded, error) {
	event := new(TokenPoolOffRampAdded)
	if err := _TokenPool.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOffRampConfiguredIterator struct {
	Event *TokenPoolOffRampConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOffRampConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOffRampConfigured)
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
		it.Event = new(TokenPoolOffRampConfigured)
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

func (it *TokenPoolOffRampConfiguredIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOffRampConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOffRampConfigured struct {
	OffRamp           common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOffRampConfigured(opts *bind.FilterOpts) (*TokenPoolOffRampConfiguredIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OffRampConfigured")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOffRampConfiguredIterator{contract: _TokenPool.contract, event: "OffRampConfigured", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOffRampConfigured(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampConfigured) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OffRampConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOffRampConfigured)
				if err := _TokenPool.contract.UnpackLog(event, "OffRampConfigured", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOffRampConfigured(log types.Log) (*TokenPoolOffRampConfigured, error) {
	event := new(TokenPoolOffRampConfigured)
	if err := _TokenPool.contract.UnpackLog(event, "OffRampConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOffRampRemovedIterator struct {
	Event *TokenPoolOffRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOffRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOffRampRemoved)
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
		it.Event = new(TokenPoolOffRampRemoved)
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

func (it *TokenPoolOffRampRemovedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOffRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOffRampRemoved struct {
	OffRamp common.Address
	Raw     types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOffRampRemoved(opts *bind.FilterOpts) (*TokenPoolOffRampRemovedIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OffRampRemoved")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOffRampRemovedIterator{contract: _TokenPool.contract, event: "OffRampRemoved", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampRemoved) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OffRampRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOffRampRemoved)
				if err := _TokenPool.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOffRampRemoved(log types.Log) (*TokenPoolOffRampRemoved, error) {
	event := new(TokenPoolOffRampRemoved)
	if err := _TokenPool.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOnRampAddedIterator struct {
	Event *TokenPoolOnRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOnRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOnRampAdded)
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
		it.Event = new(TokenPoolOnRampAdded)
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

func (it *TokenPoolOnRampAddedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOnRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOnRampAdded struct {
	OnRamp            common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOnRampAdded(opts *bind.FilterOpts) (*TokenPoolOnRampAddedIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OnRampAdded")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOnRampAddedIterator{contract: _TokenPool.contract, event: "OnRampAdded", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOnRampAdded(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampAdded) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OnRampAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOnRampAdded)
				if err := _TokenPool.contract.UnpackLog(event, "OnRampAdded", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOnRampAdded(log types.Log) (*TokenPoolOnRampAdded, error) {
	event := new(TokenPoolOnRampAdded)
	if err := _TokenPool.contract.UnpackLog(event, "OnRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOnRampConfiguredIterator struct {
	Event *TokenPoolOnRampConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOnRampConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOnRampConfigured)
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
		it.Event = new(TokenPoolOnRampConfigured)
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

func (it *TokenPoolOnRampConfiguredIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOnRampConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOnRampConfigured struct {
	OnRamp            common.Address
	RateLimiterConfig RateLimiterConfig
	Raw               types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOnRampConfigured(opts *bind.FilterOpts) (*TokenPoolOnRampConfiguredIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OnRampConfigured")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOnRampConfiguredIterator{contract: _TokenPool.contract, event: "OnRampConfigured", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOnRampConfigured(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampConfigured) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OnRampConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOnRampConfigured)
				if err := _TokenPool.contract.UnpackLog(event, "OnRampConfigured", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOnRampConfigured(log types.Log) (*TokenPoolOnRampConfigured, error) {
	event := new(TokenPoolOnRampConfigured)
	if err := _TokenPool.contract.UnpackLog(event, "OnRampConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOnRampRemovedIterator struct {
	Event *TokenPoolOnRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOnRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOnRampRemoved)
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
		it.Event = new(TokenPoolOnRampRemoved)
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

func (it *TokenPoolOnRampRemovedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOnRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOnRampRemoved struct {
	OnRamp common.Address
	Raw    types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOnRampRemoved(opts *bind.FilterOpts) (*TokenPoolOnRampRemovedIterator, error) {

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OnRampRemoved")
	if err != nil {
		return nil, err
	}
	return &TokenPoolOnRampRemovedIterator{contract: _TokenPool.contract, event: "OnRampRemoved", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOnRampRemoved(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampRemoved) (event.Subscription, error) {

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OnRampRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOnRampRemoved)
				if err := _TokenPool.contract.UnpackLog(event, "OnRampRemoved", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOnRampRemoved(log types.Log) (*TokenPoolOnRampRemoved, error) {
	event := new(TokenPoolOnRampRemoved)
	if err := _TokenPool.contract.UnpackLog(event, "OnRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOwnershipTransferRequestedIterator struct {
	Event *TokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOwnershipTransferRequested)
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
		it.Event = new(TokenPoolOwnershipTransferRequested)
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

func (it *TokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolOwnershipTransferRequestedIterator{contract: _TokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOwnershipTransferRequested)
				if err := _TokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*TokenPoolOwnershipTransferRequested, error) {
	event := new(TokenPoolOwnershipTransferRequested)
	if err := _TokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolOwnershipTransferredIterator struct {
	Event *TokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolOwnershipTransferred)
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
		it.Event = new(TokenPoolOwnershipTransferred)
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

func (it *TokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolOwnershipTransferredIterator{contract: _TokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolOwnershipTransferred)
				if err := _TokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*TokenPoolOwnershipTransferred, error) {
	event := new(TokenPoolOwnershipTransferred)
	if err := _TokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenPoolReleasedIterator struct {
	Event *TokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPoolReleased)
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
		it.Event = new(TokenPoolReleased)
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

func (it *TokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *TokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_TokenPool *TokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*TokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _TokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &TokenPoolReleasedIterator{contract: _TokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_TokenPool *TokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *TokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _TokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenPoolReleased)
				if err := _TokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_TokenPool *TokenPoolFilterer) ParseReleased(log types.Log) (*TokenPoolReleased, error) {
	event := new(TokenPoolReleased)
	if err := _TokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TokenPool *TokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TokenPool.abi.Events["AllowListAdd"].ID:
		return _TokenPool.ParseAllowListAdd(log)
	case _TokenPool.abi.Events["AllowListRemove"].ID:
		return _TokenPool.ParseAllowListRemove(log)
	case _TokenPool.abi.Events["Burned"].ID:
		return _TokenPool.ParseBurned(log)
	case _TokenPool.abi.Events["Locked"].ID:
		return _TokenPool.ParseLocked(log)
	case _TokenPool.abi.Events["Minted"].ID:
		return _TokenPool.ParseMinted(log)
	case _TokenPool.abi.Events["OffRampAdded"].ID:
		return _TokenPool.ParseOffRampAdded(log)
	case _TokenPool.abi.Events["OffRampConfigured"].ID:
		return _TokenPool.ParseOffRampConfigured(log)
	case _TokenPool.abi.Events["OffRampRemoved"].ID:
		return _TokenPool.ParseOffRampRemoved(log)
	case _TokenPool.abi.Events["OnRampAdded"].ID:
		return _TokenPool.ParseOnRampAdded(log)
	case _TokenPool.abi.Events["OnRampConfigured"].ID:
		return _TokenPool.ParseOnRampConfigured(log)
	case _TokenPool.abi.Events["OnRampRemoved"].ID:
		return _TokenPool.ParseOnRampRemoved(log)
	case _TokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _TokenPool.ParseOwnershipTransferRequested(log)
	case _TokenPool.abi.Events["OwnershipTransferred"].ID:
		return _TokenPool.ParseOwnershipTransferred(log)
	case _TokenPool.abi.Events["Released"].ID:
		return _TokenPool.ParseReleased(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (TokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (TokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (TokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (TokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (TokenPoolOffRampAdded) Topic() common.Hash {
	return common.HexToHash("0x395b7374909d2b54e5796f53c898ebf41d767c86c78ea86519acf2b805852d88")
}

func (TokenPoolOffRampConfigured) Topic() common.Hash {
	return common.HexToHash("0xb3ba339cfbb8ef80d7a29ce5493051cb90e64fcfa85d7124efc1adfa4c68399f")
}

func (TokenPoolOffRampRemoved) Topic() common.Hash {
	return common.HexToHash("0xcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c")
}

func (TokenPoolOnRampAdded) Topic() common.Hash {
	return common.HexToHash("0x0b594bb0555ff7b252e0c789ccc9d8903fec294172064308727d570505cee1ac")
}

func (TokenPoolOnRampConfigured) Topic() common.Hash {
	return common.HexToHash("0x578db78e348076074dbff64a94073a83e9a65aa6766b8c75fdc89282b0e30ed6")
}

func (TokenPoolOnRampRemoved) Topic() common.Hash {
	return common.HexToHash("0x7fd064821314ad863a0714a3f1229375ace6b6427ed5544b7b2ba1c47b1b5294")
}

func (TokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (_TokenPool *TokenPool) Address() common.Address {
	return _TokenPool.address
}

type TokenPoolInterface interface {
	CurrentOffRampRateLimiterState(opts *bind.CallOpts, offRamp common.Address) (RateLimiterTokenBucket, error)

	CurrentOnRampRateLimiterState(opts *bind.CallOpts, onRamp common.Address) (RateLimiterTokenBucket, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetArmProxy(opts *bind.CallOpts) (common.Address, error)

	GetOffRamps(opts *bind.CallOpts) ([]common.Address, error)

	GetOnRamps(opts *bind.CallOpts) ([]common.Address, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)

	IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyRampUpdates(opts *bind.TransactOpts, onRamps []TokenPoolRampUpdate, offRamps []TokenPoolRampUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, originalSender common.Address, receiver []byte, amount *big.Int, destChainSelector uint64, extraArgs []byte) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, originalSender []byte, receiver common.Address, amount *big.Int, sourceChainSelector uint64, extraData []byte) (*types.Transaction, error)

	SetOffRampRateLimiterConfig(opts *bind.TransactOpts, offRamp common.Address, config RateLimiterConfig) (*types.Transaction, error)

	SetOnRampRateLimiterConfig(opts *bind.TransactOpts, onRamp common.Address, config RateLimiterConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*TokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *TokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*TokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*TokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *TokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*TokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*TokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *TokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*TokenPoolBurned, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*TokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *TokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*TokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*TokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *TokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*TokenPoolMinted, error)

	FilterOffRampAdded(opts *bind.FilterOpts) (*TokenPoolOffRampAddedIterator, error)

	WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampAdded) (event.Subscription, error)

	ParseOffRampAdded(log types.Log) (*TokenPoolOffRampAdded, error)

	FilterOffRampConfigured(opts *bind.FilterOpts) (*TokenPoolOffRampConfiguredIterator, error)

	WatchOffRampConfigured(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampConfigured) (event.Subscription, error)

	ParseOffRampConfigured(log types.Log) (*TokenPoolOffRampConfigured, error)

	FilterOffRampRemoved(opts *bind.FilterOpts) (*TokenPoolOffRampRemovedIterator, error)

	WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *TokenPoolOffRampRemoved) (event.Subscription, error)

	ParseOffRampRemoved(log types.Log) (*TokenPoolOffRampRemoved, error)

	FilterOnRampAdded(opts *bind.FilterOpts) (*TokenPoolOnRampAddedIterator, error)

	WatchOnRampAdded(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampAdded) (event.Subscription, error)

	ParseOnRampAdded(log types.Log) (*TokenPoolOnRampAdded, error)

	FilterOnRampConfigured(opts *bind.FilterOpts) (*TokenPoolOnRampConfiguredIterator, error)

	WatchOnRampConfigured(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampConfigured) (event.Subscription, error)

	ParseOnRampConfigured(log types.Log) (*TokenPoolOnRampConfigured, error)

	FilterOnRampRemoved(opts *bind.FilterOpts) (*TokenPoolOnRampRemovedIterator, error)

	WatchOnRampRemoved(opts *bind.WatchOpts, sink chan<- *TokenPoolOnRampRemoved) (event.Subscription, error)

	ParseOnRampRemoved(log types.Log) (*TokenPoolOnRampRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TokenPoolOwnershipTransferred, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*TokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *TokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*TokenPoolReleased, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
