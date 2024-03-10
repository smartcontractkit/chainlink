// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_consumer_performance_wrapper

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

var KeeperConsumerPerformanceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_performGasToBurn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"eligible\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060075534801561001e57600080fd5b5060405161070438038061070483398101604081905261003d91610054565b60029390935560039190915560045560055561008a565b6000806000806080858703121561006a57600080fd5b505082516020840151604085015160609095015191969095509092509050565b61066b806100996000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637145f11b11610097578063b30566b411610066578063b30566b4146101f6578063c228a98e146101ff578063d826f88f14610207578063e303666f1461021457600080fd5b80637145f11b146101985780637f407edf146101cb578063926f086e146101e4578063a9a4c57c146101ed57600080fd5b80634585e33b116100d35780634585e33b14610152578063523d9b8a146101655780636250a13a1461016e5780636e04ff0d1461017757600080fd5b806306661abd1461010557806313bda75b146101215780632555d2cf146101365780632ff3617d14610149575b600080fd5b61010e60075481565b6040519081526020015b60405180910390f35b61013461012f366004610430565b600455565b005b610134610144366004610430565b600555565b61010e60045481565b610134610160366004610449565b61021c565b61010e60015481565b61010e60025481565b61018a610185366004610449565b610342565b6040516101189291906104bb565b6101bb6101a6366004610430565b60066020526000908152604090205460ff1681565b6040519015158152602001610118565b6101346101d9366004610531565b600291909155600355565b61010e60005481565b61010e60035481565b61010e60055481565b6101bb6103b7565b6101346000808055600755565b60075461010e565b60005a9050600061022b6103c6565b60005460015460408051841515815232602082015290810192909252606082015243608082018190529192507fbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc09060a00160405180910390a18161028e57600080fd5b60005460000361029e5760008190555b6003546102ac906002610582565b6102b46103f2565b6102be91906105bf565b6102c890826105fa565b6102d39060016105fa565b600155600780549060006102e683610613565b91905055505b6005545a6102fa908561064b565b101561033b574340600090815260066020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556102ec565b5050505050565b6000606060005a905060005b6004545a61035c908461064b565b10156103865780801561037f5750434060009081526006602052604090205460ff165b905061034e565b61038e6103c6565b604080518315156020820152016040516020818303038152906040529350935050509250929050565b60006103c16103c6565b905090565b6000805415806103c157506002546000546103e1904361064b565b1080156103c1575050600154431190565b60006103ff60014361064b565b604080519140602083015230908201526060016040516020818303038152906040528051906020012060001c905090565b60006020828403121561044257600080fd5b5035919050565b6000806020838503121561045c57600080fd5b823567ffffffffffffffff8082111561047457600080fd5b818501915085601f83011261048857600080fd5b81358181111561049757600080fd5b8660208285010111156104a957600080fd5b60209290920196919550909350505050565b821515815260006020604081840152835180604085015260005b818110156104f1578581018301518582016060015282016104d5565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b6000806040838503121561054457600080fd5b50508035926020909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156105ba576105ba610553565b500290565b6000826105f5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b8082018082111561060d5761060d610553565b92915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361064457610644610553565b5060010190565b8181038181111561060d5761060d61055356fea164736f6c6343000810000a",
}

var KeeperConsumerPerformanceABI = KeeperConsumerPerformanceMetaData.ABI

var KeeperConsumerPerformanceBin = KeeperConsumerPerformanceMetaData.Bin

func DeployKeeperConsumerPerformance(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _averageEligibilityCadence *big.Int, _checkGasToBurn *big.Int, _performGasToBurn *big.Int) (common.Address, *types.Transaction, *KeeperConsumerPerformance, error) {
	parsed, err := KeeperConsumerPerformanceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperConsumerPerformanceBin), backend, _testRange, _averageEligibilityCadence, _checkGasToBurn, _performGasToBurn)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperConsumerPerformance{address: address, abi: *parsed, KeeperConsumerPerformanceCaller: KeeperConsumerPerformanceCaller{contract: contract}, KeeperConsumerPerformanceTransactor: KeeperConsumerPerformanceTransactor{contract: contract}, KeeperConsumerPerformanceFilterer: KeeperConsumerPerformanceFilterer{contract: contract}}, nil
}

type KeeperConsumerPerformance struct {
	address common.Address
	abi     abi.ABI
	KeeperConsumerPerformanceCaller
	KeeperConsumerPerformanceTransactor
	KeeperConsumerPerformanceFilterer
}

type KeeperConsumerPerformanceCaller struct {
	contract *bind.BoundContract
}

type KeeperConsumerPerformanceTransactor struct {
	contract *bind.BoundContract
}

type KeeperConsumerPerformanceFilterer struct {
	contract *bind.BoundContract
}

type KeeperConsumerPerformanceSession struct {
	Contract     *KeeperConsumerPerformance
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperConsumerPerformanceCallerSession struct {
	Contract *KeeperConsumerPerformanceCaller
	CallOpts bind.CallOpts
}

type KeeperConsumerPerformanceTransactorSession struct {
	Contract     *KeeperConsumerPerformanceTransactor
	TransactOpts bind.TransactOpts
}

type KeeperConsumerPerformanceRaw struct {
	Contract *KeeperConsumerPerformance
}

type KeeperConsumerPerformanceCallerRaw struct {
	Contract *KeeperConsumerPerformanceCaller
}

type KeeperConsumerPerformanceTransactorRaw struct {
	Contract *KeeperConsumerPerformanceTransactor
}

func NewKeeperConsumerPerformance(address common.Address, backend bind.ContractBackend) (*KeeperConsumerPerformance, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperConsumerPerformanceABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperConsumerPerformance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformance{address: address, abi: abi, KeeperConsumerPerformanceCaller: KeeperConsumerPerformanceCaller{contract: contract}, KeeperConsumerPerformanceTransactor: KeeperConsumerPerformanceTransactor{contract: contract}, KeeperConsumerPerformanceFilterer: KeeperConsumerPerformanceFilterer{contract: contract}}, nil
}

func NewKeeperConsumerPerformanceCaller(address common.Address, caller bind.ContractCaller) (*KeeperConsumerPerformanceCaller, error) {
	contract, err := bindKeeperConsumerPerformance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceCaller{contract: contract}, nil
}

func NewKeeperConsumerPerformanceTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperConsumerPerformanceTransactor, error) {
	contract, err := bindKeeperConsumerPerformance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceTransactor{contract: contract}, nil
}

func NewKeeperConsumerPerformanceFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperConsumerPerformanceFilterer, error) {
	contract, err := bindKeeperConsumerPerformance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformanceFilterer{contract: contract}, nil
}

func bindKeeperConsumerPerformance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperConsumerPerformanceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceTransactor.contract.Transfer(opts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.KeeperConsumerPerformanceTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperConsumerPerformance.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.contract.Transfer(opts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "averageEligibilityCadence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.AverageEligibilityCadence(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) AverageEligibilityCadence() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.AverageEligibilityCadence(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) CheckEligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "checkEligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) CheckEligible() (bool, error) {
	return _KeeperConsumerPerformance.Contract.CheckEligible(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) CheckEligible() (bool, error) {
	return _KeeperConsumerPerformance.Contract.CheckEligible(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "checkGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) CheckGasToBurn() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.CheckGasToBurn(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) CheckGasToBurn() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.CheckGasToBurn(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _KeeperConsumerPerformance.Contract.CheckUpkeep(&_KeeperConsumerPerformance.CallOpts, data)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _KeeperConsumerPerformance.Contract.CheckUpkeep(&_KeeperConsumerPerformance.CallOpts, data)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) Count() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.Count(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) Count() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.Count(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _KeeperConsumerPerformance.Contract.DummyMap(&_KeeperConsumerPerformance.CallOpts, arg0)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _KeeperConsumerPerformance.Contract.DummyMap(&_KeeperConsumerPerformance.CallOpts, arg0)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) GetCountPerforms(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "getCountPerforms")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.GetCountPerforms(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) GetCountPerforms() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.GetCountPerforms(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) InitialCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "initialCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.InitialCall(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) InitialCall() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.InitialCall(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) NextEligible(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "nextEligible")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.NextEligible(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) NextEligible() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.NextEligible(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "performGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) PerformGasToBurn() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.PerformGasToBurn(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) PerformGasToBurn() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.PerformGasToBurn(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperConsumerPerformance.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.TestRange(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceCallerSession) TestRange() (*big.Int, error) {
	return _KeeperConsumerPerformance.Contract.TestRange(&_KeeperConsumerPerformance.CallOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "performUpkeep", data)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.PerformUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) PerformUpkeep(data []byte) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.PerformUpkeep(&_KeeperConsumerPerformance.TransactOpts, data)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "reset")
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.Reset(&_KeeperConsumerPerformance.TransactOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) Reset() (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.Reset(&_KeeperConsumerPerformance.TransactOpts)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "setCheckGasToBurn", value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetCheckGasToBurn(&_KeeperConsumerPerformance.TransactOpts, value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetCheckGasToBurn(&_KeeperConsumerPerformance.TransactOpts, value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "setPerformGasToBurn", value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetPerformGasToBurn(&_KeeperConsumerPerformance.TransactOpts, value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetPerformGasToBurn(&_KeeperConsumerPerformance.TransactOpts, value)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactor) SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.contract.Transact(opts, "setSpread", _newTestRange, _newAverageEligibilityCadence)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetSpread(&_KeeperConsumerPerformance.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceTransactorSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _KeeperConsumerPerformance.Contract.SetSpread(&_KeeperConsumerPerformance.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

type KeeperConsumerPerformancePerformingUpkeepIterator struct {
	Event *KeeperConsumerPerformancePerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperConsumerPerformancePerformingUpkeep)
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
		it.Event = new(KeeperConsumerPerformancePerformingUpkeep)
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

func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *KeeperConsumerPerformancePerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperConsumerPerformancePerformingUpkeep struct {
	Eligible     bool
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*KeeperConsumerPerformancePerformingUpkeepIterator, error) {

	logs, sub, err := _KeeperConsumerPerformance.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &KeeperConsumerPerformancePerformingUpkeepIterator{contract: _KeeperConsumerPerformance.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *KeeperConsumerPerformancePerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _KeeperConsumerPerformance.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperConsumerPerformancePerformingUpkeep)
				if err := _KeeperConsumerPerformance.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_KeeperConsumerPerformance *KeeperConsumerPerformanceFilterer) ParsePerformingUpkeep(log types.Log) (*KeeperConsumerPerformancePerformingUpkeep, error) {
	event := new(KeeperConsumerPerformancePerformingUpkeep)
	if err := _KeeperConsumerPerformance.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformance) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperConsumerPerformance.abi.Events["PerformingUpkeep"].ID:
		return _KeeperConsumerPerformance.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperConsumerPerformancePerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0")
}

func (_KeeperConsumerPerformance *KeeperConsumerPerformance) Address() common.Address {
	return _KeeperConsumerPerformance.address
}

type KeeperConsumerPerformanceInterface interface {
	AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error)

	CheckEligible(opts *bind.CallOpts) (bool, error)

	CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Count(opts *bind.CallOpts) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	GetCountPerforms(opts *bind.CallOpts) (*big.Int, error)

	InitialCall(opts *bind.CallOpts) (*big.Int, error)

	NextEligible(opts *bind.CallOpts) (*big.Int, error)

	PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, data []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error)

	SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*KeeperConsumerPerformancePerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *KeeperConsumerPerformancePerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*KeeperConsumerPerformancePerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
