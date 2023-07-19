// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_malicious_consumer_v2_plus

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

var VRFMaliciousConsumerV2PlusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlySubOwnerCanSetVRFCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610e95380380610e9583398101604081905261002f916100d5565b816001600160a01b0381166100795760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640160405180910390fd5b600080546001600160a01b03199081166001600160a01b039384161790915560048054821694831694909417909355600580549093169116179055610108565b80516001600160a01b03811681146100d057600080fd5b919050565b600080604083850312156100e857600080fd5b6100f1836100b9565b91506100ff602084016100b9565b90509250929050565b610d7e806101176000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063706da1ca11610076578063e89e106a1161005b578063e89e106a14610161578063f08c5daa1461016a578063f6eaffc81461017357600080fd5b8063706da1ca14610109578063cf62c8ab1461014e57600080fd5b80631fe543e3146100a857806336bfffed146100bd57806344ff81ce146100d05780635e3b709f146100e3575b600080fd5b6100bb6100b6366004610a9e565b610186565b005b6100bb6100cb3660046109a6565b61020c565b6100bb6100de366004610984565b610394565b6100f66100f1366004610a6c565b61044e565b6040519081526020015b60405180910390f35b6005546101359074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610100565b6100bb61015c366004610b6c565b610548565b6100f660035481565b6100f660065481565b6100f6610181366004610a6c565b6107ca565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101fe576000546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61020882826107eb565b5050565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff16610297576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f74207365740000000000000000000000000000000000000060448201526064016101f5565b60005b815181101561020857600454600554835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff16908590859081106102ff576102ff610d13565b60200260200101516040518363ffffffff1660e01b815260040161034f92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561036957600080fd5b505af115801561037d573d6000803e3d6000fd5b50505050808061038c90610cb3565b91505061029a565b60015473ffffffffffffffffffffffffffffffffffffffff163314610407576001546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044016101f5565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6007819055600480546005546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815292830184905274010000000000000000000000000000000000000000900467ffffffffffffffff1660248301526001604483018190526207a12060648401526084830152600060a483018190529173ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b15801561050a57600080fd5b505af115801561051e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105429190610a85565b92915050565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff166106f65760048054604080517fa21a23e4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9092169263a21a23e49282820192602092908290030181600087803b1580156105db57600080fd5b505af11580156105ef573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106139190610b42565b600580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff93841681029190911791829055600480546040517f7341c10c000000000000000000000000000000000000000000000000000000008152929093049093169281019290925230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b1580156106dd57600080fd5b505af11580156106f1573d6000803e3d6000fd5b505050505b6005546004546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161077893929190610b9a565b602060405180830381600087803b15801561079257600080fd5b505af11580156107a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102089190610a4a565b600281815481106107da57600080fd5b600091825260209091200154905081565b5a60065580516108029060029060208401906108fb565b506003829055600480546007546005546040517fefcf1d940000000000000000000000000000000000000000000000000000000081529384019190915274010000000000000000000000000000000000000000900467ffffffffffffffff16602483015260016044830181905262030d4060648401526084830152600060a483015273ffffffffffffffffffffffffffffffffffffffff169063efcf1d949060c401602060405180830381600087803b1580156108be57600080fd5b505af11580156108d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108f69190610a85565b505050565b828054828255906000526020600020908101928215610936579160200282015b8281111561093657825182559160200191906001019061091b565b50610942929150610946565b5090565b5b808211156109425760008155600101610947565b803573ffffffffffffffffffffffffffffffffffffffff8116811461097f57600080fd5b919050565b60006020828403121561099657600080fd5b61099f8261095b565b9392505050565b600060208083850312156109b957600080fd5b823567ffffffffffffffff8111156109d057600080fd5b8301601f810185136109e157600080fd5b80356109f46109ef82610c8f565b610c40565b80828252848201915084840188868560051b8701011115610a1457600080fd5b600094505b83851015610a3e57610a2a8161095b565b835260019490940193918501918501610a19565b50979650505050505050565b600060208284031215610a5c57600080fd5b8151801515811461099f57600080fd5b600060208284031215610a7e57600080fd5b5035919050565b600060208284031215610a9757600080fd5b5051919050565b60008060408385031215610ab157600080fd5b8235915060208084013567ffffffffffffffff811115610ad057600080fd5b8401601f81018613610ae157600080fd5b8035610aef6109ef82610c8f565b80828252848201915084840189868560051b8701011115610b0f57600080fd5b600094505b83851015610b32578035835260019490940193918501918501610b14565b5080955050505050509250929050565b600060208284031215610b5457600080fd5b815167ffffffffffffffff8116811461099f57600080fd5b600060208284031215610b7e57600080fd5b81356bffffffffffffffffffffffff8116811461099f57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610bf857858101830151858201608001528201610bdc565b81811115610c0a576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c8757610c87610d42565b604052919050565b600067ffffffffffffffff821115610ca957610ca9610d42565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610d0c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFMaliciousConsumerV2PlusABI = VRFMaliciousConsumerV2PlusMetaData.ABI

var VRFMaliciousConsumerV2PlusBin = VRFMaliciousConsumerV2PlusMetaData.Bin

func DeployVRFMaliciousConsumerV2Plus(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFMaliciousConsumerV2Plus, error) {
	parsed, err := VRFMaliciousConsumerV2PlusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFMaliciousConsumerV2PlusBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFMaliciousConsumerV2Plus{VRFMaliciousConsumerV2PlusCaller: VRFMaliciousConsumerV2PlusCaller{contract: contract}, VRFMaliciousConsumerV2PlusTransactor: VRFMaliciousConsumerV2PlusTransactor{contract: contract}, VRFMaliciousConsumerV2PlusFilterer: VRFMaliciousConsumerV2PlusFilterer{contract: contract}}, nil
}

type VRFMaliciousConsumerV2Plus struct {
	address common.Address
	abi     abi.ABI
	VRFMaliciousConsumerV2PlusCaller
	VRFMaliciousConsumerV2PlusTransactor
	VRFMaliciousConsumerV2PlusFilterer
}

type VRFMaliciousConsumerV2PlusCaller struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusTransactor struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusFilterer struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2PlusSession struct {
	Contract     *VRFMaliciousConsumerV2Plus
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2PlusCallerSession struct {
	Contract *VRFMaliciousConsumerV2PlusCaller
	CallOpts bind.CallOpts
}

type VRFMaliciousConsumerV2PlusTransactorSession struct {
	Contract     *VRFMaliciousConsumerV2PlusTransactor
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2PlusRaw struct {
	Contract *VRFMaliciousConsumerV2Plus
}

type VRFMaliciousConsumerV2PlusCallerRaw struct {
	Contract *VRFMaliciousConsumerV2PlusCaller
}

type VRFMaliciousConsumerV2PlusTransactorRaw struct {
	Contract *VRFMaliciousConsumerV2PlusTransactor
}

func NewVRFMaliciousConsumerV2Plus(address common.Address, backend bind.ContractBackend) (*VRFMaliciousConsumerV2Plus, error) {
	abi, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2PlusABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFMaliciousConsumerV2Plus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Plus{address: address, abi: abi, VRFMaliciousConsumerV2PlusCaller: VRFMaliciousConsumerV2PlusCaller{contract: contract}, VRFMaliciousConsumerV2PlusTransactor: VRFMaliciousConsumerV2PlusTransactor{contract: contract}, VRFMaliciousConsumerV2PlusFilterer: VRFMaliciousConsumerV2PlusFilterer{contract: contract}}, nil
}

func NewVRFMaliciousConsumerV2PlusCaller(address common.Address, caller bind.ContractCaller) (*VRFMaliciousConsumerV2PlusCaller, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusCaller{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2PlusTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFMaliciousConsumerV2PlusTransactor, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusTransactor{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2PlusFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFMaliciousConsumerV2PlusFilterer, error) {
	contract, err := bindVRFMaliciousConsumerV2Plus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2PlusFilterer{contract: contract}, nil
}

func bindVRFMaliciousConsumerV2Plus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFMaliciousConsumerV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusCaller.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusTransactor.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.VRFMaliciousConsumerV2PlusTransactor.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SGasAvailable(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SGasAvailable(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRandomWords(&_VRFMaliciousConsumerV2Plus.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRandomWords(&_VRFMaliciousConsumerV2Plus.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRequestId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SRequestId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2Plus.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SSubId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusCallerSession) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SSubId(&_VRFMaliciousConsumerV2Plus.CallOpts)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2Plus.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.CreateSubscriptionAndFund(&_VRFMaliciousConsumerV2Plus.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2Plus.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2Plus.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "requestRandomness", keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) RequestRandomness(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RequestRandomness(&_VRFMaliciousConsumerV2Plus.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) RequestRandomness(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.RequestRandomness(&_VRFMaliciousConsumerV2Plus.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "setVRFCoordinator", _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetVRFCoordinator(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.SetVRFCoordinator(&_VRFMaliciousConsumerV2Plus.TransactOpts, _vrfCoordinator)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2Plus.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2PlusTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2Plus.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2Plus.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2Plus *VRFMaliciousConsumerV2Plus) Address() common.Address {
	return _VRFMaliciousConsumerV2Plus.address
}

type VRFMaliciousConsumerV2PlusInterface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
