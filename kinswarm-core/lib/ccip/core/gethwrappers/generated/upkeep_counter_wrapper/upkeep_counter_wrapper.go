// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_counter_wrapper

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

var UpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161048338038061048383398101604081905261002f9161004d565b60009182556001556003819055426002556004819055600555610071565b6000806040838503121561006057600080fd5b505080516020909101519092909150565b610403806100806000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637f407edf11610076578063947a36fb1161005b578063947a36fb14610150578063d6d1417114610159578063d832d92f1461016257600080fd5b80637f407edf14610127578063917d895f1461014757600080fd5b806361bc221a116100a757806361bc221a146100f45780636250a13a146100fd5780636e04ff0d1461010657600080fd5b806319d8ac61146100c35780634585e33b146100df575b600080fd5b6100cc60025481565b6040519081526020015b60405180910390f35b6100f26100ed366004610291565b61017a565b005b6100cc60055481565b6100cc60005481565b610119610114366004610291565b6101fd565b6040516100d6929190610303565b6100f2610135366004610379565b60009182556001556004819055600555565b6100cc60035481565b6100cc60015481565b6100cc60045481565b61016a61024f565b60405190151581526020016100d6565b60045460000361018957426004555b4260025560055461019b9060016103ca565b600581905560045460025460035460408051938452602084019290925290820152606081019190915232907f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa9060800160405180910390a25050600254600355565b6000606061020961024f565b848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b60006004546000036102615750600190565b60005460045461027190426103e3565b10801561028c575060015460025461028990426103e3565b10155b905090565b600080602083850312156102a457600080fd5b823567ffffffffffffffff808211156102bc57600080fd5b818501915085601f8301126102d057600080fd5b8135818111156102df57600080fd5b8660208285010111156102f157600080fd5b60209290920196919550909350505050565b821515815260006020604081840152835180604085015260005b818110156103395785810183015185820160600152820161031d565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b6000806040838503121561038c57600080fd5b50508035926020909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156103dd576103dd61039b565b92915050565b818103818111156103dd576103dd61039b56fea164736f6c6343000810000a",
}

var UpkeepCounterABI = UpkeepCounterMetaData.ABI

var UpkeepCounterBin = UpkeepCounterMetaData.Bin

func DeployUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepCounter, error) {
	parsed, err := UpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepCounterBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepCounter{address: address, abi: *parsed, UpkeepCounterCaller: UpkeepCounterCaller{contract: contract}, UpkeepCounterTransactor: UpkeepCounterTransactor{contract: contract}, UpkeepCounterFilterer: UpkeepCounterFilterer{contract: contract}}, nil
}

type UpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	UpkeepCounterCaller
	UpkeepCounterTransactor
	UpkeepCounterFilterer
}

type UpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type UpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type UpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type UpkeepCounterSession struct {
	Contract     *UpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepCounterCallerSession struct {
	Contract *UpkeepCounterCaller
	CallOpts bind.CallOpts
}

type UpkeepCounterTransactorSession struct {
	Contract     *UpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepCounterRaw struct {
	Contract *UpkeepCounter
}

type UpkeepCounterCallerRaw struct {
	Contract *UpkeepCounterCaller
}

type UpkeepCounterTransactorRaw struct {
	Contract *UpkeepCounterTransactor
}

func NewUpkeepCounter(address common.Address, backend bind.ContractBackend) (*UpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounter{address: address, abi: abi, UpkeepCounterCaller: UpkeepCounterCaller{contract: contract}, UpkeepCounterTransactor: UpkeepCounterTransactor{contract: contract}, UpkeepCounterFilterer: UpkeepCounterFilterer{contract: contract}}, nil
}

func NewUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*UpkeepCounterCaller, error) {
	contract, err := bindUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterCaller{contract: contract}, nil
}

func NewUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepCounterTransactor, error) {
	contract, err := bindUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterTransactor{contract: contract}, nil
}

func NewUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepCounterFilterer, error) {
	contract, err := bindUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterFilterer{contract: contract}, nil
}

func bindUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_UpkeepCounter *UpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounter.Contract.UpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepCounter *UpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.UpkeepCounterTransactor.contract.Transfer(opts)
}

func (_UpkeepCounter *UpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.UpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepCounter *UpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepCounter *UpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.contract.Transfer(opts)
}

func (_UpkeepCounter *UpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepCounter *UpkeepCounterCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepCounter *UpkeepCounterSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounter.Contract.CheckUpkeep(&_UpkeepCounter.CallOpts, data)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounter.Contract.CheckUpkeep(&_UpkeepCounter.CallOpts, data)
}

func (_UpkeepCounter *UpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) Counter() (*big.Int, error) {
	return _UpkeepCounter.Contract.Counter(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _UpkeepCounter.Contract.Counter(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) Eligible() (bool, error) {
	return _UpkeepCounter.Contract.Eligible(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) Eligible() (bool, error) {
	return _UpkeepCounter.Contract.Eligible(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) InitialTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "initialTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) InitialTimestamp() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialTimestamp(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) InitialTimestamp() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialTimestamp(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) Interval() (*big.Int, error) {
	return _UpkeepCounter.Contract.Interval(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) Interval() (*big.Int, error) {
	return _UpkeepCounter.Contract.Interval(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) LastTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "lastTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) LastTimestamp() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastTimestamp(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) LastTimestamp() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastTimestamp(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.PreviousPerformBlock(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.PreviousPerformBlock(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) TestRange() (*big.Int, error) {
	return _UpkeepCounter.Contract.TestRange(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepCounter.Contract.TestRange(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_UpkeepCounter *UpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, performData)
}

func (_UpkeepCounter *UpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, performData)
}

func (_UpkeepCounter *UpkeepCounterTransactor) SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.contract.Transact(opts, "setSpread", _testRange, _interval)
}

func (_UpkeepCounter *UpkeepCounterSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.SetSpread(&_UpkeepCounter.TransactOpts, _testRange, _interval)
}

func (_UpkeepCounter *UpkeepCounterTransactorSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.SetSpread(&_UpkeepCounter.TransactOpts, _testRange, _interval)
}

type UpkeepCounterPerformingUpkeepIterator struct {
	Event *UpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepCounterPerformingUpkeep)
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
		it.Event = new(UpkeepCounterPerformingUpkeep)
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

func (it *UpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepCounterPerformingUpkeep struct {
	From             common.Address
	InitialTimestamp *big.Int
	LastTimestamp    *big.Int
	PreviousBlock    *big.Int
	Counter          *big.Int
	Raw              types.Log
}

func (_UpkeepCounter *UpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformingUpkeepIterator{contract: _UpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepCounter *UpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepCounterPerformingUpkeep)
				if err := _UpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_UpkeepCounter *UpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformingUpkeep, error) {
	event := new(UpkeepCounterPerformingUpkeep)
	if err := _UpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepCounter *UpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _UpkeepCounter.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa")
}

func (_UpkeepCounter *UpkeepCounter) Address() common.Address {
	return _UpkeepCounter.address
}

type UpkeepCounterInterface interface {
	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	InitialTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastTimestamp(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
