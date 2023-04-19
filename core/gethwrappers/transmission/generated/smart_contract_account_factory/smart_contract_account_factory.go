// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smart_contract_account_factory

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

var SmartContractAccountFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DeploymentFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"scaAddress\",\"type\":\"address\"}],\"name\":\"ContractCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"abiEncodedOwnerAddress\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deploySmartContractAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"scaAddress\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061021e806100206000396000f3fe60806040526004361061001e5760003560e01c80630af4926f14610023575b600080fd5b610036610031366004610138565b61005f565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6000828251836020016000f5905073ffffffffffffffffffffffffffffffffffffffff81166100ba576040517f3011642500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff821681527fcf78cf0d6f3d8371e1075c69c492ab4ec5d8cf23a1a239b6a51a1d00be7ca3129060200160405180910390a192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561014b57600080fd5b82359150602083013567ffffffffffffffff8082111561016a57600080fd5b818501915085601f83011261017e57600080fd5b81358181111561019057610190610109565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156101d6576101d6610109565b816040528281528860208487010111156101ef57600080fd5b826020860160208301376000602084830101528095505050505050925092905056fea164736f6c634300080f000a",
}

var SmartContractAccountFactoryABI = SmartContractAccountFactoryMetaData.ABI

var SmartContractAccountFactoryBin = SmartContractAccountFactoryMetaData.Bin

func DeploySmartContractAccountFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SmartContractAccountFactory, error) {
	parsed, err := SmartContractAccountFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SmartContractAccountFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SmartContractAccountFactory{SmartContractAccountFactoryCaller: SmartContractAccountFactoryCaller{contract: contract}, SmartContractAccountFactoryTransactor: SmartContractAccountFactoryTransactor{contract: contract}, SmartContractAccountFactoryFilterer: SmartContractAccountFactoryFilterer{contract: contract}}, nil
}

type SmartContractAccountFactory struct {
	address common.Address
	abi     abi.ABI
	SmartContractAccountFactoryCaller
	SmartContractAccountFactoryTransactor
	SmartContractAccountFactoryFilterer
}

type SmartContractAccountFactoryCaller struct {
	contract *bind.BoundContract
}

type SmartContractAccountFactoryTransactor struct {
	contract *bind.BoundContract
}

type SmartContractAccountFactoryFilterer struct {
	contract *bind.BoundContract
}

type SmartContractAccountFactorySession struct {
	Contract     *SmartContractAccountFactory
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SmartContractAccountFactoryCallerSession struct {
	Contract *SmartContractAccountFactoryCaller
	CallOpts bind.CallOpts
}

type SmartContractAccountFactoryTransactorSession struct {
	Contract     *SmartContractAccountFactoryTransactor
	TransactOpts bind.TransactOpts
}

type SmartContractAccountFactoryRaw struct {
	Contract *SmartContractAccountFactory
}

type SmartContractAccountFactoryCallerRaw struct {
	Contract *SmartContractAccountFactoryCaller
}

type SmartContractAccountFactoryTransactorRaw struct {
	Contract *SmartContractAccountFactoryTransactor
}

func NewSmartContractAccountFactory(address common.Address, backend bind.ContractBackend) (*SmartContractAccountFactory, error) {
	abi, err := abi.JSON(strings.NewReader(SmartContractAccountFactoryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSmartContractAccountFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountFactory{address: address, abi: abi, SmartContractAccountFactoryCaller: SmartContractAccountFactoryCaller{contract: contract}, SmartContractAccountFactoryTransactor: SmartContractAccountFactoryTransactor{contract: contract}, SmartContractAccountFactoryFilterer: SmartContractAccountFactoryFilterer{contract: contract}}, nil
}

func NewSmartContractAccountFactoryCaller(address common.Address, caller bind.ContractCaller) (*SmartContractAccountFactoryCaller, error) {
	contract, err := bindSmartContractAccountFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountFactoryCaller{contract: contract}, nil
}

func NewSmartContractAccountFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SmartContractAccountFactoryTransactor, error) {
	contract, err := bindSmartContractAccountFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountFactoryTransactor{contract: contract}, nil
}

func NewSmartContractAccountFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SmartContractAccountFactoryFilterer, error) {
	contract, err := bindSmartContractAccountFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountFactoryFilterer{contract: contract}, nil
}

func bindSmartContractAccountFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SmartContractAccountFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SmartContractAccountFactory.Contract.SmartContractAccountFactoryCaller.contract.Call(opts, result, method, params...)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.SmartContractAccountFactoryTransactor.contract.Transfer(opts)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.SmartContractAccountFactoryTransactor.contract.Transact(opts, method, params...)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SmartContractAccountFactory.Contract.contract.Call(opts, result, method, params...)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.contract.Transfer(opts)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.contract.Transact(opts, method, params...)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryTransactor) DeploySmartContractAccount(opts *bind.TransactOpts, abiEncodedOwnerAddress [32]byte, initCode []byte) (*types.Transaction, error) {
	return _SmartContractAccountFactory.contract.Transact(opts, "deploySmartContractAccount", abiEncodedOwnerAddress, initCode)
}

func (_SmartContractAccountFactory *SmartContractAccountFactorySession) DeploySmartContractAccount(abiEncodedOwnerAddress [32]byte, initCode []byte) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.DeploySmartContractAccount(&_SmartContractAccountFactory.TransactOpts, abiEncodedOwnerAddress, initCode)
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryTransactorSession) DeploySmartContractAccount(abiEncodedOwnerAddress [32]byte, initCode []byte) (*types.Transaction, error) {
	return _SmartContractAccountFactory.Contract.DeploySmartContractAccount(&_SmartContractAccountFactory.TransactOpts, abiEncodedOwnerAddress, initCode)
}

type SmartContractAccountFactoryContractCreatedIterator struct {
	Event *SmartContractAccountFactoryContractCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SmartContractAccountFactoryContractCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmartContractAccountFactoryContractCreated)
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
		it.Event = new(SmartContractAccountFactoryContractCreated)
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

func (it *SmartContractAccountFactoryContractCreatedIterator) Error() error {
	return it.fail
}

func (it *SmartContractAccountFactoryContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SmartContractAccountFactoryContractCreated struct {
	ScaAddress common.Address
	Raw        types.Log
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryFilterer) FilterContractCreated(opts *bind.FilterOpts) (*SmartContractAccountFactoryContractCreatedIterator, error) {

	logs, sub, err := _SmartContractAccountFactory.contract.FilterLogs(opts, "ContractCreated")
	if err != nil {
		return nil, err
	}
	return &SmartContractAccountFactoryContractCreatedIterator{contract: _SmartContractAccountFactory.contract, event: "ContractCreated", logs: logs, sub: sub}, nil
}

func (_SmartContractAccountFactory *SmartContractAccountFactoryFilterer) WatchContractCreated(opts *bind.WatchOpts, sink chan<- *SmartContractAccountFactoryContractCreated) (event.Subscription, error) {

	logs, sub, err := _SmartContractAccountFactory.contract.WatchLogs(opts, "ContractCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SmartContractAccountFactoryContractCreated)
				if err := _SmartContractAccountFactory.contract.UnpackLog(event, "ContractCreated", log); err != nil {
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

func (_SmartContractAccountFactory *SmartContractAccountFactoryFilterer) ParseContractCreated(log types.Log) (*SmartContractAccountFactoryContractCreated, error) {
	event := new(SmartContractAccountFactoryContractCreated)
	if err := _SmartContractAccountFactory.contract.UnpackLog(event, "ContractCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SmartContractAccountFactory *SmartContractAccountFactory) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SmartContractAccountFactory.abi.Events["ContractCreated"].ID:
		return _SmartContractAccountFactory.ParseContractCreated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SmartContractAccountFactoryContractCreated) Topic() common.Hash {
	return common.HexToHash("0xcf78cf0d6f3d8371e1075c69c492ab4ec5d8cf23a1a239b6a51a1d00be7ca312")
}

func (_SmartContractAccountFactory *SmartContractAccountFactory) Address() common.Address {
	return _SmartContractAccountFactory.address
}

type SmartContractAccountFactoryInterface interface {
	DeploySmartContractAccount(opts *bind.TransactOpts, abiEncodedOwnerAddress [32]byte, initCode []byte) (*types.Transaction, error)

	FilterContractCreated(opts *bind.FilterOpts) (*SmartContractAccountFactoryContractCreatedIterator, error)

	WatchContractCreated(opts *bind.WatchOpts, sink chan<- *SmartContractAccountFactoryContractCreated) (event.Subscription, error)

	ParseContractCreated(log types.Log) (*SmartContractAccountFactoryContractCreated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
