// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_consumer_interface

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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VRFConsumerABI is the input ABI used to generate the binding from.
const VRFConsumerABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"name\":\"_randomness\",\"type\":\"uint256\"}],\"name\":\"fulfillRandomness\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"randomnessOutput\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"makeRequestId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"name\":\"_fee\",\"type\":\"uint256\"},{\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRandomness\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"name\":\"_link\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// VRFConsumerBin is the compiled bytecode used for deploying new contracts.
var VRFConsumerBin = "0x608060405234801561001057600080fd5b506040516040806104d38339810180604052604081101561003057600080fd5b810190808051906020019092919080519060200190929190505050818181600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050506103f2806100e16000396000f3fe60806040526004361061006c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680626d6cae146100715780631f1f897f1461009c5780632f47fd86146100e1578063916bf6c71461010c578063dc6cfe1014610165575b600080fd5b34801561007d57600080fd5b506100866101b4565b6040518082815260200191505060405180910390f35b3480156100a857600080fd5b506100df600480360360408110156100bf57600080fd5b8101908080359060200190929190803590602001909291905050506101ba565b005b3480156100ed57600080fd5b506100f66101cc565b6040518082815260200191505060405180910390f35b34801561011857600080fd5b5061014f6004803603604081101561012f57600080fd5b8101908080359060200190929190803590602001909291905050506101d2565b6040518082815260200191505060405180910390f35b34801561017157600080fd5b506101b26004803603606081101561018857600080fd5b8101908080359060200190929190803590602001909291908035906020019092919050505061020b565b005b60035481565b80600281905550816003819055505050565b60025481565b60008282604051602001808381526020018281526020019250505060405160208183030381529060405280519060200120905092915050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634000aea0600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684868560405160200180838152602001828152602001925050506040516020818303038152906040526040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561033757808201518184015260208101905061031c565b50505050905090810190601f1680156103645780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561038557600080fd5b505af1158015610399573d6000803e3d6000fd5b505050506040513d60208110156103af57600080fd5b81019080805190602001909291905050505050505056fea165627a7a7230582093113798ec9eef072d65f19f21295baa80ce2abd50c2c165211606fdc28f19960029"

// DeployVRFConsumer deploys a new Ethereum contract, binding an instance of VRFConsumer to it.
func DeployVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFConsumer, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFConsumerBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumer{VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

// VRFConsumer is an auto generated Go binding around an Ethereum contract.
type VRFConsumer struct {
	VRFConsumerCaller     // Read-only binding to the contract
	VRFConsumerTransactor // Write-only binding to the contract
	VRFConsumerFilterer   // Log filterer for contract events
}

// VRFConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFConsumerSession struct {
	Contract     *VRFConsumer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFConsumerCallerSession struct {
	Contract *VRFConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// VRFConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFConsumerTransactorSession struct {
	Contract     *VRFConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// VRFConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFConsumerRaw struct {
	Contract *VRFConsumer // Generic contract binding to access the raw methods on
}

// VRFConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFConsumerCallerRaw struct {
	Contract *VRFConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// VRFConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFConsumerTransactorRaw struct {
	Contract *VRFConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFConsumer creates a new instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumer(address common.Address, backend bind.ContractBackend) (*VRFConsumer, error) {
	contract, err := bindVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumer{VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

// NewVRFConsumerCaller creates a new read-only instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerCaller, error) {
	contract, err := bindVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerCaller{contract: contract}, nil
}

// NewVRFConsumerTransactor creates a new write-only instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerTransactor, error) {
	contract, err := bindVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerTransactor{contract: contract}, nil
}

// NewVRFConsumerFilterer creates a new log filterer instance of VRFConsumer, bound to a specific deployed contract.
func NewVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerFilterer, error) {
	contract, err := bindVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerFilterer{contract: contract}, nil
}

// bindVRFConsumer binds a generic wrapper to an already deployed contract.
func bindVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumer *VRFConsumerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.VRFConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumer *VRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumer *VRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumer *VRFConsumerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumer *VRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumer *VRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transact(opts, method, params...)
}

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFConsumer *VRFConsumerCaller) MakeRequestId(opts *bind.CallOpts, _keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRFConsumer.contract.Call(opts, out, "makeRequestId", _keyHash, _seed)
	return *ret0, err
}

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFConsumer *VRFConsumerSession) MakeRequestId(_keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	return _VRFConsumer.Contract.MakeRequestId(&_VRFConsumer.CallOpts, _keyHash, _seed)
}

// MakeRequestId is a free data retrieval call binding the contract method 0x916bf6c7.
//
// Solidity: function makeRequestId(bytes32 _keyHash, uint256 _seed) constant returns(bytes32)
func (_VRFConsumer *VRFConsumerCallerSession) MakeRequestId(_keyHash [32]byte, _seed *big.Int) ([32]byte, error) {
	return _VRFConsumer.Contract.MakeRequestId(&_VRFConsumer.CallOpts, _keyHash, _seed)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() constant returns(uint256)
func (_VRFConsumer *VRFConsumerCaller) RandomnessOutput(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFConsumer.contract.Call(opts, out, "randomnessOutput")
	return *ret0, err
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() constant returns(uint256)
func (_VRFConsumer *VRFConsumerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() constant returns(uint256)
func (_VRFConsumer *VRFConsumerCallerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() constant returns(bytes32)
func (_VRFConsumer *VRFConsumerCaller) RequestId(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRFConsumer.contract.Call(opts, out, "requestId")
	return *ret0, err
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() constant returns(bytes32)
func (_VRFConsumer *VRFConsumerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() constant returns(bytes32)
func (_VRFConsumer *VRFConsumerCallerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// FulfillRandomness is a paid mutator transaction binding the contract method 0x1f1f897f.
//
// Solidity: function fulfillRandomness(bytes32 _requestId, uint256 _randomness) returns()
func (_VRFConsumer *VRFConsumerTransactor) FulfillRandomness(opts *bind.TransactOpts, _requestId [32]byte, _randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "fulfillRandomness", _requestId, _randomness)
}

// FulfillRandomness is a paid mutator transaction binding the contract method 0x1f1f897f.
//
// Solidity: function fulfillRandomness(bytes32 _requestId, uint256 _randomness) returns()
func (_VRFConsumer *VRFConsumerSession) FulfillRandomness(_requestId [32]byte, _randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.FulfillRandomness(&_VRFConsumer.TransactOpts, _requestId, _randomness)
}

// FulfillRandomness is a paid mutator transaction binding the contract method 0x1f1f897f.
//
// Solidity: function fulfillRandomness(bytes32 _requestId, uint256 _randomness) returns()
func (_VRFConsumer *VRFConsumerTransactorSession) FulfillRandomness(_requestId [32]byte, _randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.FulfillRandomness(&_VRFConsumer.TransactOpts, _requestId, _randomness)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns()
func (_VRFConsumer *VRFConsumerTransactor) RequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "requestRandomness", _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns()
func (_VRFConsumer *VRFConsumerSession) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns()
func (_VRFConsumer *VRFConsumerTransactorSession) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}
