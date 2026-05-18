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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nativeUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"automationForwarderLogic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedReadOnlyAddress\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"payoutMode\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"wrappedNativeTokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLinkLiquidity\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"overrides\",\"type\":\"tuple\"}],\"name\":\"BillingConfigOverridden\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"BillingConfigOverrideRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payments\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint96\",\"name\":\"gasChargeInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInBillingToken\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"gasReimbursementInJuels\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"premiumInJuels\",\"type\":\"uint96\"},{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"linkUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeUSD\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"billingUSD\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.PaymentReceipt\",\"name\":\"receipt\",\"type\":\"tuple\"}],\"name\":\"UpkeepCharged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableOffchainPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"getAdminPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedReadOnlyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAutomationForwarderLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getAvailableERC20ForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getBillingConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingOverrides\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingOverrides\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingOverridesEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingToken\",\"outputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBillingTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMilliCents\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingTokens\",\"outputs\":[{\"internalType\":\"contractIERC20Metadata[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCancellationDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainModule\",\"outputs\":[{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConditionalGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackNativePrice\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.OnchainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackNativePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"contractIAutomationForwarder\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHotVars\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuard\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.HotVars\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLogGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNativeUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumUpkeeps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPayoutMode\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerPerformByteGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerSignerGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReorgProtectionEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStorage\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"}],\"internalType\":\"structAutomationRegistryBase2_3.Storage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmittersWithPayees\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payeeAddress\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.TransmitterPayeeInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNativeTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setAdminPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settleNOPsOffchain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"supportsBillingToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b5060405162005c8738038062005c87833981016040819052620000359162000309565b87878787878787873380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c78162000241565b5050506001600160a01b0380891660805287811660a05286811660c05285811660e052848116610100528316610120526025805483919060ff191660018381811115620001185762000118620003b6565b0217905550806001600160a01b0316610140816001600160a01b03168152505060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000179573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200019f9190620003cc565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002099190620003cc565b60ff16146200022b576040516301f86e1760e41b815260040160405180910390fd5b50505050505050505050505050505050620003f8565b336001600160a01b038216036200029b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200030457600080fd5b919050565b600080600080600080600080610100898b0312156200032757600080fd5b6200033289620002ec565b97506200034260208a01620002ec565b96506200035260408a01620002ec565b95506200036260608a01620002ec565b94506200037260808a01620002ec565b93506200038260a08a01620002ec565b925060c0890151600281106200039757600080fd5b9150620003a760e08a01620002ec565b90509295985092959890939650565b634e487b7160e01b600052602160045260246000fd5b600060208284031215620003df57600080fd5b815160ff81168114620003f157600080fd5b9392505050565b60805160a05160c05160e051610100516101205161014051615803620004846000396000610b7001526000610ad2015260006107be0152600081816109a001526135b60152600081816109540152613e3201526000818161051d0152613690015260008181610c1801528181611df00152818161256b015281816125c80152613af601526158036000f3fe608060405234801561001057600080fd5b50600436106103d05760003560e01c80638ed02bab116101ff578063c3f909d41161011a578063eb5dcd6c116100ad578063f5a418461161007c578063f5a4184614610eb4578063f777ff0614610f35578063faa3e99614610f3c578063ffd242bd14610f8257600080fd5b8063eb5dcd6c14610de5578063ec4de4ba14610df8578063ed56b3e114610e2e578063f2fde38b14610ea157600080fd5b8063d09dc339116100e9578063d09dc33914610c3c578063d763264814610c44578063d85aa07c14610c57578063e80a3d4414610c5f57600080fd5b8063c3f909d414610bd6578063c5b964e014610beb578063c7c3a19a14610bf6578063ca30e60314610c1657600080fd5b8063aab9edd611610192578063b3596c2311610161578063b3596c23146107e2578063b6511a2a14610ba7578063b657bc9c14610bae578063ba87666814610bc157600080fd5b8063aab9edd614610b57578063abc76ae014610b66578063ac4dc59a14610b6e578063b121e14714610b9457600080fd5b8063a08714c0116101ce578063a08714c014610ad0578063a538b2eb14610af6578063a710b22114610b3c578063a87f45fe14610b4f57600080fd5b80638ed02bab14610a5c5780639089daa414610a7a57806393f6ebcf14610a8f5780639e0a99ed14610ac857600080fd5b806343cc055c116102ef5780636209e1e91161028257806379ba50971161025157806379ba5097146109ea57806379ea9943146109f25780638456cb5914610a365780638da5cb5b14610a3e57600080fd5b80636209e1e91461098b5780636709d0e51461099e578063671d36ed146109c45780636eec02a2146109d757600080fd5b80635425d8ac116102be5780635425d8ac146107bc57806357359584146107e2578063614486af146109525780636181d82d1461097857600080fd5b806343cc055c1461074957806344cb70b8146107705780634ca16c52146107935780635147cd591461079c57600080fd5b8063207b6516116103675780633b9cce59116103365780633b9cce59146106c05780633f4ba83a146106d3578063421d183b146106db57806343b46e5f1461074157600080fd5b8063207b651614610508578063226cf83c1461051b578063232c1cc5146105625780633408f73a1461056957600080fd5b8063187256e8116103a3578063187256e81461043b57806319d97a941461044e5780631e0104391461046e5780631efcf646146104d057600080fd5b8063050ee65d146103d557806306e3b632146103ed5780630b7d33e61461040d5780631865c57d14610422575b600080fd5b6201e26c5b6040519081526020015b60405180910390f35b6104006103fb36600461447a565b610f8a565b6040516103e491906144d7565b61042061041b366004614533565b6110a7565b005b61042a611108565b6040516103e495949392919061472b565b61042061044936600461486c565b611508565b61046161045c3660046148a9565b611579565b6040516103e49190614926565b6104b361047c3660046148a9565b60009081526004602052604090206001015470010000000000000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff90911681526020016103e4565b6104f86104de3660046148a9565b600090815260046020526040902054610100900460ff1690565b60405190151581526020016103e4565b6104616105163660046148a9565b61161b565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016103e4565b60186103da565b6106b36040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081019190915250604080516101608101825260165473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008084048216602086015278010000000000000000000000000000000000000000000000008085048316968601969096527c010000000000000000000000000000000000000000000000000000000093849004821660608601526017548084166080870152818104831660a0870152868104831660c087015293909304811660e085015260185491821661010085015291810482166101208401529290920490911661014082015290565b6040516103e49190614939565b6104206106ce366004614a63565b611638565b61042061188e565b6106ee6106e9366004614ad8565b6118f4565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a0016103e4565b610420611a13565b6014547801000000000000000000000000000000000000000000000000900460ff166104f8565b6104f861077e3660046148a9565b60009081526008602052604090205460ff1690565b62017f986103da565b6107af6107aa3660046148a9565b611eaa565b6040516103e49190614b34565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6108d46107f0366004614ad8565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091525073ffffffffffffffffffffffffffffffffffffffff908116600090815260226020908152604091829020825160c081018452815463ffffffff81168252640100000000810462ffffff16938201939093526701000000000000008304909416928401929092527b01000000000000000000000000000000000000000000000000000000900460ff16606083015260018101546080830152600201546bffffffffffffffffffffffff1660a082015290565b6040516103e49190600060c08201905063ffffffff835116825262ffffff602084015116602083015273ffffffffffffffffffffffffffffffffffffffff604084015116604083015260ff6060840151166060830152608083015160808301526bffffffffffffffffffffffff60a08401511660a083015292915050565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104b3610986366004614b47565b611eb5565b610461610999366004614ad8565b612010565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104206109d2366004614ba7565b612043565b6103da6109e5366004614ad8565b6120c3565b61042061216d565b61053d610a003660046148a9565b6000908152600460205260409020546a0100000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b61042061226f565b60005473ffffffffffffffffffffffffffffffffffffffff1661053d565b60155473ffffffffffffffffffffffffffffffffffffffff1661053d565b610a826122e8565b6040516103e49190614be3565b61053d610a9d3660046148a9565b60009081526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff1690565b6103a46103da565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6104f8610b04366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff908116600090815260226020526040902054670100000000000000900416151590565b610420610b4a366004614c4c565b6123fa565b610420612728565b604051600481526020016103e4565b6115e06103da565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b610420610ba2366004614ad8565b61275a565b60326103da565b6104b3610bbc3660046148a9565b612852565b610bc9612979565b6040516103e49190614c7a565b610bde6129e8565b6040516103e49190614cd4565b60255460ff166107af565b610c09610c043660046148a9565b612bd0565b6040516103e49190614e5b565b7f000000000000000000000000000000000000000000000000000000000000000061053d565b6103da612fe0565b6104b3610c523660046148a9565b612fef565b601b546103da565b610dd86040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101919091525060408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff1661010082015290565b6040516103e49190614f92565b610420610df3366004614c4c565b612ffa565b6103da610e06366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff1660009081526021602052604090205490565b610e88610e3c366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff9091166020830152016103e4565b610420610eaf366004614ad8565b613159565b610f0f610ec23660046148a9565b60408051808201909152600080825260208201525060009081526023602090815260409182902082518084019093525463ffffffff81168352640100000000900462ffffff169082015290565b60408051825163ffffffff16815260209283015162ffffff1692810192909252016103e4565b60406103da565b610f75610f4a366004614ad8565b73ffffffffffffffffffffffffffffffffffffffff166000908152601c602052604090205460ff1690565b6040516103e49190615062565b6103da61316d565b60606000610f986002613175565b9050808410610fd3576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610fdf84866150a5565b905081811180610fed575083155b610ff75780610ff9565b815b9050600061100786836150b8565b67ffffffffffffffff81111561101f5761101f6150cb565b604051908082528060200260200182016040528015611048578160200160208202803683370190505b50905060005b815181101561109b5761106c61106488836150a5565b60029061317f565b82828151811061107e5761107e6150fa565b60209081029190910101528061109381615129565b91505061104e565b50925050505b92915050565b6110af61318b565b6000838152601f602052604090206110c8828483615203565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae776983836040516110fb92919061531e565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c081019190915260408051610140810182526016547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1681526000602082018190529181018290526014546bffffffffffffffffffffffff16606080830191909152918291608081016112416002613175565b81526017547401000000000000000000000000000000000000000080820463ffffffff908116602080860191909152780100000000000000000000000000000000000000000000000080850483166040808801919091526013546060808901919091526014546c01000000000000000000000000810486166080808b0191909152760100000000000000000000000000000000000000000000820460ff16151560a09a8b015283516101e0810185526000808252968101879052601654898104891695820195909552700100000000000000000000000000000000830462ffffff169381019390935273010000000000000000000000000000000000000090910461ffff169082015296870192909252808204831660c08701527c0100000000000000000000000000000000000000000000000000000000909404821660e086015260185492830482166101008601529290910416610120830152601954610140830152601a5461016083015273ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a081016113dc60096131de565b815260175473ffffffffffffffffffffffffffffffffffffffff16602091820152601454600d80546040805182860281018601909152818152949850899489949293600e93750100000000000000000000000000000000000000000090910460ff1692859183018282801561148757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161145c575b50505050509250818054806020026020016040519081016040528092919081815260200182805480156114f057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116114c5575b50505050509150945094509450945094509091929394565b6115106131eb565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601c6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561157057611570614af5565b02179055505050565b6000818152601f6020526040902080546060919061159690615161565b80601f01602080910402602001604051908101604052809291908181526020018280546115c290615161565b801561160f5780601f106115e45761010080835404028352916020019161160f565b820191906000526020600020905b8154815290600101906020018083116115f257829003601f168201915b50505050509050919050565b6000818152601d6020526040902080546060919061159690615161565b6116406131eb565b600e54811461167b576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e5481101561184d576000600e828154811061169d5761169d6150fa565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168084526011909252604083205491935016908585858181106116e7576116e76150fa565b90506020020160208101906116fc9190614ad8565b905073ffffffffffffffffffffffffffffffffffffffff8116158061178f575073ffffffffffffffffffffffffffffffffffffffff82161580159061176d57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b801561178f575073ffffffffffffffffffffffffffffffffffffffff81811614155b156117c6576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146118375773ffffffffffffffffffffffffffffffffffffffff838116600090815260116020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061184590615129565b91505061167e565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e83836040516118829392919061536b565b60405180910390a15050565b6118966131eb565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015282918291829182919082906119ba5760608201516014546000916119a6916bffffffffffffffffffffffff1661541d565b600e549091506119b69082615471565b9150505b8151602083015160408401516119d190849061549c565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b16600090815260116020526040902054929b919a9499509750921694509092505050565b611a1b61326c565b600060255460ff166001811115611a3457611a34614af5565b03611a6b576040517fe0262d7400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e546bffffffffffffffffffffffff909116906000611a8e600f613175565b90506000611a9c82846150a5565b905060008167ffffffffffffffff811115611ab957611ab96150cb565b604051908082528060200260200182016040528015611ae2578160200160208202803683370190505b50905060008267ffffffffffffffff811115611b0057611b006150cb565b604051908082528060200260200182016040528015611b29578160200160208202803683370190505b50905060005b85811015611c5b576000600e8281548110611b4c57611b4c6150fa565b600091825260208220015473ffffffffffffffffffffffffffffffffffffffff169150611b7a828a8a6132bd565b9050806bffffffffffffffffffffffff16858481518110611b9d57611b9d6150fa565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff808416600090815260119092526040909120548551911690859085908110611beb57611beb6150fa565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920181019190915292166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611c5381615129565b915050611b2f565b5060005b84811015611dd8576000611c74600f8361317f565b73ffffffffffffffffffffffffffffffffffffffff8082166000818152600b602090815260408083208151608081018352905460ff80821615158352610100820416828501526bffffffffffffffffffffffff6201000082048116838501526e0100000000000000000000000000009091041660608201529383526011909152902054929350911684611d078a866150a5565b81518110611d1757611d176150fa565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260408101516bffffffffffffffffffffffff1685611d5a8a866150a5565b81518110611d6a57611d6a6150fa565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff9092166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611dd081615129565b915050611c5f565b5073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166000908152602160205260408120819055611e2b600f613175565b90505b8015611e6857611e55611e4d611e456001846150b8565b600f9061317f565b600f906134c5565b5080611e60816154c1565b915050611e2e565b507f5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad8183604051611e9a9291906154f6565b60405180910390a1505050505050565b60006110a1826134e7565b60408051610120810182526014546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260155473ffffffffffffffffffffffffffffffffffffffff16610100820152600090818080611fed84613592565b92509250925061200389858a8a8787878d613784565b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152602080526040902080546060919061159690615161565b61204b61318b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602080526040902061207a828483615203565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d283836040516110fb92919061531e565b73ffffffffffffffffffffffffffffffffffffffff81166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa15801561213f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121639190615524565b6110a191906150b8565b60015473ffffffffffffffffffffffffffffffffffffffff1633146121f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6122776131eb565b601480547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff167601000000000000000000000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258906020016118ea565b600e5460609060008167ffffffffffffffff811115612309576123096150cb565b60405190808252806020026020018201604052801561234e57816020015b60408051808201909152600080825260208201528152602001906001900390816123275790505b50905060005b828110156123f3576000600e8281548110612371576123716150fa565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff9081168084526011835260409384902054845180860190955281855290911691830182905285519093509091908590859081106123d3576123d36150fa565b6020026020010181905250505080806123eb90615129565b915050612354565b5092915050565b73ffffffffffffffffffffffffffffffffffffffff8116612447576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160255460ff16600181111561246057612460614af5565b03612497576040517f4a3578fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601160205260409020541633146124f7576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601454600e5460009161251a9185916bffffffffffffffffffffffff16906132bd565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600b6020908152604080832080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff1690557f000000000000000000000000000000000000000000000000000000000000000090931682526021905220549091506125b1906bffffffffffffffffffffffff8316906150b8565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081166000818152602160205260408082209490945592517fa9059cbb00000000000000000000000000000000000000000000000000000000815291851660048301526bffffffffffffffffffffffff841660248301529063a9059cbb906044016020604051808303816000875af1158015612666573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061268a919061553d565b9050806126c3576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405133815273ffffffffffffffffffffffffffffffffffffffff808516916bffffffffffffffffffffffff8516918716907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a450505050565b6127306131eb565b602580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152601260205260409020541633146127ba576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526011602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556012909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a01000000000000000000009091048116608083015260018301546fffffffffffffffffffffffffffffffff811660a08401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08401527c0100000000000000000000000000000000000000000000000000000000900490931660e08201526002909101549091169181019190915261297283612962816134e7565b8360400151846101000151611eb5565b9392505050565b606060248054806020026020016040519081016040528092919081815260200182805480156129de57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116129b3575b5050505050905090565b604080516102008101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a082018390526101c08201526101e0810191909152604080516102008101825260165463ffffffff74010000000000000000000000000000000000000000808304821684527801000000000000000000000000000000000000000000000000808404831660208601526017547c0100000000000000000000000000000000000000000000000000000000810484169686019690965273ffffffffffffffffffffffffffffffffffffffff938416606086015260145460ff828204161515608087015262ffffff70010000000000000000000000000000000082041660a0870152601854928304841660c087015290820490921660e085015293821661010084015261ffff7301000000000000000000000000000000000000009091041661012083015291909116610140820152601954610160820152601a54610180820152601b546101a08201526101c08101612baa60096131de565b815260155473ffffffffffffffffffffffffffffffffffffffff16602090910152919050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815161012081018352815460ff8082161515835261010080830490911615159583019590955263ffffffff620100008204811694830194909452660100000000000081048416606083015273ffffffffffffffffffffffffffffffffffffffff6a010000000000000000000090910481166080830181905260018401546fffffffffffffffffffffffffffffffff811660a08501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660c08501527c0100000000000000000000000000000000000000000000000000000000900490941660e08301526002909201549091169281019290925290919015612da557816080015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612d7c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612da0919061555f565b612da8565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff168152602001600760008781526020019081526020016000208054612e0090615161565b80601f0160208091040260200160405190810160405280929190818152602001828054612e2c90615161565b8015612e795780601f10612e4e57610100808354040283529160200191612e79565b820191906000526020600020905b815481529060010190602001808311612e5c57829003601f168201915b505050505081526020018360c001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836060015163ffffffff1667ffffffffffffffff1681526020018360e0015163ffffffff1681526020018360a001516bffffffffffffffffffffffff168152602001836000015115158152602001601e60008781526020019081526020016000208054612f5690615161565b80601f0160208091040260200160405190810160405280929190818152602001828054612f8290615161565b8015612fcf5780601f10612fa457610100808354040283529160200191612fcf565b820191906000526020600020905b815481529060010190602001808311612fb257829003601f168201915b505050505081525092505050919050565b6000612fea613af4565b905090565b60006110a182612852565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526011602052604090205416331461305a576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216036130a9576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601260205260409020548116908216146131555773ffffffffffffffffffffffffffffffffffffffff82811660008181526012602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45b5050565b6131616131eb565b61316a81613bbe565b50565b6000612fea60025b60006110a1825490565b60006129728383613cb3565b60175473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6060600061297283613cdd565b60005473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016121ea565b60185473ffffffffffffffffffffffffffffffffffffffff1633146131dc576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906134b9576000816060015185613355919061541d565b905060006133638583615471565b90508083604001818151613377919061549c565b6bffffffffffffffffffffffff16905250613392858261557c565b836060018181516133a3919061549c565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b60006129728373ffffffffffffffffffffffffffffffffffffffff8416613d38565b6000818160045b600f811015613574577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061352c5761352c6150fa565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461356257506000949350505050565b8061356c81615129565b9150506134ee565b5081600f1a600181111561358a5761358a614af5565b949350505050565b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561361f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061364391906155c3565b509450909250505060008113158061365a57508142105b8061367b575082801561367b575061367282426150b8565b8463ffffffff16105b1561368a57601954965061368e565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156136f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061371d91906155c3565b509450909250505060008113158061373457508142105b806137555750828015613755575061374c82426150b8565b8463ffffffff16105b1561376457601a549550613768565b8095505b86866137738a613e2b565b965096509650505050509193909250565b600080808089600181111561379b5761379b614af5565b036137aa575062017f986137ff565b60018960018111156137be576137be614af5565b036137cd57506201e26c6137ff565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a6080015160016138129190615613565b6138209060ff16604061562c565b60185461384e906103a49074010000000000000000000000000000000000000000900463ffffffff166150a5565b61385891906150a5565b601554604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa1580156138c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138ed9190615643565b909250905081836138ff8360186150a5565b613909919061562c565b60808f0151613919906001615613565b6139289060ff166115e061562c565b61393291906150a5565b61393c91906150a5565b61394690856150a5565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa1580156139b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139dd9190615524565b8d6060015161ffff166139f0919061562c565b94505050506000613a018b86613f15565b60008d815260046020526040902054909150610100900460ff1615613a665760008c81526023602090815260409182902082518084018452905463ffffffff811680835264010000000090910462ffffff908116928401928352928501525116908201525b6000613ad08c6040518061012001604052808d63ffffffff1681526020018681526020018781526020018c81526020018b81526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff16815260200185815260200160001515815250614091565b60208101518151919250613ae39161549c565b9d9c50505050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152602160205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015613b90573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bb49190615524565b612fea9190615667565b3373ffffffffffffffffffffffffffffffffffffffff821603613c3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016121ea565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613cca57613cca6150fa565b9060005260206000200154905092915050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561160f57602002820191906000526020600020905b815481526020019060010190808311613d195750505050509050919050565b60008181526001830160205260408120548015613e21576000613d5c6001836150b8565b8554909150600090613d70906001906150b8565b9050818114613dd5576000866000018281548110613d9057613d906150fa565b9060005260206000200154905080876000018481548110613db357613db36150fa565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613de657613de6615687565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506110a1565b60009150506110a1565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613e9b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ebf91906155c3565b50935050925050600082131580613ed557508042105b80613f0557506000846040015162ffffff16118015613f055750613ef981426150b8565b846040015162ffffff16105b156123f3575050601b5492915050565b60408051608081018252600080825260208083018281528385018381526060850184905273ffffffffffffffffffffffffffffffffffffffff878116855260229093528584208054640100000000810462ffffff1690925263ffffffff82169092527b01000000000000000000000000000000000000000000000000000000810460ff16855285517ffeaf968c00000000000000000000000000000000000000000000000000000000815295519495919484936701000000000000009092049091169163feaf968c9160048083019260a09291908290030181865afa158015614002573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061402691906155c3565b5093505092505060008213158061403c57508042105b8061406c57506000866040015162ffffff1611801561406c575061406081426150b8565b866040015162ffffff16105b156140805760018301546060850152614088565b606084018290525b50505092915050565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081019190915260008260e001516000015160ff1690506000846060015161ffff1684606001516140fc919061562c565b9050836101000151801561410f5750803a105b1561411757503a5b60006012831161412857600161413e565b6141336012846150b8565b61413e90600a6157d6565b9050600060128410614151576001614167565b61415c8460126150b8565b61416790600a6157d6565b905060008660a0015187604001518860200151896000015161418991906150a5565b614193908761562c565b61419d91906150a5565b6141a7919061562c565b9050614203828860e00151606001516141c0919061562c565b6001848a60e00151606001516141d6919061562c565b6141e091906150b8565b6141ea868561562c565b6141f491906150a5565b6141fe91906157e2565b6143d8565b6bffffffffffffffffffffffff1686526080870151614226906141fe90836157e2565b6bffffffffffffffffffffffff1660408088019190915260e0880151015160009061425f9062ffffff16683635c9adc5dea0000061562c565b9050600081633b9aca008a60a001518b60e001516020015163ffffffff168c604001518d600001518b614292919061562c565b61429c91906150a5565b6142a6919061562c565b6142b0919061562c565b6142ba91906157e2565b6142c491906150a5565b9050614307848a60e00151606001516142dd919061562c565b6001868c60e00151606001516142f3919061562c565b6142fd91906150b8565b6141ea888561562c565b6bffffffffffffffffffffffff166020890152608089015161432d906141fe90836157e2565b6bffffffffffffffffffffffff16606089015260c089015173ffffffffffffffffffffffffffffffffffffffff166080808a0191909152890151614370906143d8565b6bffffffffffffffffffffffff1660a0808a0191909152890151614393906143d8565b6bffffffffffffffffffffffff1660c089015260e0890151606001516143b8906143d8565b6bffffffffffffffffffffffff1660e08901525050505050505092915050565b60006bffffffffffffffffffffffff821115614476576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016121ea565b5090565b6000806040838503121561448d57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156144cc578151875295820195908201906001016144b0565b509495945050505050565b602081526000612972602083018461449c565b60008083601f8401126144fc57600080fd5b50813567ffffffffffffffff81111561451457600080fd5b60208301915083602082850101111561452c57600080fd5b9250929050565b60008060006040848603121561454857600080fd5b83359250602084013567ffffffffffffffff81111561456657600080fd5b614572868287016144ea565b9497909650939450505050565b600081518084526020808501945080840160005b838110156144cc57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614593565b805163ffffffff16825260006101e060208301516145eb602086018263ffffffff169052565b506040830151614603604086018263ffffffff169052565b50606083015161461a606086018262ffffff169052565b506080830151614630608086018261ffff169052565b5060a083015161465060a08601826bffffffffffffffffffffffff169052565b5060c083015161466860c086018263ffffffff169052565b5060e083015161468060e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a0808401518186018390526146f58387018261457f565b925050506101c0808401516147218287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161475960208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161478360608501826bffffffffffffffffffffffff169052565b506080880151608084015260a08801516147a560a085018263ffffffff169052565b5060c08801516147bd60c085018263ffffffff169052565b5060e088015160e0840152610100808901516147e08286018263ffffffff169052565b5050610120888101511515908401526101408301819052614803818401886145c5565b9050828103610160840152614818818761457f565b905082810361018084015261482d818661457f565b9150506148406101a083018460ff169052565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461316a57600080fd5b6000806040838503121561487f57600080fd5b823561488a8161484a565b915060208301356004811061489e57600080fd5b809150509250929050565b6000602082840312156148bb57600080fd5b5035919050565b6000815180845260005b818110156148e8576020818501810151868301820152016148cc565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061297260208301846148c2565b815173ffffffffffffffffffffffffffffffffffffffff1681526101608101602083015161496f602084018263ffffffff169052565b506040830151614987604084018263ffffffff169052565b50606083015161499f606084018263ffffffff169052565b5060808301516149c7608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a08301516149df60a084018263ffffffff169052565b5060c08301516149f760c084018263ffffffff169052565b5060e0830151614a0f60e084018263ffffffff169052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff81168483015250506101208381015163ffffffff81168483015250506101408381015163ffffffff8116848301525b505092915050565b60008060208385031215614a7657600080fd5b823567ffffffffffffffff80821115614a8e57600080fd5b818501915085601f830112614aa257600080fd5b813581811115614ab157600080fd5b8660208260051b8501011115614ac657600080fd5b60209290920196919550909350505050565b600060208284031215614aea57600080fd5b81356129728161484a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6002811061316a5761316a614af5565b60208101614b4183614b24565b91905290565b60008060008060808587031215614b5d57600080fd5b84359350602085013560028110614b7357600080fd5b9250604085013563ffffffff81168114614b8c57600080fd5b91506060850135614b9c8161484a565b939692955090935050565b600080600060408486031215614bbc57600080fd5b8335614bc78161484a565b9250602084013567ffffffffffffffff81111561456657600080fd5b602080825282518282018190526000919060409081850190868401855b82811015614c3f578151805173ffffffffffffffffffffffffffffffffffffffff90811686529087015116868501529284019290850190600101614c00565b5091979650505050505050565b60008060408385031215614c5f57600080fd5b8235614c6a8161484a565b9150602083013561489e8161484a565b6020808252825182820181905260009190848201906040850190845b81811015614cc857835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101614c96565b50909695505050505050565b60208152614ceb60208201835163ffffffff169052565b60006020830151614d04604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015173ffffffffffffffffffffffffffffffffffffffff8116608084015250608083015180151560a08401525060a083015162ffffff811660c08401525060c083015163ffffffff811660e08401525060e0830151610100614d838185018363ffffffff169052565b8401519050610120614dac8482018373ffffffffffffffffffffffffffffffffffffffff169052565b8401519050610140614dc38482018361ffff169052565b8401519050610160614dec8482018373ffffffffffffffffffffffffffffffffffffffff169052565b840151610180848101919091528401516101a0808501919091528401516101c0808501919091528401516102006101e080860182905291925090614e3461022086018461457f565b9086015173ffffffffffffffffffffffffffffffffffffffff811683870152909250614721565b60208152614e8260208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614e9b604084018263ffffffff169052565b506040830151610140806060850152614eb86101608501836148c2565b91506060850151614ed960808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614f45818701836bffffffffffffffffffffffff169052565b8601519050610120614f5a8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061484083826148c2565b6000610120820190506bffffffffffffffffffffffff835116825263ffffffff60208401511660208301526040830151614fd3604084018262ffffff169052565b506060830151614fe9606084018261ffff169052565b506080830151614ffe608084018260ff169052565b5060a083015161501260a084018215159052565b5060c083015161502660c084018215159052565b5060e083015161503a60e084018215159052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff811684830152614a5b565b6020810160048310614b4157614b41614af5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156110a1576110a1615076565b818103818111156110a1576110a1615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361515a5761515a615076565b5060010190565b600181811c9082168061517557607f821691505b6020821081036151ae577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156151fe57600081815260208120601f850160051c810160208610156151db5750805b601f850160051c820191505b818110156151fa578281556001016151e7565b5050505b505050565b67ffffffffffffffff83111561521b5761521b6150cb565b61522f836152298354615161565b836151b4565b6000601f841160018114615281576000851561524b5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355615317565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156152d057868501358255602094850194600190920191016152b0565b508682101561530b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000604082016040835280865480835260608501915087600052602092508260002060005b828110156153c257815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101615390565b505050838103828501528481528590820160005b868110156154115782356153e98161484a565b73ffffffffffffffffffffffffffffffffffffffff16825291830191908301906001016153d6565b50979650505050505050565b6bffffffffffffffffffffffff8281168282160390808211156123f3576123f3615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff8084168061549057615490615442565b92169190910492915050565b6bffffffffffffffffffffffff8181168382160190808211156123f3576123f3615076565b6000816154d0576154d0615076565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b604081526000615509604083018561457f565b828103602084015261551b818561449c565b95945050505050565b60006020828403121561553657600080fd5b5051919050565b60006020828403121561554f57600080fd5b8151801515811461297257600080fd5b60006020828403121561557157600080fd5b81516129728161484a565b6bffffffffffffffffffffffff818116838216028082169190828114614a5b57614a5b615076565b805169ffffffffffffffffffff811681146155be57600080fd5b919050565b600080600080600060a086880312156155db57600080fd5b6155e4866155a4565b9450602086015193506040860151925060608601519150615607608087016155a4565b90509295509295909350565b60ff81811683821601908111156110a1576110a1615076565b80820281158282048414176110a1576110a1615076565b6000806040838503121561565657600080fd5b505080516020909101519092909150565b81810360008312801583831316838312821617156123f3576123f3615076565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b600181815b8085111561570f57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156156f5576156f5615076565b8085161561570257918102915b93841c93908002906156bb565b509250929050565b600082615726575060016110a1565b81615733575060006110a1565b816001811461574957600281146157535761576f565b60019150506110a1565b60ff84111561576457615764615076565b50506001821b6110a1565b5060208310610133831016604e8410600b8410161715615792575081810a6110a1565b61579c83836156b6565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156157ce576157ce615076565b029392505050565b60006129728383615717565b6000826157f1576157f1615442565b50049056fea164736f6c6343000813000a",
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

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCaller) GetBillingOverridesEnabled(opts *bind.CallOpts, upkeepID *big.Int) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicC.contract.Call(opts, &out, "getBillingOverridesEnabled", upkeepID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCSession) GetBillingOverridesEnabled(upkeepID *big.Int) (bool, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingOverridesEnabled(&_AutomationRegistryLogicC.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicC *AutomationRegistryLogicCCallerSession) GetBillingOverridesEnabled(upkeepID *big.Int) (bool, error) {
	return _AutomationRegistryLogicC.Contract.GetBillingOverridesEnabled(&_AutomationRegistryLogicC.CallOpts, upkeepID)
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
