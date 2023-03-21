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

// KeeperConsumerMetaData contains all meta data concerning the KeeperConsumer contract.
var KeeperConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updateInterval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimeStamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516103583803806103588339818101604052602081101561003357600080fd5b505160805242600155600080556080516102fe61005a6000398061025452506102fe6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80633f3b3b271461005c5780634585e33b1461007657806361bc221a146100e65780636e04ff0d146100ee578063947a36fb146101dd575b600080fd5b6100646101e5565b60408051918252519081900360200190f35b6100e46004803603602081101561008c57600080fd5b810190602081018135600160201b8111156100a657600080fd5b8201836020820111156100b857600080fd5b803590602001918460018302840111600160201b831117156100d957600080fd5b5090925090506101eb565b005b6100646101f8565b61015c6004803603602081101561010457600080fd5b810190602081018135600160201b81111561011e57600080fd5b82018360208201111561013057600080fd5b803590602001918460018302840111600160201b8311171561015157600080fd5b5090925090506101fe565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156101a1578181015183820152602001610189565b50505050905090810190601f1680156101ce5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b610064610252565b60015481565b5050600080546001019055565b60005481565b6000606061020a610276565b6001848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b32156102c6576040805162461bcd60e51b815260206004820152601a6024820152791bdb9b1e48199bdc881cda5b5d5b185d195908189858dad95b9960321b604482015290519081900360640190fd5b56fea2646970667358221220c0e089efa59b00d8b131c6b0456904c0ef8f5646c27f81de540a4cc400cff70c64736f6c63430007060033",
}

// KeeperConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperConsumerMetaData.ABI instead.
var KeeperConsumerABI = KeeperConsumerMetaData.ABI

// KeeperConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperConsumerMetaData.Bin instead.
var KeeperConsumerBin = KeeperConsumerMetaData.Bin

// DeployKeeperConsumer deploys a new Ethereum contract, binding an instance of KeeperConsumer to it.
func DeployKeeperConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, updateInterval *big.Int) (common.Address, *types.Transaction, *KeeperConsumer, error) {
	parsed, err := KeeperConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperConsumerBin), backend, updateInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumer{KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

// KeeperConsumer is an auto generated Go binding around an Ethereum contract.
type KeeperConsumer struct {
	KeeperConsumerCaller     // Read-only binding to the contract
	KeeperConsumerTransactor // Write-only binding to the contract
	KeeperConsumerFilterer   // Log filterer for contract events
}

// KeeperConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperConsumerSession struct {
	Contract     *KeeperConsumer   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperConsumerCallerSession struct {
	Contract *KeeperConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// KeeperConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperConsumerTransactorSession struct {
	Contract     *KeeperConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// KeeperConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperConsumerRaw struct {
	Contract *KeeperConsumer // Generic contract binding to access the raw methods on
}

// KeeperConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperConsumerCallerRaw struct {
	Contract *KeeperConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperConsumerTransactorRaw struct {
	Contract *KeeperConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperConsumer creates a new instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumer(address common.Address, backend bind.ContractBackend) (*KeeperConsumer, error) {
	contract, err := bindKeeperConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumer{KeeperConsumerCaller: KeeperConsumerCaller{contract: contract}, KeeperConsumerTransactor: KeeperConsumerTransactor{contract: contract}, KeeperConsumerFilterer: KeeperConsumerFilterer{contract: contract}}, nil
}

// NewKeeperConsumerCaller creates a new read-only instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerCaller, error) {
	contract, err := bindKeeperConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerCaller{contract: contract}, nil
}

// NewKeeperConsumerTransactor creates a new write-only instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerTransactor, error) {
	contract, err := bindKeeperConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerTransactor{contract: contract}, nil
}

// NewKeeperConsumerFilterer creates a new log filterer instance of KeeperConsumer, bound to a specific deployed contract.
func NewKeeperConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerFilterer, error) {
	contract, err := bindKeeperConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerFilterer{contract: contract}, nil
}

// bindKeeperConsumer binds a generic wrapper to an already deployed contract.
func bindKeeperConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumer *KeeperConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.KeeperConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumer *KeeperConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumer *KeeperConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.KeeperConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperConsumer *KeeperConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperConsumer *KeeperConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.contract.Transact(opts, method, params...)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) view returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerCaller) CheckUpkeep(opts *bind.CallOpts, checkData []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "checkUpkeep", checkData)

	outstruct := new(struct {
		UpkeepNeeded bool
		PerformData  []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) view returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerSession) CheckUpkeep(checkData []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.CallOpts, checkData)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes checkData) view returns(bool upkeepNeeded, bytes performData)
func (_KeeperConsumer *KeeperConsumerCallerSession) CheckUpkeep(checkData []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error) {
	return _KeeperConsumer.Contract.CheckUpkeep(&_KeeperConsumer.CallOpts, checkData)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) Counter() (*big.Int, error) {
	return _KeeperConsumer.Contract.Counter(&_KeeperConsumer.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

// Interval is a free data retrieval call binding the contract method 0x947a36fb.
//
// Solidity: function interval() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) Interval() (*big.Int, error) {
	return _KeeperConsumer.Contract.Interval(&_KeeperConsumer.CallOpts)
}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCaller) LastTimeStamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumer.contract.Call(opts, &out, "lastTimeStamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

// LastTimeStamp is a free data retrieval call binding the contract method 0x3f3b3b27.
//
// Solidity: function lastTimeStamp() view returns(uint256)
func (_KeeperConsumer *KeeperConsumerCallerSession) LastTimeStamp() (*big.Int, error) {
	return _KeeperConsumer.Contract.LastTimeStamp(&_KeeperConsumer.CallOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.contract.Transact(opts, "performUpkeep", performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_KeeperConsumer *KeeperConsumerTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _KeeperConsumer.Contract.PerformUpkeep(&_KeeperConsumer.TransactOpts, performData)
}
