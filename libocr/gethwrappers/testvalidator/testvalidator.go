// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testvalidator

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

// AggregatorValidatorInterfaceMetaData contains all meta data concerning the AggregatorValidatorInterface contract.
var AggregatorValidatorInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AggregatorValidatorInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use AggregatorValidatorInterfaceMetaData.ABI instead.
var AggregatorValidatorInterfaceABI = AggregatorValidatorInterfaceMetaData.ABI

// AggregatorValidatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorValidatorInterface struct {
	AggregatorValidatorInterfaceCaller     // Read-only binding to the contract
	AggregatorValidatorInterfaceTransactor // Write-only binding to the contract
	AggregatorValidatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorValidatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorValidatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorValidatorInterfaceSession struct {
	Contract     *AggregatorValidatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorValidatorInterfaceCallerSession struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// AggregatorValidatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorValidatorInterfaceTransactorSession struct {
	Contract     *AggregatorValidatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceRaw struct {
	Contract *AggregatorValidatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCallerRaw struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactorRaw struct {
	Contract *AggregatorValidatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorValidatorInterface creates a new instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorValidatorInterface, error) {
	contract, err := bindAggregatorValidatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterface{AggregatorValidatorInterfaceCaller: AggregatorValidatorInterfaceCaller{contract: contract}, AggregatorValidatorInterfaceTransactor: AggregatorValidatorInterfaceTransactor{contract: contract}, AggregatorValidatorInterfaceFilterer: AggregatorValidatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorValidatorInterfaceCaller creates a new read-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorValidatorInterfaceCaller, error) {
	contract, err := bindAggregatorValidatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceTransactor creates a new write-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorValidatorInterfaceTransactor, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceFilterer creates a new log filterer instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorValidatorInterfaceFilterer, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorValidatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorValidatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorValidatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transact(opts, method, params...)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// TestValidatorMetaData contains all meta data concerning the TestValidator contract.
var TestValidatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"}],\"name\":\"Validated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestRoundId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"minGasUse\",\"type\":\"uint32\"}],\"name\":\"setMinGasUse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610178806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806311a8f41314610046578063beed9b5114610060578063c4792df4146100a3575b600080fd5b61004e6100c8565b60408051918252519081900360200190f35b61008f6004803603608081101561007657600080fd5b50803590602081013590604081013590606001356100ce565b604080519115158252519081900360200190f35b6100c6600480360360208110156100b957600080fd5b503563ffffffff1661014f565b005b60015490565b6000805a6040805188815260208101889052808201879052606081018690526080810183905290519192507fdb623f4f39d41e75ae1cbe50460c3d1496b6cf9a0db391b7197f82cab2744d21919081900360a00190a1600184905560005463ffffffff165b805a8303101561014257610133565b5060019695505050505050565b6000805463ffffffff191663ffffffff9290921691909117905556fea164736f6c6343000706000a",
}

// TestValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use TestValidatorMetaData.ABI instead.
var TestValidatorABI = TestValidatorMetaData.ABI

// TestValidatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestValidatorMetaData.Bin instead.
var TestValidatorBin = TestValidatorMetaData.Bin

// DeployTestValidator deploys a new Ethereum contract, binding an instance of TestValidator to it.
func DeployTestValidator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestValidator, error) {
	parsed, err := TestValidatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestValidatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestValidator{TestValidatorCaller: TestValidatorCaller{contract: contract}, TestValidatorTransactor: TestValidatorTransactor{contract: contract}, TestValidatorFilterer: TestValidatorFilterer{contract: contract}}, nil
}

// TestValidator is an auto generated Go binding around an Ethereum contract.
type TestValidator struct {
	TestValidatorCaller     // Read-only binding to the contract
	TestValidatorTransactor // Write-only binding to the contract
	TestValidatorFilterer   // Log filterer for contract events
}

// TestValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestValidatorSession struct {
	Contract     *TestValidator    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestValidatorCallerSession struct {
	Contract *TestValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TestValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestValidatorTransactorSession struct {
	Contract     *TestValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestValidatorRaw struct {
	Contract *TestValidator // Generic contract binding to access the raw methods on
}

// TestValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestValidatorCallerRaw struct {
	Contract *TestValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// TestValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestValidatorTransactorRaw struct {
	Contract *TestValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestValidator creates a new instance of TestValidator, bound to a specific deployed contract.
func NewTestValidator(address common.Address, backend bind.ContractBackend) (*TestValidator, error) {
	contract, err := bindTestValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestValidator{TestValidatorCaller: TestValidatorCaller{contract: contract}, TestValidatorTransactor: TestValidatorTransactor{contract: contract}, TestValidatorFilterer: TestValidatorFilterer{contract: contract}}, nil
}

// NewTestValidatorCaller creates a new read-only instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorCaller(address common.Address, caller bind.ContractCaller) (*TestValidatorCaller, error) {
	contract, err := bindTestValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorCaller{contract: contract}, nil
}

// NewTestValidatorTransactor creates a new write-only instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestValidatorTransactor, error) {
	contract, err := bindTestValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorTransactor{contract: contract}, nil
}

// NewTestValidatorFilterer creates a new log filterer instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestValidatorFilterer, error) {
	contract, err := bindTestValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestValidatorFilterer{contract: contract}, nil
}

// bindTestValidator binds a generic wrapper to an already deployed contract.
func bindTestValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestValidator *TestValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.TestValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestValidator *TestValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestValidator *TestValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestValidator *TestValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestValidator *TestValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestValidator *TestValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transact(opts, method, params...)
}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint256)
func (_TestValidator *TestValidatorCaller) LatestRoundId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestValidator.contract.Call(opts, &out, "latestRoundId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint256)
func (_TestValidator *TestValidatorSession) LatestRoundId() (*big.Int, error) {
	return _TestValidator.Contract.LatestRoundId(&_TestValidator.CallOpts)
}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint256)
func (_TestValidator *TestValidatorCallerSession) LatestRoundId() (*big.Int, error) {
	return _TestValidator.Contract.LatestRoundId(&_TestValidator.CallOpts)
}

// SetMinGasUse is a paid mutator transaction binding the contract method 0xc4792df4.
//
// Solidity: function setMinGasUse(uint32 minGasUse) returns()
func (_TestValidator *TestValidatorTransactor) SetMinGasUse(opts *bind.TransactOpts, minGasUse uint32) (*types.Transaction, error) {
	return _TestValidator.contract.Transact(opts, "setMinGasUse", minGasUse)
}

// SetMinGasUse is a paid mutator transaction binding the contract method 0xc4792df4.
//
// Solidity: function setMinGasUse(uint32 minGasUse) returns()
func (_TestValidator *TestValidatorSession) SetMinGasUse(minGasUse uint32) (*types.Transaction, error) {
	return _TestValidator.Contract.SetMinGasUse(&_TestValidator.TransactOpts, minGasUse)
}

// SetMinGasUse is a paid mutator transaction binding the contract method 0xc4792df4.
//
// Solidity: function setMinGasUse(uint32 minGasUse) returns()
func (_TestValidator *TestValidatorTransactorSession) SetMinGasUse(minGasUse uint32) (*types.Transaction, error) {
	return _TestValidator.Contract.SetMinGasUse(&_TestValidator.TransactOpts, minGasUse)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_TestValidator *TestValidatorTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _TestValidator.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_TestValidator *TestValidatorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_TestValidator *TestValidatorTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// TestValidatorValidatedIterator is returned from FilterValidated and is used to iterate over the raw logs and unpacked data for Validated events raised by the TestValidator contract.
type TestValidatorValidatedIterator struct {
	Event *TestValidatorValidated // Event containing the contract specifics and raw log

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
func (it *TestValidatorValidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestValidatorValidated)
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
		it.Event = new(TestValidatorValidated)
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
func (it *TestValidatorValidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestValidatorValidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestValidatorValidated represents a Validated event raised by the TestValidator contract.
type TestValidatorValidated struct {
	PreviousRoundId *big.Int
	PreviousAnswer  *big.Int
	CurrentRoundId  *big.Int
	CurrentAnswer   *big.Int
	InitialGas      *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidated is a free log retrieval operation binding the contract event 0xdb623f4f39d41e75ae1cbe50460c3d1496b6cf9a0db391b7197f82cab2744d21.
//
// Solidity: event Validated(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer, uint256 initialGas)
func (_TestValidator *TestValidatorFilterer) FilterValidated(opts *bind.FilterOpts) (*TestValidatorValidatedIterator, error) {

	logs, sub, err := _TestValidator.contract.FilterLogs(opts, "Validated")
	if err != nil {
		return nil, err
	}
	return &TestValidatorValidatedIterator{contract: _TestValidator.contract, event: "Validated", logs: logs, sub: sub}, nil
}

// WatchValidated is a free log subscription operation binding the contract event 0xdb623f4f39d41e75ae1cbe50460c3d1496b6cf9a0db391b7197f82cab2744d21.
//
// Solidity: event Validated(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer, uint256 initialGas)
func (_TestValidator *TestValidatorFilterer) WatchValidated(opts *bind.WatchOpts, sink chan<- *TestValidatorValidated) (event.Subscription, error) {

	logs, sub, err := _TestValidator.contract.WatchLogs(opts, "Validated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestValidatorValidated)
				if err := _TestValidator.contract.UnpackLog(event, "Validated", log); err != nil {
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

// ParseValidated is a log parse operation binding the contract event 0xdb623f4f39d41e75ae1cbe50460c3d1496b6cf9a0db391b7197f82cab2744d21.
//
// Solidity: event Validated(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer, uint256 initialGas)
func (_TestValidator *TestValidatorFilterer) ParseValidated(log types.Log) (*TestValidatorValidated, error) {
	event := new(TestValidatorValidated)
	if err := _TestValidator.contract.UnpackLog(event, "Validated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
