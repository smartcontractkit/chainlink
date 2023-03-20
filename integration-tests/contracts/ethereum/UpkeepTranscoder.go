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

// UpkeepTranscoderMetaData contains all meta data concerning the UpkeepTranscoder contract.
var UpkeepTranscoderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidTranscoding\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"fromVersion\",\"type\":\"uint8\"},{\"internalType\":\"enumUpkeepFormat\",\"name\":\"toVersion\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"transcodeUpkeeps\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610276806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063181f5a771461003b578063c71249ab14610086575b600080fd5b61007060405180604001604052806016815260200175055706b6565705472616e73636f64657220312e302e360541b81525081565b60405161007d919061016e565b60405180910390f35b61007061009436600461019c565b606060008580156100a7576100a761022a565b1415806100c3575060008480156100c0576100c061022a565b14155b156100e1576040516390aaccc360e01b815260040160405180910390fd5b82828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509298975050505050505050565b6000815180845260005b818110156101475760208185018101518683018201520161012b565b81811115610159576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006101816020830184610121565b9392505050565b80356001811061019757600080fd5b919050565b600080600080606085870312156101b257600080fd5b6101bb85610188565b93506101c960208601610188565b9250604085013567ffffffffffffffff808211156101e657600080fd5b818701915087601f8301126101fa57600080fd5b81358181111561020957600080fd5b88602082850101111561021b57600080fd5b95989497505060200194505050565b634e487b7160e01b600052602160045260246000fdfea26469706673582212200831a58eaa3d637278be672e8f608672a90ebe54a5ae1e6fdb7b544152d6b00b64736f6c634300080d0033",
}

// UpkeepTranscoderABI is the input ABI used to generate the binding from.
// Deprecated: Use UpkeepTranscoderMetaData.ABI instead.
var UpkeepTranscoderABI = UpkeepTranscoderMetaData.ABI

// UpkeepTranscoderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpkeepTranscoderMetaData.Bin instead.
var UpkeepTranscoderBin = UpkeepTranscoderMetaData.Bin

// DeployUpkeepTranscoder deploys a new Ethereum contract, binding an instance of UpkeepTranscoder to it.
func DeployUpkeepTranscoder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UpkeepTranscoder, error) {
	parsed, err := UpkeepTranscoderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepTranscoderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepTranscoder{UpkeepTranscoderCaller: UpkeepTranscoderCaller{contract: contract}, UpkeepTranscoderTransactor: UpkeepTranscoderTransactor{contract: contract}, UpkeepTranscoderFilterer: UpkeepTranscoderFilterer{contract: contract}}, nil
}

// UpkeepTranscoder is an auto generated Go binding around an Ethereum contract.
type UpkeepTranscoder struct {
	UpkeepTranscoderCaller     // Read-only binding to the contract
	UpkeepTranscoderTransactor // Write-only binding to the contract
	UpkeepTranscoderFilterer   // Log filterer for contract events
}

// UpkeepTranscoderCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpkeepTranscoderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepTranscoderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpkeepTranscoderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepTranscoderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpkeepTranscoderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepTranscoderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpkeepTranscoderSession struct {
	Contract     *UpkeepTranscoder // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpkeepTranscoderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpkeepTranscoderCallerSession struct {
	Contract *UpkeepTranscoderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// UpkeepTranscoderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpkeepTranscoderTransactorSession struct {
	Contract     *UpkeepTranscoderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UpkeepTranscoderRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpkeepTranscoderRaw struct {
	Contract *UpkeepTranscoder // Generic contract binding to access the raw methods on
}

// UpkeepTranscoderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpkeepTranscoderCallerRaw struct {
	Contract *UpkeepTranscoderCaller // Generic read-only contract binding to access the raw methods on
}

// UpkeepTranscoderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpkeepTranscoderTransactorRaw struct {
	Contract *UpkeepTranscoderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpkeepTranscoder creates a new instance of UpkeepTranscoder, bound to a specific deployed contract.
func NewUpkeepTranscoder(address common.Address, backend bind.ContractBackend) (*UpkeepTranscoder, error) {
	contract, err := bindUpkeepTranscoder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoder{UpkeepTranscoderCaller: UpkeepTranscoderCaller{contract: contract}, UpkeepTranscoderTransactor: UpkeepTranscoderTransactor{contract: contract}, UpkeepTranscoderFilterer: UpkeepTranscoderFilterer{contract: contract}}, nil
}

// NewUpkeepTranscoderCaller creates a new read-only instance of UpkeepTranscoder, bound to a specific deployed contract.
func NewUpkeepTranscoderCaller(address common.Address, caller bind.ContractCaller) (*UpkeepTranscoderCaller, error) {
	contract, err := bindUpkeepTranscoder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderCaller{contract: contract}, nil
}

// NewUpkeepTranscoderTransactor creates a new write-only instance of UpkeepTranscoder, bound to a specific deployed contract.
func NewUpkeepTranscoderTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepTranscoderTransactor, error) {
	contract, err := bindUpkeepTranscoder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderTransactor{contract: contract}, nil
}

// NewUpkeepTranscoderFilterer creates a new log filterer instance of UpkeepTranscoder, bound to a specific deployed contract.
func NewUpkeepTranscoderFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepTranscoderFilterer, error) {
	contract, err := bindUpkeepTranscoder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepTranscoderFilterer{contract: contract}, nil
}

// bindUpkeepTranscoder binds a generic wrapper to an already deployed contract.
func bindUpkeepTranscoder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepTranscoderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepTranscoder *UpkeepTranscoderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepTranscoder *UpkeepTranscoderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepTranscoder *UpkeepTranscoderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.UpkeepTranscoderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepTranscoder *UpkeepTranscoderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepTranscoder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepTranscoder *UpkeepTranscoderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepTranscoder *UpkeepTranscoderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepTranscoder.Contract.contract.Transact(opts, method, params...)
}

// TranscodeUpkeeps is a free data retrieval call binding the contract method 0xc71249ab.
//
// Solidity: function transcodeUpkeeps(uint8 fromVersion, uint8 toVersion, bytes encodedUpkeeps) view returns(bytes)
func (_UpkeepTranscoder *UpkeepTranscoderCaller) TranscodeUpkeeps(opts *bind.CallOpts, fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	var out []interface{}
	err := _UpkeepTranscoder.contract.Call(opts, &out, "transcodeUpkeeps", fromVersion, toVersion, encodedUpkeeps)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TranscodeUpkeeps is a free data retrieval call binding the contract method 0xc71249ab.
//
// Solidity: function transcodeUpkeeps(uint8 fromVersion, uint8 toVersion, bytes encodedUpkeeps) view returns(bytes)
func (_UpkeepTranscoder *UpkeepTranscoderSession) TranscodeUpkeeps(fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	return _UpkeepTranscoder.Contract.TranscodeUpkeeps(&_UpkeepTranscoder.CallOpts, fromVersion, toVersion, encodedUpkeeps)
}

// TranscodeUpkeeps is a free data retrieval call binding the contract method 0xc71249ab.
//
// Solidity: function transcodeUpkeeps(uint8 fromVersion, uint8 toVersion, bytes encodedUpkeeps) view returns(bytes)
func (_UpkeepTranscoder *UpkeepTranscoderCallerSession) TranscodeUpkeeps(fromVersion uint8, toVersion uint8, encodedUpkeeps []byte) ([]byte, error) {
	return _UpkeepTranscoder.Contract.TranscodeUpkeeps(&_UpkeepTranscoder.CallOpts, fromVersion, toVersion, encodedUpkeeps)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_UpkeepTranscoder *UpkeepTranscoderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepTranscoder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_UpkeepTranscoder *UpkeepTranscoderSession) TypeAndVersion() (string, error) {
	return _UpkeepTranscoder.Contract.TypeAndVersion(&_UpkeepTranscoder.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_UpkeepTranscoder *UpkeepTranscoderCallerSession) TypeAndVersion() (string, error) {
	return _UpkeepTranscoder.Contract.TypeAndVersion(&_UpkeepTranscoder.CallOpts)
}
