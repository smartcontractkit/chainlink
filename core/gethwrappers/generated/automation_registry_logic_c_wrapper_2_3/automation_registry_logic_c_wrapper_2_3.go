// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_registry_logic_c_wrapper_2_3

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

type AutomationRegistryBase23BillingConfig struct {
	GasFeePPB         uint32
	FlatFeeMilliCents *big.Int
	PriceFeed         common.Address
	Decimals          uint8
	FallbackPrice     *big.Int
	MinSpend          *big.Int
}

type AutomationRegistryBase23BillingOverrides struct {
	GasFeePPB         uint32
	FlatFeeMilliCents *big.Int
}

type AutomationRegistryBase23HotVars struct {
	TotalPremium           *big.Int
	LatestEpoch            uint32
	StalenessSeconds       *big.Int
	GasCeilingMultiplier   uint16
	F                      uint8
	Paused                 bool
	ReentrancyGuard        bool
	ReorgProtectionEnabled bool
	ChainModule            common.Address
}

type AutomationRegistryBase23OnchainConfig struct {
	CheckGasLimit          uint32
	MaxPerformGas          uint32
	MaxCheckDataSize       uint32
	Transcoder             common.Address
	ReorgProtectionEnabled bool
	StalenessSeconds       *big.Int
	MaxPerformDataSize     uint32
	MaxRevertDataSize      uint32
	UpkeepPrivilegeManager common.Address
	GasCeilingMultiplier   uint16
	FinanceAdmin           common.Address
	FallbackGasPrice       *big.Int
	FallbackLinkPrice      *big.Int
	FallbackNativePrice    *big.Int
	Registrars             []common.Address
	ChainModule            common.Address
}

type AutomationRegistryBase23PaymentReceipt struct {
	GasChargeInBillingToken *big.Int
	PremiumInBillingToken   *big.Int
	GasReimbursementInJuels *big.Int
	PremiumInJuels          *big.Int
	BillingToken            common.Address
	LinkUSD                 *big.Int
	NativeUSD               *big.Int
	BillingUSD              *big.Int
}

type AutomationRegistryBase23Storage struct {
	Transcoder              common.Address
	CheckGasLimit           uint32
	MaxPerformGas           uint32
	Nonce                   uint32
	UpkeepPrivilegeManager  common.Address
	ConfigCount             uint32
	LatestConfigBlockNumber uint32
	MaxCheckDataSize        uint32
	FinanceAdmin            common.Address
	MaxPerformDataSize      uint32
	MaxRevertDataSize       uint32
}

type AutomationRegistryBase23TransmitterPayeeInfo struct {
	TransmitterAddress common.Address
	PayeeAddress       common.Address
}

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

var AutomationRegistryLogicC23MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"linkUSDFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nativeUSDFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fastGasFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"automationForwarderLogic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedReadOnlyAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payoutMode\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\"},{\"name\":\"wrappedNativeTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"disableOffchainPayments\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveUpkeepIDs\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowedReadOnlyAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAutomationForwarderLogic\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAvailableERC20ForPayment\",\"inputs\":[{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingConfig\",\"inputs\":[{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingOverrides\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingOverridesEnabled\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingToken\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingTokenConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"contractIERC20Metadata[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCancellationDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getChainModule\",\"inputs\":[],\"outputs\":[{\"name\":\"chainModule\",\"type\":\"address\",\"internalType\":\"contractIChainModule\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConditionalGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"components\":[{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"financeAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackNativePrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"chainModule\",\"type\":\"address\",\"internalType\":\"contractIChainModule\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFallbackNativePrice\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFastGasFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getForwarder\",\"inputs\":[{\"name\":\"upkeepID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAutomationForwarder\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHotVars\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.HotVars\",\"components\":[{\"name\":\"totalPremium\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"latestEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reentrancyGuard\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"chainModule\",\"type\":\"address\",\"internalType\":\"contractIChainModule\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkUSDFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLogGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getMaxPaymentForGas\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"maxPayment\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalance\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalanceForUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"minBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNativeUSDFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumUpkeeps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPayoutMode\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPerPerformByteGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPerSignerGasOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getReorgProtectionEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReserveAmount\",\"inputs\":[{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSignerInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ownerLinkBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"expectedLinkBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalPremium\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"numUpkeeps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"latestEpoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.Storage\",\"components\":[{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"financeAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTransmitterInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"lastCollected\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmittersWithPayees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structAutomationRegistryBase2_3.TransmitterPayeeInfo[]\",\"components\":[{\"name\":\"transmitterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payeeAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTriggerType\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepInfo\",\"type\":\"tuple\",\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"amountSpent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWrappedNativeTokenAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasDedupKey\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"linkAvailableForPayment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdminPrivilegeConfig\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepPrivilegeConfig\",\"inputs\":[{\"name\":\"upkeepId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleNOPsOffchain\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsBillingToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upkeepVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverridden\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"overrides\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverrideRemoved\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigSet\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainSpecificModuleUpdated\",\"inputs\":[{\"name\":\"newModule\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeesWithdrawn\",\"inputs\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NOPsSettledOffchain\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payments\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCharged\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receipt\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"components\":[{\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"linkUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"nativeUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLinkLiquidity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOffchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOnchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyFinanceAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x6101606040523480156200001257600080fd5b5060405162005c8738038062005c87833981016040819052620000359162000309565b87878787878787873380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c78162000241565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff191660018381811115620001185762000118620003b6565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000179573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200019f9190620003cc565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002099190620003cc565b60ff16146200022b576040516301f86e1760e41b815260040160405180910390fd5b50505050505050505050505050505050620003f8565b336001600160a01b038216036200029b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200030457600080fd5b919050565b600080600080600080600080610100898b0312156200032757600080fd5b6200033289620002ec565b97506200034260208a01620002ec565b96506200035260408a01620002ec565b95506200036260608a01620002ec565b94506200037260808a01620002ec565b93506200038260a08a01620002ec565b925060c0890151600281106200039757600080fd5b9150620003a760e08a01620002ec565b90509295985092959890939650565b634e487b7160e01b600052602160045260246000fd5b600060208284031215620003df57600080fd5b815160ff81168114620003f157600080fd5b9392505050565b60805160a05160c05160e051610100516101205161014051615803620004846000396000610b7001526000610ad2015260006107be0152600081816109a001526135b60152600081816109540152613e3201526000818161051d0152613690015260008181610c1801528181611df00152818161256b015281816125c80152613af601526158036000f3fe608060405234801561001057600080fd5b50600436106103d05760003560e01c80638ed02bab116101ff578063c3f909d41161011a578063eb5dcd6c116100ad578063f5a418461161007c578063f5a4184614610eb4578063f777ff0614610f35578063faa3e99614610f3c578063ffd242bd14610f8257600080fd5b8063eb5dcd6c14610de5578063ec4de4ba14610df8578063ed56b3e114610e2e578063f2fde38b14610ea157600080fd5b8063d09dc339116100e9578063d09dc33914610c3c578063d763264814610c44578063d85aa07c14610c57578063e80a3d4414610c5f57600080fd5b8063c3f909d414610bd6578063c5b964e014610beb578063c7c3a19a14610bf6578063ca30e60314610c1657600080fd5b8063aab9edd611610192578063b3596c2311610161578063b3596c23146107e2578063b6511a2a14610ba7578063b657bc9c14610bae578063ba87666814610bc157600080fd5b8063aab9edd614610b57578063abc76ae014610b66578063ac4dc59a14610b6e578063b121e14714610b9457600080fd5b8063a08714c0116101ce578063a08714c014610ad0578063a538b2eb14610af6578063a710b22114610b3c578063a87f45fe14610b4f57600080fd5b80638ed02bab14610a5c5780639089daa414610a7a57806393f6ebcf14610a8f5780639e0a99ed14610ac857600080fd5b806343cc055c116102ef5780636209e1e91161028257806379ba50971161025157806379ba5097146109ea57806379ea9943146109f25780638456cb5914610a365780638da5cb5b14610a3e57600080fd5b80636209e1e91461098b5780636709d0e51461099e578063671d36ed146109c45780636eec02a2146109d757600080fd5b80635425d8ac116102be5780635425d8ac146107bc57806357359584146107e2578063614486af146109525780636181d82d1461097857600080fd5b806343cc055c1461074957806344cb70b8146107705780634ca16c52146107935780635147cd591461079c57600080fd5b8063207b6516116103675780633b9cce59116103365780633b9cce59146106c05780633f4ba83a146106d3578063421d183b146106db57806343b46e5f1461074157600080fd5b8063207b651614610508578063226cf83c1461051b578063232c1cc5146105625780633408f73a1461056957600080fd5b8063187256e8116103a3578063187256e81461043b57806319d97a941461044e5780631e0104391461046e5780631efcf646146104d057600080fd5b8063050ee65d146103d557806306e3b632146103ed5780630b7d33e61461040d5780631865c57d14610422575b600080fd5b6201e26c5b6040519081526020015b60405180910390f35b6104006103fb36600461447a565b610f8a565b6040516103e491906144d7565b61042061041b366004614533565b6110a7565b005b61042a611108565b6040516103e495949392919061472b565b61042061044936600461486c565b611508565b61046161045c3660046148a9565b611579565b6040516103e49190614926565b6104b361047c3660046148a9565b60009081526004602052604090206001015470010000000000000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff90911681526020016103e4565b6104f86104de3660046148a9565b600090815260046020526040902054610100900460ff1690565b60405190151581526020016103e4565b6104616105163660046148a9565b61161b565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016103e4565b60186103da565b6106b36040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081019190915250604080516101608101825260165473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008084048216602086015278010000000000000000000000000000000000000000000000008085048316968601969096527c010000000000000000000000000000000000000000000000000000000093849004821660608601526017548084166080870152818104831660a0870152868104831660c087015293909304811660e085015260185491821661010085015291810482166101208401529290920490911661014082015290565b6040516103e49190614939565b6104206106ce366004614a63565b611638565b61042061188e565b6106ee6106e9366004614ad8565b6118f4565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a0016103e4565b610420611a13565b6014547801000000000000000000000000000000000000000000000000900460ff166104f8565b6104f861077e3660046148a9565b60009081526008602052604090205460ff1690565b62017f986103da565b6107af6107aa3660046148a9565b611eaa565b6040516103e49190614b34565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6108d46107f0366004614ad8565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091525073ffffffffffffffffffffffffffffffffffffffff908116600090815260226020908152604091829020825160c081018452815463ffffffff81168252640100000000810462ffffff16938201939093526701000000000000008304909416928401929092527b01000000000000000000000000000000000000000000000000000000900460ff16606083015260018101546080830152600201546bffffffffffffffffffffffff1660a082015290565b6040516103e49190600060c08201905063ffffffff835116825262ffffff602084015116602083015273ffffffffffffffffffffffffffffffffffffffff604084015116604083015260ff6060840151166060830152608083015160808301526bffffffffffffffffffffffff60a08401511660a083015292915050565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104b3610986366004614b47565b611eb5565b610461610999366004614ad8565b612010565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104206109d2366004614ba7565b612043565b6103da6109e5366004614ad8565b6120c3565b61042061216d565b61053d610a003660046148a9565b6000908152600460205260409020546a0100000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b61042061226f565b60005473ffffffffffffffffffffffffffffffffffffffff1661053d565b60155473ffffffffffffffffffffffffffffffffffffffff1661053d565b610a826122e8565b6040516103e49190614be3565b61053d610a9d3660046148a9565b60009081526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff1690565b6103a46103da565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104f8610b04366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff908116600090815260226020526040902054670100000000000000900416151590565b610420610b4a366004614c4c565b6123fa565b610420612728565b604051600481526020016103e4565b6115e06103da565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b610420610ba2366004614ad8565b61275a565b60326103da565b6104b3610bbc3660046148a9565b612852565b610bc9612979565b6040516103e49190614c7a565b610bde6129e8565b6040516103e49190614cd4565b60255460ff166107af565b610c09610c043660046148a9565b612bd0565b6040516103e49190614e5b565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6103da612fe0565b6104b3610c523660046148a9565b612fef565b601b546103da565b610dd86040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101919091525060408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff1661010082015290565b6040516103e49190614f92565b610420610df3366004614c4c565b612ffa565b6103da610e06366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff1660009081526021602052604090205490565b610e88610e3c366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff9091166020830152016103e4565b610420610eaf366004614ad8565b613159565b610f0f610ec23660046148a9565b60408051808201909152600080825260208201525060009081526023602090815260409182902082518084019093525463ffffffff81168352640100000000900462ffffff169082015290565b60408051825163ffffffff16815260209283015162ffffff1692810192909252016103e4565b60406103da565b610f75610f4a366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff166000908152601c602052604090205460ff1690565b6040516103e49190615062565b6103da61316d565b60606000610f986002613175565b9050808410610fd3576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610fdf84866150a5565b905081811180610fed575083155b610ff75780610ff9565b815b9050600061100786836150b8565b67ffffffffffffffff81111561101f5761101f6150cb565b604051908082528060200260200182016040528015611048578160200160208202803683370190505b50905060005b815181101561109b5761106c61106488836150a5565b60029061317f565b82828151811061107e5761107e6150fa565b60209081029190910101528061109381615129565b91505061104e565b50925050505b92915050565b6110af61318b565b6000838152601f602052604090206110c8828483615203565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae776983836040516110fb92919061531e565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c081019190915260408051610140810182526016547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1681526000602082018190529181018290526014546bffffffffffffffffffffffff16606080830191909152918291608081016112416002613175565b81526017547401000000000000000000000000000000000000000080820463ffffffff908116602080860191909152780100000000000000000000000000000000000000000000000080850483166040808801919091526013546060808901919091526014546c01000000000000000000000000810486166080808b0191909152760100000000000000000000000000000000000000000000820460ff16151560a09a8b015283516101e0810185526000808252968101879052601654898104891695820195909552700100000000000000000000000000000000830462ffffff169381019390935273010000000000000000000000000000000000000090910461ffff169082015296870192909252808204831660c08701527c0100000000000000000000000000000000000000000000000000000000909404821660e086015260185492830482166101008601529290910416610120830152601954610140830152601a5461016083015273ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a081016113dc60096131de565b815260175473ffffffffffffffffffffffffffffffffffffffff16602091820152601454600d80546040805182860281018601909152818152949850899489949293600e93750100000000000000000000000000000000000000000090910460ff1692859183018282801561148757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161145c575b50505050509250818054806020026020016040519081016040528092919081815260200182805480156114f057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116114c5575b50505050509150945094509450945094509091929394565b6115106131eb565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601c6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561157057611570614af5565b02179055505050565b6000818152601f6020526040902080546060919061159690615161565b80601f01602080910402602001604051908101604052809291908181526020018280546115c290615161565b801561160f5780601f106115e45761010080835404028352916020019161160f565b820191906000526020600020905b8154815290600101906020018083116115f257829003601f168201915b50505050509050919050565b6000818152601d6020526040902080546060919061159690615161565b6116406131eb565b600e54811461167b576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e5481101561184d576000600e828154811061169d5761169d6150fa565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168084526011909252604083205491935016908585858181106116e7576116e76150fa565b90506020020160208101906116fc9190614ad8565b905073ffffffffffffffffffffffffffffffffffffffff8116158061178f575073ffffffffffffffffffffffffffffffffffffffff82161580159061176d57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b801561178f575073ffffffffffffffffffffffffffffffffffffffff81811614155b156117c6576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146118375773ffffffffffffffffffffffffffffffffffffffff838116600090815260116020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061184590615129565b91505061167e565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e83836040516118829392919061536b565b60405180910390a15050565b6118966131eb565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015282918291829182919082906119ba5760608201516014546000916119a6916bffffffffffffffffffffffff1661541d565b600e549091506119b69082615471565b9150505b8151602083015160408401516119d190849061549c565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b16600090815260116020526040902054929b919a9499509750921694509092505050565b611a1b61326c565b600060255460ff166001811115611a3457611a34614af5565b03611a6b576040517fe0262d7400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e546bffffffffffffffffffffffff909116906000611a8e600f613175565b90506000611a9c82846150a5565b905060008167ffffffffffffffff811115611ab957611ab96150cb565b604051908082528060200260200182016040528015611ae2578160200160208202803683370190505b50905060008267ffffffffffffffff811115611b0057611b006150cb565b604051908082528060200260200182016040528015611b29578160200160208202803683370190505b50905060005b85811015611c5b576000600e8281548110611b4c57611b4c6150fa565b600091825260208220015473ffffffffffffffffffffffffffffffffffffffff169150611b7a828a8a6132bd565b9050806bffffffffffffffffffffffff16858481518110611b9d57611b9d6150fa565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff808416600090815260119092526040909120548551911690859085908110611beb57611beb6150fa565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920181019190915292166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611c5381615129565b915050611b2f565b5060005b84811015611dd8576000611c74600f8361317f565b73ffffffffffffffffffffffffffffffffffffffff8082166000818152600b602090815260408083208151608081018352905460ff80821615158352610100820416828501526bffffffffffffffffffffffff6201000082048116838501526e0100000000000000000000000000009091041660608201529383526011909152902054929350911684611d078a866150a5565b81518110611d1757611d176150fa565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260408101516bffffffffffffffffffffffff1685611d5a8a866150a5565b81518110611d6a57611d6a6150fa565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff9092166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611dd081615129565b915050611c5f565b5073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166000908152602160205260408120819055611e2b600f613175565b90505b8015611e6857611e55611e4d611e456001846150b8565b600f9061317f565b600f906134c5565b5080611e60816154c1565b915050611e2e565b507f5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad8183604051611e9a9291906154f6565b60405180910390a1505050505050565b60006110a1826134e7565b60408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff16610100820152600090818080611fed84613592565b92509250925061200389858a8a8787878d613784565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152602080526040902080546060919061159690615161565b61204b61318b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602080526040902061207a828483615203565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d283836040516110fb92919061531e565b73ffffffffffffffffffffffffffffffffffffffff81166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa15801561213f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121639190615524565b6110a191906150b8565b60015473ffffffffffffffffffffffffffffffffffffffff1633146121f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6122776131eb565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff167601000000000000000000000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020016118ea565b600e5460609060008167ffffffffffffffff811115612309576123096150cb565b60405190808252806020026020018201604052801561234e57816020015b60408051808201909152600080825260208201528152602001906001900390816123275790505b50905060005b828110156123f3576000600e8281548110612371576123716150fa565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168084526011835260409384902054845180860190955281855290911691830182905285519093509091908590859081106123d3576123d36150fa565b6020026020010181905250505080806123eb90615129565b915050612354565b5092915050565b73ffffffffffffffffffffffffffffffffffffffff8116612447576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160255460ff16600181111561246057612460614af5565b03612497576040517f4a3578fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601160205260409020541633146124f7576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e5460009161251a9185916bffffffffffffffffffffffff16906132bd565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600b6020908152604080832080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff1690557f000000000000000000000000000000000000000000000000000000000000000090931682526021905220549091506125b1906bffffffffffffffffffffffff8316906150b8565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081166000818152602160205260408082209490945592517fa9059cbb00000000000000000000000000000000000000000000000000000000815291851660048301526bffffffffffffffffffffffff841660248301529063a9059cbb906044016020604051808303816000875af1158015612666573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061268a919061553d565b9050806126c3576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405133815273ffffffffffffffffffffffffffffffffffffffff808516916bffffffffffffffffffffffff8516918716907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a450505050565b6127306131eb565b602580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152601260205260409020541633146127ba576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526011602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556012909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490931660e08201526002909101549091169181019190915261297283612962816134e7565b8360400151846101000151611eb5565b9392505050565b606060248054806020026020016040519081016040528092919081815260200182805480156129de57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116129b3575b5050505050905090565b604080516102008101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a082018390526101c08201526101e0810191909152604080516102008101825260165463ffffffff74010000000000000000000000000000000000000000808304821684527801000000000000000000000000000000000000000000000000808404831660208601526017547c0100000000000000000000000000000000000000000000000000000000810484169686019690965273ffffffffffffffffffffffffffffffffffffffff938416606086015260145460ff828204161515608087015262ffffff70010000000000000000000000000000000082041660a0870152601854928304841660c087015290820490921660e085015293821661010084015261ffff7301000000000000000000000000000000000000009091041661012083015291909116610140820152601954610160820152601a54610180820152601b546101a08201526101c08101612baa60096131de565b815260155473ffffffffffffffffffffffffffffffffffffffff16602090910152919050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a010000000000000000000090910481166080830181905260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c0100000000000000000000000000000000000000000000000000000000900490941660e08301526002909201549091169281019290925290919015612da557816080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612d7c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612da0919061555f565b612da8565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff168152602001600760008781526020019081526020016000208054612e0090615161565b80601f0160208091040260200160405190810160405280929190818152602001828054612e2c90615161565b8015612e795780601f10612e4e57610100808354040283529160200191612e79565b820191906000526020600020905b815481529060010190602001808311612e5c57829003601f168201915b505050505081526020018360c001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836060015163ffffffff1667ffffffffffffffff1681526020018360e0015163ffffffff1681526020018360a001516bffffffffffffffffffffffff168152602001836000015115158152602001601e60008781526020019081526020016000208054612f5690615161565b80601f0160208091040260200160405190810160405280929190818152602001828054612f8290615161565b8015612fcf5780601f10612fa457610100808354040283529160200191612fcf565b820191906000526020600020905b815481529060010190602001808311612fb257829003601f168201915b505050505081525092505050919050565b6000612fea613af4565b905090565b60006110a182612852565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526011602052604090205416331461305a576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216036130a9576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601260205260409020548116908216146131555773ffffffffffffffffffffffffffffffffffffffff82811660008181526012602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45b5050565b6131616131eb565b61316a81613bbe565b50565b6000612fea60025b60006110a1825490565b60006129728383613cb3565b60175473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6060600061297283613cdd565b60005473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016121ea565b60185473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906134b9576000816060015185613355919061541d565b905060006133638583615471565b90508083604001818151613377919061549c565b6bffffffffffffffffffffffff16905250613392858261557c565b836060018181516133a3919061549c565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b60006129728373ffffffffffffffffffffffffffffffffffffffff8416613d38565b6000818160045b600f811015613574577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061352c5761352c6150fa565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461356257506000949350505050565b8061356c81615129565b9150506134ee565b5081600f1a600181111561358a5761358a614af5565b949350505050565b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561361f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061364391906155c3565b509450909250505060008113158061365a57508142105b8061367b575082801561367b575061367282426150b8565b8463ffffffff16105b1561368a57601954965061368e565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156136f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061371d91906155c3565b509450909250505060008113158061373457508142105b806137555750828015613755575061374c82426150b8565b8463ffffffff16105b1561376457601a549550613768565b8095505b86866137738a613e2b565b965096509650505050509193909250565b600080808089600181111561379b5761379b614af5565b036137aa575062017f986137ff565b60018960018111156137be576137be614af5565b036137cd57506201e26c6137ff565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a6080015160016138129190615613565b6138209060ff16604061562c565b60185461384e906103a49074010000000000000000000000000000000000000000900463ffffffff166150a5565b61385891906150a5565b601554604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa1580156138c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138ed9190615643565b909250905081836138ff8360186150a5565b613909919061562c565b60808f0151613919906001615613565b6139289060ff166115e061562c565b61393291906150a5565b61393c91906150a5565b61394690856150a5565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa1580156139b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139dd9190615524565b8d6060015161ffff166139f0919061562c565b94505050506000613a018b86613f15565b60008d815260046020526040902054909150610100900460ff1615613a665760008c81526023602090815260409182902082518084018452905463ffffffff811680835264010000000090910462ffffff908116928401928352928501525116908201525b6000613ad08c6040518061012001604052808d63ffffffff1681526020018681526020018781526020018c81526020018b81526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff16815260200185815260200160001515815250614091565b60208101518151919250613ae39161549c565b9d9c50505050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015613b90573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bb49190615524565b612fea9190615667565b3373ffffffffffffffffffffffffffffffffffffffff821603613c3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016121ea565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613cca57613cca6150fa565b9060005260206000200154905092915050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561160f57602002820191906000526020600020905b815481526020019060010190808311613d195750505050509050919050565b60008181526001830160205260408120548015613e21576000613d5c6001836150b8565b8554909150600090613d70906001906150b8565b9050818114613dd5576000866000018281548110613d9057613d906150fa565b9060005260206000200154905080876000018481548110613db357613db36150fa565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613de657613de6615687565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506110a1565b60009150506110a1565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613e9b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ebf91906155c3565b50935050925050600082131580613ed557508042105b80613f0557506000846040015162ffffff16118015613f055750613ef981426150b8565b846040015162ffffff16105b156123f3575050601b5492915050565b60408051608081018252600080825260208083018281528385018381526060850184905273ffffffffffffffffffffffffffffffffffffffff878116855260229093528584208054640100000000810462ffffff1690925263ffffffff82169092527b01000000000000000000000000000000000000000000000000000000810460ff16855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495919484936701000000000000009092049091169163feaf968c9160048083019260a09291908290030181865afa158015614002573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061402691906155c3565b5093505092505060008213158061403c57508042105b8061406c57506000866040015162ffffff1611801561406c575061406081426150b8565b866040015162ffffff16105b156140805760018301546060850152614088565b606084018290525b50505092915050565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081019190915260008260e001516000015160ff1690506000846060015161ffff1684606001516140fc919061562c565b9050836101000151801561410f5750803a105b1561411757503a5b60006012831161412857600161413e565b6141336012846150b8565b61413e90600a6157d6565b9050600060128410614151576001614167565b61415c8460126150b8565b61416790600a6157d6565b905060008660a0015187604001518860200151896000015161418991906150a5565b614193908761562c565b61419d91906150a5565b6141a7919061562c565b9050614203828860e00151606001516141c0919061562c565b6001848a60e00151606001516141d6919061562c565b6141e091906150b8565b6141ea868561562c565b6141f491906150a5565b6141fe91906157e2565b6143d8565b6bffffffffffffffffffffffff1686526080870151614226906141fe90836157e2565b6bffffffffffffffffffffffff1660408088019190915260e0880151015160009061425f9062ffffff16683635c9adc5dea0000061562c565b9050600081633b9aca008a60a001518b60e001516020015163ffffffff168c604001518d600001518b614292919061562c565b61429c91906150a5565b6142a6919061562c565b6142b0919061562c565b6142ba91906157e2565b6142c491906150a5565b9050614307848a60e00151606001516142dd919061562c565b6001868c60e00151606001516142f3919061562c565b6142fd91906150b8565b6141ea888561562c565b6bffffffffffffffffffffffff166020890152608089015161432d906141fe90836157e2565b6bffffffffffffffffffffffff16606089015260c089015173ffffffffffffffffffffffffffffffffffffffff166080808a0191909152890151614370906143d8565b6bffffffffffffffffffffffff1660a0808a0191909152890151614393906143d8565b6bffffffffffffffffffffffff1660c089015260e0890151606001516143b8906143d8565b6bffffffffffffffffffffffff1660e08901525050505050505092915050565b60006bffffffffffffffffffffffff821115614476576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016121ea565b5090565b6000806040838503121561448d57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156144cc578151875295820195908201906001016144b0565b509495945050505050565b602081526000612972602083018461449c565b60008083601f8401126144fc57600080fd5b50813567ffffffffffffffff81111561451457600080fd5b60208301915083602082850101111561452c57600080fd5b9250929050565b60008060006040848603121561454857600080fd5b83359250602084013567ffffffffffffffff81111561456657600080fd5b614572868287016144ea565b9497909650939450505050565b600081518084526020808501945080840160005b838110156144cc57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614593565b805163ffffffff16825260006101e060208301516145eb602086018263ffffffff169052565b506040830151614603604086018263ffffffff169052565b50606083015161461a606086018262ffffff169052565b506080830151614630608086018261ffff169052565b5060a083015161465060a08601826bffffffffffffffffffffffff169052565b5060c083015161466860c086018263ffffffff169052565b5060e083015161468060e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a0808401518186018390526146f58387018261457f565b925050506101c0808401516147218287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161475960208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161478360608501826bffffffffffffffffffffffff169052565b506080880151608084015260a08801516147a560a085018263ffffffff169052565b5060c08801516147bd60c085018263ffffffff169052565b5060e088015160e0840152610100808901516147e08286018263ffffffff169052565b5050610120888101511515908401526101408301819052614803818401886145c5565b9050828103610160840152614818818761457f565b905082810361018084015261482d818661457f565b9150506148406101a083018460ff169052565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461316a57600080fd5b6000806040838503121561487f57600080fd5b823561488a8161484a565b915060208301356004811061489e57600080fd5b809150509250929050565b6000602082840312156148bb57600080fd5b5035919050565b6000815180845260005b818110156148e8576020818501810151868301820152016148cc565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061297260208301846148c2565b815173ffffffffffffffffffffffffffffffffffffffff1681526101608101602083015161496f602084018263ffffffff169052565b506040830151614987604084018263ffffffff169052565b50606083015161499f606084018263ffffffff169052565b5060808301516149c7608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a08301516149df60a084018263ffffffff169052565b5060c08301516149f760c084018263ffffffff169052565b5060e0830151614a0f60e084018263ffffffff169052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff81168483015250506101208381015163ffffffff81168483015250506101408381015163ffffffff8116848301525b505092915050565b60008060208385031215614a7657600080fd5b823567ffffffffffffffff80821115614a8e57600080fd5b818501915085601f830112614aa257600080fd5b813581811115614ab157600080fd5b8660208260051b8501011115614ac657600080fd5b60209290920196919550909350505050565b600060208284031215614aea57600080fd5b81356129728161484a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6002811061316a5761316a614af5565b60208101614b4183614b24565b91905290565b60008060008060808587031215614b5d57600080fd5b84359350602085013560028110614b7357600080fd5b9250604085013563ffffffff81168114614b8c57600080fd5b91506060850135614b9c8161484a565b939692955090935050565b600080600060408486031215614bbc57600080fd5b8335614bc78161484a565b9250602084013567ffffffffffffffff81111561456657600080fd5b602080825282518282018190526000919060409081850190868401855b82811015614c3f578151805173ffffffffffffffffffffffffffffffffffffffff90811686529087015116868501529284019290850190600101614c00565b5091979650505050505050565b60008060408385031215614c5f57600080fd5b8235614c6a8161484a565b9150602083013561489e8161484a565b6020808252825182820181905260009190848201906040850190845b81811015614cc857835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101614c96565b50909695505050505050565b60208152614ceb60208201835163ffffffff169052565b60006020830151614d04604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015173ffffffffffffffffffffffffffffffffffffffff8116608084015250608083015180151560a08401525060a083015162ffffff811660c08401525060c083015163ffffffff811660e08401525060e0830151610100614d838185018363ffffffff169052565b8401519050610120614dac8482018373ffffffffffffffffffffffffffffffffffffffff169052565b8401519050610140614dc38482018361ffff169052565b8401519050610160614dec8482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151610180848101919091528401516101a0808501919091528401516101c0808501919091528401516102006101e080860182905291925090614e3461022086018461457f565b9086015173ffffffffffffffffffffffffffffffffffffffff811683870152909250614721565b60208152614e8260208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614e9b604084018263ffffffff169052565b506040830151610140806060850152614eb86101608501836148c2565b91506060850151614ed960808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614f45818701836bffffffffffffffffffffffff169052565b8601519050610120614f5a8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061484083826148c2565b6000610120820190506bffffffffffffffffffffffff835116825263ffffffff60208401511660208301526040830151614fd3604084018262ffffff169052565b506060830151614fe9606084018261ffff169052565b506080830151614ffe608084018260ff169052565b5060a083015161501260a084018215159052565b5060c083015161502660c084018215159052565b5060e083015161503a60e084018215159052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff811684830152614a5b565b6020810160048310614b4157614b41614af5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156110a1576110a1615076565b818103818111156110a1576110a1615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361515a5761515a615076565b5060010190565b600181811c9082168061517557607f821691505b6020821081036151ae577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156151fe57600081815260208120601f850160051c810160208610156151db5750805b601f850160051c820191505b818110156151fa578281556001016151e7565b5050505b505050565b67ffffffffffffffff83111561521b5761521b6150cb565b61522f836152298354615161565b836151b4565b6000601f841160018114615281576000851561524b5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355615317565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156152d057868501358255602094850194600190920191016152b0565b508682101561530b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000604082016040835280865480835260608501915087600052602092508260002060005b828110156153c257815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101615390565b505050838103828501528481528590820160005b868110156154115782356153e98161484a565b73ffffffffffffffffffffffffffffffffffffffff16825291830191908301906001016153d6565b50979650505050505050565b6bffffffffffffffffffffffff8281168282160390808211156123f3576123f3615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff8084168061549057615490615442565b92169190910492915050565b6bffffffffffffffffffffffff8181168382160190808211156123f3576123f3615076565b6000816154d0576154d0615076565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b604081526000615509604083018561457f565b828103602084015261551b818561449c565b95945050505050565b60006020828403121561553657600080fd5b5051919050565b60006020828403121561554f57600080fd5b8151801515811461297257600080fd5b60006020828403121561557157600080fd5b81516129728161484a565b6bffffffffffffffffffffffff818116838216028082169190828114614a5b57614a5b615076565b805169ffffffffffffffffffff811681146155be57600080fd5b919050565b600080600080600060a086880312156155db57600080fd5b6155e4866155a4565b9450602086015193506040860151925060608601519150615607608087016155a4565b90509295509295909350565b60ff81811683821601908111156110a1576110a1615076565b80820281158282048414176110a1576110a1615076565b6000806040838503121561565657600080fd5b505080516020909101519092909150565b81810360008312801583831316838312821617156123f3576123f3615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600181815b8085111561570f57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156156f5576156f5615076565b8085161561570257918102915b93841c93908002906156bb565b509250929050565b600082615726575060016110a1565b81615733575060006110a1565b816001811461574957600281146157535761576f565b60019150506110a1565b60ff84111561576457615764615076565b50506001821b6110a1565b5060208310610133831016604e8410600b8410161715615792575081810a6110a1565b61579c83836156b6565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156157ce576157ce615076565b029392505050565b60006129728383615717565b6000826157f1576157f1615442565b50049056fea164736f6c6343000813000a",
}

var AutomationRegistryLogicC23ABI = AutomationRegistryLogicC23MetaData.ABI

var AutomationRegistryLogicC23Bin = AutomationRegistryLogicC23MetaData.Bin

func DeployAutomationRegistryLogicC23(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkUSDFeed common.Address, nativeUSDFeed common.Address, fastGasFeed common.Address, automationForwarderLogic common.Address, allowedReadOnlyAddress common.Address, payoutMode uint8, wrappedNativeTokenAddress common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicC23, error) {
	parsed, err := AutomationRegistryLogicC23MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicC23Bin), backend, link, linkUSDFeed, nativeUSDFeed, fastGasFeed, automationForwarderLogic, allowedReadOnlyAddress, payoutMode, wrappedNativeTokenAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicC23{address: address, abi: *parsed, AutomationRegistryLogicC23Caller: AutomationRegistryLogicC23Caller{contract: contract}, AutomationRegistryLogicC23Transactor: AutomationRegistryLogicC23Transactor{contract: contract}, AutomationRegistryLogicC23Filterer: AutomationRegistryLogicC23Filterer{contract: contract}}, nil
}

type AutomationRegistryLogicC23 struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicC23Caller
	AutomationRegistryLogicC23Transactor
	AutomationRegistryLogicC23Filterer
}

type AutomationRegistryLogicC23Caller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicC23Transactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicC23Filterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicC23Session struct {
	Contract     *AutomationRegistryLogicC23
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicC23CallerSession struct {
	Contract *AutomationRegistryLogicC23Caller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicC23TransactorSession struct {
	Contract     *AutomationRegistryLogicC23Transactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicC23Raw struct {
	Contract *AutomationRegistryLogicC23
}

type AutomationRegistryLogicC23CallerRaw struct {
	Contract *AutomationRegistryLogicC23Caller
}

type AutomationRegistryLogicC23TransactorRaw struct {
	Contract *AutomationRegistryLogicC23Transactor
}

func NewAutomationRegistryLogicC23(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicC23, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicC23ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicC23(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23{address: address, abi: abi, AutomationRegistryLogicC23Caller: AutomationRegistryLogicC23Caller{contract: contract}, AutomationRegistryLogicC23Transactor: AutomationRegistryLogicC23Transactor{contract: contract}, AutomationRegistryLogicC23Filterer: AutomationRegistryLogicC23Filterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicC23Caller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicC23Caller, error) {
	contract, err := bindAutomationRegistryLogicC23(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23Caller{contract: contract}, nil
}

func NewAutomationRegistryLogicC23Transactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicC23Transactor, error) {
	contract, err := bindAutomationRegistryLogicC23(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23Transactor{contract: contract}, nil
}

func NewAutomationRegistryLogicC23Filterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicC23Filterer, error) {
	contract, err := bindAutomationRegistryLogicC23(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23Filterer{contract: contract}, nil
}

func bindAutomationRegistryLogicC23(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicC23MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicC23.Contract.AutomationRegistryLogicC23Caller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AutomationRegistryLogicC23Transactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AutomationRegistryLogicC23Transactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicC23.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicC23.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicC23.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicC23.CallOpts, admin)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicC23.CallOpts, admin)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetAvailableERC20ForPayment(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getAvailableERC20ForPayment", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetAvailableERC20ForPayment(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetAvailableERC20ForPayment(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetBalance(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetBalance(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingConfig(opts *bind.CallOpts, billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingConfig", billingToken)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingConfig(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingConfig(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingOverrides(opts *bind.CallOpts, upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingOverrides", upkeepID)

	if err != nil {
		return *new(AutomationRegistryBase23BillingOverrides), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingOverrides)).(*AutomationRegistryBase23BillingOverrides)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingOverrides(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingOverrides(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingOverridesEnabled(opts *bind.CallOpts, upkeepID *big.Int) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingOverridesEnabled", upkeepID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingOverridesEnabled(upkeepID *big.Int) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingOverridesEnabled(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingOverridesEnabled(upkeepID *big.Int) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingOverridesEnabled(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingToken", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingToken(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingToken(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingTokenConfig", token)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicC23.CallOpts, token)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicC23.CallOpts, token)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getBillingTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingTokens(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetBillingTokens(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetCancellationDelay(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetCancellationDelay(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetChainModule(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetChainModule(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetConfig(opts *bind.CallOpts) (AutomationRegistryBase23OnchainConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(AutomationRegistryBase23OnchainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23OnchainConfig)).(*AutomationRegistryBase23OnchainConfig)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetConfig(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _AutomationRegistryLogicC23.Contract.GetConfig(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getFallbackNativePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetForwarder(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetForwarder(&_AutomationRegistryLogicC23.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getHotVars")

	if err != nil {
		return *new(AutomationRegistryBase23HotVars), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23HotVars)).(*AutomationRegistryBase23HotVars)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicC23.Contract.GetHotVars(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicC23.Contract.GetHotVars(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetLinkAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetLinkAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getLinkUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetLogGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetLogGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetMaxPaymentForGas(opts *bind.CallOpts, id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getMaxPaymentForGas", id, triggerType, gasLimit, billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicC23.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicC23.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMinBalance(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMinBalance(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetNativeUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getNativeUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetNumUpkeeps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getNumUpkeeps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetNumUpkeeps(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetNumUpkeeps(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetPayoutMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getPayoutMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetPayoutMode(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetPayoutMode(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC23.CallOpts, peer)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC23.CallOpts, peer)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicC23.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicC23.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetReserveAmount(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getReserveAmount", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetReserveAmount(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetReserveAmount(&_AutomationRegistryLogicC23.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetSignerInfo(&_AutomationRegistryLogicC23.CallOpts, query)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetSignerInfo(&_AutomationRegistryLogicC23.CallOpts, query)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getState")

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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetState(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetState(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetStorage(opts *bind.CallOpts) (AutomationRegistryBase23Storage, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getStorage")

	if err != nil {
		return *new(AutomationRegistryBase23Storage), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23Storage)).(*AutomationRegistryBase23Storage)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicC23.Contract.GetStorage(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicC23.Contract.GetStorage(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitterInfo(&_AutomationRegistryLogicC23.CallOpts, query)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmitterInfo(&_AutomationRegistryLogicC23.CallOpts, query)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetTransmittersWithPayees(opts *bind.CallOpts) ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getTransmittersWithPayees")

	if err != nil {
		return *new([]AutomationRegistryBase23TransmitterPayeeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutomationRegistryBase23TransmitterPayeeInfo)).(*[]AutomationRegistryBase23TransmitterPayeeInfo)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmittersWithPayees(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _AutomationRegistryLogicC23.Contract.GetTransmittersWithPayees(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetTriggerType(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.GetTriggerType(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeep(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeep(&_AutomationRegistryLogicC23.CallOpts, id)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC23.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicC23.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) GetWrappedNativeTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "getWrappedNativeTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetWrappedNativeTokenAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.GetWrappedNativeTokenAddress(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.HasDedupKey(&_AutomationRegistryLogicC23.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.HasDedupKey(&_AutomationRegistryLogicC23.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicC23.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) Owner() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.Owner(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicC23.Contract.Owner(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "supportsBillingToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.SupportsBillingToken(&_AutomationRegistryLogicC23.CallOpts, token)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicC23.Contract.SupportsBillingToken(&_AutomationRegistryLogicC23.CallOpts, token)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Caller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC23.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.UpkeepVersion(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23CallerSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicC23.Contract.UpkeepVersion(&_AutomationRegistryLogicC23.CallOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AcceptOwnership(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AcceptOwnership(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AcceptPayeeship(&_AutomationRegistryLogicC23.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.AcceptPayeeship(&_AutomationRegistryLogicC23.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "disableOffchainPayments")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.DisableOffchainPayments(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.DisableOffchainPayments(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "pause")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.Pause(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.Pause(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicC23.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicC23.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "setPayees", payees)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetPayees(&_AutomationRegistryLogicC23.TransactOpts, payees)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetPayees(&_AutomationRegistryLogicC23.TransactOpts, payees)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC23.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC23.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC23.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC23.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "settleNOPsOffchain")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.TransferOwnership(&_AutomationRegistryLogicC23.TransactOpts, to)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.TransferOwnership(&_AutomationRegistryLogicC23.TransactOpts, to)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.TransferPayeeship(&_AutomationRegistryLogicC23.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.TransferPayeeship(&_AutomationRegistryLogicC23.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "unpause")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.Unpause(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.Unpause(&_AutomationRegistryLogicC23.TransactOpts)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.WithdrawPayment(&_AutomationRegistryLogicC23.TransactOpts, from, to)
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC23.Contract.WithdrawPayment(&_AutomationRegistryLogicC23.TransactOpts, from, to)
}

type AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicC23AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23AdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicC23AdminPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23AdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicC23AdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicC23AdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23BillingConfigOverriddenIterator struct {
	Event *AutomationRegistryLogicC23BillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23BillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23BillingConfigOverridden)
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
		it.Event = new(AutomationRegistryLogicC23BillingConfigOverridden)
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

func (it *AutomationRegistryLogicC23BillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23BillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23BillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23BillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23BillingConfigOverriddenIterator{contract: _AutomationRegistryLogicC23.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23BillingConfigOverridden)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicC23BillingConfigOverridden, error) {
	event := new(AutomationRegistryLogicC23BillingConfigOverridden)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistryLogicC23BillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23BillingConfigOverrideRemoved)
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
		it.Event = new(AutomationRegistryLogicC23BillingConfigOverrideRemoved)
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

func (it *AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23BillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator{contract: _AutomationRegistryLogicC23.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23BillingConfigOverrideRemoved)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicC23BillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistryLogicC23BillingConfigOverrideRemoved)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23BillingConfigSetIterator struct {
	Event *AutomationRegistryLogicC23BillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23BillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23BillingConfigSet)
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
		it.Event = new(AutomationRegistryLogicC23BillingConfigSet)
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

func (it *AutomationRegistryLogicC23BillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23BillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23BillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicC23BillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23BillingConfigSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23BillingConfigSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicC23BillingConfigSet, error) {
	event := new(AutomationRegistryLogicC23BillingConfigSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23CancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicC23CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23CancelledUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicC23CancelledUpkeepReport)
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

func (it *AutomationRegistryLogicC23CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23CancelledUpkeepReportIterator{contract: _AutomationRegistryLogicC23.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23CancelledUpkeepReport)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicC23CancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicC23CancelledUpkeepReport)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicC23ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23ChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistryLogicC23ChainSpecificModuleUpdated)
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

func (it *AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicC23.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23ChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicC23ChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicC23ChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23DedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicC23DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23DedupKeyAdded)
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
		it.Event = new(AutomationRegistryLogicC23DedupKeyAdded)
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

func (it *AutomationRegistryLogicC23DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicC23DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23DedupKeyAddedIterator{contract: _AutomationRegistryLogicC23.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23DedupKeyAdded)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicC23DedupKeyAdded, error) {
	event := new(AutomationRegistryLogicC23DedupKeyAdded)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23FeesWithdrawnIterator struct {
	Event *AutomationRegistryLogicC23FeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23FeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23FeesWithdrawn)
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
		it.Event = new(AutomationRegistryLogicC23FeesWithdrawn)
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

func (it *AutomationRegistryLogicC23FeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23FeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23FeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicC23FeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23FeesWithdrawnIterator{contract: _AutomationRegistryLogicC23.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23FeesWithdrawn)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicC23FeesWithdrawn, error) {
	event := new(AutomationRegistryLogicC23FeesWithdrawn)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23FundsAddedIterator struct {
	Event *AutomationRegistryLogicC23FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23FundsAdded)
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
		it.Event = new(AutomationRegistryLogicC23FundsAdded)
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

func (it *AutomationRegistryLogicC23FundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicC23FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23FundsAddedIterator{contract: _AutomationRegistryLogicC23.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23FundsAdded)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicC23FundsAdded, error) {
	event := new(AutomationRegistryLogicC23FundsAdded)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23FundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicC23FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23FundsWithdrawn)
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
		it.Event = new(AutomationRegistryLogicC23FundsWithdrawn)
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

func (it *AutomationRegistryLogicC23FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23FundsWithdrawnIterator{contract: _AutomationRegistryLogicC23.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23FundsWithdrawn)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicC23FundsWithdrawn, error) {
	event := new(AutomationRegistryLogicC23FundsWithdrawn)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicC23InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23InsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicC23InsufficientFundsUpkeepReport)
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

func (it *AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicC23.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23InsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicC23InsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicC23InsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23NOPsSettledOffchainIterator struct {
	Event *AutomationRegistryLogicC23NOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23NOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23NOPsSettledOffchain)
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
		it.Event = new(AutomationRegistryLogicC23NOPsSettledOffchain)
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

func (it *AutomationRegistryLogicC23NOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23NOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23NOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicC23NOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23NOPsSettledOffchainIterator{contract: _AutomationRegistryLogicC23.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23NOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23NOPsSettledOffchain)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicC23NOPsSettledOffchain, error) {
	event := new(AutomationRegistryLogicC23NOPsSettledOffchain)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23OwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicC23OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23OwnershipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicC23OwnershipTransferRequested)
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

func (it *AutomationRegistryLogicC23OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23OwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicC23.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23OwnershipTransferRequested)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicC23OwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicC23OwnershipTransferRequested)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23OwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicC23OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23OwnershipTransferred)
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
		it.Event = new(AutomationRegistryLogicC23OwnershipTransferred)
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

func (it *AutomationRegistryLogicC23OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23OwnershipTransferredIterator{contract: _AutomationRegistryLogicC23.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23OwnershipTransferred)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicC23OwnershipTransferred, error) {
	event := new(AutomationRegistryLogicC23OwnershipTransferred)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23PausedIterator struct {
	Event *AutomationRegistryLogicC23Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23Paused)
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
		it.Event = new(AutomationRegistryLogicC23Paused)
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

func (it *AutomationRegistryLogicC23PausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicC23PausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23PausedIterator{contract: _AutomationRegistryLogicC23.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23Paused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23Paused)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParsePaused(log types.Log) (*AutomationRegistryLogicC23Paused, error) {
	event := new(AutomationRegistryLogicC23Paused)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23PayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicC23PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23PayeesUpdated)
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
		it.Event = new(AutomationRegistryLogicC23PayeesUpdated)
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

func (it *AutomationRegistryLogicC23PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicC23PayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23PayeesUpdatedIterator{contract: _AutomationRegistryLogicC23.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23PayeesUpdated)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicC23PayeesUpdated, error) {
	event := new(AutomationRegistryLogicC23PayeesUpdated)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23PayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicC23PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23PayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicC23PayeeshipTransferRequested)
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

func (it *AutomationRegistryLogicC23PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23PayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicC23.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23PayeeshipTransferRequested)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicC23PayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicC23PayeeshipTransferRequested)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23PayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicC23PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23PayeeshipTransferred)
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
		it.Event = new(AutomationRegistryLogicC23PayeeshipTransferred)
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

func (it *AutomationRegistryLogicC23PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23PayeeshipTransferredIterator{contract: _AutomationRegistryLogicC23.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23PayeeshipTransferred)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicC23PayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicC23PayeeshipTransferred)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23PaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicC23PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23PaymentWithdrawn)
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
		it.Event = new(AutomationRegistryLogicC23PaymentWithdrawn)
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

func (it *AutomationRegistryLogicC23PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicC23PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23PaymentWithdrawnIterator{contract: _AutomationRegistryLogicC23.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23PaymentWithdrawn)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicC23PaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicC23PaymentWithdrawn)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23ReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicC23ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23ReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicC23ReorgedUpkeepReport)
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

func (it *AutomationRegistryLogicC23ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23ReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicC23.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23ReorgedUpkeepReport)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicC23ReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicC23ReorgedUpkeepReport)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23StaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicC23StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23StaleUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicC23StaleUpkeepReport)
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

func (it *AutomationRegistryLogicC23StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23StaleUpkeepReportIterator{contract: _AutomationRegistryLogicC23.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23StaleUpkeepReport)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicC23StaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicC23StaleUpkeepReport)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UnpausedIterator struct {
	Event *AutomationRegistryLogicC23Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23Unpaused)
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
		it.Event = new(AutomationRegistryLogicC23Unpaused)
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

func (it *AutomationRegistryLogicC23UnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicC23UnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UnpausedIterator{contract: _AutomationRegistryLogicC23.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23Unpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23Unpaused)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicC23Unpaused, error) {
	event := new(AutomationRegistryLogicC23Unpaused)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepAdminTransferRequested)
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

func (it *AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicC23UpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicC23UpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicC23UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepAdminTransferred)
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

func (it *AutomationRegistryLogicC23UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepAdminTransferred)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicC23UpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicC23UpkeepAdminTransferred)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicC23UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepCanceled)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepCanceled)
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

func (it *AutomationRegistryLogicC23UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicC23UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepCanceledIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepCanceled)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicC23UpkeepCanceled, error) {
	event := new(AutomationRegistryLogicC23UpkeepCanceled)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepChargedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepCharged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepChargedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepCharged)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepCharged)
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

func (it *AutomationRegistryLogicC23UpkeepChargedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepChargedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepCharged struct {
	Id      *big.Int
	Receipt AutomationRegistryBase23PaymentReceipt
	Raw     types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepChargedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepChargedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepCharged", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCharged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepCharged)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicC23UpkeepCharged, error) {
	event := new(AutomationRegistryLogicC23UpkeepCharged)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicC23UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepCheckDataSet)
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

func (it *AutomationRegistryLogicC23UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepCheckDataSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicC23UpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicC23UpkeepCheckDataSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicC23UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepGasLimitSet)
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

func (it *AutomationRegistryLogicC23UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepGasLimitSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicC23UpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicC23UpkeepGasLimitSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepMigrated)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepMigrated)
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

func (it *AutomationRegistryLogicC23UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepMigratedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepMigrated)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicC23UpkeepMigrated, error) {
	event := new(AutomationRegistryLogicC23UpkeepMigrated)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicC23UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepOffchainConfigSet)
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

func (it *AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicC23UpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepPausedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepPaused)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepPaused)
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

func (it *AutomationRegistryLogicC23UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepPausedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepPaused)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicC23UpkeepPaused, error) {
	event := new(AutomationRegistryLogicC23UpkeepPaused)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepPerformed)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepPerformed)
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

func (it *AutomationRegistryLogicC23UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicC23UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepPerformedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepPerformed)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicC23UpkeepPerformed, error) {
	event := new(AutomationRegistryLogicC23UpkeepPerformed)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicC23UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicC23UpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepReceived)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepReceived)
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

func (it *AutomationRegistryLogicC23UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepReceivedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepReceived)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicC23UpkeepReceived, error) {
	event := new(AutomationRegistryLogicC23UpkeepReceived)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicC23UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepRegistered)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepRegistered)
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

func (it *AutomationRegistryLogicC23UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepRegisteredIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepRegistered)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicC23UpkeepRegistered, error) {
	event := new(AutomationRegistryLogicC23UpkeepRegistered)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicC23UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepTriggerConfigSet)
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

func (it *AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicC23UpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicC23UpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicC23UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicC23UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicC23UpkeepUnpaused)
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
		it.Event = new(AutomationRegistryLogicC23UpkeepUnpaused)
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

func (it *AutomationRegistryLogicC23UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicC23UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicC23UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC23UpkeepUnpausedIterator{contract: _AutomationRegistryLogicC23.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC23.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicC23UpkeepUnpaused)
				if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23Filterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicC23UpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicC23UpkeepUnpaused)
	if err := _AutomationRegistryLogicC23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicC23.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicC23.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicC23.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistryLogicC23.ParseBillingConfigOverridden(log)
	case _AutomationRegistryLogicC23.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistryLogicC23.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistryLogicC23.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistryLogicC23.ParseBillingConfigSet(log)
	case _AutomationRegistryLogicC23.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicC23.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicC23.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicC23.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicC23.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicC23.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicC23.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistryLogicC23.ParseFeesWithdrawn(log)
	case _AutomationRegistryLogicC23.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicC23.ParseFundsAdded(log)
	case _AutomationRegistryLogicC23.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicC23.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicC23.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicC23.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicC23.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistryLogicC23.ParseNOPsSettledOffchain(log)
	case _AutomationRegistryLogicC23.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicC23.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicC23.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicC23.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicC23.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicC23.ParsePaused(log)
	case _AutomationRegistryLogicC23.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicC23.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicC23.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicC23.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicC23.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicC23.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicC23.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicC23.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicC23.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicC23.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicC23.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicC23.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicC23.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicC23.ParseUnpaused(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepCharged"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepCharged(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicC23.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicC23.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicC23AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicC23BillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistryLogicC23BillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistryLogicC23BillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba3")
}

func (AutomationRegistryLogicC23CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicC23ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicC23DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicC23FeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistryLogicC23FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicC23FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicC23InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicC23NOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistryLogicC23OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicC23OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicC23Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicC23PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicC23PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicC23PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicC23PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicC23ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicC23StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicC23Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicC23UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicC23UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicC23UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicC23UpkeepCharged) Topic() common.Hash {
	return common.HexToHash("0x801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc")
}

func (AutomationRegistryLogicC23UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicC23UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicC23UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicC23UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicC23UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicC23UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicC23UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicC23UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicC23UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicC23UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicC23UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicC23 *AutomationRegistryLogicC23) Address() common.Address {
	return _AutomationRegistryLogicC23.address
}

type AutomationRegistryLogicC23Interface interface {
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error)

	GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error)

	GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error)

	GetAvailableERC20ForPayment(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetBillingConfig(opts *bind.CallOpts, billingToken common.Address) (AutomationRegistryBase23BillingConfig, error)

	GetBillingOverrides(opts *bind.CallOpts, upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error)

	GetBillingOverridesEnabled(opts *bind.CallOpts, upkeepID *big.Int) (bool, error)

	GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error)

	GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error)

	GetChainModule(opts *bind.CallOpts) (common.Address, error)

	GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (AutomationRegistryBase23OnchainConfig, error)

	GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error)

	GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetNativeUSDFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetNumUpkeeps(opts *bind.CallOpts) (*big.Int, error)

	GetPayoutMode(opts *bind.CallOpts) (uint8, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error)

	GetReserveAmount(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error)

	GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

		error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetStorage(opts *bind.CallOpts) (AutomationRegistryBase23Storage, error)

	GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

		error)

	GetTransmittersWithPayees(opts *bind.CallOpts) ([]AutomationRegistryBase23TransmitterPayeeInfo, error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetWrappedNativeTokenAddress(opts *bind.CallOpts) (common.Address, error)

	HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicC23AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicC23AdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23BillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicC23BillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23BillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicC23BillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicC23BillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23BillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicC23BillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicC23CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicC23ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicC23ChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicC23DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicC23DedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicC23FeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicC23FeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicC23FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicC23FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicC23FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicC23InsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicC23NOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23NOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicC23NOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicC23OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicC23OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicC23PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicC23Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicC23PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicC23PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicC23PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicC23PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicC23PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicC23PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicC23ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicC23StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicC23UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicC23Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicC23UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicC23UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicC23UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicC23UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicC23UpkeepCanceled, error)

	FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepChargedIterator, error)

	WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCharged, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicC23UpkeepCharged, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicC23UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicC23UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicC23UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicC23UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicC23UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicC23UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicC23UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicC23UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicC23UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicC23UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicC23UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicC23UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
