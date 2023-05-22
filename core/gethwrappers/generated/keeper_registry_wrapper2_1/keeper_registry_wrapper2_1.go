// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper2_1

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

type OnchainConfig struct {
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

type State struct {
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

type UpkeepInfo struct {
	Target                 common.Address
	Forwarder              common.Address
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

var KeeperRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"trigger\",\"type\":\"uint8\"}],\"name\":\"getActiveUpkeepIDsByType\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structOnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structUpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_next\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"forwarderAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b506040516200ed023803806200ed02833981016040819052620000359162000246565b8383838360405162000047906200021b565b620000569493929190620002a8565b604051809103906000f08015801562000073573d6000803e3d6000fd5b50848484843380600081620000cf5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001025762000102816200016f565b5050508360028111156200011a576200011a620002f2565b60e0816002811115620001315762000131620002f2565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b9091166101005250620003089350505050565b6001600160a01b038116331415620001ca5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000c6565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b618dc78062005f3b83390190565b80516001600160a01b03811681146200024157600080fd5b919050565b600080600080608085870312156200025d57600080fd5b8451600381106200026d57600080fd5b93506200027d6020860162000229565b92506200028d6040860162000229565b91506200029d6060860162000229565b905092959194509250565b6080810160038610620002cb57634e487b7160e01b600052602160045260246000fd5b9481526001600160a01b039384166020820152918316604083015290911660609091015290565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c615ba762000394600039600081816101740152818161020f015261028b01526000818161025701528181612a0701528181612ad401528181613e390152614006015260006102d50152600061046f0152600081816104ad01528181610ae30152610bdb0152615ba76000f3fe60806040523480156200001157600080fd5b5060043610620001725760003560e01c80638da5cb5b11620000e1578063b10b673c116200008d578063ca30e603116200006f578063ca30e60314620004ab578063e3d0e71214620004d2578063f2fde38b14620004e95762000172565b8063b10b673c146200046d578063b1dc65a414620004945762000172565b8063a4c0ed3611620000c3578063a4c0ed3614620003dd578063aed2e92914620003f4578063afcb95d714620004235762000172565b80638da5cb5b14620003a7578063948108f714620003c65762000172565b80635147cd5911620001415780636ded9eae11620001235780636ded9eae14620002fa57806379ba5097146200033c57806381ff704814620003485762000172565b80635147cd5914620002ad5780636709d0e514620002d35762000172565b8063181f5a7714620001b8578063349e8cca146200020d5780634b4fd03b14620002555780634ff597c11462000285575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015620001b3573d6000f35b3d6000fd5b620001f56040518060400160405280601481526020017f4b6565706572526567697374727920322e312e3000000000000000000000000081525081565b60405162000204919062004f25565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000204565b7f000000000000000000000000000000000000000000000000000000000000000060405162000204919062004f3a565b6200022f7f000000000000000000000000000000000000000000000000000000000000000081565b620002c4620002be36600462004c3f565b50600090565b60405162000204919062004f57565b7f00000000000000000000000000000000000000000000000000000000000000006200022f565b620003116200030b36600462004845565b62000500565b6040805192835273ffffffffffffffffffffffffffffffffffffffff90911660208301520162000204565b6200034662000808565b005b62000383601254600e5463ffffffff6c0100000000000000000000000083048116937001000000000000000000000000000000009093041691565b6040805163ffffffff94851681529390921660208401529082015260600162000204565b60005473ffffffffffffffffffffffffffffffffffffffff166200022f565b62000346620003d736600462004db3565b6200090b565b62000346620003ee366004620047e3565b62000bc3565b6200040b6200040536600462004c59565b62000de7565b60408051921515835260208301919091520162000204565b600e54600f54604080516000815260208101939093527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169082015260600162000204565b7f00000000000000000000000000000000000000000000000000000000000000006200022f565b62000346620004a5366004620049d8565b62000f64565b7f00000000000000000000000000000000000000000000000000000000000000006200022f565b62000346620004e3366004620048fa565b62001bb3565b62000346620004fa366004620047c3565b620029e9565b6000806200052360005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156200058757506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b15620005bf576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620005e16001620005cf62002a01565b620005db9190620053be565b62002ace565b601254604080516020810193909352309083015268010000000000000000900463ffffffff1660608201526080016040516020818303038152906040528051906020012060001c91506000896040516200063b9062004322565b73ffffffffffffffffffffffffffffffffffffffff9091168152602001604051809103906000f08015801562000675573d6000803e3d6000fd5b50905062000708838b8b8b60008c8c8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505060008c8c8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c925062002c5e915050565b6012805468010000000000000000900463ffffffff169060086200072c83620054cd565b825463ffffffff9182166101009390930a92830291909202199091161790555060008381526017602052604090206200076790868662004330565b506040805163ffffffff8b16815273ffffffffffffffffffffffffffffffffffffffff8a16602082015284917fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012910160405180910390a2827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508686604051620007f292919062004ed8565b60405180910390a2905097509795505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146200088f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600082815260046020908152604091829020825161010081018452815463ffffffff80821683526401000000008204811694830185905260ff6801000000000000000083041615159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004851660c083015260029092015490911660e0820152911462000a1e576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a0015162000a30919062005286565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560155462000a989184169062005221565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b15801562000b3d57600080fd5b505af115801562000b52573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000b78919062004a9a565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161462000c33576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811462000c6e576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600062000c7e8284018462004c3f565b600081815260046020526040902054909150640100000000900463ffffffff9081161462000cd8576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090206001015462000d159085906c0100000000000000000000000090046bffffffffffffffffffffffff1662005286565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff90921691909117905560155462000d8290859062005221565b6015556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008062000df46200312b565b600f546e010000000000000000000000000000900460ff161562000e44576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825161010081018452815463ffffffff8082168352640100000000820481168386015260ff6801000000000000000083041615158387015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606084015260018401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c083015260029092015490911660e08201528251601f87018390048302810183019093528583529162000f579183919088908890819084018382808284376000920191909152506200316692505050565b9250925050935093915050565b60005a6040805161012081018252600f5460ff808216835261010080830463ffffffff90811660208601526501000000000084048116958501959095526901000000000000000000830462ffffff1660608501526c01000000000000000000000000830461ffff1660808501526e0100000000000000000000000000008304821615801560a08601526f010000000000000000000000000000008404909216151560c085015270010000000000000000000000000000000083046bffffffffffffffffffffffff1660e08501527c01000000000000000000000000000000000000000000000000000000009092049093169082015291925062001093576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3360009081526008602052604090205460ff16620010dd576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000620011208a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200338992505050565b9050600081604001515167ffffffffffffffff811115620011455762001145620055b0565b6040519080825280602002602001820160405280156200120357816020015b604080516101c081018252600060c0820181815260e08301829052610100830182905261012083018290526101408301829052610160830182905261018083018290526101a0830182905282526020808301829052928201819052606082018190526080820181905260a082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620011645790505b5090506000805b836040015151811015620014e457600460008560400151838151811062001235576200123562005581565b6020908102919091018101518252818101929092526040908101600020815161010081018352815463ffffffff8082168352640100000000820481169583019590955260ff6801000000000000000082041615159382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152600201541660e082015283518490839081106200132f576200132f62005581565b6020026020010151600001819052506200139c8584838151811062001358576200135862005581565b602002602001015160000151600001518660800151848151811062001381576200138162005581565b6020026020010151518760000151886020015160016200345c565b838281518110620013b157620013b162005581565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff16815250506200146d84604001518281518110620013fc57620013fc62005581565b6020026020010151856060015183815181106200141d576200141d62005581565b60200260200101518584815181106200143a576200143a62005581565b6020026020010151600001518685815181106200145b576200145b62005581565b602002602001015160400151620034ab565b83828151811062001482576200148262005581565b60200260200101516020019015159081151581525050828181518110620014ad57620014ad62005581565b60200260200101516020015115620014cf57620014cc600183620051f8565b91505b80620014db8162005491565b9150506200120a565b5061ffff8116620014fa57505050505062001ba9565b600e548d351462001537576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8351620015469060016200525e565b60ff1689141580620015585750888714155b1562001590576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620015a28d8d8d8d8d8d8d8d620035a5565b60005b836040015151811015620017ae57828181518110620015c857620015c862005581565b602002602001015160200151156200179957620015e462002a01565b63ffffffff16600460008660400151848151811062001607576200160762005581565b6020026020010151815260200190815260200160002060010160189054906101000a900463ffffffff1663ffffffff16141562001670576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620016bc83828151811062001689576200168962005581565b60200260200101516000015185608001518381518110620016ae57620016ae62005581565b602002602001015162003166565b848381518110620016d157620016d162005581565b6020026020010151606001858481518110620016f157620016f162005581565b602002602001015160800182815250821515151581525050508281815181106200171f576200171f62005581565b60200260200101516080015186620017389190620053be565b95506200174462002a01565b600460008660400151848151811062001761576200176162005581565b6020026020010151815260200190815260200160002060010160186101000a81548163ffffffff021916908363ffffffff1602179055505b80620017a58162005491565b915050620015a5565b508351620017be9060016200525e565b620017cf9060ff1661044c620052f5565b616914620017df8d6010620052f5565b5a620017ec9089620053be565b620017f8919062005221565b62001804919062005221565b62001810919062005221565b94506116a86200182561ffff831687620052b0565b62001831919062005221565b945060008060008060005b87604001515181101562001a43578681815181106200185f576200185f62005581565b6020026020010151602001511562001a2e57620018a38a896080015183815181106200188f576200188f62005581565b6020026020010151518b6000015162003825565b878281518110620018b857620018b862005581565b602002602001015160a00181815250506200191c8989604001518381518110620018e657620018e662005581565b602002602001015189848151811062001903576200190362005581565b60200260200101518b600001518c602001518b62003846565b90935091506200192d828562005286565b93506200193b838662005286565b945086818151811062001952576200195262005581565b60200260200101516060015115158860400151828151811062001979576200197962005581565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b8486620019b0919062005286565b8a8581518110620019c557620019c562005581565b6020026020010151608001518b8681518110620019e657620019e662005581565b602002602001015160a001518d60600151878151811062001a0b5762001a0b62005581565b602002602001015160405162001a259493929190620050ca565b60405180910390a35b8062001a3a8162005491565b9150506200183c565b5050336000908152600860205260409020805484925060029062001a7d9084906201000090046bffffffffffffffffffffffff1662005286565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600f60000160108282829054906101000a90046bffffffffffffffffffffffff1662001ad9919062005286565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f60016003811062001b1f5762001b1f62005581565b602002013560001c9050600060088264ffffffffff16901c905087610100015163ffffffff168163ffffffff16111562001b9f57600f80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b62001bbd62003941565b601f8651111562001bfa576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff841662001c35576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8451865114158062001c57575062001c4f84600362005364565b60ff16865111155b1562001c8f576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b547001000000000000000000000000000000009091046bffffffffffffffffffffffff169060005b816bffffffffffffffffffffffff1681101562001d2d5762001d17600b828154811062001ced5762001ced62005581565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484620039c4565b508062001d248162005491565b91505062001cbc565b5060008060005b836bffffffffffffffffffffffff1681101562001e4057600a818154811062001d615762001d6162005581565b600091825260209091200154600b805473ffffffffffffffffffffffffffffffffffffffff9092169450908290811062001d9f5762001d9f62005581565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff868116845260098352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600890925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508062001e378162005491565b91505062001d34565b5062001e4f600a6000620043dd565b62001e5d600b6000620043dd565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015620021f057600960008e838151811062001ea65762001ea662005581565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161562001f12576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600960008f848151811062001f465762001f4662005581565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c908290811062001ff15762001ff162005581565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff81166000908152600883526040908190208151608081018352905460ff80821615801584526101008304909116958301959095526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015294509250620020b7576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff871660009081526008909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580620021e78162005491565b91505062001e83565b50508a51620022089150600a9060208d0190620043fd565b5088516200221e90600b9060208c0190620043fd565b5060008780602001905181019062002237919062004b36565b60125460c082015191925063ffffffff6401000000009091048116911610156200228d576040517f39abc10400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125460e082015163ffffffff74010000000000000000000000000000000000000000909204821691161015620022f0576040517f1fa9bdcb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125461010082015163ffffffff780100000000000000000000000000000000000000000000000090920482169116101562002358576040517fd1d5faa800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061012001604052808a60ff168152602001826000015163ffffffff168152602001826020015163ffffffff168152602001826060015162ffffff168152602001826080015161ffff168152602001600015158152602001600015158152602001866bffffffffffffffffffffffff168152602001600063ffffffff16815250600f60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160096101000a81548162ffffff021916908362ffffff160217905550608082015181600001600c6101000a81548161ffff021916908361ffff16021790555060a082015181600001600e6101000a81548160ff02191690831515021790555060c082015181600001600f6101000a81548160ff02191690831515021790555060e08201518160000160106101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061010082015181600001601c6101000a81548163ffffffff021916908363ffffffff1602179055509050506040518061016001604052808260a001516bffffffffffffffffffffffff16815260200182610160015173ffffffffffffffffffffffffffffffffffffffff168152602001601060010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200182610180015173ffffffffffffffffffffffffffffffffffffffff168152602001826040015163ffffffff1681526020018260c0015163ffffffff168152602001601060020160089054906101000a900463ffffffff1663ffffffff1681526020016010600201600c9054906101000a900463ffffffff1663ffffffff168152602001601060020160109054906101000a900463ffffffff1663ffffffff1681526020018260e0015163ffffffff16815260200182610100015163ffffffff16815250601060008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160020160046101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600201600c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160106101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160146101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160186101000a81548163ffffffff021916908363ffffffff1602179055509050508061012001516013819055508061014001516014819055506000601060020160109054906101000a900463ffffffff169050620028c062002a01565b601280547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000063ffffffff9384160217808255600192600c91620029299185916c010000000000000000000000009004166200523c565b92506101000a81548163ffffffff021916908363ffffffff1602179055506200297546306010600201600c9054906101000a900463ffffffff1663ffffffff168f8f8f8f8f8f62003bf7565b600e819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600e546010600201600c9054906101000a900463ffffffff168f8f8f8f8f8f604051620029d3999897969594939291906200501b565b60405180910390a1505050505050505050505050565b620029f362003941565b620029fe8162003ca3565b50565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562002a3a5762002a3a62005552565b141562002ac957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801562002a8957600080fd5b505afa15801562002a9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002ac4919062004abe565b905090565b504390565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562002b075762002b0762005552565b141562002c54576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801562002b5857600080fd5b505afa15801562002b6d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b93919062004abe565b9050808310158062002bb1575061010062002baf8483620053be565b115b1562002bc05750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a829060240160206040518083038186803b15801562002c1257600080fd5b505afa15801562002c27573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002c4d919062004abe565b9392505050565b504090565b919050565b600f546e010000000000000000000000000000900460ff161562002cae576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff88163b62002cfd576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125484517401000000000000000000000000000000000000000090910463ffffffff16101562002d5a576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8763ffffffff16108062002d84575060125463ffffffff6401000000009091048116908816115b1562002dbc576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008981526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562002e1c576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061010001604052808863ffffffff16815260200163ffffffff8016815260200184151581526020018273ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001866bffffffffffffffffffffffff168152602001600063ffffffff1681526020018973ffffffffffffffffffffffffffffffffffffffff16815250600460008b815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505085600560008b815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550846bffffffffffffffffffffffff16601554620030cb919062005221565b60155560008981526007602090815260409091208551620030ef928701906200447a565b506000898152601760209081526040909120835162003111928501906200447a565b506200311f60028a62003d9b565b50505050505050505050565b321562003164576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b600f5460009081906f01000000000000000000000000000000900460ff1615620031bc576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f010000000000000000000000000000001790555a90506000634585e33b60e01b8460405160240162003216919062004f25565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909416939093179092526060870151875192517f79188d1600000000000000000000000000000000000000000000000000000000815291935073ffffffffffffffffffffffffffffffffffffffff16916379188d1691620032f091908590600401620050a9565b602060405180830381600087803b1580156200330b57600080fd5b505af115801562003320573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003346919062004a9a565b92505a620033559083620053be565b915050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff16905590939092509050565b620033bc6040518060a0016040528060008152602001600081526020016060815260200160608152602001606081525090565b600080600080600086806020019051810190620033da919062004ca9565b9450945094509450945081518351141580620033f857508051835114155b1562003430576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160a0810182529586526020860194909452928401919091526060830152608082015292915050565b6000806200346f86896000015162003db2565b90506000806200348c8a8a63ffffffff16858a8a60018b62003e00565b90925090506200349d818362005286565b9a9950505050505050505050565b600080600085806020019051810190620034c6919062004ad8565b9050620034d587828762004221565b620034e6576000925050506200359d565b50620034f162002a01565b846020015163ffffffff1611620035385760405186907fd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f90600090a260009150506200359d565b826bffffffffffffffffffffffff168460a001516bffffffffffffffffffffffff161015620035975760405186907f7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb9690600090a260009150506200359d565b60019150505b949350505050565b60008787604051620035b992919062004e85565b604051908190038120620035d2918b9060200162004ebe565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b88811015620037bb5760018587836020811062003642576200364262005581565b6200365191901a601b6200525e565b8c8c8581811062003666576200366662005581565b905060200201358b8b8681811062003682576200368262005581565b9050602002013560405160008152602001604052604051620036c0949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015620036e3573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526009602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905062003792576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b840193508080620037b29062005491565b91505062003621565b50827e0101010101010101010101010101010101010101010101010101010101010184161462003817576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b600062003833838362003db2565b90508084101562002c4d57509192915050565b600080620038638887608001518860a00151888888600162003e00565b9092509050600062003876828462005286565b600089815260046020526040902060010180549192508291600c90620038bc9084906c0100000000000000000000000090046bffffffffffffffffffffffff16620053d8565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092620039079185911662005286565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003164576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000886565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e0100000000000000000000000000009092041660608201819052829062003a519086620053d8565b9050600062003a618583620052c7565b9050808360400181815162003a77919062005286565b6bffffffffffffffffffffffff908116909152871660608501525062003a9e858262005390565b62003aaa9083620053d8565b6011805460009062003acc9084906bffffffffffffffffffffffff1662005286565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b6000808a8a8a8a8a8a8a8a8a60405160200162003c1d9998979695949392919062004f6e565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff811633141562003d25576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000886565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062003da98383620042d0565b90505b92915050565b600062003dc763ffffffff84166014620052f5565b62003dd48360016200525e565b62003de59060ff16611d4c620052f5565b62003df490620124f862005221565b62003da9919062005221565b6000806000896080015161ffff168762003e1b9190620052f5565b905083801562003e2a5750803a105b1562003e3357503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003e6c5762003e6c62005552565b14156200400257604080516000815260208101909152851562003ed15760003660405180608001604052806048815260200162005b536048913960405160200162003eba9392919062004e95565b604051602081830303815290604052905062003f53565b60125462003f03907801000000000000000000000000000000000000000000000000900463ffffffff16600462005335565b63ffffffff1667ffffffffffffffff81111562003f245762003f24620055b0565b6040519080825280601f01601f19166020018201604052801562003f4f576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9062003fa590849060040162004f25565b60206040518083038186803b15801562003fbe57600080fd5b505afa15801562003fd3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ff9919062004abe565b915050620040c6565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562004039576200403962005552565b1415620040c657606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200408857600080fd5b505afa1580156200409d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620040c3919062004abe565b90505b84620040e557808b6080015161ffff16620040e29190620052f5565b90505b620040f561ffff871682620052b0565b905060008782620041078c8e62005221565b620041139086620052f5565b6200411f919062005221565b6200413390670de0b6b3a7640000620052f5565b6200413f9190620052b0565b905060008c6040015163ffffffff1664e8d4a51000620041609190620052f5565b898e6020015163ffffffff16858f886200417b9190620052f5565b62004187919062005221565b6200419790633b9aca00620052f5565b620041a39190620052f5565b620041af9190620052b0565b620041bb919062005221565b90506b033b2e3c9fd0803ce8000000620041d6828462005221565b11156200420f576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60008160c0015163ffffffff16836000015163ffffffff161015620042745760405184907f5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a8990600090a250600062002c4d565b602083015183516200428c9063ffffffff1662002ace565b14620042c65760405184907f561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc1390600090a250600062002c4d565b5060019392505050565b6000818152600183016020526040812054620043195750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562003dac565b50600062003dac565b610522806200563183390190565b8280546200433e906200543b565b90600052602060002090601f016020900481019282620043625760008555620043cb565b82601f106200439b578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555620043cb565b82800160010185558215620043cb579182015b82811115620043cb578235825591602001919060010190620043ae565b50620043d9929150620044f7565b5090565b5080546000825590600052602060002090810190620029fe9190620044f7565b828054828255906000526020600020908101928215620043cb579160200282015b82811115620043cb57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906200441e565b82805462004488906200543b565b90600052602060002090601f016020900481019282620044ac5760008555620043cb565b82601f10620044c757805160ff1916838001178555620043cb565b82800160010185558215620043cb579182015b82811115620043cb578251825591602001919060010190620044da565b5b80821115620043d95760008155600101620044f8565b805162002c5981620055df565b600082601f8301126200452d57600080fd5b8135602062004546620045408362005188565b62005136565b80838252828201915082860187848660051b89010111156200456757600080fd5b60005b85811015620045935781356200458081620055df565b845292840192908401906001016200456a565b5090979650505050505050565b60008083601f840112620045b357600080fd5b50813567ffffffffffffffff811115620045cc57600080fd5b6020830191508360208260051b8501011115620045e857600080fd5b9250929050565b600082601f8301126200460157600080fd5b8151602062004614620045408362005188565b80838252828201915082860187848660051b89010111156200463557600080fd5b60005b858110156200459357815167ffffffffffffffff8111156200465957600080fd5b8801603f81018a136200466b57600080fd5b858101516040620046806200454083620051af565b8281528c828486010111156200469557600080fd5b620046a6838a830184870162005408565b8752505050928401929084019060010162004638565b60008083601f840112620046cf57600080fd5b50813567ffffffffffffffff811115620046e857600080fd5b602083019150836020828501011115620045e857600080fd5b600082601f8301126200471357600080fd5b8135620047246200454082620051af565b8181528460208386010111156200473a57600080fd5b816020850160208301376000918101602001919091529392505050565b805161ffff8116811462002c5957600080fd5b805162ffffff8116811462002c5957600080fd5b805162002c598162005602565b803567ffffffffffffffff8116811462002c5957600080fd5b803560ff8116811462002c5957600080fd5b805162002c598162005615565b600060208284031215620047d657600080fd5b813562002c4d81620055df565b60008060008060608587031215620047fa57600080fd5b84356200480781620055df565b935060208501359250604085013567ffffffffffffffff8111156200482b57600080fd5b6200483987828801620046bc565b95989497509550505050565b600080600080600080600060a0888a0312156200486157600080fd5b87356200486e81620055df565b96506020880135620048808162005602565b955060408801356200489281620055df565b9450606088013567ffffffffffffffff80821115620048b057600080fd5b620048be8b838c01620046bc565b909650945060808a0135915080821115620048d857600080fd5b50620048e78a828b01620046bc565b989b979a50959850939692959293505050565b60008060008060008060c087890312156200491457600080fd5b863567ffffffffffffffff808211156200492d57600080fd5b6200493b8a838b016200451b565b975060208901359150808211156200495257600080fd5b620049608a838b016200451b565b96506200497060408a01620047a4565b955060608901359150808211156200498757600080fd5b620049958a838b0162004701565b9450620049a560808a016200478b565b935060a0890135915080821115620049bc57600080fd5b50620049cb89828a0162004701565b9150509295509295509295565b60008060008060008060008060e0898b031215620049f557600080fd5b606089018a81111562004a0757600080fd5b8998503567ffffffffffffffff8082111562004a2257600080fd5b62004a308c838d01620046bc565b909950975060808b013591508082111562004a4a57600080fd5b62004a588c838d01620045a0565b909750955060a08b013591508082111562004a7257600080fd5b5062004a818b828c01620045a0565b999c989b50969995989497949560c00135949350505050565b60006020828403121562004aad57600080fd5b8151801515811462002c4d57600080fd5b60006020828403121562004ad157600080fd5b5051919050565b60006040828403121562004aeb57600080fd5b6040516040810181811067ffffffffffffffff8211171562004b115762004b11620055b0565b604052825162004b218162005602565b81526020928301519281019290925250919050565b60006101a0828403121562004b4a57600080fd5b62004b5462005109565b62004b5f836200477e565b815262004b6f602084016200477e565b602082015262004b82604084016200477e565b604082015262004b95606084016200476a565b606082015262004ba86080840162004757565b608082015262004bbb60a08401620047b6565b60a082015262004bce60c084016200477e565b60c082015262004be160e084016200477e565b60e082015261010062004bf68185016200477e565b908201526101208381015190820152610140808401519082015261016062004c208185016200450e565b9082015261018062004c348482016200450e565b908201529392505050565b60006020828403121562004c5257600080fd5b5035919050565b60008060006040848603121562004c6f57600080fd5b83359250602084013567ffffffffffffffff81111562004c8e57600080fd5b62004c9c86828701620046bc565b9497909650939450505050565b600080600080600060a0868803121562004cc257600080fd5b855194506020808701519450604087015167ffffffffffffffff8082111562004cea57600080fd5b818901915089601f83011262004cff57600080fd5b815162004d10620045408262005188565b8082825285820191508585018d878560051b880101111562004d3157600080fd5b600095505b8386101562004d5657805183526001959095019491860191860162004d36565b5060608c0151909850945050508083111562004d7157600080fd5b62004d7f8a848b01620045ef565b9450608089015192508083111562004d9657600080fd5b505062004da688828901620045ef565b9150509295509295909350565b6000806040838503121562004dc757600080fd5b82359150602083013562004ddb8162005615565b809150509250929050565b600081518084526020808501945080840160005b8381101562004e2e57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004dfa565b509495945050505050565b6000815180845262004e5381602086016020860162005408565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8183823760009101908152919050565b82848237600083820160008152835162004eb481836020880162005408565b0195945050505050565b828152608081016060836020840137600081529392505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600062003da9602083018462004e39565b602081016003831062004f515762004f5162005552565b91905290565b602081016004831062004f515762004f5162005552565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b16604085015281606085015262004fb78285018b62004de6565b9150838203608085015262004fcd828a62004de6565b915060ff881660a085015283820360c085015262004fec828862004e39565b90861660e085015283810361010085015290506200500b818562004e39565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526200504d8184018a62004de6565b9050828103608084015262005063818962004de6565b905060ff871660a084015282810360c084015262005082818762004e39565b905067ffffffffffffffff851660e08401528281036101008401526200500b818562004e39565b63ffffffff831681526040602082015260006200359d604083018462004e39565b6bffffffffffffffffffffffff85168152836020820152826040820152608060608201526000620050ff608083018462004e39565b9695505050505050565b6040516101a0810167ffffffffffffffff81118282101715620051305762005130620055b0565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715620051805762005180620055b0565b604052919050565b600067ffffffffffffffff821115620051a557620051a5620055b0565b5060051b60200190565b600067ffffffffffffffff821115620051cc57620051cc620055b0565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600061ffff808316818516808303821115620052185762005218620054f4565b01949350505050565b60008219821115620052375762005237620054f4565b500190565b600063ffffffff808316818516808303821115620052185762005218620054f4565b600060ff821660ff84168060ff038211156200527e576200527e620054f4565b019392505050565b60006bffffffffffffffffffffffff808316818516808303821115620052185762005218620054f4565b600082620052c257620052c262005523565b500490565b60006bffffffffffffffffffffffff80841680620052e957620052e962005523565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620053305762005330620054f4565b500290565b600063ffffffff808316818516818304811182151516156200535b576200535b620054f4565b02949350505050565b600060ff821660ff84168160ff0481118215151615620053885762005388620054f4565b029392505050565b60006bffffffffffffffffffffffff808316818516818304811182151516156200535b576200535b620054f4565b600082821015620053d357620053d3620054f4565b500390565b60006bffffffffffffffffffffffff83811690831681811015620054005762005400620054f4565b039392505050565b60005b83811015620054255781810151838201526020016200540b565b8381111562005435576000848401525b50505050565b600181811c908216806200545057607f821691505b602082108114156200548b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620054c657620054c6620054f4565b5060010190565b600063ffffffff80831681811415620054ea57620054ea620054f4565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff81168114620029fe57600080fd5b63ffffffff81168114620029fe57600080fd5b6bffffffffffffffffffffffff81168114620029fe57600080fdfe60a060405234801561001057600080fd5b5060405161052238038061052283398101604081905261002f91610058565b60008054336001600160a01b031990911617905560601b6001600160601b031916608052610088565b60006020828403121561006a57600080fd5b81516001600160a01b038116811461008157600080fd5b9392505050565b60805160601c6104756100ad60003960008181610137015261024701526104756000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80635ab1bd53116100505780635ab1bd53146100d357806379188d1614610112578063f00e6a2a1461013557600080fd5b8063181f5a771461006c5780631a5da6c8146100be575b600080fd5b6100a86040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100b591906103c6565b60405180910390f35b6100d16100cc3660046102b0565b61015b565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100b5565b6101256101203660046102ed565b6101f3565b60405190151581526020016100b5565b7f00000000000000000000000000000000000000000000000000000000000000006100ed565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101ac576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff163314610245576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a61138881101561027657600080fd5b61138881039050846040820482031161028e57600080fd5b50803b61029a57600080fd5b60008084516020860160008589f1949350505050565b6000602082840312156102c257600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146102e657600080fd5b9392505050565b6000806040838503121561030057600080fd5b82359150602083013567ffffffffffffffff8082111561031f57600080fd5b818501915085601f83011261033357600080fd5b81358181111561034557610345610439565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561038b5761038b610439565b816040528281528860208487010111156103a457600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600060208083528351808285015260005b818110156103f3578581018301518582016040015282016103d7565b81811115610405576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a6101206040523480156200001257600080fd5b5060405162008dc738038062008dc7833981016040819052620000359162000246565b8383838360405162000047906200021b565b620000569493929190620002a8565b604051809103906000f08015801562000073573d6000803e3d6000fd5b50848484843380600081620000cf5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001025762000102816200016f565b5050508360028111156200011a576200011a620002f2565b60e0816002811115620001315762000131620002f2565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b9091166101005250620003089350505050565b6001600160a01b038116331415620001ca5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000c6565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b612e548062005f7383390190565b80516001600160a01b03811681146200024157600080fd5b919050565b600080600080608085870312156200025d57600080fd5b8451600381106200026d57600080fd5b93506200027d6020860162000229565b92506200028d6040860162000229565b91506200029d6060860162000229565b905092959194509250565b6080810160038610620002cb57634e487b7160e01b600052602160045260246000fd5b9481526001600160a01b039384166020820152918316604083015290911660609091015290565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c615bd1620003a2600039600081816101f20152818161025101526102e00152600081816102b50152818161350601528181613f2f01526140fc01526000818161034101526132c201526000818161049001526133b60152600081816104fc015281816111ab01528181611df501526126ad0152615bd16000f3fe60806040523480156200001157600080fd5b5060043610620001f05760003560e01c80638dcf0fe71162000111578063b121e14711620000a5578063ca30e603116200007b578063ca30e60314620004fa578063eb5dcd6c1462000521578063f2fde38b1462000538578063f7d334ba146200054f57620001f0565b8063b121e14714620004b5578063b148ab6b14620004cc578063c804802214620004e357620001f0565b8063a710b22111620000e7578063a710b2211462000444578063a72aa27e146200045b578063aab9edd61462000472578063b10b673c146200048e57620001f0565b80638dcf0fe714620003ff5780638e86139b14620004165780639fab4386146200042d57620001f0565b80636709d0e5116200018957806379ba5097116200015f57806379ba509714620003a857806385c1b0ba14620003b25780638765ecbe14620003c95780638da5cb5b14620003e057620001f0565b80636709d0e5146200033f57806371791aa01462000366578063744bfe61146200039157620001f0565b80634b4fd03b11620001cb5780634b4fd03b14620002b35780634ff597c114620002da5780635147cd5914620003025780635165f2f5146200032857620001f0565b80631a2af0111462000236578063349e8cca146200024f57806348013d7b146200029b575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e80801562000231573d6000f35b3d6000fd5b6200024d6200024736600462004c99565b62000566565b005b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b620002a4600081565b60405162000292919062005149565b7f0000000000000000000000000000000000000000000000000000000000000000620002a4565b620002717f000000000000000000000000000000000000000000000000000000000000000081565b620003196200031336600462004c65565b50600090565b6040516200029291906200515e565b6200024d6200033936600462004c65565b620006bf565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200037d6200037736600462004d11565b6200084e565b604051620002929695949392919062005092565b6200024d620003a236600462004c99565b62000e91565b6200024d620012ba565b6200024d620003c336600462004a00565b620013bd565b6200024d620003da36600462004c65565b62001e85565b60005473ffffffffffffffffffffffffffffffffffffffff1662000271565b6200024d6200041036600462004cc1565b62002021565b6200024d6200042736600462004be6565b62002089565b6200024d6200043e36600462004cc1565b62002447565b6200024d62000455366004620049c2565b620024fd565b6200024d6200046c36600462004d52565b62002792565b6200047b600281565b60405160ff909116815260200162000292565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200024d620004c63660046200499b565b62002878565b6200024d620004dd36600462004c65565b62002971565b6200024d620004f436600462004c65565b62002b76565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200024d62000532366004620049c2565b62002f51565b6200024d620005493660046200499b565b620030b3565b6200037d6200056036600462004c65565b620030cb565b620005718262003198565b73ffffffffffffffffffffffffffffffffffffffff8116331415620005c2576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811662000610576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff828116911614620006bb5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b620006ca8162003198565b600081815260046020908152604091829020825161010081018452815463ffffffff8082168352640100000000820481169483019490945260ff68010000000000000000820416151594820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490931660c082015260029091015490911660e082015290620007dd576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1690556200081e6002836200324d565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b600060606000806000806200086262003264565b6000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008b8152602001908152602001600020604051806101000160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816020015163ffffffff161462000bd057505060408051602081019091526000808252965094506001935085915081905062000e87565b80604001511562000c0157505060408051602081019091526000808252965094506002935085915081905062000e87565b62000c0c826200329f565b825160125492965090945060009162000c4c9185917801000000000000000000000000000000000000000000000000900463ffffffff16888886620034b1565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000c9a57600060405180602001604052806000815250600698509850985050505062000e87565b5a95506000636e04ff0d60e01b8b60405160240162000cba919062005134565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260e08501516012549251919350600092839273ffffffffffffffffffffffffffffffffffffffff9092169163ffffffff9091169062000d7890869062004f0c565b60006040518083038160008787f1925050503d806000811462000db8576040519150601f19603f3d011682016040523d82523d6000602084013e62000dbd565b606091505b50915091505a62000dcf908a620053ec565b98508162000de1576003995062000e7d565b8080602001905181019062000df7919062004b8f565b909c5090508b62000e285760006040518060200160405280600081525060049b509b509b5050505050505062000e87565b6012548151780100000000000000000000000000000000000000000000000090910463ffffffff16101562000e7d5760006040518060200160405280600081525060059b509b509b5050505050505062000e87565b909a509850505050505b9295509295509295565b600f546f01000000000000000000000000000000900460ff161562000ee2576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff811662000f6a576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815161010081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606084015260018401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c0830152600290920154821660e08201528685526005909352922054909116331462001083576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200108d62003500565b816020015163ffffffff161115620010d1576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546015546c010000000000000000000000009091046bffffffffffffffffffffffff169062001113908290620053ec565b60155560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015620011f257600080fd5b505af115801562001207573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200122d919062004b6f565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001341576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff166003811115620013fc57620013fc62005559565b14158015620014485750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff16600381111562001445576200144562005559565b14155b1562001480576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16620014e0576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162001518576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115620015775762001577620055e6565b604051908082528060200260200182016040528015620015ac57816020015b6060815260200190600190039081620015965790505b50905060008667ffffffffffffffff811115620015cd57620015cd620055e6565b604051908082528060200260200182016040528015620015f7578160200160208202803683370190505b50905060008767ffffffffffffffff811115620016185762001618620055e6565b604051908082528060200260200182016040528015620016a757816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620016375790505b50905060008867ffffffffffffffff811115620016c857620016c8620055e6565b604051908082528060200260200182016040528015620016fd57816020015b6060815260200190600190039081620016e75790505b50905060005b8981101562001b9e578a8a82818110620017215762001721620055b7565b60209081029290920135600081815260048452604090819020815161010081018352815463ffffffff8082168352640100000000820481169783019790975260ff6801000000000000000082041615159382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490951660c0860152600201541660e084015299509097506200181290508862003198565b86838281518110620018285762001828620055b7565b6020026020010181905250600760008981526020019081526020016000208054620018539062005469565b80601f0160208091040260200160405190810160405280929190818152602001828054620018819062005469565b8015620018d25780601f10620018a657610100808354040283529160200191620018d2565b820191906000526020600020905b815481529060010190602001808311620018b457829003601f168201915b5050505050858281518110620018ec57620018ec620055b7565b60200260200101819052506005600089815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16848281518110620019405762001940620055b7565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101820152600089815260179091526040902080546200197f9062005469565b80601f0160208091040260200160405190810160405280929190818152602001828054620019ad9062005469565b8015620019fe5780601f10620019d257610100808354040283529160200191620019fe565b820191906000526020600020905b815481529060010190602001808311620019e057829003601f168201915b505050505082828151811062001a185762001a18620055b7565b60200260200101819052508660a001516bffffffffffffffffffffffff168662001a43919062005294565b600089815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091975062001ae391906200441b565b600088815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b24600289620035cd565b5060a0870151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8b16602083015289917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001b9581620054bf565b91505062001703565b508460155462001baf9190620053ec565b60155560405160009062001bd2908c908c90869089908990889060200162004f2a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808b1691638e86139b916c010000000000000000000000009091041663c71249ab60028d73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381600087803b15801562001c8757600080fd5b505af115801562001c9c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001cc2919062004dd6565b866040518463ffffffff1660e01b815260040162001ce39392919062005175565b60006040518083038186803b15801562001cfc57600080fd5b505afa15801562001d11573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001d59919081019062004c2c565b6040518263ffffffff1660e01b815260040162001d77919062005134565b600060405180830381600087803b15801562001d9257600080fd5b505af115801562001da7573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c81166004830152602482018a90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b15801562001e3c57600080fd5b505af115801562001e51573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001e77919062004b6f565b505050505050505050505050565b62001e908162003198565b600081815260046020908152604091829020825161010081018452815463ffffffff8082168352640100000000820481169483019490945260ff680100000000000000008204161580159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490931660c082015260029091015490911660e08201529062001fa5576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff166801000000000000000017905562001ff1600283620035cd565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b6200202c8362003198565b6000838152601760205260409020620020479083836200445a565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485083836040516200207c929190620050e7565b60405180910390a2505050565b60023360009081526016602052604090205460ff166003811115620020b257620020b262005559565b14158015620020e8575060033360009081526016602052604090205460ff166003811115620020e557620020e562005559565b14155b1562002120576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080620021348688018862004a8d565b9450945094509450945060005b85518110156200243d57600073ffffffffffffffffffffffffffffffffffffffff16858281518110620021785762002178620055b7565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1614156200225f57848181518110620021b657620021b6620055b7565b602002602001015160e00151604051620021d09062004507565b73ffffffffffffffffffffffffffffffffffffffff9091168152602001604051809103906000f0801580156200220a573d6000803e3d6000fd5b50858281518110620022205762002220620055b7565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002382868281518110620022785762002278620055b7565b6020026020010151868381518110620022955762002295620055b7565b602002602001015160e00151878481518110620022b657620022b6620055b7565b602002602001015160000151868581518110620022d757620022d7620055b7565b6020026020010151898681518110620022f457620022f4620055b7565b602002602001015160a00151898781518110620023155762002315620055b7565b60200260200101518b8881518110620023325762002332620055b7565b602002602001015160400151898981518110620023535762002353620055b7565b60200260200101518d8a81518110620023705762002370620055b7565b602002602001015160600151620035db565b858181518110620023975762002397620055b7565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71868381518110620023d557620023d5620055b7565b602002602001015160a0015133604051620024209291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2806200243481620054bf565b91505062002141565b5050505050505050565b620024528362003198565b60125474010000000000000000000000000000000000000000900463ffffffff16811115620024ad576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020620024c89083836200445a565b50827f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf83836040516200207c929190620050e7565b73ffffffffffffffffffffffffffffffffffffffff81166200254b576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314620025ac576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b54600091620025e591859170010000000000000000000000000000000090046bffffffffffffffffffffffff169062003aa8565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260086020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905560155490915062002651906bffffffffffffffffffffffff831690620053ec565b6015556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b158015620026f257600080fd5b505af115801562002707573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200272d919062004b6f565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff161080620027bc575060125463ffffffff6401000000009091048116908216115b15620027f4576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620027ff8262003198565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8516908117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600d6020526040902054163314620028d9576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600c602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600d909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825161010081018452815463ffffffff80821683526401000000008204811694830185905260ff6801000000000000000083041615159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004851660c083015260029092015490911660e0820152911462002a84576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462002ae2576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b6000818152600460209081526040808320815161010081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004841660c083015260029092015490911660e0820152929114159062002c7460005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905081801562002ccf575080801562002ccd575062002cc062003500565b836020015163ffffffff16115b155b1562002d07576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8015801562002d3a575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002d72576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600062002d7e62003500565b90508162002d965762002d9360328262005294565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790915562002df1906002908790620035cd16565b5060105460808501516bffffffffffffffffffffffff918216916000911682111562002e5a57608086015162002e28908362005406565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002e5a575060a08501515b808660a0015162002e6c919062005406565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560115462002ed491839116620052d7565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c602052604090205416331462002fb2576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811633141562003003576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600d6020526040902054811690821614620006bb5773ffffffffffffffffffffffffffffffffffffffff8281166000818152600d602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b620030bd62003cdb565b620030c88162003d5e565b50565b600060606000806000806200318487600760008a81526020019081526020016000208054620030fa9062005469565b80601f0160208091040260200160405190810160405280929190818152602001828054620031289062005469565b8015620031795780601f106200314d5761010080835404028352916020019162003179565b820191906000526020600020905b8154815290600101906020018083116200315b57829003601f168201915b50505050506200084e565b949c939b5091995097509550909350915050565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031f6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff90811614620030c8576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006200325b838362003e56565b90505b92915050565b32156200329d576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200332757600080fd5b505afa1580156200333c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003362919062004d81565b50945090925050506000811315806200337a57508142105b806200339f57508280156200339f5750620033968242620053ec565b8463ffffffff16105b15620033b0576013549550620033b4565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200341b57600080fd5b505afa15801562003430573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003456919062004d81565b50945090925050506000811315806200346e57508142105b806200349357508280156200349357506200348a8242620053ec565b8463ffffffff16105b15620034a4576014549450620034a8565b8094505b50505050915091565b600080620034c486896000015162003ea8565b9050600080620034e18a8a63ffffffff16858a8a60018b62003ef6565b9092509050620034f28183620052d7565b9a9950505050505050505050565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003539576200353962005559565b1415620035c857606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200358857600080fd5b505afa1580156200359d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035c3919062004c7f565b905090565b504390565b60006200325b838362004317565b600f546e010000000000000000000000000000900460ff16156200362b576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff88163b6200367a576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125484517401000000000000000000000000000000000000000090910463ffffffff161015620036d7576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8763ffffffff16108062003701575060125463ffffffff6401000000009091048116908816115b1562003739576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008981526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562003799576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061010001604052808863ffffffff16815260200163ffffffff8016815260200184151581526020018273ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001866bffffffffffffffffffffffff168152602001600063ffffffff1681526020018973ffffffffffffffffffffffffffffffffffffffff16815250600460008b815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505085600560008b815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550846bffffffffffffffffffffffff1660155462003a48919062005294565b6015556000898152600760209081526040909120855162003a6c9287019062004515565b506000898152601760209081526040909120835162003a8e9285019062004515565b5062003a9c60028a6200324d565b50505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e0100000000000000000000000000009092041660608201819052829062003b35908662005406565b9050600062003b45858362005321565b9050808360400181815162003b5b9190620052d7565b6bffffffffffffffffffffffff908116909152871660608501525062003b828582620053be565b62003b8e908362005406565b6011805460009062003bb09084906bffffffffffffffffffffffff16620052d7565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146200329d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162001338565b73ffffffffffffffffffffffffffffffffffffffff811633141562003de0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162001338565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600081815260018301602052604081205462003e9f575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200325e565b5060006200325e565b600062003ebd63ffffffff841660146200534f565b62003eca836001620052af565b62003edb9060ff16611d4c6200534f565b62003eea90620124f862005294565b6200325b919062005294565b6000806000896080015161ffff168762003f1191906200534f565b905083801562003f205750803a105b1562003f2957503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003f625762003f6262005559565b1415620040f857604080516000815260208101909152851562003fc75760003660405180608001604052806048815260200162005b7d6048913960405160200162003fb09392919062004ee3565b604051602081830303815290604052905062004049565b60125462003ff9907801000000000000000000000000000000000000000000000000900463ffffffff1660046200538f565b63ffffffff1667ffffffffffffffff8111156200401a576200401a620055e6565b6040519080825280601f01601f19166020018201604052801562004045576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906200409b90849060040162005134565b60206040518083038186803b158015620040b457600080fd5b505afa158015620040c9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620040ef919062004c7f565b915050620041bc565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156200412f576200412f62005559565b1415620041bc57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200417e57600080fd5b505afa15801562004193573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620041b9919062004c7f565b90505b84620041db57808b6080015161ffff16620041d891906200534f565b90505b620041eb61ffff8716826200530a565b905060008782620041fd8c8e62005294565b6200420990866200534f565b62004215919062005294565b6200422990670de0b6b3a76400006200534f565b6200423591906200530a565b905060008c6040015163ffffffff1664e8d4a510006200425691906200534f565b898e6020015163ffffffff16858f886200427191906200534f565b6200427d919062005294565b6200428d90633b9aca006200534f565b6200429991906200534f565b620042a591906200530a565b620042b1919062005294565b90506b033b2e3c9fd0803ce8000000620042cc828462005294565b111562004305576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60008181526001830160205260408120548015620044105760006200433e600183620053ec565b85549091506000906200435490600190620053ec565b9050818114620043c0576000866000018281548110620043785762004378620055b7565b90600052602060002001549050808760000184815481106200439e576200439e620055b7565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620043d457620043d462005588565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200325e565b60009150506200325e565b508054620044299062005469565b6000825580601f106200443a575050565b601f016020900490600052602060002090810190620030c8919062004592565b828054620044689062005469565b90600052602060002090601f0160209004810192826200448c5760008555620044f5565b82601f10620044c5578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555620044f5565b82800160010185558215620044f5579182015b82811115620044f5578235825591602001919060010190620044d8565b506200450392915062004592565b5090565b610522806200565b83390190565b828054620045239062005469565b90600052602060002090601f016020900481019282620045475760008555620044f5565b82601f106200456257805160ff1916838001178555620044f5565b82800160010185558215620044f5579182015b82811115620044f557825182559160200191906001019062004575565b5b8082111562004503576000815560010162004593565b8035620045b68162005628565b919050565b600082601f830112620045cd57600080fd5b81356020620045e6620045e08362005224565b620051d2565b80838252828201915082860187848660051b89010111156200460757600080fd5b60005b8581101562004633578135620046208162005628565b845292840192908401906001016200460a565b5090979650505050505050565b600082601f8301126200465257600080fd5b8135602062004665620045e08362005224565b80838252828201915082860187848660051b89010111156200468657600080fd5b6000805b86811015620046ce57823567ffffffffffffffff811115620046aa578283fd5b620046ba8b88838d0101620048a4565b86525093850193918501916001016200468a565b509198975050505050505050565b600082601f830112620046ee57600080fd5b8135602062004701620045e08362005224565b80838252828201915082860187848660081b89010111156200472257600080fd5b6000805b86811015620046ce5761010080848c03121562004741578283fd5b6200474b620051a5565b62004756856200494e565b8152620047658886016200494e565b888201526040808601356200477a816200564b565b9082015260606200478d868201620045a9565b908201526080620047a08682016200497e565b9082015260a0620047b38682016200497e565b9082015260c0620047c68682016200494e565b9082015260e0620047d9868201620045a9565b90820152865294860194929092019160010162004726565b600082601f8301126200480357600080fd5b8135602062004816620045e08362005224565b80838252828201915082860187848660051b89010111156200483757600080fd5b60005b8581101562004633578135845292840192908401906001016200483a565b60008083601f8401126200486b57600080fd5b50813567ffffffffffffffff8111156200488457600080fd5b6020830191508360208285010111156200489d57600080fd5b9250929050565b600082601f830112620048b657600080fd5b8135620048c7620045e0826200524b565b818152846020838601011115620048dd57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126200490c57600080fd5b81516200491d620045e0826200524b565b8181528460208386010111156200493357600080fd5b6200494682602083016020870162005436565b949350505050565b803563ffffffff81168114620045b657600080fd5b805169ffffffffffffffffffff81168114620045b657600080fd5b80356bffffffffffffffffffffffff81168114620045b657600080fd5b600060208284031215620049ae57600080fd5b8135620049bb8162005628565b9392505050565b60008060408385031215620049d657600080fd5b8235620049e38162005628565b91506020830135620049f58162005628565b809150509250929050565b60008060006040848603121562004a1657600080fd5b833567ffffffffffffffff8082111562004a2f57600080fd5b818601915086601f83011262004a4457600080fd5b81358181111562004a5457600080fd5b8760208260051b850101111562004a6a57600080fd5b6020928301955093505084013562004a828162005628565b809150509250925092565b600080600080600060a0868803121562004aa657600080fd5b853567ffffffffffffffff8082111562004abf57600080fd5b62004acd89838a01620047f1565b9650602088013591508082111562004ae457600080fd5b62004af289838a01620046dc565b9550604088013591508082111562004b0957600080fd5b62004b1789838a0162004640565b9450606088013591508082111562004b2e57600080fd5b62004b3c89838a01620045bb565b9350608088013591508082111562004b5357600080fd5b5062004b628882890162004640565b9150509295509295909350565b60006020828403121562004b8257600080fd5b8151620049bb816200564b565b6000806040838503121562004ba357600080fd5b825162004bb0816200564b565b602084015190925067ffffffffffffffff81111562004bce57600080fd5b62004bdc85828601620048fa565b9150509250929050565b6000806020838503121562004bfa57600080fd5b823567ffffffffffffffff81111562004c1257600080fd5b62004c208582860162004858565b90969095509350505050565b60006020828403121562004c3f57600080fd5b815167ffffffffffffffff81111562004c5757600080fd5b6200494684828501620048fa565b60006020828403121562004c7857600080fd5b5035919050565b60006020828403121562004c9257600080fd5b5051919050565b6000806040838503121562004cad57600080fd5b823591506020830135620049f58162005628565b60008060006040848603121562004cd757600080fd5b83359250602084013567ffffffffffffffff81111562004cf657600080fd5b62004d048682870162004858565b9497909650939450505050565b6000806040838503121562004d2557600080fd5b82359150602083013567ffffffffffffffff81111562004d4457600080fd5b62004bdc85828601620048a4565b6000806040838503121562004d6657600080fd5b8235915062004d78602084016200494e565b90509250929050565b600080600080600060a0868803121562004d9a57600080fd5b62004da58662004963565b945060208601519350604086015192506060860151915062004dca6080870162004963565b90509295509295909350565b60006020828403121562004de957600080fd5b815160ff81168114620049bb57600080fd5b600081518084526020808501945080840160005b8381101562004e4357815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004e0f565b509495945050505050565b6000815180845260208085019450848260051b860182860160005b858110156200463357838303895262004e8483835162004e97565b9885019892509084019060010162004e69565b6000815180845262004eb181602086016020860162005436565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b82848237600083820160008152835162004f0281836020880162005436565b0195945050505050565b6000825162004f2081846020870162005436565b9190910192915050565b600060a0808352878184015260c07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff89111562004f6657600080fd5b8860051b808b838701378085019050818101600081526020838784030181880152818b5180845260e093508385019150828d01945060005b818110156200503e578551805163ffffffff90811685528582015181168686015260408083015115159086015260608083015173ffffffffffffffffffffffffffffffffffffffff908116918701919091526080808401516bffffffffffffffffffffffff908116918801919091528b840151168b8701528983015190911689860152908601511685840152948301946101009092019160010162004f9e565b5050878103604089015262005054818c62004e4e565b9550505050505082810360608401526200506f818662004dfb565b9050828103608084015262005085818562004e4e565b9998505050505050505050565b861515815260c060208201526000620050af60c083018862004e97565b905060078610620050c457620050c462005559565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6020815260006200325b602083018462004e97565b60208101620051588362005615565b91905290565b602081016004831062005158576200515862005559565b60ff8416815260ff831660208201526060604082015260006200519c606083018462004e97565b95945050505050565b604051610100810167ffffffffffffffff81118282101715620051cc57620051cc620055e6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156200521c576200521c620055e6565b604052919050565b600067ffffffffffffffff821115620052415762005241620055e6565b5060051b60200190565b600067ffffffffffffffff821115620052685762005268620055e6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60008219821115620052aa57620052aa620054fb565b500190565b600060ff821660ff84168060ff03821115620052cf57620052cf620054fb565b019392505050565b60006bffffffffffffffffffffffff808316818516808303821115620053015762005301620054fb565b01949350505050565b6000826200531c576200531c6200552a565b500490565b60006bffffffffffffffffffffffff808416806200534357620053436200552a565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156200538a576200538a620054fb565b500290565b600063ffffffff80831681851681830481118215151615620053b557620053b5620054fb565b02949350505050565b60006bffffffffffffffffffffffff80831681851681830481118215151615620053b557620053b5620054fb565b600082821015620054015762005401620054fb565b500390565b60006bffffffffffffffffffffffff838116908316818110156200542e576200542e620054fb565b039392505050565b60005b838110156200545357818101518382015260200162005439565b8381111562005463576000848401525b50505050565b600181811c908216806200547e57607f821691505b60208210811415620054b9577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620054f457620054f4620054fb565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60038110620030c857620030c862005559565b73ffffffffffffffffffffffffffffffffffffffff81168114620030c857600080fd5b8015158114620030c857600080fdfe60a060405234801561001057600080fd5b5060405161052238038061052283398101604081905261002f91610058565b60008054336001600160a01b031990911617905560601b6001600160601b031916608052610088565b60006020828403121561006a57600080fd5b81516001600160a01b038116811461008157600080fd5b9392505050565b60805160601c6104756100ad60003960008181610137015261024701526104756000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80635ab1bd53116100505780635ab1bd53146100d357806379188d1614610112578063f00e6a2a1461013557600080fd5b8063181f5a771461006c5780631a5da6c8146100be575b600080fd5b6100a86040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100b591906103c6565b60405180910390f35b6100d16100cc3660046102b0565b61015b565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100b5565b6101256101203660046102ed565b6101f3565b60405190151581526020016100b5565b7f00000000000000000000000000000000000000000000000000000000000000006100ed565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101ac576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff163314610245576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a61138881101561027657600080fd5b61138881039050846040820482031161028e57600080fd5b50803b61029a57600080fd5b60008084516020860160008589f1949350505050565b6000602082840312156102c257600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146102e657600080fd5b9392505050565b6000806040838503121561030057600080fd5b82359150602083013567ffffffffffffffff8082111561031f57600080fd5b818501915085601f83011261033357600080fd5b81358181111561034557610345610439565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561038b5761038b610439565b816040528281528860208487010111156103a457600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600060208083528351808285015260005b818110156103f3578581018301518582016040015282016103d7565b81811115610405576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a6101006040523480156200001257600080fd5b5060405162002e5438038062002e548339810160408190526200003591620001ef565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000126565b505050836002811115620000db57620000db62000251565b60e0816002811115620000f257620000f262000251565b60f81b9052506001600160601b0319606093841b811660805291831b821660a05290911b1660c05250620002679350505050565b6001600160a01b038116331415620001815760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ea57600080fd5b919050565b600080600080608085870312156200020657600080fd5b8451600381106200021657600080fd5b93506200022660208601620001d2565b92506200023660408601620001d2565b91506200024660608601620001d2565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c612b76620002de600039600081816102af01528181611ca30152611e5a0152600081816102fe015261186301526000818161037b015261194c0152600081816103dc0152818161122c0152818161139001526114230152612b766000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806379ba5097116100e3578063b79550be1161008c578063ed56b3e111610066578063ed56b3e114610400578063f2fde38b14610473578063faa3e9961461048657600080fd5b8063b79550be146103b2578063c7c3a19a146103ba578063ca30e603146103da57600080fd5b80638da5cb5b116100bd5780638da5cb5b1461035b578063b10b673c14610379578063b657bc9c1461039f57600080fd5b806379ba5097146103435780637d9b97e01461034b5780638456cb591461035357600080fd5b80633f4ba83a116101455780634b4fd03b1161011f5780634b4fd03b146102ad5780635147cd59146102db5780636709d0e5146102fc57600080fd5b80633f4ba83a1461022c5780634184e12c14610234578063421d183b1461024757600080fd5b80631865c57d116101765780631865c57d146101eb578063187256e8146102045780633b9cce591461021957600080fd5b806306e3b632146101925780630e08ae84146101bb575b600080fd5b6101a56101a03660046121ab565b6104bf565b6040516101b2919061251c565b60405180910390f35b6101ce6101c9366004612206565b6105c8565b6040516bffffffffffffffffffffffff90911681526020016101b2565b6101f361070b565b6040516101b295949392919061259a565b6102176102123660046120ab565b610ace565b005b6102176102273660046120e2565b610b3f565b610217610d95565b6101a56102423660046121cd565b610dfb565b61025a610255366004612090565b610f34565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a0016101b2565b7f00000000000000000000000000000000000000000000000000000000000000006040516101b29190612586565b6102ef6102e9366004612179565b50600090565b6040516101b29190612573565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101b2565b610217611052565b610217611154565b6102176112c2565b60005473ffffffffffffffffffffffffffffffffffffffff1661031e565b7f000000000000000000000000000000000000000000000000000000000000000061031e565b6101ce6103ad366004612179565b611333565b610217611357565b6103cd6103c8366004612179565b6114c2565b6040516101b291906126a7565b7f000000000000000000000000000000000000000000000000000000000000000061031e565b61045a61040e366004612090565b73ffffffffffffffffffffffffffffffffffffffff1660009081526009602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff9091166020830152016101b2565b610217610481366004612090565b61180f565b6102ef610494366004612090565b73ffffffffffffffffffffffffffffffffffffffff1660009081526016602052604090205460ff1690565b606060006104cd6002611823565b90508083116104dc57826104de565b805b92508284111561051a576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610526858561291c565b905060008167ffffffffffffffff81111561054357610543612ad5565b60405190808252806020026020018201604052801561056c578160200160208202803683370190505b50905060005b828110156105be5761058f6105878883612807565b60029061182d565b8282815181106105a1576105a1612aa6565b6020908102919091010152806105b6816129e0565b915050610572565b5095945050505050565b6040805161012081018252600f5460ff808216835263ffffffff6101008084048216602086015265010000000000840482169585019590955262ffffff6901000000000000000000840416606085015261ffff6c0100000000000000000000000084041660808501526e01000000000000000000000000000083048216151560a08501526f010000000000000000000000000000008304909116151560c08401526bffffffffffffffffffffffff70010000000000000000000000000000000083041660e08401527c0100000000000000000000000000000000000000000000000000000000909104169181019190915260009081806106c783611840565b601254919350915061070290849087907801000000000000000000000000000000000000000000000000900463ffffffff1685856000611a3c565b95945050505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101a081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101829052610180810191909152604080516101408101825260125468010000000000000000900463ffffffff1681526011546bffffffffffffffffffffffff908116602083015260155492820192909252600f547001000000000000000000000000000000009004909116606080830191909152908190600090608081016108406002611823565b815260125463ffffffff6c01000000000000000000000000808304821660208086019190915270010000000000000000000000000000000084048316604080870191909152600e54606080880191909152600f547c0100000000000000000000000000000000000000000000000000000000810486166080808a019190915260ff6e01000000000000000000000000000083048116151560a09a8b015284516101a0810186526101008085048a1682526501000000000085048a1682890152898b168288015262ffffff69010000000000000000008604169582019590955261ffff88850416928101929092526010546bffffffffffffffffffffffff81169a83019a909a526401000000008904881660c0830152740100000000000000000000000000000000000000008904881660e083015278010000000000000000000000000000000000000000000000009098049096169186019190915260135461012086015260145461014086015273ffffffffffffffffffffffffffffffffffffffff96849004871661016086015260115493909304909516610180840152600a8054865181840281018401909752808752969b509299508a958a959394600b9493169291859190830182828015610a4d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610a22575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610ab657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610a8b575b50505050509150945094509450945094509091929394565b610ad6611a87565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260166020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115610b3657610b36612a77565b02179055505050565b610b47611a87565b600b548114610b82576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600b54811015610d54576000600b8281548110610ba457610ba4612aa6565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600c90925260408320549193501690858585818110610bee57610bee612aa6565b9050602002016020810190610c039190612090565b905073ffffffffffffffffffffffffffffffffffffffff81161580610c96575073ffffffffffffffffffffffffffffffffffffffff821615801590610c7457508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015610c96575073ffffffffffffffffffffffffffffffffffffffff81811614155b15610ccd576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811614610d3e5773ffffffffffffffffffffffffffffffffffffffff8381166000908152600c6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b5050508080610d4c906129e0565b915050610b85565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600b8383604051610d899392919061246c565b60405180910390a15050565b610d9d611a87565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b60606000610e096002611823565b9050808411610e185783610e1a565b805b935083851115610e56576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610e62868661291c565b90506000808267ffffffffffffffff811115610e8057610e80612ad5565b604051908082528060200260200182016040528015610ea9578160200160208202803683370190505b50905060005b83811015610f1d576000610ec66105878b84612807565b9050876003811115610eda57610eda612a77565b610f0a5780838581518110610ef157610ef1612aa6565b602090810291909101015283610f06816129e0565b9450505b5080610f15816129e0565b915050610eaf565b50828214610f29578181525b979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e010000000000000000000000000000909204811660608301819052600f5484938493849384938492610fe09291700100000000000000000000000000000000900416612933565b600b54909150600090610ff39083612888565b90508260000151836020015182856040015161100f9190612844565b6060959095015173ffffffffffffffffffffffffffffffffffffffff9b8c166000908152600c6020526040902054929c919b959a50985093169550919350505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146110d8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61115c611a87565b6011546015546bffffffffffffffffffffffff9091169061117e90829061291c565b601555601180547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b602060405180830381600087803b15801561128657600080fd5b505af115801561129a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112be9190612157565b5050565b6112ca611a87565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001610df1565b6000818152600460205260408120546113519063ffffffff166105c8565b92915050565b61135f611a87565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b1580156113e757600080fd5b505afa1580156113fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141f9190612192565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb336015548461146c919061291c565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff9092166004830152602482015260440161126c565b604080516101608101825260008082526020820181905291810182905260608082018190526080820183905260a0820183905260c0820183905260e0820183905261010082018390526101208201929092526101408101919091526000828152600460209081526040808320815161010081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606080850191825260018601546bffffffffffffffffffffffff80821660808801526c0100000000000000000000000082041660a087015278010000000000000000000000000000000000000000000000009004831660c0860152600290950154831660e0850190815286516101608101885290518416815290519092168287015282511681850152878652600790945291909320805491938301916116309061298c565b80601f016020809104026020016040519081016040528092919081815260200182805461165c9061298c565b80156116a95780601f1061167e576101008083540402835291602001916116a9565b820191906000526020600020905b81548152906001019060200180831161168c57829003601f168201915b505050505081526020018260a001516bffffffffffffffffffffffff1681526020016005600086815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001826020015163ffffffff1667ffffffffffffffff1681526020018260c0015163ffffffff16815260200182608001516bffffffffffffffffffffffff1681526020018260400151151581526020016017600086815260200190815260200160002080546117869061298c565b80601f01602080910402602001604051908101604052809291908181526020018280546117b29061298c565b80156117ff5780601f106117d4576101008083540402835291602001916117ff565b820191906000526020600020905b8154815290600101906020018083116117e257829003601f168201915b5050505050815250915050919050565b611817611a87565b61182081611b0a565b50565b6000611351825490565b60006118398383611c00565b9392505050565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156118c757600080fd5b505afa1580156118db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118ff919061222c565b509450909250505060008113158061191657508142105b806119375750828015611937575061192e824261291c565b8463ffffffff16105b1561194657601354955061194a565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156119b057600080fd5b505afa1580156119c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119e8919061222c565b50945090925050506000811315806119ff57508142105b80611a205750828015611a205750611a17824261291c565b8463ffffffff16105b15611a2f576014549450611a33565b8094505b50505050915091565b600080611a4d868960000151611c2a565b9050600080611a688a8a63ffffffff16858a8a60018b611c6e565b9092509050611a778183612844565b93505050505b9695505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611b08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016110cf565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415611b8a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016110cf565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110611c1757611c17612aa6565b9060005260206000200154905092915050565b6000611c3d63ffffffff841660146128b3565b611c4883600161281f565b611c579060ff16611d4c6128b3565b611c6490620124f8612807565b6118399190612807565b6000806000896080015161ffff1687611c8791906128b3565b9050838015611c955750803a105b15611c9d57503a5b600060027f00000000000000000000000000000000000000000000000000000000000000006002811115611cd357611cd3612a77565b1415611e56576040805160008152602081019091528515611d3257600036604051806080016040528060488152602001612b2260489139604051602001611d1c93929190612445565b6040516020818303038152906040529050611dae565b601254611d62907801000000000000000000000000000000000000000000000000900463ffffffff1660046128f0565b63ffffffff1667ffffffffffffffff811115611d8057611d80612ad5565b6040519080825280601f01601f191660200182016040528015611daa576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90611dfe908490600401612560565b60206040518083038186803b158015611e1657600080fd5b505afa158015611e2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e4e9190612192565b915050611f12565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115611e8a57611e8a612a77565b1415611f1257606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015611ed757600080fd5b505afa158015611eeb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f0f9190612192565b90505b84611f2e57808b6080015161ffff16611f2b91906128b3565b90505b611f3c61ffff871682612874565b905060008782611f4c8c8e612807565b611f5690866128b3565b611f609190612807565b611f7290670de0b6b3a76400006128b3565b611f7c9190612874565b905060008c6040015163ffffffff1664e8d4a51000611f9b91906128b3565b898e6020015163ffffffff16858f88611fb491906128b3565b611fbe9190612807565b611fcc90633b9aca006128b3565b611fd691906128b3565b611fe09190612874565b611fea9190612807565b90506b033b2e3c9fd0803ce80000006120038284612807565b111561203b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461207157600080fd5b919050565b805169ffffffffffffffffffff8116811461207157600080fd5b6000602082840312156120a257600080fd5b6118398261204d565b600080604083850312156120be57600080fd5b6120c78361204d565b915060208301356120d781612b14565b809150509250929050565b600080602083850312156120f557600080fd5b823567ffffffffffffffff8082111561210d57600080fd5b818501915085601f83011261212157600080fd5b81358181111561213057600080fd5b8660208260051b850101111561214557600080fd5b60209290920196919550909350505050565b60006020828403121561216957600080fd5b8151801515811461183957600080fd5b60006020828403121561218b57600080fd5b5035919050565b6000602082840312156121a457600080fd5b5051919050565b600080604083850312156121be57600080fd5b50508035926020909101359150565b6000806000606084860312156121e257600080fd5b833592506020840135915060408401356121fb81612b14565b809150509250925092565b60006020828403121561221857600080fd5b813563ffffffff8116811461183957600080fd5b600080600080600060a0868803121561224457600080fd5b61224d86612076565b945060208601519350604086015192506060860151915061227060808701612076565b90509295509295909350565b600081518084526020808501945080840160005b838110156122c257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101612290565b509495945050505050565b600081518084526122e5816020860160208601612960565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b805163ffffffff1682526020810151612338602084018263ffffffff169052565b506040810151612350604084018263ffffffff169052565b506060810151612367606084018262ffffff169052565b50608081015161237d608084018261ffff169052565b5060a081015161239d60a08401826bffffffffffffffffffffffff169052565b5060c08101516123b560c084018263ffffffff169052565b5060e08101516123cd60e084018263ffffffff169052565b506101008181015163ffffffff8116848301525050610120818101519083015261014080820151908301526101608082015173ffffffffffffffffffffffffffffffffffffffff81168285015250506101808181015173ffffffffffffffffffffffffffffffffffffffff8116848301525b50505050565b828482376000838201600081528351612462818360208801612960565b0195945050505050565b6000604082016040835280865480835260608501915087600052602092508260002060005b828110156124c357815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101612491565b505050838103828501528481528590820160005b868110156125105773ffffffffffffffffffffffffffffffffffffffff6124fd8461204d565b16825291830191908301906001016124d7565b50979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561255457835183529284019291840191600101612538565b50909695505050505050565b60208152600061183960208301846122cd565b6020810161258083612b04565b91905290565b602081016003831061258057612580612a77565b855163ffffffff168152600061034060208801516125c860208501826bffffffffffffffffffffffff169052565b506040880151604084015260608801516125f260608501826bffffffffffffffffffffffff169052565b506080880151608084015260a088015161261460a085018263ffffffff169052565b5060c088015161262c60c085018263ffffffff169052565b5060e088015160e08401526101008089015161264f8286018263ffffffff169052565b50506101208881015115159084015261266c610140840188612317565b806102e084015261267f8184018761227c565b9050828103610300840152612694818661227c565b915050611a7d61032083018460ff169052565b602081526126ce60208201835173ffffffffffffffffffffffffffffffffffffffff169052565b600060208301516126f7604084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604083015163ffffffff811660608401525060608301516101608060808501526127266101808501836122cd565b9150608085015161274760a08601826bffffffffffffffffffffffff169052565b5060a085015173ffffffffffffffffffffffffffffffffffffffff811660c08601525060c085015167ffffffffffffffff811660e08601525060e08501516101006127998187018363ffffffff169052565b86015190506101206127ba868201836bffffffffffffffffffffffff169052565b86015190506101406127cf8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001838701529050611a7d83826122cd565b6000821982111561281a5761281a612a19565b500190565b600060ff821660ff84168060ff0382111561283c5761283c612a19565b019392505050565b60006bffffffffffffffffffffffff80831681851680830382111561286b5761286b612a19565b01949350505050565b60008261288357612883612a48565b500490565b60006bffffffffffffffffffffffff808416806128a7576128a7612a48565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156128eb576128eb612a19565b500290565b600063ffffffff8083168185168183048111821515161561291357612913612a19565b02949350505050565b60008282101561292e5761292e612a19565b500390565b60006bffffffffffffffffffffffff8381169083168181101561295857612958612a19565b039392505050565b60005b8381101561297b578181015183820152602001612963565b8381111561243f5750506000910152565b600181811c908216806129a057607f821691505b602082108114156129da577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415612a1257612a12612a19565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6004811061182057611820612a77565b6004811061182057600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var KeeperRegistryABI = KeeperRegistryMetaData.ABI

var KeeperRegistryBin = KeeperRegistryMetaData.Bin

func DeployKeeperRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistry, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry{KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

type KeeperRegistry struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCaller
	KeeperRegistryTransactor
	KeeperRegistryFilterer
}

type KeeperRegistryCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistrySession struct {
	Contract     *KeeperRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCallerSession struct {
	Contract *KeeperRegistryCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryTransactorSession struct {
	Contract     *KeeperRegistryTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryRaw struct {
	Contract *KeeperRegistry
}

type KeeperRegistryCallerRaw struct {
	Contract *KeeperRegistryCaller
}

type KeeperRegistryTransactorRaw struct {
	Contract *KeeperRegistryTransactor
}

func NewKeeperRegistry(address common.Address, backend bind.ContractBackend) (*KeeperRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry{address: address, abi: abi, KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCaller, error) {
	contract, err := bindKeeperRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCaller{contract: contract}, nil
}

func NewKeeperRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryTransactor, error) {
	contract, err := bindKeeperRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryTransactor{contract: contract}, nil
}

func NewKeeperRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryFilterer, error) {
	contract, err := bindKeeperRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFilterer{contract: contract}, nil
}

func bindKeeperRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistry *KeeperRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.KeeperRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry *KeeperRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transfer(opts)
}

func (_KeeperRegistry *KeeperRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry *KeeperRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transfer(opts)
}

func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry *KeeperRegistryCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) FallbackTo() (common.Address, error) {
	return _KeeperRegistry.Contract.FallbackTo(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) FallbackTo() (common.Address, error) {
	return _KeeperRegistry.Contract.FallbackTo(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, endIndex)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, endIndex)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, endIndex)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getActiveUpkeepIDsByType", startIndex, endIndex, trigger)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDsByType(&_KeeperRegistry.CallOpts, startIndex, endIndex, trigger)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDsByType(&_KeeperRegistry.CallOpts, startIndex, endIndex, trigger)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetFastGasFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetFastGasFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMode() (uint8, error) {
	return _KeeperRegistry.Contract.GetMode(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMode() (uint8, error) {
	return _KeeperRegistry.Contract.GetMode(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistry.Contract.GetSignerInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistry.Contract.GetSignerInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State)).(*State)
	outstruct.Config = *abi.ConvertType(out[1], new(OnchainConfig)).(*OnchainConfig)
	outstruct.Signers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.Transmitters = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	outstruct.F = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetState() (GetState,

	error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetState() (GetState,

	error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_KeeperRegistry *KeeperRegistrySession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistry.Contract.GetTransmitterInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistry.Contract.GetTransmitterInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistry.Contract.GetTriggerType(&_KeeperRegistry.CallOpts, upkeepId)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistry.Contract.GetTriggerType(&_KeeperRegistry.CallOpts, upkeepId)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (UpkeepInfo, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(UpkeepInfo)).(*UpkeepInfo)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetUpkeep(id *big.Int) (UpkeepInfo, error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCaller) INext(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "i_next")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) INext() (common.Address, error) {
	return _KeeperRegistry.Contract.INext(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) INext() (common.Address, error) {
	return _KeeperRegistry.Contract.INext(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDetails(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDetails(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) UpkeepVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) UpkeepVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistry *KeeperRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistry *KeeperRegistrySession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, transmitter)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, transmitter)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistry *KeeperRegistrySession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptUpkeepAdmin(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptUpkeepAdmin(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistry *KeeperRegistrySession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

func (_KeeperRegistry *KeeperRegistryTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistry *KeeperRegistrySession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "checkUpkeep", id, checkData)
}

func (_KeeperRegistry *KeeperRegistrySession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, checkData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, checkData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "checkUpkeep0", id)
}

func (_KeeperRegistry *KeeperRegistrySession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep0(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep0(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistry *KeeperRegistrySession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
}

func (_KeeperRegistry *KeeperRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistrySession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "pause")
}

func (_KeeperRegistry *KeeperRegistrySession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistry *KeeperRegistrySession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PauseUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PauseUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistrySession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistryTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistry *KeeperRegistrySession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistrySession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistrySession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setPayees", payees)
}

func (_KeeperRegistry *KeeperRegistrySession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPayees(&_KeeperRegistry.TransactOpts, payees)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPayees(&_KeeperRegistry.TransactOpts, payees)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistry *KeeperRegistrySession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistry *KeeperRegistrySession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistry *KeeperRegistrySession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry.TransactOpts, id, config)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepOffchainConfig(&_KeeperRegistry.TransactOpts, id, config)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_KeeperRegistry *KeeperRegistrySession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SimulatePerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SimulatePerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistry *KeeperRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

func (_KeeperRegistry *KeeperRegistryTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistry *KeeperRegistrySession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistry *KeeperRegistrySession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferUpkeepAdmin(&_KeeperRegistry.TransactOpts, id, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferUpkeepAdmin(&_KeeperRegistry.TransactOpts, id, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistrySession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Transmit(&_KeeperRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Transmit(&_KeeperRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "unpause")
}

func (_KeeperRegistry *KeeperRegistrySession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistry *KeeperRegistrySession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.UnpauseUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.UnpauseUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

func (_KeeperRegistry *KeeperRegistrySession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.UpdateCheckData(&_KeeperRegistry.TransactOpts, id, newCheckData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.UpdateCheckData(&_KeeperRegistry.TransactOpts, id, newCheckData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.RawTransact(opts, calldata)
}

func (_KeeperRegistry *KeeperRegistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Fallback(&_KeeperRegistry.TransactOpts, calldata)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Fallback(&_KeeperRegistry.TransactOpts, calldata)
}

type KeeperRegistryCancelledUpkeepReportIterator struct {
	Event *KeeperRegistryCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCancelledUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCancelledUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCancelledUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryCancelledUpkeepReport, error) {
	event := new(KeeperRegistryCancelledUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryConfigSetIterator struct {
	Event *KeeperRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryConfigSet struct {
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

func (_KeeperRegistry *KeeperRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryConfigSetIterator{contract: _KeeperRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error) {
	event := new(KeeperRegistryConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryFundsAddedIterator struct {
	Event *KeeperRegistryFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsAddedIterator{contract: _KeeperRegistry.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryFundsAdded)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error) {
	event := new(KeeperRegistryFundsAdded)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryFundsWithdrawnIterator struct {
	Event *KeeperRegistryFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error) {
	event := new(KeeperRegistryFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryInsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryInsufficientFundsUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryInsufficientFundsUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryInsufficientFundsUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryInsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryInsufficientFundsUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnerFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnerFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryOwnerFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnershipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryOwnershipTransferRequested, error) {
	event := new(KeeperRegistryOwnershipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnershipTransferredIterator struct {
	Event *KeeperRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferredIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnershipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryOwnershipTransferred, error) {
	event := new(KeeperRegistryOwnershipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPausedIterator struct {
	Event *KeeperRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryPausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPausedIterator{contract: _KeeperRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaused(log types.Log) (*KeeperRegistryPaused, error) {
	event := new(KeeperRegistryPaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeesUpdatedIterator struct {
	Event *KeeperRegistryPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryPayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeesUpdatedIterator{contract: _KeeperRegistry.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeesUpdated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryPayeesUpdated, error) {
	event := new(KeeperRegistryPayeesUpdated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeeshipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryPayeeshipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeeshipTransferredIterator struct {
	Event *KeeperRegistryPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferredIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeeshipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error) {
	event := new(KeeperRegistryPayeeshipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPaymentWithdrawnIterator struct {
	Event *KeeperRegistryPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPaymentWithdrawnIterator{contract: _KeeperRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPaymentWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error) {
	event := new(KeeperRegistryPaymentWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryReorgedUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryReorgedUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryReorgedUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryReorgedUpkeepReport, error) {
	event := new(KeeperRegistryReorgedUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryStaleUpkeepReportIterator struct {
	Event *KeeperRegistryStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryStaleUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryStaleUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryStaleUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryStaleUpkeepReport, error) {
	event := new(KeeperRegistryStaleUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryTransmittedIterator struct {
	Event *KeeperRegistryTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryTransmittedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistryTransmittedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryTransmittedIterator{contract: _KeeperRegistry.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistryTransmitted) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryTransmitted)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseTransmitted(log types.Log) (*KeeperRegistryTransmitted, error) {
	event := new(KeeperRegistryTransmitted)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUnpausedIterator struct {
	Event *KeeperRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUnpausedIterator{contract: _KeeperRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUnpaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error) {
	event := new(KeeperRegistryUnpaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepAdminTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryUpkeepAdminTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepAdminTransferredIterator{contract: _KeeperRegistry.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepAdminTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryUpkeepAdminTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepCanceledIterator struct {
	Event *KeeperRegistryUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepCanceledIterator{contract: _KeeperRegistry.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepCanceled)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error) {
	event := new(KeeperRegistryUpkeepCanceled)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistryUpkeepCheckDataUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepCheckDataUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepCheckDataUpdatedIterator{contract: _KeeperRegistry.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepCheckDataUpdated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryUpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistryUpkeepCheckDataUpdated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepGasLimitSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepGasLimitSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryUpkeepGasLimitSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepMigratedIterator struct {
	Event *KeeperRegistryUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepMigratedIterator{contract: _KeeperRegistry.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepMigrated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error) {
	event := new(KeeperRegistryUpkeepMigrated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepOffchainConfigSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepOffchainConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryUpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryUpkeepOffchainConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepPausedIterator struct {
	Event *KeeperRegistryUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPausedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepPaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryUpkeepPaused, error) {
	event := new(KeeperRegistryUpkeepPaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepPerformedIterator struct {
	Event *KeeperRegistryUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPerformedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepPerformed)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error) {
	event := new(KeeperRegistryUpkeepPerformed)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepReceivedIterator struct {
	Event *KeeperRegistryUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepReceivedIterator{contract: _KeeperRegistry.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepReceived)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error) {
	event := new(KeeperRegistryUpkeepReceived)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepRegisteredIterator struct {
	Event *KeeperRegistryUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepRegisteredIterator{contract: _KeeperRegistry.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepRegistered)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error) {
	event := new(KeeperRegistryUpkeepRegistered)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepUnpausedIterator struct {
	Event *KeeperRegistryUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepUnpausedIterator{contract: _KeeperRegistry.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepUnpaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryUpkeepUnpaused, error) {
	event := new(KeeperRegistryUpkeepUnpaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
	State        State
	Config       OnchainConfig
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

func (_KeeperRegistry *KeeperRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistry.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistry.ParseCancelledUpkeepReport(log)
	case _KeeperRegistry.abi.Events["ConfigSet"].ID:
		return _KeeperRegistry.ParseConfigSet(log)
	case _KeeperRegistry.abi.Events["FundsAdded"].ID:
		return _KeeperRegistry.ParseFundsAdded(log)
	case _KeeperRegistry.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistry.ParseFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistry.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistry.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistry.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistry.ParseOwnershipTransferRequested(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistry.ParseOwnershipTransferred(log)
	case _KeeperRegistry.abi.Events["Paused"].ID:
		return _KeeperRegistry.ParsePaused(log)
	case _KeeperRegistry.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistry.ParsePayeesUpdated(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferred(log)
	case _KeeperRegistry.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistry.ParsePaymentWithdrawn(log)
	case _KeeperRegistry.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistry.ParseReorgedUpkeepReport(log)
	case _KeeperRegistry.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistry.ParseStaleUpkeepReport(log)
	case _KeeperRegistry.abi.Events["Transmitted"].ID:
		return _KeeperRegistry.ParseTransmitted(log)
	case _KeeperRegistry.abi.Events["Unpaused"].ID:
		return _KeeperRegistry.ParseUnpaused(log)
	case _KeeperRegistry.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistry.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistry.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistry.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistry.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistry.ParseUpkeepCanceled(log)
	case _KeeperRegistry.abi.Events["UpkeepCheckDataUpdated"].ID:
		return _KeeperRegistry.ParseUpkeepCheckDataUpdated(log)
	case _KeeperRegistry.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistry.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistry.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistry.ParseUpkeepMigrated(log)
	case _KeeperRegistry.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistry.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistry.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistry.ParseUpkeepPaused(log)
	case _KeeperRegistry.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistry.ParseUpkeepPerformed(log)
	case _KeeperRegistry.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistry.ParseUpkeepReceived(log)
	case _KeeperRegistry.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistry.ParseUpkeepRegistered(log)
	case _KeeperRegistry.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistry.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f")
}

func (KeeperRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (KeeperRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96")
}

func (KeeperRegistryOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13")
}

func (KeeperRegistryStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89")
}

func (KeeperRegistryTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (KeeperRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryUpkeepCheckDataUpdated) Topic() common.Hash {
	return common.HexToHash("0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf")
}

func (KeeperRegistryUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistry *KeeperRegistry) Address() common.Address {
	return _KeeperRegistry.address
}

type KeeperRegistryInterface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error)

	GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error)

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

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (UpkeepInfo, error)

	INext(opts *bind.CallOpts) (common.Address, error)

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

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistryTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistryTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*KeeperRegistryTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error)

	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepCheckDataUpdatedIterator, error)

	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryUpkeepCheckDataUpdated, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
