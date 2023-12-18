// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple_log_upkeep_counter_wrapper

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

type CheckData struct {
	CheckBurnAmount   *big.Int
	PerformBurnAmount *big.Int
	EventSig          [32]byte
}

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var SimpleLogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeToPerform\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"checkBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"eventSig\",\"type\":\"bytes32\"}],\"internalType\":\"structCheckData\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_checkDataConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeToPerform\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060006002819055436001556003819055600455610c63806100336000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80636977947311610076578063806b984f1161005b578063806b984f14610147578063917d895f14610150578063c6066f0d1461015957600080fd5b806369779473146101035780637145f11b1461011457600080fd5b80632cb15864146100a857806340691db4146100c45780634585e33b146100e557806361bc221a146100fa575b600080fd5b6100b160035481565b6040519081526020015b60405180910390f35b6100d76100d236600461062d565b610162565b6040516100bb92919061088b565b6100f86100f336600461058f565b610296565b005b6100b160045481565b6100f86101113660046105d1565b50565b610137610122366004610576565b60006020819052908152604090205460ff1681565b60405190151581526020016100bb565b6100b160015481565b6100b160025481565b6100b160055481565b6000606081808061017586880188610799565b92509250925060005a9050600061018d600143610bb2565b409050600085156101fc575b855a6101a59085610bb2565b10156101fc578080156101c6575060008281526020819052604090205460ff165b60408051602081018590523091810191909152909150606001604051602081830303815290604052805190602001209150610199565b8361020a60c08d018d6109ee565b600281811061021b5761021b610bf8565b90506020020135141561025d5760018b8b8b60405160200161023f93929190610908565b6040516020818303038152906040529750975050505050505061028e565b60008b8b8b60405160200161027493929190610908565b604051602081830303815290604052975097505050505050505b935093915050565b6003546102a257436003555b4360019081556004546102b491610b9a565b6004556001546002556000806102cc8385018561069f565b915091508160200151426102e09190610bb2565b6005819055506000806000838060200190518101906102ff91906107c5565b92509250925060005a90506000610317600143610bb2565b40905060008415610386575b845a61032f9085610bb2565b101561038657808015610350575060008281526020819052604090205460ff165b60408051602081018590523091810191909152909150606001604051602081830303815290604052805190602001209150610323565b600354600154600254600454600554604080519586526020860194909452928401919091526060830152608082015232907f4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd659060a00160405180910390a250505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461041457600080fd5b919050565b600082601f83011261042a57600080fd5b8135602067ffffffffffffffff82111561044657610446610c27565b8160051b610455828201610a80565b83815282810190868401838801850189101561047057600080fd5b600093505b85841015610493578035835260019390930192918401918401610475565b50979650505050505050565b60008083601f8401126104b157600080fd5b50813567ffffffffffffffff8111156104c957600080fd5b6020830191508360208285010111156104e157600080fd5b9250929050565b600082601f8301126104f957600080fd5b813567ffffffffffffffff81111561051357610513610c27565b61054460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610a80565b81815284602083860101111561055957600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561058857600080fd5b5035919050565b600080602083850312156105a257600080fd5b823567ffffffffffffffff8111156105b957600080fd5b6105c58582860161049f565b90969095509350505050565b6000606082840312156105e357600080fd5b6040516060810181811067ffffffffffffffff8211171561060657610606610c27565b80604052508235815260208301356020820152604083013560408201528091505092915050565b60008060006040848603121561064257600080fd5b833567ffffffffffffffff8082111561065a57600080fd5b90850190610100828803121561066f57600080fd5b9093506020850135908082111561068557600080fd5b506106928682870161049f565b9497909650939450505050565b600080604083850312156106b257600080fd5b823567ffffffffffffffff808211156106ca57600080fd5b9084019061010082870312156106df57600080fd5b6106e7610a56565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261071f60a084016103f0565b60a082015260c08301358281111561073657600080fd5b61074288828601610419565b60c08301525060e08301358281111561075a57600080fd5b610766888286016104e8565b60e0830152509350602085013591508082111561078257600080fd5b5061078f858286016104e8565b9150509250929050565b6000806000606084860312156107ae57600080fd5b505081359360208301359350604090920135919050565b6000806000606084860312156107da57600080fd5b8351925060208401519150604084015190509250925092565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561082557600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b821515815260006020604081840152835180604085015260005b818110156108c1578581018301518582016060015282016108a5565b818111156108d3576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b60408152833560408201526020840135606082015260408401356080820152606084013560a0820152608084013560c082015273ffffffffffffffffffffffffffffffffffffffff61095c60a086016103f0565b1660e0820152600061097160c0860186610acf565b61010084810152610987610140850182846107f3565b91505061099760e0870187610b36565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016101208601526109cd838284610842565b9250505082810360208401526109e4818587610842565b9695505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610a2357600080fd5b83018035915067ffffffffffffffff821115610a3e57600080fd5b6020019150600581901b36038213156104e157600080fd5b604051610100810167ffffffffffffffff81118282101715610a7a57610a7a610c27565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610ac757610ac7610c27565b604052919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610b0457600080fd5b830160208101925035905067ffffffffffffffff811115610b2457600080fd5b8060051b36038313156104e157600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610b6b57600080fd5b830160208101925035905067ffffffffffffffff811115610b8b57600080fd5b8036038313156104e157600080fd5b60008219821115610bad57610bad610bc9565b500190565b600082821015610bc457610bc4610bc9565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var SimpleLogUpkeepCounterABI = SimpleLogUpkeepCounterMetaData.ABI

var SimpleLogUpkeepCounterBin = SimpleLogUpkeepCounterMetaData.Bin

func DeploySimpleLogUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleLogUpkeepCounter, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleLogUpkeepCounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleLogUpkeepCounter{address: address, abi: *parsed, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

type SimpleLogUpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	SimpleLogUpkeepCounterCaller
	SimpleLogUpkeepCounterTransactor
	SimpleLogUpkeepCounterFilterer
}

type SimpleLogUpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterSession struct {
	Contract     *SimpleLogUpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterCallerSession struct {
	Contract *SimpleLogUpkeepCounterCaller
	CallOpts bind.CallOpts
}

type SimpleLogUpkeepCounterTransactorSession struct {
	Contract     *SimpleLogUpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterRaw struct {
	Contract *SimpleLogUpkeepCounter
}

type SimpleLogUpkeepCounterCallerRaw struct {
	Contract *SimpleLogUpkeepCounterCaller
}

type SimpleLogUpkeepCounterTransactorRaw struct {
	Contract *SimpleLogUpkeepCounterTransactor
}

func NewSimpleLogUpkeepCounter(address common.Address, backend bind.ContractBackend) (*SimpleLogUpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleLogUpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSimpleLogUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounter{address: address, abi: abi, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

func NewSimpleLogUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*SimpleLogUpkeepCounterCaller, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterCaller{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleLogUpkeepCounterTransactor, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterTransactor{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleLogUpkeepCounterFilterer, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterFilterer{contract: contract}, nil
}

func bindSimpleLogUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkLog", log, checkData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) TimeToPerform(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "timeToPerform")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "_checkDataConfig", arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

type SimpleLogUpkeepCounterPerformingUpkeepIterator struct {
	Event *SimpleLogUpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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
		it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleLogUpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	TimeToPerform *big.Int
	Raw           types.Log
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterPerformingUpkeepIterator{contract: _SimpleLogUpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleLogUpkeepCounterPerformingUpkeep)
				if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error) {
	event := new(SimpleLogUpkeepCounterPerformingUpkeep)
	if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SimpleLogUpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _SimpleLogUpkeepCounter.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SimpleLogUpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x4874b8dd61a40fe23599b4360a9a824d7081742fca9f555bcee3d389c4f4bd65")
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) Address() common.Address {
	return _SimpleLogUpkeepCounter.address
}

type SimpleLogUpkeepCounterInterface interface {
	CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TimeToPerform(opts *bind.CallOpts) (*big.Int, error)

	CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
