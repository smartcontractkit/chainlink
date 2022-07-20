// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package authorized_forwarder

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var AuthorizedForwarderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"OwnershipTransferRequestedWithMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"ownerForward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"transferOwnershipWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200126138038062001261833981810160405260808110156200003757600080fd5b8151602083015160408085015160608601805192519496939591949391820192846401000000008211156200006b57600080fd5b9083019060208201858111156200008157600080fd5b82516401000000008111828201881017156200009c57600080fd5b82525081516020918201929091019080838360005b83811015620000cb578181015183820152602001620000b1565b50505050905090810190601f168015620000f95780820380516001836020036101000a031916815260200191505b50604052508491508390506001600160a01b03821662000160576040805162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f0000000000000000604482015290519081900360640190fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200019357620001938162000272565b50506001600160601b0319606085901b166080526001600160a01b038216156200026857816001600160a01b0316836001600160a01b03167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e836040518080602001828103825283818151815260200191508051906020019080838360005b838110156200022c57818101518382015260200162000212565b50505050905090810190601f1680156200025a5780820380516001836020036101000a031916815260200191505b509250505060405180910390a35b5050505062000322565b6001600160a01b038116331415620002d1576040805162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60805160601c610f1c6200034560003980610491528061061d5250610f1c6000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636fadcf7211610081578063ee56997b1161005b578063ee56997b1461038d578063f2fde38b146103fd578063fa00763a14610430576100c9565b80636fadcf72146102f057806379ba50971461037d5780638da5cb5b14610385576100c9565b8063181f5a77116100b2578063181f5a771461018e5780632408afaa1461020b5780634d3e232314610263576100c9565b8063033f49f7146100ce578063165d35e11461015d575b600080fd5b61015b600480360360408110156100e457600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561011c57600080fd5b82018360208201111561012e57600080fd5b8035906020019184600183028401116401000000008311171561015057600080fd5b509092509050610477565b005b61016561048f565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101966104b3565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101d05781810151838201526020016101b8565b50505050905090810190601f1680156101fd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102136104ea565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561024f578181015183820152602001610237565b505050509050019250505060405180910390f35b61015b6004803603604081101561027957600080fd5b73ffffffffffffffffffffffffffffffffffffffff82351691908101906040810160208201356401000000008111156102b157600080fd5b8201836020820111156102c357600080fd5b803590602001918460018302840111640100000000831117156102e557600080fd5b509092509050610559565b61015b6004803603604081101561030657600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561033e57600080fd5b82018360208201111561035057600080fd5b8035906020019184600183028401116401000000008311171561037257600080fd5b509092509050610613565b61015b6106d6565b6101656107d8565b61015b600480360360208110156103a357600080fd5b8101906020810181356401000000008111156103be57600080fd5b8201836020820111156103d057600080fd5b803590602001918460208302840111640100000000831117156103f257600080fd5b5090925090506107f4565b61015b6004803603602081101561041357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610a79565b6104636004803603602081101561044657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610a8d565b604080519115158252519081900360200190f35b61047f610ab8565b61048a838383610b40565b505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60408051808201909152601981527f417574686f72697a6564466f7277617264657220312e302e3000000000000000602082015290565b6060600380548060200260200160405190810160405280929190818152602001828054801561054f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610524575b5050505050905090565b61056283610a79565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e848460405180806020018281038252848482818152602001925080828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039550909350505050a3505050565b61061b610cb0565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561047f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f43616e6e6f7420666f727761726420746f204c696e6b20746f6b656e00000000604482015290519081900360640190fd5b60015473ffffffffffffffffffffffffffffffffffffffff16331461075c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b6107fc610d24565b61086757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f43616e6e6f742073657420617574686f72697a65642073656e64657273000000604482015290519081900360640190fd5b806108bd576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526026815260200180610eea6026913960400191505060405180910390fd5b60035460005b8181101561094557600060026000600384815481106108de57fe5b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556001016108c3565b5060005b828110156109c75760016002600086868581811061096357fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055600101610949565b506109d460038484610e4c565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a083833360405180806020018373ffffffffffffffffffffffffffffffffffffffff1681526020018281038252858582818152602001925060200280828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201829003965090945050505050a1505050565b610a81610ab8565b610a8a81610d4b565b50565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b3e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b565b610b5f8373ffffffffffffffffffffffffffffffffffffffff16610e46565b610bca57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d75737420666f727761726420746f206120636f6e7472616374000000000000604482015290519081900360640190fd5b60008373ffffffffffffffffffffffffffffffffffffffff168383604051808383808284376040519201945060009350909150508083038183865af19150503d8060008114610c35576040519150601f19603f3d011682016040523d82523d6000602084013e610c3a565b606091505b5050905080610caa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f466f727761726465642063616c6c206661696c65640000000000000000000000604482015290519081900360640190fd5b50505050565b610cb933610a8d565b610b3e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f4e6f7420617574686f72697a65642073656e6465720000000000000000000000604482015290519081900360640190fd5b600033610d2f6107d8565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b73ffffffffffffffffffffffffffffffffffffffff8116331415610dd057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b3b151590565b828054828255906000526020600020908101928215610ec4579160200282015b82811115610ec45781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610e6c565b50610ed0929150610ed4565b5090565b5b80821115610ed05760008155600101610ed556fe4d7573742068617665206174206c65617374203120617574686f72697a65642073656e646572a164736f6c6343000706000a",
}

var AuthorizedForwarderABI = AuthorizedForwarderMetaData.ABI

var AuthorizedForwarderBin = AuthorizedForwarderMetaData.Bin

func DeployAuthorizedForwarder(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, owner common.Address, recipient common.Address, message []byte) (common.Address, *types.Transaction, *AuthorizedForwarder, error) {
	parsed, err := AuthorizedForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AuthorizedForwarderBin), backend, link, owner, recipient, message)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AuthorizedForwarder{AuthorizedForwarderCaller: AuthorizedForwarderCaller{contract: contract}, AuthorizedForwarderTransactor: AuthorizedForwarderTransactor{contract: contract}, AuthorizedForwarderFilterer: AuthorizedForwarderFilterer{contract: contract}}, nil
}

type AuthorizedForwarder struct {
	address common.Address
	abi     abi.ABI
	AuthorizedForwarderCaller
	AuthorizedForwarderTransactor
	AuthorizedForwarderFilterer
}

type AuthorizedForwarderCaller struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderTransactor struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderFilterer struct {
	contract *bind.BoundContract
}

type AuthorizedForwarderSession struct {
	Contract     *AuthorizedForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AuthorizedForwarderCallerSession struct {
	Contract *AuthorizedForwarderCaller
	CallOpts bind.CallOpts
}

type AuthorizedForwarderTransactorSession struct {
	Contract     *AuthorizedForwarderTransactor
	TransactOpts bind.TransactOpts
}

type AuthorizedForwarderRaw struct {
	Contract *AuthorizedForwarder
}

type AuthorizedForwarderCallerRaw struct {
	Contract *AuthorizedForwarderCaller
}

type AuthorizedForwarderTransactorRaw struct {
	Contract *AuthorizedForwarderTransactor
}

func NewAuthorizedForwarder(address common.Address, backend bind.ContractBackend) (*AuthorizedForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(AuthorizedForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAuthorizedForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarder{address: address, abi: abi, AuthorizedForwarderCaller: AuthorizedForwarderCaller{contract: contract}, AuthorizedForwarderTransactor: AuthorizedForwarderTransactor{contract: contract}, AuthorizedForwarderFilterer: AuthorizedForwarderFilterer{contract: contract}}, nil
}

func NewAuthorizedForwarderCaller(address common.Address, caller bind.ContractCaller) (*AuthorizedForwarderCaller, error) {
	contract, err := bindAuthorizedForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderCaller{contract: contract}, nil
}

func NewAuthorizedForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*AuthorizedForwarderTransactor, error) {
	contract, err := bindAuthorizedForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderTransactor{contract: contract}, nil
}

func NewAuthorizedForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*AuthorizedForwarderFilterer, error) {
	contract, err := bindAuthorizedForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderFilterer{contract: contract}, nil
}

func bindAuthorizedForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuthorizedForwarderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderTransactor.contract.Transfer(opts)
}

func (_AuthorizedForwarder *AuthorizedForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AuthorizedForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.contract.Transfer(opts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedForwarder.Contract.GetAuthorizedSenders(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedForwarder.Contract.GetAuthorizedSenders(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) GetChainlinkToken() (common.Address, error) {
	return _AuthorizedForwarder.Contract.GetChainlinkToken(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) GetChainlinkToken() (common.Address, error) {
	return _AuthorizedForwarder.Contract.GetChainlinkToken(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedForwarder.Contract.IsAuthorizedSender(&_AuthorizedForwarder.CallOpts, sender)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedForwarder.Contract.IsAuthorizedSender(&_AuthorizedForwarder.CallOpts, sender)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) Owner() (common.Address, error) {
	return _AuthorizedForwarder.Contract.Owner(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) Owner() (common.Address, error) {
	return _AuthorizedForwarder.Contract.Owner(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AuthorizedForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TypeAndVersion() (string, error) {
	return _AuthorizedForwarder.Contract.TypeAndVersion(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderCallerSession) TypeAndVersion() (string, error) {
	return _AuthorizedForwarder.Contract.TypeAndVersion(&_AuthorizedForwarder.CallOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AcceptOwnership(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.AcceptOwnership(&_AuthorizedForwarder.TransactOpts)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "forward", to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.Forward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) Forward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.Forward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) OwnerForward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "ownerForward", to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) OwnerForward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.OwnerForward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) OwnerForward(to common.Address, data []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.OwnerForward(&_AuthorizedForwarder.TransactOpts, to, data)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.SetAuthorizedSenders(&_AuthorizedForwarder.TransactOpts, senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.SetAuthorizedSenders(&_AuthorizedForwarder.TransactOpts, senders)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnership(&_AuthorizedForwarder.TransactOpts, to)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnership(&_AuthorizedForwarder.TransactOpts, to)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactor) TransferOwnershipWithMessage(opts *bind.TransactOpts, to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.contract.Transact(opts, "transferOwnershipWithMessage", to, message)
}

func (_AuthorizedForwarder *AuthorizedForwarderSession) TransferOwnershipWithMessage(to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnershipWithMessage(&_AuthorizedForwarder.TransactOpts, to, message)
}

func (_AuthorizedForwarder *AuthorizedForwarderTransactorSession) TransferOwnershipWithMessage(to common.Address, message []byte) (*types.Transaction, error) {
	return _AuthorizedForwarder.Contract.TransferOwnershipWithMessage(&_AuthorizedForwarder.TransactOpts, to, message)
}

type AuthorizedForwarderAuthorizedSendersChangedIterator struct {
	Event *AuthorizedForwarderAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderAuthorizedSendersChanged)
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
		it.Event = new(AuthorizedForwarderAuthorizedSendersChanged)
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

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedForwarderAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderAuthorizedSendersChangedIterator{contract: _AuthorizedForwarder.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderAuthorizedSendersChanged)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedForwarderAuthorizedSendersChanged, error) {
	event := new(AuthorizedForwarderAuthorizedSendersChanged)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferRequestedIterator struct {
	Event *AuthorizedForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferRequested)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferRequested)
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

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferRequestedIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferRequested)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*AuthorizedForwarderOwnershipTransferRequested, error) {
	event := new(AuthorizedForwarderOwnershipTransferRequested)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator struct {
	Event *AuthorizedForwarderOwnershipTransferRequestedWithMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
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

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferRequestedWithMessage struct {
	From    common.Address
	To      common.Address
	Message []byte
	Raw     types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferRequestedWithMessage(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferRequestedWithMessage", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferRequestedWithMessage", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferRequestedWithMessage(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequestedWithMessage, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferRequestedWithMessage", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequestedWithMessage", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferRequestedWithMessage(log types.Log) (*AuthorizedForwarderOwnershipTransferRequestedWithMessage, error) {
	event := new(AuthorizedForwarderOwnershipTransferRequestedWithMessage)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferRequestedWithMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AuthorizedForwarderOwnershipTransferredIterator struct {
	Event *AuthorizedForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedForwarderOwnershipTransferred)
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
		it.Event = new(AuthorizedForwarderOwnershipTransferred)
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

func (it *AuthorizedForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AuthorizedForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AuthorizedForwarderOwnershipTransferredIterator{contract: _AuthorizedForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AuthorizedForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedForwarderOwnershipTransferred)
				if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AuthorizedForwarder *AuthorizedForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*AuthorizedForwarderOwnershipTransferred, error) {
	event := new(AuthorizedForwarderOwnershipTransferred)
	if err := _AuthorizedForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AuthorizedForwarder *AuthorizedForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AuthorizedForwarder.abi.Events["AuthorizedSendersChanged"].ID:
		return _AuthorizedForwarder.ParseAuthorizedSendersChanged(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferRequested(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferRequestedWithMessage"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferRequestedWithMessage(log)
	case _AuthorizedForwarder.abi.Events["OwnershipTransferred"].ID:
		return _AuthorizedForwarder.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AuthorizedForwarderAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (AuthorizedForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AuthorizedForwarderOwnershipTransferRequestedWithMessage) Topic() common.Hash {
	return common.HexToHash("0x4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e")
}

func (AuthorizedForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_AuthorizedForwarder *AuthorizedForwarder) Address() common.Address {
	return _AuthorizedForwarder.address
}

type AuthorizedForwarderInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetChainlinkToken(opts *bind.CallOpts) (common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Forward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

	OwnerForward(opts *bind.TransactOpts, to common.Address, data []byte) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferOwnershipWithMessage(opts *bind.TransactOpts, to common.Address, message []byte) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedForwarderAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedForwarderAuthorizedSendersChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AuthorizedForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferRequestedWithMessage(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferRequestedWithMessageIterator, error)

	WatchOwnershipTransferRequestedWithMessage(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferRequestedWithMessage, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequestedWithMessage(log types.Log) (*AuthorizedForwarderOwnershipTransferRequestedWithMessage, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AuthorizedForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuthorizedForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AuthorizedForwarderOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
