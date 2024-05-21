// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_automation_v21_plus_common

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

type IAutomationV21PlusCommonOnchainConfigLegacy struct {
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

type IAutomationV21PlusCommonStateLegacy struct {
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

type IAutomationV21PlusCommonUpkeepInfoLegacy struct {
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

var IAutomationV21PlusCommonMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var IAutomationV21PlusCommonABI = IAutomationV21PlusCommonMetaData.ABI

type IAutomationV21PlusCommon struct {
	address common.Address
	abi     abi.ABI
	IAutomationV21PlusCommonCaller
	IAutomationV21PlusCommonTransactor
	IAutomationV21PlusCommonFilterer
}

type IAutomationV21PlusCommonCaller struct {
	contract *bind.BoundContract
}

type IAutomationV21PlusCommonTransactor struct {
	contract *bind.BoundContract
}

type IAutomationV21PlusCommonFilterer struct {
	contract *bind.BoundContract
}

type IAutomationV21PlusCommonSession struct {
	Contract     *IAutomationV21PlusCommon
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAutomationV21PlusCommonCallerSession struct {
	Contract *IAutomationV21PlusCommonCaller
	CallOpts bind.CallOpts
}

type IAutomationV21PlusCommonTransactorSession struct {
	Contract     *IAutomationV21PlusCommonTransactor
	TransactOpts bind.TransactOpts
}

type IAutomationV21PlusCommonRaw struct {
	Contract *IAutomationV21PlusCommon
}

type IAutomationV21PlusCommonCallerRaw struct {
	Contract *IAutomationV21PlusCommonCaller
}

type IAutomationV21PlusCommonTransactorRaw struct {
	Contract *IAutomationV21PlusCommonTransactor
}

func NewIAutomationV21PlusCommon(address common.Address, backend bind.ContractBackend) (*IAutomationV21PlusCommon, error) {
	abi, err := abi.JSON(strings.NewReader(IAutomationV21PlusCommonABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAutomationV21PlusCommon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommon{address: address, abi: abi, IAutomationV21PlusCommonCaller: IAutomationV21PlusCommonCaller{contract: contract}, IAutomationV21PlusCommonTransactor: IAutomationV21PlusCommonTransactor{contract: contract}, IAutomationV21PlusCommonFilterer: IAutomationV21PlusCommonFilterer{contract: contract}}, nil
}

func NewIAutomationV21PlusCommonCaller(address common.Address, caller bind.ContractCaller) (*IAutomationV21PlusCommonCaller, error) {
	contract, err := bindIAutomationV21PlusCommon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonCaller{contract: contract}, nil
}

func NewIAutomationV21PlusCommonTransactor(address common.Address, transactor bind.ContractTransactor) (*IAutomationV21PlusCommonTransactor, error) {
	contract, err := bindIAutomationV21PlusCommon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonTransactor{contract: contract}, nil
}

func NewIAutomationV21PlusCommonFilterer(address common.Address, filterer bind.ContractFilterer) (*IAutomationV21PlusCommonFilterer, error) {
	contract, err := bindIAutomationV21PlusCommon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonFilterer{contract: contract}, nil
}

func bindIAutomationV21PlusCommon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAutomationV21PlusCommonMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationV21PlusCommon.Contract.IAutomationV21PlusCommonCaller.contract.Call(opts, result, method, params...)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.IAutomationV21PlusCommonTransactor.contract.Transfer(opts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.IAutomationV21PlusCommonTransactor.contract.Transact(opts, method, params...)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationV21PlusCommon.Contract.contract.Call(opts, result, method, params...)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.contract.Transfer(opts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.contract.Transact(opts, method, params...)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "checkCallback", id, values, extraData)

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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckCallback(&_IAutomationV21PlusCommon.CallOpts, id, values, extraData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckCallback(&_IAutomationV21PlusCommon.CallOpts, id, values, extraData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "checkUpkeep", id, triggerData)

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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckUpkeep(&_IAutomationV21PlusCommon.CallOpts, id, triggerData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckUpkeep(&_IAutomationV21PlusCommon.CallOpts, id, triggerData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

	error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "checkUpkeep0", id)

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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckUpkeep0(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationV21PlusCommon.Contract.CheckUpkeep0(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationV21PlusCommon.Contract.GetActiveUpkeepIDs(&_IAutomationV21PlusCommon.CallOpts, startIndex, maxCount)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationV21PlusCommon.Contract.GetActiveUpkeepIDs(&_IAutomationV21PlusCommon.CallOpts, startIndex, maxCount)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationV21PlusCommon.Contract.GetMinBalance(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationV21PlusCommon.Contract.GetMinBalance(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(IAutomationV21PlusCommonStateLegacy)).(*IAutomationV21PlusCommonStateLegacy)
	outstruct.Config = *abi.ConvertType(out[1], new(IAutomationV21PlusCommonOnchainConfigLegacy)).(*IAutomationV21PlusCommonOnchainConfigLegacy)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetState() (GetState,

	error) {
	return _IAutomationV21PlusCommon.Contract.GetState(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetState() (GetState,

	error) {
	return _IAutomationV21PlusCommon.Contract.GetState(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationV21PlusCommon.Contract.GetTriggerType(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationV21PlusCommon.Contract.GetTriggerType(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeep(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeep(&_IAutomationV21PlusCommon.CallOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeepPrivilegeConfig(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeepPrivilegeConfig(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeepTriggerConfig(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationV21PlusCommon.Contract.GetUpkeepTriggerConfig(&_IAutomationV21PlusCommon.CallOpts, upkeepId)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationV21PlusCommon.Contract.HasDedupKey(&_IAutomationV21PlusCommon.CallOpts, dedupKey)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationV21PlusCommon.Contract.HasDedupKey(&_IAutomationV21PlusCommon.CallOpts, dedupKey)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) Owner() (common.Address, error) {
	return _IAutomationV21PlusCommon.Contract.Owner(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) Owner() (common.Address, error) {
	return _IAutomationV21PlusCommon.Contract.Owner(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "simulatePerformUpkeep", id, performData)

	outstruct := new(SimulatePerformUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GasUsed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationV21PlusCommon.Contract.SimulatePerformUpkeep(&_IAutomationV21PlusCommon.CallOpts, id, performData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationV21PlusCommon.Contract.SimulatePerformUpkeep(&_IAutomationV21PlusCommon.CallOpts, id, performData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAutomationV21PlusCommon.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) TypeAndVersion() (string, error) {
	return _IAutomationV21PlusCommon.Contract.TypeAndVersion(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonCallerSession) TypeAndVersion() (string, error) {
	return _IAutomationV21PlusCommon.Contract.TypeAndVersion(&_IAutomationV21PlusCommon.CallOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "addFunds", id, amount)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.AddFunds(&_IAutomationV21PlusCommon.TransactOpts, id, amount)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.AddFunds(&_IAutomationV21PlusCommon.TransactOpts, id, amount)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.CancelUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.CancelUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.ExecuteCallback(&_IAutomationV21PlusCommon.TransactOpts, id, payload)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.ExecuteCallback(&_IAutomationV21PlusCommon.TransactOpts, id, payload)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "pause")
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) Pause() (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.Pause(&_IAutomationV21PlusCommon.TransactOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) Pause() (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.Pause(&_IAutomationV21PlusCommon.TransactOpts)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.PauseUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.PauseUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepCheckData(&_IAutomationV21PlusCommon.TransactOpts, id, newCheckData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepCheckData(&_IAutomationV21PlusCommon.TransactOpts, id, newCheckData)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepGasLimit(&_IAutomationV21PlusCommon.TransactOpts, id, gasLimit)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepGasLimit(&_IAutomationV21PlusCommon.TransactOpts, id, gasLimit)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepPrivilegeConfig(&_IAutomationV21PlusCommon.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepPrivilegeConfig(&_IAutomationV21PlusCommon.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepTriggerConfig(&_IAutomationV21PlusCommon.TransactOpts, id, triggerConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.SetUpkeepTriggerConfig(&_IAutomationV21PlusCommon.TransactOpts, id, triggerConfig)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.UnpauseUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationV21PlusCommon.Contract.UnpauseUpkeep(&_IAutomationV21PlusCommon.TransactOpts, id)
}

type IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator struct {
	Event *IAutomationV21PlusCommonAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonAdminPrivilegeConfigSet)
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
		it.Event = new(IAutomationV21PlusCommonAdminPrivilegeConfigSet)
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

func (it *IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonAdminPrivilegeConfigSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationV21PlusCommonAdminPrivilegeConfigSet, error) {
	event := new(IAutomationV21PlusCommonAdminPrivilegeConfigSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonCancelledUpkeepReportIterator struct {
	Event *IAutomationV21PlusCommonCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonCancelledUpkeepReport)
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
		it.Event = new(IAutomationV21PlusCommonCancelledUpkeepReport)
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

func (it *IAutomationV21PlusCommonCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonCancelledUpkeepReportIterator{contract: _IAutomationV21PlusCommon.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonCancelledUpkeepReport)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseCancelledUpkeepReport(log types.Log) (*IAutomationV21PlusCommonCancelledUpkeepReport, error) {
	event := new(IAutomationV21PlusCommonCancelledUpkeepReport)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonConfigSetIterator struct {
	Event *IAutomationV21PlusCommonConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonConfigSet)
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
		it.Event = new(IAutomationV21PlusCommonConfigSet)
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

func (it *IAutomationV21PlusCommonConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterConfigSet(opts *bind.FilterOpts) (*IAutomationV21PlusCommonConfigSetIterator, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonConfigSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonConfigSet) (event.Subscription, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonConfigSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseConfigSet(log types.Log) (*IAutomationV21PlusCommonConfigSet, error) {
	event := new(IAutomationV21PlusCommonConfigSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonDedupKeyAddedIterator struct {
	Event *IAutomationV21PlusCommonDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonDedupKeyAdded)
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
		it.Event = new(IAutomationV21PlusCommonDedupKeyAdded)
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

func (it *IAutomationV21PlusCommonDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationV21PlusCommonDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonDedupKeyAddedIterator{contract: _IAutomationV21PlusCommon.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonDedupKeyAdded)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseDedupKeyAdded(log types.Log) (*IAutomationV21PlusCommonDedupKeyAdded, error) {
	event := new(IAutomationV21PlusCommonDedupKeyAdded)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonFundsAddedIterator struct {
	Event *IAutomationV21PlusCommonFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonFundsAdded)
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
		it.Event = new(IAutomationV21PlusCommonFundsAdded)
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

func (it *IAutomationV21PlusCommonFundsAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationV21PlusCommonFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonFundsAddedIterator{contract: _IAutomationV21PlusCommon.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonFundsAdded)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseFundsAdded(log types.Log) (*IAutomationV21PlusCommonFundsAdded, error) {
	event := new(IAutomationV21PlusCommonFundsAdded)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonFundsWithdrawnIterator struct {
	Event *IAutomationV21PlusCommonFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonFundsWithdrawn)
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
		it.Event = new(IAutomationV21PlusCommonFundsWithdrawn)
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

func (it *IAutomationV21PlusCommonFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonFundsWithdrawnIterator{contract: _IAutomationV21PlusCommon.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonFundsWithdrawn)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseFundsWithdrawn(log types.Log) (*IAutomationV21PlusCommonFundsWithdrawn, error) {
	event := new(IAutomationV21PlusCommonFundsWithdrawn)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator struct {
	Event *IAutomationV21PlusCommonInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonInsufficientFundsUpkeepReport)
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
		it.Event = new(IAutomationV21PlusCommonInsufficientFundsUpkeepReport)
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

func (it *IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator{contract: _IAutomationV21PlusCommon.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonInsufficientFundsUpkeepReport)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationV21PlusCommonInsufficientFundsUpkeepReport, error) {
	event := new(IAutomationV21PlusCommonInsufficientFundsUpkeepReport)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonOwnershipTransferRequestedIterator struct {
	Event *IAutomationV21PlusCommonOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonOwnershipTransferRequested)
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
		it.Event = new(IAutomationV21PlusCommonOwnershipTransferRequested)
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

func (it *IAutomationV21PlusCommonOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonOwnershipTransferRequestedIterator{contract: _IAutomationV21PlusCommon.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonOwnershipTransferRequested)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseOwnershipTransferRequested(log types.Log) (*IAutomationV21PlusCommonOwnershipTransferRequested, error) {
	event := new(IAutomationV21PlusCommonOwnershipTransferRequested)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonOwnershipTransferredIterator struct {
	Event *IAutomationV21PlusCommonOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonOwnershipTransferred)
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
		it.Event = new(IAutomationV21PlusCommonOwnershipTransferred)
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

func (it *IAutomationV21PlusCommonOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonOwnershipTransferredIterator{contract: _IAutomationV21PlusCommon.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonOwnershipTransferred)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseOwnershipTransferred(log types.Log) (*IAutomationV21PlusCommonOwnershipTransferred, error) {
	event := new(IAutomationV21PlusCommonOwnershipTransferred)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonPausedIterator struct {
	Event *IAutomationV21PlusCommonPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonPaused)
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
		it.Event = new(IAutomationV21PlusCommonPaused)
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

func (it *IAutomationV21PlusCommonPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterPaused(opts *bind.FilterOpts) (*IAutomationV21PlusCommonPausedIterator, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonPausedIterator{contract: _IAutomationV21PlusCommon.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonPaused)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParsePaused(log types.Log) (*IAutomationV21PlusCommonPaused, error) {
	event := new(IAutomationV21PlusCommonPaused)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonPayeesUpdatedIterator struct {
	Event *IAutomationV21PlusCommonPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonPayeesUpdated)
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
		it.Event = new(IAutomationV21PlusCommonPayeesUpdated)
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

func (it *IAutomationV21PlusCommonPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationV21PlusCommonPayeesUpdatedIterator, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonPayeesUpdatedIterator{contract: _IAutomationV21PlusCommon.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonPayeesUpdated)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParsePayeesUpdated(log types.Log) (*IAutomationV21PlusCommonPayeesUpdated, error) {
	event := new(IAutomationV21PlusCommonPayeesUpdated)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonPayeeshipTransferRequestedIterator struct {
	Event *IAutomationV21PlusCommonPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonPayeeshipTransferRequested)
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
		it.Event = new(IAutomationV21PlusCommonPayeeshipTransferRequested)
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

func (it *IAutomationV21PlusCommonPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonPayeeshipTransferRequestedIterator{contract: _IAutomationV21PlusCommon.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonPayeeshipTransferRequested)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParsePayeeshipTransferRequested(log types.Log) (*IAutomationV21PlusCommonPayeeshipTransferRequested, error) {
	event := new(IAutomationV21PlusCommonPayeeshipTransferRequested)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonPayeeshipTransferredIterator struct {
	Event *IAutomationV21PlusCommonPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonPayeeshipTransferred)
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
		it.Event = new(IAutomationV21PlusCommonPayeeshipTransferred)
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

func (it *IAutomationV21PlusCommonPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonPayeeshipTransferredIterator{contract: _IAutomationV21PlusCommon.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonPayeeshipTransferred)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParsePayeeshipTransferred(log types.Log) (*IAutomationV21PlusCommonPayeeshipTransferred, error) {
	event := new(IAutomationV21PlusCommonPayeeshipTransferred)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonPaymentWithdrawnIterator struct {
	Event *IAutomationV21PlusCommonPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonPaymentWithdrawn)
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
		it.Event = new(IAutomationV21PlusCommonPaymentWithdrawn)
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

func (it *IAutomationV21PlusCommonPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationV21PlusCommonPaymentWithdrawnIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonPaymentWithdrawnIterator{contract: _IAutomationV21PlusCommon.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonPaymentWithdrawn)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParsePaymentWithdrawn(log types.Log) (*IAutomationV21PlusCommonPaymentWithdrawn, error) {
	event := new(IAutomationV21PlusCommonPaymentWithdrawn)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonReorgedUpkeepReportIterator struct {
	Event *IAutomationV21PlusCommonReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonReorgedUpkeepReport)
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
		it.Event = new(IAutomationV21PlusCommonReorgedUpkeepReport)
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

func (it *IAutomationV21PlusCommonReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonReorgedUpkeepReportIterator{contract: _IAutomationV21PlusCommon.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonReorgedUpkeepReport)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseReorgedUpkeepReport(log types.Log) (*IAutomationV21PlusCommonReorgedUpkeepReport, error) {
	event := new(IAutomationV21PlusCommonReorgedUpkeepReport)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonStaleUpkeepReportIterator struct {
	Event *IAutomationV21PlusCommonStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonStaleUpkeepReport)
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
		it.Event = new(IAutomationV21PlusCommonStaleUpkeepReport)
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

func (it *IAutomationV21PlusCommonStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonStaleUpkeepReportIterator{contract: _IAutomationV21PlusCommon.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonStaleUpkeepReport)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseStaleUpkeepReport(log types.Log) (*IAutomationV21PlusCommonStaleUpkeepReport, error) {
	event := new(IAutomationV21PlusCommonStaleUpkeepReport)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonTransmittedIterator struct {
	Event *IAutomationV21PlusCommonTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonTransmitted)
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
		it.Event = new(IAutomationV21PlusCommonTransmitted)
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

func (it *IAutomationV21PlusCommonTransmittedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterTransmitted(opts *bind.FilterOpts) (*IAutomationV21PlusCommonTransmittedIterator, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonTransmittedIterator{contract: _IAutomationV21PlusCommon.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonTransmitted) (event.Subscription, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonTransmitted)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseTransmitted(log types.Log) (*IAutomationV21PlusCommonTransmitted, error) {
	event := new(IAutomationV21PlusCommonTransmitted)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUnpausedIterator struct {
	Event *IAutomationV21PlusCommonUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUnpaused)
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
		it.Event = new(IAutomationV21PlusCommonUnpaused)
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

func (it *IAutomationV21PlusCommonUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IAutomationV21PlusCommonUnpausedIterator, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUnpausedIterator{contract: _IAutomationV21PlusCommon.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUnpaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUnpaused)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUnpaused(log types.Log) (*IAutomationV21PlusCommonUnpaused, error) {
	event := new(IAutomationV21PlusCommonUnpaused)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepAdminTransferRequested)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepAdminTransferRequested)
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

func (it *IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepAdminTransferRequested)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationV21PlusCommonUpkeepAdminTransferRequested, error) {
	event := new(IAutomationV21PlusCommonUpkeepAdminTransferRequested)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepAdminTransferredIterator struct {
	Event *IAutomationV21PlusCommonUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepAdminTransferred)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepAdminTransferred)
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

func (it *IAutomationV21PlusCommonUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonUpkeepAdminTransferredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepAdminTransferredIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepAdminTransferred)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepAdminTransferred(log types.Log) (*IAutomationV21PlusCommonUpkeepAdminTransferred, error) {
	event := new(IAutomationV21PlusCommonUpkeepAdminTransferred)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepCanceledIterator struct {
	Event *IAutomationV21PlusCommonUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepCanceled)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepCanceled)
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

func (it *IAutomationV21PlusCommonUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationV21PlusCommonUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepCanceledIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepCanceled)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepCanceled(log types.Log) (*IAutomationV21PlusCommonUpkeepCanceled, error) {
	event := new(IAutomationV21PlusCommonUpkeepCanceled)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepCheckDataSetIterator struct {
	Event *IAutomationV21PlusCommonUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepCheckDataSet)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepCheckDataSet)
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

func (it *IAutomationV21PlusCommonUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepCheckDataSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepCheckDataSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepCheckDataSet(log types.Log) (*IAutomationV21PlusCommonUpkeepCheckDataSet, error) {
	event := new(IAutomationV21PlusCommonUpkeepCheckDataSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepGasLimitSetIterator struct {
	Event *IAutomationV21PlusCommonUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepGasLimitSet)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepGasLimitSet)
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

func (it *IAutomationV21PlusCommonUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepGasLimitSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepGasLimitSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepGasLimitSet(log types.Log) (*IAutomationV21PlusCommonUpkeepGasLimitSet, error) {
	event := new(IAutomationV21PlusCommonUpkeepGasLimitSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepMigratedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepMigrated)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepMigrated)
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

func (it *IAutomationV21PlusCommonUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepMigratedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepMigrated)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepMigrated(log types.Log) (*IAutomationV21PlusCommonUpkeepMigrated, error) {
	event := new(IAutomationV21PlusCommonUpkeepMigrated)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator struct {
	Event *IAutomationV21PlusCommonUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepOffchainConfigSet)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepOffchainConfigSet)
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

func (it *IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepOffchainConfigSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepOffchainConfigSet, error) {
	event := new(IAutomationV21PlusCommonUpkeepOffchainConfigSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepPausedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepPaused)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepPaused)
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

func (it *IAutomationV21PlusCommonUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepPausedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepPaused)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepPaused(log types.Log) (*IAutomationV21PlusCommonUpkeepPaused, error) {
	event := new(IAutomationV21PlusCommonUpkeepPaused)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepPerformedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepPerformed)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepPerformed)
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

func (it *IAutomationV21PlusCommonUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationV21PlusCommonUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepPerformedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepPerformed)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepPerformed(log types.Log) (*IAutomationV21PlusCommonUpkeepPerformed, error) {
	event := new(IAutomationV21PlusCommonUpkeepPerformed)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator struct {
	Event *IAutomationV21PlusCommonUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepPrivilegeConfigSet)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepPrivilegeConfigSet)
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

func (it *IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepPrivilegeConfigSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepPrivilegeConfigSet, error) {
	event := new(IAutomationV21PlusCommonUpkeepPrivilegeConfigSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepReceivedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepReceived)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepReceived)
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

func (it *IAutomationV21PlusCommonUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepReceivedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepReceived)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepReceived(log types.Log) (*IAutomationV21PlusCommonUpkeepReceived, error) {
	event := new(IAutomationV21PlusCommonUpkeepReceived)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepRegisteredIterator struct {
	Event *IAutomationV21PlusCommonUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepRegistered)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepRegistered)
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

func (it *IAutomationV21PlusCommonUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepRegisteredIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepRegistered)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepRegistered(log types.Log) (*IAutomationV21PlusCommonUpkeepRegistered, error) {
	event := new(IAutomationV21PlusCommonUpkeepRegistered)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator struct {
	Event *IAutomationV21PlusCommonUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepTriggerConfigSet)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepTriggerConfigSet)
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

func (it *IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepTriggerConfigSet)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepTriggerConfigSet, error) {
	event := new(IAutomationV21PlusCommonUpkeepTriggerConfigSet)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV21PlusCommonUpkeepUnpausedIterator struct {
	Event *IAutomationV21PlusCommonUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV21PlusCommonUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV21PlusCommonUpkeepUnpaused)
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
		it.Event = new(IAutomationV21PlusCommonUpkeepUnpaused)
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

func (it *IAutomationV21PlusCommonUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV21PlusCommonUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV21PlusCommonUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV21PlusCommonUpkeepUnpausedIterator{contract: _IAutomationV21PlusCommon.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV21PlusCommon.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV21PlusCommonUpkeepUnpaused)
				if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommonFilterer) ParseUpkeepUnpaused(log types.Log) (*IAutomationV21PlusCommonUpkeepUnpaused, error) {
	event := new(IAutomationV21PlusCommonUpkeepUnpaused)
	if err := _IAutomationV21PlusCommon.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
	State        IAutomationV21PlusCommonStateLegacy
	Config       IAutomationV21PlusCommonOnchainConfigLegacy
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}
type SimulatePerformUpkeep struct {
	Success bool
	GasUsed *big.Int
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommon) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IAutomationV21PlusCommon.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _IAutomationV21PlusCommon.ParseAdminPrivilegeConfigSet(log)
	case _IAutomationV21PlusCommon.abi.Events["CancelledUpkeepReport"].ID:
		return _IAutomationV21PlusCommon.ParseCancelledUpkeepReport(log)
	case _IAutomationV21PlusCommon.abi.Events["ConfigSet"].ID:
		return _IAutomationV21PlusCommon.ParseConfigSet(log)
	case _IAutomationV21PlusCommon.abi.Events["DedupKeyAdded"].ID:
		return _IAutomationV21PlusCommon.ParseDedupKeyAdded(log)
	case _IAutomationV21PlusCommon.abi.Events["FundsAdded"].ID:
		return _IAutomationV21PlusCommon.ParseFundsAdded(log)
	case _IAutomationV21PlusCommon.abi.Events["FundsWithdrawn"].ID:
		return _IAutomationV21PlusCommon.ParseFundsWithdrawn(log)
	case _IAutomationV21PlusCommon.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IAutomationV21PlusCommon.ParseInsufficientFundsUpkeepReport(log)
	case _IAutomationV21PlusCommon.abi.Events["OwnershipTransferRequested"].ID:
		return _IAutomationV21PlusCommon.ParseOwnershipTransferRequested(log)
	case _IAutomationV21PlusCommon.abi.Events["OwnershipTransferred"].ID:
		return _IAutomationV21PlusCommon.ParseOwnershipTransferred(log)
	case _IAutomationV21PlusCommon.abi.Events["Paused"].ID:
		return _IAutomationV21PlusCommon.ParsePaused(log)
	case _IAutomationV21PlusCommon.abi.Events["PayeesUpdated"].ID:
		return _IAutomationV21PlusCommon.ParsePayeesUpdated(log)
	case _IAutomationV21PlusCommon.abi.Events["PayeeshipTransferRequested"].ID:
		return _IAutomationV21PlusCommon.ParsePayeeshipTransferRequested(log)
	case _IAutomationV21PlusCommon.abi.Events["PayeeshipTransferred"].ID:
		return _IAutomationV21PlusCommon.ParsePayeeshipTransferred(log)
	case _IAutomationV21PlusCommon.abi.Events["PaymentWithdrawn"].ID:
		return _IAutomationV21PlusCommon.ParsePaymentWithdrawn(log)
	case _IAutomationV21PlusCommon.abi.Events["ReorgedUpkeepReport"].ID:
		return _IAutomationV21PlusCommon.ParseReorgedUpkeepReport(log)
	case _IAutomationV21PlusCommon.abi.Events["StaleUpkeepReport"].ID:
		return _IAutomationV21PlusCommon.ParseStaleUpkeepReport(log)
	case _IAutomationV21PlusCommon.abi.Events["Transmitted"].ID:
		return _IAutomationV21PlusCommon.ParseTransmitted(log)
	case _IAutomationV21PlusCommon.abi.Events["Unpaused"].ID:
		return _IAutomationV21PlusCommon.ParseUnpaused(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepAdminTransferRequested(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepAdminTransferred"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepAdminTransferred(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepCanceled"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepCanceled(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepCheckDataSet"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepCheckDataSet(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepGasLimitSet"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepGasLimitSet(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepMigrated"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepMigrated(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepOffchainConfigSet(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepPaused"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepPaused(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepPerformed"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepPerformed(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepPrivilegeConfigSet(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepReceived"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepReceived(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepRegistered"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepRegistered(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepTriggerConfigSet(log)
	case _IAutomationV21PlusCommon.abi.Events["UpkeepUnpaused"].ID:
		return _IAutomationV21PlusCommon.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IAutomationV21PlusCommonAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (IAutomationV21PlusCommonCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IAutomationV21PlusCommonConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IAutomationV21PlusCommonDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IAutomationV21PlusCommonFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IAutomationV21PlusCommonFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IAutomationV21PlusCommonInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IAutomationV21PlusCommonOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IAutomationV21PlusCommonOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IAutomationV21PlusCommonPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IAutomationV21PlusCommonPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IAutomationV21PlusCommonPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IAutomationV21PlusCommonPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IAutomationV21PlusCommonPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IAutomationV21PlusCommonReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IAutomationV21PlusCommonStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IAutomationV21PlusCommonTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IAutomationV21PlusCommonUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IAutomationV21PlusCommonUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IAutomationV21PlusCommonUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IAutomationV21PlusCommonUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IAutomationV21PlusCommonUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IAutomationV21PlusCommonUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IAutomationV21PlusCommonUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IAutomationV21PlusCommonUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IAutomationV21PlusCommonUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IAutomationV21PlusCommonUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IAutomationV21PlusCommonUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (IAutomationV21PlusCommonUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IAutomationV21PlusCommonUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IAutomationV21PlusCommonUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IAutomationV21PlusCommonUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IAutomationV21PlusCommon *IAutomationV21PlusCommon) Address() common.Address {
	return _IAutomationV21PlusCommon.address
}

type IAutomationV21PlusCommonInterface interface {
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

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error)

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

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationV21PlusCommonAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationV21PlusCommonAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IAutomationV21PlusCommonCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IAutomationV21PlusCommonConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IAutomationV21PlusCommonConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationV21PlusCommonDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IAutomationV21PlusCommonDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationV21PlusCommonFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IAutomationV21PlusCommonFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IAutomationV21PlusCommonFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationV21PlusCommonInsufficientFundsUpkeepReport, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IAutomationV21PlusCommonOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IAutomationV21PlusCommonOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IAutomationV21PlusCommonPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IAutomationV21PlusCommonPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationV21PlusCommonPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IAutomationV21PlusCommonPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IAutomationV21PlusCommonPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IAutomationV21PlusCommonPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationV21PlusCommonPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IAutomationV21PlusCommonPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IAutomationV21PlusCommonReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IAutomationV21PlusCommonStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IAutomationV21PlusCommonTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IAutomationV21PlusCommonTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IAutomationV21PlusCommonUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IAutomationV21PlusCommonUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationV21PlusCommonUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV21PlusCommonUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IAutomationV21PlusCommonUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationV21PlusCommonUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IAutomationV21PlusCommonUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IAutomationV21PlusCommonUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IAutomationV21PlusCommonUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IAutomationV21PlusCommonUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IAutomationV21PlusCommonUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationV21PlusCommonUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IAutomationV21PlusCommonUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IAutomationV21PlusCommonUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IAutomationV21PlusCommonUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationV21PlusCommonUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV21PlusCommonUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV21PlusCommonUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IAutomationV21PlusCommonUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
