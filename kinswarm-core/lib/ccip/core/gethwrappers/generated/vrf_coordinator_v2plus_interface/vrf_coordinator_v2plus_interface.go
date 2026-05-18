// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2plus_interface

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

type IVRFCoordinatorV2PlusInternalProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

type IVRFCoordinatorV2PlusInternalRequestCommitment struct {
	BlockNum         uint64
	SubId            *big.Int
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	ExtraArgs        []byte
}

type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

var IVRFCoordinatorV2PlusInternalMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structIVRFCoordinatorV2PlusInternal.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structIVRFCoordinatorV2PlusInternal.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IVRFCoordinatorV2PlusInternalABI = IVRFCoordinatorV2PlusInternalMetaData.ABI

type IVRFCoordinatorV2PlusInternal struct {
	address common.Address
	abi     abi.ABI
	IVRFCoordinatorV2PlusInternalCaller
	IVRFCoordinatorV2PlusInternalTransactor
	IVRFCoordinatorV2PlusInternalFilterer
}

type IVRFCoordinatorV2PlusInternalCaller struct {
	contract *bind.BoundContract
}

type IVRFCoordinatorV2PlusInternalTransactor struct {
	contract *bind.BoundContract
}

type IVRFCoordinatorV2PlusInternalFilterer struct {
	contract *bind.BoundContract
}

type IVRFCoordinatorV2PlusInternalSession struct {
	Contract     *IVRFCoordinatorV2PlusInternal
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IVRFCoordinatorV2PlusInternalCallerSession struct {
	Contract *IVRFCoordinatorV2PlusInternalCaller
	CallOpts bind.CallOpts
}

type IVRFCoordinatorV2PlusInternalTransactorSession struct {
	Contract     *IVRFCoordinatorV2PlusInternalTransactor
	TransactOpts bind.TransactOpts
}

type IVRFCoordinatorV2PlusInternalRaw struct {
	Contract *IVRFCoordinatorV2PlusInternal
}

type IVRFCoordinatorV2PlusInternalCallerRaw struct {
	Contract *IVRFCoordinatorV2PlusInternalCaller
}

type IVRFCoordinatorV2PlusInternalTransactorRaw struct {
	Contract *IVRFCoordinatorV2PlusInternalTransactor
}

func NewIVRFCoordinatorV2PlusInternal(address common.Address, backend bind.ContractBackend) (*IVRFCoordinatorV2PlusInternal, error) {
	abi, err := abi.JSON(strings.NewReader(IVRFCoordinatorV2PlusInternalABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIVRFCoordinatorV2PlusInternal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternal{address: address, abi: abi, IVRFCoordinatorV2PlusInternalCaller: IVRFCoordinatorV2PlusInternalCaller{contract: contract}, IVRFCoordinatorV2PlusInternalTransactor: IVRFCoordinatorV2PlusInternalTransactor{contract: contract}, IVRFCoordinatorV2PlusInternalFilterer: IVRFCoordinatorV2PlusInternalFilterer{contract: contract}}, nil
}

func NewIVRFCoordinatorV2PlusInternalCaller(address common.Address, caller bind.ContractCaller) (*IVRFCoordinatorV2PlusInternalCaller, error) {
	contract, err := bindIVRFCoordinatorV2PlusInternal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternalCaller{contract: contract}, nil
}

func NewIVRFCoordinatorV2PlusInternalTransactor(address common.Address, transactor bind.ContractTransactor) (*IVRFCoordinatorV2PlusInternalTransactor, error) {
	contract, err := bindIVRFCoordinatorV2PlusInternal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternalTransactor{contract: contract}, nil
}

func NewIVRFCoordinatorV2PlusInternalFilterer(address common.Address, filterer bind.ContractFilterer) (*IVRFCoordinatorV2PlusInternalFilterer, error) {
	contract, err := bindIVRFCoordinatorV2PlusInternal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternalFilterer{contract: contract}, nil
}

func bindIVRFCoordinatorV2PlusInternal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IVRFCoordinatorV2PlusInternalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IVRFCoordinatorV2PlusInternal.Contract.IVRFCoordinatorV2PlusInternalCaller.contract.Call(opts, result, method, params...)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.IVRFCoordinatorV2PlusInternalTransactor.contract.Transfer(opts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.IVRFCoordinatorV2PlusInternalTransactor.contract.Transact(opts, method, params...)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IVRFCoordinatorV2PlusInternal.Contract.contract.Call(opts, result, method, params...)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.contract.Transfer(opts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.contract.Transact(opts, method, params...)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCaller) LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IVRFCoordinatorV2PlusInternal.contract.Call(opts, &out, "LINK_NATIVE_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) LINKNATIVEFEED() (common.Address, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.LINKNATIVEFEED(&_IVRFCoordinatorV2PlusInternal.CallOpts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerSession) LINKNATIVEFEED() (common.Address, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.LINKNATIVEFEED(&_IVRFCoordinatorV2PlusInternal.CallOpts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCaller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IVRFCoordinatorV2PlusInternal.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.GetActiveSubscriptionIds(&_IVRFCoordinatorV2PlusInternal.CallOpts, startIndex, maxCount)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.GetActiveSubscriptionIds(&_IVRFCoordinatorV2PlusInternal.CallOpts, startIndex, maxCount)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _IVRFCoordinatorV2PlusInternal.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.GetSubscription(&_IVRFCoordinatorV2PlusInternal.CallOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.GetSubscription(&_IVRFCoordinatorV2PlusInternal.CallOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCaller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _IVRFCoordinatorV2PlusInternal.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.PendingRequestExists(&_IVRFCoordinatorV2PlusInternal.CallOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.PendingRequestExists(&_IVRFCoordinatorV2PlusInternal.CallOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCaller) SRequestCommitments(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _IVRFCoordinatorV2PlusInternal.contract.Call(opts, &out, "s_requestCommitments", requestID)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) SRequestCommitments(requestID *big.Int) ([32]byte, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.SRequestCommitments(&_IVRFCoordinatorV2PlusInternal.CallOpts, requestID)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalCallerSession) SRequestCommitments(requestID *big.Int) ([32]byte, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.SRequestCommitments(&_IVRFCoordinatorV2PlusInternal.CallOpts, requestID)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.AcceptSubscriptionOwnerTransfer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.AcceptSubscriptionOwnerTransfer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.AddConsumer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.AddConsumer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.CancelSubscription(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, to)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.CancelSubscription(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, to)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "createSubscription")
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) CreateSubscription() (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.CreateSubscription(&_IVRFCoordinatorV2PlusInternal.TransactOpts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.CreateSubscription(&_IVRFCoordinatorV2PlusInternal.TransactOpts)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof IVRFCoordinatorV2PlusInternalProof, rc IVRFCoordinatorV2PlusInternalRequestCommitment, onlyPremium bool) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "fulfillRandomWords", proof, rc, onlyPremium)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) FulfillRandomWords(proof IVRFCoordinatorV2PlusInternalProof, rc IVRFCoordinatorV2PlusInternalRequestCommitment, onlyPremium bool) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.FulfillRandomWords(&_IVRFCoordinatorV2PlusInternal.TransactOpts, proof, rc, onlyPremium)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) FulfillRandomWords(proof IVRFCoordinatorV2PlusInternalProof, rc IVRFCoordinatorV2PlusInternalRequestCommitment, onlyPremium bool) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.FulfillRandomWords(&_IVRFCoordinatorV2PlusInternal.TransactOpts, proof, rc, onlyPremium)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "fundSubscriptionWithNative", subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.FundSubscriptionWithNative(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.FundSubscriptionWithNative(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RemoveConsumer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RemoveConsumer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, consumer)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "requestRandomWords", req)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RequestRandomWords(&_IVRFCoordinatorV2PlusInternal.TransactOpts, req)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RequestRandomWords(&_IVRFCoordinatorV2PlusInternal.TransactOpts, req)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RequestSubscriptionOwnerTransfer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, newOwner)
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _IVRFCoordinatorV2PlusInternal.Contract.RequestSubscriptionOwnerTransfer(&_IVRFCoordinatorV2PlusInternal.TransactOpts, subId, newOwner)
}

type IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator struct {
	Event *IVRFCoordinatorV2PlusInternalRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IVRFCoordinatorV2PlusInternalRandomWordsFulfilled)
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

	select {
	case log := <-it.logs:
		it.Event = new(IVRFCoordinatorV2PlusInternalRandomWordsFulfilled)
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

func (it *IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IVRFCoordinatorV2PlusInternalRandomWordsFulfilled struct {
	RequestId   *big.Int
	OutputSeed  *big.Int
	SubId       *big.Int
	Payment     *big.Int
	Success     bool
	OnlyPremium bool
	Raw         types.Log
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _IVRFCoordinatorV2PlusInternal.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator{contract: _IVRFCoordinatorV2PlusInternal.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *IVRFCoordinatorV2PlusInternalRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _IVRFCoordinatorV2PlusInternal.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IVRFCoordinatorV2PlusInternalRandomWordsFulfilled)
				if err := _IVRFCoordinatorV2PlusInternal.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) ParseRandomWordsFulfilled(log types.Log) (*IVRFCoordinatorV2PlusInternalRandomWordsFulfilled, error) {
	event := new(IVRFCoordinatorV2PlusInternalRandomWordsFulfilled)
	if err := _IVRFCoordinatorV2PlusInternal.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator struct {
	Event *IVRFCoordinatorV2PlusInternalRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IVRFCoordinatorV2PlusInternalRandomWordsRequested)
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

	select {
	case log := <-it.logs:
		it.Event = new(IVRFCoordinatorV2PlusInternalRandomWordsRequested)
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

func (it *IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IVRFCoordinatorV2PlusInternalRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       *big.Int
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	ExtraArgs                   []byte
	Sender                      common.Address
	Raw                         types.Log
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IVRFCoordinatorV2PlusInternal.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator{contract: _IVRFCoordinatorV2PlusInternal.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *IVRFCoordinatorV2PlusInternalRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IVRFCoordinatorV2PlusInternal.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IVRFCoordinatorV2PlusInternalRandomWordsRequested)
				if err := _IVRFCoordinatorV2PlusInternal.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternalFilterer) ParseRandomWordsRequested(log types.Log) (*IVRFCoordinatorV2PlusInternalRandomWordsRequested, error) {
	event := new(IVRFCoordinatorV2PlusInternalRandomWordsRequested)
	if err := _IVRFCoordinatorV2PlusInternal.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance       *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	Owner         common.Address
	Consumers     []common.Address
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternal) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IVRFCoordinatorV2PlusInternal.abi.Events["RandomWordsFulfilled"].ID:
		return _IVRFCoordinatorV2PlusInternal.ParseRandomWordsFulfilled(log)
	case _IVRFCoordinatorV2PlusInternal.abi.Events["RandomWordsRequested"].ID:
		return _IVRFCoordinatorV2PlusInternal.ParseRandomWordsRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IVRFCoordinatorV2PlusInternalRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x6c6b5394380e16e41988d8383648010de6f5c2e4814803be5de1c6b1c852db55")
}

func (IVRFCoordinatorV2PlusInternalRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (_IVRFCoordinatorV2PlusInternal *IVRFCoordinatorV2PlusInternal) Address() common.Address {
	return _IVRFCoordinatorV2PlusInternal.address
}

type IVRFCoordinatorV2PlusInternalInterface interface {
	LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error)

	GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error)

	SRequestCommitments(opts *bind.CallOpts, requestID *big.Int) ([32]byte, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof IVRFCoordinatorV2PlusInternalProof, rc IVRFCoordinatorV2PlusInternalRequestCommitment, onlyPremium bool) (*types.Transaction, error)

	FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*IVRFCoordinatorV2PlusInternalRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *IVRFCoordinatorV2PlusInternalRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*IVRFCoordinatorV2PlusInternalRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*IVRFCoordinatorV2PlusInternalRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *IVRFCoordinatorV2PlusInternalRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*IVRFCoordinatorV2PlusInternalRandomWordsRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
