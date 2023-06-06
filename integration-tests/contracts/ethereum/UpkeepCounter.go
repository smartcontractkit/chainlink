// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// UpkeepCounterMetaData contains all meta data concerning the UpkeepCounter contract.
var UpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516104423803806104428339818101604052604081101561003357600080fd5b508051602090910151600091825560015560038190554360025560048190556005556103de806100646000396000f3fe608060405234801561001057600080fd5b506004361061008e5760003560e01c80632cb15864146100935780634585e33b146100ad57806361bc221a1461011d5780636250a13a146101255780636e04ff0d1461012d5780637f407edf1461021c578063806b984f1461023f578063917d895f14610247578063947a36fb1461024f578063d832d92f14610257575b600080fd5b61009b610273565b60408051918252519081900360200190f35b61011b600480360360208110156100c357600080fd5b810190602081018135600160201b8111156100dd57600080fd5b8201836020820111156100ef57600080fd5b803590602001918460018302840111600160201b8311171561011057600080fd5b509092509050610279565b005b61009b6102f0565b61009b6102f6565b61019b6004803603602081101561014357600080fd5b810190602081018135600160201b81111561015d57600080fd5b82018360208201111561016f57600080fd5b803590602001918460018302840111600160201b8311171561019057600080fd5b5090925090506102fc565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156101e05781810151838201526020016101c8565b50505050905090810190601f16801561020d5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61011b6004803603604081101561023257600080fd5b508035906020013561034e565b61009b610360565b61009b610366565b61009b61036c565b61025f610372565b604080519115158252519081900360200190f35b60045481565b60045461028557436004555b4360028190556005805460010190819055600454600354604080519283526020830194909452818401526060810191909152905132917f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa919081900360800190a25050600254600355565b60055481565b60005481565b60006060610308610372565b848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b60009182556001556004819055600555565b60025481565b60035481565b60015481565b600060045460001415610387575060016103a5565b60005460045443031080156103a25750600154600254430310155b90505b9056fea264697066735822122040ec1d65fcf245bb817ffdc5ef880653b8cc880e41da1037b27143dcb90127c864736f6c63430007060033",
}

// UpkeepCounterABI is the input ABI used to generate the binding from.
// Deprecated: Use UpkeepCounterMetaData.ABI instead.
var UpkeepCounterABI = UpkeepCounterMetaData.ABI

// UpkeepCounterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpkeepCounterMetaData.Bin instead.
var UpkeepCounterBin = UpkeepCounterMetaData.Bin

// DeployUpkeepCounter deploys a new Ethereum contract, binding an instance of UpkeepCounter to it.
func DeployUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepCounter, error) {
	parsed, err := UpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepCounterBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepCounter{UpkeepCounterCaller: UpkeepCounterCaller{contract: contract}, UpkeepCounterTransactor: UpkeepCounterTransactor{contract: contract}, UpkeepCounterFilterer: UpkeepCounterFilterer{contract: contract}}, nil
}

// UpkeepCounter is an auto generated Go binding around an Ethereum contract.
type UpkeepCounter struct {
	UpkeepCounterCaller     // Read-only binding to the contract
	UpkeepCounterTransactor // Write-only binding to the contract
	UpkeepCounterFilterer   // Log filterer for contract events
}

// UpkeepCounterCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpkeepCounterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepCounterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpkeepCounterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepCounterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpkeepCounterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepCounterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpkeepCounterSession struct {
	Contract     *UpkeepCounter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpkeepCounterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpkeepCounterCallerSession struct {
	Contract *UpkeepCounterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// UpkeepCounterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpkeepCounterTransactorSession struct {
	Contract     *UpkeepCounterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// UpkeepCounterRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpkeepCounterRaw struct {
	Contract *UpkeepCounter // Generic contract binding to access the raw methods on
}

// UpkeepCounterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpkeepCounterCallerRaw struct {
	Contract *UpkeepCounterCaller // Generic read-only contract binding to access the raw methods on
}

// UpkeepCounterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpkeepCounterTransactorRaw struct {
	Contract *UpkeepCounterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpkeepCounter creates a new instance of UpkeepCounter, bound to a specific deployed contract.
func NewUpkeepCounter(address common.Address, backend bind.ContractBackend) (*UpkeepCounter, error) {
	contract, err := bindUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounter{UpkeepCounterCaller: UpkeepCounterCaller{contract: contract}, UpkeepCounterTransactor: UpkeepCounterTransactor{contract: contract}, UpkeepCounterFilterer: UpkeepCounterFilterer{contract: contract}}, nil
}

// NewUpkeepCounterCaller creates a new read-only instance of UpkeepCounter, bound to a specific deployed contract.
func NewUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*UpkeepCounterCaller, error) {
	contract, err := bindUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterCaller{contract: contract}, nil
}

// NewUpkeepCounterTransactor creates a new write-only instance of UpkeepCounter, bound to a specific deployed contract.
func NewUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepCounterTransactor, error) {
	contract, err := bindUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterTransactor{contract: contract}, nil
}

// NewUpkeepCounterFilterer creates a new log filterer instance of UpkeepCounter, bound to a specific deployed contract.
func NewUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepCounterFilterer, error) {
	contract, err := bindUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterFilterer{contract: contract}, nil
}

// bindUpkeepCounter binds a generic wrapper to an already deployed contract.
func bindUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepCounter *UpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounter.Contract.UpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepCounter *UpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.UpkeepCounterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepCounter *UpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.UpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepCounter *UpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepCounter *UpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepCounter *UpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_UpkeepCounter *UpkeepCounterCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_UpkeepCounter *UpkeepCounterSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounter.Contract.CheckUpkeep(&_UpkeepCounter.CallOpts, data)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_UpkeepCounter *UpkeepCounterCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounter.Contract.CheckUpkeep(&_UpkeepCounter.CallOpts, data)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) Counter() (*big.Int, error) {
	return _UpkeepCounter.Contract.Counter(&_UpkeepCounter.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _UpkeepCounter.Contract.Counter(&_UpkeepCounter.CallOpts)
}

// Eligible is a free data retrieval call binding the contract method 0xd832d92f.
//
// Solidity: function eligible() view returns(bool)
func (_UpkeepCounter *UpkeepCounterCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Eligible is a free data retrieval call binding the contract method 0xd832d92f.
//
// Solidity: function eligible() view returns(bool)
func (_UpkeepCounter *UpkeepCounterSession) Eligible() (bool, error) {
	return _UpkeepCounter.Contract.Eligible(&_UpkeepCounter.CallOpts)
}

// Eligible is a free data retrieval call binding the contract method 0xd832d92f.
//
// Solidity: function eligible() view returns(bool)
func (_UpkeepCounter *UpkeepCounterCallerSession) Eligible() (bool, error) {
	return _UpkeepCounter.Contract.Eligible(&_UpkeepCounter.CallOpts)
}

// InitialBlock is a free data retrieval call binding the contract method 0x2cb15864.
//
// Solidity: function initialBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialBlock is a free data retrieval call binding the contract method 0x2cb15864.
//
// Solidity: function initialBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialBlock(&_UpkeepCounter.CallOpts)
}

// InitialBlock is a free data retrieval call binding the contract method 0x2cb15864.
//
// Solidity: function initialBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialBlock(&_UpkeepCounter.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) Interval() (*big.Int, error) {
	return _UpkeepCounter.Contract.Interval(&_UpkeepCounter.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) Interval() (*big.Int, error) {
	return _UpkeepCounter.Contract.Interval(&_UpkeepCounter.CallOpts)
}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastBlock(&_UpkeepCounter.CallOpts)
}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastBlock(&_UpkeepCounter.CallOpts)
}

// PreviousPerformBlock is a free data retrieval call binding the contract method 0x917d895f.
//
// Solidity: function previousPerformBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviousPerformBlock is a free data retrieval call binding the contract method 0x917d895f.
//
// Solidity: function previousPerformBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.PreviousPerformBlock(&_UpkeepCounter.CallOpts)
}

// PreviousPerformBlock is a free data retrieval call binding the contract method 0x917d895f.
//
// Solidity: function previousPerformBlock() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.PreviousPerformBlock(&_UpkeepCounter.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterSession) TestRange() (*big.Int, error) {
	return _UpkeepCounter.Contract.TestRange(&_UpkeepCounter.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepCounter *UpkeepCounterCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepCounter.Contract.TestRange(&_UpkeepCounter.CallOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_UpkeepCounter *UpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_UpkeepCounter *UpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_UpkeepCounter *UpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, performData)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _testRange, uint256 _interval) returns()
func (_UpkeepCounter *UpkeepCounterTransactor) SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.contract.Transact(opts, "setSpread", _testRange, _interval)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _testRange, uint256 _interval) returns()
func (_UpkeepCounter *UpkeepCounterSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.SetSpread(&_UpkeepCounter.TransactOpts, _testRange, _interval)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _testRange, uint256 _interval) returns()
func (_UpkeepCounter *UpkeepCounterTransactorSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.SetSpread(&_UpkeepCounter.TransactOpts, _testRange, _interval)
}

// UpkeepCounterPerformingUpkeepIterator is returned from FilterPerformingUpkeep and is used to iterate over the raw logs and unpacked data for PerformingUpkeep events raised by the UpkeepCounter contract.
type UpkeepCounterPerformingUpkeepIterator struct {
	Event *UpkeepCounterPerformingUpkeep // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepCounterPerformingUpkeepIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepCounterPerformingUpkeep)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpkeepCounterPerformingUpkeep)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepCounterPerformingUpkeep represents a PerformingUpkeep event raised by the UpkeepCounter contract.
type UpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPerformingUpkeep is a free log retrieval operation binding the contract event 0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa.
//
// Solidity: event PerformingUpkeep(address indexed from, uint256 initialBlock, uint256 lastBlock, uint256 previousBlock, uint256 counter)
func (_UpkeepCounter *UpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformingUpkeepIterator{contract: _UpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

// WatchPerformingUpkeep is a free log subscription operation binding the contract event 0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa.
//
// Solidity: event PerformingUpkeep(address indexed from, uint256 initialBlock, uint256 lastBlock, uint256 previousBlock, uint256 counter)
func (_UpkeepCounter *UpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpkeepCounterPerformingUpkeep)
				if err := _UpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

// ParsePerformingUpkeep is a log parse operation binding the contract event 0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa.
//
// Solidity: event PerformingUpkeep(address indexed from, uint256 initialBlock, uint256 lastBlock, uint256 previousBlock, uint256 counter)
func (_UpkeepCounter *UpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformingUpkeep, error) {
	event := new(UpkeepCounterPerformingUpkeep)
	if err := _UpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
