// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_request_id

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

// VRFRequestIDBaseTestHelperABI is the input ABI used to generate the binding from.
const VRFRequestIDBaseTestHelperABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"name\":\"_userSeed\",\"type\":\"uint256\"},{\"name\":\"_requester\",\"type\":\"address\"},{\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"makeVRFInputSeed_\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"name\":\"_vRFInputSeed\",\"type\":\"uint256\"}],\"name\":\"makeRequestId_\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFRequestIDBaseTestHelperBin is the compiled bytecode used for deploying new contracts.
var VRFRequestIDBaseTestHelperBin = "0x608060405234801561001057600080fd5b50610239806100206000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806337ab429a14610051578063bda087ae146100d4575b600080fd5b34801561005d57600080fd5b506100be6004803603608081101561007457600080fd5b810190808035906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061012d565b6040518082815260200191505060405180910390f35b3480156100e057600080fd5b50610117600480360360408110156100f757600080fd5b810190808035906020019092919080359060200190929190505050610145565b6040518082815260200191505060405180910390f35b600061013b85858585610159565b9050949350505050565b600061015183836101d4565b905092915050565b600084848484604051602001808581526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200194505050505060405160208183030381529060405280519060200120600190049050949350505050565b6000828260405160200180838152602001828152602001925050506040516020818303038152906040528051906020012090509291505056fea165627a7a7230582089f764eb1c43d139339fda73af1a473feb1897b612bc02aacd77b95796e80a830029"

// DeployVRFRequestIDBaseTestHelper deploys a new Ethereum contract, binding an instance of VRFRequestIDBaseTestHelper to it.
func DeployVRFRequestIDBaseTestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFRequestIDBaseTestHelper, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFRequestIDBaseTestHelperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFRequestIDBaseTestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFRequestIDBaseTestHelper{VRFRequestIDBaseTestHelperCaller: VRFRequestIDBaseTestHelperCaller{contract: contract}, VRFRequestIDBaseTestHelperTransactor: VRFRequestIDBaseTestHelperTransactor{contract: contract}, VRFRequestIDBaseTestHelperFilterer: VRFRequestIDBaseTestHelperFilterer{contract: contract}}, nil
}

// VRFRequestIDBaseTestHelper is an auto generated Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelper struct {
	VRFRequestIDBaseTestHelperCaller     // Read-only binding to the contract
	VRFRequestIDBaseTestHelperTransactor // Write-only binding to the contract
	VRFRequestIDBaseTestHelperFilterer   // Log filterer for contract events
}

// VRFRequestIDBaseTestHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseTestHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseTestHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFRequestIDBaseTestHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseTestHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFRequestIDBaseTestHelperSession struct {
	Contract     *VRFRequestIDBaseTestHelper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// VRFRequestIDBaseTestHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFRequestIDBaseTestHelperCallerSession struct {
	Contract *VRFRequestIDBaseTestHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// VRFRequestIDBaseTestHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFRequestIDBaseTestHelperTransactorSession struct {
	Contract     *VRFRequestIDBaseTestHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// VRFRequestIDBaseTestHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelperRaw struct {
	Contract *VRFRequestIDBaseTestHelper // Generic contract binding to access the raw methods on
}

// VRFRequestIDBaseTestHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelperCallerRaw struct {
	Contract *VRFRequestIDBaseTestHelperCaller // Generic read-only contract binding to access the raw methods on
}

// VRFRequestIDBaseTestHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTestHelperTransactorRaw struct {
	Contract *VRFRequestIDBaseTestHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFRequestIDBaseTestHelper creates a new instance of VRFRequestIDBaseTestHelper, bound to a specific deployed contract.
func NewVRFRequestIDBaseTestHelper(address common.Address, backend bind.ContractBackend) (*VRFRequestIDBaseTestHelper, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelper{VRFRequestIDBaseTestHelperCaller: VRFRequestIDBaseTestHelperCaller{contract: contract}, VRFRequestIDBaseTestHelperTransactor: VRFRequestIDBaseTestHelperTransactor{contract: contract}, VRFRequestIDBaseTestHelperFilterer: VRFRequestIDBaseTestHelperFilterer{contract: contract}}, nil
}

// NewVRFRequestIDBaseTestHelperCaller creates a new read-only instance of VRFRequestIDBaseTestHelper, bound to a specific deployed contract.
func NewVRFRequestIDBaseTestHelperCaller(address common.Address, caller bind.ContractCaller) (*VRFRequestIDBaseTestHelperCaller, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperCaller{contract: contract}, nil
}

// NewVRFRequestIDBaseTestHelperTransactor creates a new write-only instance of VRFRequestIDBaseTestHelper, bound to a specific deployed contract.
func NewVRFRequestIDBaseTestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFRequestIDBaseTestHelperTransactor, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperTransactor{contract: contract}, nil
}

// NewVRFRequestIDBaseTestHelperFilterer creates a new log filterer instance of VRFRequestIDBaseTestHelper, bound to a specific deployed contract.
func NewVRFRequestIDBaseTestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFRequestIDBaseTestHelperFilterer, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperFilterer{contract: contract}, nil
}

// bindVRFRequestIDBaseTestHelper binds a generic wrapper to an already deployed contract.
func bindVRFRequestIDBaseTestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFRequestIDBaseTestHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Transact(opts, method, params...)
}

// MakeRequestId is a free data retrieval call binding the contract method 0xbda087ae.
//
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) constant returns(bytes32)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCaller) MakeRequestId(opts *bind.CallOpts, _keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _VRFRequestIDBaseTestHelper.contract.Call(opts, out, "makeRequestId_", _keyHash, _vRFInputSeed)
	return *ret0, err
}

// MakeRequestId is a free data retrieval call binding the contract method 0xbda087ae.
//
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) constant returns(bytes32)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

// MakeRequestId is a free data retrieval call binding the contract method 0xbda087ae.
//
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) constant returns(bytes32)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

// MakeVRFInputSeed is a free data retrieval call binding the contract method 0x37ab429a.
//
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) constant returns(uint256)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCaller) MakeVRFInputSeed(opts *bind.CallOpts, _keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _VRFRequestIDBaseTestHelper.contract.Call(opts, out, "makeVRFInputSeed_", _keyHash, _userSeed, _requester, _nonce)
	return *ret0, err
}

// MakeVRFInputSeed is a free data retrieval call binding the contract method 0x37ab429a.
//
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) constant returns(uint256)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}

// MakeVRFInputSeed is a free data retrieval call binding the contract method 0x37ab429a.
//
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) constant returns(uint256)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}
