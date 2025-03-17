// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_keeper_registry_master_wrapper_2_1

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

var IKeeperRegistryMasterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"fallback\",\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"values\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fallbackTo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveUpkeepIDs\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAutomationForwarderLogic\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCancellationDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getConditionalGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getFastGasFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForwarder\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkNativeFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLogGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getMaxPaymentForGas\",\"inputs\":[{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"maxPayment\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalanceForUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"minBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMode\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPerPerformByteGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPerSignerGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getSignerInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ownerLinkBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"expectedLinkBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalPremium\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"numUpkeeps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"latestEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitterInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"lastCollected\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTriggerType\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepInfo\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"amountSpent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasDedupKey\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[],\"outputs\":[{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDigestAndEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"scanLogs\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateUpkeeps\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTokenTransfer\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"receiveUpkeeps\",\"inputs\":[{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"onchainConfigBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfigTypeSafe\",\"inputs\":[{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepOffchainConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"simulatePerformUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transmit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[3]\",\"internalType\":\"bytes32[3]\"},{\"name\":\"rawReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upkeepTranscoderVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"upkeepVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawOwnerFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"configCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
}

var IKeeperRegistryMasterABI = IKeeperRegistryMasterMetaData.ABI

type IKeeperRegistryMaster struct {
	address common.Address
	abi     abi.ABI
	IKeeperRegistryMasterCaller
	IKeeperRegistryMasterTransactor
	IKeeperRegistryMasterFilterer
}

type IKeeperRegistryMasterCaller struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterTransactor struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterFilterer struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterSession struct {
	Contract     *IKeeperRegistryMaster
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IKeeperRegistryMasterCallerSession struct {
	Contract *IKeeperRegistryMasterCaller
	CallOpts bind.CallOpts
}

type IKeeperRegistryMasterTransactorSession struct {
	Contract     *IKeeperRegistryMasterTransactor
	TransactOpts bind.TransactOpts
}

type IKeeperRegistryMasterRaw struct {
	Contract *IKeeperRegistryMaster
}

type IKeeperRegistryMasterCallerRaw struct {
	Contract *IKeeperRegistryMasterCaller
}

type IKeeperRegistryMasterTransactorRaw struct {
	Contract *IKeeperRegistryMasterTransactor
}

func NewIKeeperRegistryMaster(address common.Address, backend bind.ContractBackend) (*IKeeperRegistryMaster, error) {
	abi, err := abi.JSON(strings.NewReader(IKeeperRegistryMasterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIKeeperRegistryMaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMaster{address: address, abi: abi, IKeeperRegistryMasterCaller: IKeeperRegistryMasterCaller{contract: contract}, IKeeperRegistryMasterTransactor: IKeeperRegistryMasterTransactor{contract: contract}, IKeeperRegistryMasterFilterer: IKeeperRegistryMasterFilterer{contract: contract}}, nil
}

func NewIKeeperRegistryMasterCaller(address common.Address, caller bind.ContractCaller) (*IKeeperRegistryMasterCaller, error) {
	contract, err := bindIKeeperRegistryMaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterCaller{contract: contract}, nil
}

func NewIKeeperRegistryMasterTransactor(address common.Address, transactor bind.ContractTransactor) (*IKeeperRegistryMasterTransactor, error) {
	contract, err := bindIKeeperRegistryMaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterTransactor{contract: contract}, nil
}

func NewIKeeperRegistryMasterFilterer(address common.Address, filterer bind.ContractFilterer) (*IKeeperRegistryMasterFilterer, error) {
	contract, err := bindIKeeperRegistryMaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFilterer{contract: contract}, nil
}

func bindIKeeperRegistryMaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IKeeperRegistryMasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterCaller.contract.Call(opts, result, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterTransactor.contract.Transfer(opts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterTransactor.contract.Transact(opts, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKeeperRegistryMaster.Contract.contract.Call(opts, result, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.contract.Transfer(opts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.contract.Transact(opts, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "checkCallback", id, values, extraData)

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckCallback(&_IKeeperRegistryMaster.CallOpts, id, values, extraData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckCallback(&_IKeeperRegistryMaster.CallOpts, id, values, extraData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "checkUpkeep", id, triggerData)

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep(&_IKeeperRegistryMaster.CallOpts, id, triggerData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep(&_IKeeperRegistryMaster.CallOpts, id, triggerData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "checkUpkeep0", id)

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep0(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep0(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) FallbackTo() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.FallbackTo(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) FallbackTo() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.FallbackTo(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDs(&_IKeeperRegistryMaster.CallOpts, startIndex, maxCount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDs(&_IKeeperRegistryMaster.CallOpts, startIndex, maxCount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetAdminPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, admin)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetAdminPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, admin)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetAutomationForwarderLogic(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetAutomationForwarderLogic(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetBalance(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetBalance(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetCancellationDelay() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetCancellationDelay(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetCancellationDelay() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetCancellationDelay(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetConditionalGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetConditionalGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetFastGasFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetFastGasFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetForwarder(&_IKeeperRegistryMaster.CallOpts, upkeepID)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetForwarder(&_IKeeperRegistryMaster.CallOpts, upkeepID)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLinkAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLinkAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLogGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetLogGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetLogGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMaxPaymentForGas(&_IKeeperRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMaxPaymentForGas(&_IKeeperRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalance(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalance(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMode() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetMode(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMode() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetMode(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.CallOpts, peer)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.CallOpts, peer)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetPerPerformByteGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetPerPerformByteGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetPerSignerGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetPerSignerGasOverhead(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetSignerInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetSignerInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getState")

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetState() (GetState,

	error) {
	return _IKeeperRegistryMaster.Contract.GetState(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetState() (GetState,

	error) {
	return _IKeeperRegistryMaster.Contract.GetState(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetTransmitterInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetTransmitterInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetTriggerType(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetTriggerType(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IKeeperRegistryMaster.Contract.HasDedupKey(&_IKeeperRegistryMaster.CallOpts, dedupKey)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IKeeperRegistryMaster.Contract.HasDedupKey(&_IKeeperRegistryMaster.CallOpts, dedupKey)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDetails(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDetails(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Owner() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.Owner(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) Owner() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.Owner(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "simulatePerformUpkeep", id, performData)

	outstruct := new(SimulatePerformUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GasUsed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IKeeperRegistryMaster.Contract.SimulatePerformUpkeep(&_IKeeperRegistryMaster.CallOpts, id, performData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IKeeperRegistryMaster.Contract.SimulatePerformUpkeep(&_IKeeperRegistryMaster.CallOpts, id, performData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TypeAndVersion() (string, error) {
	return _IKeeperRegistryMaster.Contract.TypeAndVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) TypeAndVersion() (string, error) {
	return _IKeeperRegistryMaster.Contract.TypeAndVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepTranscoderVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepTranscoderVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UpkeepVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) UpkeepVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptOwnership")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptOwnership() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptOwnership(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptOwnership(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "addFunds", id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AddFunds(&_IKeeperRegistryMaster.TransactOpts, id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AddFunds(&_IKeeperRegistryMaster.TransactOpts, id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CancelUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CancelUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ExecuteCallback(&_IKeeperRegistryMaster.TransactOpts, id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ExecuteCallback(&_IKeeperRegistryMaster.TransactOpts, id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.MigrateUpkeeps(&_IKeeperRegistryMaster.TransactOpts, ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.MigrateUpkeeps(&_IKeeperRegistryMaster.TransactOpts, ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.OnTokenTransfer(&_IKeeperRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.OnTokenTransfer(&_IKeeperRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "pause")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Pause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Pause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Pause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Pause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.PauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.PauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ReceiveUpkeeps(&_IKeeperRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ReceiveUpkeeps(&_IKeeperRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "recoverFunds")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RecoverFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RecoverFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RecoverFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep0(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep0(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetAdminPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, admin, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetAdminPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, admin, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfig(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfig(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig IAutomationV21PlusCommonOnchainConfigLegacy, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig IAutomationV21PlusCommonOnchainConfigLegacy, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfigTypeSafe(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig IAutomationV21PlusCommonOnchainConfigLegacy, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfigTypeSafe(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setPayees", payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPayees(&_IKeeperRegistryMaster.TransactOpts, payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPayees(&_IKeeperRegistryMaster.TransactOpts, payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.TransactOpts, peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.TransactOpts, peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepCheckData(&_IKeeperRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepCheckData(&_IKeeperRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepGasLimit(&_IKeeperRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepGasLimit(&_IKeeperRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IKeeperRegistryMaster.TransactOpts, id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IKeeperRegistryMaster.TransactOpts, id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IKeeperRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IKeeperRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferOwnership", to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferOwnership(&_IKeeperRegistryMaster.TransactOpts, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferOwnership(&_IKeeperRegistryMaster.TransactOpts, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Transmit(&_IKeeperRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Transmit(&_IKeeperRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "unpause")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Unpause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Unpause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Unpause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Unpause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.UnpauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.UnpauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawFunds(&_IKeeperRegistryMaster.TransactOpts, id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawFunds(&_IKeeperRegistryMaster.TransactOpts, id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawOwnerFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawOwnerFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawPayment(&_IKeeperRegistryMaster.TransactOpts, from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawPayment(&_IKeeperRegistryMaster.TransactOpts, from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.RawTransact(opts, calldata)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Fallback(&_IKeeperRegistryMaster.TransactOpts, calldata)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Fallback(&_IKeeperRegistryMaster.TransactOpts, calldata)
}

type IKeeperRegistryMasterAdminPrivilegeConfigSetIterator struct {
	Event *IKeeperRegistryMasterAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterAdminPrivilegeConfigSet)
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
		it.Event = new(IKeeperRegistryMasterAdminPrivilegeConfigSet)
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

func (it *IKeeperRegistryMasterAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IKeeperRegistryMasterAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterAdminPrivilegeConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterAdminPrivilegeConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterAdminPrivilegeConfigSet, error) {
	event := new(IKeeperRegistryMasterAdminPrivilegeConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterCancelledUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterCancelledUpkeepReport)
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
		it.Event = new(IKeeperRegistryMasterCancelledUpkeepReport)
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

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterCancelledUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterCancelledUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseCancelledUpkeepReport(log types.Log) (*IKeeperRegistryMasterCancelledUpkeepReport, error) {
	event := new(IKeeperRegistryMasterCancelledUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterConfigSetIterator struct {
	Event *IKeeperRegistryMasterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterConfigSet)
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
		it.Event = new(IKeeperRegistryMasterConfigSet)
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

func (it *IKeeperRegistryMasterConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterConfigSet struct {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterConfigSet(opts *bind.FilterOpts) (*IKeeperRegistryMasterConfigSetIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterConfigSet) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseConfigSet(log types.Log) (*IKeeperRegistryMasterConfigSet, error) {
	event := new(IKeeperRegistryMasterConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterDedupKeyAddedIterator struct {
	Event *IKeeperRegistryMasterDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterDedupKeyAdded)
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
		it.Event = new(IKeeperRegistryMasterDedupKeyAdded)
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

func (it *IKeeperRegistryMasterDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IKeeperRegistryMasterDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterDedupKeyAddedIterator{contract: _IKeeperRegistryMaster.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterDedupKeyAdded)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseDedupKeyAdded(log types.Log) (*IKeeperRegistryMasterDedupKeyAdded, error) {
	event := new(IKeeperRegistryMasterDedupKeyAdded)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterFundsAddedIterator struct {
	Event *IKeeperRegistryMasterFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterFundsAdded)
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
		it.Event = new(IKeeperRegistryMasterFundsAdded)
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

func (it *IKeeperRegistryMasterFundsAddedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IKeeperRegistryMasterFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFundsAddedIterator{contract: _IKeeperRegistryMaster.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterFundsAdded)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseFundsAdded(log types.Log) (*IKeeperRegistryMasterFundsAdded, error) {
	event := new(IKeeperRegistryMasterFundsAdded)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterFundsWithdrawnIterator struct {
	Event *IKeeperRegistryMasterFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterFundsWithdrawn)
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
		it.Event = new(IKeeperRegistryMasterFundsWithdrawn)
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

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFundsWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterFundsWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterFundsWithdrawn, error) {
	event := new(IKeeperRegistryMasterFundsWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
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
		it.Event = new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
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

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IKeeperRegistryMasterInsufficientFundsUpkeepReport, error) {
	event := new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnerFundsWithdrawnIterator struct {
	Event *IKeeperRegistryMasterOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnerFundsWithdrawn)
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
		it.Event = new(IKeeperRegistryMasterOwnerFundsWithdrawn)
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

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IKeeperRegistryMasterOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnerFundsWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnerFundsWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterOwnerFundsWithdrawn, error) {
	event := new(IKeeperRegistryMasterOwnerFundsWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnershipTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnershipTransferRequested)
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
		it.Event = new(IKeeperRegistryMasterOwnershipTransferRequested)
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

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnershipTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnershipTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnershipTransferRequested(log types.Log) (*IKeeperRegistryMasterOwnershipTransferRequested, error) {
	event := new(IKeeperRegistryMasterOwnershipTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnershipTransferredIterator struct {
	Event *IKeeperRegistryMasterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnershipTransferred)
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
		it.Event = new(IKeeperRegistryMasterOwnershipTransferred)
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

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnershipTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnershipTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnershipTransferred(log types.Log) (*IKeeperRegistryMasterOwnershipTransferred, error) {
	event := new(IKeeperRegistryMasterOwnershipTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPausedIterator struct {
	Event *IKeeperRegistryMasterPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPaused)
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
		it.Event = new(IKeeperRegistryMasterPaused)
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

func (it *IKeeperRegistryMasterPausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterPausedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPausedIterator{contract: _IKeeperRegistryMaster.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaused) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePaused(log types.Log) (*IKeeperRegistryMasterPaused, error) {
	event := new(IKeeperRegistryMasterPaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeesUpdatedIterator struct {
	Event *IKeeperRegistryMasterPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeesUpdated)
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
		it.Event = new(IKeeperRegistryMasterPayeesUpdated)
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

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IKeeperRegistryMasterPayeesUpdatedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeesUpdatedIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeesUpdated)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeesUpdated(log types.Log) (*IKeeperRegistryMasterPayeesUpdated, error) {
	event := new(IKeeperRegistryMasterPayeesUpdated)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeeshipTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeeshipTransferRequested)
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
		it.Event = new(IKeeperRegistryMasterPayeeshipTransferRequested)
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

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeeshipTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeeshipTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeeshipTransferRequested(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferRequested, error) {
	event := new(IKeeperRegistryMasterPayeeshipTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeeshipTransferredIterator struct {
	Event *IKeeperRegistryMasterPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeeshipTransferred)
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
		it.Event = new(IKeeperRegistryMasterPayeeshipTransferred)
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

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeeshipTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeeshipTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeeshipTransferred(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferred, error) {
	event := new(IKeeperRegistryMasterPayeeshipTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPaymentWithdrawnIterator struct {
	Event *IKeeperRegistryMasterPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPaymentWithdrawn)
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
		it.Event = new(IKeeperRegistryMasterPaymentWithdrawn)
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

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IKeeperRegistryMasterPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPaymentWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPaymentWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePaymentWithdrawn(log types.Log) (*IKeeperRegistryMasterPaymentWithdrawn, error) {
	event := new(IKeeperRegistryMasterPaymentWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterReorgedUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterReorgedUpkeepReport)
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
		it.Event = new(IKeeperRegistryMasterReorgedUpkeepReport)
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

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterReorgedUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterReorgedUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseReorgedUpkeepReport(log types.Log) (*IKeeperRegistryMasterReorgedUpkeepReport, error) {
	event := new(IKeeperRegistryMasterReorgedUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterStaleUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterStaleUpkeepReport)
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
		it.Event = new(IKeeperRegistryMasterStaleUpkeepReport)
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

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterStaleUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterStaleUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseStaleUpkeepReport(log types.Log) (*IKeeperRegistryMasterStaleUpkeepReport, error) {
	event := new(IKeeperRegistryMasterStaleUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterTransmittedIterator struct {
	Event *IKeeperRegistryMasterTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterTransmitted)
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
		it.Event = new(IKeeperRegistryMasterTransmitted)
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

func (it *IKeeperRegistryMasterTransmittedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterTransmitted(opts *bind.FilterOpts) (*IKeeperRegistryMasterTransmittedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterTransmittedIterator{contract: _IKeeperRegistryMaster.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterTransmitted) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterTransmitted)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseTransmitted(log types.Log) (*IKeeperRegistryMasterTransmitted, error) {
	event := new(IKeeperRegistryMasterTransmitted)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUnpausedIterator struct {
	Event *IKeeperRegistryMasterUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUnpaused)
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
		it.Event = new(IKeeperRegistryMasterUnpaused)
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

func (it *IKeeperRegistryMasterUnpausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterUnpausedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUnpausedIterator{contract: _IKeeperRegistryMaster.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUnpaused) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUnpaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUnpaused(log types.Log) (*IKeeperRegistryMasterUnpaused, error) {
	event := new(IKeeperRegistryMasterUnpaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
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
		it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
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

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferRequested, error) {
	event := new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepAdminTransferredIterator struct {
	Event *IKeeperRegistryMasterUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferred)
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
		it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferred)
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

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepAdminTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepAdminTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepAdminTransferred(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferred, error) {
	event := new(IKeeperRegistryMasterUpkeepAdminTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepCanceledIterator struct {
	Event *IKeeperRegistryMasterUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepCanceled)
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
		it.Event = new(IKeeperRegistryMasterUpkeepCanceled)
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

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IKeeperRegistryMasterUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepCanceledIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepCanceled)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepCanceled(log types.Log) (*IKeeperRegistryMasterUpkeepCanceled, error) {
	event := new(IKeeperRegistryMasterUpkeepCanceled)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepCheckDataSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepCheckDataSet)
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
		it.Event = new(IKeeperRegistryMasterUpkeepCheckDataSet)
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

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepCheckDataSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepCheckDataSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepCheckDataSet(log types.Log) (*IKeeperRegistryMasterUpkeepCheckDataSet, error) {
	event := new(IKeeperRegistryMasterUpkeepCheckDataSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepGasLimitSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepGasLimitSet)
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
		it.Event = new(IKeeperRegistryMasterUpkeepGasLimitSet)
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

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepGasLimitSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepGasLimitSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepGasLimitSet(log types.Log) (*IKeeperRegistryMasterUpkeepGasLimitSet, error) {
	event := new(IKeeperRegistryMasterUpkeepGasLimitSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepMigratedIterator struct {
	Event *IKeeperRegistryMasterUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepMigrated)
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
		it.Event = new(IKeeperRegistryMasterUpkeepMigrated)
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

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepMigratedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepMigrated)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepMigrated(log types.Log) (*IKeeperRegistryMasterUpkeepMigrated, error) {
	event := new(IKeeperRegistryMasterUpkeepMigrated)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepOffchainConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
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
		it.Event = new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
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

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepOffchainConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepOffchainConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPausedIterator struct {
	Event *IKeeperRegistryMasterUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPaused)
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
		it.Event = new(IKeeperRegistryMasterUpkeepPaused)
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

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPausedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPaused(log types.Log) (*IKeeperRegistryMasterUpkeepPaused, error) {
	event := new(IKeeperRegistryMasterUpkeepPaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPerformedIterator struct {
	Event *IKeeperRegistryMasterUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPerformed)
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
		it.Event = new(IKeeperRegistryMasterUpkeepPerformed)
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

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IKeeperRegistryMasterUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPerformedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPerformed)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPerformed(log types.Log) (*IKeeperRegistryMasterUpkeepPerformed, error) {
	event := new(IKeeperRegistryMasterUpkeepPerformed)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
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
		it.Event = new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
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

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepReceivedIterator struct {
	Event *IKeeperRegistryMasterUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepReceived)
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
		it.Event = new(IKeeperRegistryMasterUpkeepReceived)
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

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepReceivedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepReceived)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepReceived(log types.Log) (*IKeeperRegistryMasterUpkeepReceived, error) {
	event := new(IKeeperRegistryMasterUpkeepReceived)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepRegisteredIterator struct {
	Event *IKeeperRegistryMasterUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepRegistered)
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
		it.Event = new(IKeeperRegistryMasterUpkeepRegistered)
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

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepRegisteredIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepRegistered)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepRegistered(log types.Log) (*IKeeperRegistryMasterUpkeepRegistered, error) {
	event := new(IKeeperRegistryMasterUpkeepRegistered)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepTriggerConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
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
		it.Event = new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
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

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepTriggerConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepTriggerConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepUnpausedIterator struct {
	Event *IKeeperRegistryMasterUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepUnpaused)
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
		it.Event = new(IKeeperRegistryMasterUpkeepUnpaused)
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

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepUnpausedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepUnpaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepUnpaused(log types.Log) (*IKeeperRegistryMasterUpkeepUnpaused, error) {
	event := new(IKeeperRegistryMasterUpkeepUnpaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}
type SimulatePerformUpkeep struct {
	Success bool
	GasUsed *big.Int
}

func (_IKeeperRegistryMaster *IKeeperRegistryMaster) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IKeeperRegistryMaster.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseAdminPrivilegeConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["CancelledUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseCancelledUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["ConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["DedupKeyAdded"].ID:
		return _IKeeperRegistryMaster.ParseDedupKeyAdded(log)
	case _IKeeperRegistryMaster.abi.Events["FundsAdded"].ID:
		return _IKeeperRegistryMaster.ParseFundsAdded(log)
	case _IKeeperRegistryMaster.abi.Events["FundsWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParseFundsWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseInsufficientFundsUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["OwnerFundsWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParseOwnerFundsWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["OwnershipTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParseOwnershipTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["OwnershipTransferred"].ID:
		return _IKeeperRegistryMaster.ParseOwnershipTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["Paused"].ID:
		return _IKeeperRegistryMaster.ParsePaused(log)
	case _IKeeperRegistryMaster.abi.Events["PayeesUpdated"].ID:
		return _IKeeperRegistryMaster.ParsePayeesUpdated(log)
	case _IKeeperRegistryMaster.abi.Events["PayeeshipTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParsePayeeshipTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["PayeeshipTransferred"].ID:
		return _IKeeperRegistryMaster.ParsePayeeshipTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["PaymentWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParsePaymentWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["ReorgedUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseReorgedUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["StaleUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseStaleUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["Transmitted"].ID:
		return _IKeeperRegistryMaster.ParseTransmitted(log)
	case _IKeeperRegistryMaster.abi.Events["Unpaused"].ID:
		return _IKeeperRegistryMaster.ParseUnpaused(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepAdminTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepAdminTransferred"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepAdminTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepCanceled"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepCanceled(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepCheckDataSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepCheckDataSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepGasLimitSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepGasLimitSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepMigrated"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepMigrated(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepOffchainConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPaused"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPaused(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPerformed"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPerformed(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPrivilegeConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepReceived"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepReceived(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepRegistered"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepRegistered(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepTriggerConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepUnpaused"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IKeeperRegistryMasterAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (IKeeperRegistryMasterCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IKeeperRegistryMasterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IKeeperRegistryMasterDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IKeeperRegistryMasterFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IKeeperRegistryMasterFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IKeeperRegistryMasterInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IKeeperRegistryMasterOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (IKeeperRegistryMasterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IKeeperRegistryMasterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IKeeperRegistryMasterPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IKeeperRegistryMasterPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IKeeperRegistryMasterPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IKeeperRegistryMasterPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IKeeperRegistryMasterPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IKeeperRegistryMasterReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IKeeperRegistryMasterStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IKeeperRegistryMasterTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IKeeperRegistryMasterUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IKeeperRegistryMasterUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IKeeperRegistryMasterUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IKeeperRegistryMasterUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IKeeperRegistryMasterUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IKeeperRegistryMasterUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IKeeperRegistryMasterUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IKeeperRegistryMasterUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IKeeperRegistryMasterUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IKeeperRegistryMasterUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IKeeperRegistryMasterUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (IKeeperRegistryMasterUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IKeeperRegistryMasterUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IKeeperRegistryMasterUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IKeeperRegistryMasterUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMaster) Address() common.Address {
	return _IKeeperRegistryMaster.address
}

type IKeeperRegistryMasterInterface interface {
	CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

		error)

	CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

		error)

	CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

		error)

	FallbackTo(opts *bind.CallOpts) (common.Address, error)

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

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

		error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig IAutomationV21PlusCommonOnchainConfigLegacy, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IKeeperRegistryMasterAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IKeeperRegistryMasterCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IKeeperRegistryMasterConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IKeeperRegistryMasterConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IKeeperRegistryMasterDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IKeeperRegistryMasterDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IKeeperRegistryMasterFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IKeeperRegistryMasterFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IKeeperRegistryMasterInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IKeeperRegistryMasterOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IKeeperRegistryMasterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IKeeperRegistryMasterOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IKeeperRegistryMasterPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IKeeperRegistryMasterPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IKeeperRegistryMasterPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IKeeperRegistryMasterPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IKeeperRegistryMasterPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IKeeperRegistryMasterReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IKeeperRegistryMasterStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IKeeperRegistryMasterTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IKeeperRegistryMasterTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IKeeperRegistryMasterUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IKeeperRegistryMasterUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IKeeperRegistryMasterUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IKeeperRegistryMasterUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IKeeperRegistryMasterUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IKeeperRegistryMasterUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IKeeperRegistryMasterUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IKeeperRegistryMasterUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IKeeperRegistryMasterUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IKeeperRegistryMasterUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IKeeperRegistryMasterUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IKeeperRegistryMasterUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
