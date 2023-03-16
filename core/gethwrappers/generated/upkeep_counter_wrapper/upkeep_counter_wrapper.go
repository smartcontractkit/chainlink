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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

var UpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161044d38038061044d8339818101604052604081101561003357600080fd5b508051602090910151600091825560015560038190554360025560048190556005556103e9806100646000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637f407edf11610076578063917d895f1161005b578063917d895f1461027b578063947a36fb14610283578063d832d92f1461028b576100be565b80637f407edf14610250578063806b984f14610273576100be565b806361bc221a116100a757806361bc221a1461014f5780636250a13a146101575780636e04ff0d1461015f576100be565b80632cb15864146100c35780634585e33b146100dd575b600080fd5b6100cb6102a7565b60408051918252519081900360200190f35b61014d600480360360208110156100f357600080fd5b81019060208101813564010000000081111561010e57600080fd5b82018360208201111561012057600080fd5b8035906020019184600183028401116401000000008311171561014257600080fd5b5090925090506102ad565b005b6100cb610324565b6100cb61032a565b6101cf6004803603602081101561017557600080fd5b81019060208101813564010000000081111561019057600080fd5b8201836020820111156101a257600080fd5b803590602001918460018302840111640100000000831117156101c457600080fd5b509092509050610330565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102145781810151838201526020016101fc565b50505050905090810190601f1680156102415780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61014d6004803603604081101561026657600080fd5b5080359060200135610382565b6100cb610394565b6100cb61039a565b6100cb6103a0565b6102936103a6565b604080519115158252519081900360200190f35b60045481565b6004546102b957436004555b4360028190556005805460010190819055600454600354604080519283526020830194909452818401526060810191909152905132917f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa919081900360800190a25050600254600355565b60055481565b60005481565b6000606061033c6103a6565b848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b60009182556001556004819055600555565b60025481565b60035481565b60015481565b6000600454600014156103bb575060016103d9565b60005460045443031080156103d65750600154600254430310155b90505b9056fea164736f6c6343000706000a",
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
	return address, tx, &UpkeepCounter{UpkeepCounterCaller: UpkeepCounterCaller{contract: contract}, UpkeepCounterTransactor: UpkeepCounterTransactor{contract: contract}, UpkeepCounterFilterer: UpkeepCounterFilterer{contract: contract}}, nil
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
	parsed, err := abi.JSON(strings.NewReader(UpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_UpkeepCounter *UpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialBlock(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.InitialBlock(&_UpkeepCounter.CallOpts)
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

func (_UpkeepCounter *UpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastBlock(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounter.Contract.LastBlock(&_UpkeepCounter.CallOpts)
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
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Raw           types.Log
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

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

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
