// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token_admin_registry

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

type TokenAdminRegistryTokenConfig struct {
	Administrator        common.Address
	PendingAdministrator common.Address
	TokenPool            common.Address
}

var TokenAdminRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"AlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidTokenPoolToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"OnlyAdministrator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"OnlyPendingAdministrator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"OnlyRegistryModuleOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"currentAdmin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdministratorTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdministratorTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousPool\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"PoolSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"RegistryModuleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"RegistryModuleRemoved\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"name\":\"acceptAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"addRegistryModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxCount\",\"type\":\"uint64\"}],\"name\":\"getAllConfiguredTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getPools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pendingAdministrator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenPool\",\"type\":\"address\"}],\"internalType\":\"structTokenAdminRegistry.TokenConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"isAdministrator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"isRegistryModule\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"proposeAdministrator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"removeRegistryModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"setPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60808060405234603d573315602c57600180546001600160a01b031916331790556112c790816100438239f35b639b15e16f60e01b60005260046000fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806310cbcf1814610cc0578063156194da14610ba1578063181f5a7714610aad5780633dc4577214610a3a5780634e847fc71461086c5780635e63547a1461079257806372d64a811461073b57806379ba5097146106705780637d3f25521461061f5780638da5cb5b146105eb578063bbe4f6db1461058b578063c1af6e031461050f578063cb67e3b11461045b578063ddadfa8e14610349578063e677ae37146101a35763f2fde38b146100cc57600080fd5b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff6100fa610d2d565b610102611007565b1633811461017457807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b600080fd5b3461019e57604060031936011261019e576101bc610d2d565b6101c4610d50565b6101db336000526006602052604060002054151590565b1580610327575b6102f95773ffffffffffffffffffffffffffffffffffffffff169081156102cf5773ffffffffffffffffffffffffffffffffffffffff16806000526002602052604060002073ffffffffffffffffffffffffffffffffffffffff8154166102a15760010180547fffffffffffffffffffffffff0000000000000000000000000000000000000000168317905560009061027a81611260565b507fc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b7168280a4005b507f45ed80e90000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7fd92e233d0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f51ca1ec3000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b5073ffffffffffffffffffffffffffffffffffffffff600154163314156101e2565b3461019e57604060031936011261019e57610362610d2d565b73ffffffffffffffffffffffffffffffffffffffff61037f610d50565b91169081600052600260205273ffffffffffffffffffffffffffffffffffffffff6040600020541633036104295773ffffffffffffffffffffffffffffffffffffffff9082600052600260205260016040600020018282167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055169033907fc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b716600080a4005b507fed5d85b5000000000000000000000000000000000000000000000000000000006000523360045260245260446000fd5b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff610489610d2d565b60006040805161049881610d73565b82815282602082015201521660005260026020526060604060002073ffffffffffffffffffffffffffffffffffffffff6040516104d481610d73565b818084541693848352604082600281600185015416936020870194855201541693019283526040519485525116602084015251166040820152f35b3461019e57604060031936011261019e57602061052a610d2d565b73ffffffffffffffffffffffffffffffffffffffff610547610d50565b91166000526002825273ffffffffffffffffffffffffffffffffffffffff6040600020541673ffffffffffffffffffffffffffffffffffffffff6040519216148152f35b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff6105b9610d2d565b166000526002602052602073ffffffffffffffffffffffffffffffffffffffff60026040600020015416604051908152f35b3461019e57600060031936011261019e57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b3461019e57602060031936011261019e57602061066673ffffffffffffffffffffffffffffffffffffffff610652610d2d565b166000526006602052604060002054151590565b6040519015158152f35b3461019e57600060031936011261019e5760005473ffffffffffffffffffffffffffffffffffffffff81163303610711577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b3461019e57604060031936011261019e5760043567ffffffffffffffff8116810361019e576024359067ffffffffffffffff8216820361019e5761078e9161078291610f06565b60405191829182610dd0565b0390f35b3461019e57602060031936011261019e5760043567ffffffffffffffff811161019e573660238201121561019e57806004013567ffffffffffffffff811161019e573660248260051b8401011161019e576107ec81610e38565b60005b8281101561085e57600060248260051b8601013573ffffffffffffffffffffffffffffffffffffffff811680910361085a57600260408373ffffffffffffffffffffffffffffffffffffffff93600196955282602052200154166108538285610e87565b52016107ef565b5080fd5b6040518061078e8482610dd0565b3461019e57604060031936011261019e57610885610d2d565b73ffffffffffffffffffffffffffffffffffffffff6108a2610d50565b91169081600052600260205273ffffffffffffffffffffffffffffffffffffffff6040600020541633036104295773ffffffffffffffffffffffffffffffffffffffff1690811515806109a4575b6109775780600052600260205260026040600020019073ffffffffffffffffffffffffffffffffffffffff82541691837fffffffffffffffffffffffff000000000000000000000000000000000000000082541617905582820361095057005b7f754449ec3aff3bd528bfce43ae9319c4a381b67fcd1d20097b3b24dacaecc35d600080a4005b7f962b60e60000000000000000000000000000000000000000000000000000000060005260045260246000fd5b506040517f240028e8000000000000000000000000000000000000000000000000000000008152816004820152602081602481865afa908115610a2e576000916109f0575b50156108f0565b6020813d602011610a26575b81610a0960209383610d8f565b8101031261085a5751908115158203610a235750836109e9565b80fd5b3d91506109fc565b6040513d6000823e3d90fd5b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff610a68610d2d565b610a70611007565b16610a7a81611200565b610a8057005b60207f3cabf004338366bfeaeb610ad827cb58d16b588017c509501f2c97c83caae7b291604051908152a1005b3461019e57600060031936011261019e576040516040810181811067ffffffffffffffff821117610b7257604052601881527f546f6b656e41646d696e526567697374727920312e352e300000000000000000602082015260405190602082528181519182602083015260005b838110610b5a5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f836000604080968601015201168101030190f35b60208282018101516040878401015285935001610b1a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff610bcf610d2d565b168060005260026020526040600020600181019073ffffffffffffffffffffffffffffffffffffffff8254163303610c8e5773ffffffffffffffffffffffffffffffffffffffff33167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557fffffffffffffffffffffffff0000000000000000000000000000000000000000815416905533907f399b55200f7f639a63d76efe3dcfa9156ce367058d6b673041b84a628885f5a7600080a3005b827f3edffe75000000000000000000000000000000000000000000000000000000006000523360045260245260446000fd5b3461019e57602060031936011261019e5773ffffffffffffffffffffffffffffffffffffffff610cee610d2d565b610cf6611007565b16610d008161106a565b610d0657005b7f93eaa26dcb9275e56bacb1d33fdbf402262da6f0f4baf2a6e2cd154b73f387f8600080a2005b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361019e57565b6024359073ffffffffffffffffffffffffffffffffffffffff8216820361019e57565b6060810190811067ffffffffffffffff821117610b7257604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610b7257604052565b602060408183019282815284518094520192019060005b818110610df45750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101610de7565b67ffffffffffffffff8111610b725760051b60200190565b90610e4282610e20565b610e4f6040519182610d8f565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0610e7d8294610e20565b0190602036910137565b8051821015610e9b5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b91908201809211610ed757565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff60035491169181831015610fff5767ffffffffffffffff169080610f338484610eca565b11610fee575b5090610f4482610e38565b91600091600354925b828110610f5b575050505090565b610f658183610eca565b600085821015610fc157600390527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b01546001919073ffffffffffffffffffffffffffffffffffffffff16610fba8288610e87565b5201610f4d565b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526032600452fd5b82810391508111610ed75738610f39565b505050606090565b73ffffffffffffffffffffffffffffffffffffffff60015416330361102857565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b8054821015610e9b5760005260206000200190600090565b60008181526006602052604090205480156111f9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111610ed757600554907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211610ed75781810361118a575b505050600554801561115b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01611118816005611052565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600555600052600660205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6111e161119b6111ac936005611052565b90549060031b1c9283926005611052565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905560005260066020526040600020553880806110df565b5050600090565b8060005260066020526040600020541560001461125a5760055468010000000000000000811015610b72576112416111ac8260018594016005556005611052565b9055600554906000526006602052604060002055600190565b50600090565b8060005260046020526040600020541560001461125a5760035468010000000000000000811015610b72576112a16111ac8260018594016003556003611052565b905560035490600052600460205260406000205560019056fea164736f6c634300081a000a",
}

var TokenAdminRegistryABI = TokenAdminRegistryMetaData.ABI

var TokenAdminRegistryBin = TokenAdminRegistryMetaData.Bin

func DeployTokenAdminRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenAdminRegistry, error) {
	parsed, err := TokenAdminRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenAdminRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenAdminRegistry{address: address, abi: *parsed, TokenAdminRegistryCaller: TokenAdminRegistryCaller{contract: contract}, TokenAdminRegistryTransactor: TokenAdminRegistryTransactor{contract: contract}, TokenAdminRegistryFilterer: TokenAdminRegistryFilterer{contract: contract}}, nil
}

type TokenAdminRegistry struct {
	address common.Address
	abi     abi.ABI
	TokenAdminRegistryCaller
	TokenAdminRegistryTransactor
	TokenAdminRegistryFilterer
}

type TokenAdminRegistryCaller struct {
	contract *bind.BoundContract
}

type TokenAdminRegistryTransactor struct {
	contract *bind.BoundContract
}

type TokenAdminRegistryFilterer struct {
	contract *bind.BoundContract
}

type TokenAdminRegistrySession struct {
	Contract     *TokenAdminRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TokenAdminRegistryCallerSession struct {
	Contract *TokenAdminRegistryCaller
	CallOpts bind.CallOpts
}

type TokenAdminRegistryTransactorSession struct {
	Contract     *TokenAdminRegistryTransactor
	TransactOpts bind.TransactOpts
}

type TokenAdminRegistryRaw struct {
	Contract *TokenAdminRegistry
}

type TokenAdminRegistryCallerRaw struct {
	Contract *TokenAdminRegistryCaller
}

type TokenAdminRegistryTransactorRaw struct {
	Contract *TokenAdminRegistryTransactor
}

func NewTokenAdminRegistry(address common.Address, backend bind.ContractBackend) (*TokenAdminRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(TokenAdminRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTokenAdminRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistry{address: address, abi: abi, TokenAdminRegistryCaller: TokenAdminRegistryCaller{contract: contract}, TokenAdminRegistryTransactor: TokenAdminRegistryTransactor{contract: contract}, TokenAdminRegistryFilterer: TokenAdminRegistryFilterer{contract: contract}}, nil
}

func NewTokenAdminRegistryCaller(address common.Address, caller bind.ContractCaller) (*TokenAdminRegistryCaller, error) {
	contract, err := bindTokenAdminRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryCaller{contract: contract}, nil
}

func NewTokenAdminRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenAdminRegistryTransactor, error) {
	contract, err := bindTokenAdminRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryTransactor{contract: contract}, nil
}

func NewTokenAdminRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenAdminRegistryFilterer, error) {
	contract, err := bindTokenAdminRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryFilterer{contract: contract}, nil
}

func bindTokenAdminRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenAdminRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryTransactor.contract.Transfer(opts)
}

func (_TokenAdminRegistry *TokenAdminRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TokenAdminRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenAdminRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.contract.Transfer(opts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getAllConfiguredTokens", startIndex, maxCount)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetAllConfiguredTokens(startIndex uint64, maxCount uint64) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetAllConfiguredTokens(&_TokenAdminRegistry.CallOpts, startIndex, maxCount)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetAllConfiguredTokens(startIndex uint64, maxCount uint64) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetAllConfiguredTokens(&_TokenAdminRegistry.CallOpts, startIndex, maxCount)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetPool(opts *bind.CallOpts, token common.Address) (common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getPool", token)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetPool(token common.Address) (common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPool(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetPool(token common.Address) (common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPool(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetPools(opts *bind.CallOpts, tokens []common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getPools", tokens)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetPools(tokens []common.Address) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPools(&_TokenAdminRegistry.CallOpts, tokens)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetPools(tokens []common.Address) ([]common.Address, error) {
	return _TokenAdminRegistry.Contract.GetPools(&_TokenAdminRegistry.CallOpts, tokens)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) GetTokenConfig(opts *bind.CallOpts, token common.Address) (TokenAdminRegistryTokenConfig, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "getTokenConfig", token)

	if err != nil {
		return *new(TokenAdminRegistryTokenConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(TokenAdminRegistryTokenConfig)).(*TokenAdminRegistryTokenConfig)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) GetTokenConfig(token common.Address) (TokenAdminRegistryTokenConfig, error) {
	return _TokenAdminRegistry.Contract.GetTokenConfig(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) GetTokenConfig(token common.Address) (TokenAdminRegistryTokenConfig, error) {
	return _TokenAdminRegistry.Contract.GetTokenConfig(&_TokenAdminRegistry.CallOpts, token)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) IsAdministrator(opts *bind.CallOpts, localToken common.Address, administrator common.Address) (bool, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "isAdministrator", localToken, administrator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) IsAdministrator(localToken common.Address, administrator common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsAdministrator(&_TokenAdminRegistry.CallOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) IsAdministrator(localToken common.Address, administrator common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsAdministrator(&_TokenAdminRegistry.CallOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) IsRegistryModule(opts *bind.CallOpts, module common.Address) (bool, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "isRegistryModule", module)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) IsRegistryModule(module common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsRegistryModule(&_TokenAdminRegistry.CallOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) IsRegistryModule(module common.Address) (bool, error) {
	return _TokenAdminRegistry.Contract.IsRegistryModule(&_TokenAdminRegistry.CallOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) Owner() (common.Address, error) {
	return _TokenAdminRegistry.Contract.Owner(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) Owner() (common.Address, error) {
	return _TokenAdminRegistry.Contract.Owner(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenAdminRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TypeAndVersion() (string, error) {
	return _TokenAdminRegistry.Contract.TypeAndVersion(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryCallerSession) TypeAndVersion() (string, error) {
	return _TokenAdminRegistry.Contract.TypeAndVersion(&_TokenAdminRegistry.CallOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AcceptAdminRole(opts *bind.TransactOpts, localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "acceptAdminRole", localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AcceptAdminRole(localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptAdminRole(&_TokenAdminRegistry.TransactOpts, localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AcceptAdminRole(localToken common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptAdminRole(&_TokenAdminRegistry.TransactOpts, localToken)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptOwnership(&_TokenAdminRegistry.TransactOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AcceptOwnership(&_TokenAdminRegistry.TransactOpts)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) AddRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "addRegistryModule", module)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) AddRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AddRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) AddRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.AddRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) ProposeAdministrator(opts *bind.TransactOpts, localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "proposeAdministrator", localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) ProposeAdministrator(localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.ProposeAdministrator(&_TokenAdminRegistry.TransactOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) ProposeAdministrator(localToken common.Address, administrator common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.ProposeAdministrator(&_TokenAdminRegistry.TransactOpts, localToken, administrator)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) RemoveRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "removeRegistryModule", module)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) RemoveRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.RemoveRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) RemoveRegistryModule(module common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.RemoveRegistryModule(&_TokenAdminRegistry.TransactOpts, module)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) SetPool(opts *bind.TransactOpts, localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "setPool", localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) SetPool(localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.SetPool(&_TokenAdminRegistry.TransactOpts, localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) SetPool(localToken common.Address, pool common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.SetPool(&_TokenAdminRegistry.TransactOpts, localToken, pool)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) TransferAdminRole(opts *bind.TransactOpts, localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "transferAdminRole", localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TransferAdminRole(localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferAdminRole(&_TokenAdminRegistry.TransactOpts, localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) TransferAdminRole(localToken common.Address, newAdmin common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferAdminRole(&_TokenAdminRegistry.TransactOpts, localToken, newAdmin)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_TokenAdminRegistry *TokenAdminRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferOwnership(&_TokenAdminRegistry.TransactOpts, to)
}

func (_TokenAdminRegistry *TokenAdminRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TokenAdminRegistry.Contract.TransferOwnership(&_TokenAdminRegistry.TransactOpts, to)
}

type TokenAdminRegistryAdministratorTransferRequestedIterator struct {
	Event *TokenAdminRegistryAdministratorTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryAdministratorTransferRequested)
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
		it.Event = new(TokenAdminRegistryAdministratorTransferRequested)
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

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryAdministratorTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryAdministratorTransferRequested struct {
	Token        common.Address
	CurrentAdmin common.Address
	NewAdmin     common.Address
	Raw          types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterAdministratorTransferRequested(opts *bind.FilterOpts, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferRequestedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var currentAdminRule []interface{}
	for _, currentAdminItem := range currentAdmin {
		currentAdminRule = append(currentAdminRule, currentAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "AdministratorTransferRequested", tokenRule, currentAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryAdministratorTransferRequestedIterator{contract: _TokenAdminRegistry.contract, event: "AdministratorTransferRequested", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchAdministratorTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferRequested, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var currentAdminRule []interface{}
	for _, currentAdminItem := range currentAdmin {
		currentAdminRule = append(currentAdminRule, currentAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "AdministratorTransferRequested", tokenRule, currentAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryAdministratorTransferRequested)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferRequested", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseAdministratorTransferRequested(log types.Log) (*TokenAdminRegistryAdministratorTransferRequested, error) {
	event := new(TokenAdminRegistryAdministratorTransferRequested)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryAdministratorTransferredIterator struct {
	Event *TokenAdminRegistryAdministratorTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryAdministratorTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryAdministratorTransferred)
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
		it.Event = new(TokenAdminRegistryAdministratorTransferred)
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

func (it *TokenAdminRegistryAdministratorTransferredIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryAdministratorTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryAdministratorTransferred struct {
	Token    common.Address
	NewAdmin common.Address
	Raw      types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterAdministratorTransferred(opts *bind.FilterOpts, token []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "AdministratorTransferred", tokenRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryAdministratorTransferredIterator{contract: _TokenAdminRegistry.contract, event: "AdministratorTransferred", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchAdministratorTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferred, token []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "AdministratorTransferred", tokenRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryAdministratorTransferred)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferred", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseAdministratorTransferred(log types.Log) (*TokenAdminRegistryAdministratorTransferred, error) {
	event := new(TokenAdminRegistryAdministratorTransferred)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "AdministratorTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryOwnershipTransferRequestedIterator struct {
	Event *TokenAdminRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryOwnershipTransferRequested)
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
		it.Event = new(TokenAdminRegistryOwnershipTransferRequested)
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

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryOwnershipTransferRequestedIterator{contract: _TokenAdminRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryOwnershipTransferRequested)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*TokenAdminRegistryOwnershipTransferRequested, error) {
	event := new(TokenAdminRegistryOwnershipTransferRequested)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryOwnershipTransferredIterator struct {
	Event *TokenAdminRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryOwnershipTransferred)
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
		it.Event = new(TokenAdminRegistryOwnershipTransferred)
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

func (it *TokenAdminRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryOwnershipTransferredIterator{contract: _TokenAdminRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryOwnershipTransferred)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*TokenAdminRegistryOwnershipTransferred, error) {
	event := new(TokenAdminRegistryOwnershipTransferred)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryPoolSetIterator struct {
	Event *TokenAdminRegistryPoolSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryPoolSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryPoolSet)
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
		it.Event = new(TokenAdminRegistryPoolSet)
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

func (it *TokenAdminRegistryPoolSetIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryPoolSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryPoolSet struct {
	Token        common.Address
	PreviousPool common.Address
	NewPool      common.Address
	Raw          types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterPoolSet(opts *bind.FilterOpts, token []common.Address, previousPool []common.Address, newPool []common.Address) (*TokenAdminRegistryPoolSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var previousPoolRule []interface{}
	for _, previousPoolItem := range previousPool {
		previousPoolRule = append(previousPoolRule, previousPoolItem)
	}
	var newPoolRule []interface{}
	for _, newPoolItem := range newPool {
		newPoolRule = append(newPoolRule, newPoolItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "PoolSet", tokenRule, previousPoolRule, newPoolRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryPoolSetIterator{contract: _TokenAdminRegistry.contract, event: "PoolSet", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchPoolSet(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryPoolSet, token []common.Address, previousPool []common.Address, newPool []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var previousPoolRule []interface{}
	for _, previousPoolItem := range previousPool {
		previousPoolRule = append(previousPoolRule, previousPoolItem)
	}
	var newPoolRule []interface{}
	for _, newPoolItem := range newPool {
		newPoolRule = append(newPoolRule, newPoolItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "PoolSet", tokenRule, previousPoolRule, newPoolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryPoolSet)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "PoolSet", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParsePoolSet(log types.Log) (*TokenAdminRegistryPoolSet, error) {
	event := new(TokenAdminRegistryPoolSet)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "PoolSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryRegistryModuleAddedIterator struct {
	Event *TokenAdminRegistryRegistryModuleAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryRegistryModuleAdded)
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
		it.Event = new(TokenAdminRegistryRegistryModuleAdded)
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

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryRegistryModuleAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryRegistryModuleAdded struct {
	Module common.Address
	Raw    types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterRegistryModuleAdded(opts *bind.FilterOpts) (*TokenAdminRegistryRegistryModuleAddedIterator, error) {

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "RegistryModuleAdded")
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryRegistryModuleAddedIterator{contract: _TokenAdminRegistry.contract, event: "RegistryModuleAdded", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchRegistryModuleAdded(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleAdded) (event.Subscription, error) {

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "RegistryModuleAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryRegistryModuleAdded)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleAdded", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseRegistryModuleAdded(log types.Log) (*TokenAdminRegistryRegistryModuleAdded, error) {
	event := new(TokenAdminRegistryRegistryModuleAdded)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TokenAdminRegistryRegistryModuleRemovedIterator struct {
	Event *TokenAdminRegistryRegistryModuleRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAdminRegistryRegistryModuleRemoved)
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
		it.Event = new(TokenAdminRegistryRegistryModuleRemoved)
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

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Error() error {
	return it.fail
}

func (it *TokenAdminRegistryRegistryModuleRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TokenAdminRegistryRegistryModuleRemoved struct {
	Module common.Address
	Raw    types.Log
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) FilterRegistryModuleRemoved(opts *bind.FilterOpts, module []common.Address) (*TokenAdminRegistryRegistryModuleRemovedIterator, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.FilterLogs(opts, "RegistryModuleRemoved", moduleRule)
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistryRegistryModuleRemovedIterator{contract: _TokenAdminRegistry.contract, event: "RegistryModuleRemoved", logs: logs, sub: sub}, nil
}

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) WatchRegistryModuleRemoved(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleRemoved, module []common.Address) (event.Subscription, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _TokenAdminRegistry.contract.WatchLogs(opts, "RegistryModuleRemoved", moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TokenAdminRegistryRegistryModuleRemoved)
				if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleRemoved", log); err != nil {
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

func (_TokenAdminRegistry *TokenAdminRegistryFilterer) ParseRegistryModuleRemoved(log types.Log) (*TokenAdminRegistryRegistryModuleRemoved, error) {
	event := new(TokenAdminRegistryRegistryModuleRemoved)
	if err := _TokenAdminRegistry.contract.UnpackLog(event, "RegistryModuleRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TokenAdminRegistry *TokenAdminRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TokenAdminRegistry.abi.Events["AdministratorTransferRequested"].ID:
		return _TokenAdminRegistry.ParseAdministratorTransferRequested(log)
	case _TokenAdminRegistry.abi.Events["AdministratorTransferred"].ID:
		return _TokenAdminRegistry.ParseAdministratorTransferred(log)
	case _TokenAdminRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _TokenAdminRegistry.ParseOwnershipTransferRequested(log)
	case _TokenAdminRegistry.abi.Events["OwnershipTransferred"].ID:
		return _TokenAdminRegistry.ParseOwnershipTransferred(log)
	case _TokenAdminRegistry.abi.Events["PoolSet"].ID:
		return _TokenAdminRegistry.ParsePoolSet(log)
	case _TokenAdminRegistry.abi.Events["RegistryModuleAdded"].ID:
		return _TokenAdminRegistry.ParseRegistryModuleAdded(log)
	case _TokenAdminRegistry.abi.Events["RegistryModuleRemoved"].ID:
		return _TokenAdminRegistry.ParseRegistryModuleRemoved(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TokenAdminRegistryAdministratorTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xc54c3051ff16e63bb9203214432372aca006c589e3653619b577a3265675b716")
}

func (TokenAdminRegistryAdministratorTransferred) Topic() common.Hash {
	return common.HexToHash("0x399b55200f7f639a63d76efe3dcfa9156ce367058d6b673041b84a628885f5a7")
}

func (TokenAdminRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TokenAdminRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TokenAdminRegistryPoolSet) Topic() common.Hash {
	return common.HexToHash("0x754449ec3aff3bd528bfce43ae9319c4a381b67fcd1d20097b3b24dacaecc35d")
}

func (TokenAdminRegistryRegistryModuleAdded) Topic() common.Hash {
	return common.HexToHash("0x3cabf004338366bfeaeb610ad827cb58d16b588017c509501f2c97c83caae7b2")
}

func (TokenAdminRegistryRegistryModuleRemoved) Topic() common.Hash {
	return common.HexToHash("0x93eaa26dcb9275e56bacb1d33fdbf402262da6f0f4baf2a6e2cd154b73f387f8")
}

func (_TokenAdminRegistry *TokenAdminRegistry) Address() common.Address {
	return _TokenAdminRegistry.address
}

type TokenAdminRegistryInterface interface {
	GetAllConfiguredTokens(opts *bind.CallOpts, startIndex uint64, maxCount uint64) ([]common.Address, error)

	GetPool(opts *bind.CallOpts, token common.Address) (common.Address, error)

	GetPools(opts *bind.CallOpts, tokens []common.Address) ([]common.Address, error)

	GetTokenConfig(opts *bind.CallOpts, token common.Address) (TokenAdminRegistryTokenConfig, error)

	IsAdministrator(opts *bind.CallOpts, localToken common.Address, administrator common.Address) (bool, error)

	IsRegistryModule(opts *bind.CallOpts, module common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptAdminRole(opts *bind.TransactOpts, localToken common.Address) (*types.Transaction, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error)

	ProposeAdministrator(opts *bind.TransactOpts, localToken common.Address, administrator common.Address) (*types.Transaction, error)

	RemoveRegistryModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error)

	SetPool(opts *bind.TransactOpts, localToken common.Address, pool common.Address) (*types.Transaction, error)

	TransferAdminRole(opts *bind.TransactOpts, localToken common.Address, newAdmin common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAdministratorTransferRequested(opts *bind.FilterOpts, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferRequestedIterator, error)

	WatchAdministratorTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferRequested, token []common.Address, currentAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseAdministratorTransferRequested(log types.Log) (*TokenAdminRegistryAdministratorTransferRequested, error)

	FilterAdministratorTransferred(opts *bind.FilterOpts, token []common.Address, newAdmin []common.Address) (*TokenAdminRegistryAdministratorTransferredIterator, error)

	WatchAdministratorTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryAdministratorTransferred, token []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseAdministratorTransferred(log types.Log) (*TokenAdminRegistryAdministratorTransferred, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TokenAdminRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenAdminRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TokenAdminRegistryOwnershipTransferred, error)

	FilterPoolSet(opts *bind.FilterOpts, token []common.Address, previousPool []common.Address, newPool []common.Address) (*TokenAdminRegistryPoolSetIterator, error)

	WatchPoolSet(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryPoolSet, token []common.Address, previousPool []common.Address, newPool []common.Address) (event.Subscription, error)

	ParsePoolSet(log types.Log) (*TokenAdminRegistryPoolSet, error)

	FilterRegistryModuleAdded(opts *bind.FilterOpts) (*TokenAdminRegistryRegistryModuleAddedIterator, error)

	WatchRegistryModuleAdded(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleAdded) (event.Subscription, error)

	ParseRegistryModuleAdded(log types.Log) (*TokenAdminRegistryRegistryModuleAdded, error)

	FilterRegistryModuleRemoved(opts *bind.FilterOpts, module []common.Address) (*TokenAdminRegistryRegistryModuleRemovedIterator, error)

	WatchRegistryModuleRemoved(opts *bind.WatchOpts, sink chan<- *TokenAdminRegistryRegistryModuleRemoved, module []common.Address) (event.Subscription, error)

	ParseRegistryModuleRemoved(log types.Log) (*TokenAdminRegistryRegistryModuleRemoved, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
