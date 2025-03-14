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
	_ = abi.ConvertType
)

var VRFConsumerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"vrfCoordinator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"doRequestRandomness\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"randomnessOutput\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rawFulfillRandomness\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"randomness\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516104e33803806104e383398101604081905261002f91610069565b6001600160601b0319606092831b811660a052911b1660805261009c565b80516001600160a01b038116811461006457600080fd5b919050565b6000806040838503121561007c57600080fd5b6100858361004d565b91506100936020840161004d565b90509250929050565b60805160601c60a05160601c6104166100cd6000396000818160c701526101980152600061015c01526104166000f3fe608060405234801561001057600080fd5b506004361061004b5760003560e01c80626d6cae146100505780631a8da9761461006b5780632f47fd861461007e57806394985ddd14610087575b600080fd5b61005960025481565b60405190815260200160405180910390f35b610059610079366004610310565b61009c565b61005960015481565b61009a610095366004610310565b6100af565b005b60006100a88383610158565b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610152576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c00604482015260640160405180910390fd5b60015550565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f0000000000000000000000000000000000000000000000000000000000000000848660006040516020016101d5929190918252602082015260400190565b6040516020818303038152906040526040518463ffffffff1660e01b815260040161020293929190610332565b602060405180830381600087803b15801561021c57600080fd5b505af1158015610230573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025491906102ee565b5060008381526020818152604080832054815180840188905280830185905230606082015260808082018390528351808303909101815260a0909101909252815191830191909120868452929091526102ae9060016103ca565b6000858152602081815260409182902092909255805180830187905280820184905281518082038301815260609091019091528051910120949350505050565b60006020828403121561030057600080fd5b815180151581146100a857600080fd5b6000806040838503121561032357600080fd5b50508035926020909101359150565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561038257858101830151858201608001528201610366565b81811115610394576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b60008219821115610404577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c6343000806000a",
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

func (_VRFConsumer *VRFConsumerTransactor) DoRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.contract.Transact(opts, "doRequestRandomness", keyHash, fee)
}

func (_VRFConsumer *VRFConsumerSession) DoRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.DoRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
}

func (_VRFConsumer *VRFConsumerTransactorSession) DoRequestRandomness(keyHash [32]byte, fee *big.Int) (*types.Transaction, error) {
	return _VRFConsumer.Contract.DoRequestRandomness(&_VRFConsumer.TransactOpts, keyHash, fee)
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

func (_VRFConsumer *VRFConsumer) Address() common.Address {
	return _VRFConsumer.address
}

type VRFConsumerInterface interface {
	RandomnessOutput(opts *bind.CallOpts) (*big.Int, error)

	RequestId(opts *bind.CallOpts) ([32]byte, error)

	DoRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, fee *big.Int) (*types.Transaction, error)

	RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error)

	Address() common.Address
}
