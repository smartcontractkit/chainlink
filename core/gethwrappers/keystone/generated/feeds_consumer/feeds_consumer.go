// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feeds_consumer

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

var KeystoneFeedsConsumerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPrice\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"_allowedSendersList\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_allowedWorkflowOwnersList\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_allowedWorkflowNamesList\",\"type\":\"bytes10[]\",\"internalType\":\"bytes10[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeedReceived\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"price\",\"type\":\"uint224\",\"indexed\":false,\"internalType\":\"uint224\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"UnauthorizedSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedWorkflowName\",\"inputs\":[{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedWorkflowOwner\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f80fd5b5033805f816100655760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b5f80546001600160a01b0319166001600160a01b0384811691909117909155811615610094576100948161009c565b505050610144565b336001600160a01b038216036100f45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005c565b600180546001600160a01b0319166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611278806101515f395ff3fe608060405234801561000f575f80fd5b506004361061007a575f3560e01c8063805f213211610058578063805f2132146101645780638da5cb5b14610177578063e34017111461019e578063f2fde38b146101b1575f80fd5b806301ffc9a71461007e57806331d98b3f146100a657806379ba50971461015a575b5f80fd5b61009161008c366004610dfc565b6101c4565b60405190151581526020015b60405180910390f35b6101216100b4366004610e42565b5f908152600260209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff90911660208301520161009d565b61016261025c565b005b610162610172366004610e9e565b61035d565b5f5460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161009d565b6101626101ac366004610f46565b6106d2565b6101626101bf366004610fd9565b610afd565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000148061025657507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102e2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b5f8054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b335f9081526004602052604090205460ff166103a7576040517f3fcc3f170000000000000000000000000000000000000000000000000000000081523360048201526024016102d9565b5f806103e786868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250610b1192505050565b7fffffffffffffffffffff0000000000000000000000000000000000000000000082165f90815260086020526040902054919350915060ff1661047a576040517f4b942f800000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffff00000000000000000000000000000000000000000000831660048201526024016102d9565b73ffffffffffffffffffffffffffffffffffffffff81165f9081526006602052604090205460ff166104f0576040517fbf24162300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016102d9565b5f6104fd848601866110b1565b90505f5b81518110156106c8576040518060400160405280838381518110610527576105276111b8565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001838381518110610568576105686111b8565b60200260200101516040015163ffffffff1681525060025f848481518110610592576105926111b8565b602090810291909101810151518252818101929092526040015f208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790558151829082908110610613576106136111b8565b60200260200101515f01517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d838381518110610651576106516111b8565b60200260200101516020015184848151811061066f5761066f6111b8565b6020026020010151604001516040516106b89291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a2600101610501565b5050505050505050565b6106da610b27565b5f5b60035463ffffffff82161015610777575f60045f60038463ffffffff1681548110610709576107096111b8565b5f9182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610770816111e5565b90506106dc565b505f5b63ffffffff811686111561081c57600160045f89898563ffffffff168181106107a5576107a56111b8565b90506020020160208101906107ba9190610fd9565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040015f2080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610815816111e5565b905061077a565b5061082960038787610c9d565b505f5b60055463ffffffff821610156108c7575f60065f60058463ffffffff1681548110610859576108596111b8565b5f9182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556108c0816111e5565b905061082c565b505f5b63ffffffff811684111561096c57600160065f87878563ffffffff168181106108f5576108f56111b8565b905060200201602081019061090a9190610fd9565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040015f2080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610965816111e5565b90506108ca565b5061097960058585610c9d565b505f5b60075463ffffffff82161015610a36575f60085f60078463ffffffff16815481106109a9576109a96111b8565b5f91825260208083206003808404909101549206600a026101000a90910460b01b7fffffffffffffffffffff00000000000000000000000000000000000000000000168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610a2f816111e5565b905061097c565b505f5b63ffffffff8116821115610ae757600160085f85858563ffffffff16818110610a6457610a646111b8565b9050602002016020810190610a79919061122c565b7fffffffffffffffffffff0000000000000000000000000000000000000000000016815260208101919091526040015f2080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610ae0816111e5565b9050610a39565b50610af460078383610d23565b50505050505050565b610b05610b27565b610b0e81610ba9565b50565b6040810151604a90910151909160609190911c90565b5f5473ffffffffffffffffffffffffffffffffffffffff163314610ba7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102d9565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610c28576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102d9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255905f5260205f20908101928215610d13579160200282015b82811115610d135781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610cbb565b50610d1f929150610de8565b5090565b828054828255905f5260205f2090600201600390048101928215610d13579160200282015f5b83821115610da957833575ffffffffffffffffffffffffffffffffffffffffffff191683826101000a81548169ffffffffffffffffffff021916908360b01c02179055509260200192600a01602081600901049283019260010302610d49565b8015610ddf5782816101000a81549069ffffffffffffffffffff0219169055600a01602081600901049283019260010302610da9565b5050610d1f9291505b5b80821115610d1f575f8155600101610de9565b5f60208284031215610e0c575f80fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610e3b575f80fd5b9392505050565b5f60208284031215610e52575f80fd5b5035919050565b5f8083601f840112610e69575f80fd5b50813567ffffffffffffffff811115610e80575f80fd5b602083019150836020828501011115610e97575f80fd5b9250929050565b5f805f8060408587031215610eb1575f80fd5b843567ffffffffffffffff80821115610ec8575f80fd5b610ed488838901610e59565b90965094506020870135915080821115610eec575f80fd5b50610ef987828801610e59565b95989497509550505050565b5f8083601f840112610f15575f80fd5b50813567ffffffffffffffff811115610f2c575f80fd5b6020830191508360208260051b8501011115610e97575f80fd5b5f805f805f8060608789031215610f5b575f80fd5b863567ffffffffffffffff80821115610f72575f80fd5b610f7e8a838b01610f05565b90985096506020890135915080821115610f96575f80fd5b610fa28a838b01610f05565b90965094506040890135915080821115610fba575f80fd5b50610fc789828a01610f05565b979a9699509497509295939492505050565b5f60208284031215610fe9575f80fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610e3b575f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040516060810167ffffffffffffffff8111828210171561105c5761105c61100c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156110a9576110a961100c565b604052919050565b5f60208083850312156110c2575f80fd5b823567ffffffffffffffff808211156110d9575f80fd5b818501915085601f8301126110ec575f80fd5b8135818111156110fe576110fe61100c565b61110c848260051b01611062565b8181528481019250606091820284018501918883111561112a575f80fd5b938501935b828510156111ac5780858a031215611145575f80fd5b61114d611039565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461117e575f80fd5b8188015260408681013563ffffffff81168114611199575f80fd5b908201528452938401939285019261112f565b50979650505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f63ffffffff808316818103611222577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b6001019392505050565b5f6020828403121561123c575f80fd5b81357fffffffffffffffffffff0000000000000000000000000000000000000000000081168114610e3b575f80fdfea164736f6c6343000818000a",
}

var KeystoneFeedsConsumerABI = KeystoneFeedsConsumerMetaData.ABI

var KeystoneFeedsConsumerBin = KeystoneFeedsConsumerMetaData.Bin

func DeployKeystoneFeedsConsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneFeedsConsumer, error) {
	parsed, err := KeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneFeedsConsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneFeedsConsumer{address: address, abi: *parsed, KeystoneFeedsConsumerCaller: KeystoneFeedsConsumerCaller{contract: contract}, KeystoneFeedsConsumerTransactor: KeystoneFeedsConsumerTransactor{contract: contract}, KeystoneFeedsConsumerFilterer: KeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

type KeystoneFeedsConsumer struct {
	address common.Address
	abi     abi.ABI
	KeystoneFeedsConsumerCaller
	KeystoneFeedsConsumerTransactor
	KeystoneFeedsConsumerFilterer
}

type KeystoneFeedsConsumerCaller struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerTransactor struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerFilterer struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerSession struct {
	Contract     *KeystoneFeedsConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneFeedsConsumerCallerSession struct {
	Contract *KeystoneFeedsConsumerCaller
	CallOpts bind.CallOpts
}

type KeystoneFeedsConsumerTransactorSession struct {
	Contract     *KeystoneFeedsConsumerTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneFeedsConsumerRaw struct {
	Contract *KeystoneFeedsConsumer
}

type KeystoneFeedsConsumerCallerRaw struct {
	Contract *KeystoneFeedsConsumerCaller
}

type KeystoneFeedsConsumerTransactorRaw struct {
	Contract *KeystoneFeedsConsumerTransactor
}

func NewKeystoneFeedsConsumer(address common.Address, backend bind.ContractBackend) (*KeystoneFeedsConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneFeedsConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneFeedsConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumer{address: address, abi: abi, KeystoneFeedsConsumerCaller: KeystoneFeedsConsumerCaller{contract: contract}, KeystoneFeedsConsumerTransactor: KeystoneFeedsConsumerTransactor{contract: contract}, KeystoneFeedsConsumerFilterer: KeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

func NewKeystoneFeedsConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeystoneFeedsConsumerCaller, error) {
	contract, err := bindKeystoneFeedsConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerCaller{contract: contract}, nil
}

func NewKeystoneFeedsConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneFeedsConsumerTransactor, error) {
	contract, err := bindKeystoneFeedsConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerTransactor{contract: contract}, nil
}

func NewKeystoneFeedsConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneFeedsConsumerFilterer, error) {
	contract, err := bindKeystoneFeedsConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerFilterer{contract: contract}, nil
}

func bindKeystoneFeedsConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerTransactor.contract.Transfer(opts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneFeedsConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.contract.Transfer(opts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "getPrice", feedId)

	if err != nil {
		return *new(*big.Int), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetPrice(&_KeystoneFeedsConsumer.CallOpts, feedId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetPrice(&_KeystoneFeedsConsumer.CallOpts, feedId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) Owner() (common.Address, error) {
	return _KeystoneFeedsConsumer.Contract.Owner(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) Owner() (common.Address, error) {
	return _KeystoneFeedsConsumer.Contract.Owner(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KeystoneFeedsConsumer.Contract.SupportsInterface(&_KeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KeystoneFeedsConsumer.Contract.SupportsInterface(&_KeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.AcceptOwnership(&_KeystoneFeedsConsumer.TransactOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.AcceptOwnership(&_KeystoneFeedsConsumer.TransactOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "onReport", metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.OnReport(&_KeystoneFeedsConsumer.TransactOpts, metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.OnReport(&_KeystoneFeedsConsumer.TransactOpts, metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) SetConfig(opts *bind.TransactOpts, _allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "setConfig", _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.SetConfig(&_KeystoneFeedsConsumer.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.SetConfig(&_KeystoneFeedsConsumer.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.TransferOwnership(&_KeystoneFeedsConsumer.TransactOpts, to)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.TransferOwnership(&_KeystoneFeedsConsumer.TransactOpts, to)
}

type KeystoneFeedsConsumerFeedReceivedIterator struct {
	Event *KeystoneFeedsConsumerFeedReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerFeedReceived)
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
		it.Event = new(KeystoneFeedsConsumerFeedReceived)
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

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerFeedReceived struct {
	FeedId    [32]byte
	Price     *big.Int
	Timestamp uint32
	Raw       types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*KeystoneFeedsConsumerFeedReceivedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerFeedReceivedIterator{contract: _KeystoneFeedsConsumer.contract, event: "FeedReceived", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerFeedReceived)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseFeedReceived(log types.Log) (*KeystoneFeedsConsumerFeedReceived, error) {
	event := new(KeystoneFeedsConsumerFeedReceived)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneFeedsConsumerOwnershipTransferRequestedIterator struct {
	Event *KeystoneFeedsConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerOwnershipTransferRequested)
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
		it.Event = new(KeystoneFeedsConsumerOwnershipTransferRequested)
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

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerOwnershipTransferRequestedIterator{contract: _KeystoneFeedsConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerOwnershipTransferRequested)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferRequested, error) {
	event := new(KeystoneFeedsConsumerOwnershipTransferRequested)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneFeedsConsumerOwnershipTransferredIterator struct {
	Event *KeystoneFeedsConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerOwnershipTransferred)
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
		it.Event = new(KeystoneFeedsConsumerOwnershipTransferred)
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

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerOwnershipTransferredIterator{contract: _KeystoneFeedsConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerOwnershipTransferred)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferred, error) {
	event := new(KeystoneFeedsConsumerOwnershipTransferred)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneFeedsConsumer.abi.Events["FeedReceived"].ID:
		return _KeystoneFeedsConsumer.ParseFeedReceived(log)
	case _KeystoneFeedsConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneFeedsConsumer.ParseOwnershipTransferRequested(log)
	case _KeystoneFeedsConsumer.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneFeedsConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneFeedsConsumerFeedReceived) Topic() common.Hash {
	return common.HexToHash("0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d")
}

func (KeystoneFeedsConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneFeedsConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumer) Address() common.Address {
	return _KeystoneFeedsConsumer.address
}

type KeystoneFeedsConsumerInterface interface {
	GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*KeystoneFeedsConsumerFeedReceivedIterator, error)

	WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error)

	ParseFeedReceived(log types.Log) (*KeystoneFeedsConsumerFeedReceived, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
