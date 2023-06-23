// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_sub_owner

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

var VRFV2PlusExternalSubOwnerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161074338038061074383398101604081905261002f9161009e565b6001600160601b0319606083901b16608052600080546001600160a01b03199081166001600160a01b039485161790915560018054929093169181169190911790915560048054909116331790556100d1565b80516001600160a01b038116811461009957600080fd5b919050565b600080604083850312156100b157600080fd5b6100ba83610082565b91506100c860208401610082565b90509250929050565b60805160601c61064e6100f56000396000818160ed0152610155015261064e6000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063f2fde38b11610050578063f2fde38b1461009c578063f6eaffc8146100af578063f793a70e146100c257600080fd5b80631fe543e31461006c578063e89e106a14610081575b600080fd5b61007f61007a366004610495565b6100d5565b005b61008a60035481565b60405190815260200160405180910390f35b61007f6100aa366004610426565b610195565b61008a6100bd366004610463565b610200565b61007f6100d0366004610584565b610221565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610187576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b610191828261032a565b5050565b60045473ffffffffffffffffffffffffffffffffffffffff1633146101b957600080fd5b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6002818154811061021057600080fd5b600091825260209091200154905081565b60045473ffffffffffffffffffffffffffffffffffffffff16331461024557600080fd5b6000546040517fefcf1d940000000000000000000000000000000000000000000000000000000081526004810184905267ffffffffffffffff8816602482015261ffff8616604482015263ffffffff80881660648301528516608482015282151560a482015273ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b1580156102e757600080fd5b505af11580156102fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031f919061047c565b600355505050505050565b6003548214610395576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161017e565b80516103a89060029060208401906103ad565b505050565b8280548282559060005260206000209081019282156103e8579160200282015b828111156103e85782518255916020019190600101906103cd565b506103f49291506103f8565b5090565b5b808211156103f457600081556001016103f9565b803563ffffffff8116811461042157600080fd5b919050565b60006020828403121561043857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461045c57600080fd5b9392505050565b60006020828403121561047557600080fd5b5035919050565b60006020828403121561048e57600080fd5b5051919050565b600080604083850312156104a857600080fd5b8235915060208084013567ffffffffffffffff808211156104c857600080fd5b818601915086601f8301126104dc57600080fd5b8135818111156104ee576104ee610612565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561053157610531610612565b604052828152858101935084860182860187018b101561055057600080fd5b600095505b83861015610573578035855260019590950194938601938601610555565b508096505050505050509250929050565b60008060008060008060c0878903121561059d57600080fd5b863567ffffffffffffffff811681146105b557600080fd5b95506105c36020880161040d565b9450604087013561ffff811681146105da57600080fd5b93506105e86060880161040d565b92506080870135915060a0870135801515811461060457600080fd5b809150509295509295509295565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusExternalSubOwnerExampleABI = VRFV2PlusExternalSubOwnerExampleMetaData.ABI

var VRFV2PlusExternalSubOwnerExampleBin = VRFV2PlusExternalSubOwnerExampleMetaData.Bin

func DeployVRFV2PlusExternalSubOwnerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusExternalSubOwnerExample, error) {
	parsed, err := VRFV2PlusExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusExternalSubOwnerExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusExternalSubOwnerExample{VRFV2PlusExternalSubOwnerExampleCaller: VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, VRFV2PlusExternalSubOwnerExampleTransactor: VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, VRFV2PlusExternalSubOwnerExampleFilterer: VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusExternalSubOwnerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusExternalSubOwnerExampleCaller
	VRFV2PlusExternalSubOwnerExampleTransactor
	VRFV2PlusExternalSubOwnerExampleFilterer
}

type VRFV2PlusExternalSubOwnerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusExternalSubOwnerExampleSession struct {
	Contract     *VRFV2PlusExternalSubOwnerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusExternalSubOwnerExampleCallerSession struct {
	Contract *VRFV2PlusExternalSubOwnerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusExternalSubOwnerExampleTransactorSession struct {
	Contract     *VRFV2PlusExternalSubOwnerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusExternalSubOwnerExampleRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExample
}

type VRFV2PlusExternalSubOwnerExampleCallerRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExampleCaller
}

type VRFV2PlusExternalSubOwnerExampleTransactorRaw struct {
	Contract *VRFV2PlusExternalSubOwnerExampleTransactor
}

func NewVRFV2PlusExternalSubOwnerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusExternalSubOwnerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExample{address: address, abi: abi, VRFV2PlusExternalSubOwnerExampleCaller: VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, VRFV2PlusExternalSubOwnerExampleTransactor: VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, VRFV2PlusExternalSubOwnerExampleFilterer: VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusExternalSubOwnerExampleCaller, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusExternalSubOwnerExampleTransactor, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusExternalSubOwnerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusExternalSubOwnerExampleFilterer, error) {
	contract, err := bindVRFV2PlusExternalSubOwnerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusExternalSubOwnerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusExternalSubOwnerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.VRFV2PlusExternalSubOwnerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusExternalSubOwnerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRandomWords(&_VRFV2PlusExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRandomWords(&_VRFV2PlusExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusExternalSubOwnerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRequestId(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.SRequestId(&_VRFV2PlusExternalSubOwnerExample.CallOpts)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.TransferOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, newOwner)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExampleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _VRFV2PlusExternalSubOwnerExample.Contract.TransferOwnership(&_VRFV2PlusExternalSubOwnerExample.TransactOpts, newOwner)
}

func (_VRFV2PlusExternalSubOwnerExample *VRFV2PlusExternalSubOwnerExample) Address() common.Address {
	return _VRFV2PlusExternalSubOwnerExample.address
}

type VRFV2PlusExternalSubOwnerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Address() common.Address
}
