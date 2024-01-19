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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ONE_THOUSAND_BYTES_PADDING\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610be1380380610be183398101604081905261002f9161004d565b60009182556001556003819055436002556004819055600555610071565b6000806040838503121561006057600080fd5b505080516020909101519092909150565b610b61806100806000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80637f407edf11610081578063947a36fb1161005b578063947a36fb14610184578063c719ae101461018d578063d832d92f1461019557600080fd5b80637f407edf14610152578063806b984f14610172578063917d895f1461017b57600080fd5b806361bc221a116100b257806361bc221a1461011f5780636250a13a146101285780636e04ff0d1461013157600080fd5b80632cb15864146100d95780634585e33b146100f557806356ed5d7e1461010a575b600080fd5b6100e260045481565b6040519081526020015b60405180910390f35b6101086101033660046103c0565b6101ad565b005b61011261027f565b6040516100ec91906104a0565b6100e260055481565b6100e260005481565b61014461013f3660046103c0565b61029e565b6040516100ec9291906104ba565b6101086101603660046104dd565b60009182556001556004819055600555565b6100e260025481565b6100e260035481565b6100e260015481565b6101126102f0565b61019d61037e565b60405190151581526020016100ec565b6004546000036101bc57436004555b436002556005546101ce90600161052e565b60058190555081816040518061042001604052806103e8815260200161076d6103e8913960405160200161020493929190610547565b60405160208183030381529060405260069081610221919061063f565b5060045460025460035460055460408051948552602085019390935291830152606082015232907f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa9060800160405180910390a25050600254600355565b6040518061042001604052806103e8815260200161076d6103e8913981565b600060606102aa61037e565b848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959a92995091975050505050505050565b600680546102fd9061059d565b80601f01602080910402602001604051908101604052809291908181526020018280546103299061059d565b80156103765780601f1061034b57610100808354040283529160200191610376565b820191906000526020600020905b81548152906001019060200180831161035957829003601f168201915b505050505081565b60006004546000036103905750600190565b6000546004546103a09043610759565b1080156103bb57506001546002546103b89043610759565b10155b905090565b600080602083850312156103d357600080fd5b823567ffffffffffffffff808211156103eb57600080fd5b818501915085601f8301126103ff57600080fd5b81358181111561040e57600080fd5b86602082850101111561042057600080fd5b60209290920196919550909350505050565b60005b8381101561044d578181015183820152602001610435565b50506000910152565b6000815180845261046e816020860160208601610432565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006104b36020830184610456565b9392505050565b82151581526040602082015260006104d56040830184610456565b949350505050565b600080604083850312156104f057600080fd5b50508035926020909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115610541576105416104ff565b92915050565b828482376000838201600081528351610564818360208801610432565b0195945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600181811c908216806105b157607f821691505b6020821081036105ea577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561063a57600081815260208120601f850160051c810160208610156106175750805b601f850160051c820191505b8181101561063657828155600101610623565b5050505b505050565b815167ffffffffffffffff8111156106595761065961056e565b61066d81610667845461059d565b846105f0565b602080601f8311600181146106c0576000841561068a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610636565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561070d578886015182559484019460019091019084016106ee565b508582101561074957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115610541576105416104ff56feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000810000a",
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

func (_UpkeepCounter *UpkeepCounterCaller) ONETHOUSANDBYTESPADDING(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "ONE_THOUSAND_BYTES_PADDING")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) ONETHOUSANDBYTESPADDING() ([]byte, error) {
	return _UpkeepCounter.Contract.ONETHOUSANDBYTESPADDING(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) ONETHOUSANDBYTESPADDING() ([]byte, error) {
	return _UpkeepCounter.Contract.ONETHOUSANDBYTESPADDING(&_UpkeepCounter.CallOpts)
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

func (_UpkeepCounter *UpkeepCounterCaller) PerformData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _UpkeepCounter.contract.Call(opts, &out, "performData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_UpkeepCounter *UpkeepCounterSession) PerformData() ([]byte, error) {
	return _UpkeepCounter.Contract.PerformData(&_UpkeepCounter.CallOpts)
}

func (_UpkeepCounter *UpkeepCounterCallerSession) PerformData() ([]byte, error) {
	return _UpkeepCounter.Contract.PerformData(&_UpkeepCounter.CallOpts)
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

func (_UpkeepCounter *UpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, _performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.contract.Transact(opts, "performUpkeep", _performData)
}

func (_UpkeepCounter *UpkeepCounterSession) PerformUpkeep(_performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, _performData)
}

func (_UpkeepCounter *UpkeepCounterTransactorSession) PerformUpkeep(_performData []byte) (*types.Transaction, error) {
	return _UpkeepCounter.Contract.PerformUpkeep(&_UpkeepCounter.TransactOpts, _performData)
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
	ONETHOUSANDBYTESPADDING(opts *bind.CallOpts) ([]byte, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PerformData(opts *bind.CallOpts) ([]byte, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, _performData []byte) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
