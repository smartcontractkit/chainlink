// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_external_sub_owner_example

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
)

var VRFExternalSubOwnerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"setSubscriptionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516107f83803806107f883398101604081905261002f91610146565b606086811b6001600160601b0319166080908152600080546001600160a01b03998a166001600160a01b031991821617825560018054821699909a16989098179098556040805160a08101825298895263ffffffff96871660208a0181905261ffff969096169089018190529390951690870181905295909301839052600280546001600160701b0319166801000000000000000090930261ffff60601b1916929092176c010000000000000000000000009091021763ffffffff60701b1916600160701b90940293909317909255600391909155600680543392169190911790556101bc565b80516001600160a01b038116811461012d57600080fd5b919050565b805163ffffffff8116811461012d57600080fd5b60008060008060008060c0878903121561015f57600080fd5b61016887610116565b955061017660208801610116565b945061018460408801610132565b9350606087015161ffff8116811461019b57600080fd5b92506101a960808801610132565b915060a087015190509295509295509295565b60805160601c6106176101e1600039600081816101e8015261025001526106176000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063e0c8628911610050578063e0c862891461013f578063e89e106a14610147578063f6eaffc81461015e57600080fd5b80631edebf0b146100775780631fe543e31461008c5780637db9263f1461009f575b600080fd5b61008a6100853660046105aa565b610171565b005b61008a61009a3660046104bb565b6101d0565b6002546003546100f89167ffffffffffffffff81169163ffffffff68010000000000000000830481169261ffff6c01000000000000000000000000820416926e0100000000000000000000000000009091049091169085565b6040805167ffffffffffffffff909616865263ffffffff948516602087015261ffff90931692850192909252919091166060830152608082015260a0015b60405180910390f35b61008a61028f565b61015060055481565b604051908152602001610136565b61015061016c366004610489565b6103f0565b60065473ffffffffffffffffffffffffffffffffffffffff16331461019557600080fd5b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610281576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b61028b8282610411565b5050565b60065473ffffffffffffffffffffffffffffffffffffffff1633146102b357600080fd5b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e010000000000000000000000000000909404909116606085018190526003546080860181905260005496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101919091526024810193909352604483019390935260648201526084810191909152909173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156103b257600080fd5b505af11580156103c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ea91906104a2565b60055550565b6004818154811061040057600080fd5b600091825260209091200154905081565b8051610424906004906020840190610429565b505050565b828054828255906000526020600020908101928215610464579160200282015b82811115610464578251825591602001919060010190610449565b50610470929150610474565b5090565b5b808211156104705760008155600101610475565b60006020828403121561049b57600080fd5b5035919050565b6000602082840312156104b457600080fd5b5051919050565b600080604083850312156104ce57600080fd5b8235915060208084013567ffffffffffffffff808211156104ee57600080fd5b818601915086601f83011261050257600080fd5b813581811115610514576105146105db565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610557576105576105db565b604052828152858101935084860182860187018b101561057657600080fd5b600095505b8386101561059957803585526001959095019493860193860161057b565b508096505050505050509250929050565b6000602082840312156105bc57600080fd5b813567ffffffffffffffff811681146105d457600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFExternalSubOwnerExampleABI = VRFExternalSubOwnerExampleMetaData.ABI

var VRFExternalSubOwnerExampleBin = VRFExternalSubOwnerExampleMetaData.Bin

func DeployVRFExternalSubOwnerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (common.Address, *types.Transaction, *VRFExternalSubOwnerExample, error) {
	parsed, err := VRFExternalSubOwnerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFExternalSubOwnerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFExternalSubOwnerExample{VRFExternalSubOwnerExampleCaller: VRFExternalSubOwnerExampleCaller{contract: contract}, VRFExternalSubOwnerExampleTransactor: VRFExternalSubOwnerExampleTransactor{contract: contract}, VRFExternalSubOwnerExampleFilterer: VRFExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

type VRFExternalSubOwnerExample struct {
	address common.Address
	abi     abi.ABI
	VRFExternalSubOwnerExampleCaller
	VRFExternalSubOwnerExampleTransactor
	VRFExternalSubOwnerExampleFilterer
}

type VRFExternalSubOwnerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFExternalSubOwnerExampleSession struct {
	Contract     *VRFExternalSubOwnerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFExternalSubOwnerExampleCallerSession struct {
	Contract *VRFExternalSubOwnerExampleCaller
	CallOpts bind.CallOpts
}

type VRFExternalSubOwnerExampleTransactorSession struct {
	Contract     *VRFExternalSubOwnerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFExternalSubOwnerExampleRaw struct {
	Contract *VRFExternalSubOwnerExample
}

type VRFExternalSubOwnerExampleCallerRaw struct {
	Contract *VRFExternalSubOwnerExampleCaller
}

type VRFExternalSubOwnerExampleTransactorRaw struct {
	Contract *VRFExternalSubOwnerExampleTransactor
}

func NewVRFExternalSubOwnerExample(address common.Address, backend bind.ContractBackend) (*VRFExternalSubOwnerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFExternalSubOwnerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExample{address: address, abi: abi, VRFExternalSubOwnerExampleCaller: VRFExternalSubOwnerExampleCaller{contract: contract}, VRFExternalSubOwnerExampleTransactor: VRFExternalSubOwnerExampleTransactor{contract: contract}, VRFExternalSubOwnerExampleFilterer: VRFExternalSubOwnerExampleFilterer{contract: contract}}, nil
}

func NewVRFExternalSubOwnerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFExternalSubOwnerExampleCaller, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleCaller{contract: contract}, nil
}

func NewVRFExternalSubOwnerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFExternalSubOwnerExampleTransactor, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleTransactor{contract: contract}, nil
}

func NewVRFExternalSubOwnerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFExternalSubOwnerExampleFilterer, error) {
	contract, err := bindVRFExternalSubOwnerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFExternalSubOwnerExampleFilterer{contract: contract}, nil
}

func bindVRFExternalSubOwnerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFExternalSubOwnerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleTransactor.contract.Transfer(opts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.VRFExternalSubOwnerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFExternalSubOwnerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.contract.Transfer(opts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFExternalSubOwnerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRandomWords(&_VRFExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRandomWords(&_VRFExternalSubOwnerExample.CallOpts, arg0)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCaller) SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

	error) {
	var out []interface{}
	err := _VRFExternalSubOwnerExample.contract.Call(opts, &out, "s_requestConfig")

	outstruct := new(SRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.NumWords = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.KeyHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestConfig(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestConfig(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFExternalSubOwnerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestId(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFExternalSubOwnerExample.Contract.SRequestId(&_VRFExternalSubOwnerExample.CallOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RawFulfillRandomWords(&_VRFExternalSubOwnerExample.TransactOpts, requestId, randomWords)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFExternalSubOwnerExample.TransactOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.RequestRandomWords(&_VRFExternalSubOwnerExample.TransactOpts)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactor) SetSubscriptionID(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.contract.Transact(opts, "setSubscriptionID", subId)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleSession) SetSubscriptionID(subId uint64) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.SetSubscriptionID(&_VRFExternalSubOwnerExample.TransactOpts, subId)
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExampleTransactorSession) SetSubscriptionID(subId uint64) (*types.Transaction, error) {
	return _VRFExternalSubOwnerExample.Contract.SetSubscriptionID(&_VRFExternalSubOwnerExample.TransactOpts, subId)
}

type SRequestConfig struct {
	SubId                uint64
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
}

func (_VRFExternalSubOwnerExample *VRFExternalSubOwnerExample) Address() common.Address {
	return _VRFExternalSubOwnerExample.address
}

type VRFExternalSubOwnerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	SetSubscriptionID(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	Address() common.Address
}
