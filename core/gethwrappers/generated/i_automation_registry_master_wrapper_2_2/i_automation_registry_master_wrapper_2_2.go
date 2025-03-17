// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_automation_registry_master_wrapper_2_2

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

type AutomationRegistryBase22OnchainConfig struct {
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
	ChainModule            common.Address
	ReorgProtectionEnabled bool
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

var IAutomationRegistryMasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"getAdminPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedReadOnlyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAutomationForwarderLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCancellationDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainModule\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConditionalGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLogGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerPerformByteGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPerSignerGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReorgProtectionEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structIAutomationV21PlusCommon.StateLegacy\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structIAutomationV21PlusCommon.OnchainConfigLegacy\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataFixedBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitCalldataPerSignerBytesOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"hasDedupKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setAdminPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"chainModule\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_2.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var IAutomationRegistryMasterABI = IAutomationRegistryMasterMetaData.ABI

type IAutomationRegistryMaster struct {
	address common.Address
	abi     abi.ABI
	IAutomationRegistryMasterCaller
	IAutomationRegistryMasterTransactor
	IAutomationRegistryMasterFilterer
}

type IAutomationRegistryMasterCaller struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMasterTransactor struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMasterFilterer struct {
	contract *bind.BoundContract
}

type IAutomationRegistryMasterSession struct {
	Contract     *IAutomationRegistryMaster
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IAutomationRegistryMasterCallerSession struct {
	Contract *IAutomationRegistryMasterCaller
	CallOpts bind.CallOpts
}

type IAutomationRegistryMasterTransactorSession struct {
	Contract     *IAutomationRegistryMasterTransactor
	TransactOpts bind.TransactOpts
}

type IAutomationRegistryMasterRaw struct {
	Contract *IAutomationRegistryMaster
}

type IAutomationRegistryMasterCallerRaw struct {
	Contract *IAutomationRegistryMasterCaller
}

type IAutomationRegistryMasterTransactorRaw struct {
	Contract *IAutomationRegistryMasterTransactor
}

func NewIAutomationRegistryMaster(address common.Address, backend bind.ContractBackend) (*IAutomationRegistryMaster, error) {
	abi, err := abi.JSON(strings.NewReader(IAutomationRegistryMasterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIAutomationRegistryMaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMaster{address: address, abi: abi, IAutomationRegistryMasterCaller: IAutomationRegistryMasterCaller{contract: contract}, IAutomationRegistryMasterTransactor: IAutomationRegistryMasterTransactor{contract: contract}, IAutomationRegistryMasterFilterer: IAutomationRegistryMasterFilterer{contract: contract}}, nil
}

func NewIAutomationRegistryMasterCaller(address common.Address, caller bind.ContractCaller) (*IAutomationRegistryMasterCaller, error) {
	contract, err := bindIAutomationRegistryMaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterCaller{contract: contract}, nil
}

func NewIAutomationRegistryMasterTransactor(address common.Address, transactor bind.ContractTransactor) (*IAutomationRegistryMasterTransactor, error) {
	contract, err := bindIAutomationRegistryMaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterTransactor{contract: contract}, nil
}

func NewIAutomationRegistryMasterFilterer(address common.Address, filterer bind.ContractFilterer) (*IAutomationRegistryMasterFilterer, error) {
	contract, err := bindIAutomationRegistryMaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterFilterer{contract: contract}, nil
}

func bindIAutomationRegistryMaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAutomationRegistryMasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationRegistryMaster.Contract.IAutomationRegistryMasterCaller.contract.Call(opts, result, method, params...)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.IAutomationRegistryMasterTransactor.contract.Transfer(opts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.IAutomationRegistryMasterTransactor.contract.Transact(opts, method, params...)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutomationRegistryMaster.Contract.contract.Call(opts, result, method, params...)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.contract.Transfer(opts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.contract.Transact(opts, method, params...)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "checkCallback", id, values, extraData)

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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckCallback(&_IAutomationRegistryMaster.CallOpts, id, values, extraData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckCallback(&_IAutomationRegistryMaster.CallOpts, id, values, extraData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "checkUpkeep", id, triggerData)

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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckUpkeep(&_IAutomationRegistryMaster.CallOpts, id, triggerData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) CheckUpkeep(id *big.Int, triggerData []byte) (CheckUpkeep,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckUpkeep(&_IAutomationRegistryMaster.CallOpts, id, triggerData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) CheckUpkeep0(opts *bind.CallOpts, id *big.Int) (CheckUpkeep0,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "checkUpkeep0", id)

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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckUpkeep0(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) CheckUpkeep0(id *big.Int) (CheckUpkeep0,

	error) {
	return _IAutomationRegistryMaster.Contract.CheckUpkeep0(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) FallbackTo() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.FallbackTo(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) FallbackTo() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.FallbackTo(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetActiveUpkeepIDs(&_IAutomationRegistryMaster.CallOpts, startIndex, maxCount)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetActiveUpkeepIDs(&_IAutomationRegistryMaster.CallOpts, startIndex, maxCount)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetAdminPrivilegeConfig(opts *bind.CallOpts, admin common.Address) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getAdminPrivilegeConfig", admin)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetAdminPrivilegeConfig(&_IAutomationRegistryMaster.CallOpts, admin)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetAdminPrivilegeConfig(admin common.Address) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetAdminPrivilegeConfig(&_IAutomationRegistryMaster.CallOpts, admin)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetAllowedReadOnlyAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getAllowedReadOnlyAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetAllowedReadOnlyAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetAllowedReadOnlyAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetAllowedReadOnlyAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetAutomationForwarderLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getAutomationForwarderLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetAutomationForwarderLogic(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetAutomationForwarderLogic() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetAutomationForwarderLogic(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetBalance(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetBalance(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getCancellationDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetCancellationDelay() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetCancellationDelay(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetCancellationDelay() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetCancellationDelay(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetChainModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getChainModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetChainModule() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetChainModule(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetChainModule() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetChainModule(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getConditionalGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetConditionalGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetConditionalGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetConditionalGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetFastGasFeedAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetFastGasFeedAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetForwarder(&_IAutomationRegistryMaster.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetForwarder(&_IAutomationRegistryMaster.CallOpts, upkeepID)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetLinkAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetLinkAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetLinkAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetLinkAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getLogGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetLogGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetLogGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetLogGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetLogGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMaxPaymentForGas(&_IAutomationRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMaxPaymentForGas(&_IAutomationRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getMinBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMinBalance(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetMinBalance(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMinBalance(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IAutomationRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster.CallOpts, peer)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IAutomationRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster.CallOpts, peer)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getPerPerformByteGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetPerPerformByteGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetPerPerformByteGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetPerPerformByteGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getPerSignerGasOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetPerSignerGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetPerSignerGasOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetPerSignerGasOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getReorgProtectionEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetReorgProtectionEnabled() (bool, error) {
	return _IAutomationRegistryMaster.Contract.GetReorgProtectionEnabled(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetReorgProtectionEnabled() (bool, error) {
	return _IAutomationRegistryMaster.Contract.GetReorgProtectionEnabled(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IAutomationRegistryMaster.Contract.GetSignerInfo(&_IAutomationRegistryMaster.CallOpts, query)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IAutomationRegistryMaster.Contract.GetSignerInfo(&_IAutomationRegistryMaster.CallOpts, query)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getState")

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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetState() (GetState,

	error) {
	return _IAutomationRegistryMaster.Contract.GetState(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetState() (GetState,

	error) {
	return _IAutomationRegistryMaster.Contract.GetState(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getTransmitCalldataFixedBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitCalldataFixedBytesOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetTransmitCalldataFixedBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitCalldataFixedBytesOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getTransmitCalldataPerSignerBytesOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetTransmitCalldataPerSignerBytesOverhead() (*big.Int, error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitCalldataPerSignerBytesOverhead(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitterInfo(&_IAutomationRegistryMaster.CallOpts, query)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IAutomationRegistryMaster.Contract.GetTransmitterInfo(&_IAutomationRegistryMaster.CallOpts, query)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationRegistryMaster.Contract.GetTriggerType(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IAutomationRegistryMaster.Contract.GetTriggerType(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeep(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetUpkeep(id *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeep(&_IAutomationRegistryMaster.CallOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IAutomationRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IAutomationRegistryMaster.CallOpts, upkeepId)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) HasDedupKey(opts *bind.CallOpts, dedupKey [32]byte) (bool, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "hasDedupKey", dedupKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationRegistryMaster.Contract.HasDedupKey(&_IAutomationRegistryMaster.CallOpts, dedupKey)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) HasDedupKey(dedupKey [32]byte) (bool, error) {
	return _IAutomationRegistryMaster.Contract.HasDedupKey(&_IAutomationRegistryMaster.CallOpts, dedupKey)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IAutomationRegistryMaster.Contract.LatestConfigDetails(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IAutomationRegistryMaster.Contract.LatestConfigDetails(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IAutomationRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IAutomationRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) Owner() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.Owner(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) Owner() (common.Address, error) {
	return _IAutomationRegistryMaster.Contract.Owner(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) SimulatePerformUpkeep(opts *bind.CallOpts, id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "simulatePerformUpkeep", id, performData)

	outstruct := new(SimulatePerformUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GasUsed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationRegistryMaster.Contract.SimulatePerformUpkeep(&_IAutomationRegistryMaster.CallOpts, id, performData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (SimulatePerformUpkeep,

	error) {
	return _IAutomationRegistryMaster.Contract.SimulatePerformUpkeep(&_IAutomationRegistryMaster.CallOpts, id, performData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) TypeAndVersion() (string, error) {
	return _IAutomationRegistryMaster.Contract.TypeAndVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) TypeAndVersion() (string, error) {
	return _IAutomationRegistryMaster.Contract.TypeAndVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IAutomationRegistryMaster.Contract.UpkeepTranscoderVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IAutomationRegistryMaster.Contract.UpkeepTranscoderVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IAutomationRegistryMaster.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) UpkeepVersion() (uint8, error) {
	return _IAutomationRegistryMaster.Contract.UpkeepVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterCallerSession) UpkeepVersion() (uint8, error) {
	return _IAutomationRegistryMaster.Contract.UpkeepVersion(&_IAutomationRegistryMaster.CallOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "acceptOwnership")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) AcceptOwnership() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptOwnership(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptOwnership(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptPayeeship(&_IAutomationRegistryMaster.TransactOpts, transmitter)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptPayeeship(&_IAutomationRegistryMaster.TransactOpts, transmitter)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptUpkeepAdmin(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AcceptUpkeepAdmin(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "addFunds", id, amount)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AddFunds(&_IAutomationRegistryMaster.TransactOpts, id, amount)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.AddFunds(&_IAutomationRegistryMaster.TransactOpts, id, amount)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.CancelUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.CancelUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.ExecuteCallback(&_IAutomationRegistryMaster.TransactOpts, id, payload)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.ExecuteCallback(&_IAutomationRegistryMaster.TransactOpts, id, payload)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.MigrateUpkeeps(&_IAutomationRegistryMaster.TransactOpts, ids, destination)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.MigrateUpkeeps(&_IAutomationRegistryMaster.TransactOpts, ids, destination)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.OnTokenTransfer(&_IAutomationRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.OnTokenTransfer(&_IAutomationRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "pause")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) Pause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Pause(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) Pause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Pause(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.PauseUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.PauseUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.ReceiveUpkeeps(&_IAutomationRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.ReceiveUpkeeps(&_IAutomationRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "recoverFunds")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) RecoverFunds() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RecoverFunds(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RecoverFunds(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RegisterUpkeep(&_IAutomationRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RegisterUpkeep(&_IAutomationRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RegisterUpkeep0(&_IAutomationRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.RegisterUpkeep0(&_IAutomationRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetAdminPrivilegeConfig(opts *bind.TransactOpts, admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setAdminPrivilegeConfig", admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetAdminPrivilegeConfig(&_IAutomationRegistryMaster.TransactOpts, admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetAdminPrivilegeConfig(admin common.Address, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetAdminPrivilegeConfig(&_IAutomationRegistryMaster.TransactOpts, admin, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetConfig(&_IAutomationRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetConfig(&_IAutomationRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetConfigTypeSafe(&_IAutomationRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetConfigTypeSafe(&_IAutomationRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setPayees", payees)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetPayees(&_IAutomationRegistryMaster.TransactOpts, payees)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetPayees(&_IAutomationRegistryMaster.TransactOpts, payees)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster.TransactOpts, peer, permission)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IAutomationRegistryMaster.TransactOpts, peer, permission)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepCheckData(&_IAutomationRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepCheckData(&_IAutomationRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepGasLimit(&_IAutomationRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepGasLimit(&_IAutomationRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IAutomationRegistryMaster.TransactOpts, id, config)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IAutomationRegistryMaster.TransactOpts, id, config)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IAutomationRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IAutomationRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IAutomationRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "transferOwnership", to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferOwnership(&_IAutomationRegistryMaster.TransactOpts, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferOwnership(&_IAutomationRegistryMaster.TransactOpts, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferPayeeship(&_IAutomationRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferPayeeship(&_IAutomationRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferUpkeepAdmin(&_IAutomationRegistryMaster.TransactOpts, id, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.TransferUpkeepAdmin(&_IAutomationRegistryMaster.TransactOpts, id, proposed)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Transmit(&_IAutomationRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Transmit(&_IAutomationRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "unpause")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) Unpause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Unpause(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) Unpause() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Unpause(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.UnpauseUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.UnpauseUpkeep(&_IAutomationRegistryMaster.TransactOpts, id)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawFunds(&_IAutomationRegistryMaster.TransactOpts, id, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawFunds(&_IAutomationRegistryMaster.TransactOpts, id, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawOwnerFunds(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawOwnerFunds(&_IAutomationRegistryMaster.TransactOpts)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawPayment(&_IAutomationRegistryMaster.TransactOpts, from, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.WithdrawPayment(&_IAutomationRegistryMaster.TransactOpts, from, to)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.contract.RawTransact(opts, calldata)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Fallback(&_IAutomationRegistryMaster.TransactOpts, calldata)
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IAutomationRegistryMaster.Contract.Fallback(&_IAutomationRegistryMaster.TransactOpts, calldata)
}

type IAutomationRegistryMasterAdminPrivilegeConfigSetIterator struct {
	Event *IAutomationRegistryMasterAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterAdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterAdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationRegistryMasterAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterAdminPrivilegeConfigSetIterator{contract: _IAutomationRegistryMaster.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterAdminPrivilegeConfigSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMasterAdminPrivilegeConfigSet, error) {
	event := new(IAutomationRegistryMasterAdminPrivilegeConfigSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterCancelledUpkeepReportIterator struct {
	Event *IAutomationRegistryMasterCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterCancelledUpkeepReportIterator{contract: _IAutomationRegistryMaster.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterCancelledUpkeepReport)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseCancelledUpkeepReport(log types.Log) (*IAutomationRegistryMasterCancelledUpkeepReport, error) {
	event := new(IAutomationRegistryMasterCancelledUpkeepReport)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterChainSpecificModuleUpdatedIterator struct {
	Event *IAutomationRegistryMasterChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMasterChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterChainSpecificModuleUpdatedIterator{contract: _IAutomationRegistryMaster.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterChainSpecificModuleUpdated)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*IAutomationRegistryMasterChainSpecificModuleUpdated, error) {
	event := new(IAutomationRegistryMasterChainSpecificModuleUpdated)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterConfigSetIterator struct {
	Event *IAutomationRegistryMasterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterConfigSet struct {
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

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterConfigSet(opts *bind.FilterOpts) (*IAutomationRegistryMasterConfigSetIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterConfigSetIterator{contract: _IAutomationRegistryMaster.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterConfigSet) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterConfigSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseConfigSet(log types.Log) (*IAutomationRegistryMasterConfigSet, error) {
	event := new(IAutomationRegistryMasterConfigSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterDedupKeyAddedIterator struct {
	Event *IAutomationRegistryMasterDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterDedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterDedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationRegistryMasterDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterDedupKeyAddedIterator{contract: _IAutomationRegistryMaster.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterDedupKeyAdded)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseDedupKeyAdded(log types.Log) (*IAutomationRegistryMasterDedupKeyAdded, error) {
	event := new(IAutomationRegistryMasterDedupKeyAdded)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterFundsAddedIterator struct {
	Event *IAutomationRegistryMasterFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterFundsAddedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationRegistryMasterFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterFundsAddedIterator{contract: _IAutomationRegistryMaster.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterFundsAdded)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseFundsAdded(log types.Log) (*IAutomationRegistryMasterFundsAdded, error) {
	event := new(IAutomationRegistryMasterFundsAdded)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterFundsWithdrawnIterator struct {
	Event *IAutomationRegistryMasterFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterFundsWithdrawnIterator{contract: _IAutomationRegistryMaster.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterFundsWithdrawn)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseFundsWithdrawn(log types.Log) (*IAutomationRegistryMasterFundsWithdrawn, error) {
	event := new(IAutomationRegistryMasterFundsWithdrawn)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator struct {
	Event *IAutomationRegistryMasterInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator{contract: _IAutomationRegistryMaster.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterInsufficientFundsUpkeepReport)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationRegistryMasterInsufficientFundsUpkeepReport, error) {
	event := new(IAutomationRegistryMasterInsufficientFundsUpkeepReport)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterOwnerFundsWithdrawnIterator struct {
	Event *IAutomationRegistryMasterOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterOwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterOwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IAutomationRegistryMasterOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterOwnerFundsWithdrawnIterator{contract: _IAutomationRegistryMaster.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterOwnerFundsWithdrawn)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*IAutomationRegistryMasterOwnerFundsWithdrawn, error) {
	event := new(IAutomationRegistryMasterOwnerFundsWithdrawn)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterOwnershipTransferRequestedIterator struct {
	Event *IAutomationRegistryMasterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMasterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterOwnershipTransferRequestedIterator{contract: _IAutomationRegistryMaster.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterOwnershipTransferRequested)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseOwnershipTransferRequested(log types.Log) (*IAutomationRegistryMasterOwnershipTransferRequested, error) {
	event := new(IAutomationRegistryMasterOwnershipTransferRequested)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterOwnershipTransferredIterator struct {
	Event *IAutomationRegistryMasterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMasterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterOwnershipTransferredIterator{contract: _IAutomationRegistryMaster.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterOwnershipTransferred)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseOwnershipTransferred(log types.Log) (*IAutomationRegistryMasterOwnershipTransferred, error) {
	event := new(IAutomationRegistryMasterOwnershipTransferred)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterPausedIterator struct {
	Event *IAutomationRegistryMasterPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterPaused(opts *bind.FilterOpts) (*IAutomationRegistryMasterPausedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterPausedIterator{contract: _IAutomationRegistryMaster.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterPaused)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParsePaused(log types.Log) (*IAutomationRegistryMasterPaused, error) {
	event := new(IAutomationRegistryMasterPaused)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterPayeesUpdatedIterator struct {
	Event *IAutomationRegistryMasterPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMasterPayeesUpdatedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterPayeesUpdatedIterator{contract: _IAutomationRegistryMaster.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterPayeesUpdated)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParsePayeesUpdated(log types.Log) (*IAutomationRegistryMasterPayeesUpdated, error) {
	event := new(IAutomationRegistryMasterPayeesUpdated)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterPayeeshipTransferRequestedIterator struct {
	Event *IAutomationRegistryMasterPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMasterPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterPayeeshipTransferRequestedIterator{contract: _IAutomationRegistryMaster.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterPayeeshipTransferRequested)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParsePayeeshipTransferRequested(log types.Log) (*IAutomationRegistryMasterPayeeshipTransferRequested, error) {
	event := new(IAutomationRegistryMasterPayeeshipTransferRequested)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterPayeeshipTransferredIterator struct {
	Event *IAutomationRegistryMasterPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMasterPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterPayeeshipTransferredIterator{contract: _IAutomationRegistryMaster.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterPayeeshipTransferred)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParsePayeeshipTransferred(log types.Log) (*IAutomationRegistryMasterPayeeshipTransferred, error) {
	event := new(IAutomationRegistryMasterPayeeshipTransferred)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterPaymentWithdrawnIterator struct {
	Event *IAutomationRegistryMasterPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationRegistryMasterPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterPaymentWithdrawnIterator{contract: _IAutomationRegistryMaster.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterPaymentWithdrawn)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParsePaymentWithdrawn(log types.Log) (*IAutomationRegistryMasterPaymentWithdrawn, error) {
	event := new(IAutomationRegistryMasterPaymentWithdrawn)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterReorgedUpkeepReportIterator struct {
	Event *IAutomationRegistryMasterReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterReorgedUpkeepReportIterator{contract: _IAutomationRegistryMaster.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterReorgedUpkeepReport)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseReorgedUpkeepReport(log types.Log) (*IAutomationRegistryMasterReorgedUpkeepReport, error) {
	event := new(IAutomationRegistryMasterReorgedUpkeepReport)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterStaleUpkeepReportIterator struct {
	Event *IAutomationRegistryMasterStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterStaleUpkeepReportIterator{contract: _IAutomationRegistryMaster.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterStaleUpkeepReport)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseStaleUpkeepReport(log types.Log) (*IAutomationRegistryMasterStaleUpkeepReport, error) {
	event := new(IAutomationRegistryMasterStaleUpkeepReport)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterTransmittedIterator struct {
	Event *IAutomationRegistryMasterTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterTransmittedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterTransmitted(opts *bind.FilterOpts) (*IAutomationRegistryMasterTransmittedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterTransmittedIterator{contract: _IAutomationRegistryMaster.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterTransmitted) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterTransmitted)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseTransmitted(log types.Log) (*IAutomationRegistryMasterTransmitted, error) {
	event := new(IAutomationRegistryMasterTransmitted)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUnpausedIterator struct {
	Event *IAutomationRegistryMasterUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IAutomationRegistryMasterUnpausedIterator, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUnpausedIterator{contract: _IAutomationRegistryMaster.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUnpaused) (event.Subscription, error) {

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUnpaused)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUnpaused(log types.Log) (*IAutomationRegistryMasterUnpaused, error) {
	event := new(IAutomationRegistryMasterUnpaused)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator struct {
	Event *IAutomationRegistryMasterUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepAdminTransferRequested)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationRegistryMasterUpkeepAdminTransferRequested, error) {
	event := new(IAutomationRegistryMasterUpkeepAdminTransferRequested)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepAdminTransferredIterator struct {
	Event *IAutomationRegistryMasterUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMasterUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepAdminTransferredIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepAdminTransferred)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepAdminTransferred(log types.Log) (*IAutomationRegistryMasterUpkeepAdminTransferred, error) {
	event := new(IAutomationRegistryMasterUpkeepAdminTransferred)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepCanceledIterator struct {
	Event *IAutomationRegistryMasterUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationRegistryMasterUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepCanceledIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepCanceled)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepCanceled(log types.Log) (*IAutomationRegistryMasterUpkeepCanceled, error) {
	event := new(IAutomationRegistryMasterUpkeepCanceled)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepCheckDataSetIterator struct {
	Event *IAutomationRegistryMasterUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepCheckDataSetIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepCheckDataSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepCheckDataSet(log types.Log) (*IAutomationRegistryMasterUpkeepCheckDataSet, error) {
	event := new(IAutomationRegistryMasterUpkeepCheckDataSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepGasLimitSetIterator struct {
	Event *IAutomationRegistryMasterUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepGasLimitSetIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepGasLimitSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepGasLimitSet(log types.Log) (*IAutomationRegistryMasterUpkeepGasLimitSet, error) {
	event := new(IAutomationRegistryMasterUpkeepGasLimitSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepMigratedIterator struct {
	Event *IAutomationRegistryMasterUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepMigratedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepMigrated)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepMigrated(log types.Log) (*IAutomationRegistryMasterUpkeepMigrated, error) {
	event := new(IAutomationRegistryMasterUpkeepMigrated)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepOffchainConfigSetIterator struct {
	Event *IAutomationRegistryMasterUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepOffchainConfigSetIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepOffchainConfigSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepOffchainConfigSet, error) {
	event := new(IAutomationRegistryMasterUpkeepOffchainConfigSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepPausedIterator struct {
	Event *IAutomationRegistryMasterUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepPausedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepPaused)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepPaused(log types.Log) (*IAutomationRegistryMasterUpkeepPaused, error) {
	event := new(IAutomationRegistryMasterUpkeepPaused)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepPerformedIterator struct {
	Event *IAutomationRegistryMasterUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationRegistryMasterUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepPerformedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepPerformed)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepPerformed(log types.Log) (*IAutomationRegistryMasterUpkeepPerformed, error) {
	event := new(IAutomationRegistryMasterUpkeepPerformed)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator struct {
	Event *IAutomationRegistryMasterUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepPrivilegeConfigSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepPrivilegeConfigSet, error) {
	event := new(IAutomationRegistryMasterUpkeepPrivilegeConfigSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepReceivedIterator struct {
	Event *IAutomationRegistryMasterUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepReceivedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepReceived)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepReceived(log types.Log) (*IAutomationRegistryMasterUpkeepReceived, error) {
	event := new(IAutomationRegistryMasterUpkeepReceived)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepRegisteredIterator struct {
	Event *IAutomationRegistryMasterUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepRegisteredIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepRegistered)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepRegistered(log types.Log) (*IAutomationRegistryMasterUpkeepRegistered, error) {
	event := new(IAutomationRegistryMasterUpkeepRegistered)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepTriggerConfigSetIterator struct {
	Event *IAutomationRegistryMasterUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepTriggerConfigSetIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepTriggerConfigSet)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepTriggerConfigSet, error) {
	event := new(IAutomationRegistryMasterUpkeepTriggerConfigSet)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IAutomationRegistryMasterUpkeepUnpausedIterator struct {
	Event *IAutomationRegistryMasterUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IAutomationRegistryMasterUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAutomationRegistryMasterUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IAutomationRegistryMasterUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IAutomationRegistryMasterUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IAutomationRegistryMasterUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IAutomationRegistryMasterUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IAutomationRegistryMasterUpkeepUnpausedIterator{contract: _IAutomationRegistryMaster.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IAutomationRegistryMaster.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IAutomationRegistryMasterUpkeepUnpaused)
				if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IAutomationRegistryMaster *IAutomationRegistryMasterFilterer) ParseUpkeepUnpaused(log types.Log) (*IAutomationRegistryMasterUpkeepUnpaused, error) {
	event := new(IAutomationRegistryMasterUpkeepUnpaused)
	if err := _IAutomationRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_IAutomationRegistryMaster *IAutomationRegistryMaster) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IAutomationRegistryMaster.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _IAutomationRegistryMaster.ParseAdminPrivilegeConfigSet(log)
	case _IAutomationRegistryMaster.abi.Events["CancelledUpkeepReport"].ID:
		return _IAutomationRegistryMaster.ParseCancelledUpkeepReport(log)
	case _IAutomationRegistryMaster.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _IAutomationRegistryMaster.ParseChainSpecificModuleUpdated(log)
	case _IAutomationRegistryMaster.abi.Events["ConfigSet"].ID:
		return _IAutomationRegistryMaster.ParseConfigSet(log)
	case _IAutomationRegistryMaster.abi.Events["DedupKeyAdded"].ID:
		return _IAutomationRegistryMaster.ParseDedupKeyAdded(log)
	case _IAutomationRegistryMaster.abi.Events["FundsAdded"].ID:
		return _IAutomationRegistryMaster.ParseFundsAdded(log)
	case _IAutomationRegistryMaster.abi.Events["FundsWithdrawn"].ID:
		return _IAutomationRegistryMaster.ParseFundsWithdrawn(log)
	case _IAutomationRegistryMaster.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IAutomationRegistryMaster.ParseInsufficientFundsUpkeepReport(log)
	case _IAutomationRegistryMaster.abi.Events["OwnerFundsWithdrawn"].ID:
		return _IAutomationRegistryMaster.ParseOwnerFundsWithdrawn(log)
	case _IAutomationRegistryMaster.abi.Events["OwnershipTransferRequested"].ID:
		return _IAutomationRegistryMaster.ParseOwnershipTransferRequested(log)
	case _IAutomationRegistryMaster.abi.Events["OwnershipTransferred"].ID:
		return _IAutomationRegistryMaster.ParseOwnershipTransferred(log)
	case _IAutomationRegistryMaster.abi.Events["Paused"].ID:
		return _IAutomationRegistryMaster.ParsePaused(log)
	case _IAutomationRegistryMaster.abi.Events["PayeesUpdated"].ID:
		return _IAutomationRegistryMaster.ParsePayeesUpdated(log)
	case _IAutomationRegistryMaster.abi.Events["PayeeshipTransferRequested"].ID:
		return _IAutomationRegistryMaster.ParsePayeeshipTransferRequested(log)
	case _IAutomationRegistryMaster.abi.Events["PayeeshipTransferred"].ID:
		return _IAutomationRegistryMaster.ParsePayeeshipTransferred(log)
	case _IAutomationRegistryMaster.abi.Events["PaymentWithdrawn"].ID:
		return _IAutomationRegistryMaster.ParsePaymentWithdrawn(log)
	case _IAutomationRegistryMaster.abi.Events["ReorgedUpkeepReport"].ID:
		return _IAutomationRegistryMaster.ParseReorgedUpkeepReport(log)
	case _IAutomationRegistryMaster.abi.Events["StaleUpkeepReport"].ID:
		return _IAutomationRegistryMaster.ParseStaleUpkeepReport(log)
	case _IAutomationRegistryMaster.abi.Events["Transmitted"].ID:
		return _IAutomationRegistryMaster.ParseTransmitted(log)
	case _IAutomationRegistryMaster.abi.Events["Unpaused"].ID:
		return _IAutomationRegistryMaster.ParseUnpaused(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepAdminTransferRequested(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepAdminTransferred"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepAdminTransferred(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepCanceled"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepCanceled(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepCheckDataSet"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepCheckDataSet(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepGasLimitSet"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepGasLimitSet(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepMigrated"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepMigrated(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepOffchainConfigSet(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepPaused"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepPaused(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepPerformed"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepPerformed(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepPrivilegeConfigSet(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepReceived"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepReceived(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepRegistered"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepRegistered(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepTriggerConfigSet(log)
	case _IAutomationRegistryMaster.abi.Events["UpkeepUnpaused"].ID:
		return _IAutomationRegistryMaster.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IAutomationRegistryMasterAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (IAutomationRegistryMasterCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IAutomationRegistryMasterChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (IAutomationRegistryMasterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IAutomationRegistryMasterDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (IAutomationRegistryMasterFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IAutomationRegistryMasterFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IAutomationRegistryMasterInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IAutomationRegistryMasterOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (IAutomationRegistryMasterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IAutomationRegistryMasterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IAutomationRegistryMasterPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IAutomationRegistryMasterPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IAutomationRegistryMasterPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IAutomationRegistryMasterPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IAutomationRegistryMasterPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IAutomationRegistryMasterReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IAutomationRegistryMasterStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IAutomationRegistryMasterTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IAutomationRegistryMasterUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IAutomationRegistryMasterUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IAutomationRegistryMasterUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IAutomationRegistryMasterUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IAutomationRegistryMasterUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IAutomationRegistryMasterUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IAutomationRegistryMasterUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IAutomationRegistryMasterUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IAutomationRegistryMasterUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IAutomationRegistryMasterUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IAutomationRegistryMasterUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (IAutomationRegistryMasterUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IAutomationRegistryMasterUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IAutomationRegistryMasterUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IAutomationRegistryMasterUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IAutomationRegistryMaster *IAutomationRegistryMaster) Address() common.Address {
	return _IAutomationRegistryMaster.address
}

type IAutomationRegistryMasterInterface interface {
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

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetCancellationDelay(opts *bind.CallOpts) (*big.Int, error)

	GetChainModule(opts *bind.CallOpts) (common.Address, error)

	GetConditionalGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error)

	GetMinBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetPerPerformByteGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetPerSignerGasOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetReorgProtectionEnabled(opts *bind.CallOpts) (bool, error)

	GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

		error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetTransmitCalldataFixedBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

	GetTransmitCalldataPerSignerBytesOverhead(opts *bind.CallOpts) (*big.Int, error)

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

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*IAutomationRegistryMasterAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMasterAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IAutomationRegistryMasterCancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMasterChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*IAutomationRegistryMasterChainSpecificModuleUpdated, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IAutomationRegistryMasterConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IAutomationRegistryMasterConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*IAutomationRegistryMasterDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*IAutomationRegistryMasterDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IAutomationRegistryMasterFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IAutomationRegistryMasterFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IAutomationRegistryMasterFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IAutomationRegistryMasterInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IAutomationRegistryMasterOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*IAutomationRegistryMasterOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMasterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IAutomationRegistryMasterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IAutomationRegistryMasterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IAutomationRegistryMasterOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IAutomationRegistryMasterPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IAutomationRegistryMasterPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IAutomationRegistryMasterPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IAutomationRegistryMasterPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMasterPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IAutomationRegistryMasterPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IAutomationRegistryMasterPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IAutomationRegistryMasterPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IAutomationRegistryMasterPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IAutomationRegistryMasterPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IAutomationRegistryMasterReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IAutomationRegistryMasterStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IAutomationRegistryMasterTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IAutomationRegistryMasterTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IAutomationRegistryMasterUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IAutomationRegistryMasterUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMasterUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IAutomationRegistryMasterUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IAutomationRegistryMasterUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IAutomationRegistryMasterUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IAutomationRegistryMasterUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IAutomationRegistryMasterUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IAutomationRegistryMasterUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IAutomationRegistryMasterUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IAutomationRegistryMasterUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IAutomationRegistryMasterUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IAutomationRegistryMasterUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IAutomationRegistryMasterUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IAutomationRegistryMasterUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IAutomationRegistryMasterUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IAutomationRegistryMasterUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IAutomationRegistryMasterUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IAutomationRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IAutomationRegistryMasterUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
