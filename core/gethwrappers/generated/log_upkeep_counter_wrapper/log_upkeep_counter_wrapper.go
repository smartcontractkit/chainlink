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

type Log struct {
	Index       *big.Int
	TxIndex     *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var LogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"}],\"name\":\"Trigger\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040527f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d6000557f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da6001557f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c86002557f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e716003553480156100a057600080fd5b50604051610f17380380610f178339810160408190526100bf916100da565b600455600060068190554360055560078190556008556100f3565b6000602082840312156100ec57600080fd5b5051919050565b610e15806101026000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063917d895f11610076578063be61b7751161005b578063be61b77514610135578063be9a655514610156578063d832d92f1461015e57600080fd5b8063917d895f1461010f578063b66a261c1461011857600080fd5b806361bc221a116100a757806361bc221a146100f45780636250a13a146100fd578063806b984f1461010657600080fd5b80632cb15864146100c35780634585e33b146100df575b600080fd5b6100cc60075481565b6040519081526020015b60405180910390f35b6100f26100ed366004610810565b610176565b005b6100cc60085481565b6100cc60045481565b6100cc60055481565b6100cc60065481565b6100f261012636600461099a565b60045560006007819055600855565b610148610143366004610882565b6103ed565b6040516100d6929190610a4b565b6100f26105d0565b6101666106aa565b60405190151581526020016100d6565b60075461018257436007555b43600555600854610194906001610d4c565b60085560055460065560006101ab828401846108c5565b90506000548160c001516000815181106101c7576101c7610daa565b60200260200101511415610203576040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a1610395565b6001548160c0015160008151811061021d5761021d610daa565b6020026020010151141561026557604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da906020015b60405180910390a1610395565b6002548160c0015160008151811061027f5761027f610daa565b602002602001015114156102c4576040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c89101610258565b6003548160c001516000815181106102de576102de610daa565b6020026020010151141561032e576040805160018152600260208201526003918101919091527f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e7190606001610258565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f636f756c64206e6f742066696e64206d61746368696e6720736967000000000060448201526064015b60405180910390fd5b60075460055460065460085460408051948552602085019390935291830152606082015232907f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa9060800160405180910390a2505050565b600060606103f96106aa565b61045f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f6e6f7420656c696769626c650000000000000000000000000000000000000000604482015260640161038c565b60005461046f60c0850185610b99565b600081811061048057610480610daa565b9050602002013514806104b8575060015461049e60c0850185610b99565b60008181106104af576104af610daa565b90506020020135145b806104e857506002546104ce60c0850185610b99565b60008181106104df576104df610daa565b90506020020135145b8061051857506003546104fe60c0850185610b99565b600081811061050f5761050f610daa565b90506020020135145b15610548576001836040516020016105309190610ac8565b60405160208183030381529060405291509150915091565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f6700000000000000000000000000000000000000000000000000000000000000606482015260840161038c565b6040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a1604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da9060200160405180910390a16040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c8910160405180910390a160408051600181526002602082015260038183015290517f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e719181900360600190a1565b6000600754600014156106bd5750600190565b6004546007546106cd9043610d64565b10905090565b803573ffffffffffffffffffffffffffffffffffffffff811681146106f757600080fd5b919050565b600082601f83011261070d57600080fd5b8135602067ffffffffffffffff82111561072957610729610dd9565b8160051b610738828201610c32565b83815282810190868401838801850189101561075357600080fd5b600093505b85841015610776578035835260019390930192918401918401610758565b50979650505050505050565b600082601f83011261079357600080fd5b813567ffffffffffffffff8111156107ad576107ad610dd9565b6107de60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610c32565b8181528460208386010111156107f357600080fd5b816020850160208301376000918101602001919091529392505050565b6000806020838503121561082357600080fd5b823567ffffffffffffffff8082111561083b57600080fd5b818501915085601f83011261084f57600080fd5b81358181111561085e57600080fd5b86602082850101111561087057600080fd5b60209290920196919550909350505050565b60006020828403121561089457600080fd5b813567ffffffffffffffff8111156108ab57600080fd5b820161010081850312156108be57600080fd5b9392505050565b6000602082840312156108d757600080fd5b813567ffffffffffffffff808211156108ef57600080fd5b90830190610100828603121561090457600080fd5b61090c610c08565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261094460a084016106d3565b60a082015260c08301358281111561095b57600080fd5b610967878286016106fc565b60c08301525060e08301358281111561097f57600080fd5b61098b87828601610782565b60e08301525095945050505050565b6000602082840312156109ac57600080fd5b5035919050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8311156109e557600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b821515815260006020604081840152835180604085015260005b81811015610a8157858101830151858201606001528201610a65565b81811115610a93576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b6020815281356020820152602082013560408201526040820135606082015260608201356080820152608082013560a082015273ffffffffffffffffffffffffffffffffffffffff610b1c60a084016106d3565b1660c08201526000610b3160c0840184610c81565b6101008060e0860152610b49610120860183856109b3565b9250610b5860e0870187610ce8565b92507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08685030182870152610b8e848483610a02565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610bce57600080fd5b83018035915067ffffffffffffffff821115610be957600080fd5b6020019150600581901b3603821315610c0157600080fd5b9250929050565b604051610100810167ffffffffffffffff81118282101715610c2c57610c2c610dd9565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610c7957610c79610dd9565b604052919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610cb657600080fd5b830160208101925035905067ffffffffffffffff811115610cd657600080fd5b8060051b3603831315610c0157600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610d1d57600080fd5b830160208101925035905067ffffffffffffffff811115610d3d57600080fd5b803603831315610c0157600080fd5b60008219821115610d5f57610d5f610d7b565b500190565b600082821015610d7657610d76610d7b565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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

func (_LogUpkeepCounter *LogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log) (bool, []byte, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "checkLog", log)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) CheckLog(log Log) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckLog(&_LogUpkeepCounter.CallOpts, log)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) CheckLog(log Log) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckLog(&_LogUpkeepCounter.CallOpts, log)
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
	CheckLog(opts *bind.CallOpts, log Log) (bool, []byte, error)

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
