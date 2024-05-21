// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_automation_registry_master_wrapper_2_3

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

var IAutomationRegistryMaster23MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLinkLiquidity\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"overrides\",\"type\":\"tuple\"}],\"name\":\"BillingConfigOverridden\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"BillingConfigOverrideRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payments\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInJuels\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"linkUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"billingUSD\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"name\":\"receipt\",\"type\":\"tuple\"}],\"name\":\"UpkeepCharged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableOffchainPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"getAdminPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedReadOnlyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAutomationForwarderLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getAvailableERC20ForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getBillingConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingOverrides\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBillingTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCancellationDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainModule\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConditionalGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackNativePrice\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackNativePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHotVars\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuard\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.HotVars\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLogGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNativeUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumUpkeeps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPayoutMode\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerPerformByteGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerSignerGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReorgProtectionEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStorage\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"}],\"internalType\":\"structAutomationRegistryBase2_3.Storage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmittersWithPayees\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payeeAddress\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.TransmitterPayeeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNativeTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"removeBillingOverrides\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setAdminPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"billingOverrides\",\"type\":\"tuple\"}],\"name\":\"setBillingOverrides\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackNativePrice\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"address[]\",\"name\":\"billingTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig[]\",\"name\":\"billingConfigs\",\"type\":\"tuple[]\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settleNOPsOffchain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"supportsBillingToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20Fees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var IAutomationRegistryMaster23ABI = IAutomationRegistryMaster23MetaData.ABI

type IAutomationRegistryMaster23 struct {
	address common.Address
	abi     abi.ABI
	IAutomationRegistryMaster23Caller
	IAutomationRegistryMaster23Transactor
	IAutomationRegistryMaster23Filterer
}

type IAutomationRegistryMaster23Caller struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMaster23Transactor struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMaster23Filterer struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMaster23Session struct {
	Contract     *IAutomationRegistryMaster23
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAutomationRegistryMaster23CallerSession struct {
	Contract *IAutomationRegistryMaster23Caller
	CallOpts bind.CallOpts
}

type IAutomationRegistryMaster23TransactorSession struct {
	Contract     *IAutomationRegistryMaster23Transactor
	TransactOpts bind.TransactOpts
}

type IAutomationRegistryMaster23Raw struct {
	Contract *IAutomationRegistryMaster23
}

type IAutomationRegistryMaster23CallerRaw struct {
	Contract *IAutomationRegistryMaster23Caller
}

type IAutomationRegistryMaster23TransactorRaw struct {
	Contract *IAutomationRegistryMaster23Transactor
}

func NewIAutomationRegistryMaster23(address common.Address, backend bind.ContractBackend) (*IAutomationRegistryMaster23, error) {
	abi, err := abi.JSON(strings.NewReader(IAutomationRegistryMaster23ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAutomationRegistryMaster23(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23{address: address, abi: abi, IAutomationRegistryMaster23Caller: IAutomationRegistryMaster23Caller{contract: contract}, IAutomationRegistryMaster23Transactor: IAutomationRegistryMaster23Transactor{contract: contract}, IAutomationRegistryMaster23Filterer: IAutomationRegistryMaster23Filterer{contract: contract}}, nil
}

func NewIAutomationRegistryMaster23Caller(address common.Address, caller bind.ContractCaller) (*IAutomationRegistryMaster23Caller, error) {
	contract, err := bindIAutomationRegistryMaster23(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23Caller{contract: contract}, nil
}

func NewIAutomationRegistryMaster23Transactor(address common.Address, transactor bind.ContractTransactor) (*IAutomationRegistryMaster23Transactor, error) {
	contract, err := bindIAutomationRegistryMaster23(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23Transactor{contract: contract}, nil
}

func NewIAutomationRegistryMaster23Filterer(address common.Address, filterer bind.ContractFilterer) (*IAutomationRegistryMaster23Filterer, error) {
	contract, err := bindIAutomationRegistryMaster23(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23Filterer{contract: contract}, nil
}

func bindIAutomationRegistryMaster23(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAutomationRegistryMaster23MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationRegistryMaster23.Contract.IAutomationRegistryMaster23Caller.contract.Call(opts, result, method, params...)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.IAutomationRegistryMaster23Transactor.contract.Transfer(opts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.IAutomationRegistryMaster23Transactor.contract.Transact(opts, method, params...)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationRegistryMaster23.Contract.contract.Call(opts, result, method, params...)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.contract.Transfer(opts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.contract.Transact(opts, method, params...)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "checkCallback", id, values, extraData)

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

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckCallback(&_IAutomationRegistryMaster23.CallOpts, id, values, extraData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckCallback(&_IAutomationRegistryMaster23.CallOpts, id, values, extraData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "checkUpkeep", id, triggerData)

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
	outstruct.LinkUSD = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckUpkeep(&_IAutomationRegistryMaster23.CallOpts, id, triggerData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckUpkeep(&_IAutomationRegistryMaster23.CallOpts, id, triggerData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "checkUpkeep0", id)

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
	outstruct.LinkUSD = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckUpkeep0(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationRegistryMaster23.Contract.CheckUpkeep0(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) FallbackTo() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.FallbackTo(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) FallbackTo() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.FallbackTo(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetActiveUpkeepIDs(&_IAutomationRegistryMaster23.CallOpts, startIndex, maxCount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetActiveUpkeepIDs(&_IAutomationRegistryMaster23.CallOpts, startIndex, maxCount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetAdminPrivilegeConfig(&_IAutomationRegistryMaster23.CallOpts, admin)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetAdminPrivilegeConfig(&_IAutomationRegistryMaster23.CallOpts, admin)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetAllowedReadOnlyAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetAllowedReadOnlyAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetAutomationForwarderLogic() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetAutomationForwarderLogic(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetAutomationForwarderLogic(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetAvailableERC20ForPayment(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getAvailableERC20ForPayment", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetAvailableERC20ForPayment(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetAvailableERC20ForPayment(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetBalance(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetBalance(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBillingConfig(opts *bind.CallOpts, billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBillingConfig", billingToken)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingConfig(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingConfig(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBillingOverrides(opts *bind.CallOpts, upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBillingOverrides", upkeepID)

	if err != nil {
		return *new(AutomationRegistryBase23BillingOverrides), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingOverrides)).(*AutomationRegistryBase23BillingOverrides)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingOverrides(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingOverrides(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBillingToken", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingToken(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingToken(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBillingTokenConfig", token)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingTokenConfig(&_IAutomationRegistryMaster23.CallOpts, token)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingTokenConfig(&_IAutomationRegistryMaster23.CallOpts, token)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getBillingTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetBillingTokens() ([]common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingTokens(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetBillingTokens() ([]common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetBillingTokens(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetCancellationDelay() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetCancellationDelay(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetCancellationDelay() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetCancellationDelay(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetChainModule() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetChainModule(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetChainModule() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetChainModule(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetConditionalGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetConditionalGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetConditionalGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetConfig(opts *bind.CallOpts) (AutomationRegistryBase23OnchainConfig, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(AutomationRegistryBase23OnchainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23OnchainConfig)).(*AutomationRegistryBase23OnchainConfig)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetConfig(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _IAutomationRegistryMaster23.Contract.GetConfig(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getFallbackNativePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetFallbackNativePrice() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetFallbackNativePrice(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetFallbackNativePrice() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetFallbackNativePrice(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetFastGasFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetFastGasFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetFastGasFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetForwarder(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetForwarder(&_IAutomationRegistryMaster23.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getHotVars")

	if err != nil {
		return *new(AutomationRegistryBase23HotVars), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23HotVars)).(*AutomationRegistryBase23HotVars)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _IAutomationRegistryMaster23.Contract.GetHotVars(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _IAutomationRegistryMaster23.Contract.GetHotVars(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetLinkAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetLinkAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetLinkAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetLinkAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getLinkUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetLinkUSDFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetLinkUSDFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetLinkUSDFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetLogGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetLogGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetLogGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetMaxPaymentForGas(opts *bind.CallOpts, id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getMaxPaymentForGas", id, triggerType, gasLimit, billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMaxPaymentForGas(&_IAutomationRegistryMaster23.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMaxPaymentForGas(&_IAutomationRegistryMaster23.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMinBalance(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMinBalance(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMinBalanceForUpkeep(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetMinBalanceForUpkeep(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetNativeUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getNativeUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetNativeUSDFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetNativeUSDFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetNativeUSDFeedAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetNumUpkeeps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getNumUpkeeps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetNumUpkeeps() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetNumUpkeeps(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetNumUpkeeps() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetNumUpkeeps(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetPayoutMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getPayoutMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetPayoutMode() (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetPayoutMode(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetPayoutMode() (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetPayoutMode(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster23.CallOpts, peer)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster23.CallOpts, peer)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetPerPerformByteGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetPerPerformByteGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetPerSignerGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetPerSignerGasOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetReorgProtectionEnabled() (bool, error) {
	return _IAutomationRegistryMaster23.Contract.GetReorgProtectionEnabled(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _IAutomationRegistryMaster23.Contract.GetReorgProtectionEnabled(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetReserveAmount(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getReserveAmount", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetReserveAmount(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetReserveAmount(&_IAutomationRegistryMaster23.CallOpts, billingToken)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetSignerInfo(&_IAutomationRegistryMaster23.CallOpts, query)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetSignerInfo(&_IAutomationRegistryMaster23.CallOpts, query)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getState")

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

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetState() (GetState,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetState(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetState() (GetState,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetState(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetStorage(opts *bind.CallOpts) (AutomationRegistryBase23Storage, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getStorage")

	if err != nil {
		return *new(AutomationRegistryBase23Storage), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23Storage)).(*AutomationRegistryBase23Storage)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _IAutomationRegistryMaster23.Contract.GetStorage(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _IAutomationRegistryMaster23.Contract.GetStorage(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitCalldataFixedBytesOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitCalldataFixedBytesOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitterInfo(&_IAutomationRegistryMaster23.CallOpts, query)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmitterInfo(&_IAutomationRegistryMaster23.CallOpts, query)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetTransmittersWithPayees(opts *bind.CallOpts) ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getTransmittersWithPayees")

	if err != nil {
		return *new([]AutomationRegistryBase23TransmitterPayeeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutomationRegistryBase23TransmitterPayeeInfo)).(*[]AutomationRegistryBase23TransmitterPayeeInfo)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmittersWithPayees(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _IAutomationRegistryMaster23.Contract.GetTransmittersWithPayees(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetTriggerType(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.GetTriggerType(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeep(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeep(&_IAutomationRegistryMaster23.CallOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeepTriggerConfig(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster23.Contract.GetUpkeepTriggerConfig(&_IAutomationRegistryMaster23.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) GetWrappedNativeTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "getWrappedNativeTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetWrappedNativeTokenAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.GetWrappedNativeTokenAddress(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationRegistryMaster23.Contract.HasDedupKey(&_IAutomationRegistryMaster23.CallOpts, dedupKey)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationRegistryMaster23.Contract.HasDedupKey(&_IAutomationRegistryMaster23.CallOpts, dedupKey)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IAutomationRegistryMaster23.Contract.LatestConfigDetails(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IAutomationRegistryMaster23.Contract.LatestConfigDetails(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IAutomationRegistryMaster23.Contract.LatestConfigDigestAndEpoch(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IAutomationRegistryMaster23.Contract.LatestConfigDigestAndEpoch(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) LinkAvailableForPayment() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.LinkAvailableForPayment(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _IAutomationRegistryMaster23.Contract.LinkAvailableForPayment(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) Owner() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.Owner(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) Owner() (common.Address, error) {
	return _IAutomationRegistryMaster23.Contract.Owner(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "simulatePerformUpkeep", id, performData)

	outstruct := new(SimulatePerformUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GasUsed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationRegistryMaster23.Contract.SimulatePerformUpkeep(&_IAutomationRegistryMaster23.CallOpts, id, performData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationRegistryMaster23.Contract.SimulatePerformUpkeep(&_IAutomationRegistryMaster23.CallOpts, id, performData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "supportsBillingToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SupportsBillingToken(token common.Address) (bool, error) {
	return _IAutomationRegistryMaster23.Contract.SupportsBillingToken(&_IAutomationRegistryMaster23.CallOpts, token)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _IAutomationRegistryMaster23.Contract.SupportsBillingToken(&_IAutomationRegistryMaster23.CallOpts, token)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) TypeAndVersion() (string, error) {
	return _IAutomationRegistryMaster23.Contract.TypeAndVersion(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) TypeAndVersion() (string, error) {
	return _IAutomationRegistryMaster23.Contract.TypeAndVersion(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Caller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster23.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) UpkeepVersion() (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.UpkeepVersion(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23CallerSession) UpkeepVersion() (uint8, error) {
	return _IAutomationRegistryMaster23.Contract.UpkeepVersion(&_IAutomationRegistryMaster23.CallOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "acceptOwnership")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) AcceptOwnership() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptOwnership(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptOwnership(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptPayeeship(&_IAutomationRegistryMaster23.TransactOpts, transmitter)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptPayeeship(&_IAutomationRegistryMaster23.TransactOpts, transmitter)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptUpkeepAdmin(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AcceptUpkeepAdmin(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "addFunds", id, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AddFunds(&_IAutomationRegistryMaster23.TransactOpts, id, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.AddFunds(&_IAutomationRegistryMaster23.TransactOpts, id, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.CancelUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.CancelUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "disableOffchainPayments")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) DisableOffchainPayments() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.DisableOffchainPayments(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.DisableOffchainPayments(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.ExecuteCallback(&_IAutomationRegistryMaster23.TransactOpts, id, payload)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.ExecuteCallback(&_IAutomationRegistryMaster23.TransactOpts, id, payload)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.MigrateUpkeeps(&_IAutomationRegistryMaster23.TransactOpts, ids, destination)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.MigrateUpkeeps(&_IAutomationRegistryMaster23.TransactOpts, ids, destination)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.OnTokenTransfer(&_IAutomationRegistryMaster23.TransactOpts, sender, amount, data)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.OnTokenTransfer(&_IAutomationRegistryMaster23.TransactOpts, sender, amount, data)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "pause")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) Pause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Pause(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) Pause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Pause(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.PauseUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.PauseUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.ReceiveUpkeeps(&_IAutomationRegistryMaster23.TransactOpts, encodedUpkeeps)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.ReceiveUpkeeps(&_IAutomationRegistryMaster23.TransactOpts, encodedUpkeeps)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, billingToken common.Address, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, billingToken, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, billingToken common.Address, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.RegisterUpkeep(&_IAutomationRegistryMaster23.TransactOpts, target, gasLimit, admin, triggerType, billingToken, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, billingToken common.Address, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.RegisterUpkeep(&_IAutomationRegistryMaster23.TransactOpts, target, gasLimit, admin, triggerType, billingToken, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) RemoveBillingOverrides(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "removeBillingOverrides", id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.RemoveBillingOverrides(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) RemoveBillingOverrides(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.RemoveBillingOverrides(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetAdminPrivilegeConfig(&_IAutomationRegistryMaster23.TransactOpts, admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetAdminPrivilegeConfig(&_IAutomationRegistryMaster23.TransactOpts, admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetBillingOverrides(opts *bind.TransactOpts, id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setBillingOverrides", id, billingOverrides)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetBillingOverrides(&_IAutomationRegistryMaster23.TransactOpts, id, billingOverrides)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetBillingOverrides(id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetBillingOverrides(&_IAutomationRegistryMaster23.TransactOpts, id, billingOverrides)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetConfig(&_IAutomationRegistryMaster23.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetConfig(&_IAutomationRegistryMaster23.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetConfigTypeSafe(&_IAutomationRegistryMaster23.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetConfigTypeSafe(&_IAutomationRegistryMaster23.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, billingTokens, billingConfigs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setPayees", payees)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetPayees(&_IAutomationRegistryMaster23.TransactOpts, payees)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetPayees(&_IAutomationRegistryMaster23.TransactOpts, payees)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster23.TransactOpts, peer, permission)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster23.TransactOpts, peer, permission)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepCheckData(&_IAutomationRegistryMaster23.TransactOpts, id, newCheckData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepCheckData(&_IAutomationRegistryMaster23.TransactOpts, id, newCheckData)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepGasLimit(&_IAutomationRegistryMaster23.TransactOpts, id, gasLimit)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepGasLimit(&_IAutomationRegistryMaster23.TransactOpts, id, gasLimit)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepOffchainConfig(&_IAutomationRegistryMaster23.TransactOpts, id, config)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepOffchainConfig(&_IAutomationRegistryMaster23.TransactOpts, id, config)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster23.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster23.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepTriggerConfig(&_IAutomationRegistryMaster23.TransactOpts, id, triggerConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SetUpkeepTriggerConfig(&_IAutomationRegistryMaster23.TransactOpts, id, triggerConfig)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "settleNOPsOffchain")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) SettleNOPsOffchain() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SettleNOPsOffchain(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.SettleNOPsOffchain(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "transferOwnership", to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferOwnership(&_IAutomationRegistryMaster23.TransactOpts, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferOwnership(&_IAutomationRegistryMaster23.TransactOpts, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferPayeeship(&_IAutomationRegistryMaster23.TransactOpts, transmitter, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferPayeeship(&_IAutomationRegistryMaster23.TransactOpts, transmitter, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferUpkeepAdmin(&_IAutomationRegistryMaster23.TransactOpts, id, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.TransferUpkeepAdmin(&_IAutomationRegistryMaster23.TransactOpts, id, proposed)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Transmit(&_IAutomationRegistryMaster23.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Transmit(&_IAutomationRegistryMaster23.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "unpause")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) Unpause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Unpause(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) Unpause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Unpause(&_IAutomationRegistryMaster23.TransactOpts)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.UnpauseUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.UnpauseUpkeep(&_IAutomationRegistryMaster23.TransactOpts, id)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) WithdrawERC20Fees(opts *bind.TransactOpts, asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "withdrawERC20Fees", asset, to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawERC20Fees(&_IAutomationRegistryMaster23.TransactOpts, asset, to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) WithdrawERC20Fees(asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawERC20Fees(&_IAutomationRegistryMaster23.TransactOpts, asset, to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawFunds(&_IAutomationRegistryMaster23.TransactOpts, id, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawFunds(&_IAutomationRegistryMaster23.TransactOpts, id, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) WithdrawLink(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "withdrawLink", to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawLink(&_IAutomationRegistryMaster23.TransactOpts, to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) WithdrawLink(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawLink(&_IAutomationRegistryMaster23.TransactOpts, to, amount)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawPayment(&_IAutomationRegistryMaster23.TransactOpts, from, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.WithdrawPayment(&_IAutomationRegistryMaster23.TransactOpts, from, to)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.contract.RawTransact(opts, calldata)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Fallback(&_IAutomationRegistryMaster23.TransactOpts, calldata)
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster23.Contract.Fallback(&_IAutomationRegistryMaster23.TransactOpts, calldata)
}

type IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator struct {
	Event *IAutomationRegistryMaster23AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23AdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23AdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23AdminPrivilegeConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMaster23AdminPrivilegeConfigSet, error) {
	event := new(IAutomationRegistryMaster23AdminPrivilegeConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23BillingConfigOverriddenIterator struct {
	Event *IAutomationRegistryMaster23BillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23BillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23BillingConfigOverridden)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23BillingConfigOverridden)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23BillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23BillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23BillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23BillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23BillingConfigOverriddenIterator{contract: _IAutomationRegistryMaster23.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23BillingConfigOverridden)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseBillingConfigOverridden(log types.Log) (*IAutomationRegistryMaster23BillingConfigOverridden, error) {
	event := new(IAutomationRegistryMaster23BillingConfigOverridden)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator struct {
	Event *IAutomationRegistryMaster23BillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23BillingConfigOverrideRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23BillingConfigOverrideRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23BillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator{contract: _IAutomationRegistryMaster23.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23BillingConfigOverrideRemoved)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseBillingConfigOverrideRemoved(log types.Log) (*IAutomationRegistryMaster23BillingConfigOverrideRemoved, error) {
	event := new(IAutomationRegistryMaster23BillingConfigOverrideRemoved)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23BillingConfigSetIterator struct {
	Event *IAutomationRegistryMaster23BillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23BillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23BillingConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23BillingConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23BillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23BillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23BillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*IAutomationRegistryMaster23BillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23BillingConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23BillingConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseBillingConfigSet(log types.Log) (*IAutomationRegistryMaster23BillingConfigSet, error) {
	event := new(IAutomationRegistryMaster23BillingConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23CancelledUpkeepReportIterator struct {
	Event *IAutomationRegistryMaster23CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23CancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23CancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23CancelledUpkeepReportIterator{contract: _IAutomationRegistryMaster23.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23CancelledUpkeepReport)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseCancelledUpkeepReport(log types.Log) (*IAutomationRegistryMaster23CancelledUpkeepReport, error) {
	event := new(IAutomationRegistryMaster23CancelledUpkeepReport)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator struct {
	Event *IAutomationRegistryMaster23ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23ChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23ChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator{contract: _IAutomationRegistryMaster23.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23ChainSpecificModuleUpdated)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*IAutomationRegistryMaster23ChainSpecificModuleUpdated, error) {
	event := new(IAutomationRegistryMaster23ChainSpecificModuleUpdated)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23ConfigSetIterator struct {
	Event *IAutomationRegistryMaster23ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23ConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23ConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23ConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23ConfigSet struct {
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

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterConfigSet(opts *bind.FilterOpts) (*IAutomationRegistryMaster23ConfigSetIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23ConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ConfigSet) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23ConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseConfigSet(log types.Log) (*IAutomationRegistryMaster23ConfigSet, error) {
	event := new(IAutomationRegistryMaster23ConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23DedupKeyAddedIterator struct {
	Event *IAutomationRegistryMaster23DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23DedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23DedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationRegistryMaster23DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23DedupKeyAddedIterator{contract: _IAutomationRegistryMaster23.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23DedupKeyAdded)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseDedupKeyAdded(log types.Log) (*IAutomationRegistryMaster23DedupKeyAdded, error) {
	event := new(IAutomationRegistryMaster23DedupKeyAdded)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23FeesWithdrawnIterator struct {
	Event *IAutomationRegistryMaster23FeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23FeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23FeesWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23FeesWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23FeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23FeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23FeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*IAutomationRegistryMaster23FeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23FeesWithdrawnIterator{contract: _IAutomationRegistryMaster23.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23FeesWithdrawn)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseFeesWithdrawn(log types.Log) (*IAutomationRegistryMaster23FeesWithdrawn, error) {
	event := new(IAutomationRegistryMaster23FeesWithdrawn)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23FundsAddedIterator struct {
	Event *IAutomationRegistryMaster23FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23FundsAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationRegistryMaster23FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23FundsAddedIterator{contract: _IAutomationRegistryMaster23.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23FundsAdded)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseFundsAdded(log types.Log) (*IAutomationRegistryMaster23FundsAdded, error) {
	event := new(IAutomationRegistryMaster23FundsAdded)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23FundsWithdrawnIterator struct {
	Event *IAutomationRegistryMaster23FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23FundsWithdrawnIterator{contract: _IAutomationRegistryMaster23.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23FundsWithdrawn)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseFundsWithdrawn(log types.Log) (*IAutomationRegistryMaster23FundsWithdrawn, error) {
	event := new(IAutomationRegistryMaster23FundsWithdrawn)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator struct {
	Event *IAutomationRegistryMaster23InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23InsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23InsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator{contract: _IAutomationRegistryMaster23.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23InsufficientFundsUpkeepReport)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationRegistryMaster23InsufficientFundsUpkeepReport, error) {
	event := new(IAutomationRegistryMaster23InsufficientFundsUpkeepReport)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23NOPsSettledOffchainIterator struct {
	Event *IAutomationRegistryMaster23NOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23NOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23NOPsSettledOffchain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23NOPsSettledOffchain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23NOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23NOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23NOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*IAutomationRegistryMaster23NOPsSettledOffchainIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23NOPsSettledOffchainIterator{contract: _IAutomationRegistryMaster23.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23NOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23NOPsSettledOffchain)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseNOPsSettledOffchain(log types.Log) (*IAutomationRegistryMaster23NOPsSettledOffchain, error) {
	event := new(IAutomationRegistryMaster23NOPsSettledOffchain)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23OwnershipTransferRequestedIterator struct {
	Event *IAutomationRegistryMaster23OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23OwnershipTransferRequestedIterator{contract: _IAutomationRegistryMaster23.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23OwnershipTransferRequested)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseOwnershipTransferRequested(log types.Log) (*IAutomationRegistryMaster23OwnershipTransferRequested, error) {
	event := new(IAutomationRegistryMaster23OwnershipTransferRequested)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23OwnershipTransferredIterator struct {
	Event *IAutomationRegistryMaster23OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23OwnershipTransferredIterator{contract: _IAutomationRegistryMaster23.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23OwnershipTransferred)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseOwnershipTransferred(log types.Log) (*IAutomationRegistryMaster23OwnershipTransferred, error) {
	event := new(IAutomationRegistryMaster23OwnershipTransferred)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23PausedIterator struct {
	Event *IAutomationRegistryMaster23Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23PausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterPaused(opts *bind.FilterOpts) (*IAutomationRegistryMaster23PausedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23PausedIterator{contract: _IAutomationRegistryMaster23.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Paused) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23Paused)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParsePaused(log types.Log) (*IAutomationRegistryMaster23Paused, error) {
	event := new(IAutomationRegistryMaster23Paused)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23PayeesUpdatedIterator struct {
	Event *IAutomationRegistryMaster23PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23PayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23PayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMaster23PayeesUpdatedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23PayeesUpdatedIterator{contract: _IAutomationRegistryMaster23.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23PayeesUpdated)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParsePayeesUpdated(log types.Log) (*IAutomationRegistryMaster23PayeesUpdated, error) {
	event := new(IAutomationRegistryMaster23PayeesUpdated)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23PayeeshipTransferRequestedIterator struct {
	Event *IAutomationRegistryMaster23PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23PayeeshipTransferRequestedIterator{contract: _IAutomationRegistryMaster23.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23PayeeshipTransferRequested)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParsePayeeshipTransferRequested(log types.Log) (*IAutomationRegistryMaster23PayeeshipTransferRequested, error) {
	event := new(IAutomationRegistryMaster23PayeeshipTransferRequested)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23PayeeshipTransferredIterator struct {
	Event *IAutomationRegistryMaster23PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23PayeeshipTransferredIterator{contract: _IAutomationRegistryMaster23.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23PayeeshipTransferred)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParsePayeeshipTransferred(log types.Log) (*IAutomationRegistryMaster23PayeeshipTransferred, error) {
	event := new(IAutomationRegistryMaster23PayeeshipTransferred)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23PaymentWithdrawnIterator struct {
	Event *IAutomationRegistryMaster23PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationRegistryMaster23PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23PaymentWithdrawnIterator{contract: _IAutomationRegistryMaster23.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23PaymentWithdrawn)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParsePaymentWithdrawn(log types.Log) (*IAutomationRegistryMaster23PaymentWithdrawn, error) {
	event := new(IAutomationRegistryMaster23PaymentWithdrawn)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23ReorgedUpkeepReportIterator struct {
	Event *IAutomationRegistryMaster23ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23ReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23ReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23ReorgedUpkeepReportIterator{contract: _IAutomationRegistryMaster23.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23ReorgedUpkeepReport)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseReorgedUpkeepReport(log types.Log) (*IAutomationRegistryMaster23ReorgedUpkeepReport, error) {
	event := new(IAutomationRegistryMaster23ReorgedUpkeepReport)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23StaleUpkeepReportIterator struct {
	Event *IAutomationRegistryMaster23StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23StaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23StaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23StaleUpkeepReportIterator{contract: _IAutomationRegistryMaster23.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23StaleUpkeepReport)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseStaleUpkeepReport(log types.Log) (*IAutomationRegistryMaster23StaleUpkeepReport, error) {
	event := new(IAutomationRegistryMaster23StaleUpkeepReport)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23TransmittedIterator struct {
	Event *IAutomationRegistryMaster23Transmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23TransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23Transmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23Transmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23TransmittedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23Transmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterTransmitted(opts *bind.FilterOpts) (*IAutomationRegistryMaster23TransmittedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23TransmittedIterator{contract: _IAutomationRegistryMaster23.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Transmitted) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23Transmitted)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseTransmitted(log types.Log) (*IAutomationRegistryMaster23Transmitted, error) {
	event := new(IAutomationRegistryMaster23Transmitted)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UnpausedIterator struct {
	Event *IAutomationRegistryMaster23Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUnpaused(opts *bind.FilterOpts) (*IAutomationRegistryMaster23UnpausedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UnpausedIterator{contract: _IAutomationRegistryMaster23.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Unpaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23Unpaused)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUnpaused(log types.Log) (*IAutomationRegistryMaster23Unpaused, error) {
	event := new(IAutomationRegistryMaster23Unpaused)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepAdminTransferRequested)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationRegistryMaster23UpkeepAdminTransferRequested, error) {
	event := new(IAutomationRegistryMaster23UpkeepAdminTransferRequested)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepAdminTransferredIterator struct {
	Event *IAutomationRegistryMaster23UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepAdminTransferredIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepAdminTransferred)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepAdminTransferred(log types.Log) (*IAutomationRegistryMaster23UpkeepAdminTransferred, error) {
	event := new(IAutomationRegistryMaster23UpkeepAdminTransferred)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepCanceledIterator struct {
	Event *IAutomationRegistryMaster23UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationRegistryMaster23UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepCanceledIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepCanceled)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepCanceled(log types.Log) (*IAutomationRegistryMaster23UpkeepCanceled, error) {
	event := new(IAutomationRegistryMaster23UpkeepCanceled)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepChargedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepCharged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepChargedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepCharged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepCharged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepChargedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepChargedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepCharged struct {
	Id      *big.Int
	Receipt AutomationRegistryBase23PaymentReceipt
	Raw     types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepChargedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepChargedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepCharged", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCharged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepCharged)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepCharged(log types.Log) (*IAutomationRegistryMaster23UpkeepCharged, error) {
	event := new(IAutomationRegistryMaster23UpkeepCharged)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepCheckDataSetIterator struct {
	Event *IAutomationRegistryMaster23UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepCheckDataSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepCheckDataSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepCheckDataSet(log types.Log) (*IAutomationRegistryMaster23UpkeepCheckDataSet, error) {
	event := new(IAutomationRegistryMaster23UpkeepCheckDataSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepGasLimitSetIterator struct {
	Event *IAutomationRegistryMaster23UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepGasLimitSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepGasLimitSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepGasLimitSet(log types.Log) (*IAutomationRegistryMaster23UpkeepGasLimitSet, error) {
	event := new(IAutomationRegistryMaster23UpkeepGasLimitSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepMigratedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepMigratedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepMigrated)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepMigrated(log types.Log) (*IAutomationRegistryMaster23UpkeepMigrated, error) {
	event := new(IAutomationRegistryMaster23UpkeepMigrated)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator struct {
	Event *IAutomationRegistryMaster23UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepOffchainConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepOffchainConfigSet, error) {
	event := new(IAutomationRegistryMaster23UpkeepOffchainConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepPausedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepPausedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepPaused)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepPaused(log types.Log) (*IAutomationRegistryMaster23UpkeepPaused, error) {
	event := new(IAutomationRegistryMaster23UpkeepPaused)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepPerformedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationRegistryMaster23UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepPerformedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepPerformed)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepPerformed(log types.Log) (*IAutomationRegistryMaster23UpkeepPerformed, error) {
	event := new(IAutomationRegistryMaster23UpkeepPerformed)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator struct {
	Event *IAutomationRegistryMaster23UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepPrivilegeConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepPrivilegeConfigSet, error) {
	event := new(IAutomationRegistryMaster23UpkeepPrivilegeConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepReceivedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepReceivedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepReceived)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepReceived(log types.Log) (*IAutomationRegistryMaster23UpkeepReceived, error) {
	event := new(IAutomationRegistryMaster23UpkeepReceived)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepRegisteredIterator struct {
	Event *IAutomationRegistryMaster23UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepRegisteredIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepRegistered)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepRegistered(log types.Log) (*IAutomationRegistryMaster23UpkeepRegistered, error) {
	event := new(IAutomationRegistryMaster23UpkeepRegistered)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator struct {
	Event *IAutomationRegistryMaster23UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepTriggerConfigSet)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepTriggerConfigSet, error) {
	event := new(IAutomationRegistryMaster23UpkeepTriggerConfigSet)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMaster23UpkeepUnpausedIterator struct {
	Event *IAutomationRegistryMaster23UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMaster23UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMaster23UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMaster23UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMaster23UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMaster23UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMaster23UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster23UpkeepUnpausedIterator{contract: _IAutomationRegistryMaster23.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster23.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMaster23UpkeepUnpaused)
				if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23Filterer) ParseUpkeepUnpaused(log types.Log) (*IAutomationRegistryMaster23UpkeepUnpaused, error) {
	event := new(IAutomationRegistryMaster23UpkeepUnpaused)
	if err := _IAutomationRegistryMaster23.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
	LinkUSD             *big.Int
}
type CheckUpkeep0 struct {
	UpkeepNeeded        bool
	PerformData         []byte
	UpkeepFailureReason uint8
	GasUsed             *big.Int
	GasLimit            *big.Int
	FastGasWei          *big.Int
	LinkUSD             *big.Int
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

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IAutomationRegistryMaster23.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseAdminPrivilegeConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["BillingConfigOverridden"].ID:
		return _IAutomationRegistryMaster23.ParseBillingConfigOverridden(log)
	case _IAutomationRegistryMaster23.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _IAutomationRegistryMaster23.ParseBillingConfigOverrideRemoved(log)
	case _IAutomationRegistryMaster23.abi.Events["BillingConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseBillingConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["CancelledUpkeepReport"].ID:
		return _IAutomationRegistryMaster23.ParseCancelledUpkeepReport(log)
	case _IAutomationRegistryMaster23.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _IAutomationRegistryMaster23.ParseChainSpecificModuleUpdated(log)
	case _IAutomationRegistryMaster23.abi.Events["ConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["DedupKeyAdded"].ID:
		return _IAutomationRegistryMaster23.ParseDedupKeyAdded(log)
	case _IAutomationRegistryMaster23.abi.Events["FeesWithdrawn"].ID:
		return _IAutomationRegistryMaster23.ParseFeesWithdrawn(log)
	case _IAutomationRegistryMaster23.abi.Events["FundsAdded"].ID:
		return _IAutomationRegistryMaster23.ParseFundsAdded(log)
	case _IAutomationRegistryMaster23.abi.Events["FundsWithdrawn"].ID:
		return _IAutomationRegistryMaster23.ParseFundsWithdrawn(log)
	case _IAutomationRegistryMaster23.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IAutomationRegistryMaster23.ParseInsufficientFundsUpkeepReport(log)
	case _IAutomationRegistryMaster23.abi.Events["NOPsSettledOffchain"].ID:
		return _IAutomationRegistryMaster23.ParseNOPsSettledOffchain(log)
	case _IAutomationRegistryMaster23.abi.Events["OwnershipTransferRequested"].ID:
		return _IAutomationRegistryMaster23.ParseOwnershipTransferRequested(log)
	case _IAutomationRegistryMaster23.abi.Events["OwnershipTransferred"].ID:
		return _IAutomationRegistryMaster23.ParseOwnershipTransferred(log)
	case _IAutomationRegistryMaster23.abi.Events["Paused"].ID:
		return _IAutomationRegistryMaster23.ParsePaused(log)
	case _IAutomationRegistryMaster23.abi.Events["PayeesUpdated"].ID:
		return _IAutomationRegistryMaster23.ParsePayeesUpdated(log)
	case _IAutomationRegistryMaster23.abi.Events["PayeeshipTransferRequested"].ID:
		return _IAutomationRegistryMaster23.ParsePayeeshipTransferRequested(log)
	case _IAutomationRegistryMaster23.abi.Events["PayeeshipTransferred"].ID:
		return _IAutomationRegistryMaster23.ParsePayeeshipTransferred(log)
	case _IAutomationRegistryMaster23.abi.Events["PaymentWithdrawn"].ID:
		return _IAutomationRegistryMaster23.ParsePaymentWithdrawn(log)
	case _IAutomationRegistryMaster23.abi.Events["ReorgedUpkeepReport"].ID:
		return _IAutomationRegistryMaster23.ParseReorgedUpkeepReport(log)
	case _IAutomationRegistryMaster23.abi.Events["StaleUpkeepReport"].ID:
		return _IAutomationRegistryMaster23.ParseStaleUpkeepReport(log)
	case _IAutomationRegistryMaster23.abi.Events["Transmitted"].ID:
		return _IAutomationRegistryMaster23.ParseTransmitted(log)
	case _IAutomationRegistryMaster23.abi.Events["Unpaused"].ID:
		return _IAutomationRegistryMaster23.ParseUnpaused(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepAdminTransferRequested(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepAdminTransferred"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepAdminTransferred(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepCanceled"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepCanceled(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepCharged"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepCharged(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepCheckDataSet"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepCheckDataSet(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepGasLimitSet"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepGasLimitSet(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepMigrated"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepMigrated(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepOffchainConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepPaused"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepPaused(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepPerformed"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepPerformed(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepPrivilegeConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepReceived"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepReceived(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepRegistered"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepRegistered(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepTriggerConfigSet(log)
	case _IAutomationRegistryMaster23.abi.Events["UpkeepUnpaused"].ID:
		return _IAutomationRegistryMaster23.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IAutomationRegistryMaster23AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (IAutomationRegistryMaster23BillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (IAutomationRegistryMaster23BillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (IAutomationRegistryMaster23BillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba3")
}

func (IAutomationRegistryMaster23CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IAutomationRegistryMaster23ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (IAutomationRegistryMaster23ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IAutomationRegistryMaster23DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IAutomationRegistryMaster23FeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (IAutomationRegistryMaster23FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IAutomationRegistryMaster23FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IAutomationRegistryMaster23InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IAutomationRegistryMaster23NOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (IAutomationRegistryMaster23OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IAutomationRegistryMaster23OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IAutomationRegistryMaster23Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IAutomationRegistryMaster23PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IAutomationRegistryMaster23PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IAutomationRegistryMaster23PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IAutomationRegistryMaster23PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IAutomationRegistryMaster23ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IAutomationRegistryMaster23StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IAutomationRegistryMaster23Transmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IAutomationRegistryMaster23Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IAutomationRegistryMaster23UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IAutomationRegistryMaster23UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IAutomationRegistryMaster23UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IAutomationRegistryMaster23UpkeepCharged) Topic() common.Hash {
	return common.HexToHash("0x801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc")
}

func (IAutomationRegistryMaster23UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IAutomationRegistryMaster23UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IAutomationRegistryMaster23UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IAutomationRegistryMaster23UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IAutomationRegistryMaster23UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IAutomationRegistryMaster23UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IAutomationRegistryMaster23UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (IAutomationRegistryMaster23UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IAutomationRegistryMaster23UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IAutomationRegistryMaster23UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IAutomationRegistryMaster23UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IAutomationRegistryMaster23 *IAutomationRegistryMaster23) Address() common.Address {
	return _IAutomationRegistryMaster23.address
}

type IAutomationRegistryMaster23Interface interface {
	CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

		error)

	CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

		error)

	CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

		error)

	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error)

	GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error)

	GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error)

	GetAvailableERC20ForPayment(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetBillingConfig(opts *bind.CallOpts, billingToken common.Address) (AutomationRegistryBase23BillingConfig, error)

	GetBillingOverrides(opts *bind.CallOpts, upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error)

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

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

		error)

	SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, billingToken common.Address, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	RemoveBillingOverrides(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetBillingOverrides(opts *bind.TransactOpts, id *big.Int, billingOverrides AutomationRegistryBase23BillingOverrides) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase23OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte, billingTokens []common.Address, billingConfigs []AutomationRegistryBase23BillingConfig) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawERC20Fees(opts *bind.TransactOpts, asset common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationRegistryMaster23AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMaster23AdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23BillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*IAutomationRegistryMaster23BillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23BillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*IAutomationRegistryMaster23BillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*IAutomationRegistryMaster23BillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23BillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*IAutomationRegistryMaster23BillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IAutomationRegistryMaster23CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMaster23ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*IAutomationRegistryMaster23ChainSpecificModuleUpdated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IAutomationRegistryMaster23ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IAutomationRegistryMaster23ConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationRegistryMaster23DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IAutomationRegistryMaster23DedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*IAutomationRegistryMaster23FeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*IAutomationRegistryMaster23FeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationRegistryMaster23FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IAutomationRegistryMaster23FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IAutomationRegistryMaster23FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationRegistryMaster23InsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*IAutomationRegistryMaster23NOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23NOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*IAutomationRegistryMaster23NOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IAutomationRegistryMaster23OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IAutomationRegistryMaster23OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IAutomationRegistryMaster23PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IAutomationRegistryMaster23Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMaster23PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IAutomationRegistryMaster23PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IAutomationRegistryMaster23PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IAutomationRegistryMaster23PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationRegistryMaster23PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IAutomationRegistryMaster23PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IAutomationRegistryMaster23ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IAutomationRegistryMaster23StaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IAutomationRegistryMaster23TransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Transmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IAutomationRegistryMaster23Transmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IAutomationRegistryMaster23UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IAutomationRegistryMaster23Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationRegistryMaster23UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMaster23UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IAutomationRegistryMaster23UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationRegistryMaster23UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IAutomationRegistryMaster23UpkeepCanceled, error)

	FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepChargedIterator, error)

	WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCharged, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCharged(log types.Log) (*IAutomationRegistryMaster23UpkeepCharged, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IAutomationRegistryMaster23UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IAutomationRegistryMaster23UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IAutomationRegistryMaster23UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IAutomationRegistryMaster23UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationRegistryMaster23UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IAutomationRegistryMaster23UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IAutomationRegistryMaster23UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IAutomationRegistryMaster23UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationRegistryMaster23UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMaster23UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMaster23UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IAutomationRegistryMaster23UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
