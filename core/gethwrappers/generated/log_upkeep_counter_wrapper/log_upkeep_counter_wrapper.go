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
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var LogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_testRange\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"autoExecution\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkLog\",\"inputs\":[{\"name\":\"log\",\"type\":\"tuple\",\"internalType\":\"structLog\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"source\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"topics\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"counter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eligible\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"performUpkeep\",\"inputs\":[{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"previousPerformBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAuto\",\"inputs\":[{\"name\":\"_auto\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSpread\",\"inputs\":[{\"name\":\"_testRange\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"start\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRange\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"PerformingUpkeep\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"initialBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lastBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"previousBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"counter\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Trigger\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Trigger\",\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Trigger\",\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"b\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Trigger\",\"inputs\":[{\"name\":\"a\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"b\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"c\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040527f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d6000557f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da6001557f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c86002557f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e716003553480156100a057600080fd5b50604051610fe5380380610fe58339810160408190526100bf916100e7565b600455600060068190554360055560078190556008556009805460ff19166001179055610100565b6000602082840312156100f957600080fd5b5051919050565b610ed68061010f6000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063806b984f11610081578063be9a65551161005b578063be9a655514610189578063cf129c7f14610191578063d832d92f146101d057600080fd5b8063806b984f1461015a578063917d895f14610163578063b66a261c1461016c57600080fd5b806353bebdf3116100b257806353bebdf31461012b57806361bc221a146101485780636250a13a1461015157600080fd5b80632cb15864146100d957806340691db4146100f55780634585e33b14610116575b600080fd5b6100e260075481565b6040519081526020015b60405180910390f35b610108610103366004610920565b6101d8565b6040516100ec929190610b13565b6101296101243660046108ae565b6103c7565b005b6009546101389060ff1681565b60405190151581526020016100ec565b6100e260085481565b6100e260045481565b6100e260055481565b6100e260065481565b61012961017a366004610a62565b60045560006007819055600855565b610129610645565b61012961019f366004610885565b600980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b61013861071f565b600060606101e461071f565b61024f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f6e6f7420656c696769626c65000000000000000000000000000000000000000060448201526064015b60405180910390fd5b60005461025f60c0860186610c61565b600081811061027057610270610e6b565b9050602002013514806102a8575060015461028e60c0860186610c61565b600081811061029f5761029f610e6b565b90506020020135145b806102d857506002546102be60c0860186610c61565b60008181106102cf576102cf610e6b565b90506020020135145b8061030857506003546102ee60c0860186610c61565b60008181106102ff576102ff610e6b565b90506020020135145b15610338576001846040516020016103209190610b90565b604051602081830303815290604052915091506103c0565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401610246565b9250929050565b6007546103d357436007555b436005556008546103e5906001610e0d565b60085560055460065560006103fc8284018461098d565b60095490915060ff16156105ed576000548160c0015160008151811061042457610424610e6b565b60200260200101511415610460576040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a16105ed565b6001548160c0015160008151811061047a5761047a610e6b565b602002602001015114156104c257604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da906020015b60405180910390a16105ed565b6002548160c001516000815181106104dc576104dc610e6b565b60200260200101511415610521576040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c891016104b5565b6003548160c0015160008151811061053b5761053b610e6b565b6020026020010151141561058b576040805160018152600260208201526003918101919091527f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e71906060016104b5565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f636f756c64206e6f742066696e64206d61746368696e672073696700000000006044820152606401610246565b60075460055460065460085460408051948552602085019390935291830152606082015232907f8e8112f20a2134e18e591d2cdd68cd86a95d06e6328ede501fc6314f4a5075fa9060800160405180910390a2505050565b6040517f3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d90600090a1604051600181527f57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da9060200160405180910390a16040805160018152600260208201527f1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c8910160405180910390a160408051600181526002602082015260038183015290517f5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e719181900360600190a1565b6000600754600014156107325750600190565b6004546007546107429043610e25565b10905090565b803573ffffffffffffffffffffffffffffffffffffffff8116811461076c57600080fd5b919050565b600082601f83011261078257600080fd5b8135602067ffffffffffffffff82111561079e5761079e610e9a565b8160051b6107ad828201610cf3565b8381528281019086840183880185018910156107c857600080fd5b600093505b858410156107eb5780358352600193909301929184019184016107cd565b50979650505050505050565b600082601f83011261080857600080fd5b813567ffffffffffffffff81111561082257610822610e9a565b61085360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610cf3565b81815284602083860101111561086857600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561089757600080fd5b813580151581146108a757600080fd5b9392505050565b600080602083850312156108c157600080fd5b823567ffffffffffffffff808211156108d957600080fd5b818501915085601f8301126108ed57600080fd5b8135818111156108fc57600080fd5b86602082850101111561090e57600080fd5b60209290920196919550909350505050565b6000806040838503121561093357600080fd5b823567ffffffffffffffff8082111561094b57600080fd5b90840190610100828703121561096057600080fd5b9092506020840135908082111561097657600080fd5b50610983858286016107f7565b9150509250929050565b60006020828403121561099f57600080fd5b813567ffffffffffffffff808211156109b757600080fd5b9083019061010082860312156109cc57600080fd5b6109d4610cc9565b8235815260208301356020820152604083013560408201526060830135606082015260808301356080820152610a0c60a08401610748565b60a082015260c083013582811115610a2357600080fd5b610a2f87828601610771565b60c08301525060e083013582811115610a4757600080fd5b610a53878286016107f7565b60e08301525095945050505050565b600060208284031215610a7457600080fd5b5035919050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610aad57600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b821515815260006020604081840152835180604085015260005b81811015610b4957858101830151858201606001528201610b2d565b81811115610b5b576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b6020815281356020820152602082013560408201526040820135606082015260608201356080820152608082013560a082015273ffffffffffffffffffffffffffffffffffffffff610be460a08401610748565b1660c08201526000610bf960c0840184610d42565b6101008060e0860152610c1161012086018385610a7b565b9250610c2060e0870187610da9565b92507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08685030182870152610c56848483610aca565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610c9657600080fd5b83018035915067ffffffffffffffff821115610cb157600080fd5b6020019150600581901b36038213156103c057600080fd5b604051610100810167ffffffffffffffff81118282101715610ced57610ced610e9a565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d3a57610d3a610e9a565b604052919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610d7757600080fd5b830160208101925035905067ffffffffffffffff811115610d9757600080fd5b8060051b36038313156103c057600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610dde57600080fd5b830160208101925035905067ffffffffffffffff811115610dfe57600080fd5b8036038313156103c057600080fd5b60008219821115610e2057610e20610e3c565b500190565b600082821015610e3757610e37610e3c565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &LogUpkeepCounter{address: address, abi: *parsed, LogUpkeepCounterCaller: LogUpkeepCounterCaller{contract: contract}, LogUpkeepCounterTransactor: LogUpkeepCounterTransactor{contract: contract}, LogUpkeepCounterFilterer: LogUpkeepCounterFilterer{contract: contract}}, nil
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

func (_LogUpkeepCounter *LogUpkeepCounterCaller) AutoExecution(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "autoExecution")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) AutoExecution() (bool, error) {
	return _LogUpkeepCounter.Contract.AutoExecution(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) AutoExecution() (bool, error) {
	return _LogUpkeepCounter.Contract.AutoExecution(&_LogUpkeepCounter.CallOpts)
}

func (_LogUpkeepCounter *LogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log, arg1 []byte) (bool, []byte, error) {
	var out []interface{}
	err := _LogUpkeepCounter.contract.Call(opts, &out, "checkLog", log, arg1)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_LogUpkeepCounter *LogUpkeepCounterSession) CheckLog(log Log, arg1 []byte) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckLog(&_LogUpkeepCounter.CallOpts, log, arg1)
}

func (_LogUpkeepCounter *LogUpkeepCounterCallerSession) CheckLog(log Log, arg1 []byte) (bool, []byte, error) {
	return _LogUpkeepCounter.Contract.CheckLog(&_LogUpkeepCounter.CallOpts, log, arg1)
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

func (_LogUpkeepCounter *LogUpkeepCounterTransactor) SetAuto(opts *bind.TransactOpts, _auto bool) (*types.Transaction, error) {
	return _LogUpkeepCounter.contract.Transact(opts, "setAuto", _auto)
}

func (_LogUpkeepCounter *LogUpkeepCounterSession) SetAuto(_auto bool) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.SetAuto(&_LogUpkeepCounter.TransactOpts, _auto)
}

func (_LogUpkeepCounter *LogUpkeepCounterTransactorSession) SetAuto(_auto bool) (*types.Transaction, error) {
	return _LogUpkeepCounter.Contract.SetAuto(&_LogUpkeepCounter.TransactOpts, _auto)
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
	AutoExecution(opts *bind.CallOpts) (bool, error)

	CheckLog(opts *bind.CallOpts, log Log, arg1 []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetAuto(opts *bind.TransactOpts, _auto bool) (*types.Transaction, error)

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
