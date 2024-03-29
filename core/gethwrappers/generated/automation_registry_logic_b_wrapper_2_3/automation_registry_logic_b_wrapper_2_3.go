// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_registry_logic_b_wrapper_2_3

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

var AutomationRegistryLogicBMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicC2_3\",\"name\":\"logicC\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"overrides\",\"type\":\"tuple\"}],\"name\":\"BillingConfigOverridden\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"BillingConfigOverrideRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payments\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"removeBillingOverrides\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"billingOverrides\",\"type\":\"tuple\"}],\"name\":\"setBillingOverrides\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20Fees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101806040523480156200001257600080fd5b50604051620055bb380380620055bb83398101604081905262000035916200062f565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b91906200062f565b826001600160a01b031663226cf83c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200062f565b836001600160a01b031663614486af6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200062f565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200062f565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200062f565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200062f565b876001600160a01b031663c5b964e06040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002d3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002f9919062000656565b886001600160a01b031663ac4dc59a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000338573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200035e91906200062f565b3380600081620003b55760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620003e857620003e8816200056b565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff19166001838181111562000439576200043962000679565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200049a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004c091906200068f565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000504573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200052a91906200068f565b60ff16146200054c576040516301f86e1760e41b815260040160405180910390fd5b5050506001600160a01b039095166101605250620006b4945050505050565b336001600160a01b03821603620005c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620003ac565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200062c57600080fd5b50565b6000602082840312156200064257600080fd5b81516200064f8162000616565b9392505050565b6000602082840312156200066957600080fd5b8151600281106200064f57600080fd5b634e487b7160e01b600052602160045260246000fd5b600060208284031215620006a257600080fd5b815160ff811681146200064f57600080fd5b60805160a05160c05160e05161010051610120516101405161016051614e8b62000730600039600081816101770152610224015260005050600061276f0152600050506000612a5e015260006136cc01526000612b38015260008181610c2901528181610cea01528181610db501526128300152614e8b6000f3fe6080604052600436106101755760003560e01c80638765ecbe116100cb578063aed2e9291161007f578063ce7dc5b411610059578063ce7dc5b414610493578063f2fde38b146104b3578063f7d334ba146104d357610175565b8063aed2e9291461041c578063b148ab6b14610453578063cd7f71b51461047357610175565b80638dcf0fe7116100b05780638dcf0fe7146103bc578063a72aa27e146103dc578063a86e1781146103fc57610175565b80638765ecbe146103715780638da5cb5b1461039157610175565b806354b7faae1161012d57806371fae17f1161010757806371fae17f1461031c578063744bfe611461033c57806379ba50971461035c57610175565b806354b7faae146102a957806368d369d8146102c957806371791aa0146102e957610175565b8063349e8cca1161015e578063349e8cca146102155780634ee88d35146102695780635165f2f51461028957610175565b80631a2af011146101bc57806329c5efad146101dc575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156101b5573d6000f35b3d6000fd5b005b3480156101c857600080fd5b506101ba6101d7366004613e95565b6104f3565b3480156101e857600080fd5b506101fc6101f7366004614009565b6105f9565b60405161020c9493929190614128565b60405180910390f35b34801561022157600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161020c565b34801561027557600080fd5b506101ba6102843660046141aa565b6108d7565b34801561029557600080fd5b506101ba6102a43660046141f6565b610939565b3480156102b557600080fd5b506101ba6102c436600461420f565b610aeb565b3480156102d557600080fd5b506101ba6102e436600461423b565b610d5e565b3480156102f557600080fd5b50610309610304366004614009565b610fa8565b60405161020c979695949392919061427c565b34801561032857600080fd5b506101ba6103373660046141f6565b61170c565b34801561034857600080fd5b506101ba610357366004613e95565b6117a2565b34801561036857600080fd5b506101ba611c1a565b34801561037d57600080fd5b506101ba61038c3660046141f6565b611d17565b34801561039d57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610244565b3480156103c857600080fd5b506101ba6103d73660046141aa565b611ecc565b3480156103e857600080fd5b506101ba6103f73660046142dc565b611f21565b34801561040857600080fd5b506101ba610417366004614301565b612020565b34801561042857600080fd5b5061043c6104373660046141aa565b612101565b60408051921515835260208301919091520161020c565b34801561045f57600080fd5b506101ba61046e3660046141f6565b6122ac565b34801561047f57600080fd5b506101ba61048e3660046141aa565b6124d9565b34801561049f57600080fd5b506101fc6104ae36600461437b565b612590565b3480156104bf57600080fd5b506101ba6104ce36600461445d565b612652565b3480156104df57600080fd5b506103096104ee3660046141f6565b612666565b6104fc826126a2565b3373ffffffffffffffffffffffffffffffffffffffff82160361054b576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116146105f55760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b60006060600080610608612757565b600086815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615159483019490945263ffffffff620100008204811695830195909552660100000000000081048516606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490941660e0820152600290910154909216908201525a9150600080826080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561075e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610782919061448a565b73ffffffffffffffffffffffffffffffffffffffff16601660000160149054906101000a900463ffffffff1663ffffffff16896040516107c291906144a7565b60006040518083038160008787f1925050503d8060008114610800576040519150601f19603f3d011682016040523d82523d6000602084013e610805565b606091505b50915091505a61081590856144f2565b93508161083e5760006040518060200160405280600081525060079650965096505050506108ce565b80806020019051810190610852919061455a565b90975095508661087e5760006040518060200160405280600081525060049650965096505050506108ce565b60185486517401000000000000000000000000000000000000000090910463ffffffff1610156108ca5760006040518060200160405280600081525060059650965096505050506108ce565b5050505b92959194509250565b6108e0836126a2565b6000838152601d602052604090206108f982848361463f565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664838360405161092c9291906147a3565b60405180910390a2505050565b610942816126a2565b600081815260046020908152604091829020825161012081018452815460ff808216151580845261010080840490921615159584019590955262010000820463ffffffff9081169684019690965266010000000000008204861660608401526a010000000000000000000090910473ffffffffffffffffffffffffffffffffffffffff908116608084015260018401546fffffffffffffffffffffffffffffffff811660a085015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08501527c0100000000000000000000000000000000000000000000000000000000900490951660e083015260029092015490931690830152610a7c576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055610abb6002836127c8565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b610af36127dd565b73ffffffffffffffffffffffffffffffffffffffff8216610b40576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610b4a61282e565b90506000811215610b96576040517fcf47918100000000000000000000000000000000000000000000000000000000815260006004820152602481018390526044015b60405180910390fd5b80821115610bda576040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610b8d565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490526000917f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015610c74573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c9891906147b7565b905080610cd1576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8373ffffffffffffffffffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa885604051610d5091815260200190565b60405180910390a350505050565b610d666127dd565b73ffffffffffffffffffffffffffffffffffffffff8216610db3576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610e38576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015610eb4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ed891906147d2565b610ee291906144f2565b905080821115610f28576040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610b8d565b610f4973ffffffffffffffffffffffffffffffffffffffff851684846128fd565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa884604051610d5091815260200190565b600060606000806000806000610fbc612757565b6000610fc78a61298f565b905060006014604051806101200160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900462ffffff1662ffffff1662ffffff1681526020016000820160139054906101000a900461ffff1661ffff1661ffff1681526020016000820160159054906101000a900460ff1660ff1660ff1681526020016000820160169054906101000a900460ff161515151581526020016000820160179054906101000a900460ff161515151581526020016000820160189054906101000a900460ff161515151581526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d8152602001908152602001600020604051806101200160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900460ff161515151581526020016000820160029054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160069054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200160008201600a9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020016001820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000808360a001511561138a5750506040805160208101825260008082529290910151919a5098506009975089965063ffffffff169450859350839250611700915050565b606083015163ffffffff908116146113d45750506040805160208101825260008082529290910151919a5098506001975089965063ffffffff169450859350839250611700915050565b8251156114135750506040805160208101825260008082529290910151919a5098506002975089965063ffffffff169450859350839250611700915050565b61141c84612a3a565b8094508198508299505050506114418e858786604001518b8b888a6101000151612c2c565b9050806bffffffffffffffffffffffff168360c001516bffffffffffffffffffffffff1610156114a35750506040805160208101825260008082529290910151919a5098506006975089965063ffffffff169450859350839250611700915050565b505060006114b28d858e612f9c565b90505a9750600080836080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611509573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061152d919061448a565b73ffffffffffffffffffffffffffffffffffffffff16601660000160149054906101000a900463ffffffff1663ffffffff168460405161156d91906144a7565b60006040518083038160008787f1925050503d80600081146115ab576040519150601f19603f3d011682016040523d82523d6000602084013e6115b0565b606091505b50915091505a6115c0908b6144f2565b995081611646576018548151780100000000000000000000000000000000000000000000000090910463ffffffff1610156116255750506040805160208101825260008082529390910151929b509950600898505063ffffffff169450611700915050565b604090930151929a50600399505063ffffffff909116955061170092505050565b8080602001905181019061165a919061455a565b909d509b508c6116945750506040805160208101825260008082529390910151929b509950600498505063ffffffff169450611700915050565b6018548c517401000000000000000000000000000000000000000090910463ffffffff1610156116ee5750506040805160208101825260008082529390910151929b509950600598505063ffffffff169450611700915050565b5050506040015163ffffffff16945050505b92959891949750929550565b61171461317c565b600081815260046020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055602390915280822080547fffffffffffffffffffffffffffffffffffffffffffffffffff000000000000001690555182917f97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a91a250565b60145477010000000000000000000000000000000000000000000000900460ff16156117fa576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff167701000000000000000000000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116611889576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615158387015263ffffffff620100008304811684870152660100000000000083048116606085015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009093048316608085015260018501546fffffffffffffffffffffffffffffffff811660a08601526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08601527c010000000000000000000000000000000000000000000000000000000090041660e0840152600290930154811692820192909252868552600590935292205490911633146119c9576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554604080517f57e871e7000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff909216916357e871e7916004808201926020929091908290030181865afa158015611a39573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a5d91906147d2565b816060015163ffffffff161115611aa0576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526004602090815260408083206001015461010085015173ffffffffffffffffffffffffffffffffffffffff1684526021909252909120547001000000000000000000000000000000009091046bffffffffffffffffffffffff1690611b0b9082906144f2565b6101008301805173ffffffffffffffffffffffffffffffffffffffff908116600090815260216020908152604080832095909555888252600490529290922060010180547fffffffff000000000000000000000000ffffffffffffffffffffffffffffffff16905551611b8e9116846bffffffffffffffffffffffff84166128fd565b604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611c9b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610b8d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611d20816126a2565b600081815260046020908152604091829020825161012081018452815460ff808216158015845261010080840490921615159584019590955262010000820463ffffffff9081169684019690965266010000000000008204861660608401526a010000000000000000000090910473ffffffffffffffffffffffffffffffffffffffff908116608084015260018401546fffffffffffffffffffffffffffffffff811660a085015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08501527c0100000000000000000000000000000000000000000000000000000000900490951660e083015260029092015490931690830152611e5a576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055611e9c6002836131cd565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b611ed5836126a2565b6000838152601e60205260409020611eee82848361463f565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850838360405161092c9291906147a3565b6108fc8163ffffffff161080611f5e575060165463ffffffff78010000000000000000000000000000000000000000000000009091048116908216115b15611f95576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611f9e826126a2565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffff166201000063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c91015b60405180910390a25050565b61202861317c565b6000828152600460205260409020546601000000000000900463ffffffff90811614612080576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790556023909152902081906120ce82826147fc565b905050817fd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340826040516120149190614883565b60008061210c612757565b601454760100000000000000000000000000000000000000000000900460ff1615612163576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615158386015262010000820463ffffffff90811684880181905266010000000000008404821660608601526a010000000000000000000090930473ffffffffffffffffffffffffffffffffffffffff9081166080860181905260018701546fffffffffffffffffffffffffffffffff811660a088015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08801527c0100000000000000000000000000000000000000000000000000000000900490921660e0860152600290950154909416908301528451601f8901859004850281018501909552878552909361229f9389908990819084018382808284376000920191909152506131d992505050565b9097909650945050505050565b600081815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615159483019490945263ffffffff6201000082048116958301959095526601000000000000810485166060830181905273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009092048216608084015260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c01000000000000000000000000000000000000000000000000000000009004861660e084015260029093015416928101929092529091146123e8576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612445576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b6124e2836126a2565b6017547c0100000000000000000000000000000000000000000000000000000000900463ffffffff16811115612544576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083815260076020526040902061255d82848361463f565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d838360405161092c9291906147a3565b600060606000806000634b56a42e60e01b8888886040516024016125b6939291906148ba565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905061263f89826105f9565b929c919b50995090975095505050505050565b61265a6133f4565b61266381613475565b50565b60006060600080600080600061268b8860405180602001604052806000815250610fa8565b959e949d50929b5090995097509550909350915050565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146126ff576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020546601000000000000900463ffffffff90811614612663576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146127c6576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60006127d4838361356a565b90505b92915050565b60185473ffffffffffffffffffffffffffffffffffffffff1633146127c6576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa1580156128ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128ee91906147d2565b6128f89190614950565b905090565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb0000000000000000000000000000000000000000000000000000000017905261298a9084906135b9565b505050565b6000818160045b600f811015612a1c577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129d4576129d4614970565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a0a57506000949350505050565b80612a148161499f565b915050612996565b5081600f1a6001811115612a3257612a326140be565b949350505050565b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015612ac7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612aeb91906149f1565b5094509092505050600081131580612b0257508142105b80612b235750828015612b235750612b1a82426144f2565b8463ffffffff16105b15612b32576019549650612b36565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015612ba1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612bc591906149f1565b5094509092505050600081131580612bdc57508142105b80612bfd5750828015612bfd5750612bf482426144f2565b8463ffffffff16105b15612c0c57601a549550612c10565b8095505b8686612c1b8a6136c5565b965096509650505050509193909250565b6000808080896001811115612c4357612c436140be565b03612c52575062016b48612ca7565b6001896001811115612c6657612c666140be565b03612c7557506201ccf0612ca7565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a608001516001612cba9190614a41565b612cc89060ff166040614a5a565b601854612cf6906103a49074010000000000000000000000000000000000000000900463ffffffff16614a71565b612d009190614a71565b601554604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa158015612d71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d959190614a84565b90925090508183612da7836018614a71565b612db19190614a5a565b60808f0151612dc1906001614a41565b612dd09060ff166115e0614a5a565b612dda9190614a71565b612de49190614a71565b612dee9085614a71565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa158015612e61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e8591906147d2565b8d6060015161ffff16612e989190614a5a565b94505050506000612ea98b866137b6565b60008d815260046020526040902054909150610100900460ff1615612f0e5760008c81526023602090815260409182902082518084019093525463ffffffff811680845262ffffff640100000000909204821693830193845284529151909116908201525b6000612f788c6040518061012001604052808d63ffffffff1681526020018681526020018781526020018c81526020018b81526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff16815260200185815260200160001515815250613907565b60208101518151919250612f8b91614aa8565b9d9c50505050505050505050505050565b60606000836001811115612fb257612fb26140be565b0361307b576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091612ff791602401614b68565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050613175565b600183600181111561308f5761308f6140be565b03612c75576000828060200190518101906130aa9190614be1565b6000868152600760205260409081902090519192507f40691db400000000000000000000000000000000000000000000000000000000916130ef918491602401614cf1565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915291506131759050565b9392505050565b60175473ffffffffffffffffffffffffffffffffffffffff1633146127c6576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006127d48383613ad7565b601454600090819077010000000000000000000000000000000000000000000000900460ff1615613236576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16770100000000000000000000000000000000000000000000001790556040517f4585e33b00000000000000000000000000000000000000000000000000000000906132ab908590602401614dbc565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061337e9087908790600401614dcf565b60408051808303816000875af115801561339c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133c09190614de8565b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146127c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b8d565b3373ffffffffffffffffffffffffffffffffffffffff8216036134f4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b8d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008181526001830160205260408120546135b1575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556127d7565b5060006127d7565b600061361b826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16613bca9092919063ffffffff16565b80519091501561298a578080602001905181019061363991906147b7565b61298a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610b8d565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613735573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061375991906149f1565b5093505092505060008213158061376f57508042105b8061379f57506000846040015162ffffff1611801561379f575061379381426144f2565b846040015162ffffff16105b156137af575050601b5492915050565b5092915050565b604080516060810182526000808252602080830182815283850183905273ffffffffffffffffffffffffffffffffffffffff868116845260229092528483208054640100000000810462ffffff1690925263ffffffff8216855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495909484936701000000000000009093049092169163feaf968c9160048082019260a0929091908290030181865afa158015613878573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061389c91906149f1565b509350509250506000821315806138b257508042105b806138e257506000866040015162ffffff161180156138e257506138d681426144f2565b866040015162ffffff16105b156138f657600183015460408501526138fe565b604084018290525b50505092915050565b6040805160808101825260008082526020820181905291810182905260608101919091526000836060015161ffff1683606001516139459190614a5a565b905082610100015180156139585750803a105b1561396057503a5b60008360a001518460400151856020015186600001516139809190614a71565b61398a9085614a5a565b6139949190614a71565b61399e9190614a5a565b90506139bc8460e0015160400151826139b79190614e14565b613bd9565b6bffffffffffffffffffffffff16835260808401516139df906139b79083614e14565b6bffffffffffffffffffffffff16604084015260e084015160200151600090613a169062ffffff16683635c9adc5dea00000614a5a565b9050600081633b9aca008760a001518860e001516000015163ffffffff1689604001518a6000015189613a499190614a5a565b613a539190614a71565b613a5d9190614a5a565b613a679190614a5a565b613a719190614e14565b613a7b9190614a71565b9050613a948660e0015160400151826139b79190614e14565b6bffffffffffffffffffffffff1660208601526080860151613aba906139b79083614e14565b6bffffffffffffffffffffffff1660608601525050505092915050565b60008181526001830160205260408120548015613bc0576000613afb6001836144f2565b8554909150600090613b0f906001906144f2565b9050818114613b74576000866000018281548110613b2f57613b2f614970565b9060005260206000200154905080876000018481548110613b5257613b52614970565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613b8557613b85614e4f565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506127d7565b60009150506127d7565b6060612a328484600085613c7b565b60006bffffffffffffffffffffffff821115613c77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610b8d565b5090565b606082471015613d0d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610b8d565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051613d3691906144a7565b60006040518083038185875af1925050503d8060008114613d73576040519150601f19603f3d011682016040523d82523d6000602084013e613d78565b606091505b5091509150613d8987838387613d94565b979650505050505050565b60608315613e2a578251600003613e235773ffffffffffffffffffffffffffffffffffffffff85163b613e23576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610b8d565b5081612a32565b612a328383815115613e3f5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b8d9190614dbc565b73ffffffffffffffffffffffffffffffffffffffff8116811461266357600080fd5b60008060408385031215613ea857600080fd5b823591506020830135613eba81613e73565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff81118282101715613f1857613f18613ec5565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613f6557613f65613ec5565b604052919050565b600067ffffffffffffffff821115613f8757613f87613ec5565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613fc457600080fd5b8135613fd7613fd282613f6d565b613f1e565b818152846020838601011115613fec57600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561401c57600080fd5b82359150602083013567ffffffffffffffff81111561403a57600080fd5b61404685828601613fb3565b9150509250929050565b60005b8381101561406b578181015183820152602001614053565b50506000910152565b6000815180845261408c816020860160208601614050565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110614124577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b84151581526080602082015260006141436080830186614074565b905061415260408301856140ed565b82606083015295945050505050565b60008083601f84011261417357600080fd5b50813567ffffffffffffffff81111561418b57600080fd5b6020830191508360208285010111156141a357600080fd5b9250929050565b6000806000604084860312156141bf57600080fd5b83359250602084013567ffffffffffffffff8111156141dd57600080fd5b6141e986828701614161565b9497909650939450505050565b60006020828403121561420857600080fd5b5035919050565b6000806040838503121561422257600080fd5b823561422d81613e73565b946020939093013593505050565b60008060006060848603121561425057600080fd5b833561425b81613e73565b9250602084013561426b81613e73565b929592945050506040919091013590565b871515815260e06020820152600061429760e0830189614074565b90506142a660408301886140ed565b8560608301528460808301528360a08301528260c083015298975050505050505050565b63ffffffff8116811461266357600080fd5b600080604083850312156142ef57600080fd5b823591506020830135613eba816142ca565b600080828403606081121561431557600080fd5b8335925060407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08201121561434957600080fd5b506020830190509250929050565b600067ffffffffffffffff82111561437157614371613ec5565b5060051b60200190565b6000806000806060858703121561439157600080fd5b8435935060208086013567ffffffffffffffff808211156143b157600080fd5b818801915088601f8301126143c557600080fd5b81356143d3613fd282614357565b81815260059190911b8301840190848101908b8311156143f257600080fd5b8585015b8381101561442a5780358581111561440e5760008081fd5b61441c8e89838a0101613fb3565b8452509186019186016143f6565b5097505050604088013592508083111561444357600080fd5b505061445187828801614161565b95989497509550505050565b60006020828403121561446f57600080fd5b813561317581613e73565b805161448581613e73565b919050565b60006020828403121561449c57600080fd5b815161317581613e73565b600082516144b9818460208701614050565b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156127d7576127d76144c3565b8051801515811461448557600080fd5b600082601f83011261452657600080fd5b8151614534613fd282613f6d565b81815284602083860101111561454957600080fd5b612a32826020830160208701614050565b6000806040838503121561456d57600080fd5b61457683614505565b9150602083015167ffffffffffffffff81111561459257600080fd5b61404685828601614515565b600181811c908216806145b257607f821691505b6020821081036145eb577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561298a57600081815260208120601f850160051c810160208610156146185750805b601f850160051c820191505b8181101561463757828155600101614624565b505050505050565b67ffffffffffffffff83111561465757614657613ec5565b61466b83614665835461459e565b836145f1565b6000601f8411600181146146bd57600085156146875750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614753565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561470c57868501358255602094850194600190920191016146ec565b5086821015614747577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000612a3260208301848661475a565b6000602082840312156147c957600080fd5b6127d482614505565b6000602082840312156147e457600080fd5b5051919050565b62ffffff8116811461266357600080fd5b8135614807816142ca565b63ffffffff811690508154817fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000082161783556020840135614847816147eb565b66ffffff000000008160201b16837fffffffffffffffffffffffffffffffffffffffffffffffffff000000000000008416171784555050505050565b604081018235614892816142ca565b63ffffffff16825260208301356148a8816147eb565b62ffffff811660208401525092915050565b6000604082016040835280865180835260608501915060608160051b8601019250602080890160005b8381101561492f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261491d868351614074565b955093820193908201906001016148e3565b50508584038187015250505061494681858761475a565b9695505050505050565b81810360008312801583831316838312821617156137af576137af6144c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036149d0576149d06144c3565b5060010190565b805169ffffffffffffffffffff8116811461448557600080fd5b600080600080600060a08688031215614a0957600080fd5b614a12866149d7565b9450602086015193506040860151925060608601519150614a35608087016149d7565b90509295509295909350565b60ff81811683821601908111156127d7576127d76144c3565b80820281158282048414176127d7576127d76144c3565b808201808211156127d7576127d76144c3565b60008060408385031215614a9757600080fd5b505080516020909101519092909150565b6bffffffffffffffffffffffff8181168382160190808211156137af576137af6144c3565b60008154614ada8161459e565b808552602060018381168015614af75760018114614b2f57614b5d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550614b5d565b866000528260002060005b85811015614b555781548a8201860152908301908401614b3a565b890184019650505b505050505092915050565b6020815260006127d46020830184614acd565b600082601f830112614b8c57600080fd5b81516020614b9c613fd283614357565b82815260059290921b84018101918181019086841115614bbb57600080fd5b8286015b84811015614bd65780518352918301918301614bbf565b509695505050505050565b600060208284031215614bf357600080fd5b815167ffffffffffffffff80821115614c0b57600080fd5b908301906101008286031215614c2057600080fd5b614c28613ef4565b8251815260208301516020820152604083015160408201526060830151606082015260808301516080820152614c6060a0840161447a565b60a082015260c083015182811115614c7757600080fd5b614c8387828601614b7b565b60c08301525060e083015182811115614c9b57600080fd5b614ca787828601614515565b60e08301525095945050505050565b600081518084526020808501945080840160005b83811015614ce657815187529582019590820190600101614cca565b509495945050505050565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c0840151610100808185015250614d62610140840182614cb6565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc084830301610120850152614d9e8282614074565b9150508281036020840152614db38185614acd565b95945050505050565b6020815260006127d46020830184614074565b828152604060208201526000612a326040830184614074565b60008060408385031215614dfb57600080fd5b614e0483614505565b9150602083015190509250929050565b600082614e4a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var AutomationRegistryLogicBABI = AutomationRegistryLogicBMetaData.ABI

var AutomationRegistryLogicBBin = AutomationRegistryLogicBMetaData.Bin

func DeployAutomationRegistryLogicB(auth *bind.TransactOpts, backend bind.ContractBackend, logicC common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicB, error) {
	parsed, err := AutomationRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicBBin), backend, logicC)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicB{address: address, abi: *parsed, AutomationRegistryLogicBCaller: AutomationRegistryLogicBCaller{contract: contract}, AutomationRegistryLogicBTransactor: AutomationRegistryLogicBTransactor{contract: contract}, AutomationRegistryLogicBFilterer: AutomationRegistryLogicBFilterer{contract: contract}}, nil
}

type AutomationRegistryLogicB struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicBCaller
	AutomationRegistryLogicBTransactor
	AutomationRegistryLogicBFilterer
}

type AutomationRegistryLogicBCaller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicBTransactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicBFilterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicBSession struct {
	Contract     *AutomationRegistryLogicB
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicBCallerSession struct {
	Contract *AutomationRegistryLogicBCaller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicBTransactorSession struct {
	Contract     *AutomationRegistryLogicBTransactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicBRaw struct {
	Contract *AutomationRegistryLogicB
}

type AutomationRegistryLogicBCallerRaw struct {
	Contract *AutomationRegistryLogicBCaller
}

type AutomationRegistryLogicBTransactorRaw struct {
	Contract *AutomationRegistryLogicBTransactor
}

func NewAutomationRegistryLogicB(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicB, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicBABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicB(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB{address: address, abi: abi, AutomationRegistryLogicBCaller: AutomationRegistryLogicBCaller{contract: contract}, AutomationRegistryLogicBTransactor: AutomationRegistryLogicBTransactor{contract: contract}, AutomationRegistryLogicBFilterer: AutomationRegistryLogicBFilterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicBCaller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicBCaller, error) {
	contract, err := bindAutomationRegistryLogicB(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBCaller{contract: contract}, nil
}

func NewAutomationRegistryLogicBTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicBTransactor, error) {
	contract, err := bindAutomationRegistryLogicB(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBTransactor{contract: contract}, nil
}

func NewAutomationRegistryLogicBFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicBFilterer, error) {
	contract, err := bindAutomationRegistryLogicB(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBFilterer{contract: contract}, nil
}

func bindAutomationRegistryLogicB(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB.Contract.AutomationRegistryLogicBCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AutomationRegistryLogicBTransactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AutomationRegistryLogicBTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.FallbackTo(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.FallbackTo(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.Owner(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.Owner(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptOwnership(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptOwnership(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckCallback(&_AutomationRegistryLogicB.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckCallback(&_AutomationRegistryLogicB.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckUpkeep(&_AutomationRegistryLogicB.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckUpkeep(&_AutomationRegistryLogicB.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "checkUpkeep0", id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckUpkeep0(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.CheckUpkeep0(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "executeCallback", id, payload)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.ExecuteCallback(&_AutomationRegistryLogicB.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.ExecuteCallback(&_AutomationRegistryLogicB.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "pauseUpkeep", id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.PauseUpkeep(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.PauseUpkeep(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) RemoveBillingOverrides(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "removeBillingOverrides", id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.RemoveBillingOverrides(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.RemoveBillingOverrides(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetBillingOverrides(opts *bind.TransactOpts, id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setBillingOverrides", id, billingOverrides)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetBillingOverrides(&_AutomationRegistryLogicB.TransactOpts, id, billingOverrides)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetBillingOverrides(&_AutomationRegistryLogicB.TransactOpts, id, billingOverrides)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicB.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicB.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SimulatePerformUpkeep(&_AutomationRegistryLogicB.TransactOpts, id, performData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SimulatePerformUpkeep(&_AutomationRegistryLogicB.TransactOpts, id, performData)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferOwnership(&_AutomationRegistryLogicB.TransactOpts, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferOwnership(&_AutomationRegistryLogicB.TransactOpts, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB.TransactOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawERC20Fees(opts *bind.TransactOpts, asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawERC20Fees", asset, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB.TransactOpts, asset, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB.TransactOpts, asset, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawFunds(&_AutomationRegistryLogicB.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawFunds(&_AutomationRegistryLogicB.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawLink(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawLink", to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawLink(&_AutomationRegistryLogicB.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawLink(&_AutomationRegistryLogicB.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Fallback(&_AutomationRegistryLogicB.TransactOpts, calldata)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Fallback(&_AutomationRegistryLogicB.TransactOpts, calldata)
}

type AutomationRegistryLogicBAdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicBAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBAdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBAdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicBAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBAdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBAdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicBAdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicBAdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBBillingConfigOverriddenIterator struct {
	Event *AutomationRegistryLogicBBillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBBillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBBillingConfigOverridden)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBBillingConfigOverridden)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBBillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBBillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBBillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBBillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBBillingConfigOverriddenIterator{contract: _AutomationRegistryLogicB.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBBillingConfigOverridden)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicBBillingConfigOverridden, error) {
	event := new(AutomationRegistryLogicBBillingConfigOverridden)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBBillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistryLogicBBillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBBillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBBillingConfigOverrideRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBBillingConfigOverrideRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBBillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBBillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBBillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBBillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBBillingConfigOverrideRemovedIterator{contract: _AutomationRegistryLogicB.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBBillingConfigOverrideRemoved)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicBBillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistryLogicBBillingConfigOverrideRemoved)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBBillingConfigSetIterator struct {
	Event *AutomationRegistryLogicBBillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBBillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBBillingConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBBillingConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBBillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBBillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBBillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicBBillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBBillingConfigSetIterator{contract: _AutomationRegistryLogicB.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBBillingConfigSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicBBillingConfigSet, error) {
	event := new(AutomationRegistryLogicBBillingConfigSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBCancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicBCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBCancelledUpkeepReportIterator{contract: _AutomationRegistryLogicB.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBCancelledUpkeepReport)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicBCancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicBCancelledUpkeepReport)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicBChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicBChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicB.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicBChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicBChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBDedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicBDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBDedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBDedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicBDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBDedupKeyAddedIterator{contract: _AutomationRegistryLogicB.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBDedupKeyAdded)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicBDedupKeyAdded, error) {
	event := new(AutomationRegistryLogicBDedupKeyAdded)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBFeesWithdrawnIterator struct {
	Event *AutomationRegistryLogicBFeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBFeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBFeesWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBFeesWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBFeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBFeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBFeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicBFeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBFeesWithdrawnIterator{contract: _AutomationRegistryLogicB.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBFeesWithdrawn)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicBFeesWithdrawn, error) {
	event := new(AutomationRegistryLogicBFeesWithdrawn)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBFundsAddedIterator struct {
	Event *AutomationRegistryLogicBFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBFundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicBFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBFundsAddedIterator{contract: _AutomationRegistryLogicB.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBFundsAdded)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicBFundsAdded, error) {
	event := new(AutomationRegistryLogicBFundsAdded)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicBFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBFundsWithdrawnIterator{contract: _AutomationRegistryLogicB.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBFundsWithdrawn)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicBFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicBFundsWithdrawn)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicBInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicB.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBInsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicBInsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicBInsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBNOPsSettledOffchainIterator struct {
	Event *AutomationRegistryLogicBNOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBNOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBNOPsSettledOffchain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBNOPsSettledOffchain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBNOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBNOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBNOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicBNOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBNOPsSettledOffchainIterator{contract: _AutomationRegistryLogicB.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBNOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBNOPsSettledOffchain)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicBNOPsSettledOffchain, error) {
	event := new(AutomationRegistryLogicBNOPsSettledOffchain)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBOwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicBOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicBOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBOwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicB.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBOwnershipTransferRequested)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicBOwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicBOwnershipTransferRequested)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBOwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicBOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicBOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBOwnershipTransferredIterator{contract: _AutomationRegistryLogicB.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBOwnershipTransferred)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicBOwnershipTransferred, error) {
	event := new(AutomationRegistryLogicBOwnershipTransferred)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBPausedIterator struct {
	Event *AutomationRegistryLogicBPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicBPausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBPausedIterator{contract: _AutomationRegistryLogicB.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBPaused)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParsePaused(log types.Log) (*AutomationRegistryLogicBPaused, error) {
	event := new(AutomationRegistryLogicBPaused)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBPayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicBPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicBPayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBPayeesUpdatedIterator{contract: _AutomationRegistryLogicB.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBPayeesUpdated)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicBPayeesUpdated, error) {
	event := new(AutomationRegistryLogicBPayeesUpdated)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBPayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicBPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicBPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBPayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicB.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBPayeeshipTransferRequested)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicBPayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicBPayeeshipTransferRequested)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBPayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicBPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicBPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBPayeeshipTransferredIterator{contract: _AutomationRegistryLogicB.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBPayeeshipTransferred)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicBPayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicBPayeeshipTransferred)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBPaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicBPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicBPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBPaymentWithdrawnIterator{contract: _AutomationRegistryLogicB.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBPaymentWithdrawn)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicBPaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicBPaymentWithdrawn)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicBReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicB.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBReorgedUpkeepReport)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicBReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicBReorgedUpkeepReport)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBStaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicBStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBStaleUpkeepReportIterator{contract: _AutomationRegistryLogicB.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBStaleUpkeepReport)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicBStaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicBStaleUpkeepReport)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUnpausedIterator struct {
	Event *AutomationRegistryLogicBUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicBUnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUnpausedIterator{contract: _AutomationRegistryLogicB.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUnpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUnpaused)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicBUnpaused, error) {
	event := new(AutomationRegistryLogicBUnpaused)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicBUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicBUpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicBUpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicBUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicBUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepAdminTransferred)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicBUpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicBUpkeepAdminTransferred)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicBUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicBUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepCanceledIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepCanceled)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicBUpkeepCanceled, error) {
	event := new(AutomationRegistryLogicBUpkeepCanceled)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicBUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepCheckDataSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicBUpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicBUpkeepCheckDataSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicBUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepGasLimitSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicBUpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicBUpkeepGasLimitSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicBUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepMigratedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepMigrated)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicBUpkeepMigrated, error) {
	event := new(AutomationRegistryLogicBUpkeepMigrated)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicBUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicBUpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepPausedIterator struct {
	Event *AutomationRegistryLogicBUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepPausedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepPaused)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicBUpkeepPaused, error) {
	event := new(AutomationRegistryLogicBUpkeepPaused)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicBUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicBUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepPerformedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepPerformed)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicBUpkeepPerformed, error) {
	event := new(AutomationRegistryLogicBUpkeepPerformed)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicBUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicBUpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicBUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepReceivedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepReceived)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicBUpkeepReceived, error) {
	event := new(AutomationRegistryLogicBUpkeepReceived)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicBUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepRegisteredIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepRegistered)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicBUpkeepRegistered, error) {
	event := new(AutomationRegistryLogicBUpkeepRegistered)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicBUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicBUpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicBUpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicBUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicBUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicBUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicBUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicBUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicBUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicBUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBUpkeepUnpausedIterator{contract: _AutomationRegistryLogicB.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicBUpkeepUnpaused)
				if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicBUpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicBUpkeepUnpaused)
	if err := _AutomationRegistryLogicB.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicB) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicB.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistryLogicB.ParseBillingConfigOverridden(log)
	case _AutomationRegistryLogicB.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistryLogicB.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistryLogicB.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseBillingConfigSet(log)
	case _AutomationRegistryLogicB.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicB.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicB.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicB.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicB.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicB.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicB.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistryLogicB.ParseFeesWithdrawn(log)
	case _AutomationRegistryLogicB.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicB.ParseFundsAdded(log)
	case _AutomationRegistryLogicB.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicB.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicB.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicB.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicB.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistryLogicB.ParseNOPsSettledOffchain(log)
	case _AutomationRegistryLogicB.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicB.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicB.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicB.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicB.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicB.ParsePaused(log)
	case _AutomationRegistryLogicB.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicB.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicB.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicB.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicB.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicB.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicB.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicB.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicB.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicB.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicB.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicB.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicB.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicB.ParseUnpaused(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicB.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicB.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicBAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicBBillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistryLogicBBillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistryLogicBBillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0x720a5849025dc4fd0061aed1bb30efd713cde64ce7f8d807953ecca27c8f143c")
}

func (AutomationRegistryLogicBCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicBChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicBDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicBFeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistryLogicBFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicBFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicBInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicBNOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistryLogicBOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicBOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicBPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicBPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicBPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicBPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicBPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicBReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicBStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicBUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicBUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicBUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicBUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicBUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicBUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicBUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicBUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicBUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicBUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicBUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicBUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicBUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicBUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicBUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicB) Address() common.Address {
	return _AutomationRegistryLogicB.address
}

type AutomationRegistryLogicBInterface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	RemoveBillingOverrides(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	SetBillingOverrides(opts *bind.TransactOpts, id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawERC20Fees(opts *bind.TransactOpts, asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicBAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicBAdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBBillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicBBillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBBillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicBBillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicBBillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBBillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicBBillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicBCancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicBChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicBChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicBDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicBDedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicBFeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicBFeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicBFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicBFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicBFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicBInsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicBNOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBNOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicBNOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicBOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicBOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicBOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicBOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicBPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicBPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicBPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicBPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicBPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicBPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicBPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicBPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicBPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicBPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicBReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicBStaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicBUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicBUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicBUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicBUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicBUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicBUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicBUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicBUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicBUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicBUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicBUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicBUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicBUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicBUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicBUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicBUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicBUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicBUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicBUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
