// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_request_id_v08

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

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

var VRFRequestIDBaseTestHelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_vRFInputSeed\",\"type\":\"uint256\"}],\"name\":\"makeRequestId_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_userSeed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"makeVRFInputSeed_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610170806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806337ab429a1461003b578063bda087ae146100af575b600080fd5b61009d61004936600461010b565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b60405190815260200160405180910390f35b61009d6100bd3660046100e9565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b600080604083850312156100fc57600080fd5b50508035926020909101359150565b6000806000806080858703121561012157600080fd5b8435935060208501359250604085013573ffffffffffffffffffffffffffffffffffffffff8116811461015357600080fd5b939692955092936060013592505056fea164736f6c6343000806000a",
}

var VRFRequestIDBaseTestHelperABI = VRFRequestIDBaseTestHelperMetaData.ABI

var VRFRequestIDBaseTestHelperBin = VRFRequestIDBaseTestHelperMetaData.Bin

func DeployVRFRequestIDBaseTestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFRequestIDBaseTestHelper, error) {
	parsed, err := VRFRequestIDBaseTestHelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFRequestIDBaseTestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFRequestIDBaseTestHelper{address: address, abi: *parsed, VRFRequestIDBaseTestHelperCaller: VRFRequestIDBaseTestHelperCaller{contract: contract}, VRFRequestIDBaseTestHelperTransactor: VRFRequestIDBaseTestHelperTransactor{contract: contract}, VRFRequestIDBaseTestHelperFilterer: VRFRequestIDBaseTestHelperFilterer{contract: contract}}, nil
}

type VRFRequestIDBaseTestHelper struct {
	address common.Address
	abi     abi.ABI
	VRFRequestIDBaseTestHelperCaller
	VRFRequestIDBaseTestHelperTransactor
	VRFRequestIDBaseTestHelperFilterer
}

type VRFRequestIDBaseTestHelperCaller struct {
	contract *bind.BoundContract
}

type VRFRequestIDBaseTestHelperTransactor struct {
	contract *bind.BoundContract
}

type VRFRequestIDBaseTestHelperFilterer struct {
	contract *bind.BoundContract
}

type VRFRequestIDBaseTestHelperSession struct {
	Contract     *VRFRequestIDBaseTestHelper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFRequestIDBaseTestHelperCallerSession struct {
	Contract *VRFRequestIDBaseTestHelperCaller
	CallOpts bind.CallOpts
}

type VRFRequestIDBaseTestHelperTransactorSession struct {
	Contract     *VRFRequestIDBaseTestHelperTransactor
	TransactOpts bind.TransactOpts
}

type VRFRequestIDBaseTestHelperRaw struct {
	Contract *VRFRequestIDBaseTestHelper
}

type VRFRequestIDBaseTestHelperCallerRaw struct {
	Contract *VRFRequestIDBaseTestHelperCaller
}

type VRFRequestIDBaseTestHelperTransactorRaw struct {
	Contract *VRFRequestIDBaseTestHelperTransactor
}

func NewVRFRequestIDBaseTestHelper(address common.Address, backend bind.ContractBackend) (*VRFRequestIDBaseTestHelper, error) {
	abi, err := abi.JSON(strings.NewReader(VRFRequestIDBaseTestHelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFRequestIDBaseTestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelper{address: address, abi: abi, VRFRequestIDBaseTestHelperCaller: VRFRequestIDBaseTestHelperCaller{contract: contract}, VRFRequestIDBaseTestHelperTransactor: VRFRequestIDBaseTestHelperTransactor{contract: contract}, VRFRequestIDBaseTestHelperFilterer: VRFRequestIDBaseTestHelperFilterer{contract: contract}}, nil
}

func NewVRFRequestIDBaseTestHelperCaller(address common.Address, caller bind.ContractCaller) (*VRFRequestIDBaseTestHelperCaller, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperCaller{contract: contract}, nil
}

func NewVRFRequestIDBaseTestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFRequestIDBaseTestHelperTransactor, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperTransactor{contract: contract}, nil
}

func NewVRFRequestIDBaseTestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFRequestIDBaseTestHelperFilterer, error) {
	contract, err := bindVRFRequestIDBaseTestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTestHelperFilterer{contract: contract}, nil
}

func bindVRFRequestIDBaseTestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFRequestIDBaseTestHelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperCaller.contract.Call(opts, result, method, params...)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperTransactor.contract.Transfer(opts)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.VRFRequestIDBaseTestHelperTransactor.contract.Transact(opts, method, params...)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Transfer(opts)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBaseTestHelper.Contract.contract.Transact(opts, method, params...)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCaller) MakeRequestId(opts *bind.CallOpts, _keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFRequestIDBaseTestHelper.contract.Call(opts, &out, "makeRequestId_", _keyHash, _vRFInputSeed)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeRequestId(_keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeRequestId(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _vRFInputSeed)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCaller) MakeVRFInputSeed(opts *bind.CallOpts, _keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFRequestIDBaseTestHelper.contract.Call(opts, &out, "makeVRFInputSeed_", _keyHash, _userSeed, _requester, _nonce)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelperCallerSession) MakeVRFInputSeed(_keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error) {
	return _VRFRequestIDBaseTestHelper.Contract.MakeVRFInputSeed(&_VRFRequestIDBaseTestHelper.CallOpts, _keyHash, _userSeed, _requester, _nonce)
}

func (_VRFRequestIDBaseTestHelper *VRFRequestIDBaseTestHelper) Address() common.Address {
	return _VRFRequestIDBaseTestHelper.address
}

type VRFRequestIDBaseTestHelperInterface interface {
	MakeRequestId(opts *bind.CallOpts, _keyHash [32]byte, _vRFInputSeed *big.Int) ([32]byte, error)

	MakeVRFInputSeed(opts *bind.CallOpts, _keyHash [32]byte, _userSeed *big.Int, _requester common.Address, _nonce *big.Int) (*big.Int, error)

	Address() common.Address
}
