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

// Config1_3 is an auto generated low-level Go binding around an user-defined struct.
type Config1_3 struct {
	PaymentPremiumPPB    uint32
	FlatFeeMicroLink     uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	MinUpkeepSpend       *big.Int
	MaxPerformGas        uint32
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Transcoder           common.Address
	Registrar            common.Address
}

// State1_3 is an auto generated low-level Go binding around an user-defined struct.
type State1_3 struct {
	Nonce               uint32
	OwnerLinkBalance    *big.Int
	ExpectedLinkBalance *big.Int
	NumUpkeeps          *big.Int
}

// KeeperRegistry13MetaData contains all meta data concerning the KeeperRegistry13 contract.
var KeeperRegistry13MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogic1_3\",\"name\":\"keeperRegistryLogic\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_3\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig1_3\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"ARB_NITRO_ORACLE\",\"outputs\":[{\"internalType\":\"contractArbGasInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"KEEPER_REGISTRY_LOGIC\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMISM_ORACLE\",\"outputs\":[{\"internalType\":\"contractOVM_GasPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAYMENT_MODEL\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase1_3.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REGISTRY_GAS_OVERHEAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getKeeperInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase1_3.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"}],\"internalType\":\"structState1_3\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_3\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_3\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase1_3.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6101806040527f420000000000000000000000000000000000000f00000000000000000000000060e0526c6c000000000000000000000000610100523480156200004857600080fd5b5060405162003f6538038062003f658339810160408190526200006b91620008a7565b816001600160a01b0316638811cbe86040518163ffffffff1660e01b815260040160206040518083038186803b158015620000a557600080fd5b505afa158015620000ba573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000e09190620009ce565b826001600160a01b0316635077b2106040518163ffffffff1660e01b815260040160206040518083038186803b1580156200011a57600080fd5b505afa1580156200012f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001559190620009f1565b836001600160a01b0316631b6b6d236040518163ffffffff1660e01b815260040160206040518083038186803b1580156200018f57600080fd5b505afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000880565b846001600160a01b031663ad1783616040518163ffffffff1660e01b815260040160206040518083038186803b1580156200020457600080fd5b505afa15801562000219573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200023f919062000880565b856001600160a01b0316634584a4196040518163ffffffff1660e01b815260040160206040518083038186803b1580156200027957600080fd5b505afa1580156200028e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002b4919062000880565b33806000816200030b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200033e576200033e816200041b565b5050600160029081556003805460ff1916905586915081111562000366576200036662000b42565b6101208160028111156200037e576200037e62000b42565b60f81b9052506101408490526001600160a01b0383161580620003a857506001600160a01b038216155b80620003bb57506001600160a01b038116155b15620003da57604051637138356f60e01b815260040160405180910390fd5b6001600160601b0319606093841b811660805291831b821660a052821b811660c0529085901b16610160525062000413905081620004c7565b505062000b71565b6001600160a01b038116331415620004765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000302565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620004d1620007bc565b600e5460e082015163ffffffff918216911610156200050357604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600d60010160049054906101000a900463ffffffff1663ffffffff16815250600d60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600f81905550806101200151601081905550806101400151601360006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601460006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051620007b1919062000a0b565b60405180910390a150565b6000546001600160a01b03163314620008185760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000302565b565b8051620008278162000b58565b919050565b805161ffff811681146200082757600080fd5b805162ffffff811681146200082757600080fd5b805163ffffffff811681146200082757600080fd5b80516001600160601b03811681146200082757600080fd5b6000602082840312156200089357600080fd5b8151620008a08162000b58565b9392505050565b6000808284036101a0811215620008bd57600080fd5b8351620008ca8162000b58565b9250610180601f198201811315620008e157600080fd5b620008eb62000b0a565b9150620008fb6020860162000853565b82526200090b6040860162000853565b60208301526200091e606086016200083f565b6040830152620009316080860162000853565b60608301526200094460a086016200083f565b60808301526200095760c086016200082c565b60a08301526200096a60e0860162000868565b60c08301526101006200097f81870162000853565b60e084015261012080870151828501526101409150818701518185015250610160620009ad8188016200081a565b82850152620009be8388016200081a565b9084015250929590945092505050565b600060208284031215620009e157600080fd5b815160038110620008a057600080fd5b60006020828403121562000a0457600080fd5b5051919050565b815163ffffffff1681526101808101602083015162000a32602084018263ffffffff169052565b50604083015162000a4a604084018262ffffff169052565b50606083015162000a63606084018263ffffffff169052565b50608083015162000a7b608084018262ffffff169052565b5060a083015162000a9260a084018261ffff169052565b5060c083015162000aae60c08401826001600160601b03169052565b5060e083015162000ac760e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60405161018081016001600160401b038111828210171562000b3c57634e487b7160e01b600052604160045260246000fd5b60405290565b634e487b7160e01b600052602160045260246000fd5b6001600160a01b038116811462000b6e57600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160601c6101005160601c6101205160f81c610140516101605160601c61333a62000c2b600039600081816102a0015261092f0152600081816104db0152611f5601526000818161063e01528181611fa901526120ff01526000818161059601526121370152600081816105ca015261206e0152600081816104850152611ce30152600081816107810152611db7015260008181610398015261112a015261333a6000f3fe6080604052600436106102415760003560e01c80638811cbe81161012f578063b148ab6b116100b1578063b148ab6b146107c3578063b657bc9c146107de578063b79550be1461045e578063b7fdb436146107fe578063c41b813a1461081e578063c7c3a19a1461084f578063c8048022146107c3578063da5c674114610884578063eb5dcd6c14610739578063ef47a0ce146108a4578063f2fde38b146108c4578063faa3e996146108e457610250565b80638811cbe81461062c5780638da5cb5b1461066d5780638e86139b1461068b57806393f0c1fc146106a6578063948108f7146106de5780639fab4386146106f9578063a4c0ed3614610719578063a710b22114610739578063a72aa27e14610754578063ad1783611461076f578063b121e147146107a357610250565b80635077b210116101c35780635077b210146104c95780635165f2f51461050b5780635c975abb1461052b578063744bfe611461036b57806379ba50971461054f5780637bbaf1ea146105645780637d9b97e01461045e5780637f37618e146105845780638456cb591461045e578063850cce34146105b857806385c1b0ba146105ec5780638765ecbe1461060c57610250565b806306e3b632146102585780630852c7c91461028e578063181f5a77146102da5780631865c57d14610327578063187256e81461034b5780631a2af0111461036b5780631b6b6d23146103865780631e12b8a5146103ba5780633f4ba83a1461045e5780634584a4191461047357806348013d7b146104a757610250565b366102505761024e61092a565b005b61024e61092a565b34801561026457600080fd5b50610278610273366004612b78565b610955565b6040516102859190612dfc565b60405180910390f35b34801561029a57600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610285565b3480156102e657600080fd5b5061031a6040518060400160405280601481526020017304b6565706572526567697374727920312e332e360641b81525081565b6040516102859190612e6f565b34801561033357600080fd5b5061033c610a37565b60405161028593929190612f0a565b34801561035757600080fd5b5061024e6103663660046127fe565b610c91565b34801561037757600080fd5b5061024e610366366004612b0a565b34801561039257600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b3480156103c657600080fd5b506104306103d53660046127b0565b6001600160a01b0390811660009081526008602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b031692810183905260019091015460ff16151592018290529192909190565b604080516001600160a01b03909416845291151560208401526001600160601b031690820152606001610285565b34801561046a57600080fd5b5061024e610c9d565b34801561047f57600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b3480156104b357600080fd5b506104bc600181565b6040516102859190612ee7565b3480156104d557600080fd5b506104fd7f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610285565b34801561051757600080fd5b5061024e610526366004612ad8565b610ca5565b34801561053757600080fd5b5060035460ff165b6040519015158152602001610285565b34801561055b57600080fd5b5061024e610dc3565b34801561057057600080fd5b5061053f61057f366004612b2d565b610e72565b34801561059057600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b3480156105c457600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b3480156105f857600080fd5b5061024e610607366004612966565b610ed0565b34801561061857600080fd5b5061024e610627366004612ad8565b610edd565b34801561063857600080fd5b506106607f000000000000000000000000000000000000000000000000000000000000000081565b6040516102859190612ed3565b34801561067957600080fd5b506000546001600160a01b03166102c2565b34801561069757600080fd5b5061024e6103663660046129b9565b3480156106b257600080fd5b506106c66106c1366004612ad8565b611002565b6040516001600160601b039091168152602001610285565b3480156106ea57600080fd5b5061024e610366366004612bbd565b34801561070557600080fd5b5061024e610714366004612b2d565b611020565b34801561072557600080fd5b5061024e610734366004612839565b61111f565b34801561074557600080fd5b5061024e6103663660046127cb565b34801561076057600080fd5b5061024e610366366004612b9a565b34801561077b57600080fd5b506102c27f000000000000000000000000000000000000000000000000000000000000000081565b3480156107af57600080fd5b5061024e6107be3660046127b0565b611289565b3480156107cf57600080fd5b5061024e6107be366004612ad8565b3480156107ea57600080fd5b506106c66107f9366004612ad8565b611294565b34801561080a57600080fd5b5061024e610819366004612907565b6112b5565b34801561082a57600080fd5b5061083e610839366004612b0a565b6112c3565b604051610285959493929190612e82565b34801561085b57600080fd5b5061086f61086a366004612ad8565b6112e5565b60405161028599989796959493929190612d7b565b34801561089057600080fd5b506104fd61089f366004612892565b61147b565b3480156108b057600080fd5b5061024e6108bf3660046129fa565b61148e565b3480156108d057600080fd5b5061024e6108df3660046127b0565b61177e565b3480156108f057600080fd5b5061091d6108ff3660046127b0565b6001600160a01b03166000908152600c602052604090205460ff1690565b6040516102859190612eb9565b6109537f000000000000000000000000000000000000000000000000000000000000000061178f565b565b6060600061096360056117b3565b905080841061098557604051631390f2a160e01b815260040160405180910390fd5b8261099757610994848261307d565b92505b6000836001600160401b038111156109b1576109b1613196565b6040519080825280602002602001820160405280156109da578160200160208202803683370190505b50905060005b84811015610a2c576109fd6109f58288612ff9565b6005906117bd565b828281518110610a0f57610a0f613180565b602090810291909101015280610a2481613123565b9150506109e0565b509150505b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600d5463ffffffff8082168352600160201b808304821660208086019190915262ffffff600160401b8504811686880152600160581b85048416606087810191909152600160781b8604909116608087015261ffff600160901b86041660a08701526001600160601b03600160a01b909504851660c0870152600e5480851660e088015292909204909216610100850181905287526011549092169086015260125492850192909252610b7a60056117b3565b606080860191909152815163ffffffff908116855260208084015182168187015260408085015162ffffff90811682890152858501518416948801949094526080808601519094169387019390935260a08085015161ffff169087015260c0808501516001600160601b03169087015260e08085015190921691860191909152600f546101008601526010546101208601526013546001600160a01b0390811661014087015260145416610160860152600480548351818402810184019094528084528793879390918391830182828015610c7e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610c60575b5050505050905093509350935050909192565b610c9961092a565b5050565b61095361092a565b60008181526007602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e0820152610d4a816117d0565b8060e00151610d6c576040516306e229e160e21b815260040160405180910390fd5b6000828152600760205260409020600201805460ff60e01b19169055610d93600583611831565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6001546001600160a01b03163314610e1b5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610e7c61183d565b610ec8610ec3338686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250611883915050565b61194e565b949350505050565b610ed861092a565b505050565b60008181526007602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e0820152610f82816117d0565b8060e0015115610fa557604051631452db0960e21b815260040160405180910390fd5b6000828152600760205260409020600201805460ff60e01b1916600160e01b179055610fd2600583611ca4565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b600080600061100f611cb0565b91509150610ec88483836000611e91565b60008381526007602090815260409182902082516101008101845281546001600160601b0380821683526001600160a01b03600160601b92839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff8082166080850152600160201b82041660a0840152600160401b810490911660c083015260ff600160e01b90910416151560e08201526110c5816117d0565b6000848152600b602052604090206110de908484612605565b50837f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8484604051611111929190612e40565b60405180910390a250505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146111685760405163c8bad78d60e01b815260040160405180910390fd5b6020811461118957604051630dfe930960e41b815260040160405180910390fd5b600061119782840184612ad8565b600081815260076020526040902060020154909150600160201b900463ffffffff908116146111d957604051634e0041d160e11b815260040160405180910390fd5b6000818152600760205260409020546111fc9085906001600160601b0316613011565b600082815260076020526040902080546001600160601b0319166001600160601b0392909216919091179055601254611236908590612ff9565b6012556040516001600160601b03851681526001600160a01b0386169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b61129161092a565b50565b600081815260076020526040812060020154610a319063ffffffff16611002565b6112bd61092a565b50505050565b60606000806000806112d3612277565b6112db61092a565b9295509295909350565b600081815260076020908152604080832081516101008101835281546001600160601b038082168352600160601b918290046001600160a01b039081168488019081526001860154928316858801908152939092048116606080860191825260029096015463ffffffff80821660808801819052600160201b830490911660a08801908152600160401b830490941660c08801819052600160e01b90920460ff16151560e088019081528c8c52600b909a52978a2086519451925193519551995181548c9b999a8c9a8b9a8b9a8b9a8b9a8b9a93999895979596919593949093909187906113d2906130e8565b80601f01602080910402602001604051908101604052809291908181526020018280546113fe906130e8565b801561144b5780601f106114205761010080835404028352916020019161144b565b820191906000526020600020905b81548152906001019060200180831161142e57829003601f168201915b505050505096508263ffffffff169250995099509950995099509950995099509950509193959799909294969850565b600061148561092a565b95945050505050565b611496612296565b600e5460e082015163ffffffff918216911610156114c757604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600d60010160049054906101000a900463ffffffff1663ffffffff16815250600d60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600f81905550806101200151601081905550806101400151601360006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601460006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325816040516117739190612efb565b60405180910390a150565b611786612296565b611291816122e9565b3660008037600080366000845af43d6000803e8080156117ae573d6000f35b3d6000fd5b6000610a31825490565b60006117c9838361238d565b9392505050565b80606001516001600160a01b0316336001600160a01b0316146118065760405163523e0b8360e11b815260040160405180910390fd5b60a081015163ffffffff9081161461129157604051634e0041d160e11b815260040160405180910390fd5b60006117c983836123b7565b60035460ff16156109535760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b6044820152606401610e12565b6118cc6040518060e0016040528060006001600160a01b031681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526007602052604081206002015463ffffffff1690806118ee611cb0565b91509150600061190084848489611e91565b6040805160e0810182526001600160a01b03909b168b5260208b0199909952978901969096526001600160601b039096166060880152608087019190915260a086015250505060c082015290565b60006002805414156119a25760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610e12565b600280805560208084015160009081526007825260409081902081516101008101835281546001600160601b0380821683526001600160a01b03600160601b92839004811696840196909652600184015490811694830194909452909204831660608301529092015463ffffffff8082166080850152600160201b82041660a08401819052600160401b820490921660c084015260ff600160e01b90910416151560e08301524310611a6757604051634e0041d160e11b815260040160405180910390fd5b611a7a8184600001518560600151612406565b60005a90506000634585e33b60e01b8560400151604051602401611a9e9190612e6f565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b0383818316178352505050509050611ae685608001518460c00151836124c6565b93505a611af3908361307d565b91506000611b0c838760a001518860c001516001611e91565b602080880151600090815260079091526040902054909150611b389082906001600160601b0316613094565b6020878101805160009081526007909252604080832080546001600160601b0319166001600160601b0395861617905590518252902060010154611b7e91839116613011565b60208781018051600090815260078352604080822060010180546001600160601b0319166001600160601b039687161790558a519251825280822080548616600160601b6001600160a01b03958616021790558a5190921681526008909252902054611bf3918391600160a01b900416613011565b6008600088600001516001600160a01b03166001600160a01b0316815260200190815260200160002060000160146101000a8154816001600160601b0302191690836001600160601b0316021790555085600001516001600160a01b031685151587602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848a60400151604051611c8e929190612f9e565b60405180910390a4505050506001600255919050565b60006117c98383612512565b6000806000600d600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611d3a57600080fd5b505afa158015611d4e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d729190612be0565b509450909250849150508015611d965750611d8d824261307d565b8463ffffffff16105b80611da2575060008113155b15611db157600f549550611db5565b8095505b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015611e0e57600080fd5b505afa158015611e22573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e469190612be0565b509450909250849150508015611e6a5750611e61824261307d565b8463ffffffff16105b80611e76575060008113155b15611e85576010549450611e89565b8094505b505050509091565b6040805161012081018252600d5463ffffffff8082168352600160201b80830482166020850152600160401b830462ffffff90811695850195909552600160581b830482166060850152600160781b83049094166080840152600160901b820461ffff1660a08401819052600160a01b9092046001600160601b031660c0840152600e5480821660e0850152939093049092166101008201526000918290611f39908761305e565b9050838015611f475750803a105b15611f4f57503a5b6000611f7b7f000000000000000000000000000000000000000000000000000000000000000089612ff9565b611f85908361305e565b8351909150600090611fa19063ffffffff16633b9aca00612ff9565b9050600060027f00000000000000000000000000000000000000000000000000000000000000006002811115611fd957611fd9613154565b14156120fb576040805160008152602081019091528715612038576000366040518060800160405280604881526020016131ad6048913960405160200161202293929190612d54565b6040516020818303038152906040529050612057565b60405180610140016040528061011081526020016131f5610110913990505b6040516324ca470760e11b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906349948e0e906120a3908490600401612e6f565b60206040518083038186803b1580156120bb57600080fd5b505afa1580156120cf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120f39190612af1565b9150506121c9565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561212f5761212f613154565b14156121c9577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561218e57600080fd5b505afa1580156121a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121c69190612af1565b90505b866121e557808560a0015161ffff166121e2919061305e565b90505b6000856020015163ffffffff1664e8d4a51000612202919061305e565b898461220e8588612ff9565b61221c90633b9aca0061305e565b612226919061305e565b612230919061303c565b61223a9190612ff9565b90506b033b2e3c9fd0803ce80000008111156122695760405163156baa3d60e11b815260040160405180910390fd5b9a9950505050505050505050565b32156109535760405163b60ac5db60e01b815260040160405180910390fd5b6000546001600160a01b031633146109535760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610e12565b6001600160a01b03811633141561233c5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610e12565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106123a4576123a4613180565b9060005260206000200154905092915050565b60008181526001830160205260408120546123fe57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610a31565b506000610a31565b8260e001511561242957604051631452db0960e21b815260040160405180910390fd5b6001600160a01b03821660009081526008602052604090206001015460ff16612465576040516319f759fb60e31b815260040160405180910390fd5b82516001600160601b03168111156124905760405163356680b760e01b815260040160405180910390fd5b816001600160a01b031683602001516001600160a01b03161415610ed857604051621af04160e61b815260040160405180910390fd5b60005a6113888110156124d857600080fd5b6113888103905084604082048203116124f057600080fd5b50823b6124fc57600080fd5b60008083516020850160008789f1949350505050565b600081815260018301602052604081205480156125fb57600061253660018361307d565b855490915060009061254a9060019061307d565b90508181146125af57600086600001828154811061256a5761256a613180565b906000526020600020015490508087600001848154811061258d5761258d613180565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806125c0576125c061316a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610a31565b6000915050610a31565b828054612611906130e8565b90600052602060002090601f0160209004810192826126335760008555612679565b82601f1061264c5782800160ff19823516178555612679565b82800160010185558215612679579182015b8281111561267957823582559160200191906001019061265e565b50612685929150612689565b5090565b5b80821115612685576000815560010161268a565b80356001600160a01b03811681146126b557600080fd5b919050565b60008083601f8401126126cc57600080fd5b5081356001600160401b038111156126e357600080fd5b6020830191508360208260051b85010111156126fe57600080fd5b9250929050565b60008083601f84011261271757600080fd5b5081356001600160401b0381111561272e57600080fd5b6020830191508360208285010111156126fe57600080fd5b803561ffff811681146126b557600080fd5b803562ffffff811681146126b557600080fd5b803563ffffffff811681146126b557600080fd5b805169ffffffffffffffffffff811681146126b557600080fd5b80356001600160601b03811681146126b557600080fd5b6000602082840312156127c257600080fd5b6117c98261269e565b600080604083850312156127de57600080fd5b6127e78361269e565b91506127f56020840161269e565b90509250929050565b6000806040838503121561281157600080fd5b61281a8361269e565b915060208301356004811061282e57600080fd5b809150509250929050565b6000806000806060858703121561284f57600080fd5b6128588561269e565b93506020850135925060408501356001600160401b0381111561287a57600080fd5b61288687828801612705565b95989497509550505050565b6000806000806000608086880312156128aa57600080fd5b6128b38661269e565b94506128c16020870161276b565b93506128cf6040870161269e565b925060608601356001600160401b038111156128ea57600080fd5b6128f688828901612705565b969995985093965092949392505050565b6000806000806040858703121561291d57600080fd5b84356001600160401b038082111561293457600080fd5b612940888389016126ba565b9096509450602087013591508082111561295957600080fd5b50612886878288016126ba565b60008060006040848603121561297b57600080fd5b83356001600160401b0381111561299157600080fd5b61299d868287016126ba565b90945092506129b090506020850161269e565b90509250925092565b600080602083850312156129cc57600080fd5b82356001600160401b038111156129e257600080fd5b6129ee85828601612705565b90969095509350505050565b60006101808284031215612a0d57600080fd5b612a15612fc2565b612a1e8361276b565b8152612a2c6020840161276b565b6020820152612a3d60408401612758565b6040820152612a4e6060840161276b565b6060820152612a5f60808401612758565b6080820152612a7060a08401612746565b60a0820152612a8160c08401612799565b60c0820152612a9260e0840161276b565b60e082015261010083810135908201526101208084013590820152610140612abb81850161269e565b90820152610160612acd84820161269e565b908201529392505050565b600060208284031215612aea57600080fd5b5035919050565b600060208284031215612b0357600080fd5b5051919050565b60008060408385031215612b1d57600080fd5b823591506127f56020840161269e565b600080600060408486031215612b4257600080fd5b8335925060208401356001600160401b03811115612b5f57600080fd5b612b6b86828701612705565b9497909650939450505050565b60008060408385031215612b8b57600080fd5b50508035926020909101359150565b60008060408385031215612bad57600080fd5b823591506127f56020840161276b565b60008060408385031215612bd057600080fd5b823591506127f560208401612799565b600080600080600060a08688031215612bf857600080fd5b612c018661277f565b9450602086015193506040860151925060608601519150612c246080870161277f565b90509295509295909350565b60008151808452612c488160208601602086016130bc565b601f01601f19169290920160200192915050565b805163ffffffff1682526020810151612c7d602084018263ffffffff169052565b506040810151612c94604084018262ffffff169052565b506060810151612cac606084018263ffffffff169052565b506080810151612cc3608084018262ffffff169052565b5060a0810151612cd960a084018261ffff169052565b5060c0810151612cf460c08401826001600160601b03169052565b5060e0810151612d0c60e084018263ffffffff169052565b5061010081810151908301526101208082015190830152610140808201516001600160a01b038116828501525050610160818101516001600160a01b038116848301526112bd565b828482376000838201600081528351612d718183602088016130bc565b0195945050505050565b6001600160a01b038a8116825263ffffffff8a16602083015261012060408301819052600091612dad8483018c612c30565b6001600160601b039a8b16606086015298811660808501529690961660a0830152506001600160401b039390931660c0840152941660e082015292151561010090930192909252949350505050565b6020808252825182820181905260009190848201906040850190845b81811015612e3457835183529284019291840191600101612e18565b50909695505050505050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b6020815260006117c96020830184612c30565b60a081526000612e9560a0830188612c30565b90508560208301528460408301528360608301528260808301529695505050505050565b6020810160048310612ecd57612ecd613154565b91905290565b6020810160038310612ecd57612ecd613154565b6020810160028310612ecd57612ecd613154565b6101808101610a318284612c5c565b600061022080830163ffffffff8751168452602060018060601b038189015116818601526040880151604086015260608801516060860152612f4f6080860188612c5c565b6102008501929092528451908190526102408401918086019160005b81811015612f905783516001600160a01b031685529382019392820192600101612f6b565b509298975050505050505050565b6001600160601b0383168152604060208201819052600090610ec890830184612c30565b60405161018081016001600160401b0381118282101715612ff357634e487b7160e01b600052604160045260246000fd5b60405290565b6000821982111561300c5761300c61313e565b500190565b60006001600160601b038281168482168083038211156130335761303361313e565b01949350505050565b60008261305957634e487b7160e01b600052601260045260246000fd5b500490565b60008160001904831182151516156130785761307861313e565b500290565b60008282101561308f5761308f61313e565b500390565b60006001600160601b03838116908316818110156130b4576130b461313e565b039392505050565b60005b838110156130d75781810151838201526020016130bf565b838111156112bd5750506000910152565b600181811c908216806130fc57607f821691505b6020821081141561311d57634e487b7160e01b600052602260045260246000fd5b50919050565b60006000198214156131375761313761313e565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfe3078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a26469706673582212208c3adaa7d9d4bc9c5e08f49dd52f046b48ae150608bc3a5add1d2e6ddbb9166664736f6c63430008060033",
}

// KeeperRegistry13ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistry13MetaData.ABI instead.
var KeeperRegistry13ABI = KeeperRegistry13MetaData.ABI

// KeeperRegistry13Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistry13MetaData.Bin instead.
var KeeperRegistry13Bin = KeeperRegistry13MetaData.Bin

// DeployKeeperRegistry13 deploys a new Ethereum contract, binding an instance of KeeperRegistry13 to it.
func DeployKeeperRegistry13(auth *bind.TransactOpts, backend bind.ContractBackend, keeperRegistryLogic common.Address, config Config1_3) (common.Address, *types.Transaction, *KeeperRegistry13, error) {
	parsed, err := KeeperRegistry13MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry13Bin), backend, keeperRegistryLogic, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry13{KeeperRegistry13Caller: KeeperRegistry13Caller{contract: contract}, KeeperRegistry13Transactor: KeeperRegistry13Transactor{contract: contract}, KeeperRegistry13Filterer: KeeperRegistry13Filterer{contract: contract}}, nil
}

// KeeperRegistry13 is an auto generated Go binding around an Ethereum contract.
type KeeperRegistry13 struct {
	KeeperRegistry13Caller     // Read-only binding to the contract
	KeeperRegistry13Transactor // Write-only binding to the contract
	KeeperRegistry13Filterer   // Log filterer for contract events
}

// KeeperRegistry13Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistry13Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry13Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistry13Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry13Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistry13Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry13Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistry13Session struct {
	Contract     *KeeperRegistry13 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistry13CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistry13CallerSession struct {
	Contract *KeeperRegistry13Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// KeeperRegistry13TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistry13TransactorSession struct {
	Contract     *KeeperRegistry13Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// KeeperRegistry13Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistry13Raw struct {
	Contract *KeeperRegistry13 // Generic contract binding to access the raw methods on
}

// KeeperRegistry13CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistry13CallerRaw struct {
	Contract *KeeperRegistry13Caller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistry13TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistry13TransactorRaw struct {
	Contract *KeeperRegistry13Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistry13 creates a new instance of KeeperRegistry13, bound to a specific deployed contract.
func NewKeeperRegistry13(address common.Address, backend bind.ContractBackend) (*KeeperRegistry13, error) {
	contract, err := bindKeeperRegistry13(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13{KeeperRegistry13Caller: KeeperRegistry13Caller{contract: contract}, KeeperRegistry13Transactor: KeeperRegistry13Transactor{contract: contract}, KeeperRegistry13Filterer: KeeperRegistry13Filterer{contract: contract}}, nil
}

// NewKeeperRegistry13Caller creates a new read-only instance of KeeperRegistry13, bound to a specific deployed contract.
func NewKeeperRegistry13Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistry13Caller, error) {
	contract, err := bindKeeperRegistry13(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Caller{contract: contract}, nil
}

// NewKeeperRegistry13Transactor creates a new write-only instance of KeeperRegistry13, bound to a specific deployed contract.
func NewKeeperRegistry13Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistry13Transactor, error) {
	contract, err := bindKeeperRegistry13(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Transactor{contract: contract}, nil
}

// NewKeeperRegistry13Filterer creates a new log filterer instance of KeeperRegistry13, bound to a specific deployed contract.
func NewKeeperRegistry13Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistry13Filterer, error) {
	contract, err := bindKeeperRegistry13(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Filterer{contract: contract}, nil
}

// bindKeeperRegistry13 binds a generic wrapper to an already deployed contract.
func bindKeeperRegistry13(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistry13ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry13 *KeeperRegistry13Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry13.Contract.KeeperRegistry13Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry13 *KeeperRegistry13Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.KeeperRegistry13Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry13 *KeeperRegistry13Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.KeeperRegistry13Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry13 *KeeperRegistry13CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry13.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry13 *KeeperRegistry13TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry13 *KeeperRegistry13TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.contract.Transact(opts, method, params...)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "ARB_NITRO_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.ARBNITROORACLE(&_KeeperRegistry13.CallOpts)
}

// ARBNITROORACLE is a free data retrieval call binding the contract method 0x7f37618e.
//
// Solidity: function ARB_NITRO_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.ARBNITROORACLE(&_KeeperRegistry13.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.FASTGASFEED(&_KeeperRegistry13.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.FASTGASFEED(&_KeeperRegistry13.CallOpts)
}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) KEEPERREGISTRYLOGIC(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "KEEPER_REGISTRY_LOGIC")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry13.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry13.CallOpts)
}

// KEEPERREGISTRYLOGIC is a free data retrieval call binding the contract method 0x0852c7c9.
//
// Solidity: function KEEPER_REGISTRY_LOGIC() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry13.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry13.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) LINK() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINK(&_KeeperRegistry13.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINK(&_KeeperRegistry13.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINKETHFEED(&_KeeperRegistry13.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINKETHFEED(&_KeeperRegistry13.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "OPTIMISM_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.OPTIMISMORACLE(&_KeeperRegistry13.CallOpts)
}

// OPTIMISMORACLE is a free data retrieval call binding the contract method 0x850cce34.
//
// Solidity: function OPTIMISM_ORACLE() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.OPTIMISMORACLE(&_KeeperRegistry13.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Caller) PAYMENTMODEL(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "PAYMENT_MODEL")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Session) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry13.Contract.PAYMENTMODEL(&_KeeperRegistry13.CallOpts)
}

// PAYMENTMODEL is a free data retrieval call binding the contract method 0x8811cbe8.
//
// Solidity: function PAYMENT_MODEL() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry13.Contract.PAYMENTMODEL(&_KeeperRegistry13.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry13 *KeeperRegistry13Caller) REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "REGISTRY_GAS_OVERHEAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry13 *KeeperRegistry13Session) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry13.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry13.CallOpts)
}

// REGISTRYGASOVERHEAD is a free data retrieval call binding the contract method 0x5077b210.
//
// Solidity: function REGISTRY_GAS_OVERHEAD() view returns(uint256)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry13.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry13.CallOpts)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry13 *KeeperRegistry13Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry13.Contract.GetActiveUpkeepIDs(&_KeeperRegistry13.CallOpts, startIndex, maxCount)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry13.Contract.GetActiveUpkeepIDs(&_KeeperRegistry13.CallOpts, startIndex, maxCount)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getKeeperInfo", query)

	outstruct := new(struct {
		Payee   common.Address
		Active  bool
		Balance *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Payee = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Active = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry13.Contract.GetKeeperInfo(&_KeeperRegistry13.CallOpts, query)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry13.Contract.GetKeeperInfo(&_KeeperRegistry13.CallOpts, query)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMaxPaymentForGas(&_KeeperRegistry13.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMaxPaymentForGas(&_KeeperRegistry13.CallOpts, gasLimit)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry13.CallOpts, id)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry13.CallOpts, id)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry13.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry13.CallOpts, peer)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry13.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry13.CallOpts, peer)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetState(opts *bind.CallOpts) (struct {
	State   State1_3
	Config  Config1_3
	Keepers []common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		State   State1_3
		Config  Config1_3
		Keepers []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State1_3)).(*State1_3)
	outstruct.Config = *abi.ConvertType(out[1], new(Config1_3)).(*Config1_3)
	outstruct.Keepers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetState() (struct {
	State   State1_3
	Config  Config1_3
	Keepers []common.Address
}, error) {
	return _KeeperRegistry13.Contract.GetState(&_KeeperRegistry13.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetState() (struct {
	State   State1_3
	Config  Config1_3
	Keepers []common.Address
}, error) {
	return _KeeperRegistry13.Contract.GetState(&_KeeperRegistry13.CallOpts)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry13 *KeeperRegistry13Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getUpkeep", id)

	outstruct := new(struct {
		Target              common.Address
		ExecuteGas          uint32
		CheckData           []byte
		Balance             *big.Int
		LastKeeper          common.Address
		Admin               common.Address
		MaxValidBlocknumber uint64
		AmountSpent         *big.Int
		Paused              bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ExecuteGas = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CheckData = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Balance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LastKeeper = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Admin = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.MaxValidBlocknumber = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.AmountSpent = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Paused = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry13 *KeeperRegistry13Session) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	return _KeeperRegistry13.Contract.GetUpkeep(&_KeeperRegistry13.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent, bool paused)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
	Paused              bool
}, error) {
	return _KeeperRegistry13.Contract.GetUpkeep(&_KeeperRegistry13.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13Session) Owner() (common.Address, error) {
	return _KeeperRegistry13.Contract.Owner(&_KeeperRegistry13.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry13.Contract.Owner(&_KeeperRegistry13.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry13 *KeeperRegistry13Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry13 *KeeperRegistry13Session) Paused() (bool, error) {
	return _KeeperRegistry13.Contract.Paused(&_KeeperRegistry13.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) Paused() (bool, error) {
	return _KeeperRegistry13.Contract.Paused(&_KeeperRegistry13.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry13 *KeeperRegistry13Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry13 *KeeperRegistry13Session) TypeAndVersion() (string, error) {
	return _KeeperRegistry13.Contract.TypeAndVersion(&_KeeperRegistry13.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry13.Contract.TypeAndVersion(&_KeeperRegistry13.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry13.Contract.UpkeepTranscoderVersion(&_KeeperRegistry13.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry13 *KeeperRegistry13CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry13.Contract.UpkeepTranscoderVersion(&_KeeperRegistry13.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptOwnership(&_KeeperRegistry13.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptOwnership(&_KeeperRegistry13.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptPayeeship", keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptPayeeship(&_KeeperRegistry13.TransactOpts, keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptPayeeship(&_KeeperRegistry13.TransactOpts, keeper)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id)
}

// AcceptUpkeepAdmin is a paid mutator transaction binding the contract method 0xb148ab6b.
//
// Solidity: function acceptUpkeepAdmin(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "addFunds", id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AddFunds(&_KeeperRegistry13.TransactOpts, id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AddFunds(&_KeeperRegistry13.TransactOpts, id, amount)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "cancelUpkeep", id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CancelUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CancelUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry13 *KeeperRegistry13Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "checkUpkeep", id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry13 *KeeperRegistry13Session) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CheckUpkeep(&_KeeperRegistry13.TransactOpts, id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CheckUpkeep(&_KeeperRegistry13.TransactOpts, id, from)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.MigrateUpkeeps(&_KeeperRegistry13.TransactOpts, ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.MigrateUpkeeps(&_KeeperRegistry13.TransactOpts, ids, destination)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.OnTokenTransfer(&_KeeperRegistry13.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.OnTokenTransfer(&_KeeperRegistry13.TransactOpts, sender, amount, data)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Pause(&_KeeperRegistry13.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Pause(&_KeeperRegistry13.TransactOpts)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "pauseUpkeep", id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// PauseUpkeep is a paid mutator transaction binding the contract method 0x8765ecbe.
//
// Solidity: function pauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry13 *KeeperRegistry13Transactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "performUpkeep", id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry13 *KeeperRegistry13Session) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PerformUpkeep(&_KeeperRegistry13.TransactOpts, id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PerformUpkeep(&_KeeperRegistry13.TransactOpts, id, performData)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.ReceiveUpkeeps(&_KeeperRegistry13.TransactOpts, encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.ReceiveUpkeeps(&_KeeperRegistry13.TransactOpts, encodedUpkeeps)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "recoverFunds")
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RecoverFunds(&_KeeperRegistry13.TransactOpts)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RecoverFunds(&_KeeperRegistry13.TransactOpts)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry13 *KeeperRegistry13Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry13 *KeeperRegistry13Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RegisterUpkeep(&_KeeperRegistry13.TransactOpts, target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RegisterUpkeep(&_KeeperRegistry13.TransactOpts, target, gasLimit, admin, checkData)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetConfig(opts *bind.TransactOpts, config Config1_3) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setConfig", config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) SetConfig(config Config1_3) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetConfig(&_KeeperRegistry13.TransactOpts, config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetConfig(config Config1_3) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetConfig(&_KeeperRegistry13.TransactOpts, config)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setKeepers", keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetKeepers(&_KeeperRegistry13.TransactOpts, keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetKeepers(&_KeeperRegistry13.TransactOpts, keepers, payees)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry13.TransactOpts, peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry13.TransactOpts, peer, permission)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetUpkeepGasLimit(&_KeeperRegistry13.TransactOpts, id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetUpkeepGasLimit(&_KeeperRegistry13.TransactOpts, id, gasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferOwnership(&_KeeperRegistry13.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferOwnership(&_KeeperRegistry13.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferPayeeship(&_KeeperRegistry13.TransactOpts, keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferPayeeship(&_KeeperRegistry13.TransactOpts, keeper, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id, proposed)
}

// TransferUpkeepAdmin is a paid mutator transaction binding the contract method 0x1a2af011.
//
// Solidity: function transferUpkeepAdmin(uint256 id, address proposed) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id, proposed)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Unpause(&_KeeperRegistry13.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Unpause(&_KeeperRegistry13.TransactOpts)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "unpauseUpkeep", id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UnpauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// UnpauseUpkeep is a paid mutator transaction binding the contract method 0x5165f2f5.
//
// Solidity: function unpauseUpkeep(uint256 id) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UnpauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UpdateCheckData(&_KeeperRegistry13.TransactOpts, id, newCheckData)
}

// UpdateCheckData is a paid mutator transaction binding the contract method 0x9fab4386.
//
// Solidity: function updateCheckData(uint256 id, bytes newCheckData) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UpdateCheckData(&_KeeperRegistry13.TransactOpts, id, newCheckData)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawFunds", id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawFunds(&_KeeperRegistry13.TransactOpts, id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawFunds(&_KeeperRegistry13.TransactOpts, id, to)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawOwnerFunds")
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawOwnerFunds(&_KeeperRegistry13.TransactOpts)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawOwnerFunds(&_KeeperRegistry13.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawPayment", from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawPayment(&_KeeperRegistry13.TransactOpts, from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawPayment(&_KeeperRegistry13.TransactOpts, from, to)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Fallback(&_KeeperRegistry13.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Fallback(&_KeeperRegistry13.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13Session) Receive() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Receive(&_KeeperRegistry13.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Receive() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Receive(&_KeeperRegistry13.TransactOpts)
}

// KeeperRegistry13ConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistry13 contract.
type KeeperRegistry13ConfigSetIterator struct {
	Event *KeeperRegistry13ConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13ConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13ConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13ConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13ConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13ConfigSet represents a ConfigSet event raised by the KeeperRegistry13 contract.
type KeeperRegistry13ConfigSet struct {
	Config Config1_3
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry13ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13ConfigSetIterator{contract: _KeeperRegistry13.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13ConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13ConfigSet)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigSet is a log parse operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseConfigSet(log types.Log) (*KeeperRegistry13ConfigSet, error) {
	event := new(KeeperRegistry13ConfigSet)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13FundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistry13 contract.
type KeeperRegistry13FundsAddedIterator struct {
	Event *KeeperRegistry13FundsAdded // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13FundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13FundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13FundsAdded represents a FundsAdded event raised by the KeeperRegistry13 contract.
type KeeperRegistry13FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry13FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13FundsAddedIterator{contract: _KeeperRegistry13.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13FundsAdded)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistry13FundsAdded, error) {
	event := new(KeeperRegistry13FundsAdded)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13FundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistry13 contract.
type KeeperRegistry13FundsWithdrawnIterator struct {
	Event *KeeperRegistry13FundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13FundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13FundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13FundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistry13 contract.
type KeeperRegistry13FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13FundsWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13FundsWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistry13FundsWithdrawn, error) {
	event := new(KeeperRegistry13FundsWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13KeepersUpdatedIterator is returned from FilterKeepersUpdated and is used to iterate over the raw logs and unpacked data for KeepersUpdated events raised by the KeeperRegistry13 contract.
type KeeperRegistry13KeepersUpdatedIterator struct {
	Event *KeeperRegistry13KeepersUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13KeepersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13KeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13KeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13KeepersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13KeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13KeepersUpdated represents a KeepersUpdated event raised by the KeeperRegistry13 contract.
type KeeperRegistry13KeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKeepersUpdated is a free log retrieval operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistry13KeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13KeepersUpdatedIterator{contract: _KeeperRegistry13.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepersUpdated is a free log subscription operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13KeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13KeepersUpdated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseKeepersUpdated is a log parse operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistry13KeepersUpdated, error) {
	event := new(KeeperRegistry13KeepersUpdated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13OwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistry13OwnerFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13OwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry13OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnerFundsWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13OwnerFundsWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry13OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistry13OwnerFundsWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistry13OwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnershipTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13OwnershipTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry13OwnershipTransferRequested, error) {
	event := new(KeeperRegistry13OwnershipTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnershipTransferredIterator struct {
	Event *KeeperRegistry13OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13OwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistry13 contract.
type KeeperRegistry13OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnershipTransferredIterator{contract: _KeeperRegistry13.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13OwnershipTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistry13OwnershipTransferred, error) {
	event := new(KeeperRegistry13OwnershipTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistry13 contract.
type KeeperRegistry13PausedIterator struct {
	Event *KeeperRegistry13Paused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13Paused represents a Paused event raised by the KeeperRegistry13 contract.
type KeeperRegistry13Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry13PausedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PausedIterator{contract: _KeeperRegistry13.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13Paused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePaused(log types.Log) (*KeeperRegistry13Paused, error) {
	event := new(KeeperRegistry13Paused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13PayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistry13 contract.
type KeeperRegistry13PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistry13PayeeshipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13PayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistry13 contract.
type KeeperRegistry13PayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferRequestedIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PayeeshipTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13PayeeshipTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry13PayeeshipTransferRequested, error) {
	event := new(KeeperRegistry13PayeeshipTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13PayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistry13 contract.
type KeeperRegistry13PayeeshipTransferredIterator struct {
	Event *KeeperRegistry13PayeeshipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13PayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13PayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13PayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistry13 contract.
type KeeperRegistry13PayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferredIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PayeeshipTransferredIterator{contract: _KeeperRegistry13.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13PayeeshipTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry13PayeeshipTransferred, error) {
	event := new(KeeperRegistry13PayeeshipTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13PaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistry13 contract.
type KeeperRegistry13PaymentWithdrawnIterator struct {
	Event *KeeperRegistry13PaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13PaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13PaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13PaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistry13 contract.
type KeeperRegistry13PaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry13PaymentWithdrawnIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PaymentWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13PaymentWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry13PaymentWithdrawn, error) {
	event := new(KeeperRegistry13PaymentWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UnpausedIterator struct {
	Event *KeeperRegistry13Unpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13Unpaused represents a Unpaused event raised by the KeeperRegistry13 contract.
type KeeperRegistry13Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry13UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UnpausedIterator{contract: _KeeperRegistry13.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13Unpaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUnpaused(log types.Log) (*KeeperRegistry13Unpaused, error) {
	event := new(KeeperRegistry13Unpaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepAdminTransferRequestedIterator is returned from FilterUpkeepAdminTransferRequested and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferRequested events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistry13UpkeepAdminTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepAdminTransferRequested represents a UpkeepAdminTransferRequested event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferRequested is a free log retrieval operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferRequested is a free log subscription operation binding the contract event 0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35.
//
// Solidity: event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepAdminTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistry13UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistry13UpkeepAdminTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepAdminTransferredIterator is returned from FilterUpkeepAdminTransferred and is used to iterate over the raw logs and unpacked data for UpkeepAdminTransferred events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistry13UpkeepAdminTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepAdminTransferred represents a UpkeepAdminTransferred event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpkeepAdminTransferred is a free log retrieval operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepAdminTransferredIterator{contract: _KeeperRegistry13.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

// WatchUpkeepAdminTransferred is a free log subscription operation binding the contract event 0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c.
//
// Solidity: event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepAdminTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistry13UpkeepAdminTransferred, error) {
	event := new(KeeperRegistry13UpkeepAdminTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepCanceledIterator struct {
	Event *KeeperRegistry13UpkeepCanceled // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry13UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepCanceledIterator{contract: _KeeperRegistry13.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepCanceled)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistry13UpkeepCanceled, error) {
	event := new(KeeperRegistry13UpkeepCanceled)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepCheckDataUpdatedIterator is returned from FilterUpkeepCheckDataUpdated and is used to iterate over the raw logs and unpacked data for UpkeepCheckDataUpdated events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistry13UpkeepCheckDataUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepCheckDataUpdated represents a UpkeepCheckDataUpdated event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCheckDataUpdated is a free log retrieval operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepCheckDataUpdatedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

// WatchUpkeepCheckDataUpdated is a free log subscription operation binding the contract event 0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf.
//
// Solidity: event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepCheckDataUpdated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistry13UpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistry13UpkeepCheckDataUpdated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistry13UpkeepGasLimitSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepGasLimitSetIterator{contract: _KeeperRegistry13.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepGasLimitSet)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry13UpkeepGasLimitSet, error) {
	event := new(KeeperRegistry13UpkeepGasLimitSet)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepMigratedIterator struct {
	Event *KeeperRegistry13UpkeepMigrated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepMigratedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepMigrated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistry13UpkeepMigrated, error) {
	event := new(KeeperRegistry13UpkeepMigrated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepPausedIterator is returned from FilterUpkeepPaused and is used to iterate over the raw logs and unpacked data for UpkeepPaused events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepPausedIterator struct {
	Event *KeeperRegistry13UpkeepPaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepPaused represents a UpkeepPaused event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPaused is a free log retrieval operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepPausedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepPaused is a free log subscription operation binding the contract event 0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f.
//
// Solidity: event UpkeepPaused(uint256 indexed id)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepPaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistry13UpkeepPaused, error) {
	event := new(KeeperRegistry13UpkeepPaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepPerformedIterator struct {
	Event *KeeperRegistry13UpkeepPerformed // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepPerformed struct {
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpkeepPerformed is a free log retrieval operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistry13UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepPerformedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepPerformed)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpkeepPerformed is a log parse operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistry13UpkeepPerformed, error) {
	event := new(KeeperRegistry13UpkeepPerformed)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepReceivedIterator struct {
	Event *KeeperRegistry13UpkeepReceived // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepReceivedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepReceived)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistry13UpkeepReceived, error) {
	event := new(KeeperRegistry13UpkeepReceived)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepRegisteredIterator struct {
	Event *KeeperRegistry13UpkeepRegistered // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepRegisteredIterator{contract: _KeeperRegistry13.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepRegistered)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistry13UpkeepRegistered, error) {
	event := new(KeeperRegistry13UpkeepRegistered)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry13UpkeepUnpausedIterator is returned from FilterUpkeepUnpaused and is used to iterate over the raw logs and unpacked data for UpkeepUnpaused events raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepUnpausedIterator struct {
	Event *KeeperRegistry13UpkeepUnpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry13UpkeepUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry13UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry13UpkeepUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry13UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry13UpkeepUnpaused represents a UpkeepUnpaused event raised by the KeeperRegistry13 contract.
type KeeperRegistry13UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUpkeepUnpaused is a free log retrieval operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepUnpausedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

// WatchUpkeepUnpaused is a free log subscription operation binding the contract event 0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456.
//
// Solidity: event UpkeepUnpaused(uint256 indexed id)
func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry13UpkeepUnpaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistry13UpkeepUnpaused, error) {
	event := new(KeeperRegistry13UpkeepUnpaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
