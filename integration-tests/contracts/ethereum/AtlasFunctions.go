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

// AtlasFunctionsMetaData contains all meta data concerning the AtlasFunctions contract.
var AtlasFunctionsMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"UserCallbackError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"UserCallbackRawError\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"fireOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"fireOracleResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"fireSubscriptionFunded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"fireUserCallbackError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"fireUserCallbackRawError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104c7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806338ddd89c1461005c5780634bc2065f146100715780635fa55d55146100845780637c89a07f146100975780638b2d730b146100aa575b600080fd5b61006f61006a3660046102ab565b6100bd565b005b61006f61007f366004610301565b610102565b61006f61009236600461031a565b610130565b61006f6100a536600461034d565b610171565b61006f6100b8366004610393565b6101ad565b827fa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c30338533866040516100f5959493929190610430565b60405180910390a2505050565b60405181907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a250565b60408051838152602081018390526001600160401b038516917fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f891016100f5565b817fe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2826040516101a1919061047e565b60405180910390a25050565b817fb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c826040516101a1919061047e565b80356001600160401b03811681146101f457600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b0380841115610229576102296101f9565b604051601f8501601f19908116603f01168101908282118183101715610251576102516101f9565b8160405280935085815286868601111561026a57600080fd5b858560208301376000602087830101525050509392505050565b600082601f83011261029557600080fd5b6102a48383356020850161020f565b9392505050565b6000806000606084860312156102c057600080fd5b833592506102d0602085016101dd565b915060408401356001600160401b038111156102eb57600080fd5b6102f786828701610284565b9150509250925092565b60006020828403121561031357600080fd5b5035919050565b60008060006060848603121561032f57600080fd5b610338846101dd565b95602085013595506040909401359392505050565b6000806040838503121561036057600080fd5b8235915060208301356001600160401b0381111561037d57600080fd5b61038985828601610284565b9150509250929050565b600080604083850312156103a657600080fd5b8235915060208301356001600160401b038111156103c357600080fd5b8301601f810185136103d457600080fd5b6103898582356020840161020f565b6000815180845260005b81811015610409576020818501810151868301820152016103ed565b8181111561041b576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b03868116825285811660208301526001600160401b03851660408301528316606082015260a060808201819052600090610473908301846103e3565b979650505050505050565b6020815260006102a460208301846103e356fea2646970667358221220fcf04cf28d335033e93f0442c18b019ba1043c4c073653033ea339a4dfce7c7764736f6c634300080d0033",
}

// AtlasFunctionsABI is the input ABI used to generate the binding from.
// Deprecated: Use AtlasFunctionsMetaData.ABI instead.
var AtlasFunctionsABI = AtlasFunctionsMetaData.ABI

// AtlasFunctionsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AtlasFunctionsMetaData.Bin instead.
var AtlasFunctionsBin = AtlasFunctionsMetaData.Bin

// DeployAtlasFunctions deploys a new Ethereum contract, binding an instance of AtlasFunctions to it.
func DeployAtlasFunctions(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtlasFunctions, error) {
	parsed, err := AtlasFunctionsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AtlasFunctionsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtlasFunctions{AtlasFunctionsCaller: AtlasFunctionsCaller{contract: contract}, AtlasFunctionsTransactor: AtlasFunctionsTransactor{contract: contract}, AtlasFunctionsFilterer: AtlasFunctionsFilterer{contract: contract}}, nil
}

// AtlasFunctions is an auto generated Go binding around an Ethereum contract.
type AtlasFunctions struct {
	AtlasFunctionsCaller     // Read-only binding to the contract
	AtlasFunctionsTransactor // Write-only binding to the contract
	AtlasFunctionsFilterer   // Log filterer for contract events
}

// AtlasFunctionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtlasFunctionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtlasFunctionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtlasFunctionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtlasFunctionsSession struct {
	Contract     *AtlasFunctions   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtlasFunctionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtlasFunctionsCallerSession struct {
	Contract *AtlasFunctionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AtlasFunctionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtlasFunctionsTransactorSession struct {
	Contract     *AtlasFunctionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AtlasFunctionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtlasFunctionsRaw struct {
	Contract *AtlasFunctions // Generic contract binding to access the raw methods on
}

// AtlasFunctionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtlasFunctionsCallerRaw struct {
	Contract *AtlasFunctionsCaller // Generic read-only contract binding to access the raw methods on
}

// AtlasFunctionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtlasFunctionsTransactorRaw struct {
	Contract *AtlasFunctionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtlasFunctions creates a new instance of AtlasFunctions, bound to a specific deployed contract.
func NewAtlasFunctions(address common.Address, backend bind.ContractBackend) (*AtlasFunctions, error) {
	contract, err := bindAtlasFunctions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctions{AtlasFunctionsCaller: AtlasFunctionsCaller{contract: contract}, AtlasFunctionsTransactor: AtlasFunctionsTransactor{contract: contract}, AtlasFunctionsFilterer: AtlasFunctionsFilterer{contract: contract}}, nil
}

// NewAtlasFunctionsCaller creates a new read-only instance of AtlasFunctions, bound to a specific deployed contract.
func NewAtlasFunctionsCaller(address common.Address, caller bind.ContractCaller) (*AtlasFunctionsCaller, error) {
	contract, err := bindAtlasFunctions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsCaller{contract: contract}, nil
}

// NewAtlasFunctionsTransactor creates a new write-only instance of AtlasFunctions, bound to a specific deployed contract.
func NewAtlasFunctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*AtlasFunctionsTransactor, error) {
	contract, err := bindAtlasFunctions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTransactor{contract: contract}, nil
}

// NewAtlasFunctionsFilterer creates a new log filterer instance of AtlasFunctions, bound to a specific deployed contract.
func NewAtlasFunctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*AtlasFunctionsFilterer, error) {
	contract, err := bindAtlasFunctions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsFilterer{contract: contract}, nil
}

// bindAtlasFunctions binds a generic wrapper to an already deployed contract.
func bindAtlasFunctions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtlasFunctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtlasFunctions *AtlasFunctionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtlasFunctions.Contract.AtlasFunctionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtlasFunctions *AtlasFunctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.AtlasFunctionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtlasFunctions *AtlasFunctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.AtlasFunctionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtlasFunctions *AtlasFunctionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtlasFunctions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtlasFunctions *AtlasFunctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtlasFunctions *AtlasFunctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.contract.Transact(opts, method, params...)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctions *AtlasFunctionsTransactor) FireOracleRequest(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctions.contract.Transact(opts, "fireOracleRequest", requestId, subscriptionId, data)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctions *AtlasFunctionsSession) FireOracleRequest(requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireOracleRequest(&_AtlasFunctions.TransactOpts, requestId, subscriptionId, data)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctions *AtlasFunctionsTransactorSession) FireOracleRequest(requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireOracleRequest(&_AtlasFunctions.TransactOpts, requestId, subscriptionId, data)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctions *AtlasFunctionsTransactor) FireOracleResponse(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctions.contract.Transact(opts, "fireOracleResponse", requestId)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctions *AtlasFunctionsSession) FireOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireOracleResponse(&_AtlasFunctions.TransactOpts, requestId)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctions *AtlasFunctionsTransactorSession) FireOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireOracleResponse(&_AtlasFunctions.TransactOpts, requestId)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctions *AtlasFunctionsTransactor) FireSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctions.contract.Transact(opts, "fireSubscriptionFunded", subscriptionId, oldBalance, newBalance)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctions *AtlasFunctionsSession) FireSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireSubscriptionFunded(&_AtlasFunctions.TransactOpts, subscriptionId, oldBalance, newBalance)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctions *AtlasFunctionsTransactorSession) FireSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireSubscriptionFunded(&_AtlasFunctions.TransactOpts, subscriptionId, oldBalance, newBalance)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctions *AtlasFunctionsTransactor) FireUserCallbackError(opts *bind.TransactOpts, requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctions.contract.Transact(opts, "fireUserCallbackError", requestId, reason)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctions *AtlasFunctionsSession) FireUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireUserCallbackError(&_AtlasFunctions.TransactOpts, requestId, reason)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctions *AtlasFunctionsTransactorSession) FireUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireUserCallbackError(&_AtlasFunctions.TransactOpts, requestId, reason)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctions *AtlasFunctionsTransactor) FireUserCallbackRawError(opts *bind.TransactOpts, requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctions.contract.Transact(opts, "fireUserCallbackRawError", requestId, lowLevelData)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctions *AtlasFunctionsSession) FireUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireUserCallbackRawError(&_AtlasFunctions.TransactOpts, requestId, lowLevelData)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctions *AtlasFunctionsTransactorSession) FireUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctions.Contract.FireUserCallbackRawError(&_AtlasFunctions.TransactOpts, requestId, lowLevelData)
}

// AtlasFunctionsOracleRequestIterator is returned from FilterOracleRequest and is used to iterate over the raw logs and unpacked data for OracleRequest events raised by the AtlasFunctions contract.
type AtlasFunctionsOracleRequestIterator struct {
	Event *AtlasFunctionsOracleRequest // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsOracleRequest)
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
		it.Event = new(AtlasFunctionsOracleRequest)
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
func (it *AtlasFunctionsOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsOracleRequest represents a OracleRequest event raised by the AtlasFunctions contract.
type AtlasFunctionsOracleRequest struct {
	RequestId          [32]byte
	RequestingContract common.Address
	RequestInitiator   common.Address
	SubscriptionId     uint64
	SubscriptionOwner  common.Address
	Data               []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOracleRequest is a free log retrieval operation binding the contract event 0xa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c.
//
// Solidity: event OracleRequest(bytes32 indexed requestId, address requestingContract, address requestInitiator, uint64 subscriptionId, address subscriptionOwner, bytes data)
func (_AtlasFunctions *AtlasFunctionsFilterer) FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.FilterLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsOracleRequestIterator{contract: _AtlasFunctions.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

// WatchOracleRequest is a free log subscription operation binding the contract event 0xa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c.
//
// Solidity: event OracleRequest(bytes32 indexed requestId, address requestingContract, address requestInitiator, uint64 subscriptionId, address subscriptionOwner, bytes data)
func (_AtlasFunctions *AtlasFunctionsFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.WatchLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsOracleRequest)
				if err := _AtlasFunctions.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

// ParseOracleRequest is a log parse operation binding the contract event 0xa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c.
//
// Solidity: event OracleRequest(bytes32 indexed requestId, address requestingContract, address requestInitiator, uint64 subscriptionId, address subscriptionOwner, bytes data)
func (_AtlasFunctions *AtlasFunctionsFilterer) ParseOracleRequest(log types.Log) (*AtlasFunctionsOracleRequest, error) {
	event := new(AtlasFunctionsOracleRequest)
	if err := _AtlasFunctions.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsOracleResponseIterator is returned from FilterOracleResponse and is used to iterate over the raw logs and unpacked data for OracleResponse events raised by the AtlasFunctions contract.
type AtlasFunctionsOracleResponseIterator struct {
	Event *AtlasFunctionsOracleResponse // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsOracleResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsOracleResponse)
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
		it.Event = new(AtlasFunctionsOracleResponse)
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
func (it *AtlasFunctionsOracleResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsOracleResponse represents a OracleResponse event raised by the AtlasFunctions contract.
type AtlasFunctionsOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleResponse is a free log retrieval operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_AtlasFunctions *AtlasFunctionsFilterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsOracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsOracleResponseIterator{contract: _AtlasFunctions.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

// WatchOracleResponse is a free log subscription operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_AtlasFunctions *AtlasFunctionsFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsOracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsOracleResponse)
				if err := _AtlasFunctions.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

// ParseOracleResponse is a log parse operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_AtlasFunctions *AtlasFunctionsFilterer) ParseOracleResponse(log types.Log) (*AtlasFunctionsOracleResponse, error) {
	event := new(AtlasFunctionsOracleResponse)
	if err := _AtlasFunctions.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsSubscriptionFundedIterator is returned from FilterSubscriptionFunded and is used to iterate over the raw logs and unpacked data for SubscriptionFunded events raised by the AtlasFunctions contract.
type AtlasFunctionsSubscriptionFundedIterator struct {
	Event *AtlasFunctionsSubscriptionFunded // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsSubscriptionFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsSubscriptionFunded)
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
		it.Event = new(AtlasFunctionsSubscriptionFunded)
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
func (it *AtlasFunctionsSubscriptionFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsSubscriptionFunded represents a SubscriptionFunded event raised by the AtlasFunctions contract.
type AtlasFunctionsSubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionFunded is a free log retrieval operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance)
func (_AtlasFunctions *AtlasFunctionsFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*AtlasFunctionsSubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsSubscriptionFundedIterator{contract: _AtlasFunctions.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

// WatchSubscriptionFunded is a free log subscription operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance)
func (_AtlasFunctions *AtlasFunctionsFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsSubscriptionFunded)
				if err := _AtlasFunctions.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

// ParseSubscriptionFunded is a log parse operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance)
func (_AtlasFunctions *AtlasFunctionsFilterer) ParseSubscriptionFunded(log types.Log) (*AtlasFunctionsSubscriptionFunded, error) {
	event := new(AtlasFunctionsSubscriptionFunded)
	if err := _AtlasFunctions.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsUserCallbackErrorIterator is returned from FilterUserCallbackError and is used to iterate over the raw logs and unpacked data for UserCallbackError events raised by the AtlasFunctions contract.
type AtlasFunctionsUserCallbackErrorIterator struct {
	Event *AtlasFunctionsUserCallbackError // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsUserCallbackErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsUserCallbackError)
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
		it.Event = new(AtlasFunctionsUserCallbackError)
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
func (it *AtlasFunctionsUserCallbackErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsUserCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsUserCallbackError represents a UserCallbackError event raised by the AtlasFunctions contract.
type AtlasFunctionsUserCallbackError struct {
	RequestId [32]byte
	Reason    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUserCallbackError is a free log retrieval operation binding the contract event 0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c.
//
// Solidity: event UserCallbackError(bytes32 indexed requestId, string reason)
func (_AtlasFunctions *AtlasFunctionsFilterer) FilterUserCallbackError(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsUserCallbackErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.FilterLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsUserCallbackErrorIterator{contract: _AtlasFunctions.contract, event: "UserCallbackError", logs: logs, sub: sub}, nil
}

// WatchUserCallbackError is a free log subscription operation binding the contract event 0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c.
//
// Solidity: event UserCallbackError(bytes32 indexed requestId, string reason)
func (_AtlasFunctions *AtlasFunctionsFilterer) WatchUserCallbackError(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsUserCallbackError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.WatchLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsUserCallbackError)
				if err := _AtlasFunctions.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
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

// ParseUserCallbackError is a log parse operation binding the contract event 0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c.
//
// Solidity: event UserCallbackError(bytes32 indexed requestId, string reason)
func (_AtlasFunctions *AtlasFunctionsFilterer) ParseUserCallbackError(log types.Log) (*AtlasFunctionsUserCallbackError, error) {
	event := new(AtlasFunctionsUserCallbackError)
	if err := _AtlasFunctions.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsUserCallbackRawErrorIterator is returned from FilterUserCallbackRawError and is used to iterate over the raw logs and unpacked data for UserCallbackRawError events raised by the AtlasFunctions contract.
type AtlasFunctionsUserCallbackRawErrorIterator struct {
	Event *AtlasFunctionsUserCallbackRawError // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsUserCallbackRawErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsUserCallbackRawError)
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
		it.Event = new(AtlasFunctionsUserCallbackRawError)
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
func (it *AtlasFunctionsUserCallbackRawErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsUserCallbackRawErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsUserCallbackRawError represents a UserCallbackRawError event raised by the AtlasFunctions contract.
type AtlasFunctionsUserCallbackRawError struct {
	RequestId    [32]byte
	LowLevelData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUserCallbackRawError is a free log retrieval operation binding the contract event 0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2.
//
// Solidity: event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData)
func (_AtlasFunctions *AtlasFunctionsFilterer) FilterUserCallbackRawError(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsUserCallbackRawErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.FilterLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsUserCallbackRawErrorIterator{contract: _AtlasFunctions.contract, event: "UserCallbackRawError", logs: logs, sub: sub}, nil
}

// WatchUserCallbackRawError is a free log subscription operation binding the contract event 0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2.
//
// Solidity: event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData)
func (_AtlasFunctions *AtlasFunctionsFilterer) WatchUserCallbackRawError(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsUserCallbackRawError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctions.contract.WatchLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsUserCallbackRawError)
				if err := _AtlasFunctions.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
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

// ParseUserCallbackRawError is a log parse operation binding the contract event 0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2.
//
// Solidity: event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData)
func (_AtlasFunctions *AtlasFunctionsFilterer) ParseUserCallbackRawError(log types.Log) (*AtlasFunctionsUserCallbackRawError, error) {
	event := new(AtlasFunctionsUserCallbackRawError)
	if err := _AtlasFunctions.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
