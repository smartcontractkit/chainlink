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

type OnchainConfigV21 struct {
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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV2Common.State\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfigV21\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV2Common.UpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
	outstruct.Config = *abi.ConvertType(out[1], new(OnchainConfigV21)).(*OnchainConfigV21)
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

type IAutomationV2CommonAdminPrivilegeConfigSetIterator struct {
	Event *IAutomationV2CommonAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonAdminPrivilegeConfigSet)
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
		it.Event = new(IAutomationV2CommonAdminPrivilegeConfigSet)
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

func (it *IAutomationV2CommonAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationV2CommonAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonAdminPrivilegeConfigSetIterator{contract: _IAutomationV2Common.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonAdminPrivilegeConfigSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationV2CommonAdminPrivilegeConfigSet, error) {
	event := new(IAutomationV2CommonAdminPrivilegeConfigSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonCancelledUpkeepReportIterator struct {
	Event *IAutomationV2CommonCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonCancelledUpkeepReport)
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
		it.Event = new(IAutomationV2CommonCancelledUpkeepReport)
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

func (it *IAutomationV2CommonCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonCancelledUpkeepReportIterator{contract: _IAutomationV2Common.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonCancelledUpkeepReport)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseCancelledUpkeepReport(log types.Log) (*IAutomationV2CommonCancelledUpkeepReport, error) {
	event := new(IAutomationV2CommonCancelledUpkeepReport)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonConfigSetIterator struct {
	Event *IAutomationV2CommonConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonConfigSet)
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
		it.Event = new(IAutomationV2CommonConfigSet)
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

func (it *IAutomationV2CommonConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonConfigSet struct {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterConfigSet(opts *bind.FilterOpts) (*IAutomationV2CommonConfigSetIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonConfigSetIterator{contract: _IAutomationV2Common.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonConfigSet) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonConfigSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseConfigSet(log types.Log) (*IAutomationV2CommonConfigSet, error) {
	event := new(IAutomationV2CommonConfigSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

type IAutomationV2CommonFundsAddedIterator struct {
	Event *IAutomationV2CommonFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonFundsAdded)
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
		it.Event = new(IAutomationV2CommonFundsAdded)
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

func (it *IAutomationV2CommonFundsAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationV2CommonFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonFundsAddedIterator{contract: _IAutomationV2Common.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonFundsAdded)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseFundsAdded(log types.Log) (*IAutomationV2CommonFundsAdded, error) {
	event := new(IAutomationV2CommonFundsAdded)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonFundsWithdrawnIterator struct {
	Event *IAutomationV2CommonFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonFundsWithdrawn)
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
		it.Event = new(IAutomationV2CommonFundsWithdrawn)
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

func (it *IAutomationV2CommonFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonFundsWithdrawnIterator{contract: _IAutomationV2Common.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonFundsWithdrawn)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseFundsWithdrawn(log types.Log) (*IAutomationV2CommonFundsWithdrawn, error) {
	event := new(IAutomationV2CommonFundsWithdrawn)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

type IAutomationV2CommonOwnerFundsWithdrawnIterator struct {
	Event *IAutomationV2CommonOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonOwnerFundsWithdrawn)
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
		it.Event = new(IAutomationV2CommonOwnerFundsWithdrawn)
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

func (it *IAutomationV2CommonOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IAutomationV2CommonOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonOwnerFundsWithdrawnIterator{contract: _IAutomationV2Common.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonOwnerFundsWithdrawn)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*IAutomationV2CommonOwnerFundsWithdrawn, error) {
	event := new(IAutomationV2CommonOwnerFundsWithdrawn)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonOwnershipTransferRequestedIterator struct {
	Event *IAutomationV2CommonOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonOwnershipTransferRequested)
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
		it.Event = new(IAutomationV2CommonOwnershipTransferRequested)
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

func (it *IAutomationV2CommonOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV2CommonOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonOwnershipTransferRequestedIterator{contract: _IAutomationV2Common.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonOwnershipTransferRequested)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseOwnershipTransferRequested(log types.Log) (*IAutomationV2CommonOwnershipTransferRequested, error) {
	event := new(IAutomationV2CommonOwnershipTransferRequested)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonOwnershipTransferredIterator struct {
	Event *IAutomationV2CommonOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonOwnershipTransferred)
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
		it.Event = new(IAutomationV2CommonOwnershipTransferred)
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

func (it *IAutomationV2CommonOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV2CommonOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonOwnershipTransferredIterator{contract: _IAutomationV2Common.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonOwnershipTransferred)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseOwnershipTransferred(log types.Log) (*IAutomationV2CommonOwnershipTransferred, error) {
	event := new(IAutomationV2CommonOwnershipTransferred)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonPausedIterator struct {
	Event *IAutomationV2CommonPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonPaused)
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
		it.Event = new(IAutomationV2CommonPaused)
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

func (it *IAutomationV2CommonPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterPaused(opts *bind.FilterOpts) (*IAutomationV2CommonPausedIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonPausedIterator{contract: _IAutomationV2Common.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonPaused)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParsePaused(log types.Log) (*IAutomationV2CommonPaused, error) {
	event := new(IAutomationV2CommonPaused)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonPayeesUpdatedIterator struct {
	Event *IAutomationV2CommonPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonPayeesUpdated)
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
		it.Event = new(IAutomationV2CommonPayeesUpdated)
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

func (it *IAutomationV2CommonPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationV2CommonPayeesUpdatedIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonPayeesUpdatedIterator{contract: _IAutomationV2Common.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonPayeesUpdated)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParsePayeesUpdated(log types.Log) (*IAutomationV2CommonPayeesUpdated, error) {
	event := new(IAutomationV2CommonPayeesUpdated)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonPayeeshipTransferRequestedIterator struct {
	Event *IAutomationV2CommonPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonPayeeshipTransferRequested)
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
		it.Event = new(IAutomationV2CommonPayeeshipTransferRequested)
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

func (it *IAutomationV2CommonPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV2CommonPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonPayeeshipTransferRequestedIterator{contract: _IAutomationV2Common.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonPayeeshipTransferRequested)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParsePayeeshipTransferRequested(log types.Log) (*IAutomationV2CommonPayeeshipTransferRequested, error) {
	event := new(IAutomationV2CommonPayeeshipTransferRequested)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonPayeeshipTransferredIterator struct {
	Event *IAutomationV2CommonPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonPayeeshipTransferred)
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
		it.Event = new(IAutomationV2CommonPayeeshipTransferred)
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

func (it *IAutomationV2CommonPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV2CommonPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonPayeeshipTransferredIterator{contract: _IAutomationV2Common.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonPayeeshipTransferred)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParsePayeeshipTransferred(log types.Log) (*IAutomationV2CommonPayeeshipTransferred, error) {
	event := new(IAutomationV2CommonPayeeshipTransferred)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonPaymentWithdrawnIterator struct {
	Event *IAutomationV2CommonPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonPaymentWithdrawn)
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
		it.Event = new(IAutomationV2CommonPaymentWithdrawn)
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

func (it *IAutomationV2CommonPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationV2CommonPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonPaymentWithdrawnIterator{contract: _IAutomationV2Common.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonPaymentWithdrawn)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParsePaymentWithdrawn(log types.Log) (*IAutomationV2CommonPaymentWithdrawn, error) {
	event := new(IAutomationV2CommonPaymentWithdrawn)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

type IAutomationV2CommonTransmittedIterator struct {
	Event *IAutomationV2CommonTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonTransmitted)
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
		it.Event = new(IAutomationV2CommonTransmitted)
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

func (it *IAutomationV2CommonTransmittedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterTransmitted(opts *bind.FilterOpts) (*IAutomationV2CommonTransmittedIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonTransmittedIterator{contract: _IAutomationV2Common.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonTransmitted) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonTransmitted)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseTransmitted(log types.Log) (*IAutomationV2CommonTransmitted, error) {
	event := new(IAutomationV2CommonTransmitted)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUnpausedIterator struct {
	Event *IAutomationV2CommonUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUnpaused)
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
		it.Event = new(IAutomationV2CommonUnpaused)
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

func (it *IAutomationV2CommonUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IAutomationV2CommonUnpausedIterator, error) {

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUnpausedIterator{contract: _IAutomationV2Common.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUnpaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUnpaused)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUnpaused(log types.Log) (*IAutomationV2CommonUnpaused, error) {
	event := new(IAutomationV2CommonUnpaused)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepAdminTransferRequestedIterator struct {
	Event *IAutomationV2CommonUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepAdminTransferRequested)
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
		it.Event = new(IAutomationV2CommonUpkeepAdminTransferRequested)
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

func (it *IAutomationV2CommonUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV2CommonUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepAdminTransferRequestedIterator{contract: _IAutomationV2Common.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepAdminTransferRequested)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationV2CommonUpkeepAdminTransferRequested, error) {
	event := new(IAutomationV2CommonUpkeepAdminTransferRequested)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepAdminTransferredIterator struct {
	Event *IAutomationV2CommonUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepAdminTransferred)
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
		it.Event = new(IAutomationV2CommonUpkeepAdminTransferred)
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

func (it *IAutomationV2CommonUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV2CommonUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepAdminTransferredIterator{contract: _IAutomationV2Common.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepAdminTransferred)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepAdminTransferred(log types.Log) (*IAutomationV2CommonUpkeepAdminTransferred, error) {
	event := new(IAutomationV2CommonUpkeepAdminTransferred)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

type IAutomationV2CommonUpkeepCheckDataSetIterator struct {
	Event *IAutomationV2CommonUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepCheckDataSet)
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
		it.Event = new(IAutomationV2CommonUpkeepCheckDataSet)
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

func (it *IAutomationV2CommonUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepCheckDataSetIterator{contract: _IAutomationV2Common.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepCheckDataSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepCheckDataSet(log types.Log) (*IAutomationV2CommonUpkeepCheckDataSet, error) {
	event := new(IAutomationV2CommonUpkeepCheckDataSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationV2CommonUpkeepGasLimitSetIterator struct {
	Event *IAutomationV2CommonUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepGasLimitSet)
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
		it.Event = new(IAutomationV2CommonUpkeepGasLimitSet)
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

func (it *IAutomationV2CommonUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepGasLimitSetIterator{contract: _IAutomationV2Common.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepGasLimitSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepGasLimitSet(log types.Log) (*IAutomationV2CommonUpkeepGasLimitSet, error) {
	event := new(IAutomationV2CommonUpkeepGasLimitSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

type IAutomationV2CommonUpkeepOffchainConfigSetIterator struct {
	Event *IAutomationV2CommonUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepOffchainConfigSet)
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
		it.Event = new(IAutomationV2CommonUpkeepOffchainConfigSet)
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

func (it *IAutomationV2CommonUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepOffchainConfigSetIterator{contract: _IAutomationV2Common.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepOffchainConfigSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationV2CommonUpkeepOffchainConfigSet, error) {
	event := new(IAutomationV2CommonUpkeepOffchainConfigSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

type IAutomationV2CommonUpkeepPrivilegeConfigSetIterator struct {
	Event *IAutomationV2CommonUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationV2CommonUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationV2CommonUpkeepPrivilegeConfigSet)
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
		it.Event = new(IAutomationV2CommonUpkeepPrivilegeConfigSet)
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

func (it *IAutomationV2CommonUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationV2CommonUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationV2CommonUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationV2CommonUpkeepPrivilegeConfigSetIterator{contract: _IAutomationV2Common.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationV2Common *IAutomationV2CommonFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationV2Common.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationV2CommonUpkeepPrivilegeConfigSet)
				if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_IAutomationV2Common *IAutomationV2CommonFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationV2CommonUpkeepPrivilegeConfigSet, error) {
	event := new(IAutomationV2CommonUpkeepPrivilegeConfigSet)
	if err := _IAutomationV2Common.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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
	Config       OnchainConfigV21
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
	case _IAutomationV2Common.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _IAutomationV2Common.ParseAdminPrivilegeConfigSet(log)
	case _IAutomationV2Common.abi.Events["CancelledUpkeepReport"].ID:
		return _IAutomationV2Common.ParseCancelledUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["ConfigSet"].ID:
		return _IAutomationV2Common.ParseConfigSet(log)
	case _IAutomationV2Common.abi.Events["DedupKeyAdded"].ID:
		return _IAutomationV2Common.ParseDedupKeyAdded(log)
	case _IAutomationV2Common.abi.Events["FundsAdded"].ID:
		return _IAutomationV2Common.ParseFundsAdded(log)
	case _IAutomationV2Common.abi.Events["FundsWithdrawn"].ID:
		return _IAutomationV2Common.ParseFundsWithdrawn(log)
	case _IAutomationV2Common.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IAutomationV2Common.ParseInsufficientFundsUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["OwnerFundsWithdrawn"].ID:
		return _IAutomationV2Common.ParseOwnerFundsWithdrawn(log)
	case _IAutomationV2Common.abi.Events["OwnershipTransferRequested"].ID:
		return _IAutomationV2Common.ParseOwnershipTransferRequested(log)
	case _IAutomationV2Common.abi.Events["OwnershipTransferred"].ID:
		return _IAutomationV2Common.ParseOwnershipTransferred(log)
	case _IAutomationV2Common.abi.Events["Paused"].ID:
		return _IAutomationV2Common.ParsePaused(log)
	case _IAutomationV2Common.abi.Events["PayeesUpdated"].ID:
		return _IAutomationV2Common.ParsePayeesUpdated(log)
	case _IAutomationV2Common.abi.Events["PayeeshipTransferRequested"].ID:
		return _IAutomationV2Common.ParsePayeeshipTransferRequested(log)
	case _IAutomationV2Common.abi.Events["PayeeshipTransferred"].ID:
		return _IAutomationV2Common.ParsePayeeshipTransferred(log)
	case _IAutomationV2Common.abi.Events["PaymentWithdrawn"].ID:
		return _IAutomationV2Common.ParsePaymentWithdrawn(log)
	case _IAutomationV2Common.abi.Events["ReorgedUpkeepReport"].ID:
		return _IAutomationV2Common.ParseReorgedUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["StaleUpkeepReport"].ID:
		return _IAutomationV2Common.ParseStaleUpkeepReport(log)
	case _IAutomationV2Common.abi.Events["Transmitted"].ID:
		return _IAutomationV2Common.ParseTransmitted(log)
	case _IAutomationV2Common.abi.Events["Unpaused"].ID:
		return _IAutomationV2Common.ParseUnpaused(log)
	case _IAutomationV2Common.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IAutomationV2Common.ParseUpkeepAdminTransferRequested(log)
	case _IAutomationV2Common.abi.Events["UpkeepAdminTransferred"].ID:
		return _IAutomationV2Common.ParseUpkeepAdminTransferred(log)
	case _IAutomationV2Common.abi.Events["UpkeepCanceled"].ID:
		return _IAutomationV2Common.ParseUpkeepCanceled(log)
	case _IAutomationV2Common.abi.Events["UpkeepCheckDataSet"].ID:
		return _IAutomationV2Common.ParseUpkeepCheckDataSet(log)
	case _IAutomationV2Common.abi.Events["UpkeepGasLimitSet"].ID:
		return _IAutomationV2Common.ParseUpkeepGasLimitSet(log)
	case _IAutomationV2Common.abi.Events["UpkeepMigrated"].ID:
		return _IAutomationV2Common.ParseUpkeepMigrated(log)
	case _IAutomationV2Common.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IAutomationV2Common.ParseUpkeepOffchainConfigSet(log)
	case _IAutomationV2Common.abi.Events["UpkeepPaused"].ID:
		return _IAutomationV2Common.ParseUpkeepPaused(log)
	case _IAutomationV2Common.abi.Events["UpkeepPerformed"].ID:
		return _IAutomationV2Common.ParseUpkeepPerformed(log)
	case _IAutomationV2Common.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IAutomationV2Common.ParseUpkeepPrivilegeConfigSet(log)
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

func (IAutomationV2CommonAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (IAutomationV2CommonCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IAutomationV2CommonConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IAutomationV2CommonDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IAutomationV2CommonFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IAutomationV2CommonFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IAutomationV2CommonInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IAutomationV2CommonOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (IAutomationV2CommonOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IAutomationV2CommonOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IAutomationV2CommonPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IAutomationV2CommonPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IAutomationV2CommonPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IAutomationV2CommonPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IAutomationV2CommonPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IAutomationV2CommonReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IAutomationV2CommonStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IAutomationV2CommonTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IAutomationV2CommonUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IAutomationV2CommonUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IAutomationV2CommonUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IAutomationV2CommonUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IAutomationV2CommonUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IAutomationV2CommonUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IAutomationV2CommonUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IAutomationV2CommonUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IAutomationV2CommonUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IAutomationV2CommonUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IAutomationV2CommonUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationV2CommonAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationV2CommonAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IAutomationV2CommonCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IAutomationV2CommonConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IAutomationV2CommonConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationV2CommonDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IAutomationV2CommonDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationV2CommonFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IAutomationV2CommonFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IAutomationV2CommonFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationV2CommonInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IAutomationV2CommonOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*IAutomationV2CommonOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV2CommonOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IAutomationV2CommonOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationV2CommonOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IAutomationV2CommonOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IAutomationV2CommonPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IAutomationV2CommonPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationV2CommonPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IAutomationV2CommonPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV2CommonPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IAutomationV2CommonPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationV2CommonPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IAutomationV2CommonPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationV2CommonPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IAutomationV2CommonPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IAutomationV2CommonReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IAutomationV2CommonStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IAutomationV2CommonTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IAutomationV2CommonTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IAutomationV2CommonUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IAutomationV2CommonUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV2CommonUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationV2CommonUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationV2CommonUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IAutomationV2CommonUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationV2CommonUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IAutomationV2CommonUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IAutomationV2CommonUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IAutomationV2CommonUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IAutomationV2CommonUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationV2CommonUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IAutomationV2CommonUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationV2CommonUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IAutomationV2CommonUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationV2CommonUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationV2CommonUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationV2CommonUpkeepPrivilegeConfigSet, error)

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
