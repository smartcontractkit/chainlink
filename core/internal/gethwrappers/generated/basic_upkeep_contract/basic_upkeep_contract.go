// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package basic_upkeep_contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BasicUpkeepContractABI is the input ABI used to generate the binding from.
const BasicUpkeepContractABI = "[{\"inputs\":[],\"name\":\"bytesToSend\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receivedBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytes\",\"type\":\"bytes\"}],\"name\":\"setBytesToSend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_should\",\"type\":\"bool\"}],\"name\":\"setShouldPerformUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldPerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BasicUpkeepContractBin is the compiled bytecode used for deploying new contracts.
var BasicUpkeepContractBin = "0x608060405234801561001057600080fd5b5061072a806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80634585e33b1161005b5780634585e33b146101af5780636e04ff0d1461021f57806384aadecc14610310578063bb4462681461032f5761007d565b80630427e4b7146100825780632c3b84ac146100ff57806333437c77146101a7575b600080fd5b61008a61034b565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100c45781810151838201526020016100ac565b50505050905090810190601f1680156100f15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101a56004803603602081101561011557600080fd5b81019060208101813564010000000081111561013057600080fd5b82018360208201111561014257600080fd5b8035906020019184600183028401116401000000008311171561016457600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506103f4945050505050565b005b61008a61040b565b6101a5600480360360208110156101c557600080fd5b8101906020810181356401000000008111156101e057600080fd5b8201836020820111156101f257600080fd5b8035906020019184600183028401116401000000008311171561021457600080fd5b509092509050610483565b61028f6004803603602081101561023557600080fd5b81019060208101813564010000000081111561025057600080fd5b82018360208201111561026257600080fd5b8035906020019184600183028401116401000000008311171561028457600080fd5b5090925090506104bc565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102d45781810151838201526020016102bc565b50505050905090810190601f1680156103015780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101a56004803603602081101561032657600080fd5b5035151561057f565b6103376105b0565b604080519115158252519081900360200190f35b600280546040805160206001841615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01909316849004601f810184900484028201840190925281815292918301828280156103ec5780601f106103c1576101008083540402835291602001916103ec565b820191906000526020600020905b8154815290600101906020018083116103cf57829003601f168201915b505050505081565b80516104079060019060208401906105b9565b5050565b60018054604080516020600284861615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f810184900484028201840190925281815292918301828280156103ec5780601f106103c1576101008083540402835291602001916103ec565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556104b760028383610645565b505050565b6000805460018054604080516020600261010085871615027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f810184900484028201840190925281815260609460ff1693929091839183018282801561056d5780601f106105425761010080835404028352916020019161056d565b820191906000526020600020905b81548152906001019060200180831161055057829003601f168201915b50505050509050915091509250929050565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005460ff1681565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826105ef5760008555610635565b82601f1061060857805160ff1916838001178555610635565b82800160010185558215610635579182015b8281111561063557825182559160200191906001019061061a565b506106419291506106df565b5090565b828054600181600116156101000203166002900490600052602060002090601f01602090048101928261067b5760008555610635565b82601f106106b2578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555610635565b82800160010185558215610635579182015b828111156106355782358255916020019190600101906106c4565b5b8082111561064157600081556001016106e056fea2646970667358221220d14dfba5a1efc4145905e37243ce77c189567ba18352066c0fe8a0c3ba53d27e64736f6c63430007060033"

// DeployBasicUpkeepContract deploys a new Ethereum contract, binding an instance of BasicUpkeepContract to it.
func DeployBasicUpkeepContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BasicUpkeepContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BasicUpkeepContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BasicUpkeepContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BasicUpkeepContract{BasicUpkeepContractCaller: BasicUpkeepContractCaller{contract: contract}, BasicUpkeepContractTransactor: BasicUpkeepContractTransactor{contract: contract}, BasicUpkeepContractFilterer: BasicUpkeepContractFilterer{contract: contract}}, nil
}

// BasicUpkeepContract is an auto generated Go binding around an Ethereum contract.
type BasicUpkeepContract struct {
	BasicUpkeepContractCaller     // Read-only binding to the contract
	BasicUpkeepContractTransactor // Write-only binding to the contract
	BasicUpkeepContractFilterer   // Log filterer for contract events
}

// BasicUpkeepContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type BasicUpkeepContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicUpkeepContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BasicUpkeepContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicUpkeepContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BasicUpkeepContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BasicUpkeepContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BasicUpkeepContractSession struct {
	Contract     *BasicUpkeepContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BasicUpkeepContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BasicUpkeepContractCallerSession struct {
	Contract *BasicUpkeepContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// BasicUpkeepContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BasicUpkeepContractTransactorSession struct {
	Contract     *BasicUpkeepContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// BasicUpkeepContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type BasicUpkeepContractRaw struct {
	Contract *BasicUpkeepContract // Generic contract binding to access the raw methods on
}

// BasicUpkeepContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BasicUpkeepContractCallerRaw struct {
	Contract *BasicUpkeepContractCaller // Generic read-only contract binding to access the raw methods on
}

// BasicUpkeepContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BasicUpkeepContractTransactorRaw struct {
	Contract *BasicUpkeepContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBasicUpkeepContract creates a new instance of BasicUpkeepContract, bound to a specific deployed contract.
func NewBasicUpkeepContract(address common.Address, backend bind.ContractBackend) (*BasicUpkeepContract, error) {
	contract, err := bindBasicUpkeepContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BasicUpkeepContract{BasicUpkeepContractCaller: BasicUpkeepContractCaller{contract: contract}, BasicUpkeepContractTransactor: BasicUpkeepContractTransactor{contract: contract}, BasicUpkeepContractFilterer: BasicUpkeepContractFilterer{contract: contract}}, nil
}

// NewBasicUpkeepContractCaller creates a new read-only instance of BasicUpkeepContract, bound to a specific deployed contract.
func NewBasicUpkeepContractCaller(address common.Address, caller bind.ContractCaller) (*BasicUpkeepContractCaller, error) {
	contract, err := bindBasicUpkeepContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BasicUpkeepContractCaller{contract: contract}, nil
}

// NewBasicUpkeepContractTransactor creates a new write-only instance of BasicUpkeepContract, bound to a specific deployed contract.
func NewBasicUpkeepContractTransactor(address common.Address, transactor bind.ContractTransactor) (*BasicUpkeepContractTransactor, error) {
	contract, err := bindBasicUpkeepContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BasicUpkeepContractTransactor{contract: contract}, nil
}

// NewBasicUpkeepContractFilterer creates a new log filterer instance of BasicUpkeepContract, bound to a specific deployed contract.
func NewBasicUpkeepContractFilterer(address common.Address, filterer bind.ContractFilterer) (*BasicUpkeepContractFilterer, error) {
	contract, err := bindBasicUpkeepContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BasicUpkeepContractFilterer{contract: contract}, nil
}

// bindBasicUpkeepContract binds a generic wrapper to an already deployed contract.
func bindBasicUpkeepContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BasicUpkeepContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BasicUpkeepContract *BasicUpkeepContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BasicUpkeepContract.Contract.BasicUpkeepContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BasicUpkeepContract *BasicUpkeepContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.BasicUpkeepContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BasicUpkeepContract *BasicUpkeepContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.BasicUpkeepContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BasicUpkeepContract *BasicUpkeepContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BasicUpkeepContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BasicUpkeepContract *BasicUpkeepContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BasicUpkeepContract *BasicUpkeepContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.contract.Transact(opts, method, params...)
}

// BytesToSend is a free data retrieval call binding the contract method 0x33437c77.
//
// Solidity: function bytesToSend() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractCaller) BytesToSend(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BasicUpkeepContract.contract.Call(opts, &out, "bytesToSend")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// BytesToSend is a free data retrieval call binding the contract method 0x33437c77.
//
// Solidity: function bytesToSend() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractSession) BytesToSend() ([]byte, error) {
	return _BasicUpkeepContract.Contract.BytesToSend(&_BasicUpkeepContract.CallOpts)
}

// BytesToSend is a free data retrieval call binding the contract method 0x33437c77.
//
// Solidity: function bytesToSend() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractCallerSession) BytesToSend() ([]byte, error) {
	return _BasicUpkeepContract.Contract.BytesToSend(&_BasicUpkeepContract.CallOpts)
}

// ReceivedBytes is a free data retrieval call binding the contract method 0x0427e4b7.
//
// Solidity: function receivedBytes() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractCaller) ReceivedBytes(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BasicUpkeepContract.contract.Call(opts, &out, "receivedBytes")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ReceivedBytes is a free data retrieval call binding the contract method 0x0427e4b7.
//
// Solidity: function receivedBytes() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractSession) ReceivedBytes() ([]byte, error) {
	return _BasicUpkeepContract.Contract.ReceivedBytes(&_BasicUpkeepContract.CallOpts)
}

// ReceivedBytes is a free data retrieval call binding the contract method 0x0427e4b7.
//
// Solidity: function receivedBytes() view returns(bytes)
func (_BasicUpkeepContract *BasicUpkeepContractCallerSession) ReceivedBytes() ([]byte, error) {
	return _BasicUpkeepContract.Contract.ReceivedBytes(&_BasicUpkeepContract.CallOpts)
}

// ShouldPerformUpkeep is a free data retrieval call binding the contract method 0xbb446268.
//
// Solidity: function shouldPerformUpkeep() view returns(bool)
func (_BasicUpkeepContract *BasicUpkeepContractCaller) ShouldPerformUpkeep(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BasicUpkeepContract.contract.Call(opts, &out, "shouldPerformUpkeep")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ShouldPerformUpkeep is a free data retrieval call binding the contract method 0xbb446268.
//
// Solidity: function shouldPerformUpkeep() view returns(bool)
func (_BasicUpkeepContract *BasicUpkeepContractSession) ShouldPerformUpkeep() (bool, error) {
	return _BasicUpkeepContract.Contract.ShouldPerformUpkeep(&_BasicUpkeepContract.CallOpts)
}

// ShouldPerformUpkeep is a free data retrieval call binding the contract method 0xbb446268.
//
// Solidity: function shouldPerformUpkeep() view returns(bool)
func (_BasicUpkeepContract *BasicUpkeepContractCallerSession) ShouldPerformUpkeep() (bool, error) {
	return _BasicUpkeepContract.Contract.ShouldPerformUpkeep(&_BasicUpkeepContract.CallOpts)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_BasicUpkeepContract *BasicUpkeepContractTransactor) CheckUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.contract.Transact(opts, "checkUpkeep", data)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_BasicUpkeepContract *BasicUpkeepContractSession) CheckUpkeep(data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.CheckUpkeep(&_BasicUpkeepContract.TransactOpts, data)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) returns(bool, bytes)
func (_BasicUpkeepContract *BasicUpkeepContractTransactorSession) CheckUpkeep(data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.CheckUpkeep(&_BasicUpkeepContract.TransactOpts, data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactor) PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.contract.Transact(opts, "performUpkeep", data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_BasicUpkeepContract *BasicUpkeepContractSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.PerformUpkeep(&_BasicUpkeepContract.TransactOpts, data)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes data) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactorSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.PerformUpkeep(&_BasicUpkeepContract.TransactOpts, data)
}

// SetBytesToSend is a paid mutator transaction binding the contract method 0x2c3b84ac.
//
// Solidity: function setBytesToSend(bytes _bytes) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactor) SetBytesToSend(opts *bind.TransactOpts, _bytes []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.contract.Transact(opts, "setBytesToSend", _bytes)
}

// SetBytesToSend is a paid mutator transaction binding the contract method 0x2c3b84ac.
//
// Solidity: function setBytesToSend(bytes _bytes) returns()
func (_BasicUpkeepContract *BasicUpkeepContractSession) SetBytesToSend(_bytes []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.SetBytesToSend(&_BasicUpkeepContract.TransactOpts, _bytes)
}

// SetBytesToSend is a paid mutator transaction binding the contract method 0x2c3b84ac.
//
// Solidity: function setBytesToSend(bytes _bytes) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactorSession) SetBytesToSend(_bytes []byte) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.SetBytesToSend(&_BasicUpkeepContract.TransactOpts, _bytes)
}

// SetShouldPerformUpkeep is a paid mutator transaction binding the contract method 0x84aadecc.
//
// Solidity: function setShouldPerformUpkeep(bool _should) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactor) SetShouldPerformUpkeep(opts *bind.TransactOpts, _should bool) (*types.Transaction, error) {
	return _BasicUpkeepContract.contract.Transact(opts, "setShouldPerformUpkeep", _should)
}

// SetShouldPerformUpkeep is a paid mutator transaction binding the contract method 0x84aadecc.
//
// Solidity: function setShouldPerformUpkeep(bool _should) returns()
func (_BasicUpkeepContract *BasicUpkeepContractSession) SetShouldPerformUpkeep(_should bool) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.SetShouldPerformUpkeep(&_BasicUpkeepContract.TransactOpts, _should)
}

// SetShouldPerformUpkeep is a paid mutator transaction binding the contract method 0x84aadecc.
//
// Solidity: function setShouldPerformUpkeep(bool _should) returns()
func (_BasicUpkeepContract *BasicUpkeepContractTransactorSession) SetShouldPerformUpkeep(_should bool) (*types.Transaction, error) {
	return _BasicUpkeepContract.Contract.SetShouldPerformUpkeep(&_BasicUpkeepContract.TransactOpts, _should)
}
