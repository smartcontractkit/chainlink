// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_testnet_d20

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

// VRFTestnetD20ABI is the input ABI used to generate the binding from.
const VRFTestnetD20ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"d20Results\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoll\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"d20result\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_seed\",\"type\":\"uint256\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userProvidedSeed\",\"type\":\"uint256\"}],\"name\":\"rollDice\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VRFTestnetD20Bin is the compiled bytecode used for deploying new contracts.
var VRFTestnetD20Bin = "0x608060405234801561001057600080fd5b506040516107b53803806107b58339818101604052606081101561003357600080fd5b508051602082015160409092015160018054600080546001600160a01b039586166001600160a01b031993841681178416179093559390941690841681179093169092179055600455670de0b6b3a764000060055561071e806100976000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063acfff37711610050578063acfff377146100e8578063ae383a4d14610105578063dc6cfe101461010d57610072565b80634ab5fc501461007757806394985ddd146100a65780639e317f12146100cb575b600080fd5b6100946004803603602081101561008d57600080fd5b5035610136565b60408051918252519081900360200190f35b6100c9600480360360408110156100bc57600080fd5b5080359060200135610154565b005b610094600480360360208110156100e157600080fd5b50356101e8565b610094600480360360208110156100fe57600080fd5b50356101fa565b61009461030d565b6100946004803603606081101561012357600080fd5b5080359060208101359060400135610351565b6003818154811061014357fe5b600091825260209091200154905081565b60015473ffffffffffffffffffffffffffffffffffffffff1633146101da57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015290519081900360640190fd5b6101e482826104ee565b5050565b60026020526000908152604090205481565b60055460008054604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905192939273ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b15801561027257600080fd5b505afa158015610286573d6000803e3d6000fd5b505050506040513d602081101561029c57600080fd5b505110156102f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b8152602001806106be602b913960400191505060405180910390fd5b600061030660045460055485610351565b9392505050565b60038054600091907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810190811061034157fe5b9060005260206000200154905090565b600080546001546040805160208082018990528183018790528251808303840181526060808401948590527f4000aea00000000000000000000000000000000000000000000000000000000090945273ffffffffffffffffffffffffffffffffffffffff94851660648401818152608485018b905260a48501958652825160c486015282519690971696634000aea09691958b9593949193909260e4909101918501908083838d5b838110156104115781810151838201526020016103f9565b50505050905090810190601f16801561043e5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561045f57600080fd5b505af1158015610473573d6000803e3d6000fd5b505050506040513d602081101561048957600080fd5b50506000848152600260205260408120546104a99086908590309061054a565b6000868152600260205260409020549091506104cc90600163ffffffff61059e16565b6000868152600260205260409020556104e58582610612565b95945050505050565b6000610512600161050684601463ffffffff61063e16565b9063ffffffff61059e16565b600380546001810182556000919091527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b0155505050565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60008282018381101561030657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b6000816106ac57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f536166654d6174683a206d6f64756c6f206279207a65726f0000000000000000604482015290519081900360640190fd5b8183816106b557fe5b06939250505056fe4e6f7420656e6f756768204c494e4b202d2066696c6c20636f6e7472616374207769746820666175636574a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

// DeployVRFTestnetD20 deploys a new Ethereum contract, binding an instance of VRFTestnetD20 to it.
func DeployVRFTestnetD20(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address, _keyHash [32]byte) (common.Address, *types.Transaction, *VRFTestnetD20, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestnetD20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFTestnetD20Bin), backend, _vrfCoordinator, _link, _keyHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFTestnetD20{VRFTestnetD20Caller: VRFTestnetD20Caller{contract: contract}, VRFTestnetD20Transactor: VRFTestnetD20Transactor{contract: contract}, VRFTestnetD20Filterer: VRFTestnetD20Filterer{contract: contract}}, nil
}

// VRFTestnetD20 is an auto generated Go binding around an Ethereum contract.
type VRFTestnetD20 struct {
	VRFTestnetD20Caller     // Read-only binding to the contract
	VRFTestnetD20Transactor // Write-only binding to the contract
	VRFTestnetD20Filterer   // Log filterer for contract events
}

// VRFTestnetD20Caller is an auto generated read-only Go binding around an Ethereum contract.
type VRFTestnetD20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestnetD20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFTestnetD20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestnetD20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFTestnetD20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFTestnetD20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFTestnetD20Session struct {
	Contract     *VRFTestnetD20    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFTestnetD20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFTestnetD20CallerSession struct {
	Contract *VRFTestnetD20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VRFTestnetD20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFTestnetD20TransactorSession struct {
	Contract     *VRFTestnetD20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VRFTestnetD20Raw is an auto generated low-level Go binding around an Ethereum contract.
type VRFTestnetD20Raw struct {
	Contract *VRFTestnetD20 // Generic contract binding to access the raw methods on
}

// VRFTestnetD20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFTestnetD20CallerRaw struct {
	Contract *VRFTestnetD20Caller // Generic read-only contract binding to access the raw methods on
}

// VRFTestnetD20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFTestnetD20TransactorRaw struct {
	Contract *VRFTestnetD20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFTestnetD20 creates a new instance of VRFTestnetD20, bound to a specific deployed contract.
func NewVRFTestnetD20(address common.Address, backend bind.ContractBackend) (*VRFTestnetD20, error) {
	contract, err := bindVRFTestnetD20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFTestnetD20{VRFTestnetD20Caller: VRFTestnetD20Caller{contract: contract}, VRFTestnetD20Transactor: VRFTestnetD20Transactor{contract: contract}, VRFTestnetD20Filterer: VRFTestnetD20Filterer{contract: contract}}, nil
}

// NewVRFTestnetD20Caller creates a new read-only instance of VRFTestnetD20, bound to a specific deployed contract.
func NewVRFTestnetD20Caller(address common.Address, caller bind.ContractCaller) (*VRFTestnetD20Caller, error) {
	contract, err := bindVRFTestnetD20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestnetD20Caller{contract: contract}, nil
}

// NewVRFTestnetD20Transactor creates a new write-only instance of VRFTestnetD20, bound to a specific deployed contract.
func NewVRFTestnetD20Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFTestnetD20Transactor, error) {
	contract, err := bindVRFTestnetD20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTestnetD20Transactor{contract: contract}, nil
}

// NewVRFTestnetD20Filterer creates a new log filterer instance of VRFTestnetD20, bound to a specific deployed contract.
func NewVRFTestnetD20Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFTestnetD20Filterer, error) {
	contract, err := bindVRFTestnetD20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFTestnetD20Filterer{contract: contract}, nil
}

// bindVRFTestnetD20 binds a generic wrapper to an already deployed contract.
func bindVRFTestnetD20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFTestnetD20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestnetD20 *VRFTestnetD20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestnetD20.Contract.VRFTestnetD20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestnetD20 *VRFTestnetD20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.VRFTestnetD20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestnetD20 *VRFTestnetD20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.VRFTestnetD20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFTestnetD20 *VRFTestnetD20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFTestnetD20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFTestnetD20 *VRFTestnetD20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFTestnetD20 *VRFTestnetD20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.contract.Transact(opts, method, params...)
}

// D20Results is a free data retrieval call binding the contract method 0x4ab5fc50.
//
// Solidity: function d20Results(uint256 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20Caller) D20Results(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestnetD20.contract.Call(opts, out, "d20Results", arg0)
	return *ret0, err
}

// D20Results is a free data retrieval call binding the contract method 0x4ab5fc50.
//
// Solidity: function d20Results(uint256 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20Session) D20Results(arg0 *big.Int) (*big.Int, error) {
	return _VRFTestnetD20.Contract.D20Results(&_VRFTestnetD20.CallOpts, arg0)
}

// D20Results is a free data retrieval call binding the contract method 0x4ab5fc50.
//
// Solidity: function d20Results(uint256 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20CallerSession) D20Results(arg0 *big.Int) (*big.Int, error) {
	return _VRFTestnetD20.Contract.D20Results(&_VRFTestnetD20.CallOpts, arg0)
}

// LatestRoll is a free data retrieval call binding the contract method 0xae383a4d.
//
// Solidity: function latestRoll() view returns(uint256 d20result)
func (_VRFTestnetD20 *VRFTestnetD20Caller) LatestRoll(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestnetD20.contract.Call(opts, out, "latestRoll")
	return *ret0, err
}

// LatestRoll is a free data retrieval call binding the contract method 0xae383a4d.
//
// Solidity: function latestRoll() view returns(uint256 d20result)
func (_VRFTestnetD20 *VRFTestnetD20Session) LatestRoll() (*big.Int, error) {
	return _VRFTestnetD20.Contract.LatestRoll(&_VRFTestnetD20.CallOpts)
}

// LatestRoll is a free data retrieval call binding the contract method 0xae383a4d.
//
// Solidity: function latestRoll() view returns(uint256 d20result)
func (_VRFTestnetD20 *VRFTestnetD20CallerSession) LatestRoll() (*big.Int, error) {
	return _VRFTestnetD20.Contract.LatestRoll(&_VRFTestnetD20.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20Caller) Nonces(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFTestnetD20.contract.Call(opts, out, "nonces", arg0)
	return *ret0, err
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20Session) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFTestnetD20.Contract.Nonces(&_VRFTestnetD20.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x9e317f12.
//
// Solidity: function nonces(bytes32 ) view returns(uint256)
func (_VRFTestnetD20 *VRFTestnetD20CallerSession) Nonces(arg0 [32]byte) (*big.Int, error) {
	return _VRFTestnetD20.Contract.Nonces(&_VRFTestnetD20.CallOpts, arg0)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFTestnetD20 *VRFTestnetD20Transactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFTestnetD20 *VRFTestnetD20Session) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RawFulfillRandomness(&_VRFTestnetD20.TransactOpts, requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFTestnetD20 *VRFTestnetD20TransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RawFulfillRandomness(&_VRFTestnetD20.TransactOpts, requestId, randomness)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20Transactor) RequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.contract.Transact(opts, "requestRandomness", _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20Session) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RequestRandomness(&_VRFTestnetD20.TransactOpts, _keyHash, _fee, _seed)
}

// RequestRandomness is a paid mutator transaction binding the contract method 0xdc6cfe10.
//
// Solidity: function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20TransactorSession) RequestRandomness(_keyHash [32]byte, _fee *big.Int, _seed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RequestRandomness(&_VRFTestnetD20.TransactOpts, _keyHash, _fee, _seed)
}

// RollDice is a paid mutator transaction binding the contract method 0xacfff377.
//
// Solidity: function rollDice(uint256 userProvidedSeed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20Transactor) RollDice(opts *bind.TransactOpts, userProvidedSeed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.contract.Transact(opts, "rollDice", userProvidedSeed)
}

// RollDice is a paid mutator transaction binding the contract method 0xacfff377.
//
// Solidity: function rollDice(uint256 userProvidedSeed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20Session) RollDice(userProvidedSeed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RollDice(&_VRFTestnetD20.TransactOpts, userProvidedSeed)
}

// RollDice is a paid mutator transaction binding the contract method 0xacfff377.
//
// Solidity: function rollDice(uint256 userProvidedSeed) returns(bytes32 requestId)
func (_VRFTestnetD20 *VRFTestnetD20TransactorSession) RollDice(userProvidedSeed *big.Int) (*types.Transaction, error) {
	return _VRFTestnetD20.Contract.RollDice(&_VRFTestnetD20.TransactOpts, userProvidedSeed)
}
