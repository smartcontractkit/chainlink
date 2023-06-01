// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package log_upkeep_counter_wrapper

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

var LogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"logData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610a76380380610a7683398101604081905261002f9161004a565b60009081556002819055436001556003819055600455610063565b60006020828403121561005c57600080fd5b5051919050565b610a04806100726000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063806b984f11610076578063b66a261c1161005b578063b66a261c14610139578063be9a655514610156578063d832d92f1461015e57600080fd5b8063806b984f14610127578063917d895f1461013057600080fd5b806361bc221a116100a757806361bc221a146100f45780636250a13a146100fd5780636e04ff0d1461010657600080fd5b80632cb15864146100c35780634585e33b146100df575b600080fd5b6100cc60035481565b6040519081526020015b60405180910390f35b6100f26100ed366004610867565b610176565b005b6100cc60045481565b6100cc60005481565b610119610114366004610867565b610485565b6040516100d69291906108f2565b6100cc60015481565b6100cc60025481565b6100f26101473660046108d9565b60009081556003819055600455565b6100f261071b565b6101666107f5565b60405190151581526020016100d6565b60035461018257436003555b43600190815560045461019491610999565b60049081556001546002556000906101ae9082848661096f565b8101906101bb919061081e565b90507f3d53a395000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000082161415610235576040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a161042d565b7f57b1de35000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000821614156102b957604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da906020015b60405180910390a161042d565b7f1da9f70f000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216141561033a576040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c891016102ac565b7f5121119b000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000821614156103c6576040805160018152600260208201526003918101919091527f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e71906060016102ac565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f636f756c64206e6f742066696e64206d61746368696e6720736967000000000060448201526064015b60405180910390fd5b60035460015460025460045460408051948552602085019390935291830152606082015232907f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa9060800160405180910390a2505050565b600060606104916107f5565b6104f7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f6e6f7420656c696769626c6500000000000000000000000000000000000000006044820152606401610424565b6000610506600482868861096f565b810190610513919061081e565b90507f3d53a395000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000821614806105a657507f57b1de35000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216145b806105f257507f1da9f70f000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216145b8061063e57507f5121119b000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216145b1561068c576001858581818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525095985091965061071495505050505050565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401610424565b9250929050565b6040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a1604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da9060200160405180910390a16040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c8910160405180910390a160408051600181526002602082015260038183015290517f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e719181900360600190a1565b6000600354600014156108085750600190565b60005460035461081890436109b1565b10905090565b60006020828403121561083057600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461086057600080fd5b9392505050565b6000806020838503121561087a57600080fd5b823567ffffffffffffffff8082111561089257600080fd5b818501915085601f8301126108a657600080fd5b8135818111156108b557600080fd5b8660208285010111156108c757600080fd5b60209290920196919550909350505050565b6000602082840312156108eb57600080fd5b5035919050565b821515815260006020604081840152835180604085015260005b818110156109285785810183015185820160600152820161090c565b8181111561093a576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b6000808585111561097f57600080fd5b8386111561098c57600080fd5b5050820193919092039150565b600082198211156109ac576109ac6109c8565b500190565b6000828210156109c3576109c36109c8565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000806000a",
}

var LogUpkeepCounterABI = LogUpkeepCounterMetaData.ABI

var LogUpkeepCounterBin = LogUpkeepCounterMetaData.Bin

func DeployLogUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int) (common.Address, *types.Transaction, *LogUpkeepCounter, error) {
	parsed, err := LogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogUpkeepCounterBin), backend, _testRange)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LogUpkeepCounter{LogUpkeepCounterCaller: LogUpkeepCounterCaller{contract: contract}, LogUpkeepCounterTransactor: LogUpkeepCounterTransactor{contract: contract}, LogUpkeepCounterFilterer: LogUpkeepCounterFilterer{contract: contract}}, nil
}

type LogUpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	LogUpkeepCounterCaller
	LogUpkeepCounterTransactor
	LogUpkeepCounterFilterer
}

type LogUpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type LogUpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type LogUpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type LogUpkeepCounterSession struct {
	Contract     *LogUpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LogUpkeepCounterCallerSession struct {
	Contract *LogUpkeepCounterCaller
	CallOpts bind.CallOpts
}

type LogUpkeepCounterTransactorSession struct {
	Contract     *LogUpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type LogUpkeepCounterRaw struct {
	Contract *LogUpkeepCounter
}

type LogUpkeepCounterCallerRaw struct {
	Contract *LogUpkeepCounterCaller
}

type LogUpkeepCounterTransactorRaw struct {
	Contract *LogUpkeepCounterTransactor
}

func NewLogUpkeepCounter(address common.Address, backend bind.ContractBackend) (*LogUpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(LogUpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLogUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounter{address: address, abi: abi, LogUpkeepCounterCaller: LogUpkeepCounterCaller{contract: contract}, LogUpkeepCounterTransactor: LogUpkeepCounterTransactor{contract: contract}, LogUpkeepCounterFilterer: LogUpkeepCounterFilterer{contract: contract}}, nil
}

func NewLogUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*LogUpkeepCounterCaller, error) {
	contract, err := bindLogUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterCaller{contract: contract}, nil
}

func NewLogUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*LogUpkeepCounterTransactor, error) {
	contract, err := bindLogUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterTransactor{contract: contract}, nil
}

func NewLogUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*LogUpkeepCounterFilterer, error) {
	contract, err := bindLogUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterFilterer{contract: contract}, nil
}

func bindLogUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LogUpkeepCounter *LogUpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogUpkeepCounter.Contract.LogUpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_LogUpkeepCounter *LogUpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.LogUpkeepCounterTransactor.contract.Transfer(opts)
}

func (_LogUpkeepCounter *LogUpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.LogUpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogUpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.contract.Transfer(opts)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) CheckUpkeep(opts *bind.CallOpts, logData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "checkUpkeep", logData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) CheckUpkeep(logData []byte) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckUpkeep(&_LogUpkeepCounter.CallOpts, logData)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) CheckUpkeep(logData []byte) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckUpkeep(&_LogUpkeepCounter.CallOpts, logData)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) Counter() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.Counter(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.Counter(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) Eligible() (bool, error) {
	return _LogUpkeepCounter.Contract.Eligible(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) Eligible() (bool, error) {
	return _LogUpkeepCounter.Contract.Eligible(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.InitialBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.InitialBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.LastBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.LastBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.PreviousPerformBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.PreviousPerformBlock(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) TestRange() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.TestRange(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) TestRange() (*big.Int, error) {
	return _LogUpkeepCounter.Contract.TestRange(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _LogUpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_LogUpkeepCounter *LogUpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.PerformUpkeep(&_LogUpkeepCounter.TransactOpts, performData)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.PerformUpkeep(&_LogUpkeepCounter.TransactOpts, performData)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactor) SetSpread(opts *bind.TransactOpts, _testRange *big.Int) (*types.Transaction, error) {
	return _LogUpkeepCounter.contract.Transact(opts, "setSpread", _testRange)
}

func (_LogUpkeepCounter *LogUpkeepCounterSession) SetSpread(_testRange *big.Int) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.SetSpread(&_LogUpkeepCounter.TransactOpts, _testRange)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorSession) SetSpread(_testRange *big.Int) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.SetSpread(&_LogUpkeepCounter.TransactOpts, _testRange)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactor) Start(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogUpkeepCounter.contract.Transact(opts, "start")
}

func (_LogUpkeepCounter *LogUpkeepCounterSession) Start() (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.Start(&_LogUpkeepCounter.TransactOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorSession) Start() (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.Start(&_LogUpkeepCounter.TransactOpts)
}

type LogUpkeepCounterPerformingUpkeepIterator struct {
	Event *LogUpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogUpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogUpkeepCounterPerformingUpkeep)
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
		it.Event = new(LogUpkeepCounterPerformingUpkeep)
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

func (it *LogUpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *LogUpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogUpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	Raw           types.Log
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogUpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogUpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterPerformingUpkeepIterator{contract: _LogUpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogUpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogUpkeepCounterPerformingUpkeep)
				if err := _LogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*LogUpkeepCounterPerformingUpkeep, error) {
	event := new(LogUpkeepCounterPerformingUpkeep)
	if err := _LogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogUpkeepCounterTriggerIterator struct {
	Event *LogUpkeepCounterTrigger

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogUpkeepCounterTriggerIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogUpkeepCounterTrigger)
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
		it.Event = new(LogUpkeepCounterTrigger)
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

func (it *LogUpkeepCounterTriggerIterator) Error() error {
	return it.fail
}

func (it *LogUpkeepCounterTriggerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogUpkeepCounterTrigger struct {
	Raw types.Log
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) FilterTrigger(opts *bind.FilterOpts) (*LogUpkeepCounterTriggerIterator, error) {

	logs, sub, err := _LogUpkeepCounter.contract.FilterLogs(opts, "Trigger")
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterTriggerIterator{contract: _LogUpkeepCounter.contract, event: "Trigger", logs: logs, sub: sub}, nil
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) WatchTrigger(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger) (event.Subscription, error) {

	logs, sub, err := _LogUpkeepCounter.contract.WatchLogs(opts, "Trigger")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogUpkeepCounterTrigger)
				if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger", log); err != nil {
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

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) ParseTrigger(log types.Log) (*LogUpkeepCounterTrigger, error) {
	event := new(LogUpkeepCounterTrigger)
	if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogUpkeepCounterTrigger0Iterator struct {
	Event *LogUpkeepCounterTrigger0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogUpkeepCounterTrigger0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogUpkeepCounterTrigger0)
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
		it.Event = new(LogUpkeepCounterTrigger0)
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

func (it *LogUpkeepCounterTrigger0Iterator) Error() error {
	return it.fail
}

func (it *LogUpkeepCounterTrigger0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogUpkeepCounterTrigger0 struct {
	A   *big.Int
	Raw types.Log
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) FilterTrigger0(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger0Iterator, error) {

	logs, sub, err := _LogUpkeepCounter.contract.FilterLogs(opts, "Trigger0")
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterTrigger0Iterator{contract: _LogUpkeepCounter.contract, event: "Trigger0", logs: logs, sub: sub}, nil
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) WatchTrigger0(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger0) (event.Subscription, error) {

	logs, sub, err := _LogUpkeepCounter.contract.WatchLogs(opts, "Trigger0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogUpkeepCounterTrigger0)
				if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger0", log); err != nil {
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

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) ParseTrigger0(log types.Log) (*LogUpkeepCounterTrigger0, error) {
	event := new(LogUpkeepCounterTrigger0)
	if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogUpkeepCounterTrigger1Iterator struct {
	Event *LogUpkeepCounterTrigger1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogUpkeepCounterTrigger1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogUpkeepCounterTrigger1)
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
		it.Event = new(LogUpkeepCounterTrigger1)
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

func (it *LogUpkeepCounterTrigger1Iterator) Error() error {
	return it.fail
}

func (it *LogUpkeepCounterTrigger1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogUpkeepCounterTrigger1 struct {
	A   *big.Int
	B   *big.Int
	Raw types.Log
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) FilterTrigger1(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger1Iterator, error) {

	logs, sub, err := _LogUpkeepCounter.contract.FilterLogs(opts, "Trigger1")
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterTrigger1Iterator{contract: _LogUpkeepCounter.contract, event: "Trigger1", logs: logs, sub: sub}, nil
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) WatchTrigger1(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger1) (event.Subscription, error) {

	logs, sub, err := _LogUpkeepCounter.contract.WatchLogs(opts, "Trigger1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogUpkeepCounterTrigger1)
				if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger1", log); err != nil {
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

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) ParseTrigger1(log types.Log) (*LogUpkeepCounterTrigger1, error) {
	event := new(LogUpkeepCounterTrigger1)
	if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogUpkeepCounterTrigger2Iterator struct {
	Event *LogUpkeepCounterTrigger2

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogUpkeepCounterTrigger2Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogUpkeepCounterTrigger2)
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
		it.Event = new(LogUpkeepCounterTrigger2)
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

func (it *LogUpkeepCounterTrigger2Iterator) Error() error {
	return it.fail
}

func (it *LogUpkeepCounterTrigger2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogUpkeepCounterTrigger2 struct {
	A   *big.Int
	B   *big.Int
	C   *big.Int
	Raw types.Log
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) FilterTrigger2(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger2Iterator, error) {

	logs, sub, err := _LogUpkeepCounter.contract.FilterLogs(opts, "Trigger2")
	if err != nil {
		return nil, err
	}
	return &LogUpkeepCounterTrigger2Iterator{contract: _LogUpkeepCounter.contract, event: "Trigger2", logs: logs, sub: sub}, nil
}

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) WatchTrigger2(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger2) (event.Subscription, error) {

	logs, sub, err := _LogUpkeepCounter.contract.WatchLogs(opts, "Trigger2")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogUpkeepCounterTrigger2)
				if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger2", log); err != nil {
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

func (_LogUpkeepCounter *LogUpkeepCounterFilterer) ParseTrigger2(log types.Log) (*LogUpkeepCounterTrigger2, error) {
	event := new(LogUpkeepCounterTrigger2)
	if err := _LogUpkeepCounter.contract.UnpackLog(event, "Trigger2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LogUpkeepCounter *LogUpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogUpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _LogUpkeepCounter.ParsePerformingUpkeep(log)
	case _LogUpkeepCounter.abi.Events["Trigger"].ID:
		return _LogUpkeepCounter.ParseTrigger(log)
	case _LogUpkeepCounter.abi.Events["Trigger0"].ID:
		return _LogUpkeepCounter.ParseTrigger0(log)
	case _LogUpkeepCounter.abi.Events["Trigger1"].ID:
		return _LogUpkeepCounter.ParseTrigger1(log)
	case _LogUpkeepCounter.abi.Events["Trigger2"].ID:
		return _LogUpkeepCounter.ParseTrigger2(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogUpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa")
}

func (LogUpkeepCounterTrigger) Topic() common.Hash {
	return common.HexToHash("0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d")
}

func (LogUpkeepCounterTrigger0) Topic() common.Hash {
	return common.HexToHash("0x57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da")
}

func (LogUpkeepCounterTrigger1) Topic() common.Hash {
	return common.HexToHash("0x1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c8")
}

func (LogUpkeepCounterTrigger2) Topic() common.Hash {
	return common.HexToHash("0x5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e71")
}

func (_LogUpkeepCounter *LogUpkeepCounter) Address() common.Address {
	return _LogUpkeepCounter.address
}

type LogUpkeepCounterInterface interface {
	CheckUpkeep(opts *bind.CallOpts, logData []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _testRange *big.Int) (*types.Transaction, error)

	Start(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*LogUpkeepCounterPerformingUpkeep, error)

	FilterTrigger(opts *bind.FilterOpts) (*LogUpkeepCounterTriggerIterator, error)

	WatchTrigger(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger) (event.Subscription, error)

	ParseTrigger(log types.Log) (*LogUpkeepCounterTrigger, error)

	FilterTrigger0(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger0Iterator, error)

	WatchTrigger0(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger0) (event.Subscription, error)

	ParseTrigger0(log types.Log) (*LogUpkeepCounterTrigger0, error)

	FilterTrigger1(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger1Iterator, error)

	WatchTrigger1(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger1) (event.Subscription, error)

	ParseTrigger1(log types.Log) (*LogUpkeepCounterTrigger1, error)

	FilterTrigger2(opts *bind.FilterOpts) (*LogUpkeepCounterTrigger2Iterator, error)

	WatchTrigger2(opts *bind.WatchOpts, sink chan<- *LogUpkeepCounterTrigger2) (event.Subscription, error)

	ParseTrigger2(log types.Log) (*LogUpkeepCounterTrigger2, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
