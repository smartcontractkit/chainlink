// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2abstract

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

// OCR2AbstractMetaData contains all meta data concerning the OCR2Abstract contract.
var OCR2AbstractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR2AbstractABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2AbstractMetaData.ABI instead.
var OCR2AbstractABI = OCR2AbstractMetaData.ABI

// OCR2Abstract is an auto generated Go binding around an Ethereum contract.
type OCR2Abstract struct {
	OCR2AbstractCaller     // Read-only binding to the contract
	OCR2AbstractTransactor // Write-only binding to the contract
	OCR2AbstractFilterer   // Log filterer for contract events
}

// OCR2AbstractCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2AbstractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2AbstractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2AbstractSession struct {
	Contract     *OCR2Abstract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2AbstractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2AbstractCallerSession struct {
	Contract *OCR2AbstractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OCR2AbstractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2AbstractTransactorSession struct {
	Contract     *OCR2AbstractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OCR2AbstractRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2AbstractRaw struct {
	Contract *OCR2Abstract // Generic contract binding to access the raw methods on
}

// OCR2AbstractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2AbstractCallerRaw struct {
	Contract *OCR2AbstractCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2AbstractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactorRaw struct {
	Contract *OCR2AbstractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Abstract creates a new instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2Abstract(address common.Address, backend bind.ContractBackend) (*OCR2Abstract, error) {
	contract, err := bindOCR2Abstract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Abstract{OCR2AbstractCaller: OCR2AbstractCaller{contract: contract}, OCR2AbstractTransactor: OCR2AbstractTransactor{contract: contract}, OCR2AbstractFilterer: OCR2AbstractFilterer{contract: contract}}, nil
}

// NewOCR2AbstractCaller creates a new read-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractCaller(address common.Address, caller bind.ContractCaller) (*OCR2AbstractCaller, error) {
	contract, err := bindOCR2Abstract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractCaller{contract: contract}, nil
}

// NewOCR2AbstractTransactor creates a new write-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AbstractTransactor, error) {
	contract, err := bindOCR2Abstract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransactor{contract: contract}, nil
}

// NewOCR2AbstractFilterer creates a new log filterer instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AbstractFilterer, error) {
	contract, err := bindOCR2Abstract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractFilterer{contract: contract}, nil
}

// bindOCR2Abstract binds a generic wrapper to an already deployed contract.
func bindOCR2Abstract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2AbstractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.OCR2AbstractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2AbstractConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Abstract contract.
type OCR2AbstractConfigSetIterator struct {
	Event *OCR2AbstractConfigSet // Event containing the contract specifics and raw log

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
func (it *OCR2AbstractConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractConfigSet)
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
		it.Event = new(OCR2AbstractConfigSet)
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
func (it *OCR2AbstractConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractConfigSet represents a ConfigSet event raised by the OCR2Abstract contract.
type OCR2AbstractConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AbstractConfigSetIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractConfigSetIterator{contract: _OCR2Abstract.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AbstractConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractConfigSet)
				if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseConfigSet(log types.Log) (*OCR2AbstractConfigSet, error) {
	event := new(OCR2AbstractConfigSet)
	if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AbstractTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Abstract contract.
type OCR2AbstractTransmittedIterator struct {
	Event *OCR2AbstractTransmitted // Event containing the contract specifics and raw log

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
func (it *OCR2AbstractTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractTransmitted)
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
		it.Event = new(OCR2AbstractTransmitted)
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
func (it *OCR2AbstractTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractTransmitted represents a Transmitted event raised by the OCR2Abstract contract.
type OCR2AbstractTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AbstractTransmittedIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransmittedIterator{contract: _OCR2Abstract.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AbstractTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractTransmitted)
				if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseTransmitted(log types.Log) (*OCR2AbstractTransmitted, error) {
	event := new(OCR2AbstractTransmitted)
	if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TypeAndVersionInterfaceMetaData contains all meta data concerning the TypeAndVersionInterface contract.
var TypeAndVersionInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// TypeAndVersionInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use TypeAndVersionInterfaceMetaData.ABI instead.
var TypeAndVersionInterfaceABI = TypeAndVersionInterfaceMetaData.ABI

// TypeAndVersionInterface is an auto generated Go binding around an Ethereum contract.
type TypeAndVersionInterface struct {
	TypeAndVersionInterfaceCaller     // Read-only binding to the contract
	TypeAndVersionInterfaceTransactor // Write-only binding to the contract
	TypeAndVersionInterfaceFilterer   // Log filterer for contract events
}

// TypeAndVersionInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TypeAndVersionInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TypeAndVersionInterfaceSession struct {
	Contract     *TypeAndVersionInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TypeAndVersionInterfaceCallerSession struct {
	Contract *TypeAndVersionInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TypeAndVersionInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TypeAndVersionInterfaceTransactorSession struct {
	Contract     *TypeAndVersionInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type TypeAndVersionInterfaceRaw struct {
	Contract *TypeAndVersionInterface // Generic contract binding to access the raw methods on
}

// TypeAndVersionInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCallerRaw struct {
	Contract *TypeAndVersionInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// TypeAndVersionInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactorRaw struct {
	Contract *TypeAndVersionInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTypeAndVersionInterface creates a new instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterface(address common.Address, backend bind.ContractBackend) (*TypeAndVersionInterface, error) {
	contract, err := bindTypeAndVersionInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterface{TypeAndVersionInterfaceCaller: TypeAndVersionInterfaceCaller{contract: contract}, TypeAndVersionInterfaceTransactor: TypeAndVersionInterfaceTransactor{contract: contract}, TypeAndVersionInterfaceFilterer: TypeAndVersionInterfaceFilterer{contract: contract}}, nil
}

// NewTypeAndVersionInterfaceCaller creates a new read-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceCaller(address common.Address, caller bind.ContractCaller) (*TypeAndVersionInterfaceCaller, error) {
	contract, err := bindTypeAndVersionInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceCaller{contract: contract}, nil
}

// NewTypeAndVersionInterfaceTransactor creates a new write-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*TypeAndVersionInterfaceTransactor, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceTransactor{contract: contract}, nil
}

// NewTypeAndVersionInterfaceFilterer creates a new log filterer instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*TypeAndVersionInterfaceFilterer, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceFilterer{contract: contract}, nil
}

// bindTypeAndVersionInterface binds a generic wrapper to an already deployed contract.
func bindTypeAndVersionInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TypeAndVersionInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transact(opts, method, params...)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TypeAndVersionInterface.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}
