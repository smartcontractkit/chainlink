// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_consumer_interface

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

var VRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"currentRoundID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c0604052600060015534801561001557600080fd5b506040516105363803806105368339818101604052604081101561003857600080fd5b5080516020909101516001600160601b0319606092831b811660a052911b1660805260805160601c60a05160601c6104ae6100886000398061011452806101f45250806101b852506104ae6000f3fe608060405234801561001057600080fd5b50600436106100665760003560e01c8063866ee74811610050578063866ee7481461008d57806394985ddd146100b0578063a312c4f2146100d557610066565b80626d6cae1461006b5780632f47fd8614610085575b600080fd5b6100736100dd565b60408051918252519081900360200190f35b6100736100e3565b610073600480360360408110156100a357600080fd5b50803590602001356100e9565b6100d3600480360360408110156100c657600080fd5b50803590602001356100fc565b005b6100736101ae565b60035481565b60025481565b60006100f583836101b4565b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146101a057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015290519081900360640190fd5b6101aa828261039d565b5050565b60015481565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f00000000000000000000000000000000000000000000000000000000000000008486600060405160200180838152602001828152602001925050506040516020818303038152906040526040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102c15781810151838201526020016102a9565b50505050905090810190601f1680156102ee5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561030f57600080fd5b505af1158015610323573d6000803e3d6000fd5b505050506040513d602081101561033957600080fd5b5050600083815260208190526040812054610359908590839030906103ad565b60008581526020819052604090205490915061037c90600163ffffffff61040116565b6000858152602081905260409020556103958482610475565b949350505050565b6002556003556001805481019055565b604080516020808201969096528082019490945273ffffffffffffffffffffffffffffffffffffffff9290921660608401526080808401919091528151808403909101815260a09092019052805191012090565b6000828201838110156100f557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b60408051602080820194909452808201929092528051808303820181526060909201905280519101209056fea164736f6c6343000606000a",
}

var VRFConsumerABI = VRFConsumerMetaData.ABI

var VRFConsumerBin = VRFConsumerMetaData.Bin

func DeployVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFConsumer, error) {
	parsed, err := VRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumer{address: address, abi: *parsed, VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

type VRFConsumer struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerCaller
	VRFConsumerTransactor
	VRFConsumerFilterer
}

type VRFConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFConsumerSession struct {
	Contract     *VRFConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerCallerSession struct {
	Contract *VRFConsumerCaller
	CallOpts bind.CallOpts
}

type VRFConsumerTransactorSession struct {
	Contract     *VRFConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerRaw struct {
	Contract *VRFConsumer
}

type VRFConsumerCallerRaw struct {
	Contract *VRFConsumerCaller
}

type VRFConsumerTransactorRaw struct {
	Contract *VRFConsumerTransactor
}

func NewVRFConsumer(address common.Address, backend bind.ContractBackend) (*VRFConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumer{address: address, abi: abi, VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
}

func NewVRFConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerCaller, error) {
	contract, err := bindVRFConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerCaller{contract: contract}, nil
}

func NewVRFConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerTransactor, error) {
	contract, err := bindVRFConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerTransactor{contract: contract}, nil
}

func NewVRFConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerFilterer, error) {
	contract, err := bindVRFConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerFilterer{contract: contract}, nil
}

func bindVRFConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFConsumer *VRFConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.VRFConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumer *VRFConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transfer(opts)
}

func (_VRFConsumer *VRFConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.VRFConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumer *VRFConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumer *VRFConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transfer(opts)
}

func (_VRFConsumer *VRFConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumer *VRFConsumerCaller) CurrentRoundID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "currentRoundID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumer *VRFConsumerSession) CurrentRoundID() (*big.Int, error) {
	return _VRFConsumer.Contract.CurrentRoundID(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCallerSession) CurrentRoundID() (*big.Int, error) {
	return _VRFConsumer.Contract.CurrentRoundID(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCaller) RandomnessOutput(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "randomnessOutput")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumer *VRFConsumerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCallerSession) RandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.RandomnessOutput(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCaller) RequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "requestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFConsumer *VRFConsumerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCallerSession) RequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.RequestId(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

func (_VRFConsumer *VRFConsumerSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

func (_VRFConsumer *VRFConsumerTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.RawFulfillRandomness(&_VRFConsumer.TransactOpts, requestId, randomness)
}

func (_VRFConsumer *VRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "testRequestRandomness", _keyHash, _fee)
}

func (_VRFConsumer *VRFConsumerSession) TestRequestRandomness(_keyHash [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee)
}

func (_VRFConsumer *VRFConsumerTransactorSession) TestRequestRandomness(_keyHash [32]byte, _fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, _keyHash, _fee)
}

func (_VRFConsumer *VRFConsumer) Address() common.Address {
	return _VRFConsumer.address
}

type VRFConsumerInterface interface {
	CurrentRoundID(opts *bind.CallOpts) (*big.Int, error)

	RandomnessOutput(opts *bind.CallOpts) (*big.Int, error)

	RequestId(opts *bind.CallOpts) ([32]byte, error)

	RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _fee *big.Int) (*types.Transaction, error)

	Address() common.Address
}
