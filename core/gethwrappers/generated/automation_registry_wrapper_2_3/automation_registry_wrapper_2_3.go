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

var AutomationRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicA2_3\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"overrides\",\"type\":\"tuple\"}],\"name\":\"BillingConfigOverridden\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"BillingConfigOverrideRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payments\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackNativePrice\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"contractIERC20[]\",\"name\":\"billingTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig[]\",\"name\":\"billingConfigs\",\"type\":\"tuple[]\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101806040523480156200001257600080fd5b50604051620064ff380380620064ff83398101604081905262000035916200062f565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b91906200062f565b826001600160a01b031663226cf83c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200062f565b836001600160a01b031663614486af6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200062f565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200062f565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200062f565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200062f565b876001600160a01b031663c5b964e06040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002d3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002f9919062000656565b886001600160a01b031663ac4dc59a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000338573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200035e91906200062f565b3380600081620003b55760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620003e857620003e8816200056b565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff19166001838181111562000439576200043962000679565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200049a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004c091906200068f565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000504573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200052a91906200068f565b60ff16146200054c576040516301f86e1760e41b815260040160405180910390fd5b5050506001600160a01b039095166101605250620006b4945050505050565b336001600160a01b03821603620005c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620003ac565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200062c57600080fd5b50565b6000602082840312156200064257600080fd5b81516200064f8162000616565b9392505050565b6000602082840312156200066957600080fd5b8151600281106200064f57600080fd5b634e487b7160e01b600052602160045260246000fd5b600060208284031215620006a257600080fd5b815160ff811681146200064f57600080fd5b60805160a05160c05160e05161010051610120516101405161016051615dc8620007376000396000818160c9015261017c015260005050600061226a015260005050600050506000613a2801526000505060008181610df301528181610f010152818161102c015281816110760152818161176f015261336e0152615dc86000f3fe6080604052600436106100c75760003560e01c8063a4c0ed3611610074578063b1dc65a41161004e578063b1dc65a414610338578063e3d0e71214610358578063f2fde38b14610378576100c7565b8063a4c0ed361461029b578063aed2e929146102bb578063afcb95d7146102f2576100c7565b80637bee169d116100a55780637bee169d146101d657806381ff7048146101f65780638da5cb5b14610270576100c7565b8063181f5a771461010e578063349e8cca1461016d57806379ba5097146101c1575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610107573d6000f35b3d6000fd5b005b34801561011a57600080fd5b506101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e332e30000000000000000081525081565b604051610164919061476f565b60405180910390f35b34801561017957600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b3480156101cd57600080fd5b5061010c610398565b3480156101e257600080fd5b5061010c6101f1366004614ce4565b61049a565b34801561020257600080fd5b5061024d60175460135463ffffffff74010000000000000000000000000000000000000000830481169378010000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b34801561027c57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661019c565b3480156102a757600080fd5b5061010c6102b6366004614e45565b610ddb565b3480156102c757600080fd5b506102db6102d6366004614ea1565b6110f7565b604080519215158352602083019190915201610164565b3480156102fe57600080fd5b50601354601454604080516000815260208101939093526c0100000000000000000000000090910463ffffffff1690820152606001610164565b34801561034457600080fd5b5061010c610353366004614f32565b6112a2565b34801561036457600080fd5b5061010c610373366004614fe9565b611583565b34801561038457600080fd5b5061010c6103933660046150b6565b6115bd565b60015473ffffffffffffffffffffffffffffffffffffffff16331461041e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104a26115d1565b601f885111156104de576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8560ff1660000361051b576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8651885114158061053a5750610532866003615109565b60ff16885111155b15610571576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80518251146105ac576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105b68282611654565b6105c08888611abc565b604051806101200160405280601460000160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001600063ffffffff1681526020018660a0015162ffffff16815260200186610120015161ffff1681526020018760ff168152602001601460000160169054906101000a900460ff1615158152602001601460000160179054906101000a900460ff1615158152602001866080015115158152602001866101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548162ffffff021916908362ffffff16021790555060608201518160000160136101000a81548161ffff021916908361ffff16021790555060808201518160000160156101000a81548160ff021916908360ff16021790555060a08201518160000160166101000a81548160ff02191690831515021790555060c08201518160000160176101000a81548160ff02191690831515021790555060e08201518160000160186101000a81548160ff0219169083151502179055506101008201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550905050604051806101600160405280866060015173ffffffffffffffffffffffffffffffffffffffff168152602001866000015163ffffffff168152602001866020015163ffffffff1681526020016016600001601c9054906101000a900463ffffffff1663ffffffff16815260200186610100015173ffffffffffffffffffffffffffffffffffffffff168152602001601660010160149054906101000a900463ffffffff1663ffffffff168152602001601660010160189054906101000a900463ffffffff1663ffffffff168152602001866040015163ffffffff16815260200186610140015173ffffffffffffffffffffffffffffffffffffffff1681526020018660c0015163ffffffff1681526020018660e0015163ffffffff16815250601660008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160186101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101208201518160020160146101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160186101000a81548163ffffffff021916908363ffffffff160217905550905050846101600151601981905550846101800151601a81905550846101a00151601b819055506000601660010160189054906101000a900463ffffffff169050856101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610bc6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bea9190615125565b601780547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff9384160217808255600192601491610c619185917401000000000000000000000000000000000000000090041661513e565b92506101000a81548163ffffffff021916908363ffffffff160217905550600086604051602001610c9291906151ac565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152919052601754909150610cf7904690309074010000000000000000000000000000000000000000900463ffffffff168d8d8d878d8d61217d565b601355600960008181610d0a8282614656565b5050505060005b876101c0015151811015610d6457610d51886101c001518281518110610d3957610d39615332565b6020026020010151600961222790919063ffffffff16565b5080610d5c81615361565b915050610d11565b506013546017546040517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0592610dc79286927401000000000000000000000000000000000000000090910463ffffffff16908f908f908f9089908f908f90615399565b60405180910390a150505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610e4a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610e84576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610e928284018461542f565b6000818152600460205260409020549091506601000000000000900463ffffffff90811614610eed576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600201547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff908116911614610f71576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902060010154610fb090859070010000000000000000000000000000000090046bffffffffffffffffffffffff16615448565b600082815260046020908152604080832060010180546bffffffffffffffffffffffff95909516700100000000000000000000000000000000027fffffffff000000000000000000000000ffffffffffffffffffffffffffffffff9095169490941790935573ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016825260219052205461105f90859061546d565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660009081526021602090815260409182902093909355516bffffffffffffffffffffffff871681529087169183917fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203910160405180910390a35050505050565b600080611102612252565b601454760100000000000000000000000000000000000000000000900460ff1615611159576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615158386015262010000820463ffffffff90811684880181905266010000000000008404821660608601526a010000000000000000000090930473ffffffffffffffffffffffffffffffffffffffff9081166080860181905260018701546fffffffffffffffffffffffffffffffff811660a088015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08801527c0100000000000000000000000000000000000000000000000000000000900490921660e0860152600290950154909416908301528451601f890185900485028101850190955287855290936112959389908990819084018382808284376000920191909152506122c192505050565b9097909650945050505050565b60005a60408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff75010000000000000000000000000000000000000000008304811660808301527601000000000000000000000000000000000000000000008304811615801560a08401527701000000000000000000000000000000000000000000000084048216151560c0840152780100000000000000000000000000000000000000000000000090930416151560e082015260155473ffffffffffffffffffffffffffffffffffffffff16610100820152919250611406576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff1661144f576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013548a351461148b576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608081015161149b906001615480565b60ff16861415806114ac5750858414155b156114e3576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6114f38a8a8a8a8a8a8a8a6124dc565b60006114ff8a8a612745565b905060208b0135600881901c63ffffffff1661151c8484876127fe565b836020015163ffffffff168163ffffffff16111561157457601480547fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c0100000000000000000000000063ffffffff8416021790555b50505050505050505050505050565b60008060008580602001905181019061159c9190615657565b9250925092506115b2898989868989888861049a565b505050505050505050565b6115c56115d1565b6115ce816133e4565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611652576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610415565b565b60005b60245481101561171257602260006024838154811061167857611678615332565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001812080547fffffffffff000000000000000000000000000000000000000000000000000000168155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690558061170a81615361565b915050611657565b5061171f60246000614656565b60255460ff1660005b8351811015611ab657600084828151811061174557611745615332565b60200260200101519050600084838151811061176357611763615332565b602002602001015190507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161480156117d9575060018460018111156117d7576117d7615802565b145b15611810576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216158061184b5750604081015173ffffffffffffffffffffffffffffffffffffffff16155b15611882576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff828116600090815260226020526040902054670100000000000000900416156118ec576040517f357d0cc400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6024805460018082019092557f7cd332d19b93bcabe3cce7ca0c18a052f57e5fd03b4758a09f30f5ddc4b22ec40180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8581169182179092556000818152602260209081526040918290208651815488840180518a8701805163ffffffff9095167fffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000909416841764010000000062ffffff93841602177fffffffffff0000000000000000000000000000000000000000ffffffffffffff16670100000000000000958b16959095029490941785556060808c0180519b87019b909b556080808d018051600290980180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff998a1617905589519586529351909216968401969096529251909716948101949094529551918301919091529251909216928201929092527f720a5849025dc4fd0061aed1bb30efd713cde64ce7f8d807953ecca27c8f143c9060a00160405180910390a250508080611aae90615361565b915050611728565b50505050565b60005b600e54811015611b3357611b20600e8281548110611adf57611adf615332565b600091825260209091200154601454600e5473ffffffffffffffffffffffffffffffffffffffff909216916bffffffffffffffffffffffff909116906134d9565b5080611b2b81615361565b915050611abf565b5060255460009060ff16815b600e54811015611ca457600e8181548110611b5c57611b5c615332565b6000918252602082200154600d805473ffffffffffffffffffffffffffffffffffffffff9092169550600c929184908110611b9957611b99615332565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168452838201949094526040928301822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559286168152600b909252902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556001826001811115611c3b57611c3b615802565b148015611c80575073ffffffffffffffffffffffffffffffffffffffff83166000908152600b60205260409020546201000090046bffffffffffffffffffffffff1615155b15611c9257611c90600f84612227565b505b80611c9c81615361565b915050611b3f565b50611cb1600d6000614656565b611cbd600e6000614656565b6040805160808101825260008082526020820181905291810182905260608101829052905b855181101561214d57600c6000878381518110611d0157611d01615332565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615611d6c576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16868281518110611d9657611d96615332565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611deb576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c6000888481518110611e1c57611e1c615332565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558451859082908110611ec457611ec4615332565b60200260200101519350600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603611f34576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e01000000000000000000000000000090049092166060830152909250611fef576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180835260ff80831660208086019182526014546bffffffffffffffffffffffff9081166060880190815273ffffffffffffffffffffffffffffffffffffffff8a166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905583600181111561212957612129615802565b0361213b57612139600f856136e1565b505b8061214581615361565b915050611ce2565b50845161216190600d906020880190614674565b50835161217590600e906020870190614674565b505050505050565b6000808a8a8a8a8a8a8a8a8a6040516020016121a199989796959493929190615831565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b60006122498373ffffffffffffffffffffffffffffffffffffffff8416613703565b90505b92915050565b3273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611652576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600090819077010000000000000000000000000000000000000000000000900460ff161561231e576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16770100000000000000000000000000000000000000000000001790556040517f4585e33b000000000000000000000000000000000000000000000000000000009061239390859060240161476f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061246690879087906004016158c6565b60408051808303816000875af1158015612484573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124a891906158df565b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b600087876040516124ee92919061590d565b604051908190038120612505918b9060200161591d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b888110156126dc5760018587836020811061257157612571615332565b61257e91901a601b615480565b8c8c8581811061259057612590615332565b905060200201358b8b868181106125a9576125a9615332565b90506020020135604051600081526020016040526040516125e6949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612608573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550935090506126b6576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b8401935080806126d490615361565b915050612554565b50827e01010101010101010101010101010101010101010101010101010101010101841614612737576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b61277e6040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061278c83850185615a0e565b60408101515160608201515191925090811415806127af57508082608001515114155b806127bf5750808260a001515114155b156127f6576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b509392505050565b600082604001515167ffffffffffffffff81111561281e5761281e614782565b6040519080825280602002602001820160405280156128ea57816020015b6040805161020081018252600060e08201818152610100830182905261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c083018290526101e0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191018161283c5790505b50905060006040518060800160405280600061ffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160008152509050600085610100015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015612988573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129ac9190615125565b9050600086610100015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa158015612a00573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a249190615125565b905060005b866040015151811015612eb4576004600088604001518381518110612a5057612a50615332565b6020908102919091018101518252818101929092526040908101600020815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490931660e0820152600290910154909116918101919091528551869083908110612b7357612b73615332565b602002602001015160000181905250612ba887604001518281518110612b9b57612b9b615332565b6020026020010151613752565b858281518110612bba57612bba615332565b6020026020010151606001906001811115612bd757612bd7615802565b90816001811115612bea57612bea615802565b81525050612c4e87604001518281518110612c0757612c07615332565b60200260200101518489608001518481518110612c2657612c26615332565b6020026020010151888581518110612c4057612c40615332565b60200260200101518c6137fd565b868381518110612c6057612c60615332565b6020026020010151602001878481518110612c7d57612c7d615332565b602002602001015160c0018281525082151515158152505050848181518110612ca857612ca8615332565b60200260200101516020015115612cd857600184600001818151612ccc9190615afb565b61ffff16905250612cdd565b612ea2565b612d43858281518110612cf257612cf2615332565b6020026020010151600001516080015188606001518381518110612d1857612d18615332565b60200260200101518960a001518481518110612d3657612d36615332565b60200260200101516122c1565b868381518110612d5557612d55615332565b6020026020010151604001878481518110612d7257612d72615332565b6020026020010151608001828152508215151515815250505087608001516001612d9c9190615480565b612daa9060ff166040615b16565b6103a48860a001518381518110612dc357612dc3615332565b602002602001015151612dd6919061546d565b612de0919061546d565b858281518110612df257612df2615332565b602002602001015160a0018181525050848181518110612e1457612e14615332565b602002602001015160a0015184606001818151612e31919061546d565b9052508451859082908110612e4857612e48615332565b60200260200101516080015186612e5f9190615b2d565b9550612ea287604001518281518110612e7a57612e7a615332565b602002602001015184878481518110612e9557612e95615332565b602002602001015161391c565b80612eac81615361565b915050612a29565b50825161ffff16600003612ecb5750505050505050565b61c738612ed9366010615b16565b5a612ee49088615b2d565b612eee919061546d565b612ef8919061546d565b835190955061232890612f0f9061ffff1687615b6f565b612f19919061546d565b60408051606081018252600080825260208201819052918101829052919650612f4189613a21565b905060005b88604001515181101561327d57868181518110612f6557612f65615332565b6020026020010151602001511561326b57801580612ffd575086612f8a600183615b2d565b81518110612f9a57612f9a615332565b602002602001015160000151610100015173ffffffffffffffffffffffffffffffffffffffff16878281518110612fd357612fd3615332565b602002602001015160000151610100015173ffffffffffffffffffffffffffffffffffffffff1614155b156130315761302e8a88838151811061301857613018615332565b6020026020010151600001516101000151613b12565b92505b600061314f8b6040518061012001604052808b868151811061305557613055615332565b60200260200101516080015181526020018c81526020018a606001518c878151811061308357613083615332565b602002602001015160a001518a61309a9190615b16565b6130a49190615b6f565b81526020018d6000015181526020018d6020015181526020018681526020018b86815181106130d5576130d5615332565b602002602001015160000151610100015173ffffffffffffffffffffffffffffffffffffffff168152602001878152602001600115158152508c60400151858151811061312457613124615332565b60200260200101518b868151811061313e5761313e615332565b602002602001015160000151613c63565b90508060600151876040018181516131679190615448565b6bffffffffffffffffffffffff16905250604081015160208801805161318e908390615448565b6bffffffffffffffffffffffff1690525087518890839081106131b3576131b3615332565b60200260200101516040015115158a6040015183815181106131d7576131d7615332565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b836060015184604001516132149190615448565b8b868151811061322657613226615332565b6020026020010151608001518d8f60800151888151811061324957613249615332565b60200260200101516040516132619493929190615b83565b60405180910390a3505b8061327581615361565b915050612f46565b505050602083810151336000908152600b909252604090912080546002906132ba9084906201000090046bffffffffffffffffffffffff16615448565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508260400151601460000160008282829054906101000a90046bffffffffffffffffffffffff166133189190615448565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508260400151836020015161335a9190615448565b6bffffffffffffffffffffffff16602160007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546133d6919061546d565b909155505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613463576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610415565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906136d55760008160600151856135719190615bc0565b9050600061357f8583615be5565b905080836040018181516135939190615448565b6bffffffffffffffffffffffff169052506135ae8582615c10565b836060018181516135bf9190615448565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b60006122498373ffffffffffffffffffffffffffffffffffffffff8416613f5f565b600081815260018301602052604081205461374a5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561224c565b50600061224c565b6000818160045b600f8110156137df577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061379757613797615332565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146137cd57506000949350505050565b806137d781615361565b915050613759565b5081600f1a60018111156137f5576137f5615802565b949350505050565b60008080808560600151600181111561381857613818615802565b0361383e5761382a8888888888614052565b61383957600092509050613912565b6138b6565b60018560600151600181111561385657613856615802565b03613884576000613869898989886141dc565b925090508061387e5750600092509050613912565b506138b6565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516060015163ffffffff16871061390b57877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636876040516138f8919061476f565b60405180910390a2600092509050613912565b6001925090505b9550959350505050565b60008160600151600181111561393457613934615802565b0361399857600083815260046020526040902060010180547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff851602179055505050565b6001816060015160018111156139b0576139b0615802565b03613a1c5760c08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613a91573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ab59190615c5a565b50935050925050600082131580613acb57508042105b80613afb57506000846040015162ffffff16118015613afb5750613aef8142615b2d565b846040015162ffffff16105b15613b0b575050601b5492915050565b5092915050565b604080516060810182526000808252602080830182815283850183905273ffffffffffffffffffffffffffffffffffffffff868116845260229092528483208054640100000000810462ffffff1690925263ffffffff8216855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495909484936701000000000000009093049092169163feaf968c9160048082019260a0929091908290030181865afa158015613bd4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bf89190615c5a565b50935050925050600082131580613c0e57508042105b80613c3e57506000866040015162ffffff16118015613c3e5750613c328142615b2d565b866040015162ffffff16105b15613c525760018301546040850152613c5a565b604084018290525b50505092915050565b604080516080810182526000808252602080830182905292820181905260608201529082015115613cdc5760008381526023602090815260409182902082518084019093525463ffffffff811680845262ffffff640100000000909204821693830193845260e088018051919091529251925192169101525b6000613ce886866143e9565b60c0840151602082015182519293509091600091613d0591615448565b905082600001516bffffffffffffffffffffffff16826bffffffffffffffffffffffff161015613d8a57819050613d6b87608001518860e0015160400151846bffffffffffffffffffffffff16613d5c9190615b16565b613d669190615b6f565b6145b4565b6bffffffffffffffffffffffff16604084015260006060840152613e16565b806bffffffffffffffffffffffff16826bffffffffffffffffffffffff161015613e1657819050613e0283604001516bffffffffffffffffffffffff1688608001518960e0015160400151856bffffffffffffffffffffffff16613dee9190615b16565b613df89190615b6f565b613d669190615b2d565b6bffffffffffffffffffffffff1660608401525b60008681526004602052604090206001018054829190601090613e5c90849070010000000000000000000000000000000090046bffffffffffffffffffffffff16615bc0565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008881526004602052604081206001018054928516935091613eb79084906fffffffffffffffffffffffffffffffff16615caa565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550806bffffffffffffffffffffffff16602160008960c0015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254613f4e9190615b2d565b909155509298975050505050505050565b60008181526001830160205260408120548015614048576000613f83600183615b2d565b8554909150600090613f9790600190615b2d565b9050818114613ffc576000866000018281548110613fb757613fb7615332565b9060005260206000200154905080876000018481548110613fda57613fda615332565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061400d5761400d615cd3565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061224c565b600091505061224c565b600080848060200190518101906140699190615d02565b845160e00151815191925063ffffffff908116911610156140c657867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516140b4919061476f565b60405180910390a260009150506141d3565b8260e00151801561418657506020810151158015906141865750602081015161010084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561415f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141839190615125565b14155b806141985750805163ffffffff168611155b156141cd57867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516140b4919061476f565b60019150505b95945050505050565b6000806000848060200190518101906141f59190615d5a565b905060008782600001518360200151846040015160405160200161425794939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b6040516020818303038152906040528051906020012090508460e00151801561433257506080820151158015906143325750608082015161010086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561430b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061432f9190615125565b14155b80614347575086826060015163ffffffff1610155b1561439157877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc3018760405161437c919061476f565b60405180910390a26000935091506143e09050565b60008181526008602052604090205460ff16156143d857877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e88760405161437c919061476f565b600193509150505b94509492505050565b6040805160808101825260008082526020820181905291810182905260608101919091526000836060015161ffff1683606001516144279190615b16565b9050826101000151801561443a5750803a105b1561444257503a5b60008360a00151846040015185602001518660000151614462919061546d565b61446c9085615b16565b614476919061546d565b6144809190615b16565b90506144998460e001516040015182613d669190615b6f565b6bffffffffffffffffffffffff16835260808401516144bc90613d669083615b6f565b6bffffffffffffffffffffffff16604084015260e0840151602001516000906144f39062ffffff16683635c9adc5dea00000615b16565b9050600081633b9aca008760a001518860e001516000015163ffffffff1689604001518a60000151896145269190615b16565b614530919061546d565b61453a9190615b16565b6145449190615b16565b61454e9190615b6f565b614558919061546d565b90506145718660e001516040015182613d669190615b6f565b6bffffffffffffffffffffffff166020860152608086015161459790613d669083615b6f565b6bffffffffffffffffffffffff1660608601525050505092915050565b60006bffffffffffffffffffffffff821115614652576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610415565b5090565b50805460008255906000526020600020908101906115ce91906146f6565b8280548282559060005260206000209081019282156146ee579160200282015b828111156146ee57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614694565b506146529291505b5b8082111561465257600081556001016146f7565b6000815180845260005b8181101561473157602081850181015186830182015201614715565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000612249602083018461470b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610200810167ffffffffffffffff811182821017156147d5576147d5614782565b60405290565b60405160a0810167ffffffffffffffff811182821017156147d5576147d5614782565b60405160c0810167ffffffffffffffff811182821017156147d5576147d5614782565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561486857614868614782565b604052919050565b600067ffffffffffffffff82111561488a5761488a614782565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff811681146115ce57600080fd5b80356148c181614894565b919050565b600082601f8301126148d757600080fd5b813560206148ec6148e783614870565b614821565b82815260059290921b8401810191818101908684111561490b57600080fd5b8286015b8481101561492f57803561492281614894565b835291830191830161490f565b509695505050505050565b803560ff811681146148c157600080fd5b63ffffffff811681146115ce57600080fd5b80356148c18161494b565b80151581146115ce57600080fd5b80356148c181614968565b62ffffff811681146115ce57600080fd5b80356148c181614981565b61ffff811681146115ce57600080fd5b80356148c18161499d565b600061020082840312156149cb57600080fd5b6149d36147b1565b90506149de8261495d565b81526149ec6020830161495d565b60208201526149fd6040830161495d565b6040820152614a0e606083016148b6565b6060820152614a1f60808301614976565b6080820152614a3060a08301614992565b60a0820152614a4160c0830161495d565b60c0820152614a5260e0830161495d565b60e0820152610100614a658184016148b6565b90820152610120614a778382016149ad565b90820152610140614a898382016148b6565b90820152610160828101359082015261018080830135908201526101a080830135908201526101c08083013567ffffffffffffffff811115614aca57600080fd5b614ad6858286016148c6565b8284015250506101e0614aea8184016148b6565b9082015292915050565b803567ffffffffffffffff811681146148c157600080fd5b600082601f830112614b1d57600080fd5b813567ffffffffffffffff811115614b3757614b37614782565b614b6860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614821565b818152846020838601011115614b7d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112614bab57600080fd5b81356020614bbb6148e783614870565b82815260059290921b84018101918181019086841115614bda57600080fd5b8286015b8481101561492f578035614bf181614894565b8352918301918301614bde565b6bffffffffffffffffffffffff811681146115ce57600080fd5b600082601f830112614c2957600080fd5b81356020614c396148e783614870565b82815260a09283028501820192828201919087851115614c5857600080fd5b8387015b85811015614cd75781818a031215614c745760008081fd5b614c7c6147db565b8135614c878161494b565b815281860135614c9681614981565b81870152604082810135614ca981614894565b9082015260608281013590820152608080830135614cc681614bfe565b908201528452928401928101614c5c565b5090979650505050505050565b600080600080600080600080610100898b031215614d0157600080fd5b883567ffffffffffffffff80821115614d1957600080fd5b614d258c838d016148c6565b995060208b0135915080821115614d3b57600080fd5b614d478c838d016148c6565b9850614d5560408c0161493a565b975060608b0135915080821115614d6b57600080fd5b614d778c838d016149b8565b9650614d8560808c01614af4565b955060a08b0135915080821115614d9b57600080fd5b614da78c838d01614b0c565b945060c08b0135915080821115614dbd57600080fd5b614dc98c838d01614b9a565b935060e08b0135915080821115614ddf57600080fd5b50614dec8b828c01614c18565b9150509295985092959890939650565b60008083601f840112614e0e57600080fd5b50813567ffffffffffffffff811115614e2657600080fd5b602083019150836020828501011115614e3e57600080fd5b9250929050565b60008060008060608587031215614e5b57600080fd5b8435614e6681614894565b935060208501359250604085013567ffffffffffffffff811115614e8957600080fd5b614e9587828801614dfc565b95989497509550505050565b600080600060408486031215614eb657600080fd5b83359250602084013567ffffffffffffffff811115614ed457600080fd5b614ee086828701614dfc565b9497909650939450505050565b60008083601f840112614eff57600080fd5b50813567ffffffffffffffff811115614f1757600080fd5b6020830191508360208260051b8501011115614e3e57600080fd5b60008060008060008060008060e0898b031215614f4e57600080fd5b606089018a811115614f5f57600080fd5b8998503567ffffffffffffffff80821115614f7957600080fd5b614f858c838d01614dfc565b909950975060808b0135915080821115614f9e57600080fd5b614faa8c838d01614eed565b909750955060a08b0135915080821115614fc357600080fd5b50614fd08b828c01614eed565b999c989b50969995989497949560c00135949350505050565b60008060008060008060c0878903121561500257600080fd5b863567ffffffffffffffff8082111561501a57600080fd5b6150268a838b016148c6565b9750602089013591508082111561503c57600080fd5b6150488a838b016148c6565b965061505660408a0161493a565b9550606089013591508082111561506c57600080fd5b6150788a838b01614b0c565b945061508660808a01614af4565b935060a089013591508082111561509c57600080fd5b506150a989828a01614b0c565b9150509295509295509295565b6000602082840312156150c857600080fd5b81356150d381614894565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160290811690818114613b0b57613b0b6150da565b60006020828403121561513757600080fd5b5051919050565b63ffffffff818116838216019080821115613b0b57613b0b6150da565b600081518084526020808501945080840160005b838110156151a157815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161516f565b509495945050505050565b602081526151c360208201835163ffffffff169052565b600060208301516151dc604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015173ffffffffffffffffffffffffffffffffffffffff8116608084015250608083015180151560a08401525060a083015162ffffff811660c08401525060c083015163ffffffff811660e08401525060e083015161010061525b8185018363ffffffff169052565b84015190506101206152848482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151905061014061529b8482018361ffff169052565b84015190506101606152c48482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151610180848101919091528401516101a0808501919091528401516101c0808501919091528401516102006101e08086018290529192509061530c61022086018461515b565b95015173ffffffffffffffffffffffffffffffffffffffff169301929092525090919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203615392576153926150da565b5060010190565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526153c98184018a61515b565b905082810360808401526153dd818961515b565b905060ff871660a084015282810360c08401526153fa818761470b565b905067ffffffffffffffff851660e084015282810361010084015261541f818561470b565b9c9b505050505050505050505050565b60006020828403121561544157600080fd5b5035919050565b6bffffffffffffffffffffffff818116838216019080821115613b0b57613b0b6150da565b8082018082111561224c5761224c6150da565b60ff818116838216019081111561224c5761224c6150da565b80516148c18161494b565b80516148c181614894565b80516148c181614968565b80516148c181614981565b80516148c18161499d565b600082601f8301126154e157600080fd5b815160206154f16148e783614870565b82815260059290921b8401810191818101908684111561551057600080fd5b8286015b8481101561492f57805161552781614894565b8352918301918301615514565b600082601f83011261554557600080fd5b815160206155556148e783614870565b82815260059290921b8401810191818101908684111561557457600080fd5b8286015b8481101561492f57805161558b81614894565b8352918301918301615578565b600082601f8301126155a957600080fd5b815160206155b96148e783614870565b82815260a092830285018201928282019190878511156155d857600080fd5b8387015b85811015614cd75781818a0312156155f45760008081fd5b6155fc6147db565b81516156078161494b565b81528186015161561681614981565b8187015260408281015161562981614894565b908201526060828101519082015260808083015161564681614bfe565b9082015284529284019281016155dc565b60008060006060848603121561566c57600080fd5b835167ffffffffffffffff8082111561568457600080fd5b90850190610200828803121561569957600080fd5b6156a16147b1565b6156aa83615499565b81526156b860208401615499565b60208201526156c960408401615499565b60408201526156da606084016154a4565b60608201526156eb608084016154af565b60808201526156fc60a084016154ba565b60a082015261570d60c08401615499565b60c082015261571e60e08401615499565b60e08201526101006157318185016154a4565b908201526101206157438482016154c5565b908201526101406157558482016154a4565b90820152610160838101519082015261018080840151908201526101a080840151908201526101c0808401518381111561578e57600080fd5b61579a8a8287016154d0565b8284015250506101e06157ae8185016154a4565b9082015260208701519095509150808211156157c957600080fd5b6157d587838801615534565b935060408601519150808211156157eb57600080fd5b506157f886828701615598565b9150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526158788285018b61515b565b9150838203608085015261588c828a61515b565b915060ff881660a085015283820360c08501526158a9828861470b565b90861660e0850152838103610100850152905061541f818561470b565b8281526040602082015260006137f5604083018461470b565b600080604083850312156158f257600080fd5b82516158fd81614968565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f83011261594457600080fd5b813560206159546148e783614870565b82815260059290921b8401810191818101908684111561597357600080fd5b8286015b8481101561492f5780358352918301918301615977565b600082601f83011261599f57600080fd5b813560206159af6148e783614870565b82815260059290921b840181019181810190868411156159ce57600080fd5b8286015b8481101561492f57803567ffffffffffffffff8111156159f25760008081fd5b615a008986838b0101614b0c565b8452509183019183016159d2565b600060208284031215615a2057600080fd5b813567ffffffffffffffff80821115615a3857600080fd5b9083019060c08286031215615a4c57600080fd5b615a546147fe565b8235815260208301356020820152604083013582811115615a7457600080fd5b615a8087828601615933565b604083015250606083013582811115615a9857600080fd5b615aa487828601615933565b606083015250608083013582811115615abc57600080fd5b615ac88782860161598e565b60808301525060a083013582811115615ae057600080fd5b615aec8782860161598e565b60a08301525095945050505050565b61ffff818116838216019080821115613b0b57613b0b6150da565b808202811582820484141761224c5761224c6150da565b8181038181111561224c5761224c6150da565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082615b7e57615b7e615b40565b500490565b6bffffffffffffffffffffffff85168152836020820152826040820152608060608201526000615bb6608083018461470b565b9695505050505050565b6bffffffffffffffffffffffff828116828216039080821115613b0b57613b0b6150da565b60006bffffffffffffffffffffffff80841680615c0457615c04615b40565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114615c3857615c386150da565b505092915050565b805169ffffffffffffffffffff811681146148c157600080fd5b600080600080600060a08688031215615c7257600080fd5b615c7b86615c40565b9450602086015193506040860151925060608601519150615c9e60808701615c40565b90509295509295909350565b6fffffffffffffffffffffffffffffffff818116838216019080821115613b0b57613b0b6150da565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600060408284031215615d1457600080fd5b6040516040810181811067ffffffffffffffff82111715615d3757615d37614782565b6040528251615d458161494b565b81526020928301519281019290925250919050565b600060a08284031215615d6c57600080fd5b615d746147db565b82518152602083015160208201526040830151615d908161494b565b60408201526060830151615da38161494b565b6060820152608092830151928101929092525091905056fea164736f6c6343000813000a",
}

var AutomationRegistryABI = AutomationRegistryMetaData.ABI

var AutomationRegistryBin = AutomationRegistryMetaData.Bin

func DeployAutomationRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, logicA common.Address) (common.Address, *types.Transaction, *AutomationRegistry, error) {
	parsed, err := AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryBin), backend, logicA)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistry{address: address, abi: *parsed, AutomationRegistryCaller: AutomationRegistryCaller{contract: contract}, AutomationRegistryTransactor: AutomationRegistryTransactor{contract: contract}, AutomationRegistryFilterer: AutomationRegistryFilterer{contract: contract}}, nil
}

type AutomationRegistry struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryCaller
	AutomationRegistryTransactor
	AutomationRegistryFilterer
}

type AutomationRegistryCaller struct {
	contract *bind.BoundContract
}

type AutomationRegistryTransactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryFilterer struct {
	contract *bind.BoundContract
}

type AutomationRegistrySession struct {
	Contract     *AutomationRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryCallerSession struct {
	Contract *AutomationRegistryCaller
	CallOpts bind.CallOpts
}

type AutomationRegistryTransactorSession struct {
	Contract     *AutomationRegistryTransactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryRaw struct {
	Contract *AutomationRegistry
}

type AutomationRegistryCallerRaw struct {
	Contract *AutomationRegistryCaller
}

type AutomationRegistryTransactorRaw struct {
	Contract *AutomationRegistryTransactor
}

func NewAutomationRegistry(address common.Address, backend bind.ContractBackend) (*AutomationRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry{address: address, abi: abi, AutomationRegistryCaller: AutomationRegistryCaller{contract: contract}, AutomationRegistryTransactor: AutomationRegistryTransactor{contract: contract}, AutomationRegistryFilterer: AutomationRegistryFilterer{contract: contract}}, nil
}

func NewAutomationRegistryCaller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryCaller, error) {
	contract, err := bindAutomationRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryCaller{contract: contract}, nil
}

func NewAutomationRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryTransactor, error) {
	contract, err := bindAutomationRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryTransactor{contract: contract}, nil
}

func NewAutomationRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryFilterer, error) {
	contract, err := bindAutomationRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFilterer{contract: contract}, nil
}

func bindAutomationRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistry *AutomationRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry.Contract.AutomationRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry *AutomationRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AutomationRegistryTransactor.contract.Transfer(opts)
}

func (_AutomationRegistry *AutomationRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AutomationRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry *AutomationRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry *AutomationRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.contract.Transfer(opts)
}

func (_AutomationRegistry *AutomationRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry *AutomationRegistryCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) FallbackTo() (common.Address, error) {
	return _AutomationRegistry.Contract.FallbackTo(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistry.Contract.FallbackTo(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_AutomationRegistry *AutomationRegistrySession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDetails(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDetails(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_AutomationRegistry *AutomationRegistrySession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) Owner() (common.Address, error) {
	return _AutomationRegistry.Contract.Owner(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) Owner() (common.Address, error) {
	return _AutomationRegistry.Contract.Owner(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) TypeAndVersion() (string, error) {
	return _AutomationRegistry.Contract.TypeAndVersion(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) TypeAndVersion() (string, error) {
	return _AutomationRegistry.Contract.TypeAndVersion(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistry *AutomationRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AcceptOwnership(&_AutomationRegistry.TransactOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AcceptOwnership(&_AutomationRegistry.TransactOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistrySession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.OnTokenTransfer(&_AutomationRegistry.TransactOpts, sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.OnTokenTransfer(&_AutomationRegistry.TransactOpts, sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistrySession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfig(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfig(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry *AutomationRegistrySession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfigTypeSafe(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfigTypeSafe(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_AutomationRegistry *AutomationRegistrySession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SimulatePerformUpkeep(&_AutomationRegistry.TransactOpts, id, performData)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SimulatePerformUpkeep(&_AutomationRegistry.TransactOpts, id, performData)
}

func (_AutomationRegistry *AutomationRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistry *AutomationRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.TransferOwnership(&_AutomationRegistry.TransactOpts, to)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.TransferOwnership(&_AutomationRegistry.TransactOpts, to)
}

func (_AutomationRegistry *AutomationRegistryTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistrySession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Transmit(&_AutomationRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Transmit(&_AutomationRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistry *AutomationRegistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Fallback(&_AutomationRegistry.TransactOpts, calldata)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Fallback(&_AutomationRegistry.TransactOpts, calldata)
}

type AutomationRegistryAdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryAdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryAdminPrivilegeConfigSet)
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

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryAdminPrivilegeConfigSetIterator{contract: _AutomationRegistry.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryAdminPrivilegeConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryAdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryAdminPrivilegeConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryBillingConfigOverriddenIterator struct {
	Event *AutomationRegistryBillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryBillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryBillingConfigOverridden)
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
		it.Event = new(AutomationRegistryBillingConfigOverridden)
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

func (it *AutomationRegistryBillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryBillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryBillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryBillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryBillingConfigOverriddenIterator{contract: _AutomationRegistry.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryBillingConfigOverridden)
				if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryBillingConfigOverridden, error) {
	event := new(AutomationRegistryBillingConfigOverridden)
	if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryBillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistryBillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryBillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryBillingConfigOverrideRemoved)
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
		it.Event = new(AutomationRegistryBillingConfigOverrideRemoved)
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

func (it *AutomationRegistryBillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryBillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryBillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryBillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryBillingConfigOverrideRemovedIterator{contract: _AutomationRegistry.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryBillingConfigOverrideRemoved)
				if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryBillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistryBillingConfigOverrideRemoved)
	if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryBillingConfigSetIterator struct {
	Event *AutomationRegistryBillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryBillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryBillingConfigSet)
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
		it.Event = new(AutomationRegistryBillingConfigSet)
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

func (it *AutomationRegistryBillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryBillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryBillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryBillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryBillingConfigSetIterator{contract: _AutomationRegistry.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryBillingConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistryBillingConfigSet, error) {
	event := new(AutomationRegistryBillingConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryCancelledUpkeepReportIterator struct {
	Event *AutomationRegistryCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryCancelledUpkeepReport)
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
		it.Event = new(AutomationRegistryCancelledUpkeepReport)
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

func (it *AutomationRegistryCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryCancelledUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryCancelledUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryCancelledUpkeepReport, error) {
	event := new(AutomationRegistryCancelledUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistryChainSpecificModuleUpdated)
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

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryChainSpecificModuleUpdatedIterator{contract: _AutomationRegistry.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryChainSpecificModuleUpdated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryChainSpecificModuleUpdated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryConfigSetIterator struct {
	Event *AutomationRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryConfigSet)
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
		it.Event = new(AutomationRegistryConfigSet)
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

func (it *AutomationRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryConfigSet struct {
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

func (_AutomationRegistry *AutomationRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistryConfigSetIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryConfigSetIterator{contract: _AutomationRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseConfigSet(log types.Log) (*AutomationRegistryConfigSet, error) {
	event := new(AutomationRegistryConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryDedupKeyAddedIterator struct {
	Event *AutomationRegistryDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryDedupKeyAdded)
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
		it.Event = new(AutomationRegistryDedupKeyAdded)
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

func (it *AutomationRegistryDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryDedupKeyAddedIterator{contract: _AutomationRegistry.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryDedupKeyAdded)
				if err := _AutomationRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryDedupKeyAdded, error) {
	event := new(AutomationRegistryDedupKeyAdded)
	if err := _AutomationRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryFeesWithdrawnIterator struct {
	Event *AutomationRegistryFeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryFeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryFeesWithdrawn)
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
		it.Event = new(AutomationRegistryFeesWithdrawn)
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

func (it *AutomationRegistryFeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryFeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryFeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryFeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFeesWithdrawnIterator{contract: _AutomationRegistry.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryFeesWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistryFeesWithdrawn, error) {
	event := new(AutomationRegistryFeesWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryFundsAddedIterator struct {
	Event *AutomationRegistryFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryFundsAdded)
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
		it.Event = new(AutomationRegistryFundsAdded)
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

func (it *AutomationRegistryFundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFundsAddedIterator{contract: _AutomationRegistry.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryFundsAdded)
				if err := _AutomationRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseFundsAdded(log types.Log) (*AutomationRegistryFundsAdded, error) {
	event := new(AutomationRegistryFundsAdded)
	if err := _AutomationRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryFundsWithdrawnIterator struct {
	Event *AutomationRegistryFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryFundsWithdrawn)
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
		it.Event = new(AutomationRegistryFundsWithdrawn)
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

func (it *AutomationRegistryFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFundsWithdrawnIterator{contract: _AutomationRegistry.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryFundsWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryFundsWithdrawn, error) {
	event := new(AutomationRegistryFundsWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryInsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryInsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistryInsufficientFundsUpkeepReport)
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

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryInsufficientFundsUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryInsufficientFundsUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryInsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryInsufficientFundsUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryNOPsSettledOffchainIterator struct {
	Event *AutomationRegistryNOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryNOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryNOPsSettledOffchain)
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
		it.Event = new(AutomationRegistryNOPsSettledOffchain)
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

func (it *AutomationRegistryNOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryNOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryNOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryNOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryNOPsSettledOffchainIterator{contract: _AutomationRegistry.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryNOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryNOPsSettledOffchain)
				if err := _AutomationRegistry.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryNOPsSettledOffchain, error) {
	event := new(AutomationRegistryNOPsSettledOffchain)
	if err := _AutomationRegistry.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryOwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryOwnershipTransferRequested)
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
		it.Event = new(AutomationRegistryOwnershipTransferRequested)
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

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryOwnershipTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryOwnershipTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryOwnershipTransferRequested, error) {
	event := new(AutomationRegistryOwnershipTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryOwnershipTransferredIterator struct {
	Event *AutomationRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryOwnershipTransferred)
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
		it.Event = new(AutomationRegistryOwnershipTransferred)
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

func (it *AutomationRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryOwnershipTransferredIterator{contract: _AutomationRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryOwnershipTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryOwnershipTransferred, error) {
	event := new(AutomationRegistryOwnershipTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPausedIterator struct {
	Event *AutomationRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPaused)
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
		it.Event = new(AutomationRegistryPaused)
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

func (it *AutomationRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryPausedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPausedIterator{contract: _AutomationRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePaused(log types.Log) (*AutomationRegistryPaused, error) {
	event := new(AutomationRegistryPaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeesUpdatedIterator struct {
	Event *AutomationRegistryPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeesUpdated)
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
		it.Event = new(AutomationRegistryPayeesUpdated)
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

func (it *AutomationRegistryPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryPayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeesUpdatedIterator{contract: _AutomationRegistry.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeesUpdated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryPayeesUpdated, error) {
	event := new(AutomationRegistryPayeesUpdated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistryPayeeshipTransferRequested)
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

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeeshipTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeeshipTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryPayeeshipTransferRequested, error) {
	event := new(AutomationRegistryPayeeshipTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeeshipTransferredIterator struct {
	Event *AutomationRegistryPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeeshipTransferred)
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
		it.Event = new(AutomationRegistryPayeeshipTransferred)
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

func (it *AutomationRegistryPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeeshipTransferredIterator{contract: _AutomationRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeeshipTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryPayeeshipTransferred, error) {
	event := new(AutomationRegistryPayeeshipTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPaymentWithdrawnIterator struct {
	Event *AutomationRegistryPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPaymentWithdrawn)
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
		it.Event = new(AutomationRegistryPaymentWithdrawn)
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

func (it *AutomationRegistryPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPaymentWithdrawnIterator{contract: _AutomationRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPaymentWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryPaymentWithdrawn, error) {
	event := new(AutomationRegistryPaymentWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistryReorgedUpkeepReport)
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

func (it *AutomationRegistryReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryReorgedUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryReorgedUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryReorgedUpkeepReport, error) {
	event := new(AutomationRegistryReorgedUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryStaleUpkeepReportIterator struct {
	Event *AutomationRegistryStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryStaleUpkeepReport)
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
		it.Event = new(AutomationRegistryStaleUpkeepReport)
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

func (it *AutomationRegistryStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryStaleUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryStaleUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryStaleUpkeepReport, error) {
	event := new(AutomationRegistryStaleUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryTransmittedIterator struct {
	Event *AutomationRegistryTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryTransmitted)
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
		it.Event = new(AutomationRegistryTransmitted)
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

func (it *AutomationRegistryTransmittedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistryTransmittedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryTransmittedIterator{contract: _AutomationRegistry.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistryTransmitted) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryTransmitted)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseTransmitted(log types.Log) (*AutomationRegistryTransmitted, error) {
	event := new(AutomationRegistryTransmitted)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUnpausedIterator struct {
	Event *AutomationRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUnpaused)
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
		it.Event = new(AutomationRegistryUnpaused)
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

func (it *AutomationRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryUnpausedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUnpausedIterator{contract: _AutomationRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUnpaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUnpaused(log types.Log) (*AutomationRegistryUnpaused, error) {
	event := new(AutomationRegistryUnpaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistryUpkeepAdminTransferRequested)
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

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepAdminTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepAdminTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryUpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryUpkeepAdminTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistryUpkeepAdminTransferred)
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

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepAdminTransferredIterator{contract: _AutomationRegistry.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepAdminTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryUpkeepAdminTransferred, error) {
	event := new(AutomationRegistryUpkeepAdminTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepCanceledIterator struct {
	Event *AutomationRegistryUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepCanceled)
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
		it.Event = new(AutomationRegistryUpkeepCanceled)
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

func (it *AutomationRegistryUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepCanceledIterator{contract: _AutomationRegistry.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepCanceled)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryUpkeepCanceled, error) {
	event := new(AutomationRegistryUpkeepCanceled)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistryUpkeepCheckDataSet)
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

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepCheckDataSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepCheckDataSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryUpkeepCheckDataSet, error) {
	event := new(AutomationRegistryUpkeepCheckDataSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistryUpkeepGasLimitSet)
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

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepGasLimitSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepGasLimitSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryUpkeepGasLimitSet, error) {
	event := new(AutomationRegistryUpkeepGasLimitSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepMigratedIterator struct {
	Event *AutomationRegistryUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepMigrated)
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
		it.Event = new(AutomationRegistryUpkeepMigrated)
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

func (it *AutomationRegistryUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepMigratedIterator{contract: _AutomationRegistry.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepMigrated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryUpkeepMigrated, error) {
	event := new(AutomationRegistryUpkeepMigrated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepOffchainConfigSet)
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

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepOffchainConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepOffchainConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryUpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryUpkeepOffchainConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPausedIterator struct {
	Event *AutomationRegistryUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPaused)
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
		it.Event = new(AutomationRegistryUpkeepPaused)
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

func (it *AutomationRegistryUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPausedIterator{contract: _AutomationRegistry.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryUpkeepPaused, error) {
	event := new(AutomationRegistryUpkeepPaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPerformedIterator struct {
	Event *AutomationRegistryUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPerformed)
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
		it.Event = new(AutomationRegistryUpkeepPerformed)
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

func (it *AutomationRegistryUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPerformedIterator{contract: _AutomationRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPerformed)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryUpkeepPerformed, error) {
	event := new(AutomationRegistryUpkeepPerformed)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPrivilegeConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryUpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryUpkeepPrivilegeConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepReceivedIterator struct {
	Event *AutomationRegistryUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepReceived)
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
		it.Event = new(AutomationRegistryUpkeepReceived)
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

func (it *AutomationRegistryUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepReceivedIterator{contract: _AutomationRegistry.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepReceived)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryUpkeepReceived, error) {
	event := new(AutomationRegistryUpkeepReceived)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepRegisteredIterator struct {
	Event *AutomationRegistryUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepRegistered)
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
		it.Event = new(AutomationRegistryUpkeepRegistered)
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

func (it *AutomationRegistryUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepRegisteredIterator{contract: _AutomationRegistry.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepRegistered)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryUpkeepRegistered, error) {
	event := new(AutomationRegistryUpkeepRegistered)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepTriggerConfigSet)
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

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepTriggerConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepTriggerConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryUpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryUpkeepTriggerConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepUnpausedIterator struct {
	Event *AutomationRegistryUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepUnpaused)
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
		it.Event = new(AutomationRegistryUpkeepUnpaused)
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

func (it *AutomationRegistryUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepUnpausedIterator{contract: _AutomationRegistry.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepUnpaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryUpkeepUnpaused, error) {
	event := new(AutomationRegistryUpkeepUnpaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistry.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistry.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistry.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistry.ParseBillingConfigOverridden(log)
	case _AutomationRegistry.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistry.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistry.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistry.ParseBillingConfigSet(log)
	case _AutomationRegistry.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistry.ParseCancelledUpkeepReport(log)
	case _AutomationRegistry.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistry.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistry.abi.Events["ConfigSet"].ID:
		return _AutomationRegistry.ParseConfigSet(log)
	case _AutomationRegistry.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistry.ParseDedupKeyAdded(log)
	case _AutomationRegistry.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistry.ParseFeesWithdrawn(log)
	case _AutomationRegistry.abi.Events["FundsAdded"].ID:
		return _AutomationRegistry.ParseFundsAdded(log)
	case _AutomationRegistry.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistry.ParseFundsWithdrawn(log)
	case _AutomationRegistry.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistry.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistry.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistry.ParseNOPsSettledOffchain(log)
	case _AutomationRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistry.ParseOwnershipTransferRequested(log)
	case _AutomationRegistry.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistry.ParseOwnershipTransferred(log)
	case _AutomationRegistry.abi.Events["Paused"].ID:
		return _AutomationRegistry.ParsePaused(log)
	case _AutomationRegistry.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistry.ParsePayeesUpdated(log)
	case _AutomationRegistry.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistry.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistry.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistry.ParsePayeeshipTransferred(log)
	case _AutomationRegistry.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistry.ParsePaymentWithdrawn(log)
	case _AutomationRegistry.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistry.ParseReorgedUpkeepReport(log)
	case _AutomationRegistry.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistry.ParseStaleUpkeepReport(log)
	case _AutomationRegistry.abi.Events["Transmitted"].ID:
		return _AutomationRegistry.ParseTransmitted(log)
	case _AutomationRegistry.abi.Events["Unpaused"].ID:
		return _AutomationRegistry.ParseUnpaused(log)
	case _AutomationRegistry.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistry.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistry.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistry.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistry.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistry.ParseUpkeepCanceled(log)
	case _AutomationRegistry.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistry.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistry.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistry.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistry.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistry.ParseUpkeepMigrated(log)
	case _AutomationRegistry.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistry.ParseUpkeepPaused(log)
	case _AutomationRegistry.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistry.ParseUpkeepPerformed(log)
	case _AutomationRegistry.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistry.ParseUpkeepReceived(log)
	case _AutomationRegistry.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistry.ParseUpkeepRegistered(log)
	case _AutomationRegistry.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistry.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryBillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistryBillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistryBillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0x720a5849025dc4fd0061aed1bb30efd713cde64ce7f8d807953ecca27c8f143c")
}

func (AutomationRegistryCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (AutomationRegistryDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryFeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryNOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (AutomationRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistry *AutomationRegistry) Address() common.Address {
	return _AutomationRegistry.address
}

type AutomationRegistryInterface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryAdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryBillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryBillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryBillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryBillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryBillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryBillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistryBillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryCancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryChainSpecificModuleUpdated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*AutomationRegistryConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryDedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryFeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistryFeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryInsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryNOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryNOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryNOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistryTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistryTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*AutomationRegistryTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
