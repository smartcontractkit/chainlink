// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic1_3

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

type Config struct {
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

var KeeperRegistryLogicMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase.PaymentModel\",\"name\":\"paymentModel\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"registryGasOverhead\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ARB_NITRO_ORACLE\",\"outputs\":[{\"internalType\":\"contractArbGasInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_FEE_DATA_PADDING\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_INPUT_DATA\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMISM_ORACLE\",\"outputs\":[{\"internalType\":\"contractOVM_GasPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAYMENT_MODEL\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REGISTRY_GAS_OVERHEAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101e060405260486101608181529062005a3a6101803980516200002c916004916020909101906200027e565b50604051806101400160405280610110815260200162005a826101109139805162000060916005916020909101906200027e565b507f420000000000000000000000000000000000000f00000000000000000000000060e0526c6c00000000000000000000000061010052348015620000a457600080fd5b5060405162005b9238038062005b92833981016040819052620000c79162000341565b84848484843380600081620001235760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000156576200015681620001d2565b5050600160029081556003805460ff191690558691508111156200017e576200017e620003eb565b610120816002811115620001965762000196620003eb565b60f81b905250610140939093526001600160601b0319606092831b811660805290821b811660a05291901b1660c0525062000401945050505050565b6001600160a01b0381163314156200022d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200011a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280546200028c90620003ae565b90600052602060002090601f016020900481019282620002b05760008555620002fb565b82601f10620002cb57805160ff1916838001178555620002fb565b82800160010185558215620002fb579182015b82811115620002fb578251825591602001919060010190620002de565b50620003099291506200030d565b5090565b5b808211156200030957600081556001016200030e565b80516001600160a01b03811681146200033c57600080fd5b919050565b600080600080600060a086880312156200035a57600080fd5b8551600381106200036a57600080fd5b60208701519095509350620003826040870162000324565b9250620003926060870162000324565b9150620003a26080870162000324565b90509295509295909350565b600181811c90821680620003c357607f821691505b60208210811415620003e557634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160601c6101005160601c6101205160f81c6101405161556a620004d0600039600081816102a20152613fed015260008181610379015281816140400152614214015260008181610310015261424c01526000818161033f015261418301526000818161027b0152613d2101526000818161042c0152613e0201526000818161022201528181610a7901528181610cd601528181611543015281816119a801528181611d260152818161233601526123c9015261556a6000f3fe608060405234801561001057600080fd5b50600436106101f05760003560e01c80638e86139b1161010f578063b79550be116100a2578063c804802211610071578063c8048022146104bb578063da5c6741146104ce578063eb5dcd6c146104e1578063f2fde38b146104f457600080fd5b8063b79550be14610474578063b7fdb4361461047c578063be3081771461048f578063c41b813a1461049757600080fd5b8063a72aa27e116100de578063a72aa27e14610414578063ad17836114610427578063b121e1471461044e578063b148ab6b1461046157600080fd5b80638e86139b146103c6578063948108f7146103d9578063a0ad00cf146103ec578063a710b2211461040157600080fd5b806379ba509711610187578063850cce3411610156578063850cce341461033a57806385c1b0ba146103615780638811cbe8146103745780638da5cb5b146103a857600080fd5b806379ba5097146102fb5780637d9b97e0146103035780637f37618e1461030b5780638456cb591461033257600080fd5b80634584a419116101c35780634584a419146102765780635077b2101461029d5780635c975abb146102d2578063744bfe61146102e857600080fd5b8063187256e8146101f55780631a2af0111461020a5780631b6b6d231461021d5780633f4ba83a1461026e575b600080fd5b6102086102033660046148c4565b610507565b005b610208610218366004614c58565b610578565b6102447f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6102086107a3565b6102447f000000000000000000000000000000000000000000000000000000000000000081565b6102c47f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610265565b60035460ff166040519015158152602001610265565b6102086102f6366004614c58565b6107b5565b610208610afc565b610208610bfe565b6102447f000000000000000000000000000000000000000000000000000000000000000081565b610208610d6c565b6102447f000000000000000000000000000000000000000000000000000000000000000081565b61020861036f3660046149e1565b610d7c565b61039b7f000000000000000000000000000000000000000000000000000000000000000081565b6040516102659190615150565b60005473ffffffffffffffffffffffffffffffffffffffff16610244565b6102086103d4366004614b8e565b6115cd565b6102086103e7366004614c9e565b6117eb565b6103f4611a84565b6040516102659190615063565b61020861040f366004614891565b611b12565b610208610422366004614c7b565b611daa565b6102447f000000000000000000000000000000000000000000000000000000000000000081565b61020861045c366004614876565b611fe3565b61020861046f366004614c26565b6120db565b6102086122fd565b61020861048a366004614975565b612468565b6103f46127c9565b6104aa6104a5366004614c58565b6127d6565b604051610265959493929190615076565b6102086104c9366004614c26565b612ae8565b6102c46104dc3660046148ff565b612e62565b6102086104ef366004614891565b613059565b610208610502366004614876565b6131b8565b61050f6131cc565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600e6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600183600381111561056f5761056f615493565b02179055505050565b60008281526009602090815260409182902082516101008101845281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff808216608085015264010000000082041660a084015268010000000000000000810490911660c083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560e08201526106578161324d565b73ffffffffffffffffffffffffffffffffffffffff82163314156106a7576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166106f4576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600c602052604090205473ffffffffffffffffffffffffffffffffffffffff83811691161461079e576000838152600c602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff861690811790915590519091339186917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b505050565b6107ab6131cc565b6107b36132fa565b565b73ffffffffffffffffffffffffffffffffffffffff8116610802576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526009602090815260409182902082516101008101845281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff9081169584019590955260018401549081169583019590955290930482166060840181905260029091015463ffffffff808216608086015264010000000082041660a085015268010000000000000000810490921660c08401527c010000000000000000000000000000000000000000000000000000000090910460ff16151560e08301523314610915576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b438160a0015163ffffffff161115610959576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600960205260409020546014546bffffffffffffffffffffffff90911690610987908290615343565b60145560008481526009602090815260409182902080547fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016905581516bffffffffffffffffffffffff8416815273ffffffffffffffffffffffffffffffffffffffff86169181019190915285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a26040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b158015610abd57600080fd5b505af1158015610ad1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610af59190614b19565b5050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b82576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610c066131cc565b6013546014546bffffffffffffffffffffffff90911690610c28908290615343565b601455601380547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b602060405180830381600087803b158015610d3057600080fd5b505af1158015610d44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d689190614b19565b5050565b610d746131cc565b6107b36133db565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152600e602052604090205460ff166003811115610db857610db8615493565b14158015610e005750600373ffffffffffffffffffffffffffffffffffffffff82166000908152600e602052604090205460ff166003811115610dfd57610dfd615493565b14155b15610e37576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60155473ffffffffffffffffffffffffffffffffffffffff16610e86576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81610ebd576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115610f1957610f19615520565b604051908082528060200260200182016040528015610f4c57816020015b6060815260200190600190039081610f375790505b50905060008667ffffffffffffffff811115610f6a57610f6a615520565b604051908082528060200260200182016040528015610ff757816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610f885790505b50905060005b8781101561130157888882818110611017576110176154f1565b6020908102929092013560008181526009845260409081902081516101008101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811698840198909852600184015490811694830194909452909204851660608301526002015463ffffffff808216608084015264010000000082041660a083015268010000000000000000810490941660c08201527c010000000000000000000000000000000000000000000000000000000090930460ff16151560e0840152975090955061110590508561324d565b84828281518110611118576111186154f1565b6020026020010181905250600d60008781526020019081526020016000208054611141906153b3565b80601f016020809104026020016040519081016040528092919081815260200182805461116d906153b3565b80156111ba5780601f1061118f576101008083540402835291602001916111ba565b820191906000526020600020905b81548152906001019060200180831161119d57829003601f168201915b50505050508382815181106111d1576111d16154f1565b602090810291909101015284516111f6906bffffffffffffffffffffffff1685615283565b60008781526009602090815260408083208381556001810184905560020180547fffffff0000000000000000000000000000000000000000000000000000000000169055600d909152812091955061124e9190614401565b6000868152600c6020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905561128d60078761349b565b508451604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8916602083015287917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a2806112f981615407565b915050610ffd565b50826014546113109190615343565b60145560405160009061132d908a908a9085908790602001614f05565b60405160208183030381529060405290508673ffffffffffffffffffffffffffffffffffffffff16638e86139b601560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60018b73ffffffffffffffffffffffffffffffffffffffff166348013d7b6040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156113e257600080fd5b505af11580156113f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141a9190614c05565b866040518463ffffffff1660e01b81526004016114399392919061516a565b60006040518083038186803b15801561145157600080fd5b505afa158015611465573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114ab9190810190614bd0565b6040518263ffffffff1660e01b81526004016114c79190615063565b600060405180830381600087803b1580156114e157600080fd5b505af11580156114f5573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a81166004830152602482018890527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b15801561158957600080fd5b505af115801561159d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c19190614b19565b50505050505050505050565b6002336000908152600e602052604090205460ff1660038111156115f3576115f3615493565b1415801561162557506003336000908152600e602052604090205460ff16600381111561162257611622615493565b14155b1561165c576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808061166c84860186614a35565b92509250925060005b83518110156117e357611750848281518110611693576116936154f1565b60200260200101518483815181106116ad576116ad6154f1565b602002602001015160c001518584815181106116cb576116cb6154f1565b6020026020010151608001518685815181106116e9576116e96154f1565b602002602001015160600151878681518110611707576117076154f1565b602002602001015160000151878781518110611725576117256154f1565b602002602001015189888151811061173f5761173f6154f1565b602002602001015160e001516134b0565b838181518110611762576117626154f1565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7184838151811061179d5761179d6154f1565b60209081029190910181015151604080516bffffffffffffffffffffffff909216825233928201929092520160405180910390a2806117db81615407565b915050611675565b505050505050565b60008281526009602090815260409182902082516101008101845281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff80821660808501526401000000008204811660a0850181905268010000000000000000830490931660c08501527c010000000000000000000000000000000000000000000000000000000090910460ff16151560e0840152146118fc576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805161190990839061529b565b600084815260096020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92831617905560145461195d91841690615283565b6014556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b158015611a0157600080fd5b505af1158015611a15573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a399190614b19565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b60058054611a91906153b3565b80601f0160208091040260200160405190810160405280929190818152602001828054611abd906153b3565b8015611b0a5780601f10611adf57610100808354040283529160200191611b0a565b820191906000526020600020905b815481529060010190602001808311611aed57829003601f168201915b505050505081565b73ffffffffffffffffffffffffffffffffffffffff8116611b5f576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600a602090815260409182902082516060810184528154948516808252740100000000000000000000000000000000000000009095046bffffffffffffffffffffffff16928101929092526001015460ff16151591810191909152903314611c10576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8084166000908152600a602090815260409091208054909216909155810151601454611c5f916bffffffffffffffffffffffff1690615343565b60145560208082015160405133815273ffffffffffffffffffffffffffffffffffffffff808616936bffffffffffffffffffffffff90931692908716917f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698910160405180910390a460208101516040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301526bffffffffffffffffffffffff90921660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015611d6c57600080fd5b505af1158015611d80573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611da49190614b19565b50505050565b6108fc8163ffffffff161080611dcb575060105463ffffffff908116908216115b15611e02576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526009602090815260409182902082516101008101845281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff80821660808501526401000000008204811660a0850181905268010000000000000000830490931660c08501527c010000000000000000000000000000000000000000000000000000000090910460ff16151560e084015214611f13576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606081015173ffffffffffffffffffffffffffffffffffffffff163314611f66576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526009602090815260409182902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8616908117909155915191825284917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a2505050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600b6020526040902054163314612043576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600a602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600b909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60008181526009602090815260409182902082516101008101845281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff80821660808501526401000000008204811660a0850181905268010000000000000000830490931660c08501527c010000000000000000000000000000000000000000000000000000000090910460ff16151560e0840152146121ec576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600c602052604090205473ffffffffffffffffffffffffffffffffffffffff163314612249576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060810151600083815260096020908152604080832060010180546bffffffffffffffffffffffff16336c01000000000000000000000000810291909117909155600c90925280832080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905551909173ffffffffffffffffffffffffffffffffffffffff84169186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b6123056131cc565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561238d57600080fd5b505afa1580156123a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123c59190614c3f565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33601454846124129190615343565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401610d16565b6124706131cc565b828114158061247f5750600283105b156124b6576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600654811015612542576000600682815481106124d8576124d86154f1565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168252600a90526040902060010180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055508061253a81615407565b9150506124b9565b5060005b83811015612778576000858583818110612562576125626154f1565b90506020020160208101906125779190614876565b73ffffffffffffffffffffffffffffffffffffffff8082166000908152600a60205260408120805493945092909116908686868181106125b9576125b96154f1565b90506020020160208101906125ce9190614876565b905073ffffffffffffffffffffffffffffffffffffffff81161580612661575073ffffffffffffffffffffffffffffffffffffffff82161580159061263f57508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015612661575073ffffffffffffffffffffffffffffffffffffffff81811614155b15612698576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600183015460ff16156126d7576040517f357d0cc400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600183810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905573ffffffffffffffffffffffffffffffffffffffff818116146127615782547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff82161783555b50505050808061277090615407565b915050612546565b506127856006858561443b565b507f056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f848484846040516127bb9493929190614ed3565b60405180910390a150505050565b60048054611a91906153b3565b60606000806000806127e6613886565b600087815260096020908152604080832081516101008101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff908116848801526001850154918216848701529190048116606083015260029092015463ffffffff808216608084015264010000000082041660a083015268010000000000000000810490921660c08201527c010000000000000000000000000000000000000000000000000000000090910460ff16151560e08201528a8452600d90925280832090519192917f6e04ff0d00000000000000000000000000000000000000000000000000000000916128f6916024016150ad565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000808360c0015173ffffffffffffffffffffffffffffffffffffffff16600f600001600b9054906101000a900463ffffffff1663ffffffff168460405161299d9190614eb7565b60006040518083038160008787f1925050503d80600081146129db576040519150601f19603f3d011682016040523d82523d6000602084013e6129e0565b606091505b509150915081612a1e57806040517f96c36235000000000000000000000000000000000000000000000000000000008152600401610b799190615063565b80806020019051810190612a329190614b3d565b9950915081612a6d576040517f865676e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612a7c8b8d8c60006138be565b9050612a9185826000015183606001516139a8565b60608101516080820151600f5460a08401518d939291612acc91720100000000000000000000000000000000000090910461ffff1690615306565b60c090940151929f919e509c50919a5098509650505050505050565b600081815260096020908152604080832081516101008101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff90811696840196909652600184015490811694830194909452909204831660608301526002015463ffffffff80821660808401526401000000008204811660a08401819052680100000000000000008304851660c08501527c010000000000000000000000000000000000000000000000000000000090920460ff16151560e08401529354919314801592919091163314908290612bee5750808015612bec5750438360a0015163ffffffff16115b155b15612c25576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158015612c635750826060015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15612c9a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b4381612cae57612cab603282615283565b90505b6000858152600960205260409020600201805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117909155612d0a90600790879061349b16565b50600f5460408501516bffffffffffffffffffffffff7401000000000000000000000000000000000000000090920482169160009116821115612d84576040860151612d56908361535a565b905085600001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115612d84575084515b8551612d9190829061535a565b600088815260096020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff928316179055601354612de59183911661529b565b601380547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590612ea3575060165473ffffffffffffffffffffffffffffffffffffffff163314155b15612eda576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612ee5600143615343565b60105460408051924060208401523060601b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690830152640100000000900460e01b7fffffffff000000000000000000000000000000000000000000000000000000001660548201526058016040516020818303038152906040528051906020012060001c9050612fb281878787600088888080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525092506134b0915050565b60108054640100000000900463ffffffff16906004612fd083615440565b91906101000a81548163ffffffff021916908363ffffffff16021790555050807fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012868660405161304892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a295945050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600a60205260409020541633146130b9576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116331415613109576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600b6020526040902054811690821614610d685773ffffffffffffffffffffffffffffffffffffffff8281166000818152600b602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6131c06131cc565b6131c981613af9565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107b3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b79565b806060015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146132b6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a081015163ffffffff908116146131c9576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60035460ff16613366576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610b79565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b60035460ff1615613448576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610b79565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586133b13390565b60006134a78383613bef565b90505b92915050565b60035460ff161561351d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610b79565b73ffffffffffffffffffffffffffffffffffffffff86163b61356b576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8563ffffffff16108061358c575060105463ffffffff908116908616115b156135c3576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806101000160405280846bffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1681526020018663ffffffff16815260200163ffffffff801681526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018215158152506009600089815260200190815260200160002060008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160020160046101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160020160086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060e082015181600201601c6101000a81548160ff021916908315150217905550905050826bffffffffffffffffffffffff1660145461384e9190615283565b6014556000878152600d602090815260409091208351613870928501906144c3565b5061387c600788613ce2565b5050505050505050565b32156107b3576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6139146040518060e00160405280600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526009602052604081206002015463ffffffff169080613936613cee565b91509150600061394884848489613ee9565b6040805160e08101825273ffffffffffffffffffffffffffffffffffffffff909b168b5260208b0199909952978901969096526bffffffffffffffffffffffff9096166060880152608087019190915260a086015250505060c082015290565b8260e00151156139e4576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600a602052604090206001015460ff16613a46576040517fcfbacfd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82516bffffffffffffffffffffffff16811115613a8f576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff16836020015173ffffffffffffffffffffffffffffffffffffffff16141561079e576040517f06bc104000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116331415613b79576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b79565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008181526001830160205260408120548015613cd8576000613c13600183615343565b8554909150600090613c2790600190615343565b9050818114613c8c576000866000018281548110613c4757613c476154f1565b9060005260206000200154905080876000018481548110613c6a57613c6a6154f1565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613c9d57613c9d6154c2565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506134aa565b60009150506134aa565b60006134a783836143b2565b6000806000600f600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015613d8557600080fd5b505afa158015613d99573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dbd9190614cc1565b509450909250849150508015613de15750613dd88242615343565b8463ffffffff16105b80613ded575060008113155b15613dfc576011549550613e00565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015613e6657600080fd5b505afa158015613e7a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e9e9190614cc1565b509450909250849150508015613ec25750613eb98242615343565b8463ffffffff16105b80613ece575060008113155b15613edd576012549450613ee1565b8094505b505050509091565b6040805161012081018252600f5463ffffffff80821683526401000000008083048216602085015268010000000000000000830462ffffff908116958501959095526b0100000000000000000000008304821660608501526f01000000000000000000000000000000830490941660808401527201000000000000000000000000000000000000820461ffff1660a08401819052740100000000000000000000000000000000000000009092046bffffffffffffffffffffffff1660c084015260105480821660e0850152939093049092166101008201526000918290613fd09087615306565b9050838015613fde5750803a105b15613fe657503a5b60006140127f000000000000000000000000000000000000000000000000000000000000000089615283565b61401c9083615306565b83519091506000906140389063ffffffff16633b9aca00615283565b9050600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561407057614070615493565b14156142105760408051600081526020810190915287156140b85760003660046040516020016140a293929190614e18565b6040516020818303038152906040529050614146565b600580546140c5906153b3565b80601f01602080910402602001604051908101604052809291908181526020018280546140f1906153b3565b801561413e5780601f106141135761010080835404028352916020019161413e565b820191906000526020600020905b81548152906001019060200180831161412157829003601f168201915b505050505090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906349948e0e906141b8908490600401615063565b60206040518083038186803b1580156141d057600080fd5b505afa1580156141e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142089190614c3f565b9150506142eb565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561424457614244615493565b14156142eb577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156142b057600080fd5b505afa1580156142c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142e89190614c3f565b90505b8661430757808560a0015161ffff166143049190615306565b90505b6000856020015163ffffffff1664e8d4a510006143249190615306565b89846143308588615283565b61433e90633b9aca00615306565b6143489190615306565b61435291906152cb565b61435c9190615283565b90506b033b2e3c9fd0803ce80000008111156143a4576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9a9950505050505050505050565b60008181526001830160205260408120546143f9575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556134aa565b5060006134aa565b50805461440d906153b3565b6000825580601f1061441d575050565b601f0160209004906000526020600020908101906131c99190614537565b8280548282559060005260206000209081019282156144b3579160200282015b828111156144b35781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84351617825560209092019160019091019061445b565b506144bf929150614537565b5090565b8280546144cf906153b3565b90600052602060002090601f0160209004810192826144f157600085556144b3565b82601f1061450a57805160ff19168380011785556144b3565b828001600101855582156144b3579182015b828111156144b357825182559160200191906001019061451c565b5b808211156144bf5760008155600101614538565b803573ffffffffffffffffffffffffffffffffffffffff8116811461457057600080fd5b919050565b60008083601f84011261458757600080fd5b50813567ffffffffffffffff81111561459f57600080fd5b6020830191508360208260051b85010111156145ba57600080fd5b9250929050565b600082601f8301126145d257600080fd5b813560206145e76145e283615219565b6151ca565b80838252828201915082860187848660051b890101111561460757600080fd5b60005b8581101561468757813567ffffffffffffffff81111561462957600080fd5b8801603f81018a1361463a57600080fd5b85810135604061464c6145e28361523d565b8281528c8284860101111561466057600080fd5b828285018a830137600092810189019290925250855250928401929084019060010161460a565b5090979650505050505050565b600082601f8301126146a557600080fd5b813560206146b56145e283615219565b80838252828201915082860187848660081b89010111156146d557600080fd5b6000805b8681101561478f5761010080848c0312156146f2578283fd5b6146fa6151a0565b6147038561485a565b815261471088860161454c565b88820152604061472181870161485a565b90820152606061473286820161454c565b90820152608061474386820161482c565b9082015260a061475486820161482c565b9082015260c061476586820161454c565b9082015260e0858101356147788161554f565b9082015286529486019492909201916001016146d9565b509198975050505050505050565b60008083601f8401126147af57600080fd5b50813567ffffffffffffffff8111156147c757600080fd5b6020830191508360208285010111156145ba57600080fd5b600082601f8301126147f057600080fd5b81516147fe6145e28261523d565b81815284602083860101111561481357600080fd5b614824826020830160208701615387565b949350505050565b803563ffffffff8116811461457057600080fd5b805169ffffffffffffffffffff8116811461457057600080fd5b80356bffffffffffffffffffffffff8116811461457057600080fd5b60006020828403121561488857600080fd5b6134a78261454c565b600080604083850312156148a457600080fd5b6148ad8361454c565b91506148bb6020840161454c565b90509250929050565b600080604083850312156148d757600080fd5b6148e08361454c565b91506020830135600481106148f457600080fd5b809150509250929050565b60008060008060006080868803121561491757600080fd5b6149208661454c565b945061492e6020870161482c565b935061493c6040870161454c565b9250606086013567ffffffffffffffff81111561495857600080fd5b6149648882890161479d565b969995985093965092949392505050565b6000806000806040858703121561498b57600080fd5b843567ffffffffffffffff808211156149a357600080fd5b6149af88838901614575565b909650945060208701359150808211156149c857600080fd5b506149d587828801614575565b95989497509550505050565b6000806000604084860312156149f657600080fd5b833567ffffffffffffffff811115614a0d57600080fd5b614a1986828701614575565b9094509250614a2c90506020850161454c565b90509250925092565b600080600060608486031215614a4a57600080fd5b833567ffffffffffffffff80821115614a6257600080fd5b818601915086601f830112614a7657600080fd5b81356020614a866145e283615219565b8083825282820191508286018b848660051b8901011115614aa657600080fd5b600096505b84871015614ac9578035835260019690960195918301918301614aab565b5097505087013592505080821115614ae057600080fd5b614aec87838801614694565b93506040860135915080821115614b0257600080fd5b50614b0f868287016145c1565b9150509250925092565b600060208284031215614b2b57600080fd5b8151614b368161554f565b9392505050565b60008060408385031215614b5057600080fd5b8251614b5b8161554f565b602084015190925067ffffffffffffffff811115614b7857600080fd5b614b84858286016147df565b9150509250929050565b60008060208385031215614ba157600080fd5b823567ffffffffffffffff811115614bb857600080fd5b614bc48582860161479d565b90969095509350505050565b600060208284031215614be257600080fd5b815167ffffffffffffffff811115614bf957600080fd5b614824848285016147df565b600060208284031215614c1757600080fd5b815160028110614b3657600080fd5b600060208284031215614c3857600080fd5b5035919050565b600060208284031215614c5157600080fd5b5051919050565b60008060408385031215614c6b57600080fd5b823591506148bb6020840161454c565b60008060408385031215614c8e57600080fd5b823591506148bb6020840161482c565b60008060408385031215614cb157600080fd5b823591506148bb6020840161485a565b600080600080600060a08688031215614cd957600080fd5b614ce286614840565b9450602086015193506040860151925060608601519150614d0560808701614840565b90509295509295909350565b8183526000602080850194508260005b85811015614d5a5773ffffffffffffffffffffffffffffffffffffffff614d478361454c565b1687529582019590820190600101614d21565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b85811015614dad578284038952614d9b848351614dba565b98850198935090840190600101614d83565b5091979650505050505050565b60008151808452614dd2816020860160208601615387565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60028110614e1457614e14615493565b9052565b8284823760008382016000815260008454614e32816153b3565b60018281168015614e4a5760018114614e7957614ea8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841686528286019450614ea8565b8860005260208060002060005b85811015614e9f57815489820152908401908201614e86565b50505082860194505b50929998505050505050505050565b60008251614ec9818460208701615387565b9190910192915050565b604081526000614ee7604083018688614d11565b8281036020840152614efa818587614d11565b979650505050505050565b60006060808352858184015260807f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff871115614f4057600080fd5b8660051b808983870137808501905081810160008152602083878403018188015281895180845260a093508385019150828b01945060005b8181101561503f5785516bffffffffffffffffffffffff80825116855273ffffffffffffffffffffffffffffffffffffffff868301511686860152604081818401511681870152505088810151614fe68a86018273ffffffffffffffffffffffffffffffffffffffff169052565b508781015163ffffffff908116858a015286820151168685015260c08082015173ffffffffffffffffffffffffffffffffffffffff169085015260e0908101511515908401529483019461010090920191600101614f78565b50508781036040890152615053818a614d65565b9c9b505050505050505050505050565b6020815260006134a76020830184614dba565b60a08152600061508960a0830188614dba565b90508560208301528460408301528360608301528260808301529695505050505050565b60006020808352600084546150c1816153b3565b808487015260406001808416600081146150e2576001811461511457615142565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838a0152606089019550615142565b896000528660002060005b8581101561513a5781548b820186015290830190880161511f565b8a0184019650505b509398975050505050505050565b602081016003831061516457615164615493565b91905290565b6151748185614e04565b6151816020820184614e04565b6060604082015260006151976060830184614dba565b95945050505050565b604051610100810167ffffffffffffffff811182821017156151c4576151c4615520565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561521157615211615520565b604052919050565b600067ffffffffffffffff82111561523357615233615520565b5060051b60200190565b600067ffffffffffffffff82111561525757615257615520565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111561529657615296615464565b500190565b60006bffffffffffffffffffffffff8083168185168083038211156152c2576152c2615464565b01949350505050565b600082615301577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561533e5761533e615464565b500290565b60008282101561535557615355615464565b500390565b60006bffffffffffffffffffffffff8381169083168181101561537f5761537f615464565b039392505050565b60005b838110156153a257818101518382015260200161538a565b83811115611da45750506000910152565b600181811c908216806153c757607f821691505b60208210811415615401577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561543957615439615464565b5060010190565b600063ffffffff8083168181141561545a5761545a615464565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80151581146131c957600080fdfea164736f6c6343000806000a3078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
}

var KeeperRegistryLogicABI = KeeperRegistryLogicMetaData.ABI

var KeeperRegistryLogicBin = KeeperRegistryLogicMetaData.Bin

func DeployKeeperRegistryLogic(auth *bind.TransactOpts, backend bind.ContractBackend, paymentModel uint8, registryGasOverhead *big.Int, link common.Address, linkEthFeed common.Address, fastGasFeed common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogic, error) {
	parsed, err := KeeperRegistryLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicBin), backend, paymentModel, registryGasOverhead, link, linkEthFeed, fastGasFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogic{KeeperRegistryLogicCaller: KeeperRegistryLogicCaller{contract: contract}, KeeperRegistryLogicTransactor: KeeperRegistryLogicTransactor{contract: contract}, KeeperRegistryLogicFilterer: KeeperRegistryLogicFilterer{contract: contract}}, nil
}

type KeeperRegistryLogic struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogicCaller
	KeeperRegistryLogicTransactor
	KeeperRegistryLogicFilterer
}

type KeeperRegistryLogicCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicSession struct {
	Contract     *KeeperRegistryLogic
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicCallerSession struct {
	Contract *KeeperRegistryLogicCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogicTransactorSession struct {
	Contract     *KeeperRegistryLogicTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicRaw struct {
	Contract *KeeperRegistryLogic
}

type KeeperRegistryLogicCallerRaw struct {
	Contract *KeeperRegistryLogicCaller
}

type KeeperRegistryLogicTransactorRaw struct {
	Contract *KeeperRegistryLogicTransactor
}

func NewKeeperRegistryLogic(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogic, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogicABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic{address: address, abi: abi, KeeperRegistryLogicCaller: KeeperRegistryLogicCaller{contract: contract}, KeeperRegistryLogicTransactor: KeeperRegistryLogicTransactor{contract: contract}, KeeperRegistryLogicFilterer: KeeperRegistryLogicFilterer{contract: contract}}, nil
}

func NewKeeperRegistryLogicCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogicCaller, error) {
	contract, err := bindKeeperRegistryLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicCaller{contract: contract}, nil
}

func NewKeeperRegistryLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogicTransactor, error) {
	contract, err := bindKeeperRegistryLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicTransactor{contract: contract}, nil
}

func NewKeeperRegistryLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogicFilterer, error) {
	contract, err := bindKeeperRegistryLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicFilterer{contract: contract}, nil
}

func bindKeeperRegistryLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistryLogicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogic.Contract.KeeperRegistryLogicCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.KeeperRegistryLogicTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.KeeperRegistryLogicTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogic.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "ARB_NITRO_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.ARBNITROORACLE(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.ARBNITROORACLE(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.FASTGASFEED(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.FASTGASFEED(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) L1FEEDATAPADDING(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "L1_FEE_DATA_PADDING")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) L1FEEDATAPADDING() ([]byte, error) {
	return _KeeperRegistryLogic.Contract.L1FEEDATAPADDING(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) L1FEEDATAPADDING() ([]byte, error) {
	return _KeeperRegistryLogic.Contract.L1FEEDATAPADDING(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) LINK() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.LINK(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.LINK(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.LINKETHFEED(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.LINKETHFEED(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) MAXINPUTDATA(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "MAX_INPUT_DATA")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) MAXINPUTDATA() ([]byte, error) {
	return _KeeperRegistryLogic.Contract.MAXINPUTDATA(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) MAXINPUTDATA() ([]byte, error) {
	return _KeeperRegistryLogic.Contract.MAXINPUTDATA(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "OPTIMISM_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.OPTIMISMORACLE(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.OPTIMISMORACLE(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) PAYMENTMODEL(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "PAYMENT_MODEL")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistryLogic.Contract.PAYMENTMODEL(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistryLogic.Contract.PAYMENTMODEL(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "REGISTRY_GAS_OVERHEAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistryLogic.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistryLogic.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.Owner(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.Owner(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) Paused() (bool, error) {
	return _KeeperRegistryLogic.Contract.Paused(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) Paused() (bool, error) {
	return _KeeperRegistryLogic.Contract.Paused(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptOwnership(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptOwnership(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "acceptPayeeship", keeper)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptPayeeship(&_KeeperRegistryLogic.TransactOpts, keeper)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptPayeeship(&_KeeperRegistryLogic.TransactOpts, keeper)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AddFunds(&_KeeperRegistryLogic.TransactOpts, id, amount)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AddFunds(&_KeeperRegistryLogic.TransactOpts, id, amount)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CancelUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CancelUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "checkUpkeep", id, from)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CheckUpkeep(&_KeeperRegistryLogic.TransactOpts, id, from)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CheckUpkeep(&_KeeperRegistryLogic.TransactOpts, id, from)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.MigrateUpkeeps(&_KeeperRegistryLogic.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.MigrateUpkeeps(&_KeeperRegistryLogic.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "pause")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.Pause(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.Pause(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.ReceiveUpkeeps(&_KeeperRegistryLogic.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.ReceiveUpkeeps(&_KeeperRegistryLogic.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RecoverFunds(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RecoverFunds(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RegisterUpkeep(&_KeeperRegistryLogic.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RegisterUpkeep(&_KeeperRegistryLogic.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "setKeepers", keepers, payees)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetKeepers(&_KeeperRegistryLogic.TransactOpts, keepers, payees)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetKeepers(&_KeeperRegistryLogic.TransactOpts, keepers, payees)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogic.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogic.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogic.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogic.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferOwnership(&_KeeperRegistryLogic.TransactOpts, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferOwnership(&_KeeperRegistryLogic.TransactOpts, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferPayeeship(&_KeeperRegistryLogic.TransactOpts, keeper, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferPayeeship(&_KeeperRegistryLogic.TransactOpts, keeper, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogic.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogic.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "unpause")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.Unpause(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.Unpause(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawFunds(&_KeeperRegistryLogic.TransactOpts, id, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawFunds(&_KeeperRegistryLogic.TransactOpts, id, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawPayment(&_KeeperRegistryLogic.TransactOpts, from, to)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.WithdrawPayment(&_KeeperRegistryLogic.TransactOpts, from, to)
}

type KeeperRegistryLogicConfigSetIterator struct {
	Event *KeeperRegistryLogicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicConfigSet struct {
	Config Config
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryLogicConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicConfigSetIterator{contract: _KeeperRegistryLogic.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicConfigSet)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryLogicConfigSet, error) {
	event := new(KeeperRegistryLogicConfigSet)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicFundsAddedIterator struct {
	Event *KeeperRegistryLogicFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicFundsAddedIterator{contract: _KeeperRegistryLogic.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicFundsAdded)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogicFundsAdded, error) {
	event := new(KeeperRegistryLogicFundsAdded)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicFundsWithdrawnIterator{contract: _KeeperRegistryLogic.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicFundsWithdrawn)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicFundsWithdrawn)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicKeepersUpdatedIterator struct {
	Event *KeeperRegistryLogicKeepersUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicKeepersUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicKeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicKeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicKeepersUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicKeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicKeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicKeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicKeepersUpdatedIterator{contract: _KeeperRegistryLogic.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicKeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicKeepersUpdated)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistryLogicKeepersUpdated, error) {
	event := new(KeeperRegistryLogicKeepersUpdated)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicOwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicOwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicOwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogic.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicOwnerFundsWithdrawn)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicOwnerFundsWithdrawn)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicOwnershipTransferRequestedIterator{contract: _KeeperRegistryLogic.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicOwnershipTransferRequested)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicOwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogicOwnershipTransferRequested)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicOwnershipTransferredIterator struct {
	Event *KeeperRegistryLogicOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicOwnershipTransferredIterator{contract: _KeeperRegistryLogic.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicOwnershipTransferred)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicOwnershipTransferred, error) {
	event := new(KeeperRegistryLogicOwnershipTransferred)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicPausedIterator struct {
	Event *KeeperRegistryLogicPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicPausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPausedIterator{contract: _KeeperRegistryLogic.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicPaused)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParsePaused(log types.Log) (*KeeperRegistryLogicPaused, error) {
	event := new(KeeperRegistryLogicPaused)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicPayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogic.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicPayeeshipTransferRequested)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogicPayeeshipTransferRequested)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicPayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogicPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicPayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPayeeshipTransferredIterator{contract: _KeeperRegistryLogic.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicPayeeshipTransferred)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicPayeeshipTransferred, error) {
	event := new(KeeperRegistryLogicPayeeshipTransferred)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicPaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogicPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicPaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPaymentWithdrawnIterator{contract: _KeeperRegistryLogic.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicPaymentWithdrawn)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicPaymentWithdrawn, error) {
	event := new(KeeperRegistryLogicPaymentWithdrawn)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUnpausedIterator struct {
	Event *KeeperRegistryLogicUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUnpausedIterator{contract: _KeeperRegistryLogic.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUnpaused)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogicUnpaused, error) {
	event := new(KeeperRegistryLogicUnpaused)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogicUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogicUpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogicUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepAdminTransferredIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepAdminTransferred)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogicUpkeepAdminTransferred)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepCanceledIterator struct {
	Event *KeeperRegistryLogicUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepCanceledIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepCanceled)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicUpkeepCanceled, error) {
	event := new(KeeperRegistryLogicUpkeepCanceled)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistryLogicUpkeepCheckDataUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepCheckDataUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepCheckDataUpdatedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepCheckDataUpdated)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogicUpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistryLogicUpkeepCheckDataUpdated)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogicUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepGasLimitSetIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepGasLimitSet)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogicUpkeepGasLimitSet)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepMigratedIterator struct {
	Event *KeeperRegistryLogicUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepMigratedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepMigrated)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicUpkeepMigrated, error) {
	event := new(KeeperRegistryLogicUpkeepMigrated)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepPausedIterator struct {
	Event *KeeperRegistryLogicUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepPausedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepPaused)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicUpkeepPaused, error) {
	event := new(KeeperRegistryLogicUpkeepPaused)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepPerformedIterator struct {
	Event *KeeperRegistryLogicUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepPerformed struct {
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryLogicUpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepPerformedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepPerformed)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicUpkeepPerformed, error) {
	event := new(KeeperRegistryLogicUpkeepPerformed)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepReceivedIterator struct {
	Event *KeeperRegistryLogicUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepReceivedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepReceived)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicUpkeepReceived, error) {
	event := new(KeeperRegistryLogicUpkeepReceived)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogicUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepRegisteredIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepRegistered)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicUpkeepRegistered, error) {
	event := new(KeeperRegistryLogicUpkeepRegistered)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicUpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogicUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepUnpausedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepUnpaused)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicUpkeepUnpaused, error) {
	event := new(KeeperRegistryLogicUpkeepUnpaused)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogic) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogic.abi.Events["ConfigSet"].ID:
		return _KeeperRegistryLogic.ParseConfigSet(log)
	case _KeeperRegistryLogic.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogic.ParseFundsAdded(log)
	case _KeeperRegistryLogic.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogic.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["KeepersUpdated"].ID:
		return _KeeperRegistryLogic.ParseKeepersUpdated(log)
	case _KeeperRegistryLogic.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogic.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogic.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogic.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogic.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogic.abi.Events["Paused"].ID:
		return _KeeperRegistryLogic.ParsePaused(log)
	case _KeeperRegistryLogic.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogic.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogic.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogic.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogic.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogic.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogic.ParseUnpaused(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogic.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogic.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogic.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepCheckDataUpdated"].ID:
		return _KeeperRegistryLogic.ParseUpkeepCheckDataUpdated(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogic.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogic.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogic.ParseUpkeepPaused(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogic.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogic.ParseUpkeepReceived(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogic.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogic.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogic.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325")
}

func (KeeperRegistryLogicFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicKeepersUpdated) Topic() common.Hash {
	return common.HexToHash("0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f")
}

func (KeeperRegistryLogicOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogicOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogicOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogicPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogicPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogicPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogicPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogicUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogicUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogicUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogicUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogicUpkeepCheckDataUpdated) Topic() common.Hash {
	return common.HexToHash("0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf")
}

func (KeeperRegistryLogicUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogicUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogicUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6")
}

func (KeeperRegistryLogicUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogicUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogic *KeeperRegistryLogic) Address() common.Address {
	return _KeeperRegistryLogic.address
}

type KeeperRegistryLogicInterface interface {
	ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error)

	FASTGASFEED(opts *bind.CallOpts) (common.Address, error)

	L1FEEDATAPADDING(opts *bind.CallOpts) ([]byte, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	MAXINPUTDATA(opts *bind.CallOpts) ([]byte, error)

	OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error)

	PAYMENTMODEL(opts *bind.CallOpts) (uint8, error)

	REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error)

	SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryLogicConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryLogicConfigSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicFundsWithdrawn, error)

	FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicKeepersUpdatedIterator, error)

	WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicKeepersUpdated) (event.Subscription, error)

	ParseKeepersUpdated(log types.Log) (*KeeperRegistryLogicKeepersUpdated, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogicPaused, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicPaymentWithdrawn, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogicUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicUpkeepCanceled, error)

	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepCheckDataUpdatedIterator, error)

	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogicUpkeepCheckDataUpdated, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicUpkeepMigrated, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryLogicUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicUpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicUpkeepRegistered, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
