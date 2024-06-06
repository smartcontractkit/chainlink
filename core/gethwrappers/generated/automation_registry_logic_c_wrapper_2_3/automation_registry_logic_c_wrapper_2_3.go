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

var AutomationRegistryLogicCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nativeUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"automationForwarderLogic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedReadOnlyAddress\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"payoutMode\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"wrappedNativeTokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLinkLiquidity\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"overrides\",\"type\":\"tuple\"}],\"name\":\"BillingConfigOverridden\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"BillingConfigOverrideRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payments\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInJuels\",\"type\":\"uint96\"},{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"linkUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"billingUSD\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"name\":\"receipt\",\"type\":\"tuple\"}],\"name\":\"UpkeepCharged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableOffchainPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"getAdminPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedReadOnlyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAutomationForwarderLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getAvailableERC20ForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getBillingConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingOverrides\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingToken\",\"outputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBillingTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingTokens\",\"outputs\":[{\"internalType\":\"contractIERC20Metadata[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCancellationDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainModule\",\"outputs\":[{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConditionalGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackNativePrice\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackNativePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"contractIAutomationForwarder\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHotVars\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuard\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.HotVars\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLogGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNativeUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumUpkeeps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPayoutMode\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerPerformByteGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerSignerGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReorgProtectionEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStorage\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"}],\"internalType\":\"structAutomationRegistryBase2_3.Storage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmittersWithPayees\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payeeAddress\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.TransmitterPayeeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNativeTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setAdminPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settleNOPsOffchain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"supportsBillingToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b5060405162005c5038038062005c50833981016040819052620000359162000309565b87878787878787873380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c78162000241565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff191660018381811115620001185762000118620003b6565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000179573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200019f9190620003cc565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002099190620003cc565b60ff16146200022b576040516301f86e1760e41b815260040160405180910390fd5b50505050505050505050505050505050620003f8565b336001600160a01b038216036200029b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200030457600080fd5b919050565b600080600080600080600080610100898b0312156200032757600080fd5b6200033289620002ec565b97506200034260208a01620002ec565b96506200035260408a01620002ec565b95506200036260608a01620002ec565b94506200037260808a01620002ec565b93506200038260a08a01620002ec565b925060c0890151600281106200039757600080fd5b9150620003a760e08a01620002ec565b90509295985092959890939650565b634e487b7160e01b600052602160045260246000fd5b600060208284031215620003df57600080fd5b815160ff81168114620003f157600080fd5b9392505050565b60805160a05160c05160e0516101005161012051610140516157cc620004846000396000610b3901526000610a9b01526000610787015260008181610969015261357f01526000818161091d0152613dfb0152600081816104da0152613659015260008181610be101528181611db901528181612534015281816125910152613abf01526157cc6000f3fe608060405234801561001057600080fd5b50600436106103c55760003560e01c80638ed02bab116101ff578063c3f909d41161011a578063eb5dcd6c116100ad578063f5a418461161007c578063f5a4184614610e7d578063f777ff0614610efe578063faa3e99614610f05578063ffd242bd14610f4b57600080fd5b8063eb5dcd6c14610dae578063ec4de4ba14610dc1578063ed56b3e114610df7578063f2fde38b14610e6a57600080fd5b8063d09dc339116100e9578063d09dc33914610c05578063d763264814610c0d578063d85aa07c14610c20578063e80a3d4414610c2857600080fd5b8063c3f909d414610b9f578063c5b964e014610bb4578063c7c3a19a14610bbf578063ca30e60314610bdf57600080fd5b8063aab9edd611610192578063b3596c2311610161578063b3596c23146107ab578063b6511a2a14610b70578063b657bc9c14610b77578063ba87666814610b8a57600080fd5b8063aab9edd614610b20578063abc76ae014610b2f578063ac4dc59a14610b37578063b121e14714610b5d57600080fd5b8063a08714c0116101ce578063a08714c014610a99578063a538b2eb14610abf578063a710b22114610b05578063a87f45fe14610b1857600080fd5b80638ed02bab14610a255780639089daa414610a4357806393f6ebcf14610a585780639e0a99ed14610a9157600080fd5b806343cc055c116102ef5780636209e1e91161028257806379ba50971161025157806379ba5097146109b357806379ea9943146109bb5780638456cb59146109ff5780638da5cb5b14610a0757600080fd5b80636209e1e9146109545780636709d0e514610967578063671d36ed1461098d5780636eec02a2146109a057600080fd5b80635425d8ac116102be5780635425d8ac1461078557806357359584146107ab578063614486af1461091b5780636181d82d1461094157600080fd5b806343cc055c1461070657806344cb70b8146107395780634ca16c521461075c5780635147cd591461076557600080fd5b8063207b6516116103675780633b9cce59116103365780633b9cce591461067d5780633f4ba83a14610690578063421d183b1461069857806343b46e5f146106fe57600080fd5b8063207b6516146104c5578063226cf83c146104d8578063232c1cc51461051f5780633408f73a1461052657600080fd5b80631865c57d116103a35780631865c57d14610417578063187256e81461043057806319d97a94146104435780631e0104391461046357600080fd5b8063050ee65d146103ca57806306e3b632146103e25780630b7d33e614610402575b600080fd5b6201de845b6040519081526020015b60405180910390f35b6103f56103f0366004614443565b610f53565b6040516103d991906144a0565b6104156104103660046144fc565b611070565b005b61041f6110d1565b6040516103d99594939291906146f4565b61041561043e366004614835565b6114d1565b610456610451366004614872565b611542565b6040516103d991906148ef565b6104a8610471366004614872565b60009081526004602052604090206001015470010000000000000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff90911681526020016103d9565b6104566104d3366004614872565b6115e4565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016103d9565b60186103cf565b6106706040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081019190915250604080516101608101825260165473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008084048216602086015278010000000000000000000000000000000000000000000000008085048316968601969096527c010000000000000000000000000000000000000000000000000000000093849004821660608601526017548084166080870152818104831660a0870152868104831660c087015293909304811660e085015260185491821661010085015291810482166101208401529290920490911661014082015290565b6040516103d99190614902565b61041561068b366004614a2c565b611601565b610415611857565b6106ab6106a6366004614aa1565b6118bd565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a0016103d9565b6104156119dc565b6014547801000000000000000000000000000000000000000000000000900460ff165b60405190151581526020016103d9565b610729610747366004614872565b60009081526008602052604090205460ff1690565b62017f986103cf565b610778610773366004614872565b611e73565b6040516103d99190614afd565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b61089d6107b9366004614aa1565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091525073ffffffffffffffffffffffffffffffffffffffff908116600090815260226020908152604091829020825160c081018452815463ffffffff81168252640100000000810462ffffff16938201939093526701000000000000008304909416928401929092527b01000000000000000000000000000000000000000000000000000000900460ff16606083015260018101546080830152600201546bffffffffffffffffffffffff1660a082015290565b6040516103d99190600060c08201905063ffffffff835116825262ffffff602084015116602083015273ffffffffffffffffffffffffffffffffffffffff604084015116604083015260ff6060840151166060830152608083015160808301526bffffffffffffffffffffffff60a08401511660a083015292915050565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b6104a861094f366004614b10565b611e7e565b610456610962366004614aa1565b611fd9565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b61041561099b366004614b70565b61200c565b6103cf6109ae366004614aa1565b61208c565b610415612136565b6104fa6109c9366004614872565b6000908152600460205260409020546a0100000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b610415612238565b60005473ffffffffffffffffffffffffffffffffffffffff166104fa565b60155473ffffffffffffffffffffffffffffffffffffffff166104fa565b610a4b6122b1565b6040516103d99190614bac565b6104fa610a66366004614872565b60009081526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff1690565b6103a46103cf565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b610729610acd366004614aa1565b73ffffffffffffffffffffffffffffffffffffffff908116600090815260226020526040902054670100000000000000900416151590565b610415610b13366004614c15565b6123c3565b6104156126f1565b604051600481526020016103d9565b6115e06103cf565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b610415610b6b366004614aa1565b612723565b60326103cf565b6104a8610b85366004614872565b61281b565b610b92612942565b6040516103d99190614c43565b610ba76129b1565b6040516103d99190614c9d565b60255460ff16610778565b610bd2610bcd366004614872565b612b99565b6040516103d99190614e24565b7f00000000000000000000000000000000000000000000000000000000000000006104fa565b6103cf612fa9565b6104a8610c1b366004614872565b612fb8565b601b546103cf565b610da16040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101919091525060408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff1661010082015290565b6040516103d99190614f5b565b610415610dbc366004614c15565b612fc3565b6103cf610dcf366004614aa1565b73ffffffffffffffffffffffffffffffffffffffff1660009081526021602052604090205490565b610e51610e05366004614aa1565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff9091166020830152016103d9565b610415610e78366004614aa1565b613122565b610ed8610e8b366004614872565b60408051808201909152600080825260208201525060009081526023602090815260409182902082518084019093525463ffffffff81168352640100000000900462ffffff169082015290565b60408051825163ffffffff16815260209283015162ffffff1692810192909252016103d9565b60406103cf565b610f3e610f13366004614aa1565b73ffffffffffffffffffffffffffffffffffffffff166000908152601c602052604090205460ff1690565b6040516103d9919061502b565b6103cf613136565b60606000610f61600261313e565b9050808410610f9c576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610fa8848661506e565b905081811180610fb6575083155b610fc05780610fc2565b815b90506000610fd08683615081565b67ffffffffffffffff811115610fe857610fe8615094565b604051908082528060200260200182016040528015611011578160200160208202803683370190505b50905060005b81518110156110645761103561102d888361506e565b600290613148565b828281518110611047576110476150c3565b60209081029190910101528061105c816150f2565b915050611017565b50925050505b92915050565b611078613154565b6000838152601f602052604090206110918284836151cc565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae776983836040516110c49291906152e7565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c081019190915260408051610140810182526016547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1681526000602082018190529181018290526014546bffffffffffffffffffffffff166060808301919091529182916080810161120a600261313e565b81526017547401000000000000000000000000000000000000000080820463ffffffff908116602080860191909152780100000000000000000000000000000000000000000000000080850483166040808801919091526013546060808901919091526014546c01000000000000000000000000810486166080808b0191909152760100000000000000000000000000000000000000000000820460ff16151560a09a8b015283516101e0810185526000808252968101879052601654898104891695820195909552700100000000000000000000000000000000830462ffffff169381019390935273010000000000000000000000000000000000000090910461ffff169082015296870192909252808204831660c08701527c0100000000000000000000000000000000000000000000000000000000909404821660e086015260185492830482166101008601529290910416610120830152601954610140830152601a5461016083015273ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a081016113a560096131a7565b815260175473ffffffffffffffffffffffffffffffffffffffff16602091820152601454600d80546040805182860281018601909152818152949850899489949293600e93750100000000000000000000000000000000000000000090910460ff1692859183018282801561145057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611425575b50505050509250818054806020026020016040519081016040528092919081815260200182805480156114b957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161148e575b50505050509150945094509450945094509091929394565b6114d96131b4565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601c6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561153957611539614abe565b02179055505050565b6000818152601f6020526040902080546060919061155f9061512a565b80601f016020809104026020016040519081016040528092919081815260200182805461158b9061512a565b80156115d85780601f106115ad576101008083540402835291602001916115d8565b820191906000526020600020905b8154815290600101906020018083116115bb57829003601f168201915b50505050509050919050565b6000818152601d6020526040902080546060919061155f9061512a565b6116096131b4565b600e548114611644576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e54811015611816576000600e8281548110611666576116666150c3565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168084526011909252604083205491935016908585858181106116b0576116b06150c3565b90506020020160208101906116c59190614aa1565b905073ffffffffffffffffffffffffffffffffffffffff81161580611758575073ffffffffffffffffffffffffffffffffffffffff82161580159061173657508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015611758575073ffffffffffffffffffffffffffffffffffffffff81811614155b1561178f576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146118005773ffffffffffffffffffffffffffffffffffffffff838116600090815260116020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061180e906150f2565b915050611647565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e838360405161184b93929190615334565b60405180910390a15050565b61185f6131b4565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152829182918291829190829061198357606082015160145460009161196f916bffffffffffffffffffffffff166153e6565b600e5490915061197f908261543a565b9150505b81516020830151604084015161199a908490615465565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b16600090815260116020526040902054929b919a9499509750921694509092505050565b6119e4613235565b600060255460ff1660018111156119fd576119fd614abe565b03611a34576040517fe0262d7400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e546bffffffffffffffffffffffff909116906000611a57600f61313e565b90506000611a65828461506e565b905060008167ffffffffffffffff811115611a8257611a82615094565b604051908082528060200260200182016040528015611aab578160200160208202803683370190505b50905060008267ffffffffffffffff811115611ac957611ac9615094565b604051908082528060200260200182016040528015611af2578160200160208202803683370190505b50905060005b85811015611c24576000600e8281548110611b1557611b156150c3565b600091825260208220015473ffffffffffffffffffffffffffffffffffffffff169150611b43828a8a613286565b9050806bffffffffffffffffffffffff16858481518110611b6657611b666150c3565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff808416600090815260119092526040909120548551911690859085908110611bb457611bb46150c3565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920181019190915292166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611c1c816150f2565b915050611af8565b5060005b84811015611da1576000611c3d600f83613148565b73ffffffffffffffffffffffffffffffffffffffff8082166000818152600b602090815260408083208151608081018352905460ff80821615158352610100820416828501526bffffffffffffffffffffffff6201000082048116838501526e0100000000000000000000000000009091041660608201529383526011909152902054929350911684611cd08a8661506e565b81518110611ce057611ce06150c3565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260408101516bffffffffffffffffffffffff1685611d238a8661506e565b81518110611d3357611d336150c3565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff9092166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611d99816150f2565b915050611c28565b5073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166000908152602160205260408120819055611df4600f61313e565b90505b8015611e3157611e1e611e16611e0e600184615081565b600f90613148565b600f9061348e565b5080611e298161548a565b915050611df7565b507f5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad8183604051611e639291906154bf565b60405180910390a1505050505050565b600061106a826134b0565b60408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff16610100820152600090818080611fb68461355b565b925092509250611fcc89858a8a8787878d61374d565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152602080526040902080546060919061155f9061512a565b612014613154565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260208052604090206120438284836151cc565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d283836040516110c49291906152e7565b73ffffffffffffffffffffffffffffffffffffffff81166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015612108573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061212c91906154ed565b61106a9190615081565b60015473ffffffffffffffffffffffffffffffffffffffff1633146121bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6122406131b4565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff167601000000000000000000000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020016118b3565b600e5460609060008167ffffffffffffffff8111156122d2576122d2615094565b60405190808252806020026020018201604052801561231757816020015b60408051808201909152600080825260208201528152602001906001900390816122f05790505b50905060005b828110156123bc576000600e828154811061233a5761233a6150c3565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff90811680845260118352604093849020548451808601909552818552909116918301829052855190935090919085908590811061239c5761239c6150c3565b6020026020010181905250505080806123b4906150f2565b91505061231d565b5092915050565b73ffffffffffffffffffffffffffffffffffffffff8116612410576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160255460ff16600181111561242957612429614abe565b03612460576040517f4a3578fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601160205260409020541633146124c0576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e546000916124e39185916bffffffffffffffffffffffff1690613286565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600b6020908152604080832080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff1690557f0000000000000000000000000000000000000000000000000000000000000000909316825260219052205490915061257a906bffffffffffffffffffffffff831690615081565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081166000818152602160205260408082209490945592517fa9059cbb00000000000000000000000000000000000000000000000000000000815291851660048301526bffffffffffffffffffffffff841660248301529063a9059cbb906044016020604051808303816000875af115801561262f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126539190615506565b90508061268c576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405133815273ffffffffffffffffffffffffffffffffffffffff808516916bffffffffffffffffffffffff8516918716907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a450505050565b6126f96131b4565b602580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260126020526040902054163314612783576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526011602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556012909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490931660e08201526002909101549091169181019190915261293b8361292b816134b0565b8360400151846101000151611e7e565b9392505050565b606060248054806020026020016040519081016040528092919081815260200182805480156129a757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161297c575b5050505050905090565b604080516102008101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a082018390526101c08201526101e0810191909152604080516102008101825260165463ffffffff74010000000000000000000000000000000000000000808304821684527801000000000000000000000000000000000000000000000000808404831660208601526017547c0100000000000000000000000000000000000000000000000000000000810484169686019690965273ffffffffffffffffffffffffffffffffffffffff938416606086015260145460ff828204161515608087015262ffffff70010000000000000000000000000000000082041660a0870152601854928304841660c087015290820490921660e085015293821661010084015261ffff7301000000000000000000000000000000000000009091041661012083015291909116610140820152601954610160820152601a54610180820152601b546101a08201526101c08101612b7360096131a7565b815260155473ffffffffffffffffffffffffffffffffffffffff16602090910152919050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a010000000000000000000090910481166080830181905260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c0100000000000000000000000000000000000000000000000000000000900490941660e08301526002909201549091169281019290925290919015612d6e57816080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612d45573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d699190615528565b612d71565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff168152602001600760008781526020019081526020016000208054612dc99061512a565b80601f0160208091040260200160405190810160405280929190818152602001828054612df59061512a565b8015612e425780601f10612e1757610100808354040283529160200191612e42565b820191906000526020600020905b815481529060010190602001808311612e2557829003601f168201915b505050505081526020018360c001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836060015163ffffffff1667ffffffffffffffff1681526020018360e0015163ffffffff1681526020018360a001516bffffffffffffffffffffffff168152602001836000015115158152602001601e60008781526020019081526020016000208054612f1f9061512a565b80601f0160208091040260200160405190810160405280929190818152602001828054612f4b9061512a565b8015612f985780601f10612f6d57610100808354040283529160200191612f98565b820191906000526020600020905b815481529060010190602001808311612f7b57829003601f168201915b505050505081525092505050919050565b6000612fb3613abd565b905090565b600061106a8261281b565b73ffffffffffffffffffffffffffffffffffffffff828116600090815260116020526040902054163314613023576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603613072576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526012602052604090205481169082161461311e5773ffffffffffffffffffffffffffffffffffffffff82811660008181526012602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45b5050565b61312a6131b4565b61313381613b87565b50565b6000612fb360025b600061106a825490565b600061293b8383613c7c565b60175473ffffffffffffffffffffffffffffffffffffffff1633146131a5576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6060600061293b83613ca6565b60005473ffffffffffffffffffffffffffffffffffffffff1633146131a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016121b3565b60185473ffffffffffffffffffffffffffffffffffffffff1633146131a5576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e010000000000000000000000000000900490911660608201529061348257600081606001518561331e91906153e6565b9050600061332c858361543a565b905080836040018181516133409190615465565b6bffffffffffffffffffffffff1690525061335b8582615545565b8360600181815161336c9190615465565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b600061293b8373ffffffffffffffffffffffffffffffffffffffff8416613d01565b6000818160045b600f81101561353d577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106134f5576134f56150c3565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461352b57506000949350505050565b80613535816150f2565b9150506134b7565b5081600f1a600181111561355357613553614abe565b949350505050565b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156135e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061360c919061558c565b509450909250505060008113158061362357508142105b806136445750828015613644575061363b8242615081565b8463ffffffff16105b15613653576019549650613657565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156136c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136e6919061558c565b50945090925050506000811315806136fd57508142105b8061371e575082801561371e57506137158242615081565b8463ffffffff16105b1561372d57601a549550613731565b8095505b868661373c8a613df4565b965096509650505050509193909250565b600080808089600181111561376457613764614abe565b03613773575062017f986137c8565b600189600181111561378757613787614abe565b0361379657506201de846137c8565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a6080015160016137db91906155dc565b6137e99060ff1660406155f5565b601854613817906103a49074010000000000000000000000000000000000000000900463ffffffff1661506e565b613821919061506e565b601554604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa158015613892573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138b6919061560c565b909250905081836138c883601861506e565b6138d291906155f5565b60808f01516138e29060016155dc565b6138f19060ff166115e06155f5565b6138fb919061506e565b613905919061506e565b61390f908561506e565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa158015613982573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139a691906154ed565b8d6060015161ffff166139b991906155f5565b945050505060006139ca8b86613ede565b60008d815260046020526040902054909150610100900460ff1615613a2f5760008c81526023602090815260409182902082518084018452905463ffffffff811680835264010000000090910462ffffff908116928401928352928501525116908201525b6000613a998c6040518061012001604052808d63ffffffff1681526020018681526020018781526020018c81526020018b81526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff1681526020018581526020016000151581525061405a565b60208101518151919250613aac91615465565b9d9c50505050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015613b59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b7d91906154ed565b612fb39190615630565b3373ffffffffffffffffffffffffffffffffffffffff821603613c06576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016121b3565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613c9357613c936150c3565b9060005260206000200154905092915050565b6060816000018054806020026020016040519081016040528092919081815260200182805480156115d857602002820191906000526020600020905b815481526020019060010190808311613ce25750505050509050919050565b60008181526001830160205260408120548015613dea576000613d25600183615081565b8554909150600090613d3990600190615081565b9050818114613d9e576000866000018281548110613d5957613d596150c3565b9060005260206000200154905080876000018481548110613d7c57613d7c6150c3565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613daf57613daf615650565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061106a565b600091505061106a565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613e64573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e88919061558c565b50935050925050600082131580613e9e57508042105b80613ece57506000846040015162ffffff16118015613ece5750613ec28142615081565b846040015162ffffff16105b156123bc575050601b5492915050565b60408051608081018252600080825260208083018281528385018381526060850184905273ffffffffffffffffffffffffffffffffffffffff878116855260229093528584208054640100000000810462ffffff1690925263ffffffff82169092527b01000000000000000000000000000000000000000000000000000000810460ff16855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495919484936701000000000000009092049091169163feaf968c9160048083019260a09291908290030181865afa158015613fcb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613fef919061558c565b5093505092505060008213158061400557508042105b8061403557506000866040015162ffffff1611801561403557506140298142615081565b866040015162ffffff16105b156140495760018301546060850152614051565b606084018290525b50505092915050565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081019190915260008260e001516000015160ff1690506000846060015161ffff1684606001516140c591906155f5565b905083610100015180156140d85750803a105b156140e057503a5b6000601283116140f1576001614107565b6140fc601284615081565b61410790600a61579f565b905060006012841061411a576001614130565b614125846012615081565b61413090600a61579f565b905060008660a00151876040015188602001518960000151614152919061506e565b61415c90876155f5565b614166919061506e565b61417091906155f5565b90506141cc828860e001516060015161418991906155f5565b6001848a60e001516060015161419f91906155f5565b6141a99190615081565b6141b386856155f5565b6141bd919061506e565b6141c791906157ab565b6143a1565b6bffffffffffffffffffffffff16865260808701516141ef906141c790836157ab565b6bffffffffffffffffffffffff1660408088019190915260e088015101516000906142289062ffffff16683635c9adc5dea000006155f5565b9050600081633b9aca008a60a001518b60e001516020015163ffffffff168c604001518d600001518b61425b91906155f5565b614265919061506e565b61426f91906155f5565b61427991906155f5565b61428391906157ab565b61428d919061506e565b90506142d0848a60e00151606001516142a691906155f5565b6001868c60e00151606001516142bc91906155f5565b6142c69190615081565b6141b388856155f5565b6bffffffffffffffffffffffff16602089015260808901516142f6906141c790836157ab565b6bffffffffffffffffffffffff16606089015260c089015173ffffffffffffffffffffffffffffffffffffffff166080808a0191909152890151614339906143a1565b6bffffffffffffffffffffffff1660a0808a019190915289015161435c906143a1565b6bffffffffffffffffffffffff1660c089015260e089015160600151614381906143a1565b6bffffffffffffffffffffffff1660e08901525050505050505092915050565b60006bffffffffffffffffffffffff82111561443f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016121b3565b5090565b6000806040838503121561445657600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b8381101561449557815187529582019590820190600101614479565b509495945050505050565b60208152600061293b6020830184614465565b60008083601f8401126144c557600080fd5b50813567ffffffffffffffff8111156144dd57600080fd5b6020830191508360208285010111156144f557600080fd5b9250929050565b60008060006040848603121561451157600080fd5b83359250602084013567ffffffffffffffff81111561452f57600080fd5b61453b868287016144b3565b9497909650939450505050565b600081518084526020808501945080840160005b8381101561449557815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161455c565b805163ffffffff16825260006101e060208301516145b4602086018263ffffffff169052565b5060408301516145cc604086018263ffffffff169052565b5060608301516145e3606086018262ffffff169052565b5060808301516145f9608086018261ffff169052565b5060a083015161461960a08601826bffffffffffffffffffffffff169052565b5060c083015161463160c086018263ffffffff169052565b5060e083015161464960e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a0808401518186018390526146be83870182614548565b925050506101c0808401516146ea8287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161472260208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161474c60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a088015161476e60a085018263ffffffff169052565b5060c088015161478660c085018263ffffffff169052565b5060e088015160e0840152610100808901516147a98286018263ffffffff169052565b50506101208881015115159084015261014083018190526147cc8184018861458e565b90508281036101608401526147e18187614548565b90508281036101808401526147f68186614548565b9150506148096101a083018460ff169052565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461313357600080fd5b6000806040838503121561484857600080fd5b823561485381614813565b915060208301356004811061486757600080fd5b809150509250929050565b60006020828403121561488457600080fd5b5035919050565b6000815180845260005b818110156148b157602081850181015186830182015201614895565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061293b602083018461488b565b815173ffffffffffffffffffffffffffffffffffffffff16815261016081016020830151614938602084018263ffffffff169052565b506040830151614950604084018263ffffffff169052565b506060830151614968606084018263ffffffff169052565b506080830151614990608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a08301516149a860a084018263ffffffff169052565b5060c08301516149c060c084018263ffffffff169052565b5060e08301516149d860e084018263ffffffff169052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff81168483015250506101208381015163ffffffff81168483015250506101408381015163ffffffff8116848301525b505092915050565b60008060208385031215614a3f57600080fd5b823567ffffffffffffffff80821115614a5757600080fd5b818501915085601f830112614a6b57600080fd5b813581811115614a7a57600080fd5b8660208260051b8501011115614a8f57600080fd5b60209290920196919550909350505050565b600060208284031215614ab357600080fd5b813561293b81614813565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6002811061313357613133614abe565b60208101614b0a83614aed565b91905290565b60008060008060808587031215614b2657600080fd5b84359350602085013560028110614b3c57600080fd5b9250604085013563ffffffff81168114614b5557600080fd5b91506060850135614b6581614813565b939692955090935050565b600080600060408486031215614b8557600080fd5b8335614b9081614813565b9250602084013567ffffffffffffffff81111561452f57600080fd5b602080825282518282018190526000919060409081850190868401855b82811015614c08578151805173ffffffffffffffffffffffffffffffffffffffff90811686529087015116868501529284019290850190600101614bc9565b5091979650505050505050565b60008060408385031215614c2857600080fd5b8235614c3381614813565b9150602083013561486781614813565b6020808252825182820181905260009190848201906040850190845b81811015614c9157835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101614c5f565b50909695505050505050565b60208152614cb460208201835163ffffffff169052565b60006020830151614ccd604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015173ffffffffffffffffffffffffffffffffffffffff8116608084015250608083015180151560a08401525060a083015162ffffff811660c08401525060c083015163ffffffff811660e08401525060e0830151610100614d4c8185018363ffffffff169052565b8401519050610120614d758482018373ffffffffffffffffffffffffffffffffffffffff169052565b8401519050610140614d8c8482018361ffff169052565b8401519050610160614db58482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151610180848101919091528401516101a0808501919091528401516101c0808501919091528401516102006101e080860182905291925090614dfd610220860184614548565b9086015173ffffffffffffffffffffffffffffffffffffffff8116838701529092506146ea565b60208152614e4b60208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614e64604084018263ffffffff169052565b506040830151610140806060850152614e8161016085018361488b565b91506060850151614ea260808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614f0e818701836bffffffffffffffffffffffff169052565b8601519050610120614f238682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001838701529050614809838261488b565b6000610120820190506bffffffffffffffffffffffff835116825263ffffffff60208401511660208301526040830151614f9c604084018262ffffff169052565b506060830151614fb2606084018261ffff169052565b506080830151614fc7608084018260ff169052565b5060a0830151614fdb60a084018215159052565b5060c0830151614fef60c084018215159052565b5060e083015161500360e084018215159052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff811684830152614a24565b6020810160048310614b0a57614b0a614abe565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561106a5761106a61503f565b8181038181111561106a5761106a61503f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036151235761512361503f565b5060010190565b600181811c9082168061513e57607f821691505b602082108103615177577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156151c757600081815260208120601f850160051c810160208610156151a45750805b601f850160051c820191505b818110156151c3578281556001016151b0565b5050505b505050565b67ffffffffffffffff8311156151e4576151e4615094565b6151f8836151f2835461512a565b8361517d565b6000601f84116001811461524a57600085156152145750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556152e0565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156152995786850135825560209485019460019092019101615279565b50868210156152d4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000604082016040835280865480835260608501915087600052602092508260002060005b8281101561538b57815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101615359565b505050838103828501528481528590820160005b868110156153da5782356153b281614813565b73ffffffffffffffffffffffffffffffffffffffff168252918301919083019060010161539f565b50979650505050505050565b6bffffffffffffffffffffffff8281168282160390808211156123bc576123bc61503f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff808416806154595761545961540b565b92169190910492915050565b6bffffffffffffffffffffffff8181168382160190808211156123bc576123bc61503f565b6000816154995761549961503f565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b6040815260006154d26040830185614548565b82810360208401526154e48185614465565b95945050505050565b6000602082840312156154ff57600080fd5b5051919050565b60006020828403121561551857600080fd5b8151801515811461293b57600080fd5b60006020828403121561553a57600080fd5b815161293b81614813565b6bffffffffffffffffffffffff818116838216028082169190828114614a2457614a2461503f565b805169ffffffffffffffffffff8116811461558757600080fd5b919050565b600080600080600060a086880312156155a457600080fd5b6155ad8661556d565b94506020860151935060408601519250606086015191506155d06080870161556d565b90509295509295909350565b60ff818116838216019081111561106a5761106a61503f565b808202811582820484141761106a5761106a61503f565b6000806040838503121561561f57600080fd5b505080516020909101519092909150565b81810360008312801583831316838312821617156123bc576123bc61503f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600181815b808511156156d857817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156156be576156be61503f565b808516156156cb57918102915b93841c9390800290615684565b509250929050565b6000826156ef5750600161106a565b816156fc5750600061106a565b8160018114615712576002811461571c57615738565b600191505061106a565b60ff84111561572d5761572d61503f565b50506001821b61106a565b5060208310610133831016604e8410600b841016171561575b575081810a61106a565b615765838361567f565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156157975761579761503f565b029392505050565b600061293b83836156e0565b6000826157ba576157ba61540b565b50049056fea164736f6c6343000813000a",
}

var AutomationRegistryLogicCABI = AutomationRegistryLogicCMetaData.ABI

var AutomationRegistryLogicCBin = AutomationRegistryLogicCMetaData.Bin

func DeployAutomationRegistryLogicC(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkUSDFeed common.Address, nativeUSDFeed common.Address, fastGasFeed common.Address, automationForwarderLogic common.Address, allowedReadOnlyAddress common.Address, payoutMode uint8, wrappedNativeTokenAddress common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicC, error) {
	parsed, err := AutomationRegistryLogicCMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicCBin), backend, link, linkUSDFeed, nativeUSDFeed, fastGasFeed, automationForwarderLogic, allowedReadOnlyAddress, payoutMode, wrappedNativeTokenAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicC{address: address, abi: *parsed, AutomationRegistryLogicCCaller: AutomationRegistryLogicCCaller{contract: contract}, AutomationRegistryLogicCTransactor: AutomationRegistryLogicCTransactor{contract: contract}, AutomationRegistryLogicCFilterer: AutomationRegistryLogicCFilterer{contract: contract}}, nil
}

type AutomationRegistryLogicC struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicCCaller
	AutomationRegistryLogicCTransactor
	AutomationRegistryLogicCFilterer
}

type AutomationRegistryLogicCCaller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicCTransactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicCFilterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicCSession struct {
	Contract     *AutomationRegistryLogicC
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicCCallerSession struct {
	Contract *AutomationRegistryLogicCCaller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicCTransactorSession struct {
	Contract     *AutomationRegistryLogicCTransactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicCRaw struct {
	Contract *AutomationRegistryLogicC
}

type AutomationRegistryLogicCCallerRaw struct {
	Contract *AutomationRegistryLogicCCaller
}

type AutomationRegistryLogicCTransactorRaw struct {
	Contract *AutomationRegistryLogicCTransactor
}

func NewAutomationRegistryLogicC(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicC, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicCABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicC{address: address, abi: abi, AutomationRegistryLogicCCaller: AutomationRegistryLogicCCaller{contract: contract}, AutomationRegistryLogicCTransactor: AutomationRegistryLogicCTransactor{contract: contract}, AutomationRegistryLogicCFilterer: AutomationRegistryLogicCFilterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicCCaller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicCCaller, error) {
	contract, err := bindAutomationRegistryLogicC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCCaller{contract: contract}, nil
}

func NewAutomationRegistryLogicCTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicCTransactor, error) {
	contract, err := bindAutomationRegistryLogicC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCTransactor{contract: contract}, nil
}

func NewAutomationRegistryLogicCFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicCFilterer, error) {
	contract, err := bindAutomationRegistryLogicC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCFilterer{contract: contract}, nil
}

func bindAutomationRegistryLogicC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicC.Contract.AutomationRegistryLogicCCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AutomationRegistryLogicCTransactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AutomationRegistryLogicCTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicC.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicC.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicC.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicC.CallOpts, admin)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicC.CallOpts, admin)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetAvailableERC20ForPayment(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getAvailableERC20ForPayment", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetAvailableERC20ForPayment(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetAvailableERC20ForPayment(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetAvailableERC20ForPayment(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetBalance(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetBalance(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingConfig(opts *bind.CallOpts, billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingConfig", billingToken)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingConfig(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingConfig(billingToken common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingConfig(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingOverrides(opts *bind.CallOpts, upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingOverrides", upkeepID)

	if err != nil {
		return *new(AutomationRegistryBase23BillingOverrides), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingOverrides)).(*AutomationRegistryBase23BillingOverrides)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingOverrides(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingOverrides(upkeepID *big.Int) (AutomationRegistryBase23BillingOverrides, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingOverrides(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingToken", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingToken(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingToken(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingTokenConfig", token)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicC.CallOpts, token)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicC.CallOpts, token)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingTokens(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingTokens(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetCancellationDelay(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetCancellationDelay(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetChainModule(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetChainModule(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetConfig(opts *bind.CallOpts) (AutomationRegistryBase23OnchainConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(AutomationRegistryBase23OnchainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23OnchainConfig)).(*AutomationRegistryBase23OnchainConfig)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetConfig(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetConfig() (AutomationRegistryBase23OnchainConfig, error) {
	return _AutomationRegistryLogicC.Contract.GetConfig(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getFallbackNativePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetForwarder(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetForwarder(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getHotVars")

	if err != nil {
		return *new(AutomationRegistryBase23HotVars), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23HotVars)).(*AutomationRegistryBase23HotVars)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicC.Contract.GetHotVars(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicC.Contract.GetHotVars(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetLinkAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetLinkAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getLinkUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetLogGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetLogGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetMaxPaymentForGas(opts *bind.CallOpts, id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getMaxPaymentForGas", id, triggerType, gasLimit, billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicC.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetMaxPaymentForGas(id *big.Int, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicC.CallOpts, id, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMinBalance(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMinBalance(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetNativeUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getNativeUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetNumUpkeeps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getNumUpkeeps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetNumUpkeeps(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetNumUpkeeps(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetPayoutMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getPayoutMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetPayoutMode(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetPayoutMode(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC.CallOpts, peer)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC.CallOpts, peer)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicC.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicC.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetReserveAmount(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getReserveAmount", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetReserveAmount(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetReserveAmount(&_AutomationRegistryLogicC.CallOpts, billingToken)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicC.Contract.GetSignerInfo(&_AutomationRegistryLogicC.CallOpts, query)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicC.Contract.GetSignerInfo(&_AutomationRegistryLogicC.CallOpts, query)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getState")

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

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicC.Contract.GetState(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicC.Contract.GetState(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetStorage(opts *bind.CallOpts) (AutomationRegistryBase23Storage, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getStorage")

	if err != nil {
		return *new(AutomationRegistryBase23Storage), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23Storage)).(*AutomationRegistryBase23Storage)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicC.Contract.GetStorage(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicC.Contract.GetStorage(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitterInfo(&_AutomationRegistryLogicC.CallOpts, query)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicC.Contract.GetTransmitterInfo(&_AutomationRegistryLogicC.CallOpts, query)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetTransmittersWithPayees(opts *bind.CallOpts) ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getTransmittersWithPayees")

	if err != nil {
		return *new([]AutomationRegistryBase23TransmitterPayeeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutomationRegistryBase23TransmitterPayeeInfo)).(*[]AutomationRegistryBase23TransmitterPayeeInfo)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmittersWithPayees(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetTransmittersWithPayees() ([]AutomationRegistryBase23TransmitterPayeeInfo, error) {
	return _AutomationRegistryLogicC.Contract.GetTransmittersWithPayees(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetTriggerType(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicC.Contract.GetTriggerType(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeep(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeep(&_AutomationRegistryLogicC.CallOpts, id)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicC.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicC.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetWrappedNativeTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getWrappedNativeTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetWrappedNativeTokenAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetWrappedNativeTokenAddress() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.GetWrappedNativeTokenAddress(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicC.Contract.HasDedupKey(&_AutomationRegistryLogicC.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicC.Contract.HasDedupKey(&_AutomationRegistryLogicC.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicC.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.Owner(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicC.Contract.Owner(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "supportsBillingToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicC.Contract.SupportsBillingToken(&_AutomationRegistryLogicC.CallOpts, token)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicC.Contract.SupportsBillingToken(&_AutomationRegistryLogicC.CallOpts, token)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicC.Contract.UpkeepVersion(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicC.Contract.UpkeepVersion(&_AutomationRegistryLogicC.CallOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AcceptOwnership(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AcceptOwnership(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AcceptPayeeship(&_AutomationRegistryLogicC.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.AcceptPayeeship(&_AutomationRegistryLogicC.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "disableOffchainPayments")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.DisableOffchainPayments(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.DisableOffchainPayments(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "pause")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.Pause(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.Pause(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicC.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicC.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "setPayees", payees)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetPayees(&_AutomationRegistryLogicC.TransactOpts, payees)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetPayees(&_AutomationRegistryLogicC.TransactOpts, payees)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicC.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicC.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "settleNOPsOffchain")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.TransferOwnership(&_AutomationRegistryLogicC.TransactOpts, to)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.TransferOwnership(&_AutomationRegistryLogicC.TransactOpts, to)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.TransferPayeeship(&_AutomationRegistryLogicC.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.TransferPayeeship(&_AutomationRegistryLogicC.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "unpause")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.Unpause(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.Unpause(&_AutomationRegistryLogicC.TransactOpts)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.WithdrawPayment(&_AutomationRegistryLogicC.TransactOpts, from, to)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicC.Contract.WithdrawPayment(&_AutomationRegistryLogicC.TransactOpts, from, to)
}

type AutomationRegistryLogicCAdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicCAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCAdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCAdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicCAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCAdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicC.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCAdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicCAdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicCAdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCBillingConfigOverriddenIterator struct {
	Event *AutomationRegistryLogicCBillingConfigOverridden

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCBillingConfigOverriddenIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCBillingConfigOverridden)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCBillingConfigOverridden)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCBillingConfigOverriddenIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCBillingConfigOverriddenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCBillingConfigOverridden struct {
	Id        *big.Int
	Overrides AutomationRegistryBase23BillingOverrides
	Raw       types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCBillingConfigOverriddenIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCBillingConfigOverriddenIterator{contract: _AutomationRegistryLogicC.contract, event: "BillingConfigOverridden", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigOverridden, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "BillingConfigOverridden", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCBillingConfigOverridden)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicCBillingConfigOverridden, error) {
	event := new(AutomationRegistryLogicCBillingConfigOverridden)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigOverridden", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCBillingConfigOverrideRemovedIterator struct {
	Event *AutomationRegistryLogicCBillingConfigOverrideRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCBillingConfigOverrideRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCBillingConfigOverrideRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCBillingConfigOverrideRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCBillingConfigOverrideRemovedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCBillingConfigOverrideRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCBillingConfigOverrideRemoved struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCBillingConfigOverrideRemovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCBillingConfigOverrideRemovedIterator{contract: _AutomationRegistryLogicC.contract, event: "BillingConfigOverrideRemoved", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "BillingConfigOverrideRemoved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCBillingConfigOverrideRemoved)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicCBillingConfigOverrideRemoved, error) {
	event := new(AutomationRegistryLogicCBillingConfigOverrideRemoved)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigOverrideRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCBillingConfigSetIterator struct {
	Event *AutomationRegistryLogicCBillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCBillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCBillingConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCBillingConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCBillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCBillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCBillingConfigSet struct {
	Token  common.Address
	Config AutomationRegistryBase23BillingConfig
	Raw    types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicCBillingConfigSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCBillingConfigSetIterator{contract: _AutomationRegistryLogicC.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigSet, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "BillingConfigSet", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCBillingConfigSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicCBillingConfigSet, error) {
	event := new(AutomationRegistryLogicCBillingConfigSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCCancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicCCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCCancelledUpkeepReportIterator{contract: _AutomationRegistryLogicC.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCCancelledUpkeepReport)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicCCancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicCCancelledUpkeepReport)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicCChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicCChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicC.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicCChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicCChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCDedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicCDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCDedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCDedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicCDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCDedupKeyAddedIterator{contract: _AutomationRegistryLogicC.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCDedupKeyAdded)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicCDedupKeyAdded, error) {
	event := new(AutomationRegistryLogicCDedupKeyAdded)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCFeesWithdrawnIterator struct {
	Event *AutomationRegistryLogicCFeesWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCFeesWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCFeesWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCFeesWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCFeesWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCFeesWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCFeesWithdrawn struct {
	AssetAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicCFeesWithdrawnIterator, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCFeesWithdrawnIterator{contract: _AutomationRegistryLogicC.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error) {

	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "FeesWithdrawn", assetAddressRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCFeesWithdrawn)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicCFeesWithdrawn, error) {
	event := new(AutomationRegistryLogicCFeesWithdrawn)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FeesWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCFundsAddedIterator struct {
	Event *AutomationRegistryLogicCFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCFundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicCFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCFundsAddedIterator{contract: _AutomationRegistryLogicC.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCFundsAdded)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicCFundsAdded, error) {
	event := new(AutomationRegistryLogicCFundsAdded)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicCFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCFundsWithdrawnIterator{contract: _AutomationRegistryLogicC.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCFundsWithdrawn)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicCFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicCFundsWithdrawn)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicCInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicC.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCInsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicCInsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicCInsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCNOPsSettledOffchainIterator struct {
	Event *AutomationRegistryLogicCNOPsSettledOffchain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCNOPsSettledOffchainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCNOPsSettledOffchain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCNOPsSettledOffchain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCNOPsSettledOffchainIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCNOPsSettledOffchainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCNOPsSettledOffchain struct {
	Payees   []common.Address
	Payments []*big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicCNOPsSettledOffchainIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCNOPsSettledOffchainIterator{contract: _AutomationRegistryLogicC.contract, event: "NOPsSettledOffchain", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCNOPsSettledOffchain) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "NOPsSettledOffchain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCNOPsSettledOffchain)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicCNOPsSettledOffchain, error) {
	event := new(AutomationRegistryLogicCNOPsSettledOffchain)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "NOPsSettledOffchain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCOwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicCOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicCOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCOwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicC.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCOwnershipTransferRequested)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicCOwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicCOwnershipTransferRequested)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCOwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicCOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicCOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCOwnershipTransferredIterator{contract: _AutomationRegistryLogicC.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCOwnershipTransferred)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicCOwnershipTransferred, error) {
	event := new(AutomationRegistryLogicCOwnershipTransferred)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCPausedIterator struct {
	Event *AutomationRegistryLogicCPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicCPausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCPausedIterator{contract: _AutomationRegistryLogicC.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCPaused)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParsePaused(log types.Log) (*AutomationRegistryLogicCPaused, error) {
	event := new(AutomationRegistryLogicCPaused)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCPayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicCPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicCPayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCPayeesUpdatedIterator{contract: _AutomationRegistryLogicC.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCPayeesUpdated)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicCPayeesUpdated, error) {
	event := new(AutomationRegistryLogicCPayeesUpdated)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCPayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicCPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicCPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCPayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicC.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCPayeeshipTransferRequested)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicCPayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicCPayeeshipTransferRequested)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCPayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicCPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicCPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCPayeeshipTransferredIterator{contract: _AutomationRegistryLogicC.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCPayeeshipTransferred)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicCPayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicCPayeeshipTransferred)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCPaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicCPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicCPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCPaymentWithdrawnIterator{contract: _AutomationRegistryLogicC.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCPaymentWithdrawn)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicCPaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicCPaymentWithdrawn)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicCReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicC.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCReorgedUpkeepReport)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicCReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicCReorgedUpkeepReport)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCStaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicCStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCStaleUpkeepReportIterator{contract: _AutomationRegistryLogicC.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCStaleUpkeepReport)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicCStaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicCStaleUpkeepReport)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUnpausedIterator struct {
	Event *AutomationRegistryLogicCUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicCUnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUnpausedIterator{contract: _AutomationRegistryLogicC.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUnpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUnpaused)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicCUnpaused, error) {
	event := new(AutomationRegistryLogicCUnpaused)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicCUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicCUpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicCUpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicCUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicCUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepAdminTransferred)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicCUpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicCUpkeepAdminTransferred)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicCUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicCUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepCanceledIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepCanceled)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicCUpkeepCanceled, error) {
	event := new(AutomationRegistryLogicCUpkeepCanceled)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepChargedIterator struct {
	Event *AutomationRegistryLogicCUpkeepCharged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepChargedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepCharged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepCharged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepChargedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepChargedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepCharged struct {
	Id      *big.Int
	Receipt AutomationRegistryBase23PaymentReceipt
	Raw     types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepChargedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepChargedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepCharged", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCharged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepCharged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepCharged)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicCUpkeepCharged, error) {
	event := new(AutomationRegistryLogicCUpkeepCharged)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCharged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicCUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepCheckDataSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicCUpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicCUpkeepCheckDataSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicCUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepGasLimitSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicCUpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicCUpkeepGasLimitSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicCUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepMigratedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepMigrated)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicCUpkeepMigrated, error) {
	event := new(AutomationRegistryLogicCUpkeepMigrated)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicCUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicCUpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepPausedIterator struct {
	Event *AutomationRegistryLogicCUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepPausedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepPaused)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicCUpkeepPaused, error) {
	event := new(AutomationRegistryLogicCUpkeepPaused)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicCUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicCUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepPerformedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepPerformed)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicCUpkeepPerformed, error) {
	event := new(AutomationRegistryLogicCUpkeepPerformed)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicCUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicCUpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicCUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepReceivedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepReceived)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicCUpkeepReceived, error) {
	event := new(AutomationRegistryLogicCUpkeepReceived)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicCUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepRegisteredIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepRegistered)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicCUpkeepRegistered, error) {
	event := new(AutomationRegistryLogicCUpkeepRegistered)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicCUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicCUpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicCUpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicCUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicCUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicCUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicCUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicCUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicCUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicCUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicCUpkeepUnpausedIterator{contract: _AutomationRegistryLogicC.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicC.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicCUpkeepUnpaused)
				if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCFilterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicCUpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicCUpkeepUnpaused)
	if err := _AutomationRegistryLogicC.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistryLogicC *AutomationRegistryLogicC) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicC.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicC.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicC.abi.Events["BillingConfigOverridden"].ID:
		return _AutomationRegistryLogicC.ParseBillingConfigOverridden(log)
	case _AutomationRegistryLogicC.abi.Events["BillingConfigOverrideRemoved"].ID:
		return _AutomationRegistryLogicC.ParseBillingConfigOverrideRemoved(log)
	case _AutomationRegistryLogicC.abi.Events["BillingConfigSet"].ID:
		return _AutomationRegistryLogicC.ParseBillingConfigSet(log)
	case _AutomationRegistryLogicC.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicC.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicC.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicC.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicC.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicC.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicC.abi.Events["FeesWithdrawn"].ID:
		return _AutomationRegistryLogicC.ParseFeesWithdrawn(log)
	case _AutomationRegistryLogicC.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicC.ParseFundsAdded(log)
	case _AutomationRegistryLogicC.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicC.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicC.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicC.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicC.abi.Events["NOPsSettledOffchain"].ID:
		return _AutomationRegistryLogicC.ParseNOPsSettledOffchain(log)
	case _AutomationRegistryLogicC.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicC.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicC.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicC.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicC.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicC.ParsePaused(log)
	case _AutomationRegistryLogicC.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicC.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicC.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicC.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicC.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicC.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicC.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicC.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicC.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicC.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicC.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicC.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicC.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicC.ParseUnpaused(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepCharged"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepCharged(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicC.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicC.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicCAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicCBillingConfigOverridden) Topic() common.Hash {
	return common.HexToHash("0xd8a6d79d170a55968079d3a89b960d86b4442aef6aac1d01e644c32b9e38b340")
}

func (AutomationRegistryLogicCBillingConfigOverrideRemoved) Topic() common.Hash {
	return common.HexToHash("0x97d0ef3f46a56168af653f547bdb6f77ec2b1d7d9bc6ba0193c2b340ec68064a")
}

func (AutomationRegistryLogicCBillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0xca93cbe727c73163ec538f71be6c0a64877d7f1f6dd35d5ca7cbaef3a3e34ba3")
}

func (AutomationRegistryLogicCCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicCChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicCDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicCFeesWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa8")
}

func (AutomationRegistryLogicCFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicCFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicCInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicCNOPsSettledOffchain) Topic() common.Hash {
	return common.HexToHash("0x5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad")
}

func (AutomationRegistryLogicCOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicCOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicCPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicCPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicCPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicCPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicCPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicCReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicCStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicCUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicCUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicCUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicCUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicCUpkeepCharged) Topic() common.Hash {
	return common.HexToHash("0x801ba6ed51146ffe3e99d1dbd9dd0f4de6292e78a9a34c39c0183de17b3f40fc")
}

func (AutomationRegistryLogicCUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicCUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicCUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicCUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicCUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicCUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicCUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicCUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicCUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicCUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicCUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicC) Address() common.Address {
	return _AutomationRegistryLogicC.address
}

type AutomationRegistryLogicCInterface interface {
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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicCAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicCAdminPrivilegeConfigSet, error)

	FilterBillingConfigOverridden(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCBillingConfigOverriddenIterator, error)

	WatchBillingConfigOverridden(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigOverridden, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverridden(log types.Log) (*AutomationRegistryLogicCBillingConfigOverridden, error)

	FilterBillingConfigOverrideRemoved(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCBillingConfigOverrideRemovedIterator, error)

	WatchBillingConfigOverrideRemoved(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigOverrideRemoved, id []*big.Int) (event.Subscription, error)

	ParseBillingConfigOverrideRemoved(log types.Log) (*AutomationRegistryLogicCBillingConfigOverrideRemoved, error)

	FilterBillingConfigSet(opts *bind.FilterOpts, token []common.Address) (*AutomationRegistryLogicCBillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCBillingConfigSet, token []common.Address) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*AutomationRegistryLogicCBillingConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicCCancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicCChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicCChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicCDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicCDedupKeyAdded, error)

	FilterFeesWithdrawn(opts *bind.FilterOpts, assetAddress []common.Address, recipient []common.Address) (*AutomationRegistryLogicCFeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFeesWithdrawn, assetAddress []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseFeesWithdrawn(log types.Log) (*AutomationRegistryLogicCFeesWithdrawn, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicCFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicCFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicCFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicCInsufficientFundsUpkeepReport, error)

	FilterNOPsSettledOffchain(opts *bind.FilterOpts) (*AutomationRegistryLogicCNOPsSettledOffchainIterator, error)

	WatchNOPsSettledOffchain(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCNOPsSettledOffchain) (event.Subscription, error)

	ParseNOPsSettledOffchain(log types.Log) (*AutomationRegistryLogicCNOPsSettledOffchain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicCOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicCOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicCOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicCOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicCPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicCPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicCPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicCPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicCPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicCPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicCPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicCPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicCPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicCPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicCReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicCStaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicCUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicCUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicCUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicCUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicCUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicCUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicCUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicCUpkeepCanceled, error)

	FilterUpkeepCharged(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepChargedIterator, error)

	WatchUpkeepCharged(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCharged, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCharged(log types.Log) (*AutomationRegistryLogicCUpkeepCharged, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicCUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicCUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicCUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicCUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicCUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicCUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicCUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicCUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicCUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicCUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicCUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicCUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
