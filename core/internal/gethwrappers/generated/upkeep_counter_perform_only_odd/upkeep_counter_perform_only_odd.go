// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_counter_perform_only_odd

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

var UpkeepCounterPerformOnlyOddMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516104d13803806104d18339818101604052604081101561003357600080fd5b5080516020909101516000918255600155600381905543600255600481905560055561046d806100646000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637f407edf11610076578063917d895f1161005b578063917d895f1461027b578063947a36fb14610283578063d832d92f1461028b576100be565b80637f407edf14610250578063806b984f14610273576100be565b806361bc221a116100a757806361bc221a1461014f5780636250a13a146101575780636e04ff0d1461015f576100be565b80632cb15864146100c35780634585e33b146100dd575b600080fd5b6100cb6102a7565b60408051918252519081900360200190f35b61014d600480360360208110156100f357600080fd5b81019060208101813564010000000081111561010e57600080fd5b82018360208201111561012057600080fd5b8035906020019184600183028401116401000000008311171561014257600080fd5b5090925090506102ad565b005b6100cb61037e565b6100cb610384565b6101cf6004803603602081101561017557600080fd5b81019060208101813564010000000081111561019057600080fd5b8201836020820111156101a257600080fd5b803590602001918460018302840111640100000000831117156101c457600080fd5b50909250905061038a565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102145781810151838201526020016101fc565b50505050905090810190601f1680156102415780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61014d6004803603604081101561026657600080fd5b50803590602001356103dc565b6100cb6103ee565b6100cb6103f4565b6100cb6103fa565b610293610400565b604080519115158252519081900360200190f35b60045481565b6002430615610307576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180610437602a913960400191505060405180910390fd5b60045461031357436004555b4360028190556005805460010190819055600454600354604080519283526020830194909452818401526060810191909152905132917f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa919081900360800190a25050600254600355565b60055481565b60005481565b60006060610396610400565b848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b60009182556001556004819055600555565b60025481565b60035481565b60015481565b60006004546000141561041557506001610433565b60005460045443031080156104305750600154600254430310155b90505b9056fe706572666f726d20646f65736e277420776f726b206f6e206576656e20626c6f636b206e756d62657273a164736f6c6343000706000a",
}

var UpkeepCounterPerformOnlyOddABI = UpkeepCounterPerformOnlyOddMetaData.ABI

var UpkeepCounterPerformOnlyOddBin = UpkeepCounterPerformOnlyOddMetaData.Bin

func DeployUpkeepCounterPerformOnlyOdd(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepCounterPerformOnlyOdd, error) {
	parsed, err := UpkeepCounterPerformOnlyOddMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepCounterPerformOnlyOddBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepCounterPerformOnlyOdd{UpkeepCounterPerformOnlyOddCaller: UpkeepCounterPerformOnlyOddCaller{contract: contract}, UpkeepCounterPerformOnlyOddTransactor: UpkeepCounterPerformOnlyOddTransactor{contract: contract}, UpkeepCounterPerformOnlyOddFilterer: UpkeepCounterPerformOnlyOddFilterer{contract: contract}}, nil
}

type UpkeepCounterPerformOnlyOdd struct {
	address common.Address
	abi     abi.ABI
	UpkeepCounterPerformOnlyOddCaller
	UpkeepCounterPerformOnlyOddTransactor
	UpkeepCounterPerformOnlyOddFilterer
}

type UpkeepCounterPerformOnlyOddCaller struct {
	contract *bind.BoundContract
}

type UpkeepCounterPerformOnlyOddTransactor struct {
	contract *bind.BoundContract
}

type UpkeepCounterPerformOnlyOddFilterer struct {
	contract *bind.BoundContract
}

type UpkeepCounterPerformOnlyOddSession struct {
	Contract     *UpkeepCounterPerformOnlyOdd
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepCounterPerformOnlyOddCallerSession struct {
	Contract *UpkeepCounterPerformOnlyOddCaller
	CallOpts bind.CallOpts
}

type UpkeepCounterPerformOnlyOddTransactorSession struct {
	Contract     *UpkeepCounterPerformOnlyOddTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepCounterPerformOnlyOddRaw struct {
	Contract *UpkeepCounterPerformOnlyOdd
}

type UpkeepCounterPerformOnlyOddCallerRaw struct {
	Contract *UpkeepCounterPerformOnlyOddCaller
}

type UpkeepCounterPerformOnlyOddTransactorRaw struct {
	Contract *UpkeepCounterPerformOnlyOddTransactor
}

func NewUpkeepCounterPerformOnlyOdd(address common.Address, backend bind.ContractBackend) (*UpkeepCounterPerformOnlyOdd, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepCounterPerformOnlyOddABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepCounterPerformOnlyOdd(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformOnlyOdd{address: address, abi: abi, UpkeepCounterPerformOnlyOddCaller: UpkeepCounterPerformOnlyOddCaller{contract: contract}, UpkeepCounterPerformOnlyOddTransactor: UpkeepCounterPerformOnlyOddTransactor{contract: contract}, UpkeepCounterPerformOnlyOddFilterer: UpkeepCounterPerformOnlyOddFilterer{contract: contract}}, nil
}

func NewUpkeepCounterPerformOnlyOddCaller(address common.Address, caller bind.ContractCaller) (*UpkeepCounterPerformOnlyOddCaller, error) {
	contract, err := bindUpkeepCounterPerformOnlyOdd(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformOnlyOddCaller{contract: contract}, nil
}

func NewUpkeepCounterPerformOnlyOddTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepCounterPerformOnlyOddTransactor, error) {
	contract, err := bindUpkeepCounterPerformOnlyOdd(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformOnlyOddTransactor{contract: contract}, nil
}

func NewUpkeepCounterPerformOnlyOddFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepCounterPerformOnlyOddFilterer, error) {
	contract, err := bindUpkeepCounterPerformOnlyOdd(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformOnlyOddFilterer{contract: contract}, nil
}

func bindUpkeepCounterPerformOnlyOdd(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepCounterPerformOnlyOddABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounterPerformOnlyOdd.Contract.UpkeepCounterPerformOnlyOddCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.UpkeepCounterPerformOnlyOddTransactor.contract.Transfer(opts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.UpkeepCounterPerformOnlyOddTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepCounterPerformOnlyOdd.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.contract.Transfer(opts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.CheckUpkeep(&_UpkeepCounterPerformOnlyOdd.CallOpts, data)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.CheckUpkeep(&_UpkeepCounterPerformOnlyOdd.CallOpts, data)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) Counter() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Counter(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) Counter() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Counter(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) Eligible() (bool, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Eligible(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) Eligible() (bool, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Eligible(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.InitialBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.InitialBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) Interval() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Interval(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) Interval() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.Interval(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.LastBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.LastBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.PreviousPerformBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.PreviousPerformBlock(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepCounterPerformOnlyOdd.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) TestRange() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.TestRange(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.TestRange(&_UpkeepCounterPerformOnlyOdd.CallOpts)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.contract.Transact(opts, "performUpkeep", performData)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.PerformUpkeep(&_UpkeepCounterPerformOnlyOdd.TransactOpts, performData)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.PerformUpkeep(&_UpkeepCounterPerformOnlyOdd.TransactOpts, performData)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactor) SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.contract.Transact(opts, "setSpread", _testRange, _interval)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.SetSpread(&_UpkeepCounterPerformOnlyOdd.TransactOpts, _testRange, _interval)
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddTransactorSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepCounterPerformOnlyOdd.Contract.SetSpread(&_UpkeepCounterPerformOnlyOdd.TransactOpts, _testRange, _interval)
}

type UpkeepCounterPerformOnlyOddPerformingUpkeepIterator struct {
	Event *UpkeepCounterPerformOnlyOddPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepCounterPerformOnlyOddPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepCounterPerformOnlyOddPerformingUpkeep)
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
		it.Event = new(UpkeepCounterPerformOnlyOddPerformingUpkeep)
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

func (it *UpkeepCounterPerformOnlyOddPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepCounterPerformOnlyOddPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepCounterPerformOnlyOddPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Raw           types.Log
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformOnlyOddPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounterPerformOnlyOdd.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepCounterPerformOnlyOddPerformingUpkeepIterator{contract: _UpkeepCounterPerformOnlyOdd.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformOnlyOddPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepCounterPerformOnlyOdd.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepCounterPerformOnlyOddPerformingUpkeep)
				if err := _UpkeepCounterPerformOnlyOdd.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOddFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformOnlyOddPerformingUpkeep, error) {
	event := new(UpkeepCounterPerformOnlyOddPerformingUpkeep)
	if err := _UpkeepCounterPerformOnlyOdd.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOdd) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepCounterPerformOnlyOdd.abi.Events["PerformingUpkeep"].ID:
		return _UpkeepCounterPerformOnlyOdd.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepCounterPerformOnlyOddPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa")
}

func (_UpkeepCounterPerformOnlyOdd *UpkeepCounterPerformOnlyOdd) Address() common.Address {
	return _UpkeepCounterPerformOnlyOdd.address
}

type UpkeepCounterPerformOnlyOddInterface interface {
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

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformOnlyOddPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformOnlyOddPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformOnlyOddPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
