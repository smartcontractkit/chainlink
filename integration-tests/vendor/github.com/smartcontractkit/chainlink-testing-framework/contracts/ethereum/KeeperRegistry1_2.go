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

// Config1_2 is an auto generated low-level Go binding around an user-defined struct.
type Config1_2 struct {
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

// State1_2 is an auto generated low-level Go binding around an user-defined struct.
type State1_2 struct {
	Nonce               uint32
	OwnerLinkBalance    *big.Int
	ExpectedLinkBalance *big.Int
	NumUpkeeps          *big.Int
}

// KeeperRegistry12MetaData contains all meta data concerning the KeeperRegistry12 contract.
var KeeperRegistry12MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_2\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig1_2\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getKeeperInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistry1_2.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"expectedLinkBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"}],\"internalType\":\"structState1_2\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_2\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig1_2\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistry1_2.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200558d3803806200558d833981016040819052620000349162000577565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000107565b50506001600255506003805460ff191690556001600160601b0319606085811b821660805284811b821660a05283901b1660c052620000fd81620001b3565b50505050620007fa565b6001600160a01b038116331415620001625760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001bd620004a8565b600d5460e082015163ffffffff91821691161015620001ef57604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601360006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325816040516200049d9190620006c3565b60405180910390a150565b6000546001600160a01b03163314620005045760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200051e57600080fd5b919050565b805161ffff811681146200051e57600080fd5b805162ffffff811681146200051e57600080fd5b805163ffffffff811681146200051e57600080fd5b80516001600160601b03811681146200051e57600080fd5b6000806000808486036101e08112156200059057600080fd5b6200059b8662000506565b9450620005ab6020870162000506565b9350620005bb6040870162000506565b925061018080605f1983011215620005d257600080fd5b620005dc620007c2565b9150620005ec606088016200054a565b8252620005fc608088016200054a565b60208301526200060f60a0880162000536565b60408301526200062260c088016200054a565b60608301526200063560e0880162000536565b60808301526101006200064a81890162000523565b60a08401526101206200065f818a016200055f565b60c085015261014062000674818b016200054a565b60e0860152610160808b015184870152848b0151838701526200069b6101a08c0162000506565b82870152620006ae6101c08c0162000506565b90860152509699959850939650909450505050565b815163ffffffff16815261018081016020830151620006ea602084018263ffffffff169052565b50604083015162000702604084018262ffffff169052565b5060608301516200071b606084018263ffffffff169052565b50608083015162000733608084018262ffffff169052565b5060a08301516200074a60a084018261ffff169052565b5060c08301516200076660c08401826001600160601b03169052565b5060e08301516200077f60e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60405161018081016001600160401b0381118282101715620007f457634e487b7160e01b600052604160045260246000fd5b60405290565b60805160601c60a05160601c60c05160601c614d14620008796000396000818161034e01526132e001526000818161047901526133b401526000818161027b01528181610b3901528181610d62015281816113bf015281816116dc015281816117b301528181611a6601528181611d0d0152611d960152614d146000f3fe608060405234801561001057600080fd5b50600436106101dc5760003560e01c806393f0c1fc11610105578063b7fdb4361161009d578063b7fdb436146104c9578063c41b813a146104dc578063c7c3a19a14610500578063c804802214610527578063da5c67411461053a578063eb5dcd6c1461055b578063ef47a0ce1461056e578063f2fde38b14610581578063faa3e9961461059457600080fd5b806393f0c1fc14610408578063948108f714610428578063a4c0ed361461043b578063a710b2211461044e578063a72aa27e14610461578063ad17836114610474578063b121e1471461049b578063b657bc9c146104ae578063b79550be146104c157600080fd5b80635c975abb116101785780635c975abb14610385578063744bfe611461039c57806379ba5097146103af5780637bbaf1ea146103b75780637d9b97e0146103ca5780638456cb59146103d257806385c1b0ba146103da5780638da5cb5b146103ed5780638e86139b146103f557600080fd5b806306e3b632146101e1578063181f5a771461020a5780631865c57d1461024a578063187256e8146102615780631b6b6d23146102765780631e12b8a5146102aa5780633f4ba83a146103415780634584a4191461034957806348013d7b14610370575b600080fd5b6101f46101ef3660046142f9565b6105cd565b60405161020191906147ba565b60405180910390f35b61023d6040518060400160405280601481526020017304b6565706572526567697374727920312e322e360641b81525081565b60405161020191906147fe565b6102526106af565b60405161020193929190614954565b61027461026f366004613ddf565b610909565b005b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b60405161020191906145ab565b6103136102b8366004613d91565b6001600160a01b0390811660009081526008602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b031692810183905260019091015460ff16151592018290529192909190565b604080516001600160a01b03909416845291151560208401526001600160601b031690820152606001610201565b61027461094f565b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b610378600081565b604051610201919061490a565b60035460ff165b6040519015158152602001610201565b6102746103aa36600461428b565b610961565b610274610bce565b61038c6103c53660046142ae565b610c7d565b610274610cdb565b610274610def565b6102746103e8366004613f47565b610dff565b61029d6113f8565b6102746104033660046140e5565b611407565b61041b610416366004614259565b6115dd565b60405161020191906149e8565b61027461043636600461433e565b611611565b610274610449366004613e1a565b6117a8565b61027461045c366004613dac565b611911565b61027461046f36600461431b565b611af4565b61029d7f000000000000000000000000000000000000000000000000000000000000000081565b6102746104a9366004613d91565b611c1d565b61041b6104bc366004614259565b611cca565b610274611ceb565b6102746104d7366004613ee8565b611def565b6104ef6104ea36600461428b565b61203c565b604051610201959493929190614811565b61051361050e366004614259565b612252565b6040516102019897969594939291906145d8565b610274610535366004614259565b6123c0565b61054d610548366004613e73565b612524565b604051908152602001610201565b610274610569366004613dac565b6126b8565b61027461057c36600461417b565b612799565b61027461058f366004613d91565b612a89565b6105c06105a2366004613d91565b6001600160a01b03166000908152600b602052604090205460ff1690565b60405161020191906148f0565b606060006105db6005612a9d565b90508084106105fd57604051631390f2a160e01b815260040160405180910390fd5b8261060f5761060c8482614b8b565b92505b6000836001600160401b0381111561062957610629614cc8565b604051908082528060200260200182016040528015610652578160200160208202803683370190505b50905060005b848110156106a45761067561066d8288614b07565b600590612aa7565b82828151811061068757610687614cb2565b60209081029190910101528061069c81614c31565b915050610658565b509150505b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600c5463ffffffff8082168352600160201b808304821660208086019190915262ffffff600160401b8504811686880152600160581b85048416606087810191909152600160781b8604909116608087015261ffff600160901b86041660a08701526001600160601b03600160a01b909504851660c0870152600d5480851660e0880152929092049092166101008501819052875260105490921690860152601154928501929092526107f26005612a9d565b606080860191909152815163ffffffff908116855260208084015182168187015260408085015162ffffff90811682890152858501518416948801949094526080808601519094169387019390935260a08085015161ffff169087015260c0808501516001600160601b03169087015260e08085015190921691860191909152600e54610100860152600f546101208601526012546001600160a01b03908116610140870152601354166101608601526004805483518184028101840190945280845287938793909183918301828280156108f657602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116108d8575b5050505050905093509350935050909192565b610911612aba565b6001600160a01b0382166000908152600b60205260409020805482919060ff1916600183600381111561094657610946614c86565b02179055505050565b610957612aba565b61095f612b0d565b565b806001600160a01b03811661098957604051634e46966960e11b815260040160405180910390fd5b6000838152600760205260409020600201548390600160601b90046001600160a01b031633146109cc5760405163523e0b8360e11b815260040160405180910390fd5b60008481526007602052604090206001015443600160201b9091046001600160401b03161115610a12576040516001627b1a2360e01b0319815260040160405180910390fd5b600c54600085815260076020526040812080546002909101546001600160601b03600160a01b9094048416939182169291169083821015610a7657610a578285614ba2565b9050826001600160601b0316816001600160601b03161115610a765750815b6000610a828285614ba2565b60008a815260076020526040902080546001600160601b0319169055601054909150610ab89083906001600160601b0316614b1f565b601080546001600160601b0319166001600160601b03928316179055601154610ae391831690614b8b565b60115560405189907ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a6831890610b1a9084908c906149fc565b60405180910390a260405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90610b70908b908590600401614651565b602060405180830381600087803b158015610b8a57600080fd5b505af1158015610b9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bc2919061407d565b50505050505050505050565b6001546001600160a01b03163314610c265760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610c87612b59565b610cd3610cce338686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250612b9f915050565b612c7a565b949350505050565b610ce3612aba565b6010546011546001600160601b0390911690610d00908290614b8b565b601155601080546001600160601b03191690556040517f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f190610d439083906149e8565b60405180910390a160405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90610d999033908590600401614651565b602060405180830381600087803b158015610db357600080fd5b505af1158015610dc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610deb919061407d565b5050565b610df7612aba565b61095f612feb565b60016001600160a01b0382166000908152600b602052604090205460ff166003811115610e2e57610e2e614c86565b14158015610e69575060036001600160a01b0382166000908152600b602052604090205460ff166003811115610e6657610e66614c86565b14155b15610e87576040516303afbb0f60e21b815260040160405180910390fd5b6012546001600160a01b0316610eb05760405163d12d7d8d60e01b815260040160405180910390fd5b81610ece57604051632c2fc94160e01b815260040160405180910390fd5b6000610ed86138f7565b600080856001600160401b03811115610ef357610ef3614cc8565b604051908082528060200260200182016040528015610f2657816020015b6060815260200190600190039081610f115790505b5090506000866001600160401b03811115610f4357610f43614cc8565b604051908082528060200260200182016040528015610f7c57816020015b610f696138f7565b815260200190600190039081610f615790505b50905060005b8781101561120457888882818110610f9c57610f9c614cb2565b60209081029290920135600081815260078452604090819020815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811698840198909852600184015463ffffffff8116958401959095526001600160401b03600160201b8604166060840152938190048716608083015260029092015492831660a0820152910490931660c084018190529098509196505033146110595760405163523e0b8360e11b815260040160405180910390fd5b60608501516001600160401b039081161461108757604051633425886760e21b815260040160405180910390fd5b8482828151811061109a5761109a614cb2565b6020026020010181905250600a600087815260200190815260200160002080546110c390614bf6565b80601f01602080910402602001604051908101604052809291908181526020018280546110ef90614bf6565b801561113c5780601f106111115761010080835404028352916020019161113c565b820191906000526020600020905b81548152906001019060200180831161111f57829003601f168201915b505050505083828151811061115357611153614cb2565b60209081029190910101528451611173906001600160601b031685614b07565b600087815260076020908152604080832083815560018101849055600201839055600a90915281209195506111a89190613933565b6111b3600587613028565b50857fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff8660000151896040516111ea9291906149fc565b60405180910390a2806111fc81614c31565b915050610f82565b50826011546112139190614b8b565b601155604051600090611230908a908a90859087906020016146a5565b6040516020818303038152906040529050866001600160a01b0316638e86139b601260009054906101000a90046001600160a01b03166001600160a01b031663c71249ab60008b6001600160a01b03166348013d7b6040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156112b157600080fd5b505af11580156112c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112e9919061415a565b866040518463ffffffff1660e01b815260040161130893929190614918565b60006040518083038186803b15801561132057600080fd5b505afa158015611334573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261135c9190810190614126565b6040518263ffffffff1660e01b815260040161137891906147fe565b600060405180830381600087803b15801561139257600080fd5b505af11580156113a6573d6000803e3d6000fd5b505060405163a9059cbb60e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150610b70908a9088906004016145bf565b6000546001600160a01b031690565b6002336000908152600b602052604090205460ff16600381111561142d5761142d614c86565b1415801561145f57506003336000908152600b602052604090205460ff16600381111561145c5761145c614c86565b14155b1561147d576040516303afbb0f60e21b815260040160405180910390fd5b6000808061148d84860186613f9a565b92509250925060005b83518110156115d5576115538482815181106114b4576114b4614cb2565b60200260200101518483815181106114ce576114ce614cb2565b6020026020010151608001518584815181106114ec576114ec614cb2565b60200260200101516040015186858151811061150a5761150a614cb2565b602002602001015160c0015187868151811061152857611528614cb2565b60200260200101516000015187878151811061154657611546614cb2565b6020026020010151613034565b83818151811061156557611565614cb2565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a718483815181106115a0576115a0614cb2565b602002602001015160000151336040516115bb9291906149fc565b60405180910390a2806115cd81614c31565b915050611496565b505050505050565b60008060006115ea6132ad565b9150915060006115fb83600061348e565b90506116088582846134c4565b95945050505050565b6000828152600760205260409020600101548290600160201b90046001600160401b039081161461165557604051633425886760e21b815260040160405180910390fd5b6000838152600760205260409020546116789083906001600160601b0316614b1f565b600084815260076020526040902080546001600160601b0319166001600160601b039283161790556011546116af91841690614b07565b6011556040516323b872dd60e01b81523360048201523060248201526001600160601b03831660448201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906323b872dd90606401602060405180830381600087803b15801561172857600080fd5b505af115801561173c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611760919061407d565b50336001600160a01b0316837fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062038460405161179b91906149e8565b60405180910390a3505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146117f15760405163c8bad78d60e01b815260040160405180910390fd5b6020811461181257604051630dfe930960e41b815260040160405180910390fd5b600061182082840184614259565b600081815260076020526040902060010154909150600160201b90046001600160401b039081161461186557604051633425886760e21b815260040160405180910390fd5b6000818152600760205260409020546118889085906001600160601b0316614b1f565b600082815260076020526040902080546001600160601b0319166001600160601b03929092169190911790556011546118c2908590614b07565b601181905550846001600160a01b0316817fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062038660405161190291906149e8565b60405180910390a35050505050565b806001600160a01b03811661193957604051634e46966960e11b815260040160405180910390fd5b6001600160a01b0383811660009081526008602090815260409182902082516060810184528154948516808252600160a01b9095046001600160601b0316928101929092526001015460ff161515918101919091529033146119ae5760405163cebf515b60e01b815260040160405180910390fd5b6001600160a01b03808516600090815260086020908152604090912080549092169091558101516011546119eb916001600160601b031690614b8b565b601181905550826001600160a01b031681602001516001600160601b0316856001600160a01b03167f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f4069833604051611a4291906145ab565b60405180910390a4602081015160405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163a9059cbb91611a9b918791600401614651565b602060405180830381600087803b158015611ab557600080fd5b505af1158015611ac9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611aed919061407d565b5050505050565b6000828152600760205260409020600101548290600160201b90046001600160401b0390811614611b3857604051633425886760e21b815260040160405180910390fd5b6000838152600760205260409020600201548390600160601b90046001600160a01b03163314611b7b5760405163523e0b8360e11b815260040160405180910390fd5b6108fc8363ffffffff161080611b9c5750600d5463ffffffff908116908416115b15611bba576040516314c237fb60e01b815260040160405180910390fd5b600084815260076020908152604091829020600101805463ffffffff191663ffffffff8716908117909155915191825285917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a250505050565b6001600160a01b03818116600090815260096020526040902054163314611c57576040516333a973d560e11b815260040160405180910390fd5b6001600160a01b0381811660008181526008602090815260408083208054336001600160a01b031980831682179093556009909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600760205260408120600101546106a99063ffffffff166115dd565b611cf3612aba565b6040516370a0823160e01b81526000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190611d429030906004016145ab565b60206040518083038186803b158015611d5a57600080fd5b505afa158015611d6e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d929190614272565b90507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb3360115484611dd29190614b8b565b6040518363ffffffff1660e01b8152600401610d999291906145bf565b611df7612aba565b8281141580611e065750600283105b15611e24576040516367aa603560e11b815260040160405180910390fd5b60005b600454811015611e8557600060048281548110611e4657611e46614cb2565b60009182526020808320909101546001600160a01b031682526008905260409020600101805460ff191690555080611e7d81614c31565b915050611e27565b5060005b83811015611feb576000858583818110611ea557611ea5614cb2565b9050602002016020810190611eba9190613d91565b6001600160a01b03808216600090815260086020526040812080549394509290911690868686818110611eef57611eef614cb2565b9050602002016020810190611f049190613d91565b90506001600160a01b0381161580611f5657506001600160a01b03821615801590611f415750806001600160a01b0316826001600160a01b031614155b8015611f5657506001600160a01b0381811614155b15611f7457604051631670f44760e31b815260040160405180910390fd5b600183015460ff1615611f9a57604051630d5f433160e21b815260040160405180910390fd5b6001838101805460ff191690911790556001600160a01b0381811614611fd45782546001600160a01b0319166001600160a01b0382161783555b505050508080611fe390614c31565b915050611e89565b50611ff86004858561396d565b507f056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f8484848460405161202e9493929190614673565b60405180910390a150505050565b606060008060008061204c613587565b6000878152600760209081526040808320815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811684880152600185015463ffffffff8116858801526001600160401b03600160201b82041660608601528390048116608085015260029094015490811660a08401520490911660c08201528a8452600a9092528083209051919291636e04ff0d60e01b916120f691602401614848565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b038381831617835250505050905060008083608001516001600160a01b0316600c600001600b9054906101000a900463ffffffff1663ffffffff1684604051612166919061458f565b60006040518083038160008787f1925050503d80600081146121a4576040519150601f19603f3d011682016040523d82523d6000602084013e6121a9565b606091505b5091509150816121ce57806040516396c3623560e01b8152600401610c1d91906147fe565b808060200190518101906121e29190614098565b99509150816122045760405163865676e360e01b815260040160405180910390fd5b60006122138b8d8c6000612b9f565b905061222885826000015183606001516135a6565b6060810151608082015160a083015160c0909301519b9e919d509b50909998509650505050505050565b6000818152600760209081526040808320815160e08101835281546001600160601b0380821683526001600160a01b03600160601b928390048116848801908152600186015463ffffffff81168689018190526001600160401b03600160201b83041660608881019182529287900485166080890181905260029099015495861660a089019081529690950490931660c087019081528b8b52600a9099529689208551915198519351945181548b9a8b998a998a998a998a999298939792969295939493909290869061232490614bf6565b80601f016020809104026020016040519081016040528092919081815260200182805461235090614bf6565b801561239d5780601f106123725761010080835404028352916020019161239d565b820191906000526020600020905b81548152906001019060200180831161238057829003601f168201915b505050505095509850985098509850985098509850985050919395975091939597565b6000818152600760205260408120600101546001600160401b03600160201b90910481169190821415906123f26113f8565b6001600160a01b0316336001600160a01b03161490508180156124275750808015612425575043836001600160401b0316115b155b1561244557604051631f7806af60e31b815260040160405180910390fd5b801580156124745750600084815260076020526040902060020154600160601b90046001600160a01b03163314155b1561249257604051637dedc72b60e11b815260040160405180910390fd5b43816124a6576124a3603282614b07565b90505b600085815260076020526040902060010180546bffffffffffffffff000000001916600160201b6001600160401b038416021790556124e6600586613028565b506040516001600160401b0382169086907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a35050505050565b600061252e6113f8565b6001600160a01b0316336001600160a01b03161415801561255a57506013546001600160a01b03163314155b156125785760405163d48b678b60e01b815260040160405180910390fd5b612583600143614b8b565b600d5460408051924060208401523060601b6001600160601b03191690830152600160201b900460e01b6001600160e01b03191660548201526058016040516020818303038152906040528051906020012060001c905061261f81878787600088888080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061303492505050565b600d8054600160201b900463ffffffff1690600461263c83614c4c565b91906101000a81548163ffffffff021916908363ffffffff16021790555050807fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d01286866040516126a792919063ffffffff9290921682526001600160a01b0316602082015260400190565b60405180910390a295945050505050565b6001600160a01b038281166000908152600860205260409020541633146126f25760405163cebf515b60e01b815260040160405180910390fd5b6001600160a01b03811633141561271c57604051638c8728c760e01b815260040160405180910390fd5b6001600160a01b03828116600090815260096020526040902054811690821614610deb576001600160a01b0382811660008181526009602052604080822080546001600160a01b0319169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6127a1612aba565b600d5460e082015163ffffffff918216911610156127d257604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601360006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051612a7e9190614945565b60405180910390a150565b612a91612aba565b612a9a81613648565b50565b60006106a9825490565b6000612ab383836136ec565b9392505050565b6000546001600160a01b0316331461095f5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610c1d565b612b15613716565b6003805460ff191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b604051612b4f91906145ab565b60405180910390a1565b60035460ff161561095f5760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b6044820152606401610c1d565b612be86040518060e0016040528060006001600160a01b031681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526007602052604081206001015463ffffffff169080612c0a6132ad565b915091506000612c1a838761348e565b90506000612c298583856134c4565b6040805160e0810182526001600160a01b038d168152602081018c90529081018a90526001600160601b03909116606082015260808101959095525060a084015260c0830152509050949350505050565b6000600280541415612cce5760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610c1d565b600280556020828101516000818152600790925260409091206001015443600160201b9091046001600160401b031611612d1b57604051633425886760e21b815260040160405180910390fd5b602080840151600090815260078252604090819020815160e08101835281546001600160601b0380821683526001600160a01b03600160601b92839004811696840196909652600184015463ffffffff8116958401959095526001600160401b03600160201b860416606080850191909152948290048616608084015260029093015492831660a083015290910490921660c0830152845190850151612dc29183916135a6565b60005a90506000634585e33b60e01b8660400151604051602401612de691906147fe565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b0383818316178352505050509050612e2e866080015184608001518361375f565b94505a612e3b9083614b8b565b91506000612e52838860a001518960c001516134c4565b602080890151600090815260079091526040902054909150612e7e9082906001600160601b0316614ba2565b6020888101805160009081526007909252604080832080546001600160601b0319166001600160601b0395861617905590518252902060020154612ec491839116614b1f565b60208881018051600090815260078352604080822060020180546001600160601b0319166001600160601b039687161790558b519251825280822080548616600160601b6001600160a01b03958616021790558b5190921681526008909252902054612f39918391600160a01b900416614b1f565b6008600089600001516001600160a01b03166001600160a01b0316815260200190815260200160002060000160146101000a8154816001600160601b0302191690836001600160601b0316021790555086600001516001600160a01b031686151588602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848b60400151604051612fd4929190614a1e565b60405180910390a450505050506001600255919050565b612ff3612b59565b6003805460ff191660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258612b423390565b6000612ab383836137ab565b61303c612b59565b6001600160a01b0385163b613064576040516309ee12d560e01b815260040160405180910390fd5b6108fc8463ffffffff1610806130855750600d5463ffffffff908116908516115b156130a3576040516314c237fb60e01b815260040160405180910390fd5b6040518060e00160405280836001600160601b0316815260200160006001600160a01b031681526020018563ffffffff1681526020016001600160401b0380168152602001866001600160a01b0316815260200160006001600160601b03168152602001846001600160a01b03168152506007600088815260200190815260200160002060008201518160000160006101000a8154816001600160601b0302191690836001600160601b03160217905550602082015181600001600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160006101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160046101000a8154816001600160401b0302191690836001600160401b03160217905550608082015181600101600c6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060a08201518160020160006101000a8154816001600160601b0302191690836001600160601b0316021790555060c082015181600201600c6101000a8154816001600160a01b0302191690836001600160a01b03160217905550905050816001600160601b03166011546132769190614b07565b6011556000868152600a602090815260409091208251613298928401906139d0565b506132a460058761389e565b50505050505050565b6000806000600c600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561333757600080fd5b505afa15801561334b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061336f9190614361565b509450909250849150508015613393575061338a8242614b8b565b8463ffffffff16105b8061339f575060008113155b156133ae57600e5495506133b2565b8095505b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561340b57600080fd5b505afa15801561341f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134439190614361565b509450909250849150508015613467575061345e8242614b8b565b8463ffffffff16105b80613473575060008113155b1561348257600f549450613486565b8094505b505050509091565b600c546000906134a990600160901b900461ffff1684614b6c565b90508180156134b75750803a105b156106a957503a92915050565b6000806134d46201388086614b07565b6134de9085614b6c565b600c549091506000906134fb9063ffffffff16633b9aca00614b07565b600c5490915060009061352090600160201b900463ffffffff1664e8d4a51000614b6c565b858361353086633b9aca00614b6c565b61353a9190614b6c565b6135449190614b4a565b61354e9190614b07565b90506b033b2e3c9fd0803ce800000081111561357d5760405163156baa3d60e11b815260040160405180910390fd5b9695505050505050565b321561095f5760405163b60ac5db60e01b815260040160405180910390fd5b6001600160a01b03821660009081526008602052604090206001015460ff166135e2576040516319f759fb60e31b815260040160405180910390fd5b82516001600160601b031681111561360d5760405163356680b760e01b815260040160405180910390fd5b816001600160a01b031683602001516001600160a01b0316141561364357604051621af04160e61b815260040160405180910390fd5b505050565b6001600160a01b03811633141561369b5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610c1d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061370357613703614cb2565b9060005260206000200154905092915050565b60035460ff1661095f5760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b6044820152606401610c1d565b60005a61138881101561377157600080fd5b61138881039050846040820482031161378957600080fd5b50823b61379557600080fd5b60008083516020850160008789f1949350505050565b600081815260018301602052604081205480156138945760006137cf600183614b8b565b85549091506000906137e390600190614b8b565b905081811461384857600086600001828154811061380357613803614cb2565b906000526020600020015490508087600001848154811061382657613826614cb2565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061385957613859614c9c565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106a9565b60009150506106a9565b6000818152600183016020526040812054612ab3908490849084906138ef575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106a9565b5060006106a9565b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081019190915290565b50805461393f90614bf6565b6000825580601f1061394f575050565b601f016020900490600052602060002090810190612a9a9190613a44565b8280548282559060005260206000209081019282156139c0579160200282015b828111156139c05781546001600160a01b0319166001600160a01b0384351617825560209092019160019091019061398d565b506139cc929150613a44565b5090565b8280546139dc90614bf6565b90600052602060002090601f0160209004810192826139fe57600085556139c0565b82601f10613a1757805160ff19168380011785556139c0565b828001600101855582156139c0579182015b828111156139c0578251825591602001919060010190613a29565b5b808211156139cc5760008155600101613a45565b80356001600160a01b0381168114613a7057600080fd5b919050565b60008083601f840112613a8757600080fd5b5081356001600160401b03811115613a9e57600080fd5b6020830191508360208260051b8501011115613ab957600080fd5b9250929050565b600082601f830112613ad157600080fd5b81356020613ae6613ae183614abd565b614a8d565b80838252828201915082860187848660051b8901011115613b0657600080fd5b60005b85811015613b855781356001600160401b03811115613b2757600080fd5b8801603f81018a13613b3857600080fd5b858101356040613b4a613ae183614ae0565b8281528c82848601011115613b5e57600080fd5b828285018a8301376000928101890192909252508552509284019290840190600101613b09565b5090979650505050505050565b600082601f830112613ba357600080fd5b81356020613bb3613ae183614abd565b8281528181019085830160e080860288018501891015613bd257600080fd5b60005b86811015613c835781838b031215613bec57600080fd5b613bf4614a42565b613bfd84613d7a565b8152613c0a878501613a59565b878201526040613c1b818601613d4c565b908201526060848101356001600160401b0381168114613c3a57600080fd5b908201526080613c4b858201613a59565b9082015260a0613c5c858201613d7a565b9082015260c0613c6d858201613a59565b9082015285529385019391810191600101613bd5565b509198975050505050505050565b80518015158114613a7057600080fd5b60008083601f840112613cb357600080fd5b5081356001600160401b03811115613cca57600080fd5b602083019150836020828501011115613ab957600080fd5b600082601f830112613cf357600080fd5b8151613d01613ae182614ae0565b818152846020838601011115613d1657600080fd5b610cd3826020830160208701614bca565b803561ffff81168114613a7057600080fd5b803562ffffff81168114613a7057600080fd5b803563ffffffff81168114613a7057600080fd5b805169ffffffffffffffffffff81168114613a7057600080fd5b80356001600160601b0381168114613a7057600080fd5b600060208284031215613da357600080fd5b612ab382613a59565b60008060408385031215613dbf57600080fd5b613dc883613a59565b9150613dd660208401613a59565b90509250929050565b60008060408385031215613df257600080fd5b613dfb83613a59565b9150602083013560048110613e0f57600080fd5b809150509250929050565b60008060008060608587031215613e3057600080fd5b613e3985613a59565b93506020850135925060408501356001600160401b03811115613e5b57600080fd5b613e6787828801613ca1565b95989497509550505050565b600080600080600060808688031215613e8b57600080fd5b613e9486613a59565b9450613ea260208701613d4c565b9350613eb060408701613a59565b925060608601356001600160401b03811115613ecb57600080fd5b613ed788828901613ca1565b969995985093965092949392505050565b60008060008060408587031215613efe57600080fd5b84356001600160401b0380821115613f1557600080fd5b613f2188838901613a75565b90965094506020870135915080821115613f3a57600080fd5b50613e6787828801613a75565b600080600060408486031215613f5c57600080fd5b83356001600160401b03811115613f7257600080fd5b613f7e86828701613a75565b9094509250613f91905060208501613a59565b90509250925092565b600080600060608486031215613faf57600080fd5b83356001600160401b0380821115613fc657600080fd5b818601915086601f830112613fda57600080fd5b81356020613fea613ae183614abd565b8083825282820191508286018b848660051b890101111561400a57600080fd5b600096505b8487101561402d57803583526001969096019591830191830161400f565b509750508701359250508082111561404457600080fd5b61405087838801613b92565b9350604086013591508082111561406657600080fd5b5061407386828701613ac0565b9150509250925092565b60006020828403121561408f57600080fd5b612ab382613c91565b600080604083850312156140ab57600080fd5b6140b483613c91565b915060208301516001600160401b038111156140cf57600080fd5b6140db85828601613ce2565b9150509250929050565b600080602083850312156140f857600080fd5b82356001600160401b0381111561410e57600080fd5b61411a85828601613ca1565b90969095509350505050565b60006020828403121561413857600080fd5b81516001600160401b0381111561414e57600080fd5b610cd384828501613ce2565b60006020828403121561416c57600080fd5b815160028110612ab357600080fd5b6000610180828403121561418e57600080fd5b614196614a6a565b61419f83613d4c565b81526141ad60208401613d4c565b60208201526141be60408401613d39565b60408201526141cf60608401613d4c565b60608201526141e060808401613d39565b60808201526141f160a08401613d27565b60a082015261420260c08401613d7a565b60c082015261421360e08401613d4c565b60e08201526101008381013590820152610120808401359082015261014061423c818501613a59565b9082015261016061424e848201613a59565b908201529392505050565b60006020828403121561426b57600080fd5b5035919050565b60006020828403121561428457600080fd5b5051919050565b6000806040838503121561429e57600080fd5b82359150613dd660208401613a59565b6000806000604084860312156142c357600080fd5b8335925060208401356001600160401b038111156142e057600080fd5b6142ec86828701613ca1565b9497909650939450505050565b6000806040838503121561430c57600080fd5b50508035926020909101359150565b6000806040838503121561432e57600080fd5b82359150613dd660208401613d4c565b6000806040838503121561435157600080fd5b82359150613dd660208401613d7a565b600080600080600060a0868803121561437957600080fd5b61438286613d60565b94506020860151935060408601519250606086015191506143a560808701613d60565b90509295509295909350565b6001600160a01b03169052565b8183526000602080850194508260005b858110156143fa576001600160a01b036143e783613a59565b16875295820195908201906001016143ce565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b8581101561444d57828403895261443b84835161445a565b98850198935090840190600101614423565b5091979650505050505050565b60008151808452614472816020860160208601614bca565b601f01601f19169290920160200192915050565b6002811061449657614496614c86565b9052565b805163ffffffff16825260208101516144bb602084018263ffffffff169052565b5060408101516144d2604084018262ffffff169052565b5060608101516144ea606084018263ffffffff169052565b506080810151614501608084018262ffffff169052565b5060a081015161451760a084018261ffff169052565b5060c081015161453260c08401826001600160601b03169052565b5060e081015161454a60e084018263ffffffff169052565b506101008181015190830152610120808201519083015261014080820151614574828501826143b1565b505061016080820151614589828501826143b1565b50505050565b600082516145a1818460208701614bca565b9190910192915050565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b03898116825263ffffffff891660208301526101006040830181905260009161460a8483018b61445a565b6001600160601b03998a16606086015297811660808501529590951660a0830152506001600160401b039290921660c083015290931660e090930192909252949350505050565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6040815260006146876040830186886143be565b828103602084015261469a8185876143be565b979650505050505050565b606080825281810185905260009060806001600160fb1b038711156146c957600080fd5b8660051b808983870137808501905081810160008152602083878403018188015281895180845260a093508385019150828b01945060005b8181101561479657855180516001600160601b03168452848101516001600160a01b039081168686015260408083015163ffffffff1690860152898201516001600160401b03168a8601528882015116888501528581015161476d878601826001600160601b03169052565b5060c09081015190614781858201836143b1565b50509483019460e09290920191600101614701565b505087810360408901526147aa818a614405565b9c9b505050505050505050505050565b6020808252825182820181905260009190848201906040850190845b818110156147f2578351835292840192918401916001016147d6565b50909695505050505050565b602081526000612ab3602083018461445a565b60a08152600061482460a083018861445a565b90508560208301528460408301528360608301528260808301529695505050505050565b600060208083526000845481600182811c91508083168061486a57607f831692505b85831081141561488857634e487b7160e01b85526022600452602485fd5b8786018381526020018180156148a557600181146148b6576148e1565b60ff198616825287820196506148e1565b60008b81526020902060005b868110156148db578154848201529085019089016148c2565b83019750505b50949998505050505050505050565b602081016004831061490457614904614c86565b91905290565b602081016106a98284614486565b6149228185614486565b61492f6020820184614486565b606060408201526000611608606083018461445a565b61018081016106a9828461449a565b600061022080830163ffffffff8751168452602060018060601b038189015116818601526040880151604086015260608801516060860152614999608086018861449a565b6102008501929092528451908190526102408401918086019160005b818110156149da5783516001600160a01b0316855293820193928201926001016149b5565b509298975050505050505050565b6001600160601b0391909116815260200190565b6001600160601b039290921682526001600160a01b0316602082015260400190565b6001600160601b0383168152604060208201819052600090610cd39083018461445a565b60405160e081016001600160401b0381118282101715614a6457614a64614cc8565b60405290565b60405161018081016001600160401b0381118282101715614a6457614a64614cc8565b604051601f8201601f191681016001600160401b0381118282101715614ab557614ab5614cc8565b604052919050565b60006001600160401b03821115614ad657614ad6614cc8565b5060051b60200190565b60006001600160401b03821115614af957614af9614cc8565b50601f01601f191660200190565b60008219821115614b1a57614b1a614c70565b500190565b60006001600160601b03828116848216808303821115614b4157614b41614c70565b01949350505050565b600082614b6757634e487b7160e01b600052601260045260246000fd5b500490565b6000816000190483118215151615614b8657614b86614c70565b500290565b600082821015614b9d57614b9d614c70565b500390565b60006001600160601b0383811690831681811015614bc257614bc2614c70565b039392505050565b60005b83811015614be5578181015183820152602001614bcd565b838111156145895750506000910152565b600181811c90821680614c0a57607f821691505b60208210811415614c2b57634e487b7160e01b600052602260045260246000fd5b50919050565b6000600019821415614c4557614c45614c70565b5060010190565b600063ffffffff80831681811415614c6657614c66614c70565b6001019392505050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea264697066735822122057730712b9822162aa1238d197ed30a9fc536715972647f911c347bd60c76a7964736f6c63430008060033",
}

// KeeperRegistry12ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistry12MetaData.ABI instead.
var KeeperRegistry12ABI = KeeperRegistry12MetaData.ABI

// KeeperRegistry12Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistry12MetaData.Bin instead.
var KeeperRegistry12Bin = KeeperRegistry12MetaData.Bin

// DeployKeeperRegistry12 deploys a new Ethereum contract, binding an instance of KeeperRegistry12 to it.
func DeployKeeperRegistry12(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkEthFeed common.Address, fastGasFeed common.Address, config Config1_2) (common.Address, *types.Transaction, *KeeperRegistry12, error) {
	parsed, err := KeeperRegistry12MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry12Bin), backend, link, linkEthFeed, fastGasFeed, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry12{KeeperRegistry12Caller: KeeperRegistry12Caller{contract: contract}, KeeperRegistry12Transactor: KeeperRegistry12Transactor{contract: contract}, KeeperRegistry12Filterer: KeeperRegistry12Filterer{contract: contract}}, nil
}

// KeeperRegistry12 is an auto generated Go binding around an Ethereum contract.
type KeeperRegistry12 struct {
	KeeperRegistry12Caller     // Read-only binding to the contract
	KeeperRegistry12Transactor // Write-only binding to the contract
	KeeperRegistry12Filterer   // Log filterer for contract events
}

// KeeperRegistry12Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistry12Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry12Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistry12Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry12Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistry12Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistry12Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistry12Session struct {
	Contract     *KeeperRegistry12 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistry12CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistry12CallerSession struct {
	Contract *KeeperRegistry12Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// KeeperRegistry12TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistry12TransactorSession struct {
	Contract     *KeeperRegistry12Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// KeeperRegistry12Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistry12Raw struct {
	Contract *KeeperRegistry12 // Generic contract binding to access the raw methods on
}

// KeeperRegistry12CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistry12CallerRaw struct {
	Contract *KeeperRegistry12Caller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistry12TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistry12TransactorRaw struct {
	Contract *KeeperRegistry12Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistry12 creates a new instance of KeeperRegistry12, bound to a specific deployed contract.
func NewKeeperRegistry12(address common.Address, backend bind.ContractBackend) (*KeeperRegistry12, error) {
	contract, err := bindKeeperRegistry12(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12{KeeperRegistry12Caller: KeeperRegistry12Caller{contract: contract}, KeeperRegistry12Transactor: KeeperRegistry12Transactor{contract: contract}, KeeperRegistry12Filterer: KeeperRegistry12Filterer{contract: contract}}, nil
}

// NewKeeperRegistry12Caller creates a new read-only instance of KeeperRegistry12, bound to a specific deployed contract.
func NewKeeperRegistry12Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistry12Caller, error) {
	contract, err := bindKeeperRegistry12(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12Caller{contract: contract}, nil
}

// NewKeeperRegistry12Transactor creates a new write-only instance of KeeperRegistry12, bound to a specific deployed contract.
func NewKeeperRegistry12Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistry12Transactor, error) {
	contract, err := bindKeeperRegistry12(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12Transactor{contract: contract}, nil
}

// NewKeeperRegistry12Filterer creates a new log filterer instance of KeeperRegistry12, bound to a specific deployed contract.
func NewKeeperRegistry12Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistry12Filterer, error) {
	contract, err := bindKeeperRegistry12(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12Filterer{contract: contract}, nil
}

// bindKeeperRegistry12 binds a generic wrapper to an already deployed contract.
func bindKeeperRegistry12(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistry12ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry12 *KeeperRegistry12Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry12.Contract.KeeperRegistry12Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry12 *KeeperRegistry12Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.KeeperRegistry12Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry12 *KeeperRegistry12Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.KeeperRegistry12Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistry12 *KeeperRegistry12CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry12.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistry12 *KeeperRegistry12TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistry12 *KeeperRegistry12TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.contract.Transact(opts, method, params...)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Caller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Session) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry12.Contract.FASTGASFEED(&_KeeperRegistry12.CallOpts)
}

// FASTGASFEED is a free data retrieval call binding the contract method 0x4584a419.
//
// Solidity: function FAST_GAS_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry12.Contract.FASTGASFEED(&_KeeperRegistry12.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Session) LINK() (common.Address, error) {
	return _KeeperRegistry12.Contract.LINK(&_KeeperRegistry12.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry12.Contract.LINK(&_KeeperRegistry12.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Session) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry12.Contract.LINKETHFEED(&_KeeperRegistry12.CallOpts)
}

// LINKETHFEED is a free data retrieval call binding the contract method 0xad178361.
//
// Solidity: function LINK_ETH_FEED() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry12.Contract.LINKETHFEED(&_KeeperRegistry12.CallOpts)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry12 *KeeperRegistry12Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry12.Contract.GetActiveUpkeepIDs(&_KeeperRegistry12.CallOpts, startIndex, maxCount)
}

// GetActiveUpkeepIDs is a free data retrieval call binding the contract method 0x06e3b632.
//
// Solidity: function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) view returns(uint256[])
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry12.Contract.GetActiveUpkeepIDs(&_KeeperRegistry12.CallOpts, startIndex, maxCount)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getKeeperInfo", query)

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
func (_KeeperRegistry12 *KeeperRegistry12Session) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry12.Contract.GetKeeperInfo(&_KeeperRegistry12.CallOpts, query)
}

// GetKeeperInfo is a free data retrieval call binding the contract method 0x1e12b8a5.
//
// Solidity: function getKeeperInfo(address query) view returns(address payee, bool active, uint96 balance)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetKeeperInfo(query common.Address) (struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}, error) {
	return _KeeperRegistry12.Contract.GetKeeperInfo(&_KeeperRegistry12.CallOpts, query)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry12 *KeeperRegistry12Session) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry12.Contract.GetMaxPaymentForGas(&_KeeperRegistry12.CallOpts, gasLimit)
}

// GetMaxPaymentForGas is a free data retrieval call binding the contract method 0x93f0c1fc.
//
// Solidity: function getMaxPaymentForGas(uint256 gasLimit) view returns(uint96 maxPayment)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry12.Contract.GetMaxPaymentForGas(&_KeeperRegistry12.CallOpts, gasLimit)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry12 *KeeperRegistry12Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry12.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry12.CallOpts, id)
}

// GetMinBalanceForUpkeep is a free data retrieval call binding the contract method 0xb657bc9c.
//
// Solidity: function getMinBalanceForUpkeep(uint256 id) view returns(uint96 minBalance)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry12.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry12.CallOpts, id)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry12.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry12.CallOpts, peer)
}

// GetPeerRegistryMigrationPermission is a free data retrieval call binding the contract method 0xfaa3e996.
//
// Solidity: function getPeerRegistryMigrationPermission(address peer) view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry12.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry12.CallOpts, peer)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetState(opts *bind.CallOpts) (struct {
	State   State1_2
	Config  Config1_2
	Keepers []common.Address
}, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		State   State1_2
		Config  Config1_2
		Keepers []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State1_2)).(*State1_2)
	outstruct.Config = *abi.ConvertType(out[1], new(Config1_2)).(*Config1_2)
	outstruct.Keepers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry12 *KeeperRegistry12Session) GetState() (struct {
	State   State1_2
	Config  Config1_2
	Keepers []common.Address
}, error) {
	return _KeeperRegistry12.Contract.GetState(&_KeeperRegistry12.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns((uint32,uint96,uint256,uint256) state, (uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config, address[] keepers)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetState() (struct {
	State   State1_2
	Config  Config1_2
	Keepers []common.Address
}, error) {
	return _KeeperRegistry12.Contract.GetState(&_KeeperRegistry12.CallOpts)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry12 *KeeperRegistry12Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "getUpkeep", id)

	outstruct := new(struct {
		Target              common.Address
		ExecuteGas          uint32
		CheckData           []byte
		Balance             *big.Int
		LastKeeper          common.Address
		Admin               common.Address
		MaxValidBlocknumber uint64
		AmountSpent         *big.Int
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

	return *outstruct, err

}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry12 *KeeperRegistry12Session) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	return _KeeperRegistry12.Contract.GetUpkeep(&_KeeperRegistry12.CallOpts, id)
}

// GetUpkeep is a free data retrieval call binding the contract method 0xc7c3a19a.
//
// Solidity: function getUpkeep(uint256 id) view returns(address target, uint32 executeGas, bytes checkData, uint96 balance, address lastKeeper, address admin, uint64 maxValidBlocknumber, uint96 amountSpent)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) GetUpkeep(id *big.Int) (struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}, error) {
	return _KeeperRegistry12.Contract.GetUpkeep(&_KeeperRegistry12.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12Session) Owner() (common.Address, error) {
	return _KeeperRegistry12.Contract.Owner(&_KeeperRegistry12.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry12.Contract.Owner(&_KeeperRegistry12.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry12 *KeeperRegistry12Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry12 *KeeperRegistry12Session) Paused() (bool, error) {
	return _KeeperRegistry12.Contract.Paused(&_KeeperRegistry12.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) Paused() (bool, error) {
	return _KeeperRegistry12.Contract.Paused(&_KeeperRegistry12.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry12 *KeeperRegistry12Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry12 *KeeperRegistry12Session) TypeAndVersion() (string, error) {
	return _KeeperRegistry12.Contract.TypeAndVersion(&_KeeperRegistry12.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry12.Contract.TypeAndVersion(&_KeeperRegistry12.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry12.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry12.Contract.UpkeepTranscoderVersion(&_KeeperRegistry12.CallOpts)
}

// UpkeepTranscoderVersion is a free data retrieval call binding the contract method 0x48013d7b.
//
// Solidity: function upkeepTranscoderVersion() view returns(uint8)
func (_KeeperRegistry12 *KeeperRegistry12CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry12.Contract.UpkeepTranscoderVersion(&_KeeperRegistry12.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AcceptOwnership(&_KeeperRegistry12.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AcceptOwnership(&_KeeperRegistry12.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "acceptPayeeship", keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AcceptPayeeship(&_KeeperRegistry12.TransactOpts, keeper)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address keeper) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AcceptPayeeship(&_KeeperRegistry12.TransactOpts, keeper)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "addFunds", id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AddFunds(&_KeeperRegistry12.TransactOpts, id, amount)
}

// AddFunds is a paid mutator transaction binding the contract method 0x948108f7.
//
// Solidity: function addFunds(uint256 id, uint96 amount) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.AddFunds(&_KeeperRegistry12.TransactOpts, id, amount)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "cancelUpkeep", id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.CancelUpkeep(&_KeeperRegistry12.TransactOpts, id)
}

// CancelUpkeep is a paid mutator transaction binding the contract method 0xc8048022.
//
// Solidity: function cancelUpkeep(uint256 id) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.CancelUpkeep(&_KeeperRegistry12.TransactOpts, id)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry12 *KeeperRegistry12Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "checkUpkeep", id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry12 *KeeperRegistry12Session) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.CheckUpkeep(&_KeeperRegistry12.TransactOpts, id, from)
}

// CheckUpkeep is a paid mutator transaction binding the contract method 0xc41b813a.
//
// Solidity: function checkUpkeep(uint256 id, address from) returns(bytes performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth)
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.CheckUpkeep(&_KeeperRegistry12.TransactOpts, id, from)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.MigrateUpkeeps(&_KeeperRegistry12.TransactOpts, ids, destination)
}

// MigrateUpkeeps is a paid mutator transaction binding the contract method 0x85c1b0ba.
//
// Solidity: function migrateUpkeeps(uint256[] ids, address destination) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.MigrateUpkeeps(&_KeeperRegistry12.TransactOpts, ids, destination)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.OnTokenTransfer(&_KeeperRegistry12.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.OnTokenTransfer(&_KeeperRegistry12.TransactOpts, sender, amount, data)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.Pause(&_KeeperRegistry12.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.Pause(&_KeeperRegistry12.TransactOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry12 *KeeperRegistry12Transactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "performUpkeep", id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry12 *KeeperRegistry12Session) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.PerformUpkeep(&_KeeperRegistry12.TransactOpts, id, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x7bbaf1ea.
//
// Solidity: function performUpkeep(uint256 id, bytes performData) returns(bool success)
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.PerformUpkeep(&_KeeperRegistry12.TransactOpts, id, performData)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.ReceiveUpkeeps(&_KeeperRegistry12.TransactOpts, encodedUpkeeps)
}

// ReceiveUpkeeps is a paid mutator transaction binding the contract method 0x8e86139b.
//
// Solidity: function receiveUpkeeps(bytes encodedUpkeeps) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.ReceiveUpkeeps(&_KeeperRegistry12.TransactOpts, encodedUpkeeps)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "recoverFunds")
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.RecoverFunds(&_KeeperRegistry12.TransactOpts)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0xb79550be.
//
// Solidity: function recoverFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.RecoverFunds(&_KeeperRegistry12.TransactOpts)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry12 *KeeperRegistry12Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry12 *KeeperRegistry12Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.RegisterUpkeep(&_KeeperRegistry12.TransactOpts, target, gasLimit, admin, checkData)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0xda5c6741.
//
// Solidity: function registerUpkeep(address target, uint32 gasLimit, address admin, bytes checkData) returns(uint256 id)
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.RegisterUpkeep(&_KeeperRegistry12.TransactOpts, target, gasLimit, admin, checkData)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) SetConfig(opts *bind.TransactOpts, config Config1_2) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "setConfig", config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) SetConfig(config Config1_2) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetConfig(&_KeeperRegistry12.TransactOpts, config)
}

// SetConfig is a paid mutator transaction binding the contract method 0xef47a0ce.
//
// Solidity: function setConfig((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) SetConfig(config Config1_2) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetConfig(&_KeeperRegistry12.TransactOpts, config)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "setKeepers", keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetKeepers(&_KeeperRegistry12.TransactOpts, keepers, payees)
}

// SetKeepers is a paid mutator transaction binding the contract method 0xb7fdb436.
//
// Solidity: function setKeepers(address[] keepers, address[] payees) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetKeepers(&_KeeperRegistry12.TransactOpts, keepers, payees)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry12.TransactOpts, peer, permission)
}

// SetPeerRegistryMigrationPermission is a paid mutator transaction binding the contract method 0x187256e8.
//
// Solidity: function setPeerRegistryMigrationPermission(address peer, uint8 permission) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry12.TransactOpts, peer, permission)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetUpkeepGasLimit(&_KeeperRegistry12.TransactOpts, id, gasLimit)
}

// SetUpkeepGasLimit is a paid mutator transaction binding the contract method 0xa72aa27e.
//
// Solidity: function setUpkeepGasLimit(uint256 id, uint32 gasLimit) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.SetUpkeepGasLimit(&_KeeperRegistry12.TransactOpts, id, gasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.TransferOwnership(&_KeeperRegistry12.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.TransferOwnership(&_KeeperRegistry12.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.TransferPayeeship(&_KeeperRegistry12.TransactOpts, keeper, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address keeper, address proposed) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.TransferPayeeship(&_KeeperRegistry12.TransactOpts, keeper, proposed)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.Unpause(&_KeeperRegistry12.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.Unpause(&_KeeperRegistry12.TransactOpts)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "withdrawFunds", id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawFunds(&_KeeperRegistry12.TransactOpts, id, to)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x744bfe61.
//
// Solidity: function withdrawFunds(uint256 id, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawFunds(&_KeeperRegistry12.TransactOpts, id, to)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "withdrawOwnerFunds")
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawOwnerFunds(&_KeeperRegistry12.TransactOpts)
}

// WithdrawOwnerFunds is a paid mutator transaction binding the contract method 0x7d9b97e0.
//
// Solidity: function withdrawOwnerFunds() returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawOwnerFunds(&_KeeperRegistry12.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.contract.Transact(opts, "withdrawPayment", from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawPayment(&_KeeperRegistry12.TransactOpts, from, to)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0xa710b221.
//
// Solidity: function withdrawPayment(address from, address to) returns()
func (_KeeperRegistry12 *KeeperRegistry12TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry12.Contract.WithdrawPayment(&_KeeperRegistry12.TransactOpts, from, to)
}

// KeeperRegistry12ConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the KeeperRegistry12 contract.
type KeeperRegistry12ConfigSetIterator struct {
	Event *KeeperRegistry12ConfigSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12ConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12ConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12ConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12ConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12ConfigSet represents a ConfigSet event raised by the KeeperRegistry12 contract.
type KeeperRegistry12ConfigSet struct {
	Config Config1_2
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry12ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12ConfigSetIterator{contract: _KeeperRegistry12.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325.
//
// Solidity: event ConfigSet((uint32,uint32,uint24,uint32,uint24,uint16,uint96,uint32,uint256,uint256,address,address) config)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12ConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12ConfigSet)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseConfigSet(log types.Log) (*KeeperRegistry12ConfigSet, error) {
	event := new(KeeperRegistry12ConfigSet)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12FundsAddedIterator is returned from FilterFundsAdded and is used to iterate over the raw logs and unpacked data for FundsAdded events raised by the KeeperRegistry12 contract.
type KeeperRegistry12FundsAddedIterator struct {
	Event *KeeperRegistry12FundsAdded // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12FundsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12FundsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12FundsAdded represents a FundsAdded event raised by the KeeperRegistry12 contract.
type KeeperRegistry12FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsAdded is a free log retrieval operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry12FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12FundsAddedIterator{contract: _KeeperRegistry12.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

// WatchFundsAdded is a free log subscription operation binding the contract event 0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203.
//
// Solidity: event FundsAdded(uint256 indexed id, address indexed from, uint96 amount)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12FundsAdded)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistry12FundsAdded, error) {
	event := new(KeeperRegistry12FundsAdded)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12FundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the KeeperRegistry12 contract.
type KeeperRegistry12FundsWithdrawnIterator struct {
	Event *KeeperRegistry12FundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12FundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12FundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12FundsWithdrawn represents a FundsWithdrawn event raised by the KeeperRegistry12 contract.
type KeeperRegistry12FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry12FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12FundsWithdrawnIterator{contract: _KeeperRegistry12.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318.
//
// Solidity: event FundsWithdrawn(uint256 indexed id, uint256 amount, address to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12FundsWithdrawn)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistry12FundsWithdrawn, error) {
	event := new(KeeperRegistry12FundsWithdrawn)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12KeepersUpdatedIterator is returned from FilterKeepersUpdated and is used to iterate over the raw logs and unpacked data for KeepersUpdated events raised by the KeeperRegistry12 contract.
type KeeperRegistry12KeepersUpdatedIterator struct {
	Event *KeeperRegistry12KeepersUpdated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12KeepersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12KeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12KeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12KeepersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12KeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12KeepersUpdated represents a KeepersUpdated event raised by the KeeperRegistry12 contract.
type KeeperRegistry12KeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterKeepersUpdated is a free log retrieval operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistry12KeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12KeepersUpdatedIterator{contract: _KeeperRegistry12.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepersUpdated is a free log subscription operation binding the contract event 0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f.
//
// Solidity: event KeepersUpdated(address[] keepers, address[] payees)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12KeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12KeepersUpdated)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistry12KeepersUpdated, error) {
	event := new(KeeperRegistry12KeepersUpdated)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12OwnerFundsWithdrawnIterator is returned from FilterOwnerFundsWithdrawn and is used to iterate over the raw logs and unpacked data for OwnerFundsWithdrawn events raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistry12OwnerFundsWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12OwnerFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12OwnerFundsWithdrawn represents a OwnerFundsWithdrawn event raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerFundsWithdrawn is a free log retrieval operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry12OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12OwnerFundsWithdrawnIterator{contract: _KeeperRegistry12.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchOwnerFundsWithdrawn is a free log subscription operation binding the contract event 0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1.
//
// Solidity: event OwnerFundsWithdrawn(uint96 amount)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12OwnerFundsWithdrawn)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry12OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistry12OwnerFundsWithdrawn)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistry12OwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry12OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12OwnershipTransferRequestedIterator{contract: _KeeperRegistry12.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12OwnershipTransferRequested)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry12OwnershipTransferRequested, error) {
	event := new(KeeperRegistry12OwnershipTransferRequested)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnershipTransferredIterator struct {
	Event *KeeperRegistry12OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12OwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistry12 contract.
type KeeperRegistry12OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry12OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12OwnershipTransferredIterator{contract: _KeeperRegistry12.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12OwnershipTransferred)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistry12OwnershipTransferred, error) {
	event := new(KeeperRegistry12OwnershipTransferred)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the KeeperRegistry12 contract.
type KeeperRegistry12PausedIterator struct {
	Event *KeeperRegistry12Paused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12Paused represents a Paused event raised by the KeeperRegistry12 contract.
type KeeperRegistry12Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry12PausedIterator, error) {

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12PausedIterator{contract: _KeeperRegistry12.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12Paused)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParsePaused(log types.Log) (*KeeperRegistry12Paused, error) {
	event := new(KeeperRegistry12Paused)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12PayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the KeeperRegistry12 contract.
type KeeperRegistry12PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistry12PayeeshipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12PayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12PayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the KeeperRegistry12 contract.
type KeeperRegistry12PayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry12PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12PayeeshipTransferRequestedIterator{contract: _KeeperRegistry12.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12PayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12PayeeshipTransferRequested)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry12PayeeshipTransferRequested, error) {
	event := new(KeeperRegistry12PayeeshipTransferRequested)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12PayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the KeeperRegistry12 contract.
type KeeperRegistry12PayeeshipTransferredIterator struct {
	Event *KeeperRegistry12PayeeshipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12PayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12PayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12PayeeshipTransferred represents a PayeeshipTransferred event raised by the KeeperRegistry12 contract.
type KeeperRegistry12PayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry12PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12PayeeshipTransferredIterator{contract: _KeeperRegistry12.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12PayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12PayeeshipTransferred)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry12PayeeshipTransferred, error) {
	event := new(KeeperRegistry12PayeeshipTransferred)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12PaymentWithdrawnIterator is returned from FilterPaymentWithdrawn and is used to iterate over the raw logs and unpacked data for PaymentWithdrawn events raised by the KeeperRegistry12 contract.
type KeeperRegistry12PaymentWithdrawnIterator struct {
	Event *KeeperRegistry12PaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12PaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12PaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12PaymentWithdrawn represents a PaymentWithdrawn event raised by the KeeperRegistry12 contract.
type KeeperRegistry12PaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaymentWithdrawn is a free log retrieval operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry12PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12PaymentWithdrawnIterator{contract: _KeeperRegistry12.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

// WatchPaymentWithdrawn is a free log subscription operation binding the contract event 0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698.
//
// Solidity: event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12PaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12PaymentWithdrawn)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry12PaymentWithdrawn, error) {
	event := new(KeeperRegistry12PaymentWithdrawn)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UnpausedIterator struct {
	Event *KeeperRegistry12Unpaused // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12Unpaused represents a Unpaused event raised by the KeeperRegistry12 contract.
type KeeperRegistry12Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry12UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UnpausedIterator{contract: _KeeperRegistry12.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12Unpaused)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUnpaused(log types.Log) (*KeeperRegistry12Unpaused, error) {
	event := new(KeeperRegistry12Unpaused)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepCanceledIterator is returned from FilterUpkeepCanceled and is used to iterate over the raw logs and unpacked data for UpkeepCanceled events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepCanceledIterator struct {
	Event *KeeperRegistry12UpkeepCanceled // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepCanceled represents a UpkeepCanceled event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpkeepCanceled is a free log retrieval operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry12UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepCanceledIterator{contract: _KeeperRegistry12.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

// WatchUpkeepCanceled is a free log subscription operation binding the contract event 0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181.
//
// Solidity: event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepCanceled)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistry12UpkeepCanceled, error) {
	event := new(KeeperRegistry12UpkeepCanceled)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepGasLimitSetIterator is returned from FilterUpkeepGasLimitSet and is used to iterate over the raw logs and unpacked data for UpkeepGasLimitSet events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistry12UpkeepGasLimitSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepGasLimitSet represents a UpkeepGasLimitSet event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpkeepGasLimitSet is a free log retrieval operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry12UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepGasLimitSetIterator{contract: _KeeperRegistry12.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchUpkeepGasLimitSet is a free log subscription operation binding the contract event 0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c.
//
// Solidity: event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepGasLimitSet)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry12UpkeepGasLimitSet, error) {
	event := new(KeeperRegistry12UpkeepGasLimitSet)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepMigratedIterator is returned from FilterUpkeepMigrated and is used to iterate over the raw logs and unpacked data for UpkeepMigrated events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepMigratedIterator struct {
	Event *KeeperRegistry12UpkeepMigrated // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepMigrated represents a UpkeepMigrated event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpkeepMigrated is a free log retrieval operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry12UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepMigratedIterator{contract: _KeeperRegistry12.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

// WatchUpkeepMigrated is a free log subscription operation binding the contract event 0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff.
//
// Solidity: event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepMigrated)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistry12UpkeepMigrated, error) {
	event := new(KeeperRegistry12UpkeepMigrated)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepPerformedIterator is returned from FilterUpkeepPerformed and is used to iterate over the raw logs and unpacked data for UpkeepPerformed events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepPerformedIterator struct {
	Event *KeeperRegistry12UpkeepPerformed // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepPerformedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepPerformedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepPerformed represents a UpkeepPerformed event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepPerformed struct {
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistry12UpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepPerformedIterator{contract: _KeeperRegistry12.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

// WatchUpkeepPerformed is a free log subscription operation binding the contract event 0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6.
//
// Solidity: event UpkeepPerformed(uint256 indexed id, bool indexed success, address indexed from, uint96 payment, bytes performData)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepPerformed)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistry12UpkeepPerformed, error) {
	event := new(KeeperRegistry12UpkeepPerformed)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepReceivedIterator is returned from FilterUpkeepReceived and is used to iterate over the raw logs and unpacked data for UpkeepReceived events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepReceivedIterator struct {
	Event *KeeperRegistry12UpkeepReceived // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepReceived represents a UpkeepReceived event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpkeepReceived is a free log retrieval operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry12UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepReceivedIterator{contract: _KeeperRegistry12.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

// WatchUpkeepReceived is a free log subscription operation binding the contract event 0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71.
//
// Solidity: event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepReceived)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistry12UpkeepReceived, error) {
	event := new(KeeperRegistry12UpkeepReceived)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistry12UpkeepRegisteredIterator is returned from FilterUpkeepRegistered and is used to iterate over the raw logs and unpacked data for UpkeepRegistered events raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepRegisteredIterator struct {
	Event *KeeperRegistry12UpkeepRegistered // Event containing the contract specifics and raw log

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
func (it *KeeperRegistry12UpkeepRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry12UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistry12UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
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
func (it *KeeperRegistry12UpkeepRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistry12UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistry12UpkeepRegistered represents a UpkeepRegistered event raised by the KeeperRegistry12 contract.
type KeeperRegistry12UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpkeepRegistered is a free log retrieval operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry12UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry12UpkeepRegisteredIterator{contract: _KeeperRegistry12.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

// WatchUpkeepRegistered is a free log subscription operation binding the contract event 0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012.
//
// Solidity: event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin)
func (_KeeperRegistry12 *KeeperRegistry12Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry12UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry12.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistry12UpkeepRegistered)
				if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
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
func (_KeeperRegistry12 *KeeperRegistry12Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistry12UpkeepRegistered, error) {
	event := new(KeeperRegistry12UpkeepRegistered)
	if err := _KeeperRegistry12.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
