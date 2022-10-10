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

// UpkeepResetterMetaData contains all meta data concerning the UpkeepResetter contract.
var UpkeepResetterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"upkeepAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averageEligibilityCadence\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"firstEligibleBuffer\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"ResetManyConsumerBenchmark\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061033f806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80639c7d39ad14610030575b600080fd5b6100e8600480360360c081101561004657600080fd5b810190602081018135600160201b81111561006057600080fd5b82018360208201111561007257600080fd5b803590602001918460208302840111600160201b8311171561009357600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092955050823593505050602081013590604081013590606081013590608001356100ea565b005b60005b865181101561011f5761011787828151811061010557fe5b60200260200101518787878787610128565b6001016100ed565b50505050505050565b6000869050806001600160a01b031663ad0d5c4d856040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561017357600080fd5b505af1158015610187573d6000803e3d6000fd5b50505050806001600160a01b0316637f407edf87876040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b1580156101d957600080fd5b505af11580156101ed573d6000803e3d6000fd5b50505050806001600160a01b03166313bda75b846040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561023757600080fd5b505af115801561024b573d6000803e3d6000fd5b50505050806001600160a01b0316632555d2cf836040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561029557600080fd5b505af11580156102a9573d6000803e3d6000fd5b50505050806001600160a01b031663d826f88f6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102e857600080fd5b505af11580156102fc573d6000803e3d6000fd5b505050505050505050505056fea2646970667358221220bbcbe6fd33f4a3980198632dd5b31899f329b28d1d5803b8376a8a44456358ba64736f6c63430007060033",
}

// UpkeepResetterABI is the input ABI used to generate the binding from.
// Deprecated: Use UpkeepResetterMetaData.ABI instead.
var UpkeepResetterABI = UpkeepResetterMetaData.ABI

// UpkeepResetterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpkeepResetterMetaData.Bin instead.
var UpkeepResetterBin = UpkeepResetterMetaData.Bin

// DeployUpkeepResetter deploys a new Ethereum contract, binding an instance of UpkeepResetter to it.
func DeployUpkeepResetter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UpkeepResetter, error) {
	parsed, err := UpkeepResetterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepResetterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepResetter{UpkeepResetterCaller: UpkeepResetterCaller{contract: contract}, UpkeepResetterTransactor: UpkeepResetterTransactor{contract: contract}, UpkeepResetterFilterer: UpkeepResetterFilterer{contract: contract}}, nil
}

// UpkeepResetter is an auto generated Go binding around an Ethereum contract.
type UpkeepResetter struct {
	UpkeepResetterCaller     // Read-only binding to the contract
	UpkeepResetterTransactor // Write-only binding to the contract
	UpkeepResetterFilterer   // Log filterer for contract events
}

// UpkeepResetterCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpkeepResetterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepResetterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpkeepResetterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepResetterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpkeepResetterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpkeepResetterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpkeepResetterSession struct {
	Contract     *UpkeepResetter   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpkeepResetterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpkeepResetterCallerSession struct {
	Contract *UpkeepResetterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// UpkeepResetterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpkeepResetterTransactorSession struct {
	Contract     *UpkeepResetterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// UpkeepResetterRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpkeepResetterRaw struct {
	Contract *UpkeepResetter // Generic contract binding to access the raw methods on
}

// UpkeepResetterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpkeepResetterCallerRaw struct {
	Contract *UpkeepResetterCaller // Generic read-only contract binding to access the raw methods on
}

// UpkeepResetterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpkeepResetterTransactorRaw struct {
	Contract *UpkeepResetterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpkeepResetter creates a new instance of UpkeepResetter, bound to a specific deployed contract.
func NewUpkeepResetter(address common.Address, backend bind.ContractBackend) (*UpkeepResetter, error) {
	contract, err := bindUpkeepResetter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepResetter{UpkeepResetterCaller: UpkeepResetterCaller{contract: contract}, UpkeepResetterTransactor: UpkeepResetterTransactor{contract: contract}, UpkeepResetterFilterer: UpkeepResetterFilterer{contract: contract}}, nil
}

// NewUpkeepResetterCaller creates a new read-only instance of UpkeepResetter, bound to a specific deployed contract.
func NewUpkeepResetterCaller(address common.Address, caller bind.ContractCaller) (*UpkeepResetterCaller, error) {
	contract, err := bindUpkeepResetter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepResetterCaller{contract: contract}, nil
}

// NewUpkeepResetterTransactor creates a new write-only instance of UpkeepResetter, bound to a specific deployed contract.
func NewUpkeepResetterTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepResetterTransactor, error) {
	contract, err := bindUpkeepResetter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepResetterTransactor{contract: contract}, nil
}

// NewUpkeepResetterFilterer creates a new log filterer instance of UpkeepResetter, bound to a specific deployed contract.
func NewUpkeepResetterFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepResetterFilterer, error) {
	contract, err := bindUpkeepResetter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepResetterFilterer{contract: contract}, nil
}

// bindUpkeepResetter binds a generic wrapper to an already deployed contract.
func bindUpkeepResetter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepResetterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepResetter *UpkeepResetterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepResetter.Contract.UpkeepResetterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepResetter *UpkeepResetterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.UpkeepResetterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepResetter *UpkeepResetterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.UpkeepResetterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpkeepResetter *UpkeepResetterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepResetter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpkeepResetter *UpkeepResetterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpkeepResetter *UpkeepResetterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.contract.Transact(opts, method, params...)
}

// ResetManyConsumerBenchmark is a paid mutator transaction binding the contract method 0x9c7d39ad.
//
// Solidity: function ResetManyConsumerBenchmark(address[] upkeepAddresses, uint256 testRange, uint256 averageEligibilityCadence, uint256 firstEligibleBuffer, uint256 checkGasToBurn, uint256 performGasToBurn) returns()
func (_UpkeepResetter *UpkeepResetterTransactor) ResetManyConsumerBenchmark(opts *bind.TransactOpts, upkeepAddresses []common.Address, testRange *big.Int, averageEligibilityCadence *big.Int, firstEligibleBuffer *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _UpkeepResetter.contract.Transact(opts, "ResetManyConsumerBenchmark", upkeepAddresses, testRange, averageEligibilityCadence, firstEligibleBuffer, checkGasToBurn, performGasToBurn)
}

// ResetManyConsumerBenchmark is a paid mutator transaction binding the contract method 0x9c7d39ad.
//
// Solidity: function ResetManyConsumerBenchmark(address[] upkeepAddresses, uint256 testRange, uint256 averageEligibilityCadence, uint256 firstEligibleBuffer, uint256 checkGasToBurn, uint256 performGasToBurn) returns()
func (_UpkeepResetter *UpkeepResetterSession) ResetManyConsumerBenchmark(upkeepAddresses []common.Address, testRange *big.Int, averageEligibilityCadence *big.Int, firstEligibleBuffer *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.ResetManyConsumerBenchmark(&_UpkeepResetter.TransactOpts, upkeepAddresses, testRange, averageEligibilityCadence, firstEligibleBuffer, checkGasToBurn, performGasToBurn)
}

// ResetManyConsumerBenchmark is a paid mutator transaction binding the contract method 0x9c7d39ad.
//
// Solidity: function ResetManyConsumerBenchmark(address[] upkeepAddresses, uint256 testRange, uint256 averageEligibilityCadence, uint256 firstEligibleBuffer, uint256 checkGasToBurn, uint256 performGasToBurn) returns()
func (_UpkeepResetter *UpkeepResetterTransactorSession) ResetManyConsumerBenchmark(upkeepAddresses []common.Address, testRange *big.Int, averageEligibilityCadence *big.Int, firstEligibleBuffer *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _UpkeepResetter.Contract.ResetManyConsumerBenchmark(&_UpkeepResetter.TransactOpts, upkeepAddresses, testRange, averageEligibilityCadence, firstEligibleBuffer, checkGasToBurn, performGasToBurn)
}
