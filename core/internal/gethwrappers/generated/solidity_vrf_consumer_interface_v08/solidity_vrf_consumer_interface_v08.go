// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity_vrf_consumer_interface_v08

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

var VRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516104ef3803806104ef83398101604081905261002f91610069565b6001600160601b0319606092831b811660a052911b1660805261009c565b80516001600160a01b038116811461006457600080fd5b919050565b6000806040838503121561007c57600080fd5b6100858361004d565b91506100936020840161004d565b90509250929050565b60805160601c60a05160601c6104226100cd6000396000818160c801526101a40152600061016801526104226000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80635eb7978314610051578063866ee7481461006c57806394985ddd1461007f578063e89e106a14610094575b600080fd5b61005a60015481565b60405190815260200160405180910390f35b61005a61007a36600461031c565b61009d565b61009261008d36600461031c565b6100b0565b005b61005a60025481565b60006100a98383610164565b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610153576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015260640160405180910390fd5b6101608282600155600255565b5050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f0000000000000000000000000000000000000000000000000000000000000000848660006040516020016101e1929190918252602082015260400190565b6040516020818303038152906040526040518463ffffffff1660e01b815260040161020e9392919061033e565b602060405180830381600087803b15801561022857600080fd5b505af115801561023c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061026091906102fa565b5060008381526020818152604080832054815180840188905280830185905230606082015260808082018390528351808303909101815260a0909101909252815191830191909120868452929091526102ba9060016103d6565b6000858152602081815260409182902092909255805180830187905280820184905281518082038301815260609091019091528051910120949350505050565b60006020828403121561030c57600080fd5b815180151581146100a957600080fd5b6000806040838503121561032f57600080fd5b50508035926020909101359150565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561038e57858101830151858201608001528201610372565b818111156103a0576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b60008219821115610410577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c6343000806000a",
}

var VRFConsumerABI = VRFConsumerMetaData.ABI

var VRFConsumerBin = VRFConsumerMetaData.Bin

func DeployVRFConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFConsumer, error) {
	parsed, err := VRFConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumer{VRFConsumerCaller: VRFConsumerCaller{contract: contract}, VRFConsumerTransactor: VRFConsumerTransactor{contract: contract}, VRFConsumerFilterer: VRFConsumerFilterer{contract: contract}}, nil
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
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_VRFConsumer *VRFConsumerCaller) SRandomnessOutput(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "s_randomnessOutput")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumer *VRFConsumerSession) SRandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.SRandomnessOutput(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCallerSession) SRandomnessOutput() (*big.Int, error) {
	return _VRFConsumer.Contract.SRandomnessOutput(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCaller) SRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFConsumer.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFConsumer *VRFConsumerSession) SRequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.SRequestId(&_VRFConsumer.CallOpts)
}

func (_VRFConsumer *VRFConsumerCallerSession) SRequestId() ([32]byte, error) {
	return _VRFConsumer.Contract.SRequestId(&_VRFConsumer.CallOpts)
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

func (_VRFConsumer *VRFConsumerTransactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "testRequestRandomness", keyHash, fee)
}

func (_VRFConsumer *VRFConsumerSession) TestRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
}

func (_VRFConsumer *VRFConsumerTransactorSession) TestRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.TestRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
}

func (_VRFConsumer *VRFConsumer) Address() common.Address {
	return _VRFConsumer.address
}

type VRFConsumerInterface interface {
	SRandomnessOutput(opts *bind.CallOpts) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) ([32]byte, error)

	RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, fee *big.Int) (*types.Transaction, error)

	Address() common.Address
}
