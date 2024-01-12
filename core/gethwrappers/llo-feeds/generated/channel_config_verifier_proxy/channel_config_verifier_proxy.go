// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package channel_config_verifier_proxy

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

var ChannelVerifierProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessForbidden\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"payloads\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"parameterPayload\",\"type\":\"bytes\"}],\"name\":\"verifyBulk\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"verifiedReports\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6109cb806101576000396000f3fe6080604052600436106100705760003560e01c80638da5cb5b1161004e5780638da5cb5b14610152578063f2fde38b14610187578063f7e83aee146101a7578063f873a61c146101ba57600080fd5b806301ffc9a714610075578063181f5a77146100ec57806379ba50971461013b575b600080fd5b34801561008157600080fd5b506100d76100903660046105f0565b7fffffffff00000000000000000000000000000000000000000000000000000000167f0f9b9cf2000000000000000000000000000000000000000000000000000000001490565b60405190151581526020015b60405180910390f35b3480156100f857600080fd5b5060408051808201909152601a81527f4368616e6e656c566572696669657250726f787920302e302e3000000000000060208201525b6040516100e3919061069d565b34801561014757600080fd5b506101506101da565b005b34801561015e57600080fd5b5060005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e3565b34801561019357600080fd5b506101506101a23660046106b0565b6102dc565b61012e6101b536600461072f565b6102f0565b6101cd6101c836600461079b565b61031c565b6040516100e3919061081c565b60015473ffffffffffffffffffffffffffffffffffffffff163314610260576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6102e46103ee565b6102ed81610471565b50565b60606102fe85858585610566565b505060408051600080825260208201909252905b5095945050505050565b60608367ffffffffffffffff8111156103375761033761089c565b60405190808252806020026020018201604052801561036a57816020015b60608152602001906001900390816103555790505b50905060005b848110156103b9576103a686868381811061038d5761038d6108cb565b905060200281019061039f91906108fa565b8686610566565b5050806103b29061095f565b9050610370565b506040805160008082526020820190925290610312565b60608152602001906001900390816103d05750909695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461046f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610257565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036104f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610257565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6060808585858583838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516020601f87018190048102820181019092528581529397509495509293919250849184915081908401838280828437600092019190915250959c929b50919950505050505050505050565b60006020828403121561060257600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461063257600080fd5b9392505050565b6000815180845260005b8181101561065f57602081850181015186830182015201610643565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006106326020830184610639565b6000602082840312156106c257600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461063257600080fd5b60008083601f8401126106f857600080fd5b50813567ffffffffffffffff81111561071057600080fd5b60208301915083602082850101111561072857600080fd5b9250929050565b6000806000806040858703121561074557600080fd5b843567ffffffffffffffff8082111561075d57600080fd5b610769888389016106e6565b9096509450602087013591508082111561078257600080fd5b5061078f878288016106e6565b95989497509550505050565b600080600080604085870312156107b157600080fd5b843567ffffffffffffffff808211156107c957600080fd5b818701915087601f8301126107dd57600080fd5b8135818111156107ec57600080fd5b8860208260051b850101111561080157600080fd5b60209283019650945090860135908082111561078257600080fd5b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561088f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261087d858351610639565b94509285019290850190600101610843565b5092979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261092f57600080fd5b83018035915067ffffffffffffffff82111561094a57600080fd5b60200191503681900382131561072857600080fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036109b7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea164736f6c6343000810000a",
}

var ChannelVerifierProxyABI = ChannelVerifierProxyMetaData.ABI

var ChannelVerifierProxyBin = ChannelVerifierProxyMetaData.Bin

func DeployChannelVerifierProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChannelVerifierProxy, error) {
	parsed, err := ChannelVerifierProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChannelVerifierProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChannelVerifierProxy{address: address, abi: *parsed, ChannelVerifierProxyCaller: ChannelVerifierProxyCaller{contract: contract}, ChannelVerifierProxyTransactor: ChannelVerifierProxyTransactor{contract: contract}, ChannelVerifierProxyFilterer: ChannelVerifierProxyFilterer{contract: contract}}, nil
}

type ChannelVerifierProxy struct {
	address common.Address
	abi     abi.ABI
	ChannelVerifierProxyCaller
	ChannelVerifierProxyTransactor
	ChannelVerifierProxyFilterer
}

type ChannelVerifierProxyCaller struct {
	contract *bind.BoundContract
}

type ChannelVerifierProxyTransactor struct {
	contract *bind.BoundContract
}

type ChannelVerifierProxyFilterer struct {
	contract *bind.BoundContract
}

type ChannelVerifierProxySession struct {
	Contract     *ChannelVerifierProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ChannelVerifierProxyCallerSession struct {
	Contract *ChannelVerifierProxyCaller
	CallOpts bind.CallOpts
}

type ChannelVerifierProxyTransactorSession struct {
	Contract     *ChannelVerifierProxyTransactor
	TransactOpts bind.TransactOpts
}

type ChannelVerifierProxyRaw struct {
	Contract *ChannelVerifierProxy
}

type ChannelVerifierProxyCallerRaw struct {
	Contract *ChannelVerifierProxyCaller
}

type ChannelVerifierProxyTransactorRaw struct {
	Contract *ChannelVerifierProxyTransactor
}

func NewChannelVerifierProxy(address common.Address, backend bind.ContractBackend) (*ChannelVerifierProxy, error) {
	abi, err := abi.JSON(strings.NewReader(ChannelVerifierProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindChannelVerifierProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxy{address: address, abi: abi, ChannelVerifierProxyCaller: ChannelVerifierProxyCaller{contract: contract}, ChannelVerifierProxyTransactor: ChannelVerifierProxyTransactor{contract: contract}, ChannelVerifierProxyFilterer: ChannelVerifierProxyFilterer{contract: contract}}, nil
}

func NewChannelVerifierProxyCaller(address common.Address, caller bind.ContractCaller) (*ChannelVerifierProxyCaller, error) {
	contract, err := bindChannelVerifierProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxyCaller{contract: contract}, nil
}

func NewChannelVerifierProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ChannelVerifierProxyTransactor, error) {
	contract, err := bindChannelVerifierProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxyTransactor{contract: contract}, nil
}

func NewChannelVerifierProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ChannelVerifierProxyFilterer, error) {
	contract, err := bindChannelVerifierProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxyFilterer{contract: contract}, nil
}

func bindChannelVerifierProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChannelVerifierProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ChannelVerifierProxy *ChannelVerifierProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelVerifierProxy.Contract.ChannelVerifierProxyCaller.contract.Call(opts, result, method, params...)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.ChannelVerifierProxyTransactor.contract.Transfer(opts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.ChannelVerifierProxyTransactor.contract.Transact(opts, method, params...)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChannelVerifierProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.contract.Transfer(opts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.contract.Transact(opts, method, params...)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChannelVerifierProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) Owner() (common.Address, error) {
	return _ChannelVerifierProxy.Contract.Owner(&_ChannelVerifierProxy.CallOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCallerSession) Owner() (common.Address, error) {
	return _ChannelVerifierProxy.Contract.Owner(&_ChannelVerifierProxy.CallOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ChannelVerifierProxy.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelVerifierProxy.Contract.SupportsInterface(&_ChannelVerifierProxy.CallOpts, interfaceId)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ChannelVerifierProxy.Contract.SupportsInterface(&_ChannelVerifierProxy.CallOpts, interfaceId)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ChannelVerifierProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) TypeAndVersion() (string, error) {
	return _ChannelVerifierProxy.Contract.TypeAndVersion(&_ChannelVerifierProxy.CallOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyCallerSession) TypeAndVersion() (string, error) {
	return _ChannelVerifierProxy.Contract.TypeAndVersion(&_ChannelVerifierProxy.CallOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChannelVerifierProxy.contract.Transact(opts, "acceptOwnership")
}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.AcceptOwnership(&_ChannelVerifierProxy.TransactOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.AcceptOwnership(&_ChannelVerifierProxy.TransactOpts)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ChannelVerifierProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.TransferOwnership(&_ChannelVerifierProxy.TransactOpts, to)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.TransferOwnership(&_ChannelVerifierProxy.TransactOpts, to)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactor) Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.contract.Transact(opts, "verify", payload, parameterPayload)
}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.Verify(&_ChannelVerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorSession) Verify(payload []byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.Verify(&_ChannelVerifierProxy.TransactOpts, payload, parameterPayload)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactor) VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.contract.Transact(opts, "verifyBulk", payloads, parameterPayload)
}

func (_ChannelVerifierProxy *ChannelVerifierProxySession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.VerifyBulk(&_ChannelVerifierProxy.TransactOpts, payloads, parameterPayload)
}

func (_ChannelVerifierProxy *ChannelVerifierProxyTransactorSession) VerifyBulk(payloads [][]byte, parameterPayload []byte) (*types.Transaction, error) {
	return _ChannelVerifierProxy.Contract.VerifyBulk(&_ChannelVerifierProxy.TransactOpts, payloads, parameterPayload)
}

type ChannelVerifierProxyOwnershipTransferRequestedIterator struct {
	Event *ChannelVerifierProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierProxyOwnershipTransferRequested)
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
		it.Event = new(ChannelVerifierProxyOwnershipTransferRequested)
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

func (it *ChannelVerifierProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxyOwnershipTransferRequestedIterator{contract: _ChannelVerifierProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierProxyOwnershipTransferRequested)
				if err := _ChannelVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*ChannelVerifierProxyOwnershipTransferRequested, error) {
	event := new(ChannelVerifierProxyOwnershipTransferRequested)
	if err := _ChannelVerifierProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ChannelVerifierProxyOwnershipTransferredIterator struct {
	Event *ChannelVerifierProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ChannelVerifierProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChannelVerifierProxyOwnershipTransferred)
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
		it.Event = new(ChannelVerifierProxyOwnershipTransferred)
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

func (it *ChannelVerifierProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *ChannelVerifierProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ChannelVerifierProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifierProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ChannelVerifierProxyOwnershipTransferredIterator{contract: _ChannelVerifierProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ChannelVerifierProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ChannelVerifierProxyOwnershipTransferred)
				if err := _ChannelVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_ChannelVerifierProxy *ChannelVerifierProxyFilterer) ParseOwnershipTransferred(log types.Log) (*ChannelVerifierProxyOwnershipTransferred, error) {
	event := new(ChannelVerifierProxyOwnershipTransferred)
	if err := _ChannelVerifierProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ChannelVerifierProxy *ChannelVerifierProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ChannelVerifierProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _ChannelVerifierProxy.ParseOwnershipTransferRequested(log)
	case _ChannelVerifierProxy.abi.Events["OwnershipTransferred"].ID:
		return _ChannelVerifierProxy.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ChannelVerifierProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (ChannelVerifierProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_ChannelVerifierProxy *ChannelVerifierProxy) Address() common.Address {
	return _ChannelVerifierProxy.address
}

type ChannelVerifierProxyInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Verify(opts *bind.TransactOpts, payload []byte, parameterPayload []byte) (*types.Transaction, error)

	VerifyBulk(opts *bind.TransactOpts, payloads [][]byte, parameterPayload []byte) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ChannelVerifierProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*ChannelVerifierProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ChannelVerifierProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ChannelVerifierProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*ChannelVerifierProxyOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
