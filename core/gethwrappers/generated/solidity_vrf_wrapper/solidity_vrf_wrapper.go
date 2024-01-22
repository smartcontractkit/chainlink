// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_wrapper

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

var VRFMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"PROOF_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060588061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063e911439c14602d575b600080fd5b60336045565b60408051918252519081900360200190f35b6101a08156fea164736f6c6343000606000a",
}

var VRFABI = VRFMetaData.ABI

var VRFBin = VRFMetaData.Bin

func DeployVRF(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRF, error) {
	parsed, err := VRFMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRF{address: address, abi: *parsed, VRFCaller: VRFCaller{contract: contract}, VRFTransactor: VRFTransactor{contract: contract}, VRFFilterer: VRFFilterer{contract: contract}}, nil
}

type VRF struct {
	address common.Address
	abi     abi.ABI
	VRFCaller
	VRFTransactor
	VRFFilterer
}

type VRFCaller struct {
	contract *bind.BoundContract
}

type VRFTransactor struct {
	contract *bind.BoundContract
}

type VRFFilterer struct {
	contract *bind.BoundContract
}

type VRFSession struct {
	Contract     *VRF
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCallerSession struct {
	Contract *VRFCaller
	CallOpts bind.CallOpts
}

type VRFTransactorSession struct {
	Contract     *VRFTransactor
	TransactOpts bind.TransactOpts
}

type VRFRaw struct {
	Contract *VRF
}

type VRFCallerRaw struct {
	Contract *VRFCaller
}

type VRFTransactorRaw struct {
	Contract *VRFTransactor
}

func NewVRF(address common.Address, backend bind.ContractBackend) (*VRF, error) {
	abi, err := abi.JSON(strings.NewReader(VRFABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRF(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRF{address: address, abi: abi, VRFCaller: VRFCaller{contract: contract}, VRFTransactor: VRFTransactor{contract: contract}, VRFFilterer: VRFFilterer{contract: contract}}, nil
}

func NewVRFCaller(address common.Address, caller bind.ContractCaller) (*VRFCaller, error) {
	contract, err := bindVRF(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCaller{contract: contract}, nil
}

func NewVRFTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFTransactor, error) {
	contract, err := bindVRF(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFTransactor{contract: contract}, nil
}

func NewVRFFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFFilterer, error) {
	contract, err := bindVRF(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFFilterer{contract: contract}, nil
}

func bindVRF(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRF *VRFRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRF.Contract.VRFCaller.contract.Call(opts, result, method, params...)
}

func (_VRF *VRFRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRF.Contract.VRFTransactor.contract.Transfer(opts)
}

func (_VRF *VRFRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRF.Contract.VRFTransactor.contract.Transact(opts, method, params...)
}

func (_VRF *VRFCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRF.Contract.contract.Call(opts, result, method, params...)
}

func (_VRF *VRFTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRF.Contract.contract.Transfer(opts)
}

func (_VRF *VRFTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRF.Contract.contract.Transact(opts, method, params...)
}

func (_VRF *VRFCaller) PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRF.contract.Call(opts, &out, "PROOF_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRF *VRFSession) PROOFLENGTH() (*big.Int, error) {
	return _VRF.Contract.PROOFLENGTH(&_VRF.CallOpts)
}

func (_VRF *VRFCallerSession) PROOFLENGTH() (*big.Int, error) {
	return _VRF.Contract.PROOFLENGTH(&_VRF.CallOpts)
}

func (_VRF *VRF) Address() common.Address {
	return _VRF.address
}

type VRFInterface interface {
	PROOFLENGTH(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
