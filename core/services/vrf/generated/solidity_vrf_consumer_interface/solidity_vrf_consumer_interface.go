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
const VRFConsumerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_randomness\",\"type\":\"uint256\"}],\"name\":\"fulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VRFConsumerBin is the compiled bytecode used for deploying new contracts.
var VRFConsumerBin = "0x608060405234801561001057600080fd5b506040516104663803806104668339818101604052604081101561003357600080fd5b508051602090910151600180546001600160a01b039384166001600160a01b031991821617909155600080549390921692169190911790556103ec8061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100665760003560e01c80632f47fd86116100505780632f47fd86146100aa5780639e317f12146100b2578063dc6cfe10146100cf57610066565b80626d6cae1461006b5780631f1f897f14610085575b600080fd5b6100736100f8565b60408051918252519081900360200190f35b6100a86004803603604081101561009b57600080fd5b50803590602001356100fe565b005b610073610106565b610073600480360360208110156100c857600080fd5b503561010c565b610073600480360360608110156100e557600080fd5b508035906020810135906040013561011e565b60045481565b600355600455565b60035481565b60026020526000908152604090205481565b600080546001546040805160208082018990528183018790528251808303840181526060808401948590527f4000aea00000000000000000000000000000000000000000000000000000000090945273ffffffffffffffffffffffffffffffffffffffff94851660648401818152608485018b905260a48501958652825160c486015282519690971696634000aea09691958b9593949193909260e4909101918501908083838d5b838110156101de5781810151838201526020016101c6565b50505050905090810190601f16801561020b5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561022c57600080fd5b505af1158015610240573d6000803e3d6000fd5b505050506040513d602081101561025657600080fd5b5050600084815260026020526040812054610276908690859030906102bb565b60008681526002602052604090205490915061029990600163ffffffff61030f16565b6000868152600260205260409020556102b2858261038a565b95945050505050565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60008282018381101561038357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b60408051602080820194909452808201929092528051808303820181526060909201905280519101209056fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"

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

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) constant returns(uint256)
func (_VRFConsumer *VRFConsumerCaller) Nonces(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFConsumer.contract.Call(opts, out, "nonces", arg0)
	return *ret0, err
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) constant returns(uint256)
func (_VRFConsumer *VRFConsumerSession) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFConsumer.Contract.Nonces(&_VRFConsumer.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) constant returns(uint256)
func (_VRFConsumer *VRFConsumerCallerSession) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFConsumer.Contract.Nonces(&_VRFConsumer.CallOpts, arg0)
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
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerTransactor) RequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "requestRandomness", _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerSession) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerTransactorSession) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}
