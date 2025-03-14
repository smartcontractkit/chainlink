// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_registry_logic_b_wrapper_2_2

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

var AutomationRegistryLogicB22MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"linkNativeFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fastGasFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"automationForwarderLogic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedReadOnlyAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveUpkeepIDs\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedReadOnlyAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAutomationForwarderLogic\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCancellationDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getChainModule\",\"inputs\":[],\"outputs\":[{\"name\":\"chainModule\",\"type\":\"address\",\"internalType\":\"contractIChainModule\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConditionalGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getFastGasFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForwarder\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAutomationForwarder\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkNativeFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLogGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getMaxPaymentForGas\",\"inputs\":[{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"maxPayment\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalanceForUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"minBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.MigrationPermission\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPerPerformByteGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPerSignerGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getReorgProtectionEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSignerInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ownerLinkBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"expectedLinkBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalPremium\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"numUpkeeps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"latestEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTransmitterInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"lastCollected\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTriggerType\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepInfo\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"amountSpent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasDedupKey\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.MigrationPermission\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepOffchainConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upkeepTranscoderVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFormat\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"upkeepVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawOwnerFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainSpecificModuleUpdated\",\"inputs\":[{\"name\":\"newModule\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162004daf38038062004daf8339810160408190526200003591620001bf565b84848484843380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c481620000f7565b5050506001600160a01b0394851660805292841660a05290831660c052821660e0521661010052506200022f9350505050565b336001600160a01b03821603620001515760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ba57600080fd5b919050565b600080600080600060a08688031215620001d857600080fd5b620001e386620001a2565b9450620001f360208701620001a2565b93506200020360408701620001a2565b92506200021360608701620001a2565b91506200022360808701620001a2565b90509295509295909350565b60805160a05160c05160e05161010051614b0a620002a56000396000610715015260006105920152600081816105ff01526132df01526000818161077801526133b901526000818161080601528181611d490152818161201e015281816124870152818161299a0152612a1e0152614b0a6000f3fe608060405234801561001057600080fd5b50600436106103625760003560e01c80637d9b97e0116101c8578063b121e14711610104578063cd7f71b5116100a2578063ed56b3e11161007c578063ed56b3e114610863578063f2fde38b146108d6578063f777ff06146108e9578063faa3e996146108f057600080fd5b8063cd7f71b51461082a578063d76326481461083d578063eb5dcd6c1461085057600080fd5b8063b657bc9c116100de578063b657bc9c146107c9578063b79550be146107dc578063c7c3a19a146107e4578063ca30e6031461080457600080fd5b8063b121e1471461079c578063b148ab6b146107af578063b6511a2a146107c257600080fd5b80639e0a99ed11610171578063a72aa27e1161014b578063a72aa27e1461074c578063aab9edd61461075f578063abc76ae01461076e578063b10b673c1461077657600080fd5b80639e0a99ed1461070b578063a08714c014610713578063a710b2211461073957600080fd5b80638da5cb5b116101a25780638da5cb5b146106b75780638dcf0fe7146106d55780638ed02bab146106e857600080fd5b80637d9b97e0146106945780638456cb591461069c5780638765ecbe146106a457600080fd5b806343cc055c116102a25780635b6aa71c11610240578063671d36ed1161021a578063671d36ed14610623578063744bfe611461063657806379ba50971461064957806379ea99431461065157600080fd5b80635b6aa71c146105d75780636209e1e9146105ea5780636709d0e5146105fd57600080fd5b80634ca16c521161027c5780634ca16c52146105555780635147cd591461055d5780635165f2f51461057d5780635425d8ac1461059057600080fd5b806343cc055c1461050c57806344cb70b81461052357806348013d7b1461054657600080fd5b80631a2af0111161030f578063232c1cc5116102e9578063232c1cc5146104845780633b9cce591461048b5780633f4ba83a1461049e578063421d183b146104a657600080fd5b80631a2af011146104005780631e01043914610413578063207b65161461047157600080fd5b80631865c57d116103405780631865c57d146103b4578063187256e8146103cd57806319d97a94146103e057600080fd5b8063050ee65d1461036757806306e3b6321461037f5780630b7d33e61461039f575b600080fd5b62014c085b6040519081526020015b60405180910390f35b61039261038d366004613d35565b610936565b6040516103769190613d57565b6103b26103ad366004613de4565b610a53565b005b6103bc610b0d565b604051610376959493929190613fe7565b6103b26103db36600461411e565b610f57565b6103f36103ee36600461415b565b610fc8565b60405161037691906141d8565b6103b261040e3660046141eb565b61106a565b61045461042136600461415b565b6000908152600460205260409020600101546c0100000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff9091168152602001610376565b6103f361047f36600461415b565b611170565b601861036c565b6103b2610499366004614210565b61118d565b6103b26113e3565b6104b96104b4366004614285565b611449565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a001610376565b60135460ff165b6040519015158152602001610376565b61051361053136600461415b565b60009081526008602052604090205460ff1690565b600060405161037691906142d1565b61ea6061036c565b61057061056b36600461415b565b611568565b60405161037691906142eb565b6103b261058b36600461415b565b611573565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610376565b6104546105e5366004614318565b6116ea565b6103f36105f8366004614285565b61188f565b7f00000000000000000000000000000000000000000000000000000000000000006105b2565b6103b2610631366004614351565b6118c3565b6103b26106443660046141eb565b61199d565b6103b2611e44565b6105b261065f36600461415b565b6000908152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b6103b2611f46565b6103b26120a1565b6103b26106b236600461415b565b612122565b60005473ffffffffffffffffffffffffffffffffffffffff166105b2565b6103b26106e3366004613de4565b61229c565b601354610100900473ffffffffffffffffffffffffffffffffffffffff166105b2565b6103a461036c565b7f00000000000000000000000000000000000000000000000000000000000000006105b2565b6103b261074736600461438d565b6122f1565b6103b261075a3660046143bb565b612559565b60405160038152602001610376565b6115e061036c565b7f00000000000000000000000000000000000000000000000000000000000000006105b2565b6103b26107aa366004614285565b61264e565b6103b26107bd36600461415b565b612746565b603261036c565b6104546107d736600461415b565b612934565b6103b2612961565b6107f76107f236600461415b565b612abd565b60405161037691906143de565b7f00000000000000000000000000000000000000000000000000000000000000006105b2565b6103b2610838366004613de4565b612e90565b61045461084b36600461415b565b612f27565b6103b261085e36600461438d565b612f32565b6108bd610871366004614285565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff909116602083015201610376565b6103b26108e4366004614285565b613090565b604061036c565b6109296108fe366004614285565b73ffffffffffffffffffffffffffffffffffffffff166000908152601a602052604090205460ff1690565b6040516103769190614515565b6060600061094460026130a4565b905080841061097f576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061098b8486614558565b905081811180610999575083155b6109a357806109a5565b815b905060006109b3868361456b565b67ffffffffffffffff8111156109cb576109cb61457e565b6040519080825280602002602001820160405280156109f4578160200160208202803683370190505b50905060005b8151811015610a4757610a18610a108883614558565b6002906130ae565b828281518110610a2a57610a2a6145ad565b602090810291909101015280610a3f816145dc565b9150506109fa565b50925050505b92915050565b6016546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314610ab4576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152601d60205260409020610acd8284836146b6565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae77698383604051610b009291906147d1565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c0810191909152604080516101408101825260155463ffffffff7401000000000000000000000000000000000000000082041682526bffffffffffffffffffffffff90811660208301526019549282019290925260125490911660608281019190915290819060009060808101610c4660026130a4565b8152601554780100000000000000000000000000000000000000000000000080820463ffffffff9081166020808601919091527c01000000000000000000000000000000000000000000000000000000008404821660408087019190915260115460608088019190915260125474010000000000000000000000000000000000000000810485166080808a01919091527e01000000000000000000000000000000000000000000000000000000000000820460ff16151560a0998a015283516101e0810185526c0100000000000000000000000080840488168252700100000000000000000000000000000000808504891697830197909752808a0488169582019590955296820462ffffff16928701929092527b01000000000000000000000000000000000000000000000000000000900461ffff16908501526014546bffffffffffffffffffffffff8116968501969096529304811660c083015260165480821660e08401526401000000008104821661010084015268010000000000000000900416610120820152601754610140820152601854610160820152910473ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a08101610e1360096130c1565b81526016546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16602091820152601254600d80546040805182860281018601909152818152949850899489949293600e937d01000000000000000000000000000000000000000000000000000000000090910460ff16928591830182828015610ed657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610eab575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610f3f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f14575b50505050509150945094509450945094509091929394565b610f5f6130ce565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601a6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115610fbf57610fbf6142a2565b02179055505050565b6000818152601d60205260409020805460609190610fe590614614565b80601f016020809104026020016040519081016040528092919081815260200182805461101190614614565b801561105e5780601f106110335761010080835404028352916020019161105e565b820191906000526020600020905b81548152906001019060200180831161104157829003601f168201915b50505050509050919050565b61107382613151565b3373ffffffffffffffffffffffffffffffffffffffff8216036110c2576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff82811691161461116c5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b6000818152601b60205260409020805460609190610fe590614614565b6111956130ce565b600e5481146111d0576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e548110156113a2576000600e82815481106111f2576111f26145ad565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600f9092526040832054919350169085858581811061123c5761123c6145ad565b90506020020160208101906112519190614285565b905073ffffffffffffffffffffffffffffffffffffffff811615806112e4575073ffffffffffffffffffffffffffffffffffffffff8216158015906112c257508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156112e4575073ffffffffffffffffffffffffffffffffffffffff81811614155b1561131b576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181161461138c5773ffffffffffffffffffffffffffffffffffffffff8381166000908152600f6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061139a906145dc565b9150506111d3565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e83836040516113d79392919061481e565b60405180910390a15050565b6113eb6130ce565b601280547fff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152829182918291829190829061150f5760608201516012546000916114fb916bffffffffffffffffffffffff166148d0565b600e5490915061150b9082614924565b9150505b81516020830151604084015161152690849061494f565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b166000908152600f6020526040902054929b919a9499509750921694509092505050565b6000610a4d82613205565b61157c81613151565b600081815260046020908152604091829020825160e081018452815460ff8116151580835263ffffffff610100830481169584019590955265010000000000820485169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201529061167b576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556116ba6002836132b0565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e0100000000000000000000000000000000000000000000000000000000000083048116151560e08301527f01000000000000000000000000000000000000000000000000000000000000009092048216151561010080830191909152601354928316151561012083015273ffffffffffffffffffffffffffffffffffffffff9204919091166101408201526000908180611874836132bc565b91509150611885838787858561349a565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152601e60205260409020805460609190610fe590614614565b6016546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314611924576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152601e602052604090206119548284836146b6565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d28383604051610b009291906147d1565b6012547f0100000000000000000000000000000000000000000000000000000000000000900460ff16156119fd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f010000000000000000000000000000000000000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116611a93576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815160e081018352815460ff81161515825263ffffffff610100820481168387015265010000000000820481168386015273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909204821660608401526001909301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c082015286855260059093529220549091163314611b9a576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015611c0a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c2e9190614974565b816040015163ffffffff161115611c71576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546019546c010000000000000000000000009091046bffffffffffffffffffffffff1690611cb190829061456b565b60195560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015611d94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611db8919061498d565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611eca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611f4e6130ce565b6015546019546bffffffffffffffffffffffff90911690611f7090829061456b565b601955601580547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b6020604051808303816000875af115801561207d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061116c919061498d565b6120a96130ce565b601280547fff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e010000000000000000000000000000000000000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589060200161143f565b61212b81613151565b600081815260046020908152604091829020825160e081018452815460ff8116158015835263ffffffff610100830481169584019590955265010000000000820485169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201529061222a576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905561226c600283613722565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b6122a583613151565b6000838152601c602052604090206122be8284836146b6565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508383604051610b009291906147d1565b73ffffffffffffffffffffffffffffffffffffffff811661233e576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f602052604090205416331461239e576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546000916123c19185916bffffffffffffffffffffffff169061372e565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905560195490915061242b906bffffffffffffffffffffffff83169061456b565b6019556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af11580156124d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124f4919061498d565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff16108061258e575060155463ffffffff7001000000000000000000000000000000009091048116908216115b156125c5576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6125ce82613151565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff1661010063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152601060205260409020541633146126ae576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600f602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556010909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529114612843576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146128a0576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b6000610a4d61294283613205565b600084815260046020526040902054610100900463ffffffff166116ea565b6129696130ce565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156129f6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a1a9190614974565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360195484612a67919061456b565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff9092166004830152602482015260440161205e565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606083018190526001909101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c0820152919015612c5557816060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c2c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c5091906149af565b612c58565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836020015163ffffffff168152602001600760008781526020019081526020016000208054612cb090614614565b80601f0160208091040260200160405190810160405280929190818152602001828054612cdc90614614565b8015612d295780601f10612cfe57610100808354040283529160200191612d29565b820191906000526020600020905b815481529060010190602001808311612d0c57829003601f168201915b505050505081526020018360a001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff1667ffffffffffffffff1681526020018360c0015163ffffffff16815260200183608001516bffffffffffffffffffffffff168152602001836000015115158152602001601c60008781526020019081526020016000208054612e0690614614565b80601f0160208091040260200160405190810160405280929190818152602001828054612e3290614614565b8015612e7f5780601f10612e5457610100808354040283529160200191612e7f565b820191906000526020600020905b815481529060010190602001808311612e6257829003601f168201915b505050505081525092505050919050565b612e9983613151565b60165463ffffffff16811115612edb576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020612ef48284836146b6565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d8383604051610b009291906147d1565b6000610a4d82612934565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612f92576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603612fe1576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526010602052604090205481169082161461116c5773ffffffffffffffffffffffffffffffffffffffff82811660008181526010602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6130986130ce565b6130a181613936565b50565b6000610a4d825490565b60006130ba8383613a2b565b9392505050565b606060006130ba83613a55565b60005473ffffffffffffffffffffffffffffffffffffffff16331461314f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611ec1565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146131ae576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff908116146130a1576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015613292577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061324a5761324a6145ad565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461328057506000949350505050565b8061328a816145dc565b91505061320c565b5081600f1a60018111156132a8576132a86142a2565b949350505050565b60006130ba8383613ab0565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613348573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061336c91906149e6565b509450909250505060008113158061338357508142105b806133a457508280156133a4575061339b824261456b565b8463ffffffff16105b156133b35760175495506133b7565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613422573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061344691906149e6565b509450909250505060008113158061345d57508142105b8061347e575082801561347e5750613475824261456b565b8463ffffffff16105b1561348d576018549450613491565b8094505b50505050915091565b600080808660018111156134b0576134b06142a2565b036134be575061ea60613513565b60018660018111156134d2576134d26142a2565b036134e1575062014c08613513565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008760c0015160016135269190614a36565b6135349060ff166040614a4f565b601654613552906103a490640100000000900463ffffffff16614558565b61355c9190614558565b601354604080517fde9ee35e00000000000000000000000000000000000000000000000000000000815281519394506000938493610100900473ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa1580156135d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135f69190614a66565b90925090508183613608836018614558565b6136129190614a4f565b60c08c0151613622906001614a36565b6136319060ff166115e0614a4f565b61363b9190614558565b6136459190614558565b61364f9085614558565b935060008a610140015173ffffffffffffffffffffffffffffffffffffffff166312544140856040518263ffffffff1660e01b815260040161369391815260200190565b602060405180830381865afa1580156136b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136d49190614974565b8b60a0015161ffff166136e79190614a4f565b90506000806137028d8c63ffffffff1689868e8e6000613aff565b9092509050613711818361494f565b9d9c50505050505050505050505050565b60006130ba8383613c3b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e010000000000000000000000000000900490911660608201529061392a5760008160600151856137c691906148d0565b905060006137d48583614924565b905080836040018181516137e8919061494f565b6bffffffffffffffffffffffff169052506138038582614a8a565b83606001818151613814919061494f565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036139b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611ec1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613a4257613a426145ad565b9060005260206000200154905092915050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561105e57602002820191906000526020600020905b815481526020019060010190808311613a915750505050509050919050565b6000818152600183016020526040812054613af757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610a4d565b506000610a4d565b60008060008960a0015161ffff1686613b189190614a4f565b9050838015613b265750803a105b15613b2e57503a5b60008588613b3c8b8d614558565b613b469085614a4f565b613b509190614558565b613b6290670de0b6b3a7640000614a4f565b613b6c9190614aba565b905060008b6040015163ffffffff1664e8d4a51000613b8b9190614a4f565b60208d0151889063ffffffff168b613ba38f88614a4f565b613bad9190614558565b613bbb90633b9aca00614a4f565b613bc59190614a4f565b613bcf9190614aba565b613bd99190614558565b90506b033b2e3c9fd0803ce8000000613bf28284614558565b1115613c2a576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909b909a5098505050505050505050565b60008181526001830160205260408120548015613d24576000613c5f60018361456b565b8554909150600090613c739060019061456b565b9050818114613cd8576000866000018281548110613c9357613c936145ad565b9060005260206000200154905080876000018481548110613cb657613cb66145ad565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613ce957613ce9614ace565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610a4d565b6000915050610a4d565b5092915050565b60008060408385031215613d4857600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613d8f57835183529284019291840191600101613d73565b50909695505050505050565b60008083601f840112613dad57600080fd5b50813567ffffffffffffffff811115613dc557600080fd5b602083019150836020828501011115613ddd57600080fd5b9250929050565b600080600060408486031215613df957600080fd5b83359250602084013567ffffffffffffffff811115613e1757600080fd5b613e2386828701613d9b565b9497909650939450505050565b600081518084526020808501945080840160005b83811015613e7657815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613e44565b509495945050505050565b805163ffffffff16825260006101e06020830151613ea7602086018263ffffffff169052565b506040830151613ebf604086018263ffffffff169052565b506060830151613ed6606086018262ffffff169052565b506080830151613eec608086018261ffff169052565b5060a0830151613f0c60a08601826bffffffffffffffffffffffff169052565b5060c0830151613f2460c086018263ffffffff169052565b5060e0830151613f3c60e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a080840151818601839052613fb183870182613e30565b925050506101c080840151613fdd8287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161401560208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161403f60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a088015161406160a085018263ffffffff169052565b5060c088015161407960c085018263ffffffff169052565b5060e088015160e08401526101008089015161409c8286018263ffffffff169052565b50506101208881015115159084015261014083018190526140bf81840188613e81565b90508281036101608401526140d48187613e30565b90508281036101808401526140e98186613e30565b9150506118856101a083018460ff169052565b73ffffffffffffffffffffffffffffffffffffffff811681146130a157600080fd5b6000806040838503121561413157600080fd5b823561413c816140fc565b915060208301356004811061415057600080fd5b809150509250929050565b60006020828403121561416d57600080fd5b5035919050565b6000815180845260005b8181101561419a5760208185018101518683018201520161417e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006130ba6020830184614174565b600080604083850312156141fe57600080fd5b823591506020830135614150816140fc565b6000806020838503121561422357600080fd5b823567ffffffffffffffff8082111561423b57600080fd5b818501915085601f83011261424f57600080fd5b81358181111561425e57600080fd5b8660208260051b850101111561427357600080fd5b60209290920196919550909350505050565b60006020828403121561429757600080fd5b81356130ba816140fc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60208101600383106142e5576142e56142a2565b91905290565b60208101600283106142e5576142e56142a2565b803563ffffffff8116811461431357600080fd5b919050565b6000806040838503121561432b57600080fd5b82356002811061433a57600080fd5b9150614348602084016142ff565b90509250929050565b60008060006040848603121561436657600080fd5b8335614371816140fc565b9250602084013567ffffffffffffffff811115613e1757600080fd5b600080604083850312156143a057600080fd5b82356143ab816140fc565b91506020830135614150816140fc565b600080604083850312156143ce57600080fd5b82359150614348602084016142ff565b6020815261440560208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161441e604084018263ffffffff169052565b50604083015161014080606085015261443b610160850183614174565b9150606085015161445c60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e08501516101006144c8818701836bffffffffffffffffffffffff169052565b86015190506101206144dd8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506118858382614174565b60208101600483106142e5576142e56142a2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115610a4d57610a4d614529565b81810381811115610a4d57610a4d614529565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361460d5761460d614529565b5060010190565b600181811c9082168061462857607f821691505b602082108103614661577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156146b157600081815260208120601f850160051c8101602086101561468e5750805b601f850160051c820191505b818110156146ad5782815560010161469a565b5050505b505050565b67ffffffffffffffff8311156146ce576146ce61457e565b6146e2836146dc8354614614565b83614667565b6000601f84116001811461473457600085156146fe5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556147ca565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156147835786850135825560209485019460019092019101614763565b50868210156147be577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000604082016040835280865480835260608501915087600052602092508260002060005b8281101561487557815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101614843565b505050838103828501528481528590820160005b868110156148c457823561489c816140fc565b73ffffffffffffffffffffffffffffffffffffffff1682529183019190830190600101614889565b50979650505050505050565b6bffffffffffffffffffffffff828116828216039080821115613d2e57613d2e614529565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680614943576149436148f5565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115613d2e57613d2e614529565b60006020828403121561498657600080fd5b5051919050565b60006020828403121561499f57600080fd5b815180151581146130ba57600080fd5b6000602082840312156149c157600080fd5b81516130ba816140fc565b805169ffffffffffffffffffff8116811461431357600080fd5b600080600080600060a086880312156149fe57600080fd5b614a07866149cc565b9450602086015193506040860151925060608601519150614a2a608087016149cc565b90509295509295909350565b60ff8181168382160190811115610a4d57610a4d614529565b8082028115828204841417610a4d57610a4d614529565b60008060408385031215614a7957600080fd5b505080516020909101519092909150565b6bffffffffffffffffffffffff818116838216028082169190828114614ab257614ab2614529565b505092915050565b600082614ac957614ac96148f5565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var AutomationRegistryLogicB22ABI = AutomationRegistryLogicB22MetaData.ABI

var AutomationRegistryLogicB22Bin = AutomationRegistryLogicB22MetaData.Bin

func DeployAutomationRegistryLogicB22(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address, automationForwarderLogic common.Address, allowedReadOnlyAddress common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicB22, error) {
	parsed, err := AutomationRegistryLogicB22MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicB22Bin), backend, link, linkNativeFeed, fastGasFeed, automationForwarderLogic, allowedReadOnlyAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicB22{address: address, abi: *parsed, AutomationRegistryLogicB22Caller: AutomationRegistryLogicB22Caller{contract: contract}, AutomationRegistryLogicB22Transactor: AutomationRegistryLogicB22Transactor{contract: contract}, AutomationRegistryLogicB22Filterer: AutomationRegistryLogicB22Filterer{contract: contract}}, nil
}

type AutomationRegistryLogicB22 struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicB22Caller
	AutomationRegistryLogicB22Transactor
	AutomationRegistryLogicB22Filterer
}

type AutomationRegistryLogicB22Caller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB22Transactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB22Filterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB22Session struct {
	Contract     *AutomationRegistryLogicB22
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicB22CallerSession struct {
	Contract *AutomationRegistryLogicB22Caller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicB22TransactorSession struct {
	Contract     *AutomationRegistryLogicB22Transactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicB22Raw struct {
	Contract *AutomationRegistryLogicB22
}

type AutomationRegistryLogicB22CallerRaw struct {
	Contract *AutomationRegistryLogicB22Caller
}

type AutomationRegistryLogicB22TransactorRaw struct {
	Contract *AutomationRegistryLogicB22Transactor
}

func NewAutomationRegistryLogicB22(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicB22, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicB22ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicB22(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22{address: address, abi: abi, AutomationRegistryLogicB22Caller: AutomationRegistryLogicB22Caller{contract: contract}, AutomationRegistryLogicB22Transactor: AutomationRegistryLogicB22Transactor{contract: contract}, AutomationRegistryLogicB22Filterer: AutomationRegistryLogicB22Filterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicB22Caller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicB22Caller, error) {
	contract, err := bindAutomationRegistryLogicB22(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22Caller{contract: contract}, nil
}

func NewAutomationRegistryLogicB22Transactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicB22Transactor, error) {
	contract, err := bindAutomationRegistryLogicB22(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22Transactor{contract: contract}, nil
}

func NewAutomationRegistryLogicB22Filterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicB22Filterer, error) {
	contract, err := bindAutomationRegistryLogicB22(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22Filterer{contract: contract}, nil
}

func bindAutomationRegistryLogicB22(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicB22MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB22.Contract.AutomationRegistryLogicB22Caller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AutomationRegistryLogicB22Transactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AutomationRegistryLogicB22Transactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB22.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicB22.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicB22.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicB22.CallOpts, admin)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicB22.CallOpts, admin)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetBalance(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetBalance(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetCancellationDelay(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetCancellationDelay(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetChainModule(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetChainModule(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetForwarder(&_AutomationRegistryLogicB22.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetForwarder(&_AutomationRegistryLogicB22.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetLinkAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetLinkAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetLinkNativeFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetLinkNativeFeedAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.GetLinkNativeFeedAddress(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetLogGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetLogGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicB22.CallOpts, triggerType, gasLimit)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicB22.CallOpts, triggerType, gasLimit)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMinBalance(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMinBalance(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB22.CallOpts, peer)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB22.CallOpts, peer)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicB22.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicB22.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetSignerInfo(&_AutomationRegistryLogicB22.CallOpts, query)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetSignerInfo(&_AutomationRegistryLogicB22.CallOpts, query)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getState")

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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetState(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetState(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getTransmitterInfo", query)

	outstruct := new(GetTransmitterInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastCollected = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Payee = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitterInfo(&_AutomationRegistryLogicB22.CallOpts, query)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicB22.Contract.GetTransmitterInfo(&_AutomationRegistryLogicB22.CallOpts, query)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.GetTriggerType(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.GetTriggerType(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeep(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeep(&_AutomationRegistryLogicB22.CallOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB22.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicB22.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicB22.Contract.HasDedupKey(&_AutomationRegistryLogicB22.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicB22.Contract.HasDedupKey(&_AutomationRegistryLogicB22.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.Owner(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB22.Contract.Owner(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) UpkeepTranscoderVersion() (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.UpkeepTranscoderVersion(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.UpkeepTranscoderVersion(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Caller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB22.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.UpkeepVersion(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22CallerSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicB22.Contract.UpkeepVersion(&_AutomationRegistryLogicB22.CallOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptOwnership(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptOwnership(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptPayeeship(&_AutomationRegistryLogicB22.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptPayeeship(&_AutomationRegistryLogicB22.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "pause")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.Pause(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.Pause(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "pauseUpkeep", id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.PauseUpkeep(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.PauseUpkeep(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "recoverFunds")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) RecoverFunds() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.RecoverFunds(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.RecoverFunds(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicB22.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicB22.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setPayees", payees)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetPayees(&_AutomationRegistryLogicB22.TransactOpts, payees)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetPayees(&_AutomationRegistryLogicB22.TransactOpts, payees)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB22.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB22.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB22.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB22.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB22.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB22.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB22.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB22.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB22.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB22.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferOwnership(&_AutomationRegistryLogicB22.TransactOpts, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferOwnership(&_AutomationRegistryLogicB22.TransactOpts, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferPayeeship(&_AutomationRegistryLogicB22.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferPayeeship(&_AutomationRegistryLogicB22.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB22.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB22.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "unpause")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.Unpause(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.Unpause(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB22.TransactOpts, id)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawFunds(&_AutomationRegistryLogicB22.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawFunds(&_AutomationRegistryLogicB22.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawOwnerFunds(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawOwnerFunds(&_AutomationRegistryLogicB22.TransactOpts)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawPayment(&_AutomationRegistryLogicB22.TransactOpts, from, to)
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB22.Contract.WithdrawPayment(&_AutomationRegistryLogicB22.TransactOpts, from, to)
}

type AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicB22AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22AdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicB22AdminPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22AdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB22AdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicB22AdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22CancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB22CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22CancelledUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicB22CancelledUpkeepReport)
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

func (it *AutomationRegistryLogicB22CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22CancelledUpkeepReportIterator{contract: _AutomationRegistryLogicB22.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22CancelledUpkeepReport)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicB22CancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicB22CancelledUpkeepReport)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicB22ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22ChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistryLogicB22ChainSpecificModuleUpdated)
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

func (it *AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicB22.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22ChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicB22ChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicB22ChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22DedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicB22DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22DedupKeyAdded)
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
		it.Event = new(AutomationRegistryLogicB22DedupKeyAdded)
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

func (it *AutomationRegistryLogicB22DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicB22DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22DedupKeyAddedIterator{contract: _AutomationRegistryLogicB22.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22DedupKeyAdded)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicB22DedupKeyAdded, error) {
	event := new(AutomationRegistryLogicB22DedupKeyAdded)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22FundsAddedIterator struct {
	Event *AutomationRegistryLogicB22FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22FundsAdded)
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
		it.Event = new(AutomationRegistryLogicB22FundsAdded)
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

func (it *AutomationRegistryLogicB22FundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicB22FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22FundsAddedIterator{contract: _AutomationRegistryLogicB22.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22FundsAdded)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicB22FundsAdded, error) {
	event := new(AutomationRegistryLogicB22FundsAdded)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22FundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicB22FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22FundsWithdrawn)
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
		it.Event = new(AutomationRegistryLogicB22FundsWithdrawn)
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

func (it *AutomationRegistryLogicB22FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22FundsWithdrawnIterator{contract: _AutomationRegistryLogicB22.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22FundsWithdrawn)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB22FundsWithdrawn, error) {
	event := new(AutomationRegistryLogicB22FundsWithdrawn)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB22InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22InsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicB22InsufficientFundsUpkeepReport)
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

func (it *AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicB22.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22InsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicB22InsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicB22InsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22OwnerFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicB22OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22OwnerFundsWithdrawn)
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
		it.Event = new(AutomationRegistryLogicB22OwnerFundsWithdrawn)
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

func (it *AutomationRegistryLogicB22OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicB22OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22OwnerFundsWithdrawnIterator{contract: _AutomationRegistryLogicB22.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22OwnerFundsWithdrawn)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB22OwnerFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicB22OwnerFundsWithdrawn)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22OwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB22OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22OwnershipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicB22OwnershipTransferRequested)
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

func (it *AutomationRegistryLogicB22OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22OwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicB22.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22OwnershipTransferRequested)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicB22OwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicB22OwnershipTransferRequested)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22OwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicB22OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22OwnershipTransferred)
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
		it.Event = new(AutomationRegistryLogicB22OwnershipTransferred)
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

func (it *AutomationRegistryLogicB22OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22OwnershipTransferredIterator{contract: _AutomationRegistryLogicB22.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22OwnershipTransferred)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicB22OwnershipTransferred, error) {
	event := new(AutomationRegistryLogicB22OwnershipTransferred)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22PausedIterator struct {
	Event *AutomationRegistryLogicB22Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22Paused)
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
		it.Event = new(AutomationRegistryLogicB22Paused)
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

func (it *AutomationRegistryLogicB22PausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB22PausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22PausedIterator{contract: _AutomationRegistryLogicB22.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22Paused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22Paused)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParsePaused(log types.Log) (*AutomationRegistryLogicB22Paused, error) {
	event := new(AutomationRegistryLogicB22Paused)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22PayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicB22PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22PayeesUpdated)
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
		it.Event = new(AutomationRegistryLogicB22PayeesUpdated)
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

func (it *AutomationRegistryLogicB22PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB22PayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22PayeesUpdatedIterator{contract: _AutomationRegistryLogicB22.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22PayeesUpdated)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicB22PayeesUpdated, error) {
	event := new(AutomationRegistryLogicB22PayeesUpdated)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22PayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB22PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22PayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicB22PayeeshipTransferRequested)
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

func (it *AutomationRegistryLogicB22PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22PayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicB22.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22PayeeshipTransferRequested)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicB22PayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicB22PayeeshipTransferRequested)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22PayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicB22PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22PayeeshipTransferred)
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
		it.Event = new(AutomationRegistryLogicB22PayeeshipTransferred)
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

func (it *AutomationRegistryLogicB22PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22PayeeshipTransferredIterator{contract: _AutomationRegistryLogicB22.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22PayeeshipTransferred)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicB22PayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicB22PayeeshipTransferred)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22PaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicB22PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22PaymentWithdrawn)
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
		it.Event = new(AutomationRegistryLogicB22PaymentWithdrawn)
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

func (it *AutomationRegistryLogicB22PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicB22PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22PaymentWithdrawnIterator{contract: _AutomationRegistryLogicB22.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22PaymentWithdrawn)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicB22PaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicB22PaymentWithdrawn)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22ReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB22ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22ReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicB22ReorgedUpkeepReport)
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

func (it *AutomationRegistryLogicB22ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22ReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicB22.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22ReorgedUpkeepReport)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicB22ReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicB22ReorgedUpkeepReport)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22StaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB22StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22StaleUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicB22StaleUpkeepReport)
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

func (it *AutomationRegistryLogicB22StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22StaleUpkeepReportIterator{contract: _AutomationRegistryLogicB22.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22StaleUpkeepReport)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicB22StaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicB22StaleUpkeepReport)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UnpausedIterator struct {
	Event *AutomationRegistryLogicB22Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22Unpaused)
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
		it.Event = new(AutomationRegistryLogicB22Unpaused)
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

func (it *AutomationRegistryLogicB22UnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB22UnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UnpausedIterator{contract: _AutomationRegistryLogicB22.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22Unpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22Unpaused)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicB22Unpaused, error) {
	event := new(AutomationRegistryLogicB22Unpaused)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepAdminTransferRequested)
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

func (it *AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicB22UpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicB22UpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicB22UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepAdminTransferred)
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

func (it *AutomationRegistryLogicB22UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepAdminTransferred)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicB22UpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicB22UpkeepAdminTransferred)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicB22UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepCanceled)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepCanceled)
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

func (it *AutomationRegistryLogicB22UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicB22UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepCanceledIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepCanceled)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicB22UpkeepCanceled, error) {
	event := new(AutomationRegistryLogicB22UpkeepCanceled)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicB22UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepCheckDataSet)
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

func (it *AutomationRegistryLogicB22UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepCheckDataSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicB22UpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicB22UpkeepCheckDataSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicB22UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepGasLimitSet)
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

func (it *AutomationRegistryLogicB22UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepGasLimitSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicB22UpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicB22UpkeepGasLimitSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepMigrated)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepMigrated)
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

func (it *AutomationRegistryLogicB22UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepMigratedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepMigrated)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicB22UpkeepMigrated, error) {
	event := new(AutomationRegistryLogicB22UpkeepMigrated)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicB22UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepOffchainConfigSet)
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

func (it *AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicB22UpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepPausedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepPaused)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepPaused)
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

func (it *AutomationRegistryLogicB22UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepPausedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepPaused)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicB22UpkeepPaused, error) {
	event := new(AutomationRegistryLogicB22UpkeepPaused)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepPerformed)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepPerformed)
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

func (it *AutomationRegistryLogicB22UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicB22UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepPerformedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepPerformed)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicB22UpkeepPerformed, error) {
	event := new(AutomationRegistryLogicB22UpkeepPerformed)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicB22UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicB22UpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepReceived)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepReceived)
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

func (it *AutomationRegistryLogicB22UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepReceivedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepReceived)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicB22UpkeepReceived, error) {
	event := new(AutomationRegistryLogicB22UpkeepReceived)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicB22UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepRegistered)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepRegistered)
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

func (it *AutomationRegistryLogicB22UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepRegisteredIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepRegistered)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicB22UpkeepRegistered, error) {
	event := new(AutomationRegistryLogicB22UpkeepRegistered)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicB22UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepTriggerConfigSet)
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

func (it *AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicB22UpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB22UpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicB22UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB22UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB22UpkeepUnpaused)
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
		it.Event = new(AutomationRegistryLogicB22UpkeepUnpaused)
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

func (it *AutomationRegistryLogicB22UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB22UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB22UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB22UpkeepUnpausedIterator{contract: _AutomationRegistryLogicB22.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB22.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB22UpkeepUnpaused)
				if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22Filterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicB22UpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicB22UpkeepUnpaused)
	if err := _AutomationRegistryLogicB22.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSignerInfo struct {
	Active bool
	Index  uint8
}
type GetState struct {
	State        IAutomationV21PlusCommonStateLegacy
	Config       IAutomationV21PlusCommonOnchainConfigLegacy
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}
type GetTransmitterInfo struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicB22.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB22.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB22.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicB22.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicB22.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicB22.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicB22.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicB22.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicB22.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicB22.ParseFundsAdded(log)
	case _AutomationRegistryLogicB22.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicB22.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicB22.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicB22.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicB22.abi.Events["OwnerFundsWithdrawn"].ID:
		return _AutomationRegistryLogicB22.ParseOwnerFundsWithdrawn(log)
	case _AutomationRegistryLogicB22.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicB22.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicB22.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicB22.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicB22.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicB22.ParsePaused(log)
	case _AutomationRegistryLogicB22.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicB22.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicB22.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicB22.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicB22.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicB22.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicB22.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicB22.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicB22.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicB22.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicB22.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicB22.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicB22.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicB22.ParseUnpaused(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicB22.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicB22.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicB22AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicB22CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicB22ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicB22DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicB22FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicB22FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicB22InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicB22OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (AutomationRegistryLogicB22OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicB22OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicB22Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicB22PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicB22PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicB22PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicB22PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicB22ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicB22StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicB22Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicB22UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicB22UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicB22UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicB22UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicB22UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicB22UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicB22UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicB22UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicB22UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicB22UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicB22UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicB22UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicB22UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicB22UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicB22 *AutomationRegistryLogicB22) Address() common.Address {
	return _AutomationRegistryLogicB22.address
}

type AutomationRegistryLogicB22Interface interface {
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error)

	GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error)

	GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error)

	GetChainModule(opts *bind.CallOpts) (common.Address, error)

	GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error)

	GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error)

	GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

		error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

		error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicB22AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB22AdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicB22CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB22ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicB22ChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicB22DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicB22DedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicB22FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicB22FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB22FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicB22InsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicB22OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB22OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicB22OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicB22OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB22PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicB22Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB22PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicB22PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicB22PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicB22PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicB22PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicB22PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicB22ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicB22StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB22UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicB22Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicB22UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB22UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicB22UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicB22UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicB22UpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicB22UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicB22UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicB22UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicB22UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicB22UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicB22UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicB22UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicB22UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicB22UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB22UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB22UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicB22UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
