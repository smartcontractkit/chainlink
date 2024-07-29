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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"UnauthorizedWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"}],\"name\":\"UnauthorizedWorkflowOwner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"price\",\"type\":\"uint224\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"FeedReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_allowedSendersList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_allowedWorkflowOwnersList\",\"type\":\"address[]\"},{\"internalType\":\"bytes10[]\",\"name\":\"_allowedWorkflowNamesList\",\"type\":\"bytes10[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6111cf806101576000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c80638da5cb5b116100505780638da5cb5b1461014e578063e340171114610176578063f2fde38b1461018957600080fd5b806331d98b3f1461007757806379ba509714610131578063805f21321461013b575b600080fd5b6100f3610085366004610d64565b6000908152600260209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff9091166020830152015b60405180910390f35b61013961019c565b005b610139610149366004610dc6565b61029e565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610128565b610139610184366004610e77565b61061d565b610139610197366004610f11565b610a5d565b60015473ffffffffffffffffffffffffffffffffffffffff163314610222576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081526004602052604090205460ff166102e9576040517f3fcc3f17000000000000000000000000000000000000000000000000000000008152336004820152602401610219565b60008061032b86868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610a7192505050565b7fffffffffffffffffffff000000000000000000000000000000000000000000008216600090815260086020526040902054919350915060ff166103bf576040517f4b942f800000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffff0000000000000000000000000000000000000000000083166004820152602401610219565b73ffffffffffffffffffffffffffffffffffffffff811660009081526006602052604090205460ff16610436576040517fbf24162300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610219565b600061044484860186610ff5565b905060005b815181101561061357604051806040016040528083838151811061046f5761046f611107565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018383815181106104b0576104b0611107565b60200260200101516040015163ffffffff16815250600260008484815181106104db576104db611107565b602090810291909101810151518252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055815182908290811061055d5761055d611107565b6020026020010151600001517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d83838151811061059c5761059c611107565b6020026020010151602001518484815181106105ba576105ba611107565b6020026020010151604001516040516106039291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a2600101610449565b5050505050505050565b610625610a87565b60005b60035463ffffffff821610156106c65760006004600060038463ffffffff168154811061065757610657611107565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556106bf81611136565b9050610628565b5060005b63ffffffff811686111561076e5760016004600089898563ffffffff168181106106f6576106f6611107565b905060200201602081019061070b9190610f11565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905561076781611136565b90506106ca565b5061077b60038787610bff565b5060005b60055463ffffffff8216101561081d5760006006600060058463ffffffff16815481106107ae576107ae611107565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905561081681611136565b905061077f565b5060005b63ffffffff81168411156108c55760016006600087878563ffffffff1681811061084d5761084d611107565b90506020020160208101906108629190610f11565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556108be81611136565b9050610821565b506108d260058585610bff565b5060005b60075463ffffffff821610156109935760006008600060078463ffffffff168154811061090557610905611107565b600091825260208083206003808404909101549206600a026101000a90910460b01b7fffffffffffffffffffff00000000000000000000000000000000000000000000168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905561098c81611136565b90506108d6565b5060005b63ffffffff8116821115610a475760016008600085858563ffffffff168181106109c3576109c3611107565b90506020020160208101906109d89190611180565b7fffffffffffffffffffff00000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610a4081611136565b9050610997565b50610a5460078383610c87565b50505050505050565b610a65610a87565b610a6e81610b0a565b50565b6040810151604a90910151909160609190911c90565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610219565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610b89576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610219565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610c77579160200282015b82811115610c775781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610c1f565b50610c83929150610d4f565b5090565b82805482825590600052602060002090600201600390048101928215610c775791602002820160005b83821115610d1057833575ffffffffffffffffffffffffffffffffffffffffffff191683826101000a81548169ffffffffffffffffffff021916908360b01c02179055509260200192600a01602081600901049283019260010302610cb0565b8015610d465782816101000a81549069ffffffffffffffffffff0219169055600a01602081600901049283019260010302610d10565b5050610c839291505b5b80821115610c835760008155600101610d50565b600060208284031215610d7657600080fd5b5035919050565b60008083601f840112610d8f57600080fd5b50813567ffffffffffffffff811115610da757600080fd5b602083019150836020828501011115610dbf57600080fd5b9250929050565b60008060008060408587031215610ddc57600080fd5b843567ffffffffffffffff80821115610df457600080fd5b610e0088838901610d7d565b90965094506020870135915080821115610e1957600080fd5b50610e2687828801610d7d565b95989497509550505050565b60008083601f840112610e4457600080fd5b50813567ffffffffffffffff811115610e5c57600080fd5b6020830191508360208260051b8501011115610dbf57600080fd5b60008060008060008060608789031215610e9057600080fd5b863567ffffffffffffffff80821115610ea857600080fd5b610eb48a838b01610e32565b90985096506020890135915080821115610ecd57600080fd5b610ed98a838b01610e32565b90965094506040890135915080821115610ef257600080fd5b50610eff89828a01610e32565b979a9699509497509295939492505050565b600060208284031215610f2357600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610f4757600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715610fa057610fa0610f4e565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610fed57610fed610f4e565b604052919050565b6000602080838503121561100857600080fd5b823567ffffffffffffffff8082111561102057600080fd5b818501915085601f83011261103457600080fd5b81358181111561104657611046610f4e565b611054848260051b01610fa6565b8181528481019250606091820284018501918883111561107357600080fd5b938501935b828510156110fb5780858a0312156110905760008081fd5b611098610f7d565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146110cb5760008081fd5b8188015260408681013563ffffffff811681146110e85760008081fd5b9082015284529384019392850192611078565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600063ffffffff808316818103611176577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b60006020828403121561119257600080fd5b81357fffffffffffffffffffff0000000000000000000000000000000000000000000081168114610f4757600080fdfea164736f6c6343000818000a",
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
