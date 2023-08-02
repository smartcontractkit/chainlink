// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic_b_wrapper_2_1

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

var KeeperRegistryLogicBMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"},{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"trigger\",\"type\":\"uint8\"}],\"name\":\"getActiveUpkeepIDsByType\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getConditionalTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkCadance\",\"type\":\"uint32\"}],\"internalType\":\"structKeeperRegistryBase2_1.ConditionalTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structKeeperRegistryBase2_1.LogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structKeeperRegistryBase2_1.State\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structKeeperRegistryBase2_1.OnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPrivilegeConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004fd638038062004fd68339810160408190526200003591620001e1565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000119565b505050836002811115620000db57620000db62000243565b60e0816002811115620000f257620000f262000243565b9052506001600160a01b0392831660805290821660a0521660c05250620002599350505050565b336001600160a01b03821603620001735760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001dc57600080fd5b919050565b60008060008060808587031215620001f857600080fd5b8451600381106200020857600080fd5b93506200021860208601620001c4565b92506200022860408601620001c4565b91506200023860608601620001c4565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051614cfd620002d960003960008181610470015281816133c70152818161394a0152613add01526000818161050101526131bb015260008181610605015261329501526000818161068c01528181611d460152818161201c01528181612494015281816129b60152612a3a0152614cfd6000f3fe608060405234801561001057600080fd5b50600436106102ad5760003560e01c806379ba50971161017b578063b148ab6b116100d8578063cd7f71b51161008c578063ed56b3e111610071578063ed56b3e1146106d6578063f2fde38b14610749578063faa3e9961461075c57600080fd5b8063cd7f71b5146106b0578063eb5dcd6c146106c357600080fd5b8063b79550be116100bd578063b79550be14610662578063c7c3a19a1461066a578063ca30e6031461068a57600080fd5b8063b148ab6b1461063c578063b657bc9c1461064f57600080fd5b80638dcf0fe71161012f578063a72aa27e11610114578063a72aa27e146105f0578063b10b673c14610603578063b121e1471461062957600080fd5b80638dcf0fe7146105ca578063a710b221146105dd57600080fd5b80638456cb59116101605780638456cb59146105915780638765ecbe146105995780638da5cb5b146105ac57600080fd5b806379ba5097146105815780637d9b97e01461058957600080fd5b80633f4ba83a116102295780635165f2f5116101dd5780636709d0e5116101c25780636709d0e5146104ff578063708a056114610546578063744bfe611461056e57600080fd5b80635165f2f5146104bc5780635b6aa71c146104cf57600080fd5b8063421d183b1161020e578063421d183b146104085780634b4fd03b1461046e5780635147cd591461049c57600080fd5b80633f4ba83a146103ed5780634184e12c146103f557600080fd5b8063187256e8116102805780631a2af011116102655780631a2af011146103b4578063207b6516146103c75780633b9cce59146103da57600080fd5b8063187256e81461038157806319d97a941461039457600080fd5b806306e3b632146102b25780630b7d33e6146102db5780630d4a4fb1146102f05780631865c57d14610368575b600080fd5b6102c56102c0366004613dba565b6107a2565b6040516102d29190613ddc565b60405180910390f35b6102ee6102e9366004613e20565b6108ad565b005b6103036102fe366004613e9c565b610967565b6040516102d29190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b610370610a69565b6040516102d295949392919061406c565b6102ee61038f3660046141a3565b610e82565b6103a76103a2366004613e9c565b610ef3565b6040516102d2919061424e565b6102ee6103c2366004614261565b610f95565b6103a76103d5366004613e9c565b61109b565b6102ee6103e8366004614286565b6110b8565b6102ee61130e565b6102c561040336600461430f565b611374565b61041b610416366004614344565b6114c8565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a0016102d2565b7f00000000000000000000000000000000000000000000000000000000000000006040516102d29190614390565b6104af6104aa366004613e9c565b6115fb565b6040516102d291906143aa565b6102ee6104ca366004613e9c565b6116a6565b6104e26104dd3660046143d0565b611828565b6040516bffffffffffffffffffffffff90911681526020016102d2565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016102d2565b610559610554366004613e9c565b61195a565b604051905163ffffffff1681526020016102d2565b6102ee61057c366004614261565b611a39565b6102ee611e42565b6102ee611f44565b6102ee61209f565b6102ee6105a7366004613e9c565b612110565b60005473ffffffffffffffffffffffffffffffffffffffff16610521565b6102ee6105d8366004613e20565b612295565b6102ee6105eb3660046143fc565b6122ea565b6102ee6105fe36600461442a565b612566565b7f0000000000000000000000000000000000000000000000000000000000000000610521565b6102ee610637366004614344565b61265b565b6102ee61064a366004613e9c565b612753565b6104e261065d366004613e9c565b612950565b6102ee61297d565b61067d610678366004613e9c565b612ad9565b6040516102d2919061444f565b7f0000000000000000000000000000000000000000000000000000000000000000610521565b6102ee6106be366004613e20565b612e22565b6102ee6106d13660046143fc565b612eb9565b6107306106e4366004614344565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff9091166020830152016102d2565b6102ee610757366004614344565b613017565b61079561076a366004614344565b73ffffffffffffffffffffffffffffffffffffffff1660009081526019602052604090205460ff1690565b6040516102d291906145af565b606060006107b0600261302b565b90508083116107bf57826107c1565b805b9250828411156107fd576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061080985856145f2565b905060008167ffffffffffffffff81111561082657610826614605565b60405190808252806020026020018201604052801561084f578160200160208202803683370190505b50905060005b828110156108a15761087261086a8883614634565b600290613035565b82828151811061088457610884614647565b60209081029190910101528061089981614676565b915050610855565b50925050505b92915050565b6015546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16331461090e576040517f77c3599200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152601c60205260409020610927828483614750565b50827f2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769838360405161095a92919061486b565b60405180910390a2505050565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a081019190915260016109a4836115fb565b60018111156109b5576109b5614361565b146109bf57600080fd5b6000828152601a6020526040902080546109d8906146ae565b80601f0160208091040260200160405190810160405280929190818152602001828054610a04906146ae565b8015610a515780601f10610a2657610100808354040283529160200191610a51565b820191906000526020600020905b815481529060010190602001808311610a3457829003601f168201915b50505050508060200190518101906108a791906148b8565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101e08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820183905261018082018390526101a08201526101c0810191909152604080516101408101825260145463ffffffff7401000000000000000000000000000000000000000082041682526bffffffffffffffffffffffff908116602083015260185492820192909252601254700100000000000000000000000000000000900490911660608281019190915290819060009060808101610bb6600261302b565b81526014547801000000000000000000000000000000000000000000000000810463ffffffff9081166020808501919091527c0100000000000000000000000000000000000000000000000000000000808404831660408087019190915260115460608088019190915260125492830485166080808901919091526e010000000000000000000000000000840460ff16151560a09889015282516101e081018452610100808604881682526501000000000086048816968201969096526c010000000000000000000000008089048816948201949094526901000000000000000000850462ffffff16928101929092529282900461ffff16928101929092526013546bffffffffffffffffffffffff811696830196909652700100000000000000000000000000000000909404831660c082015260155480841660e0830152640100000000810484169282019290925268010000000000000000909104909116610120820152601654610140820152601754610160820152910473ffffffffffffffffffffffffffffffffffffffff166101808201529095506101a08101610d5e6009613048565b81526015546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16602091820152601254600d80546040805182860281018601909152818152949850899489949293600e9360ff909116928591830182828015610e0157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610dd6575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610e6a57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610e3f575b50505050509150945094509450945094509091929394565b610e8a613055565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260196020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115610eea57610eea614361565b02179055505050565b6000818152601c60205260409020805460609190610f10906146ae565b80601f0160208091040260200160405190810160405280929190818152602001828054610f3c906146ae565b8015610f895780601f10610f5e57610100808354040283529160200191610f89565b820191906000526020600020905b815481529060010190602001808311610f6c57829003601f168201915b50505050509050919050565b610f9e826130d8565b3373ffffffffffffffffffffffffffffffffffffffff821603610fed576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff8281169116146110975760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b6000818152601a60205260409020805460609190610f10906146ae565b6110c0613055565b600e5481146110fb576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e548110156112cd576000600e828154811061111d5761111d614647565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600f9092526040832054919350169085858581811061116757611167614647565b905060200201602081019061117c9190614344565b905073ffffffffffffffffffffffffffffffffffffffff8116158061120f575073ffffffffffffffffffffffffffffffffffffffff8216158015906111ed57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b801561120f575073ffffffffffffffffffffffffffffffffffffffff81811614155b15611246576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146112b75773ffffffffffffffffffffffffffffffffffffffff8381166000908152600f6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b50505080806112c590614676565b9150506110fe565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e83836040516113029392919061494a565b60405180910390a15050565b611316613055565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b60606000611382600261302b565b90508084116113915783611393565b805b9350838511156113cf576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006113db86866145f2565b90506000808267ffffffffffffffff8111156113f9576113f9614605565b604051908082528060200260200182016040528015611422578160200160208202803683370190505b50905060005b838110156114b157600061143f61086a8b84614634565b905087600181111561145357611453614361565b61145c826115fb565b600181111561146d5761146d614361565b0361149e578083858151811061148557611485614647565b60209081029190910101528361149a81614676565b9450505b50806114a981614676565b915050611428565b508282146114bd578181525b979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015282918291829182919082906115a257606082015160125460009161158e9170010000000000000000000000000000000090046bffffffffffffffffffffffff166149fc565b600e5490915061159e9082614a50565b9150505b8151602083015160408401516115b9908490614a7b565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b166000908152600f6020526040902054929b919a9499509750921694509092505050565b6000818160045b600f811015611688577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061164057611640614647565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461167657506000949350505050565b8061168081614676565b915050611602565b5081600f1a600181111561169e5761169e614361565b949350505050565b6116af816130d8565b60008181526004602090815260409182902082516101008082018552825460ff8116151580845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c082015260029091015490921660e08301526117b9576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556117f860028361318c565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b604080516101208101825260125460ff808216835263ffffffff6101008084048216602086015265010000000000840482169585019590955262ffffff6901000000000000000000840416606085015261ffff6c0100000000000000000000000084041660808501526e01000000000000000000000000000083048216151560a08501526f010000000000000000000000000000008304909116151560c08401526bffffffffffffffffffffffff70010000000000000000000000000000000083041660e08401527c01000000000000000000000000000000000000000000000000000000009091041691810191909152600090818061192783613198565b91509150611950838787601360020160049054906101000a900463ffffffff1686866000613376565b9695505050505050565b6040805160208101909152600081526001611974836115fb565b600181111561198557611985614361565b1461198f57600080fd5b6000828152601a6020526040902080546119a8906146ae565b80601f01602080910402602001604051908101604052809291908181526020018280546119d4906146ae565b8015611a215780601f106119f657610100808354040283529160200191611a21565b820191906000526020600020905b815481529060010190602001808311611a0457829003601f168201915b50505050508060200190518101906108a79190614aa0565b6012546f01000000000000000000000000000000900460ff1615611a89576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116611b10576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff91810482168387015265010000000000810482168386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606084015260018401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a08501527801000000000000000000000000000000000000000000000000900490911660c0830152600290920154821660e082015286855260059093529220549091163314611c23576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611c2b6133c1565b816040015163ffffffff161115611c6e576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546018546c010000000000000000000000009091046bffffffffffffffffffffffff1690611cae9082906145f2565b60185560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015611d91573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611db59190614aec565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611ec8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611f4c613055565b6014546018546bffffffffffffffffffffffff90911690611f6e9082906145f2565b601855601480547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b6020604051808303816000875af115801561207b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110979190614aec565b6120a7613055565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589060200161136a565b612119816130d8565b60008181526004602090815260409182902082516101008082018552825460ff8116158015845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c082015260029091015490921660e0830152612223576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055612265600283613476565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b61229e836130d8565b6000838152601b602052604090206122b7828483614750565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850838360405161095a92919061486b565b73ffffffffffffffffffffffffffffffffffffffff8116612337576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612397576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546000916123ce91859170010000000000000000000000000000000090046bffffffffffffffffffffffff1690613482565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601854909150612438906bffffffffffffffffffffffff8316906145f2565b6018556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af11580156124dd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125019190614aec565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff16108061259b575060145463ffffffff7001000000000000000000000000000000009091048116908216115b156125d2576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6125db826130d8565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff1661010063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152601060205260409020541633146126bb576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600f602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556010909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60008181526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821694830194909452650100000000008404811694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490921660e083015290911461285f576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146128bc576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b60006108a761295e836115fb565b600084815260046020526040902054610100900463ffffffff16611828565b612985613055565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612a12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a369190614b0e565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360185484612a8391906145f2565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff9092166004830152602482015260440161205c565b604080516101608101825260008082526020820181905291810182905260608082018190526080820183905260a0820183905260c0820183905260e082018390526101008201839052610120820192909252610140810191909152600082815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff918104821683870190815265010000000000820483168487015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606080860191825260018701546bffffffffffffffffffffffff80821660808901526c0100000000000000000000000082041660a088015278010000000000000000000000000000000000000000000000009004851660c0870152600290960154831660e08601908152875161016081018952905184168152905190921682880152519091168185015287865260079094529190932080549193830191612c43906146ae565b80601f0160208091040260200160405190810160405280929190818152602001828054612c6f906146ae565b8015612cbc5780601f10612c9157610100808354040283529160200191612cbc565b820191906000526020600020905b815481529060010190602001808311612c9f57829003601f168201915b505050505081526020018260a001516bffffffffffffffffffffffff1681526020016005600086815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001826040015163ffffffff1667ffffffffffffffff1681526020018260c0015163ffffffff16815260200182608001516bffffffffffffffffffffffff168152602001826000015115158152602001601b60008681526020019081526020016000208054612d99906146ae565b80601f0160208091040260200160405190810160405280929190818152602001828054612dc5906146ae565b8015612e125780601f10612de757610100808354040283529160200191612e12565b820191906000526020600020905b815481529060010190602001808311612df557829003601f168201915b5050505050815250915050919050565b612e2b836130d8565b60155463ffffffff16811115612e6d576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020612e86828483614750565b50827fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d838360405161095a92919061486b565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612f19576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603612f68576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152601060205260409020548116908216146110975773ffffffffffffffffffffffffffffffffffffffff82811660008181526010602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b61301f613055565b6130288161368a565b50565b60006108a7825490565b6000613041838361377f565b9392505050565b60606000613041836137a9565b60005473ffffffffffffffffffffffffffffffffffffffff1633146130d6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611ebf565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613135576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff90811614613028576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006130418383613804565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613224573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132489190614b41565b509450909250505060008113158061325f57508142105b806132805750828015613280575061327782426145f2565b8463ffffffff16105b1561328f576016549550613293565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156132fe573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133229190614b41565b509450909250505060008113158061333957508142105b8061335a575082801561335a575061335182426145f2565b8463ffffffff16105b1561336957601754945061336d565b8094505b50505050915091565b60008061338888878b60000151613853565b90506000806133a38b8a63ffffffff16858a8a60018b613915565b90925090506133b28183614a7b565b9b9a5050505050505050505050565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156133f7576133f7614361565b0361347157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613448573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061346c9190614b0e565b905090565b504390565b60006130418383613cc0565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e010000000000000000000000000000900490911660608201529061367e57600081606001518561351a91906149fc565b905060006135288583614a50565b9050808360400181815161353c9190614a7b565b6bffffffffffffffffffffffff169052506135578582614b91565b836060018181516135689190614a7b565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613709576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611ebf565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061379657613796614647565b9060005260206000200154905092915050565b606081600001805480602002602001604051908101604052809291908181526020018280548015610f8957602002820191906000526020600020905b8154815260200190600101908083116137e55750505050509050919050565b600081815260018301602052604081205461384b575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556108a7565b5060006108a7565b6000808085600181111561386957613869614361565b036138785750620138806138cd565b600185600181111561388c5761388c614361565b0361389b5750620186a06138cd565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6138de63ffffffff85166014614bc5565b6138e9846001614c02565b6138f89060ff16611d4c614bc5565b6139029083614634565b61390c9190614634565b95945050505050565b6000806000896080015161ffff168761392e9190614bc5565b905083801561393c5750803a105b1561394457503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561397a5761397a614361565b03613ad95760408051600081526020810190915285156139d857600036604051806080016040528060488152602001614ca9604891396040516020016139c293929190614c1b565b6040516020818303038152906040529050613a40565b6015546139f490640100000000900463ffffffff166004614c42565b63ffffffff1667ffffffffffffffff811115613a1257613a12614605565b6040519080825280601f01601f191660200182016040528015613a3c576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90613a9090849060040161424e565b602060405180830381865afa158015613aad573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ad19190614b0e565b915050613b85565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115613b0d57613b0d614361565b03613b8557606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613b5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b829190614b0e565b90505b84613ba157808b6080015161ffff16613b9e9190614bc5565b90505b613baf61ffff871682614c65565b905060008782613bbf8c8e614634565b613bc99086614bc5565b613bd39190614634565b613be590670de0b6b3a7640000614bc5565b613bef9190614c65565b905060008c6040015163ffffffff1664e8d4a51000613c0e9190614bc5565b898e6020015163ffffffff16858f88613c279190614bc5565b613c319190614634565b613c3f90633b9aca00614bc5565b613c499190614bc5565b613c539190614c65565b613c5d9190614634565b90506b033b2e3c9fd0803ce8000000613c768284614634565b1115613cae576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60008181526001830160205260408120548015613da9576000613ce46001836145f2565b8554909150600090613cf8906001906145f2565b9050818114613d5d576000866000018281548110613d1857613d18614647565b9060005260206000200154905080876000018481548110613d3b57613d3b614647565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613d6e57613d6e614c79565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506108a7565b60009150506108a7565b5092915050565b60008060408385031215613dcd57600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613e1457835183529284019291840191600101613df8565b50909695505050505050565b600080600060408486031215613e3557600080fd5b83359250602084013567ffffffffffffffff80821115613e5457600080fd5b818601915086601f830112613e6857600080fd5b813581811115613e7757600080fd5b876020828501011115613e8957600080fd5b6020830194508093505050509250925092565b600060208284031215613eae57600080fd5b5035919050565b600081518084526020808501945080840160005b83811015613efb57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613ec9565b509495945050505050565b805163ffffffff16825260006101e06020830151613f2c602086018263ffffffff169052565b506040830151613f44604086018263ffffffff169052565b506060830151613f5b606086018262ffffff169052565b506080830151613f71608086018261ffff169052565b5060a0830151613f9160a08601826bffffffffffffffffffffffff169052565b5060c0830151613fa960c086018263ffffffff169052565b5060e0830151613fc160e086018263ffffffff169052565b506101008381015163ffffffff908116918601919091526101208085015190911690850152610140808401519085015261016080840151908501526101808084015173ffffffffffffffffffffffffffffffffffffffff16908501526101a08084015181860183905261403683870182613eb5565b925050506101c0808401516140628287018273ffffffffffffffffffffffffffffffffffffffff169052565b5090949350505050565b855163ffffffff16815260006101c0602088015161409a60208501826bffffffffffffffffffffffff169052565b506040880151604084015260608801516140c460608501826bffffffffffffffffffffffff169052565b506080880151608084015260a08801516140e660a085018263ffffffff169052565b5060c08801516140fe60c085018263ffffffff169052565b5060e088015160e0840152610100808901516141218286018263ffffffff169052565b505061012088810151151590840152610140830181905261414481840188613f06565b90508281036101608401526141598187613eb5565b905082810361018084015261416e8186613eb5565b9150506119506101a083018460ff169052565b73ffffffffffffffffffffffffffffffffffffffff8116811461302857600080fd5b600080604083850312156141b657600080fd5b82356141c181614181565b91506020830135600481106141d557600080fd5b809150509250929050565b60005b838110156141fb5781810151838201526020016141e3565b50506000910152565b6000815180845261421c8160208601602086016141e0565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006130416020830184614204565b6000806040838503121561427457600080fd5b8235915060208301356141d581614181565b6000806020838503121561429957600080fd5b823567ffffffffffffffff808211156142b157600080fd5b818501915085601f8301126142c557600080fd5b8135818111156142d457600080fd5b8660208260051b85010111156142e957600080fd5b60209290920196919550909350505050565b80356002811061430a57600080fd5b919050565b60008060006060848603121561432457600080fd5b833592506020840135915061433b604085016142fb565b90509250925092565b60006020828403121561435657600080fd5b813561304181614181565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60208101600383106143a4576143a4614361565b91905290565b60208101600283106143a4576143a4614361565b63ffffffff8116811461302857600080fd5b600080604083850312156143e357600080fd5b6143ec836142fb565b915060208301356141d5816143be565b6000806040838503121561440f57600080fd5b823561441a81614181565b915060208301356141d581614181565b6000806040838503121561443d57600080fd5b8235915060208301356141d5816143be565b6020815261447660208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161449f604084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604083015163ffffffff811660608401525060608301516101608060808501526144ce610180850183614204565b915060808501516144ef60a08601826bffffffffffffffffffffffff169052565b5060a085015173ffffffffffffffffffffffffffffffffffffffff811660c08601525060c085015167ffffffffffffffff811660e08601525060e08501516101006145418187018363ffffffff169052565b8601519050610120614562868201836bffffffffffffffffffffffff169052565b86015190506101406145778682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506119508382614204565b60208101600483106143a4576143a4614361565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156108a7576108a76145c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b808201808211156108a7576108a76145c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036146a7576146a76145c3565b5060010190565b600181811c908216806146c257607f821691505b6020821081036146fb577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561474b57600081815260208120601f850160051c810160208610156147285750805b601f850160051c820191505b8181101561474757828155600101614734565b5050505b505050565b67ffffffffffffffff83111561476857614768614605565b61477c8361477683546146ae565b83614701565b6000601f8411600181146147ce57600085156147985750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614864565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561481d57868501358255602094850194600190920191016147fd565b5086821015614858577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b600060c082840312156148ca57600080fd5b60405160c0810181811067ffffffffffffffff821117156148ed576148ed614605565b60405282516148fb81614181565b8152602083015160ff8116811461491157600080fd5b8060208301525060408301516040820152606083015160608201526080830151608082015260a083015160a08201528091505092915050565b6000604082016040835280865480835260608501915087600052602092508260002060005b828110156149a157815473ffffffffffffffffffffffffffffffffffffffff168452928401926001918201910161496f565b505050838103828501528481528590820160005b868110156149f05782356149c881614181565b73ffffffffffffffffffffffffffffffffffffffff16825291830191908301906001016149b5565b50979650505050505050565b6bffffffffffffffffffffffff828116828216039080821115613db357613db36145c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680614a6f57614a6f614a21565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115613db357613db36145c3565b600060208284031215614ab257600080fd5b6040516020810181811067ffffffffffffffff82111715614ad557614ad5614605565b6040528251614ae3816143be565b81529392505050565b600060208284031215614afe57600080fd5b8151801515811461304157600080fd5b600060208284031215614b2057600080fd5b5051919050565b805169ffffffffffffffffffff8116811461430a57600080fd5b600080600080600060a08688031215614b5957600080fd5b614b6286614b27565b9450602086015193506040860151925060608601519150614b8560808701614b27565b90509295509295909350565b60006bffffffffffffffffffffffff80831681851681830481118215151615614bbc57614bbc6145c3565b02949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614bfd57614bfd6145c3565b500290565b60ff81811683821601908111156108a7576108a76145c3565b828482376000838201600081528351614c388183602088016141e0565b0195945050505050565b600063ffffffff80831681851681830481118215151615614bbc57614bbc6145c3565b600082614c7457614c74614a21565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
}

var KeeperRegistryLogicBABI = KeeperRegistryLogicBMetaData.ABI

var KeeperRegistryLogicBBin = KeeperRegistryLogicBMetaData.Bin

func DeployKeeperRegistryLogicB(auth *bind.TransactOpts, backend bind.ContractBackend, mode uint8, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogicB, error) {
	parsed, err := KeeperRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicBBin), backend, mode, link, linkNativeFeed, fastGasFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogicB{KeeperRegistryLogicBCaller: KeeperRegistryLogicBCaller{contract: contract}, KeeperRegistryLogicBTransactor: KeeperRegistryLogicBTransactor{contract: contract}, KeeperRegistryLogicBFilterer: KeeperRegistryLogicBFilterer{contract: contract}}, nil
}

type KeeperRegistryLogicB struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogicBCaller
	KeeperRegistryLogicBTransactor
	KeeperRegistryLogicBFilterer
}

type KeeperRegistryLogicBCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicBTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicBFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicBSession struct {
	Contract     *KeeperRegistryLogicB
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicBCallerSession struct {
	Contract *KeeperRegistryLogicBCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogicBTransactorSession struct {
	Contract     *KeeperRegistryLogicBTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicBRaw struct {
	Contract *KeeperRegistryLogicB
}

type KeeperRegistryLogicBCallerRaw struct {
	Contract *KeeperRegistryLogicBCaller
}

type KeeperRegistryLogicBTransactorRaw struct {
	Contract *KeeperRegistryLogicBTransactor
}

func NewKeeperRegistryLogicB(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogicB, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogicBABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogicB(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicB{address: address, abi: abi, KeeperRegistryLogicBCaller: KeeperRegistryLogicBCaller{contract: contract}, KeeperRegistryLogicBTransactor: KeeperRegistryLogicBTransactor{contract: contract}, KeeperRegistryLogicBFilterer: KeeperRegistryLogicBFilterer{contract: contract}}, nil
}

func NewKeeperRegistryLogicBCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogicBCaller, error) {
	contract, err := bindKeeperRegistryLogicB(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBCaller{contract: contract}, nil
}

func NewKeeperRegistryLogicBTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogicBTransactor, error) {
	contract, err := bindKeeperRegistryLogicB(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBTransactor{contract: contract}, nil
}

func NewKeeperRegistryLogicBFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogicBFilterer, error) {
	contract, err := bindKeeperRegistryLogicB(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBFilterer{contract: contract}, nil
}

func bindKeeperRegistryLogicB(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryLogicBMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicB.Contract.KeeperRegistryLogicBCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.KeeperRegistryLogicBTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.KeeperRegistryLogicBTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicB.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, endIndex)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetActiveUpkeepIDs(&_KeeperRegistryLogicB.CallOpts, startIndex, endIndex)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, endIndex *big.Int) ([]*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetActiveUpkeepIDs(&_KeeperRegistryLogicB.CallOpts, startIndex, endIndex)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getActiveUpkeepIDsByType", startIndex, endIndex, trigger)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetActiveUpkeepIDsByType(&_KeeperRegistryLogicB.CallOpts, startIndex, endIndex, trigger)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetActiveUpkeepIDsByType(startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetActiveUpkeepIDsByType(&_KeeperRegistryLogicB.CallOpts, startIndex, endIndex, trigger)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetConditionalTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getConditionalTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21ConditionalTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21ConditionalTriggerConfig)).(*KeeperRegistryBase21ConditionalTriggerConfig)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetConditionalTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetConditionalTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetConditionalTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21ConditionalTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetConditionalTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetLinkAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetLinkAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getLogTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21LogTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21LogTriggerConfig)).(*KeeperRegistryBase21LogTriggerConfig)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetLogTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetLogTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetLogTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetLogTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetMaxPaymentForGas(opts *bind.CallOpts, triggerType uint8, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getMaxPaymentForGas", triggerType, gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB.CallOpts, triggerType, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetMaxPaymentForGas(triggerType uint8, gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB.CallOpts, triggerType, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryLogicB.CallOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryLogicB.CallOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetMode(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetMode(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB.CallOpts, peer)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB.CallOpts, peer)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetSignerInfo(opts *bind.CallOpts, query common.Address) (GetSignerInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getSignerInfo", query)

	outstruct := new(GetSignerInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Active = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Index = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistryLogicB.Contract.GetSignerInfo(&_KeeperRegistryLogicB.CallOpts, query)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetSignerInfo(query common.Address) (GetSignerInfo,

	error) {
	return _KeeperRegistryLogicB.Contract.GetSignerInfo(&_KeeperRegistryLogicB.CallOpts, query)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getState")

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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetState() (GetState,

	error) {
	return _KeeperRegistryLogicB.Contract.GetState(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetState() (GetState,

	error) {
	return _KeeperRegistryLogicB.Contract.GetState(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (GetTransmitterInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getTransmitterInfo", query)

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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistryLogicB.Contract.GetTransmitterInfo(&_KeeperRegistryLogicB.CallOpts, query)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetTransmitterInfo(query common.Address) (GetTransmitterInfo,

	error) {
	return _KeeperRegistryLogicB.Contract.GetTransmitterInfo(&_KeeperRegistryLogicB.CallOpts, query)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetTriggerType(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicB.Contract.GetTriggerType(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getUpkeep", id)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetUpkeep(id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeep(&_KeeperRegistryLogicB.CallOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetUpkeep(id *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeep(&_KeeperRegistryLogicB.CallOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.Owner(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.Owner(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptOwnership(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptOwnership(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptPayeeship(&_KeeperRegistryLogicB.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptPayeeship(&_KeeperRegistryLogicB.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "pause")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.Pause(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.Pause(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.PauseUpkeep(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.PauseUpkeep(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.RecoverFunds(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.RecoverFunds(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setPayees", payees)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetPayees(&_KeeperRegistryLogicB.TransactOpts, payees)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetPayees(&_KeeperRegistryLogicB.TransactOpts, payees)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogicB.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepCheckData", id, newCheckData)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepCheckData(&_KeeperRegistryLogicB.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepCheckData(&_KeeperRegistryLogicB.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicB.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicB.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicB.TransactOpts, id, config)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicB.TransactOpts, id, config)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepPrivilegeConfig(&_KeeperRegistryLogicB.TransactOpts, upkeepId, newPrivilegeConfig)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferOwnership(&_KeeperRegistryLogicB.TransactOpts, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferOwnership(&_KeeperRegistryLogicB.TransactOpts, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferPayeeship(&_KeeperRegistryLogicB.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferPayeeship(&_KeeperRegistryLogicB.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicB.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicB.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "unpause")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.Unpause(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.Unpause(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.UnpauseUpkeep(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.UnpauseUpkeep(&_KeeperRegistryLogicB.TransactOpts, id)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawFunds(&_KeeperRegistryLogicB.TransactOpts, id, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawFunds(&_KeeperRegistryLogicB.TransactOpts, id, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogicB.TransactOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawPayment(&_KeeperRegistryLogicB.TransactOpts, from, to)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.WithdrawPayment(&_KeeperRegistryLogicB.TransactOpts, from, to)
}

type KeeperRegistryLogicBCancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogicBCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBCancelledUpkeepReportIterator{contract: _KeeperRegistryLogicB.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBCancelledUpkeepReport)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicBCancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogicBCancelledUpkeepReport)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBFundsAddedIterator struct {
	Event *KeeperRegistryLogicBFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicBFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBFundsAddedIterator{contract: _KeeperRegistryLogicB.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBFundsAdded)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogicBFundsAdded, error) {
	event := new(KeeperRegistryLogicBFundsAdded)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicBFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBFundsWithdrawnIterator{contract: _KeeperRegistryLogicB.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBFundsWithdrawn)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicBFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicBFundsWithdrawn)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogicBInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogicB.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBInsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicBInsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogicBInsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicBOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBOwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBOwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicBOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBOwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogicB.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBOwnerFundsWithdrawn)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicBOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicBOwnerFundsWithdrawn)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicBOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicBOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBOwnershipTransferRequestedIterator{contract: _KeeperRegistryLogicB.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBOwnershipTransferRequested)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicBOwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogicBOwnershipTransferRequested)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBOwnershipTransferredIterator struct {
	Event *KeeperRegistryLogicBOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicBOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBOwnershipTransferredIterator{contract: _KeeperRegistryLogicB.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBOwnershipTransferred)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicBOwnershipTransferred, error) {
	event := new(KeeperRegistryLogicBOwnershipTransferred)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBPausedIterator struct {
	Event *KeeperRegistryLogicBPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicBPausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBPausedIterator{contract: _KeeperRegistryLogicB.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBPaused)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParsePaused(log types.Log) (*KeeperRegistryLogicBPaused, error) {
	event := new(KeeperRegistryLogicBPaused)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBPayeesUpdatedIterator struct {
	Event *KeeperRegistryLogicBPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicBPayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBPayeesUpdatedIterator{contract: _KeeperRegistryLogicB.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBPayeesUpdated)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicBPayeesUpdated, error) {
	event := new(KeeperRegistryLogicBPayeesUpdated)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicBPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicBPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBPayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogicB.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBPayeeshipTransferRequested)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicBPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogicBPayeeshipTransferRequested)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBPayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogicBPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicBPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBPayeeshipTransferredIterator{contract: _KeeperRegistryLogicB.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBPayeeshipTransferred)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicBPayeeshipTransferred, error) {
	event := new(KeeperRegistryLogicBPayeeshipTransferred)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBPaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogicBPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicBPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBPaymentWithdrawnIterator{contract: _KeeperRegistryLogicB.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBPaymentWithdrawn)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicBPaymentWithdrawn, error) {
	event := new(KeeperRegistryLogicBPaymentWithdrawn)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogicBReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBReorgedUpkeepReportIterator{contract: _KeeperRegistryLogicB.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBReorgedUpkeepReport)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicBReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogicBReorgedUpkeepReport)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBStaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogicBStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBStaleUpkeepReportIterator{contract: _KeeperRegistryLogicB.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBStaleUpkeepReport)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicBStaleUpkeepReport, error) {
	event := new(KeeperRegistryLogicBStaleUpkeepReport)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUnpausedIterator struct {
	Event *KeeperRegistryLogicBUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicBUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUnpausedIterator{contract: _KeeperRegistryLogicB.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUnpaused)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogicBUnpaused, error) {
	event := new(KeeperRegistryLogicBUnpaused)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogicBUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogicBUpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogicBUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepAdminTransferredIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepAdminTransferred)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogicBUpkeepAdminTransferred)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepCanceledIterator struct {
	Event *KeeperRegistryLogicBUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicBUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepCanceledIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepCanceled)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicBUpkeepCanceled, error) {
	event := new(KeeperRegistryLogicBUpkeepCanceled)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepCheckDataSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepCheckDataSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepCheckDataSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicBUpkeepCheckDataSet, error) {
	event := new(KeeperRegistryLogicBUpkeepCheckDataSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepGasLimitSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepGasLimitSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicBUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogicBUpkeepGasLimitSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepMigratedIterator struct {
	Event *KeeperRegistryLogicBUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepMigratedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepMigrated)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicBUpkeepMigrated, error) {
	event := new(KeeperRegistryLogicBUpkeepMigrated)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicBUpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepPausedIterator struct {
	Event *KeeperRegistryLogicBUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepPausedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepPaused)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicBUpkeepPaused, error) {
	event := new(KeeperRegistryLogicBUpkeepPaused)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepPerformedIterator struct {
	Event *KeeperRegistryLogicBUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicBUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepPerformedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepPerformed)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicBUpkeepPerformed, error) {
	event := new(KeeperRegistryLogicBUpkeepPerformed)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepPrivilegeConfigSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicBUpkeepPrivilegeConfigSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepReceivedIterator struct {
	Event *KeeperRegistryLogicBUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepReceivedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepReceived)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicBUpkeepReceived, error) {
	event := new(KeeperRegistryLogicBUpkeepReceived)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogicBUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepRegisteredIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepRegistered)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicBUpkeepRegistered, error) {
	event := new(KeeperRegistryLogicBUpkeepRegistered)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepTriggerConfigSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepTriggerConfigSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepTriggerConfigSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepTriggerConfigSet, error) {
	event := new(KeeperRegistryLogicBUpkeepTriggerConfigSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicBUpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogicBUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepUnpausedIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepUnpaused)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicBUpkeepUnpaused, error) {
	event := new(KeeperRegistryLogicBUpkeepUnpaused)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicB) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogicB.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogicB.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogicB.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogicB.ParseFundsAdded(log)
	case _KeeperRegistryLogicB.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogicB.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogicB.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogicB.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogicB.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogicB.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogicB.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogicB.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogicB.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogicB.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogicB.abi.Events["Paused"].ID:
		return _KeeperRegistryLogicB.ParsePaused(log)
	case _KeeperRegistryLogicB.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogicB.ParsePayeesUpdated(log)
	case _KeeperRegistryLogicB.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogicB.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogicB.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogicB.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogicB.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogicB.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogicB.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogicB.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogicB.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogicB.ParseStaleUpkeepReport(log)
	case _KeeperRegistryLogicB.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogicB.ParseUnpaused(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepCheckDataSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepCheckDataSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepPaused(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepPrivilegeConfigSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepReceived(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepTriggerConfigSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicBCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (KeeperRegistryLogicBFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicBFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicBInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (KeeperRegistryLogicBOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogicBOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogicBOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogicBPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogicBPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryLogicBPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogicBPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogicBPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogicBReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (KeeperRegistryLogicBStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (KeeperRegistryLogicBUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogicBUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogicBUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogicBUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogicBUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (KeeperRegistryLogicBUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogicBUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogicBUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogicBUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicBUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryLogicBUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (KeeperRegistryLogicBUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicBUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogicBUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (KeeperRegistryLogicBUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicB) Address() common.Address {
	return _KeeperRegistryLogicB.address
}

type KeeperRegistryLogicBInterface interface {
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

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, newPrivilegeConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicBCancelledUpkeepReport, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicBFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicBFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicBFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicBInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicBOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicBOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicBOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicBOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicBOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicBOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicBPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogicBPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicBPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicBPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicBPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicBPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicBPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicBPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicBPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicBPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicBReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicBStaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicBUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogicBUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicBUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicBUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicBUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicBUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicBUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicBUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicBUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicBUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicBUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicBUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicBUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
