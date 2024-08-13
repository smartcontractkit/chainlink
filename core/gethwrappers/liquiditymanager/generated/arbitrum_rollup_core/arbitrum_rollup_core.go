// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_rollup_core

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

type Assertion struct {
	BeforeState ExecutionState
	AfterState  ExecutionState
	NumBlocks   uint64
}

type ExecutionState struct {
	GlobalState   GlobalState
	MachineStatus uint8
}

type GlobalState struct {
	Bytes32Vals [2][32]byte
	U64Vals     [2]uint64
}

type IRollupCoreStaker struct {
	AmountStaked     *big.Int
	Index            uint64
	LatestStakedNode uint64
	CurrentChallenge uint64
	IsStaked         bool
}

type Node struct {
	StateHash                   [32]byte
	ChallengeHash               [32]byte
	ConfirmData                 [32]byte
	PrevNum                     uint64
	DeadlineBlock               uint64
	NoChildConfirmedBeforeBlock uint64
	StakerCount                 uint64
	ChildStakerCount            uint64
	FirstChildBlock             uint64
	LatestChildNumber           uint64
	CreatedAtBlock              uint64
	NodeHash                    [32]byte
}

var ArbRollupCoreMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"sendRoot\",\"type\":\"bytes32\"}],\"name\":\"NodeConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"parentNodeHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"nodeHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"executionHash\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32[2]\",\"name\":\"bytes32Vals\",\"type\":\"bytes32[2]\"},{\"internalType\":\"uint64[2]\",\"name\":\"u64Vals\",\"type\":\"uint64[2]\"}],\"internalType\":\"structGlobalState\",\"name\":\"globalState\",\"type\":\"tuple\"},{\"internalType\":\"enumMachineStatus\",\"name\":\"machineStatus\",\"type\":\"uint8\"}],\"internalType\":\"structExecutionState\",\"name\":\"beforeState\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32[2]\",\"name\":\"bytes32Vals\",\"type\":\"bytes32[2]\"},{\"internalType\":\"uint64[2]\",\"name\":\"u64Vals\",\"type\":\"uint64[2]\"}],\"internalType\":\"structGlobalState\",\"name\":\"globalState\",\"type\":\"tuple\"},{\"internalType\":\"enumMachineStatus\",\"name\":\"machineStatus\",\"type\":\"uint8\"}],\"internalType\":\"structExecutionState\",\"name\":\"afterState\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"numBlocks\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structAssertion\",\"name\":\"assertion\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"afterInboxBatchAcc\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"wasmModuleRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"inboxMaxCount\",\"type\":\"uint256\"}],\"name\":\"NodeCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"}],\"name\":\"NodeRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"challengeIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asserter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"challenger\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"challengedNode\",\"type\":\"uint64\"}],\"name\":\"RollupChallengeStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"machineHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"RollupInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalBalance\",\"type\":\"uint256\"}],\"name\":\"UserStakeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalBalance\",\"type\":\"uint256\"}],\"name\":\"UserWithdrawableFundsUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"amountStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractIBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengeManager\",\"outputs\":[{\"internalType\":\"contractIChallengeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"confirmPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"currentChallenge\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"extraChallengeTimeBlocks\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"firstUnresolvedNode\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"challengeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"confirmData\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"prevNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"deadlineBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"noChildConfirmedBeforeBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"stakerCount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"childStakerCount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"firstChildBlock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestChildNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"createdAtBlock\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"nodeHash\",\"type\":\"bytes32\"}],\"internalType\":\"structNode\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"}],\"name\":\"getNodeCreationBlockForLogLookup\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"getStaker\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amountStaked\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"latestStakedNode\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"currentChallenge\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isStaked\",\"type\":\"bool\"}],\"internalType\":\"structIRollupCore.Staker\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"stakerNum\",\"type\":\"uint64\"}],\"name\":\"getStakerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"isStaked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"isZombie\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastStakeBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfirmed\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestNodeCreated\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"latestStakedNode\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"loserStakeEscrow\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumAssertionPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nodeNum\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"nodeHasStaker\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outbox\",\"outputs\":[{\"internalType\":\"contractIOutbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupEventInbox\",\"outputs\":[{\"internalType\":\"contractIRollupEventInbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sequencerInbox\",\"outputs\":[{\"internalType\":\"contractISequencerInbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakerCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorWhitelistDisabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wasmModuleRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"withdrawableFunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"zombieNum\",\"type\":\"uint256\"}],\"name\":\"zombieAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zombieCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"zombieNum\",\"type\":\"uint256\"}],\"name\":\"zombieLatestStakedNode\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var ArbRollupCoreABI = ArbRollupCoreMetaData.ABI

type ArbRollupCore struct {
	address common.Address
	abi     abi.ABI
	ArbRollupCoreCaller
	ArbRollupCoreTransactor
	ArbRollupCoreFilterer
}

type ArbRollupCoreCaller struct {
	contract *bind.BoundContract
}

type ArbRollupCoreTransactor struct {
	contract *bind.BoundContract
}

type ArbRollupCoreFilterer struct {
	contract *bind.BoundContract
}

type ArbRollupCoreSession struct {
	Contract     *ArbRollupCore
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbRollupCoreCallerSession struct {
	Contract *ArbRollupCoreCaller
	CallOpts bind.CallOpts
}

type ArbRollupCoreTransactorSession struct {
	Contract     *ArbRollupCoreTransactor
	TransactOpts bind.TransactOpts
}

type ArbRollupCoreRaw struct {
	Contract *ArbRollupCore
}

type ArbRollupCoreCallerRaw struct {
	Contract *ArbRollupCoreCaller
}

type ArbRollupCoreTransactorRaw struct {
	Contract *ArbRollupCoreTransactor
}

func NewArbRollupCore(address common.Address, backend bind.ContractBackend) (*ArbRollupCore, error) {
	abi, err := abi.JSON(strings.NewReader(ArbRollupCoreABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbRollupCore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCore{address: address, abi: abi, ArbRollupCoreCaller: ArbRollupCoreCaller{contract: contract}, ArbRollupCoreTransactor: ArbRollupCoreTransactor{contract: contract}, ArbRollupCoreFilterer: ArbRollupCoreFilterer{contract: contract}}, nil
}

func NewArbRollupCoreCaller(address common.Address, caller bind.ContractCaller) (*ArbRollupCoreCaller, error) {
	contract, err := bindArbRollupCore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreCaller{contract: contract}, nil
}

func NewArbRollupCoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbRollupCoreTransactor, error) {
	contract, err := bindArbRollupCore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreTransactor{contract: contract}, nil
}

func NewArbRollupCoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbRollupCoreFilterer, error) {
	contract, err := bindArbRollupCore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreFilterer{contract: contract}, nil
}

func bindArbRollupCore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbRollupCoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbRollupCore *ArbRollupCoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbRollupCore.Contract.ArbRollupCoreCaller.contract.Call(opts, result, method, params...)
}

func (_ArbRollupCore *ArbRollupCoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbRollupCore.Contract.ArbRollupCoreTransactor.contract.Transfer(opts)
}

func (_ArbRollupCore *ArbRollupCoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbRollupCore.Contract.ArbRollupCoreTransactor.contract.Transact(opts, method, params...)
}

func (_ArbRollupCore *ArbRollupCoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbRollupCore.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbRollupCore *ArbRollupCoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbRollupCore.Contract.contract.Transfer(opts)
}

func (_ArbRollupCore *ArbRollupCoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbRollupCore.Contract.contract.Transact(opts, method, params...)
}

func (_ArbRollupCore *ArbRollupCoreCaller) AmountStaked(opts *bind.CallOpts, staker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "amountStaked", staker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) AmountStaked(staker common.Address) (*big.Int, error) {
	return _ArbRollupCore.Contract.AmountStaked(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) AmountStaked(staker common.Address) (*big.Int, error) {
	return _ArbRollupCore.Contract.AmountStaked(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) BaseStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "baseStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) BaseStake() (*big.Int, error) {
	return _ArbRollupCore.Contract.BaseStake(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) BaseStake() (*big.Int, error) {
	return _ArbRollupCore.Contract.BaseStake(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) Bridge() (common.Address, error) {
	return _ArbRollupCore.Contract.Bridge(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) Bridge() (common.Address, error) {
	return _ArbRollupCore.Contract.Bridge(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "chainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ChainId() (*big.Int, error) {
	return _ArbRollupCore.Contract.ChainId(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ChainId() (*big.Int, error) {
	return _ArbRollupCore.Contract.ChainId(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ChallengeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "challengeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ChallengeManager() (common.Address, error) {
	return _ArbRollupCore.Contract.ChallengeManager(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ChallengeManager() (common.Address, error) {
	return _ArbRollupCore.Contract.ChallengeManager(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ConfirmPeriodBlocks(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "confirmPeriodBlocks")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ConfirmPeriodBlocks() (uint64, error) {
	return _ArbRollupCore.Contract.ConfirmPeriodBlocks(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ConfirmPeriodBlocks() (uint64, error) {
	return _ArbRollupCore.Contract.ConfirmPeriodBlocks(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) CurrentChallenge(opts *bind.CallOpts, staker common.Address) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "currentChallenge", staker)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) CurrentChallenge(staker common.Address) (uint64, error) {
	return _ArbRollupCore.Contract.CurrentChallenge(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) CurrentChallenge(staker common.Address) (uint64, error) {
	return _ArbRollupCore.Contract.CurrentChallenge(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ExtraChallengeTimeBlocks(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "extraChallengeTimeBlocks")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ExtraChallengeTimeBlocks() (uint64, error) {
	return _ArbRollupCore.Contract.ExtraChallengeTimeBlocks(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ExtraChallengeTimeBlocks() (uint64, error) {
	return _ArbRollupCore.Contract.ExtraChallengeTimeBlocks(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) FirstUnresolvedNode(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "firstUnresolvedNode")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) FirstUnresolvedNode() (uint64, error) {
	return _ArbRollupCore.Contract.FirstUnresolvedNode(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) FirstUnresolvedNode() (uint64, error) {
	return _ArbRollupCore.Contract.FirstUnresolvedNode(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) GetNode(opts *bind.CallOpts, nodeNum uint64) (Node, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "getNode", nodeNum)

	if err != nil {
		return *new(Node), err
	}

	out0 := *abi.ConvertType(out[0], new(Node)).(*Node)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) GetNode(nodeNum uint64) (Node, error) {
	return _ArbRollupCore.Contract.GetNode(&_ArbRollupCore.CallOpts, nodeNum)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) GetNode(nodeNum uint64) (Node, error) {
	return _ArbRollupCore.Contract.GetNode(&_ArbRollupCore.CallOpts, nodeNum)
}

func (_ArbRollupCore *ArbRollupCoreCaller) GetNodeCreationBlockForLogLookup(opts *bind.CallOpts, nodeNum uint64) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "getNodeCreationBlockForLogLookup", nodeNum)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) GetNodeCreationBlockForLogLookup(nodeNum uint64) (*big.Int, error) {
	return _ArbRollupCore.Contract.GetNodeCreationBlockForLogLookup(&_ArbRollupCore.CallOpts, nodeNum)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) GetNodeCreationBlockForLogLookup(nodeNum uint64) (*big.Int, error) {
	return _ArbRollupCore.Contract.GetNodeCreationBlockForLogLookup(&_ArbRollupCore.CallOpts, nodeNum)
}

func (_ArbRollupCore *ArbRollupCoreCaller) GetStaker(opts *bind.CallOpts, staker common.Address) (IRollupCoreStaker, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "getStaker", staker)

	if err != nil {
		return *new(IRollupCoreStaker), err
	}

	out0 := *abi.ConvertType(out[0], new(IRollupCoreStaker)).(*IRollupCoreStaker)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) GetStaker(staker common.Address) (IRollupCoreStaker, error) {
	return _ArbRollupCore.Contract.GetStaker(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) GetStaker(staker common.Address) (IRollupCoreStaker, error) {
	return _ArbRollupCore.Contract.GetStaker(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) GetStakerAddress(opts *bind.CallOpts, stakerNum uint64) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "getStakerAddress", stakerNum)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) GetStakerAddress(stakerNum uint64) (common.Address, error) {
	return _ArbRollupCore.Contract.GetStakerAddress(&_ArbRollupCore.CallOpts, stakerNum)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) GetStakerAddress(stakerNum uint64) (common.Address, error) {
	return _ArbRollupCore.Contract.GetStakerAddress(&_ArbRollupCore.CallOpts, stakerNum)
}

func (_ArbRollupCore *ArbRollupCoreCaller) IsStaked(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "isStaked", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) IsStaked(staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsStaked(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) IsStaked(staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsStaked(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) IsValidator(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "isValidator", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) IsValidator(arg0 common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsValidator(&_ArbRollupCore.CallOpts, arg0)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) IsValidator(arg0 common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsValidator(&_ArbRollupCore.CallOpts, arg0)
}

func (_ArbRollupCore *ArbRollupCoreCaller) IsZombie(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "isZombie", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) IsZombie(staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsZombie(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) IsZombie(staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.IsZombie(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) LastStakeBlock(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "lastStakeBlock")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) LastStakeBlock() (uint64, error) {
	return _ArbRollupCore.Contract.LastStakeBlock(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) LastStakeBlock() (uint64, error) {
	return _ArbRollupCore.Contract.LastStakeBlock(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) LatestConfirmed(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "latestConfirmed")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) LatestConfirmed() (uint64, error) {
	return _ArbRollupCore.Contract.LatestConfirmed(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) LatestConfirmed() (uint64, error) {
	return _ArbRollupCore.Contract.LatestConfirmed(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) LatestNodeCreated(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "latestNodeCreated")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) LatestNodeCreated() (uint64, error) {
	return _ArbRollupCore.Contract.LatestNodeCreated(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) LatestNodeCreated() (uint64, error) {
	return _ArbRollupCore.Contract.LatestNodeCreated(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) LatestStakedNode(opts *bind.CallOpts, staker common.Address) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "latestStakedNode", staker)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) LatestStakedNode(staker common.Address) (uint64, error) {
	return _ArbRollupCore.Contract.LatestStakedNode(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) LatestStakedNode(staker common.Address) (uint64, error) {
	return _ArbRollupCore.Contract.LatestStakedNode(&_ArbRollupCore.CallOpts, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) LoserStakeEscrow(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "loserStakeEscrow")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) LoserStakeEscrow() (common.Address, error) {
	return _ArbRollupCore.Contract.LoserStakeEscrow(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) LoserStakeEscrow() (common.Address, error) {
	return _ArbRollupCore.Contract.LoserStakeEscrow(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) MinimumAssertionPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "minimumAssertionPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) MinimumAssertionPeriod() (*big.Int, error) {
	return _ArbRollupCore.Contract.MinimumAssertionPeriod(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) MinimumAssertionPeriod() (*big.Int, error) {
	return _ArbRollupCore.Contract.MinimumAssertionPeriod(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) NodeHasStaker(opts *bind.CallOpts, nodeNum uint64, staker common.Address) (bool, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "nodeHasStaker", nodeNum, staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) NodeHasStaker(nodeNum uint64, staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.NodeHasStaker(&_ArbRollupCore.CallOpts, nodeNum, staker)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) NodeHasStaker(nodeNum uint64, staker common.Address) (bool, error) {
	return _ArbRollupCore.Contract.NodeHasStaker(&_ArbRollupCore.CallOpts, nodeNum, staker)
}

func (_ArbRollupCore *ArbRollupCoreCaller) Outbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "outbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) Outbox() (common.Address, error) {
	return _ArbRollupCore.Contract.Outbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) Outbox() (common.Address, error) {
	return _ArbRollupCore.Contract.Outbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) RollupEventInbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "rollupEventInbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) RollupEventInbox() (common.Address, error) {
	return _ArbRollupCore.Contract.RollupEventInbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) RollupEventInbox() (common.Address, error) {
	return _ArbRollupCore.Contract.RollupEventInbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) SequencerInbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "sequencerInbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) SequencerInbox() (common.Address, error) {
	return _ArbRollupCore.Contract.SequencerInbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) SequencerInbox() (common.Address, error) {
	return _ArbRollupCore.Contract.SequencerInbox(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) StakeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "stakeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) StakeToken() (common.Address, error) {
	return _ArbRollupCore.Contract.StakeToken(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) StakeToken() (common.Address, error) {
	return _ArbRollupCore.Contract.StakeToken(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) StakerCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "stakerCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) StakerCount() (uint64, error) {
	return _ArbRollupCore.Contract.StakerCount(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) StakerCount() (uint64, error) {
	return _ArbRollupCore.Contract.StakerCount(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ValidatorWhitelistDisabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "validatorWhitelistDisabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ValidatorWhitelistDisabled() (bool, error) {
	return _ArbRollupCore.Contract.ValidatorWhitelistDisabled(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ValidatorWhitelistDisabled() (bool, error) {
	return _ArbRollupCore.Contract.ValidatorWhitelistDisabled(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) WasmModuleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "wasmModuleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) WasmModuleRoot() ([32]byte, error) {
	return _ArbRollupCore.Contract.WasmModuleRoot(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) WasmModuleRoot() ([32]byte, error) {
	return _ArbRollupCore.Contract.WasmModuleRoot(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) WithdrawableFunds(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "withdrawableFunds", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) WithdrawableFunds(owner common.Address) (*big.Int, error) {
	return _ArbRollupCore.Contract.WithdrawableFunds(&_ArbRollupCore.CallOpts, owner)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) WithdrawableFunds(owner common.Address) (*big.Int, error) {
	return _ArbRollupCore.Contract.WithdrawableFunds(&_ArbRollupCore.CallOpts, owner)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ZombieAddress(opts *bind.CallOpts, zombieNum *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "zombieAddress", zombieNum)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ZombieAddress(zombieNum *big.Int) (common.Address, error) {
	return _ArbRollupCore.Contract.ZombieAddress(&_ArbRollupCore.CallOpts, zombieNum)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ZombieAddress(zombieNum *big.Int) (common.Address, error) {
	return _ArbRollupCore.Contract.ZombieAddress(&_ArbRollupCore.CallOpts, zombieNum)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ZombieCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "zombieCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ZombieCount() (*big.Int, error) {
	return _ArbRollupCore.Contract.ZombieCount(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ZombieCount() (*big.Int, error) {
	return _ArbRollupCore.Contract.ZombieCount(&_ArbRollupCore.CallOpts)
}

func (_ArbRollupCore *ArbRollupCoreCaller) ZombieLatestStakedNode(opts *bind.CallOpts, zombieNum *big.Int) (uint64, error) {
	var out []interface{}
	err := _ArbRollupCore.contract.Call(opts, &out, "zombieLatestStakedNode", zombieNum)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_ArbRollupCore *ArbRollupCoreSession) ZombieLatestStakedNode(zombieNum *big.Int) (uint64, error) {
	return _ArbRollupCore.Contract.ZombieLatestStakedNode(&_ArbRollupCore.CallOpts, zombieNum)
}

func (_ArbRollupCore *ArbRollupCoreCallerSession) ZombieLatestStakedNode(zombieNum *big.Int) (uint64, error) {
	return _ArbRollupCore.Contract.ZombieLatestStakedNode(&_ArbRollupCore.CallOpts, zombieNum)
}

type ArbRollupCoreNodeConfirmedIterator struct {
	Event *ArbRollupCoreNodeConfirmed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreNodeConfirmedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreNodeConfirmed)
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
		it.Event = new(ArbRollupCoreNodeConfirmed)
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

func (it *ArbRollupCoreNodeConfirmedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreNodeConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreNodeConfirmed struct {
	NodeNum   uint64
	BlockHash [32]byte
	SendRoot  [32]byte
	Raw       types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterNodeConfirmed(opts *bind.FilterOpts, nodeNum []uint64) (*ArbRollupCoreNodeConfirmedIterator, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "NodeConfirmed", nodeNumRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreNodeConfirmedIterator{contract: _ArbRollupCore.contract, event: "NodeConfirmed", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchNodeConfirmed(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeConfirmed, nodeNum []uint64) (event.Subscription, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "NodeConfirmed", nodeNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreNodeConfirmed)
				if err := _ArbRollupCore.contract.UnpackLog(event, "NodeConfirmed", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseNodeConfirmed(log types.Log) (*ArbRollupCoreNodeConfirmed, error) {
	event := new(ArbRollupCoreNodeConfirmed)
	if err := _ArbRollupCore.contract.UnpackLog(event, "NodeConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreNodeCreatedIterator struct {
	Event *ArbRollupCoreNodeCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreNodeCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreNodeCreated)
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
		it.Event = new(ArbRollupCoreNodeCreated)
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

func (it *ArbRollupCoreNodeCreatedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreNodeCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreNodeCreated struct {
	NodeNum            uint64
	ParentNodeHash     [32]byte
	NodeHash           [32]byte
	ExecutionHash      [32]byte
	Assertion          Assertion
	AfterInboxBatchAcc [32]byte
	WasmModuleRoot     [32]byte
	InboxMaxCount      *big.Int
	Raw                types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterNodeCreated(opts *bind.FilterOpts, nodeNum []uint64, parentNodeHash [][32]byte, nodeHash [][32]byte) (*ArbRollupCoreNodeCreatedIterator, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}
	var parentNodeHashRule []interface{}
	for _, parentNodeHashItem := range parentNodeHash {
		parentNodeHashRule = append(parentNodeHashRule, parentNodeHashItem)
	}
	var nodeHashRule []interface{}
	for _, nodeHashItem := range nodeHash {
		nodeHashRule = append(nodeHashRule, nodeHashItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "NodeCreated", nodeNumRule, parentNodeHashRule, nodeHashRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreNodeCreatedIterator{contract: _ArbRollupCore.contract, event: "NodeCreated", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeCreated, nodeNum []uint64, parentNodeHash [][32]byte, nodeHash [][32]byte) (event.Subscription, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}
	var parentNodeHashRule []interface{}
	for _, parentNodeHashItem := range parentNodeHash {
		parentNodeHashRule = append(parentNodeHashRule, parentNodeHashItem)
	}
	var nodeHashRule []interface{}
	for _, nodeHashItem := range nodeHash {
		nodeHashRule = append(nodeHashRule, nodeHashItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "NodeCreated", nodeNumRule, parentNodeHashRule, nodeHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreNodeCreated)
				if err := _ArbRollupCore.contract.UnpackLog(event, "NodeCreated", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseNodeCreated(log types.Log) (*ArbRollupCoreNodeCreated, error) {
	event := new(ArbRollupCoreNodeCreated)
	if err := _ArbRollupCore.contract.UnpackLog(event, "NodeCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreNodeRejectedIterator struct {
	Event *ArbRollupCoreNodeRejected

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreNodeRejectedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreNodeRejected)
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
		it.Event = new(ArbRollupCoreNodeRejected)
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

func (it *ArbRollupCoreNodeRejectedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreNodeRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreNodeRejected struct {
	NodeNum uint64
	Raw     types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterNodeRejected(opts *bind.FilterOpts, nodeNum []uint64) (*ArbRollupCoreNodeRejectedIterator, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "NodeRejected", nodeNumRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreNodeRejectedIterator{contract: _ArbRollupCore.contract, event: "NodeRejected", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchNodeRejected(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeRejected, nodeNum []uint64) (event.Subscription, error) {

	var nodeNumRule []interface{}
	for _, nodeNumItem := range nodeNum {
		nodeNumRule = append(nodeNumRule, nodeNumItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "NodeRejected", nodeNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreNodeRejected)
				if err := _ArbRollupCore.contract.UnpackLog(event, "NodeRejected", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseNodeRejected(log types.Log) (*ArbRollupCoreNodeRejected, error) {
	event := new(ArbRollupCoreNodeRejected)
	if err := _ArbRollupCore.contract.UnpackLog(event, "NodeRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreRollupChallengeStartedIterator struct {
	Event *ArbRollupCoreRollupChallengeStarted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreRollupChallengeStartedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreRollupChallengeStarted)
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
		it.Event = new(ArbRollupCoreRollupChallengeStarted)
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

func (it *ArbRollupCoreRollupChallengeStartedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreRollupChallengeStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreRollupChallengeStarted struct {
	ChallengeIndex uint64
	Asserter       common.Address
	Challenger     common.Address
	ChallengedNode uint64
	Raw            types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterRollupChallengeStarted(opts *bind.FilterOpts, challengeIndex []uint64) (*ArbRollupCoreRollupChallengeStartedIterator, error) {

	var challengeIndexRule []interface{}
	for _, challengeIndexItem := range challengeIndex {
		challengeIndexRule = append(challengeIndexRule, challengeIndexItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "RollupChallengeStarted", challengeIndexRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreRollupChallengeStartedIterator{contract: _ArbRollupCore.contract, event: "RollupChallengeStarted", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchRollupChallengeStarted(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreRollupChallengeStarted, challengeIndex []uint64) (event.Subscription, error) {

	var challengeIndexRule []interface{}
	for _, challengeIndexItem := range challengeIndex {
		challengeIndexRule = append(challengeIndexRule, challengeIndexItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "RollupChallengeStarted", challengeIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreRollupChallengeStarted)
				if err := _ArbRollupCore.contract.UnpackLog(event, "RollupChallengeStarted", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseRollupChallengeStarted(log types.Log) (*ArbRollupCoreRollupChallengeStarted, error) {
	event := new(ArbRollupCoreRollupChallengeStarted)
	if err := _ArbRollupCore.contract.UnpackLog(event, "RollupChallengeStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreRollupInitializedIterator struct {
	Event *ArbRollupCoreRollupInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreRollupInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreRollupInitialized)
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
		it.Event = new(ArbRollupCoreRollupInitialized)
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

func (it *ArbRollupCoreRollupInitializedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreRollupInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreRollupInitialized struct {
	MachineHash [32]byte
	ChainId     *big.Int
	Raw         types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterRollupInitialized(opts *bind.FilterOpts) (*ArbRollupCoreRollupInitializedIterator, error) {

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "RollupInitialized")
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreRollupInitializedIterator{contract: _ArbRollupCore.contract, event: "RollupInitialized", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchRollupInitialized(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreRollupInitialized) (event.Subscription, error) {

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "RollupInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreRollupInitialized)
				if err := _ArbRollupCore.contract.UnpackLog(event, "RollupInitialized", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseRollupInitialized(log types.Log) (*ArbRollupCoreRollupInitialized, error) {
	event := new(ArbRollupCoreRollupInitialized)
	if err := _ArbRollupCore.contract.UnpackLog(event, "RollupInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreUserStakeUpdatedIterator struct {
	Event *ArbRollupCoreUserStakeUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreUserStakeUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreUserStakeUpdated)
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
		it.Event = new(ArbRollupCoreUserStakeUpdated)
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

func (it *ArbRollupCoreUserStakeUpdatedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreUserStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreUserStakeUpdated struct {
	User           common.Address
	InitialBalance *big.Int
	FinalBalance   *big.Int
	Raw            types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterUserStakeUpdated(opts *bind.FilterOpts, user []common.Address) (*ArbRollupCoreUserStakeUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "UserStakeUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreUserStakeUpdatedIterator{contract: _ArbRollupCore.contract, event: "UserStakeUpdated", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchUserStakeUpdated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreUserStakeUpdated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "UserStakeUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreUserStakeUpdated)
				if err := _ArbRollupCore.contract.UnpackLog(event, "UserStakeUpdated", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseUserStakeUpdated(log types.Log) (*ArbRollupCoreUserStakeUpdated, error) {
	event := new(ArbRollupCoreUserStakeUpdated)
	if err := _ArbRollupCore.contract.UnpackLog(event, "UserStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbRollupCoreUserWithdrawableFundsUpdatedIterator struct {
	Event *ArbRollupCoreUserWithdrawableFundsUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbRollupCoreUserWithdrawableFundsUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbRollupCoreUserWithdrawableFundsUpdated)
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
		it.Event = new(ArbRollupCoreUserWithdrawableFundsUpdated)
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

func (it *ArbRollupCoreUserWithdrawableFundsUpdatedIterator) Error() error {
	return it.fail
}

func (it *ArbRollupCoreUserWithdrawableFundsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbRollupCoreUserWithdrawableFundsUpdated struct {
	User           common.Address
	InitialBalance *big.Int
	FinalBalance   *big.Int
	Raw            types.Log
}

func (_ArbRollupCore *ArbRollupCoreFilterer) FilterUserWithdrawableFundsUpdated(opts *bind.FilterOpts, user []common.Address) (*ArbRollupCoreUserWithdrawableFundsUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArbRollupCore.contract.FilterLogs(opts, "UserWithdrawableFundsUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return &ArbRollupCoreUserWithdrawableFundsUpdatedIterator{contract: _ArbRollupCore.contract, event: "UserWithdrawableFundsUpdated", logs: logs, sub: sub}, nil
}

func (_ArbRollupCore *ArbRollupCoreFilterer) WatchUserWithdrawableFundsUpdated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreUserWithdrawableFundsUpdated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArbRollupCore.contract.WatchLogs(opts, "UserWithdrawableFundsUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbRollupCoreUserWithdrawableFundsUpdated)
				if err := _ArbRollupCore.contract.UnpackLog(event, "UserWithdrawableFundsUpdated", log); err != nil {
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

func (_ArbRollupCore *ArbRollupCoreFilterer) ParseUserWithdrawableFundsUpdated(log types.Log) (*ArbRollupCoreUserWithdrawableFundsUpdated, error) {
	event := new(ArbRollupCoreUserWithdrawableFundsUpdated)
	if err := _ArbRollupCore.contract.UnpackLog(event, "UserWithdrawableFundsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ArbRollupCore *ArbRollupCore) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ArbRollupCore.abi.Events["NodeConfirmed"].ID:
		return _ArbRollupCore.ParseNodeConfirmed(log)
	case _ArbRollupCore.abi.Events["NodeCreated"].ID:
		return _ArbRollupCore.ParseNodeCreated(log)
	case _ArbRollupCore.abi.Events["NodeRejected"].ID:
		return _ArbRollupCore.ParseNodeRejected(log)
	case _ArbRollupCore.abi.Events["RollupChallengeStarted"].ID:
		return _ArbRollupCore.ParseRollupChallengeStarted(log)
	case _ArbRollupCore.abi.Events["RollupInitialized"].ID:
		return _ArbRollupCore.ParseRollupInitialized(log)
	case _ArbRollupCore.abi.Events["UserStakeUpdated"].ID:
		return _ArbRollupCore.ParseUserStakeUpdated(log)
	case _ArbRollupCore.abi.Events["UserWithdrawableFundsUpdated"].ID:
		return _ArbRollupCore.ParseUserWithdrawableFundsUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ArbRollupCoreNodeConfirmed) Topic() common.Hash {
	return common.HexToHash("0x22ef0479a7ff660660d1c2fe35f1b632cf31675c2d9378db8cec95b00d8ffa3c")
}

func (ArbRollupCoreNodeCreated) Topic() common.Hash {
	return common.HexToHash("0x4f4caa9e67fb994e349dd35d1ad0ce23053d4323f83ce11dc817b5435031d096")
}

func (ArbRollupCoreNodeRejected) Topic() common.Hash {
	return common.HexToHash("0xeaffa3d968707ec919a2fc9f31d5ab2b86c905881ff561725d5a82fc95ad4640")
}

func (ArbRollupCoreRollupChallengeStarted) Topic() common.Hash {
	return common.HexToHash("0x6db7dc2f507647d135035469b27aa79cea90582779d084a7821d6cd092cbd873")
}

func (ArbRollupCoreRollupInitialized) Topic() common.Hash {
	return common.HexToHash("0xfc1b83c11d99d08a938e0b82a0bd45f822f71ff5abf23f999c93c4533d752464")
}

func (ArbRollupCoreUserStakeUpdated) Topic() common.Hash {
	return common.HexToHash("0xebd093d389ab57f3566918d2c379a2b4d9539e8eb95efad9d5e465457833fde6")
}

func (ArbRollupCoreUserWithdrawableFundsUpdated) Topic() common.Hash {
	return common.HexToHash("0xa740af14c56e4e04a617b1de1eb20de73270decbaaead14f142aabf3038e5ae2")
}

func (_ArbRollupCore *ArbRollupCore) Address() common.Address {
	return _ArbRollupCore.address
}

type ArbRollupCoreInterface interface {
	AmountStaked(opts *bind.CallOpts, staker common.Address) (*big.Int, error)

	BaseStake(opts *bind.CallOpts) (*big.Int, error)

	Bridge(opts *bind.CallOpts) (common.Address, error)

	ChainId(opts *bind.CallOpts) (*big.Int, error)

	ChallengeManager(opts *bind.CallOpts) (common.Address, error)

	ConfirmPeriodBlocks(opts *bind.CallOpts) (uint64, error)

	CurrentChallenge(opts *bind.CallOpts, staker common.Address) (uint64, error)

	ExtraChallengeTimeBlocks(opts *bind.CallOpts) (uint64, error)

	FirstUnresolvedNode(opts *bind.CallOpts) (uint64, error)

	GetNode(opts *bind.CallOpts, nodeNum uint64) (Node, error)

	GetNodeCreationBlockForLogLookup(opts *bind.CallOpts, nodeNum uint64) (*big.Int, error)

	GetStaker(opts *bind.CallOpts, staker common.Address) (IRollupCoreStaker, error)

	GetStakerAddress(opts *bind.CallOpts, stakerNum uint64) (common.Address, error)

	IsStaked(opts *bind.CallOpts, staker common.Address) (bool, error)

	IsValidator(opts *bind.CallOpts, arg0 common.Address) (bool, error)

	IsZombie(opts *bind.CallOpts, staker common.Address) (bool, error)

	LastStakeBlock(opts *bind.CallOpts) (uint64, error)

	LatestConfirmed(opts *bind.CallOpts) (uint64, error)

	LatestNodeCreated(opts *bind.CallOpts) (uint64, error)

	LatestStakedNode(opts *bind.CallOpts, staker common.Address) (uint64, error)

	LoserStakeEscrow(opts *bind.CallOpts) (common.Address, error)

	MinimumAssertionPeriod(opts *bind.CallOpts) (*big.Int, error)

	NodeHasStaker(opts *bind.CallOpts, nodeNum uint64, staker common.Address) (bool, error)

	Outbox(opts *bind.CallOpts) (common.Address, error)

	RollupEventInbox(opts *bind.CallOpts) (common.Address, error)

	SequencerInbox(opts *bind.CallOpts) (common.Address, error)

	StakeToken(opts *bind.CallOpts) (common.Address, error)

	StakerCount(opts *bind.CallOpts) (uint64, error)

	ValidatorWhitelistDisabled(opts *bind.CallOpts) (bool, error)

	WasmModuleRoot(opts *bind.CallOpts) ([32]byte, error)

	WithdrawableFunds(opts *bind.CallOpts, owner common.Address) (*big.Int, error)

	ZombieAddress(opts *bind.CallOpts, zombieNum *big.Int) (common.Address, error)

	ZombieCount(opts *bind.CallOpts) (*big.Int, error)

	ZombieLatestStakedNode(opts *bind.CallOpts, zombieNum *big.Int) (uint64, error)

	FilterNodeConfirmed(opts *bind.FilterOpts, nodeNum []uint64) (*ArbRollupCoreNodeConfirmedIterator, error)

	WatchNodeConfirmed(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeConfirmed, nodeNum []uint64) (event.Subscription, error)

	ParseNodeConfirmed(log types.Log) (*ArbRollupCoreNodeConfirmed, error)

	FilterNodeCreated(opts *bind.FilterOpts, nodeNum []uint64, parentNodeHash [][32]byte, nodeHash [][32]byte) (*ArbRollupCoreNodeCreatedIterator, error)

	WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeCreated, nodeNum []uint64, parentNodeHash [][32]byte, nodeHash [][32]byte) (event.Subscription, error)

	ParseNodeCreated(log types.Log) (*ArbRollupCoreNodeCreated, error)

	FilterNodeRejected(opts *bind.FilterOpts, nodeNum []uint64) (*ArbRollupCoreNodeRejectedIterator, error)

	WatchNodeRejected(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreNodeRejected, nodeNum []uint64) (event.Subscription, error)

	ParseNodeRejected(log types.Log) (*ArbRollupCoreNodeRejected, error)

	FilterRollupChallengeStarted(opts *bind.FilterOpts, challengeIndex []uint64) (*ArbRollupCoreRollupChallengeStartedIterator, error)

	WatchRollupChallengeStarted(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreRollupChallengeStarted, challengeIndex []uint64) (event.Subscription, error)

	ParseRollupChallengeStarted(log types.Log) (*ArbRollupCoreRollupChallengeStarted, error)

	FilterRollupInitialized(opts *bind.FilterOpts) (*ArbRollupCoreRollupInitializedIterator, error)

	WatchRollupInitialized(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreRollupInitialized) (event.Subscription, error)

	ParseRollupInitialized(log types.Log) (*ArbRollupCoreRollupInitialized, error)

	FilterUserStakeUpdated(opts *bind.FilterOpts, user []common.Address) (*ArbRollupCoreUserStakeUpdatedIterator, error)

	WatchUserStakeUpdated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreUserStakeUpdated, user []common.Address) (event.Subscription, error)

	ParseUserStakeUpdated(log types.Log) (*ArbRollupCoreUserStakeUpdated, error)

	FilterUserWithdrawableFundsUpdated(opts *bind.FilterOpts, user []common.Address) (*ArbRollupCoreUserWithdrawableFundsUpdatedIterator, error)

	WatchUserWithdrawableFundsUpdated(opts *bind.WatchOpts, sink chan<- *ArbRollupCoreUserWithdrawableFundsUpdated, user []common.Address) (event.Subscription, error)

	ParseUserWithdrawableFundsUpdated(log types.Log) (*ArbRollupCoreUserWithdrawableFundsUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
