// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_reverting_example

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

var VRFV2PlusRevertingExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlySubOwnerCanSetVRFCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610eb5380380610eb583398101604081905261002f916100d5565b816001600160a01b0381166100795760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640160405180910390fd5b600080546001600160a01b03199081166001600160a01b039384161790915560048054821694831694909417909355600580549093169116179055610108565b80516001600160a01b03811681146100d057600080fd5b919050565b600080604083850312156100e857600080fd5b6100f1836100b9565b91506100ff602084016100b9565b90509250929050565b610d9e806101176000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063706da1ca11610076578063e89e106a1161005b578063e89e106a1461018f578063f08c5daa14610198578063f6eaffc8146101a157600080fd5b8063706da1ca14610137578063cf62c8ab1461017c57600080fd5b80632fa4e442116100a75780632fa4e442146100fe57806336bfffed1461011157806344ff81ce1461012457600080fd5b8063177b9692146100c35780631fe543e3146100e9575b600080fd5b6100d66100d1366004610a17565b6101b4565b6040519081526020015b60405180910390f35b6100fc6100f7366004610ab2565b61029f565b005b6100fc61010c366004610b73565b610325565b6100fc61011f366004610951565b610485565b6100fc61013236600461092f565b61060d565b6005546101639074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100e0565b6100fc61018a366004610b73565b6106c7565b6100d660035481565b6100d660065481565b6100d66101af366004610a80565b6108d1565b600480546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815291820187905267ffffffffffffffff8616602483015261ffff8516604483015263ffffffff808516606484015283166084830152600060a483018190529173ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b15801561025857600080fd5b505af115801561026c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102909190610a99565b60038190559695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610317576000546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6103218282600080fd5b5050565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff166103b0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161030e565b6005546004546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161043393929190610ba1565b602060405180830381600087803b15801561044d57600080fd5b505af1158015610461573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032191906109f5565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff16610510576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161030e565b60005b815181101561032157600454600554835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061057857610578610d1a565b60200260200101516040518363ffffffff1660e01b81526004016105c892919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105e257600080fd5b505af11580156105f6573d6000803e3d6000fd5b50505050808061060590610cba565b915050610513565b60015473ffffffffffffffffffffffffffffffffffffffff163314610680576001546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff909116602482015260440161030e565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff166103b05760048054604080517fa21a23e4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9092169263a21a23e49282820192602092908290030181600087803b15801561075a57600080fd5b505af115801561076e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107929190610b56565b600580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff93841681029190911791829055600480546040517f7341c10c000000000000000000000000000000000000000000000000000000008152929093049093169281019290925230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b15801561085c57600080fd5b505af1158015610870573d6000803e3d6000fd5b50506005546004546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610406565b600281815481106108e157600080fd5b600091825260209091200154905081565b803573ffffffffffffffffffffffffffffffffffffffff8116811461091657600080fd5b919050565b803563ffffffff8116811461091657600080fd5b60006020828403121561094157600080fd5b61094a826108f2565b9392505050565b6000602080838503121561096457600080fd5b823567ffffffffffffffff81111561097b57600080fd5b8301601f8101851361098c57600080fd5b803561099f61099a82610c96565b610c47565b80828252848201915084840188868560051b87010111156109bf57600080fd5b600094505b838510156109e9576109d5816108f2565b8352600194909401939185019185016109c4565b50979650505050505050565b600060208284031215610a0757600080fd5b8151801515811461094a57600080fd5b600080600080600060a08688031215610a2f57600080fd5b853594506020860135610a4181610d78565b9350604086013561ffff81168114610a5857600080fd5b9250610a666060870161091b565b9150610a746080870161091b565b90509295509295909350565b600060208284031215610a9257600080fd5b5035919050565b600060208284031215610aab57600080fd5b5051919050565b60008060408385031215610ac557600080fd5b8235915060208084013567ffffffffffffffff811115610ae457600080fd5b8401601f81018613610af557600080fd5b8035610b0361099a82610c96565b80828252848201915084840189868560051b8701011115610b2357600080fd5b600094505b83851015610b46578035835260019490940193918501918501610b28565b5080955050505050509250929050565b600060208284031215610b6857600080fd5b815161094a81610d78565b600060208284031215610b8557600080fd5b81356bffffffffffffffffffffffff8116811461094a57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610bff57858101830151858201608001528201610be3565b81811115610c11576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c8e57610c8e610d49565b604052919050565b600067ffffffffffffffff821115610cb057610cb0610d49565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610d13577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610d8e57600080fd5b5056fea164736f6c6343000806000a",
}

var VRFV2PlusRevertingExampleABI = VRFV2PlusRevertingExampleMetaData.ABI

var VRFV2PlusRevertingExampleBin = VRFV2PlusRevertingExampleMetaData.Bin

func DeployVRFV2PlusRevertingExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusRevertingExample, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusRevertingExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusRevertingExample{VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusRevertingExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusRevertingExampleCaller
	VRFV2PlusRevertingExampleTransactor
	VRFV2PlusRevertingExampleFilterer
}

type VRFV2PlusRevertingExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleSession struct {
	Contract     *VRFV2PlusRevertingExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleCallerSession struct {
	Contract *VRFV2PlusRevertingExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusRevertingExampleTransactorSession struct {
	Contract     *VRFV2PlusRevertingExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleRaw struct {
	Contract *VRFV2PlusRevertingExample
}

type VRFV2PlusRevertingExampleCallerRaw struct {
	Contract *VRFV2PlusRevertingExampleCaller
}

type VRFV2PlusRevertingExampleTransactorRaw struct {
	Contract *VRFV2PlusRevertingExampleTransactor
}

func NewVRFV2PlusRevertingExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusRevertingExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusRevertingExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusRevertingExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExample{address: address, abi: abi, VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusRevertingExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusRevertingExampleCaller, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusRevertingExampleTransactor, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusRevertingExampleFilterer, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusRevertingExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "setVRFCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetVRFCoordinator(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetVRFCoordinator(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) Address() common.Address {
	return _VRFV2PlusRevertingExample.address
}

type VRFV2PlusRevertingExampleInterface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
