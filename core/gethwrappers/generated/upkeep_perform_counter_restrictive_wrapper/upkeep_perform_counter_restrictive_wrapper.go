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

var UpkeepPerformCounterRestrictiveMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"eligible\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060075534801561001e57600080fd5b5060405161074238038061074283398101604081905261003d9161004b565b60029190915560035561006f565b6000806040838503121561005e57600080fd5b505080516020909101519092909150565b6106c48061007e6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80637145f11b11610097578063b30566b411610066578063b30566b4146101e2578063c228a98e146101eb578063d826f88f146101f3578063e303666f1461020057600080fd5b80637145f11b146101845780637f407edf146101b7578063926f086e146101d0578063a9a4c57c146101d957600080fd5b80634585e33b116100d35780634585e33b1461013e578063523d9b8a146101515780636250a13a1461015a5780636e04ff0d1461016357600080fd5b806313bda75b146100fa5780632555d2cf1461010f5780632ff3617d14610122575b600080fd5b61010d610108366004610454565b600455565b005b61010d61011d366004610454565b600555565b61012b60045481565b6040519081526020015b60405180910390f35b61010d61014c36600461046d565b610208565b61012b60015481565b61012b60025481565b61017661017136600461046d565b610349565b6040516101359291906104df565b6101a7610192366004610454565b60066020526000908152604090205460ff1681565b6040519015158152602001610135565b61010d6101c5366004610555565b600291909155600355565b61012b60005481565b61012b60035481565b61012b60055481565b6101a76103db565b61010d6000808055600755565b60075461012b565b60005a905060006102176103ea565b60005460015460408051841515815232602082015290810192909252606082015243608082018190529192507fbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc09060a00160405180910390a18161027a57600080fd5b60005460000361028a5760008190555b6003546102989060026105a6565b6102a0610416565b6102aa91906105e3565b6102b4908261061e565b6102bf90600161061e565b600155600780549060006102d283610637565b919050555080806102e29061066f565b9150505b6005545a6102f490856106a4565b1015610342578040600090815260066020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558061033a8161066f565b9150506102e6565b5050505050565b6000606060005a9050600061035f6001436106a4565b905060005b6004545a61037290856106a4565b10156103a9578080156103955750814060009081526006602052604090205460ff165b9050816103a18161066f565b925050610364565b6103b16103ea565b60408051831515602082015201604051602081830303815290604052945094505050509250929050565b60006103e56103ea565b905090565b6000805415806103e5575060025460005461040590436106a4565b1080156103e5575050600154431190565b60006104236001436106a4565b604080519140602083015230908201526060016040516020818303038152906040528051906020012060001c905090565b60006020828403121561046657600080fd5b5035919050565b6000806020838503121561048057600080fd5b823567ffffffffffffffff8082111561049857600080fd5b818501915085601f8301126104ac57600080fd5b8135818111156104bb57600080fd5b8660208285010111156104cd57600080fd5b60209290920196919550909350505050565b821515815260006020604081840152835180604085015260005b81811015610515578581018301518582016060015282016104f9565b5060006060828601015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509392505050565b6000806040838503121561056857600080fd5b50508035926020909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156105de576105de610577565b500290565b600082610619577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b8082018082111561063157610631610577565b92915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361066857610668610577565b5060010190565b60008161067e5761067e610577565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b818103818111156106315761063161057756fea164736f6c6343000810000a",
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
	return address, tx, &UpkeepPerformCounterRestrictive{address: address, abi: *parsed, UpkeepPerformCounterRestrictiveCaller: UpkeepPerformCounterRestrictiveCaller{contract: contract}, UpkeepPerformCounterRestrictiveTransactor: UpkeepPerformCounterRestrictiveTransactor{contract: contract}, UpkeepPerformCounterRestrictiveFilterer: UpkeepPerformCounterRestrictiveFilterer{contract: contract}}, nil
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
	parsed, err := UpkeepPerformCounterRestrictiveMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "checkGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.DummyMap(&_UpkeepPerformCounterRestrictive.CallOpts, arg0)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.DummyMap(&_UpkeepPerformCounterRestrictive.CallOpts, arg0)
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "performGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) PerformGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) PerformGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "performUpkeep", arg0)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, arg0)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, arg0)
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

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setCheckGasToBurn", value)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetCheckGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetCheckGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setPerformGasToBurn", value)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetPerformGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetPerformGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
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

	CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	GetCountPerforms(opts *bind.CallOpts) (*big.Int, error)

	InitialCall(opts *bind.CallOpts) (*big.Int, error)

	NextEligible(opts *bind.CallOpts) (*big.Int, error)

	PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error)

	SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*UpkeepPerformCounterRestrictivePerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepPerformCounterRestrictivePerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*UpkeepPerformCounterRestrictivePerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
