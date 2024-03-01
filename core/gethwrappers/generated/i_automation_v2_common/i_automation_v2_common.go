// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_automation_v2_common

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

type IAutomationV2CommonState struct {
	Nonce                   uint32
	OwnerLinkBalance        *big.Int
	ExpectedLinkBalance     *big.Int
	TotalPremium            *big.Int
	NumUpkeeps              *big.Int
	ConfigCount             uint32
	LatestConfigBlockNumber uint32
	LatestConfigDigest      [32]byte
	LatestEpoch             uint32
	Paused                  bool
}

type IAutomationV2CommonUpkeepInfo struct {
	Target                   common.Address
	PerformGas               uint32
	CheckData                []byte
	Balance                  *big.Int
	Admin                    common.Address
	MaxValidBlocknumber      uint64
	LastPerformedBlockNumber uint32
	AmountSpent              *big.Int
	Paused                   bool
	OffchainConfig           []byte
}

type OnchainConfigLegacy struct {
	PaymentPremiumPPB      uint32
	FlatFeeMicroLink       uint32
	CheckGasLimit          uint32
	StalenessSeconds       *big.Int
	GasCeilingMultiplier   uint16
	MinUpkeepSpend         *big.Int
	MaxPerformGas          uint32
	MaxCheckDataSize       uint32
	MaxPerformDataSize     uint32
	MaxRevertDataSize      uint32
	FallbackGasPrice       *big.Int
	FallbackLinkPrice      *big.Int
	Transcoder             common.Address
	Registrars             []common.Address
	UpkeepPrivilegeManager common.Address
}

var IAutomationV2CommonMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV2Common.State\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV2Common.UpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var IAutomationV2CommonABI = IAutomationV2CommonMetaData.ABI

type IAutomationV2Common struct {
	address common.Address
	abi     abi.ABI
	IAutomationV2CommonCaller
	IAutomationV2CommonTransactor
	IAutomationV2CommonFilterer
}

type IAutomationV2CommonCaller struct {
	contract *bind.BoundContract
}

type IAutomationV2CommonTransactor struct {
	contract *bind.BoundContract
}

type IAutomationV2CommonFilterer struct {
	contract *bind.BoundContract
}

type IAutomationV2CommonSession struct {
	Contract     *IAutomationV2Common
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAutomationV2CommonCallerSession struct {
	Contract *IAutomationV2CommonCaller
	CallOpts bind.CallOpts
}

type IAutomationV2CommonTransactorSession struct {
	Contract     *IAutomationV2CommonTransactor
	TransactOpts bind.TransactOpts
}

type IAutomationV2CommonRaw struct {
	Contract *IAutomationV2Common
}

type IAutomationV2CommonCallerRaw struct {
	Contract *IAutomationV2CommonCaller
}

type IAutomationV2CommonTransactorRaw struct {
	Contract *IAutomationV2CommonTransactor
}

func NewIAutomationV2Common(address common.Address, backend bind.ContractBackend) (*IAutomationV2Common, error) {
	abi, err := abi.JSON(strings.NewReader(IAutomationV2CommonABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAutomationV2Common(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2Common{address: address, abi: abi, IAutomationV2CommonCaller: IAutomationV2CommonCaller{contract: contract}, IAutomationV2CommonTransactor: IAutomationV2CommonTransactor{contract: contract}, IAutomationV2CommonFilterer: IAutomationV2CommonFilterer{contract: contract}}, nil
}

func NewIAutomationV2CommonCaller(address common.Address, caller bind.ContractCaller) (*IAutomationV2CommonCaller, error) {
	contract, err := bindIAutomationV2Common(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonCaller{contract: contract}, nil
}

func NewIAutomationV2CommonTransactor(address common.Address, transactor bind.ContractTransactor) (*IAutomationV2CommonTransactor, error) {
	contract, err := bindIAutomationV2Common(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonTransactor{contract: contract}, nil
}

func NewIAutomationV2CommonFilterer(address common.Address, filterer bind.ContractFilterer) (*IAutomationV2CommonFilterer, error) {
	contract, err := bindIAutomationV2Common(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonFilterer{contract: contract}, nil
}

func bindIAutomationV2Common(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAutomationV2CommonMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IAutomationV2Common *IAutomationV2CommonRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationV2Common.Contract.IAutomationV2CommonCaller.contract.Call(opts, result, method, params...)
}

func (_IAutomationV2Common *IAutomationV2CommonRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.IAutomationV2CommonTransactor.contract.Transfer(opts)
}

func (_IAutomationV2Common *IAutomationV2CommonRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.IAutomationV2CommonTransactor.contract.Transact(opts, method, params...)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationV2Common.Contract.contract.Call(opts, result, method, params...)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.contract.Transfer(opts)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.contract.Transact(opts, method, params...)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "checkCallback", id, values, extraData)

	outstruct := new(CheckCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.UpkeepFailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationV2Common.Contract.CheckCallback(&_IAutomationV2Common.CallOpts, id, values, extraData)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationV2Common.Contract.CheckCallback(&_IAutomationV2Common.CallOpts, id, values, extraData)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "checkUpkeep", id, triggerData)

	outstruct := new(CheckUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.UpkeepFailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GasLimit = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FastGasWei = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.LinkNative = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationV2Common.Contract.CheckUpkeep(&_IAutomationV2Common.CallOpts, id, triggerData)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationV2Common.Contract.CheckUpkeep(&_IAutomationV2Common.CallOpts, id, triggerData)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

	error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "checkUpkeep0", id)

	outstruct := new(CheckUpkeep0)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.UpkeepFailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GasLimit = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FastGasWei = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.LinkNative = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationV2Common.Contract.CheckUpkeep0(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationV2Common.Contract.CheckUpkeep0(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationV2Common.Contract.GetActiveUpkeepIDs(&_IAutomationV2Common.CallOpts, startIndex, maxCount)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationV2Common.Contract.GetActiveUpkeepIDs(&_IAutomationV2Common.CallOpts, startIndex, maxCount)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationV2Common.Contract.GetMinBalance(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationV2Common.Contract.GetMinBalance(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(IAutomationV2CommonState)).(*IAutomationV2CommonState)
	outstruct.Config = *abi.ConvertType(out[1], new(OnchainConfigLegacy)).(*OnchainConfigLegacy)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetState() (GetState,

	error) {
	return _IAutomationV2Common.Contract.GetState(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetState() (GetState,

	error) {
	return _IAutomationV2Common.Contract.GetState(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationV2Common.Contract.GetTriggerType(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationV2Common.Contract.GetTriggerType(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV2CommonUpkeepInfo, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV2CommonUpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV2CommonUpkeepInfo)).(*IAutomationV2CommonUpkeepInfo)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetUpkeep(id *big.Int) (IAutomationV2CommonUpkeepInfo, error) {
	return _IAutomationV2Common.Contract.GetUpkeep(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetUpkeep(id *big.Int) (IAutomationV2CommonUpkeepInfo, error) {
	return _IAutomationV2Common.Contract.GetUpkeep(&_IAutomationV2Common.CallOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV2Common.Contract.GetUpkeepPrivilegeConfig(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV2Common.Contract.GetUpkeepPrivilegeConfig(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV2Common.Contract.GetUpkeepTriggerConfig(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV2Common.Contract.GetUpkeepTriggerConfig(&_IAutomationV2Common.CallOpts, upkeepId)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationV2Common.Contract.HasDedupKey(&_IAutomationV2Common.CallOpts, dedupKey)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationV2Common.Contract.HasDedupKey(&_IAutomationV2Common.CallOpts, dedupKey)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) Owner() (common.Address, error) {
	return _IAutomationV2Common.Contract.Owner(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) Owner() (common.Address, error) {
	return _IAutomationV2Common.Contract.Owner(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "simulatePerformUpkeep", id, performData)

	outstruct := new(SimulatePerformUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GasUsed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationV2Common.Contract.SimulatePerformUpkeep(&_IAutomationV2Common.CallOpts, id, performData)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationV2Common.Contract.SimulatePerformUpkeep(&_IAutomationV2Common.CallOpts, id, performData)
}

func (_IAutomationV2Common *IAutomationV2CommonCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAutomationV2Common.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAutomationV2Common *IAutomationV2CommonSession) TypeAndVersion() (string, error) {
	return _IAutomationV2Common.Contract.TypeAndVersion(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonCallerSession) TypeAndVersion() (string, error) {
	return _IAutomationV2Common.Contract.TypeAndVersion(&_IAutomationV2Common.CallOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "addFunds", id, amount)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.AddFunds(&_IAutomationV2Common.TransactOpts, id, amount)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.AddFunds(&_IAutomationV2Common.TransactOpts, id, amount)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.CancelUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.CancelUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.ExecuteCallback(&_IAutomationV2Common.TransactOpts, id, payload)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.ExecuteCallback(&_IAutomationV2Common.TransactOpts, id, payload)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "pause")
}

func (_IAutomationV2Common *IAutomationV2CommonSession) Pause() (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.Pause(&_IAutomationV2Common.TransactOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) Pause() (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.Pause(&_IAutomationV2Common.TransactOpts)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.PauseUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.PauseUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.RegisterUpkeep(&_IAutomationV2Common.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.RegisterUpkeep(&_IAutomationV2Common.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetConfig(&_IAutomationV2Common.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetConfig(&_IAutomationV2Common.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepCheckData(&_IAutomationV2Common.TransactOpts, id, newCheckData)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepCheckData(&_IAutomationV2Common.TransactOpts, id, newCheckData)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepGasLimit(&_IAutomationV2Common.TransactOpts, id, gasLimit)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepGasLimit(&_IAutomationV2Common.TransactOpts, id, gasLimit)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepPrivilegeConfig(&_IAutomationV2Common.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepPrivilegeConfig(&_IAutomationV2Common.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepTriggerConfig(&_IAutomationV2Common.TransactOpts, id, triggerConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.SetUpkeepTriggerConfig(&_IAutomationV2Common.TransactOpts, id, triggerConfig)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IAutomationV2Common *IAutomationV2CommonSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.UnpauseUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

func (_IAutomationV2Common *IAutomationV2CommonTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV2Common.Contract.UnpauseUpkeep(&_IAutomationV2Common.TransactOpts, id)
}

type IAutomationV2CommonDedupKeyAddedIterator struct {
	Event *IAutomationV2CommonDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonDedupKeyAdded)
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
		it.Event = new(IAutomationV2CommonDedupKeyAdded)
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

func (it *IAutomationV2CommonDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationV2CommonDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonDedupKeyAddedIterator{contract: _IAutomationV2Common.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonDedupKeyAdded)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseDedupKeyAdded(log types.Log) (*IAutomationV2CommonDedupKeyAdded, error) {
	event := new(IAutomationV2CommonDedupKeyAdded)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonInsufficientFundsUpkeepReportIterator struct {
	Event *IAutomationV2CommonInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonInsufficientFundsUpkeepReport)
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
		it.Event = new(IAutomationV2CommonInsufficientFundsUpkeepReport)
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

func (it *IAutomationV2CommonInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonInsufficientFundsUpkeepReportIterator{contract: _IAutomationV2Common.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonInsufficientFundsUpkeepReport)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationV2CommonInsufficientFundsUpkeepReport, error) {
	event := new(IAutomationV2CommonInsufficientFundsUpkeepReport)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonReorgedUpkeepReportIterator struct {
	Event *IAutomationV2CommonReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonReorgedUpkeepReport)
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
		it.Event = new(IAutomationV2CommonReorgedUpkeepReport)
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

func (it *IAutomationV2CommonReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonReorgedUpkeepReportIterator{contract: _IAutomationV2Common.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonReorgedUpkeepReport)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseReorgedUpkeepReport(log types.Log) (*IAutomationV2CommonReorgedUpkeepReport, error) {
	event := new(IAutomationV2CommonReorgedUpkeepReport)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonStaleUpkeepReportIterator struct {
	Event *IAutomationV2CommonStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonStaleUpkeepReport)
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
		it.Event = new(IAutomationV2CommonStaleUpkeepReport)
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

func (it *IAutomationV2CommonStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonStaleUpkeepReportIterator{contract: _IAutomationV2Common.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonStaleUpkeepReport)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseStaleUpkeepReport(log types.Log) (*IAutomationV2CommonStaleUpkeepReport, error) {
	event := new(IAutomationV2CommonStaleUpkeepReport)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepCanceledIterator struct {
	Event *IAutomationV2CommonUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepCanceled)
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
		it.Event = new(IAutomationV2CommonUpkeepCanceled)
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

func (it *IAutomationV2CommonUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationV2CommonUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepCanceledIterator{contract: _IAutomationV2Common.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepCanceled)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepCanceled(log types.Log) (*IAutomationV2CommonUpkeepCanceled, error) {
	event := new(IAutomationV2CommonUpkeepCanceled)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepMigratedIterator struct {
	Event *IAutomationV2CommonUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepMigrated)
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
		it.Event = new(IAutomationV2CommonUpkeepMigrated)
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

func (it *IAutomationV2CommonUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepMigratedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepMigrated)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepMigrated(log types.Log) (*IAutomationV2CommonUpkeepMigrated, error) {
	event := new(IAutomationV2CommonUpkeepMigrated)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepPausedIterator struct {
	Event *IAutomationV2CommonUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepPaused)
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
		it.Event = new(IAutomationV2CommonUpkeepPaused)
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

func (it *IAutomationV2CommonUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepPausedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepPaused)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepPaused(log types.Log) (*IAutomationV2CommonUpkeepPaused, error) {
	event := new(IAutomationV2CommonUpkeepPaused)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepPerformedIterator struct {
	Event *IAutomationV2CommonUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepPerformed)
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
		it.Event = new(IAutomationV2CommonUpkeepPerformed)
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

func (it *IAutomationV2CommonUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationV2CommonUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepPerformedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepPerformed)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepPerformed(log types.Log) (*IAutomationV2CommonUpkeepPerformed, error) {
	event := new(IAutomationV2CommonUpkeepPerformed)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepReceivedIterator struct {
	Event *IAutomationV2CommonUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepReceived)
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
		it.Event = new(IAutomationV2CommonUpkeepReceived)
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

func (it *IAutomationV2CommonUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepReceivedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepReceived)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepReceived(log types.Log) (*IAutomationV2CommonUpkeepReceived, error) {
	event := new(IAutomationV2CommonUpkeepReceived)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepRegisteredIterator struct {
	Event *IAutomationV2CommonUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepRegistered)
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
		it.Event = new(IAutomationV2CommonUpkeepRegistered)
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

func (it *IAutomationV2CommonUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepRegisteredIterator{contract: _IAutomationV2Common.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepRegistered)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepRegistered(log types.Log) (*IAutomationV2CommonUpkeepRegistered, error) {
	event := new(IAutomationV2CommonUpkeepRegistered)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepTriggerConfigSetIterator struct {
	Event *IAutomationV2CommonUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepTriggerConfigSet)
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
		it.Event = new(IAutomationV2CommonUpkeepTriggerConfigSet)
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

func (it *IAutomationV2CommonUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepTriggerConfigSetIterator{contract: _IAutomationV2Common.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepTriggerConfigSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationV2CommonUpkeepTriggerConfigSet, error) {
	event := new(IAutomationV2CommonUpkeepTriggerConfigSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepUnpausedIterator struct {
	Event *IAutomationV2CommonUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepUnpaused)
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
		it.Event = new(IAutomationV2CommonUpkeepUnpaused)
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

func (it *IAutomationV2CommonUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepUnpausedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepUnpaused)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepUnpaused(log types.Log) (*IAutomationV2CommonUpkeepUnpaused, error) {
	event := new(IAutomationV2CommonUpkeepUnpaused)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckCallback struct {
	UpkeepNeeded        bool
	PerformData         []byte
	UpkeepFailureReason uint8
	GasUsed             *big.Int
}
type CheckUpkeep struct {
	UpkeepNeeded        bool
	PerformData         []byte
	UpkeepFailureReason uint8
	GasUsed             *big.Int
	GasLimit            *big.Int
	FastGasWei          *big.Int
	LinkNative          *big.Int
}
type CheckUpkeep0 struct {
	UpkeepNeeded        bool
	PerformData         []byte
	UpkeepFailureReason uint8
	GasUsed             *big.Int
	GasLimit            *big.Int
	FastGasWei          *big.Int
	LinkNative          *big.Int
}
type GetState struct {
	State        IAutomationV2CommonState
	Config       OnchainConfigLegacy
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}
type SimulatePerformUpkeep struct {
	Success bool
	GasUsed *big.Int
}

func (_IAutomationV2Common *IAutomationV2Common) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IAutomationV2Common.abi.Events["DedupKeyAdded"].ID:
		return _IAutomationV2Common.ParseDedupKeyAdded(log)
	case _IAutomationV2Common.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IAutomationV2Common.ParseInsufficientFundsUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["ReorgedUpkeepReport"].ID:
		return _IAutomationV2Common.ParseReorgedUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["StaleUpkeepReport"].ID:
		return _IAutomationV2Common.ParseStaleUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["UpkeepCanceled"].ID:
		return _IAutomationV2Common.ParseUpkeepCanceled(log)
	case _IAutomationV2Common.abi.Events["UpkeepMigrated"].ID:
		return _IAutomationV2Common.ParseUpkeepMigrated(log)
	case _IAutomationV2Common.abi.Events["UpkeepPaused"].ID:
		return _IAutomationV2Common.ParseUpkeepPaused(log)
	case _IAutomationV2Common.abi.Events["UpkeepPerformed"].ID:
		return _IAutomationV2Common.ParseUpkeepPerformed(log)
	case _IAutomationV2Common.abi.Events["UpkeepReceived"].ID:
		return _IAutomationV2Common.ParseUpkeepReceived(log)
	case _IAutomationV2Common.abi.Events["UpkeepRegistered"].ID:
		return _IAutomationV2Common.ParseUpkeepRegistered(log)
	case _IAutomationV2Common.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IAutomationV2Common.ParseUpkeepTriggerConfigSet(log)
	case _IAutomationV2Common.abi.Events["UpkeepUnpaused"].ID:
		return _IAutomationV2Common.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IAutomationV2CommonDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IAutomationV2CommonInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IAutomationV2CommonReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IAutomationV2CommonStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IAutomationV2CommonUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IAutomationV2CommonUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IAutomationV2CommonUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IAutomationV2CommonUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IAutomationV2CommonUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IAutomationV2CommonUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IAutomationV2CommonUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IAutomationV2CommonUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IAutomationV2Common *IAutomationV2Common) Address() common.Address {
	return _IAutomationV2Common.address
}

type IAutomationV2CommonInterface interface {
	CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

		error)

	CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

		error)

	CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

		error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV2CommonUpkeepInfo, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

		error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationV2CommonDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IAutomationV2CommonDedupKeyAdded, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationV2CommonInsufficientFundsUpkeepReport, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IAutomationV2CommonReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IAutomationV2CommonStaleUpkeepReport, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationV2CommonUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IAutomationV2CommonUpkeepCanceled, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IAutomationV2CommonUpkeepMigrated, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IAutomationV2CommonUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationV2CommonUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IAutomationV2CommonUpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IAutomationV2CommonUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IAutomationV2CommonUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationV2CommonUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IAutomationV2CommonUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
