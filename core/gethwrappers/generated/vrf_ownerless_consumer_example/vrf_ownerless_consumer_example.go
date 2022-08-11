// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_ownerless_consumer_example

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

var VRFOwnerlessConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_randomnessOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516106a73803806106a783398101604081905261002f91610069565b6001600160601b0319606092831b811660a052911b1660805261009c565b80516001600160a01b038116811461006457600080fd5b919050565b6000806040838503121561007c57600080fd5b6100858361004d565b91506100936020840161004d565b90509250929050565b60805160601c60a05160601c6105d36100d46000396000818160b50152610293015260008181610167015261025701526105d36000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80635eb797831461005157806394985ddd1461006c578063a4c0ed3614610081578063e89e106a14610094575b600080fd5b61005a60015481565b60405190815260200160405180910390f35b61007f61007a3660046104cd565b61009d565b005b61007f61008f3660046103e9565b61014f565b61005a60025481565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610141576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f4f6e6c7920565246436f6f7264696e61746f722063616e2066756c66696c6c0060448201526064015b60405180910390fd5b61014b82826101e2565b5050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146101be576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006101cc828401846104b4565b90506101d88185610253565b6002555050505050565b600254821461024d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610138565b60015550565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f0000000000000000000000000000000000000000000000000000000000000000848660006040516020016102d0929190918252602082015260400190565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016102fd939291906104ef565b602060405180830381600087803b15801561031757600080fd5b505af115801561032b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061034f919061048b565b5060008381526020818152604080832054815180840188905280830185905230606082015260808082018390528351808303909101815260a0909101909252815191830191909120868452929091526103a9906001610587565b6000858152602081815260409182902092909255805180830187905280820184905281518082038301815260609091019091528051910120949350505050565b600080600080606085870312156103ff57600080fd5b843573ffffffffffffffffffffffffffffffffffffffff8116811461042357600080fd5b935060208501359250604085013567ffffffffffffffff8082111561044757600080fd5b818701915087601f83011261045b57600080fd5b81358181111561046a57600080fd5b88602082850101111561047c57600080fd5b95989497505060200194505050565b60006020828403121561049d57600080fd5b815180151581146104ad57600080fd5b9392505050565b6000602082840312156104c657600080fd5b5035919050565b600080604083850312156104e057600080fd5b50508035926020909101359150565b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b8181101561053f57858101830151858201608001528201610523565b81811115610551576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b600082198211156105c1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c6343000806000a",
}

var VRFOwnerlessConsumerExampleABI = VRFOwnerlessConsumerExampleMetaData.ABI

var VRFOwnerlessConsumerExampleBin = VRFOwnerlessConsumerExampleMetaData.Bin

func DeployVRFOwnerlessConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address, _link common.Address) (common.Address, *types.Transaction, *VRFOwnerlessConsumerExample, error) {
	parsed, err := VRFOwnerlessConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFOwnerlessConsumerExampleBin), backend, _vrfCoordinator, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFOwnerlessConsumerExample{VRFOwnerlessConsumerExampleCaller: VRFOwnerlessConsumerExampleCaller{contract: contract}, VRFOwnerlessConsumerExampleTransactor: VRFOwnerlessConsumerExampleTransactor{contract: contract}, VRFOwnerlessConsumerExampleFilterer: VRFOwnerlessConsumerExampleFilterer{contract: contract}}, nil
}

type VRFOwnerlessConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFOwnerlessConsumerExampleCaller
	VRFOwnerlessConsumerExampleTransactor
	VRFOwnerlessConsumerExampleFilterer
}

type VRFOwnerlessConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFOwnerlessConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFOwnerlessConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFOwnerlessConsumerExampleSession struct {
	Contract     *VRFOwnerlessConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFOwnerlessConsumerExampleCallerSession struct {
	Contract *VRFOwnerlessConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFOwnerlessConsumerExampleTransactorSession struct {
	Contract     *VRFOwnerlessConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFOwnerlessConsumerExampleRaw struct {
	Contract *VRFOwnerlessConsumerExample
}

type VRFOwnerlessConsumerExampleCallerRaw struct {
	Contract *VRFOwnerlessConsumerExampleCaller
}

type VRFOwnerlessConsumerExampleTransactorRaw struct {
	Contract *VRFOwnerlessConsumerExampleTransactor
}

func NewVRFOwnerlessConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFOwnerlessConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFOwnerlessConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFOwnerlessConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerlessConsumerExample{address: address, abi: abi, VRFOwnerlessConsumerExampleCaller: VRFOwnerlessConsumerExampleCaller{contract: contract}, VRFOwnerlessConsumerExampleTransactor: VRFOwnerlessConsumerExampleTransactor{contract: contract}, VRFOwnerlessConsumerExampleFilterer: VRFOwnerlessConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFOwnerlessConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFOwnerlessConsumerExampleCaller, error) {
	contract, err := bindVRFOwnerlessConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerlessConsumerExampleCaller{contract: contract}, nil
}

func NewVRFOwnerlessConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFOwnerlessConsumerExampleTransactor, error) {
	contract, err := bindVRFOwnerlessConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerlessConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFOwnerlessConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFOwnerlessConsumerExampleFilterer, error) {
	contract, err := bindVRFOwnerlessConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFOwnerlessConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFOwnerlessConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFOwnerlessConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFOwnerlessConsumerExample.Contract.VRFOwnerlessConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.VRFOwnerlessConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.VRFOwnerlessConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFOwnerlessConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleCaller) SRandomnessOutput(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFOwnerlessConsumerExample.contract.Call(opts, &out, "s_randomnessOutput")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleSession) SRandomnessOutput() (*big.Int, error) {
	return _VRFOwnerlessConsumerExample.Contract.SRandomnessOutput(&_VRFOwnerlessConsumerExample.CallOpts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleCallerSession) SRandomnessOutput() (*big.Int, error) {
	return _VRFOwnerlessConsumerExample.Contract.SRandomnessOutput(&_VRFOwnerlessConsumerExample.CallOpts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleCaller) SRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFOwnerlessConsumerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleSession) SRequestId() ([32]byte, error) {
	return _VRFOwnerlessConsumerExample.Contract.SRequestId(&_VRFOwnerlessConsumerExample.CallOpts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleCallerSession) SRequestId() ([32]byte, error) {
	return _VRFOwnerlessConsumerExample.Contract.SRequestId(&_VRFOwnerlessConsumerExample.CallOpts)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.contract.Transact(opts, "onTokenTransfer", arg0, _amount, _data)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.OnTokenTransfer(&_VRFOwnerlessConsumerExample.TransactOpts, arg0, _amount, _data)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactorSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.OnTokenTransfer(&_VRFOwnerlessConsumerExample.TransactOpts, arg0, _amount, _data)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.RawFulfillRandomness(&_VRFOwnerlessConsumerExample.TransactOpts, requestId, randomness)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExampleTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFOwnerlessConsumerExample.Contract.RawFulfillRandomness(&_VRFOwnerlessConsumerExample.TransactOpts, requestId, randomness)
}

func (_VRFOwnerlessConsumerExample *VRFOwnerlessConsumerExample) Address() common.Address {
	return _VRFOwnerlessConsumerExample.address
}

type VRFOwnerlessConsumerExampleInterface interface {
	SRandomnessOutput(opts *bind.CallOpts) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) ([32]byte, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error)

	Address() common.Address
}
