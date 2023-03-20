// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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

// UpkeepPerformCounterRestrictiveMetaData contains all meta data concerning the UpkeepPerformCounterRestrictive contract.
var UpkeepPerformCounterRestrictiveMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_averageEligibilityCadence\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"eligible\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialCall\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextEligible\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"averageEligibilityCadence\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCountPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextEligible\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performGasToBurn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newTestRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newAverageEligibilityCadence\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155600060075534801561001e57600080fd5b506040516105c93803806105c98339818101604052604081101561004157600080fd5b508051602090910151600291909155600355610567806100626000396000f3fe608060405234801561001057600080fd5b50600436106100c55760003560e01c806313bda75b146100ca5780632555d2cf146100e95780632ff3617d146101065780634585e33b14610120578063523d9b8a1461018e5780636250a13a146101965780636e04ff0d1461019e5780637145f11b1461028d5780637f407edf146102be578063926f086e146102e1578063a9a4c57c146102e9578063b30566b4146102f1578063c228a98e146102f9578063d826f88f14610301578063e303666f14610309575b600080fd5b6100e7600480360360208110156100e057600080fd5b5035610311565b005b6100e7600480360360208110156100ff57600080fd5b5035610316565b61010e61031b565b60408051918252519081900360200190f35b6100e76004803603602081101561013657600080fd5b810190602081018135600160201b81111561015057600080fd5b82018360208201111561016257600080fd5b803590602001918460018302840111600160201b8311171561018357600080fd5b509092509050610321565b61010e610404565b61010e61040a565b61020c600480360360208110156101b457600080fd5b810190602081018135600160201b8111156101ce57600080fd5b8201836020820111156101e057600080fd5b803590602001918460018302840111600160201b8311171561020157600080fd5b509092509050610410565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610251578181015183820152602001610239565b50505050905090810190601f16801561027e5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6102aa600480360360208110156102a357600080fd5b503561048b565b604080519115158252519081900360200190f35b6100e7600480360360408110156102d457600080fd5b50803590602001356104a0565b61010e6104ab565b61010e6104b1565b61010e6104b7565b6102aa6104bd565b6100e76104cc565b61010e6104d6565b600455565b600555565b60045481565b60005a905060006103306104dc565b60005460015460408051841515815232602082015280820193909352606083019190915243608083018190529051929350917fbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc09181900360a00190a18161039657600080fd5b6000546103a35760008190555b6003546002026103b1610500565b816103b857fe5b06810160019081018155600780549091019055600019015b6005545a840310156103fd5780406000908152600660205260409020805460ff19169055600019016103d0565b5050505050565b60015481565b60025481565b6000606060005a9050600019430160005b6004545a840310156104585780801561044a5750814060009081526006602052604090205460ff165b600019909201919050610421565b6104606104dc565b6040805192151560208085019190915281518085039091018152928101905297909650945050505050565b60066020526000908152604090205460ff1681565b600291909155600355565b60005481565b60035481565b60055481565b60006104c76104dc565b905090565b6000808055600755565b60075490565b6000805415806104c7575060025460005443031080156104c7575050600154431190565b604080516000194301406020808301919091523082840152825180830384018152606090920190925280519101209056fea2646970667358221220a317dab4792a9ae36241f654e15b5fbb29c4a249e54654ea0b02219f739b347a64736f6c63430007060033",
}

// UpkeepPerformCounterRestrictiveABI is the input ABI used to generate the binding from.
// Deprecated: Use UpkeepPerformCounterRestrictiveMetaData.ABI instead.
var UpkeepPerformCounterRestrictiveABI = UpkeepPerformCounterRestrictiveMetaData.ABI

// UpkeepPerformCounterRestrictiveBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpkeepPerformCounterRestrictiveMetaData.Bin instead.
var UpkeepPerformCounterRestrictiveBin = UpkeepPerformCounterRestrictiveMetaData.Bin

// DeployUpkeepPerformCounterRestrictive deploys a new Ethereum contract, binding an instance of UpkeepPerformCounterRestrictive to it.
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

// UpkeepPerformCounterRestrictive is an auto generated Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictive struct {
	UpkeepPerformCounterRestrictiveCaller     // Read-only binding to the contract
	UpkeepPerformCounterRestrictiveTransactor // Write-only binding to the contract
	UpkeepPerformCounterRestrictiveFilterer   // Log filterer for contract events
}

// UpkeepPerformCounterRestrictiveCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictiveCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepPerformCounterRestrictiveTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictiveTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepPerformCounterRestrictiveFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpkeepPerformCounterRestrictiveFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepPerformCounterRestrictiveSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpkeepPerformCounterRestrictiveSession struct {
	Contract     *UpkeepPerformCounterRestrictive // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// UpkeepPerformCounterRestrictiveCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpkeepPerformCounterRestrictiveCallerSession struct {
	Contract *UpkeepPerformCounterRestrictiveCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// UpkeepPerformCounterRestrictiveTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpkeepPerformCounterRestrictiveTransactorSession struct {
	Contract     *UpkeepPerformCounterRestrictiveTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// UpkeepPerformCounterRestrictiveRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictiveRaw struct {
	Contract *UpkeepPerformCounterRestrictive // Generic contract binding to access the raw methods on
}

// UpkeepPerformCounterRestrictiveCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictiveCallerRaw struct {
	Contract *UpkeepPerformCounterRestrictiveCaller // Generic read-only contract binding to access the raw methods on
}

// UpkeepPerformCounterRestrictiveTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpkeepPerformCounterRestrictiveTransactorRaw struct {
	Contract *UpkeepPerformCounterRestrictiveTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpkeepPerformCounterRestrictive creates a new instance of UpkeepPerformCounterRestrictive, bound to a specific deployed contract.
func NewUpkeepPerformCounterRestrictive(address common.Address, backend bind.ContractBackend) (*UpkeepPerformCounterRestrictive, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictive{UpkeepPerformCounterRestrictiveCaller: UpkeepPerformCounterRestrictiveCaller{contract: contract}, UpkeepPerformCounterRestrictiveTransactor: UpkeepPerformCounterRestrictiveTransactor{contract: contract}, UpkeepPerformCounterRestrictiveFilterer: UpkeepPerformCounterRestrictiveFilterer{contract: contract}}, nil
}

// NewUpkeepPerformCounterRestrictiveCaller creates a new read-only instance of UpkeepPerformCounterRestrictive, bound to a specific deployed contract.
func NewUpkeepPerformCounterRestrictiveCaller(address common.Address, caller bind.ContractCaller) (*UpkeepPerformCounterRestrictiveCaller, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveCaller{contract: contract}, nil
}

// NewUpkeepPerformCounterRestrictiveTransactor creates a new write-only instance of UpkeepPerformCounterRestrictive, bound to a specific deployed contract.
func NewUpkeepPerformCounterRestrictiveTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepPerformCounterRestrictiveTransactor, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveTransactor{contract: contract}, nil
}

// NewUpkeepPerformCounterRestrictiveFilterer creates a new log filterer instance of UpkeepPerformCounterRestrictive, bound to a specific deployed contract.
func NewUpkeepPerformCounterRestrictiveFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepPerformCounterRestrictiveFilterer, error) {
	contract, err := bindUpkeepPerformCounterRestrictive(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictiveFilterer{contract: contract}, nil
}

// bindUpkeepPerformCounterRestrictive binds a generic wrapper to an already deployed contract.
func bindUpkeepPerformCounterRestrictive(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepPerformCounterRestrictiveABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.UpkeepPerformCounterRestrictiveTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.contract.Transact(opts, method, params...)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) AverageEligibilityCadence(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "averageEligibilityCadence")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) AverageEligibilityCadence() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.AverageEligibilityCadence(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// AverageEligibilityCadence is a free data retrieval call binding the contract method 0xa9a4c57c.
//
// Solidity: function averageEligibilityCadence() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) AverageEligibilityCadence() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.AverageEligibilityCadence(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) CheckEligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "checkEligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckEligible() (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// CheckEligible is a free data retrieval call binding the contract method 0xc228a98e.
//
// Solidity: function checkEligible() view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckEligible() (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) CheckGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "checkGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// CheckGasToBurn is a free data retrieval call binding the contract method 0x2ff3617d.
//
// Solidity: function checkGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
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

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckUpkeep(&_UpkeepPerformCounterRestrictive.CallOpts, data)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes data) view returns(bool, bytes)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepPerformCounterRestrictive.Contract.CheckUpkeep(&_UpkeepPerformCounterRestrictive.CallOpts, data)
}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.DummyMap(&_UpkeepPerformCounterRestrictive.CallOpts, arg0)
}

// DummyMap is a free data retrieval call binding the contract method 0x7145f11b.
//
// Solidity: function dummyMap(bytes32 ) view returns(bool)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _UpkeepPerformCounterRestrictive.Contract.DummyMap(&_UpkeepPerformCounterRestrictive.CallOpts, arg0)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) GetCountPerforms(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "getCountPerforms")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) GetCountPerforms() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.GetCountPerforms(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// GetCountPerforms is a free data retrieval call binding the contract method 0xe303666f.
//
// Solidity: function getCountPerforms() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) GetCountPerforms() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.GetCountPerforms(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) InitialCall(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "initialCall")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) InitialCall() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.InitialCall(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// InitialCall is a free data retrieval call binding the contract method 0x926f086e.
//
// Solidity: function initialCall() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) InitialCall() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.InitialCall(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) NextEligible(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "nextEligible")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) NextEligible() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.NextEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// NextEligible is a free data retrieval call binding the contract method 0x523d9b8a.
//
// Solidity: function nextEligible() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) NextEligible() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.NextEligible(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) PerformGasToBurn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "performGasToBurn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) PerformGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// PerformGasToBurn is a free data retrieval call binding the contract method 0xb30566b4.
//
// Solidity: function performGasToBurn() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) PerformGasToBurn() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformGasToBurn(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepPerformCounterRestrictive.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) TestRange() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.TestRange(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// TestRange is a free data retrieval call binding the contract method 0x6250a13a.
//
// Solidity: function testRange() view returns(uint256)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepPerformCounterRestrictive.Contract.TestRange(&_UpkeepPerformCounterRestrictive.CallOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) PerformUpkeep(opts *bind.TransactOpts, arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "performUpkeep", arg0)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, arg0)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes ) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) PerformUpkeep(arg0 []byte) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.PerformUpkeep(&_UpkeepPerformCounterRestrictive.TransactOpts, arg0)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) Reset() (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.Reset(&_UpkeepPerformCounterRestrictive.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) Reset() (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.Reset(&_UpkeepPerformCounterRestrictive.TransactOpts)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setCheckGasToBurn", value)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetCheckGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

// SetCheckGasToBurn is a paid mutator transaction binding the contract method 0x13bda75b.
//
// Solidity: function setCheckGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetCheckGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetCheckGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setPerformGasToBurn", value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetPerformGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

// SetPerformGasToBurn is a paid mutator transaction binding the contract method 0x2555d2cf.
//
// Solidity: function setPerformGasToBurn(uint256 value) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetPerformGasToBurn(value *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetPerformGasToBurn(&_UpkeepPerformCounterRestrictive.TransactOpts, value)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactor) SetSpread(opts *bind.TransactOpts, _newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.contract.Transact(opts, "setSpread", _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetSpread(&_UpkeepPerformCounterRestrictive.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// SetSpread is a paid mutator transaction binding the contract method 0x7f407edf.
//
// Solidity: function setSpread(uint256 _newTestRange, uint256 _newAverageEligibilityCadence) returns()
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveTransactorSession) SetSpread(_newTestRange *big.Int, _newAverageEligibilityCadence *big.Int) (*types.Transaction, error) {
	return _UpkeepPerformCounterRestrictive.Contract.SetSpread(&_UpkeepPerformCounterRestrictive.TransactOpts, _newTestRange, _newAverageEligibilityCadence)
}

// UpkeepPerformCounterRestrictivePerformingUpkeepIterator is returned from FilterPerformingUpkeep and is used to iterate over the raw logs and unpacked data for PerformingUpkeep events raised by the UpkeepPerformCounterRestrictive contract.
type UpkeepPerformCounterRestrictivePerformingUpkeepIterator struct {
	Event *UpkeepPerformCounterRestrictivePerformingUpkeep // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpkeepPerformCounterRestrictivePerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpkeepPerformCounterRestrictivePerformingUpkeep represents a PerformingUpkeep event raised by the UpkeepPerformCounterRestrictive contract.
type UpkeepPerformCounterRestrictivePerformingUpkeep struct {
	Eligible     bool
	From         common.Address
	InitialCall  *big.Int
	NextEligible *big.Int
	BlockNumber  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPerformingUpkeep is a free log retrieval operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*UpkeepPerformCounterRestrictivePerformingUpkeepIterator, error) {

	logs, sub, err := _UpkeepPerformCounterRestrictive.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &UpkeepPerformCounterRestrictivePerformingUpkeepIterator{contract: _UpkeepPerformCounterRestrictive.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

// WatchPerformingUpkeep is a free log subscription operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePerformingUpkeep is a log parse operation binding the contract event 0xbd6b6608a51477954e8b498c633bda87e5cd555e06ead50486398d9e3b9cebc0.
//
// Solidity: event PerformingUpkeep(bool eligible, address from, uint256 initialCall, uint256 nextEligible, uint256 blockNumber)
func (_UpkeepPerformCounterRestrictive *UpkeepPerformCounterRestrictiveFilterer) ParsePerformingUpkeep(log types.Log) (*UpkeepPerformCounterRestrictivePerformingUpkeep, error) {
	event := new(UpkeepPerformCounterRestrictivePerformingUpkeep)
	if err := _UpkeepPerformCounterRestrictive.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
