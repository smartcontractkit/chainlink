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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// VRFConsumerABI is the input ABI used to generate the binding from.
const VRFConsumerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VRFConsumerBin is the compiled bytecode used for deploying new contracts.
var VRFConsumerBin = "0x60c060405234801561001057600080fd5b506040516105823803806105828339818101604052604081101561003357600080fd5b5080516020909101516001600160601b0319606092831b811660a052911b1660805260805160601c60a05160601c6104ff61008360003980610131528061021d5250806101e152506104ff6000f3fe608060405234801561001057600080fd5b50600436106100665760003560e01c80638a5f002b116100505780638a5f002b1461008d57806394985ddd146100b65780639e317f12146100db57610066565b80626d6cae1461006b5780632f47fd8614610085575b600080fd5b6100736100f8565b60408051918252519081900360200190f35b6100736100fe565b610073600480360360608110156100a357600080fd5b5080359060208101359060400135610104565b6100d9600480360360408110156100cc57600080fd5b5080359060200135610119565b005b610073600480360360208110156100f157600080fd5b50356101cb565b60025481565b60015481565b60006101118484846101dd565b949350505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146101bd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015290519081900360640190fd5b6101c782826103c6565b5050565b60006020819052908152604090205481565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f000000000000000000000000000000000000000000000000000000000000000085878660405160200180838152602001828152602001925050506040516020818303038152906040526040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102e95781810151838201526020016102d1565b50505050905090810190601f1680156103165780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561033757600080fd5b505af115801561034b573d6000803e3d6000fd5b505050506040513d602081101561036157600080fd5b5050600084815260208190526040812054610381908690859030906103ce565b6000868152602081905260409020549091506103a490600163ffffffff61042216565b6000868152602081905260409020556103bd858261049d565b95945050505050565b600155600255565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60008282018381101561049657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b60408051602080820194909452808201929092528051808303820181526060909201905280519101209056fea264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

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
// Solidity: function nonces(bytes32 ) view returns(uint256)
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
// Solidity: function nonces(bytes32 ) view returns(uint256)
func (_VRFConsumer *VRFConsumerSession) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFConsumer.Contract.Nonces(&_VRFConsumer.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) view returns(uint256)
func (_VRFConsumer *VRFConsumerCallerSession) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFConsumer.Contract.Nonces(&_VRFConsumer.CallOpts, arg0)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() view returns(uint256)
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
// Solidity: function randomnessOutput() view returns(uint256)
func (_VRFConsumer *VRFConsumerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RandomnessOutput is a free data retrieval call binding the contract method 0x2f47fd86.
//
// Solidity: function randomnessOutput() view returns(uint256)
func (_VRFConsumer *VRFConsumerCallerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() view returns(bytes32)
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
// Solidity: function requestId() view returns(bytes32)
func (_VRFConsumer *VRFConsumerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// RequestId is a free data retrieval call binding the contract method 0x006d6cae.
//
// Solidity: function requestId() view returns(bytes32)
func (_VRFConsumer *VRFConsumerCallerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumer *VRFConsumerTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x8a5f002b.
//
// Solidity: function testRequestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "testRequestRandomness", _keyHash, _fee, _seed)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x8a5f002b.
//
// Solidity: function testRequestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerSession) TestRequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}

// TestRequestRandomness is a paid mutator transaction binding the contract method 0x8a5f002b.
//
// Solidity: function testRequestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFConsumer *VRFConsumerTransactorSession) TestRequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee, _seed)
}
