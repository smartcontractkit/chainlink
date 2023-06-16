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

type KeeperRegistryBase21BlockTriggerConfig struct {
	CheckCadance uint32
}

type KeeperRegistryBase21CronTriggerConfig struct {
	Cron string
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
	Registrars           []common.Address
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
	Target                              common.Address
	Forwarder                           common.Address
	ExecuteGas                          uint32
	CheckData                           []byte
	Balance                             *big.Int
	Admin                               common.Address
	MaxValidBlocknumber                 uint64
	LastPerformedBlockNumberOrTimestamp uint32
	AmountSpent                         *big.Int
	Paused                              bool
	OffchainConfig                      []byte
}

var KeeperRegistryLogicBMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PipelineDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"adminOffchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepAdminOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPipelineData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPipelineDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endIndex\",\"type\":\"uint256\"},{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"trigger\",\"type\":\"uint8\"}],\"name\":\"getActiveUpkeepIDsByType\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBlockTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"checkCadance\",\"type\":\"uint32\"}],\"internalType\":\"structKeeperRegistryBase2_1.BlockTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getCronTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"cron\",\"type\":\"string\"}],\"internalType\":\"structKeeperRegistryBase2_1.CronTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"filterSelector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"internalType\":\"structKeeperRegistryBase2_1.LogTriggerConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getSignerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"totalPremium\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"latestConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"latestConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"latestEpoch\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"internalType\":\"structKeeperRegistryBase2_1.State\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"}],\"internalType\":\"structKeeperRegistryBase2_1.OnchainConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getTransmitterInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"lastCollected\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumberOrTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"upkeepInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepAdminOffchainConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUpkeepManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_1.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newAdminOffchainConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepAdminOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newUpkeepManager\",\"type\":\"address\"}],\"name\":\"setUpkeepManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newPipelineData\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b50604051620051d5380380620051d58339810160408190526200003591620001f3565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c3816200012b565b505050836002811115620000db57620000db62000255565b60e0816002811115620000f257620000f262000255565b9052506001600160a01b0392831660805290821660a0521660c0525050601980546001600160a01b03191633179055506200026b915050565b336001600160a01b03821603620001855760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ee57600080fd5b919050565b600080600080608085870312156200020a57600080fd5b8451600381106200021a57600080fd5b93506200022a60208601620001d6565b92506200023a60408601620001d6565b91506200024a60608601620001d6565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051614eea620002eb6000396000818161054301528181613548015281816139a30152613b3601526000818161059801526132a0015260008181610653015261337a0152600081816106da01528181611dec015281816120c20152818161253a01528181612a520152612ad60152614eea6000f3fe608060405234801561001057600080fd5b50600436106102de5760003560e01c806379ba509711610186578063b148ab6b116100e3578063d921875c11610097578063ee048eae11610071578063ee048eae14610797578063f2fde38b146107bf578063faa3e996146107d257600080fd5b8063d921875c146106fe578063eb5dcd6c14610711578063ed56b3e11461072457600080fd5b8063b79550be116100c8578063b79550be146106b0578063c7c3a19a146106b8578063ca30e603146106d857600080fd5b8063b148ab6b1461068a578063b657bc9c1461069d57600080fd5b80638dcf0fe71161013a578063a72aa27e1161011f578063a72aa27e1461063e578063b10b673c14610651578063b121e1471461067757600080fd5b80638dcf0fe714610618578063a710b2211461062b57600080fd5b80638456cb591161016b5780638456cb59146105df5780638765ecbe146105e75780638da5cb5b146105fa57600080fd5b806379ba5097146105cf5780637d9b97e0146105d757600080fd5b80633733f9cc1161023f578063421d183b116101f35780635165f2f5116101cd5780635165f2f5146105835780636709d0e514610596578063744bfe61146105bc57600080fd5b8063421d183b146104db5780634b4fd03b146105415780635147cd591461057057600080fd5b80633f4ba83a116102245780633f4ba83a146104ad5780634184e12c146104b5578063418d76b6146104c857600080fd5b80633733f9cc1461047a5780633b9cce591461049a57600080fd5b80631865c57d116102965780631a2af0111161027b5780631a2af011146104345780631fffe83514610447578063207b65161461045a57600080fd5b80631865c57d14610408578063187256e81461042157600080fd5b80630d4a4fb1116102c75780630d4a4fb1146103215780630d4cbb7f146103995780630e08ae84146103d857600080fd5b806306e3b632146102e35780630d0be1471461030c575b600080fd5b6102f66102f1366004613ebd565b610818565b6040516103039190613edf565b60405180910390f35b61031f61031a366004613f23565b610923565b005b61033461032f366004613f9f565b6109c7565b6040516103039190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b60195473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610303565b6103eb6103e6366004613fca565b610ac9565b6040516bffffffffffffffffffffffff9091168152602001610303565b610410610bf8565b604051610303959493929190614155565b61031f61042f36600461428c565b610fbe565b61031f6104423660046142c9565b61102f565b61031f6104553660046142ee565b611182565b61046d610468366004613f9f565b6111d1565b6040516103039190614379565b61048d610488366004613f9f565b611273565b604051610303919061438c565b61031f6104a83660046143a7565b611352565b61031f6115a8565b6102f66104c336600461441c565b61160e565b61046d6104d6366004613f9f565b611762565b6104ee6104e93660046142ee565b61177f565b60408051951515865260ff90941660208601526bffffffffffffffffffffffff9283169385019390935216606083015273ffffffffffffffffffffffffffffffffffffffff16608082015260a001610303565b7f00000000000000000000000000000000000000000000000000000000000000005b6040516103039190614498565b61056361057e366004613f9f565b6118b2565b61031f610591366004613f9f565b61195d565b7f00000000000000000000000000000000000000000000000000000000000000006103b3565b61031f6105ca3660046142c9565b611adf565b61031f611ee8565b61031f611fea565b61031f612145565b61031f6105f5366004613f9f565b6121b6565b60005473ffffffffffffffffffffffffffffffffffffffff166103b3565b61031f610626366004613f23565b61233b565b61031f6106393660046144ab565b612390565b61031f61064c3660046144d9565b61260c565b7f00000000000000000000000000000000000000000000000000000000000000006103b3565b61031f6106853660046142ee565b612701565b61031f610698366004613f9f565b6127f9565b6103eb6106ab366004613f9f565b6129f6565b61031f612a19565b6106cb6106c6366004613f9f565b612b75565b60405161030391906144fe565b7f00000000000000000000000000000000000000000000000000000000000000006103b3565b61031f61070c366004613f23565b612ebe565b61031f61071f3660046144ab565b612f5b565b61077e6107323660046142ee565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602090815260409182902082518084019093525460ff8082161515808552610100909204169290910182905291565b60408051921515835260ff909116602083015201610303565b6107aa6107a5366004613f9f565b6130b9565b604051905163ffffffff168152602001610303565b61031f6107cd3660046142ee565b613198565b61080b6107e03660046142ee565b73ffffffffffffffffffffffffffffffffffffffff166000908152601a602052604090205460ff1690565b604051610303919061465e565b6060600061082660026131ac565b90508083116108355782610837565b805b925082841115610873576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061087f85856146a1565b905060008167ffffffffffffffff81111561089c5761089c6146b4565b6040519080825280602002602001820160405280156108c5578160200160208202803683370190505b50905060005b82811015610917576108e86108e088836146e3565b6002906131b6565b8282815181106108fa576108fa6146f6565b60209081029190910101528061090f81614725565b9150506108cb565b50925050505b92915050565b61092c836131c9565b60155463ffffffff1681111561096e576040517f2b49197900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206109878284836147ff565b50827f787b2e556c364148324d7dc0cee4322a048aaf798bda184136c8636c72ed305683836040516109ba92919061491a565b60405180910390a2505050565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091526001610a04836118b2565b6002811115610a1557610a15614459565b14610a1f57600080fd5b6000828152601b602052604090208054610a389061475d565b80601f0160208091040260200160405190810160405280929190818152602001828054610a649061475d565b8015610ab15780601f10610a8657610100808354040283529160200191610ab1565b820191906000526020600020905b815481529060010190602001808311610a9457829003601f168201915b505050505080602001905181019061091d91906149df565b604080516101208101825260125460ff808216835263ffffffff6101008084048216602086015265010000000000840482169585019590955262ffffff6901000000000000000000840416606085015261ffff6c0100000000000000000000000084041660808501526e01000000000000000000000000000083048216151560a08501526f010000000000000000000000000000008304909116151560c08401526bffffffffffffffffffffffff70010000000000000000000000000000000083041660e08401527c010000000000000000000000000000000000000000000000000000000090910416918101919091526000908180610bc88361327d565b6015549193509150610bef9084908790640100000000900463ffffffff168585600061345b565b95945050505050565b6040805161014081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810191909152604080516101a08101825260008082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820183905261012082018390526101408201839052610160820192909252610180810191909152604080516101408101825260145463ffffffff7401000000000000000000000000000000000000000082041682526bffffffffffffffffffffffff908116602083015260185492820192909252601254700100000000000000000000000000000000900490911660608281019190915290819060009060808101610d3860026131ac565b81526014547801000000000000000000000000000000000000000000000000810463ffffffff9081166020808501919091527c0100000000000000000000000000000000000000000000000000000000808404831660408087019190915260115460608088019190915260125492830485166080808901919091526e010000000000000000000000000000840460ff16151560a09889015282516101a081018452610100808604881682526501000000000086048816968201969096526c010000000000000000000000008089048816948201949094526901000000000000000000850462ffffff16928101929092529282900461ffff16928101929092526013546bffffffffffffffffffffffff811696830196909652700100000000000000000000000000000000909404831660c082015260155480841660e0830152640100000000900490921690820152601654610120820152601754610140820152910473ffffffffffffffffffffffffffffffffffffffff166101608201529095506101808101610ec860096134a6565b9052601254600d80546040805160208084028201810190925282815294985089948994600e9360ff90911692859190830182828015610f3d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f12575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610fa657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f7b575b50505050509150945094509450945094509091929394565b610fc66134b3565b73ffffffffffffffffffffffffffffffffffffffff82166000908152601a6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561102657611026614459565b02179055505050565b611038826131c9565b3373ffffffffffffffffffffffffffffffffffffffff821603611087576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166110d4576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff82811691161461117e5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b61118a6134b3565b601980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000818152601b602052604090208054606091906111ee9061475d565b80601f016020809104026020016040519081016040528092919081815260200182805461121a9061475d565b80156112675780601f1061123c57610100808354040283529160200191611267565b820191906000526020600020905b81548152906001019060200180831161124a57829003601f168201915b50505050509050919050565b604080516020810190915260608152600261128d836118b2565b600281111561129e5761129e614459565b146112a857600080fd5b6000828152601b6020526040902080546112c19061475d565b80601f01602080910402602001604051908101604052809291908181526020018280546112ed9061475d565b801561133a5780601f1061130f5761010080835404028352916020019161133a565b820191906000526020600020905b81548152906001019060200180831161131d57829003601f168201915b505050505080602001905181019061091d9190614a71565b61135a6134b3565b600e548114611395576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600e54811015611567576000600e82815481106113b7576113b76146f6565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600f90925260408320549193501690858585818110611401576114016146f6565b905060200201602081019061141691906142ee565b905073ffffffffffffffffffffffffffffffffffffffff811615806114a9575073ffffffffffffffffffffffffffffffffffffffff82161580159061148757508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156114a9575073ffffffffffffffffffffffffffffffffffffffff81811614155b156114e0576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146115515773ffffffffffffffffffffffffffffffffffffffff8381166000908152600f6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061155f90614725565b915050611398565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600e838360405161159c93929190614b50565b60405180910390a15050565b6115b06134b3565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b6060600061161c60026131ac565b905080841161162b578361162d565b805b935083851115611669576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061167586866146a1565b90506000808267ffffffffffffffff811115611693576116936146b4565b6040519080825280602002602001820160405280156116bc578160200160208202803683370190505b50905060005b8381101561174b5760006116d96108e08b846146e3565b90508760028111156116ed576116ed614459565b6116f6826118b2565b600281111561170757611707614459565b03611738578083858151811061171f5761171f6146f6565b60209081029190910101528361173481614725565b9450505b508061174381614725565b9150506116c2565b50828214611757578181525b979650505050505050565b6000818152601d602052604090208054606091906111ee9061475d565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015282918291829182919082906118595760608201516012546000916118459170010000000000000000000000000000000090046bffffffffffffffffffffffff16614c02565b600e549091506118559082614c56565b9150505b815160208301516040840151611870908490614c81565b6060949094015173ffffffffffffffffffffffffffffffffffffffff9a8b166000908152600f6020526040902054929b919a9499509750921694509092505050565b6000818160045b600f81101561193f577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106118f7576118f76146f6565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461192d57506000949350505050565b8061193781614725565b9150506118b9565b5081600f1a600281111561195557611955614459565b949350505050565b611966816131c9565b60008181526004602090815260409182902082516101008082018552825460ff8116151580845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c082015260029091015490921660e0830152611a70576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055611aaf600283613536565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6012546f01000000000000000000000000000000900460ff1615611b2f576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116611bb6576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff91810482168387015265010000000000810482168386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606084015260018401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a08501527801000000000000000000000000000000000000000000000000900490911660c0830152600290920154821660e082015286855260059093529220549091163314611cc9576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611cd1613542565b816040015163ffffffff161115611d14576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546018546c010000000000000000000000009091046bffffffffffffffffffffffff1690611d549082906146a1565b60185560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015611e37573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e5b9190614ca6565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611f6e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611ff26134b3565b6014546018546bffffffffffffffffffffffff909116906120149082906146a1565b601855601480547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b6020604051808303816000875af1158015612121573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061117e9190614ca6565b61214d6134b3565b601280547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001611604565b6121bf816131c9565b60008181526004602090815260409182902082516101008082018552825460ff8116158015845263ffffffff92820483169584019590955265010000000000810482169583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c082015260029091015490921660e08301526122c9576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905561230b6002836135f7565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b612344836131c9565b6000838152601c6020526040902061235d8284836147ff565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485083836040516109ba92919061491a565b73ffffffffffffffffffffffffffffffffffffffff81166123dd576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f602052604090205416331461243d576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e5460009161247491859170010000000000000000000000000000000090046bffffffffffffffffffffffff1690613603565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff1690556018549091506124de906bffffffffffffffffffffffff8316906146a1565b6018556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015612583573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125a79190614ca6565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff161080612641575060145463ffffffff7001000000000000000000000000000000009091048116908216115b15612678576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612681826131c9565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff1661010063ffffffff861690810291909117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260106020526040902054163314612761576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600f602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556010909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60008181526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821694830194909452650100000000008404811694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490921660e0830152909114612905576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612962576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b60008181526004602052604081205461091d90610100900463ffffffff16610ac9565b612a216134b3565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612aae573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ad29190614cc8565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360185484612b1f91906146a1565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401612102565b604080516101608101825260008082526020820181905291810182905260608082018190526080820183905260a0820183905260c0820183905260e082018390526101008201839052610120820192909252610140810191909152600082815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff918104821683870190815265010000000000820483168487015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606080860191825260018701546bffffffffffffffffffffffff80821660808901526c0100000000000000000000000082041660a088015278010000000000000000000000000000000000000000000000009004851660c0870152600290960154831660e08601908152875161016081018952905184168152905190921682880152519091168185015287865260079094529190932080549193830191612cdf9061475d565b80601f0160208091040260200160405190810160405280929190818152602001828054612d0b9061475d565b8015612d585780601f10612d2d57610100808354040283529160200191612d58565b820191906000526020600020905b815481529060010190602001808311612d3b57829003601f168201915b505050505081526020018260a001516bffffffffffffffffffffffff1681526020016005600086815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001826040015163ffffffff1667ffffffffffffffff1681526020018260c0015163ffffffff16815260200182608001516bffffffffffffffffffffffff168152602001826000015115158152602001601c60008681526020019081526020016000208054612e359061475d565b80601f0160208091040260200160405190810160405280929190818152602001828054612e619061475d565b8015612eae5780601f10612e8357610100808354040283529160200191612eae565b820191906000526020600020905b815481529060010190602001808311612e9157829003601f168201915b5050505050815250915050919050565b60195473ffffffffffffffffffffffffffffffffffffffff163314612f0f576040517fee5dc90100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152601d60205260409020612f288284836147ff565b50827f09a658476c5597979b9948f488ec2958cfead97bc8f46b19ca0b21cdab93cdee83836040516109ba92919061491a565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600f6020526040902054163314612fbb576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82160361300a576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526010602052604090205481169082161461117e5773ffffffffffffffffffffffffffffffffffffffff82811660008181526010602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b60408051602081019091526000815260016130d3836118b2565b60028111156130e4576130e4614459565b146130ee57600080fd5b6000828152601b6020526040902080546131079061475d565b80601f01602080910402602001604051908101604052809291908181526020018280546131339061475d565b80156131805780601f1061315557610100808354040283529160200191613180565b820191906000526020600020905b81548152906001019060200180831161316357829003601f168201915b505050505080602001905181019061091d9190614ce1565b6131a06134b3565b6131a98161380b565b50565b600061091d825490565b60006131c28383613900565b9392505050565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613226576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff908116146131a9576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613309573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061332d9190614d2e565b509450909250505060008113158061334457508142105b806133655750828015613365575061335c82426146a1565b8463ffffffff16105b15613374576016549550613378565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156133e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134079190614d2e565b509450909250505060008113158061341e57508142105b8061343f575082801561343f575061343682426146a1565b8463ffffffff16105b1561344e576017549450613452565b8094505b50505050915091565b60008061346c86896000015161392a565b90506000806134878a8a63ffffffff16858a8a60018b61396e565b90925090506134968183614c81565b93505050505b9695505050505050565b606060006131c283613d19565b60005473ffffffffffffffffffffffffffffffffffffffff163314613534576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611f65565b565b60006131c28383613d74565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111561357857613578614459565b036135f257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156135c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135ed9190614cc8565b905090565b504390565b60006131c28383613dc3565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906137ff57600081606001518561369b9190614c02565b905060006136a98583614c56565b905080836040018181516136bd9190614c81565b6bffffffffffffffffffffffff169052506136d88582614d7e565b836060018181516136e99190614c81565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b3373ffffffffffffffffffffffffffffffffffffffff82160361388a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611f65565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613917576139176146f6565b9060005260206000200154905092915050565b600061393d63ffffffff84166014614db2565b613948836001614def565b6139579060ff16611d4c614db2565b61396490620138806146e3565b6131c291906146e3565b6000806000896080015161ffff16876139879190614db2565b90508380156139955750803a105b1561399d57503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156139d3576139d3614459565b03613b32576040805160008152602081019091528515613a3157600036604051806080016040528060488152602001614e9660489139604051602001613a1b93929190614e08565b6040516020818303038152906040529050613a99565b601554613a4d90640100000000900463ffffffff166004614e2f565b63ffffffff1667ffffffffffffffff811115613a6b57613a6b6146b4565b6040519080825280601f01601f191660200182016040528015613a95576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90613ae9908490600401614379565b602060405180830381865afa158015613b06573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b2a9190614cc8565b915050613bde565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115613b6657613b66614459565b03613bde57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613bb7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bdb9190614cc8565b90505b84613bfa57808b6080015161ffff16613bf79190614db2565b90505b613c0861ffff871682614e52565b905060008782613c188c8e6146e3565b613c229086614db2565b613c2c91906146e3565b613c3e90670de0b6b3a7640000614db2565b613c489190614e52565b905060008c6040015163ffffffff1664e8d4a51000613c679190614db2565b898e6020015163ffffffff16858f88613c809190614db2565b613c8a91906146e3565b613c9890633b9aca00614db2565b613ca29190614db2565b613cac9190614e52565b613cb691906146e3565b90506b033b2e3c9fd0803ce8000000613ccf82846146e3565b1115613d07576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561126757602002820191906000526020600020905b815481526020019060010190808311613d555750505050509050919050565b6000818152600183016020526040812054613dbb5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561091d565b50600061091d565b60008181526001830160205260408120548015613eac576000613de76001836146a1565b8554909150600090613dfb906001906146a1565b9050818114613e60576000866000018281548110613e1b57613e1b6146f6565b9060005260206000200154905080876000018481548110613e3e57613e3e6146f6565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613e7157613e71614e66565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061091d565b600091505061091d565b5092915050565b60008060408385031215613ed057600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613f1757835183529284019291840191600101613efb565b50909695505050505050565b600080600060408486031215613f3857600080fd5b83359250602084013567ffffffffffffffff80821115613f5757600080fd5b818601915086601f830112613f6b57600080fd5b813581811115613f7a57600080fd5b876020828501011115613f8c57600080fd5b6020830194508093505050509250925092565b600060208284031215613fb157600080fd5b5035919050565b63ffffffff811681146131a957600080fd5b600060208284031215613fdc57600080fd5b81356131c281613fb8565b600081518084526020808501945080840160005b8381101561402d57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101613ffb565b509495945050505050565b805163ffffffff16825260006101a0602083015161405e602086018263ffffffff169052565b506040830151614076604086018263ffffffff169052565b50606083015161408d606086018262ffffff169052565b5060808301516140a3608086018261ffff169052565b5060a08301516140c360a08601826bffffffffffffffffffffffff169052565b5060c08301516140db60c086018263ffffffff169052565b5060e08301516140f360e086018263ffffffff169052565b506101008381015163ffffffff1690850152610120808401519085015261014080840151908501526101608084015173ffffffffffffffffffffffffffffffffffffffff16908501526101808084015181860183905261349c83870182613fe7565b855163ffffffff16815260006101c0602088015161418360208501826bffffffffffffffffffffffff169052565b506040880151604084015260608801516141ad60608501826bffffffffffffffffffffffff169052565b506080880151608084015260a08801516141cf60a085018263ffffffff169052565b5060c08801516141e760c085018263ffffffff169052565b5060e088015160e08401526101008089015161420a8286018263ffffffff169052565b505061012088810151151590840152610140830181905261422d81840188614038565b90508281036101608401526142428187613fe7565b90508281036101808401526142578186613fe7565b91505061349c6101a083018460ff169052565b73ffffffffffffffffffffffffffffffffffffffff811681146131a957600080fd5b6000806040838503121561429f57600080fd5b82356142aa8161426a565b91506020830135600481106142be57600080fd5b809150509250929050565b600080604083850312156142dc57600080fd5b8235915060208301356142be8161426a565b60006020828403121561430057600080fd5b81356131c28161426a565b60005b8381101561432657818101518382015260200161430e565b50506000910152565b6000815180845261434781602086016020860161430b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006131c2602083018461432f565b6020815260008251602080840152611955604084018261432f565b600080602083850312156143ba57600080fd5b823567ffffffffffffffff808211156143d257600080fd5b818501915085601f8301126143e657600080fd5b8135818111156143f557600080fd5b8660208260051b850101111561440a57600080fd5b60209290920196919550909350505050565b60008060006060848603121561443157600080fd5b833592506020840135915060408401356003811061444e57600080fd5b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600381106131a9576131a9614459565b602081016144a583614488565b91905290565b600080604083850312156144be57600080fd5b82356144c98161426a565b915060208301356142be8161426a565b600080604083850312156144ec57600080fd5b8235915060208301356142be81613fb8565b6020815261452560208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161454e604084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604083015163ffffffff8116606084015250606083015161016080608085015261457d61018085018361432f565b9150608085015161459e60a08601826bffffffffffffffffffffffff169052565b5060a085015173ffffffffffffffffffffffffffffffffffffffff811660c08601525060c085015167ffffffffffffffff811660e08601525060e08501516101006145f08187018363ffffffff169052565b8601519050610120614611868201836bffffffffffffffffffffffff169052565b86015190506101406146268682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061349c838261432f565b60208101600483106144a5576144a5614459565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561091d5761091d614672565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b8082018082111561091d5761091d614672565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361475657614756614672565b5060010190565b600181811c9082168061477157607f821691505b6020821081036147aa577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156147fa57600081815260208120601f850160051c810160208610156147d75750805b601f850160051c820191505b818110156147f6578281556001016147e3565b5050505b505050565b67ffffffffffffffff831115614817576148176146b4565b61482b83614825835461475d565b836147b0565b6000601f84116001811461487d57600085156148475750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355614913565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156148cc57868501358255602094850194600190920191016148ac565b5086821015614907577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6040516020810167ffffffffffffffff8111828210171561498a5761498a6146b4565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156149d7576149d76146b4565b604052919050565b600060c082840312156149f157600080fd5b60405160c0810181811067ffffffffffffffff82111715614a1457614a146146b4565b6040528251614a228161426a565b8152602083015160ff81168114614a3857600080fd5b8060208301525060408301516040820152606083015160608201526080830151608082015260a083015160a08201528091505092915050565b60006020808385031215614a8457600080fd5b825167ffffffffffffffff80821115614a9c57600080fd5b8185019150828287031215614ab057600080fd5b614ab8614967565b825182811115614ac757600080fd5b80840193505086601f840112614adc57600080fd5b825182811115614aee57614aee6146b4565b614b1e857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601614990565b92508083528785828601011115614b3457600080fd5b614b438186850187870161430b565b5090815295945050505050565b6000604082016040835280865480835260608501915087600052602092508260002060005b82811015614ba757815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101614b75565b505050838103828501528481528590820160005b86811015614bf6578235614bce8161426a565b73ffffffffffffffffffffffffffffffffffffffff1682529183019190830190600101614bbb565b50979650505050505050565b6bffffffffffffffffffffffff828116828216039080821115613eb657613eb6614672565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680614c7557614c75614c27565b92169190910492915050565b6bffffffffffffffffffffffff818116838216019080821115613eb657613eb6614672565b600060208284031215614cb857600080fd5b815180151581146131c257600080fd5b600060208284031215614cda57600080fd5b5051919050565b600060208284031215614cf357600080fd5b614cfb614967565b8251614d0681613fb8565b81529392505050565b805169ffffffffffffffffffff81168114614d2957600080fd5b919050565b600080600080600060a08688031215614d4657600080fd5b614d4f86614d0f565b9450602086015193506040860151925060608601519150614d7260808701614d0f565b90509295509295909350565b60006bffffffffffffffffffffffff80831681851681830481118215151615614da957614da9614672565b02949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614dea57614dea614672565b500290565b60ff818116838216019081111561091d5761091d614672565b828482376000838201600081528351614e2581836020880161430b565b0195945050505050565b600063ffffffff80831681851681830481118215151615614da957614da9614672565b600082614e6157614e61614c27565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetBlockTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21BlockTriggerConfig, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getBlockTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21BlockTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21BlockTriggerConfig)).(*KeeperRegistryBase21BlockTriggerConfig)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetBlockTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21BlockTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetBlockTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetBlockTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21BlockTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetBlockTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetCronTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21CronTriggerConfig, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getCronTriggerConfig", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21CronTriggerConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21CronTriggerConfig)).(*KeeperRegistryBase21CronTriggerConfig)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetCronTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21CronTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetCronTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetCronTriggerConfig(upkeepId *big.Int) (KeeperRegistryBase21CronTriggerConfig, error) {
	return _KeeperRegistryLogicB.Contract.GetCronTriggerConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB.CallOpts, gasLimit)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetMaxPaymentForGas(gasLimit uint32) (*big.Int, error) {
	return _KeeperRegistryLogicB.Contract.GetMaxPaymentForGas(&_KeeperRegistryLogicB.CallOpts, gasLimit)
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetUpkeepAdminOffchainConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getUpkeepAdminOffchainConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetUpkeepAdminOffchainConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepAdminOffchainConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetUpkeepAdminOffchainConfig(upkeepId *big.Int) ([]byte, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepAdminOffchainConfig(&_KeeperRegistryLogicB.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCaller) GetUpkeepManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicB.contract.Call(opts, &out, "getUpkeepManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) GetUpkeepManager() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepManager(&_KeeperRegistryLogicB.CallOpts)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBCallerSession) GetUpkeepManager() (common.Address, error) {
	return _KeeperRegistryLogicB.Contract.GetUpkeepManager(&_KeeperRegistryLogicB.CallOpts)
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepAdminOffchainConfig(opts *bind.TransactOpts, upkeepId *big.Int, newAdminOffchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepAdminOffchainConfig", upkeepId, newAdminOffchainConfig)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepAdminOffchainConfig(upkeepId *big.Int, newAdminOffchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepAdminOffchainConfig(&_KeeperRegistryLogicB.TransactOpts, upkeepId, newAdminOffchainConfig)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepAdminOffchainConfig(upkeepId *big.Int, newAdminOffchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepAdminOffchainConfig(&_KeeperRegistryLogicB.TransactOpts, upkeepId, newAdminOffchainConfig)
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepManager(opts *bind.TransactOpts, newUpkeepManager common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepManager", newUpkeepManager)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepManager(newUpkeepManager common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepManager(&_KeeperRegistryLogicB.TransactOpts, newUpkeepManager)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepManager(newUpkeepManager common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepManager(&_KeeperRegistryLogicB.TransactOpts, newUpkeepManager)
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

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactor) SetUpkeepPipelineData(opts *bind.TransactOpts, id *big.Int, newPipelineData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.contract.Transact(opts, "setUpkeepPipelineData", id, newPipelineData)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBSession) SetUpkeepPipelineData(id *big.Int, newPipelineData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepPipelineData(&_KeeperRegistryLogicB.TransactOpts, id, newPipelineData)
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBTransactorSession) SetUpkeepPipelineData(id *big.Int, newPipelineData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicB.Contract.SetUpkeepPipelineData(&_KeeperRegistryLogicB.TransactOpts, id, newPipelineData)
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

type KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepAdminOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepAdminOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepAdminOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepAdminOffchainConfigSet struct {
	Id                  *big.Int
	AdminOffchainConfig []byte
	Raw                 types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepAdminOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepAdminOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepAdminOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepAdminOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepAdminOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepAdminOffchainConfigSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepAdminOffchainConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepAdminOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicBUpkeepAdminOffchainConfigSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepAdminOffchainConfigSet", log); err != nil {
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

type KeeperRegistryLogicBUpkeepPipelineDataSetIterator struct {
	Event *KeeperRegistryLogicBUpkeepPipelineDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicBUpkeepPipelineDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicBUpkeepPipelineDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicBUpkeepPipelineDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicBUpkeepPipelineDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicBUpkeepPipelineDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicBUpkeepPipelineDataSet struct {
	Id              *big.Int
	NewPipelineData []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) FilterUpkeepPipelineDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPipelineDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.FilterLogs(opts, "UpkeepPipelineDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicBUpkeepPipelineDataSetIterator{contract: _KeeperRegistryLogicB.contract, event: "UpkeepPipelineDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) WatchUpkeepPipelineDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPipelineDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicB.contract.WatchLogs(opts, "UpkeepPipelineDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicBUpkeepPipelineDataSet)
				if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPipelineDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicB *KeeperRegistryLogicBFilterer) ParseUpkeepPipelineDataSet(log types.Log) (*KeeperRegistryLogicBUpkeepPipelineDataSet, error) {
	event := new(KeeperRegistryLogicBUpkeepPipelineDataSet)
	if err := _KeeperRegistryLogicB.contract.UnpackLog(event, "UpkeepPipelineDataSet", log); err != nil {
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
	case _KeeperRegistryLogicB.abi.Events["UpkeepAdminOffchainConfigSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepAdminOffchainConfigSet(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicB.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepCanceled(log)
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
	case _KeeperRegistryLogicB.abi.Events["UpkeepPipelineDataSet"].ID:
		return _KeeperRegistryLogicB.ParseUpkeepPipelineDataSet(log)
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

func (KeeperRegistryLogicBUpkeepAdminOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x09a658476c5597979b9948f488ec2958cfead97bc8f46b19ca0b21cdab93cdee")
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

func (KeeperRegistryLogicBUpkeepPipelineDataSet) Topic() common.Hash {
	return common.HexToHash("0x787b2e556c364148324d7dc0cee4322a048aaf798bda184136c8636c72ed3056")
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

	GetBlockTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21BlockTriggerConfig, error)

	GetCronTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21CronTriggerConfig, error)

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21LogTriggerConfig, error)

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

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (KeeperRegistryBase21UpkeepInfo, error)

	GetUpkeepAdminOffchainConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepManager(opts *bind.CallOpts) (common.Address, error)

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

	SetUpkeepAdminOffchainConfig(opts *bind.TransactOpts, upkeepId *big.Int, newAdminOffchainConfig []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepManager(opts *bind.TransactOpts, newUpkeepManager common.Address) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	SetUpkeepPipelineData(opts *bind.TransactOpts, id *big.Int, newPipelineData []byte) (*types.Transaction, error)

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

	FilterUpkeepAdminOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepAdminOffchainConfigSetIterator, error)

	WatchUpkeepAdminOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepAdminOffchainConfigSet(log types.Log) (*KeeperRegistryLogicBUpkeepAdminOffchainConfigSet, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicBUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicBUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicBUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicBUpkeepCanceled, error)

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

	FilterUpkeepPipelineDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicBUpkeepPipelineDataSetIterator, error)

	WatchUpkeepPipelineDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicBUpkeepPipelineDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPipelineDataSet(log types.Log) (*KeeperRegistryLogicBUpkeepPipelineDataSet, error)

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
