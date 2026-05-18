// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package destination_verifier_proxy

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

var DestinationVerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"VerifierInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_accessController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"payloads\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verifyBulk\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"verifiedReports\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610fac806101576000396000f3fe6080604052600436106100b15760003560e01c80638da5cb5b11610069578063f2fde38b1161004e578063f2fde38b1461022d578063f7e83aee1461024d578063f873a61c1461026057600080fd5b80638da5cb5b146101ed57806394ba28461461021857600080fd5b806338416b5b1161009a57806338416b5b1461017c5780635437988d146101b657806379ba5097146101d857600080fd5b806301ffc9a7146100b6578063181f5a771461012d575b600080fd5b3480156100c257600080fd5b506101186100d13660046108f7565b7fffffffff00000000000000000000000000000000000000000000000000000000167ff7574762000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b34801561013957600080fd5b5060408051808201909152601e81527f44657374696e6174696f6e566572696669657250726f787920302e342e30000060208201525b60405161012491906109ae565b34801561018857600080fd5b50610191610280565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610124565b3480156101c257600080fd5b506101d66101d13660046109e3565b610319565b005b3480156101e457600080fd5b506101d6610469565b3480156101f957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610191565b34801561022457600080fd5b50610191610566565b34801561023957600080fd5b506101d66102483660046109e3565b6105d6565b61016f61025b366004610a49565b6105ea565b61027361026e366004610ab5565b6106b9565b6040516101249190610b36565b600254604080517f38416b5b000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff16916338416b5b9160048083019260209291908290030181865afa1580156102f0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103149190610bb6565b905090565b61032161077f565b6040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527f527146e200000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa1580156103ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103cf9190610bd3565b610422576040517f96ac86f300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610419565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600254604080517f94ba2846000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff16916394ba28469160048083019260209291908290030181865afa1580156102f0573d6000803e3d6000fd5b6105de61077f565b6105e781610802565b50565b6002546040517f294d2bb100000000000000000000000000000000000000000000000000000000815260609173ffffffffffffffffffffffffffffffffffffffff169063294d2bb190349061064b9089908990899089903390600401610c3e565b60006040518083038185885af1158015610669573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526106b09190810190610d96565b95945050505050565b6002546040517fd7c72e4e00000000000000000000000000000000000000000000000000000000815260609173ffffffffffffffffffffffffffffffffffffffff169063d7c72e4e90349061071a9089908990899089903390600401610dcb565b60006040518083038185885af1158015610738573d6000803e3d6000fd5b50505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526106b09190810190610edc565b60005473ffffffffffffffffffffffffffffffffffffffff163314610800576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610419565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610881576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610419565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561090957600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461093957600080fd5b9392505050565b60005b8381101561095b578181015183820152602001610943565b50506000910152565b6000815180845261097c816020860160208601610940565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006109396020830184610964565b73ffffffffffffffffffffffffffffffffffffffff811681146105e757600080fd5b6000602082840312156109f557600080fd5b8135610939816109c1565b60008083601f840112610a1257600080fd5b50813567ffffffffffffffff811115610a2a57600080fd5b602083019150836020828501011115610a4257600080fd5b9250929050565b60008060008060408587031215610a5f57600080fd5b843567ffffffffffffffff80821115610a7757600080fd5b610a8388838901610a00565b90965094506020870135915080821115610a9c57600080fd5b50610aa987828801610a00565b95989497509550505050565b60008060008060408587031215610acb57600080fd5b843567ffffffffffffffff80821115610ae357600080fd5b818701915087601f830112610af757600080fd5b813581811115610b0657600080fd5b8860208260051b8501011115610b1b57600080fd5b602092830196509450908601359080821115610a9c57600080fd5b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015610ba9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452610b97858351610964565b94509285019290850190600101610b5d565b5092979650505050505050565b600060208284031215610bc857600080fd5b8151610939816109c1565b600060208284031215610be557600080fd5b8151801515811461093957600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b606081526000610c52606083018789610bf5565b8281036020840152610c65818688610bf5565b91505073ffffffffffffffffffffffffffffffffffffffff831660408301529695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d0457610d04610c8e565b604052919050565b600082601f830112610d1d57600080fd5b815167ffffffffffffffff811115610d3757610d37610c8e565b610d6860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610cbd565b818152846020838601011115610d7d57600080fd5b610d8e826020830160208701610940565b949350505050565b600060208284031215610da857600080fd5b815167ffffffffffffffff811115610dbf57600080fd5b610d8e84828501610d0c565b6060808252810185905260006080600587901b8301810190830188835b89811015610e97577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8086850301835281357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18c3603018112610e4957600080fd5b8b01602081810191359067ffffffffffffffff821115610e6857600080fd5b813603831315610e7757600080fd5b610e82878385610bf5565b96509485019493909301925050600101610de8565b5050508281036020840152610ead818688610bf5565b915050610ed2604083018473ffffffffffffffffffffffffffffffffffffffff169052565b9695505050505050565b60006020808385031215610eef57600080fd5b825167ffffffffffffffff80821115610f0757600080fd5b818501915085601f830112610f1b57600080fd5b815181811115610f2d57610f2d610c8e565b8060051b610f3c858201610cbd565b9182528381018501918581019089841115610f5657600080fd5b86860192505b83831015610f9257825185811115610f745760008081fd5b610f828b89838a0101610d0c565b8352509186019190860190610f5c565b999850505050505050505056fea164736f6c6343000813000a",
}

var DestinationVerifierProxyABI = DestinationVerifierProxyMetaData.ABI

var DestinationVerifierProxyBin = DestinationVerifierProxyMetaData.Bin

func DeployDestinationVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DestinationVerifierProxy, error) {
	parsed, err := DestinationVerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DestinationVerifierProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DestinationVerifierProxy{address: address, abi: *parsed, DestinationVerifierProxyCaller: DestinationVerifierProxyCaller{contract: contract}, DestinationVerifierProxyTransactor: DestinationVerifierProxyTransactor{contract: contract}, DestinationVerifierProxyFilterer: DestinationVerifierProxyFilterer{contract: contract}}, nil
}

type DestinationVerifierProxy struct {
	address common.Address
	abi     abi.ABI
	DestinationVerifierProxyCaller
	DestinationVerifierProxyTransactor
	DestinationVerifierProxyFilterer
}

type DestinationVerifierProxyCaller struct {
	contract *bind.BoundContract
}

type DestinationVerifierProxyTransactor struct {
	contract *bind.BoundContract
}

type DestinationVerifierProxyFilterer struct {
	contract *bind.BoundContract
}

type DestinationVerifierProxySession struct {
	Contract     *DestinationVerifierProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DestinationVerifierProxyCallerSession struct {
	Contract *DestinationVerifierProxyCaller
	CallOpts bind.CallOpts
}

type DestinationVerifierProxyTransactorSession struct {
	Contract     *DestinationVerifierProxyTransactor
	TransactOpts bind.TransactOpts
}

type DestinationVerifierProxyRaw struct {
	Contract *DestinationVerifierProxy
}

type DestinationVerifierProxyCallerRaw struct {
	Contract *DestinationVerifierProxyCaller
}

type DestinationVerifierProxyTransactorRaw struct {
	Contract *DestinationVerifierProxyTransactor
}

func NewDestinationVerifierProxy(address common.Address, backend bind.ContractBackend) (*DestinationVerifierProxy, error) {
	abi, err := abi.JSON(strings.NewReader(DestinationVerifierProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDestinationVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxy{address: address, abi: abi, DestinationVerifierProxyCaller: DestinationVerifierProxyCaller{contract: contract}, DestinationVerifierProxyTransactor: DestinationVerifierProxyTransactor{contract: contract}, DestinationVerifierProxyFilterer: DestinationVerifierProxyFilterer{contract: contract}}, nil
}

func NewDestinationVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*DestinationVerifierProxyCaller, error) {
	contract, err := bindDestinationVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxyCaller{contract: contract}, nil
}

func NewDestinationVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*DestinationVerifierProxyTransactor, error) {
	contract, err := bindDestinationVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxyTransactor{contract: contract}, nil
}

func NewDestinationVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*DestinationVerifierProxyFilterer, error) {
	contract, err := bindDestinationVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxyFilterer{contract: contract}, nil
}

func bindDestinationVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DestinationVerifierProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DestinationVerifierProxy *DestinationVerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationVerifierProxy.Contract.DestinationVerifierProxyCaller.contract.Call(opts, result, method, params...)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.DestinationVerifierProxyTransactor.contract.Transfer(opts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.DestinationVerifierProxyTransactor.contract.Transact(opts, method, params...)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DestinationVerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.contract.Transfer(opts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.contract.Transact(opts, method, params...)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) Owner() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.Owner(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerSession) Owner() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.Owner(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCaller) SAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifierProxy.contract.Call(opts, &out, "s_accessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) SAccessController() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.SAccessController(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerSession) SAccessController() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.SAccessController(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCaller) SFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DestinationVerifierProxy.contract.Call(opts, &out, "s_feeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) SFeeManager() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.SFeeManager(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerSession) SFeeManager() (common.Address, error) {
	return _DestinationVerifierProxy.Contract.SFeeManager(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DestinationVerifierProxy.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationVerifierProxy.Contract.SupportsInterface(&_DestinationVerifierProxy.CallOpts, interfaceId)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DestinationVerifierProxy.Contract.SupportsInterface(&_DestinationVerifierProxy.CallOpts, interfaceId)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DestinationVerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) TypeAndVersion() (string, error) {
	return _DestinationVerifierProxy.Contract.TypeAndVersion(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _DestinationVerifierProxy.Contract.TypeAndVersion(&_DestinationVerifierProxy.CallOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DestinationVerifierProxy.contract.Transact(opts, "acceptOwnership")
}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.AcceptOwnership(&_DestinationVerifierProxy.TransactOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.AcceptOwnership(&_DestinationVerifierProxy.TransactOpts)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactor) SetVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.contract.Transact(opts, "setVerifier", verifierAddress)
}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) SetVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.SetVerifier(&_DestinationVerifierProxy.TransactOpts, verifierAddress)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorSession) SetVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.SetVerifier(&_DestinationVerifierProxy.TransactOpts, verifierAddress)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.TransferOwnership(&_DestinationVerifierProxy.TransactOpts, to)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.TransferOwnership(&_DestinationVerifierProxy.TransactOpts, to)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactor) Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.contract.Transact(opts, "verify", payload, parameterPayload)
}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.Verify(&_DestinationVerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorSession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.Verify(&_DestinationVerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactor) VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.contract.Transact(opts, "verifyBulk", payloads, parameterPayload)
}

func (_DestinationVerifierProxy *DestinationVerifierProxySession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.VerifyBulk(&_DestinationVerifierProxy.TransactOpts, payloads, parameterPayload)
}

func (_DestinationVerifierProxy *DestinationVerifierProxyTransactorSession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _DestinationVerifierProxy.Contract.VerifyBulk(&_DestinationVerifierProxy.TransactOpts, payloads, parameterPayload)
}

type DestinationVerifierProxyOwnershipTransferRequestedIterator struct {
	Event *DestinationVerifierProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierProxyOwnershipTransferRequested)
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
		it.Event = new(DestinationVerifierProxyOwnershipTransferRequested)
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

func (it *DestinationVerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxyOwnershipTransferRequestedIterator{contract: _DestinationVerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierProxyOwnershipTransferRequested)
				if err := _DestinationVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*DestinationVerifierProxyOwnershipTransferRequested, error) {
	event := new(DestinationVerifierProxyOwnershipTransferRequested)
	if err := _DestinationVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DestinationVerifierProxyOwnershipTransferredIterator struct {
	Event *DestinationVerifierProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DestinationVerifierProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DestinationVerifierProxyOwnershipTransferred)
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
		it.Event = new(DestinationVerifierProxyOwnershipTransferred)
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

func (it *DestinationVerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DestinationVerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DestinationVerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DestinationVerifierProxyOwnershipTransferredIterator{contract: _DestinationVerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DestinationVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DestinationVerifierProxyOwnershipTransferred)
				if err := _DestinationVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DestinationVerifierProxy *DestinationVerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*DestinationVerifierProxyOwnershipTransferred, error) {
	event := new(DestinationVerifierProxyOwnershipTransferred)
	if err := _DestinationVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_DestinationVerifierProxy *DestinationVerifierProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DestinationVerifierProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _DestinationVerifierProxy.ParseOwnershipTransferRequested(log)
	case _DestinationVerifierProxy.abi.Events["OwnershipTransferred"].ID:
		return _DestinationVerifierProxy.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DestinationVerifierProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DestinationVerifierProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_DestinationVerifierProxy *DestinationVerifierProxy) Address() common.Address {
	return _DestinationVerifierProxy.address
}

type DestinationVerifierProxyInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SAccessController(opts *bind.CallOpts) (common.Address, error)

	SFeeManager(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetVerifier(opts *bind.TransactOpts, verifierAddress common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error)

	VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DestinationVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DestinationVerifierProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DestinationVerifierProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DestinationVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DestinationVerifierProxyOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
