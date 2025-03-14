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
	Decimals          uint8
	FallbackPrice     *big.Int
	MinSpend          *big.Int
}

type AutomationRegistryBase23BillingOverrides struct {
	GasFeePPB         uint32
	FlatFeeMilliCents *big.Int
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

var AutomationRegistryLogicB23MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"logicC\",\"type\":\"address\",\"internalType\":\"contractAutomationRegistryLogicC2_3\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"checkCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"values\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkUSD\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkUSD\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_3.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fallbackTo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeBillingOverrides\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBillingOverrides\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"billingOverrides\",\"type\":\"tuple\",\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepOffchainConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"simulatePerformUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20Fees\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawLink\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverridden\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"overrides\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigOverrideRemoved\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingConfigSet\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"components\":[{\"name\":\"gasFeePPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"fallbackPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainSpecificModuleUpdated\",\"inputs\":[{\"name\":\"newModule\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeesWithdrawn\",\"inputs\":[{\"name\":\"assetAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NOPsSettledOffchain\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payments\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCharged\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receipt\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"components\":[{\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInBillingToken\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"premiumInJuels\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingToken\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"linkUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"nativeUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"billingUSD\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLinkLiquidity\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOffchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustSettleOnchain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyFinanceAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x6101806040523480156200001257600080fd5b5060405162005d2238038062005d2283398101604081905262000035916200062f565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b91906200062f565b826001600160a01b031663226cf83c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200062f565b836001600160a01b031663614486af6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200062f565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200062f565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200062f565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200062f565b876001600160a01b031663c5b964e06040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002d3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002f9919062000656565b886001600160a01b031663ac4dc59a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000338573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200035e91906200062f565b3380600081620003b55760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620003e857620003e8816200056b565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff19166001838181111562000439576200043962000679565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200049a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004c091906200068f565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000504573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200052a91906200068f565b60ff16146200054c576040516301f86e1760e41b815260040160405180910390fd5b5050506001600160a01b039095166101605250620006b4945050505050565b336001600160a01b03821603620005c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620003ac565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200062c57600080fd5b50565b6000602082840312156200064257600080fd5b81516200064f8162000616565b9392505050565b6000602082840312156200066957600080fd5b8151600281106200064f57600080fd5b634e487b7160e01b600052602160045260246000fd5b600060208284031215620006a257600080fd5b815160ff811681146200064f57600080fd5b60805160a05160c05160e051610100516101205161014051610160516155e86200073a60003960008181610182015261022f0152600081816120e7015261228d01526000612bb30152600050506000612ee901526000613c5d01526000612fc3015260008181610c4701528181610d0801528181610dd30152612c7401526155e86000f3fe6080604052600436106101805760003560e01c80638765ecbe116100d6578063aed2e9291161007f578063ce7dc5b411610059578063ce7dc5b4146104b1578063f2fde38b146104d1578063f7d334ba146104f157610180565b8063aed2e9291461043a578063b148ab6b14610471578063cd7f71b51461049157610180565b8063948108f7116100b0578063948108f7146103e7578063a72aa27e146103fa578063a86e17811461041a57610180565b80638765ecbe1461037c5780638da5cb5b1461039c5780638dcf0fe7146103c757610180565b806354b7faae1161013857806371fae17f1161011257806371fae17f14610327578063744bfe611461034757806379ba50971461036757610180565b806354b7faae146102b457806368d369d8146102d457806371791aa0146102f457610180565b8063349e8cca11610169578063349e8cca146102205780634ee88d35146102745780635165f2f51461029457610180565b80631a2af011146101c757806329c5efad146101e7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156101c0573d6000f35b3d6000fd5b005b3480156101d357600080fd5b506101c56101e2366004614490565b610511565b3480156101f357600080fd5b50610207610202366004614604565b610617565b6040516102179493929190614723565b60405180910390f35b34801561022c57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610217565b34801561028057600080fd5b506101c561028f3660046147a5565b6108f5565b3480156102a057600080fd5b506101c56102af3660046147f1565b610957565b3480156102c057600080fd5b506101c56102cf36600461480a565b610b09565b3480156102e057600080fd5b506101c56102ef366004614836565b610d7c565b34801561030057600080fd5b5061031461030f366004614604565b61102a565b6040516102179796959493929190614877565b34801561033357600080fd5b506101c56103423660046147f1565b61178e565b34801561035357600080fd5b506101c5610362366004614490565b611824565b34801561037357600080fd5b506101c5611c9c565b34801561038857600080fd5b506101c56103973660046147f1565b611d99565b3480156103a857600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661024f565b3480156103d357600080fd5b506101c56103e23660046147a5565b611f4e565b6101c56103f53660046148c5565b611fa3565b34801561040657600080fd5b506101c561041536600461490d565b612365565b34801561042657600080fd5b506101c5610435366004614932565b612464565b34801561044657600080fd5b5061045a6104553660046147a5565b612545565b604080519215158352602083019190915201610217565b34801561047d57600080fd5b506101c561048c3660046147f1565b6126f0565b34801561049d57600080fd5b506101c56104ac3660046147a5565b61291d565b3480156104bd57600080fd5b506102076104cc3660046149ac565b6129d4565b3480156104dd57600080fd5b506101c56104ec366004614a8e565b612a96565b3480156104fd57600080fd5b5061031461050c3660046147f1565b612aaa565b61051a82612ae6565b3373ffffffffffffffffffffffffffffffffffffffff821603610569576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116146106135760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b60006060600080610626612b9b565b600086815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615159483019490945263ffffffff620100008204811695830195909552660100000000000081048516606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490941660e0820152600290910154909216908201525a9150600080826080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561077c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107a09190614abb565b73ffffffffffffffffffffffffffffffffffffffff16601660000160149054906101000a900463ffffffff1663ffffffff16896040516107e09190614ad8565b60006040518083038160008787f1925050503d806000811461081e576040519150601f19603f3d011682016040523d82523d6000602084013e610823565b606091505b50915091505a6108339085614b23565b93508161085c5760006040518060200160405280600081525060079650965096505050506108ec565b808060200190518101906108709190614b8b565b90975095508661089c5760006040518060200160405280600081525060049650965096505050506108ec565b60185486517401000000000000000000000000000000000000000090910463ffffffff1610156108e85760006040518060200160405280600081525060059650965096505050506108ec565b5050505b92959194509250565b6108fe83612ae6565b6000838152601d60205260409020610917828483614c70565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664838360405161094a929190614dd4565b60405180910390a2505050565b61096081612ae6565b600081815260046020908152604091829020825161012081018452815460ff808216151580845261010080840490921615159584019590955262010000820463ffffffff9081169684019690965266010000000000008204861660608401526a010000000000000000000090910473ffffffffffffffffffffffffffffffffffffffff908116608084015260018401546fffffffffffffffffffffffffffffffff811660a085015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08501527c0100000000000000000000000000000000000000000000000000000000900490951660e083015260029092015490931690830152610a9a576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055610ad9600283612c0c565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b610b11612c21565b73ffffffffffffffffffffffffffffffffffffffff8216610b5e576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610b68612c72565b90506000811215610bb4576040517fcf47918100000000000000000000000000000000000000000000000000000000815260006004820152602481018390526044015b60405180910390fd5b80821115610bf8576040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610bab565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490526000917f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015610c92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cb69190614de8565b905080610cef576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8373ffffffffffffffffffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa885604051610d6e91815260200190565b60405180910390a350505050565b610d84612c21565b73ffffffffffffffffffffffffffffffffffffffff8216610dd1576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610e56576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610e60612c72565b128015610e835750600060255460ff166001811115610e8157610e816146b9565b145b15610eba576040517f981bb6a000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015610f36573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f5a9190614e03565b610f649190614b23565b905080821115610faa576040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610bab565b610fcb73ffffffffffffffffffffffffffffffffffffffff85168484612d41565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa884604051610d6e91815260200190565b60006060600080600080600061103e612b9b565b60006110498a612e1a565b905060006014604051806101200160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900462ffffff1662ffffff1662ffffff1681526020016000820160139054906101000a900461ffff1661ffff1661ffff1681526020016000820160159054906101000a900460ff1660ff1660ff1681526020016000820160169054906101000a900460ff161515151581526020016000820160179054906101000a900460ff161515151581526020016000820160189054906101000a900460ff161515151581526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d8152602001908152602001600020604051806101200160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900460ff161515151581526020016000820160029054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160069054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200160008201600a9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1681526020016001820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000808360a001511561140c5750506040805160208101825260008082529290910151919a5098506009975089965063ffffffff169450859350839250611782915050565b606083015163ffffffff908116146114565750506040805160208101825260008082529290910151919a5098506001975089965063ffffffff169450859350839250611782915050565b8251156114955750506040805160208101825260008082529290910151919a5098506002975089965063ffffffff169450859350839250611782915050565b61149e84612ec5565b8094508198508299505050506114c38e858786604001518b8b888a61010001516130b7565b9050806bffffffffffffffffffffffff168360c001516bffffffffffffffffffffffff1610156115255750506040805160208101825260008082529290910151919a5098506006975089965063ffffffff169450859350839250611782915050565b505060006115348d858e613427565b90505a9750600080836080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561158b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115af9190614abb565b73ffffffffffffffffffffffffffffffffffffffff16601660000160149054906101000a900463ffffffff1663ffffffff16846040516115ef9190614ad8565b60006040518083038160008787f1925050503d806000811461162d576040519150601f19603f3d011682016040523d82523d6000602084013e611632565b606091505b50915091505a611642908b614b23565b9950816116c8576018548151780100000000000000000000000000000000000000000000000090910463ffffffff1610156116a75750506040805160208101825260008082529390910151929b509950600898505063ffffffff169450611782915050565b604090930151929a50600399505063ffffffff909116955061178292505050565b808060200190518101906116dc9190614b8b565b909d509b508c6117165750506040805160208101825260008082529390910151929b509950600498505063ffffffff169450611782915050565b6018548c517401000000000000000000000000000000000000000090910463ffffffff1610156117705750506040805160208101825260008082529390910151929b509950600598505063ffffffff169450611782915050565b5050506040015163ffffffff16945050505b92959891949750929550565b611796613607565b600081815260046020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055602390915280822080547fffffffffffffffffffffffffffffffffffffffffffffffffff000000000000001690555182917f97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a91a250565b60145477010000000000000000000000000000000000000000000000900460ff161561187c576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff167701000000000000000000000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff811661190b576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615158387015263ffffffff620100008304811684870152660100000000000083048116606085015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009093048316608085015260018501546fffffffffffffffffffffffffffffffff811660a08601526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08601527c010000000000000000000000000000000000000000000000000000000090041660e084015260029093015481169282019290925286855260059093529220549091163314611a4b576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554604080517f57e871e7000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff909216916357e871e7916004808201926020929091908290030181865afa158015611abb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611adf9190614e03565b816060015163ffffffff161115611b22576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526004602090815260408083206001015461010085015173ffffffffffffffffffffffffffffffffffffffff1684526021909252909120547001000000000000000000000000000000009091046bffffffffffffffffffffffff1690611b8d908290614b23565b6101008301805173ffffffffffffffffffffffffffffffffffffffff908116600090815260216020908152604080832095909555888252600490529290922060010180547fffffffff000000000000000000000000ffffffffffffffffffffffffffffffff16905551611c109116846bffffffffffffffffffffffff8416612d41565b604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d1d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bab565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611da281612ae6565b600081815260046020908152604091829020825161012081018452815460ff808216158015845261010080840490921615159584019590955262010000820463ffffffff9081169684019690965266010000000000008204861660608401526a010000000000000000000090910473ffffffffffffffffffffffffffffffffffffffff908116608084015260018401546fffffffffffffffffffffffffffffffff811660a085015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08501527c0100000000000000000000000000000000000000000000000000000000900490951660e083015260029092015490931690830152611edc576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055611f1e600283613658565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b611f5783612ae6565b6000838152601e60205260409020611f70828483614c70565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850838360405161094a929190614dd4565b600082815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615159483019490945263ffffffff6201000082048116958301959095526601000000000000810485166060830181905273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009092048216608084015260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c01000000000000000000000000000000000000000000000000000000009004861660e084015260029093015416928101929092529091146120df576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b341561217b577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681610100015173ffffffffffffffffffffffffffffffffffffffff161461216f576040517fc1ab6dc100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61217834613664565b91505b818160c0015161218b9190614e1c565b600084815260046020908152604080832060010180547fffffffff000000000000000000000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000006bffffffffffffffffffffffff9687160217905561010085015173ffffffffffffffffffffffffffffffffffffffff168352602190915290205461221b91841690614e41565b61010082015173ffffffffffffffffffffffffffffffffffffffff1660009081526021602052604081209190915534900361228b576101008101516122869073ffffffffffffffffffffffffffffffffffffffff1633306bffffffffffffffffffffffff8616613706565b61231b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0836bffffffffffffffffffffffff166040518263ffffffff1660e01b81526004016000604051808303818588803b15801561230157600080fd5b505af1158015612315573d6000803e3d6000fd5b50505050505b6040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6108fc8163ffffffff1610806123a2575060165463ffffffff78010000000000000000000000000000000000000000000000009091048116908216115b156123d9576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6123e282612ae6565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffff166201000063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c91015b60405180910390a25050565b61246c613607565b6000828152600460205260409020546601000000000000900463ffffffff908116146124c4576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790556023909152902081906125128282614e65565b905050817fd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340826040516124589190614eec565b600080612550612b9b565b601454760100000000000000000000000000000000000000000000900460ff16156125a7576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615158386015262010000820463ffffffff90811684880181905266010000000000008404821660608601526a010000000000000000000090930473ffffffffffffffffffffffffffffffffffffffff9081166080860181905260018701546fffffffffffffffffffffffffffffffff811660a088015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660c08801527c0100000000000000000000000000000000000000000000000000000000900490921660e0860152600290950154909416908301528451601f890185900485028101850190955287855290936126e393899089908190840183828082843760009201919091525061376a92505050565b9097909650945050505050565b600081815260046020908152604091829020825161012081018452815460ff8082161515835261010080830490911615159483019490945263ffffffff6201000082048116958301959095526601000000000000810485166060830181905273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009092048216608084015260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c01000000000000000000000000000000000000000000000000000000009004861660e0840152600290930154169281019290925290911461282c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612889576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b61292683612ae6565b6017547c0100000000000000000000000000000000000000000000000000000000900463ffffffff16811115612988576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206129a1828483614c70565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d838360405161094a929190614dd4565b600060606000806000634b56a42e60e01b8888886040516024016129fa93929190614f23565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050612a838982610617565b929c919b50995090975095505050505050565b612a9e613985565b612aa781613a06565b50565b600060606000806000806000612acf886040518060200160405280600081525061102a565b959e949d50929b5090995097509550909350915050565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612b43576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020546601000000000000900463ffffffff90811614612aa7576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614612c0a576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000612c188383613afb565b90505b92915050565b60185473ffffffffffffffffffffffffffffffffffffffff163314612c0a576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015612d0e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d329190614e03565b612d3c9190614fb9565b905090565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052612e159084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152613b4a565b505050565b6000818160045b600f811015612ea7577fff000000000000000000000000000000000000000000000000000000000000008216838260208110612e5f57612e5f614fd9565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612e9557506000949350505050565b80612e9f81615008565b915050612e21565b5081600f1a6001811115612ebd57612ebd6146b9565b949350505050565b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015612f52573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f76919061505a565b5094509092505050600081131580612f8d57508142105b80612fae5750828015612fae5750612fa58242614b23565b8463ffffffff16105b15612fbd576019549650612fc1565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561302c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613050919061505a565b509450909250505060008113158061306757508142105b806130885750828015613088575061307f8242614b23565b8463ffffffff16105b1561309757601a54955061309b565b8095505b86866130a68a613c56565b965096509650505050509193909250565b60008080808960018111156130ce576130ce6146b9565b036130dd575062017f98613132565b60018960018111156130f1576130f16146b9565b0361310057506201e26c613132565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a60800151600161314591906150aa565b6131539060ff1660406150c3565b601854613181906103a49074010000000000000000000000000000000000000000900463ffffffff16614e41565b61318b9190614e41565b601554604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa1580156131fc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061322091906150da565b90925090508183613232836018614e41565b61323c91906150c3565b60808f015161324c9060016150aa565b61325b9060ff166115e06150c3565b6132659190614e41565b61326f9190614e41565b6132799085614e41565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa1580156132ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133109190614e03565b8d6060015161ffff1661332391906150c3565b945050505060006133348b86613d47565b60008d815260046020526040902054909150610100900460ff16156133995760008c81526023602090815260409182902082518084018452905463ffffffff811680835264010000000090910462ffffff908116928401928352928501525116908201525b60006134038c6040518061012001604052808d63ffffffff1681526020018681526020018781526020018c81526020018b81526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff16815260200185815260200160001515815250613ec3565b6020810151815191925061341691614e1c565b9d9c50505050505050505050505050565b6060600083600181111561343d5761343d6146b9565b03613506576000848152600760205260409081902090517f6e04ff0d000000000000000000000000000000000000000000000000000000009161348291602401615199565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050613600565b600183600181111561351a5761351a6146b9565b03613100576000828060200190518101906135359190615212565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009161357a918491602401615322565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915291506136009050565b9392505050565b60175473ffffffffffffffffffffffffffffffffffffffff163314612c0a576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612c18838361420a565b60006bffffffffffffffffffffffff821115613702576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f36206269747300000000000000000000000000000000000000000000000000006064820152608401610bab565b5090565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526137649085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401612d93565b50505050565b601454600090819077010000000000000000000000000000000000000000000000900460ff16156137c7576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16770100000000000000000000000000000000000000000000001790556040517f4585e33b000000000000000000000000000000000000000000000000000000009061383c9085906024016153ed565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061390f9087908790600401615400565b60408051808303816000875af115801561392d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139519190615419565b601480547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612c0a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bab565b3373ffffffffffffffffffffffffffffffffffffffff821603613a85576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bab565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000818152600183016020526040812054613b4257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612c1b565b506000612c1b565b6000613bac826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166142fd9092919063ffffffff16565b805190915015612e155780806020019051810190613bca9190614de8565b612e15576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610bab565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613cc6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613cea919061505a565b50935050925050600082131580613d0057508042105b80613d3057506000846040015162ffffff16118015613d305750613d248142614b23565b846040015162ffffff16105b15613d40575050601b5492915050565b5092915050565b60408051608081018252600080825260208083018281528385018381526060850184905273ffffffffffffffffffffffffffffffffffffffff878116855260229093528584208054640100000000810462ffffff1690925263ffffffff82169092527b01000000000000000000000000000000000000000000000000000000810460ff16855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495919484936701000000000000009092049091169163feaf968c9160048083019260a09291908290030181865afa158015613e34573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e58919061505a565b50935050925050600082131580613e6e57508042105b80613e9e57506000866040015162ffffff16118015613e9e5750613e928142614b23565b866040015162ffffff16105b15613eb25760018301546060850152613eba565b606084018290525b50505092915050565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081019190915260008260e001516000015160ff1690506000846060015161ffff168460600151613f2e91906150c3565b90508361010001518015613f415750803a105b15613f4957503a5b600060128311613f5a576001613f70565b613f65601284614b23565b613f7090600a615565565b9050600060128410613f83576001613f99565b613f8e846012614b23565b613f9990600a615565565b905060008660a00151876040015188602001518960000151613fbb9190614e41565b613fc590876150c3565b613fcf9190614e41565b613fd991906150c3565b9050614035828860e0015160600151613ff291906150c3565b6001848a60e001516060015161400891906150c3565b6140129190614b23565b61401c86856150c3565b6140269190614e41565b6140309190615571565b613664565b6bffffffffffffffffffffffff1686526080870151614058906140309083615571565b6bffffffffffffffffffffffff1660408088019190915260e088015101516000906140919062ffffff16683635c9adc5dea000006150c3565b9050600081633b9aca008a60a001518b60e001516020015163ffffffff168c604001518d600001518b6140c491906150c3565b6140ce9190614e41565b6140d891906150c3565b6140e291906150c3565b6140ec9190615571565b6140f69190614e41565b9050614139848a60e001516060015161410f91906150c3565b6001868c60e001516060015161412591906150c3565b61412f9190614b23565b61401c88856150c3565b6bffffffffffffffffffffffff166020890152608089015161415f906140309083615571565b6bffffffffffffffffffffffff16606089015260c089015173ffffffffffffffffffffffffffffffffffffffff166080808a01919091528901516141a290613664565b6bffffffffffffffffffffffff1660a0808a01919091528901516141c590613664565b6bffffffffffffffffffffffff1660c089015260e0890151606001516141ea90613664565b6bffffffffffffffffffffffff1660e08901525050505050505092915050565b600081815260018301602052604081205480156142f357600061422e600183614b23565b855490915060009061424290600190614b23565b90508181146142a757600086600001828154811061426257614262614fd9565b906000526020600020015490508087600001848154811061428557614285614fd9565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806142b8576142b86155ac565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612c1b565b6000915050612c1b565b6060612ebd8484600085856000808673ffffffffffffffffffffffffffffffffffffffff1685876040516143319190614ad8565b60006040518083038185875af1925050503d806000811461436e576040519150601f19603f3d011682016040523d82523d6000602084013e614373565b606091505b50915091506143848783838761438f565b979650505050505050565b6060831561442557825160000361441e5773ffffffffffffffffffffffffffffffffffffffff85163b61441e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610bab565b5081612ebd565b612ebd838381511561443a5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bab91906153ed565b73ffffffffffffffffffffffffffffffffffffffff81168114612aa757600080fd5b600080604083850312156144a357600080fd5b8235915060208301356144b58161446e565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff81118282101715614513576145136144c0565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614560576145606144c0565b604052919050565b600067ffffffffffffffff821115614582576145826144c0565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126145bf57600080fd5b81356145d26145cd82614568565b614519565b8181528460208386010111156145e757600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561461757600080fd5b82359150602083013567ffffffffffffffff81111561463557600080fd5b614641858286016145ae565b9150509250929050565b60005b8381101561466657818101518382015260200161464e565b50506000910152565b6000815180845261468781602086016020860161464b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a811061471f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600061473e608083018661466f565b905061474d60408301856146e8565b82606083015295945050505050565b60008083601f84011261476e57600080fd5b50813567ffffffffffffffff81111561478657600080fd5b60208301915083602082850101111561479e57600080fd5b9250929050565b6000806000604084860312156147ba57600080fd5b83359250602084013567ffffffffffffffff8111156147d857600080fd5b6147e48682870161475c565b9497909650939450505050565b60006020828403121561480357600080fd5b5035919050565b6000806040838503121561481d57600080fd5b82356148288161446e565b946020939093013593505050565b60008060006060848603121561484b57600080fd5b83356148568161446e565b925060208401356148668161446e565b929592945050506040919091013590565b871515815260e06020820152600061489260e083018961466f565b90506148a160408301886146e8565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080604083850312156148d857600080fd5b8235915060208301356bffffffffffffffffffffffff811681146144b557600080fd5b63ffffffff81168114612aa757600080fd5b6000806040838503121561492057600080fd5b8235915060208301356144b5816148fb565b600080828403606081121561494657600080fd5b8335925060407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08201121561497a57600080fd5b506020830190509250929050565b600067ffffffffffffffff8211156149a2576149a26144c0565b5060051b60200190565b600080600080606085870312156149c257600080fd5b8435935060208086013567ffffffffffffffff808211156149e257600080fd5b818801915088601f8301126149f657600080fd5b8135614a046145cd82614988565b81815260059190911b8301840190848101908b831115614a2357600080fd5b8585015b83811015614a5b57803585811115614a3f5760008081fd5b614a4d8e89838a01016145ae565b845250918601918601614a27565b50975050506040880135925080831115614a7457600080fd5b5050614a828782880161475c565b95989497509550505050565b600060208284031215614aa057600080fd5b81356136008161446e565b8051614ab68161446e565b919050565b600060208284031215614acd57600080fd5b81516136008161446e565b60008251614aea81846020870161464b565b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115612c1b57612c1b614af4565b80518015158114614ab657600080fd5b600082601f830112614b5757600080fd5b8151614b656145cd82614568565b818152846020838601011115614b7a57600080fd5b612ebd82602083016020870161464b565b60008060408385031215614b9e57600080fd5b614ba783614b36565b9150602083015167ffffffffffffffff811115614bc357600080fd5b61464185828601614b46565b600181811c90821680614be357607f821691505b602082108103614c1c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115612e1557600081815260208120601f850160051c81016020861015614c495750805b601f850160051c820191505b81811015614c6857828155600101614c55565b505050505050565b67ffffffffffffffff831115614c8857614c886144c0565b614c9c83614c968354614bcf565b83614c22565b6000601f841160018114614cee5760008515614cb85750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614d84565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614d3d5786850135825560209485019460019092019101614d1d565b5086821015614d78577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000612ebd602083018486614d8b565b600060208284031215614dfa57600080fd5b612c1882614b36565b600060208284031215614e1557600080fd5b5051919050565b6bffffffffffffffffffffffff818116838216019080821115613d4057613d40614af4565b80820180821115612c1b57612c1b614af4565b62ffffff81168114612aa757600080fd5b8135614e70816148fb565b63ffffffff811690508154817fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000082161783556020840135614eb081614e54565b66ffffff000000008160201b16837fffffffffffffffffffffffffffffffffffffffffffffffffff000000000000008416171784555050505050565b604081018235614efb816148fb565b63ffffffff1682526020830135614f1181614e54565b62ffffff811660208401525092915050565b6000604082016040835280865180835260608501915060608160051b8601019250602080890160005b83811015614f98577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552614f8686835161466f565b95509382019390820190600101614f4c565b505085840381870152505050614faf818587614d8b565b9695505050505050565b8181036000831280158383131683831282161715613d4057613d40614af4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361503957615039614af4565b5060010190565b805169ffffffffffffffffffff81168114614ab657600080fd5b600080600080600060a0868803121561507257600080fd5b61507b86615040565b945060208601519350604086015192506060860151915061509e60808701615040565b90509295509295909350565b60ff8181168382160190811115612c1b57612c1b614af4565b8082028115828204841417612c1b57612c1b614af4565b600080604083850312156150ed57600080fd5b505080516020909101519092909150565b6000815461510b81614bcf565b80855260206001838116801561512857600181146151605761518e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061518e565b866000528260002060005b858110156151865781548a820186015290830190840161516b565b890184019650505b505050505092915050565b602081526000612c1860208301846150fe565b600082601f8301126151bd57600080fd5b815160206151cd6145cd83614988565b82815260059290921b840181019181810190868411156151ec57600080fd5b8286015b8481101561520757805183529183019183016151f0565b509695505050505050565b60006020828403121561522457600080fd5b815167ffffffffffffffff8082111561523c57600080fd5b90830190610100828603121561525157600080fd5b6152596144ef565b825181526020830151602082015260408301516040820152606083015160608201526080830151608082015261529160a08401614aab565b60a082015260c0830151828111156152a857600080fd5b6152b4878286016151ac565b60c08301525060e0830151828111156152cc57600080fd5b6152d887828601614b46565b60e08301525095945050505050565b600081518084526020808501945080840160005b83811015615317578151875295820195908201906001016152fb565b509495945050505050565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506153936101408401826152e7565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526153cf828261466f565b91505082810360208401526153e481856150fe565b95945050505050565b602081526000612c18602083018461466f565b828152604060208201526000612ebd604083018461466f565b6000806040838503121561542c57600080fd5b61543583614b36565b9150602083015190509250929050565b600181815b8085111561549e57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561548457615484614af4565b8085161561549157918102915b93841c939080029061544a565b509250929050565b6000826154b557506001612c1b565b816154c257506000612c1b565b81600181146154d857600281146154e2576154fe565b6001915050612c1b565b60ff8411156154f3576154f3614af4565b50506001821b612c1b565b5060208310610133831016604e8410600b8410161715615521575081810a612c1b565b61552b8383615445565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561555d5761555d614af4565b029392505050565b6000612c1883836154a6565b6000826155a7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var AutomationRegistryLogicB23ABI = AutomationRegistryLogicB23MetaData.ABI

var AutomationRegistryLogicB23Bin = AutomationRegistryLogicB23MetaData.Bin

func DeployAutomationRegistryLogicB23(auth *bind.TransactOpts, backend bind.ContractBackend, logicC common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicB23, error) {
	parsed, err := AutomationRegistryLogicB23MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicB23Bin), backend, logicC)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicB23{address: address, abi: *parsed, AutomationRegistryLogicB23Caller: AutomationRegistryLogicB23Caller{contract: contract}, AutomationRegistryLogicB23Transactor: AutomationRegistryLogicB23Transactor{contract: contract}, AutomationRegistryLogicB23Filterer: AutomationRegistryLogicB23Filterer{contract: contract}}, nil
}

type AutomationRegistryLogicB23 struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicB23Caller
	AutomationRegistryLogicB23Transactor
	AutomationRegistryLogicB23Filterer
}

type AutomationRegistryLogicB23Caller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB23Transactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB23Filterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicB23Session struct {
	Contract     *AutomationRegistryLogicB23
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicB23CallerSession struct {
	Contract *AutomationRegistryLogicB23Caller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicB23TransactorSession struct {
	Contract     *AutomationRegistryLogicB23Transactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicB23Raw struct {
	Contract *AutomationRegistryLogicB23
}

type AutomationRegistryLogicB23CallerRaw struct {
	Contract *AutomationRegistryLogicB23Caller
}

type AutomationRegistryLogicB23TransactorRaw struct {
	Contract *AutomationRegistryLogicB23Transactor
}

func NewAutomationRegistryLogicB23(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicB23, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicB23ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicB23(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23{address: address, abi: abi, AutomationRegistryLogicB23Caller: AutomationRegistryLogicB23Caller{contract: contract}, AutomationRegistryLogicB23Transactor: AutomationRegistryLogicB23Transactor{contract: contract}, AutomationRegistryLogicB23Filterer: AutomationRegistryLogicB23Filterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicB23Caller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicB23Caller, error) {
	contract, err := bindAutomationRegistryLogicB23(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23Caller{contract: contract}, nil
}

func NewAutomationRegistryLogicB23Transactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicB23Transactor, error) {
	contract, err := bindAutomationRegistryLogicB23(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23Transactor{contract: contract}, nil
}

func NewAutomationRegistryLogicB23Filterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicB23Filterer, error) {
	contract, err := bindAutomationRegistryLogicB23(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23Filterer{contract: contract}, nil
}

func bindAutomationRegistryLogicB23(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicB23MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB23.Contract.AutomationRegistryLogicB23Caller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AutomationRegistryLogicB23Transactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AutomationRegistryLogicB23Transactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicB23.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Caller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB23.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicB23.Contract.FallbackTo(&_AutomationRegistryLogicB23.CallOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23CallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicB23.Contract.FallbackTo(&_AutomationRegistryLogicB23.CallOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB23.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB23.Contract.Owner(&_AutomationRegistryLogicB23.CallOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23CallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicB23.Contract.Owner(&_AutomationRegistryLogicB23.CallOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AcceptOwnership(&_AutomationRegistryLogicB23.TransactOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AcceptOwnership(&_AutomationRegistryLogicB23.TransactOpts)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AcceptUpkeepAdmin(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "addFunds", id, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AddFunds(&_AutomationRegistryLogicB23.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.AddFunds(&_AutomationRegistryLogicB23.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckCallback(&_AutomationRegistryLogicB23.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckCallback(&_AutomationRegistryLogicB23.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "checkUpkeep0", id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckUpkeep0(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.CheckUpkeep0(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "executeCallback", id, payload)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.ExecuteCallback(&_AutomationRegistryLogicB23.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.ExecuteCallback(&_AutomationRegistryLogicB23.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "pauseUpkeep", id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.PauseUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.PauseUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) RemoveBillingOverrides(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "removeBillingOverrides", id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.RemoveBillingOverrides(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.RemoveBillingOverrides(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SetBillingOverrides(opts *bind.TransactOpts, id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "setBillingOverrides", id, billingOverrides)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetBillingOverrides(&_AutomationRegistryLogicB23.TransactOpts, id, billingOverrides)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetBillingOverrides(&_AutomationRegistryLogicB23.TransactOpts, id, billingOverrides)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB23.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepCheckData(&_AutomationRegistryLogicB23.TransactOpts, id, newCheckData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB23.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepGasLimit(&_AutomationRegistryLogicB23.TransactOpts, id, gasLimit)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB23.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepOffchainConfig(&_AutomationRegistryLogicB23.TransactOpts, id, config)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicB23.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicB23.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SimulatePerformUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id, performData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.SimulatePerformUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id, performData)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.TransferOwnership(&_AutomationRegistryLogicB23.TransactOpts, to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.TransferOwnership(&_AutomationRegistryLogicB23.TransactOpts, to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB23.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.TransferUpkeepAdmin(&_AutomationRegistryLogicB23.TransactOpts, id, proposed)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.UnpauseUpkeep(&_AutomationRegistryLogicB23.TransactOpts, id)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) WithdrawERC20Fees(opts *bind.TransactOpts, asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "withdrawERC20Fees", asset, to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB23.TransactOpts, asset, to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB23.TransactOpts, asset, to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawFunds(&_AutomationRegistryLogicB23.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawFunds(&_AutomationRegistryLogicB23.TransactOpts, id, to)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) WithdrawLink(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.Transact(opts, "withdrawLink", to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawLink(&_AutomationRegistryLogicB23.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.WithdrawLink(&_AutomationRegistryLogicB23.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.Fallback(&_AutomationRegistryLogicB23.TransactOpts, calldata)
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB23.Contract.Fallback(&_AutomationRegistryLogicB23.TransactOpts, calldata)
}

type AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicB23AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23AdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23AdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23AdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB23AdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicB23AdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23BillingConfigOverriddenIterator struct {
	Event *AutomationRegistryLogicB23BillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23BillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23BillingConfigOverridden)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23BillingConfigOverridden)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23BillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23BillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23BillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23BillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23BillingConfigOverriddenIterator{contract: _AutomationRegistryLogicB23.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23BillingConfigOverridden)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicB23BillingConfigOverridden, error) {
	event := new(AutomationRegistryLogicB23BillingConfigOverridden)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistryLogicB23BillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23BillingConfigOverrideRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23BillingConfigOverrideRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23BillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator{contract: _AutomationRegistryLogicB23.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23BillingConfigOverrideRemoved)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicB23BillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistryLogicB23BillingConfigOverrideRemoved)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23BillingConfigSetIterator struct {
	Event *AutomationRegistryLogicB23BillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23BillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23BillingConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23BillingConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23BillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23BillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23BillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicB23BillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23BillingConfigSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23BillingConfigSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicB23BillingConfigSet, error) {
	event := new(AutomationRegistryLogicB23BillingConfigSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23CancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB23CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23CancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23CancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23CancelledUpkeepReportIterator{contract: _AutomationRegistryLogicB23.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23CancelledUpkeepReport)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicB23CancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicB23CancelledUpkeepReport)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicB23ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23ChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23ChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicB23.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23ChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicB23ChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicB23ChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23DedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicB23DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23DedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23DedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicB23DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23DedupKeyAddedIterator{contract: _AutomationRegistryLogicB23.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23DedupKeyAdded)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicB23DedupKeyAdded, error) {
	event := new(AutomationRegistryLogicB23DedupKeyAdded)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23FeesWithdrawnIterator struct {
	Event *AutomationRegistryLogicB23FeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23FeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23FeesWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23FeesWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23FeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23FeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23FeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicB23FeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23FeesWithdrawnIterator{contract: _AutomationRegistryLogicB23.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23FeesWithdrawn)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicB23FeesWithdrawn, error) {
	event := new(AutomationRegistryLogicB23FeesWithdrawn)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23FundsAddedIterator struct {
	Event *AutomationRegistryLogicB23FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23FundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicB23FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23FundsAddedIterator{contract: _AutomationRegistryLogicB23.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23FundsAdded)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicB23FundsAdded, error) {
	event := new(AutomationRegistryLogicB23FundsAdded)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23FundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicB23FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23FundsWithdrawnIterator{contract: _AutomationRegistryLogicB23.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23FundsWithdrawn)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB23FundsWithdrawn, error) {
	event := new(AutomationRegistryLogicB23FundsWithdrawn)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB23InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23InsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23InsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicB23.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23InsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicB23InsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicB23InsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23NOPsSettledOffchainIterator struct {
	Event *AutomationRegistryLogicB23NOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23NOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23NOPsSettledOffchain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23NOPsSettledOffchain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23NOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23NOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23NOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicB23NOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23NOPsSettledOffchainIterator{contract: _AutomationRegistryLogicB23.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23NOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23NOPsSettledOffchain)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicB23NOPsSettledOffchain, error) {
	event := new(AutomationRegistryLogicB23NOPsSettledOffchain)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23OwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB23OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23OwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicB23.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23OwnershipTransferRequested)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicB23OwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicB23OwnershipTransferRequested)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23OwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicB23OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23OwnershipTransferredIterator{contract: _AutomationRegistryLogicB23.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23OwnershipTransferred)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicB23OwnershipTransferred, error) {
	event := new(AutomationRegistryLogicB23OwnershipTransferred)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23PausedIterator struct {
	Event *AutomationRegistryLogicB23Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23PausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB23PausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23PausedIterator{contract: _AutomationRegistryLogicB23.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23Paused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23Paused)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParsePaused(log types.Log) (*AutomationRegistryLogicB23Paused, error) {
	event := new(AutomationRegistryLogicB23Paused)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23PayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicB23PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23PayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23PayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB23PayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23PayeesUpdatedIterator{contract: _AutomationRegistryLogicB23.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23PayeesUpdated)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicB23PayeesUpdated, error) {
	event := new(AutomationRegistryLogicB23PayeesUpdated)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23PayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB23PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23PayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicB23.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23PayeeshipTransferRequested)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicB23PayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicB23PayeeshipTransferRequested)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23PayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicB23PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23PayeeshipTransferredIterator{contract: _AutomationRegistryLogicB23.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23PayeeshipTransferred)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicB23PayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicB23PayeeshipTransferred)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23PaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicB23PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicB23PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23PaymentWithdrawnIterator{contract: _AutomationRegistryLogicB23.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23PaymentWithdrawn)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicB23PaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicB23PaymentWithdrawn)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23ReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB23ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23ReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23ReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23ReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicB23.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23ReorgedUpkeepReport)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicB23ReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicB23ReorgedUpkeepReport)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23StaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicB23StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23StaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23StaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23StaleUpkeepReportIterator{contract: _AutomationRegistryLogicB23.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23StaleUpkeepReport)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicB23StaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicB23StaleUpkeepReport)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UnpausedIterator struct {
	Event *AutomationRegistryLogicB23Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB23UnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UnpausedIterator{contract: _AutomationRegistryLogicB23.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23Unpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23Unpaused)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicB23Unpaused, error) {
	event := new(AutomationRegistryLogicB23Unpaused)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicB23UpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicB23UpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicB23UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepAdminTransferred)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicB23UpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicB23UpkeepAdminTransferred)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicB23UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicB23UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepCanceledIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepCanceled)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicB23UpkeepCanceled, error) {
	event := new(AutomationRegistryLogicB23UpkeepCanceled)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepChargedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepCharged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepChargedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepCharged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepCharged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepChargedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepChargedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepCharged struct {
	Id      *big.Int
	Receipt AutomationRegistryBase23PaymentReceipt
	Raw     types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepChargedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepChargedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepCharged", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCharged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepCharged)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicB23UpkeepCharged, error) {
	event := new(AutomationRegistryLogicB23UpkeepCharged)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicB23UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepCheckDataSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicB23UpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicB23UpkeepCheckDataSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicB23UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepGasLimitSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicB23UpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicB23UpkeepGasLimitSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepMigratedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepMigrated)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicB23UpkeepMigrated, error) {
	event := new(AutomationRegistryLogicB23UpkeepMigrated)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicB23UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicB23UpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepPausedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepPausedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepPaused)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicB23UpkeepPaused, error) {
	event := new(AutomationRegistryLogicB23UpkeepPaused)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicB23UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepPerformedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepPerformed)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicB23UpkeepPerformed, error) {
	event := new(AutomationRegistryLogicB23UpkeepPerformed)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicB23UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicB23UpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepReceivedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepReceived)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicB23UpkeepReceived, error) {
	event := new(AutomationRegistryLogicB23UpkeepReceived)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicB23UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepRegisteredIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepRegistered)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicB23UpkeepRegistered, error) {
	event := new(AutomationRegistryLogicB23UpkeepRegistered)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicB23UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicB23UpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicB23UpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicB23UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicB23UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicB23UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicB23UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicB23UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicB23UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicB23UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicB23UpkeepUnpausedIterator{contract: _AutomationRegistryLogicB23.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicB23.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicB23UpkeepUnpaused)
				if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23Filterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicB23UpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicB23UpkeepUnpaused)
	if err := _AutomationRegistryLogicB23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicB23.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB23.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB23.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistryLogicB23.ParseBillingConfigOverridden(log)
	case _AutomationRegistryLogicB23.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistryLogicB23.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistryLogicB23.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistryLogicB23.ParseBillingConfigSet(log)
	case _AutomationRegistryLogicB23.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicB23.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicB23.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicB23.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicB23.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicB23.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicB23.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistryLogicB23.ParseFeesWithdrawn(log)
	case _AutomationRegistryLogicB23.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicB23.ParseFundsAdded(log)
	case _AutomationRegistryLogicB23.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicB23.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicB23.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicB23.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicB23.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistryLogicB23.ParseNOPsSettledOffchain(log)
	case _AutomationRegistryLogicB23.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicB23.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicB23.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicB23.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicB23.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicB23.ParsePaused(log)
	case _AutomationRegistryLogicB23.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicB23.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicB23.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicB23.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicB23.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicB23.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicB23.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicB23.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicB23.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicB23.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicB23.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicB23.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicB23.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicB23.ParseUnpaused(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepCharged"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepCharged(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicB23.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicB23.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicB23AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicB23BillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistryLogicB23BillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistryLogicB23BillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba3")
}

func (AutomationRegistryLogicB23CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicB23ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicB23DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicB23FeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistryLogicB23FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicB23FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicB23InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicB23NOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistryLogicB23OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicB23OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicB23Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicB23PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicB23PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicB23PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicB23PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicB23ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicB23StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicB23Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicB23UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicB23UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicB23UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicB23UpkeepCharged) Topic() common.Hash {
	return common.HexToHash("0x801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc")
}

func (AutomationRegistryLogicB23UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicB23UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicB23UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicB23UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicB23UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicB23UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicB23UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicB23UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicB23UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicB23UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicB23UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicB23 *AutomationRegistryLogicB23) Address() common.Address {
	return _AutomationRegistryLogicB23.address
}

type AutomationRegistryLogicB23Interface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicB23AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB23AdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23BillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicB23BillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23BillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicB23BillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicB23BillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23BillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicB23BillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicB23CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB23ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicB23ChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicB23DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicB23DedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicB23FeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicB23FeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicB23FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicB23FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicB23FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicB23InsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicB23NOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23NOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicB23NOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicB23OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicB23OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB23PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicB23Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicB23PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicB23PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicB23PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicB23PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicB23PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicB23PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicB23ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicB23StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicB23UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicB23Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicB23UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicB23UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicB23UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicB23UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicB23UpkeepCanceled, error)

	FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepChargedIterator, error)

	WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCharged, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicB23UpkeepCharged, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicB23UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicB23UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicB23UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicB23UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicB23UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicB23UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicB23UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicB23UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicB23UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicB23UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicB23UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicB23UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
