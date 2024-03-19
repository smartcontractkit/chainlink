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
	GasFeePPB        uint32
	FlatFeeMicroLink *big.Int
	PriceFeed        common.Address
	FallbackPrice    *big.Int
	MinSpend         *big.Int
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

var AutomationRegistryLogicBMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nativeUSDFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"automationForwarderLogic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedReadOnlyAddress\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"payoutMode\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBillingToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFeed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOffchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSettleOnchain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFinanceAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeesWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"}],\"name\":\"NOPsSettledOffchain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableOffchainPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"getAdminPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedReadOnlyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAutomationForwarderLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getBillingToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBillingTokenConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"gasFeePPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint24\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"priceFeed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fallbackPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minSpend\",\"type\":\"uint96\"}],\"internalType\":\"structAutomationRegistryBase2_3.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingTokens\",\"outputs\":[{\"internalType\":\"contractIERC20[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCancellationDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainModule\",\"outputs\":[{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConditionalGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackNativePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"contractIAutomationForwarder\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHotVars\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuard\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"}],\"internalType\":\"structAutomationRegistryBase2_3.HotVars\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLogGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"contractIERC20\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNativeUSDFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumUpkeeps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPayoutMode\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.PayoutMode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerPerformByteGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerSignerGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReorgProtectionEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"billingToken\",\"type\":\"address\"}],\"name\":\"getReserveAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStorage\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"financeAdmin\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"}],\"internalType\":\"structAutomationRegistryBase2_3.Storage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumAutomationRegistryBase2_3.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setAdminPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_3.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settleNOPsOffchain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"supportsBillingToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20Fees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLinkFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620063ca380380620063ca8339810160408190526200003591620002eb565b868686868686863380600081620000935760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c657620000c68162000223565b5050506001600160a01b0380881660805286811660a05285811660c05284811660e052838116610100528216610120526022805482919060ff19166001838181111562000117576200011762000387565b021790555060c0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200015d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200018391906200039d565b60ff1660a0516001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001c7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ed91906200039d565b60ff16146200020f576040516301f86e1760e41b815260040160405180910390fd5b5050505050505050505050505050620003c9565b336001600160a01b038216036200027d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620002e657600080fd5b919050565b600080600080600080600060e0888a0312156200030757600080fd5b6200031288620002ce565b96506200032260208901620002ce565b95506200033260408901620002ce565b94506200034260608901620002ce565b93506200035260808901620002ce565b92506200036260a08901620002ce565b915060c0880151600281106200037757600080fd5b8091505092959891949750929550565b634e487b7160e01b600052602160045260246000fd5b600060208284031215620003b057600080fd5b815160ff81168114620003c257600080fd5b9392505050565b60805160a05160c05160e0516101005161012051615f79620004516000396000610a3a015260006108730152600081816108d20152613ef7015260008181610899015261484901526000818161057f0152613fd1015260008181610ca001528181612a3801528181612ae201528181612f3b01528181612f980152613b270152615f796000f3fe608060405234801561001057600080fd5b50600436106104575760003560e01c80638456cb5911610250578063b6511a2a11610150578063d85aa07c116100c8578063ed56b3e111610097578063f777ff061161007c578063f777ff0614610f4f578063faa3e99614610f56578063ffd242bd14610f9c57600080fd5b8063ed56b3e114610ec9578063f2fde38b14610f3c57600080fd5b8063d85aa07c14610cf2578063e80a3d4414610cfa578063eb5dcd6c14610e80578063ec4de4ba14610e9357600080fd5b8063c7c3a19a1161011f578063cd7f71b511610104578063cd7f71b514610cc4578063d09dc33914610cd7578063d763264814610cdf57600080fd5b8063c7c3a19a14610c7e578063ca30e60314610c9e57600080fd5b8063b6511a2a14610c44578063b657bc9c14610c4b578063ba87666814610c5e578063c5b964e014610c7357600080fd5b8063a538b2eb116101e3578063aab9edd6116101b2578063b121e14711610197578063b121e14714610ae8578063b148ab6b14610afb578063b3596c2314610b0e57600080fd5b8063aab9edd614610ad1578063abc76ae014610ae057600080fd5b8063a538b2eb14610a5e578063a710b22114610aa3578063a72aa27e14610ab6578063a87f45fe14610ac957600080fd5b80638ed02bab1161021f5780638ed02bab146109d957806393f6ebcf146109f75780639e0a99ed14610a30578063a08714c014610a3857600080fd5b80638456cb591461098d5780638765ecbe146109955780638da5cb5b146109a85780638dcf0fe7146109c657600080fd5b806343cc055c1161035b578063614486af116102ee57806368d369d8116102bd57806379ba5097116102a257806379ba50971461092f57806379ea9943146109375780638081fadb1461097a57600080fd5b806368d369d814610909578063744bfe611461091c57600080fd5b8063614486af146108975780636209e1e9146108bd5780636709d0e5146108d0578063671d36ed146108f657600080fd5b80634ee88d351161032a5780634ee88d351461082b5780635147cd591461083e5780635165f2f51461085e5780635425d8ac1461087157600080fd5b806343cc055c146107be57806344cb70b8146107f157806348013d7b146108145780634ca16c521461082357600080fd5b8063207b6516116103ee5780633408f73a116103bd5780633f4ba83a116103a25780633f4ba83a14610748578063421d183b1461075057806343b46e5f146107b657600080fd5b80633408f73a146105de5780633b9cce591461073557600080fd5b8063207b65161461056a578063226cf83c1461057d578063232c1cc5146105c45780632694fbd1146105cb57600080fd5b8063187256e81161042a578063187256e8146104c257806319d97a94146104d55780631a2af011146104f55780631e0104391461050857600080fd5b8063050ee65d1461045c57806306e3b632146104745780630b7d33e6146104945780631865c57d146104a9575b600080fd5b62014c085b6040519081526020015b60405180910390f35b610487610482366004614e85565b610fa4565b60405161046b9190614ee2565b6104a76104a2366004614f3e565b6110c1565b005b6104b161116b565b60405161046b959493929190615136565b6104a76104d0366004615277565b61156b565b6104e86104e33660046152b4565b6115dc565b60405161046b9190615331565b6104a7610503366004615344565b61167e565b61054d6105163660046152b4565b60009081526004602052604090206001015470010000000000000000000000000000000090046bffffffffffffffffffffffff1690565b6040516bffffffffffffffffffffffff909116815260200161046b565b6104e86105783660046152b4565b611784565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161046b565b6018610461565b61054d6105d9366004615382565b6117a1565b6107286040805161016081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081019190915250604080516101608101825260145473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008084048216602086015278010000000000000000000000000000000000000000000000008085048316968601969096527c010000000000000000000000000000000000000000000000000000000093849004821660608601526015548084166080870152818104831660a0870152868104831660c087015293909304811660e085015260165491821661010085015291810482166101208401529290920490911661014082015290565b60405161046b91906153cf565b6104a76107433660046154f9565b6118fa565b6104a7611b50565b61076361075e36600461556e565b611bb6565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a00161046b565b6104a7611cd5565b6012547801000000000000000000000000000000000000000000000000900460ff165b604051901515815260200161046b565b6107e16107ff3660046152b4565b60009081526008602052604090205460ff1690565b600060405161046b91906155ba565b61ea60610461565b6104a7610839366004614f3e565b611e9d565b61085161084c3660046152b4565b611ef2565b60405161046b91906155e4565b6104a761086c3660046152b4565b611efd565b7f000000000000000000000000000000000000000000000000000000000000000061059f565b7f000000000000000000000000000000000000000000000000000000000000000061059f565b6104e86108cb36600461556e565b612097565b7f000000000000000000000000000000000000000000000000000000000000000061059f565b6104a76109043660046155f1565b6120cb565b6104a761091736600461562d565b612195565b6104a761092a366004615344565b61232d565b6104a7612842565b61059f6109453660046152b4565b6000908152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b6104a761098836600461566e565b612944565b6104a7612b5f565b6104a76109a33660046152b4565b612bd8565b60005473ffffffffffffffffffffffffffffffffffffffff1661059f565b6104a76109d4366004614f3e565b612d75565b60135473ffffffffffffffffffffffffffffffffffffffff1661059f565b61059f610a053660046152b4565b60009081526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff1690565b6103a4610461565b7f000000000000000000000000000000000000000000000000000000000000000061059f565b6107e1610a6c36600461556e565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152602080526040902054670100000000000000900416151590565b6104a7610ab136600461569a565b612dca565b6104a7610ac43660046156c8565b6130bf565b6104a76131bc565b6040516003815260200161046b565b6115e0610461565b6104a7610af636600461556e565b6131ee565b6104a7610b093660046152b4565b6132e6565b610bd3610b1c36600461556e565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091525073ffffffffffffffffffffffffffffffffffffffff90811660009081526020808052604091829020825160a081018452815463ffffffff8116825262ffffff6401000000008204169382019390935267010000000000000090920490931691810191909152600182015460608201526002909101546bffffffffffffffffffffffff16608082015290565b60405161046b9190600060a08201905063ffffffff835116825262ffffff602084015116602083015273ffffffffffffffffffffffffffffffffffffffff6040840151166040830152606083015160608301526bffffffffffffffffffffffff608084015116608083015292915050565b6032610461565b61054d610c593660046152b4565b6134fb565b610c66613607565b60405161046b91906156f4565b60225460ff16610851565b610c91610c8c3660046152b4565b613676565b60405161046b919061574e565b7f000000000000000000000000000000000000000000000000000000000000000061059f565b6104a7610cd2366004614f3e565b613a6e565b610461613b25565b61054d610ced3660046152b4565b613bf4565b601954610461565b610e736040805161012081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101919091525060408051610120810182526012546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260135473ffffffffffffffffffffffffffffffffffffffff1661010082015290565b60405161046b9190615885565b6104a7610e8e36600461569a565b613bff565b610461610ea136600461556e565b73ffffffffffffffffffffffffffffffffffffffff166000908152601a602052604090205490565b610f23610ed736600461556e565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff90911660208301520161046b565b6104a7610f4a36600461556e565b613d5d565b6040610461565b610f8f610f6436600461556e565b73ffffffffffffffffffffffffffffffffffffffff166000908152601b602052604090205460ff1690565b60405161046b9190615955565b610461613d71565b60606000610fb26002613d79565b9050808410610fed576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610ff98486615998565b905081811180611007575083155b6110115780611013565b815b9050600061102186836159ab565b67ffffffffffffffff811115611039576110396159be565b604051908082528060200260200182016040528015611062578160200160208202803683370190505b50905060005b81518110156110b55761108661107e8883615998565b600290613d83565b828281518110611098576110986159ed565b6020908102919091010152806110ad81615a1c565b915050611068565b50925050505b92915050565b60155473ffffffffffffffffffffffffffffffffffffffff163314611112576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152601e6020526040902061112b828483615af6565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769838360405161115e929190615c11565b60405180910390a2505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c081019190915260408051610140810182526014547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1681526000602082018190529181018290526012546bffffffffffffffffffffffff16606080830191909152918291608081016112a46002613d79565b81526015547401000000000000000000000000000000000000000080820463ffffffff908116602080860191909152780100000000000000000000000000000000000000000000000080850483166040808801919091526011546060808901919091526012546c01000000000000000000000000810486166080808b0191909152760100000000000000000000000000000000000000000000820460ff16151560a09a8b015283516101e0810185526000808252968101879052601454898104891695820195909552700100000000000000000000000000000000830462ffffff169381019390935273010000000000000000000000000000000000000090910461ffff169082015296870192909252808204831660c08701527c0100000000000000000000000000000000000000000000000000000000909404821660e08601526016549283048216610100860152929091041661012083015260175461014083015260185461016083015273ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a0810161143f6009613d8f565b815260155473ffffffffffffffffffffffffffffffffffffffff16602091820152601254600d80546040805182860281018601909152818152949850899489949293600e93750100000000000000000000000000000000000000000090910460ff169285918301828280156114ea57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116114bf575b505050505092508180548060200260200160405190810160405280929190818152602001828054801561155357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611528575b50505050509150945094509450945094509091929394565b611573613d9c565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601b6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156115d3576115d361558b565b02179055505050565b6000818152601e602052604090208054606091906115f990615a54565b80601f016020809104026020016040519081016040528092919081815260200182805461162590615a54565b80156116725780601f1061164757610100808354040283529160200191611672565b820191906000526020600020905b81548152906001019060200180831161165557829003601f168201915b50505050509050919050565b61168782613e1f565b3373ffffffffffffffffffffffffffffffffffffffff8216036116d6576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116146117805760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b6000818152601c602052604090208054606091906115f990615a54565b60408051610120810182526012546bffffffffffffffffffffffff8116825263ffffffff6c01000000000000000000000000820416602083015262ffffff7001000000000000000000000000000000008204169282019290925261ffff730100000000000000000000000000000000000000830416606082015260ff750100000000000000000000000000000000000000000083048116608083015276010000000000000000000000000000000000000000000083048116151560a08301527701000000000000000000000000000000000000000000000083048116151560c08301527801000000000000000000000000000000000000000000000000909204909116151560e082015260135473ffffffffffffffffffffffffffffffffffffffff166101008201526000908180806118d984613ed3565b9250925092506118ee8489898686868c6140c5565b98975050505050505050565b611902613d9c565b600e54811461193d576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e54811015611b0f576000600e828154811061195f5761195f6159ed565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600f909252604083205491935016908585858181106119a9576119a96159ed565b90506020020160208101906119be919061556e565b905073ffffffffffffffffffffffffffffffffffffffff81161580611a51575073ffffffffffffffffffffffffffffffffffffffff821615801590611a2f57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015611a51575073ffffffffffffffffffffffffffffffffffffffff81811614155b15611a88576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811614611af95773ffffffffffffffffffffffffffffffffffffffff8381166000908152600f6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b5050508080611b0790615a1c565b915050611940565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e8383604051611b4493929190615ca9565b60405180910390a15050565b611b58613d9c565b601280547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e010000000000000000000000000000900490911660608201528291829182918291908290611c7c576060820151601254600091611c68916bffffffffffffffffffffffff16615d19565b600e54909150611c789082615d6d565b9150505b815160208301516040840151611c93908490615d98565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b166000908152600f6020526040902054929b919a9499509750921694509092505050565b611cdd6143ac565b600060225460ff166001811115611cf657611cf661558b565b03611d2d576040517fe0262d7400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e5460008167ffffffffffffffff811115611d4b57611d4b6159be565b604051908082528060200260200182016040528015611d74578160200160208202803683370190505b50905060005b82811015611e6a576000600e8281548110611d9757611d976159ed565b600091825260208220015460125473ffffffffffffffffffffffffffffffffffffffff9091169250611dd99083906bffffffffffffffffffffffff16876143fd565b9050806bffffffffffffffffffffffff16848481518110611dfc57611dfc6159ed565b60209081029190910181019190915273ffffffffffffffffffffffffffffffffffffffff9092166000908152600b909252506040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905580611e6281615a1c565b915050611d7a565b507f5af23b715253628d12b660b27a4f3fc626562ea8a55040aa99ab3dc178989fad600e82604051611b44929190615dbd565b611ea683613e1f565b6000838152601c60205260409020611ebf828483615af6565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664838360405161115e929190615c11565b60006110bb82614605565b611f0681613e1f565b60008181526004602090815260409182902082516101008082018552825460ff8116151580845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546fffffffffffffffffffffffffffffffff811660808401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660a08401527c010000000000000000000000000000000000000000000000000000000090041660c082015260029091015490921660e0830152612028576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556120676002836146b0565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152601f602052604090208054606091906115f990615a54565b60155473ffffffffffffffffffffffffffffffffffffffff16331461211c576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152601f6020526040902061214c828483615af6565b508273ffffffffffffffffffffffffffffffffffffffff167f7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2838360405161115e929190615c11565b61219d6143ac565b73ffffffffffffffffffffffffffffffffffffffff82166121ea576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018390526000919085169063a9059cbb906044016020604051808303816000875af1158015612263573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122879190615deb565b9050806122c0576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8373ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa88460405161231f91815260200190565b60405180910390a350505050565b60125477010000000000000000000000000000000000000000000000900460ff1615612385576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff167701000000000000000000000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116612414576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff91810482168387015265010000000000810482168386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606084015260018401546fffffffffffffffffffffffffffffffff811660808501526bffffffffffffffffffffffff70010000000000000000000000000000000082041660a08501527c0100000000000000000000000000000000000000000000000000000000900490911660c0830152600290920154821660e08201528685526005909352922054909116331461253f576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601354604080517f57e871e7000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff909216916357e871e7916004808201926020929091908290030181865afa1580156125af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125d39190615e0d565b816040015163ffffffff161115612616576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526004602090815260408083206001015460e085015173ffffffffffffffffffffffffffffffffffffffff168452601a909252909120547001000000000000000000000000000000009091046bffffffffffffffffffffffff16906126809082906159ab565b60e08301805173ffffffffffffffffffffffffffffffffffffffff9081166000908152601a602090815260408083209590955588825260049081905284822060010180547fffffffff000000000000000000000000ffffffffffffffffffffffffffffffff169055925193517fa9059cbb000000000000000000000000000000000000000000000000000000008152878316938101939093526bffffffffffffffffffffffff8516602484015292169063a9059cbb906044016020604051808303816000875af1158015612758573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061277c9190615deb565b9050806127b5576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516bffffffffffffffffffffffff8416815273ffffffffffffffffffffffffffffffffffffffff8616602082015286917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601280547fffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffff169055505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146128c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61294c6143ac565b73ffffffffffffffffffffffffffffffffffffffff8216612999576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006129a3613b25565b9050808211156129e9576040517fcf47918100000000000000000000000000000000000000000000000000000000815260048101829052602481018390526044016128bf565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490526000917f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015612a83573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612aa79190615deb565b905080612ae0576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f5e110f8bc8a20b65dcc87f224bdf1cc039346e267118bae2739847f07321ffa88560405161231f91815260200190565b612b67613d9c565b601280547fffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffff167601000000000000000000000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001611bac565b612be181613e1f565b60008181526004602090815260409182902082516101008082018552825460ff8116158015845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546fffffffffffffffffffffffffffffffff811660808401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660a08401527c010000000000000000000000000000000000000000000000000000000090041660c082015260029091015490921660e0830152612d03576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055612d456002836146bc565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b612d7e83613e1f565b6000838152601d60205260409020612d97828483615af6565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850838360405161115e929190615c11565b73ffffffffffffffffffffffffffffffffffffffff8116612e17576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160225460ff166001811115612e3057612e3061558b565b03612e67576040517f4a3578fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612ec7576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e54600091612eea9185916bffffffffffffffffffffffff16906143fd565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600b6020908152604080832080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff1690557f00000000000000000000000000000000000000000000000000000000000000009093168252601a90522054909150612f81906bffffffffffffffffffffffff8316906159ab565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000081166000818152601a6020526040908190209390935591517fa9059cbb00000000000000000000000000000000000000000000000000000000815290841660048201526bffffffffffffffffffffffff8316602482015263a9059cbb906044016020604051808303816000875af1158015613036573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061305a9190615deb565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff1610806130fc575060145463ffffffff78010000000000000000000000000000000000000000000000009091048116908216115b15613133576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61313c82613e1f565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff1661010063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b6131c4613d9c565b602280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526010602052604090205416331461324e576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600f602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556010909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60008181526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821694830194909452650100000000008404811694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606083015260018301546fffffffffffffffffffffffffffffffff811660808401526bffffffffffffffffffffffff70010000000000000000000000000000000082041660a08401527c01000000000000000000000000000000000000000000000000000000009004811660c083015260029092015490921660e083015290911461340a576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613467576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b600081815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff91810482169583019590955265010000000000850481169382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606082015260018201546fffffffffffffffffffffffffffffffff811660808301526bffffffffffffffffffffffff70010000000000000000000000000000000082041660a08301527c0100000000000000000000000000000000000000000000000000000000900490921660c08301526002015490911660e08201526136006135f184614605565b82602001518360e001516117a1565b9392505050565b6060602180548060200260200160405190810160405280929190818152602001828054801561366c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613641575b5050505050905090565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e08201819052610100820152610120810191909152600082815260046020908152604080832081516101008082018452825460ff81161515835290810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff9081166060840181905260018301546fffffffffffffffffffffffffffffffff8116608086015270010000000000000000000000000000000081046bffffffffffffffffffffffff1660a08601527c0100000000000000000000000000000000000000000000000000000000900490941660c08401526002909101541660e082015291901561383357816060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561380a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061382e9190615e26565b613836565b60005b90506040518061014001604052808273ffffffffffffffffffffffffffffffffffffffff168152602001836020015163ffffffff16815260200160076000878152602001908152602001600020805461388e90615a54565b80601f01602080910402602001604051908101604052809291908181526020018280546138ba90615a54565b80156139075780601f106138dc57610100808354040283529160200191613907565b820191906000526020600020905b8154815290600101906020018083116138ea57829003601f168201915b505050505081526020018360a001516bffffffffffffffffffffffff1681526020016005600087815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001836040015163ffffffff1667ffffffffffffffff1681526020018360c0015163ffffffff16815260200183608001516bffffffffffffffffffffffff168152602001836000015115158152602001601d600087815260200190815260200160002080546139e490615a54565b80601f0160208091040260200160405190810160405280929190818152602001828054613a1090615a54565b8015613a5d5780601f10613a3257610100808354040283529160200191613a5d565b820191906000526020600020905b815481529060010190602001808311613a4057829003601f168201915b505050505081525092505050919050565b613a7783613e1f565b6015547c0100000000000000000000000000000000000000000000000000000000900463ffffffff16811115613ad9576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020613af2828483615af6565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d838360405161115e929190615c11565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166000818152601a60205260408082205490517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152919290916370a0823190602401602060405180830381865afa158015613bc1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613be59190615e0d565b613bef91906159ab565b905090565b60006110bb826134fb565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314613c5f576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603613cae576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601060205260409020548116908216146117805773ffffffffffffffffffffffffffffffffffffffff82811660008181526010602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b613d65613d9c565b613d6e816146c8565b50565b6000613bef60025b60006110bb825490565b600061360083836147bd565b60606000613600836147e7565b60005473ffffffffffffffffffffffffffffffffffffffff163314613e1d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016128bf565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613e7c576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff90811614613d6e576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600080846040015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613f60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f849190615e5d565b5094509092505050600081131580613f9b57508142105b80613fbc5750828015613fbc5750613fb382426159ab565b8463ffffffff16105b15613fcb576017549650613fcf565b8096505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561403a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061405e9190615e5d565b509450909250505060008113158061407557508142105b806140965750828015614096575061408d82426159ab565b8463ffffffff16105b156140a55760185495506140a9565b8095505b86866140b48a614842565b965096509650505050509193909250565b60008080808960018111156140dc576140dc61558b565b036140ea575061ea6061413f565b60018960018111156140fe576140fe61558b565b0361410d575062014c0861413f565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a6080015160016141529190615ead565b6141609060ff166040615ec6565b60165461418e906103a49074010000000000000000000000000000000000000000900463ffffffff16615998565b6141989190615998565b601354604080517fde9ee35e0000000000000000000000000000000000000000000000000000000081528151939450600093849373ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa158015614209573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061422d9190615edd565b9092509050818361423f836018615998565b6142499190615ec6565b60808f0151614259906001615ead565b6142689060ff166115e0615ec6565b6142729190615998565b61427c9190615998565b6142869085615998565b6101008e01516040517f125441400000000000000000000000000000000000000000000000000000000081526004810186905291955073ffffffffffffffffffffffffffffffffffffffff1690631254414090602401602060405180830381865afa1580156142f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061431d9190615e0d565b8d6060015161ffff166143309190615ec6565b9450505050600061438a8b6040518061010001604052808c63ffffffff1681526020018581526020018681526020018b81526020018a815260200189815260200161437b8f8a614933565b81526000602090910152614ac6565b6020810151815191925061439d91615d98565b9b9a5050505050505050505050565b60165473ffffffffffffffffffffffffffffffffffffffff163314613e1d576040517fb6dfb7a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906145f95760008160600151856144959190615d19565b905060006144a38583615d6d565b905080836040018181516144b79190615d98565b6bffffffffffffffffffffffff169052506144d28582615f01565b836060018181516144e39190615d98565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000818160045b600f811015614692577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061464a5761464a6159ed565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461468057506000949350505050565b8061468a81615a1c565b91505061460c565b5081600f1a60018111156146a8576146a861558b565b949350505050565b60006136008383614ca1565b60006136008383614cf0565b3373ffffffffffffffffffffffffffffffffffffffff821603614747576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016128bf565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106147d4576147d46159ed565b9060005260206000200154905092915050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561167257602002820191906000526020600020905b8154815260200190600101908083116148235750505050509050919050565b60008060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156148b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906148d69190615e5d565b509350509250506000821315806148ec57508042105b8061491c57506000846040015162ffffff1611801561491c575061491081426159ab565b846040015162ffffff16105b1561492c57505060195492915050565b5092915050565b604080516060810182526000808252602082018190529181019190915273ffffffffffffffffffffffffffffffffffffffff808316600090815260208080526040808320815160a08082018452825463ffffffff808216845262ffffff6401000000008304168488018190526701000000000000009092048916848701908152600186015460608601526002909501546bffffffffffffffffffffffff1660808501529589015281519094168752905182517ffeaf968c00000000000000000000000000000000000000000000000000000000815292519195859491169263feaf968c92600482810193928290030181865afa158015614a37573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614a5b9190615e5d565b50935050925050600082131580614a7157508042105b80614aa157506000866040015162ffffff16118015614aa15750614a9581426159ab565b866040015162ffffff16105b15614ab55760608301516040850152614abd565b604084018290525b50505092915050565b6040805160808101825260008082526020820181905291810182905260608101919091526000836060015161ffff168360600151614b049190615ec6565b90508260e001518015614b165750803a105b15614b1e57503a5b60008360a00151846040015185602001518660000151614b3e9190615998565b614b489085615ec6565b614b529190615998565b614b5c9190615ec6565b9050614b7a8460c001516040015182614b759190615f29565b614de3565b6bffffffffffffffffffffffff1683526080840151614b9d90614b759083615f29565b6bffffffffffffffffffffffff166040840152608084015160c08501516020015160009190614bd69062ffffff1664e8d4a51000615ec6565b614be09190615ec6565b9050600081633b9aca008760a001518860c001516000015163ffffffff1689604001518a6000015189614c139190615ec6565b614c1d9190615998565b614c279190615ec6565b614c319190615ec6565b614c3b9190615f29565b614c459190615998565b9050614c5e8660c001516040015182614b759190615f29565b6bffffffffffffffffffffffff1660208601526080860151614c8490614b759083615f29565b6bffffffffffffffffffffffff1660608601525050505092915050565b6000818152600183016020526040812054614ce8575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556110bb565b5060006110bb565b60008181526001830160205260408120548015614dd9576000614d146001836159ab565b8554909150600090614d28906001906159ab565b9050818114614d8d576000866000018281548110614d4857614d486159ed565b9060005260206000200154905080876000018481548110614d6b57614d6b6159ed565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080614d9e57614d9e615f3d565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506110bb565b60009150506110bb565b60006bffffffffffffffffffffffff821115614e81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203960448201527f362062697473000000000000000000000000000000000000000000000000000060648201526084016128bf565b5090565b60008060408385031215614e9857600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b83811015614ed757815187529582019590820190600101614ebb565b509495945050505050565b6020815260006136006020830184614ea7565b60008083601f840112614f0757600080fd5b50813567ffffffffffffffff811115614f1f57600080fd5b602083019150836020828501011115614f3757600080fd5b9250929050565b600080600060408486031215614f5357600080fd5b83359250602084013567ffffffffffffffff811115614f7157600080fd5b614f7d86828701614ef5565b9497909650939450505050565b600081518084526020808501945080840160005b83811015614ed757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614f9e565b805163ffffffff16825260006101e06020830151614ff6602086018263ffffffff169052565b50604083015161500e604086018263ffffffff169052565b506060830151615025606086018262ffffff169052565b50608083015161503b608086018261ffff169052565b5060a083015161505b60a08601826bffffffffffffffffffffffff169052565b5060c083015161507360c086018263ffffffff169052565b5060e083015161508b60e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a08084015181860183905261510083870182614f8a565b925050506101c08084015161512c8287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161516460208501826bffffffffffffffffffffffff169052565b5060408801516040840152606088015161518e60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a08801516151b060a085018263ffffffff169052565b5060c08801516151c860c085018263ffffffff169052565b5060e088015160e0840152610100808901516151eb8286018263ffffffff169052565b505061012088810151151590840152610140830181905261520e81840188614fd0565b90508281036101608401526152238187614f8a565b90508281036101808401526152388186614f8a565b91505061524b6101a083018460ff169052565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff81168114613d6e57600080fd5b6000806040838503121561528a57600080fd5b823561529581615255565b91506020830135600481106152a957600080fd5b809150509250929050565b6000602082840312156152c657600080fd5b5035919050565b6000815180845260005b818110156152f3576020818501810151868301820152016152d7565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061360060208301846152cd565b6000806040838503121561535757600080fd5b8235915060208301356152a981615255565b803563ffffffff8116811461537d57600080fd5b919050565b60008060006060848603121561539757600080fd5b8335600281106153a657600080fd5b92506153b460208501615369565b915060408401356153c481615255565b809150509250925092565b815173ffffffffffffffffffffffffffffffffffffffff16815261016081016020830151615405602084018263ffffffff169052565b50604083015161541d604084018263ffffffff169052565b506060830151615435606084018263ffffffff169052565b50608083015161545d608084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060a083015161547560a084018263ffffffff169052565b5060c083015161548d60c084018263ffffffff169052565b5060e08301516154a560e084018263ffffffff169052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff81168483015250506101208381015163ffffffff81168483015250506101408381015163ffffffff8116848301525b505092915050565b6000806020838503121561550c57600080fd5b823567ffffffffffffffff8082111561552457600080fd5b818501915085601f83011261553857600080fd5b81358181111561554757600080fd5b8660208260051b850101111561555c57600080fd5b60209290920196919550909350505050565b60006020828403121561558057600080fd5b813561360081615255565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60208101600383106155ce576155ce61558b565b91905290565b60028110613d6e57613d6e61558b565b602081016155ce836155d4565b60008060006040848603121561560657600080fd5b833561561181615255565b9250602084013567ffffffffffffffff811115614f7157600080fd5b60008060006060848603121561564257600080fd5b833561564d81615255565b9250602084013561565d81615255565b929592945050506040919091013590565b6000806040838503121561568157600080fd5b823561568c81615255565b946020939093013593505050565b600080604083850312156156ad57600080fd5b82356156b881615255565b915060208301356152a981615255565b600080604083850312156156db57600080fd5b823591506156eb60208401615369565b90509250929050565b6020808252825182820181905260009190848201906040850190845b8181101561574257835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101615710565b50909695505050505050565b6020815261577560208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161578e604084018263ffffffff169052565b5060408301516101408060608501526157ab6101608501836152cd565b915060608501516157cc60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100615838818701836bffffffffffffffffffffffff169052565b860151905061012061584d8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061524b83826152cd565b6000610120820190506bffffffffffffffffffffffff835116825263ffffffff602084015116602083015260408301516158c6604084018262ffffff169052565b5060608301516158dc606084018261ffff169052565b5060808301516158f1608084018260ff169052565b5060a083015161590560a084018215159052565b5060c083015161591960c084018215159052565b5060e083015161592d60e084018215159052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff8116848301526154f1565b60208101600483106155ce576155ce61558b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156110bb576110bb615969565b818103818111156110bb576110bb615969565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203615a4d57615a4d615969565b5060010190565b600181811c90821680615a6857607f821691505b602082108103615aa1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115615af157600081815260208120601f850160051c81016020861015615ace5750805b601f850160051c820191505b81811015615aed57828155600101615ada565b5050505b505050565b67ffffffffffffffff831115615b0e57615b0e6159be565b615b2283615b1c8354615a54565b83615aa7565b6000601f841160018114615b745760008515615b3e5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355615c0a565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015615bc35786850135825560209485019460019092019101615ba3565b5086821015615bfe577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6000815480845260208085019450836000528060002060005b83811015614ed757815473ffffffffffffffffffffffffffffffffffffffff1687529582019560019182019101615c77565b604081526000615cbc6040830186615c5e565b8281036020848101919091528482528591810160005b86811015615d0d578335615ce581615255565b73ffffffffffffffffffffffffffffffffffffffff1682529282019290820190600101615cd2565b50979650505050505050565b6bffffffffffffffffffffffff82811682821603908082111561492c5761492c615969565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680615d8c57615d8c615d3e565b92169190910492915050565b6bffffffffffffffffffffffff81811683821601908082111561492c5761492c615969565b604081526000615dd06040830185615c5e565b8281036020840152615de28185614ea7565b95945050505050565b600060208284031215615dfd57600080fd5b8151801515811461360057600080fd5b600060208284031215615e1f57600080fd5b5051919050565b600060208284031215615e3857600080fd5b815161360081615255565b805169ffffffffffffffffffff8116811461537d57600080fd5b600080600080600060a08688031215615e7557600080fd5b615e7e86615e43565b9450602086015193506040860151925060608601519150615ea160808701615e43565b90509295509295909350565b60ff81811683821601908111156110bb576110bb615969565b80820281158282048414176110bb576110bb615969565b60008060408385031215615ef057600080fd5b505080516020909101519092909150565b6bffffffffffffffffffffffff8181168382160280821691908281146154f1576154f1615969565b600082615f3857615f38615d3e565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var AutomationRegistryLogicBABI = AutomationRegistryLogicBMetaData.ABI

var AutomationRegistryLogicBBin = AutomationRegistryLogicBMetaData.Bin

func DeployAutomationRegistryLogicB(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkUSDFeed common.Address, nativeUSDFeed common.Address, fastGasFeed common.Address, automationForwarderLogic common.Address, allowedReadOnlyAddress common.Address, payoutMode uint8) (common.Address, *types.Transaction, *AutomationRegistryLogicB, error) {
	parsed, err := AutomationRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicBBin), backend, link, linkUSDFeed, nativeUSDFeed, fastGasFeed, automationForwarderLogic, allowedReadOnlyAddress, payoutMode)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicB.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetActiveUpkeepIDs(&_AutomationRegistryLogicB.CallOpts, startIndex, maxCount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicB.CallOpts, admin)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetAdminPrivilegeConfig(&_AutomationRegistryLogicB.CallOpts, admin)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetAllowedReadOnlyAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetAutomationForwarderLogic(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetBalance(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetBalance(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getBillingToken", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingToken(&_AutomationRegistryLogicB.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetBillingToken(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingToken(&_AutomationRegistryLogicB.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getBillingTokenConfig", token)

	if err != nil {
		return *new(AutomationRegistryBase23BillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23BillingConfig)).(*AutomationRegistryBase23BillingConfig)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicB.CallOpts, token)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetBillingTokenConfig(token common.Address) (AutomationRegistryBase23BillingConfig, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingTokenConfig(&_AutomationRegistryLogicB.CallOpts, token)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getBillingTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingTokens(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetBillingTokens() ([]common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetBillingTokens(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetCancellationDelay(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetCancellationDelay() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetCancellationDelay(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetChainModule(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetChainModule() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetChainModule(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetConditionalGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getFallbackNativePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetFallbackNativePrice() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetFallbackNativePrice(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetFastGasFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetForwarder(&_AutomationRegistryLogicB.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetForwarder(&_AutomationRegistryLogicB.CallOpts, upkeepID)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getHotVars")

	if err != nil {
		return *new(AutomationRegistryBase23HotVars), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23HotVars)).(*AutomationRegistryBase23HotVars)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicB.Contract.GetHotVars(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetHotVars() (AutomationRegistryBase23HotVars, error) {
	return _AutomationRegistryLogicB.Contract.GetHotVars(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetLinkAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetLinkAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetLinkAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getLinkUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetLinkUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetLinkUSDFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetLogGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetLogGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit, billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicB.CallOpts, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMaxPaymentForGas(&_AutomationRegistryLogicB.CallOpts, triggerType, gasLimit, billingToken)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMinBalance(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMinBalance(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetMinBalanceForUpkeep(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetNativeUSDFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getNativeUSDFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetNativeUSDFeedAddress() (common.Address, error) {
	return _AutomationRegistryLogicB.Contract.GetNativeUSDFeedAddress(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetNumUpkeeps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getNumUpkeeps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetNumUpkeeps(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetNumUpkeeps() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetNumUpkeeps(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetPayoutMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getPayoutMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetPayoutMode(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetPayoutMode() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetPayoutMode(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB.CallOpts, peer)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB.CallOpts, peer)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetPerPerformByteGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetPerSignerGasOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicB.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _AutomationRegistryLogicB.Contract.GetReorgProtectionEnabled(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetReserveAmount(opts *bind.CallOpts, billingToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getReserveAmount", billingToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetReserveAmount(&_AutomationRegistryLogicB.CallOpts, billingToken)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetReserveAmount(billingToken common.Address) (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetReserveAmount(&_AutomationRegistryLogicB.CallOpts, billingToken)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicB.Contract.GetSignerInfo(&_AutomationRegistryLogicB.CallOpts, query)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _AutomationRegistryLogicB.Contract.GetSignerInfo(&_AutomationRegistryLogicB.CallOpts, query)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getState")

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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicB.Contract.GetState(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetState() (GetState,

	error) {
	return _AutomationRegistryLogicB.Contract.GetState(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetStorage(opts *bind.CallOpts) (AutomationRegistryBase23Storage, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getStorage")

	if err != nil {
		return *new(AutomationRegistryBase23Storage), err
	}

	out0 := *abi.ConvertType(out[0], new(AutomationRegistryBase23Storage)).(*AutomationRegistryBase23Storage)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicB.Contract.GetStorage(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetStorage() (AutomationRegistryBase23Storage, error) {
	return _AutomationRegistryLogicB.Contract.GetStorage(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitCalldataFixedBytesOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitterInfo(&_AutomationRegistryLogicB.CallOpts, query)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _AutomationRegistryLogicB.Contract.GetTransmitterInfo(&_AutomationRegistryLogicB.CallOpts, query)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetTriggerType(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _AutomationRegistryLogicB.Contract.GetTriggerType(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeep(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeep(&_AutomationRegistryLogicB.CallOpts, id)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _AutomationRegistryLogicB.Contract.GetUpkeepTriggerConfig(&_AutomationRegistryLogicB.CallOpts, upkeepId)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicB.Contract.HasDedupKey(&_AutomationRegistryLogicB.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _AutomationRegistryLogicB.Contract.HasDedupKey(&_AutomationRegistryLogicB.CallOpts, dedupKey)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AutomationRegistryLogicB.Contract.LinkAvailableForPayment(&_AutomationRegistryLogicB.CallOpts)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "supportsBillingToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicB.Contract.SupportsBillingToken(&_AutomationRegistryLogicB.CallOpts, token)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) SupportsBillingToken(token common.Address) (bool, error) {
	return _AutomationRegistryLogicB.Contract.SupportsBillingToken(&_AutomationRegistryLogicB.CallOpts, token)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) UpkeepTranscoderVersion() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.UpkeepTranscoderVersion(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.UpkeepTranscoderVersion(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AutomationRegistryLogicB.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.UpkeepVersion(&_AutomationRegistryLogicB.CallOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBCallerSession) UpkeepVersion() (uint8, error) {
	return _AutomationRegistryLogicB.Contract.UpkeepVersion(&_AutomationRegistryLogicB.CallOpts)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptPayeeship(&_AutomationRegistryLogicB.TransactOpts, transmitter)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.AcceptPayeeship(&_AutomationRegistryLogicB.TransactOpts, transmitter)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "disableOffchainPayments")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.DisableOffchainPayments(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) DisableOffchainPayments() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.DisableOffchainPayments(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "pause")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Pause(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) Pause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Pause(&_AutomationRegistryLogicB.TransactOpts)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicB.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetAdminPrivilegeConfig(&_AutomationRegistryLogicB.TransactOpts, admin, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setPayees", payees)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetPayees(&_AutomationRegistryLogicB.TransactOpts, payees)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetPayees(&_AutomationRegistryLogicB.TransactOpts, payees)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB.TransactOpts, peer, permission)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetPeerRegistryMigrationPermission(&_AutomationRegistryLogicB.TransactOpts, peer, permission)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SetUpkeepPrivilegeConfig(&_AutomationRegistryLogicB.TransactOpts, upkeepId, newPrivilegeConfig)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) SettleNOPsOffchain(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "settleNOPsOffchain")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) SettleNOPsOffchain() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.SettleNOPsOffchain(&_AutomationRegistryLogicB.TransactOpts)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferPayeeship(&_AutomationRegistryLogicB.TransactOpts, transmitter, proposed)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.TransferPayeeship(&_AutomationRegistryLogicB.TransactOpts, transmitter, proposed)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "unpause")
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Unpause(&_AutomationRegistryLogicB.TransactOpts)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) Unpause() (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.Unpause(&_AutomationRegistryLogicB.TransactOpts)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawERC20Fees(opts *bind.TransactOpts, assetAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawERC20Fees", assetAddress, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawERC20Fees(assetAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB.TransactOpts, assetAddress, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawERC20Fees(assetAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawERC20Fees(&_AutomationRegistryLogicB.TransactOpts, assetAddress, to, amount)
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawLinkFees(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawLinkFees", to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawLinkFees(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawLinkFees(&_AutomationRegistryLogicB.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawLinkFees(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawLinkFees(&_AutomationRegistryLogicB.TransactOpts, to, amount)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawPayment(&_AutomationRegistryLogicB.TransactOpts, from, to)
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicB.Contract.WithdrawPayment(&_AutomationRegistryLogicB.TransactOpts, from, to)
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
	Recipient    common.Address
	AssetAddress common.Address
	Amount       *big.Int
	Raw          types.Log
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) FilterFeesWithdrawn(opts *bind.FilterOpts, recipient []common.Address, assetAddress []common.Address) (*AutomationRegistryLogicBFeesWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.FilterLogs(opts, "FeesWithdrawn", recipientRule, assetAddressRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicBFeesWithdrawnIterator{contract: _AutomationRegistryLogicB.contract, event: "FeesWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicB *AutomationRegistryLogicBFilterer) WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFeesWithdrawn, recipient []common.Address, assetAddress []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var assetAddressRule []interface{}
	for _, assetAddressItem := range assetAddress {
		assetAddressRule = append(assetAddressRule, assetAddressItem)
	}

	logs, sub, err := _AutomationRegistryLogicB.contract.WatchLogs(opts, "FeesWithdrawn", recipientRule, assetAddressRule)
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
	Transmitters []common.Address
	Balances     []*big.Int
	Raw          types.Log
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

func (_AutomationRegistryLogicB *AutomationRegistryLogicB) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicB.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicB.ParseAdminPrivilegeConfigSet(log)
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
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error)

	GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error)

	GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetBillingToken(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetBillingTokenConfig(opts *bind.CallOpts, token common.Address) (AutomationRegistryBase23BillingConfig, error)

	GetBillingTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error)

	GetChainModule(opts *bind.CallOpts) (common.Address, error)

	GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetFallbackNativePrice(opts *bind.CallOpts) (*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetHotVars(opts *bind.CallOpts) (AutomationRegistryBase23HotVars, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkUSDFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32, billingToken common.Address) (*big.Int, error)

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

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsBillingToken(opts *bind.CallOpts, token common.Address) (bool, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	DisableOffchainPayments(opts *bind.TransactOpts) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error)

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

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawERC20Fees(opts *bind.TransactOpts, assetAddress common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawLinkFees(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicBAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicBAdminPrivilegeConfigSet, error)

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

	FilterFeesWithdrawn(opts *bind.FilterOpts, recipient []common.Address, assetAddress []common.Address) (*AutomationRegistryLogicBFeesWithdrawnIterator, error)

	WatchFeesWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicBFeesWithdrawn, recipient []common.Address, assetAddress []common.Address) (event.Subscription, error)

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
