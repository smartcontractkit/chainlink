// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_consumer_benchmark

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

var AutomationConsumerBenchmarkMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"range\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"firstEligibleBuffer\",\"type\":\"uint256\"}],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50436004556106ed806100246000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80637145f11b11610076578063d826f88f1161005b578063d826f88f14610177578063e81018b314610180578063f597b393146101a057600080fd5b80637145f11b14610134578063a6fe0a1e1461015757600080fd5b80633b3546c8146100a85780634585e33b146100db57806351dcee4b146100f05780636e04ff0d14610113575b600080fd5b6100c86100b63660046104dd565b60036020526000908152604090205481565b6040519081526020015b60405180910390f35b6100ee6100e93660046104f6565b6101c0565b005b6101036100fe366004610568565b610341565b60405190151581526020016100d2565b6101266101213660046104f6565b610356565b6040516100d2929190610594565b6101036101423660046104dd565b60026020526000908152604090205460ff1681565b6100c86101653660046104dd565b60016020526000908152604090205481565b6100ee43600455565b6100c861018e3660046104dd565b60009081526003602052604090205490565b6100c86101ae3660046104dd565b60006020819052908152604090205481565b600080808080806101d38789018961060a565b9550955095509550955095506101ea868583610477565b6101f357600080fd5b60005a6000888152602081905260408120549192500361021f5760008781526020819052604090204390555b610229864361067c565b6000888152600160209081526040808320939093556003905290812080549161025183610695565b90915550506000878152602081815260408083205460018352928190205481518b815232938101939093529082019290925260608101919091524360808201527f39223708d1655effd0be3f9a99a7e1d1aadd9fb456f0bfc4c2a4f50b2484a3679060a00160405180910390a160006102cb6001436106cd565b40905060005b845a6102dd90856106cd565b1015610334578080156102fe575060008281526002602052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506102d1565b5050505050505050505050565b600061034e848484610477565b949350505050565b6000606081808080808061036c898b018b61060a565b95509550955095509550955060005a9050600061038a6001436106cd565b409050600080861180156103a457506103a4898886610477565b1561040d575b855a6103b690856106cd565b101561040d578080156103d7575060008281526002602052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506103aa565b610418898886610477565b8d8d81818080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050905090509a509a505050505050505050509250929050565b600083815260208190526040812054156104c55760008481526020819052604090205483906104a690436106cd565b1080156104c0575060008481526001602052604090205443115b61034e565b6004546104d2908361067c565b431015949350505050565b6000602082840312156104ef57600080fd5b5035919050565b6000806020838503121561050957600080fd5b823567ffffffffffffffff8082111561052157600080fd5b818501915085601f83011261053557600080fd5b81358181111561054457600080fd5b86602082850101111561055657600080fd5b60209290920196919550909350505050565b60008060006060848603121561057d57600080fd5b505081359360208301359350604090920135919050565b821515815260006020604081840152835180604085015260005b818110156105ca578581018301518582016060015282016105ae565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b60008060008060008060c0878903121561062357600080fd5b505084359660208601359650604086013595606081013595506080810135945060a0013592509050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561068f5761068f61064d565b92915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036106c6576106c661064d565b5060010190565b8181038181111561068f5761068f61064d56fea164736f6c6343000810000a",
}

var AutomationConsumerBenchmarkABI = AutomationConsumerBenchmarkMetaData.ABI

var AutomationConsumerBenchmarkBin = AutomationConsumerBenchmarkMetaData.Bin

func DeployAutomationConsumerBenchmark(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AutomationConsumerBenchmark, error) {
	parsed, err := AutomationConsumerBenchmarkMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationConsumerBenchmarkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationConsumerBenchmark{address: address, abi: *parsed, AutomationConsumerBenchmarkCaller: AutomationConsumerBenchmarkCaller{contract: contract}, AutomationConsumerBenchmarkTransactor: AutomationConsumerBenchmarkTransactor{contract: contract}, AutomationConsumerBenchmarkFilterer: AutomationConsumerBenchmarkFilterer{contract: contract}}, nil
}

type AutomationConsumerBenchmark struct {
	address common.Address
	abi     abi.ABI
	AutomationConsumerBenchmarkCaller
	AutomationConsumerBenchmarkTransactor
	AutomationConsumerBenchmarkFilterer
}

type AutomationConsumerBenchmarkCaller struct {
	contract *bind.BoundContract
}

type AutomationConsumerBenchmarkTransactor struct {
	contract *bind.BoundContract
}

type AutomationConsumerBenchmarkFilterer struct {
	contract *bind.BoundContract
}

type AutomationConsumerBenchmarkSession struct {
	Contract     *AutomationConsumerBenchmark
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationConsumerBenchmarkCallerSession struct {
	Contract *AutomationConsumerBenchmarkCaller
	CallOpts bind.CallOpts
}

type AutomationConsumerBenchmarkTransactorSession struct {
	Contract     *AutomationConsumerBenchmarkTransactor
	TransactOpts bind.TransactOpts
}

type AutomationConsumerBenchmarkRaw struct {
	Contract *AutomationConsumerBenchmark
}

type AutomationConsumerBenchmarkCallerRaw struct {
	Contract *AutomationConsumerBenchmarkCaller
}

type AutomationConsumerBenchmarkTransactorRaw struct {
	Contract *AutomationConsumerBenchmarkTransactor
}

func NewAutomationConsumerBenchmark(address common.Address, backend bind.ContractBackend) (*AutomationConsumerBenchmark, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationConsumerBenchmarkABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationConsumerBenchmark(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationConsumerBenchmark{address: address, abi: abi, AutomationConsumerBenchmarkCaller: AutomationConsumerBenchmarkCaller{contract: contract}, AutomationConsumerBenchmarkTransactor: AutomationConsumerBenchmarkTransactor{contract: contract}, AutomationConsumerBenchmarkFilterer: AutomationConsumerBenchmarkFilterer{contract: contract}}, nil
}

func NewAutomationConsumerBenchmarkCaller(address common.Address, caller bind.ContractCaller) (*AutomationConsumerBenchmarkCaller, error) {
	contract, err := bindAutomationConsumerBenchmark(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationConsumerBenchmarkCaller{contract: contract}, nil
}

func NewAutomationConsumerBenchmarkTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationConsumerBenchmarkTransactor, error) {
	contract, err := bindAutomationConsumerBenchmark(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationConsumerBenchmarkTransactor{contract: contract}, nil
}

func NewAutomationConsumerBenchmarkFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationConsumerBenchmarkFilterer, error) {
	contract, err := bindAutomationConsumerBenchmark(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationConsumerBenchmarkFilterer{contract: contract}, nil
}

func bindAutomationConsumerBenchmark(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationConsumerBenchmarkMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationConsumerBenchmark.Contract.AutomationConsumerBenchmarkCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.AutomationConsumerBenchmarkTransactor.contract.Transfer(opts)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.AutomationConsumerBenchmarkTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationConsumerBenchmark.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.contract.Transfer(opts)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) CheckEligible(opts *bind.CallOpts, id *big.Int, arg1 *big.Int, firstEligibleBuffer *big.Int) (bool, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "checkEligible", id, arg1, firstEligibleBuffer)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) CheckEligible(id *big.Int, arg1 *big.Int, firstEligibleBuffer *big.Int) (bool, error) {
	return _AutomationConsumerBenchmark.Contract.CheckEligible(&_AutomationConsumerBenchmark.CallOpts, id, arg1, firstEligibleBuffer)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) CheckEligible(id *big.Int, arg1 *big.Int, firstEligibleBuffer *big.Int) (bool, error) {
	return _AutomationConsumerBenchmark.Contract.CheckEligible(&_AutomationConsumerBenchmark.CallOpts, id, arg1, firstEligibleBuffer)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) CheckUpkeep(opts *bind.CallOpts, checkData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "checkUpkeep", checkData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) CheckUpkeep(checkData []byte) (bool, []byte, error) {
	return _AutomationConsumerBenchmark.Contract.CheckUpkeep(&_AutomationConsumerBenchmark.CallOpts, checkData)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) CheckUpkeep(checkData []byte) (bool, []byte, error) {
	return _AutomationConsumerBenchmark.Contract.CheckUpkeep(&_AutomationConsumerBenchmark.CallOpts, checkData)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) Count(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "count", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) Count(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.Count(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) Count(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.Count(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _AutomationConsumerBenchmark.Contract.DummyMap(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _AutomationConsumerBenchmark.Contract.DummyMap(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) GetCountPerforms(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "getCountPerforms", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) GetCountPerforms(id *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.GetCountPerforms(&_AutomationConsumerBenchmark.CallOpts, id)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) GetCountPerforms(id *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.GetCountPerforms(&_AutomationConsumerBenchmark.CallOpts, id)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) InitialCall(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "initialCall", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) InitialCall(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.InitialCall(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) InitialCall(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.InitialCall(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCaller) NextEligible(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationConsumerBenchmark.contract.Call(opts, &out, "nextEligible", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) NextEligible(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.NextEligible(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkCallerSession) NextEligible(arg0 *big.Int) (*big.Int, error) {
	return _AutomationConsumerBenchmark.Contract.NextEligible(&_AutomationConsumerBenchmark.CallOpts, arg0)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.contract.Transact(opts, "performUpkeep", performData)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.PerformUpkeep(&_AutomationConsumerBenchmark.TransactOpts, performData)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.PerformUpkeep(&_AutomationConsumerBenchmark.TransactOpts, performData)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.contract.Transact(opts, "reset")
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkSession) Reset() (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.Reset(&_AutomationConsumerBenchmark.TransactOpts)
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkTransactorSession) Reset() (*types.Transaction, error) {
	return _AutomationConsumerBenchmark.Contract.Reset(&_AutomationConsumerBenchmark.TransactOpts)
}

type AutomationConsumerBenchmarkPerformingUpkeepIterator struct {
	Event *AutomationConsumerBenchmarkPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationConsumerBenchmarkPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationConsumerBenchmarkPerformingUpkeep)
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
		it.Event = new(AutomationConsumerBenchmarkPerformingUpkeep)
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

func (it *AutomationConsumerBenchmarkPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *AutomationConsumerBenchmarkPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationConsumerBenchmarkPerformingUpkeep struct {
	Id           *big.Int
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*AutomationConsumerBenchmarkPerformingUpkeepIterator, error) {

	logs, sub, err := _AutomationConsumerBenchmark.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &AutomationConsumerBenchmarkPerformingUpkeepIterator{contract: _AutomationConsumerBenchmark.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *AutomationConsumerBenchmarkPerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _AutomationConsumerBenchmark.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationConsumerBenchmarkPerformingUpkeep)
				if err := _AutomationConsumerBenchmark.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmarkFilterer) ParsePerformingUpkeep(log types.Log) (*AutomationConsumerBenchmarkPerformingUpkeep, error) {
	event := new(AutomationConsumerBenchmarkPerformingUpkeep)
	if err := _AutomationConsumerBenchmark.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmark) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationConsumerBenchmark.abi.Events["PerformingUpkeep"].ID:
		return _AutomationConsumerBenchmark.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationConsumerBenchmarkPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x39223708d1655effd0be3f9a99a7e1d1aadd9fb456f0bfc4c2a4f50b2484a367")
}

func (_AutomationConsumerBenchmark *AutomationConsumerBenchmark) Address() common.Address {
	return _AutomationConsumerBenchmark.address
}

type AutomationConsumerBenchmarkInterface interface {
	CheckEligible(opts *bind.CallOpts, id *big.Int, arg1 *big.Int, firstEligibleBuffer *big.Int) (bool, error)

	CheckUpkeep(opts *bind.CallOpts, checkData []byte) (bool, []byte, error)

	Count(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	GetCountPerforms(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	InitialCall(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	NextEligible(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*AutomationConsumerBenchmarkPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *AutomationConsumerBenchmarkPerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*AutomationConsumerBenchmarkPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
