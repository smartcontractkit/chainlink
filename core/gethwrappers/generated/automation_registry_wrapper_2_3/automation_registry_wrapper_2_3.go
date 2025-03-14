// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_registry_wrapper_2_3

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

var AutomationRegistry23MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"logicA\",\"type\":\"address\",\"internalType\":\"contractAutomationRegistryLogicA2_3\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fallbackTo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[],\"outputs\":[{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDigestAndEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"scanLogs\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"onchainConfigBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfigTypeSafe\",\"inputs\":[{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"components\":[{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCheckDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"maxPerformDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxRevertDataSize\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"financeAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackNativePrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registrars\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"chainModule\",\"type\":\"address\",\"internalType\":\"contractIChainModule\"}]},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"billingTokens\",\"type\":\"address[]\",\"internalType\":\"contractIERC20Metadata[]\"},{\"name\":\"billingConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig[]\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transmit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[3]\",\"internalType\":\"bytes32[3]\"},{\"name\":\"rawReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverridden\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"overrides\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverrideRemoved\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigSet\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainSpecificModuleUpdated\",\"inputs\":[{\"name\":\"newModule\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"configCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeesWithdrawn\",\"inputs\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NOPsSettledOffchain\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payments\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCharged\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receipt\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"components\":[{\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"linkUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"nativeUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLinkLiquidity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOffchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOnchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyFinanceAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x6101806040523480156200001257600080fd5b50604051620064b5380380620064b583398101604081905262000035916200062f565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b91906200062f565b826001600160a01b031663226cf83c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200062f565b836001600160a01b031663614486af6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200062f565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200062f565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200062f565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200062f565b876001600160a01b031663c5b964e06040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002d3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002f9919062000656565b886001600160a01b031663ac4dc59a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000338573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200035e91906200062f565b3380600081620003b55760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620003e857620003e8816200056b565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff19166001838181111562000439576200043962000679565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200049a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004c091906200068f565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000504573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200052a91906200068f565b60ff16146200054c576040516301f86e1760e41b815260040160405180910390fd5b5050506001600160a01b039095166101605250620006b4945050505050565b336001600160a01b03821603620005c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620003ac565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200062c57600080fd5b50565b6000602082840312156200064257600080fd5b81516200064f8162000616565b9392505050565b6000602082840312156200066957600080fd5b8151600281106200064f57600080fd5b634e487b7160e01b600052602160045260246000fd5b600060208284031215620006a257600080fd5b815160ff811681146200064f57600080fd5b60805160a05160c05160e05161010051610120516101405161016051615d9d620007186000396000818160b301526101860152600050506000505060005050600050506000613730015260005050600081816113710152612d590152615d9d6000f3fe6080604052600436106100b15760003560e01c80638da5cb5b11610069578063b1dc65a41161004e578063b1dc65a4146102cb578063e3d0e712146102eb578063f2fde38b1461030b576100b1565b80638da5cb5b1461025a578063afcb95d714610285576100b1565b8063349e8cca1161009a578063349e8cca1461017757806379ba5097146101cb57806381ff7048146101e0576100b1565b80630870d3a1146100f8578063181f5a7714610118575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156100f1573d6000f35b3d6000fd5b005b34801561010457600080fd5b506100f6610113366004614b11565b61032b565b34801561012457600080fd5b506101616040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e332e30000000000000000081525081565b60405161016e9190614c8d565b60405180910390f35b34801561018357600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161016e565b3480156101d757600080fd5b506100f6610c0f565b3480156101ec57600080fd5b5061023760175460135463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600161016e565b34801561026657600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166101a6565b34801561029157600080fd5b50601354601454604080516000815260208101939093526c0100000000000000000000000090910463ffffffff169082015260600161016e565b3480156102d757600080fd5b506100f66102e6366004614cec565b610d11565b3480156102f757600080fd5b506100f6610306366004614dd1565b611051565b34801561031757600080fd5b506100f6610326366004614e9e565b61108b565b61033361109f565b601f8851111561036f576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8560ff166000036103ac576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b865188511415806103cb57506103c3866003614eea565b60ff16885111155b15610402576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805182511461043d576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104478282611122565b610451888861175a565b604051806101200160405280601460000160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001600063ffffffff1681526020018660a0015162ffffff16815260200186610120015161ffff1681526020018760ff168152602001601460000160169054906101000a900460ff1615158152602001601460000160179054906101000a900460ff1615158152602001866080015115158152602001866101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548162ffffff021916908362ffffff16021790555060608201518160000160136101000a81548161ffff021916908361ffff16021790555060808201518160000160156101000a81548160ff021916908360ff16021790555060a08201518160000160166101000a81548160ff02191690831515021790555060c08201518160000160176101000a81548160ff02191690831515021790555060e08201518160000160186101000a81548160ff0219169083151502179055506101008201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506000601660010160189054906101000a900463ffffffff1690506000866101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610701573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107259190614f06565b6017549091506000906107579074010000000000000000000000000000000000000000900463ffffffff166001614f1f565b9050604051806101600160405280896060015173ffffffffffffffffffffffffffffffffffffffff168152602001896000015163ffffffff168152602001896020015163ffffffff1681526020016016600001601c9054906101000a900463ffffffff1663ffffffff16815260200189610100015173ffffffffffffffffffffffffffffffffffffffff1681526020018263ffffffff1681526020018363ffffffff168152602001896040015163ffffffff16815260200189610140015173ffffffffffffffffffffffffffffffffffffffff1681526020018960c0015163ffffffff1681526020018960e0015163ffffffff16815250601660008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160186101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101208201518160020160146101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160186101000a81548163ffffffff021916908363ffffffff160217905550905050876101600151601981905550876101800151601a81905550876101a00151601b81905550600088604051602001610a9a9190614f8d565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152919052601754909150610aff904690309074010000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f611e19565b6013556000610b0e6009611ec3565b90505b8015610b4b57610b38610b30610b28600184615113565b600990611ed3565b600990611ee6565b5080610b4381615126565b915050610b11565b5060005b896101c0015151811015610ba257610b8f8a6101c001518281518110610b7757610b7761515b565b60200260200101516009611f0890919063ffffffff16565b5080610b9a8161518a565b915050610b4f565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0584601354601660010160149054906101000a900463ffffffff168f8f8f878f8f604051610bf9999897969594939291906151c2565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610c95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005a90506000610d23866040615258565b610d2f8961014461526f565b610d39919061526f565b9050368114610d74576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff75010000000000000000000000000000000000000000008304811660808301527601000000000000000000000000000000000000000000008304811615801560a08401527701000000000000000000000000000000000000000000000084048216151560c0840152780100000000000000000000000000000000000000000000000090930416151560e082015260155473ffffffffffffffffffffffffffffffffffffffff1661010082015290610ed3576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff16610f1c576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013548b3514610f58576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6080810151610f68906001615282565b60ff1687141580610f795750868514155b15610fb0576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610fc08b8b8b8b8b8b8b8b611f2a565b6000610fcc8b8b612193565b905060208c0135600881901c63ffffffff16610fe984848861224c565b836020015163ffffffff168163ffffffff16111561104157601480547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c0100000000000000000000000063ffffffff8416021790555b5050505050505050505050505050565b60008060008580602001905181019061106a9190615408565b925092509250611080898989868989888861032b565b505050505050505050565b61109361109f565b61109c81612dcf565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611120576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610c8c565b565b60005b6024548110156111e05760226000602483815481106111465761114661515b565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffff00000000000000000000000000000000000000000000000000000000168155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000169055806111d88161518a565b915050611125565b506111ed60246000614565565b60255460ff1660005b83518110156117545760008482815181106112135761121361515b565b6020026020010151905060008483815181106112315761123161515b565b602002602001015190508173ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015611286573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112aa91906155b3565b60ff16816060015160ff161415806113385750806040015173ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801561130c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061133091906155b3565b60ff16600814155b1561136f576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161480156113db575060018460018111156113d9576113d96155d0565b145b15611412576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216158061144d5750604081015173ffffffffffffffffffffffffffffffffffffffff16155b15611484576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff828116600090815260226020526040902054670100000000000000900416156114ee576040517f357d0cc400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6024805460018181019092557f7cd332d19b93bcabe3cce7ca0c18a052f57e5fd03b4758a09f30f5ddc4b22ec401805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff0000000000000000000000000000000000000000909216821790925560008181526022602090815260409182902086518154928801518489015160608a015160ff167b01000000000000000000000000000000000000000000000000000000027fffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffff9190981667010000000000000002167fffffffff000000000000000000000000000000000000000000ffffffffffffff62ffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000090951663ffffffff9093169290921793909317929092169190911793909317835560808501519383019390935560a0840151600290920180546bffffffffffffffffffffffff9093167fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009093169290921790915590517fca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba390611737908490600060c08201905063ffffffff835116825262ffffff602084015116602083015273ffffffffffffffffffffffffffffffffffffffff604084015116604083015260ff6060840151166060830152608083015160808301526bffffffffffffffffffffffff60a08401511660a083015292915050565b60405180910390a25050808061174c9061518a565b9150506111f6565b50505050565b600e546014546bffffffffffffffffffffffff1660005b600e548110156117cd576117ba600e82815481106117915761179161515b565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168385612ec4565b50806117c58161518a565b915050611771565b5060255460009060ff16815b600e5481101561193e57600e81815481106117f6576117f661515b565b6000918252602082200154600d805473ffffffffffffffffffffffffffffffffffffffff9092169550600c9291849081106118335761183361515b565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168452838201949094526040928301822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559286168152600b909252902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905560018260018111156118d5576118d56155d0565b14801561191a575073ffffffffffffffffffffffffffffffffffffffff83166000908152600b60205260409020546201000090046bffffffffffffffffffffffff1615155b1561192c5761192a600f84611f08565b505b806119368161518a565b9150506117d9565b5061194b600d6000614565565b611957600e6000614565565b6040805160808101825260008082526020820181905291810182905260608101829052905b8751811015611de757600c600089838151811061199b5761199b61515b565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615611a06576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16888281518110611a3057611a3061515b565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611a85576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008a8481518110611ab657611ab661515b565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558651879082908110611b5e57611b5e61515b565b60200260200101519350600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603611bce576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e01000000000000000000000000000090049092166060830152909250611c89576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180835260ff80831660208086019182526014546bffffffffffffffffffffffff9081166060880190815273ffffffffffffffffffffffffffffffffffffffff8a166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090951694909417179190911692909217919091179055836001811115611dc357611dc36155d0565b03611dd557611dd3600f85611ee6565b505b80611ddf8161518a565b91505061197c565b508651611dfb90600d9060208a0190614583565b508551611e0f90600e906020890190614583565b5050505050505050565b6000808a8a8a8a8a8a8a8a8a604051602001611e3d999897969594939291906155ff565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b6000611ecd825490565b92915050565b6000611edf83836130cc565b9392505050565b6000611edf8373ffffffffffffffffffffffffffffffffffffffff84166130f6565b6000611edf8373ffffffffffffffffffffffffffffffffffffffff84166131f0565b60008787604051611f3c929190615694565b604051908190038120611f53918b906020016156a4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b8881101561212a57600185878360208110611fbf57611fbf61515b565b611fcc91901a601b615282565b8c8c85818110611fde57611fde61515b565b905060200201358b8b86818110611ff757611ff761515b565b9050602002013560405160008152602001604052604051612034949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612056573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955093509050612104576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b8401935080806121229061518a565b915050611fa2565b50827e01010101010101010101010101010101010101010101010101010101010101841614612185576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b6121cc6040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b60006121da83850185615795565b60408101515160608201515191925090811415806121fd57508082608001515114155b8061220d5750808260a001515114155b15612244576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b509392505050565b600082604001515167ffffffffffffffff81111561226c5761226c61461a565b60405190808252806020026020018201604052801561233857816020015b6040805161020081018252600060e08201818152610100830182905261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c083018290526101e0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191018161228a5790505b50905060006040518060800160405280600061ffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160008152509050600085610100015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156123d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123fa9190614f06565b6101008701516040517f7810d12a00000000000000000000000000000000000000000000000000000000815236600482015291925060009173ffffffffffffffffffffffffffffffffffffffff90911690637810d12a90602401602060405180830381865afa158015612471573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124959190614f06565b905060005b8660400151518110156129255760046000886040015183815181106124c1576124c161515b565b6020908102919091018101518252818101929092526040908101600020815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490931660e08201526002909101549091169181019190915285518690839081106125e4576125e461515b565b6020026020010151600001819052506126198760400151828151811061260c5761260c61515b565b602002602001015161323f565b85828151811061262b5761262b61515b565b6020026020010151606001906001811115612648576126486155d0565b9081600181111561265b5761265b6155d0565b815250506126bf876040015182815181106126785761267861515b565b602002602001015184896080015184815181106126975761269761515b565b60200260200101518885815181106126b1576126b161515b565b60200260200101518c6132ea565b8683815181106126d1576126d161515b565b60200260200101516020018784815181106126ee576126ee61515b565b602002602001015160c00182815250821515151581525050508481815181106127195761271961515b565b602002602001015160200151156127495760018460000181815161273d9190615882565b61ffff1690525061274e565b612913565b6127b48582815181106127635761276361515b565b60200260200101516000015160800151886060015183815181106127895761278961515b565b60200260200101518960a0015184815181106127a7576127a761515b565b6020026020010151613409565b8683815181106127c6576127c661515b565b60200260200101516040018784815181106127e3576127e361515b565b602002602001015160800182815250821515151581525050508760800151600161280d9190615282565b61281b9060ff166040615258565b6103a48860a0015183815181106128345761283461515b565b602002602001015151612847919061526f565b612851919061526f565b8582815181106128635761286361515b565b602002602001015160a00181815250508481815181106128855761288561515b565b602002602001015160a00151846060018181516128a2919061526f565b90525084518590829081106128b9576128b961515b565b602002602001015160800151866128d09190615113565b9550612913876040015182815181106128eb576128eb61515b565b6020026020010151848784815181106129065761290661515b565b6020026020010151613624565b8061291d8161518a565b91505061249a565b50825161ffff1660000361293c5750505050505050565b61c73861294a366010615258565b5a6129559088615113565b61295f919061526f565b612969919061526f565b8351909550613a34906129809061ffff16876158cc565b61298a919061526f565b60408051608081018252600080825260208201819052918101829052606081018290529196506129b989613729565b905060005b886040015151811015612c68578681815181106129dd576129dd61515b565b60200260200101516020015115612c5657612a1a8a888381518110612a0457612a0461515b565b6020026020010151600001516101000151613813565b92506000612b3a8b6040518061012001604052808b8681518110612a4057612a4061515b565b60200260200101516080015181526020018c81526020018a606001518c8781518110612a6e57612a6e61515b565b602002602001015160a001518a612a859190615258565b612a8f91906158cc565b81526020018d6000015181526020018d6020015181526020018681526020018b8681518110612ac057612ac061515b565b602002602001015160000151610100015173ffffffffffffffffffffffffffffffffffffffff168152602001878152602001600115158152508c604001518581518110612b0f57612b0f61515b565b60200260200101518b8681518110612b2957612b2961515b565b60200260200101516000015161398f565b9050806060015187604001818151612b5291906158e0565b6bffffffffffffffffffffffff169052506040810151602088018051612b799083906158e0565b6bffffffffffffffffffffffff169052508751889083908110612b9e57612b9e61515b565b60200260200101516040015115158a604001518381518110612bc257612bc261515b565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b83602001518460000151612bff91906158e0565b8b8681518110612c1157612c1161515b565b6020026020010151608001518d8f608001518881518110612c3457612c3461515b565b6020026020010151604051612c4c9493929190615905565b60405180910390a3505b80612c608161518a565b9150506129be565b505050602083810151336000908152600b90925260409091208054600290612ca59084906201000090046bffffffffffffffffffffffff166158e0565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508260400151601460000160008282829054906101000a90046bffffffffffffffffffffffff16612d0391906158e0565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555082604001518360200151612d4591906158e0565b6bffffffffffffffffffffffff16602160007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254612dc1919061526f565b909155505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612e4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610c8c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906130c0576000816060015185612f5c9190615942565b90506000612f6a8583615967565b90508083604001818151612f7e91906158e0565b6bffffffffffffffffffffffff16905250612f998582615992565b83606001818151612faa91906158e0565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b60008260000182815481106130e3576130e361515b565b9060005260206000200154905092915050565b600081815260018301602052604081205480156131df57600061311a600183615113565b855490915060009061312e90600190615113565b905081811461319357600086600001828154811061314e5761314e61515b565b90600052602060002001549050808760000184815481106131715761317161515b565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806131a4576131a46159c2565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611ecd565b6000915050611ecd565b5092915050565b600081815260018301602052604081205461323757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611ecd565b506000611ecd565b6000818160045b600f8110156132cc577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106132845761328461515b565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146132ba57506000949350505050565b806132c48161518a565b915050613246565b5081600f1a60018111156132e2576132e26155d0565b949350505050565b600080808085606001516001811115613305576133056155d0565b0361332b576133178888888888613dfe565b613326576000925090506133ff565b6133a3565b600185606001516001811115613343576133436155d0565b0361337157600061335689898988613f88565b925090508061336b57506000925090506133ff565b506133a3565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516060015163ffffffff1687106133f857877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636876040516133e59190614c8d565b60405180910390a26000925090506133ff565b6001925090505b9550959350505050565b601454600090819077010000000000000000000000000000000000000000000000900460ff1615613466576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16770100000000000000000000000000000000000000000000001790556040517f4585e33b00000000000000000000000000000000000000000000000000000000906134db908590602401614c8d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d16906135ae90879087906004016159f1565b60408051808303816000875af11580156135cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135f09190615a0a565b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008160600151600181111561363c5761363c6155d0565b036136a057600083815260046020526040902060010180547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff851602179055505050565b6001816060015160018111156136b8576136b86155d0565b036137245760c08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613799573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137bd9190615a52565b509350509250506000821315806137d357508042105b8061380357506000846040015162ffffff1611801561380357506137f78142615113565b846040015162ffffff16105b156131e9575050601b5492915050565b60408051608081018252600080825260208083018281528385018381526060850184905273ffffffffffffffffffffffffffffffffffffffff878116855260229093528584208054640100000000810462ffffff1690925263ffffffff82169092527b01000000000000000000000000000000000000000000000000000000810460ff16855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495919484936701000000000000009092049091169163feaf968c9160048083019260a09291908290030181865afa158015613900573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139249190615a52565b5093505092505060008213158061393a57508042105b8061396a57506000866040015162ffffff1611801561396a575061395e8142615113565b866040015162ffffff16105b1561397e5760018301546060850152613986565b606084018290525b50505092915050565b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e08201529082015115613a255760008381526023602090815260409182902082518084018452905463ffffffff811680835262ffffff640100000000909204821692840192835260e089018051909401529051915191169101525b6000613a318686614195565b60c0840151602082015182519293509091600091613a4e916158e0565b60e08801515190915060ff16600060128210613a6b576001613a81565b613a76826012615113565b613a8190600a615bc2565b9050600060128311613a94576001613aaa565b613a9f601284615113565b613aaa90600a615bc2565b905085600001516bffffffffffffffffffffffff16856bffffffffffffffffffffffff161015613b5257849350613b26818b60800151613aea9190615258565b838c60e0015160600151886bffffffffffffffffffffffff16613b0d9190615258565b613b179190615258565b613b2191906158cc565b6144c3565b6bffffffffffffffffffffffff9081166040880152600060608801819052602088015285168652613c78565b836bffffffffffffffffffffffff16856bffffffffffffffffffffffff161015613c7857849350613be086604001516bffffffffffffffffffffffff16828c60800151613b9f9190615258565b848d60e0015160600151896bffffffffffffffffffffffff16613bc29190615258565b613bcc9190615258565b613bd691906158cc565b613b219190615113565b6bffffffffffffffffffffffff1660608088019190915260e08b01510151613c6490613c0d908490615258565b6001848d60e0015160600151613c239190615258565b613c2d9190615113565b838d608001518a606001516bffffffffffffffffffffffff16613c509190615258565b613c5a9190615258565b613b17919061526f565b6bffffffffffffffffffffffff1660208701525b60008981526004602052604090206001018054859190601090613cbe90849070010000000000000000000000000000000090046bffffffffffffffffffffffff16615942565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b81526004602052604081206001018054928816935091613d199084906fffffffffffffffffffffffffffffffff16615bce565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550836bffffffffffffffffffffffff16602160008c60c0015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254613db09190615113565b92505081905550887f801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc87604051613de79190615bf7565b60405180910390a250939998505050505050505050565b60008084806020019051810190613e159190615cb7565b845160e00151815191925063ffffffff90811691161015613e7257867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e886604051613e609190614c8d565b60405180910390a26000915050613f7f565b8260e001518015613f325750602081015115801590613f325750602081015161010084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa158015613f0b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f2f9190614f06565b14155b80613f445750805163ffffffff168611155b15613f7957867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc30186604051613e609190614c8d565b60019150505b95945050505050565b600080600084806020019051810190613fa19190615d0f565b905060008782600001518360200151846040015160405160200161400394939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b6040516020818303038152906040528051906020012090508460e0015180156140de57506080820151158015906140de5750608082015161010086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156140b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140db9190614f06565b14155b806140f3575086826060015163ffffffff1610155b1561413d57877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301876040516141289190614c8d565b60405180910390a260009350915061418c9050565b60008181526008602052604090205460ff161561418457877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8876040516141289190614c8d565b600193509150505b94509492505050565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081019190915260008260e001516000015160ff1690506000846060015161ffff1684606001516142009190615258565b905083610100015180156142135750803a105b1561421b57503a5b60006012831161422c576001614242565b614237601284615113565b61424290600a615bc2565b905060006012841061425557600161426b565b614260846012615113565b61426b90600a615bc2565b905060008660a0015187604001518860200151896000015161428d919061526f565b6142979087615258565b6142a1919061526f565b6142ab9190615258565b90506142ee828860e00151606001516142c49190615258565b6001848a60e00151606001516142da9190615258565b6142e49190615113565b613c5a8685615258565b6bffffffffffffffffffffffff168652608087015161431190613b2190836158cc565b6bffffffffffffffffffffffff1660408088019190915260e0880151015160009061434a9062ffffff16683635c9adc5dea00000615258565b9050600081633b9aca008a60a001518b60e001516020015163ffffffff168c604001518d600001518b61437d9190615258565b614387919061526f565b6143919190615258565b61439b9190615258565b6143a591906158cc565b6143af919061526f565b90506143f2848a60e00151606001516143c89190615258565b6001868c60e00151606001516143de9190615258565b6143e89190615113565b613c5a8885615258565b6bffffffffffffffffffffffff166020890152608089015161441890613b2190836158cc565b6bffffffffffffffffffffffff16606089015260c089015173ffffffffffffffffffffffffffffffffffffffff166080808a019190915289015161445b906144c3565b6bffffffffffffffffffffffff1660a0808a019190915289015161447e906144c3565b6bffffffffffffffffffffffff1660c089015260e0890151606001516144a3906144c3565b6bffffffffffffffffffffffff1660e08901525050505050505092915050565b60006bffffffffffffffffffffffff821115614561576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610c8c565b5090565b508054600082559060005260206000209081019061109c9190614605565b8280548282559060005260206000209081019282156145fd579160200282015b828111156145fd57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906145a3565b506145619291505b5b808211156145615760008155600101614606565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610200810167ffffffffffffffff8111828210171561466d5761466d61461a565b60405290565b60405160c0810167ffffffffffffffff8111828210171561466d5761466d61461a565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156146dd576146dd61461a565b604052919050565b600067ffffffffffffffff8211156146ff576146ff61461a565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff8116811461109c57600080fd5b803561473681614709565b919050565b600082601f83011261474c57600080fd5b8135602061476161475c836146e5565b614696565b82815260059290921b8401810191818101908684111561478057600080fd5b8286015b848110156147a457803561479781614709565b8352918301918301614784565b509695505050505050565b60ff8116811461109c57600080fd5b8035614736816147af565b63ffffffff8116811461109c57600080fd5b8035614736816147c9565b801515811461109c57600080fd5b8035614736816147e6565b62ffffff8116811461109c57600080fd5b8035614736816147ff565b61ffff8116811461109c57600080fd5b80356147368161481b565b6000610200828403121561484957600080fd5b614851614649565b905061485c826147db565b815261486a602083016147db565b602082015261487b604083016147db565b604082015261488c6060830161472b565b606082015261489d608083016147f4565b60808201526148ae60a08301614810565b60a08201526148bf60c083016147db565b60c08201526148d060e083016147db565b60e08201526101006148e381840161472b565b908201526101206148f583820161482b565b9082015261014061490783820161472b565b90820152610160828101359082015261018080830135908201526101a080830135908201526101c08083013567ffffffffffffffff81111561494857600080fd5b6149548582860161473b565b8284015250506101e061496881840161472b565b9082015292915050565b803567ffffffffffffffff8116811461473657600080fd5b600082601f83011261499b57600080fd5b813567ffffffffffffffff8111156149b5576149b561461a565b6149e660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614696565b8181528460208386010111156149fb57600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff8116811461109c57600080fd5b600082601f830112614a4357600080fd5b81356020614a5361475c836146e5565b82815260c09283028501820192828201919087851115614a7257600080fd5b8387015b85811015614b045781818a031215614a8e5760008081fd5b614a96614673565b8135614aa1816147c9565b815281860135614ab0816147ff565b81870152604082810135614ac381614709565b90820152606082810135614ad6816147af565b908201526080828101359082015260a080830135614af381614a18565b908201528452928401928101614a76565b5090979650505050505050565b600080600080600080600080610100898b031215614b2e57600080fd5b883567ffffffffffffffff80821115614b4657600080fd5b614b528c838d0161473b565b995060208b0135915080821115614b6857600080fd5b614b748c838d0161473b565b9850614b8260408c016147be565b975060608b0135915080821115614b9857600080fd5b614ba48c838d01614836565b9650614bb260808c01614972565b955060a08b0135915080821115614bc857600080fd5b614bd48c838d0161498a565b945060c08b0135915080821115614bea57600080fd5b614bf68c838d0161473b565b935060e08b0135915080821115614c0c57600080fd5b50614c198b828c01614a32565b9150509295985092959890939650565b6000815180845260005b81811015614c4f57602081850181015186830182015201614c33565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611edf6020830184614c29565b60008083601f840112614cb257600080fd5b50813567ffffffffffffffff811115614cca57600080fd5b6020830191508360208260051b8501011115614ce557600080fd5b9250929050565b60008060008060008060008060e0898b031215614d0857600080fd5b606089018a811115614d1957600080fd5b8998503567ffffffffffffffff80821115614d3357600080fd5b818b0191508b601f830112614d4757600080fd5b813581811115614d5657600080fd5b8c6020828501011115614d6857600080fd5b6020830199508098505060808b0135915080821115614d8657600080fd5b614d928c838d01614ca0565b909750955060a08b0135915080821115614dab57600080fd5b50614db88b828c01614ca0565b999c989b50969995989497949560c00135949350505050565b60008060008060008060c08789031215614dea57600080fd5b863567ffffffffffffffff80821115614e0257600080fd5b614e0e8a838b0161473b565b97506020890135915080821115614e2457600080fd5b614e308a838b0161473b565b9650614e3e60408a016147be565b95506060890135915080821115614e5457600080fd5b614e608a838b0161498a565b9450614e6e60808a01614972565b935060a0890135915080821115614e8457600080fd5b50614e9189828a0161498a565b9150509295509295509295565b600060208284031215614eb057600080fd5b8135611edf81614709565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff81811683821602908116908181146131e9576131e9614ebb565b600060208284031215614f1857600080fd5b5051919050565b63ffffffff8181168382160190808211156131e9576131e9614ebb565b600081518084526020808501945080840160005b83811015614f8257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614f50565b509495945050505050565b60208152614fa460208201835163ffffffff169052565b60006020830151614fbd604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015173ffffffffffffffffffffffffffffffffffffffff8116608084015250608083015180151560a08401525060a083015162ffffff811660c08401525060c083015163ffffffff811660e08401525060e083015161010061503c8185018363ffffffff169052565b84015190506101206150658482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151905061014061507c8482018361ffff169052565b84015190506101606150a58482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151610180848101919091528401516101a0808501919091528401516101c0808501919091528401516102006101e0808601829052919250906150ed610220860184614f3c565b95015173ffffffffffffffffffffffffffffffffffffffff169301929092525090919050565b81810381811115611ecd57611ecd614ebb565b60008161513557615135614ebb565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036151bb576151bb614ebb565b5060010190565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526151f28184018a614f3c565b905082810360808401526152068189614f3c565b905060ff871660a084015282810360c08401526152238187614c29565b905067ffffffffffffffff851660e08401528281036101008401526152488185614c29565b9c9b505050505050505050505050565b8082028115828204841417611ecd57611ecd614ebb565b80820180821115611ecd57611ecd614ebb565b60ff8181168382160190811115611ecd57611ecd614ebb565b8051614736816147c9565b805161473681614709565b8051614736816147e6565b8051614736816147ff565b80516147368161481b565b600082601f8301126152e357600080fd5b815160206152f361475c836146e5565b82815260059290921b8401810191818101908684111561531257600080fd5b8286015b848110156147a457805161532981614709565b8352918301918301615316565b600082601f83011261534757600080fd5b8151602061535761475c836146e5565b82815260c0928302850182019282820191908785111561537657600080fd5b8387015b85811015614b045781818a0312156153925760008081fd5b61539a614673565b81516153a5816147c9565b8152818601516153b4816147ff565b818701526040828101516153c781614709565b908201526060828101516153da816147af565b908201526080828101519082015260a0808301516153f781614a18565b90820152845292840192810161537a565b60008060006060848603121561541d57600080fd5b835167ffffffffffffffff8082111561543557600080fd5b90850190610200828803121561544a57600080fd5b615452614649565b61545b8361529b565b81526154696020840161529b565b602082015261547a6040840161529b565b604082015261548b606084016152a6565b606082015261549c608084016152b1565b60808201526154ad60a084016152bc565b60a08201526154be60c0840161529b565b60c08201526154cf60e0840161529b565b60e08201526101006154e28185016152a6565b908201526101206154f48482016152c7565b908201526101406155068482016152a6565b90820152610160838101519082015261018080840151908201526101a080840151908201526101c0808401518381111561553f57600080fd5b61554b8a8287016152d2565b8284015250506101e061555f8185016152a6565b90820152602087015190955091508082111561557a57600080fd5b615586878388016152d2565b9350604086015191508082111561559c57600080fd5b506155a986828701615336565b9150509250925092565b6000602082840312156155c557600080fd5b8151611edf816147af565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526156468285018b614f3c565b9150838203608085015261565a828a614f3c565b915060ff881660a085015283820360c08501526156778288614c29565b90861660e085015283810361010085015290506152488185614c29565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f8301126156cb57600080fd5b813560206156db61475c836146e5565b82815260059290921b840181019181810190868411156156fa57600080fd5b8286015b848110156147a457803583529183019183016156fe565b600082601f83011261572657600080fd5b8135602061573661475c836146e5565b82815260059290921b8401810191818101908684111561575557600080fd5b8286015b848110156147a457803567ffffffffffffffff8111156157795760008081fd5b6157878986838b010161498a565b845250918301918301615759565b6000602082840312156157a757600080fd5b813567ffffffffffffffff808211156157bf57600080fd5b9083019060c082860312156157d357600080fd5b6157db614673565b82358152602083013560208201526040830135828111156157fb57600080fd5b615807878286016156ba565b60408301525060608301358281111561581f57600080fd5b61582b878286016156ba565b60608301525060808301358281111561584357600080fd5b61584f87828601615715565b60808301525060a08301358281111561586757600080fd5b61587387828601615715565b60a08301525095945050505050565b61ffff8181168382160190808211156131e9576131e9614ebb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826158db576158db61589d565b500490565b6bffffffffffffffffffffffff8181168382160190808211156131e9576131e9614ebb565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006159386080830184614c29565b9695505050505050565b6bffffffffffffffffffffffff8281168282160390808211156131e9576131e9614ebb565b60006bffffffffffffffffffffffff808416806159865761598661589d565b92169190910492915050565b6bffffffffffffffffffffffff8181168382160280821691908281146159ba576159ba614ebb565b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b8281526040602082015260006132e26040830184614c29565b60008060408385031215615a1d57600080fd5b8251615a28816147e6565b6020939093015192949293505050565b805169ffffffffffffffffffff8116811461473657600080fd5b600080600080600060a08688031215615a6a57600080fd5b615a7386615a38565b9450602086015193506040860151925060608601519150615a9660808701615a38565b90509295509295909350565b600181815b80851115615afb57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115615ae157615ae1614ebb565b80851615615aee57918102915b93841c9390800290615aa7565b509250929050565b600082615b1257506001611ecd565b81615b1f57506000611ecd565b8160018114615b355760028114615b3f57615b5b565b6001915050611ecd565b60ff841115615b5057615b50614ebb565b50506001821b611ecd565b5060208310610133831016604e8410600b8410161715615b7e575081810a611ecd565b615b888383615aa2565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115615bba57615bba614ebb565b029392505050565b6000611edf8383615b03565b6fffffffffffffffffffffffffffffffff8181168382160190808211156131e9576131e9614ebb565b6000610100820190506bffffffffffffffffffffffff8084511683528060208501511660208401528060408501511660408401528060608501511660608401525073ffffffffffffffffffffffffffffffffffffffff608084015116608083015260a0830151615c7760a08401826bffffffffffffffffffffffff169052565b5060c0830151615c9760c08401826bffffffffffffffffffffffff169052565b5060e08301516131e960e08401826bffffffffffffffffffffffff169052565b600060408284031215615cc957600080fd5b6040516040810181811067ffffffffffffffff82111715615cec57615cec61461a565b6040528251615cfa816147c9565b81526020928301519281019290925250919050565b600060a08284031215615d2157600080fd5b60405160a0810181811067ffffffffffffffff82111715615d4457615d4461461a565b806040525082518152602083015160208201526040830151615d65816147c9565b60408201526060830151615d78816147c9565b6060820152608092830151928101929092525091905056fea164736f6c6343000813000a",
}

var AutomationRegistry23ABI = AutomationRegistry23MetaData.ABI

var AutomationRegistry23Bin = AutomationRegistry23MetaData.Bin

func DeployAutomationRegistry23(auth *bind.TransactOpts, backend bind.ContractBackend, logicA common.Address) (common.Address, *types.Transaction, *AutomationRegistry23, error) {
	parsed, err := AutomationRegistry23MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistry23Bin), backend, logicA)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistry23{address: address, abi: *parsed, AutomationRegistry23Caller: AutomationRegistry23Caller{contract: contract}, AutomationRegistry23Transactor: AutomationRegistry23Transactor{contract: contract}, AutomationRegistry23Filterer: AutomationRegistry23Filterer{contract: contract}}, nil
}

type AutomationRegistry23 struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistry23Caller
	AutomationRegistry23Transactor
	AutomationRegistry23Filterer
}

type AutomationRegistry23Caller struct {
	contract *bind.BoundContract
}

type AutomationRegistry23Transactor struct {
	contract *bind.BoundContract
}

type AutomationRegistry23Filterer struct {
	contract *bind.BoundContract
}

type AutomationRegistry23Session struct {
	Contract     *AutomationRegistry23
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistry23CallerSession struct {
	Contract *AutomationRegistry23Caller
	CallOpts bind.CallOpts
}

type AutomationRegistry23TransactorSession struct {
	Contract     *AutomationRegistry23Transactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistry23Raw struct {
	Contract *AutomationRegistry23
}

type AutomationRegistry23CallerRaw struct {
	Contract *AutomationRegistry23Caller
}

type AutomationRegistry23TransactorRaw struct {
	Contract *AutomationRegistry23Transactor
}

func NewAutomationRegistry23(address common.Address, backend bind.ContractBackend) (*AutomationRegistry23, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistry23ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistry23(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23{address: address, abi: abi, AutomationRegistry23Caller: AutomationRegistry23Caller{contract: contract}, AutomationRegistry23Transactor: AutomationRegistry23Transactor{contract: contract}, AutomationRegistry23Filterer: AutomationRegistry23Filterer{contract: contract}}, nil
}

func NewAutomationRegistry23Caller(address common.Address, caller bind.ContractCaller) (*AutomationRegistry23Caller, error) {
	contract, err := bindAutomationRegistry23(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23Caller{contract: contract}, nil
}

func NewAutomationRegistry23Transactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistry23Transactor, error) {
	contract, err := bindAutomationRegistry23(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23Transactor{contract: contract}, nil
}

func NewAutomationRegistry23Filterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistry23Filterer, error) {
	contract, err := bindAutomationRegistry23(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23Filterer{contract: contract}, nil
}

func bindAutomationRegistry23(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistry23MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistry23 *AutomationRegistry23Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry23.Contract.AutomationRegistry23Caller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry23 *AutomationRegistry23Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.AutomationRegistry23Transactor.contract.Transfer(opts)
}

func (_AutomationRegistry23 *AutomationRegistry23Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.AutomationRegistry23Transactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry23.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.contract.Transfer(opts)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry23 *AutomationRegistry23Caller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry23.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry23 *AutomationRegistry23Session) FallbackTo() (common.Address, error) {
	return _AutomationRegistry23.Contract.FallbackTo(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistry23.Contract.FallbackTo(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Caller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _AutomationRegistry23.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_AutomationRegistry23 *AutomationRegistry23Session) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry23.Contract.LatestConfigDetails(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry23.Contract.LatestConfigDetails(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _AutomationRegistry23.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_AutomationRegistry23 *AutomationRegistry23Session) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry23.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry23.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry23.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry23 *AutomationRegistry23Session) Owner() (common.Address, error) {
	return _AutomationRegistry23.Contract.Owner(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerSession) Owner() (common.Address, error) {
	return _AutomationRegistry23.Contract.Owner(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AutomationRegistry23.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AutomationRegistry23 *AutomationRegistry23Session) TypeAndVersion() (string, error) {
	return _AutomationRegistry23.Contract.TypeAndVersion(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23CallerSession) TypeAndVersion() (string, error) {
	return _AutomationRegistry23.Contract.TypeAndVersion(&_AutomationRegistry23.CallOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistry23 *AutomationRegistry23Session) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.AcceptOwnership(&_AutomationRegistry23.TransactOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.AcceptOwnership(&_AutomationRegistry23.TransactOpts)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry23 *AutomationRegistry23Session) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.SetConfig(&_AutomationRegistry23.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.SetConfig(&_AutomationRegistry23.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry23 *AutomationRegistry23Session) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.SetConfigTypeSafe(&_AutomationRegistry23.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.SetConfigTypeSafe(&_AutomationRegistry23.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistry23 *AutomationRegistry23Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.TransferOwnership(&_AutomationRegistry23.TransactOpts, to)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.TransferOwnership(&_AutomationRegistry23.TransactOpts, to)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry23 *AutomationRegistry23Session) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.Transmit(&_AutomationRegistry23.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.Transmit(&_AutomationRegistry23.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry23 *AutomationRegistry23Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistry23 *AutomationRegistry23Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.Fallback(&_AutomationRegistry23.TransactOpts, calldata)
}

func (_AutomationRegistry23 *AutomationRegistry23TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry23.Contract.Fallback(&_AutomationRegistry23.TransactOpts, calldata)
}

type AutomationRegistry23AdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistry23AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23AdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistry23AdminPrivilegeConfigSet)
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

func (it *AutomationRegistry23AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistry23AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23AdminPrivilegeConfigSetIterator{contract: _AutomationRegistry23.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23AdminPrivilegeConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistry23AdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistry23AdminPrivilegeConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23BillingConfigOverriddenIterator struct {
	Event *AutomationRegistry23BillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23BillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23BillingConfigOverridden)
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
		it.Event = new(AutomationRegistry23BillingConfigOverridden)
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

func (it *AutomationRegistry23BillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23BillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23BillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23BillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23BillingConfigOverriddenIterator{contract: _AutomationRegistry23.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23BillingConfigOverridden)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistry23BillingConfigOverridden, error) {
	event := new(AutomationRegistry23BillingConfigOverridden)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23BillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistry23BillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23BillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23BillingConfigOverrideRemoved)
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
		it.Event = new(AutomationRegistry23BillingConfigOverrideRemoved)
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

func (it *AutomationRegistry23BillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23BillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23BillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23BillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23BillingConfigOverrideRemovedIterator{contract: _AutomationRegistry23.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23BillingConfigOverrideRemoved)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistry23BillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistry23BillingConfigOverrideRemoved)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23BillingConfigSetIterator struct {
	Event *AutomationRegistry23BillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23BillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23BillingConfigSet)
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
		it.Event = new(AutomationRegistry23BillingConfigSet)
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

func (it *AutomationRegistry23BillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23BillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23BillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistry23BillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23BillingConfigSetIterator{contract: _AutomationRegistry23.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23BillingConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistry23BillingConfigSet, error) {
	event := new(AutomationRegistry23BillingConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23CancelledUpkeepReportIterator struct {
	Event *AutomationRegistry23CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23CancelledUpkeepReport)
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
		it.Event = new(AutomationRegistry23CancelledUpkeepReport)
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

func (it *AutomationRegistry23CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23CancelledUpkeepReportIterator{contract: _AutomationRegistry23.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23CancelledUpkeepReport)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistry23CancelledUpkeepReport, error) {
	event := new(AutomationRegistry23CancelledUpkeepReport)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23ChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistry23ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23ChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistry23ChainSpecificModuleUpdated)
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

func (it *AutomationRegistry23ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistry23ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23ChainSpecificModuleUpdatedIterator{contract: _AutomationRegistry23.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23ChainSpecificModuleUpdated)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistry23ChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistry23ChainSpecificModuleUpdated)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23ConfigSetIterator struct {
	Event *AutomationRegistry23ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23ConfigSet)
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
		it.Event = new(AutomationRegistry23ConfigSet)
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

func (it *AutomationRegistry23ConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23ConfigSet struct {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistry23ConfigSetIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23ConfigSetIterator{contract: _AutomationRegistry23.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ConfigSet) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23ConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseConfigSet(log types.Log) (*AutomationRegistry23ConfigSet, error) {
	event := new(AutomationRegistry23ConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23DedupKeyAddedIterator struct {
	Event *AutomationRegistry23DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23DedupKeyAdded)
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
		it.Event = new(AutomationRegistry23DedupKeyAdded)
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

func (it *AutomationRegistry23DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistry23DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23DedupKeyAddedIterator{contract: _AutomationRegistry23.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23DedupKeyAdded)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistry23DedupKeyAdded, error) {
	event := new(AutomationRegistry23DedupKeyAdded)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23FeesWithdrawnIterator struct {
	Event *AutomationRegistry23FeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23FeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23FeesWithdrawn)
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
		it.Event = new(AutomationRegistry23FeesWithdrawn)
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

func (it *AutomationRegistry23FeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23FeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23FeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistry23FeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23FeesWithdrawnIterator{contract: _AutomationRegistry23.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23FeesWithdrawn)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistry23FeesWithdrawn, error) {
	event := new(AutomationRegistry23FeesWithdrawn)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23FundsAddedIterator struct {
	Event *AutomationRegistry23FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23FundsAdded)
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
		it.Event = new(AutomationRegistry23FundsAdded)
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

func (it *AutomationRegistry23FundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistry23FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23FundsAddedIterator{contract: _AutomationRegistry23.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23FundsAdded)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseFundsAdded(log types.Log) (*AutomationRegistry23FundsAdded, error) {
	event := new(AutomationRegistry23FundsAdded)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23FundsWithdrawnIterator struct {
	Event *AutomationRegistry23FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23FundsWithdrawn)
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
		it.Event = new(AutomationRegistry23FundsWithdrawn)
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

func (it *AutomationRegistry23FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23FundsWithdrawnIterator{contract: _AutomationRegistry23.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23FundsWithdrawn)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistry23FundsWithdrawn, error) {
	event := new(AutomationRegistry23FundsWithdrawn)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23InsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistry23InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23InsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistry23InsufficientFundsUpkeepReport)
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

func (it *AutomationRegistry23InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23InsufficientFundsUpkeepReportIterator{contract: _AutomationRegistry23.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23InsufficientFundsUpkeepReport)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistry23InsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistry23InsufficientFundsUpkeepReport)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23NOPsSettledOffchainIterator struct {
	Event *AutomationRegistry23NOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23NOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23NOPsSettledOffchain)
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
		it.Event = new(AutomationRegistry23NOPsSettledOffchain)
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

func (it *AutomationRegistry23NOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23NOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23NOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistry23NOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23NOPsSettledOffchainIterator{contract: _AutomationRegistry23.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23NOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23NOPsSettledOffchain)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistry23NOPsSettledOffchain, error) {
	event := new(AutomationRegistry23NOPsSettledOffchain)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23OwnershipTransferRequestedIterator struct {
	Event *AutomationRegistry23OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23OwnershipTransferRequested)
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
		it.Event = new(AutomationRegistry23OwnershipTransferRequested)
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

func (it *AutomationRegistry23OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistry23OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23OwnershipTransferRequestedIterator{contract: _AutomationRegistry23.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23OwnershipTransferRequested)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistry23OwnershipTransferRequested, error) {
	event := new(AutomationRegistry23OwnershipTransferRequested)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23OwnershipTransferredIterator struct {
	Event *AutomationRegistry23OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23OwnershipTransferred)
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
		it.Event = new(AutomationRegistry23OwnershipTransferred)
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

func (it *AutomationRegistry23OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistry23OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23OwnershipTransferredIterator{contract: _AutomationRegistry23.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23OwnershipTransferred)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistry23OwnershipTransferred, error) {
	event := new(AutomationRegistry23OwnershipTransferred)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23PausedIterator struct {
	Event *AutomationRegistry23Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23Paused)
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
		it.Event = new(AutomationRegistry23Paused)
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

func (it *AutomationRegistry23PausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistry23PausedIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23PausedIterator{contract: _AutomationRegistry23.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Paused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23Paused)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParsePaused(log types.Log) (*AutomationRegistry23Paused, error) {
	event := new(AutomationRegistry23Paused)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23PayeesUpdatedIterator struct {
	Event *AutomationRegistry23PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23PayeesUpdated)
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
		it.Event = new(AutomationRegistry23PayeesUpdated)
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

func (it *AutomationRegistry23PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistry23PayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23PayeesUpdatedIterator{contract: _AutomationRegistry23.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23PayeesUpdated)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistry23PayeesUpdated, error) {
	event := new(AutomationRegistry23PayeesUpdated)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23PayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistry23PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23PayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistry23PayeeshipTransferRequested)
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

func (it *AutomationRegistry23PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistry23PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23PayeeshipTransferRequestedIterator{contract: _AutomationRegistry23.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23PayeeshipTransferRequested)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistry23PayeeshipTransferRequested, error) {
	event := new(AutomationRegistry23PayeeshipTransferRequested)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23PayeeshipTransferredIterator struct {
	Event *AutomationRegistry23PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23PayeeshipTransferred)
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
		it.Event = new(AutomationRegistry23PayeeshipTransferred)
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

func (it *AutomationRegistry23PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistry23PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23PayeeshipTransferredIterator{contract: _AutomationRegistry23.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23PayeeshipTransferred)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistry23PayeeshipTransferred, error) {
	event := new(AutomationRegistry23PayeeshipTransferred)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23PaymentWithdrawnIterator struct {
	Event *AutomationRegistry23PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23PaymentWithdrawn)
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
		it.Event = new(AutomationRegistry23PaymentWithdrawn)
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

func (it *AutomationRegistry23PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistry23PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23PaymentWithdrawnIterator{contract: _AutomationRegistry23.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23PaymentWithdrawn)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistry23PaymentWithdrawn, error) {
	event := new(AutomationRegistry23PaymentWithdrawn)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23ReorgedUpkeepReportIterator struct {
	Event *AutomationRegistry23ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23ReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistry23ReorgedUpkeepReport)
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

func (it *AutomationRegistry23ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23ReorgedUpkeepReportIterator{contract: _AutomationRegistry23.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23ReorgedUpkeepReport)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistry23ReorgedUpkeepReport, error) {
	event := new(AutomationRegistry23ReorgedUpkeepReport)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23StaleUpkeepReportIterator struct {
	Event *AutomationRegistry23StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23StaleUpkeepReport)
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
		it.Event = new(AutomationRegistry23StaleUpkeepReport)
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

func (it *AutomationRegistry23StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23StaleUpkeepReportIterator{contract: _AutomationRegistry23.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23StaleUpkeepReport)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistry23StaleUpkeepReport, error) {
	event := new(AutomationRegistry23StaleUpkeepReport)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23TransmittedIterator struct {
	Event *AutomationRegistry23Transmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23TransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23Transmitted)
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
		it.Event = new(AutomationRegistry23Transmitted)
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

func (it *AutomationRegistry23TransmittedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23Transmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistry23TransmittedIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23TransmittedIterator{contract: _AutomationRegistry23.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Transmitted) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23Transmitted)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseTransmitted(log types.Log) (*AutomationRegistry23Transmitted, error) {
	event := new(AutomationRegistry23Transmitted)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UnpausedIterator struct {
	Event *AutomationRegistry23Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23Unpaused)
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
		it.Event = new(AutomationRegistry23Unpaused)
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

func (it *AutomationRegistry23UnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistry23UnpausedIterator, error) {

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UnpausedIterator{contract: _AutomationRegistry23.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Unpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23Unpaused)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUnpaused(log types.Log) (*AutomationRegistry23Unpaused, error) {
	event := new(AutomationRegistry23Unpaused)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistry23UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistry23UpkeepAdminTransferRequested)
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

func (it *AutomationRegistry23UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistry23UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepAdminTransferRequestedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepAdminTransferRequested)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistry23UpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistry23UpkeepAdminTransferRequested)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepAdminTransferredIterator struct {
	Event *AutomationRegistry23UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistry23UpkeepAdminTransferred)
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

func (it *AutomationRegistry23UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistry23UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepAdminTransferredIterator{contract: _AutomationRegistry23.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepAdminTransferred)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistry23UpkeepAdminTransferred, error) {
	event := new(AutomationRegistry23UpkeepAdminTransferred)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepCanceledIterator struct {
	Event *AutomationRegistry23UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepCanceled)
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
		it.Event = new(AutomationRegistry23UpkeepCanceled)
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

func (it *AutomationRegistry23UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistry23UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepCanceledIterator{contract: _AutomationRegistry23.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepCanceled)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistry23UpkeepCanceled, error) {
	event := new(AutomationRegistry23UpkeepCanceled)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepChargedIterator struct {
	Event *AutomationRegistry23UpkeepCharged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepChargedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepCharged)
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
		it.Event = new(AutomationRegistry23UpkeepCharged)
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

func (it *AutomationRegistry23UpkeepChargedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepChargedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepCharged struct {
	Id      *big.Int
	Receipt AutomationRegistryBase23PaymentReceipt
	Raw     types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepChargedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepChargedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepCharged", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCharged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepCharged)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepCharged(log types.Log) (*AutomationRegistry23UpkeepCharged, error) {
	event := new(AutomationRegistry23UpkeepCharged)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepCheckDataSetIterator struct {
	Event *AutomationRegistry23UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistry23UpkeepCheckDataSet)
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

func (it *AutomationRegistry23UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepCheckDataSetIterator{contract: _AutomationRegistry23.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepCheckDataSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistry23UpkeepCheckDataSet, error) {
	event := new(AutomationRegistry23UpkeepCheckDataSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepGasLimitSetIterator struct {
	Event *AutomationRegistry23UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistry23UpkeepGasLimitSet)
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

func (it *AutomationRegistry23UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepGasLimitSetIterator{contract: _AutomationRegistry23.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepGasLimitSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistry23UpkeepGasLimitSet, error) {
	event := new(AutomationRegistry23UpkeepGasLimitSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepMigratedIterator struct {
	Event *AutomationRegistry23UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepMigrated)
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
		it.Event = new(AutomationRegistry23UpkeepMigrated)
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

func (it *AutomationRegistry23UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepMigratedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepMigrated)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistry23UpkeepMigrated, error) {
	event := new(AutomationRegistry23UpkeepMigrated)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistry23UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistry23UpkeepOffchainConfigSet)
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

func (it *AutomationRegistry23UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepOffchainConfigSetIterator{contract: _AutomationRegistry23.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepOffchainConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistry23UpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistry23UpkeepOffchainConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepPausedIterator struct {
	Event *AutomationRegistry23UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepPaused)
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
		it.Event = new(AutomationRegistry23UpkeepPaused)
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

func (it *AutomationRegistry23UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepPausedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepPaused)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistry23UpkeepPaused, error) {
	event := new(AutomationRegistry23UpkeepPaused)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepPerformedIterator struct {
	Event *AutomationRegistry23UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepPerformed)
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
		it.Event = new(AutomationRegistry23UpkeepPerformed)
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

func (it *AutomationRegistry23UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistry23UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepPerformedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepPerformed)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistry23UpkeepPerformed, error) {
	event := new(AutomationRegistry23UpkeepPerformed)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistry23UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistry23UpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistry23UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistry23.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepPrivilegeConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistry23UpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistry23UpkeepPrivilegeConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepReceivedIterator struct {
	Event *AutomationRegistry23UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepReceived)
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
		it.Event = new(AutomationRegistry23UpkeepReceived)
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

func (it *AutomationRegistry23UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepReceivedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepReceived)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistry23UpkeepReceived, error) {
	event := new(AutomationRegistry23UpkeepReceived)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepRegisteredIterator struct {
	Event *AutomationRegistry23UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepRegistered)
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
		it.Event = new(AutomationRegistry23UpkeepRegistered)
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

func (it *AutomationRegistry23UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepRegisteredIterator{contract: _AutomationRegistry23.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepRegistered)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistry23UpkeepRegistered, error) {
	event := new(AutomationRegistry23UpkeepRegistered)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistry23UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistry23UpkeepTriggerConfigSet)
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

func (it *AutomationRegistry23UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepTriggerConfigSetIterator{contract: _AutomationRegistry23.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepTriggerConfigSet)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistry23UpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistry23UpkeepTriggerConfigSet)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistry23UpkeepUnpausedIterator struct {
	Event *AutomationRegistry23UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistry23UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistry23UpkeepUnpaused)
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
		it.Event = new(AutomationRegistry23UpkeepUnpaused)
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

func (it *AutomationRegistry23UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistry23UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistry23UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry23UpkeepUnpausedIterator{contract: _AutomationRegistry23.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry23 *AutomationRegistry23Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry23.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistry23UpkeepUnpaused)
				if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistry23 *AutomationRegistry23Filterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistry23UpkeepUnpaused, error) {
	event := new(AutomationRegistry23UpkeepUnpaused)
	if err := _AutomationRegistry23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

func (_AutomationRegistry23 *AutomationRegistry23) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistry23.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistry23.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistry23.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistry23.ParseBillingConfigOverridden(log)
	case _AutomationRegistry23.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistry23.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistry23.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistry23.ParseBillingConfigSet(log)
	case _AutomationRegistry23.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistry23.ParseCancelledUpkeepReport(log)
	case _AutomationRegistry23.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistry23.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistry23.abi.Events["ConfigSet"].ID:
		return _AutomationRegistry23.ParseConfigSet(log)
	case _AutomationRegistry23.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistry23.ParseDedupKeyAdded(log)
	case _AutomationRegistry23.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistry23.ParseFeesWithdrawn(log)
	case _AutomationRegistry23.abi.Events["FundsAdded"].ID:
		return _AutomationRegistry23.ParseFundsAdded(log)
	case _AutomationRegistry23.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistry23.ParseFundsWithdrawn(log)
	case _AutomationRegistry23.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistry23.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistry23.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistry23.ParseNOPsSettledOffchain(log)
	case _AutomationRegistry23.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistry23.ParseOwnershipTransferRequested(log)
	case _AutomationRegistry23.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistry23.ParseOwnershipTransferred(log)
	case _AutomationRegistry23.abi.Events["Paused"].ID:
		return _AutomationRegistry23.ParsePaused(log)
	case _AutomationRegistry23.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistry23.ParsePayeesUpdated(log)
	case _AutomationRegistry23.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistry23.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistry23.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistry23.ParsePayeeshipTransferred(log)
	case _AutomationRegistry23.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistry23.ParsePaymentWithdrawn(log)
	case _AutomationRegistry23.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistry23.ParseReorgedUpkeepReport(log)
	case _AutomationRegistry23.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistry23.ParseStaleUpkeepReport(log)
	case _AutomationRegistry23.abi.Events["Transmitted"].ID:
		return _AutomationRegistry23.ParseTransmitted(log)
	case _AutomationRegistry23.abi.Events["Unpaused"].ID:
		return _AutomationRegistry23.ParseUnpaused(log)
	case _AutomationRegistry23.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistry23.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistry23.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistry23.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistry23.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistry23.ParseUpkeepCanceled(log)
	case _AutomationRegistry23.abi.Events["UpkeepCharged"].ID:
		return _AutomationRegistry23.ParseUpkeepCharged(log)
	case _AutomationRegistry23.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistry23.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistry23.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistry23.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistry23.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistry23.ParseUpkeepMigrated(log)
	case _AutomationRegistry23.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistry23.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistry23.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistry23.ParseUpkeepPaused(log)
	case _AutomationRegistry23.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistry23.ParseUpkeepPerformed(log)
	case _AutomationRegistry23.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistry23.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistry23.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistry23.ParseUpkeepReceived(log)
	case _AutomationRegistry23.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistry23.ParseUpkeepRegistered(log)
	case _AutomationRegistry23.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistry23.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistry23.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistry23.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistry23AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistry23BillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistry23BillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistry23BillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba3")
}

func (AutomationRegistry23CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistry23ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistry23ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (AutomationRegistry23DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistry23FeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistry23FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistry23FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistry23InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistry23NOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistry23OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistry23OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistry23Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistry23PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistry23PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistry23PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistry23PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistry23ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistry23StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistry23Transmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (AutomationRegistry23Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistry23UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistry23UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistry23UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistry23UpkeepCharged) Topic() common.Hash {
	return common.HexToHash("0x801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc")
}

func (AutomationRegistry23UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistry23UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistry23UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistry23UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistry23UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistry23UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistry23UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistry23UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistry23UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistry23UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistry23UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistry23 *AutomationRegistry23) Address() common.Address {
	return _AutomationRegistry23.address
}

type AutomationRegistry23Interface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistry23AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistry23AdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23BillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistry23BillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23BillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistry23BillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistry23BillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23BillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistry23BillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistry23CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistry23ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistry23ChainSpecificModuleUpdated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistry23ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*AutomationRegistry23ConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistry23DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistry23DedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistry23FeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistry23FeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistry23FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistry23FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistry23FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistry23InsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistry23NOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23NOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistry23NOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistry23OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistry23OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistry23OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistry23OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistry23PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistry23Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistry23PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistry23PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistry23PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistry23PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistry23PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistry23PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistry23PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistry23PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistry23ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistry23StaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistry23TransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Transmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*AutomationRegistry23Transmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistry23UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistry23Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistry23UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistry23UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistry23UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistry23UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistry23UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistry23UpkeepCanceled, error)

	FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepChargedIterator, error)

	WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCharged, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCharged(log types.Log) (*AutomationRegistry23UpkeepCharged, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistry23UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistry23UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistry23UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistry23UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistry23UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistry23UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistry23UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistry23UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistry23UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistry23UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistry23UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistry23UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistry23UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistry23UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
