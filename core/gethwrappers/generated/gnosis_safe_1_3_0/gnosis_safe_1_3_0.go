// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gnosis_safe_1_3_0

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
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

// GnosisSafe130MetaData contains all meta data concerning the GnosisSafe130 contract.
var GnosisSafe130MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"approvedHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ApproveHash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"handler\",\"type\":\"address\"}],\"name\":\"ChangedFallbackHandler\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"guard\",\"type\":\"address\"}],\"name\":\"ChangedGuard\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"ChangedThreshold\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"DisabledModule\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"EnabledModule\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"ExecutionFailure\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"ExecutionFromModuleFailure\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"ExecutionFromModuleSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"ExecutionSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"owners\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initializer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fallbackHandler\",\"type\":\"address\"}],\"name\":\"SafeSetup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"msgHash\",\"type\":\"bytes32\"}],\"name\":\"SignMsg\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"addOwnerWithThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashToApprove\",\"type\":\"bytes32\"}],\"name\":\"approveHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"approvedHashes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"requiredSignatures\",\"type\":\"uint256\"}],\"name\":\"checkNSignatures\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"checkSignatures\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"prevModule\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"disableModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"domainSeparator\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"enableModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"safeTxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"gasToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refundReceiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"encodeTransactionData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"safeTxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"gasToken\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"refundReceiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"execTransaction\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"}],\"name\":\"execTransactionFromModule\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"}],\"name\":\"execTransactionFromModuleReturnData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"start\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pageSize\",\"type\":\"uint256\"}],\"name\":\"getModulesPaginated\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"array\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"next\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"getStorageAt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"safeTxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"gasToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refundReceiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"getTransactionHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"isModuleEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"prevOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"}],\"name\":\"requiredTxGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"handler\",\"type\":\"address\"}],\"name\":\"setFallbackHandler\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"guard\",\"type\":\"address\"}],\"name\":\"setGuard\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_owners\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"fallbackHandler\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"paymentReceiver\",\"type\":\"address\"}],\"name\":\"setup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"signedMessages\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetContract\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"calldataPayload\",\"type\":\"bytes\"}],\"name\":\"simulateAndRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"prevOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"swapOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// GnosisSafe130ABI is the input ABI used to generate the binding from.
// Deprecated: Use GnosisSafe130MetaData.ABI instead.
var GnosisSafe130ABI = GnosisSafe130MetaData.ABI

// GnosisSafe130 is an auto generated Go binding around an Ethereum contract.
type GnosisSafe130 struct {
	GnosisSafe130Caller     // Read-only binding to the contract
	GnosisSafe130Transactor // Write-only binding to the contract
	GnosisSafe130Filterer   // Log filterer for contract events
}

// GnosisSafe130Caller is an auto generated read-only Go binding around an Ethereum contract.
type GnosisSafe130Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GnosisSafe130Transactor is an auto generated write-only Go binding around an Ethereum contract.
type GnosisSafe130Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GnosisSafe130Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GnosisSafe130Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GnosisSafe130Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GnosisSafe130Session struct {
	Contract     *GnosisSafe130    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GnosisSafe130CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GnosisSafe130CallerSession struct {
	Contract *GnosisSafe130Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GnosisSafe130TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GnosisSafe130TransactorSession struct {
	Contract     *GnosisSafe130Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GnosisSafe130Raw is an auto generated low-level Go binding around an Ethereum contract.
type GnosisSafe130Raw struct {
	Contract *GnosisSafe130 // Generic contract binding to access the raw methods on
}

// GnosisSafe130CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GnosisSafe130CallerRaw struct {
	Contract *GnosisSafe130Caller // Generic read-only contract binding to access the raw methods on
}

// GnosisSafe130TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GnosisSafe130TransactorRaw struct {
	Contract *GnosisSafe130Transactor // Generic write-only contract binding to access the raw methods on
}

// NewGnosisSafe130 creates a new instance of GnosisSafe130, bound to a specific deployed contract.
func NewGnosisSafe130(address common.Address, backend bind.ContractBackend) (*GnosisSafe130, error) {
	contract, err := bindGnosisSafe130(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130{GnosisSafe130Caller: GnosisSafe130Caller{contract: contract}, GnosisSafe130Transactor: GnosisSafe130Transactor{contract: contract}, GnosisSafe130Filterer: GnosisSafe130Filterer{contract: contract}}, nil
}

// NewGnosisSafe130Caller creates a new read-only instance of GnosisSafe130, bound to a specific deployed contract.
func NewGnosisSafe130Caller(address common.Address, caller bind.ContractCaller) (*GnosisSafe130Caller, error) {
	contract, err := bindGnosisSafe130(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130Caller{contract: contract}, nil
}

// NewGnosisSafe130Transactor creates a new write-only instance of GnosisSafe130, bound to a specific deployed contract.
func NewGnosisSafe130Transactor(address common.Address, transactor bind.ContractTransactor) (*GnosisSafe130Transactor, error) {
	contract, err := bindGnosisSafe130(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130Transactor{contract: contract}, nil
}

// NewGnosisSafe130Filterer creates a new log filterer instance of GnosisSafe130, bound to a specific deployed contract.
func NewGnosisSafe130Filterer(address common.Address, filterer bind.ContractFilterer) (*GnosisSafe130Filterer, error) {
	contract, err := bindGnosisSafe130(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130Filterer{contract: contract}, nil
}

// bindGnosisSafe130 binds a generic wrapper to an already deployed contract.
func bindGnosisSafe130(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GnosisSafe130ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GnosisSafe130 *GnosisSafe130Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GnosisSafe130.Contract.GnosisSafe130Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GnosisSafe130 *GnosisSafe130Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.GnosisSafe130Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GnosisSafe130 *GnosisSafe130Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.GnosisSafe130Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GnosisSafe130 *GnosisSafe130CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GnosisSafe130.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GnosisSafe130 *GnosisSafe130TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GnosisSafe130 *GnosisSafe130TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_GnosisSafe130 *GnosisSafe130Caller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_GnosisSafe130 *GnosisSafe130Session) VERSION() (string, error) {
	return _GnosisSafe130.Contract.VERSION(&_GnosisSafe130.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_GnosisSafe130 *GnosisSafe130CallerSession) VERSION() (string, error) {
	return _GnosisSafe130.Contract.VERSION(&_GnosisSafe130.CallOpts)
}

// ApprovedHashes is a free data retrieval call binding the contract method 0x7d832974.
//
// Solidity: function approvedHashes(address , bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Caller) ApprovedHashes(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "approvedHashes", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApprovedHashes is a free data retrieval call binding the contract method 0x7d832974.
//
// Solidity: function approvedHashes(address , bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) ApprovedHashes(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _GnosisSafe130.Contract.ApprovedHashes(&_GnosisSafe130.CallOpts, arg0, arg1)
}

// ApprovedHashes is a free data retrieval call binding the contract method 0x7d832974.
//
// Solidity: function approvedHashes(address , bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130CallerSession) ApprovedHashes(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _GnosisSafe130.Contract.ApprovedHashes(&_GnosisSafe130.CallOpts, arg0, arg1)
}

// CheckNSignatures is a free data retrieval call binding the contract method 0x12fb68e0.
//
// Solidity: function checkNSignatures(bytes32 dataHash, bytes data, bytes signatures, uint256 requiredSignatures) view returns()
func (_GnosisSafe130 *GnosisSafe130Caller) CheckNSignatures(opts *bind.CallOpts, dataHash [32]byte, data []byte, signatures []byte, requiredSignatures *big.Int) error {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "checkNSignatures", dataHash, data, signatures, requiredSignatures)

	if err != nil {
		return err
	}

	return err

}

// CheckNSignatures is a free data retrieval call binding the contract method 0x12fb68e0.
//
// Solidity: function checkNSignatures(bytes32 dataHash, bytes data, bytes signatures, uint256 requiredSignatures) view returns()
func (_GnosisSafe130 *GnosisSafe130Session) CheckNSignatures(dataHash [32]byte, data []byte, signatures []byte, requiredSignatures *big.Int) error {
	return _GnosisSafe130.Contract.CheckNSignatures(&_GnosisSafe130.CallOpts, dataHash, data, signatures, requiredSignatures)
}

// CheckNSignatures is a free data retrieval call binding the contract method 0x12fb68e0.
//
// Solidity: function checkNSignatures(bytes32 dataHash, bytes data, bytes signatures, uint256 requiredSignatures) view returns()
func (_GnosisSafe130 *GnosisSafe130CallerSession) CheckNSignatures(dataHash [32]byte, data []byte, signatures []byte, requiredSignatures *big.Int) error {
	return _GnosisSafe130.Contract.CheckNSignatures(&_GnosisSafe130.CallOpts, dataHash, data, signatures, requiredSignatures)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x934f3a11.
//
// Solidity: function checkSignatures(bytes32 dataHash, bytes data, bytes signatures) view returns()
func (_GnosisSafe130 *GnosisSafe130Caller) CheckSignatures(opts *bind.CallOpts, dataHash [32]byte, data []byte, signatures []byte) error {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "checkSignatures", dataHash, data, signatures)

	if err != nil {
		return err
	}

	return err

}

// CheckSignatures is a free data retrieval call binding the contract method 0x934f3a11.
//
// Solidity: function checkSignatures(bytes32 dataHash, bytes data, bytes signatures) view returns()
func (_GnosisSafe130 *GnosisSafe130Session) CheckSignatures(dataHash [32]byte, data []byte, signatures []byte) error {
	return _GnosisSafe130.Contract.CheckSignatures(&_GnosisSafe130.CallOpts, dataHash, data, signatures)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x934f3a11.
//
// Solidity: function checkSignatures(bytes32 dataHash, bytes data, bytes signatures) view returns()
func (_GnosisSafe130 *GnosisSafe130CallerSession) CheckSignatures(dataHash [32]byte, data []byte, signatures []byte) error {
	return _GnosisSafe130.Contract.CheckSignatures(&_GnosisSafe130.CallOpts, dataHash, data, signatures)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130Caller) DomainSeparator(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "domainSeparator")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130Session) DomainSeparator() ([32]byte, error) {
	return _GnosisSafe130.Contract.DomainSeparator(&_GnosisSafe130.CallOpts)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130CallerSession) DomainSeparator() ([32]byte, error) {
	return _GnosisSafe130.Contract.DomainSeparator(&_GnosisSafe130.CallOpts)
}

// EncodeTransactionData is a free data retrieval call binding the contract method 0xe86637db.
//
// Solidity: function encodeTransactionData(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130Caller) EncodeTransactionData(opts *bind.CallOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([]byte, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "encodeTransactionData", to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeTransactionData is a free data retrieval call binding the contract method 0xe86637db.
//
// Solidity: function encodeTransactionData(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130Session) EncodeTransactionData(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([]byte, error) {
	return _GnosisSafe130.Contract.EncodeTransactionData(&_GnosisSafe130.CallOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)
}

// EncodeTransactionData is a free data retrieval call binding the contract method 0xe86637db.
//
// Solidity: function encodeTransactionData(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130CallerSession) EncodeTransactionData(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([]byte, error) {
	return _GnosisSafe130.Contract.EncodeTransactionData(&_GnosisSafe130.CallOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)
}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Caller) GetChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) GetChainId() (*big.Int, error) {
	return _GnosisSafe130.Contract.GetChainId(&_GnosisSafe130.CallOpts)
}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetChainId() (*big.Int, error) {
	return _GnosisSafe130.Contract.GetChainId(&_GnosisSafe130.CallOpts)
}

// GetModulesPaginated is a free data retrieval call binding the contract method 0xcc2f8452.
//
// Solidity: function getModulesPaginated(address start, uint256 pageSize) view returns(address[] array, address next)
func (_GnosisSafe130 *GnosisSafe130Caller) GetModulesPaginated(opts *bind.CallOpts, start common.Address, pageSize *big.Int) (struct {
	Array []common.Address
	Next  common.Address
}, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getModulesPaginated", start, pageSize)

	outstruct := new(struct {
		Array []common.Address
		Next  common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Array = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Next = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetModulesPaginated is a free data retrieval call binding the contract method 0xcc2f8452.
//
// Solidity: function getModulesPaginated(address start, uint256 pageSize) view returns(address[] array, address next)
func (_GnosisSafe130 *GnosisSafe130Session) GetModulesPaginated(start common.Address, pageSize *big.Int) (struct {
	Array []common.Address
	Next  common.Address
}, error) {
	return _GnosisSafe130.Contract.GetModulesPaginated(&_GnosisSafe130.CallOpts, start, pageSize)
}

// GetModulesPaginated is a free data retrieval call binding the contract method 0xcc2f8452.
//
// Solidity: function getModulesPaginated(address start, uint256 pageSize) view returns(address[] array, address next)
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetModulesPaginated(start common.Address, pageSize *big.Int) (struct {
	Array []common.Address
	Next  common.Address
}, error) {
	return _GnosisSafe130.Contract.GetModulesPaginated(&_GnosisSafe130.CallOpts, start, pageSize)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_GnosisSafe130 *GnosisSafe130Caller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_GnosisSafe130 *GnosisSafe130Session) GetOwners() ([]common.Address, error) {
	return _GnosisSafe130.Contract.GetOwners(&_GnosisSafe130.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetOwners() ([]common.Address, error) {
	return _GnosisSafe130.Contract.GetOwners(&_GnosisSafe130.CallOpts)
}

// GetStorageAt is a free data retrieval call binding the contract method 0x5624b25b.
//
// Solidity: function getStorageAt(uint256 offset, uint256 length) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130Caller) GetStorageAt(opts *bind.CallOpts, offset *big.Int, length *big.Int) ([]byte, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getStorageAt", offset, length)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStorageAt is a free data retrieval call binding the contract method 0x5624b25b.
//
// Solidity: function getStorageAt(uint256 offset, uint256 length) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130Session) GetStorageAt(offset *big.Int, length *big.Int) ([]byte, error) {
	return _GnosisSafe130.Contract.GetStorageAt(&_GnosisSafe130.CallOpts, offset, length)
}

// GetStorageAt is a free data retrieval call binding the contract method 0x5624b25b.
//
// Solidity: function getStorageAt(uint256 offset, uint256 length) view returns(bytes)
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetStorageAt(offset *big.Int, length *big.Int) ([]byte, error) {
	return _GnosisSafe130.Contract.GetStorageAt(&_GnosisSafe130.CallOpts, offset, length)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Caller) GetThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) GetThreshold() (*big.Int, error) {
	return _GnosisSafe130.Contract.GetThreshold(&_GnosisSafe130.CallOpts)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetThreshold() (*big.Int, error) {
	return _GnosisSafe130.Contract.GetThreshold(&_GnosisSafe130.CallOpts)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0xd8d11f78.
//
// Solidity: function getTransactionHash(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130Caller) GetTransactionHash(opts *bind.CallOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "getTransactionHash", to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTransactionHash is a free data retrieval call binding the contract method 0xd8d11f78.
//
// Solidity: function getTransactionHash(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130Session) GetTransactionHash(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([32]byte, error) {
	return _GnosisSafe130.Contract.GetTransactionHash(&_GnosisSafe130.CallOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0xd8d11f78.
//
// Solidity: function getTransactionHash(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns(bytes32)
func (_GnosisSafe130 *GnosisSafe130CallerSession) GetTransactionHash(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int) ([32]byte, error) {
	return _GnosisSafe130.Contract.GetTransactionHash(&_GnosisSafe130.CallOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, _nonce)
}

// IsModuleEnabled is a free data retrieval call binding the contract method 0x2d9ad53d.
//
// Solidity: function isModuleEnabled(address module) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130Caller) IsModuleEnabled(opts *bind.CallOpts, module common.Address) (bool, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "isModuleEnabled", module)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsModuleEnabled is a free data retrieval call binding the contract method 0x2d9ad53d.
//
// Solidity: function isModuleEnabled(address module) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130Session) IsModuleEnabled(module common.Address) (bool, error) {
	return _GnosisSafe130.Contract.IsModuleEnabled(&_GnosisSafe130.CallOpts, module)
}

// IsModuleEnabled is a free data retrieval call binding the contract method 0x2d9ad53d.
//
// Solidity: function isModuleEnabled(address module) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130CallerSession) IsModuleEnabled(module common.Address) (bool, error) {
	return _GnosisSafe130.Contract.IsModuleEnabled(&_GnosisSafe130.CallOpts, module)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130Caller) IsOwner(opts *bind.CallOpts, owner common.Address) (bool, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "isOwner", owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130Session) IsOwner(owner common.Address) (bool, error) {
	return _GnosisSafe130.Contract.IsOwner(&_GnosisSafe130.CallOpts, owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_GnosisSafe130 *GnosisSafe130CallerSession) IsOwner(owner common.Address) (bool, error) {
	return _GnosisSafe130.Contract.IsOwner(&_GnosisSafe130.CallOpts, owner)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Caller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) Nonce() (*big.Int, error) {
	return _GnosisSafe130.Contract.Nonce(&_GnosisSafe130.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130CallerSession) Nonce() (*big.Int, error) {
	return _GnosisSafe130.Contract.Nonce(&_GnosisSafe130.CallOpts)
}

// SignedMessages is a free data retrieval call binding the contract method 0x5ae6bd37.
//
// Solidity: function signedMessages(bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Caller) SignedMessages(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GnosisSafe130.contract.Call(opts, &out, "signedMessages", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SignedMessages is a free data retrieval call binding the contract method 0x5ae6bd37.
//
// Solidity: function signedMessages(bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) SignedMessages(arg0 [32]byte) (*big.Int, error) {
	return _GnosisSafe130.Contract.SignedMessages(&_GnosisSafe130.CallOpts, arg0)
}

// SignedMessages is a free data retrieval call binding the contract method 0x5ae6bd37.
//
// Solidity: function signedMessages(bytes32 ) view returns(uint256)
func (_GnosisSafe130 *GnosisSafe130CallerSession) SignedMessages(arg0 [32]byte) (*big.Int, error) {
	return _GnosisSafe130.Contract.SignedMessages(&_GnosisSafe130.CallOpts, arg0)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) AddOwnerWithThreshold(opts *bind.TransactOpts, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "addOwnerWithThreshold", owner, _threshold)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Session) AddOwnerWithThreshold(owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.AddOwnerWithThreshold(&_GnosisSafe130.TransactOpts, owner, _threshold)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) AddOwnerWithThreshold(owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.AddOwnerWithThreshold(&_GnosisSafe130.TransactOpts, owner, _threshold)
}

// ApproveHash is a paid mutator transaction binding the contract method 0xd4d9bdcd.
//
// Solidity: function approveHash(bytes32 hashToApprove) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) ApproveHash(opts *bind.TransactOpts, hashToApprove [32]byte) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "approveHash", hashToApprove)
}

// ApproveHash is a paid mutator transaction binding the contract method 0xd4d9bdcd.
//
// Solidity: function approveHash(bytes32 hashToApprove) returns()
func (_GnosisSafe130 *GnosisSafe130Session) ApproveHash(hashToApprove [32]byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ApproveHash(&_GnosisSafe130.TransactOpts, hashToApprove)
}

// ApproveHash is a paid mutator transaction binding the contract method 0xd4d9bdcd.
//
// Solidity: function approveHash(bytes32 hashToApprove) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) ApproveHash(hashToApprove [32]byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ApproveHash(&_GnosisSafe130.TransactOpts, hashToApprove)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) ChangeThreshold(opts *bind.TransactOpts, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "changeThreshold", _threshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Session) ChangeThreshold(_threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ChangeThreshold(&_GnosisSafe130.TransactOpts, _threshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) ChangeThreshold(_threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ChangeThreshold(&_GnosisSafe130.TransactOpts, _threshold)
}

// DisableModule is a paid mutator transaction binding the contract method 0xe009cfde.
//
// Solidity: function disableModule(address prevModule, address module) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) DisableModule(opts *bind.TransactOpts, prevModule common.Address, module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "disableModule", prevModule, module)
}

// DisableModule is a paid mutator transaction binding the contract method 0xe009cfde.
//
// Solidity: function disableModule(address prevModule, address module) returns()
func (_GnosisSafe130 *GnosisSafe130Session) DisableModule(prevModule common.Address, module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.DisableModule(&_GnosisSafe130.TransactOpts, prevModule, module)
}

// DisableModule is a paid mutator transaction binding the contract method 0xe009cfde.
//
// Solidity: function disableModule(address prevModule, address module) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) DisableModule(prevModule common.Address, module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.DisableModule(&_GnosisSafe130.TransactOpts, prevModule, module)
}

// EnableModule is a paid mutator transaction binding the contract method 0x610b5925.
//
// Solidity: function enableModule(address module) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) EnableModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "enableModule", module)
}

// EnableModule is a paid mutator transaction binding the contract method 0x610b5925.
//
// Solidity: function enableModule(address module) returns()
func (_GnosisSafe130 *GnosisSafe130Session) EnableModule(module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.EnableModule(&_GnosisSafe130.TransactOpts, module)
}

// EnableModule is a paid mutator transaction binding the contract method 0x610b5925.
//
// Solidity: function enableModule(address module) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) EnableModule(module common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.EnableModule(&_GnosisSafe130.TransactOpts, module)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0x6a761202.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, bytes signatures) payable returns(bool success)
func (_GnosisSafe130 *GnosisSafe130Transactor) ExecTransaction(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, signatures []byte) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "execTransaction", to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, signatures)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0x6a761202.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, bytes signatures) payable returns(bool success)
func (_GnosisSafe130 *GnosisSafe130Session) ExecTransaction(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, signatures []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransaction(&_GnosisSafe130.TransactOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, signatures)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0x6a761202.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, bytes signatures) payable returns(bool success)
func (_GnosisSafe130 *GnosisSafe130TransactorSession) ExecTransaction(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, signatures []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransaction(&_GnosisSafe130.TransactOpts, to, value, data, operation, safeTxGas, baseGas, gasPrice, gasToken, refundReceiver, signatures)
}

// ExecTransactionFromModule is a paid mutator transaction binding the contract method 0x468721a7.
//
// Solidity: function execTransactionFromModule(address to, uint256 value, bytes data, uint8 operation) returns(bool success)
func (_GnosisSafe130 *GnosisSafe130Transactor) ExecTransactionFromModule(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "execTransactionFromModule", to, value, data, operation)
}

// ExecTransactionFromModule is a paid mutator transaction binding the contract method 0x468721a7.
//
// Solidity: function execTransactionFromModule(address to, uint256 value, bytes data, uint8 operation) returns(bool success)
func (_GnosisSafe130 *GnosisSafe130Session) ExecTransactionFromModule(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransactionFromModule(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// ExecTransactionFromModule is a paid mutator transaction binding the contract method 0x468721a7.
//
// Solidity: function execTransactionFromModule(address to, uint256 value, bytes data, uint8 operation) returns(bool success)
func (_GnosisSafe130 *GnosisSafe130TransactorSession) ExecTransactionFromModule(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransactionFromModule(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// ExecTransactionFromModuleReturnData is a paid mutator transaction binding the contract method 0x5229073f.
//
// Solidity: function execTransactionFromModuleReturnData(address to, uint256 value, bytes data, uint8 operation) returns(bool success, bytes returnData)
func (_GnosisSafe130 *GnosisSafe130Transactor) ExecTransactionFromModuleReturnData(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "execTransactionFromModuleReturnData", to, value, data, operation)
}

// ExecTransactionFromModuleReturnData is a paid mutator transaction binding the contract method 0x5229073f.
//
// Solidity: function execTransactionFromModuleReturnData(address to, uint256 value, bytes data, uint8 operation) returns(bool success, bytes returnData)
func (_GnosisSafe130 *GnosisSafe130Session) ExecTransactionFromModuleReturnData(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransactionFromModuleReturnData(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// ExecTransactionFromModuleReturnData is a paid mutator transaction binding the contract method 0x5229073f.
//
// Solidity: function execTransactionFromModuleReturnData(address to, uint256 value, bytes data, uint8 operation) returns(bool success, bytes returnData)
func (_GnosisSafe130 *GnosisSafe130TransactorSession) ExecTransactionFromModuleReturnData(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.ExecTransactionFromModuleReturnData(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) RemoveOwner(opts *bind.TransactOpts, prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "removeOwner", prevOwner, owner, _threshold)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130Session) RemoveOwner(prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.RemoveOwner(&_GnosisSafe130.TransactOpts, prevOwner, owner, _threshold)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) RemoveOwner(prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.RemoveOwner(&_GnosisSafe130.TransactOpts, prevOwner, owner, _threshold)
}

// RequiredTxGas is a paid mutator transaction binding the contract method 0xc4ca3a9c.
//
// Solidity: function requiredTxGas(address to, uint256 value, bytes data, uint8 operation) returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Transactor) RequiredTxGas(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "requiredTxGas", to, value, data, operation)
}

// RequiredTxGas is a paid mutator transaction binding the contract method 0xc4ca3a9c.
//
// Solidity: function requiredTxGas(address to, uint256 value, bytes data, uint8 operation) returns(uint256)
func (_GnosisSafe130 *GnosisSafe130Session) RequiredTxGas(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.RequiredTxGas(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// RequiredTxGas is a paid mutator transaction binding the contract method 0xc4ca3a9c.
//
// Solidity: function requiredTxGas(address to, uint256 value, bytes data, uint8 operation) returns(uint256)
func (_GnosisSafe130 *GnosisSafe130TransactorSession) RequiredTxGas(to common.Address, value *big.Int, data []byte, operation uint8) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.RequiredTxGas(&_GnosisSafe130.TransactOpts, to, value, data, operation)
}

// SetFallbackHandler is a paid mutator transaction binding the contract method 0xf08a0323.
//
// Solidity: function setFallbackHandler(address handler) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) SetFallbackHandler(opts *bind.TransactOpts, handler common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "setFallbackHandler", handler)
}

// SetFallbackHandler is a paid mutator transaction binding the contract method 0xf08a0323.
//
// Solidity: function setFallbackHandler(address handler) returns()
func (_GnosisSafe130 *GnosisSafe130Session) SetFallbackHandler(handler common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SetFallbackHandler(&_GnosisSafe130.TransactOpts, handler)
}

// SetFallbackHandler is a paid mutator transaction binding the contract method 0xf08a0323.
//
// Solidity: function setFallbackHandler(address handler) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) SetFallbackHandler(handler common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SetFallbackHandler(&_GnosisSafe130.TransactOpts, handler)
}

// SetGuard is a paid mutator transaction binding the contract method 0xe19a9dd9.
//
// Solidity: function setGuard(address guard) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) SetGuard(opts *bind.TransactOpts, guard common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "setGuard", guard)
}

// SetGuard is a paid mutator transaction binding the contract method 0xe19a9dd9.
//
// Solidity: function setGuard(address guard) returns()
func (_GnosisSafe130 *GnosisSafe130Session) SetGuard(guard common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SetGuard(&_GnosisSafe130.TransactOpts, guard)
}

// SetGuard is a paid mutator transaction binding the contract method 0xe19a9dd9.
//
// Solidity: function setGuard(address guard) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) SetGuard(guard common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SetGuard(&_GnosisSafe130.TransactOpts, guard)
}

// Setup is a paid mutator transaction binding the contract method 0xb63e800d.
//
// Solidity: function setup(address[] _owners, uint256 _threshold, address to, bytes data, address fallbackHandler, address paymentToken, uint256 payment, address paymentReceiver) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) Setup(opts *bind.TransactOpts, _owners []common.Address, _threshold *big.Int, to common.Address, data []byte, fallbackHandler common.Address, paymentToken common.Address, payment *big.Int, paymentReceiver common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "setup", _owners, _threshold, to, data, fallbackHandler, paymentToken, payment, paymentReceiver)
}

// Setup is a paid mutator transaction binding the contract method 0xb63e800d.
//
// Solidity: function setup(address[] _owners, uint256 _threshold, address to, bytes data, address fallbackHandler, address paymentToken, uint256 payment, address paymentReceiver) returns()
func (_GnosisSafe130 *GnosisSafe130Session) Setup(_owners []common.Address, _threshold *big.Int, to common.Address, data []byte, fallbackHandler common.Address, paymentToken common.Address, payment *big.Int, paymentReceiver common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Setup(&_GnosisSafe130.TransactOpts, _owners, _threshold, to, data, fallbackHandler, paymentToken, payment, paymentReceiver)
}

// Setup is a paid mutator transaction binding the contract method 0xb63e800d.
//
// Solidity: function setup(address[] _owners, uint256 _threshold, address to, bytes data, address fallbackHandler, address paymentToken, uint256 payment, address paymentReceiver) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) Setup(_owners []common.Address, _threshold *big.Int, to common.Address, data []byte, fallbackHandler common.Address, paymentToken common.Address, payment *big.Int, paymentReceiver common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Setup(&_GnosisSafe130.TransactOpts, _owners, _threshold, to, data, fallbackHandler, paymentToken, payment, paymentReceiver)
}

// SimulateAndRevert is a paid mutator transaction binding the contract method 0xb4faba09.
//
// Solidity: function simulateAndRevert(address targetContract, bytes calldataPayload) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) SimulateAndRevert(opts *bind.TransactOpts, targetContract common.Address, calldataPayload []byte) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "simulateAndRevert", targetContract, calldataPayload)
}

// SimulateAndRevert is a paid mutator transaction binding the contract method 0xb4faba09.
//
// Solidity: function simulateAndRevert(address targetContract, bytes calldataPayload) returns()
func (_GnosisSafe130 *GnosisSafe130Session) SimulateAndRevert(targetContract common.Address, calldataPayload []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SimulateAndRevert(&_GnosisSafe130.TransactOpts, targetContract, calldataPayload)
}

// SimulateAndRevert is a paid mutator transaction binding the contract method 0xb4faba09.
//
// Solidity: function simulateAndRevert(address targetContract, bytes calldataPayload) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) SimulateAndRevert(targetContract common.Address, calldataPayload []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SimulateAndRevert(&_GnosisSafe130.TransactOpts, targetContract, calldataPayload)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) SwapOwner(opts *bind.TransactOpts, prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.contract.Transact(opts, "swapOwner", prevOwner, oldOwner, newOwner)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_GnosisSafe130 *GnosisSafe130Session) SwapOwner(prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SwapOwner(&_GnosisSafe130.TransactOpts, prevOwner, oldOwner, newOwner)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) SwapOwner(prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.SwapOwner(&_GnosisSafe130.TransactOpts, prevOwner, oldOwner, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _GnosisSafe130.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_GnosisSafe130 *GnosisSafe130Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Fallback(&_GnosisSafe130.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Fallback(&_GnosisSafe130.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GnosisSafe130 *GnosisSafe130Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GnosisSafe130.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GnosisSafe130 *GnosisSafe130Session) Receive() (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Receive(&_GnosisSafe130.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GnosisSafe130 *GnosisSafe130TransactorSession) Receive() (*types.Transaction, error) {
	return _GnosisSafe130.Contract.Receive(&_GnosisSafe130.TransactOpts)
}

// GnosisSafe130AddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the GnosisSafe130 contract.
type GnosisSafe130AddedOwnerIterator struct {
	Event *GnosisSafe130AddedOwner // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130AddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130AddedOwner)
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
		it.Event = new(GnosisSafe130AddedOwner)
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
func (it *GnosisSafe130AddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130AddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130AddedOwner represents a AddedOwner event raised by the GnosisSafe130 contract.
type GnosisSafe130AddedOwner struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterAddedOwner(opts *bind.FilterOpts) (*GnosisSafe130AddedOwnerIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130AddedOwnerIterator{contract: _GnosisSafe130.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *GnosisSafe130AddedOwner) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130AddedOwner)
				if err := _GnosisSafe130.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// ParseAddedOwner is a log parse operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseAddedOwner(log types.Log) (*GnosisSafe130AddedOwner, error) {
	event := new(GnosisSafe130AddedOwner)
	if err := _GnosisSafe130.contract.UnpackLog(event, "AddedOwner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ApproveHashIterator is returned from FilterApproveHash and is used to iterate over the raw logs and unpacked data for ApproveHash events raised by the GnosisSafe130 contract.
type GnosisSafe130ApproveHashIterator struct {
	Event *GnosisSafe130ApproveHash // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ApproveHashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ApproveHash)
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
		it.Event = new(GnosisSafe130ApproveHash)
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
func (it *GnosisSafe130ApproveHashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ApproveHashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ApproveHash represents a ApproveHash event raised by the GnosisSafe130 contract.
type GnosisSafe130ApproveHash struct {
	ApprovedHash [32]byte
	Owner        common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterApproveHash is a free log retrieval operation binding the contract event 0xf2a0eb156472d1440255b0d7c1e19cc07115d1051fe605b0dce69acfec884d9c.
//
// Solidity: event ApproveHash(bytes32 indexed approvedHash, address indexed owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterApproveHash(opts *bind.FilterOpts, approvedHash [][32]byte, owner []common.Address) (*GnosisSafe130ApproveHashIterator, error) {

	var approvedHashRule []interface{}
	for _, approvedHashItem := range approvedHash {
		approvedHashRule = append(approvedHashRule, approvedHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ApproveHash", approvedHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ApproveHashIterator{contract: _GnosisSafe130.contract, event: "ApproveHash", logs: logs, sub: sub}, nil
}

// WatchApproveHash is a free log subscription operation binding the contract event 0xf2a0eb156472d1440255b0d7c1e19cc07115d1051fe605b0dce69acfec884d9c.
//
// Solidity: event ApproveHash(bytes32 indexed approvedHash, address indexed owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchApproveHash(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ApproveHash, approvedHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var approvedHashRule []interface{}
	for _, approvedHashItem := range approvedHash {
		approvedHashRule = append(approvedHashRule, approvedHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ApproveHash", approvedHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ApproveHash)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ApproveHash", log); err != nil {
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

// ParseApproveHash is a log parse operation binding the contract event 0xf2a0eb156472d1440255b0d7c1e19cc07115d1051fe605b0dce69acfec884d9c.
//
// Solidity: event ApproveHash(bytes32 indexed approvedHash, address indexed owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseApproveHash(log types.Log) (*GnosisSafe130ApproveHash, error) {
	event := new(GnosisSafe130ApproveHash)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ApproveHash", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ChangedFallbackHandlerIterator is returned from FilterChangedFallbackHandler and is used to iterate over the raw logs and unpacked data for ChangedFallbackHandler events raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedFallbackHandlerIterator struct {
	Event *GnosisSafe130ChangedFallbackHandler // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ChangedFallbackHandlerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ChangedFallbackHandler)
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
		it.Event = new(GnosisSafe130ChangedFallbackHandler)
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
func (it *GnosisSafe130ChangedFallbackHandlerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ChangedFallbackHandlerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ChangedFallbackHandler represents a ChangedFallbackHandler event raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedFallbackHandler struct {
	Handler common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterChangedFallbackHandler is a free log retrieval operation binding the contract event 0x5ac6c46c93c8d0e53714ba3b53db3e7c046da994313d7ed0d192028bc7c228b0.
//
// Solidity: event ChangedFallbackHandler(address handler)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterChangedFallbackHandler(opts *bind.FilterOpts) (*GnosisSafe130ChangedFallbackHandlerIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ChangedFallbackHandler")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ChangedFallbackHandlerIterator{contract: _GnosisSafe130.contract, event: "ChangedFallbackHandler", logs: logs, sub: sub}, nil
}

// WatchChangedFallbackHandler is a free log subscription operation binding the contract event 0x5ac6c46c93c8d0e53714ba3b53db3e7c046da994313d7ed0d192028bc7c228b0.
//
// Solidity: event ChangedFallbackHandler(address handler)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchChangedFallbackHandler(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ChangedFallbackHandler) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ChangedFallbackHandler")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ChangedFallbackHandler)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedFallbackHandler", log); err != nil {
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

// ParseChangedFallbackHandler is a log parse operation binding the contract event 0x5ac6c46c93c8d0e53714ba3b53db3e7c046da994313d7ed0d192028bc7c228b0.
//
// Solidity: event ChangedFallbackHandler(address handler)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseChangedFallbackHandler(log types.Log) (*GnosisSafe130ChangedFallbackHandler, error) {
	event := new(GnosisSafe130ChangedFallbackHandler)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedFallbackHandler", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ChangedGuardIterator is returned from FilterChangedGuard and is used to iterate over the raw logs and unpacked data for ChangedGuard events raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedGuardIterator struct {
	Event *GnosisSafe130ChangedGuard // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ChangedGuardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ChangedGuard)
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
		it.Event = new(GnosisSafe130ChangedGuard)
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
func (it *GnosisSafe130ChangedGuardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ChangedGuardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ChangedGuard represents a ChangedGuard event raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedGuard struct {
	Guard common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterChangedGuard is a free log retrieval operation binding the contract event 0x1151116914515bc0891ff9047a6cb32cf902546f83066499bcf8ba33d2353fa2.
//
// Solidity: event ChangedGuard(address guard)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterChangedGuard(opts *bind.FilterOpts) (*GnosisSafe130ChangedGuardIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ChangedGuard")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ChangedGuardIterator{contract: _GnosisSafe130.contract, event: "ChangedGuard", logs: logs, sub: sub}, nil
}

// WatchChangedGuard is a free log subscription operation binding the contract event 0x1151116914515bc0891ff9047a6cb32cf902546f83066499bcf8ba33d2353fa2.
//
// Solidity: event ChangedGuard(address guard)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchChangedGuard(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ChangedGuard) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ChangedGuard")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ChangedGuard)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedGuard", log); err != nil {
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

// ParseChangedGuard is a log parse operation binding the contract event 0x1151116914515bc0891ff9047a6cb32cf902546f83066499bcf8ba33d2353fa2.
//
// Solidity: event ChangedGuard(address guard)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseChangedGuard(log types.Log) (*GnosisSafe130ChangedGuard, error) {
	event := new(GnosisSafe130ChangedGuard)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedGuard", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ChangedThresholdIterator is returned from FilterChangedThreshold and is used to iterate over the raw logs and unpacked data for ChangedThreshold events raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedThresholdIterator struct {
	Event *GnosisSafe130ChangedThreshold // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ChangedThresholdIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ChangedThreshold)
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
		it.Event = new(GnosisSafe130ChangedThreshold)
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
func (it *GnosisSafe130ChangedThresholdIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ChangedThresholdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ChangedThreshold represents a ChangedThreshold event raised by the GnosisSafe130 contract.
type GnosisSafe130ChangedThreshold struct {
	Threshold *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChangedThreshold is a free log retrieval operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterChangedThreshold(opts *bind.FilterOpts) (*GnosisSafe130ChangedThresholdIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ChangedThreshold")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ChangedThresholdIterator{contract: _GnosisSafe130.contract, event: "ChangedThreshold", logs: logs, sub: sub}, nil
}

// WatchChangedThreshold is a free log subscription operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchChangedThreshold(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ChangedThreshold) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ChangedThreshold")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ChangedThreshold)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedThreshold", log); err != nil {
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

// ParseChangedThreshold is a log parse operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseChangedThreshold(log types.Log) (*GnosisSafe130ChangedThreshold, error) {
	event := new(GnosisSafe130ChangedThreshold)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ChangedThreshold", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130DisabledModuleIterator is returned from FilterDisabledModule and is used to iterate over the raw logs and unpacked data for DisabledModule events raised by the GnosisSafe130 contract.
type GnosisSafe130DisabledModuleIterator struct {
	Event *GnosisSafe130DisabledModule // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130DisabledModuleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130DisabledModule)
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
		it.Event = new(GnosisSafe130DisabledModule)
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
func (it *GnosisSafe130DisabledModuleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130DisabledModuleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130DisabledModule represents a DisabledModule event raised by the GnosisSafe130 contract.
type GnosisSafe130DisabledModule struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDisabledModule is a free log retrieval operation binding the contract event 0xaab4fa2b463f581b2b32cb3b7e3b704b9ce37cc209b5fb4d77e593ace4054276.
//
// Solidity: event DisabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterDisabledModule(opts *bind.FilterOpts) (*GnosisSafe130DisabledModuleIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "DisabledModule")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130DisabledModuleIterator{contract: _GnosisSafe130.contract, event: "DisabledModule", logs: logs, sub: sub}, nil
}

// WatchDisabledModule is a free log subscription operation binding the contract event 0xaab4fa2b463f581b2b32cb3b7e3b704b9ce37cc209b5fb4d77e593ace4054276.
//
// Solidity: event DisabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchDisabledModule(opts *bind.WatchOpts, sink chan<- *GnosisSafe130DisabledModule) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "DisabledModule")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130DisabledModule)
				if err := _GnosisSafe130.contract.UnpackLog(event, "DisabledModule", log); err != nil {
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

// ParseDisabledModule is a log parse operation binding the contract event 0xaab4fa2b463f581b2b32cb3b7e3b704b9ce37cc209b5fb4d77e593ace4054276.
//
// Solidity: event DisabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseDisabledModule(log types.Log) (*GnosisSafe130DisabledModule, error) {
	event := new(GnosisSafe130DisabledModule)
	if err := _GnosisSafe130.contract.UnpackLog(event, "DisabledModule", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130EnabledModuleIterator is returned from FilterEnabledModule and is used to iterate over the raw logs and unpacked data for EnabledModule events raised by the GnosisSafe130 contract.
type GnosisSafe130EnabledModuleIterator struct {
	Event *GnosisSafe130EnabledModule // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130EnabledModuleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130EnabledModule)
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
		it.Event = new(GnosisSafe130EnabledModule)
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
func (it *GnosisSafe130EnabledModuleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130EnabledModuleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130EnabledModule represents a EnabledModule event raised by the GnosisSafe130 contract.
type GnosisSafe130EnabledModule struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEnabledModule is a free log retrieval operation binding the contract event 0xecdf3a3effea5783a3c4c2140e677577666428d44ed9d474a0b3a4c9943f8440.
//
// Solidity: event EnabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterEnabledModule(opts *bind.FilterOpts) (*GnosisSafe130EnabledModuleIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "EnabledModule")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130EnabledModuleIterator{contract: _GnosisSafe130.contract, event: "EnabledModule", logs: logs, sub: sub}, nil
}

// WatchEnabledModule is a free log subscription operation binding the contract event 0xecdf3a3effea5783a3c4c2140e677577666428d44ed9d474a0b3a4c9943f8440.
//
// Solidity: event EnabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchEnabledModule(opts *bind.WatchOpts, sink chan<- *GnosisSafe130EnabledModule) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "EnabledModule")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130EnabledModule)
				if err := _GnosisSafe130.contract.UnpackLog(event, "EnabledModule", log); err != nil {
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

// ParseEnabledModule is a log parse operation binding the contract event 0xecdf3a3effea5783a3c4c2140e677577666428d44ed9d474a0b3a4c9943f8440.
//
// Solidity: event EnabledModule(address module)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseEnabledModule(log types.Log) (*GnosisSafe130EnabledModule, error) {
	event := new(GnosisSafe130EnabledModule)
	if err := _GnosisSafe130.contract.UnpackLog(event, "EnabledModule", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ExecutionFailureIterator is returned from FilterExecutionFailure and is used to iterate over the raw logs and unpacked data for ExecutionFailure events raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFailureIterator struct {
	Event *GnosisSafe130ExecutionFailure // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ExecutionFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ExecutionFailure)
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
		it.Event = new(GnosisSafe130ExecutionFailure)
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
func (it *GnosisSafe130ExecutionFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ExecutionFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ExecutionFailure represents a ExecutionFailure event raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFailure struct {
	TxHash  [32]byte
	Payment *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterExecutionFailure is a free log retrieval operation binding the contract event 0x23428b18acfb3ea64b08dc0c1d296ea9c09702c09083ca5272e64d115b687d23.
//
// Solidity: event ExecutionFailure(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterExecutionFailure(opts *bind.FilterOpts) (*GnosisSafe130ExecutionFailureIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ExecutionFailure")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ExecutionFailureIterator{contract: _GnosisSafe130.contract, event: "ExecutionFailure", logs: logs, sub: sub}, nil
}

// WatchExecutionFailure is a free log subscription operation binding the contract event 0x23428b18acfb3ea64b08dc0c1d296ea9c09702c09083ca5272e64d115b687d23.
//
// Solidity: event ExecutionFailure(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchExecutionFailure(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ExecutionFailure) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ExecutionFailure")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ExecutionFailure)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFailure", log); err != nil {
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

// ParseExecutionFailure is a log parse operation binding the contract event 0x23428b18acfb3ea64b08dc0c1d296ea9c09702c09083ca5272e64d115b687d23.
//
// Solidity: event ExecutionFailure(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseExecutionFailure(log types.Log) (*GnosisSafe130ExecutionFailure, error) {
	event := new(GnosisSafe130ExecutionFailure)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFailure", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ExecutionFromModuleFailureIterator is returned from FilterExecutionFromModuleFailure and is used to iterate over the raw logs and unpacked data for ExecutionFromModuleFailure events raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFromModuleFailureIterator struct {
	Event *GnosisSafe130ExecutionFromModuleFailure // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ExecutionFromModuleFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ExecutionFromModuleFailure)
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
		it.Event = new(GnosisSafe130ExecutionFromModuleFailure)
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
func (it *GnosisSafe130ExecutionFromModuleFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ExecutionFromModuleFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ExecutionFromModuleFailure represents a ExecutionFromModuleFailure event raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFromModuleFailure struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExecutionFromModuleFailure is a free log retrieval operation binding the contract event 0xacd2c8702804128fdb0db2bb49f6d127dd0181c13fd45dbfe16de0930e2bd375.
//
// Solidity: event ExecutionFromModuleFailure(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterExecutionFromModuleFailure(opts *bind.FilterOpts, module []common.Address) (*GnosisSafe130ExecutionFromModuleFailureIterator, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ExecutionFromModuleFailure", moduleRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ExecutionFromModuleFailureIterator{contract: _GnosisSafe130.contract, event: "ExecutionFromModuleFailure", logs: logs, sub: sub}, nil
}

// WatchExecutionFromModuleFailure is a free log subscription operation binding the contract event 0xacd2c8702804128fdb0db2bb49f6d127dd0181c13fd45dbfe16de0930e2bd375.
//
// Solidity: event ExecutionFromModuleFailure(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchExecutionFromModuleFailure(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ExecutionFromModuleFailure, module []common.Address) (event.Subscription, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ExecutionFromModuleFailure", moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ExecutionFromModuleFailure)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFromModuleFailure", log); err != nil {
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

// ParseExecutionFromModuleFailure is a log parse operation binding the contract event 0xacd2c8702804128fdb0db2bb49f6d127dd0181c13fd45dbfe16de0930e2bd375.
//
// Solidity: event ExecutionFromModuleFailure(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseExecutionFromModuleFailure(log types.Log) (*GnosisSafe130ExecutionFromModuleFailure, error) {
	event := new(GnosisSafe130ExecutionFromModuleFailure)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFromModuleFailure", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ExecutionFromModuleSuccessIterator is returned from FilterExecutionFromModuleSuccess and is used to iterate over the raw logs and unpacked data for ExecutionFromModuleSuccess events raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFromModuleSuccessIterator struct {
	Event *GnosisSafe130ExecutionFromModuleSuccess // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ExecutionFromModuleSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ExecutionFromModuleSuccess)
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
		it.Event = new(GnosisSafe130ExecutionFromModuleSuccess)
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
func (it *GnosisSafe130ExecutionFromModuleSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ExecutionFromModuleSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ExecutionFromModuleSuccess represents a ExecutionFromModuleSuccess event raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionFromModuleSuccess struct {
	Module common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExecutionFromModuleSuccess is a free log retrieval operation binding the contract event 0x6895c13664aa4f67288b25d7a21d7aaa34916e355fb9b6fae0a139a9085becb8.
//
// Solidity: event ExecutionFromModuleSuccess(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterExecutionFromModuleSuccess(opts *bind.FilterOpts, module []common.Address) (*GnosisSafe130ExecutionFromModuleSuccessIterator, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ExecutionFromModuleSuccess", moduleRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ExecutionFromModuleSuccessIterator{contract: _GnosisSafe130.contract, event: "ExecutionFromModuleSuccess", logs: logs, sub: sub}, nil
}

// WatchExecutionFromModuleSuccess is a free log subscription operation binding the contract event 0x6895c13664aa4f67288b25d7a21d7aaa34916e355fb9b6fae0a139a9085becb8.
//
// Solidity: event ExecutionFromModuleSuccess(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchExecutionFromModuleSuccess(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ExecutionFromModuleSuccess, module []common.Address) (event.Subscription, error) {

	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ExecutionFromModuleSuccess", moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ExecutionFromModuleSuccess)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFromModuleSuccess", log); err != nil {
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

// ParseExecutionFromModuleSuccess is a log parse operation binding the contract event 0x6895c13664aa4f67288b25d7a21d7aaa34916e355fb9b6fae0a139a9085becb8.
//
// Solidity: event ExecutionFromModuleSuccess(address indexed module)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseExecutionFromModuleSuccess(log types.Log) (*GnosisSafe130ExecutionFromModuleSuccess, error) {
	event := new(GnosisSafe130ExecutionFromModuleSuccess)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionFromModuleSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130ExecutionSuccessIterator is returned from FilterExecutionSuccess and is used to iterate over the raw logs and unpacked data for ExecutionSuccess events raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionSuccessIterator struct {
	Event *GnosisSafe130ExecutionSuccess // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130ExecutionSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130ExecutionSuccess)
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
		it.Event = new(GnosisSafe130ExecutionSuccess)
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
func (it *GnosisSafe130ExecutionSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130ExecutionSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130ExecutionSuccess represents a ExecutionSuccess event raised by the GnosisSafe130 contract.
type GnosisSafe130ExecutionSuccess struct {
	TxHash  [32]byte
	Payment *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterExecutionSuccess is a free log retrieval operation binding the contract event 0x442e715f626346e8c54381002da614f62bee8d27386535b2521ec8540898556e.
//
// Solidity: event ExecutionSuccess(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterExecutionSuccess(opts *bind.FilterOpts) (*GnosisSafe130ExecutionSuccessIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "ExecutionSuccess")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130ExecutionSuccessIterator{contract: _GnosisSafe130.contract, event: "ExecutionSuccess", logs: logs, sub: sub}, nil
}

// WatchExecutionSuccess is a free log subscription operation binding the contract event 0x442e715f626346e8c54381002da614f62bee8d27386535b2521ec8540898556e.
//
// Solidity: event ExecutionSuccess(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchExecutionSuccess(opts *bind.WatchOpts, sink chan<- *GnosisSafe130ExecutionSuccess) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "ExecutionSuccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130ExecutionSuccess)
				if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionSuccess", log); err != nil {
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

// ParseExecutionSuccess is a log parse operation binding the contract event 0x442e715f626346e8c54381002da614f62bee8d27386535b2521ec8540898556e.
//
// Solidity: event ExecutionSuccess(bytes32 txHash, uint256 payment)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseExecutionSuccess(log types.Log) (*GnosisSafe130ExecutionSuccess, error) {
	event := new(GnosisSafe130ExecutionSuccess)
	if err := _GnosisSafe130.contract.UnpackLog(event, "ExecutionSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130RemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the GnosisSafe130 contract.
type GnosisSafe130RemovedOwnerIterator struct {
	Event *GnosisSafe130RemovedOwner // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130RemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130RemovedOwner)
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
		it.Event = new(GnosisSafe130RemovedOwner)
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
func (it *GnosisSafe130RemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130RemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130RemovedOwner represents a RemovedOwner event raised by the GnosisSafe130 contract.
type GnosisSafe130RemovedOwner struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterRemovedOwner(opts *bind.FilterOpts) (*GnosisSafe130RemovedOwnerIterator, error) {

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130RemovedOwnerIterator{contract: _GnosisSafe130.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *GnosisSafe130RemovedOwner) (event.Subscription, error) {

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130RemovedOwner)
				if err := _GnosisSafe130.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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

// ParseRemovedOwner is a log parse operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseRemovedOwner(log types.Log) (*GnosisSafe130RemovedOwner, error) {
	event := new(GnosisSafe130RemovedOwner)
	if err := _GnosisSafe130.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130SafeReceivedIterator is returned from FilterSafeReceived and is used to iterate over the raw logs and unpacked data for SafeReceived events raised by the GnosisSafe130 contract.
type GnosisSafe130SafeReceivedIterator struct {
	Event *GnosisSafe130SafeReceived // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130SafeReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130SafeReceived)
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
		it.Event = new(GnosisSafe130SafeReceived)
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
func (it *GnosisSafe130SafeReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130SafeReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130SafeReceived represents a SafeReceived event raised by the GnosisSafe130 contract.
type GnosisSafe130SafeReceived struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSafeReceived is a free log retrieval operation binding the contract event 0x3d0ce9bfc3ed7d6862dbb28b2dea94561fe714a1b4d019aa8af39730d1ad7c3d.
//
// Solidity: event SafeReceived(address indexed sender, uint256 value)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterSafeReceived(opts *bind.FilterOpts, sender []common.Address) (*GnosisSafe130SafeReceivedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "SafeReceived", senderRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130SafeReceivedIterator{contract: _GnosisSafe130.contract, event: "SafeReceived", logs: logs, sub: sub}, nil
}

// WatchSafeReceived is a free log subscription operation binding the contract event 0x3d0ce9bfc3ed7d6862dbb28b2dea94561fe714a1b4d019aa8af39730d1ad7c3d.
//
// Solidity: event SafeReceived(address indexed sender, uint256 value)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchSafeReceived(opts *bind.WatchOpts, sink chan<- *GnosisSafe130SafeReceived, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "SafeReceived", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130SafeReceived)
				if err := _GnosisSafe130.contract.UnpackLog(event, "SafeReceived", log); err != nil {
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

// ParseSafeReceived is a log parse operation binding the contract event 0x3d0ce9bfc3ed7d6862dbb28b2dea94561fe714a1b4d019aa8af39730d1ad7c3d.
//
// Solidity: event SafeReceived(address indexed sender, uint256 value)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseSafeReceived(log types.Log) (*GnosisSafe130SafeReceived, error) {
	event := new(GnosisSafe130SafeReceived)
	if err := _GnosisSafe130.contract.UnpackLog(event, "SafeReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130SafeSetupIterator is returned from FilterSafeSetup and is used to iterate over the raw logs and unpacked data for SafeSetup events raised by the GnosisSafe130 contract.
type GnosisSafe130SafeSetupIterator struct {
	Event *GnosisSafe130SafeSetup // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130SafeSetupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130SafeSetup)
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
		it.Event = new(GnosisSafe130SafeSetup)
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
func (it *GnosisSafe130SafeSetupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130SafeSetupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130SafeSetup represents a SafeSetup event raised by the GnosisSafe130 contract.
type GnosisSafe130SafeSetup struct {
	Initiator       common.Address
	Owners          []common.Address
	Threshold       *big.Int
	Initializer     common.Address
	FallbackHandler common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterSafeSetup is a free log retrieval operation binding the contract event 0x141df868a6331af528e38c83b7aa03edc19be66e37ae67f9285bf4f8e3c6a1a8.
//
// Solidity: event SafeSetup(address indexed initiator, address[] owners, uint256 threshold, address initializer, address fallbackHandler)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterSafeSetup(opts *bind.FilterOpts, initiator []common.Address) (*GnosisSafe130SafeSetupIterator, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "SafeSetup", initiatorRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130SafeSetupIterator{contract: _GnosisSafe130.contract, event: "SafeSetup", logs: logs, sub: sub}, nil
}

// WatchSafeSetup is a free log subscription operation binding the contract event 0x141df868a6331af528e38c83b7aa03edc19be66e37ae67f9285bf4f8e3c6a1a8.
//
// Solidity: event SafeSetup(address indexed initiator, address[] owners, uint256 threshold, address initializer, address fallbackHandler)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchSafeSetup(opts *bind.WatchOpts, sink chan<- *GnosisSafe130SafeSetup, initiator []common.Address) (event.Subscription, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "SafeSetup", initiatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130SafeSetup)
				if err := _GnosisSafe130.contract.UnpackLog(event, "SafeSetup", log); err != nil {
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

// ParseSafeSetup is a log parse operation binding the contract event 0x141df868a6331af528e38c83b7aa03edc19be66e37ae67f9285bf4f8e3c6a1a8.
//
// Solidity: event SafeSetup(address indexed initiator, address[] owners, uint256 threshold, address initializer, address fallbackHandler)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseSafeSetup(log types.Log) (*GnosisSafe130SafeSetup, error) {
	event := new(GnosisSafe130SafeSetup)
	if err := _GnosisSafe130.contract.UnpackLog(event, "SafeSetup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GnosisSafe130SignMsgIterator is returned from FilterSignMsg and is used to iterate over the raw logs and unpacked data for SignMsg events raised by the GnosisSafe130 contract.
type GnosisSafe130SignMsgIterator struct {
	Event *GnosisSafe130SignMsg // Event containing the contract specifics and raw log

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
func (it *GnosisSafe130SignMsgIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GnosisSafe130SignMsg)
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
		it.Event = new(GnosisSafe130SignMsg)
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
func (it *GnosisSafe130SignMsgIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GnosisSafe130SignMsgIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GnosisSafe130SignMsg represents a SignMsg event raised by the GnosisSafe130 contract.
type GnosisSafe130SignMsg struct {
	MsgHash [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSignMsg is a free log retrieval operation binding the contract event 0xe7f4675038f4f6034dfcbbb24c4dc08e4ebf10eb9d257d3d02c0f38d122ac6e4.
//
// Solidity: event SignMsg(bytes32 indexed msgHash)
func (_GnosisSafe130 *GnosisSafe130Filterer) FilterSignMsg(opts *bind.FilterOpts, msgHash [][32]byte) (*GnosisSafe130SignMsgIterator, error) {

	var msgHashRule []interface{}
	for _, msgHashItem := range msgHash {
		msgHashRule = append(msgHashRule, msgHashItem)
	}

	logs, sub, err := _GnosisSafe130.contract.FilterLogs(opts, "SignMsg", msgHashRule)
	if err != nil {
		return nil, err
	}
	return &GnosisSafe130SignMsgIterator{contract: _GnosisSafe130.contract, event: "SignMsg", logs: logs, sub: sub}, nil
}

// WatchSignMsg is a free log subscription operation binding the contract event 0xe7f4675038f4f6034dfcbbb24c4dc08e4ebf10eb9d257d3d02c0f38d122ac6e4.
//
// Solidity: event SignMsg(bytes32 indexed msgHash)
func (_GnosisSafe130 *GnosisSafe130Filterer) WatchSignMsg(opts *bind.WatchOpts, sink chan<- *GnosisSafe130SignMsg, msgHash [][32]byte) (event.Subscription, error) {

	var msgHashRule []interface{}
	for _, msgHashItem := range msgHash {
		msgHashRule = append(msgHashRule, msgHashItem)
	}

	logs, sub, err := _GnosisSafe130.contract.WatchLogs(opts, "SignMsg", msgHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GnosisSafe130SignMsg)
				if err := _GnosisSafe130.contract.UnpackLog(event, "SignMsg", log); err != nil {
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

// ParseSignMsg is a log parse operation binding the contract event 0xe7f4675038f4f6034dfcbbb24c4dc08e4ebf10eb9d257d3d02c0f38d122ac6e4.
//
// Solidity: event SignMsg(bytes32 indexed msgHash)
func (_GnosisSafe130 *GnosisSafe130Filterer) ParseSignMsg(log types.Log) (*GnosisSafe130SignMsg, error) {
	event := new(GnosisSafe130SignMsg)
	if err := _GnosisSafe130.contract.UnpackLog(event, "SignMsg", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
