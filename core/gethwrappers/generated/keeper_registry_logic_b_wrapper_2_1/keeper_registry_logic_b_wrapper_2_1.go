// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic_b_wrapper_2_1

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

var KeeperRegistryLogicB21MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.Mode\"},{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"linkNativeFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fastGasFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"automationForwarderLogic\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveUpkeepIDs\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAutomationForwarderLogic\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCancellationDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getConditionalGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getFastGasFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForwarder\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAutomationForwarder\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkNativeFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLogGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getMaxPaymentForGas\",\"inputs\":[{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"maxPayment\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalanceForUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"minBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMode\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.Mode\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPerPerformByteGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPerSignerGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getSignerInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ownerLinkBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"expectedLinkBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalPremium\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"numUpkeeps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"latestEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitterInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"lastCollected\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTriggerType\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepInfo\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"amountSpent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasDedupKey\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepOffchainConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upkeepTranscoderVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFormat\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"upkeepVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawOwnerFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162004fbc38038062004fbc8339810160408190526200003591620001e9565b84848484843380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c48162000121565b505050846002811115620000dc57620000dc6200025e565b60e0816002811115620000f357620000f36200025e565b9052506001600160a01b0393841660805291831660a052821660c05216610100525062000274945050505050565b336001600160a01b038216036200017b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001e457600080fd5b919050565b600080600080600060a086880312156200020257600080fd5b8551600381106200021257600080fd5b94506200022260208701620001cc565b93506200023260408701620001cc565b92506200024260608701620001cc565b91506200025260808701620001cc565b90509295509295909350565b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e05161010051614cbd620002ff600039600061058701526000818161052501528181613352015281816138d50152613a680152600081816105f4015261314601526000818161071c01526132200152600081816107aa01528181611bab01528181611e81015281816122ee0152818161280101526128850152614cbd6000f3fe608060405234801561001057600080fd5b50600436106103365760003560e01c806379ba5097116101b2578063b121e147116100f9578063ca30e603116100a2578063eb5dcd6c1161007c578063eb5dcd6c146107f4578063ed56b3e114610807578063f2fde38b1461087a578063faa3e9961461088d57600080fd5b8063ca30e603146107a8578063cd7f71b5146107ce578063d7632648146107e157600080fd5b8063b657bc9c116100d3578063b657bc9c1461076d578063b79550be14610780578063c7c3a19a1461078857600080fd5b8063b121e14714610740578063b148ab6b14610753578063b6511a2a1461076657600080fd5b80638dcf0fe71161015b578063aab9edd611610135578063aab9edd614610703578063abc76ae014610712578063b10b673c1461071a57600080fd5b80638dcf0fe7146106ca578063a710b221146106dd578063a72aa27e146106f057600080fd5b80638456cb591161018c5780638456cb59146106915780638765ecbe146106995780638da5cb5b146106ac57600080fd5b806379ba50971461063e57806379ea9943146106465780637d9b97e01461068957600080fd5b8063421d183b116102815780635165f2f51161022a5780636209e1e9116102045780636209e1e9146105df5780636709d0e5146105f2578063671d36ed14610618578063744bfe611461062b57600080fd5b80635165f2f5146105725780635425d8ac146105855780635b6aa71c146105cc57600080fd5b80634b4fd03b1161025b5780634b4fd03b146105235780634ca16c52146105495780635147cd591461055257600080fd5b8063421d183b1461047a57806344cb70b8146104e057806348013d7b1461051357600080fd5b80631a2af011116102e3578063232c1cc5116102bd578063232c1cc5146104585780633b9cce591461045f5780633f4ba83a1461047257600080fd5b80631a2af011146103d45780631e010439146103e7578063207b65161461044557600080fd5b80631865c57d116103145780631865c57d14610388578063187256e8146103a157806319d97a94146103b457600080fd5b8063050ee65d1461033b57806306e3b632146103535780630b7d33e614610373575b600080fd5b6201adb05b6040519081526020015b60405180910390f35b610366610361366004613df3565b6108d3565b60405161034a9190613e15565b610386610381366004613ea2565b6109f0565b005b610390610aaa565b60405161034a9594939291906140a5565b6103866103af3660046141dc565b610ec3565b6103c76103c2366004614219565b610f34565b60405161034a91906142a0565b6103866103e23660046142b3565b610fd6565b6104286103f5366004614219565b6000908152600460205260409020600101546c0100000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff909116815260200161034a565b6103c7610453366004614219565b6110dc565b6014610340565b61038661046d3660046142d8565b6110f9565b61038661134f565b61048d61048836600461434d565b6113b5565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a00161034a565b6105036104ee366004614219565b60009081526008602052604090205460ff1690565b604051901515815260200161034a565b60005b60405161034a91906143a9565b7f0000000000000000000000000000000000000000000000000000000000000000610516565b62015f90610340565b610565610560366004614219565b6114e8565b60405161034a91906143bc565b610386610580366004614219565b6114f3565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161034a565b6104286105da3660046143e9565b61166a565b6103c76105ed36600461434d565b61179c565b7f00000000000000000000000000000000000000000000000000000000000000006105a7565b610386610626366004614422565b6117d0565b6103866106393660046142b3565b6118aa565b610386611ca7565b6105a7610654366004614219565b6000908152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b610386611da9565b610386611f04565b6103866106a7366004614219565b611f75565b60005473ffffffffffffffffffffffffffffffffffffffff166105a7565b6103866106d8366004613ea2565b6120ef565b6103866106eb36600461445e565b612144565b6103866106fe36600461448c565b6123c0565b6040516003815260200161034a565b611d4c610340565b7f00000000000000000000000000000000000000000000000000000000000000006105a7565b61038661074e36600461434d565b6124b5565b610386610761366004614219565b6125ad565b6032610340565b61042861077b366004614219565b61279b565b6103866127c8565b61079b610796366004614219565b612924565b60405161034a91906144af565b7f00000000000000000000000000000000000000000000000000000000000000006105a7565b6103866107dc366004613ea2565b612cf7565b6104286107ef366004614219565b612d8e565b61038661080236600461445e565b612d99565b61086161081536600461434d565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff90911660208301520161034a565b61038661088836600461434d565b612ef7565b6108c661089b36600461434d565b73ffffffffffffffffffffffffffffffffffffffff1660009081526019602052604090205460ff1690565b60405161034a91906145e6565b606060006108e16002612f0b565b905080841061091c576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006109288486614629565b905081811180610936575083155b6109405780610942565b815b90506000610950868361463c565b67ffffffffffffffff8111156109685761096861464f565b604051908082528060200260200182016040528015610991578160200160208202803683370190505b50905060005b81518110156109e4576109b56109ad8883614629565b600290612f15565b8282815181106109c7576109c761467e565b6020908102919091010152806109dc816146ad565b915050610997565b50925050505b92915050565b6015546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314610a51576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152601c60205260409020610a6a828483614787565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae77698383604051610a9d9291906148a2565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c0810191909152604080516101408101825260145463ffffffff7401000000000000000000000000000000000000000082041682526bffffffffffffffffffffffff908116602083015260185492820192909252601254700100000000000000000000000000000000900490911660608281019190915290819060009060808101610bf76002612f0b565b81526014547801000000000000000000000000000000000000000000000000810463ffffffff9081166020808501919091527c0100000000000000000000000000000000000000000000000000000000808404831660408087019190915260115460608088019190915260125492830485166080808901919091526e010000000000000000000000000000840460ff16151560a09889015282516101e081018452610100808604881682526501000000000086048816968201969096526c010000000000000000000000008089048816948201949094526901000000000000000000850462ffffff16928101929092529282900461ffff16928101929092526013546bffffffffffffffffffffffff811696830196909652700100000000000000000000000000000000909404831660c082015260155480841660e0830152640100000000810484169282019290925268010000000000000000909104909116610120820152601654610140820152601754610160820152910473ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a08101610d9f6009612f28565b81526015546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16602091820152601254600d80546040805182860281018601909152818152949850899489949293600e9360ff909116928591830182828015610e4257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e17575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610eab57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e80575b50505050509150945094509450945094509091929394565b610ecb612f35565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260196020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115610f2b57610f2b61436a565b02179055505050565b6000818152601c60205260409020805460609190610f51906146e5565b80601f0160208091040260200160405190810160405280929190818152602001828054610f7d906146e5565b8015610fca5780601f10610f9f57610100808354040283529160200191610fca565b820191906000526020600020905b815481529060010190602001808311610fad57829003601f168201915b50505050509050919050565b610fdf82612fb8565b3373ffffffffffffffffffffffffffffffffffffffff82160361102e576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116146110d85760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b6000818152601a60205260409020805460609190610f51906146e5565b611101612f35565b600e54811461113c576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e5481101561130e576000600e828154811061115e5761115e61467e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600f909252604083205491935016908585858181106111a8576111a861467e565b90506020020160208101906111bd919061434d565b905073ffffffffffffffffffffffffffffffffffffffff81161580611250575073ffffffffffffffffffffffffffffffffffffffff82161580159061122e57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015611250575073ffffffffffffffffffffffffffffffffffffffff81811614155b15611287576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146112f85773ffffffffffffffffffffffffffffffffffffffff8381166000908152600f6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b5050508080611306906146ad565b91505061113f565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e8383604051611343939291906148ef565b60405180910390a15050565b611357612f35565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152829182918291829190829061148f57606082015160125460009161147b9170010000000000000000000000000000000090046bffffffffffffffffffffffff166149a1565b600e5490915061148b90826149f5565b9150505b8151602083015160408401516114a6908490614a20565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b166000908152600f6020526040902054929b919a9499509750921694509092505050565b60006109ea8261306c565b6114fc81612fb8565b600081815260046020908152604091829020825160e081018452815460ff8116151580835263ffffffff610100830481169584019590955265010000000000820485169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c0820152906115fb576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905561163a600283613117565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b604080516101208101825260125460ff808216835263ffffffff6101008084048216602086015265010000000000840482169585019590955262ffffff6901000000000000000000840416606085015261ffff6c0100000000000000000000000084041660808501526e01000000000000000000000000000083048216151560a08501526f010000000000000000000000000000008304909116151560c08401526bffffffffffffffffffffffff70010000000000000000000000000000000083041660e08401527c01000000000000000000000000000000000000000000000000000000009091041691810191909152600090818061176983613123565b91509150611792838787601360020160049054906101000a900463ffffffff1686866000613301565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152601d60205260409020805460609190610f51906146e5565b6015546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314611831576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152601d60205260409020611861828483614787565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d28383604051610a9d9291906148a2565b6012546f01000000000000000000000000000000900460ff16156118fa576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116611981576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815160e081018352815460ff81161515825263ffffffff610100820481168387015265010000000000820481168386015273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909204821660608401526001909301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c082015286855260059093529220549091163314611a88576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611a9061334c565b816040015163ffffffff161115611ad3576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546018546c010000000000000000000000009091046bffffffffffffffffffffffff1690611b1390829061463c565b60185560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015611bf6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c1a9190614a45565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d2d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611db1612f35565b6014546018546bffffffffffffffffffffffff90911690611dd390829061463c565b601855601480547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b6020604051808303816000875af1158015611ee0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110d89190614a45565b611f0c612f35565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020016113ab565b611f7e81612fb8565b600081815260046020908152604091829020825160e081018452815460ff8116158015835263ffffffff610100830481169584019590955265010000000000820485169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201529061207d576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556120bf600283613401565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b6120f883612fb8565b6000838152601b60205260409020612111828483614787565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508383604051610a9d9291906148a2565b73ffffffffffffffffffffffffffffffffffffffff8116612191576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f60205260409020541633146121f1576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e5460009161222891859170010000000000000000000000000000000090046bffffffffffffffffffffffff169061340d565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601854909150612292906bffffffffffffffffffffffff83169061463c565b6018556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015612337573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235b9190614a45565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff1610806123f5575060145463ffffffff7001000000000000000000000000000000009091048116908216115b1561242c576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61243582612fb8565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff1661010063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260106020526040902054163314612515576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600f602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556010909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146126aa576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612707576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b60006109ea6127a98361306c565b600084815260046020526040902054610100900463ffffffff1661166a565b6127d0612f35565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561285d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128819190614a67565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33601854846128ce919061463c565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401611ec1565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606083018190526001909101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c0820152919015612abc57816060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612a93573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ab79190614a80565b612abf565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836020015163ffffffff168152602001600760008781526020019081526020016000208054612b17906146e5565b80601f0160208091040260200160405190810160405280929190818152602001828054612b43906146e5565b8015612b905780601f10612b6557610100808354040283529160200191612b90565b820191906000526020600020905b815481529060010190602001808311612b7357829003601f168201915b505050505081526020018360a001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff1667ffffffffffffffff1681526020018360c0015163ffffffff16815260200183608001516bffffffffffffffffffffffff168152602001836000015115158152602001601b60008781526020019081526020016000208054612c6d906146e5565b80601f0160208091040260200160405190810160405280929190818152602001828054612c99906146e5565b8015612ce65780601f10612cbb57610100808354040283529160200191612ce6565b820191906000526020600020905b815481529060010190602001808311612cc957829003601f168201915b505050505081525092505050919050565b612d0083612fb8565b60155463ffffffff16811115612d42576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020612d5b828483614787565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d8383604051610a9d9291906148a2565b60006109ea8261279b565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612df9576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603612e48576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601060205260409020548116908216146110d85773ffffffffffffffffffffffffffffffffffffffff82811660008181526010602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b612eff612f35565b612f0881613615565b50565b60006109ea825490565b6000612f21838361370a565b9392505050565b60606000612f2183613734565b60005473ffffffffffffffffffffffffffffffffffffffff163314612fb6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611d24565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613015576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff90811614612f08576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f8110156130f9577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106130b1576130b161467e565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146130e757506000949350505050565b806130f1816146ad565b915050613073565b5081600f1a600181111561310f5761310f61436a565b949350505050565b6000612f21838361378f565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156131af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131d39190614ab7565b50945090925050506000811315806131ea57508142105b8061320b575082801561320b5750613202824261463c565b8463ffffffff16105b1561321a57601654955061321e565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613289573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132ad9190614ab7565b50945090925050506000811315806132c457508142105b806132e557508280156132e557506132dc824261463c565b8463ffffffff16105b156132f45760175494506132f8565b8094505b50505050915091565b60008061331388878b600001516137de565b905060008061332e8b8a63ffffffff16858a8a60018b6138a0565b909250905061333d8183614a20565b9b9a5050505050505050505050565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156133825761338261436a565b036133fc57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156133d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133f79190614a67565b905090565b504390565b6000612f218383613cf9565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906136095760008160600151856134a591906149a1565b905060006134b385836149f5565b905080836040018181516134c79190614a20565b6bffffffffffffffffffffffff169052506134e28582614b07565b836060018181516134f39190614a20565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613694576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611d24565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106137215761372161467e565b9060005260206000200154905092915050565b606081600001805480602002602001604051908101604052809291908181526020018280548015610fca57602002820191906000526020600020905b8154815260200190600101908083116137705750505050509050919050565b60008181526001830160205260408120546137d6575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556109ea565b5060006109ea565b600080808560018111156137f4576137f461436a565b03613803575062015f90613858565b60018560018111156138175761381761436a565b0361382657506201adb0613858565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61386963ffffffff85166014614b3b565b613874846001614b78565b6138839060ff16611d4c614b3b565b61388d9083614629565b6138979190614629565b95945050505050565b6000806000896080015161ffff16876138b99190614b3b565b90508380156138c75750803a105b156138cf57503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156139055761390561436a565b03613a6457604080516000815260208101909152851561396357600036604051806080016040528060488152602001614c696048913960405160200161394d93929190614b91565b60405160208183030381529060405290506139cb565b60155461397f90640100000000900463ffffffff166004614bb8565b63ffffffff1667ffffffffffffffff81111561399d5761399d61464f565b6040519080825280601f01601f1916602001820160405280156139c7576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90613a1b9084906004016142a0565b602060405180830381865afa158015613a38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a5c9190614a67565b915050613bbe565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115613a9857613a9861436a565b03613bbe578415613b1a57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613aef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b139190614a67565b9050613bbe565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa158015613b68573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b8c9190614bdb565b5050601554929450613baf93505050640100000000900463ffffffff1682614b3b565b613bba906010614b3b565b9150505b84613bda57808b6080015161ffff16613bd79190614b3b565b90505b613be861ffff871682614c25565b905060008782613bf88c8e614629565b613c029086614b3b565b613c0c9190614629565b613c1e90670de0b6b3a7640000614b3b565b613c289190614c25565b905060008c6040015163ffffffff1664e8d4a51000613c479190614b3b565b898e6020015163ffffffff16858f88613c609190614b3b565b613c6a9190614629565b613c7890633b9aca00614b3b565b613c829190614b3b565b613c8c9190614c25565b613c969190614629565b90506b033b2e3c9fd0803ce8000000613caf8284614629565b1115613ce7576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60008181526001830160205260408120548015613de2576000613d1d60018361463c565b8554909150600090613d319060019061463c565b9050818114613d96576000866000018281548110613d5157613d5161467e565b9060005260206000200154905080876000018481548110613d7457613d7461467e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613da757613da7614c39565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506109ea565b60009150506109ea565b5092915050565b60008060408385031215613e0657600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613e4d57835183529284019291840191600101613e31565b50909695505050505050565b60008083601f840112613e6b57600080fd5b50813567ffffffffffffffff811115613e8357600080fd5b602083019150836020828501011115613e9b57600080fd5b9250929050565b600080600060408486031215613eb757600080fd5b83359250602084013567ffffffffffffffff811115613ed557600080fd5b613ee186828701613e59565b9497909650939450505050565b600081518084526020808501945080840160005b83811015613f3457815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613f02565b509495945050505050565b805163ffffffff16825260006101e06020830151613f65602086018263ffffffff169052565b506040830151613f7d604086018263ffffffff169052565b506060830151613f94606086018262ffffff169052565b506080830151613faa608086018261ffff169052565b5060a0830151613fca60a08601826bffffffffffffffffffffffff169052565b5060c0830151613fe260c086018263ffffffff169052565b5060e0830151613ffa60e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a08084015181860183905261406f83870182613eee565b925050506101c08084015161409b8287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c060208801516140d360208501826bffffffffffffffffffffffff169052565b506040880151604084015260608801516140fd60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a088015161411f60a085018263ffffffff169052565b5060c088015161413760c085018263ffffffff169052565b5060e088015160e08401526101008089015161415a8286018263ffffffff169052565b505061012088810151151590840152610140830181905261417d81840188613f3f565b90508281036101608401526141928187613eee565b90508281036101808401526141a78186613eee565b9150506117926101a083018460ff169052565b73ffffffffffffffffffffffffffffffffffffffff81168114612f0857600080fd5b600080604083850312156141ef57600080fd5b82356141fa816141ba565b915060208301356004811061420e57600080fd5b809150509250929050565b60006020828403121561422b57600080fd5b5035919050565b60005b8381101561424d578181015183820152602001614235565b50506000910152565b6000815180845261426e816020860160208601614232565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612f216020830184614256565b600080604083850312156142c657600080fd5b82359150602083013561420e816141ba565b600080602083850312156142eb57600080fd5b823567ffffffffffffffff8082111561430357600080fd5b818501915085601f83011261431757600080fd5b81358181111561432657600080fd5b8660208260051b850101111561433b57600080fd5b60209290920196919550909350505050565b60006020828403121561435f57600080fd5b8135612f21816141ba565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60038110612f0857612f0861436a565b602081016143b683614399565b91905290565b60208101600283106143b6576143b661436a565b803563ffffffff811681146143e457600080fd5b919050565b600080604083850312156143fc57600080fd5b82356002811061440b57600080fd5b9150614419602084016143d0565b90509250929050565b60008060006040848603121561443757600080fd5b8335614442816141ba565b9250602084013567ffffffffffffffff811115613ed557600080fd5b6000806040838503121561447157600080fd5b823561447c816141ba565b9150602083013561420e816141ba565b6000806040838503121561449f57600080fd5b82359150614419602084016143d0565b602081526144d660208201835173ffffffffffffffffffffffffffffffffffffffff169052565b600060208301516144ef604084018263ffffffff169052565b50604083015161014080606085015261450c610160850183614256565b9150606085015161452d60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614599818701836bffffffffffffffffffffffff169052565b86015190506101206145ae8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506117928382614256565b60208101600483106143b6576143b661436a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156109ea576109ea6145fa565b818103818111156109ea576109ea6145fa565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036146de576146de6145fa565b5060010190565b600181811c908216806146f957607f821691505b602082108103614732577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561478257600081815260208120601f850160051c8101602086101561475f5750805b601f850160051c820191505b8181101561477e5782815560010161476b565b5050505b505050565b67ffffffffffffffff83111561479f5761479f61464f565b6147b3836147ad83546146e5565b83614738565b6000601f84116001811461480557600085156147cf5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561489b565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156148545786850135825560209485019460019092019101614834565b508682101561488f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000604082016040835280865480835260608501915087600052602092508260002060005b8281101561494657815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101614914565b505050838103828501528481528590820160005b8681101561499557823561496d816141ba565b73ffffffffffffffffffffffffffffffffffffffff168252918301919083019060010161495a565b50979650505050505050565b6bffffffffffffffffffffffff828116828216039080821115613dec57613dec6145fa565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680614a1457614a146149c6565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115613dec57613dec6145fa565b600060208284031215614a5757600080fd5b81518015158114612f2157600080fd5b600060208284031215614a7957600080fd5b5051919050565b600060208284031215614a9257600080fd5b8151612f21816141ba565b805169ffffffffffffffffffff811681146143e457600080fd5b600080600080600060a08688031215614acf57600080fd5b614ad886614a9d565b9450602086015193506040860151925060608601519150614afb60808701614a9d565b90509295509295909350565b60006bffffffffffffffffffffffff80831681851681830481118215151615614b3257614b326145fa565b02949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614b7357614b736145fa565b500290565b60ff81811683821601908111156109ea576109ea6145fa565b828482376000838201600081528351614bae818360208801614232565b0195945050505050565b600063ffffffff80831681851681830481118215151615614b3257614b326145fa565b60008060008060008060c08789031215614bf457600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b600082614c3457614c346149c6565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
}

var KeeperRegistryLogicB21ABI = KeeperRegistryLogicB21MetaData.ABI

var KeeperRegistryLogicB21Bin = KeeperRegistryLogicB21MetaData.Bin

func DeployKeeperRegistryLogicB21(auth *bind.TransactOpts, backend bind.ContractBackend, mode uint8, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address, automationForwarderLogic common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogicB21, error) {
	parsed, err := KeeperRegistryLogicB21MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicB21Bin), backend, mode, link, linkNativeFeed, fastGasFeed, automationForwarderLogic)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogicB21{address: address, abi: *parsed, KeeperRegistryLogicB21Caller: KeeperRegistryLogicB21Caller{contract: contract}, KeeperRegistryLogicB21Transactor: KeeperRegistryLogicB21Transactor{contract: contract}, KeeperRegistryLogicB21Filterer: KeeperRegistryLogicB21Filterer{contract: contract}}, nil
}

type KeeperRegistryLogicB21 struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogicB21Caller
	KeeperRegistryLogicB21Transactor
	KeeperRegistryLogicB21Filterer
}

type KeeperRegistryLogicB21Caller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicB21Transactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicB21Filterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicB21Session struct {
	Contract     *KeeperRegistryLogicB21
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicB21CallerSession struct {
	Contract *KeeperRegistryLogicB21Caller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogicB21TransactorSession struct {
	Contract     *KeeperRegistryLogicB21Transactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicB21Raw struct {
	Contract *KeeperRegistryLogicB21
}

type KeeperRegistryLogicB21CallerRaw struct {
	Contract *KeeperRegistryLogicB21Caller
}

type KeeperRegistryLogicB21TransactorRaw struct {
	Contract *KeeperRegistryLogicB21Transactor
}

func NewKeeperRegistryLogicB21(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogicB21, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogicB21ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogicB21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21{address: address, abi: abi, KeeperRegistryLogicB21Caller: KeeperRegistryLogicB21Caller{contract: contract}, KeeperRegistryLogicB21Transactor: KeeperRegistryLogicB21Transactor{contract: contract}, KeeperRegistryLogicB21Filterer: KeeperRegistryLogicB21Filterer{contract: contract}}, nil
}

func NewKeeperRegistryLogicB21Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogicB21Caller, error) {
	contract, err := bindKeeperRegistryLogicB21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21Caller{contract: contract}, nil
}

func NewKeeperRegistryLogicB21Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogicB21Transactor, error) {
	contract, err := bindKeeperRegistryLogicB21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21Transactor{contract: contract}, nil
}

func NewKeeperRegistryLogicB21Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogicB21Filterer, error) {
	contract, err := bindKeeperRegistryLogicB21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21Filterer{contract: contract}, nil
}

func bindKeeperRegistryLogicB21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryLogicB21MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicB21.Contract.KeeperRegistryLogicB21Caller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.KeeperRegistryLogicB21Transactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.KeeperRegistryLogicB21Transactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicB21.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetActiveUpkeepIDs(&_KeeperRegistryLogicB21.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetActiveUpkeepIDs(&_KeeperRegistryLogicB21.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetAdminPrivilegeConfig(&_KeeperRegistryLogicB21.CallOpts, admin)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetAdminPrivilegeConfig(&_KeeperRegistryLogicB21.CallOpts, admin)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetAutomationForwarderLogic() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetAutomationForwarderLogic(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetAutomationForwarderLogic(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetBalance(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetBalance(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetBalance(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetCancellationDelay() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetCancellationDelay(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetCancellationDelay() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetCancellationDelay(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetConditionalGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetConditionalGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetConditionalGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetForwarder(&_KeeperRegistryLogicB21.CallOpts, upkeepID)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetForwarder(&_KeeperRegistryLogicB21.CallOpts, upkeepID)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetLinkAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetLinkAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetLogGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetLogGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetLogGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB21.CallOpts, triggerType, gasLimit)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB21.CallOpts, triggerType, gasLimit)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMinBalance(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMinBalance(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetMode() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetMode(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetMode(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB21.CallOpts, peer)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB21.CallOpts, peer)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetPerPerformByteGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetPerPerformByteGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetPerSignerGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetPerSignerGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _KeeperRegistryLogicB21.Contract.GetPerSignerGasOverhead(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetSignerInfo(&_KeeperRegistryLogicB21.CallOpts, query)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetSignerInfo(&_KeeperRegistryLogicB21.CallOpts, query)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getState")

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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetState() (GetState,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetState(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetState() (GetState,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetState(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetTransmitterInfo(&_KeeperRegistryLogicB21.CallOpts, query)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistryLogicB21.Contract.GetTransmitterInfo(&_KeeperRegistryLogicB21.CallOpts, query)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetTriggerType(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.GetTriggerType(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeep(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeep(&_KeeperRegistryLogicB21.CallOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeepTriggerConfig(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB21.Contract.GetUpkeepTriggerConfig(&_KeeperRegistryLogicB21.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _KeeperRegistryLogicB21.Contract.HasDedupKey(&_KeeperRegistryLogicB21.CallOpts, dedupKey)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _KeeperRegistryLogicB21.Contract.HasDedupKey(&_KeeperRegistryLogicB21.CallOpts, dedupKey)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) Owner() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.Owner(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicB21.Contract.Owner(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.UpkeepTranscoderVersion(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.UpkeepTranscoderVersion(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Caller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB21.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) UpkeepVersion() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.UpkeepVersion(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21CallerSession) UpkeepVersion() (uint8, error) {
	return _KeeperRegistryLogicB21.Contract.UpkeepVersion(&_KeeperRegistryLogicB21.CallOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptOwnership(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptOwnership(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptPayeeship(&_KeeperRegistryLogicB21.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptPayeeship(&_KeeperRegistryLogicB21.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "pause")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.Pause(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.Pause(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.PauseUpkeep(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.PauseUpkeep(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.RecoverFunds(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.RecoverFunds(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetAdminPrivilegeConfig(&_KeeperRegistryLogicB21.TransactOpts, admin, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetAdminPrivilegeConfig(&_KeeperRegistryLogicB21.TransactOpts, admin, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setPayees", payees)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetPayees(&_KeeperRegistryLogicB21.TransactOpts, payees)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetPayees(&_KeeperRegistryLogicB21.TransactOpts, payees)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB21.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB21.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepCheckData(&_KeeperRegistryLogicB21.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepCheckData(&_KeeperRegistryLogicB21.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicB21.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicB21.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicB21.TransactOpts, id, config)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicB21.TransactOpts, id, config)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB21.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.SetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB21.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferOwnership(&_KeeperRegistryLogicB21.TransactOpts, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferOwnership(&_KeeperRegistryLogicB21.TransactOpts, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferPayeeship(&_KeeperRegistryLogicB21.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferPayeeship(&_KeeperRegistryLogicB21.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicB21.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicB21.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "unpause")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.Unpause(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.Unpause(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.UnpauseUpkeep(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.UnpauseUpkeep(&_KeeperRegistryLogicB21.TransactOpts, id)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawFunds(&_KeeperRegistryLogicB21.TransactOpts, id, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawFunds(&_KeeperRegistryLogicB21.TransactOpts, id, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogicB21.TransactOpts)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawPayment(&_KeeperRegistryLogicB21.TransactOpts, from, to)
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB21.Contract.WithdrawPayment(&_KeeperRegistryLogicB21.TransactOpts, from, to)
}

type KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicB21AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21AdminPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicB21AdminPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21AdminPrivilegeConfigSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicB21AdminPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicB21AdminPrivilegeConfigSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21CancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogicB21CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21CancelledUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicB21CancelledUpkeepReport)
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

func (it *KeeperRegistryLogicB21CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21CancelledUpkeepReportIterator{contract: _KeeperRegistryLogicB21.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21CancelledUpkeepReport)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicB21CancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogicB21CancelledUpkeepReport)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21DedupKeyAddedIterator struct {
	Event *KeeperRegistryLogicB21DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21DedupKeyAdded)
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
		it.Event = new(KeeperRegistryLogicB21DedupKeyAdded)
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

func (it *KeeperRegistryLogicB21DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicB21DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21DedupKeyAddedIterator{contract: _KeeperRegistryLogicB21.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21DedupKeyAdded)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicB21DedupKeyAdded, error) {
	event := new(KeeperRegistryLogicB21DedupKeyAdded)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21FundsAddedIterator struct {
	Event *KeeperRegistryLogicB21FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21FundsAdded)
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
		it.Event = new(KeeperRegistryLogicB21FundsAdded)
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

func (it *KeeperRegistryLogicB21FundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicB21FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21FundsAddedIterator{contract: _KeeperRegistryLogicB21.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21FundsAdded)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogicB21FundsAdded, error) {
	event := new(KeeperRegistryLogicB21FundsAdded)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21FundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicB21FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21FundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicB21FundsWithdrawn)
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

func (it *KeeperRegistryLogicB21FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21FundsWithdrawnIterator{contract: _KeeperRegistryLogicB21.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21FundsWithdrawn)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicB21FundsWithdrawn, error) {
	event := new(KeeperRegistryLogicB21FundsWithdrawn)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogicB21InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21InsufficientFundsUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicB21InsufficientFundsUpkeepReport)
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

func (it *KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogicB21.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21InsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicB21InsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogicB21InsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicB21OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21OwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicB21OwnerFundsWithdrawn)
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

func (it *KeeperRegistryLogicB21OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicB21OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21OwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogicB21.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21OwnerFundsWithdrawn)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicB21OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicB21OwnerFundsWithdrawn)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicB21OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21OwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicB21OwnershipTransferRequested)
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

func (it *KeeperRegistryLogicB21OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21OwnershipTransferRequestedIterator{contract: _KeeperRegistryLogicB21.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21OwnershipTransferRequested)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicB21OwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogicB21OwnershipTransferRequested)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21OwnershipTransferredIterator struct {
	Event *KeeperRegistryLogicB21OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21OwnershipTransferred)
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
		it.Event = new(KeeperRegistryLogicB21OwnershipTransferred)
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

func (it *KeeperRegistryLogicB21OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21OwnershipTransferredIterator{contract: _KeeperRegistryLogicB21.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21OwnershipTransferred)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicB21OwnershipTransferred, error) {
	event := new(KeeperRegistryLogicB21OwnershipTransferred)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21PausedIterator struct {
	Event *KeeperRegistryLogicB21Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21Paused)
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
		it.Event = new(KeeperRegistryLogicB21Paused)
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

func (it *KeeperRegistryLogicB21PausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicB21PausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21PausedIterator{contract: _KeeperRegistryLogicB21.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21Paused)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParsePaused(log types.Log) (*KeeperRegistryLogicB21Paused, error) {
	event := new(KeeperRegistryLogicB21Paused)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21PayeesUpdatedIterator struct {
	Event *KeeperRegistryLogicB21PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21PayeesUpdated)
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
		it.Event = new(KeeperRegistryLogicB21PayeesUpdated)
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

func (it *KeeperRegistryLogicB21PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicB21PayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21PayeesUpdatedIterator{contract: _KeeperRegistryLogicB21.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21PayeesUpdated)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicB21PayeesUpdated, error) {
	event := new(KeeperRegistryLogicB21PayeesUpdated)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicB21PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21PayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicB21PayeeshipTransferRequested)
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

func (it *KeeperRegistryLogicB21PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21PayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogicB21.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21PayeeshipTransferRequested)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicB21PayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogicB21PayeeshipTransferRequested)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21PayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogicB21PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21PayeeshipTransferred)
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
		it.Event = new(KeeperRegistryLogicB21PayeeshipTransferred)
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

func (it *KeeperRegistryLogicB21PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21PayeeshipTransferredIterator{contract: _KeeperRegistryLogicB21.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21PayeeshipTransferred)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicB21PayeeshipTransferred, error) {
	event := new(KeeperRegistryLogicB21PayeeshipTransferred)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21PaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogicB21PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21PaymentWithdrawn)
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
		it.Event = new(KeeperRegistryLogicB21PaymentWithdrawn)
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

func (it *KeeperRegistryLogicB21PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicB21PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21PaymentWithdrawnIterator{contract: _KeeperRegistryLogicB21.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21PaymentWithdrawn)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicB21PaymentWithdrawn, error) {
	event := new(KeeperRegistryLogicB21PaymentWithdrawn)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21ReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogicB21ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21ReorgedUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicB21ReorgedUpkeepReport)
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

func (it *KeeperRegistryLogicB21ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21ReorgedUpkeepReportIterator{contract: _KeeperRegistryLogicB21.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21ReorgedUpkeepReport)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicB21ReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogicB21ReorgedUpkeepReport)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21StaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogicB21StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21StaleUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicB21StaleUpkeepReport)
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

func (it *KeeperRegistryLogicB21StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21StaleUpkeepReportIterator{contract: _KeeperRegistryLogicB21.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21StaleUpkeepReport)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicB21StaleUpkeepReport, error) {
	event := new(KeeperRegistryLogicB21StaleUpkeepReport)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UnpausedIterator struct {
	Event *KeeperRegistryLogicB21Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21Unpaused)
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
		it.Event = new(KeeperRegistryLogicB21Unpaused)
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

func (it *KeeperRegistryLogicB21UnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicB21UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UnpausedIterator{contract: _KeeperRegistryLogicB21.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21Unpaused)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogicB21Unpaused, error) {
	event := new(KeeperRegistryLogicB21Unpaused)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepAdminTransferRequested)
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

func (it *KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicB21UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogicB21UpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogicB21UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepAdminTransferred)
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

func (it *KeeperRegistryLogicB21UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepAdminTransferredIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepAdminTransferred)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicB21UpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogicB21UpkeepAdminTransferred)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepCanceledIterator struct {
	Event *KeeperRegistryLogicB21UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepCanceled)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepCanceled)
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

func (it *KeeperRegistryLogicB21UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicB21UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepCanceledIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepCanceled)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicB21UpkeepCanceled, error) {
	event := new(KeeperRegistryLogicB21UpkeepCanceled)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepCheckDataSetIterator struct {
	Event *KeeperRegistryLogicB21UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepCheckDataSet)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepCheckDataSet)
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

func (it *KeeperRegistryLogicB21UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepCheckDataSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepCheckDataSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicB21UpkeepCheckDataSet, error) {
	event := new(KeeperRegistryLogicB21UpkeepCheckDataSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogicB21UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepGasLimitSet)
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

func (it *KeeperRegistryLogicB21UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepGasLimitSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepGasLimitSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicB21UpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogicB21UpkeepGasLimitSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepMigratedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepMigrated)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepMigrated)
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

func (it *KeeperRegistryLogicB21UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepMigratedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepMigrated)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicB21UpkeepMigrated, error) {
	event := new(KeeperRegistryLogicB21UpkeepMigrated)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicB21UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepOffchainConfigSet)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepOffchainConfigSet)
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

func (it *KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicB21UpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepPausedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepPaused)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepPaused)
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

func (it *KeeperRegistryLogicB21UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepPausedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepPaused)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicB21UpkeepPaused, error) {
	event := new(KeeperRegistryLogicB21UpkeepPaused)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepPerformedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepPerformed)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepPerformed)
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

func (it *KeeperRegistryLogicB21UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicB21UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepPerformedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepPerformed)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicB21UpkeepPerformed, error) {
	event := new(KeeperRegistryLogicB21UpkeepPerformed)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicB21UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepPrivilegeConfigSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicB21UpkeepPrivilegeConfigSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepReceivedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepReceived)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepReceived)
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

func (it *KeeperRegistryLogicB21UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepReceivedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepReceived)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicB21UpkeepReceived, error) {
	event := new(KeeperRegistryLogicB21UpkeepReceived)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogicB21UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepRegistered)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepRegistered)
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

func (it *KeeperRegistryLogicB21UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepRegisteredIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepRegistered)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicB21UpkeepRegistered, error) {
	event := new(KeeperRegistryLogicB21UpkeepRegistered)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator struct {
	Event *KeeperRegistryLogicB21UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepTriggerConfigSet)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepTriggerConfigSet)
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

func (it *KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepTriggerConfigSet)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepTriggerConfigSet, error) {
	event := new(KeeperRegistryLogicB21UpkeepTriggerConfigSet)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicB21UpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogicB21UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicB21UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicB21UpkeepUnpaused)
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
		it.Event = new(KeeperRegistryLogicB21UpkeepUnpaused)
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

func (it *KeeperRegistryLogicB21UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicB21UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicB21UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB21UpkeepUnpausedIterator{contract: _KeeperRegistryLogicB21.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB21.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicB21UpkeepUnpaused)
				if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicB21UpkeepUnpaused, error) {
	event := new(KeeperRegistryLogicB21UpkeepUnpaused)
	if err := _KeeperRegistryLogicB21.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogicB21.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicB21.ParseAdminPrivilegeConfigSet(log)
	case _KeeperRegistryLogicB21.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogicB21.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogicB21.abi.Events["DedupKeyAdded"].ID:
		return _KeeperRegistryLogicB21.ParseDedupKeyAdded(log)
	case _KeeperRegistryLogicB21.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogicB21.ParseFundsAdded(log)
	case _KeeperRegistryLogicB21.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogicB21.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogicB21.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogicB21.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogicB21.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogicB21.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogicB21.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogicB21.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogicB21.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogicB21.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogicB21.abi.Events["Paused"].ID:
		return _KeeperRegistryLogicB21.ParsePaused(log)
	case _KeeperRegistryLogicB21.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogicB21.ParsePayeesUpdated(log)
	case _KeeperRegistryLogicB21.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogicB21.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogicB21.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogicB21.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogicB21.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogicB21.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogicB21.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogicB21.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogicB21.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogicB21.ParseStaleUpkeepReport(log)
	case _KeeperRegistryLogicB21.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogicB21.ParseUnpaused(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepCheckDataSet"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepCheckDataSet(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepPaused(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepPrivilegeConfigSet(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepReceived(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepTriggerConfigSet(log)
	case _KeeperRegistryLogicB21.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogicB21.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicB21AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (KeeperRegistryLogicB21CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (KeeperRegistryLogicB21DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (KeeperRegistryLogicB21FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicB21FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicB21InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (KeeperRegistryLogicB21OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogicB21OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogicB21OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogicB21Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogicB21PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryLogicB21PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogicB21PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogicB21PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogicB21ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (KeeperRegistryLogicB21StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (KeeperRegistryLogicB21Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogicB21UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogicB21UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogicB21UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogicB21UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (KeeperRegistryLogicB21UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogicB21UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogicB21UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogicB21UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicB21UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryLogicB21UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (KeeperRegistryLogicB21UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicB21UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogicB21UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (KeeperRegistryLogicB21UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogicB21 *KeeperRegistryLogicB21) Address() common.Address {
	return _KeeperRegistryLogicB21.address
}

type KeeperRegistryLogicB21Interface interface {
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error)

	GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error)

	GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error)

	GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMode(opts *bind.CallOpts) (uint8, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

		error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicB21AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicB21AdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicB21CancelledUpkeepReport, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicB21DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicB21DedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicB21FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicB21FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicB21FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicB21InsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicB21OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicB21OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicB21OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicB21OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicB21PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogicB21Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicB21PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicB21PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicB21PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicB21PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicB21PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicB21PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicB21ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicB21StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicB21UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogicB21Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicB21UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicB21UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicB21UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicB21UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicB21UpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicB21UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicB21UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicB21UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicB21UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicB21UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicB21UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicB21UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicB21UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicB21UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicB21UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicB21UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicB21UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
