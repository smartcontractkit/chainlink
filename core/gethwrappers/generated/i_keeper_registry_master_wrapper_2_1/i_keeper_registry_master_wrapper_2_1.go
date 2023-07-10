// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_keeper_registry_master_wrapper_2_1

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

type KeeperRegistryBase21ConditionalTriggerConfig struct {
	CheckCadance uint32
}

type KeeperRegistryBase21LogTriggerConfig struct {
	ContractAddress common.Address
	FilterSelector  uint8
	Topic0          [32]byte
	Topic1          [32]byte
	Topic2          [32]byte
	Topic3          [32]byte
}

type KeeperRegistryBase21OnchainConfig struct {
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

type KeeperRegistryBase21State struct {
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

type KeeperRegistryBase21UpkeepInfo struct {
	Target                   common.Address
	Forwarder                common.Address
	ExecuteGas               uint32
	CheckData                []byte
	Balance                  *big.Int
	Admin                    common.Address
	MaxValidBlocknumber      uint64
	LastPerformedBlockNumber uint32
	AmountSpent              *big.Int
	Paused                   bool
	OffchainConfig           []byte
}

var IKeeperRegistryMasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"trigger\",\"type\":\"uint8\"}],\"name\":\"getActiveUpkeepIDsByType\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getConditionalTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkCadance\",\"type\":\"uint32\"}],\"internalType\":\"structKeeperRegistryBase2_1.ConditionalTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structKeeperRegistryBase2_1.LogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structKeeperRegistryBase2_1.State\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structKeeperRegistryBase2_1.OnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var IKeeperRegistryMasterABI = IKeeperRegistryMasterMetaData.ABI

type IKeeperRegistryMaster struct {
	address common.Address
	abi     abi.ABI
	IKeeperRegistryMasterCaller
	IKeeperRegistryMasterTransactor
	IKeeperRegistryMasterFilterer
}

type IKeeperRegistryMasterCaller struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterTransactor struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterFilterer struct {
	contract *bind.BoundContract
}

type IKeeperRegistryMasterSession struct {
	Contract     *IKeeperRegistryMaster
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IKeeperRegistryMasterCallerSession struct {
	Contract *IKeeperRegistryMasterCaller
	CallOpts bind.CallOpts
}

type IKeeperRegistryMasterTransactorSession struct {
	Contract     *IKeeperRegistryMasterTransactor
	TransactOpts bind.TransactOpts
}

type IKeeperRegistryMasterRaw struct {
	Contract *IKeeperRegistryMaster
}

type IKeeperRegistryMasterCallerRaw struct {
	Contract *IKeeperRegistryMasterCaller
}

type IKeeperRegistryMasterTransactorRaw struct {
	Contract *IKeeperRegistryMasterTransactor
}

func NewIKeeperRegistryMaster(address common.Address, backend bind.ContractBackend) (*IKeeperRegistryMaster, error) {
	abi, err := abi.JSON(strings.NewReader(IKeeperRegistryMasterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIKeeperRegistryMaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMaster{address: address, abi: abi, IKeeperRegistryMasterCaller: IKeeperRegistryMasterCaller{contract: contract}, IKeeperRegistryMasterTransactor: IKeeperRegistryMasterTransactor{contract: contract}, IKeeperRegistryMasterFilterer: IKeeperRegistryMasterFilterer{contract: contract}}, nil
}

func NewIKeeperRegistryMasterCaller(address common.Address, caller bind.ContractCaller) (*IKeeperRegistryMasterCaller, error) {
	contract, err := bindIKeeperRegistryMaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterCaller{contract: contract}, nil
}

func NewIKeeperRegistryMasterTransactor(address common.Address, transactor bind.ContractTransactor) (*IKeeperRegistryMasterTransactor, error) {
	contract, err := bindIKeeperRegistryMaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterTransactor{contract: contract}, nil
}

func NewIKeeperRegistryMasterFilterer(address common.Address, filterer bind.ContractFilterer) (*IKeeperRegistryMasterFilterer, error) {
	contract, err := bindIKeeperRegistryMaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFilterer{contract: contract}, nil
}

func bindIKeeperRegistryMaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IKeeperRegistryMasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterCaller.contract.Call(opts, result, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterTransactor.contract.Transfer(opts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.IKeeperRegistryMasterTransactor.contract.Transact(opts, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKeeperRegistryMaster.Contract.contract.Call(opts, result, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.contract.Transfer(opts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.contract.Transact(opts, method, params...)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) FallbackTo() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.FallbackTo(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) FallbackTo() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.FallbackTo(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, endIndex)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDs(&_IKeeperRegistryMaster.CallOpts, startIndex, endIndex)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDs(&_IKeeperRegistryMaster.CallOpts, startIndex, endIndex)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getActiveUpkeepIDsByType", startIndex, endIndex, trigger)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDsByType(&_IKeeperRegistryMaster.CallOpts, startIndex, endIndex, trigger)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetActiveUpkeepIDsByType(&_IKeeperRegistryMaster.CallOpts, startIndex, endIndex, trigger)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetConditionalTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getConditionalTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21ConditionalTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21ConditionalTriggerConfig)).(*KeeperRegistryBase21ConditionalTriggerConfig)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetConditionalTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	return _IKeeperRegistryMaster.Contract.GetConditionalTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetConditionalTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	return _IKeeperRegistryMaster.Contract.GetConditionalTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetFastGasFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetFastGasFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLinkAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLinkAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.GetLinkNativeFeedAddress(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getLogTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21LogTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21LogTriggerConfig)).(*KeeperRegistryBase21LogTriggerConfig)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetLogTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	return _IKeeperRegistryMaster.Contract.GetLogTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetLogTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	return _IKeeperRegistryMaster.Contract.GetLogTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMaxPaymentForGas(&_IKeeperRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMaxPaymentForGas(&_IKeeperRegistryMaster.CallOpts, triggerType, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _IKeeperRegistryMaster.Contract.GetMinBalanceForUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetMode() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetMode(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetMode() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetMode(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.CallOpts, peer)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.CallOpts, peer)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetSignerInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetSignerInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(KeeperRegistryBase21State)).(*KeeperRegistryBase21State)
	outstruct.Config = *abi.ConvertType(out[1], new(KeeperRegistryBase21OnchainConfig)).(*KeeperRegistryBase21OnchainConfig)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetState() (GetState,

	error) {
	return _IKeeperRegistryMaster.Contract.GetState(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetState() (GetState,

	error) {
	return _IKeeperRegistryMaster.Contract.GetState(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetTransmitterInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _IKeeperRegistryMaster.Contract.GetTransmitterInfo(&_IKeeperRegistryMaster.CallOpts, query)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetTriggerType(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _IKeeperRegistryMaster.Contract.GetTriggerType(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeep(id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeep(id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeep(&_IKeeperRegistryMaster.CallOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _IKeeperRegistryMaster.Contract.GetUpkeepTriggerConfig(&_IKeeperRegistryMaster.CallOpts, upkeepId)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDetails(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDetails(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _IKeeperRegistryMaster.Contract.LatestConfigDigestAndEpoch(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Owner() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.Owner(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) Owner() (common.Address, error) {
	return _IKeeperRegistryMaster.Contract.Owner(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TypeAndVersion() (string, error) {
	return _IKeeperRegistryMaster.Contract.TypeAndVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) TypeAndVersion() (string, error) {
	return _IKeeperRegistryMaster.Contract.TypeAndVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepTranscoderVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepTranscoderVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IKeeperRegistryMaster.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UpkeepVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterCallerSession) UpkeepVersion() (uint8, error) {
	return _IKeeperRegistryMaster.Contract.UpkeepVersion(&_IKeeperRegistryMaster.CallOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptOwnership")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptOwnership() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptOwnership(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptOwnership(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AcceptUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "addFunds", id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AddFunds(&_IKeeperRegistryMaster.TransactOpts, id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.AddFunds(&_IKeeperRegistryMaster.TransactOpts, id, amount)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "cancelUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CancelUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CancelUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckCallback(&_IKeeperRegistryMaster.TransactOpts, id, values, extraData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckCallback(&_IKeeperRegistryMaster.TransactOpts, id, values, extraData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "checkUpkeep", id, checkData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep(&_IKeeperRegistryMaster.TransactOpts, id, checkData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep(&_IKeeperRegistryMaster.TransactOpts, id, checkData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "checkUpkeep0", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep0(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.CheckUpkeep0(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "executeCallback", id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ExecuteCallback(&_IKeeperRegistryMaster.TransactOpts, id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ExecuteCallback(&_IKeeperRegistryMaster.TransactOpts, id, payload)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.MigrateUpkeeps(&_IKeeperRegistryMaster.TransactOpts, ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.MigrateUpkeeps(&_IKeeperRegistryMaster.TransactOpts, ids, destination)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.OnTokenTransfer(&_IKeeperRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.OnTokenTransfer(&_IKeeperRegistryMaster.TransactOpts, sender, amount, data)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "pause")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Pause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Pause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Pause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Pause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "pauseUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.PauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.PauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ReceiveUpkeeps(&_IKeeperRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.ReceiveUpkeeps(&_IKeeperRegistryMaster.TransactOpts, encodedUpkeeps)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "recoverFunds")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RecoverFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RecoverFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RecoverFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep0(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.RegisterUpkeep0(&_IKeeperRegistryMaster.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfig(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetConfig(&_IKeeperRegistryMaster.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setPayees", payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPayees(&_IKeeperRegistryMaster.TransactOpts, payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPayees(&_IKeeperRegistryMaster.TransactOpts, payees)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.TransactOpts, peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetPeerRegistryMigrationPermission(&_IKeeperRegistryMaster.TransactOpts, peer, permission)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepCheckData(&_IKeeperRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepCheckData(&_IKeeperRegistryMaster.TransactOpts, id, newCheckData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepGasLimit(&_IKeeperRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepGasLimit(&_IKeeperRegistryMaster.TransactOpts, id, gasLimit)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IKeeperRegistryMaster.TransactOpts, id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepOffchainConfig(&_IKeeperRegistryMaster.TransactOpts, id, config)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepPrivilegeConfig(&_IKeeperRegistryMaster.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IKeeperRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SetUpkeepTriggerConfig(&_IKeeperRegistryMaster.TransactOpts, id, triggerConfig)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SimulatePerformUpkeep(&_IKeeperRegistryMaster.TransactOpts, id, performData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.SimulatePerformUpkeep(&_IKeeperRegistryMaster.TransactOpts, id, performData)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferOwnership", to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferOwnership(&_IKeeperRegistryMaster.TransactOpts, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferOwnership(&_IKeeperRegistryMaster.TransactOpts, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferPayeeship(&_IKeeperRegistryMaster.TransactOpts, transmitter, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.TransferUpkeepAdmin(&_IKeeperRegistryMaster.TransactOpts, id, proposed)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Transmit(&_IKeeperRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Transmit(&_IKeeperRegistryMaster.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "unpause")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Unpause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Unpause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Unpause() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Unpause(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.UnpauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.UnpauseUpkeep(&_IKeeperRegistryMaster.TransactOpts, id)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawFunds(&_IKeeperRegistryMaster.TransactOpts, id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawFunds(&_IKeeperRegistryMaster.TransactOpts, id, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawOwnerFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawOwnerFunds(&_IKeeperRegistryMaster.TransactOpts)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawPayment(&_IKeeperRegistryMaster.TransactOpts, from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.WithdrawPayment(&_IKeeperRegistryMaster.TransactOpts, from, to)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.contract.RawTransact(opts, calldata)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Fallback(&_IKeeperRegistryMaster.TransactOpts, calldata)
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _IKeeperRegistryMaster.Contract.Fallback(&_IKeeperRegistryMaster.TransactOpts, calldata)
}

type IKeeperRegistryMasterCancelledUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterCancelledUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterCancelledUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseCancelledUpkeepReport(log types.Log) (*IKeeperRegistryMasterCancelledUpkeepReport, error) {
	event := new(IKeeperRegistryMasterCancelledUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterConfigSetIterator struct {
	Event *IKeeperRegistryMasterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterConfigSet struct {
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

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterConfigSet(opts *bind.FilterOpts) (*IKeeperRegistryMasterConfigSetIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterConfigSet) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseConfigSet(log types.Log) (*IKeeperRegistryMasterConfigSet, error) {
	event := new(IKeeperRegistryMasterConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterFundsAddedIterator struct {
	Event *IKeeperRegistryMasterFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterFundsAddedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IKeeperRegistryMasterFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFundsAddedIterator{contract: _IKeeperRegistryMaster.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterFundsAdded)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseFundsAdded(log types.Log) (*IKeeperRegistryMasterFundsAdded, error) {
	event := new(IKeeperRegistryMasterFundsAdded)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterFundsWithdrawnIterator struct {
	Event *IKeeperRegistryMasterFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterFundsWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterFundsWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterFundsWithdrawn, error) {
	event := new(IKeeperRegistryMasterFundsWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*IKeeperRegistryMasterInsufficientFundsUpkeepReport, error) {
	event := new(IKeeperRegistryMasterInsufficientFundsUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnerFundsWithdrawnIterator struct {
	Event *IKeeperRegistryMasterOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterOwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IKeeperRegistryMasterOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnerFundsWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnerFundsWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterOwnerFundsWithdrawn, error) {
	event := new(IKeeperRegistryMasterOwnerFundsWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnershipTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnershipTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnershipTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnershipTransferRequested(log types.Log) (*IKeeperRegistryMasterOwnershipTransferRequested, error) {
	event := new(IKeeperRegistryMasterOwnershipTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterOwnershipTransferredIterator struct {
	Event *IKeeperRegistryMasterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterOwnershipTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterOwnershipTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseOwnershipTransferred(log types.Log) (*IKeeperRegistryMasterOwnershipTransferred, error) {
	event := new(IKeeperRegistryMasterOwnershipTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPausedIterator struct {
	Event *IKeeperRegistryMasterPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterPausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterPausedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPausedIterator{contract: _IKeeperRegistryMaster.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaused) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePaused(log types.Log) (*IKeeperRegistryMasterPaused, error) {
	event := new(IKeeperRegistryMasterPaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeesUpdatedIterator struct {
	Event *IKeeperRegistryMasterPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*IKeeperRegistryMasterPayeesUpdatedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeesUpdatedIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeesUpdated)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeesUpdated(log types.Log) (*IKeeperRegistryMasterPayeesUpdated, error) {
	event := new(IKeeperRegistryMasterPayeesUpdated)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeeshipTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeeshipTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeeshipTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeeshipTransferRequested(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferRequested, error) {
	event := new(IKeeperRegistryMasterPayeeshipTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPayeeshipTransferredIterator struct {
	Event *IKeeperRegistryMasterPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPayeeshipTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPayeeshipTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePayeeshipTransferred(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferred, error) {
	event := new(IKeeperRegistryMasterPayeeshipTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterPaymentWithdrawnIterator struct {
	Event *IKeeperRegistryMasterPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IKeeperRegistryMasterPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterPaymentWithdrawnIterator{contract: _IKeeperRegistryMaster.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterPaymentWithdrawn)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParsePaymentWithdrawn(log types.Log) (*IKeeperRegistryMasterPaymentWithdrawn, error) {
	event := new(IKeeperRegistryMasterPaymentWithdrawn)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterReorgedUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterReorgedUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterReorgedUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseReorgedUpkeepReport(log types.Log) (*IKeeperRegistryMasterReorgedUpkeepReport, error) {
	event := new(IKeeperRegistryMasterReorgedUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterStaleUpkeepReportIterator struct {
	Event *IKeeperRegistryMasterStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterStaleUpkeepReportIterator{contract: _IKeeperRegistryMaster.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterStaleUpkeepReport)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseStaleUpkeepReport(log types.Log) (*IKeeperRegistryMasterStaleUpkeepReport, error) {
	event := new(IKeeperRegistryMasterStaleUpkeepReport)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterTransmittedIterator struct {
	Event *IKeeperRegistryMasterTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterTransmittedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterTransmitted(opts *bind.FilterOpts) (*IKeeperRegistryMasterTransmittedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterTransmittedIterator{contract: _IKeeperRegistryMaster.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterTransmitted) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterTransmitted)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseTransmitted(log types.Log) (*IKeeperRegistryMasterTransmitted, error) {
	event := new(IKeeperRegistryMasterTransmitted)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUnpausedIterator struct {
	Event *IKeeperRegistryMasterUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUnpausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterUnpausedIterator, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUnpausedIterator{contract: _IKeeperRegistryMaster.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUnpaused) (event.Subscription, error) {

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUnpaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUnpaused(log types.Log) (*IKeeperRegistryMasterUnpaused, error) {
	event := new(IKeeperRegistryMasterUnpaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator struct {
	Event *IKeeperRegistryMasterUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferRequested, error) {
	event := new(IKeeperRegistryMasterUpkeepAdminTransferRequested)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepAdminTransferredIterator struct {
	Event *IKeeperRegistryMasterUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepAdminTransferredIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepAdminTransferred)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepAdminTransferred(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferred, error) {
	event := new(IKeeperRegistryMasterUpkeepAdminTransferred)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepCanceledIterator struct {
	Event *IKeeperRegistryMasterUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IKeeperRegistryMasterUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepCanceledIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepCanceled)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepCanceled(log types.Log) (*IKeeperRegistryMasterUpkeepCanceled, error) {
	event := new(IKeeperRegistryMasterUpkeepCanceled)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepCheckDataSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepCheckDataSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepCheckDataSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepCheckDataSet(log types.Log) (*IKeeperRegistryMasterUpkeepCheckDataSet, error) {
	event := new(IKeeperRegistryMasterUpkeepCheckDataSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepGasLimitSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepGasLimitSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepGasLimitSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepGasLimitSet(log types.Log) (*IKeeperRegistryMasterUpkeepGasLimitSet, error) {
	event := new(IKeeperRegistryMasterUpkeepGasLimitSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepMigratedIterator struct {
	Event *IKeeperRegistryMasterUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepMigratedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepMigrated)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepMigrated(log types.Log) (*IKeeperRegistryMasterUpkeepMigrated, error) {
	event := new(IKeeperRegistryMasterUpkeepMigrated)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepOffchainConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepOffchainConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepOffchainConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepOffchainConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPausedIterator struct {
	Event *IKeeperRegistryMasterUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPausedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPaused(log types.Log) (*IKeeperRegistryMasterUpkeepPaused, error) {
	event := new(IKeeperRegistryMasterUpkeepPaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPerformedIterator struct {
	Event *IKeeperRegistryMasterUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IKeeperRegistryMasterUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPerformedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPerformed)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPerformed(log types.Log) (*IKeeperRegistryMasterUpkeepPerformed, error) {
	event := new(IKeeperRegistryMasterUpkeepPerformed)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepPrivilegeConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepReceivedIterator struct {
	Event *IKeeperRegistryMasterUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepReceivedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepReceived)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepReceived(log types.Log) (*IKeeperRegistryMasterUpkeepReceived, error) {
	event := new(IKeeperRegistryMasterUpkeepReceived)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepRegisteredIterator struct {
	Event *IKeeperRegistryMasterUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepRegisteredIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepRegistered)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepRegistered(log types.Log) (*IKeeperRegistryMasterUpkeepRegistered, error) {
	event := new(IKeeperRegistryMasterUpkeepRegistered)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepTriggerConfigSetIterator struct {
	Event *IKeeperRegistryMasterUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepTriggerConfigSetIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepTriggerConfigSet, error) {
	event := new(IKeeperRegistryMasterUpkeepTriggerConfigSet)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type IKeeperRegistryMasterUpkeepUnpausedIterator struct {
	Event *IKeeperRegistryMasterUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IKeeperRegistryMasterUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(IKeeperRegistryMasterUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *IKeeperRegistryMasterUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type IKeeperRegistryMasterUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &IKeeperRegistryMasterUpkeepUnpausedIterator{contract: _IKeeperRegistryMaster.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IKeeperRegistryMaster.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(IKeeperRegistryMasterUpkeepUnpaused)
				if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_IKeeperRegistryMaster *IKeeperRegistryMasterFilterer) ParseUpkeepUnpaused(log types.Log) (*IKeeperRegistryMasterUpkeepUnpaused, error) {
	event := new(IKeeperRegistryMasterUpkeepUnpaused)
	if err := _IKeeperRegistryMaster.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
	State        KeeperRegistryBase21State
	Config       KeeperRegistryBase21OnchainConfig
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

func (_IKeeperRegistryMaster *IKeeperRegistryMaster) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _IKeeperRegistryMaster.abi.Events["CancelledUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseCancelledUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["ConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["FundsAdded"].ID:
		return _IKeeperRegistryMaster.ParseFundsAdded(log)
	case _IKeeperRegistryMaster.abi.Events["FundsWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParseFundsWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseInsufficientFundsUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["OwnerFundsWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParseOwnerFundsWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["OwnershipTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParseOwnershipTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["OwnershipTransferred"].ID:
		return _IKeeperRegistryMaster.ParseOwnershipTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["Paused"].ID:
		return _IKeeperRegistryMaster.ParsePaused(log)
	case _IKeeperRegistryMaster.abi.Events["PayeesUpdated"].ID:
		return _IKeeperRegistryMaster.ParsePayeesUpdated(log)
	case _IKeeperRegistryMaster.abi.Events["PayeeshipTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParsePayeeshipTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["PayeeshipTransferred"].ID:
		return _IKeeperRegistryMaster.ParsePayeeshipTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["PaymentWithdrawn"].ID:
		return _IKeeperRegistryMaster.ParsePaymentWithdrawn(log)
	case _IKeeperRegistryMaster.abi.Events["ReorgedUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseReorgedUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["StaleUpkeepReport"].ID:
		return _IKeeperRegistryMaster.ParseStaleUpkeepReport(log)
	case _IKeeperRegistryMaster.abi.Events["Transmitted"].ID:
		return _IKeeperRegistryMaster.ParseTransmitted(log)
	case _IKeeperRegistryMaster.abi.Events["Unpaused"].ID:
		return _IKeeperRegistryMaster.ParseUnpaused(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepAdminTransferRequested(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepAdminTransferred"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepAdminTransferred(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepCanceled"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepCanceled(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepCheckDataSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepCheckDataSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepGasLimitSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepGasLimitSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepMigrated"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepMigrated(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepOffchainConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPaused"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPaused(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPerformed"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPerformed(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepPrivilegeConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepReceived"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepReceived(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepRegistered"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepRegistered(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepTriggerConfigSet(log)
	case _IKeeperRegistryMaster.abi.Events["UpkeepUnpaused"].ID:
		return _IKeeperRegistryMaster.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (IKeeperRegistryMasterCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (IKeeperRegistryMasterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (IKeeperRegistryMasterFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (IKeeperRegistryMasterFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (IKeeperRegistryMasterInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (IKeeperRegistryMasterOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (IKeeperRegistryMasterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (IKeeperRegistryMasterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (IKeeperRegistryMasterPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (IKeeperRegistryMasterPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (IKeeperRegistryMasterPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (IKeeperRegistryMasterPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (IKeeperRegistryMasterPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (IKeeperRegistryMasterReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (IKeeperRegistryMasterStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (IKeeperRegistryMasterTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (IKeeperRegistryMasterUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (IKeeperRegistryMasterUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (IKeeperRegistryMasterUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (IKeeperRegistryMasterUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (IKeeperRegistryMasterUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (IKeeperRegistryMasterUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (IKeeperRegistryMasterUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (IKeeperRegistryMasterUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (IKeeperRegistryMasterUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (IKeeperRegistryMasterUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (IKeeperRegistryMasterUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (IKeeperRegistryMasterUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (IKeeperRegistryMasterUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (IKeeperRegistryMasterUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (IKeeperRegistryMasterUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_IKeeperRegistryMaster *IKeeperRegistryMaster) Address() common.Address {
	return _IKeeperRegistryMaster.address
}

type IKeeperRegistryMasterInterface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error)

	GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error)

	GetConditionalTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetMode(opts *bind.CallOpts) (uint8, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

		error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

		error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (KeeperRegistryBase21UpkeepInfo, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

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

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*IKeeperRegistryMasterCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*IKeeperRegistryMasterConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*IKeeperRegistryMasterConfigSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*IKeeperRegistryMasterFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*IKeeperRegistryMasterFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*IKeeperRegistryMasterInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*IKeeperRegistryMasterOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*IKeeperRegistryMasterOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*IKeeperRegistryMasterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IKeeperRegistryMasterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*IKeeperRegistryMasterOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*IKeeperRegistryMasterPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*IKeeperRegistryMasterPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*IKeeperRegistryMasterPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*IKeeperRegistryMasterPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*IKeeperRegistryMasterPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*IKeeperRegistryMasterPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*IKeeperRegistryMasterPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*IKeeperRegistryMasterReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*IKeeperRegistryMasterStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*IKeeperRegistryMasterTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*IKeeperRegistryMasterTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*IKeeperRegistryMasterUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*IKeeperRegistryMasterUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*IKeeperRegistryMasterUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*IKeeperRegistryMasterUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*IKeeperRegistryMasterUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*IKeeperRegistryMasterUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*IKeeperRegistryMasterUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*IKeeperRegistryMasterUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*IKeeperRegistryMasterUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*IKeeperRegistryMasterUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*IKeeperRegistryMasterUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*IKeeperRegistryMasterUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*IKeeperRegistryMasterUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*IKeeperRegistryMasterUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*IKeeperRegistryMasterUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*IKeeperRegistryMasterUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *IKeeperRegistryMasterUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*IKeeperRegistryMasterUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
