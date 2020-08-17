// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_request_id

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

// VRFRequestIDBaseTestHelperABI is the input ABI used to generate the binding from.
const VRFRequestIDBaseTestHelperABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_vRFInputSeed\",\"type\":\"uint256\"}],\"name\":\"makeRequestId_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_userSeed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"makeVRFInputSeed_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// VRFRequestIDBaseTestHelperBin is the compiled bytecode used for deploying new contracts.
var VRFRequestIDBaseTestHelperBin = "0x608060405234801561001057600080fd5b50610195806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806337ab429a1461003b578063bda087ae14610092575b600080fd5b6100806004803603608081101561005157600080fd5b5080359060208101359073ffffffffffffffffffffffffffffffffffffffff60408201351690606001356100b5565b60408051918252519081900360200190f35b610080600480360360408110156100a857600080fd5b50803590602001356100cc565b60006100c3858585856100df565b95945050505050565b60006100d88383610133565b9392505050565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60408051602080820194909452808201929092528051808303820181526060909201905280519101209056fea264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

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
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) pure returns(bytes32)
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
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) pure returns(bytes32)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

// MakeRequestId is a free data retrieval call binding the contract method 0xbda087ae.
//
// Solidity: function makeRequestId_(bytes32 _keyHash, uint256 _vRFInputSeed) pure returns(bytes32)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

// MakeVRFInputSeed is a free data retrieval call binding the contract method 0x37ab429a.
//
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) pure returns(uint256)
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
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) pure returns(uint256)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}

// MakeVRFInputSeed is a free data retrieval call binding the contract method 0x37ab429a.
//
// Solidity: function makeVRFInputSeed_(bytes32 _keyHash, uint256 _userSeed, address _requester, uint256 _nonce) pure returns(uint256)
func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}
