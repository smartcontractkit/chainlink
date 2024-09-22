// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package usdc_reader_tester

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

var USDCReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_version\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_localDomain\",\"type\":\"uint32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"destinationDomain\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"destinationCaller\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"messageBody\",\"type\":\"bytes\"}],\"name\":\"emitMessageSent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"localDomain\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161042938038061042983398101604081905261002f9161005c565b63ffffffff9182166080521660a05261008f565b805163ffffffff8116811461005757600080fd5b919050565b6000806040838503121561006f57600080fd5b61007883610043565b915061008660208401610043565b90509250929050565b60805160a05161036b6100be60003960008181608b015260ea015260008181604b015260c9015261036b6000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806354fd4d50146100465780638d3638f414610086578063c48ffbd6146100ad575b600080fd5b61006d7f000000000000000000000000000000000000000000000000000000000000000081565b60405163ffffffff909116815260200160405180910390f35b61006d7f000000000000000000000000000000000000000000000000000000000000000081565b6100c06100bb3660046101cb565b6100c2565b005b600061014a7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008a87898c8c8a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061018d92505050565b90507f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b0368160405161017b91906102b3565b60405180910390a15050505050505050565b606088888888888888886040516020016101ae9897969594939291906102e6565b604051602081830303815290604052905098975050505050505050565b600080600080600080600060c0888a0312156101e657600080fd5b873563ffffffff811681146101fa57600080fd5b965060208801359550604088013594506060880135935060808801356001600160401b03808216821461022c57600080fd5b90935060a0890135908082111561024257600080fd5b818a0191508a601f83011261025657600080fd5b81358181111561026557600080fd5b8b602082850101111561027757600080fd5b60208301945080935050505092959891949750929550565b60005b838110156102aa578181015183820152602001610292565b50506000910152565b60208152600082518060208401526102d281604085016020870161028f565b601f01601f19169190910160400192915050565b6001600160e01b031960e08a811b8216835289811b8216600484015288901b1660088201526001600160c01b031960c087901b16600c820152601481018590526034810184905260548101839052815160009061034a81607485016020870161028f565b91909101607401999850505050505050505056fea164736f6c6343000818000a",
}

var USDCReaderTesterABI = USDCReaderTesterMetaData.ABI

var USDCReaderTesterBin = USDCReaderTesterMetaData.Bin

func DeployUSDCReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend, _version uint32, _localDomain uint32) (common.Address, *types.Transaction, *USDCReaderTester, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCReaderTesterBin), backend, _version, _localDomain)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &USDCReaderTester{address: address, abi: *parsed, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

type USDCReaderTester struct {
	address common.Address
	abi     abi.ABI
	USDCReaderTesterCaller
	USDCReaderTesterTransactor
	USDCReaderTesterFilterer
}

type USDCReaderTesterCaller struct {
	contract *bind.BoundContract
}

type USDCReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type USDCReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type USDCReaderTesterSession struct {
	Contract     *USDCReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterCallerSession struct {
	Contract *USDCReaderTesterCaller
	CallOpts bind.CallOpts
}

type USDCReaderTesterTransactorSession struct {
	Contract     *USDCReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterRaw struct {
	Contract *USDCReaderTester
}

type USDCReaderTesterCallerRaw struct {
	Contract *USDCReaderTesterCaller
}

type USDCReaderTesterTransactorRaw struct {
	Contract *USDCReaderTesterTransactor
}

func NewUSDCReaderTester(address common.Address, backend bind.ContractBackend) (*USDCReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(USDCReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUSDCReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTester{address: address, abi: abi, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

func NewUSDCReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*USDCReaderTesterCaller, error) {
	contract, err := bindUSDCReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterCaller{contract: contract}, nil
}

func NewUSDCReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*USDCReaderTesterTransactor, error) {
	contract, err := bindUSDCReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterTransactor{contract: contract}, nil
}

func NewUSDCReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*USDCReaderTesterFilterer, error) {
	contract, err := bindUSDCReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterFilterer{contract: contract}, nil
}

func bindUSDCReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_USDCReaderTester *USDCReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.USDCReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterCaller) LocalDomain(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCReaderTester.contract.Call(opts, &out, "localDomain")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCReaderTester *USDCReaderTesterSession) LocalDomain() (uint32, error) {
	return _USDCReaderTester.Contract.LocalDomain(&_USDCReaderTester.CallOpts)
}

func (_USDCReaderTester *USDCReaderTesterCallerSession) LocalDomain() (uint32, error) {
	return _USDCReaderTester.Contract.LocalDomain(&_USDCReaderTester.CallOpts)
}

func (_USDCReaderTester *USDCReaderTesterCaller) Version(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCReaderTester.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCReaderTester *USDCReaderTesterSession) Version() (uint32, error) {
	return _USDCReaderTester.Contract.Version(&_USDCReaderTester.CallOpts)
}

func (_USDCReaderTester *USDCReaderTesterCallerSession) Version() (uint32, error) {
	return _USDCReaderTester.Contract.Version(&_USDCReaderTester.CallOpts)
}

func (_USDCReaderTester *USDCReaderTesterTransactor) EmitMessageSent(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.contract.Transact(opts, "emitMessageSent", destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterSession) EmitMessageSent(destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterTransactorSession) EmitMessageSent(destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

type USDCReaderTesterMessageSentIterator struct {
	Event *USDCReaderTesterMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCReaderTesterMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCReaderTesterMessageSent)
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
		it.Event = new(USDCReaderTesterMessageSent)
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

func (it *USDCReaderTesterMessageSentIterator) Error() error {
	return it.fail
}

func (it *USDCReaderTesterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCReaderTesterMessageSent struct {
	Arg0 []byte
	Raw  types.Log
}

func (_USDCReaderTester *USDCReaderTesterFilterer) FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error) {

	logs, sub, err := _USDCReaderTester.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterMessageSentIterator{contract: _USDCReaderTester.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

func (_USDCReaderTester *USDCReaderTesterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error) {

	logs, sub, err := _USDCReaderTester.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCReaderTesterMessageSent)
				if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

func (_USDCReaderTester *USDCReaderTesterFilterer) ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error) {
	event := new(USDCReaderTesterMessageSent)
	if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_USDCReaderTester *USDCReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _USDCReaderTester.abi.Events["MessageSent"].ID:
		return _USDCReaderTester.ParseMessageSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (USDCReaderTesterMessageSent) Topic() common.Hash {
	return common.HexToHash("0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036")
}

func (_USDCReaderTester *USDCReaderTester) Address() common.Address {
	return _USDCReaderTester.address
}

type USDCReaderTesterInterface interface {
	LocalDomain(opts *bind.CallOpts) (uint32, error)

	Version(opts *bind.CallOpts) (uint32, error)

	EmitMessageSent(opts *bind.TransactOpts, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error)

	FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error)

	WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error)

	ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
