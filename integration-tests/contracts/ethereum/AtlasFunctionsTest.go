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

// AtlasFunctionsTestMetaData contains all meta data concerning the AtlasFunctionsTest contract.
var AtlasFunctionsTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestingContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requestInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subscriptionOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"UserCallbackError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"UserCallbackRawError\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"fireOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"fireOracleResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"fireSubscriptionFunded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"fireUserCallbackError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"lowLevelData\",\"type\":\"bytes\"}],\"name\":\"fireUserCallbackRawError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104c7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806338ddd89c1461005c5780634bc2065f146100715780635fa55d55146100845780637c89a07f146100975780638b2d730b146100aa575b600080fd5b61006f61006a3660046102ab565b6100bd565b005b61006f61007f366004610301565b610102565b61006f61009236600461031a565b610130565b61006f6100a536600461034d565b610171565b61006f6100b8366004610393565b6101ad565b827fa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c30338533866040516100f5959493929190610430565b60405180910390a2505050565b60405181907f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a6490600090a250565b60408051838152602081018390526001600160401b038516917fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f891016100f5565b817fe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2826040516101a1919061047e565b60405180910390a25050565b817fb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c826040516101a1919061047e565b80356001600160401b03811681146101f457600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b0380841115610229576102296101f9565b604051601f8501601f19908116603f01168101908282118183101715610251576102516101f9565b8160405280935085815286868601111561026a57600080fd5b858560208301376000602087830101525050509392505050565b600082601f83011261029557600080fd5b6102a48383356020850161020f565b9392505050565b6000806000606084860312156102c057600080fd5b833592506102d0602085016101dd565b915060408401356001600160401b038111156102eb57600080fd5b6102f786828701610284565b9150509250925092565b60006020828403121561031357600080fd5b5035919050565b60008060006060848603121561032f57600080fd5b610338846101dd565b95602085013595506040909401359392505050565b6000806040838503121561036057600080fd5b8235915060208301356001600160401b0381111561037d57600080fd5b61038985828601610284565b9150509250929050565b600080604083850312156103a657600080fd5b8235915060208301356001600160401b038111156103c357600080fd5b8301601f810185136103d457600080fd5b6103898582356020840161020f565b6000815180845260005b81811015610409576020818501810151868301820152016103ed565b8181111561041b576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b03868116825285811660208301526001600160401b03851660408301528316606082015260a060808201819052600090610473908301846103e3565b979650505050505050565b6020815260006102a460208301846103e356fea26469706673582212203ae7db7befe2368a32b34d8bc9445b13d2c1a6068bf4b715eeda88abb511fd7e64736f6c634300080d0033",
}

// AtlasFunctionsTestABI is the input ABI used to generate the binding from.
// Deprecated: Use AtlasFunctionsTestMetaData.ABI instead.
var AtlasFunctionsTestABI = AtlasFunctionsTestMetaData.ABI

// AtlasFunctionsTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AtlasFunctionsTestMetaData.Bin instead.
var AtlasFunctionsTestBin = AtlasFunctionsTestMetaData.Bin

// DeployAtlasFunctionsTest deploys a new Ethereum contract, binding an instance of AtlasFunctionsTest to it.
func DeployAtlasFunctionsTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtlasFunctionsTest, error) {
	parsed, err := AtlasFunctionsTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AtlasFunctionsTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtlasFunctionsTest{AtlasFunctionsTestCaller: AtlasFunctionsTestCaller{contract: contract}, AtlasFunctionsTestTransactor: AtlasFunctionsTestTransactor{contract: contract}, AtlasFunctionsTestFilterer: AtlasFunctionsTestFilterer{contract: contract}}, nil
}

// AtlasFunctionsTest is an auto generated Go binding around an Ethereum contract.
type AtlasFunctionsTest struct {
	AtlasFunctionsTestCaller     // Read-only binding to the contract
	AtlasFunctionsTestTransactor // Write-only binding to the contract
	AtlasFunctionsTestFilterer   // Log filterer for contract events
}

// AtlasFunctionsTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtlasFunctionsTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtlasFunctionsTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtlasFunctionsTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtlasFunctionsTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtlasFunctionsTestSession struct {
	Contract     *AtlasFunctionsTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AtlasFunctionsTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtlasFunctionsTestCallerSession struct {
	Contract *AtlasFunctionsTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// AtlasFunctionsTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtlasFunctionsTestTransactorSession struct {
	Contract     *AtlasFunctionsTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AtlasFunctionsTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtlasFunctionsTestRaw struct {
	Contract *AtlasFunctionsTest // Generic contract binding to access the raw methods on
}

// AtlasFunctionsTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtlasFunctionsTestCallerRaw struct {
	Contract *AtlasFunctionsTestCaller // Generic read-only contract binding to access the raw methods on
}

// AtlasFunctionsTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtlasFunctionsTestTransactorRaw struct {
	Contract *AtlasFunctionsTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtlasFunctionsTest creates a new instance of AtlasFunctionsTest, bound to a specific deployed contract.
func NewAtlasFunctionsTest(address common.Address, backend bind.ContractBackend) (*AtlasFunctionsTest, error) {
	contract, err := bindAtlasFunctionsTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTest{AtlasFunctionsTestCaller: AtlasFunctionsTestCaller{contract: contract}, AtlasFunctionsTestTransactor: AtlasFunctionsTestTransactor{contract: contract}, AtlasFunctionsTestFilterer: AtlasFunctionsTestFilterer{contract: contract}}, nil
}

// NewAtlasFunctionsTestCaller creates a new read-only instance of AtlasFunctionsTest, bound to a specific deployed contract.
func NewAtlasFunctionsTestCaller(address common.Address, caller bind.ContractCaller) (*AtlasFunctionsTestCaller, error) {
	contract, err := bindAtlasFunctionsTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestCaller{contract: contract}, nil
}

// NewAtlasFunctionsTestTransactor creates a new write-only instance of AtlasFunctionsTest, bound to a specific deployed contract.
func NewAtlasFunctionsTestTransactor(address common.Address, transactor bind.ContractTransactor) (*AtlasFunctionsTestTransactor, error) {
	contract, err := bindAtlasFunctionsTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestTransactor{contract: contract}, nil
}

// NewAtlasFunctionsTestFilterer creates a new log filterer instance of AtlasFunctionsTest, bound to a specific deployed contract.
func NewAtlasFunctionsTestFilterer(address common.Address, filterer bind.ContractFilterer) (*AtlasFunctionsTestFilterer, error) {
	contract, err := bindAtlasFunctionsTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestFilterer{contract: contract}, nil
}

// bindAtlasFunctionsTest binds a generic wrapper to an already deployed contract.
func bindAtlasFunctionsTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtlasFunctionsTestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtlasFunctionsTest *AtlasFunctionsTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtlasFunctionsTest.Contract.AtlasFunctionsTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtlasFunctionsTest *AtlasFunctionsTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.AtlasFunctionsTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtlasFunctionsTest *AtlasFunctionsTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.AtlasFunctionsTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtlasFunctionsTest *AtlasFunctionsTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtlasFunctionsTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.contract.Transact(opts, method, params...)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactor) FireOracleRequest(opts *bind.TransactOpts, requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.contract.Transact(opts, "fireOracleRequest", requestId, subscriptionId, data)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestSession) FireOracleRequest(requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireOracleRequest(&_AtlasFunctionsTest.TransactOpts, requestId, subscriptionId, data)
}

// FireOracleRequest is a paid mutator transaction binding the contract method 0x38ddd89c.
//
// Solidity: function fireOracleRequest(bytes32 requestId, uint64 subscriptionId, bytes data) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorSession) FireOracleRequest(requestId [32]byte, subscriptionId uint64, data []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireOracleRequest(&_AtlasFunctionsTest.TransactOpts, requestId, subscriptionId, data)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactor) FireOracleResponse(opts *bind.TransactOpts, requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.contract.Transact(opts, "fireOracleResponse", requestId)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestSession) FireOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireOracleResponse(&_AtlasFunctionsTest.TransactOpts, requestId)
}

// FireOracleResponse is a paid mutator transaction binding the contract method 0x4bc2065f.
//
// Solidity: function fireOracleResponse(bytes32 requestId) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorSession) FireOracleResponse(requestId [32]byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireOracleResponse(&_AtlasFunctionsTest.TransactOpts, requestId)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactor) FireSubscriptionFunded(opts *bind.TransactOpts, subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctionsTest.contract.Transact(opts, "fireSubscriptionFunded", subscriptionId, oldBalance, newBalance)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestSession) FireSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireSubscriptionFunded(&_AtlasFunctionsTest.TransactOpts, subscriptionId, oldBalance, newBalance)
}

// FireSubscriptionFunded is a paid mutator transaction binding the contract method 0x5fa55d55.
//
// Solidity: function fireSubscriptionFunded(uint64 subscriptionId, uint256 oldBalance, uint256 newBalance) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorSession) FireSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireSubscriptionFunded(&_AtlasFunctionsTest.TransactOpts, subscriptionId, oldBalance, newBalance)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactor) FireUserCallbackError(opts *bind.TransactOpts, requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctionsTest.contract.Transact(opts, "fireUserCallbackError", requestId, reason)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestSession) FireUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireUserCallbackError(&_AtlasFunctionsTest.TransactOpts, requestId, reason)
}

// FireUserCallbackError is a paid mutator transaction binding the contract method 0x8b2d730b.
//
// Solidity: function fireUserCallbackError(bytes32 requestId, string reason) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorSession) FireUserCallbackError(requestId [32]byte, reason string) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireUserCallbackError(&_AtlasFunctionsTest.TransactOpts, requestId, reason)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactor) FireUserCallbackRawError(opts *bind.TransactOpts, requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.contract.Transact(opts, "fireUserCallbackRawError", requestId, lowLevelData)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestSession) FireUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireUserCallbackRawError(&_AtlasFunctionsTest.TransactOpts, requestId, lowLevelData)
}

// FireUserCallbackRawError is a paid mutator transaction binding the contract method 0x7c89a07f.
//
// Solidity: function fireUserCallbackRawError(bytes32 requestId, bytes lowLevelData) returns()
func (_AtlasFunctionsTest *AtlasFunctionsTestTransactorSession) FireUserCallbackRawError(requestId [32]byte, lowLevelData []byte) (*types.Transaction, error) {
	return _AtlasFunctionsTest.Contract.FireUserCallbackRawError(&_AtlasFunctionsTest.TransactOpts, requestId, lowLevelData)
}

// AtlasFunctionsTestOracleRequestIterator is returned from FilterOracleRequest and is used to iterate over the raw logs and unpacked data for OracleRequest events raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestOracleRequestIterator struct {
	Event *AtlasFunctionsTestOracleRequest // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsTestOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsTestOracleRequest)
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
		it.Event = new(AtlasFunctionsTestOracleRequest)
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
func (it *AtlasFunctionsTestOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsTestOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsTestOracleRequest represents a OracleRequest event raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestOracleRequest struct {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) FilterOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsTestOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.FilterLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestOracleRequestIterator{contract: _AtlasFunctionsTest.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

// WatchOracleRequest is a free log subscription operation binding the contract event 0xa1ec73989d79578cd6f67d4f593ac3e0a4d1020e5c0164db52108d7ff785406c.
//
// Solidity: event OracleRequest(bytes32 indexed requestId, address requestingContract, address requestInitiator, uint64 subscriptionId, address subscriptionOwner, bytes data)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsTestOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.WatchLogs(opts, "OracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsTestOracleRequest)
				if err := _AtlasFunctionsTest.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) ParseOracleRequest(log types.Log) (*AtlasFunctionsTestOracleRequest, error) {
	event := new(AtlasFunctionsTestOracleRequest)
	if err := _AtlasFunctionsTest.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsTestOracleResponseIterator is returned from FilterOracleResponse and is used to iterate over the raw logs and unpacked data for OracleResponse events raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestOracleResponseIterator struct {
	Event *AtlasFunctionsTestOracleResponse // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsTestOracleResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsTestOracleResponse)
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
		it.Event = new(AtlasFunctionsTestOracleResponse)
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
func (it *AtlasFunctionsTestOracleResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsTestOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsTestOracleResponse represents a OracleResponse event raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleResponse is a free log retrieval operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) FilterOracleResponse(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsTestOracleResponseIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.FilterLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestOracleResponseIterator{contract: _AtlasFunctionsTest.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

// WatchOracleResponse is a free log subscription operation binding the contract event 0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64.
//
// Solidity: event OracleResponse(bytes32 indexed requestId)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsTestOracleResponse, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.WatchLogs(opts, "OracleResponse", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsTestOracleResponse)
				if err := _AtlasFunctionsTest.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) ParseOracleResponse(log types.Log) (*AtlasFunctionsTestOracleResponse, error) {
	event := new(AtlasFunctionsTestOracleResponse)
	if err := _AtlasFunctionsTest.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsTestSubscriptionFundedIterator is returned from FilterSubscriptionFunded and is used to iterate over the raw logs and unpacked data for SubscriptionFunded events raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestSubscriptionFundedIterator struct {
	Event *AtlasFunctionsTestSubscriptionFunded // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsTestSubscriptionFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsTestSubscriptionFunded)
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
		it.Event = new(AtlasFunctionsTestSubscriptionFunded)
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
func (it *AtlasFunctionsTestSubscriptionFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsTestSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsTestSubscriptionFunded represents a SubscriptionFunded event raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestSubscriptionFunded struct {
	SubscriptionId uint64
	OldBalance     *big.Int
	NewBalance     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionFunded is a free log retrieval operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subscriptionId []uint64) (*AtlasFunctionsTestSubscriptionFundedIterator, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.FilterLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestSubscriptionFundedIterator{contract: _AtlasFunctionsTest.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

// WatchSubscriptionFunded is a free log subscription operation binding the contract event 0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8.
//
// Solidity: event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsTestSubscriptionFunded, subscriptionId []uint64) (event.Subscription, error) {

	var subscriptionIdRule []interface{}
	for _, subscriptionIdItem := range subscriptionId {
		subscriptionIdRule = append(subscriptionIdRule, subscriptionIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.WatchLogs(opts, "SubscriptionFunded", subscriptionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsTestSubscriptionFunded)
				if err := _AtlasFunctionsTest.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) ParseSubscriptionFunded(log types.Log) (*AtlasFunctionsTestSubscriptionFunded, error) {
	event := new(AtlasFunctionsTestSubscriptionFunded)
	if err := _AtlasFunctionsTest.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsTestUserCallbackErrorIterator is returned from FilterUserCallbackError and is used to iterate over the raw logs and unpacked data for UserCallbackError events raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestUserCallbackErrorIterator struct {
	Event *AtlasFunctionsTestUserCallbackError // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsTestUserCallbackErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsTestUserCallbackError)
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
		it.Event = new(AtlasFunctionsTestUserCallbackError)
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
func (it *AtlasFunctionsTestUserCallbackErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsTestUserCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsTestUserCallbackError represents a UserCallbackError event raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestUserCallbackError struct {
	RequestId [32]byte
	Reason    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUserCallbackError is a free log retrieval operation binding the contract event 0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c.
//
// Solidity: event UserCallbackError(bytes32 indexed requestId, string reason)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) FilterUserCallbackError(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsTestUserCallbackErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.FilterLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestUserCallbackErrorIterator{contract: _AtlasFunctionsTest.contract, event: "UserCallbackError", logs: logs, sub: sub}, nil
}

// WatchUserCallbackError is a free log subscription operation binding the contract event 0xb2931868c372fe17a25643458add467d60ec5c51125a99b7309f41f5bcd2da6c.
//
// Solidity: event UserCallbackError(bytes32 indexed requestId, string reason)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) WatchUserCallbackError(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsTestUserCallbackError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.WatchLogs(opts, "UserCallbackError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsTestUserCallbackError)
				if err := _AtlasFunctionsTest.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) ParseUserCallbackError(log types.Log) (*AtlasFunctionsTestUserCallbackError, error) {
	event := new(AtlasFunctionsTestUserCallbackError)
	if err := _AtlasFunctionsTest.contract.UnpackLog(event, "UserCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtlasFunctionsTestUserCallbackRawErrorIterator is returned from FilterUserCallbackRawError and is used to iterate over the raw logs and unpacked data for UserCallbackRawError events raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestUserCallbackRawErrorIterator struct {
	Event *AtlasFunctionsTestUserCallbackRawError // Event containing the contract specifics and raw log

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
func (it *AtlasFunctionsTestUserCallbackRawErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtlasFunctionsTestUserCallbackRawError)
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
		it.Event = new(AtlasFunctionsTestUserCallbackRawError)
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
func (it *AtlasFunctionsTestUserCallbackRawErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtlasFunctionsTestUserCallbackRawErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtlasFunctionsTestUserCallbackRawError represents a UserCallbackRawError event raised by the AtlasFunctionsTest contract.
type AtlasFunctionsTestUserCallbackRawError struct {
	RequestId    [32]byte
	LowLevelData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUserCallbackRawError is a free log retrieval operation binding the contract event 0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2.
//
// Solidity: event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) FilterUserCallbackRawError(opts *bind.FilterOpts, requestId [][32]byte) (*AtlasFunctionsTestUserCallbackRawErrorIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.FilterLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &AtlasFunctionsTestUserCallbackRawErrorIterator{contract: _AtlasFunctionsTest.contract, event: "UserCallbackRawError", logs: logs, sub: sub}, nil
}

// WatchUserCallbackRawError is a free log subscription operation binding the contract event 0xe0b838ffe6ee22a0d3acf19a85db6a41b34a1ab739e2d6c759a2e42d95bdccb2.
//
// Solidity: event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData)
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) WatchUserCallbackRawError(opts *bind.WatchOpts, sink chan<- *AtlasFunctionsTestUserCallbackRawError, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _AtlasFunctionsTest.contract.WatchLogs(opts, "UserCallbackRawError", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtlasFunctionsTestUserCallbackRawError)
				if err := _AtlasFunctionsTest.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
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
func (_AtlasFunctionsTest *AtlasFunctionsTestFilterer) ParseUserCallbackRawError(log types.Log) (*AtlasFunctionsTestUserCallbackRawError, error) {
	event := new(AtlasFunctionsTestUserCallbackRawError)
	if err := _AtlasFunctionsTest.contract.UnpackLog(event, "UserCallbackRawError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
