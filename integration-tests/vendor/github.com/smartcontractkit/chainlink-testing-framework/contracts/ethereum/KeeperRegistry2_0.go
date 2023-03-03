// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OnchainConfig2_0 is an auto generated low-level Go binding around an user-defined struct.
type OnchainConfig2_0 struct {
	PaymentPremiumPPB    uint32
	FlatFeeMicroLink     uint32
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	MinUpkeepSpend       *big.Int
	MaxPerformGas        uint32
	MaxCheckDataSize     uint32
	MaxPerformDataSize   uint32
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Transcoder           common.Address
	Registrar            common.Address
}

// State2_0 is an auto generated low-level Go binding around an user-defined struct.
type State2_0 struct {
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

// UpkeepInfo is an auto generated low-level Go binding around an user-defined struct.
type UpkeepInfo struct {
	Target                 common.Address
	ExecuteGas             uint32
	CheckData              []byte
	Balance                *big.Int
	Admin                  common.Address
	MaxValidBlocknumber    uint64
	LastPerformBlockNumber uint32
	AmountSpent            *big.Int
	Paused                 bool
	OffchainConfig         []byte
}

// KeeperRegistry20MetaData contains all meta data concerning the KeeperRegistry20 contract.
var KeeperRegistry20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryBase2_0\",\"name\":\"keeperRegistryLogic\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"checkBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKeeperRegistryLogicAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPaymentModel\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structUpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b50604051620052d9380380620052d98339810160408190526200003591620003a7565b806001600160a01b031663f15701416040518163ffffffff1660e01b815260040160206040518083038186803b1580156200006f57600080fd5b505afa15801562000084573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000aa9190620003ce565b816001600160a01b031663ca30e6036040518163ffffffff1660e01b815260040160206040518083038186803b158015620000e457600080fd5b505afa158015620000f9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200011f9190620003a7565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015957600080fd5b505afa1580156200016e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001949190620003a7565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b815260040160206040518083038186803b158015620001ce57600080fd5b505afa158015620001e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002099190620003a7565b3380600081620002605760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000293576200029381620002fb565b505050836002811115620002ab57620002ab620003f1565b60e0816002811115620002c257620002c2620003f1565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b909116610100525062000420565b6001600160a01b038116331415620003565760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000257565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620003ba57600080fd5b8151620003c78162000407565b9392505050565b600060208284031215620003e157600080fd5b815160038110620003c757600080fd5b634e487b7160e01b600052602160045260246000fd5b6001600160a01b03811681146200041d57600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c614e42620004976000396000610eda01526000818161083201528181613621015261379c01526000818161048e0152612b050152600081816106b10152612be101526000818161076c0152610fd80152614e426000f3fe60806040526004361061024c5760003560e01c80638e86139b1161013a578063b657bc9c116100b1578063b657bc9c14610710578063b79550be1461039a578063c7c3a19a14610730578063c804802214610432578063ca30e6031461075d578063e3d0e71214610790578063eb5dcd6c146105f8578063ed56b3e1146107b0578063f157014114610823578063f2fde38b14610856578063f7d334ba14610876578063faa3e996146108a85761025b565b80638e86139b146105a2578063948108f7146105bd5780639fab438614610587578063a4c0ed36146105d8578063a710b221146105f8578063a72aa27e14610613578063aed2e9291461062e578063afcb95d714610665578063b10b673c146106a2578063b121e147146106d5578063b148ab6b14610432578063b1dc65a4146106f05761025b565b8063572e05e1116101ce578063572e05e1146104525780636709d0e51461047f5780636ded9eae146104b2578063744bfe611461036457806379ba5097146104e05780637d9b97e01461039a57806381ff7048146104f55780638456cb591461039a57806385c1b0ba146105495780638765ecbe146104325780638da5cb5b146105695780638dcf0fe7146105875761025b565b806306e3b632146102635780630e08ae8414610299578063181f5a77146102d15780631865c57d1461031e578063187256e8146103445780631a2af011146103645780633b9cce591461037f5780633f4ba83a1461039a578063421d183b146103af57806348013d7b146104105780635165f2f5146104325761025b565b3661025b576102596108ee565b005b6102596108ee565b34801561026f57600080fd5b5061028361027e3660046142b0565b610900565b6040516102909190614619565b60405180910390f35b3480156102a557600080fd5b506102b96102b43660046143f1565b6109e0565b6040516001600160601b039091168152602001610290565b3480156102dd57600080fd5b506103116040518060400160405280601481526020017304b6565706572526567697374727920322e302e360641b81525081565b60405161029091906146c7565b34801561032a57600080fd5b50610333610abb565b604051610290959493929190614701565b34801561035057600080fd5b5061025961035f366004613d9c565b610dd2565b34801561037057600080fd5b5061025961035f366004614240565b34801561038b57600080fd5b5061025961035f366004613ed0565b3480156103a657600080fd5b50610259610dde565b3480156103bb57600080fd5b506103cf6103ca366004613d46565b610de6565b60408051951515865260ff90941660208601526001600160601b03928316938501939093521660608301526001600160a01b0316608082015260a001610290565b34801561041c57600080fd5b50610425600281565b60405161029091906146f4565b34801561043e57600080fd5b5061025961044d36600461420e565b610ecd565b34801561045e57600080fd5b50610467610ed8565b6040516001600160a01b039091168152602001610290565b34801561048b57600080fd5b507f0000000000000000000000000000000000000000000000000000000000000000610467565b3480156104be57600080fd5b506104d26104cd366004613e29565b610efc565b604051908152602001610290565b3480156104ec57600080fd5b50610259610f11565b34801561050157600080fd5b50610526601254600e5463ffffffff600160601b8304811693600160801b9093041691565b6040805163ffffffff948516815293909216602084015290820152606001610290565b34801561055557600080fd5b50610259610564366004614093565b610fc0565b34801561057557600080fd5b506000546001600160a01b0316610467565b34801561059357600080fd5b50610259610564366004614265565b3480156105ae57600080fd5b5061025961035f3660046140e9565b3480156105c957600080fd5b5061025961035f3660046143cc565b3480156105e457600080fd5b506102596105f3366004613dce565b610fcd565b34801561060457600080fd5b5061025961035f366004613d63565b34801561061f57600080fd5b5061025961035f3660046143a7565b34801561063a57600080fd5b5061064e610649366004614265565b61114c565b604080519215158352602083019190915201610290565b34801561067157600080fd5b50600e54600f5460408051600081526020810193909352600160e01b90910463ffffffff1690820152606001610290565b3480156106ae57600080fd5b507f0000000000000000000000000000000000000000000000000000000000000000610467565b3480156106e157600080fd5b5061025961044d366004613d46565b3480156106fc57600080fd5b5061025961070b366004613fdd565b611257565b34801561071c57600080fd5b506102b961072b36600461420e565b611c8e565b34801561073c57600080fd5b5061075061074b36600461420e565b611cb2565b6040516102909190614804565b34801561076957600080fd5b507f0000000000000000000000000000000000000000000000000000000000000000610467565b34801561079c57600080fd5b506102596107ab366004613f11565b611f7c565b3480156107bc57600080fd5b5061080a6107cb366004613d46565b6001600160a01b031660009081526009602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff909116602083015201610290565b34801561082f57600080fd5b507f0000000000000000000000000000000000000000000000000000000000000000610425565b34801561086257600080fd5b50610259610871366004613d46565b612a6d565b34801561088257600080fd5b5061089661089136600461420e565b612a7e565b6040516102909695949392919061465d565b3480156108b457600080fd5b506108e16108c3366004613d46565b6001600160a01b031660009081526016602052604090205460ff1690565b60405161029091906146da565b6108fe6108f9610ed8565b612aa1565b565b6060600061090e6002612ac5565b905080841061093057604051631390f2a160e01b815260040160405180910390fd5b826109425761093f8482614c49565b92505b6000836001600160401b0381111561095c5761095c614d62565b604051908082528060200260200182016040528015610985578160200160208202803683370190505b50905060005b848110156109d7576109a86109a08288614af7565b600290612acf565b8282815181106109ba576109ba614d4c565b6020908102919091010152806109cf81614cef565b91505061098b565b50949350505050565b6040805161012081018252600f5460ff808216835263ffffffff61010080840482166020860152600160281b840482169585019590955262ffffff600160481b840416606085015261ffff600160601b8404166080850152600160701b83048216151560a0850152600160781b8304909116151560c08401526001600160601b03600160801b83041660e0840152600160e01b90910416918101919091526000908180610a8c83612ae2565b6012549193509150610ab29084908790600160c01b900463ffffffff1685856000612cc4565b95945050505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101a081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810182905261016081018290526101808101919091526040805161014081018252601254600160401b900463ffffffff1681526011546001600160601b03908116602083015260155492820192909252600f54600160801b900490911660608083019190915290819060009060808101610bd96002612ac5565b815260125463ffffffff600160601b8083048216602080860191909152600160801b84048316604080870191909152600e54606080880191909152600f54600160e01b810486166080808a019190915260ff600160701b83048116151560a09a8b015284516101a0810186526101008085048a168252600160281b85048a1682890152898b168288015262ffffff600160481b8604169582019590955261ffff88850416928101929092526010546001600160601b0381169a83019a909a52600160201b8904881660c0830152600160a01b8904881660e0830152600160c01b909804909616918601919091526013546101208601526014546101408601526001600160a01b0396849004871661016086015260115493909304909516610180840152600a8054865181840281018401909752808752969b509299508a958a959394600b9493169291859190830182828015610d5e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610d40575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610dba57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610d9c575b50505050509150945094509450945094509091929394565b610dda6108ee565b5050565b6108fe6108ee565b6001600160a01b03811660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526001600160601b03620100008404811692820192909252600160701b909204811660608301819052600f5484938493849384938492610e689291600160801b900416614c60565b600b54909150600090610e7b9083614b89565b905082600001518360200151828560400151610e979190614b53565b606095909501516001600160a01b039b8c166000908152600c6020526040902054929c919b959a50985093169550919350505050565b610ed56108ee565b50565b7f000000000000000000000000000000000000000000000000000000000000000090565b6000610f066108ee565b979650505050505050565b6001546001600160a01b03163314610f695760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610fc86108ee565b505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146110165760405163c8bad78d60e01b815260040160405180910390fd5b6020811461103757604051630dfe930960e41b815260040160405180910390fd5b60006110458284018461420e565b600081815260046020526040902054909150600160201b900463ffffffff9081161461108457604051634e0041d160e11b815260040160405180910390fd5b6000818152600460205260409020600101546110b1908590600160601b90046001600160601b0316614b53565b600082815260046020526040902060010180546001600160601b0392909216600160601b02600160601b600160c01b03199092169190911790556015546110f9908590614af7565b6015556040516001600160601b03851681526001600160a01b0386169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b600080611157612d0f565b600f54600160701b900460ff1615611182576040516309148bcd60e21b815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815463ffffffff8082168352600160201b8204811683860152600160401b820460ff16151583870152600160481b9091046001600160a01b031660608301526001909201546001600160601b038082166080840152600160601b82041660a0830152600160c01b900490911660c08201528251601f87018390048302810183019093528583529161124a918391908890889081908401838280828437600092019190915250612d2e92505050565b9250925050935093915050565b60005a6040805161012081018252600f5460ff808216835261010080830463ffffffff9081166020860152600160281b8404811695850195909552600160481b830462ffffff166060850152600160601b830461ffff166080850152600160701b8304821615801560a0860152600160781b8404909216151560c0850152600160801b83046001600160601b031660e0850152600160e01b90920490931690820152919250611319576040516309148bcd60e21b815260040160405180910390fd5b3360009081526008602052604090205460ff1661134957604051631099ed7560e01b815260040160405180910390fd5b600061138a8a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612e0392505050565b905060008160400151516001600160401b038111156113ab576113ab614d62565b60405190808252806020026020018201604052801561144157816020015b604080516101a081018252600060c0820181815260e083018290526101008301829052610120830182905261014083018290526101608301829052610180830182905282526020808301829052928201819052606082018190526080820181905260a082015282526000199092019101816113c95790505b5090506000805b8360400151518110156116a757600460008560400151838151811061146f5761146f614d4c565b6020908102919091018101518252818101929092526040908101600020815160e081018352815463ffffffff8082168352600160201b8204811695830195909552600160401b810460ff16151593820193909352600160481b9092046001600160a01b03166060830152600101546001600160601b038082166080840152600160601b82041660a0830152600160c01b900490911660c0820152835184908390811061151d5761151d614d4c565b6020026020010151600001819052506115868584838151811061154257611542614d4c565b602002602001015160000151600001518660600151848151811061156857611568614d4c565b60200260200101516040015151876000015188602001516001612cc4565b83828151811061159857611598614d4c565b6020026020010151604001906001600160601b031690816001600160601b03168152505061163c846040015182815181106115d5576115d5614d4c565b6020026020010151856060015183815181106115f3576115f3614d4c565b602002602001015185848151811061160d5761160d614d4c565b60200260200101516000015186858151811061162b5761162b614d4c565b602002602001015160400151612e96565b83828151811061164e5761164e614d4c565b6020026020010151602001901515908115158152505082818151811061167657611676614d4c565b6020026020010151602001511561169557611692600183614ad1565b91505b8061169f81614cef565b915050611448565b5061ffff81166116ca57604051637c01d16560e11b815260040160405180910390fd5b600e548d35146116ed5760405163dfdcf8e760e01b815260040160405180910390fd5b83516116fa906001614b2e565b60ff168914158061170b5750888714155b15611729576040516301227b8d60e11b815260040160405180910390fd5b6117398d8d8d8d8d8d8d8d612fce565b60005b8360400151518110156118fc5782818151811061175b5761175b614d4c565b602002602001015160200151156118ea574363ffffffff16600460008660400151848151811061178d5761178d614d4c565b6020026020010151815260200190815260200160002060010160189054906101000a900463ffffffff1663ffffffff1614156117dc57604051632d56b1d560e21b815260040160405180910390fd5b6118248382815181106117f1576117f1614d4c565b6020026020010151600001518560600151838151811061181357611813614d4c565b602002602001015160400151612d2e565b84838151811061183657611836614d4c565b602002602001015160600185848151811061185357611853614d4c565b6020026020010151608001828152508215151515815250505082818151811061187e5761187e614d4c565b602002602001015160800151866118959190614c49565b95504360046000866040015184815181106118b2576118b2614d4c565b6020026020010151815260200190815260200160002060010160186101000a81548163ffffffff021916908363ffffffff1602179055505b806118f481614cef565b91505061173c565b50835161190a906001614b2e565b6119199060ff1661044c614baf565b6169146119278d6010614baf565b5a6119329089614c49565b61193c9190614af7565b6119469190614af7565b6119509190614af7565b94506116a861196361ffff831687614b75565b61196d9190614af7565b945060008060008060005b876040015151811015611b735786818151811061199757611997614d4c565b60200260200101516020015115611b61576119d98a896060015183815181106119c2576119c2614d4c565b602002602001015160400151518b600001516131bc565b8782815181106119eb576119eb614d4c565b602002602001015160a0018181525050611a478989604001518381518110611a1557611a15614d4c565b6020026020010151898481518110611a2f57611a2f614d4c565b60200260200101518b600001518c602001518b6131da565b9093509150611a568285614b53565b9350611a628386614b53565b9450868181518110611a7657611a76614d4c565b602002602001015160600151151588604001518281518110611a9a57611a9a614d4c565b60200260200101517f29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a44208a606001518481518110611ad957611ad9614d4c565b6020026020010151600001518a8581518110611af757611af7614d4c565b6020026020010151608001518b8681518110611b1557611b15614d4c565b602002602001015160a001518789611b2d9190614b53565b6040805163ffffffff90951685526020850193909352918301526001600160601b0316606082015260800160405180910390a35b80611b6b81614cef565b915050611978565b50503360009081526008602052604090208054849250600290611ba69084906201000090046001600160601b0316614b53565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600f60000160108282829054906101000a90046001600160601b0316611bf19190614b53565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555060008f600160038110611c2a57611c2a614d4c565b602002013560001c9050600060088264ffffffffff16901c905087610100015163ffffffff168163ffffffff161115611c7b57600f80546001600160e01b0316600160e01b63ffffffff8416021790555b5050505050505050505050505050505050565b600081815260046020526040812054611cac9063ffffffff166109e0565b92915050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526000828152600460209081526040808320815160e081018352815463ffffffff8082168352600160201b8204811683870152600160401b820460ff16151583860152600160481b9091046001600160a01b03908116606084019081526001909401546001600160601b038082166080860152600160601b82041660a0850152600160c01b9004821660c0840152845161014081018652935116835281511682850152868552600790935292819020805492939291830191611dc290614cb4565b80601f0160208091040260200160405190810160405280929190818152602001828054611dee90614cb4565b8015611e3b5780601f10611e1057610100808354040283529160200191611e3b565b820191906000526020600020905b815481529060010190602001808311611e1e57829003601f168201915b505050505081526020018260a001516001600160601b031681526020016005600086815260200190815260200160002060009054906101000a90046001600160a01b03166001600160a01b03168152602001826020015163ffffffff166001600160401b031681526020018260c0015163ffffffff16815260200182608001516001600160601b03168152602001826040015115158152602001601760008681526020019081526020016000208054611ef390614cb4565b80601f0160208091040260200160405190810160405280929190818152602001828054611f1f90614cb4565b8015611f6c5780601f10611f4157610100808354040283529160200191611f6c565b820191906000526020600020905b815481529060010190602001808311611f4f57829003601f168201915b5050505050815250915050919050565b611f846132b0565b601f86511115611fa757604051630974082760e21b815260040160405180910390fd5b60ff8416611fc8576040516373bedd2b60e11b815260040160405180910390fd5b84518651141580611fe75750611fdf846003614bfa565b60ff16865111155b15612005576040516303a5a38b60e31b815260040160405180910390fd5b600f54600b54600160801b9091046001600160601b03169060005b816001600160601b031681101561207657612063600b828154811061204757612047614d4c565b6000918252602090912001546001600160a01b03168484613303565b508061206e81614cef565b915050612020565b5060008060005b836001600160601b031681101561212557600a81815481106120a1576120a1614d4c565b600091825260209091200154600b80546001600160a01b03909216945090829081106120cf576120cf614d4c565b60009182526020808320909101546001600160a01b038681168452600983526040808520805461ffff1916905591168084526008909252909120805460ff1916905591508061211d81614cef565b91505061207d565b50612132600a60006139b5565b61213e600b60006139b5565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c5181101561239d57600960008e838151811061218357612183614d4c565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156121c857604051633be7507d60e11b815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600960008f84815181106121f9576121f9614d4c565b6020908102919091018101516001600160a01b031682528181019290925260400160002082518154939092015160ff166101000261ff00199215159290921661ffff19909316929092171790558b518c908290811061225a5761225a614d4c565b6020908102919091018101516001600160a01b0381166000908152600883526040908190208151608081018352905460ff80821615801584526101008304909116958301959095526001600160601b03620100008204811693830193909352600160701b90049091166060820152945092506122e957604051636a7281ad60e01b815260040160405180910390fd5b6001835260ff80821660208086019182526001600160601b03808b16606088019081526001600160a01b03871660009081526008909352604092839020885181549551948a015192518416600160701b02600160701b600160d01b03199390941662010000029290921662010000600160d01b0319949096166101000261ff00199215159290921661ffff19909516949094171791909116929092179190911790558061239581614cef565b915050612164565b50508a516123b39150600a9060208d01906139d3565b5088516123c790600b9060208c01906139d3565b506000878060200190518101906123de919061411e565b60125460c082015191925063ffffffff600160201b90910481169116101561241957604051630e6af04160e21b815260040160405180910390fd5b60125460e082015163ffffffff600160a01b90920482169116101561245157604051631fa9bdcb60e01b815260040160405180910390fd5b60125461010082015163ffffffff600160c01b90920482169116101561248a57604051631a3abf5560e31b815260040160405180910390fd5b6040518061012001604052808a60ff168152602001826000015163ffffffff168152602001826020015163ffffffff168152602001826060015162ffffff168152602001826080015161ffff168152602001600015158152602001600015158152602001866001600160601b03168152602001600063ffffffff16815250600f60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160096101000a81548162ffffff021916908362ffffff160217905550608082015181600001600c6101000a81548161ffff021916908361ffff16021790555060a082015181600001600e6101000a81548160ff02191690831515021790555060c082015181600001600f6101000a81548160ff02191690831515021790555060e08201518160000160106101000a8154816001600160601b0302191690836001600160601b0316021790555061010082015181600001601c6101000a81548163ffffffff021916908363ffffffff1602179055509050506040518061016001604052808260a001516001600160601b031681526020018261016001516001600160a01b03168152602001601060010160009054906101000a90046001600160601b03166001600160601b031681526020018261018001516001600160a01b03168152602001826040015163ffffffff1681526020018260c0015163ffffffff168152602001601060020160089054906101000a900463ffffffff1663ffffffff1681526020016010600201600c9054906101000a900463ffffffff1663ffffffff168152602001601060020160109054906101000a900463ffffffff1663ffffffff1681526020018260e0015163ffffffff16815260200182610100015163ffffffff16815250601060008201518160000160006101000a8154816001600160601b0302191690836001600160601b03160217905550602082015181600001600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160006101000a8154816001600160601b0302191690836001600160601b03160217905550606082015181600101600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060808201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160020160046101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600201600c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160106101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160146101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160186101000a81548163ffffffff021916908363ffffffff1602179055509050508061012001516013819055508061014001516014819055506000601060020160109054906101000a900463ffffffff16905043601060020160106101000a81548163ffffffff021916908363ffffffff16021790555060016010600201600c8282829054906101000a900463ffffffff166129b19190614b0f565b92506101000a81548163ffffffff021916908363ffffffff1602179055506129fb46306010600201600c9054906101000a900463ffffffff1663ffffffff168f8f8f8f8f8f613481565b600e819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600e546010600201600c9054906101000a900463ffffffff168f8f8f8f8f8f604051612a5799989796959493929190614987565b60405180910390a1505050505050505050505050565b612a756132b0565b610ed5816134db565b60006060600080600080612a90612d0f565b612a986108ee565b91939550919395565b3660008037600080366000845af43d6000803e808015612ac0573d6000f35b3d6000fd5b6000611cac825490565b6000612adb838361357f565b9392505050565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015612b5c57600080fd5b505afa158015612b70573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b94919061440e565b5094509092505050600081131580612bab57508142105b80612bcc5750828015612bcc5750612bc38242614c49565b8463ffffffff16105b15612bdb576013549550612bdf565b8095505b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015612c3857600080fd5b505afa158015612c4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c70919061440e565b5094509092505050600081131580612c8757508142105b80612ca85750828015612ca85750612c9f8242614c49565b8463ffffffff16105b15612cb7576014549450612cbb565b8094505b50505050915091565b600080612cd58689600001516135a9565b9050600080612cf08a8a63ffffffff16858a8a60018b6135ec565b9092509050612cff8183614b53565b93505050505b9695505050505050565b32156108fe5760405163b60ac5db60e01b815260040160405180910390fd5b600f546000908190600160781b900460ff1615612d5e576040516306fda65d60e31b815260040160405180910390fd5b600f805460ff60781b1916600160781b1790555a90506000634585e33b60e01b84604051602401612d8f91906146c7565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b0383818316178352505050509050612ddd856000015163ffffffff16866060015183613969565b92505a612dea9083614c49565b915050600f805460ff60781b1916905590939092509050565b612e2e6040518060800160405280600081526020016000815260200160608152602001606081525090565b60008060008085806020019051810190612e4891906142d2565b93509350935093508051825114612e7257604051632d56b1d560e21b815260040160405180910390fd5b60408051608081018252948552602085019390935291830152606082015292915050565b60008260c0015163ffffffff16846000015163ffffffff161015612ee75760405185907f5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a8990600090a2506000612fc6565b6020840151845163ffffffff164014612f2d5760405185907f561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc1390600090a2506000612fc6565b43836020015163ffffffff1611612f715760405185907fd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f90600090a2506000612fc6565b816001600160601b03168360a001516001600160601b03161015612fc25760405185907f7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb9690600090a2506000612fc6565b5060015b949350505050565b60008787604051612fe09291906145e2565b604051908190038120612ff7918b906020016146ad565b60408051601f1981840301815282825280516020918201208383019092526000808452908301819052909250906000805b8881101561316c5760018587836020811061304557613045614d4c565b61305291901a601b614b2e565b8c8c8581811061306457613064614d4c565b905060200201358b8b8681811061307d5761307d614d4c565b90506020020135604051600081526020016040526040516130ba949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156130dc573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526009602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061314657604051630f4c073760e01b815260040160405180910390fd5b826020015160080260ff166001901b84019350808061316490614cef565b915050613028565b50827e010101010101010101010101010101010101010101010101010101010101018416146131ae57604051636081df1760e11b815260040160405180910390fd5b505050505050505050505050565b60006131c883836135a9565b905080841015612adb57509192915050565b6000806131f58887608001518860a0015188888860016135ec565b909250905060006132068284614b53565b600089815260046020526040902060010180549192508291600c9061323c908490600160601b90046001600160601b0316614c60565b82546101009290920a6001600160601b0381810219909316918316021790915560008a81526004602052604081206001018054859450909261328091859116614b53565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050965096945050505050565b6000546001600160a01b031633146108fe5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610f60565b6001600160a01b03831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526001600160601b03620100008404811692820192909252600160701b909204166060820181905282906133719086614c60565b9050600061337f8583614b89565b905080836040018181516133939190614b53565b6001600160601b0390811690915287166060850152506133b38582614c23565b6133bd9083614c60565b601180546000906133d89084906001600160601b0316614b53565b825461010092830a6001600160601b038181021990921692821602919091179092556001600160a01b039990991660009081526008602090815260409182902087518154928901519389015160609099015161ffff1990931690151561ff0019161760ff909316909b029190911762010000600160d01b0319166201000087841602600160701b600160d01b03191617600160701b919092160217909755509095945050505050565b6000808a8a8a8a8a8a8a8a8a6040516020016134a5999897969594939291906148ee565b60408051808303601f1901815291905280516020909101206001600160f01b0316600160f01b179b9a5050505050505050505050565b6001600160a01b03811633141561352e5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610f60565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061359657613596614d4c565b9060005260206000200154905092915050565b60006135bc63ffffffff84166014614baf565b6135c7836001614b2e565b6135d69060ff16611d4c614baf565b6135e29061fde8614af7565b612adb9190614af7565b6000806000896080015161ffff16876136059190614baf565b90508380156136135750803a105b1561361b57503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561365157613651614d36565b14156137985760408051600081526020810190915285156136b057600036604051806080016040528060488152602001614dc56048913960405160200161369a939291906145f2565b6040516020818303038152906040529050613716565b6012546136cb90600160c01b900463ffffffff166004614bce565b63ffffffff166001600160401b038111156136e8576136e8614d62565b6040519080825280601f01601f191660200182016040528015613712576020820181803683370190505b5090505b6040516324ca470760e11b8152600f602160991b01906349948e0e906137409084906004016146c7565b60206040518083038186803b15801561375857600080fd5b505afa15801561376c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137909190614227565b915050613847565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156137cc576137cc614d36565b141561384757606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561380c57600080fd5b505afa158015613820573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138449190614227565b90505b8461386357808b6080015161ffff166138609190614baf565b90505b61387161ffff871682614b75565b9050600087826138818c8e614af7565b61388b9086614baf565b6138959190614af7565b6138a790670de0b6b3a7640000614baf565b6138b19190614b75565b905060008c6040015163ffffffff1664e8d4a510006138d09190614baf565b898e6020015163ffffffff16858f886138e99190614baf565b6138f39190614af7565b61390190633b9aca00614baf565b61390b9190614baf565b6139159190614b75565b61391f9190614af7565b90506b033b2e3c9fd0803ce80000006139388284614af7565b11156139575760405163156baa3d60e11b815260040160405180910390fd5b909c909b509950505050505050505050565b60005a61138881101561397b57600080fd5b61138881039050846040820482031161399357600080fd5b50823b61399f57600080fd5b60008083516020850160008789f1949350505050565b5080546000825590600052602060002090810190610ed59190613a38565b828054828255906000526020600020908101928215613a28579160200282015b82811115613a2857825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906139f3565b50613a34929150613a38565b5090565b5b80821115613a345760008155600101613a39565b8051613a5881614d88565b919050565b60008083601f840112613a6f57600080fd5b5081356001600160401b03811115613a8657600080fd5b6020830191508360208260051b8501011115613aa157600080fd5b9250929050565b600082601f830112613ab957600080fd5b81356020613ace613ac983614a87565b614a57565b80838252828201915082860187848660051b8901011115613aee57600080fd5b60005b85811015613b16578135613b0481614d88565b84529284019290840190600101613af1565b5090979650505050505050565b600082601f830112613b3457600080fd5b81516020613b44613ac983614a87565b80838252828201915082860187848660051b8901011115613b6457600080fd5b60005b85811015613b165781516001600160401b0380821115613b8657600080fd5b908901906060828c03601f1901811315613b9f57600080fd5b613ba7614a0c565b88840151613bb481614d9d565b81526040848101518a830152918401519183831115613bd257600080fd5b82850194508d603f860112613be657600080fd5b898501519350613bf8613ac985614aaa565b92508383528d81858701011115613c0e57600080fd5b613c1d848b8501838801614c88565b810191909152865250509284019290840190600101613b67565b60008083601f840112613c4957600080fd5b5081356001600160401b03811115613c6057600080fd5b602083019150836020828501011115613aa157600080fd5b600082601f830112613c8957600080fd5b8135613c97613ac982614aaa565b818152846020838601011115613cac57600080fd5b816020850160208301376000918101602001919091529392505050565b805161ffff81168114613a5857600080fd5b805162ffffff81168114613a5857600080fd5b8051613a5881614d9d565b80356001600160401b0381168114613a5857600080fd5b803560ff81168114613a5857600080fd5b805169ffffffffffffffffffff81168114613a5857600080fd5b8051613a5881614daf565b600060208284031215613d5857600080fd5b8135612adb81614d88565b60008060408385031215613d7657600080fd5b8235613d8181614d88565b91506020830135613d9181614d88565b809150509250929050565b60008060408385031215613daf57600080fd5b8235613dba81614d88565b9150602083013560048110613d9157600080fd5b60008060008060608587031215613de457600080fd5b8435613def81614d88565b93506020850135925060408501356001600160401b03811115613e1157600080fd5b613e1d87828801613c37565b95989497509550505050565b600080600080600080600060a0888a031215613e4457600080fd5b8735613e4f81614d88565b96506020880135613e5f81614d9d565b95506040880135613e6f81614d88565b945060608801356001600160401b0380821115613e8b57600080fd5b613e978b838c01613c37565b909650945060808a0135915080821115613eb057600080fd5b50613ebd8a828b01613c37565b989b979a50959850939692959293505050565b60008060208385031215613ee357600080fd5b82356001600160401b03811115613ef957600080fd5b613f0585828601613a5d565b90969095509350505050565b60008060008060008060c08789031215613f2a57600080fd5b86356001600160401b0380821115613f4157600080fd5b613f4d8a838b01613aa8565b97506020890135915080821115613f6357600080fd5b613f6f8a838b01613aa8565b9650613f7d60408a01613d10565b95506060890135915080821115613f9357600080fd5b613f9f8a838b01613c78565b9450613fad60808a01613cf9565b935060a0890135915080821115613fc357600080fd5b50613fd089828a01613c78565b9150509295509295509295565b60008060008060008060008060e0898b031215613ff957600080fd5b606089018a81111561400a57600080fd5b899850356001600160401b038082111561402357600080fd5b61402f8c838d01613c37565b909950975060808b013591508082111561404857600080fd5b6140548c838d01613a5d565b909750955060a08b013591508082111561406d57600080fd5b5061407a8b828c01613a5d565b999c989b50969995989497949560c00135949350505050565b6000806000604084860312156140a857600080fd5b83356001600160401b038111156140be57600080fd5b6140ca86828701613a5d565b90945092505060208401356140de81614d88565b809150509250925092565b600080602083850312156140fc57600080fd5b82356001600160401b0381111561411257600080fd5b613f0585828601613c37565b60006101a0828403121561413157600080fd5b614139614a34565b61414283613cee565b815261415060208401613cee565b602082015261416160408401613cee565b604082015261417260608401613cdb565b606082015261418360808401613cc9565b608082015261419460a08401613d3b565b60a08201526141a560c08401613cee565b60c08201526141b660e08401613cee565b60e08201526101006141c9818501613cee565b90820152610120838101519082015261014080840151908201526101606141f1818501613a4d565b90820152610180614203848201613a4d565b908201529392505050565b60006020828403121561422057600080fd5b5035919050565b60006020828403121561423957600080fd5b5051919050565b6000806040838503121561425357600080fd5b823591506020830135613d9181614d88565b60008060006040848603121561427a57600080fd5b8335925060208401356001600160401b0381111561429757600080fd5b6142a386828701613c37565b9497909650939450505050565b600080604083850312156142c357600080fd5b50508035926020909101359150565b600080600080608085870312156142e857600080fd5b84519350602080860151935060408601516001600160401b038082111561430e57600080fd5b818801915088601f83011261432257600080fd5b8151614330613ac982614a87565b8082825285820191508585018c878560051b880101111561435057600080fd5b600095505b83861015614373578051835260019590950194918601918601614355565b5060608b0151909750945050508083111561438d57600080fd5b505061439b87828801613b23565b91505092959194509250565b600080604083850312156143ba57600080fd5b823591506020830135613d9181614d9d565b600080604083850312156143df57600080fd5b823591506020830135613d9181614daf565b60006020828403121561440357600080fd5b8135612adb81614d9d565b600080600080600060a0868803121561442657600080fd5b61442f86613d21565b945060208601519350604086015192506060860151915061445260808701613d21565b90509295509295909350565b6001600160a01b03169052565b600081518084526020808501945080840160005b838110156144a45781516001600160a01b03168752958201959082019060010161447f565b509495945050505050565b600081518084526144c7816020860160208601614c88565b601f01601f19169290920160200192915050565b805163ffffffff16825260208101516144fc602084018263ffffffff169052565b506040810151614514604084018263ffffffff169052565b50606081015161452b606084018262ffffff169052565b506080810151614541608084018261ffff169052565b5060a081015161455c60a08401826001600160601b03169052565b5060c081015161457460c084018263ffffffff169052565b5060e081015161458c60e084018263ffffffff169052565b506101008181015163ffffffff169083015261012080820151908301526101408082015190830152610160808201516145c78285018261445e565b5050610180808201516145dc8285018261445e565b50505050565b8183823760009101908152919050565b82848237600083820160008152835161460f818360208801614c88565b0195945050505050565b6020808252825182820181905260009190848201906040850190845b8181101561465157835183529284019291840191600101614635565b50909695505050505050565b861515815260c06020820152600061467860c08301886144af565b90506007861061468a5761468a614d36565b8560408301528460608301528360808301528260a0830152979650505050505050565b828152608081016060836020840137600081529392505050565b602081526000612adb60208301846144af565b60208101600483106146ee576146ee614d36565b91905290565b602081016146ee83614d78565b855163ffffffff1681526000610340602088015161472a60208501826001600160601b03169052565b5060408801516040840152606088015161474f60608501826001600160601b03169052565b506080880151608084015260a088015161477160a085018263ffffffff169052565b5060c088015161478960c085018263ffffffff169052565b5060e088015160e0840152610100808901516147ac8286018263ffffffff169052565b5050610120888101511515908401526147c96101408401886144db565b806102e08401526147dc8184018761446b565b90508281036103008401526147f1818661446b565b915050612d0561032083018460ff169052565b6020815261481660208201835161445e565b6000602083015161482f604084018263ffffffff169052565b50604083015161014080606085015261484c6101608501836144af565b9150606085015161486860808601826001600160601b03169052565b50608085015161487b60a086018261445e565b5060a08501516001600160401b03811660c08601525060c085015163ffffffff811660e08601525060e08501516101006148bf818701836001600160601b03169052565b86015190506101206148d48682018315159052565b860151858403601f1901838701529050612d0583826144af565b8981526001600160a01b03891660208201526001600160401b038881166040830152610120606083018190526000916149298483018b61446b565b9150838203608085015261493d828a61446b565b915060ff881660a085015283820360c085015261495a82886144af565b90861660e0850152838103610100850152905061497781856144af565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526149b78184018a61446b565b905082810360808401526149cb818961446b565b905060ff871660a084015282810360c08401526149e881876144af565b90506001600160401b03851660e084015282810361010084015261497781856144af565b604051606081016001600160401b0381118282101715614a2e57614a2e614d62565b60405290565b6040516101a081016001600160401b0381118282101715614a2e57614a2e614d62565b604051601f8201601f191681016001600160401b0381118282101715614a7f57614a7f614d62565b604052919050565b60006001600160401b03821115614aa057614aa0614d62565b5060051b60200190565b60006001600160401b03821115614ac357614ac3614d62565b50601f01601f191660200190565b600061ffff808316818516808303821115614aee57614aee614d0a565b01949350505050565b60008219821115614b0a57614b0a614d0a565b500190565b600063ffffffff808316818516808303821115614aee57614aee614d0a565b600060ff821660ff84168060ff03821115614b4b57614b4b614d0a565b019392505050565b60006001600160601b03828116848216808303821115614aee57614aee614d0a565b600082614b8457614b84614d20565b500490565b60006001600160601b0383811680614ba357614ba3614d20565b92169190910492915050565b6000816000190483118215151615614bc957614bc9614d0a565b500290565b600063ffffffff80831681851681830481118215151615614bf157614bf1614d0a565b02949350505050565b600060ff821660ff84168160ff0481118215151615614c1b57614c1b614d0a565b029392505050565b60006001600160601b0382811684821681151582840482111615614bf157614bf1614d0a565b600082821015614c5b57614c5b614d0a565b500390565b60006001600160601b0383811690831681811015614c8057614c80614d0a565b039392505050565b60005b83811015614ca3578181015183820152602001614c8b565b838111156145dc5750506000910152565b600181811c90821680614cc857607f821691505b60208210811415614ce957634e487b7160e01b600052602260045260246000fd5b50919050565b6000600019821415614d0357614d03614d0a565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b60038110610ed557610ed5614d36565b6001600160a01b0381168114610ed557600080fd5b63ffffffff81168114610ed557600080fd5b6001600160601b0381168114610ed557600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a264697066735822122019b44fe0e3e78a0abaaec652e8a52bdc66fcf8a2984043def62ab509a3c6d77364736f6c63430008060033",
}

// KeeperRegistry20ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistry20MetaData.ABI instead.
var KeeperRegistry20ABI = KeeperRegistry20MetaData.ABI

// KeeperRegistry20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistry20MetaData.Bin instead.
var KeeperRegistry20Bin = KeeperRegistry20MetaData.Bin

// DeployKeeperRegistry20 deploys a new Ethereum contract, binding an instance of KeeperRegistry20 to it.
func DeployKeeperRegistry20(auth *bind.TransactOpts, backend bind.ContractBackend, keeperRegistryLogic common.Address) (common.Address, *types.Transaction, *KeeperRegistry20, error) {
	parsed, err := KeeperRegistry20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry20Bin), backend, keeperRegistryLogic)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry20{KeeperRegistry20Caller: KeeperRegistry20Caller{contract: contract}, KeeperRegistry20Transactor: KeeperRegistry20Transactor{contract: contract}, KeeperRegistry20Filterer: KeeperRegistry20Filterer{contract: contract}}, nil
}

// KeeperRegistry20 is an auto generated Go binding around an Ethereum contract.
type KeeperRegistry20 struct {
	KeeperRegistry20Caller     // Read-only binding to the contract
	KeeperRegistry20Transactor // Write-only binding to the contract
	KeeperRegistry20Filterer   // Log filterer for contract events
}

// KeeperRegistry20Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistry20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistry20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistry20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistry20Session struct {
	Contract     *KeeperRegistry20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistry20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistry20CallerSession struct {
	Contract *KeeperRegistry20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// KeeperRegistry20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistry20TransactorSession struct {
	Contract     *KeeperRegistry20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// KeeperRegistry20Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistry20Raw struct {
	Contract *KeeperRegistry20 // Generic contract binding to access the raw methods on
}

// KeeperRegistry20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistry20CallerRaw struct {
	Contract *KeeperRegistry20Caller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistry20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistry20TransactorRaw struct {
	Contract *KeeperRegistry20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistry20 creates a new instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20(address common.Address, backend bind.ContractBackend) (*KeeperRegistry20, error) {
	contract, err := bindKeeperRegistry20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20{KeeperRegistry20Caller: KeeperRegistry20Caller{contract: contract}, KeeperRegistry20Transactor: KeeperRegistry20Transactor{contract: contract}, KeeperRegistry20Filterer: KeeperRegistry20Filterer{contract: contract}}, nil
}

// NewKeeperRegistry20Caller creates a new read-only instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistry20Caller, error) {
	contract, err := bindKeeperRegistry20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Caller{contract: contract}, nil
}

// NewKeeperRegistry20Transactor creates a new write-only instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistry20Transactor, error) {
	contract, err := bindKeeperRegistry20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Transactor{contract: contract}, nil
}

// NewKeeperRegistry20Filterer creates a new log filterer instance of KeeperRegistry20, bound to a specific deployed contract.
func NewKeeperRegistry20Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistry20Filterer, error) {
	contract, err := bindKeeperRegistry20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20Filterer{contract: contract}, nil
}

// bindKeeperRegistry20 binds a generic wrapper to an already deployed contract.
func bindKeeperRegistry20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistry20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry20.Contract.KeeperRegistry20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.KeeperRegistry20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry20 *KeeperRegistry20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.KeeperRegistry20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry20 *KeeperRegistry20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry20 *KeeperRegistry20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry20 *KeeperRegistry20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.contract.Transact(opts, method, params...)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry20.Contract.GetActiveUpkeepIDs(&_KeeperRegistry20.CallOpts, startIndex, maxCount)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry20.Contract.GetActiveUpkeepIDs(&_KeeperRegistry20.CallOpts, startIndex, maxCount)
}

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetFastGasFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetFastGasFeedAddress is a free data retrieval call binding the contract method 0x6709d0e5.
//
// Solidity: function getFastGasFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetFastGasFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetKeeperRegistryLogicAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getKeeperRegistryLogicAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetKeeperRegistryLogicAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetKeeperRegistryLogicAddress(&_KeeperRegistry20.CallOpts)
}

// GetKeeperRegistryLogicAddress is a free data retrieval call binding the contract method 0x572e05e1.
//
// Solidity: function getKeeperRegistryLogicAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetKeeperRegistryLogicAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetKeeperRegistryLogicAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkAddress is a free data retrieval call binding the contract method 0xca30e603.
//
// Solidity: function getLinkAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetLinkNativeFeedAddress is a free data retrieval call binding the contract method 0xb10b673c.
//
// Solidity: function getLinkNativeFeedAddress() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry20.CallOpts)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMaxPaymentForGas(&_KeeperRegistry20.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x0e08ae84.
//
// Solidity: function getMaxPaymentForGas(uint32 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMaxPaymentForGas(&_KeeperRegistry20.CallOpts, gasLimit)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry20.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetPaymentModel(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getPaymentModel")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetPaymentModel() (uint8, error) {
	return _KeeperRegistry20.Contract.GetPaymentModel(&_KeeperRegistry20.CallOpts)
}

// GetPaymentModel is a free data retrieval call binding the contract method 0xf1570141.
//
// Solidity: function getPaymentModel() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetPaymentModel() (uint8, error) {
	return _KeeperRegistry20.Contract.GetPaymentModel(&_KeeperRegistry20.CallOpts)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry20.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry20.CallOpts, peer)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry20.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry20.CallOpts, peer)
}

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(struct {
		Active bool
		Index  uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetSignerInfo(query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	return _KeeperRegistry20.Contract.GetSignerInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetSignerInfo is a free data retrieval call binding the contract method 0xed56b3e1.
//
// Solidity: function getSignerInfo(address query) view returns(bool active, uint8 index)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetSignerInfo(query common.Address) (struct {
	Active bool
	Index  uint8
}, error) {
	return _KeeperRegistry20.Contract.GetSignerInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetState(opts *bind.CallOpts) (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		State        State2_0
		Config       OnchainConfig2_0
		Signers      []common.Address
		Transmitters []common.Address
		F            uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State2_0)).(*State2_0)
	outstruct.Config = *abi.ConvertType(out[1], new(OnchainConfig2_0)).(*OnchainConfig2_0)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetState() (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint96,uint256,uint32,uint32,bytes32,uint32,bool) state, (uint32,uint32,uint32,uint24,uint16,uint96,uint32,uint32,uint32,uint256,uint256,address,address) config, address[] signers, address[] transmitters, uint8 f)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetState() (struct {
	State        State2_0
	Config       OnchainConfig2_0
	Signers      []common.Address
	Transmitters []common.Address
	F            uint8
}, error) {
	return _KeeperRegistry20.Contract.GetState(&_KeeperRegistry20.CallOpts)
}

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getTransmitterInfo", query)

	outstruct := new(struct {
		Active        bool
		Index         uint8
		Balance       *big.Int
		LastCollected *big.Int
		Payee         common.Address
	})
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

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetTransmitterInfo(query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetTransmitterInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetTransmitterInfo is a free data retrieval call binding the contract method 0x421d183b.
//
// Solidity: function getTransmitterInfo(address query) view returns(bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetTransmitterInfo(query common.Address) (struct {
	Active        bool
	Index         uint8
	Balance       *big.Int
	LastCollected *big.Int
	Payee         common.Address
}, error) {
	return _KeeperRegistry20.Contract.GetTransmitterInfo(&_KeeperRegistry20.CallOpts, query)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (UpkeepInfo, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(UpkeepInfo)).(*UpkeepInfo)

	return out0, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20Session) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns((address,uint32,bytes,uint96,address,uint64,uint32,uint96,bool,bytes) upkeepInfo)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry20.Contract.GetUpkeep(&_KeeperRegistry20.CallOpts, id)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20Session) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDetails(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDetails(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Caller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Session) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry20.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _KeeperRegistry20.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20Session) Owner() (common.Address, error) {
	return _KeeperRegistry20.Contract.Owner(&_KeeperRegistry20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry20.Contract.Owner(&_KeeperRegistry20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20Session) TypeAndVersion() (string, error) {
	return _KeeperRegistry20.Contract.TypeAndVersion(&_KeeperRegistry20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry20.Contract.TypeAndVersion(&_KeeperRegistry20.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry20.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry20.Contract.UpkeepTranscoderVersion(&_KeeperRegistry20.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry20 *KeeperRegistry20CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry20.Contract.UpkeepTranscoderVersion(&_KeeperRegistry20.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptOwnership(&_KeeperRegistry20.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptOwnership(&_KeeperRegistry20.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptPayeeship", transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptPayeeship(&_KeeperRegistry20.TransactOpts, transmitter)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AcceptUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "addFunds", id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AddFunds(&_KeeperRegistry20.TransactOpts, id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.AddFunds(&_KeeperRegistry20.TransactOpts, id, amount)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "cancelUpkeep", id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CancelUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CancelUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "checkUpkeep", id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20Session) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xf7d334ba.
//
// Solidity: function checkUpkeep(uint256 id) returns(bool upkeepNeeded, bytes performData, uint8 upkeepFailureReason, uint256 gasUsed, uint256 fastGasWei, uint256 linkNative)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.CheckUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.MigrateUpkeeps(&_KeeperRegistry20.TransactOpts, ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.MigrateUpkeeps(&_KeeperRegistry20.TransactOpts, ids, destination)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.OnTokenTransfer(&_KeeperRegistry20.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.OnTokenTransfer(&_KeeperRegistry20.TransactOpts, sender, amount, data)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Pause(&_KeeperRegistry20.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Pause(&_KeeperRegistry20.TransactOpts)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "pauseUpkeep", id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.PauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.ReceiveUpkeeps(&_KeeperRegistry20.TransactOpts, encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.ReceiveUpkeeps(&_KeeperRegistry20.TransactOpts, encodedUpkeeps)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "recoverFunds")
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RecoverFunds(&_KeeperRegistry20.TransactOpts)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RecoverFunds(&_KeeperRegistry20.TransactOpts)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, offchainConfig)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x6ded9eae.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData, bytes offchainConfig) returns(uint256 id)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.RegisterUpkeep(&_KeeperRegistry20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetConfig(&_KeeperRegistry20.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setPayees", payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPayees(&_KeeperRegistry20.TransactOpts, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x3b9cce59.
//
// Solidity: function setPayees(address[] payees) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPayees(&_KeeperRegistry20.TransactOpts, payees)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry20.TransactOpts, peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry20.TransactOpts, peer, permission)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepGasLimit(&_KeeperRegistry20.TransactOpts, id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepGasLimit(&_KeeperRegistry20.TransactOpts, id, gasLimit)
}

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry20.TransactOpts, id, config)
}

// SetUpkeepOffchainConfig is a paid mutator transaction binding the contract method 0x8dcf0fe7.
//
// Solidity: function setUpkeepOffchainConfig(uint256 id, bytes config) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry20.TransactOpts, id, config)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20Transactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20Session) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SimulatePerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
}

// SimulatePerformUpkeep is a paid mutator transaction binding the contract method 0xaed2e929.
//
// Solidity: function simulatePerformUpkeep(uint256 id, bytes performData) returns(bool success, uint256 gasUsed)
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.SimulatePerformUpkeep(&_KeeperRegistry20.TransactOpts, id, performData)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferOwnership(&_KeeperRegistry20.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferOwnership(&_KeeperRegistry20.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferPayeeship(&_KeeperRegistry20.TransactOpts, transmitter, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.TransferUpkeepAdmin(&_KeeperRegistry20.TransactOpts, id, proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Transmit(&_KeeperRegistry20.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes rawReport, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Transmit(&_KeeperRegistry20.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Unpause(&_KeeperRegistry20.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Unpause(&_KeeperRegistry20.TransactOpts)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "unpauseUpkeep", id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UnpauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UnpauseUpkeep(&_KeeperRegistry20.TransactOpts, id)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UpdateCheckData(&_KeeperRegistry20.TransactOpts, id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.UpdateCheckData(&_KeeperRegistry20.TransactOpts, id, newCheckData)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawFunds", id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawFunds(&_KeeperRegistry20.TransactOpts, id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawFunds(&_KeeperRegistry20.TransactOpts, id, to)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawOwnerFunds")
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawOwnerFunds(&_KeeperRegistry20.TransactOpts)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawOwnerFunds(&_KeeperRegistry20.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.Transact(opts, "withdrawPayment", from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawPayment(&_KeeperRegistry20.TransactOpts, from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.WithdrawPayment(&_KeeperRegistry20.TransactOpts, from, to)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Fallback(&_KeeperRegistry20.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Fallback(&_KeeperRegistry20.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry20.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20Session) Receive() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Receive(&_KeeperRegistry20.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry20 *KeeperRegistry20TransactorSession) Receive() (*types.Transaction, error) {
	return _KeeperRegistry20.Contract.Receive(&_KeeperRegistry20.TransactOpts)
}

// KeeperRegistry20CancelledUpkeepReportIterator is returned from FilterCancelledUpkeepReport and is used to iterate over the raw logs and unpacked data for CancelledUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20CancelledUpkeepReportIterator struct {
	Event *KeeperRegistry20CancelledUpkeepReport // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20CancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20CancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20CancelledUpkeepReport represents a CancelledUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20CancelledUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCancelledUpkeepReport is a free log retrieval operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20CancelledUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchCancelledUpkeepReport is a free log subscription operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20CancelledUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancelledUpkeepReport is a log parse operation binding the contract event 0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f.
//
// Solidity: event CancelledUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistry20CancelledUpkeepReport, error) {
	event := new(KeeperRegistry20CancelledUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20ConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20ConfigSetIterator struct {
	Event *KeeperRegistry20ConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20ConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20ConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20ConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20ConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20ConfigSet represents a ConfigSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20ConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry20ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20ConfigSetIterator{contract: _KeeperRegistry20.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20ConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20ConfigSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseConfigSet(log types.Log) (*KeeperRegistry20ConfigSet, error) {
	event := new(KeeperRegistry20ConfigSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20FundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsAddedIterator struct {
	Event *KeeperRegistry20FundsAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20FundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20FundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20FundsAdded represents a FundsAdded event raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry20FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20FundsAddedIterator{contract: _KeeperRegistry20.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20FundsAdded)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFundsAdded is a log parse operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistry20FundsAdded, error) {
	event := new(KeeperRegistry20FundsAdded)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20FundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsWithdrawnIterator struct {
	Event *KeeperRegistry20FundsWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20FundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20FundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20FundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20FundsWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20FundsWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFundsWithdrawn is a log parse operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistry20FundsWithdrawn, error) {
	event := new(KeeperRegistry20FundsWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20InsufficientFundsUpkeepReportIterator is returned from FilterInsufficientFundsUpkeepReport and is used to iterate over the raw logs and unpacked data for InsufficientFundsUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20InsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistry20InsufficientFundsUpkeepReport // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20InsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20InsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20InsufficientFundsUpkeepReport represents a InsufficientFundsUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20InsufficientFundsUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterInsufficientFundsUpkeepReport is a free log retrieval operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20InsufficientFundsUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchInsufficientFundsUpkeepReport is a free log subscription operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20InsufficientFundsUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInsufficientFundsUpkeepReport is a log parse operation binding the contract event 0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96.
//
// Solidity: event InsufficientFundsUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistry20InsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistry20InsufficientFundsUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistry20OwnerFundsWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry20OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnerFundsWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnerFundsWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerFundsWithdrawn is a log parse operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry20OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistry20OwnerFundsWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistry20OwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnershipTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnershipTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry20OwnershipTransferRequested, error) {
	event := new(KeeperRegistry20OwnershipTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferredIterator struct {
	Event *KeeperRegistry20OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20OwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20OwnershipTransferredIterator{contract: _KeeperRegistry20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20OwnershipTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistry20OwnershipTransferred, error) {
	event := new(KeeperRegistry20OwnershipTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PausedIterator struct {
	Event *KeeperRegistry20Paused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Paused represents a Paused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry20PausedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PausedIterator{contract: _KeeperRegistry20.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Paused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePaused(log types.Log) (*KeeperRegistry20Paused, error) {
	event := new(KeeperRegistry20Paused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PayeesUpdatedIterator is returned from FilterPayeesUpdated and is used to iterate over the raw logs and unpacked data for PayeesUpdated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeesUpdatedIterator struct {
	Event *KeeperRegistry20PayeesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20PayeesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20PayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20PayeesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeesUpdated represents a PayeesUpdated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPayeesUpdated is a free log retrieval operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistry20PayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeesUpdatedIterator{contract: _KeeperRegistry20.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

// WatchPayeesUpdated is a free log subscription operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeesUpdated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeesUpdated is a log parse operation binding the contract event 0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725.
//
// Solidity: event PayeesUpdated(address[] transmitters, address[] payees)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistry20PayeesUpdated, error) {
	event := new(KeeperRegistry20PayeesUpdated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistry20PayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeeshipTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry20PayeeshipTransferRequested, error) {
	event := new(KeeperRegistry20PayeeshipTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferredIterator struct {
	Event *KeeperRegistry20PayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20PayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20PayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry20PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PayeeshipTransferredIterator{contract: _KeeperRegistry20.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PayeeshipTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry20PayeeshipTransferred, error) {
	event := new(KeeperRegistry20PayeeshipTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20PaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistry20 contract.
type KeeperRegistry20PaymentWithdrawnIterator struct {
	Event *KeeperRegistry20PaymentWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20PaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20PaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20PaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistry20 contract.
type KeeperRegistry20PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry20PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20PaymentWithdrawnIterator{contract: _KeeperRegistry20.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20PaymentWithdrawn)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentWithdrawn is a log parse operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry20PaymentWithdrawn, error) {
	event := new(KeeperRegistry20PaymentWithdrawn)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20ReorgedUpkeepReportIterator is returned from FilterReorgedUpkeepReport and is used to iterate over the raw logs and unpacked data for ReorgedUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20ReorgedUpkeepReportIterator struct {
	Event *KeeperRegistry20ReorgedUpkeepReport // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20ReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20ReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20ReorgedUpkeepReport represents a ReorgedUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20ReorgedUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterReorgedUpkeepReport is a free log retrieval operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20ReorgedUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchReorgedUpkeepReport is a free log subscription operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20ReorgedUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReorgedUpkeepReport is a log parse operation binding the contract event 0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13.
//
// Solidity: event ReorgedUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistry20ReorgedUpkeepReport, error) {
	event := new(KeeperRegistry20ReorgedUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20StaleUpkeepReportIterator is returned from FilterStaleUpkeepReport and is used to iterate over the raw logs and unpacked data for StaleUpkeepReport events raised by the KeeperRegistry20 contract.
type KeeperRegistry20StaleUpkeepReportIterator struct {
	Event *KeeperRegistry20StaleUpkeepReport // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20StaleUpkeepReportIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20StaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20StaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20StaleUpkeepReportIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20StaleUpkeepReport represents a StaleUpkeepReport event raised by the KeeperRegistry20 contract.
type KeeperRegistry20StaleUpkeepReport struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterStaleUpkeepReport is a free log retrieval operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20StaleUpkeepReportIterator{contract: _KeeperRegistry20.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

// WatchStaleUpkeepReport is a free log subscription operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20StaleUpkeepReport)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStaleUpkeepReport is a log parse operation binding the contract event 0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89.
//
// Solidity: event StaleUpkeepReport(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistry20StaleUpkeepReport, error) {
	event := new(KeeperRegistry20StaleUpkeepReport)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20TransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the KeeperRegistry20 contract.
type KeeperRegistry20TransmittedIterator struct {
	Event *KeeperRegistry20Transmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20TransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Transmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20Transmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20TransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20TransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Transmitted represents a Transmitted event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Transmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistry20TransmittedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20TransmittedIterator{contract: _KeeperRegistry20.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Transmitted) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Transmitted)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseTransmitted(log types.Log) (*KeeperRegistry20Transmitted, error) {
	event := new(KeeperRegistry20Transmitted)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UnpausedIterator struct {
	Event *KeeperRegistry20Unpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20Unpaused represents a Unpaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry20UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UnpausedIterator{contract: _KeeperRegistry20.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20Unpaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUnpaused(log types.Log) (*KeeperRegistry20Unpaused, error) {
	event := new(KeeperRegistry20Unpaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepAdminTransferRequestedIterator is returned from FilterUpkeepAdminTransferRequested and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferRequested events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistry20UpkeepAdminTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepAdminTransferRequested represents a UpkeepAdminTransferRequested event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferRequested is a free log retrieval operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry20UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferRequested is a free log subscription operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepAdminTransferRequested)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepAdminTransferRequested is a log parse operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistry20UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistry20UpkeepAdminTransferRequested)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepAdminTransferredIterator is returned from FilterUpkeepAdminTransferred and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferred events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistry20UpkeepAdminTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepAdminTransferred represents a UpkeepAdminTransferred event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferred is a free log retrieval operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry20UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepAdminTransferredIterator{contract: _KeeperRegistry20.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferred is a free log subscription operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepAdminTransferred)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepAdminTransferred is a log parse operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistry20UpkeepAdminTransferred, error) {
	event := new(KeeperRegistry20UpkeepAdminTransferred)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCanceledIterator struct {
	Event *KeeperRegistry20UpkeepCanceled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry20UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepCanceledIterator{contract: _KeeperRegistry20.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepCanceled)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepCanceled is a log parse operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistry20UpkeepCanceled, error) {
	event := new(KeeperRegistry20UpkeepCanceled)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepCheckDataUpdatedIterator is returned from FilterUpkeepCheckDataUpdated and is used to iterate over the raw logs and unpacked data for UpkeepCheckDataUpdated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistry20UpkeepCheckDataUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepCheckDataUpdated represents a UpkeepCheckDataUpdated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCheckDataUpdated is a free log retrieval operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepCheckDataUpdatedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

// WatchUpkeepCheckDataUpdated is a free log subscription operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepCheckDataUpdated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepCheckDataUpdated is a log parse operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistry20UpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistry20UpkeepCheckDataUpdated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistry20UpkeepGasLimitSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepGasLimitSetIterator{contract: _KeeperRegistry20.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepGasLimitSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepGasLimitSet is a log parse operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry20UpkeepGasLimitSet, error) {
	event := new(KeeperRegistry20UpkeepGasLimitSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepMigratedIterator struct {
	Event *KeeperRegistry20UpkeepMigrated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepMigratedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepMigrated)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepMigrated is a log parse operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistry20UpkeepMigrated, error) {
	event := new(KeeperRegistry20UpkeepMigrated)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepOffchainConfigSetIterator is returned from FilterUpkeepOffchainConfigSet and is used to iterate over the raw logs and unpacked data for UpkeepOffchainConfigSet events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistry20UpkeepOffchainConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepOffchainConfigSet represents a UpkeepOffchainConfigSet event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpkeepOffchainConfigSet is a free log retrieval operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepOffchainConfigSetIterator{contract: _KeeperRegistry20.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepOffchainConfigSet is a free log subscription operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepOffchainConfigSet)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepOffchainConfigSet is a log parse operation binding the contract event 0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850.
//
// Solidity: event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistry20UpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistry20UpkeepOffchainConfigSet)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepPausedIterator is returned from FilterUpkeepPaused and is used to iterate over the raw logs and unpacked data for UpkeepPaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPausedIterator struct {
	Event *KeeperRegistry20UpkeepPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepPaused represents a UpkeepPaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPaused is a free log retrieval operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepPausedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepPaused is a free log subscription operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepPaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepPaused is a log parse operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistry20UpkeepPaused, error) {
	event := new(KeeperRegistry20UpkeepPaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPerformedIterator struct {
	Event *KeeperRegistry20UpkeepPerformed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepPerformed struct {
	Id               *big.Int
	Success          bool
	CheckBlockNumber uint32
	GasUsed          *big.Int
	GasOverhead      *big.Int
	TotalPayment     *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPerformed is a free log retrieval operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistry20UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepPerformedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepPerformed)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepPerformed is a log parse operation binding the contract event 0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, uint32 checkBlockNumber, uint256 gasUsed, uint256 gasOverhead, uint96 totalPayment)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistry20UpkeepPerformed, error) {
	event := new(KeeperRegistry20UpkeepPerformed)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepReceivedIterator struct {
	Event *KeeperRegistry20UpkeepReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepReceivedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepReceived)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepReceived is a log parse operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistry20UpkeepReceived, error) {
	event := new(KeeperRegistry20UpkeepReceived)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepRegisteredIterator struct {
	Event *KeeperRegistry20UpkeepRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepRegisteredIterator{contract: _KeeperRegistry20.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepRegistered)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepRegistered is a log parse operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistry20UpkeepRegistered, error) {
	event := new(KeeperRegistry20UpkeepRegistered)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry20UpkeepUnpausedIterator is returned from FilterUpkeepUnpaused and is used to iterate over the raw logs and unpacked data for UpkeepUnpaused events raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepUnpausedIterator struct {
	Event *KeeperRegistry20UpkeepUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistry20UpkeepUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry20UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry20UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistry20UpkeepUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry20UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry20UpkeepUnpaused represents a UpkeepUnpaused event raised by the KeeperRegistry20 contract.
type KeeperRegistry20UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepUnpaused is a free log retrieval operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry20UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry20UpkeepUnpausedIterator{contract: _KeeperRegistry20.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepUnpaused is a free log subscription operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry20UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry20.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry20UpkeepUnpaused)
				if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepUnpaused is a log parse operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry20 *KeeperRegistry20Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistry20UpkeepUnpaused, error) {
	event := new(KeeperRegistry20UpkeepUnpaused)
	if err := _KeeperRegistry20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
