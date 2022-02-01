// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_perform_counter_restrictive_wrapper

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

var UpkeepPerformCounterRestrictiveMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"eligible\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060045534801561001e57600080fd5b5060405161049b38038061049b8339818101604052604081101561004157600080fd5b508051602090910151600291909155600355610439806100626000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063926f086e11610076578063c228a98e1161005b578063c228a98e1461027b578063d826f88f14610297578063e303666f1461029f576100be565b8063926f086e1461026b578063a9a4c57c14610273576100be565b80636250a13a116100a75780636250a13a1461014f5780636e04ff0d146101575780637f407edf14610248576100be565b80634585e33b146100c3578063523d9b8a14610135575b600080fd5b610133600480360360208110156100d957600080fd5b8101906020810181356401000000008111156100f457600080fd5b82018360208201111561010657600080fd5b8035906020019184600183028401116401000000008311171561012857600080fd5b5090925090506102a7565b005b61013d610350565b60408051918252519081900360200190f35b61013d610356565b6101c76004803603602081101561016d57600080fd5b81019060208101813564010000000081111561018857600080fd5b82018360208201111561019a57600080fd5b803590602001918460018302840111640100000000831117156101bc57600080fd5b50909250905061035c565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561020c5781810151838201526020016101f4565b50505050905090810190601f1680156102395780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101336004803603604081101561025e57600080fd5b5080359060200135610383565b61013d61038e565b61013d610394565b61028361039a565b604080519115158252519081900360200190f35b6101336103a9565b61013d6103b3565b60006102b16103b9565b60005460015460408051841515815232602082015280820193909352606083019190915243608083018190529051929350917fbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc09181900360a00190a18161031757600080fd5b6000546103245760008190555b6003546002026103326103dd565b8161033957fe5b060160019081018155600480549091019055505050565b60015481565b60025481565b600060606103686103b9565b60405180602001604052806000815250915091509250929050565b600291909155600355565b60005481565b60035481565b60006103a46103b9565b905090565b6000808055600455565b60045490565b6000805415806103a4575060025460005443031080156103a4575050600154431190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff4301406020808301919091523082840152825180830384018152606090920190925280519101209056fea164736f6c6343000706000a",
}

var UpkeepPerformCounterRestrictiveABI = UpkeepPerformCounterRestrictiveMetaData.ABI

var UpkeepPerformCounterRestrictiveBin = UpkeepPerformCounterRestrictiveMetaData.Bin

func DeployUpkeepPerformCounterRestrictive(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _averageEligibilityCadence *big.Int) (common.Address, *types.Transaction, *UpkeepPerformCounterRestrictive, error) {
	parsed, err := UpkeepPerformCounterRestrictiveMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepPerformCounterRestrictiveBin), backend, _testRange, _averageEligibilityCadence)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepPerformCounterRestrictive{UpkeepPerformCounterRestrictiveCaller: UpkeepPerformCounterRestrictiveCaller{contract: contract}, UpkeepPerformCounterRestrictiveTransactor: UpkeepPerformCounterRestrictiveTransactor{contract: contract}, UpkeepPerformCounterRestrictiveFilterer: UpkeepPerformCounterRestrictiveFilterer{contract: contract}}, nil
}

type UpkeepPerformCounterRestrictive struct {
	address common.Address
	abi     abi.ABI
	UpkeepPerformCounterRestrictiveCaller
	UpkeepPerformCounterRestrictiveTransactor
	UpkeepPerformCounterRestrictiveFilterer
}

type UpkeepPerformCounterRestrictiveCaller struct {
	contract *bind.BoundContract
}

type UpkeepPerformCounterRestrictiveTransactor struct {
	contract *bind.BoundContract
}

type UpkeepPerformCounterRestrictiveFilterer struct {
	contract *bind.BoundContract
}

type UpkeepPerformCounterRestrictiveSession struct {
	Contract     *UpkeepPerformCounterRestrictive
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepPerformCounterRestrictiveCallerSession struct {
	Contract *UpkeepPerformCounterRestrictiveCaller
	CallOpts bind.CallOpts
}

type UpkeepPerformCounterRestrictiveTransactorSession struct {
	Contract     *UpkeepPerformCounterRestrictiveTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepPerformCounterRestrictiveRaw struct {
	Contract *UpkeepPerformCounterRestrictive
}

type UpkeepPerformCounterRestrictiveCallerRaw struct {
	Contract *UpkeepPerformCounterRestrictiveCaller
}

type UpkeepPerformCounterRestrictiveTransactorRaw struct {
	Contract *UpkeepPerformCounterRestrictiveTransactor
}

func NewUpkeepPerformCounterRestrictive(address common.Address, backend bind.ContractBackend) (*UpkeepPerformCounterRestrictive, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepPerformCounterRestrictiveABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepPerformCounterRestrictive(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictive{address: address, abi: abi, UpkeepPerformCounterRestrictiveCaller: UpkeepPerformCounterRestrictiveCaller{contract: contract}, UpkeepPerformCounterRestrictiveTransactor: UpkeepPerformCounterRestrictiveTransactor{contract: contract}, UpkeepPerformCounterRestrictiveFilterer: UpkeepPerformCounterRestrictiveFilterer{contract: contract}}, nil
}

func NewUpkeepPerformCounterRestrictiveCaller(address common.Address, caller bind.ContractCaller) (*UpkeepPerformCounterRestrictiveCaller, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveCaller{contract: contract}, nil
}

func NewUpkeepPerformCounterRestrictiveTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepPerformCounterRestrictiveTransactor, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveTransactor{contract: contract}, nil
}

func NewUpkeepPerformCounterRestrictiveFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepPerformCounterRestrictiveFilterer, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveFilterer{contract: contract}, nil
}

func bindUpkeepPerformCounterRestrictive(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepPerformCounterRestrictiveABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveTransactor.contract.Transfer(opts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Transfer(opts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "averageEligibilityCadence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) AverageEligibilityCadence() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.AverageEligibilityCadence(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) AverageEligibilityCadence() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.AverageEligibilityCadence(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) CheckEligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "checkEligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckEligible() (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckEligible() (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckUpkeep(&_UpkeepPerformCounterRestrictive.CallOpts, data)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckUpkeep(&_UpkeepPerformCounterRestrictive.CallOpts, data)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) GetCountPerforms(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "getCountPerforms")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) GetCountPerforms() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.GetCountPerforms(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) GetCountPerforms() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.GetCountPerforms(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) InitialCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "initialCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) InitialCall() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.InitialCall(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) InitialCall() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.InitialCall(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) NextEligible(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "nextEligible")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) NextEligible() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.NextEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) NextEligible() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.NextEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) TestRange() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.TestRange(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.TestRange(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "performUpkeep", data)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, data)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, data)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "reset")
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) Reset() (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.Reset(&_UpkeepPerformCounterRestrictive.TransactOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) Reset() (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.Reset(&_UpkeepPerformCounterRestrictive.TransactOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setSpread", _newTestRange, _newAverageEligibilityCadence)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetSpread(&_UpkeepPerformCounterRestrictive.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetSpread(&_UpkeepPerformCounterRestrictive.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

type UpkeepPerformCounterRestrictivePerformingUpkeepIterator struct {
	Event *UpkeepPerformCounterRestrictivePerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepPerformCounterRestrictivePerformingUpkeep)
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
		it.Event = new(UpkeepPerformCounterRestrictivePerformingUpkeep)
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

func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepPerformCounterRestrictivePerformingUpkeep struct {
	Eligible     bool
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*UpkeepPerformCounterRestrictivePerformingUpkeepIterator, error) {

	logs, sub, err := _UpkeepPerformCounterRestrictive.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictivePerformingUpkeepIterator{contract: _UpkeepPerformCounterRestrictive.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepPerformCounterRestrictivePerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _UpkeepPerformCounterRestrictive.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepPerformCounterRestrictivePerformingUpkeep)
				if err := _UpkeepPerformCounterRestrictive.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepPerformCounterRestrictivePerformingUpkeep, error) {
	event := new(UpkeepPerformCounterRestrictivePerformingUpkeep)
	if err := _UpkeepPerformCounterRestrictive.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictive) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepPerformCounterRestrictive.abi.Events["PerformingUpkeep"].ID:
		return _UpkeepPerformCounterRestrictive.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepPerformCounterRestrictivePerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0")
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictive) Address() common.Address {
	return _UpkeepPerformCounterRestrictive.address
}

type UpkeepPerformCounterRestrictiveInterface interface {
	AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error)

	CheckEligible(opts *bind.CallOpts) (bool, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	GetCountPerforms(opts *bind.CallOpts) (*big.Int, error)

	InitialCall(opts *bind.CallOpts) (*big.Int, error)

	NextEligible(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*UpkeepPerformCounterRestrictivePerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepPerformCounterRestrictivePerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepPerformCounterRestrictivePerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
